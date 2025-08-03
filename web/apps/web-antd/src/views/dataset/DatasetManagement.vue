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
              v-model:value="searchParams.type"
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
              v-model:value="searchParams.status"
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
              v-model:value="searchParams.name"
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
            <a-tag :color="getDatasetTypeColor(record.dataType)" class="type-tag">
              <component :is="getTypeIcon(record.dataType)" class="type-icon" />
              {{ getDatasetTypeText(record.dataType) }}
            </a-tag>
          </template>

          <!-- 状态列 -->
          <template #status="{ record }">
            <div class="status-wrapper">
              <a-tag :color="getDatasetStatusColor(record.status)" class="status-tag">
                <component
                  :is="getDatasetStatusIcon(record.status)"
                  class="status-icon"
                />
                {{ getDatasetStatusText(record.status) }}
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
        :model="uploadForm"
        :rules="createFormRules"
        layout="vertical"
        class="create-form"
      >
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="数据集名称" name="name">
              <a-input
                v-model:value="uploadForm.name"
                placeholder="请输入数据集名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="数据类型" name="dataType">
              <a-select
                v-model:value="uploadForm.type"
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
            v-model:value="uploadForm.description"
            placeholder="请输入数据集描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>

        <a-divider class="form-divider">数据来源</a-divider>

        <a-form-item label="数据来源方式" name="sourceType">
          <a-radio-group
            v-model:value="uploadForm.sourceType"
            class="source-radio-group"
          >
            <a-radio value="upload">本地上传</a-radio>
            <a-radio value="url">URL 导入</a-radio>
            <a-radio value="s3">S3 存储桶</a-radio>
          </a-radio-group>
        </a-form-item>

        <!-- 本地上传 -->
        <div v-if="uploadForm.sourceType === 'upload'">
          <a-form-item label="文件上传" name="files">
            <a-upload-dragger
              v-model:fileList="uploadForm.fileList"
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
        <div v-if="uploadForm.sourceType === 'url'">
          <a-form-item label="数据 URL" name="dataUrl">
            <a-input
              v-model:value="uploadForm.dataUrl"
              placeholder="请输入数据 URL"
              class="form-input"
            />
          </a-form-item>
        </div>

        <!-- S3 存储桶 -->
        <div v-if="uploadForm.sourceType === 's3'">
          <a-row :gutter="16">
            <a-col :xs="24" :sm="12">
              <a-form-item label="存储桶名称" name="s3Bucket">
                <a-input
                  v-model:value="uploadForm.s3Bucket"
                  placeholder="请输入存储桶名称"
                  class="form-input"
                />
              </a-form-item>
            </a-col>
            <a-col :xs="24" :sm="12">
              <a-form-item label="对象路径" name="s3Path">
                <a-input
                  v-model:value="uploadForm.s3Path"
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
            v-model:value="uploadForm.tags"
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
          <a-radio-group v-model:value="uploadForm.isPublic">
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
            <a-tag :color="getDatasetTypeColor(selectedDataset.type)">
              <component :is="getTypeIcon(selectedDataset.type)" />
              {{ getDatasetTypeText(selectedDataset.type) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag
              :color="getDatasetStatusColor(selectedDataset.status)"
              class="status-tag"
            >
              <component :is="getDatasetStatusIcon(selectedDataset.status)" />
              {{ getDatasetStatusText(selectedDataset.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="创建者">
            {{ selectedDataset.creatorName }}
          </a-descriptions-item>
          <a-descriptions-item label="大小">
            {{ formatSize(selectedDataset.size) }}
          </a-descriptions-item>
          <a-descriptions-item label="文件数量">
            {{ selectedDataset.sampleCount }} 个
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
  ShareAltOutlined,
  BarChartOutlined,
  HeartOutlined,
  HeartFilled,
} from '@ant-design/icons-vue';
import type {
  Dataset,
  DatasetType,
  DatasetFormat,
  DatasetStatus,
  DatasetListParams,
  DatasetUploadRequest,
} from '#/api/types';
import {
  getDatasetList,
  getDatasetStatistics,
  uploadDataset,
  deleteDataset,
  batchDeleteDatasets,
  downloadDataset,
  toggleDatasetFavorite,
  getSupportedFormats,
} from '#/api';
import { formatDateTime, formatFileSize } from '#/utils/date';

defineOptions({ name: 'DatasetManagement' });

// 响应式数据
const loading = ref<boolean>(false);
const datasetList = ref<Dataset[]>([]);
const selectedRowKeys = ref<string[]>([]);
const createModalVisible = ref<boolean>(false);
const detailDrawerVisible = ref<boolean>(false);
const versionDrawerVisible = ref<boolean>(false);
const selectedDataset = ref<Dataset | null>(null);
const supportedFormats = ref<string[]>([]);
const createLoading = ref<boolean>(false);
const previewActiveKey = ref<string>('structure');
const dataStructure = ref<any[]>([]);
const sampleData = ref<any[]>([]);

// 统计数据
const datasetStats = ref({
  totalDatasets: 0,
  publicDatasets: 0,
  privateDatasets: 0,
  totalSize: 0,
  totalSamples: 0,
  labeledDatasets: 0,
});

// 搜索和筛选参数
const searchParams = reactive<DatasetListParams>({
  name: '',
  type: undefined,
  format: undefined,
  status: undefined,
  isPublic: undefined,
  isLabeled: undefined,
  page: 1,
  pageSize: 20,
  sortBy: 'createTime',
  sortOrder: 'desc',
});

// 上传表单
interface UploadForm {
  name: string;
  description: string;
  type: DatasetType;
  format: DatasetFormat;
  workspaceId: string;
  tags: string[];
  isPublic: boolean;
  files: File[];
  labelFiles: File[];
  sourceType: 'upload' | 'url' | 's3';
  fileList: any[];
  dataUrl: string;
  s3Bucket: string;
  s3Path: string;
}

const uploadForm = reactive<UploadForm>({
  name: '',
  description: '',
  type: DatasetType.IMAGE,
  format: DatasetFormat.CUSTOM,
  workspaceId: 'workspace-001', // 默认工作空间
  tags: [],
  isPublic: false,
  files: [],
  labelFiles: [],
  sourceType: 'upload',
  fileList: [],
  dataUrl: '',
  s3Bucket: '',
  s3Path: '',
});

const createFormRef = ref<FormInstance>();

// 模拟数据
const mockDatasets: Dataset[] = [
  {
    id: 'dataset-001',
    name: 'imagenet-2012-subset',
    description: 'ImageNet 2012数据集子集，包含1000个类别的图像分类数据',
    type: DatasetType.IMAGE,
    format: DatasetFormat.IMAGENET,
    size: 5368709120, // 5GB
    sampleCount: 50000,
    filePath: '/datasets/imagenet-2012-subset',
    downloadUrl: 'https://datasets.example.com/imagenet-2012-subset.tar.gz',
    creatorId: 'user-001',
    creatorName: '张三',
    workspaceId: 'workspace-001',
    workspaceName: '默认工作空间',
    schema: {
      type: 'image_classification',
      classes: 1000,
      imageSize: [224, 224, 3],
    },
    statistics: {
      meanImageSize: 204800,
      classDistribution: 'balanced',
    },
    isLabeled: true,
    labelFormat: 'imagenet',
    labelCount: 50000,
    status: DatasetStatus.AVAILABLE,
    isPublic: true,
    downloadCount: 1258,
    useCount: 89,
    version: '1.2.0',
    parentDatasetId: null,
    versionHistory: [
      {
        version: '1.2.0',
        description: '优化数据质量，增加新类别',
        changeLog: '- 修复标注错误\n- 增加100个新类别\n- 优化图像质量',
        size: 5368709120,
        sampleCount: 50000,
        createTime: '2024-01-20 14:30:00',
        creatorId: 'user-001',
        creatorName: '张三',
      },
    ],
    tags: ['imagenet', 'classification', 'computer-vision', 'benchmark'],
    labels: {
      category: 'computer-vision',
      difficulty: 'medium',
      quality: 'high',
    },
    createTime: '2024-01-15 10:30:00',
    updateTime: '2024-01-20 14:30:00',
  },
  {
    id: 'dataset-002',
    name: 'sentiment-analysis-chinese',
    description: '中文情感分析数据集，包含电商评论的情感标注',
    type: DatasetType.TEXT,
    format: DatasetFormat.JSON,
    size: 1073741824, // 1GB
    sampleCount: 100000,
    filePath: '/datasets/sentiment-analysis-chinese',
    creatorId: 'user-002',
    creatorName: '李四',
    workspaceId: 'workspace-002',
    workspaceName: 'NLP研究室',
    schema: {
      type: 'text_classification',
      fields: ['text', 'sentiment', 'score'],
      sentiments: ['positive', 'negative', 'neutral'],
    },
    isLabeled: true,
    labelFormat: 'json',
    labelCount: 100000,
    status: DatasetStatus.AVAILABLE,
    isPublic: false,
    downloadCount: 456,
    useCount: 23,
    version: '2.1.0',
    parentDatasetId: null,
    versionHistory: [
      {
        version: '2.1.0',
        description: '增加更多评论数据',
        changeLog: '- 新增20000条评论\n- 优化标注质量\n- 修复数据格式问题',
        size: 1073741824,
        sampleCount: 100000,
        createTime: '2024-01-18 16:45:00',
        creatorId: 'user-002',
        creatorName: '李四',
      },
    ],
    tags: ['sentiment', 'chinese', 'nlp', 'ecommerce'],
    createTime: '2024-01-12 09:15:00',
    updateTime: '2024-01-18 16:45:00',
  },
  {
    id: 'dataset-003',
    name: 'speech-recognition-mandarin',
    description: '普通话语音识别数据集，多说话人录音',
    type: DatasetType.AUDIO,
    format: DatasetFormat.CUSTOM,
    size: 10737418240, // 10GB
    sampleCount: 25000,
    filePath: '/datasets/speech-recognition-mandarin',
    creatorId: 'user-003',
    creatorName: '王五',
    workspaceId: 'workspace-003',
    workspaceName: '语音实验室',
    schema: {
      type: 'speech_recognition',
      sampleRate: 16000,
      channels: 1,
      format: 'wav',
    },
    isLabeled: true,
    labelFormat: 'txt',
    labelCount: 25000,
    status: DatasetStatus.PROCESSING,
    isPublic: true,
    downloadCount: 234,
    useCount: 12,
    version: '1.0.0',
    parentDatasetId: null,
    versionHistory: [
      {
        version: '1.0.0',
        description: '初始版本',
        changeLog: '- 收集25000条语音数据\n- 完成文本标注\n- 数据质量检查',
        size: 10737418240,
        sampleCount: 25000,
        createTime: '2024-01-20 11:00:00',
        creatorId: 'user-003',
        creatorName: '王五',
      },
    ],
    tags: ['speech', 'mandarin', 'asr', 'multi-speaker'],
    createTime: '2024-01-20 11:00:00',
    updateTime: '2024-01-20 15:30:00',
  },
];

// 表格列定义
const columns: TableColumnsType<Dataset> = [
  {
    title: '数据集信息',
    key: 'datasetInfo',
    slots: { customRender: 'datasetInfo' },
    width: 300,
  },
  {
    title: '类型',
    key: 'type',
    slots: { customRender: 'type' },
    width: 120,
  },
  {
    title: '格式',
    key: 'format',
    slots: { customRender: 'format' },
    width: 100,
  },
  {
    title: '状态',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100,
  },
  {
    title: '样本数',
    key: 'sampleCount',
    slots: { customRender: 'sampleCount' },
    width: 100,
  },
  {
    title: '创建者',
    key: 'creator',
    slots: { customRender: 'creator' },
    width: 120,
  },
  {
    title: '创建时间',
    key: 'createTime',
    slots: { customRender: 'createTime' },
    width: 150,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 200,
    fixed: 'right' as const,
  },
];

// 计算属性
const filteredDatasets = computed(() => {
  let filtered = datasetList.value;
  
  // 名称搜索
  if (searchParams.name) {
    filtered = filtered.filter(dataset => 
      dataset.name.toLowerCase().includes(searchParams.name!.toLowerCase()) ||
      dataset.description?.toLowerCase().includes(searchParams.name!.toLowerCase())
    );
  }
  
  // 类型筛选
  if (searchParams.type) {
    filtered = filtered.filter(dataset => dataset.type === searchParams.type);
  }
  
  // 格式筛选
  if (searchParams.format) {
    filtered = filtered.filter(dataset => dataset.format === searchParams.format);
  }
  
  // 状态筛选
  if (searchParams.status) {
    filtered = filtered.filter(dataset => dataset.status === searchParams.status);
  }
  
  // 公开性筛选
  if (searchParams.isPublic !== undefined) {
    filtered = filtered.filter(dataset => dataset.isPublic === searchParams.isPublic);
  }
  
  // 标注状态筛选
  if (searchParams.isLabeled !== undefined) {
    filtered = filtered.filter(dataset => dataset.isLabeled === searchParams.isLabeled);
  }
  
  return filtered;
});

// 工具方法
const getDatasetTypeText = (type: DatasetType) => {
  const types = {
    [DatasetType.IMAGE]: '图像',
    [DatasetType.TEXT]: '文本',
    [DatasetType.AUDIO]: '音频',
    [DatasetType.VIDEO]: '视频',
    [DatasetType.TABULAR]: '表格',
    [DatasetType.TIME_SERIES]: '时间序列',
    [DatasetType.GRAPH]: '图网络',
    [DatasetType.CUSTOM]: '自定义',
  };
  return types[type] || type;
};

const getDatasetTypeColor = (type: DatasetType) => {
  const colors = {
    [DatasetType.IMAGE]: 'blue',
    [DatasetType.TEXT]: 'green',
    [DatasetType.AUDIO]: 'orange',
    [DatasetType.VIDEO]: 'purple',
    [DatasetType.TABULAR]: 'cyan',
    [DatasetType.TIME_SERIES]: 'magenta',
    [DatasetType.GRAPH]: 'volcano',
    [DatasetType.CUSTOM]: 'default',
  };
  return colors[type] || 'default';
};

const getDatasetFormatText = (format: DatasetFormat) => {
  const formats = {
    [DatasetFormat.CSV]: 'CSV',
    [DatasetFormat.JSON]: 'JSON',
    [DatasetFormat.PARQUET]: 'Parquet',
    [DatasetFormat.COCO]: 'COCO',
    [DatasetFormat.YOLO]: 'YOLO',
    [DatasetFormat.PASCAL_VOC]: 'Pascal VOC',
    [DatasetFormat.IMAGENET]: 'ImageNet',
    [DatasetFormat.CUSTOM]: '自定义',
  };
  return formats[format] || format;
};

const getDatasetStatusText = (status: DatasetStatus) => {
  const statuses = {
    [DatasetStatus.UPLOADING]: '上传中',
    [DatasetStatus.PROCESSING]: '处理中',
    [DatasetStatus.AVAILABLE]: '可用',
    [DatasetStatus.ERROR]: '错误',
    [DatasetStatus.DELETED]: '已删除',
  };
  return statuses[status] || status;
};

const getDatasetStatusColor = (status: DatasetStatus) => {
  const colors = {
    [DatasetStatus.UPLOADING]: 'processing',
    [DatasetStatus.PROCESSING]: 'processing',
    [DatasetStatus.AVAILABLE]: 'success',
    [DatasetStatus.ERROR]: 'error',
    [DatasetStatus.DELETED]: 'default',
  };
  return colors[status] || 'default';
};

const getTypeIcon = (type: DatasetType) => {
  const icons = {
    [DatasetType.IMAGE]: FileImageOutlined,
    [DatasetType.TEXT]: FileTextOutlined,
    [DatasetType.AUDIO]: SoundOutlined,
    [DatasetType.VIDEO]: VideoCameraOutlined,
    [DatasetType.TABULAR]: FileTextOutlined,
    [DatasetType.TIME_SERIES]: BarChartOutlined,
    [DatasetType.GRAPH]: ShareAltOutlined,
    [DatasetType.CUSTOM]: FileOutlined,
  };
  return icons[type] || FileOutlined;
};

const getDatasetStatusIcon = (status: DatasetStatus) => {
  const icons = {
    [DatasetStatus.UPLOADING]: ReloadOutlined,
    [DatasetStatus.PROCESSING]: ReloadOutlined,
    [DatasetStatus.AVAILABLE]: CheckCircleOutlined,
    [DatasetStatus.ERROR]: CloseCircleOutlined,
    [DatasetStatus.DELETED]: DeleteOutlined,
  };
  return icons[status] || CheckCircleOutlined;
};

const getAccessLevelText = (level: string) => {
  const levels = {
    'private': '私有',
    'team': '团队',
    'public': '公开'
  };
  return levels[level as keyof typeof levels] || level;
};

const getFileIcon = (fileName: string) => {
  const ext = fileName.split('.').pop()?.toLowerCase();
  if (['jpg', 'jpeg', 'png', 'gif'].includes(ext || '')) return FileImageOutlined;
  if (['txt', 'csv', 'json'].includes(ext || '')) return FileTextOutlined;
  return FileOutlined;
};

// 数据加载
const loadDatasets = async () => {
  try {
    loading.value = true;
    // const response = await getDatasetList(searchParams);
    // datasetList.value = response.data.items;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    datasetList.value = mockDatasets;
    
    updateStats();
  } catch (error) {
    message.error('加载数据集列表失败');
  } finally {
    loading.value = false;
  }
};

const loadStatistics = async () => {
  try {
    // const response = await getDatasetStatistics();
    // datasetStats.value = response.data;
    
    // 模拟统计数据
    datasetStats.value = {
      totalDatasets: mockDatasets.length,
      publicDatasets: mockDatasets.filter(d => d.isPublic).length,
      privateDatasets: mockDatasets.filter(d => !d.isPublic).length,
      totalSize: mockDatasets.reduce((sum, d) => sum + d.size, 0),
      totalSamples: mockDatasets.reduce((sum, d) => sum + d.sampleCount, 0),
      labeledDatasets: mockDatasets.filter(d => d.isLabeled).length,
    };
  } catch (error) {
    message.error('加载统计信息失败');
  }
};

const loadFormats = async () => {
  try {
    // const response = await getSupportedFormats();
    // supportedFormats.value = response.data;
    
    // 模拟格式数据
    supportedFormats.value = Object.values(DatasetFormat);
  } catch (error) {
    message.error('加载格式列表失败');
  }
};

const updateStats = () => {
  const stats = {
    totalDatasets: datasetList.value.length,
    publicDatasets: datasetList.value.filter(d => d.isPublic).length,
    privateDatasets: datasetList.value.filter(d => !d.isPublic).length,
    totalSize: datasetList.value.reduce((sum, d) => sum + d.size, 0),
    totalSamples: datasetList.value.reduce((sum, d) => sum + d.sampleCount, 0),
    labeledDatasets: datasetList.value.filter(d => d.isLabeled).length,
  };
  datasetStats.value = stats;
};

const refreshData = async () => {
  await Promise.all([
    loadDatasets(),
    loadStatistics(),
  ]);
};

const viewDetails = (record: Dataset) => {
  selectedDataset.value = record;
  detailDrawerVisible.value = true;
};

// 事件处理
const handleFilterChange = () => {
  // 筛选变化时可以重新加载数据
};

const handleSearch = () => {
  loadDatasets();
};

const handleSearchChange = () => {
  // 搜索输入变化时的处理逻辑
};

const resetFilters = () => {
  Object.assign(searchParams, {
    name: '',
    type: undefined,
    format: undefined,
    status: undefined,
    isPublic: undefined,
    isLabeled: undefined,
  });
  loadDatasets();
};

const showCreateModal = () => {
  createModalVisible.value = true;
  resetUploadForm();
};

const resetUploadForm = () => {
  Object.assign(uploadForm, {
    name: '',
    description: '',
    type: DatasetType.IMAGE,
    format: DatasetFormat.CUSTOM,
    workspaceId: 'workspace-001',
    tags: [],
    isPublic: false,
    files: [],
    labelFiles: [],
  });
};

const handleCreateSubmit = async () => {
  try {
    await createFormRef.value?.validate();
    
    if (uploadForm.files.length === 0) {
      message.error('请选择数据文件');
      return;
    }
    
    createLoading.value = true;
    
    const request: DatasetUploadRequest = {
      name: uploadForm.name,
      description: uploadForm.description,
      type: uploadForm.type,
      format: uploadForm.format,
      workspaceId: uploadForm.workspaceId,
      tags: uploadForm.tags,
      isPublic: uploadForm.isPublic,
      files: uploadForm.files,
      labelFiles: uploadForm.labelFiles,
    };
    
    // const response = await uploadDataset(request);
    
    // 模拟上传成功
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    message.success('数据集上传成功');
    createModalVisible.value = false;
    loadDatasets();
  } catch (error) {
    message.error('上传失败');
  } finally {
    createLoading.value = false;
  }
};

const handleCreateCancel = () => {
  createModalVisible.value = false;
};

const beforeUpload = (file: File, fileList: File[]) => {
  uploadForm.files = [...uploadForm.files, file];
  if (!uploadForm.name) {
    uploadForm.name = file.name.replace(/\.[^/.]+$/, '');
  }
  return false; // 阻止自动上传
};

const beforeLabelUpload = (file: File, fileList: File[]) => {
  uploadForm.labelFiles = [...uploadForm.labelFiles, file];
  return false; // 阻止自动上传
};

// 数据集操作
const viewDatasetDetail = (dataset: Dataset) => {
  selectedDataset.value = dataset;
  detailDrawerVisible.value = true;
};

const viewDatasetVersions = (dataset: Dataset) => {
  selectedDataset.value = dataset;
  versionDrawerVisible.value = true;
};

const downloadDatasetFile = async (dataset: Dataset) => {
  try {
    // const blob = await downloadDataset(dataset.id);
    // const url = window.URL.createObjectURL(blob);
    // const link = document.createElement('a');
    // link.href = url;
    // link.download = `${dataset.name}-${dataset.version}.tar.gz`;
    // link.click();
    // window.URL.revokeObjectURL(url);
    
    // 模拟下载
    message.success('数据集下载中...');
  } catch (error) {
    message.error('下载失败');
  }
};

const toggleFavorite = async (dataset: Dataset) => {
  try {
    // await toggleDatasetFavorite(dataset.id, !dataset.isFavorite);
    
    // 模拟切换收藏状态
    message.success(dataset.isFavorite ? '已取消收藏' : '已收藏');
    loadDatasets();
  } catch (error) {
    message.error('操作失败');
  }
};

const editDataset = (dataset: Dataset) => {
  message.info('编辑功能开发中');
};

const deleteDatasetItem = async (dataset: Dataset) => {
  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除数据集 "${dataset.name}" 吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        try {
          // await deleteDataset(dataset.id);
          
          // 模拟删除
          message.success('数据集删除成功');
          loadDatasets();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  });
};

const batchDelete = async () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的数据集');
    return;
  }
  
  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '批量删除',
      content: `确定要删除选中的 ${selectedRowKeys.value.length} 个数据集吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        try {
          // await batchDeleteDatasets(selectedRowKeys.value);
          
          // 模拟批量删除
          message.success('批量删除成功');
          selectedRowKeys.value = [];
          loadDatasets();
        } catch (error) {
          message.error('批量删除失败');
        }
      },
    });
  });
};

// 菜单操作处理
const handleMenuAction = (key: string, record: Dataset) => {
  const actions = {
    preview: () => viewDatasetDetail(record),
    edit: () => editDataset(record),
    copy: () => message.info('复制功能开发中'),
    version: () => viewDatasetVersions(record),
    delete: () => deleteDatasetItem(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

// 表单验证规则
const createFormRules = {
  name: [
    { required: true, message: '请输入数据集名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
  ],
  type: [
    { required: true, message: '请选择数据类型', trigger: 'change' },
  ],
  format: [
    { required: true, message: '请选择数据格式', trigger: 'change' },
  ],
  workspaceId: [
    { required: true, message: '请选择工作空间', trigger: 'change' },
  ],
};

// 格式化方法
const formatSize = (bytes: number): string => {
  return formatFileSize(bytes);
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

// 分页配置
const paginationConfig = {
  total: computed(() => filteredDatasets.value.length),
  pageSize: 20,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// 初始化
onMounted(() => {
  refreshData();
  loadFormats();
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
