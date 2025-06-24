<template>
  <div class="dataset-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <DatabaseOutlined class="title-icon" />
            <span class="title-text">数据集管理</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">管理和组织您的训练数据集</span>
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
            创建数据集
          </a-button>
        </div>
      </div>
    </div>

    <!-- 筛选器 -->
    <div class="filter-section">
      <a-card class="filter-card glass-card" :bordered="false">
        <a-row :gutter="16" align="middle">
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterType"
              placeholder="选择数据类型"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部类型</a-select-option>
              <a-select-option value="image">图像</a-select-option>
              <a-select-option value="text">文本</a-select-option>
              <a-select-option value="audio">音频</a-select-option>
              <a-select-option value="video">视频</a-select-option>
            </a-select>
          </a-col>
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
              <a-select-option value="ready">就绪</a-select-option>
              <a-select-option value="processing">处理中</a-select-option>
              <a-select-option value="error">错误</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="16" :md="8" :lg="8">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索数据集名称或描述"
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

    <!-- 数据表格 -->
    <div class="table-section">
      <a-card class="table-card glass-card" :bordered="false">
        <a-table
          :columns="columns"
          :data-source="filteredDatasets"
          :loading="loading"
          :pagination="paginationConfig"
          row-key="id"
          size="middle"
          :scroll="{ x: 'max-content' }"
          class="sci-fi-table"
        >
          <!-- 数据类型列 -->
          <template #dataType="{ record }">
            <a-tag :color="getTypeColor(record.dataType)" class="type-tag">
              <component :is="getTypeIcon(record.dataType)" class="type-icon" />
              {{ getTypeText(record.dataType) }}
            </a-tag>
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

          <!-- 大小列 -->
          <template #size="{ record }">
            <div class="size-info">
              <span class="size-value">{{ formatSize(record.size) }}</span>
              <span class="file-count">{{ record.fileCount }} 文件</span>
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
                @click="downloadDataset(record)"
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
                    <a-menu-item key="preview">
                      <EyeOutlined />
                      预览
                    </a-menu-item>
                    <a-menu-item key="edit">
                      <EditOutlined />
                      编辑
                    </a-menu-item>
                    <a-menu-item key="copy">
                      <CopyOutlined />
                      复制
                    </a-menu-item>
                    <a-menu-item key="version">
                      <BranchesOutlined />
                      版本管理
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

    <!-- 创建数据集模态框 -->
    <a-modal
      v-model:open="createModalVisible"
      title="创建数据集"
      width="800px"
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
            <a-form-item label="数据集名称" name="name">
              <a-input
                v-model:value="createForm.name"
                placeholder="请输入数据集名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="数据类型" name="dataType">
              <a-select
                v-model:value="createForm.dataType"
                placeholder="选择数据类型"
                class="form-select"
              >
                <a-select-option value="image">图像</a-select-option>
                <a-select-option value="text">文本</a-select-option>
                <a-select-option value="audio">音频</a-select-option>
                <a-select-option value="video">视频</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="createForm.description"
            placeholder="请输入数据集描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>

        <a-divider class="form-divider">数据来源</a-divider>

        <a-form-item label="数据来源方式" name="sourceType">
          <a-radio-group
            v-model:value="createForm.sourceType"
            class="source-radio-group"
          >
            <a-radio value="upload">本地上传</a-radio>
            <a-radio value="url">URL 导入</a-radio>
            <a-radio value="s3">S3 存储桶</a-radio>
          </a-radio-group>
        </a-form-item>

        <!-- 本地上传 -->
        <div v-if="createForm.sourceType === 'upload'">
          <a-form-item label="文件上传" name="files">
            <a-upload-dragger
              v-model:fileList="createForm.fileList"
              multiple
              :before-upload="beforeUpload"
              class="upload-dragger"
            >
              <p class="ant-upload-drag-icon">
                <InboxOutlined />
              </p>
              <p class="ant-upload-text">点击或拖拽文件到此区域上传</p>
              <p class="ant-upload-hint">
                支持单个或批量上传，严格禁止上传公司数据或其他敏感信息
              </p>
            </a-upload-dragger>
          </a-form-item>
        </div>

        <!-- URL 导入 -->
        <div v-if="createForm.sourceType === 'url'">
          <a-form-item label="数据 URL" name="dataUrl">
            <a-input
              v-model:value="createForm.dataUrl"
              placeholder="请输入数据 URL"
              class="form-input"
            />
          </a-form-item>
        </div>

        <!-- S3 存储桶 -->
        <div v-if="createForm.sourceType === 's3'">
          <a-row :gutter="16">
            <a-col :xs="24" :sm="12">
              <a-form-item label="存储桶名称" name="s3Bucket">
                <a-input
                  v-model:value="createForm.s3Bucket"
                  placeholder="请输入存储桶名称"
                  class="form-input"
                />
              </a-form-item>
            </a-col>
            <a-col :xs="24" :sm="12">
              <a-form-item label="对象路径" name="s3Path">
                <a-input
                  v-model:value="createForm.s3Path"
                  placeholder="请输入对象路径"
                  class="form-input"
                />
              </a-form-item>
            </a-col>
          </a-row>
        </div>

        <a-divider class="form-divider">标签和权限</a-divider>

        <a-form-item label="标签" name="tags">
          <a-select
            v-model:value="createForm.tags"
            mode="tags"
            placeholder="添加标签"
            class="form-select"
          >
            <a-select-option value="训练">训练</a-select-option>
            <a-select-option value="测试">测试</a-select-option>
            <a-select-option value="验证">验证</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="访问权限" name="accessLevel">
          <a-radio-group v-model:value="createForm.accessLevel">
            <a-radio value="private">私有</a-radio>
            <a-radio value="team">团队共享</a-radio>
            <a-radio value="public">公开</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="数据集详情"
      width="900px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedDataset" class="detail-content">
        <a-descriptions
          :column="{ xs: 1, sm: 2 }"
          bordered
          class="detail-descriptions"
        >
          <a-descriptions-item label="名称">
            {{ selectedDataset.name }}
          </a-descriptions-item>
          <a-descriptions-item label="数据类型">
            <a-tag :color="getTypeColor(selectedDataset.dataType)">
              <component :is="getTypeIcon(selectedDataset.dataType)" />
              {{ getTypeText(selectedDataset.dataType) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag
              :color="getStatusColor(selectedDataset.status)"
              class="status-tag"
            >
              <component :is="getStatusIcon(selectedDataset.status)" />
              {{ getStatusText(selectedDataset.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="创建者">
            {{ selectedDataset.creator }}
          </a-descriptions-item>
          <a-descriptions-item label="大小">
            {{ formatSize(selectedDataset.size) }}
          </a-descriptions-item>
          <a-descriptions-item label="文件数量">
            {{ selectedDataset.fileCount }} 个
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">
            {{ selectedDataset.createTime }}
          </a-descriptions-item>
          <a-descriptions-item label="最后更新">
            {{ selectedDataset.updateTime }}
          </a-descriptions-item>
          <a-descriptions-item label="版本">
            v{{ selectedDataset.version }}
          </a-descriptions-item>
          <a-descriptions-item label="访问权限">
            {{ getAccessLevelText(selectedDataset.accessLevel) }}
          </a-descriptions-item>
          <a-descriptions-item label="标签" :span="2">
            <a-space>
              <a-tag
                v-for="tag in selectedDataset.tags"
                :key="tag"
                color="blue"
              >
                {{ tag }}
              </a-tag>
            </a-space>
          </a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">
            {{ selectedDataset.description || '暂无描述' }}
          </a-descriptions-item>
        </a-descriptions>

        <!-- 数据预览 -->
        <a-divider class="preview-divider">数据预览</a-divider>
        <div class="preview-container">
          <a-tabs v-model:activeKey="previewActiveKey" class="preview-tabs">
            <a-tab-pane key="structure" tab="数据结构">
              <div class="structure-view">
                <a-tree
                  :tree-data="dataStructure"
                  :show-icon="true"
                  class="structure-tree"
                />
              </div>
            </a-tab-pane>
            <a-tab-pane key="sample" tab="样本预览">
              <div class="sample-view">
                <a-list
                  :data-source="sampleData"
                  size="small"
                  class="sample-list"
                >
                  <template #renderItem="{ item }">
                    <a-list-item>
                      <a-list-item-meta
                        :title="item.name"
                        :description="item.path"
                      >
                        <template #avatar>
                          <a-avatar
                            :src="item.preview"
                            :icon="getFileIcon(item.type)"
                            shape="square"
                          />
                        </template>
                      </a-list-item-meta>
                      <span class="file-size">{{ formatSize(item.size) }}</span>
                    </a-list-item>
                  </template>
                </a-list>
              </div>
            </a-tab-pane>
            <a-tab-pane key="stats" tab="统计信息">
              <div class="stats-view">
                <a-row :gutter="16">
                  <a-col :xs="24" :sm="12" :md="6">
                    <a-statistic
                      title="总文件数"
                      :value="selectedDataset.fileCount"
                    />
                  </a-col>
                  <a-col :xs="24" :sm="12" :md="6">
                    <a-statistic
                      title="总大小"
                      :value="formatSize(selectedDataset.size)"
                    />
                  </a-col>
                  <a-col :xs="24" :sm="12" :md="6">
                    <a-statistic
                      title="平均文件大小"
                      :value="
                        formatSize(
                          selectedDataset.size / selectedDataset.fileCount,
                        )
                      "
                    />
                  </a-col>
                  <a-col :xs="24" :sm="12" :md="6">
                    <a-statistic
                      title="版本号"
                      :value="selectedDataset.version"
                      prefix="v"
                    />
                  </a-col>
                </a-row>
              </div>
            </a-tab-pane>
          </a-tabs>
        </div>
      </div>
    </a-modal>

    <!-- 编辑模态框 -->
    <a-modal
      v-model:open="editModalVisible"
      title="编辑数据集"
      width="600px"
      :confirm-loading="editLoading"
      @ok="handleEditSubmit"
      @cancel="handleEditCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="editFormRef"
        :model="editForm"
        :rules="editFormRules"
        layout="vertical"
      >
        <a-form-item label="数据集名称" name="name">
          <a-input
            v-model:value="editForm.name"
            placeholder="请输入数据集名称"
            class="form-input"
          />
        </a-form-item>
        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="editForm.description"
            placeholder="请输入数据集描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>
        <a-form-item label="标签" name="tags">
          <a-select
            v-model:value="editForm.tags"
            mode="tags"
            placeholder="添加标签"
            class="form-select"
          >
            <a-select-option value="训练">训练</a-select-option>
            <a-select-option value="测试">测试</a-select-option>
            <a-select-option value="验证">验证</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="访问权限" name="accessLevel">
          <a-radio-group v-model:value="editForm.accessLevel">
            <a-radio value="private">私有</a-radio>
            <a-radio value="team">团队共享</a-radio>
            <a-radio value="public">公开</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 版本管理模态框 -->
    <a-modal
      v-model:open="versionModalVisible"
      title="版本管理"
      width="800px"
      :footer="null"
      class="sci-fi-modal"
    >
      <div class="version-content">
        <div class="version-header">
          <a-button
            type="primary"
            @click="createVersion"
            class="create-version-btn"
          >
            <PlusOutlined />
            创建新版本
          </a-button>
        </div>
        <a-timeline class="version-timeline">
          <a-timeline-item
            v-for="version in versions"
            :key="version.id"
            :color="version.current ? 'green' : 'blue'"
          >
            <template #dot>
              <ClockCircleOutlined v-if="!version.current" />
              <CheckCircleOutlined v-else style="color: #52c41a" />
            </template>
            <div class="version-item">
              <div class="version-header-info">
                <span class="version-number">v{{ version.version }}</span>
                <a-tag v-if="version.current" color="success">当前版本</a-tag>
                <span class="version-time">{{ version.createTime }}</span>
              </div>
              <div class="version-description">{{ version.description }}</div>
              <div class="version-stats">
                <span class="version-stat">{{ version.fileCount }} 文件</span>
                <span class="version-stat">{{ formatSize(version.size) }}</span>
                <span class="version-creator"
                  >创建者: {{ version.creator }}</span
                >
              </div>
              <div class="version-actions" v-if="!version.current">
                <a-button size="small" @click="switchVersion(version)">
                  切换到此版本
                </a-button>
                <a-button size="small" danger @click="deleteVersion(version)">
                  删除版本
                </a-button>
              </div>
            </div>
          </a-timeline-item>
        </a-timeline>
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
  EyeOutlined,
  MoreOutlined,
  EditOutlined,
  CopyOutlined,
  DeleteOutlined,
  DownloadOutlined,
  BranchesOutlined,
  InboxOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  FileImageOutlined,
  FileTextOutlined,
  SoundOutlined,
  VideoCameraOutlined,
  FolderOutlined,
  FileOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface DatasetItem {
  id: string;
  name: string;
  dataType: 'image' | 'text' | 'audio' | 'video';
  status: 'ready' | 'processing' | 'error';
  creator: string;
  description?: string;
  size: number;
  fileCount: number;
  createTime: string;
  updateTime: string;
  version: string;
  tags: string[];
  accessLevel: 'private' | 'team' | 'public';
}

interface CreateForm {
  name: string;
  dataType: 'image' | 'text' | 'audio' | 'video' | '';
  description: string;
  sourceType: 'upload' | 'url' | 's3';
  fileList: UploadFile[];
  dataUrl: string;
  s3Bucket: string;
  s3Path: string;
  tags: string[];
  accessLevel: 'private' | 'team' | 'public';
}

interface EditForm {
  name: string;
  description: string;
  tags: string[];
  accessLevel: 'private' | 'team' | 'public';
}

interface VersionItem {
  id: string;
  version: string;
  description: string;
  createTime: string;
  creator: string;
  fileCount: number;
  size: number;
  current: boolean;
}

interface FileItem {
  name: string;
  path: string;
  type: string;
  size: number;
  preview?: string;
}

interface TreeNode {
  title: string;
  key: string;
  icon?: any;
  children?: TreeNode[];
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const createModalVisible = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const editModalVisible = ref<boolean>(false);
const versionModalVisible = ref<boolean>(false);
const createLoading = ref<boolean>(false);
const editLoading = ref<boolean>(false);

const filterType = ref<string>('');
const filterStatus = ref<string>('');
const searchKeyword = ref<string>('');

const selectedDataset = ref<DatasetItem | null>(null);
const previewActiveKey = ref<string>('structure');

// ===== 表单引用 =====
const createFormRef = ref<FormInstance>();
const editFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const createForm = reactive<CreateForm>({
  name: '',
  dataType: '',
  description: '',
  sourceType: 'upload',
  fileList: [],
  dataUrl: '',
  s3Bucket: '',
  s3Path: '',
  tags: [],
  accessLevel: 'private',
});

const editForm = reactive<EditForm>({
  name: '',
  description: '',
  tags: [],
  accessLevel: 'private',
});

// ===== 配置数据 =====
const TYPE_CONFIG = {
  image: { color: 'blue', text: '图像', icon: FileImageOutlined },
  text: { color: 'green', text: '文本', icon: FileTextOutlined },
  audio: { color: 'orange', text: '音频', icon: SoundOutlined },
  video: { color: 'purple', text: '视频', icon: VideoCameraOutlined },
} as const;

const STATUS_CONFIG = {
  ready: { color: 'success', text: '就绪', icon: CheckCircleOutlined },
  processing: {
    color: 'processing',
    text: '处理中',
    icon: ClockCircleOutlined,
  },
  error: { color: 'error', text: '错误', icon: CloseCircleOutlined },
} as const;

// ===== 模拟数据 =====
const datasets = ref<DatasetItem[]>([
  {
    id: 'ds-001',
    name: 'imagenet-subset',
    dataType: 'image',
    status: 'ready',
    creator: 'admin',
    description: 'ImageNet 数据集子集，包含1000个类别的图像数据',
    size: 5368709120, // 5GB
    fileCount: 50000,
    createTime: '2024-06-20 14:30:00',
    updateTime: '2024-06-23 10:15:00',
    version: '1.2',
    tags: ['训练', '分类', '大规模'],
    accessLevel: 'team',
  },
  {
    id: 'ds-002',
    name: 'text-sentiment-analysis',
    dataType: 'text',
    status: 'ready',
    creator: 'researcher',
    description: '情感分析文本数据集，包含正面、负面、中性情感标注',
    size: 1073741824, // 1GB
    fileCount: 100000,
    createTime: '2024-06-19 09:20:00',
    updateTime: '2024-06-22 16:45:00',
    version: '2.0',
    tags: ['文本', '情感分析', '标注'],
    accessLevel: 'public',
  },
  {
    id: 'ds-003',
    name: 'speech-recognition-corpus',
    dataType: 'audio',
    status: 'processing',
    creator: 'developer',
    description: '语音识别语料库，多语言音频数据',
    size: 10737418240, // 10GB
    fileCount: 25000,
    createTime: '2024-06-21 11:00:00',
    updateTime: '2024-06-23 08:30:00',
    version: '1.0',
    tags: ['语音识别', '多语言'],
    accessLevel: 'private',
  },
  {
    id: 'ds-004',
    name: 'video-action-recognition',
    dataType: 'video',
    status: 'ready',
    creator: 'admin',
    description: '视频动作识别数据集，包含多种人体动作标注',
    size: 21474836480, // 20GB
    fileCount: 5000,
    createTime: '2024-06-18 16:15:00',
    updateTime: '2024-06-21 14:20:00',
    version: '1.1',
    tags: ['视频', '动作识别', '标注'],
    accessLevel: 'team',
  },
  {
    id: 'ds-005',
    name: 'medical-image-segmentation',
    dataType: 'image',
    status: 'error',
    creator: 'researcher',
    description: '医学图像分割数据集，CT扫描图像及标注',
    size: 3221225472, // 3GB
    fileCount: 1500,
    createTime: '2024-06-22 10:30:00',
    updateTime: '2024-06-23 09:00:00',
    version: '1.0',
    tags: ['医学', '图像分割'],
    accessLevel: 'private',
  },
]);

const versions = ref<VersionItem[]>([
  {
    id: 'v-001',
    version: '1.2',
    description: '添加了新的数据类别，修复了标注错误',
    createTime: '2024-06-23 10:15:00',
    creator: 'admin',
    fileCount: 50000,
    size: 5368709120,
    current: true,
  },
  {
    id: 'v-002',
    version: '1.1',
    description: '优化了数据质量，增加了数据验证',
    createTime: '2024-06-20 16:30:00',
    creator: 'admin',
    fileCount: 45000,
    size: 4831838208,
    current: false,
  },
  {
    id: 'v-003',
    version: '1.0',
    description: '初始版本',
    createTime: '2024-06-20 14:30:00',
    creator: 'admin',
    fileCount: 40000,
    size: 4294967296,
    current: false,
  },
]);

const dataStructure = ref<TreeNode[]>([
  {
    title: 'imagenet-subset',
    key: 'root',
    icon: FolderOutlined,
    children: [
      {
        title: 'train',
        key: 'train',
        icon: FolderOutlined,
        children: [
          { title: 'class_001', key: 'train-class_001', icon: FolderOutlined },
          { title: 'class_002', key: 'train-class_002', icon: FolderOutlined },
          { title: '...', key: 'train-more', icon: FolderOutlined },
        ],
      },
      {
        title: 'val',
        key: 'val',
        icon: FolderOutlined,
        children: [
          { title: 'class_001', key: 'val-class_001', icon: FolderOutlined },
          { title: 'class_002', key: 'val-class_002', icon: FolderOutlined },
        ],
      },
      {
        title: 'test',
        key: 'test',
        icon: FolderOutlined,
        children: [
          { title: 'images', key: 'test-images', icon: FolderOutlined },
        ],
      },
      { title: 'labels.json', key: 'labels', icon: FileOutlined },
      { title: 'metadata.json', key: 'metadata', icon: FileOutlined },
    ],
  },
]);

const sampleData = ref<FileItem[]>([
  {
    name: 'image_001.jpg',
    path: '/train/class_001/image_001.jpg',
    type: 'image',
    size: 204800,
    preview: 'https://via.placeholder.com/40x40?text=IMG',
  },
  {
    name: 'image_002.jpg',
    path: '/train/class_001/image_002.jpg',
    type: 'image',
    size: 187392,
    preview: 'https://via.placeholder.com/40x40?text=IMG',
  },
  {
    name: 'image_003.jpg',
    path: '/train/class_002/image_003.jpg',
    type: 'image',
    size: 225280,
    preview: 'https://via.placeholder.com/40x40?text=IMG',
  },
  {
    name: 'labels.json',
    path: '/labels.json',
    type: 'json',
    size: 1024,
  },
]);

// ===== 表单验证规则 =====
const createFormRules = {
  name: [
    { required: true, message: '请输入数据集名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
    {
      pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
      message: '名称只能包含小写字母、数字和连字符',
      trigger: 'blur',
    },
  ],
  dataType: [{ required: true, message: '请选择数据类型', trigger: 'change' }],
  sourceType: [
    { required: true, message: '请选择数据来源方式', trigger: 'change' },
  ],
  dataUrl: [
    { required: true, message: '请输入数据 URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的 URL', trigger: 'blur' },
  ],
  s3Bucket: [{ required: true, message: '请输入存储桶名称', trigger: 'blur' }],
  s3Path: [{ required: true, message: '请输入对象路径', trigger: 'blur' }],
  accessLevel: [
    { required: true, message: '请选择访问权限', trigger: 'change' },
  ],
};

const editFormRules = {
  name: [
    { required: true, message: '请输入数据集名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  accessLevel: [
    { required: true, message: '请选择访问权限', trigger: 'change' },
  ],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<DatasetItem> = [
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
    ellipsis: true,
  },
  {
    title: '数据类型',
    key: 'dataType',
    width: 120,
    slots: { customRender: 'dataType' },
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    slots: { customRender: 'status' },
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
    width: 100,
  },
  {
    title: '大小',
    key: 'size',
    width: 150,
    slots: { customRender: 'size' },
  },
  {
    title: '版本',
    dataIndex: 'version',
    key: 'version',
    width: 80,
    slots: { customRender: 'version' },
  },
  {
    title: '创建时间',
    key: 'createTime',
    width: 150,
    slots: { customRender: 'createTime' },
  },
  {
    title: '操作',
    key: 'action',
    width: 180,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

// ===== 分页配置 =====
const paginationConfig = {
  total: computed(() => filteredDatasets.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const filteredDatasets = computed(() => {
  let result = datasets.value;

  if (filterType.value) {
    result = result.filter((item) => item.dataType === filterType.value);
  }

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter(
      (item) =>
        item.name.toLowerCase().includes(keyword) ||
        (item.description && item.description.toLowerCase().includes(keyword)),
    );
  }

  return result;
});

// ===== 工具函数 =====
const getTypeColor = (type: string): string => {
  return TYPE_CONFIG[type as keyof typeof TYPE_CONFIG]?.color || 'default';
};

const getTypeIcon = (type: string) => {
  return TYPE_CONFIG[type as keyof typeof TYPE_CONFIG]?.icon || FileOutlined;
};

const getTypeText = (type: string): string => {
  return TYPE_CONFIG[type as keyof typeof TYPE_CONFIG]?.text || type;
};

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

const getAccessLevelText = (level: string): string => {
  const levelMap = {
    private: '私有',
    team: '团队共享',
    public: '公开',
  };
  return levelMap[level as keyof typeof levelMap] || level;
};

const getFileIcon = (type: string) => {
  const iconMap = {
    image: FileImageOutlined,
    json: FileTextOutlined,
    text: FileTextOutlined,
  };
  return iconMap[type as keyof typeof iconMap] || FileOutlined;
};

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
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

const beforeUpload = (): boolean => {
  return false; // 阻止自动上传
};

// ===== 事件处理函数 =====
const showCreateModal = (): void => {
  createModalVisible.value = true;
};

const handleCreateSubmit = async (): Promise<void> => {
  try {
    await createFormRef.value?.validate();
    createLoading.value = true;

    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 2000));

    const newDataset: DatasetItem = {
      id: `ds-${Date.now()}`,
      name: createForm.name,
      dataType: createForm.dataType as 'image' | 'text' | 'audio' | 'video',
      status: 'processing',
      creator: 'current-user',
      description: createForm.description,
      size: Math.floor(Math.random() * 10000000000), // 随机大小
      fileCount: Math.floor(Math.random() * 50000) + 1000,
      createTime: new Date().toLocaleString(),
      updateTime: new Date().toLocaleString(),
      version: '1.0',
      tags: createForm.tags,
      accessLevel: createForm.accessLevel,
    };

    datasets.value.unshift(newDataset);
    createModalVisible.value = false;
    message.success('数据集创建成功');

    // 重置表单
    createFormRef.value?.resetFields();
    Object.assign(createForm, {
      name: '',
      dataType: '',
      description: '',
      sourceType: 'upload',
      fileList: [],
      dataUrl: '',
      s3Bucket: '',
      s3Path: '',
      tags: [],
      accessLevel: 'private',
    });
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    createLoading.value = false;
  }
};

const handleCreateCancel = (): void => {
  createModalVisible.value = false;
  createFormRef.value?.resetFields();
};

const downloadDataset = (record: DatasetItem): void => {
  message.info(`开始下载数据集: ${record.name}`);
  // 模拟下载逻辑
};

const viewDetails = (record: DatasetItem): void => {
  selectedDataset.value = record;
  detailModalVisible.value = true;
  previewActiveKey.value = 'structure';
};

const handleMenuAction = (key: string, record: DatasetItem): void => {
  const actions = {
    preview: () => handlePreview(record),
    edit: () => handleEdit(record),
    copy: () => handleCopy(record),
    version: () => handleVersionManagement(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handlePreview = (record: DatasetItem): void => {
  viewDetails(record);
  previewActiveKey.value = 'sample';
};

const handleEdit = (record: DatasetItem): void => {
  editForm.name = record.name;
  editForm.description = record.description || '';
  editForm.tags = [...record.tags];
  editForm.accessLevel = record.accessLevel;
  selectedDataset.value = record;
  editModalVisible.value = true;
};

const handleEditSubmit = async (): Promise<void> => {
  try {
    await editFormRef.value?.validate();
    editLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1000));

    if (selectedDataset.value) {
      const index = datasets.value.findIndex(
        (item) => item.id === selectedDataset.value!.id,
      );
      if (index !== -1) {
        datasets.value[index] = {
          ...datasets.value[index]!,
          name: editForm.name,
          description: editForm.description,
          tags: [...editForm.tags],
          accessLevel: editForm.accessLevel,
          updateTime: new Date().toLocaleString(),
        };
      }
    }

    editModalVisible.value = false;
    message.success('数据集更新成功');
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    editLoading.value = false;
  }
};

const handleEditCancel = (): void => {
  editModalVisible.value = false;
  editFormRef.value?.resetFields();
};

const handleCopy = async (record: DatasetItem): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1500));

    const copiedDataset: DatasetItem = {
      ...record,
      id: `ds-${Date.now()}`,
      name: `${record.name}-copy`,
      status: 'processing',
      createTime: new Date().toLocaleString(),
      updateTime: new Date().toLocaleString(),
      version: '1.0',
    };

    datasets.value.unshift(copiedDataset);
    message.success('数据集复制成功');
  } catch (error) {
    message.error('复制失败');
  } finally {
    loading.value = false;
  }
};

const handleVersionManagement = (record: DatasetItem): void => {
  selectedDataset.value = record;
  versionModalVisible.value = true;
};

const createVersion = (): void => {
  message.info('创建新版本功能暂未实现');
};

const switchVersion = (version: VersionItem): void => {
  message.success(`已切换到版本 v${version.version}`);
  versions.value.forEach((v) => (v.current = false));
  version.current = true;
};

const deleteVersion = (version: VersionItem): void => {
  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除版本 v${version.version} 吗？`,
      okText: '确认',
      cancelText: '取消',
      type: 'warning',
      onOk: () => {
        const index = versions.value.findIndex((v) => v.id === version.id);
        if (index !== -1) {
          versions.value.splice(index, 1);
          message.success('版本删除成功');
        }
      },
    });
  });
};

const handleDelete = (record: DatasetItem): void => {
  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除数据集 "${record.name}" 吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      type: 'warning',
      onOk: () => {
        const index = datasets.value.findIndex((item) => item.id === record.id);
        if (index !== -1) {
          datasets.value.splice(index, 1);
          message.success('数据集删除成功');
        }
      },
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
.dataset-container {
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

/* ===== 标签样式 ===== */
.type-tag,
.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border-radius: 6px !important;
  font-weight: 500 !important;
  padding: 4px 8px !important;
}

.type-icon,
.status-icon {
  font-size: 12px;
}

.status-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.indicator-ready {
  background: #52c41a;
}

.indicator-processing {
  background: #1890ff;
}

.indicator-error {
  background: #ff4d4f;
}

/* ===== 大小信息 ===== */
.size-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.size-value {
  font-weight: 600;
  font-size: 14px;
}

.file-count {
  font-size: 12px;
  opacity: 0.8;
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
.create-form :deep(.ant-form-item-label > label) {
  font-weight: 500 !important;
}

.form-input,
.form-select,
.form-textarea {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.form-divider {
  font-weight: 500 !important;
}

.source-radio-group :deep(.ant-radio-wrapper) {
  margin-right: 16px;
}

.upload-dragger {
  border-radius: 8px !important;
}

/* ===== 详情页样式 ===== */
.detail-content {
  max-height: 700px;
  overflow-y: auto;
}

.preview-divider {
  font-weight: 500 !important;
  margin: 24px 0 16px 0 !important;
}

.preview-container {
  margin-top: 16px;
}

.preview-tabs :deep(.ant-tabs-tab) {
  font-weight: 500;
}

.structure-view,
.sample-view,
.stats-view {
  min-height: 300px;
  padding: 16px 0;
}

.structure-tree :deep(.ant-tree-node-content-wrapper) {
  padding: 4px 8px;
  border-radius: 4px;
}

.sample-list :deep(.ant-list-item) {
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.file-size {
  font-size: 12px;
  color: #8c8c8c;
}

/* ===== 版本管理样式 ===== */
.version-content {
  max-height: 600px;
  overflow-y: auto;
}

.version-header {
  margin-bottom: 24px;
  text-align: right;
}

.create-version-btn {
  border-radius: 6px !important;
}

.version-timeline :deep(.ant-timeline-item-content) {
  min-height: auto;
}

.version-item {
  padding: 12px 0;
}

.version-header-info {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.version-number {
  font-weight: 600;
  font-size: 16px;
}

.version-time {
  font-size: 12px;
  color: #8c8c8c;
}

.version-description {
  margin-bottom: 8px;
  color: #595959;
}

.version-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
  font-size: 12px;
  color: #8c8c8c;
}

.version-stat {
  display: flex;
  align-items: center;
  gap: 4px;
}

.version-creator {
  margin-left: auto;
}

.version-actions {
  display: flex;
  gap: 8px;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .dataset-container {
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

  .sci-fi-modal :deep(.ant-modal) {
    margin: 16px !important;
    max-width: calc(100vw - 32px) !important;
  }

  .version-stats {
    flex-direction: column;
    gap: 8px;
  }

  .version-creator {
    margin-left: 0;
  }

  .version-actions {
    flex-direction: column;
    gap: 4px;
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

  .size-info {
    font-size: 12px;
  }

  .action-buttons {
    flex-direction: column;
    gap: 2px;
  }

  .action-btn {
    font-size: 11px !important;
    padding: 3px 6px !important;
  }

  .version-header-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
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
