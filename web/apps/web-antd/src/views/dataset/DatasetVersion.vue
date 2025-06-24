<template>
  <div class="dataset-version-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <DatabaseOutlined class="title-icon" />
            <span class="title-text">数据集版本管理</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">管理和版本控制您的训练数据集</span>
          </p>
        </div>
        <div class="action-section">
          <a-button
            type="primary"
            size="large"
            @click="showCreateModal"
            class="create-btn"
          >
            <PlusOutlined />
            创建版本
          </a-button>
        </div>
      </div>
    </div>

    <!-- 数据集信息卡片 -->
    <div class="dataset-info-section">
      <a-card class="dataset-info-card glass-card" :bordered="false">
        <div class="dataset-header">
          <div class="dataset-title">
            <h2>{{ currentDataset.name }}</h2>
            <a-tag color="blue" class="dataset-type-tag">
              <FolderOutlined />
              {{ currentDataset.type }}
            </a-tag>
          </div>
          <div class="dataset-actions">
            <a-button @click="showUploadModal" class="upload-btn">
              <UploadOutlined />
              上传数据
            </a-button>
            <a-button @click="showSettingsModal" class="settings-btn">
              <SettingOutlined />
              设置
            </a-button>
          </div>
        </div>
        <div class="dataset-meta">
          <div class="meta-item">
            <span class="meta-label">创建者:</span>
            <span class="meta-value">{{ currentDataset.creator }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">创建时间:</span>
            <span class="meta-value">{{ currentDataset.createTime }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">总版本数:</span>
            <span class="meta-value">{{ versions.length }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">存储路径:</span>
            <span class="meta-value">{{ currentDataset.storagePath }}</span>
          </div>
        </div>
        <div class="dataset-description">
          <span class="description-label">描述:</span>
          <span class="description-content">{{
            currentDataset.description
          }}</span>
        </div>
      </a-card>
    </div>

    <!-- 筛选器 -->
    <div class="filter-section">
      <a-card class="filter-card glass-card" :bordered="false">
        <a-row :gutter="16" align="middle">
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterStatus"
              placeholder="选择状态"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部状态</a-select-option>
              <a-select-option value="active">活跃</a-select-option>
              <a-select-option value="processing">处理中</a-select-option>
              <a-select-option value="completed">已完成</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
              <a-select-option value="archived">已归档</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterTag"
              placeholder="选择标签"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部标签</a-select-option>
              <a-select-option value="stable">稳定版</a-select-option>
              <a-select-option value="beta">测试版</a-select-option>
              <a-select-option value="experimental">实验版</a-select-option>
              <a-select-option value="hotfix">修复版</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="16" :md="8" :lg="8">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索版本号或描述"
              allow-clear
              @search="handleSearch"
              @change="handleSearchChange"
              class="search-input"
            />
          </a-col>
          <a-col :xs="24" :sm="8" :md="4" :lg="4" class="refresh-btn-col">
            <a-button
              @click="refreshData"
              :loading="loading"
              class="refresh-btn"
            >
              <ReloadOutlined />
              刷新
            </a-button>
          </a-col>
        </a-row>
      </a-card>
    </div>

    <!-- 版本列表 -->
    <div class="table-section">
      <a-card class="table-card glass-card" :bordered="false">
        <a-table
          :columns="columns"
          :data-source="filteredVersions"
          :loading="loading"
          :pagination="paginationConfig"
          row-key="id"
          size="middle"
          :scroll="{ x: 'max-content' }"
          class="sci-fi-table"
        >
          <!-- 版本号列 -->
          <template #version="{ record }">
            <div class="version-wrapper">
              <span class="version-number">{{ record.version }}</span>
              <a-tag v-if="record.isLatest" color="success" class="latest-tag">
                <CrownOutlined />
                最新
              </a-tag>
            </div>
          </template>

          <!-- 状态列 -->
          <template #status="{ record }">
            <div class="status-wrapper">
              <a-tag :color="getStatusColor(record.status)" class="status-tag">
                <component
                  :is="getStatusIcon(record.status)"
                  class="status-icon"
                />
                {{ getStatusText(record.status) }}
              </a-tag>
              <div
                class="status-indicator"
                :class="`indicator-${record.status}`"
              ></div>
            </div>
          </template>

          <!-- 标签列 -->
          <template #tags="{ record }">
            <div class="tags-wrapper">
              <a-tag
                v-for="tag in record.tags"
                :key="tag"
                :color="getTagColor(tag)"
                class="version-tag"
              >
                {{ getTagText(tag) }}
              </a-tag>
            </div>
          </template>

          <!-- 文件信息列 -->
          <template #fileInfo="{ record }">
            <div class="file-info">
              <div class="file-item">
                <FileOutlined class="file-icon" />
                <span class="file-label">文件数:</span>
                <span class="file-value">{{ record.fileCount }}</span>
              </div>
              <div class="file-item">
                <CloudServerOutlined class="file-icon" />
                <span class="file-label">大小:</span>
                <span class="file-value">{{
                  formatFileSize(record.size)
                }}</span>
              </div>
            </div>
          </template>

          <!-- 创建时间列 -->
          <template #createTime="{ record }">
            <a-tooltip :title="record.createTime">
              <span class="time-text">{{
                formatRelativeTime(record.createTime)
              }}</span>
            </a-tooltip>
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space class="action-buttons">
              <a-button
                type="link"
                size="small"
                @click="downloadVersion(record)"
                :disabled="record.status !== 'completed'"
                class="action-btn"
              >
                <DownloadOutlined />
                下载
              </a-button>
              <a-button
                type="link"
                size="small"
                @click="viewDetails(record)"
                class="action-btn"
              >
                <EyeOutlined />
                详情
              </a-button>
              <a-dropdown>
                <a-button type="link" size="small" class="action-btn">
                  <MoreOutlined />
                </a-button>
                <template #overlay>
                  <a-menu
                    @click="(item: any) => handleMenuAction(item.key, record)"
                    class="action-menu"
                  >
                    <a-menu-item
                      key="setLatest"
                      :disabled="
                        record.isLatest || record.status !== 'completed'
                      "
                    >
                      <CrownOutlined />
                      设为最新
                    </a-menu-item>
                    <a-menu-item key="compare">
                      <DiffOutlined />
                      版本对比
                    </a-menu-item>
                    <a-menu-item key="clone">
                      <CopyOutlined />
                      克隆版本
                    </a-menu-item>
                    <a-menu-item
                      key="archive"
                      :disabled="record.status === 'archived'"
                    >
                      <InboxOutlined />
                      归档
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item key="delete" class="danger-item">
                      <DeleteOutlined />
                      删除
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 创建版本模态框 -->
    <a-modal
      v-model:open="createModalVisible"
      title="创建新版本"
      width="700px"
      :confirm-loading="createLoading"
      @ok="handleCreateSubmit"
      @cancel="handleCreateCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="createFormRef"
        :model="createForm"
        :rules="createFormRules"
        layout="vertical"
        class="create-form"
      >
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="版本号" name="version">
              <a-input
                v-model:value="createForm.version"
                placeholder="例如: v1.2.0"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="版本标签" name="tags">
              <a-select
                v-model:value="createForm.tags"
                mode="multiple"
                placeholder="选择标签"
                class="form-select"
              >
                <a-select-option value="stable">稳定版</a-select-option>
                <a-select-option value="beta">测试版</a-select-option>
                <a-select-option value="experimental">实验版</a-select-option>
                <a-select-option value="hotfix">修复版</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="数据来源" name="source">
          <a-radio-group v-model:value="createForm.source" class="source-radio">
            <a-radio value="upload">上传新数据</a-radio>
            <a-radio value="copy">复制现有版本</a-radio>
            <a-radio value="merge">合并多个版本</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item
          v-if="createForm.source === 'copy'"
          label="源版本"
          name="sourceVersion"
        >
          <a-select
            v-model:value="createForm.sourceVersion"
            placeholder="选择要复制的版本"
            class="form-select"
          >
            <a-select-option
              v-for="version in completedVersions"
              :key="version.id"
              :value="version.id"
            >
              {{ version.version }} - {{ version.description }}
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item
          v-if="createForm.source === 'merge'"
          label="要合并的版本"
          name="mergeVersions"
        >
          <a-select
            v-model:value="createForm.mergeVersions"
            mode="multiple"
            placeholder="选择要合并的版本"
            class="form-select"
          >
            <a-select-option
              v-for="version in completedVersions"
              :key="version.id"
              :value="version.id"
            >
              {{ version.version }} - {{ version.description }}
            </a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="版本描述" name="description">
          <a-textarea
            v-model:value="createForm.description"
            placeholder="请输入版本描述和变更说明"
            :rows="4"
            class="form-textarea"
          />
        </a-form-item>

        <a-form-item v-if="createForm.source === 'upload'" label="文件上传">
          <a-upload-dragger
            v-model:fileList="createForm.fileList"
            multiple
            :before-upload="beforeUpload"
            class="upload-dragger"
          >
            <p class="ant-upload-drag-icon">
              <CloudUploadOutlined />
            </p>
            <p class="ant-upload-text">点击或拖拽文件到此区域上传</p>
            <p class="ant-upload-hint">
              支持单个或批量上传，严格禁止上传敏感数据
            </p>
          </a-upload-dragger>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 版本详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="版本详情"
      width="900px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedVersion" class="detail-content">
        <a-descriptions
          :column="{ xs: 1, sm: 2 }"
          bordered
          class="detail-descriptions"
        >
          <a-descriptions-item label="版本号">
            <div class="version-detail">
              {{ selectedVersion.version }}
              <a-tag
                v-if="selectedVersion.isLatest"
                color="success"
                class="latest-tag"
              >
                <CrownOutlined />
                最新版本
              </a-tag>
            </div>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag
              :color="getStatusColor(selectedVersion.status)"
              class="status-tag"
            >
              <component :is="getStatusIcon(selectedVersion.status)" />
              {{ getStatusText(selectedVersion.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="创建者">
            {{ selectedVersion.creator }}
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">
            {{ selectedVersion.createTime }}
          </a-descriptions-item>
          <a-descriptions-item label="文件数量">
            {{ selectedVersion.fileCount }} 个文件
          </a-descriptions-item>
          <a-descriptions-item label="数据大小">
            {{ formatFileSize(selectedVersion.size) }}
          </a-descriptions-item>
          <a-descriptions-item label="下载次数">
            {{ selectedVersion.downloadCount }} 次
          </a-descriptions-item>
          <a-descriptions-item label="校验和">
            {{ selectedVersion.checksum }}
          </a-descriptions-item>
          <a-descriptions-item label="标签" :span="2">
            <div class="tags-detail">
              <a-tag
                v-for="tag in selectedVersion.tags"
                :key="tag"
                :color="getTagColor(tag)"
                class="version-tag"
              >
                {{ getTagText(tag) }}
              </a-tag>
            </div>
          </a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">
            {{ selectedVersion.description || '暂无描述' }}
          </a-descriptions-item>
        </a-descriptions>

        <!-- 文件列表 -->
        <a-divider class="file-divider">文件列表</a-divider>
        <div class="file-list-container">
          <a-table
            :columns="fileColumns"
            :data-source="selectedVersion.files"
            :pagination="false"
            size="small"
            class="file-table"
          >
            <template #fileName="{ record }">
              <div class="file-name-wrapper">
                <FileOutlined class="file-type-icon" />
                <span class="file-name">{{ record.name }}</span>
              </div>
            </template>
            <template #fileSize="{ record }">
              <span class="file-size">{{ formatFileSize(record.size) }}</span>
            </template>
            <template #fileAction="{ record }">
              <a-button
                type="link"
                size="small"
                @click="downloadFile(record)"
                class="file-action-btn"
              >
                <DownloadOutlined />
                下载
              </a-button>
            </template>
          </a-table>
        </div>
      </div>
    </a-modal>

    <!-- 上传数据模态框 -->
    <a-modal
      v-model:open="uploadModalVisible"
      title="上传数据文件"
      width="600px"
      :confirm-loading="uploadLoading"
      @ok="handleUploadSubmit"
      @cancel="handleUploadCancel"
      class="sci-fi-modal"
    >
      <div class="upload-content">
        <a-upload-dragger
          v-model:fileList="uploadFileList"
          multiple
          :before-upload="beforeUpload"
          class="upload-dragger"
        >
          <p class="ant-upload-drag-icon">
            <CloudUploadOutlined />
          </p>
          <p class="ant-upload-text">点击或拖拽文件到此区域上传</p>
          <p class="ant-upload-hint">支持单个或批量上传，将自动创建新版本</p>
        </a-upload-dragger>

        <a-form layout="vertical" class="upload-form">
          <a-form-item label="版本号">
            <a-input
              v-model:value="uploadForm.version"
              placeholder="自动生成或手动输入"
              class="form-input"
            />
          </a-form-item>
          <a-form-item label="变更说明">
            <a-textarea
              v-model:value="uploadForm.description"
              placeholder="请描述本次上传的变更内容"
              :rows="3"
              class="form-textarea"
            />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>

    <!-- 设置模态框 -->
    <a-modal
      v-model:open="settingsModalVisible"
      title="数据集设置"
      width="600px"
      :confirm-loading="settingsLoading"
      @ok="handleSettingsSubmit"
      @cancel="handleSettingsCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="settingsFormRef"
        :model="settingsForm"
        layout="vertical"
        class="settings-form"
      >
        <a-form-item label="数据集名称">
          <a-input
            v-model:value="settingsForm.name"
            placeholder="请输入数据集名称"
            class="form-input"
          />
        </a-form-item>
        <a-form-item label="数据集类型">
          <a-select
            v-model:value="settingsForm.type"
            placeholder="选择数据集类型"
            class="form-select"
          >
            <a-select-option value="image">图像数据</a-select-option>
            <a-select-option value="text">文本数据</a-select-option>
            <a-select-option value="audio">音频数据</a-select-option>
            <a-select-option value="video">视频数据</a-select-option>
            <a-select-option value="tabular">表格数据</a-select-option>
            <a-select-option value="mixed">混合数据</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="存储路径">
          <a-input
            v-model:value="settingsForm.storagePath"
            placeholder="请输入存储路径"
            class="form-input"
          />
        </a-form-item>
        <a-form-item label="访问权限">
          <a-select
            v-model:value="settingsForm.access"
            placeholder="选择访问权限"
            class="form-select"
          >
            <a-select-option value="private">私有</a-select-option>
            <a-select-option value="team">团队可见</a-select-option>
            <a-select-option value="public">公开</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea
            v-model:value="settingsForm.description"
            placeholder="请输入数据集描述"
            :rows="4"
            class="form-textarea"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 版本对比模态框 -->
    <a-modal
      v-model:open="compareModalVisible"
      title="版本对比"
      width="1000px"
      :footer="null"
      class="sci-fi-modal compare-modal"
    >
      <div class="compare-content">
        <div class="compare-selector">
          <a-row :gutter="16">
            <a-col :span="12">
              <a-form-item label="基准版本">
                <a-select
                  v-model:value="compareForm.baseVersion"
                  placeholder="选择基准版本"
                  class="form-select"
                >
                  <a-select-option
                    v-for="version in completedVersions"
                    :key="version.id"
                    :value="version.id"
                  >
                    {{ version.version }} - {{ version.description }}
                  </a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="对比版本">
                <a-select
                  v-model:value="compareForm.targetVersion"
                  placeholder="选择对比版本"
                  class="form-select"
                >
                  <a-select-option
                    v-for="version in completedVersions"
                    :key="version.id"
                    :value="version.id"
                  >
                    {{ version.version }} - {{ version.description }}
                  </a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <div
          v-if="compareForm.baseVersion && compareForm.targetVersion"
          class="compare-result"
        >
          <a-table
            :columns="compareColumns"
            :data-source="compareData"
            :pagination="false"
            size="small"
            class="compare-table"
          >
            <template #attribute="{ record }">
              <span class="compare-attribute">{{ record.attribute }}</span>
            </template>
            <template #baseValue="{ record }">
              <span class="compare-value">{{ record.baseValue }}</span>
            </template>
            <template #targetValue="{ record }">
              <span class="compare-value">{{ record.targetValue }}</span>
            </template>
            <template #difference="{ record }">
              <div class="difference-wrapper">
                <span
                  class="difference-value"
                  :class="getDifferenceClass(record.difference)"
                >
                  {{ record.difference }}
                </span>
              </div>
            </template>
          </a-table>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import type {
  FormInstance,
  TableColumnsType,
  UploadFile,
} from 'ant-design-vue';
import {
  DatabaseOutlined,
  PlusOutlined,
  ReloadOutlined,
  FolderOutlined,
  UploadOutlined,
  SettingOutlined,
  CrownOutlined,
  FileOutlined,
  CloudServerOutlined,
  DownloadOutlined,
  EyeOutlined,
  MoreOutlined,
  DiffOutlined,
  CopyOutlined,
  InboxOutlined,
  DeleteOutlined,
  CloudUploadOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  LoadingOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface DatasetFile {
  id: string;
  name: string;
  size: number;
  path: string;
  type: string;
}

interface DatasetVersion {
  id: string;
  version: string;
  status: 'active' | 'processing' | 'completed' | 'failed' | 'archived';
  creator: string;
  createTime: string;
  description: string;
  tags: string[];
  isLatest: boolean;
  fileCount: number;
  size: number;
  downloadCount: number;
  checksum: string;
  files: DatasetFile[];
}

interface DatasetInfo {
  id: string;
  name: string;
  type: string;
  creator: string;
  createTime: string;
  storagePath: string;
  description: string;
  access: string;
}

interface CreateForm {
  version: string;
  tags: string[];
  source: 'upload' | 'copy' | 'merge';
  sourceVersion: string;
  mergeVersions: string[];
  description: string;
  fileList: UploadFile[];
}

interface UploadForm {
  version: string;
  description: string;
}

interface SettingsForm {
  name: string;
  type: string;
  storagePath: string;
  access: string;
  description: string;
}

interface CompareForm {
  baseVersion: string;
  targetVersion: string;
}

interface CompareData {
  attribute: string;
  baseValue: string;
  targetValue: string;
  difference: string;
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const createModalVisible = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const uploadModalVisible = ref<boolean>(false);
const settingsModalVisible = ref<boolean>(false);
const compareModalVisible = ref<boolean>(false);

const createLoading = ref<boolean>(false);
const uploadLoading = ref<boolean>(false);
const settingsLoading = ref<boolean>(false);

const filterStatus = ref<string>('');
const filterTag = ref<string>('');
const searchKeyword = ref<string>('');

const selectedVersion = ref<DatasetVersion | null>(null);
const uploadFileList = ref<UploadFile[]>([]);

// ===== 表单引用 =====
const createFormRef = ref<FormInstance>();
const settingsFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const createForm = reactive<CreateForm>({
  version: '',
  tags: [],
  source: 'upload',
  sourceVersion: '',
  mergeVersions: [],
  description: '',
  fileList: [],
});

const uploadForm = reactive<UploadForm>({
  version: '',
  description: '',
});

const settingsForm = reactive<SettingsForm>({
  name: '',
  type: '',
  storagePath: '',
  access: '',
  description: '',
});

const compareForm = reactive<CompareForm>({
  baseVersion: '',
  targetVersion: '',
});

// ===== 配置数据 =====
const STATUS_CONFIG = {
  active: { color: 'success', text: '活跃', icon: CheckCircleOutlined },
  processing: { color: 'processing', text: '处理中', icon: LoadingOutlined },
  completed: { color: 'blue', text: '已完成', icon: CheckCircleOutlined },
  failed: { color: 'error', text: '失败', icon: CloseCircleOutlined },
  archived: { color: 'default', text: '已归档', icon: InboxOutlined },
} as const;

const TAG_CONFIG = {
  stable: { color: 'green', text: '稳定版' },
  beta: { color: 'orange', text: '测试版' },
  experimental: { color: 'purple', text: '实验版' },
  hotfix: { color: 'red', text: '修复版' },
} as const;

// ===== 模拟数据 =====
const currentDataset = ref<DatasetInfo>({
  id: 'ds-001',
  name: 'CIFAR-10 图像分类数据集',
  type: '图像数据',
  creator: 'admin',
  createTime: '2024-01-15 10:30:00',
  storagePath: '/data/datasets/cifar10',
  description:
    '包含10个类别的32x32彩色图像数据集，常用于机器学习和计算机视觉研究',
  access: 'team',
});

const versions = ref<DatasetVersion[]>([
  {
    id: 'v-001',
    version: 'v3.1.0',
    status: 'completed',
    creator: 'admin',
    createTime: '2024-06-23 09:30:00',
    description: '增加数据增强和清洗后的稳定版本',
    tags: ['stable'],
    isLatest: true,
    fileCount: 60000,
    size: 2147483648,
    downloadCount: 156,
    checksum: 'sha256:a1b2c3d4e5f6...',
    files: [
      {
        id: 'f1',
        name: 'train.tar.gz',
        size: 1073741824,
        path: '/train.tar.gz',
        type: 'archive',
      },
      {
        id: 'f2',
        name: 'test.tar.gz',
        size: 536870912,
        path: '/test.tar.gz',
        type: 'archive',
      },
      {
        id: 'f3',
        name: 'labels.json',
        size: 1024,
        path: '/labels.json',
        type: 'json',
      },
      {
        id: 'f4',
        name: 'metadata.yaml',
        size: 2048,
        path: '/metadata.yaml',
        type: 'yaml',
      },
    ],
  },
  {
    id: 'v-002',
    version: 'v3.0.1',
    status: 'completed',
    creator: 'researcher',
    createTime: '2024-06-20 14:15:00',
    description: '修复标签错误的热修复版本',
    tags: ['hotfix'],
    isLatest: false,
    fileCount: 60000,
    size: 2097152000,
    downloadCount: 89,
    checksum: 'sha256:b2c3d4e5f6g7...',
    files: [
      {
        id: 'f5',
        name: 'train.tar.gz',
        size: 1048576000,
        path: '/train.tar.gz',
        type: 'archive',
      },
      {
        id: 'f6',
        name: 'test.tar.gz',
        size: 524288000,
        path: '/test.tar.gz',
        type: 'archive',
      },
      {
        id: 'f7',
        name: 'labels.json',
        size: 1024,
        path: '/labels.json',
        type: 'json',
      },
    ],
  },
  {
    id: 'v-003',
    version: 'v3.0.0',
    status: 'completed',
    creator: 'developer',
    createTime: '2024-06-15 11:00:00',
    description: '重构数据格式，优化存储结构',
    tags: ['stable'],
    isLatest: false,
    fileCount: 60000,
    size: 2147483648,
    downloadCount: 234,
    checksum: 'sha256:c3d4e5f6g7h8...',
    files: [
      {
        id: 'f8',
        name: 'images.h5',
        size: 2000000000,
        path: '/images.h5',
        type: 'hdf5',
      },
      {
        id: 'f9',
        name: 'labels.csv',
        size: 147483648,
        path: '/labels.csv',
        type: 'csv',
      },
    ],
  },
  {
    id: 'v-004',
    version: 'v2.1.0-beta',
    status: 'processing',
    creator: 'admin',
    createTime: '2024-06-22 16:45:00',
    description: '实验性功能测试版本',
    tags: ['beta', 'experimental'],
    isLatest: false,
    fileCount: 55000,
    size: 1879048192,
    downloadCount: 12,
    checksum: 'processing...',
    files: [],
  },
  {
    id: 'v-005',
    version: 'v2.0.0',
    status: 'archived',
    creator: 'admin',
    createTime: '2024-05-01 09:00:00',
    description: '旧版本数据集，已归档',
    tags: ['stable'],
    isLatest: false,
    fileCount: 50000,
    size: 1610612736,
    downloadCount: 445,
    checksum: 'sha256:d4e5f6g7h8i9...',
    files: [
      {
        id: 'f10',
        name: 'data.zip',
        size: 1610612736,
        path: '/data.zip',
        type: 'archive',
      },
    ],
  },
]);

// ===== 表单验证规则 =====
const createFormRules = {
  version: [
    { required: true, message: '请输入版本号', trigger: 'blur' },
    {
      pattern: /^v?\d+\.\d+\.\d+(-\w+)?$/,
      message: '版本号格式不正确，例如: v1.0.0 或 1.0.0-beta',
      trigger: 'blur',
    },
  ],
  description: [
    { required: true, message: '请输入版本描述', trigger: 'blur' },
    { min: 10, message: '描述至少10个字符', trigger: 'blur' },
  ],
  sourceVersion: [
    { required: true, message: '请选择源版本', trigger: 'change' },
  ],
  mergeVersions: [
    { required: true, message: '请选择要合并的版本', trigger: 'change' },
    {
      type: 'array',
      min: 2,
      message: '至少选择2个版本进行合并',
      trigger: 'change',
    },
  ],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<DatasetVersion> = [
  {
    title: '版本号',
    key: 'version',
    width: 150,
    slots: { customRender: 'version' },
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    slots: { customRender: 'status' },
  },
  {
    title: '标签',
    key: 'tags',
    width: 150,
    slots: { customRender: 'tags' },
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
    width: 100,
  },
  {
    title: '文件信息',
    key: 'fileInfo',
    width: 150,
    slots: { customRender: 'fileInfo' },
  },
  {
    title: '下载次数',
    dataIndex: 'downloadCount',
    key: 'downloadCount',
    width: 100,
    sorter: (a: DatasetVersion, b: DatasetVersion) =>
      a.downloadCount - b.downloadCount,
  },
  {
    title: '创建时间',
    key: 'createTime',
    width: 150,
    slots: { customRender: 'createTime' },
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
    ellipsis: true,
    width: 200,
  },
  {
    title: '操作',
    key: 'action',
    width: 180,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

const fileColumns: TableColumnsType<DatasetFile> = [
  {
    title: '文件名',
    key: 'fileName',
    slots: { customRender: 'fileName' },
  },
  {
    title: '大小',
    key: 'fileSize',
    width: 120,
    slots: { customRender: 'fileSize' },
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type',
    width: 80,
  },
  {
    title: '操作',
    key: 'fileAction',
    width: 80,
    slots: { customRender: 'fileAction' },
  },
];

const compareColumns: TableColumnsType<CompareData> = [
  {
    title: '属性',
    key: 'attribute',
    width: 150,
    slots: { customRender: 'attribute' },
  },
  {
    title: '基准版本',
    key: 'baseValue',
    slots: { customRender: 'baseValue' },
  },
  {
    title: '对比版本',
    key: 'targetValue',
    slots: { customRender: 'targetValue' },
  },
  {
    title: '差异',
    key: 'difference',
    width: 150,
    slots: { customRender: 'difference' },
  },
];

// ===== 分页配置 =====
const paginationConfig = {
  total: computed(() => filteredVersions.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const filteredVersions = computed(() => {
  let result = versions.value;

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (filterTag.value) {
    result = result.filter((item) => item.tags.includes(filterTag.value));
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter(
      (item) =>
        item.version.toLowerCase().includes(keyword) ||
        item.description.toLowerCase().includes(keyword),
    );
  }

  return result;
});

const completedVersions = computed(() => {
  return versions.value.filter((v) => v.status === 'completed');
});

const compareData = computed<CompareData[]>(() => {
  if (!compareForm.baseVersion || !compareForm.targetVersion) {
    return [];
  }

  const baseVersion = versions.value.find(
    (v) => v.id === compareForm.baseVersion,
  );
  const targetVersion = versions.value.find(
    (v) => v.id === compareForm.targetVersion,
  );

  if (!baseVersion || !targetVersion) {
    return [];
  }

  return [
    {
      attribute: '文件数量',
      baseValue: baseVersion.fileCount.toString(),
      targetValue: targetVersion.fileCount.toString(),
      difference: (targetVersion.fileCount - baseVersion.fileCount).toString(),
    },
    {
      attribute: '数据大小',
      baseValue: formatFileSize(baseVersion.size),
      targetValue: formatFileSize(targetVersion.size),
      difference: formatFileSize(targetVersion.size - baseVersion.size),
    },
    {
      attribute: '下载次数',
      baseValue: baseVersion.downloadCount.toString(),
      targetValue: targetVersion.downloadCount.toString(),
      difference: (
        targetVersion.downloadCount - baseVersion.downloadCount
      ).toString(),
    },
  ];
});

// ===== 工具函数 =====
const getStatusColor = (status: string): string => {
  return (
    STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.color || 'default'
  );
};

const getStatusIcon = (status: string) => {
  return (
    STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.icon ||
    ClockCircleOutlined
  );
};

const getStatusText = (status: string): string => {
  return STATUS_CONFIG[status as keyof typeof STATUS_CONFIG]?.text || status;
};

const getTagColor = (tag: string): string => {
  return TAG_CONFIG[tag as keyof typeof TAG_CONFIG]?.color || 'blue';
};

const getTagText = (tag: string): string => {
  return TAG_CONFIG[tag as keyof typeof TAG_CONFIG]?.text || tag;
};

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]!;
};

const formatRelativeTime = (time: string): string => {
  const now = new Date();
  const target = new Date(time);
  const diffMs = now.getTime() - target.getTime();
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));

  if (diffHours < 1) {
    const diffMinutes = Math.floor(diffMs / (1000 * 60));
    return `${diffMinutes} 分钟前`;
  } else if (diffHours < 24) {
    return `${diffHours} 小时前`;
  } else {
    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays} 天前`;
  }
};

const getDifferenceClass = (difference: string): string => {
  const num = parseFloat(difference);
  if (num > 0) return 'positive';
  if (num < 0) return 'negative';
  return 'neutral';
};

const beforeUpload = (): boolean => {
  return false; // 阻止自动上传
};

// ===== 事件处理函数 =====
const showCreateModal = (): void => {
  createModalVisible.value = true;
};

const showUploadModal = (): void => {
  uploadModalVisible.value = true;
};

const showSettingsModal = (): void => {
  // 填充当前设置
  Object.assign(settingsForm, currentDataset.value);
  settingsModalVisible.value = true;
};

const handleCreateSubmit = async (): Promise<void> => {
  try {
    await createFormRef.value?.validate();
    createLoading.value = true;

    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 2000));

    const newVersion: DatasetVersion = {
      id: `v-${Date.now()}`,
      version: createForm.version,
      status: 'processing',
      creator: 'current-user',
      createTime: new Date().toLocaleString(),
      description: createForm.description,
      tags: createForm.tags,
      isLatest: false,
      fileCount: 0,
      size: 0,
      downloadCount: 0,
      checksum: 'generating...',
      files: [],
    };

    // 如果有新版本且当前没有最新版本标记，则设为最新
    if (
      versions.value.length === 0 ||
      !versions.value.some((v) => v.isLatest)
    ) {
      newVersion.isLatest = true;
    }

    versions.value.unshift(newVersion);
    createModalVisible.value = false;
    message.success('版本创建成功');

    // 重置表单
    createFormRef.value?.resetFields();
    createForm.fileList = [];
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    createLoading.value = false;
  }
};

const handleCreateCancel = (): void => {
  createModalVisible.value = false;
  createFormRef.value?.resetFields();
  createForm.fileList = [];
};

const handleUploadSubmit = async (): Promise<void> => {
  if (uploadFileList.value.length === 0) {
    message.error('请选择要上传的文件');
    return;
  }

  uploadLoading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 3000));

    const newVersion: DatasetVersion = {
      id: `v-${Date.now()}`,
      version: uploadForm.version || `v${versions.value.length + 1}.0.0`,
      status: 'processing',
      creator: 'current-user',
      createTime: new Date().toLocaleString(),
      description: uploadForm.description || '通过文件上传创建的版本',
      tags: ['stable'],
      isLatest: true,
      fileCount: uploadFileList.value.length,
      size: uploadFileList.value.reduce(
        (total, file) => total + (file.size || 0),
        0,
      ),
      downloadCount: 0,
      checksum: 'generating...',
      files: uploadFileList.value.map((file, index) => ({
        id: `f-${Date.now()}-${index}`,
        name: file.name,
        size: file.size || 0,
        path: `/${file.name}`,
        type: file.name.split('.').pop() || 'unknown',
      })),
    };

    // 取消其他版本的最新标记
    versions.value.forEach((v) => (v.isLatest = false));

    versions.value.unshift(newVersion);
    uploadModalVisible.value = false;
    message.success('文件上传成功，正在处理中');

    // 重置表单
    uploadFileList.value = [];
    uploadForm.version = '';
    uploadForm.description = '';
  } catch (error) {
    message.error('上传失败');
  } finally {
    uploadLoading.value = false;
  }
};

const handleUploadCancel = (): void => {
  uploadModalVisible.value = false;
  uploadFileList.value = [];
  uploadForm.version = '';
  uploadForm.description = '';
};

const handleSettingsSubmit = async (): Promise<void> => {
  settingsLoading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));

    Object.assign(currentDataset.value, settingsForm);
    settingsModalVisible.value = false;
    message.success('设置保存成功');
  } catch (error) {
    message.error('保存失败');
  } finally {
    settingsLoading.value = false;
  }
};

const handleSettingsCancel = (): void => {
  settingsModalVisible.value = false;
  settingsFormRef.value?.resetFields();
};

const downloadVersion = (record: DatasetVersion): void => {
  message.success(`开始下载版本 ${record.version}`);
  // 模拟下载逻辑
  const downloadCount = record.downloadCount + 1;
  const index = versions.value.findIndex((item) => item.id === record.id);
  if (index !== -1) {
    versions.value[index]!.downloadCount = downloadCount;
  }
};

const downloadFile = (file: DatasetFile): void => {
  message.success(`开始下载文件 ${file.name}`);
};

const viewDetails = (record: DatasetVersion): void => {
  selectedVersion.value = record;
  detailModalVisible.value = true;
};

const handleMenuAction = (key: string, record: DatasetVersion): void => {
  const actions = {
    setLatest: () => handleSetLatest(record),
    compare: () => handleCompare(record),
    clone: () => handleClone(record),
    archive: () => handleArchive(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handleSetLatest = async (record: DatasetVersion): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));

    // 取消所有版本的最新标记
    versions.value.forEach((v) => (v.isLatest = false));

    // 设置当前版本为最新
    const index = versions.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      versions.value[index]!.isLatest = true;
    }

    message.success(`版本 ${record.version} 已设为最新版本`);
  } catch (error) {
    message.error('设置失败');
  } finally {
    loading.value = false;
  }
};

const handleCompare = (record: DatasetVersion): void => {
  compareForm.targetVersion = record.id;
  compareModalVisible.value = true;
};

const handleClone = (record: DatasetVersion): void => {
  createForm.source = 'copy';
  createForm.sourceVersion = record.id;
  createForm.version = `${record.version}-copy`;
  createForm.description = `克隆自版本 ${record.version}`;
  createForm.tags = [...record.tags];
  createModalVisible.value = true;
};

const handleArchive = async (record: DatasetVersion): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const index = versions.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      versions.value[index]!.status = 'archived';
      if (versions.value[index]!.isLatest) {
        versions.value[index]!.isLatest = false;
        // 找到下一个可用版本设为最新
        const nextVersion = versions.value.find(
          (v) => v.status === 'completed' && v.id !== record.id,
        );
        if (nextVersion) {
          nextVersion.isLatest = true;
        }
      }
    }
    message.success(`版本 ${record.version} 已归档`);
  } catch (error) {
    message.error('归档失败');
  } finally {
    loading.value = false;
  }
};

const handleDelete = (record: DatasetVersion): void => {
  const deleteConfirm = () => {
    const index = versions.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      const wasLatest = versions.value[index]!.isLatest;
      versions.value.splice(index, 1);

      // 如果删除的是最新版本，设置下一个版本为最新
      if (wasLatest && versions.value.length > 0) {
        const nextVersion = versions.value.find(
          (v) => v.status === 'completed',
        );
        if (nextVersion) {
          nextVersion.isLatest = true;
        }
      }

      message.success(`版本 ${record.version} 删除成功`);
    }
  };

  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除版本 "${record.version}" 吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      type: 'warning',
      onOk: deleteConfirm,
    });
  });
};

const refreshData = async (): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    message.success('数据刷新成功');
  } catch (error) {
    message.error('刷新失败');
  } finally {
    loading.value = false;
  }
};

const handleFilterChange = (): void => {
  // 筛选变化时的处理逻辑
};

const handleSearch = (): void => {
  // 搜索处理逻辑
};

const handleSearchChange = (): void => {
  // 搜索输入变化时的处理逻辑
};

// ===== 生命周期 =====
onMounted(() => {
  refreshData();
});
</script>

<style scoped>
/* ===== 基础样式 ===== */
.dataset-version-container {
  padding: 24px;
  min-height: 100vh;
}

/* ===== 卡片样式 ===== */
.glass-card {
  border-radius: 8px !important;
}

/* ===== 页面头部 ===== */
.page-header {
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24px;
}

.title-section {
  flex: 1;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 32px;
  color: #1890ff;
}

.page-description {
  font-size: 16px;
  margin: 0;
}

/* ===== 数据集信息卡片 ===== */
.dataset-info-section {
  margin-bottom: 24px;
}

.dataset-info-card {
  border-radius: 8px !important;
}

.dataset-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.dataset-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.dataset-title h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.dataset-type-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
  border-radius: 6px !important;
}

.dataset-actions {
  display: flex;
  gap: 8px;
}

.upload-btn,
.settings-btn {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.upload-btn:hover,
.settings-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

.dataset-meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.meta-label {
  font-weight: 500;
  opacity: 0.8;
}

.meta-value {
  font-weight: 600;
}

.dataset-description {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.description-label {
  font-weight: 500;
  opacity: 0.8;
  flex-shrink: 0;
}

.description-content {
  flex: 1;
  line-height: 1.5;
}

/* ===== 按钮样式 ===== */
.create-btn {
  border: none !important;
  height: 40px !important;
  padding: 0 24px !important;
  font-size: 14px !important;
  font-weight: 500 !important;
  border-radius: 6px !important;
  transition: all 0.3s ease !important;
}

.create-btn:hover {
  background: #1890ff !important;
  transform: translateY(-1px);
}

/* ===== 筛选器样式 ===== */
.filter-section {
  margin-bottom: 24px;
}

.filter-card {
  border-radius: 8px !important;
}

.filter-select,
.search-input,
.refresh-btn {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.refresh-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

/* ===== 表格样式 ===== */
.table-section {
  margin-bottom: 24px;
}

.table-card {
  border-radius: 8px !important;
}

.sci-fi-table :deep(.ant-table-thead > tr > th) {
  font-weight: 600 !important;
}

/* ===== 版本相关样式 ===== */
.version-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.version-number {
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-weight: 600;
  font-size: 13px;
}

.latest-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border-radius: 4px !important;
  font-size: 11px !important;
  padding: 2px 6px !important;
  font-weight: 500 !important;
}

/* ===== 状态标签 ===== */
.status-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border-radius: 6px !important;
  font-weight: 500 !important;
  padding: 4px 8px !important;
}

.status-icon {
  font-size: 12px;
}

.status-indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.indicator-active {
  background: #52c41a;
}

.indicator-processing {
  background: #1890ff;
}

.indicator-completed {
  background: #1890ff;
}

.indicator-failed {
  background: #ff4d4f;
}

.indicator-archived {
  background: #8c8c8c;
}

/* ===== 标签样式 ===== */
.tags-wrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.version-tag {
  border-radius: 4px !important;
  font-size: 11px !important;
  padding: 2px 6px !important;
  font-weight: 500 !important;
}

/* ===== 文件信息 ===== */
.file-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.file-item:hover {
  color: #1890ff;
}

.file-icon {
  font-size: 12px;
  color: #1890ff;
}

.file-label {
  font-weight: 500;
}

.file-value {
  font-weight: 600;
}

/* ===== 时间显示 ===== */
.time-text {
  font-size: 12px;
}

/* ===== 操作按钮 ===== */
.action-buttons {
  display: flex;
  gap: 4px;
}

.action-btn {
  border: none !important;
  background: transparent !important;
  border-radius: 4px !important;
  padding: 4px 8px !important;
  height: auto !important;
  font-size: 12px !important;
  transition: all 0.3s ease !important;
}

.action-btn:hover {
  color: #1890ff !important;
}

.action-btn:disabled {
  color: #bfbfbf !important;
  background: transparent !important;
}

.action-menu {
  border-radius: 8px !important;
}

.action-menu :deep(.ant-menu-item) {
  border-radius: 4px !important;
  margin: 2px !important;
  transition: all 0.3s ease;
}

.action-menu :deep(.ant-menu-item:hover) {
  color: #1890ff !important;
}

.danger-item {
  color: #ff4d4f !important;
}

.danger-item:hover {
  color: #ff4d4f !important;
}

/* ===== 模态框样式 ===== */
.sci-fi-modal :deep(.ant-modal-content) {
  border-radius: 8px !important;
}

.sci-fi-modal :deep(.ant-modal-header) {
  border-radius: 8px 8px 0 0 !important;
}

.sci-fi-modal :deep(.ant-modal-title) {
  font-weight: 600 !important;
  font-size: 16px !important;
}

/* ===== 表单样式 ===== */
.create-form :deep(.ant-form-item-label > label),
.settings-form :deep(.ant-form-item-label > label) {
  font-weight: 500 !important;
}

.form-input,
.form-select,
.form-textarea,
.form-input-number {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.source-radio {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.upload-dragger {
  border-radius: 8px !important;
  transition: all 0.3s ease;
}

.upload-content {
  margin-bottom: 16px;
}

.upload-form {
  margin-top: 16px;
}

/* ===== 详情页样式 ===== */
.detail-content {
  max-height: 600px;
  overflow-y: auto;
}

.version-detail {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tags-detail {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.file-divider {
  font-weight: 500 !important;
  margin: 24px 0 16px 0 !important;
}

.file-list-container {
  margin-top: 16px;
}

.file-table {
  border-radius: 8px !important;
}

.file-name-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-type-icon {
  color: #1890ff;
  font-size: 14px;
}

.file-name {
  font-weight: 500;
}

.file-size {
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 12px;
}

.file-action-btn {
  border: none !important;
  background: transparent !important;
  border-radius: 4px !important;
  padding: 4px 8px !important;
  height: auto !important;
  font-size: 12px !important;
  transition: all 0.3s ease !important;
}

.file-action-btn:hover {
  color: #1890ff !important;
}

/* ===== 版本对比样式 ===== */
.compare-modal :deep(.ant-modal-body) {
  padding: 24px;
}

.compare-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.compare-selector {
  border-radius: 8px;
  padding: 16px;
}

.compare-result {
  flex: 1;
}

.compare-table {
  border-radius: 8px !important;
}

.compare-attribute {
  font-weight: 600;
}

.compare-value {
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 13px;
}

.difference-wrapper {
  display: flex;
  align-items: center;
}

.difference-value {
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 13px;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 4px;
}

.difference-value.positive {
  color: #52c41a;
}

.difference-value.negative {
  color: #ff4d4f;
}

.difference-value.neutral {
  color: #8c8c8c;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .dataset-version-container {
    padding: 16px;
  }

  .header-content {
    flex-direction: column;
    gap: 16px;
  }

  .page-title {
    font-size: 24px;
  }

  .title-icon {
    font-size: 28px;
  }

  .action-section {
    align-self: stretch;
  }

  .create-btn,
  .refresh-btn {
    width: 100% !important;
    justify-content: center !important;
  }

  .refresh-btn-col {
    margin-top: 12px;
  }

  .dataset-header {
    flex-direction: column;
    gap: 16px;
  }

  .dataset-actions {
    align-self: stretch;
  }

  .upload-btn,
  .settings-btn {
    flex: 1;
    justify-content: center;
  }

  .dataset-meta {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .sci-fi-modal :deep(.ant-modal) {
    margin: 16px !important;
    max-width: calc(100vw - 32px) !important;
  }
}

@media (max-width: 480px) {
  .page-title {
    font-size: 20px;
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }

  .title-icon {
    font-size: 24px;
  }

  .dataset-title {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }

  .file-info {
    gap: 2px;
  }

  .file-item {
    font-size: 11px;
    padding: 2px 4px;
  }

  .action-buttons {
    flex-direction: column;
    gap: 2px;
  }

  .action-btn {
    font-size: 11px !important;
    padding: 3px 6px !important;
  }

  .tags-wrapper {
    gap: 2px;
  }

  .version-tag {
    font-size: 10px !important;
    padding: 1px 4px !important;
  }
}

/* ===== 滚动条样式 ===== */
::-webkit-scrollbar {
  width: 4px;
}

::-webkit-scrollbar-track,
::-webkit-scrollbar-thumb {
  border-radius: 2px;
}
</style>
