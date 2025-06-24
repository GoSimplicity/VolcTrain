<template>
  <div class="training-history-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <HistoryOutlined class="title-icon" />
            <span class="title-text">训练任务历史</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text"
              >查看和管理所有训练任务的历史记录</span
            >
          </p>
        </div>
        <div class="action-section">
          <a-button
            type="primary"
            size="large"
            @click="exportHistory"
            class="export-btn"
          >
            <ExportOutlined />
            导出记录
          </a-button>
        </div>
      </div>
    </div>

    <!-- 筛选器 -->
    <div class="filter-section">
      <a-card class="filter-card glass-card" :bordered="false">
        <a-row :gutter="16" align="middle">
          <a-col :xs="24" :sm="12" :md="6" :lg="5">
            <a-select
              v-model:value="filterStatus"
              placeholder="选择状态"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部状态</a-select-option>
              <a-select-option value="completed">已完成</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
              <a-select-option value="cancelled">已取消</a-select-option>
              <a-select-option value="timeout">超时</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="5">
            <a-select
              v-model:value="filterFramework"
              placeholder="选择框架"
              allow-clear
              style="width: 100%"
              @change="handleFilterChange"
              class="filter-select"
            >
              <a-select-option value="">全部框架</a-select-option>
              <a-select-option value="tensorflow">TensorFlow</a-select-option>
              <a-select-option value="pytorch">PyTorch</a-select-option>
              <a-select-option value="mxnet">MXNet</a-select-option>
              <a-select-option value="keras">Keras</a-select-option>
            </a-select>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="5">
            <a-range-picker
              v-model:value="dateRange"
              style="width: 100%"
              @change="handleDateRangeChange"
              class="date-picker"
              placeholder="['开始时间', '结束时间']"
            />
          </a-col>
          <a-col :xs="24" :sm="12" :md="6" :lg="5">
            <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索任务名称或创建者"
              allow-clear
              @search="handleSearch"
              @change="handleSearchChange"
              class="search-input"
            />
          </a-col>
          <a-col :xs="24" :sm="24" :md="24" :lg="4" class="action-col">
            <a-space>
              <a-button
                @click="refreshData"
                :loading="loading"
                class="refresh-btn"
              >
                <ReloadOutlined />
                刷新
              </a-button>
              <a-button @click="clearFilters" class="clear-btn">
                <ClearOutlined />
                清空
              </a-button>
            </a-space>
          </a-col>
        </a-row>
      </a-card>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <a-row :gutter="16">
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <a-statistic title="总任务数" :value="totalTasks" class="stat-item">
              <template #prefix>
                <BarsOutlined class="stat-icon stat-icon-total" />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <a-statistic
              title="成功任务"
              :value="completedTasks"
              class="stat-item"
            >
              <template #prefix>
                <CheckCircleOutlined class="stat-icon stat-icon-success" />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <a-statistic
              title="失败任务"
              :value="failedTasks"
              class="stat-item"
            >
              <template #prefix>
                <CloseCircleOutlined class="stat-icon stat-icon-error" />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
        <a-col :xs="12" :sm="6" :md="6" :lg="6">
          <a-card class="stat-card glass-card" :bordered="false">
            <a-statistic
              title="成功率"
              :value="successRate"
              suffix="%"
              :precision="1"
              class="stat-item"
            >
              <template #prefix>
                <TrophyOutlined class="stat-icon stat-icon-rate" />
              </template>
            </a-statistic>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <!-- 数据表格 -->
    <div class="table-section">
      <a-card class="table-card glass-card" :bordered="false">
        <a-table
          :columns="columns"
          :data-source="filteredTasks"
          :loading="loading"
          :pagination="paginationConfig"
          row-key="id"
          size="middle"
          :scroll="{ x: 'max-content' }"
          class="sci-fi-table"
        >
          <!-- 任务名称列 -->
          <template #taskName="{ record }">
            <div class="task-name-wrapper">
              <a-button
                type="link"
                @click="viewTaskDetails(record)"
                class="task-name-link"
              >
                {{ record.taskName }}
              </a-button>
              <a-tag
                v-if="record.isTemplate"
                color="purple"
                size="small"
                class="template-tag"
              >
                模板
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

          <!-- 框架列 -->
          <template #framework="{ record }">
            <div class="framework-wrapper">
              <component
                :is="getFrameworkIcon(record.framework)"
                class="framework-icon"
              />
              <span class="framework-text">{{ record.framework }}</span>
            </div>
          </template>

          <!-- 资源使用列 -->
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

          <!-- 持续时间列 -->
          <template #duration="{ record }">
            <span class="duration-text">{{
              formatDuration(record.duration)
            }}</span>
          </template>

          <!-- 开始时间列 -->
          <template #startTime="{ record }">
            <a-tooltip :title="record.startTime">
              <span class="time-text">{{
                formatRelativeTime(record.startTime)
              }}</span>
            </a-tooltip>
          </template>

          <!-- 操作列 -->
          <template #action="{ record }">
            <a-space class="action-buttons">
              <a-button
                type="link"
                size="small"
                @click="viewTaskDetails(record)"
                class="action-btn"
              >
                <EyeOutlined />
                详情
              </a-button>
              <a-button
                type="link"
                size="small"
                @click="viewLogs(record)"
                class="action-btn"
              >
                <FileTextOutlined />
                日志
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
                    <a-menu-item key="clone">
                      <CopyOutlined />
                      复制任务
                    </a-menu-item>
                    <a-menu-item key="template">
                      <SaveOutlined />
                      存为模板
                    </a-menu-item>
                    <a-menu-item key="download">
                      <DownloadOutlined />
                      下载模型
                    </a-menu-item>
                    <a-menu-divider />
                    <a-menu-item key="delete" class="danger-item">
                      <DeleteOutlined />
                      删除记录
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </a-space>
          </template>
        </a-table>
      </a-card>
    </div>

    <!-- 任务详情模态框 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="训练任务详情"
      width="1000px"
      :footer="null"
      class="sci-fi-modal detail-modal"
    >
      <div v-if="selectedTask" class="detail-content">
        <a-tabs v-model:activeKey="activeTab" class="detail-tabs">
          <!-- 基本信息 -->
          <a-tab-pane key="basic" tab="基本信息">
            <a-descriptions
              :column="{ xs: 1, sm: 2 }"
              bordered
              class="detail-descriptions"
            >
              <a-descriptions-item label="任务名称">
                {{ selectedTask.taskName }}
              </a-descriptions-item>
              <a-descriptions-item label="状态">
                <a-tag
                  :color="getStatusColor(selectedTask.status)"
                  class="status-tag"
                >
                  <component :is="getStatusIcon(selectedTask.status)" />
                  {{ getStatusText(selectedTask.status) }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="创建者">
                {{ selectedTask.creator }}
              </a-descriptions-item>
              <a-descriptions-item label="框架">
                <div class="framework-wrapper">
                  <component
                    :is="getFrameworkIcon(selectedTask.framework)"
                    class="framework-icon"
                  />
                  {{ selectedTask.framework }}
                </div>
              </a-descriptions-item>
              <a-descriptions-item label="开始时间">
                {{ selectedTask.startTime }}
              </a-descriptions-item>
              <a-descriptions-item label="结束时间">
                {{ selectedTask.endTime || '未结束' }}
              </a-descriptions-item>
              <a-descriptions-item label="持续时间">
                {{ formatDuration(selectedTask.duration) }}
              </a-descriptions-item>
              <a-descriptions-item label="优先级">
                <a-tag :color="getPriorityColor(selectedTask.priority)">
                  {{ selectedTask.priority }}
                </a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="CPU">
                {{ selectedTask.resources.cpu }} 核
              </a-descriptions-item>
              <a-descriptions-item label="内存">
                {{ selectedTask.resources.memory }} GB
              </a-descriptions-item>
              <a-descriptions-item
                label="GPU"
                v-if="selectedTask.resources.gpu"
              >
                {{ selectedTask.resources.gpu }} 卡
              </a-descriptions-item>
              <a-descriptions-item label="存储">
                {{ selectedTask.storage }} GB
              </a-descriptions-item>
              <a-descriptions-item label="描述" :span="2">
                {{ selectedTask.description || '暂无描述' }}
              </a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <!-- 训练配置 -->
          <a-tab-pane key="config" tab="训练配置">
            <div class="config-content">
              <a-descriptions
                :column="{ xs: 1, sm: 2 }"
                bordered
                class="config-descriptions"
              >
                <a-descriptions-item label="数据集">
                  {{ selectedTask.dataset }}
                </a-descriptions-item>
                <a-descriptions-item label="模型类型">
                  {{ selectedTask.modelType }}
                </a-descriptions-item>
                <a-descriptions-item label="批次大小">
                  {{ selectedTask.batchSize }}
                </a-descriptions-item>
                <a-descriptions-item label="学习率">
                  {{ selectedTask.learningRate }}
                </a-descriptions-item>
                <a-descriptions-item label="训练轮数">
                  {{ selectedTask.epochs }}
                </a-descriptions-item>
                <a-descriptions-item label="验证频率">
                  每 {{ selectedTask.validationFreq }} 轮
                </a-descriptions-item>
              </a-descriptions>
            </div>
          </a-tab-pane>

          <!-- 训练指标 -->
          <a-tab-pane key="metrics" tab="训练指标">
            <div class="metrics-content">
              <a-row :gutter="16">
                <a-col :xs="24" :sm="12">
                  <div class="metric-card">
                    <h4 class="metric-title">最终精度</h4>
                    <div class="metric-value">
                      {{ selectedTask.finalAccuracy }}%
                    </div>
                  </div>
                </a-col>
                <a-col :xs="24" :sm="12">
                  <div class="metric-card">
                    <h4 class="metric-title">最低损失</h4>
                    <div class="metric-value">{{ selectedTask.bestLoss }}</div>
                  </div>
                </a-col>
              </a-row>
              <div class="chart-placeholder">
                <p>训练指标图表</p>
                <small
                  >在真实环境中，这里会显示训练过程中的损失和精度变化曲线</small
                >
              </div>
            </div>
          </a-tab-pane>

          <!-- 环境配置 -->
          <a-tab-pane key="environment" tab="环境配置">
            <div class="environment-content">
              <a-descriptions
                :column="{ xs: 1, sm: 1 }"
                bordered
                class="environment-descriptions"
              >
                <a-descriptions-item label="Docker 镜像">
                  {{ selectedTask.dockerImage }}
                </a-descriptions-item>
                <a-descriptions-item label="工作目录">
                  {{ selectedTask.workDir }}
                </a-descriptions-item>
                <a-descriptions-item label="命令">
                  <code class="command-code">{{ selectedTask.command }}</code>
                </a-descriptions-item>
                <a-descriptions-item label="环境变量">
                  <div class="env-vars">
                    <div
                      v-for="(value, key) in selectedTask.envVars"
                      :key="key"
                      class="env-var-item"
                    >
                      <span class="env-key">{{ key }}:</span>
                      <span class="env-value">{{ value }}</span>
                    </div>
                  </div>
                </a-descriptions-item>
              </a-descriptions>
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>
    </a-modal>

    <!-- 日志查看模态框 -->
    <a-modal
      v-model:open="logModalVisible"
      title="训练日志"
      width="800px"
      :footer="null"
      class="sci-fi-modal log-modal"
    >
      <div class="log-container">
        <div class="log-header">
          <a-space>
            <span class="log-title">训练日志</span>
            <a-select
              v-model:value="selectedLogLevel"
              style="width: 120px"
              size="small"
              @change="filterLogs"
            >
              <a-select-option value="all">全部级别</a-select-option>
              <a-select-option value="info">INFO</a-select-option>
              <a-select-option value="warning">WARNING</a-select-option>
              <a-select-option value="error">ERROR</a-select-option>
            </a-select>
          </a-space>
          <a-space>
            <a-button size="small" @click="downloadLogs" class="log-action-btn">
              <DownloadOutlined />
              下载
            </a-button>
            <a-button size="small" @click="refreshLogs" class="log-action-btn">
              <ReloadOutlined />
              刷新
            </a-button>
          </a-space>
        </div>
        <div class="log-content">
          <pre
            v-for="(log, index) in filteredLogs"
            :key="index"
            class="log-line"
            :class="`log-${log.level}`"
            >{{ log.timestamp }} [{{ log.level.toUpperCase() }}] {{
              log.message
            }}</pre
          >
        </div>
      </div>
    </a-modal>

    <!-- 复制任务模态框 -->
    <a-modal
      v-model:open="cloneModalVisible"
      title="复制训练任务"
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
        <a-form-item label="新任务名称" name="taskName">
          <a-input
            v-model:value="cloneForm.taskName"
            placeholder="请输入新任务名称"
            class="form-input"
          />
        </a-form-item>
        <a-form-item label="描述" name="description">
          <a-textarea
            v-model:value="cloneForm.description"
            placeholder="请输入任务描述"
            :rows="3"
            class="form-textarea"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import type { FormInstance, TableColumnsType } from 'ant-design-vue';
import type { Dayjs } from 'dayjs';
import {
  HistoryOutlined,
  ExportOutlined,
  ReloadOutlined,
  ClearOutlined,
  BarsOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  TrophyOutlined,
  DatabaseOutlined,
  ThunderboltOutlined,
  BugOutlined,
  EyeOutlined,
  FileTextOutlined,
  MoreOutlined,
  CopyOutlined,
  SaveOutlined,
  DownloadOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  StopOutlined,
  ClockCircleOutlined,
  ApiOutlined,
  CodeOutlined,
  FireOutlined,
  ExperimentOutlined,
} from '@ant-design/icons-vue';

// ===== 类型定义 =====
interface TaskResources {
  cpu: number;
  memory: number;
  gpu?: number;
}

interface LogEntry {
  timestamp: string;
  level: 'info' | 'warning' | 'error';
  message: string;
}

interface TrainingTask {
  id: string;
  taskName: string;
  status: 'completed' | 'failed' | 'cancelled' | 'timeout';
  creator: string;
  framework: 'tensorflow' | 'pytorch' | 'mxnet' | 'keras';
  startTime: string;
  endTime?: string;
  duration: number;
  resources: TaskResources;
  storage: number;
  priority: 'high' | 'medium' | 'low';
  dataset: string;
  modelType: string;
  batchSize: number;
  learningRate: number;
  epochs: number;
  validationFreq: number;
  finalAccuracy: number;
  bestLoss: number;
  dockerImage: string;
  workDir: string;
  command: string;
  envVars: Record<string, string>;
  description?: string;
  isTemplate?: boolean;
}

interface CloneForm {
  taskName: string;
  description: string;
}

// ===== 响应式数据 =====
const loading = ref<boolean>(false);
const detailModalVisible = ref<boolean>(false);
const logModalVisible = ref<boolean>(false);
const cloneModalVisible = ref<boolean>(false);
const cloneLoading = ref<boolean>(false);

const filterStatus = ref<string>('');
const filterFramework = ref<string>('');
const dateRange = ref<[Dayjs, Dayjs] | null>(null);
const searchKeyword = ref<string>('');
const activeTab = ref<string>('basic');
const selectedLogLevel = ref<string>('all');

const selectedTask = ref<TrainingTask | null>(null);

// ===== 表单引用 =====
const cloneFormRef = ref<FormInstance>();

// ===== 表单数据 =====
const cloneForm = reactive<CloneForm>({
  taskName: '',
  description: '',
});

// ===== 配置数据 =====
const STATUS_CONFIG = {
  completed: { color: 'success', text: '已完成', icon: CheckCircleOutlined },
  failed: { color: 'error', text: '失败', icon: CloseCircleOutlined },
  cancelled: { color: 'warning', text: '已取消', icon: StopOutlined },
  timeout: { color: 'default', text: '超时', icon: ClockCircleOutlined },
} as const;

const FRAMEWORK_ICONS = {
  tensorflow: FireOutlined,
  pytorch: ApiOutlined,
  mxnet: CodeOutlined,
  keras: ExperimentOutlined,
} as const;

const PRIORITY_COLORS = {
  high: 'red',
  medium: 'orange',
  low: 'green',
} as const;

// ===== 模拟数据 =====
const trainingTasks = ref<TrainingTask[]>([
  {
    id: 'task-001',
    taskName: 'resnet50-imagenet-training',
    status: 'completed',
    creator: 'admin',
    framework: 'tensorflow',
    startTime: '2024-06-23 08:00:00',
    endTime: '2024-06-23 14:30:00',
    duration: 23400, // 6.5小时，单位秒
    resources: { cpu: 8, memory: 32, gpu: 4 },
    storage: 500,
    priority: 'high',
    dataset: 'ImageNet-1K',
    modelType: 'ResNet-50',
    batchSize: 256,
    learningRate: 0.001,
    epochs: 100,
    validationFreq: 5,
    finalAccuracy: 76.8,
    bestLoss: 0.942,
    dockerImage: 'tensorflow/tensorflow:2.13.0-gpu',
    workDir: '/workspace/training',
    command: 'python train.py --model resnet50 --dataset imagenet',
    envVars: {
      CUDA_VISIBLE_DEVICES: '0,1,2,3',
      TF_ENABLE_GPU_MEMORY_GROWTH: 'true',
    },
    description: 'ImageNet 数据集上的 ResNet-50 模型训练',
    isTemplate: true,
  },
  {
    id: 'task-002',
    taskName: 'bert-base-finetuning',
    status: 'completed',
    creator: 'researcher',
    framework: 'pytorch',
    startTime: '2024-06-22 20:15:00',
    endTime: '2024-06-23 02:45:00',
    duration: 23400, // 6.5小时
    resources: { cpu: 16, memory: 64, gpu: 8 },
    storage: 200,
    priority: 'medium',
    dataset: 'GLUE-CoLA',
    modelType: 'BERT-Base',
    batchSize: 32,
    learningRate: 0.00002,
    epochs: 50,
    validationFreq: 2,
    finalAccuracy: 85.2,
    bestLoss: 0.387,
    dockerImage: 'pytorch/pytorch:2.0.1-cuda11.7-cudnn8-devel',
    workDir: '/workspace/nlp',
    command: 'python finetune_bert.py --task cola --epochs 50',
    envVars: {
      TOKENIZERS_PARALLELISM: 'false',
      PYTHONPATH: '/workspace/nlp',
    },
    description: 'BERT 模型在 CoLA 任务上的微调',
  },
  {
    id: 'task-003',
    taskName: 'yolo-object-detection',
    status: 'failed',
    creator: 'developer',
    framework: 'pytorch',
    startTime: '2024-06-22 16:30:00',
    endTime: '2024-06-22 18:15:00',
    duration: 6300, // 1.75小时
    resources: { cpu: 4, memory: 16, gpu: 2 },
    storage: 150,
    priority: 'low',
    dataset: 'COCO-2017',
    modelType: 'YOLOv8n',
    batchSize: 64,
    learningRate: 0.01,
    epochs: 300,
    validationFreq: 10,
    finalAccuracy: 0,
    bestLoss: 999,
    dockerImage: 'ultralytics/yolov8:latest',
    workDir: '/workspace/detection',
    command: 'yolo train data=coco.yaml model=yolov8n.pt epochs=300',
    envVars: {
      WANDB_MODE: 'offline',
    },
    description: 'YOLO 目标检测模型训练（因数据加载错误失败）',
  },
  {
    id: 'task-004',
    taskName: 'lstm-time-series',
    status: 'completed',
    creator: 'data-scientist',
    framework: 'keras',
    startTime: '2024-06-21 14:00:00',
    endTime: '2024-06-21 16:20:00',
    duration: 8400, // 2.33小时
    resources: { cpu: 2, memory: 8 },
    storage: 50,
    priority: 'medium',
    dataset: 'Stock-Prices',
    modelType: 'LSTM',
    batchSize: 128,
    learningRate: 0.002,
    epochs: 200,
    validationFreq: 20,
    finalAccuracy: 92.1,
    bestLoss: 0.156,
    dockerImage: 'tensorflow/tensorflow:2.13.0',
    workDir: '/workspace/timeseries',
    command: 'python lstm_train.py --data stock_data.csv',
    envVars: {
      TF_CPP_MIN_LOG_LEVEL: '2',
    },
    description: '股票价格预测的 LSTM 模型训练',
  },
  {
    id: 'task-005',
    taskName: 'gpt-pretraining-mini',
    status: 'cancelled',
    creator: 'admin',
    framework: 'pytorch',
    startTime: '2024-06-20 10:00:00',
    endTime: '2024-06-20 12:30:00',
    duration: 9000, // 2.5小时
    resources: { cpu: 32, memory: 128, gpu: 8 },
    storage: 1000,
    priority: 'high',
    dataset: 'OpenWebText',
    modelType: 'GPT-2-Small',
    batchSize: 16,
    learningRate: 0.0001,
    epochs: 1000,
    validationFreq: 50,
    finalAccuracy: 0,
    bestLoss: 3.842,
    dockerImage: 'pytorch/pytorch:2.0.1-cuda11.7-cudnn8-devel',
    workDir: '/workspace/gpt',
    command: 'python pretrain_gpt.py --config gpt2_small.json',
    envVars: {
      CUDA_LAUNCH_BLOCKING: '1',
      NCCL_DEBUG: 'INFO',
    },
    description: 'GPT-2 Small 模型预训练（因资源不足被取消）',
  },
  {
    id: 'task-006',
    taskName: 'mobilenet-edge-optimization',
    status: 'timeout',
    creator: 'ml-engineer',
    framework: 'tensorflow',
    startTime: '2024-06-19 22:00:00',
    endTime: '2024-06-20 06:00:00',
    duration: 28800, // 8小时
    resources: { cpu: 4, memory: 16, gpu: 1 },
    storage: 100,
    priority: 'low',
    dataset: 'CIFAR-10',
    modelType: 'MobileNetV3',
    batchSize: 512,
    learningRate: 0.005,
    epochs: 500,
    validationFreq: 25,
    finalAccuracy: 89.3,
    bestLoss: 0.298,
    dockerImage: 'tensorflow/tensorflow:2.13.0-gpu',
    workDir: '/workspace/mobile',
    command: 'python train_mobilenet.py --optimize-for-edge',
    envVars: {
      TF_ENABLE_ONEDNN_OPTS: '1',
    },
    description: 'MobileNet 边缘设备优化训练（训练超时）',
  },
]);

const logs = ref<LogEntry[]>([
  {
    timestamp: '2024-06-23 08:00:15',
    level: 'info',
    message: 'Training job started with 4 GPUs',
  },
  {
    timestamp: '2024-06-23 08:00:16',
    level: 'info',
    message: 'Loading ImageNet dataset...',
  },
  {
    timestamp: '2024-06-23 08:02:30',
    level: 'info',
    message: 'Dataset loaded successfully, 1,281,167 training samples',
  },
  {
    timestamp: '2024-06-23 08:02:31',
    level: 'info',
    message: 'Initializing ResNet-50 model...',
  },
  {
    timestamp: '2024-06-23 08:02:35',
    level: 'info',
    message: 'Model initialized, total parameters: 25,557,032',
  },
  {
    timestamp: '2024-06-23 08:02:36',
    level: 'info',
    message: 'Starting training loop...',
  },
  {
    timestamp: '2024-06-23 08:15:42',
    level: 'info',
    message: 'Epoch 1/100 - Loss: 2.845, Accuracy: 45.2%',
  },
  {
    timestamp: '2024-06-23 08:28:15',
    level: 'warning',
    message: 'Learning rate adjusted to 0.0005',
  },
  {
    timestamp: '2024-06-23 08:30:20',
    level: 'info',
    message: 'Epoch 2/100 - Loss: 2.123, Accuracy: 52.8%',
  },
  {
    timestamp: '2024-06-23 09:45:30',
    level: 'error',
    message: 'GPU memory warning: 95% utilization',
  },
]);

// ===== 表单验证规则 =====
const cloneFormRules = {
  taskName: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 3, max: 50, message: '名称长度在 3 到 50 个字符', trigger: 'blur' },
  ],
};

// ===== 表格列配置 =====
const columns: TableColumnsType<TrainingTask> = [
  {
    title: '任务名称',
    key: 'taskName',
    width: 200,
    slots: { customRender: 'taskName' },
  },
  {
    title: '状态',
    key: 'status',
    width: 120,
    slots: { customRender: 'status' },
  },
  {
    title: '框架',
    key: 'framework',
    width: 120,
    slots: { customRender: 'framework' },
  },
  {
    title: '创建者',
    dataIndex: 'creator',
    key: 'creator',
    width: 100,
  },
  {
    title: '资源配置',
    key: 'resources',
    width: 180,
    slots: { customRender: 'resources' },
  },
  {
    title: '持续时间',
    key: 'duration',
    width: 120,
    slots: { customRender: 'duration' },
  },
  {
    title: '精度',
    dataIndex: 'finalAccuracy',
    key: 'finalAccuracy',
    width: 80,
    customRender: ({ text }) => `${text}%`,
  },
  {
    title: '开始时间',
    key: 'startTime',
    width: 150,
    slots: { customRender: 'startTime' },
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
  total: computed(() => filteredTasks.value.length),
  pageSize: 10,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
};

// ===== 计算属性 =====
const filteredTasks = computed(() => {
  let result = trainingTasks.value;

  if (filterStatus.value) {
    result = result.filter((item) => item.status === filterStatus.value);
  }

  if (filterFramework.value) {
    result = result.filter((item) => item.framework === filterFramework.value);
  }

  if (dateRange.value && dateRange.value.length === 2) {
    const [start, end] = dateRange.value;
    result = result.filter((item) => {
      const taskDate = new Date(item.startTime);
      return taskDate >= start.toDate() && taskDate <= end.toDate();
    });
  }

  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter(
      (item) =>
        item.taskName.toLowerCase().includes(keyword) ||
        item.creator.toLowerCase().includes(keyword),
    );
  }

  return result;
});

const totalTasks = computed(() => trainingTasks.value.length);
const completedTasks = computed(
  () =>
    trainingTasks.value.filter((task) => task.status === 'completed').length,
);
const failedTasks = computed(
  () => trainingTasks.value.filter((task) => task.status === 'failed').length,
);
const successRate = computed(() => {
  if (totalTasks.value === 0) return 0;
  return (completedTasks.value / totalTasks.value) * 100;
});

const filteredLogs = computed(() => {
  if (selectedLogLevel.value === 'all') {
    return logs.value;
  }
  return logs.value.filter((log) => log.level === selectedLogLevel.value);
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

const getFrameworkIcon = (framework: string) => {
  return (
    FRAMEWORK_ICONS[framework as keyof typeof FRAMEWORK_ICONS] ||
    ExperimentOutlined
  );
};

const getPriorityColor = (priority: string): string => {
  return PRIORITY_COLORS[priority as keyof typeof PRIORITY_COLORS] || 'default';
};

const formatDuration = (seconds: number): string => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  if (hours > 0) {
    return `${hours}h ${minutes}m ${secs}s`;
  } else if (minutes > 0) {
    return `${minutes}m ${secs}s`;
  } else {
    return `${secs}s`;
  }
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
const exportHistory = (): void => {
  message.success('训练历史记录导出成功');
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

const clearFilters = (): void => {
  filterStatus.value = '';
  filterFramework.value = '';
  dateRange.value = null;
  searchKeyword.value = '';
  message.success('筛选条件已清空');
};

const viewTaskDetails = (record: TrainingTask): void => {
  selectedTask.value = record;
  activeTab.value = 'basic';
  detailModalVisible.value = true;
};

const viewLogs = (record: TrainingTask): void => {
  selectedTask.value = record;
  selectedLogLevel.value = 'all';
  logModalVisible.value = true;
};

const handleMenuAction = (key: string, record: TrainingTask): void => {
  const actions = {
    clone: () => handleClone(record),
    template: () => handleSaveAsTemplate(record),
    download: () => handleDownloadModel(record),
    delete: () => handleDelete(record),
  };

  const action = actions[key as keyof typeof actions];
  if (action) {
    action();
  }
};

const handleClone = (record: TrainingTask): void => {
  cloneForm.taskName = `${record.taskName}-copy`;
  cloneForm.description = `复制自任务: ${record.taskName}`;
  selectedTask.value = record;
  cloneModalVisible.value = true;
};

const handleCloneSubmit = async (): Promise<void> => {
  try {
    await cloneFormRef.value?.validate();
    cloneLoading.value = true;

    await new Promise((resolve) => setTimeout(resolve, 1500));

    message.success('任务复制成功，已添加到训练队列');
    cloneModalVisible.value = false;
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

const handleSaveAsTemplate = (record: TrainingTask): void => {
  message.success(`任务 "${record.taskName}" 已保存为模板`);
};

const handleDownloadModel = (record: TrainingTask): void => {
  if (record.status === 'completed') {
    message.success(`开始下载模型: ${record.taskName}`);
  } else {
    message.warning('只能下载已完成的训练任务模型');
  }
};

const handleDelete = (record: TrainingTask): void => {
  const deleteConfirm = () => {
    const index = trainingTasks.value.findIndex(
      (item) => item.id === record.id,
    );
    if (index !== -1) {
      trainingTasks.value.splice(index, 1);
      message.success('训练记录删除成功');
    }
  };

  import('ant-design-vue').then(({ Modal }) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除训练任务 "${record.taskName}" 的记录吗？此操作不可恢复。`,
      okText: '确认',
      cancelText: '取消',
      type: 'warning',
      onOk: deleteConfirm,
    });
  });
};

const downloadLogs = (): void => {
  if (selectedTask.value) {
    message.success(`开始下载 "${selectedTask.value.taskName}" 的训练日志`);
  }
};

const refreshLogs = async (): Promise<void> => {
  const newLogs: LogEntry[] = [
    ...logs.value,
    {
      timestamp: new Date().toLocaleString(),
      level: 'info',
      message: 'Log refreshed',
    },
  ];
  logs.value = newLogs.slice(-50);
};

const filterLogs = (): void => {
  // 日志级别筛选逻辑在计算属性中处理
};

const handleFilterChange = (): void => {
  // 筛选变化时的处理逻辑
};

const handleDateRangeChange = (): void => {
  // 日期范围变化时的处理逻辑
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
.training-history-container {
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
.export-btn {
  border: none !important;
  height: 40px !important;
  padding: 0 24px !important;
  font-size: 14px !important;
  font-weight: 500 !important;
  border-radius: 6px !important;
  transition: all 0.3s ease !important;
}

.export-btn:hover {
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
.date-picker,
.refresh-btn,
.clear-btn {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

.refresh-btn:hover,
.clear-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

.action-col {
  display: flex;
  justify-content: flex-end;
}

/* ===== 统计卡片 ===== */
.stats-section {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 8px !important;
  text-align: center;
}

.stat-item :deep(.ant-statistic-title) {
  font-size: 14px !important;
  margin-bottom: 8px !important;
}

.stat-item :deep(.ant-statistic-content) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 8px !important;
}

.stat-icon {
  font-size: 20px;
}

.stat-icon-total {
  color: #1890ff;
}

.stat-icon-success {
  color: #52c41a;
}

.stat-icon-error {
  color: #ff4d4f;
}

.stat-icon-rate {
  color: #faad14;
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

/* ===== 任务名称 ===== */
.task-name-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.task-name-link {
  border: none !important;
  background: transparent !important;
  padding: 0 !important;
  height: auto !important;
  font-weight: 500 !important;
  transition: all 0.3s ease;
}

.task-name-link:hover {
  color: #1890ff !important;
}

.template-tag {
  font-size: 11px !important;
  padding: 2px 6px !important;
  border-radius: 4px !important;
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

.indicator-completed {
  background: #52c41a;
}

.indicator-failed {
  background: #ff4d4f;
}

.indicator-cancelled {
  background: #faad14;
}

.indicator-timeout {
  background: #8c8c8c;
}

/* ===== 框架显示 ===== */
.framework-wrapper {
  display: flex;
  align-items: center;
  gap: 6px;
}

.framework-icon {
  font-size: 16px;
  color: #1890ff;
}

.framework-text {
  font-weight: 500;
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

/* ===== 持续时间 ===== */
.duration-text {
  font-weight: 500;
  font-family: monospace;
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

/* ===== 详情页样式 ===== */
.detail-content {
  max-height: 600px;
  overflow-y: auto;
}

.detail-tabs :deep(.ant-tabs-tab) {
  font-weight: 500 !important;
}

.detail-descriptions :deep(.ant-descriptions-item-label) {
  font-weight: 600 !important;
}

.config-content,
.metrics-content,
.environment-content {
  padding: 16px 0;
}

.metric-card {
  text-align: center;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  margin-bottom: 16px;
}

.metric-title {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #8c8c8c;
}

.metric-value {
  font-size: 24px;
  font-weight: 600;
  color: #1890ff;
}

.chart-placeholder {
  text-align: center;
  padding: 40px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  margin-top: 16px;
}

.command-code {
  background: #f6f8fa;
  padding: 8px 12px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
  word-break: break-all;
}

.env-vars {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.env-var-item {
  display: flex;
  gap: 8px;
  padding: 4px 8px;
  background: #f6f8fa;
  border-radius: 4px;
  font-size: 12px;
  font-family: monospace;
}

.env-key {
  font-weight: 600;
  color: #1890ff;
}

.env-value {
  color: #262626;
}

/* ===== 日志容器 ===== */
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

.log-action-btn {
  border-radius: 4px !important;
  transition: all 0.3s ease;
}

.log-action-btn:hover {
  color: #1890ff !important;
  border-color: #1890ff !important;
}

.log-content {
  border: 1px solid #d9d9d9;
  border-radius: 6px !important;
  padding: 12px !important;
  max-height: 400px;
  overflow-y: auto;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace !important;
  background: #fafafa;
}

.log-line {
  margin: 0;
  font-size: 12px;
  line-height: 1.4;
  padding: 2px 0;
}

.log-info {
  color: #262626;
}

.log-warning {
  color: #faad14;
}

.log-error {
  color: #ff4d4f;
}

/* ===== 表单样式 ===== */
.form-input,
.form-textarea {
  border-radius: 6px !important;
  transition: all 0.3s ease;
}

/* ===== 响应式设计 ===== */
@media (max-width: 768px) {
  .training-history-container {
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

  .export-btn,
  .refresh-btn,
  .clear-btn {
    width: 100% !important;
    justify-content: center !important;
  }

  .action-col {
    margin-top: 12px;
    justify-content: stretch;
  }

  .action-col .ant-space {
    width: 100%;
    display: flex;
    gap: 8px;
  }

  .action-col .ant-space .ant-btn {
    flex: 1;
  }

  .sci-fi-modal :deep(.ant-modal) {
    margin: 16px !important;
    max-width: calc(100vw - 32px) !important;
  }

  .stats-section .ant-row .ant-col {
    margin-bottom: 16px;
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

  .metric-card {
    padding: 12px;
  }

  .metric-value {
    font-size: 20px;
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
