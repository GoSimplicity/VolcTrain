<template>
  <div class="notebook-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <ExperimentOutlined class="title-icon" />
            <span class="title-text">Notebook 服务</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text"
              >管理和监控您的 Jupyter Notebook 实例</span
            >
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
            创建 Notebook
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
              v-model:value="filterStatus"
              placeholder="选择状态"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部状态</a-select-option>
              <a-select-option value="running">运行中</a-select-option>
              <a-select-option value="pending">等待中</a-select-option>
              <a-select-option value="stopped">已停止</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="6">
            <a-select
              v-model:value="filterNamespace"
              placeholder="选择命名空间"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部命名空间</a-select-option>
              <a-select-option value="default">default</a-select-option>
              <a-select-option value="ai-training">ai-training</a-select-option>
              <a-select-option value="research">research</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="16" :md="8" :lg="8">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索 Notebook 名称或创建者"
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
          :data-source="filteredNotebooks"
          :loading="loading"
          :pagination="paginationConfig"
          row-key="id"
          size="middle"
          :scroll="{ x: 'max-content' }"
          class="sci-fi-table"
        >
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

          <!-- 资源配置列 -->
          <template #resources="{ record }">
            <div class="resources-info">
              <div class="resource-item">
                <DatabaseOutlined class="resource-icon" />
                <span class="resource-label">CPU:</span>
                <span class="resource-value">{{ record.resources.cpu }}</span>
              </div>
              <div class="resource-item">
                <ThunderboltOutlined class="resource-icon" />
                <span class="resource-label">内存:</span>
                <span class="resource-value"
                  >{{ record.resources.memory }}GB</span
                >
              </div>
              <div class="resource-item" v-if="record.resources.gpu">
                <BugOutlined class="resource-icon" />
                <span class="resource-label">GPU:</span>
                <span class="resource-value">{{ record.resources.gpu }}</span>
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
                @click="openNotebook(record)"
                :disabled="record.status !== 'running'"
                class="action-btn"
              >
                <LinkOutlined />
                打开
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
                      key="start"
                      :disabled="record.status === 'running'"
                    >
                      <PlayCircleOutlined />
                      启动
                    </a-menu-item>
                    <a-menu-item
                      key="stop"
                      :disabled="record.status === 'stopped'"
                    >
                      <PauseCircleOutlined />
                      停止
                    </a-menu-item>
                    <a-menu-item
                      key="restart"
                      :disabled="record.status !== 'running'"
                    >
                      <ReloadOutlined />
                      重启
                    </a-menu-item>
                    <a-menu-item key="clone">
                      <CopyOutlined />
                      克隆
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

    <!-- 创建 Notebook 模态框 -->
    <a-modal
      v-model:open="createModalVisible"
      title="创建 Notebook"
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
            <a-form-item label="Notebook 名称" name="name">
              <a-input
                v-model:value="createForm.name"
                placeholder="请输入名称"
                class="form-input"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="命名空间" name="namespace">
              <a-select
                v-model:value="createForm.namespace"
                placeholder="选择命名空间"
                class="form-select"
              >
                <a-select-option value="default">default</a-select-option>
                <a-select-option value="ai-training"
                  >ai-training</a-select-option
                >
                <a-select-option value="research">research</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="镜像" name="image">
          <a-select
            v-model:value="createForm.image"
            placeholder="选择镜像"
            class="form-select"
          >
            <a-select-option value="jupyter/tensorflow-notebook:latest">
              TensorFlow Notebook
            </a-select-option>
            <a-select-option value="jupyter/pytorch-notebook:latest">
              PyTorch Notebook
            </a-select-option>
            <a-select-option value="jupyter/datascience-notebook:latest">
              Data Science Notebook
            </a-select-option>
            <a-select-option value="custom">自定义镜像</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item
          v-if="createForm.image === 'custom'"
          label="自定义镜像地址"
          name="customImage"
        >
          <a-input
            v-model:value="createForm.customImage"
            placeholder="请输入镜像地址"
            class="form-input"
          />
        </a-form-item>

        <a-divider class="form-divider">资源配置</a-divider>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="8">
            <a-form-item label="CPU 核数" name="cpu">
              <a-input-number
                v-model:value="createForm.cpu"
                :min="0.5"
                :max="32"
                :step="0.5"
                style="width: 100%"
                addon-after="核"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="8">
            <a-form-item label="内存" name="memory">
              <a-input-number
                v-model:value="createForm.memory"
                :min="1"
                :max="128"
                style="width: 100%"
                addon-after="GB"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="8">
            <a-form-item label="GPU 卡数" name="gpu">
              <a-input-number
                v-model:value="createForm.gpu"
                :min="0"
                :max="8"
                style="width: 100%"
                addon-after="卡"
                class="form-input-number"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="存储卷" name="storage">
          <a-input-number
            v-model:value="createForm.storage"
            :min="10"
            :max="1000"
            style="width: 100%"
            addon-after="GB"
            class="form-input-number"
          />
        </a-form-item>

        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="createForm.description"
            placeholder="请输入描述信息"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="Notebook 详情"
      width="900px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedNotebook" class="detail-content">
        <a-descriptions
          :column="{ xs: 1, sm: 2 }"
          bordered
          class="detail-descriptions"
        >
          <a-descriptions-item label="名称">
            {{ selectedNotebook.name }}
          </a-descriptions-item>
          <a-descriptions-item label="状态">
            <a-tag
              :color="getStatusColor(selectedNotebook.status)"
              class="status-tag"
            >
              <component :is="getStatusIcon(selectedNotebook.status)" />
              {{ getStatusText(selectedNotebook.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="命名空间">
            {{ selectedNotebook.namespace }}
          </a-descriptions-item>
          <a-descriptions-item label="创建者">
            {{ selectedNotebook.creator }}
          </a-descriptions-item>
          <a-descriptions-item label="镜像">
            {{ selectedNotebook.image }}
          </a-descriptions-item>
          <a-descriptions-item label="创建时间">
            {{ selectedNotebook.createTime }}
          </a-descriptions-item>
          <a-descriptions-item label="CPU">
            {{ selectedNotebook.resources.cpu }} 核
          </a-descriptions-item>
          <a-descriptions-item label="内存">
            {{ selectedNotebook.resources.memory }} GB
          </a-descriptions-item>
          <a-descriptions-item
            label="GPU"
            v-if="selectedNotebook.resources.gpu"
          >
            {{ selectedNotebook.resources.gpu }} 卡
          </a-descriptions-item>
          <a-descriptions-item label="存储">
            {{ selectedNotebook.storage }} GB
          </a-descriptions-item>
          <a-descriptions-item
            label="访问地址"
            :span="2"
            v-if="selectedNotebook.status === 'running'"
          >
            <a
              :href="selectedNotebook.url"
              target="_blank"
              class="notebook-link"
            >
              {{ selectedNotebook.url }}
              <ExportOutlined />
            </a>
          </a-descriptions-item>
          <a-descriptions-item label="描述" :span="2">
            {{ selectedNotebook.description || '暂无描述' }}
          </a-descriptions-item>
        </a-descriptions>

        <!-- 实时日志 -->
        <a-divider class="log-divider">实时日志</a-divider>
        <div class="log-container">
          <div class="log-header">
            <span class="log-title">容器日志</span>
            <a-button size="small" @click="refreshLogs" class="log-refresh-btn">
              <ReloadOutlined />
              刷新
            </a-button>
          </div>
          <div class="log-content">
            <pre v-for="(log, index) in logs" :key="index" class="log-line">{{
              log
            }}</pre>
          </div>
        </div>
      </div>
    </a-modal>

    <!-- 克隆模态框 -->
    <a-modal
      v-model:open="cloneModalVisible"
      title="克隆 Notebook"
      width="600px"
      :confirm-loading="cloneLoading"
      @ok="handleCloneSubmit"
      @cancel="handleCloneCancel"
      class="sci-fi-modal"
    >
      <a-form
        ref="cloneFormRef"
        :model="cloneForm"
        :rules="cloneFormRules"
        layout="vertical"
      >
        <a-form-item label="新 Notebook 名称" name="name">
          <a-input
            v-model:value="cloneForm.name"
            placeholder="请输入新名称"
            class="form-input"
          />
        </a-form-item>
        <a-form-item label="命名空间" name="namespace">
          <a-select
            v-model:value="cloneForm.namespace"
            placeholder="选择命名空间"
            class="form-select"
          >
            <a-select-option value="default">default</a-select-option>
            <a-select-option value="ai-training">ai-training</a-select-option>
            <a-select-option value="research">research</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import type { FormInstance, TableColumnsType } from 'ant-design-vue';
import {
  ExperimentOutlined,
  PlusOutlined,
  ReloadOutlined,
  DatabaseOutlined,
  ThunderboltOutlined,
  BugOutlined,
  LinkOutlined,
  EyeOutlined,
  MoreOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  CopyOutlined,
  DeleteOutlined,
  ExportOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  StopOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface NotebookResources {
  cpu: number;
  memory: number;
  gpu?: number;
}

interface NotebookItem {
  id: string;
  name: string;
  namespace: string;
  status: 'running' | 'pending' | 'stopped' | 'failed';
  creator: string;
  image: string;
  createTime: string;
  resources: NotebookResources;
  storage: number;
  url?: string;
  description?: string;
}

interface CreateForm {
  name: string;
  namespace: string;
  image: string;
  customImage: string;
  cpu: number;
  memory: number;
  gpu: number;
  storage: number;
  description: string;
}

interface CloneForm {
  name: string;
  namespace: string;
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const createModalVisible = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const cloneModalVisible = ref<boolean>(false);
const createLoading = ref<boolean>(false);
const cloneLoading = ref<boolean>(false);

const filterStatus = ref<string>('');
const filterNamespace = ref<string>('');
const searchKeyword = ref<string>('');

const selectedNotebook = ref<NotebookItem | null>(null);
const logs = ref<string[]>([
  '2024-06-23 10:30:15 INFO: Starting Jupyter Notebook server...',
  '2024-06-23 10:30:16 INFO: Kernel started successfully',
  '2024-06-23 10:30:17 INFO: Server is running on port 8888',
  '2024-06-23 10:30:18 INFO: Ready to accept connections',
]);

// ===== 表单引用 =====
const createFormRef = ref<FormInstance>();
const cloneFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const createForm = reactive<CreateForm>({
  name: '',
  namespace: 'default',
  image: 'jupyter/tensorflow-notebook:latest',
  customImage: '',
  cpu: 2,
  memory: 4,
  gpu: 0,
  storage: 50,
  description: '',
});

const cloneForm = reactive<CloneForm>({
  name: '',
  namespace: 'default',
});

// ===== 配置数据 =====
const STATUS_CONFIG = {
  running: { color: 'success', text: '运行中', icon: CheckCircleOutlined },
  pending: { color: 'processing', text: '等待中', icon: ClockCircleOutlined },
  stopped: { color: 'default', text: '已停止', icon: StopOutlined },
  failed: { color: 'error', text: '失败', icon: CloseCircleOutlined },
} as const;

// ===== 模拟数据 =====
const notebooks = ref<NotebookItem[]>([
  {
    id: 'nb-001',
    name: 'tensorflow-training',
    namespace: 'ai-training',
    status: 'running',
    creator: 'admin',
    image: 'jupyter/tensorflow-notebook:latest',
    createTime: '2024-06-23 09:30:00',
    resources: { cpu: 4, memory: 8, gpu: 1 },
    storage: 100,
    url: 'https://notebook-001.example.com',
    description: 'TensorFlow 模型训练环境',
  },
  {
    id: 'nb-002',
    name: 'data-analysis',
    namespace: 'research',
    status: 'running',
    creator: 'researcher',
    image: 'jupyter/datascience-notebook:latest',
    createTime: '2024-06-23 08:15:00',
    resources: { cpu: 2, memory: 4 },
    storage: 50,
    url: 'https://notebook-002.example.com',
    description: '数据分析工作环境',
  },
  {
    id: 'nb-003',
    name: 'pytorch-experiment',
    namespace: 'ai-training',
    status: 'pending',
    creator: 'developer',
    image: 'jupyter/pytorch-notebook:latest',
    createTime: '2024-06-23 10:00:00',
    resources: { cpu: 8, memory: 16, gpu: 2 },
    storage: 200,
    description: 'PyTorch 深度学习实验',
  },
  {
    id: 'nb-004',
    name: 'ml-pipeline',
    namespace: 'default',
    status: 'stopped',
    creator: 'admin',
    image: 'jupyter/tensorflow-notebook:latest',
    createTime: '2024-06-22 16:30:00',
    resources: { cpu: 4, memory: 8 },
    storage: 80,
    description: '机器学习流水线开发',
  },
  {
    id: 'nb-005',
    name: 'debug-session',
    namespace: 'default',
    status: 'failed',
    creator: 'developer',
    image: 'custom/debug-notebook:v1.0',
    createTime: '2024-06-23 09:45:00',
    resources: { cpu: 2, memory: 4 },
    storage: 30,
    description: '调试会话环境',
  },
]);

// ===== 表单验证规则 =====
const createFormRules = {
  name: [
    { required: true, message: '请输入 Notebook 名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
    {
      pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
      message: '名称只能包含小写字母、数字和连字符',
      trigger: 'blur',
    },
  ],
  namespace: [{ required: true, message: '请选择命名空间', trigger: 'change' }],
  image: [{ required: true, message: '请选择镜像', trigger: 'change' }],
  customImage: [
    { required: true, message: '请输入自定义镜像地址', trigger: 'blur' },
  ],
  cpu: [{ required: true, message: '请输入 CPU 核数', trigger: 'blur' }],
  memory: [{ required: true, message: '请输入内存大小', trigger: 'blur' }],
  storage: [{ required: true, message: '请输入存储大小', trigger: 'blur' }],
};

const cloneFormRules = {
  name: [
    { required: true, message: '请输入新 Notebook 名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
    {
      pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
      message: '名称只能包含小写字母、数字和连字符',
      trigger: 'blur',
    },
  ],
  namespace: [{ required: true, message: '请选择命名空间', trigger: 'change' }],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<NotebookItem> = [
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
    ellipsis: true,
  },
  {
    title: '命名空间',
    dataIndex: 'namespace',
    key: 'namespace',
    width: 120,
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
    title: '镜像',
    dataIndex: 'image',
    key: 'image',
    width: 200,
    ellipsis: true,
  },
  {
    title: '资源配置',
    key: 'resources',
    width: 200,
    slots: { customRender: 'resources' },
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
  total: computed(() => filteredNotebooks.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const filteredNotebooks = computed(() => {
  let result = notebooks.value;

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (filterNamespace.value) {
    result = result.filter((item) => item.namespace === filterNamespace.value);
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter(
      (item) =>
        item.name.toLowerCase().includes(keyword) ||
        item.creator.toLowerCase().includes(keyword),
    );
  }

  return result;
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

    const newNotebook: NotebookItem = {
      id: `nb-${Date.now()}`,
      name: createForm.name,
      namespace: createForm.namespace,
      status: 'pending',
      creator: 'current-user',
      image:
        createForm.image === 'custom'
          ? createForm.customImage
          : createForm.image,
      createTime: new Date().toLocaleString(),
      resources: {
        cpu: createForm.cpu,
        memory: createForm.memory,
        ...(createForm.gpu > 0 && { gpu: createForm.gpu }),
      },
      storage: createForm.storage,
      description: createForm.description,
    };

    notebooks.value.unshift(newNotebook);
    createModalVisible.value = false;
    message.success('Notebook 创建成功');

    // 重置表单
    createFormRef.value?.resetFields();
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

const openNotebook = (record: NotebookItem): void => {
  if (record.url) {
    window.open(record.url, '_blank');
  }
};

const viewDetails = (record: NotebookItem): void => {
  selectedNotebook.value = record;
  detailModalVisible.value = true;
};

const handleMenuAction = (key: string, record: NotebookItem): void => {
  const actions = {
    start: () => handleStart(record),
    stop: () => handleStop(record),
    restart: () => handleRestart(record),
    clone: () => handleClone(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handleStart = async (record: NotebookItem): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const index = notebooks.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      notebooks.value[index]!.status = 'running';
      notebooks.value[index]!.url = `https://notebook-${record.id}.example.com`;
    }
    message.success('Notebook 启动成功');
  } catch (error) {
    message.error('启动失败');
  } finally {
    loading.value = false;
  }
};

const handleStop = async (record: NotebookItem): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 1000));
    const index = notebooks.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      notebooks.value[index]!.status = 'stopped';
      delete notebooks.value[index]!.url;
    }
    message.success('Notebook 停止成功');
  } catch (error) {
    message.error('停止失败');
  } finally {
    loading.value = false;
  }
};

const handleRestart = async (record: NotebookItem): Promise<void> => {
  loading.value = true;
  try {
    await new Promise((resolve) => setTimeout(resolve, 2000));
    const index = notebooks.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      notebooks.value[index]!.status = 'running';
    }
    message.success('Notebook 重启成功');
  } catch (error) {
    message.error('重启失败');
  } finally {
    loading.value = false;
  }
};

const handleClone = (record: NotebookItem): void => {
  cloneForm.name = `${record.name}-copy`;
  cloneForm.namespace = record.namespace;
  selectedNotebook.value = record;
  cloneModalVisible.value = true;
};

const handleCloneSubmit = async (): Promise<void> => {
  try {
    await cloneFormRef.value?.validate();
    cloneLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1500));

    if (selectedNotebook.value) {
      const clonedNotebook: NotebookItem = {
        ...selectedNotebook.value,
        id: `nb-${Date.now()}`,
        name: cloneForm.name,
        namespace: cloneForm.namespace,
        status: 'pending',
        createTime: new Date().toLocaleString(),
        url: undefined,
      };

      notebooks.value.unshift(clonedNotebook);
      cloneModalVisible.value = false;
      message.success('Notebook 克隆成功');
    }
  } catch (error) {
    message.error('表单验证失败');
  } finally {
    cloneLoading.value = false;
  }
};

const handleCloneCancel = (): void => {
  cloneModalVisible.value = false;
  cloneFormRef.value?.resetFields();
};

const handleDelete = (record: NotebookItem): void => {
  const deleteConfirm = () => {
    const index = notebooks.value.findIndex((item) => item.id === record.id);
    if (index !== -1) {
      notebooks.value.splice(index, 1);
      message.success('Notebook 删除成功');
    }
  };

  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除 Notebook "${record.name}" 吗？此操作不可恢复。`,
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

const refreshLogs = async (): Promise<void> => {
  const newLogs = [
    ...logs.value,
    `${new Date().toLocaleString()} INFO: Log refreshed`,
  ];
  logs.value = newLogs.slice(-50);
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
.notebook-container {
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

.indicator-running {
  background: #52c41a;
}

.indicator-pending {
  background: #1890ff;
}

.indicator-stopped {
  background: #8c8c8c;
}

.indicator-failed {
  background: #ff4d4f;
}

/* ===== 资源信息 ===== */
.resources-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.resource-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.resource-item:hover {
  color: #1890ff;
}

.resource-icon {
  font-size: 12px;
  color: #1890ff;
}

.resource-label {
  font-weight: 500;
}

.resource-value {
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
.create-form :deep(.ant-form-item-label > label) {
  font-weight: 500 !important;
}

.form-input,
.form-select,
.form-textarea,
.form-input-number {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.form-divider {
  font-weight: 500 !important;
}

/* ===== 详情页样式 ===== */
.detail-content {
  max-height: 600px;
  overflow-y: auto;
}

.notebook-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  text-decoration: none;
  transition: all 0.3s ease;
}

.notebook-link:hover {
  text-decoration: underline;
}

/* ===== 日志容器 ===== */
.log-divider {
  font-weight: 500 !important;
  margin: 24px 0 16px 0 !important;
}

.log-container {
  margin-top: 16px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.log-title {
  font-weight: 600;
  font-size: 14px;
}

.log-refresh-btn {
  border-radius: 4px !important;
  transition: all 0.3s ease;
}

.log-refresh-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

.log-content {
  border-radius: 6px !important;
  padding: 12px !important;
  max-height: 300px;
  overflow-y: auto;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace !important;
}

.log-line {
  margin: 0;
  font-size: 12px;
  line-height: 1.4;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .notebook-container {
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

  .resources-info {
    gap: 2px;
  }

  .resource-item {
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
