<template>
  <div class="job-logs">
    <!-- 日志过滤器 -->
    <div class="logs-toolbar">
      <Row :gutter="16" align="middle">
        <Col :span="8">
          <Space>
            <Select 
              v-model:value="logLevel" 
              style="width: 120px;" 
              @change="handleLevelChange"
            >
              <Select.Option value="">全部级别</Select.Option>
              <Select.Option value="debug">Debug</Select.Option>
              <Select.Option value="info">Info</Select.Option>
              <Select.Option value="warn">Warn</Select.Option>
              <Select.Option value="error">Error</Select.Option>
            </Select>
            <InputNumber
              v-model:value="maxLines"
              :min="100"
              :max="10000"
              :step="100"
              style="width: 120px;"
              addon-after="行"
              @change="handleLinesChange"
            />
          </Space>
        </Col>
        <Col :span="8">
          <div class="time-range">
            <RangePicker
              v-model:value="timeRange"
              show-time
              format="YYYY-MM-DD HH:mm:ss"
              @change="handleTimeRangeChange"
            />
          </div>
        </Col>
        <Col :span="8" class="toolbar-actions">
          <Space>
            <Switch
              v-model:checked="autoRefresh"
              @change="handleAutoRefreshChange"
            />
            <span class="auto-refresh-label">自动刷新</span>
            <Button @click="refreshLogs">
              <RefreshCwIcon :class="{ 'animate-spin': loading }" />
              刷新
            </Button>
            <Button @click="downloadLogs">
              <DownloadIcon />
              下载
            </Button>
            <Button @click="clearLogs">
              <TrashIcon />
              清空
            </Button>
          </Space>
        </Col>
      </Row>
    </div>

    <!-- 日志状态统计 -->
    <div class="logs-stats">
      <Space size="large">
        <Statistic
          title="总日志数"
          :value="logStats.total"
          :value-style="{ color: '#3f8600' }"
        />
        <Statistic
          title="错误数"
          :value="logStats.error"
          :value-style="{ color: '#cf1322' }"
        />
        <Statistic
          title="警告数"
          :value="logStats.warn"
          :value-style="{ color: '#d48806' }"
        />
        <Statistic
          title="Info数"
          :value="logStats.info"
          :value-style="{ color: '#1890ff' }"
        />
      </Space>
    </div>

    <!-- 日志内容 -->
    <div class="logs-container">
      <div
        v-if="loading && logs.length === 0"
        class="logs-loading"
      >
        <Spin size="large" />
        <p>正在加载日志...</p>
      </div>
      
      <div v-else-if="logs.length === 0" class="logs-empty">
        <Empty description="暂无日志数据" />
      </div>
      
      <div v-else class="logs-content" ref="logsContentRef">
        <VirtualList
          :data-key="'timestamp'"
          :data-sources="filteredLogs"
          :data-component="LogItem"
          :estimate-size="32"
          :item-class-name="getLogItemClass"
          :keeps="100"
          :extra-props="{ searchKeyword: searchKeyword }"
        />
      </div>
    </div>

    <!-- 搜索浮动面板 -->
    <div class="logs-search" v-if="showSearch">
      <Input.Search
        v-model:value="searchKeyword"
        placeholder="搜索日志内容..."
        @search="handleSearch"
        @pressEnter="handleSearch"
        allowClear
      />
      <div class="search-results" v-if="searchKeyword">
        找到 {{ searchResults }} 条匹配结果
      </div>
    </div>

    <!-- 底部工具栏 -->
    <div class="logs-footer">
      <Row justify="space-between" align="middle">
        <Col>
          <Space>
            <Button 
              size="small"
              @click="showSearch = !showSearch"
            >
              <SearchIcon />
              {{ showSearch ? '隐藏搜索' : '搜索' }}
            </Button>
            <Button 
              size="small"
              @click="scrollToTop"
            >
              <ArrowUpIcon />
              回到顶部
            </Button>
            <Button 
              size="small"
              @click="scrollToBottom"
            >
              <ArrowDownIcon />
              滚动到底部
            </Button>
          </Space>
        </Col>
        <Col>
          <span class="logs-info">
            显示 {{ filteredLogs.length }} / {{ logs.length }} 条日志
          </span>
        </Col>
      </Row>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue';
import {
  Row,
  Col,
  Select,
  InputNumber,
  DatePicker,
  Switch,
  Button,
  Space,
  Statistic,
  Spin,
  Empty,
  Input,
  message,
} from 'ant-design-vue';
import { VirtualList } from '@tanstack/vue-virtual';
import { format } from 'date-fns';

// 图标
import {
  RefreshCw as RefreshCwIcon,
  Download as DownloadIcon,
  Trash as TrashIcon,
  Search as SearchIcon,
  ArrowUp as ArrowUpIcon,
  ArrowDown as ArrowDownIcon,
} from 'lucide-vue-next';

// 类型和服务
import type { TrainingApi } from '#/types/api';
import { trainingService } from '#/api/services';
import LogItem from './LogItem.vue';

const { RangePicker } = DatePicker;

interface Props {
  jobId: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  close: [];
}>();

// 响应式数据
const loading = ref(false);
const logs = ref<TrainingApi.TrainingLog[]>([]);
const logLevel = ref('');
const maxLines = ref(1000);
const timeRange = ref<[string, string] | null>(null);
const autoRefresh = ref(false);
const showSearch = ref(false);
const searchKeyword = ref('');

// DOM 引用
const logsContentRef = ref<HTMLElement>();

// 自动刷新定时器
let refreshTimer: NodeJS.Timeout | null = null;

// 日志统计
const logStats = computed(() => {
  const stats = {
    total: logs.value.length,
    error: 0,
    warn: 0,
    info: 0,
    debug: 0,
  };

  logs.value.forEach(log => {
    if (stats[log.level as keyof typeof stats] !== undefined) {
      stats[log.level as keyof typeof stats]++;
    }
  });

  return stats;
});

// 过滤后的日志
const filteredLogs = computed(() => {
  let filtered = logs.value;

  // 按级别过滤
  if (logLevel.value) {
    filtered = filtered.filter(log => log.level === logLevel.value);
  }

  // 按时间范围过滤
  if (timeRange.value && timeRange.value[0] && timeRange.value[1]) {
    const [startTime, endTime] = timeRange.value;
    filtered = filtered.filter(log => {
      const logTime = new Date(log.timestamp);
      return logTime >= new Date(startTime) && logTime <= new Date(endTime);
    });
  }

  // 按关键词过滤
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    filtered = filtered.filter(log =>
      log.message.toLowerCase().includes(keyword) ||
      log.source.toLowerCase().includes(keyword)
    );
  }

  // 限制行数
  return filtered.slice(-maxLines.value);
});

// 搜索结果数量
const searchResults = computed(() => {
  if (!searchKeyword.value) return 0;
  return filteredLogs.value.length;
});

// 获取日志项样式类
const getLogItemClass = (log: TrainingApi.TrainingLog) => {
  return `log-item log-level-${log.level}`;
};

// 加载日志数据
const loadLogs = async () => {
  try {
    loading.value = true;
    
    const params: any = {
      lines: maxLines.value,
    };

    if (logLevel.value) {
      params.level = logLevel.value;
    }

    if (timeRange.value && timeRange.value[0] && timeRange.value[1]) {
      params.startTime = timeRange.value[0];
      params.endTime = timeRange.value[1];
    }

    const newLogs = await trainingService.getJobLogs(props.jobId, params);
    logs.value = newLogs;

    // 自动滚动到底部（新日志）
    if (autoRefresh.value) {
      await nextTick();
      scrollToBottom();
    }
  } catch (error) {
    message.error('加载日志失败');
  } finally {
    loading.value = false;
  }
};

// 刷新日志
const refreshLogs = () => {
  loadLogs();
};

// 处理级别变化
const handleLevelChange = () => {
  loadLogs();
};

// 处理行数变化
const handleLinesChange = () => {
  loadLogs();
};

// 处理时间范围变化
const handleTimeRangeChange = () => {
  loadLogs();
};

// 处理自动刷新变化
const handleAutoRefreshChange = (checked: boolean) => {
  if (checked) {
    refreshTimer = setInterval(() => {
      loadLogs();
    }, 5000); // 每5秒刷新一次
  } else {
    if (refreshTimer) {
      clearInterval(refreshTimer);
      refreshTimer = null;
    }
  }
};

// 下载日志
const downloadLogs = () => {
  const logContent = filteredLogs.value.map(log => 
    `[${log.timestamp}] [${log.level.toUpperCase()}] [${log.source}] ${log.message}`
  ).join('\n');
  
  const blob = new Blob([logContent], { type: 'text/plain' });
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `job-${props.jobId}-logs-${format(new Date(), 'yyyyMMdd-HHmmss')}.log`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  window.URL.revokeObjectURL(url);
  
  message.success('日志下载成功');
};

// 清空日志显示
const clearLogs = () => {
  logs.value = [];
  message.success('日志已清空');
};

// 搜索处理
const handleSearch = () => {
  // 搜索逻辑已在 computed 中处理
  message.info(`找到 ${searchResults.value} 条匹配结果`);
};

// 滚动到顶部
const scrollToTop = () => {
  if (logsContentRef.value) {
    logsContentRef.value.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

// 滚动到底部
const scrollToBottom = () => {
  if (logsContentRef.value) {
    logsContentRef.value.scrollTo({
      top: logsContentRef.value.scrollHeight,
      behavior: 'smooth'
    });
  }
};

// 生命周期
onMounted(() => {
  loadLogs();
});

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer);
  }
});

// 监听自动刷新状态
watch(autoRefresh, handleAutoRefreshChange);
</script>

<style lang="scss" scoped>
.job-logs {
  height: 600px;
  display: flex;
  flex-direction: column;
  background: #1a1a1a;
  color: #ffffff;
  border-radius: 8px;
  overflow: hidden;

  .logs-toolbar {
    padding: 16px;
    background: #2a2a2a;
    border-bottom: 1px solid #404040;

    .toolbar-actions {
      text-align: right;

      .auto-refresh-label {
        color: #ffffff;
        font-size: 14px;
      }
    }
  }

  .logs-stats {
    padding: 12px 16px;
    background: #2a2a2a;
    border-bottom: 1px solid #404040;

    :deep(.ant-statistic-title) {
      color: #cccccc;
      font-size: 12px;
    }

    :deep(.ant-statistic-content) {
      font-size: 16px;
    }
  }

  .logs-container {
    flex: 1;
    position: relative;
    overflow: hidden;

    .logs-loading {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: #cccccc;

      :deep(.ant-spin-dot-item) {
        background-color: #1890ff;
      }
    }

    .logs-empty {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;

      :deep(.ant-empty-description) {
        color: #cccccc;
      }
    }

    .logs-content {
      height: 100%;
      overflow: auto;
      font-family: 'JetBrains Mono', 'Consolas', monospace;
      font-size: 13px;
      line-height: 1.5;

      :deep(.log-item) {
        padding: 4px 16px;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        word-wrap: break-word;

        &.log-level-error {
          background-color: rgba(255, 77, 79, 0.1);
          border-left: 3px solid #ff4d4f;
        }

        &.log-level-warn {
          background-color: rgba(250, 173, 20, 0.1);
          border-left: 3px solid #faad14;
        }

        &.log-level-info {
          background-color: rgba(24, 144, 255, 0.1);
          border-left: 3px solid #1890ff;
        }

        &.log-level-debug {
          background-color: rgba(82, 196, 26, 0.1);
          border-left: 3px solid #52c41a;
        }

        &:hover {
          background-color: rgba(255, 255, 255, 0.1);
        }
      }
    }
  }

  .logs-search {
    position: absolute;
    top: 80px;
    right: 16px;
    width: 300px;
    padding: 12px;
    background: rgba(42, 42, 42, 0.95);
    border: 1px solid #404040;
    border-radius: 6px;
    backdrop-filter: blur(10px);
    z-index: 10;

    .search-results {
      margin-top: 8px;
      font-size: 12px;
      color: #52c41a;
    }
  }

  .logs-footer {
    padding: 8px 16px;
    background: #2a2a2a;
    border-top: 1px solid #404040;

    .logs-info {
      font-size: 12px;
      color: #cccccc;
    }
  }

  // 覆盖 Ant Design 组件样式
  :deep(.ant-select),
  :deep(.ant-input-number),
  :deep(.ant-picker) {
    background-color: #3a3a3a !important;
    border-color: #505050 !important;
    color: #ffffff !important;

    .ant-select-selector,
    .ant-input-number-input,
    .ant-picker-input input {
      background-color: transparent !important;
      color: #ffffff !important;
    }

    &:hover {
      border-color: #1890ff !important;
    }

    &.ant-select-focused,
    &.ant-input-number-focused,
    &.ant-picker-focused {
      border-color: #1890ff !important;
      box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2) !important;
    }
  }

  :deep(.ant-btn) {
    &:not(.ant-btn-primary) {
      background-color: #3a3a3a;
      border-color: #505050;
      color: #ffffff;

      &:hover {
        background-color: #4a4a4a;
        border-color: #1890ff;
        color: #1890ff;
      }
    }
  }

  :deep(.ant-switch) {
    background-color: rgba(0, 0, 0, 0.25);

    &.ant-switch-checked {
      background-color: #1890ff;
    }
  }
}

// 动画效果
@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>