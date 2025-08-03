<template>
  <Drawer
    v-model:open="visible"
    title="训练任务日志"
    width="900"
    placement="right"
    class="job-logs-drawer"
  >
    <div v-if="job" class="drawer-content">
      <!-- 日志头部控制 -->
      <div class="logs-header">
        <div class="job-info">
          <h3>{{ job.name }}</h3>
          <Tag :color="getJobStatusColor(job.status)">
            {{ getJobStatusLabel(job.status) }}
          </Tag>
          <span class="job-id">ID: {{ job.id }}</span>
        </div>
        
        <div class="controls">
          <Space>
            <Select
              v-model:value="selectedPod"
              placeholder="选择Pod"
              style="width: 200px"
              @change="handlePodChange"
            >
              <Select.Option
                v-for="pod in availablePods"
                :key="pod.name"
                :value="pod.name"
              >
                <div class="pod-option">
                  <span>{{ pod.name }}</span>
                  <Tag :color="getPodStatusColor(pod.status)" size="small">
                    {{ pod.status }}
                  </Tag>
                </div>
              </Select.Option>
            </Select>
            
            <Select
              v-model:value="logLevel"
              placeholder="日志级别"
              style="width: 120px"
              @change="handleLogLevelChange"
            >
              <Select.Option value="">全部</Select.Option>
              <Select.Option value="ERROR">ERROR</Select.Option>
              <Select.Option value="WARN">WARN</Select.Option>
              <Select.Option value="INFO">INFO</Select.Option>
              <Select.Option value="DEBUG">DEBUG</Select.Option>
            </Select>
            
            <InputNumber
              v-model:value="tailLines"
              placeholder="行数"
              :min="100"
              :max="10000"
              :step="100"
              style="width: 100px"
              @change="handleTailLinesChange"
            />
            
            <Button @click="refreshLogs" :loading="refreshing">
              <ReloadOutlined />
              刷新
            </Button>
            
            <Button @click="toggleAutoRefresh" :type="autoRefresh ? 'primary' : 'default'">
              <ClockCircleOutlined />
              {{ autoRefresh ? '停止自动刷新' : '自动刷新' }}
            </Button>
            
            <Button @click="downloadLogs">
              <DownloadOutlined />
              下载
            </Button>
          </Space>
        </div>
      </div>

      <Divider style="margin: 16px 0" />

      <!-- 日志内容区域 -->
      <div class="logs-container">
        <div class="logs-toolbar">
          <div class="toolbar-left">
            <Space>
              <Button 
                size="small" 
                @click="clearLogs"
                :disabled="logs.length === 0"
              >
                <ClearOutlined />
                清空
              </Button>
              <Button 
                size="small" 
                @click="scrollToTop"
              >
                <UpOutlined />
                顶部
              </Button>
              <Button 
                size="small" 
                @click="scrollToBottom"
              >
                <DownOutlined />
                底部
              </Button>
              <Switch
                v-model:checked="followTail"
                size="small"
                checked-children="跟随"
                un-checked-children="固定"
              />
            </Space>
          </div>
          
          <div class="toolbar-right">
            <Space>
              <span class="log-stats">
                共 {{ logs.length }} 行
                <span v-if="filteredLogs.length !== logs.length">
                  (过滤后 {{ filteredLogs.length }} 行)
                </span>
              </span>
              <Input.Search
                v-model:value="searchKeyword"
                placeholder="搜索日志内容"
                style="width: 200px"
                size="small"
                allow-clear
                @search="handleSearch"
                @change="handleSearchChange"
              />
            </Space>
          </div>
        </div>

        <div
          ref="logsContent"
          class="logs-content"
          @scroll="handleScroll"
        >
          <div
            v-for="(log, index) in displayLogs"
            :key="index"
            :class="['log-line', getLogLineClass(log)]"
            @click="selectLogLine(index)"
          >
            <span class="log-timestamp">{{ extractTimestamp(log) }}</span>
            <span class="log-level">{{ extractLogLevel(log) }}</span>
            <span class="log-message">{{ extractMessage(log) }}</span>
          </div>
          
          <div v-if="loading" class="loading-indicator">
            <Spin size="small" />
            <span style="margin-left: 8px">加载更多日志...</span>
          </div>
          
          <div v-if="displayLogs.length === 0" class="empty-logs">
            <Empty
              description="暂无日志数据"
              :image="Empty.PRESENTED_IMAGE_SIMPLE"
            />
          </div>
        </div>
      </div>

      <!-- 日志详情模态框 -->
      <Modal
        v-model:open="logDetailVisible"
        title="日志详情"
        width="800px"
        :footer="null"
      >
        <div v-if="selectedLog" class="log-detail">
          <Descriptions :column="1" bordered>
            <Descriptions.Item label="时间戳">
              {{ extractTimestamp(selectedLog) }}
            </Descriptions.Item>
            <Descriptions.Item label="日志级别">
              <Tag :color="getLogLevelColor(extractLogLevel(selectedLog))">
                {{ extractLogLevel(selectedLog) }}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="Pod名称">
              {{ selectedPod }}
            </Descriptions.Item>
            <Descriptions.Item label="完整消息">
              <pre class="log-message-detail">{{ selectedLog }}</pre>
            </Descriptions.Item>
          </Descriptions>
        </div>
      </Modal>
    </div>
  </Drawer>
</template>

<script lang="ts" setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue';
import {
  Drawer,
  Tag,
  Space,
  Select,
  InputNumber,
  Button,
  Divider,
  Switch,
  Input,
  Spin,
  Empty,
  Modal,
  Descriptions,
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  ClockCircleOutlined,
  DownloadOutlined,
  ClearOutlined,
  UpOutlined,
  DownOutlined,
} from '@ant-design/icons-vue';
import type { TrainingJob } from '#/api/types';
import { getTrainingJobLogs } from '#/api';

const props = defineProps<{
  visible: boolean;
  job: TrainingJob | null;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
}>();

// 响应式数据
const loading = ref(false);
const refreshing = ref(false);
const selectedPod = ref('');
const logLevel = ref('');
const tailLines = ref(1000);
const autoRefresh = ref(false);
const followTail = ref(true);
const searchKeyword = ref('');
const logs = ref<string[]>([]);
const selectedLog = ref<string>('');
const logDetailVisible = ref(false);

// DOM引用
const logsContent = ref<HTMLElement>();

// 自动刷新定时器
let autoRefreshTimer: NodeJS.Timeout | null = null;

// 模拟可用的Pod数据
const availablePods = ref([
  { name: 'job-worker-0', status: 'Running' },
  { name: 'job-worker-1', status: 'Running' },
  { name: 'job-parameter-server', status: 'Running' },
]);

// 模拟日志数据
const mockLogs = [
  '2024-01-20 10:30:00 INFO: Starting training job initialization...',
  '2024-01-20 10:30:01 INFO: Loading configuration from /config/train.yaml',
  '2024-01-20 10:30:02 INFO: Initializing distributed training with 2 workers',
  '2024-01-20 10:30:03 INFO: Setting up GPU devices: [0, 1]',
  '2024-01-20 10:30:04 INFO: Loading training dataset from /data/train',
  '2024-01-20 10:30:05 INFO: Dataset loaded: 50000 samples',
  '2024-01-20 10:30:06 INFO: Loading validation dataset from /data/val',
  '2024-01-20 10:30:07 INFO: Validation dataset loaded: 10000 samples',
  '2024-01-20 10:30:08 INFO: Model architecture: ResNet50',
  '2024-01-20 10:30:09 INFO: Total parameters: 25,557,032',
  '2024-01-20 10:30:10 INFO: Trainable parameters: 25,557,032',
  '2024-01-20 10:30:11 INFO: Optimizer: Adam, lr=0.001',
  '2024-01-20 10:30:12 INFO: Loss function: CrossEntropyLoss',
  '2024-01-20 10:30:13 INFO: Starting training for 100 epochs',
  '2024-01-20 10:30:14 INFO: Epoch 1/100 started',
  '2024-01-20 10:31:00 INFO: Epoch 1/100 - Loss: 2.3045, Acc: 0.1234',
  '2024-01-20 10:31:30 INFO: Validation - Loss: 2.2891, Acc: 0.1456',
  '2024-01-20 10:31:31 INFO: Epoch 2/100 started',
  '2024-01-20 10:32:15 INFO: Epoch 2/100 - Loss: 2.1234, Acc: 0.2345',
  '2024-01-20 10:32:45 INFO: Validation - Loss: 2.0987, Acc: 0.2567',
  '2024-01-20 10:32:46 WARN: GPU memory usage high: 7.2GB/8GB',
  '2024-01-20 10:32:47 INFO: Epoch 3/100 started',
  '2024-01-20 10:33:30 ERROR: CUDA out of memory. Tried to allocate 256.00 MiB',
  '2024-01-20 10:33:31 INFO: Reducing batch size from 64 to 32',
  '2024-01-20 10:33:32 INFO: Resuming training with reduced batch size',
];

// 计算属性
const filteredLogs = computed(() => {
  let result = logs.value;
  
  // 按日志级别过滤
  if (logLevel.value) {
    result = result.filter(log => log.includes(logLevel.value));
  }
  
  // 按关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter(log => log.toLowerCase().includes(keyword));
  }
  
  return result;
});

const displayLogs = computed(() => {
  return filteredLogs.value.slice(-tailLines.value);
});

// 监听器
watch(() => props.visible, (newVal) => {
  if (newVal && props.job) {
    initializeLogs();
  } else {
    stopAutoRefresh();
  }
});

watch(() => props.job, (newJob) => {
  if (newJob && props.visible) {
    initializeLogs();
  }
});

// 工具方法
const getJobStatusColor = (status: string) => {
  const colors = {
    pending: 'default',
    queued: 'processing',
    running: 'success',
    completed: 'success',
    failed: 'error',
    cancelled: 'warning',
    paused: 'warning',
    stopped: 'default',
  };
  return colors[status as keyof typeof colors] || 'default';
};

const getJobStatusLabel = (status: string) => {
  const labels = {
    pending: '等待中',
    queued: '队列中',
    running: '运行中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消',
    paused: '已暂停',
    stopped: '已停止',
  };
  return labels[status as keyof typeof labels] || status;
};

const getPodStatusColor = (status: string) => {
  const colors = {
    Running: 'success',
    Pending: 'processing',
    Failed: 'error',
    Succeeded: 'success',
  };
  return colors[status as keyof typeof colors] || 'default';
};

const getLogLineClass = (log: string) => {
  if (log.includes('ERROR')) return 'log-error';
  if (log.includes('WARN')) return 'log-warning';
  if (log.includes('INFO')) return 'log-info';
  if (log.includes('DEBUG')) return 'log-debug';
  return '';
};

const getLogLevelColor = (level: string) => {
  const colors = {
    ERROR: 'red',
    WARN: 'orange',
    INFO: 'blue',
    DEBUG: 'default',
  };
  return colors[level as keyof typeof colors] || 'default';
};

const extractTimestamp = (log: string) => {
  const match = log.match(/^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})/);
  return match ? match[1] : '';
};

const extractLogLevel = (log: string) => {
  const match = log.match(/\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} (\w+):/);
  return match ? match[1] : '';
};

const extractMessage = (log: string) => {
  return log.replace(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \w+: /, '');
};

// 事件处理
const initializeLogs = () => {
  if (availablePods.value.length > 0 && !selectedPod.value) {
    selectedPod.value = availablePods.value[0].name;
  }
  loadLogs();
  if (autoRefresh.value) {
    startAutoRefresh();
  }
};

const loadLogs = async () => {
  if (!props.job || !selectedPod.value) return;
  
  try {
    loading.value = true;
    // 这里应该调用实际的API
    // const response = await getTrainingJobLogs(props.job.id, {
    //   lines: tailLines.value,
    //   pod: selectedPod.value
    // });
    // logs.value = response;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    logs.value = [...mockLogs];
    
    if (followTail.value) {
      await nextTick();
      scrollToBottom();
    }
  } catch (error) {
    message.error('加载日志失败');
  } finally {
    loading.value = false;
  }
};

const refreshLogs = async () => {
  refreshing.value = true;
  await loadLogs();
  refreshing.value = false;
  message.success('日志已刷新');
};

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value;
  if (autoRefresh.value) {
    startAutoRefresh();
    message.info('已开启自动刷新');
  } else {
    stopAutoRefresh();
    message.info('已停止自动刷新');
  }
};

const startAutoRefresh = () => {
  stopAutoRefresh();
  autoRefreshTimer = setInterval(() => {
    loadLogs();
  }, 5000);
};

const stopAutoRefresh = () => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer);
    autoRefreshTimer = null;
  }
};

const clearLogs = () => {
  logs.value = [];
  message.info('日志已清空');
};

const scrollToTop = () => {
  if (logsContent.value) {
    logsContent.value.scrollTop = 0;
  }
};

const scrollToBottom = () => {
  if (logsContent.value) {
    logsContent.value.scrollTop = logsContent.value.scrollHeight;
  }
};

const downloadLogs = () => {
  const logContent = logs.value.join('\n');
  const blob = new Blob([logContent], { type: 'text/plain' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${props.job?.name || 'training-job'}-${selectedPod.value}-logs.txt`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
  message.success('日志下载成功');
};

const selectLogLine = (index: number) => {
  selectedLog.value = displayLogs.value[index];
  logDetailVisible.value = true;
};

const handlePodChange = () => {
  loadLogs();
};

const handleLogLevelChange = () => {
  // 过滤逻辑已在计算属性中处理
};

const handleTailLinesChange = () => {
  // 显示逻辑已在计算属性中处理
};

const handleSearch = () => {
  // 搜索逻辑已在计算属性中处理
};

const handleSearchChange = () => {
  // 实时搜索逻辑已在计算属性中处理
};

const handleScroll = () => {
  if (!logsContent.value) return;
  
  const { scrollTop, scrollHeight, clientHeight } = logsContent.value;
  const isAtBottom = scrollTop + clientHeight >= scrollHeight - 10;
  
  // 如果用户滚动到底部，自动跟随
  if (isAtBottom && !followTail.value) {
    followTail.value = true;
  } else if (!isAtBottom && followTail.value) {
    followTail.value = false;
  }
};

// 生命周期
onMounted(() => {
  if (props.visible && props.job) {
    initializeLogs();
  }
});

onUnmounted(() => {
  stopAutoRefresh();
});
</script>

<style scoped lang="scss">
.job-logs-drawer {
  :deep(.ant-drawer-body) {
    padding: 0;
    display: flex;
    flex-direction: column;
  }
}

.drawer-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 24px;
}

.logs-header {
  flex-shrink: 0;
  
  .job-info {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
      color: #1890ff;
    }
    
    .job-id {
      font-size: 12px;
      color: #999;
      font-family: 'Monaco', 'Consolas', monospace;
    }
  }
  
  .controls {
    display: flex;
    justify-content: flex-end;
  }
}

.pod-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logs-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.logs-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: #fafafa;
  border-radius: 6px;
  
  .toolbar-left,
  .toolbar-right {
    display: flex;
    align-items: center;
  }
  
  .log-stats {
    font-size: 12px;
    color: #666;
    margin-right: 12px;
  }
}

.logs-content {
  flex: 1;
  background: #1e1e1e;
  border-radius: 6px;
  padding: 12px;
  overflow-y: auto;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.4;
  
  .log-line {
    color: #d4d4d4;
    margin-bottom: 2px;
    cursor: pointer;
    padding: 2px 4px;
    border-radius: 2px;
    transition: background-color 0.2s;
    
    &:hover {
      background-color: rgba(255, 255, 255, 0.1);
    }
    
    .log-timestamp {
      color: #569cd6;
      margin-right: 8px;
    }
    
    .log-level {
      margin-right: 8px;
      font-weight: bold;
    }
    
    .log-message {
      color: #d4d4d4;
    }
    
    &.log-error {
      .log-level {
        color: #f85149;
      }
      .log-message {
        color: #f85149;
      }
    }
    
    &.log-warning {
      .log-level {
        color: #d29922;
      }
      .log-message {
        color: #d29922;
      }
    }
    
    &.log-info {
      .log-level {
        color: #3794ff;
      }
    }
    
    &.log-debug {
      .log-level {
        color: #b5cea8;
      }
      .log-message {
        color: #999;
      }
    }
  }
  
  .loading-indicator {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px;
    color: #999;
  }
  
  .empty-logs {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 200px;
    color: #999;
  }
}

.log-detail {
  .log-message-detail {
    background: #f5f5f5;
    padding: 12px;
    border-radius: 4px;
    font-family: 'Monaco', 'Consolas', monospace;
    font-size: 12px;
    line-height: 1.4;
    max-height: 300px;
    overflow-y: auto;
    white-space: pre-wrap;
    word-break: break-all;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .drawer-content {
    padding: 16px;
  }
  
  .logs-header {
    .job-info {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;
    }
    
    .controls {
      :deep(.ant-space) {
        flex-wrap: wrap;
      }
    }
  }
  
  .logs-toolbar {
    flex-direction: column;
    gap: 12px;
    
    .toolbar-left,
    .toolbar-right {
      width: 100%;
      justify-content: center;
    }
  }
  
  .logs-content {
    font-size: 11px;
  }
}

// 滚动条样式
.logs-content::-webkit-scrollbar {
  width: 6px;
}

.logs-content::-webkit-scrollbar-track {
  background: #2d2d2d;
  border-radius: 3px;
}

.logs-content::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 3px;
}

.logs-content::-webkit-scrollbar-thumb:hover {
  background: #777;
}
</style>