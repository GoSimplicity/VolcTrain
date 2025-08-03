<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue';
import {
  Card,
  Row,
  Col,
  Button,
  Space,
  Statistic,
  Table,
  Tag,
  Badge,
  Input,
  Select,
  Upload,
  Modal,
  Form,
  Progress,
  Tooltip,
  Dropdown,
  Menu,
  Avatar,
  Image,
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  PlusOutlined,
  UploadOutlined,
  SearchOutlined,
  FilterOutlined,
  DownloadOutlined,
  EyeOutlined,
  EditOutlined,
  DeleteOutlined,
  ShareAltOutlined,
  HeartOutlined,
  HeartFilled,
  MoreOutlined,
  FileTextOutlined,
  RocketOutlined,
  SafetyCertificateOutlined,
  ApiOutlined,
  BarChartOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  ExpandOutlined,
  DiffOutlined,
  CheckOutlined,
  CopyOutlined,
} from '@ant-design/icons-vue';
import type { 
  Model,
  ModelType,
  ModelStatus,
  ModelListParams,
  ModelUploadRequest
} from '#/api/types';
import { 
  getModelList,
  getModelStatistics,
  deleteModel,
  batchDeleteModels,
  uploadModel,
  downloadModel,
  toggleModelFavorite,
  getSupportedFrameworks
} from '#/api';
import { formatDateTime, formatFileSize } from '#/utils/date';
import ModelDetailDrawer from './components/ModelDetailDrawer.vue';
import ModelVersionDrawer from './components/ModelVersionDrawer.vue';
import ModelDeployDrawer from './components/ModelDeployDrawer.vue';

defineOptions({ name: 'ModelList' });

// 响应式数据
const loading = ref(false);
const modelList = ref<Model[]>([]);
const selectedRowKeys = ref<string[]>([]);
const uploadModalVisible = ref(false);
const detailDrawerVisible = ref(false);
const versionDrawerVisible = ref(false);
const deployDrawerVisible = ref(false);
const selectedModel = ref<Model | null>(null);
const supportedFrameworks = ref<string[]>([]);

// 统计数据
const modelStats = ref({
  totalModels: 0,
  publicModels: 0,
  privateModels: 0,
  storageUsed: 0,
});

// 搜索和筛选参数
const searchParams = reactive<ModelListParams>({
  name: '',
  type: undefined,
  status: undefined,
  framework: '',
  isPublic: undefined,
  page: 1,
  pageSize: 20,
  sortBy: 'createTime',
  sortOrder: 'desc',
});

// 上传表单
interface UploadForm {
  name: string;
  description: string;
  type: ModelType;
  framework: string;
  workspaceId: string;
  tags: string[];
  isPublic: boolean;
  file: File | null;
}

const uploadForm = reactive<UploadForm>({
  name: '',
  description: '',
  type: ModelType.CLASSIFICATION,
  framework: '',
  workspaceId: 'workspace-001', // 默认工作空间
  tags: [],
  isPublic: false,
  file: null,
});

const uploadFormRef = ref();

// 模拟数据
const mockModels: Model[] = [
  {
    id: 'model-001',
    name: 'BERT-Base-Chinese',
    description: '中文BERT预训练模型，适用于各种NLP任务',
    version: 'v1.2.0',
    type: ModelType.NLP,
    framework: 'PyTorch',
    size: 412000000,
    accuracy: 0.95,
    filePath: '/models/bert-base-chinese/v1.2.0',
    downloadUrl: 'https://models.example.com/bert-base-chinese-v1.2.0.tar.gz',
    creatorId: 'user-001',
    creatorName: '张三',
    workspaceId: 'workspace-001',
    workspaceName: '默认工作空间',
    projectId: 'project-001',
    projectName: 'NLP研究项目',
    trainedJobId: 'job-001',
    trainingDataset: 'chinese-corpus-v1',
    hyperParameters: {
      'learning_rate': 0.00005,
      'batch_size': 32,
      'epochs': 10,
    },
    status: ModelStatus.AVAILABLE,
    isPublic: true,
    downloadCount: 1258,
    useCount: 89,
    tags: ['bert', 'chinese', 'nlp', 'pretrained'],
    labels: {
      'category': 'nlp',
      'language': 'chinese',
    },
    createTime: '2024-01-15 10:30:00',
    updateTime: '2024-01-20 14:20:00',
  },
  {
    id: 'model-002',
    name: 'ResNet-50-ImageNet',
    description: 'ResNet-50图像分类模型，在ImageNet数据集上预训练',
    version: 'v2.1.0',
    type: ModelType.CLASSIFICATION,
    framework: 'TensorFlow',
    size: 98000000,
    accuracy: 0.92,
    filePath: '/models/resnet50-imagenet/v2.1.0',
    downloadUrl: 'https://models.example.com/resnet50-imagenet-v2.1.0.h5',
    creatorId: 'user-002',
    creatorName: '李四',
    workspaceId: 'workspace-001',
    workspaceName: '默认工作空间',
    trainedJobId: 'job-002',
    trainingDataset: 'imagenet-2012',
    status: ModelStatus.AVAILABLE,
    isPublic: false,
    downloadCount: 456,
    useCount: 23,
    tags: ['resnet', 'classification', 'imagenet', 'computer-vision'],
    createTime: '2024-01-18 09:15:00',
    updateTime: '2024-01-19 16:45:00',
  },
  {
    id: 'model-003',
    name: 'YOLOv8-Object-Detection',
    description: 'YOLOv8目标检测模型，支持实时检测',
    version: 'v1.0.0',
    type: ModelType.OBJECT_DETECTION,
    framework: 'YOLOv8',
    size: 47000000,
    accuracy: 0.88,
    filePath: '/models/yolov8-detection/v1.0.0',
    creatorId: 'user-003',
    creatorName: '王五',
    workspaceId: 'workspace-002',
    workspaceName: 'CV实验室',
    status: ModelStatus.TRAINING,
    isPublic: true,
    downloadCount: 234,
    useCount: 12,
    tags: ['yolo', 'detection', 'real-time', 'computer-vision'],
    createTime: '2024-01-20 11:00:00',
    updateTime: '2024-01-20 15:30:00',
  },
];

// 表格列定义
const columns = [
  {
    title: '模型信息',
    key: 'modelInfo',
    slots: { customRender: 'modelInfo' },
    width: 300,
  },
  {
    title: '类型',
    key: 'type',
    slots: { customRender: 'type' },
    width: 120,
  },
  {
    title: '框架',
    dataIndex: 'framework',
    key: 'framework',
    width: 100,
  },
  {
    title: '状态',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100,
  },
  {
    title: '准确率',
    key: 'accuracy',
    slots: { customRender: 'accuracy' },
    width: 100,
  },
  {
    title: '下载量',
    key: 'downloads',
    slots: { customRender: 'downloads' },
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
const filteredModels = computed(() => {
  let filtered = modelList.value;
  
  // 名称搜索
  if (searchParams.name) {
    filtered = filtered.filter(model => 
      model.name.toLowerCase().includes(searchParams.name!.toLowerCase()) ||
      model.description?.toLowerCase().includes(searchParams.name!.toLowerCase())
    );
  }
  
  // 类型筛选
  if (searchParams.type) {
    filtered = filtered.filter(model => model.type === searchParams.type);
  }
  
  // 状态筛选
  if (searchParams.status) {
    filtered = filtered.filter(model => model.status === searchParams.status);
  }
  
  // 框架筛选
  if (searchParams.framework) {
    filtered = filtered.filter(model => model.framework === searchParams.framework);
  }
  
  // 公开性筛选
  if (searchParams.isPublic !== undefined) {
    filtered = filtered.filter(model => model.isPublic === searchParams.isPublic);
  }
  
  return filtered;
});

// 工具方法
const getModelTypeText = (type: ModelType) => {
  const types = {
    [ModelType.CLASSIFICATION]: '分类',
    [ModelType.REGRESSION]: '回归',
    [ModelType.OBJECT_DETECTION]: '目标检测',
    [ModelType.SEMANTIC_SEGMENTATION]: '语义分割',
    [ModelType.NLP]: '自然语言处理',
    [ModelType.RECOMMENDATION]: '推荐系统',
    [ModelType.GENERATIVE]: '生成模型',
    [ModelType.CUSTOM]: '自定义',
  };
  return types[type] || type;
};

const getModelTypeColor = (type: ModelType) => {
  const colors = {
    [ModelType.CLASSIFICATION]: 'blue',
    [ModelType.REGRESSION]: 'green',
    [ModelType.OBJECT_DETECTION]: 'orange',
    [ModelType.SEMANTIC_SEGMENTATION]: 'purple',
    [ModelType.NLP]: 'cyan',
    [ModelType.RECOMMENDATION]: 'magenta',
    [ModelType.GENERATIVE]: 'red',
    [ModelType.CUSTOM]: 'default',
  };
  return colors[type] || 'default';
};

const getModelStatusText = (status: ModelStatus) => {
  const statuses = {
    [ModelStatus.TRAINING]: '训练中',
    [ModelStatus.AVAILABLE]: '可用',
    [ModelStatus.DEPRECATED]: '已弃用',
    [ModelStatus.DELETED]: '已删除',
  };
  return statuses[status] || status;
};

const getModelStatusColor = (status: ModelStatus) => {
  const colors = {
    [ModelStatus.TRAINING]: 'processing',
    [ModelStatus.AVAILABLE]: 'success',
    [ModelStatus.DEPRECATED]: 'warning',
    [ModelStatus.DELETED]: 'error',
  };
  return colors[status] || 'default';
};

// 数据加载
const loadModels = async () => {
  try {
    loading.value = true;
    // const response = await getModelList(searchParams);
    // modelList.value = response.data.items;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    modelList.value = mockModels;
    
    updateStats();
  } catch (error) {
    message.error('加载模型列表失败');
  } finally {
    loading.value = false;
  }
};

const loadStatistics = async () => {
  try {
    // const response = await getModelStatistics();
    // modelStats.value = response.data;
    
    // 模拟统计数据
    modelStats.value = {
      totalModels: mockModels.length,
      publicModels: mockModels.filter(m => m.isPublic).length,
      privateModels: mockModels.filter(m => !m.isPublic).length,
      storageUsed: mockModels.reduce((sum, m) => sum + m.size, 0),
    };
  } catch (error) {
    message.error('加载统计信息失败');
  }
};

const loadFrameworks = async () => {
  try {
    // const response = await getSupportedFrameworks();
    // supportedFrameworks.value = response.data;
    
    // 模拟框架数据
    supportedFrameworks.value = ['PyTorch', 'TensorFlow', 'YOLOv8', 'Scikit-learn', 'Hugging Face'];
  } catch (error) {
    message.error('加载框架列表失败');
  }
};

const updateStats = () => {
  const stats = {
    totalModels: modelList.value.length,
    publicModels: modelList.value.filter(m => m.isPublic).length,
    privateModels: modelList.value.filter(m => !m.isPublic).length,
    storageUsed: modelList.value.reduce((sum, m) => sum + m.size, 0),
  };
  modelStats.value = stats;
};

const refreshData = () => {
  loadModels();
  loadStatistics();
};

// 事件处理
const handleSearch = () => {
  loadModels();
};

const resetFilters = () => {
  Object.assign(searchParams, {
    name: '',
    type: undefined,
    status: undefined,
    framework: '',
    isPublic: undefined,
  });
  loadModels();
};

const showUploadModal = () => {
  uploadModalVisible.value = true;
  resetUploadForm();
};

const resetUploadForm = () => {
  Object.assign(uploadForm, {
    name: '',
    description: '',
    type: ModelType.CLASSIFICATION,
    framework: '',
    workspaceId: 'workspace-001',
    tags: [],
    isPublic: false,
    file: null,
  });
};

const handleUploadSubmit = async () => {
  try {
    await uploadFormRef.value?.validate();
    
    if (!uploadForm.file) {
      message.error('请选择模型文件');
      return;
    }
    
    const request: ModelUploadRequest = {
      name: uploadForm.name,
      description: uploadForm.description,
      type: uploadForm.type,
      framework: uploadForm.framework,
      workspaceId: uploadForm.workspaceId,
      tags: uploadForm.tags,
      isPublic: uploadForm.isPublic,
      file: uploadForm.file,
    };
    
    // const response = await uploadModel(request);
    
    // 模拟上传成功
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    message.success('模型上传成功');
    uploadModalVisible.value = false;
    loadModels();
  } catch (error) {
    message.error('上传失败');
  }
};

const handleUploadCancel = () => {
  uploadModalVisible.value = false;
};

const beforeUpload = (file: File) => {
  uploadForm.file = file;
  if (!uploadForm.name) {
    uploadForm.name = file.name.replace(/\.[^/.]+$/, '');
  }
  return false; // 阻止自动上传
};

// 模型操作
const viewModelDetail = (model: Model) => {
  selectedModel.value = model;
  detailDrawerVisible.value = true;
};

const viewModelVersions = (model: Model) => {
  selectedModel.value = model;
  versionDrawerVisible.value = true;
};

const deployModel = (model: Model) => {
  selectedModel.value = model;
  deployDrawerVisible.value = true;
};

const downloadModelFile = async (model: Model) => {
  try {
    // const blob = await downloadModel(model.id);
    // const url = window.URL.createObjectURL(blob);
    // const link = document.createElement('a');
    // link.href = url;
    // link.download = `${model.name}-${model.version}.tar.gz`;
    // link.click();
    // window.URL.revokeObjectURL(url);
    
    // 模拟下载
    message.success('模型下载中...');
  } catch (error) {
    message.error('下载失败');
  }
};

const toggleFavorite = async (model: Model) => {
  try {
    // await toggleModelFavorite(model.id, !model.isFavorite);
    
    // 模拟切换收藏状态
    message.success(model.isFavorite ? '已取消收藏' : '已收藏');
    loadModels();
  } catch (error) {
    message.error('操作失败');
  }
};

const editModel = (model: Model) => {
  message.info('编辑功能开发中');
};

const deleteModelItem = async (model: Model) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除模型 "${model.name}" 吗？此操作不可恢复。`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await deleteModel(model.id);
        
        // 模拟删除
        message.success('模型删除成功');
        loadModels();
      } catch (error) {
        message.error('删除失败');
      }
    },
  });
};

const batchDelete = async () => {
  if (selectedRowKeys.value.length === 0) {
    message.warning('请选择要删除的模型');
    return;
  }
  
  Modal.confirm({
    title: '批量删除',
    content: `确定要删除选中的 ${selectedRowKeys.value.length} 个模型吗？此操作不可恢复。`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await batchDeleteModels(selectedRowKeys.value);
        
        // 模拟批量删除
        message.success('批量删除成功');
        selectedRowKeys.value = [];
        loadModels();
      } catch (error) {
        message.error('批量删除失败');
      }
    },
  });
};

// 表单验证规则
const uploadFormRules = {
  name: [
    { required: true, message: '请输入模型名称', trigger: 'blur' },
  ],
  type: [
    { required: true, message: '请选择模型类型', trigger: 'change' },
  ],
  framework: [
    { required: true, message: '请选择框架', trigger: 'change' },
  ],
  workspaceId: [
    { required: true, message: '请选择工作空间', trigger: 'change' },
  ],
};

// 初始化
onMounted(() => {
  loadModels();
  loadStatistics();
  loadFrameworks();
});
</script>

<template>
  <div class="model-list-container">
    <!-- 页面头部 -->
    <Card>
      <div class="page-header">
        <div class="header-left">
          <h2>模型管理</h2>
          <p>管理和部署机器学习模型</p>
        </div>
        <div class="header-right">
          <Space>
            <Button @click="refreshData">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary" @click="showUploadModal">
              <UploadOutlined />
              上传模型
            </Button>
          </Space>
        </div>
      </div>
    </Card>

    <!-- 统计卡片 -->
    <Row :gutter="16" style="margin: 16px 0">
      <Col :span="6">
        <Card>
          <Statistic
            title="总模型数"
            :value="modelStats.totalModels"
            :value-style="{ color: '#3f8600' }"
            prefix="🔮"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="公开模型"
            :value="modelStats.publicModels"
            :value-style="{ color: '#1890ff' }"
            prefix="🌐"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="私有模型"
            :value="modelStats.privateModels"
            :value-style="{ color: '#722ed1' }"
            prefix="🔒"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="存储使用"
            :value="formatFileSize(modelStats.storageUsed)"
            :value-style="{ color: '#faad14' }"
            prefix="💾"
          />
        </Card>
      </Col>
    </Row>

    <!-- 搜索和筛选 -->
    <Card style="margin-bottom: 16px">
      <div class="search-filters">
        <Row :gutter="16">
          <Col :span="6">
            <Input
              v-model:value="searchParams.name"
              placeholder="搜索模型名称或描述"
              allow-clear
              @press-enter="handleSearch"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </Input>
          </Col>
          <Col :span="4">
            <Select
              v-model:value="searchParams.type"
              placeholder="模型类型"
              allow-clear
            >
              <Select.Option :value="ModelType.CLASSIFICATION">分类</Select.Option>
              <Select.Option :value="ModelType.REGRESSION">回归</Select.Option>
              <Select.Option :value="ModelType.OBJECT_DETECTION">目标检测</Select.Option>
              <Select.Option :value="ModelType.SEMANTIC_SEGMENTATION">语义分割</Select.Option>
              <Select.Option :value="ModelType.NLP">自然语言处理</Select.Option>
              <Select.Option :value="ModelType.RECOMMENDATION">推荐系统</Select.Option>
              <Select.Option :value="ModelType.GENERATIVE">生成模型</Select.Option>
              <Select.Option :value="ModelType.CUSTOM">自定义</Select.Option>
            </Select>
          </Col>
          <Col :span="4">
            <Select
              v-model:value="searchParams.framework"
              placeholder="框架"
              allow-clear
            >
              <Select.Option 
                v-for="framework in supportedFrameworks" 
                :key="framework" 
                :value="framework"
              >
                {{ framework }}
              </Select.Option>
            </Select>
          </Col>
          <Col :span="4">
            <Select
              v-model:value="searchParams.status"
              placeholder="状态"
              allow-clear
            >
              <Select.Option :value="ModelStatus.TRAINING">训练中</Select.Option>
              <Select.Option :value="ModelStatus.AVAILABLE">可用</Select.Option>
              <Select.Option :value="ModelStatus.DEPRECATED">已弃用</Select.Option>
            </Select>
          </Col>
          <Col :span="3">
            <Select
              v-model:value="searchParams.isPublic"
              placeholder="可见性"
              allow-clear
            >
              <Select.Option :value="true">公开</Select.Option>
              <Select.Option :value="false">私有</Select.Option>
            </Select>
          </Col>
          <Col :span="3">
            <Space>
              <Button type="primary" @click="handleSearch">
                <SearchOutlined />
                搜索
              </Button>
              <Button @click="resetFilters">
                重置
              </Button>
            </Space>
          </Col>
        </Row>
      </div>
    </Card>

    <!-- 操作栏 -->
    <Card style="margin-bottom: 16px">
      <div class="toolbar">
        <div class="toolbar-left">
          <Space>
            <Button 
              danger 
              :disabled="selectedRowKeys.length === 0"
              @click="batchDelete"
            >
              <DeleteOutlined />
              批量删除 ({{ selectedRowKeys.length }})
            </Button>
          </Space>
        </div>
        <div class="toolbar-right">
          <Space>
            <span>共 {{ filteredModels.length }} 个模型</span>
          </Space>
        </div>
      </div>
    </Card>

    <!-- 模型列表 -->
    <Card>
      <Table
        :columns="columns"
        :data-source="filteredModels"
        :loading="loading"
        row-key="id"
        :pagination="{
          current: searchParams.page,
          pageSize: searchParams.pageSize,
          total: filteredModels.length,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`,
        }"
        :row-selection="{
          selectedRowKeys,
          onChange: (keys: string[]) => { selectedRowKeys = keys; },
        }"
        :scroll="{ x: 1500 }"
      >
        <!-- 模型信息 -->
        <template #modelInfo="{ record }">
          <div class="model-info">
            <div class="model-header">
              <div class="model-name">
                <span class="name-text">{{ record.name }}</span>
                <Tag v-if="record.isPublic" color="blue" size="small">公开</Tag>
                <HeartFilled 
                  v-if="record.isFavorite" 
                  class="favorite-icon" 
                  style="color: #f5222d"
                />
              </div>
              <div class="model-version">{{ record.version }}</div>
            </div>
            <div class="model-desc">{{ record.description || '无描述' }}</div>
            <div class="model-tags">
              <Tag 
                v-for="tag in record.tags.slice(0, 3)" 
                :key="tag" 
                size="small"
              >
                {{ tag }}
              </Tag>
              <Tag 
                v-if="record.tags.length > 3" 
                size="small" 
                color="default"
              >
                +{{ record.tags.length - 3 }}
              </Tag>
            </div>
            <div class="model-metrics">
              <span class="metric-item">
                📊 准确率: {{ (record.accuracy * 100).toFixed(1) }}%
              </span>
              <span class="metric-item">
                📦 大小: {{ formatFileSize(record.size) }}
              </span>
            </div>
          </div>
        </template>

        <!-- 模型类型 -->
        <template #type="{ record }">
          <Tag :color="getModelTypeColor(record.type)">
            {{ getModelTypeText(record.type) }}
          </Tag>
        </template>

        <!-- 状态 -->
        <template #status="{ record }">
          <Badge 
            :status="getModelStatusColor(record.status) as any" 
            :text="getModelStatusText(record.status)"
          />
        </template>

        <!-- 准确率 -->
        <template #accuracy="{ record }">
          <div class="accuracy-display">
            <Progress
              :percent="record.accuracy * 100"
              size="small"
              :stroke-color="record.accuracy >= 0.9 ? '#52c41a' : record.accuracy >= 0.8 ? '#faad14' : '#ff4d4f'"
            />
            <div class="accuracy-text">{{ (record.accuracy * 100).toFixed(1) }}%</div>
          </div>
        </template>

        <!-- 下载量 -->
        <template #downloads="{ record }">
          <div class="download-stats">
            <div class="download-count">
              <DownloadOutlined style="margin-right: 4px" />
              {{ record.downloadCount }}
            </div>
            <div class="use-count">使用: {{ record.useCount }}</div>
          </div>
        </template>

        <!-- 创建者 -->
        <template #creator="{ record }">
          <div class="creator-info">
            <Avatar size="small" style="margin-right: 8px">
              {{ record.creatorName?.[0] }}
            </Avatar>
            <span>{{ record.creatorName }}</span>
          </div>
        </template>

        <!-- 创建时间 -->
        <template #createTime="{ record }">
          <div class="time-info">
            <div>{{ formatDateTime(record.createTime, 'MM-DD') }}</div>
            <div class="time-detail">{{ formatDateTime(record.createTime, 'HH:mm') }}</div>
          </div>
        </template>

        <!-- 操作 -->
        <template #action="{ record }">
          <Space size="small">
            <Tooltip title="查看详情">
              <Button type="text" size="small" @click="viewModelDetail(record)">
                <EyeOutlined />
              </Button>
            </Tooltip>
            <Tooltip title="版本管理">
              <Button type="text" size="small" @click="viewModelVersions(record)">
                <FileTextOutlined />
              </Button>
            </Tooltip>
            <Tooltip title="下载模型">
              <Button type="text" size="small" @click="downloadModelFile(record)">
                <DownloadOutlined />
              </Button>
            </Tooltip>
            <Tooltip title="部署模型">
              <Button 
                type="text" 
                size="small" 
                @click="deployModel(record)"
                :disabled="record.status !== ModelStatus.AVAILABLE"
              >
                <RocketOutlined />
              </Button>
            </Tooltip>
            
            <Dropdown>
              <Button type="text" size="small">
                <MoreOutlined />
              </Button>
              <template #overlay>
                <Menu>
                  <Menu.Item @click="toggleFavorite(record)">
                    <HeartOutlined v-if="!record.isFavorite" />
                    <HeartFilled v-else style="color: #f5222d" />
                    {{ record.isFavorite ? '取消收藏' : '收藏' }}
                  </Menu.Item>
                  <Menu.Item @click="editModel(record)">
                    <EditOutlined />
                    编辑
                  </Menu.Item>
                  <Menu.Item @click="deleteModelItem(record)" danger>
                    <DeleteOutlined />
                    删除
                  </Menu.Item>
                </Menu>
              </template>
            </Dropdown>
          </Space>
        </template>
      </Table>
    </Card>

    <!-- 上传模型模态框 -->
    <Modal
      v-model:open="uploadModalVisible"
      title="上传模型"
      width="800px"
      @ok="handleUploadSubmit"
      @cancel="handleUploadCancel"
      :confirm-loading="loading"
    >
      <Form
        ref="uploadFormRef"
        :model="uploadForm"
        :rules="uploadFormRules"
        layout="vertical"
      >
        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="模型名称" name="name">
              <Input v-model:value="uploadForm.name" placeholder="请输入模型名称" />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="模型类型" name="type">
              <Select v-model:value="uploadForm.type" placeholder="选择模型类型">
                <Select.Option :value="ModelType.CLASSIFICATION">分类</Select.Option>
                <Select.Option :value="ModelType.REGRESSION">回归</Select.Option>
                <Select.Option :value="ModelType.OBJECT_DETECTION">目标检测</Select.Option>
                <Select.Option :value="ModelType.SEMANTIC_SEGMENTATION">语义分割</Select.Option>
                <Select.Option :value="ModelType.NLP">自然语言处理</Select.Option>
                <Select.Option :value="ModelType.RECOMMENDATION">推荐系统</Select.Option>
                <Select.Option :value="ModelType.GENERATIVE">生成模型</Select.Option>
                <Select.Option :value="ModelType.CUSTOM">自定义</Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="框架" name="framework">
              <Select v-model:value="uploadForm.framework" placeholder="选择框架">
                <Select.Option 
                  v-for="framework in supportedFrameworks" 
                  :key="framework" 
                  :value="framework"
                >
                  {{ framework }}
                </Select.Option>
              </Select>
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="工作空间" name="workspaceId">
              <Select v-model:value="uploadForm.workspaceId" placeholder="选择工作空间">
                <Select.Option value="workspace-001">默认工作空间</Select.Option>
                <Select.Option value="workspace-002">CV实验室</Select.Option>
                <Select.Option value="workspace-003">NLP研究室</Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Form.Item label="描述" name="description">
          <Input.TextArea
            v-model:value="uploadForm.description"
            placeholder="请输入模型描述"
            :rows="3"
          />
        </Form.Item>

        <Form.Item label="标签" name="tags">
          <Select
            v-model:value="uploadForm.tags"
            mode="tags"
            placeholder="输入标签，按回车添加"
            style="width: 100%"
          />
        </Form.Item>

        <Form.Item label="模型文件" required>
          <Upload
            :before-upload="beforeUpload"
            :file-list="uploadForm.file ? [{ uid: '1', name: uploadForm.file.name, status: 'done' }] : []"
            accept=".h5,.pkl,.pth,.pt,.onnx,.pb,.tflite,.tar.gz,.zip"
          >
            <Button>
              <UploadOutlined />
              选择文件
            </Button>
            <div style="margin-top: 8px; color: #999; font-size: 12px">
              支持格式: .h5, .pkl, .pth, .pt, .onnx, .pb, .tflite, .tar.gz, .zip
            </div>
          </Upload>
        </Form.Item>

        <Form.Item name="isPublic">
          <Space>
            <span>可见性：</span>
            <Select v-model:value="uploadForm.isPublic" style="width: 120px">
              <Select.Option :value="false">私有</Select.Option>
              <Select.Option :value="true">公开</Select.Option>
            </Select>
          </Space>
        </Form.Item>
      </Form>
    </Modal>

    <!-- 模型详情抽屉 -->
    <ModelDetailDrawer
      v-model:visible="detailDrawerVisible"
      :model="selectedModel"
    />

    <!-- 版本管理抽屉 -->
    <ModelVersionDrawer
      v-model:visible="versionDrawerVisible"
      :model="selectedModel"
    />

    <!-- 部署管理抽屉 -->
    <ModelDeployDrawer
      v-model:visible="deployDrawerVisible"
      :model="selectedModel"
    />
  </div>
</template>

<style scoped lang="scss">
.model-list-container {
  padding: 24px;
  background: #f5f5f5;
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .header-left {
    h2 {
      margin: 0;
      color: #1890ff;
    }
    
    p {
      margin: 8px 0 0 0;
      color: #666;
    }
  }
}

.search-filters {
  .ant-row {
    align-items: center;
  }
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.model-info {
  .model-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 8px;
    
    .model-name {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .name-text {
        font-weight: 600;
        font-size: 14px;
        color: #1890ff;
      }
      
      .favorite-icon {
        font-size: 12px;
      }
    }
    
    .model-version {
      font-size: 12px;
      color: #999;
      font-family: 'Monaco', 'Consolas', monospace;
    }
  }
  
  .model-desc {
    color: #666;
    font-size: 12px;
    margin-bottom: 8px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  
  .model-tags {
    margin-bottom: 8px;
  }
  
  .model-metrics {
    display: flex;
    gap: 16px;
    
    .metric-item {
      font-size: 11px;
      color: #999;
    }
  }
}

.accuracy-display {
  .accuracy-text {
    font-size: 12px;
    text-align: center;
    margin-top: 4px;
  }
}

.download-stats {
  .download-count {
    font-weight: 500;
    margin-bottom: 2px;
  }
  
  .use-count {
    font-size: 12px;
    color: #999;
  }
}

.creator-info {
  display: flex;
  align-items: center;
  
  span {
    font-size: 12px;
  }
}

.time-info {
  .time-detail {
    font-size: 12px;
    color: #999;
    margin-top: 2px;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .model-list-container {
    padding: 16px;
  }
  
  .search-filters {
    .ant-col {
      margin-bottom: 12px;
    }
  }
  
  .toolbar {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
}
</style>
