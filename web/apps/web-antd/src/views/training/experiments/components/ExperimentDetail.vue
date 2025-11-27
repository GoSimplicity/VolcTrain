<template>
  <div class="experiment-detail">
    <!-- 基本信息 -->
    <Descriptions :column="2" bordered>
      <Descriptions.Item label="实验名称">
        {{ experiment.name }}
      </Descriptions.Item>
      <Descriptions.Item label="状态">
        <Badge
          :status="getStatusBadgeType(experiment.status)"
          :text="getStatusText(experiment.status)"
        />
      </Descriptions.Item>
      <Descriptions.Item label="创建者">
        {{ experiment.createdBy }}
      </Descriptions.Item>
      <Descriptions.Item label="总运行次数">
        {{ experiment.totalRuns || 0 }}
      </Descriptions.Item>
      <Descriptions.Item label="最佳指标" v-if="experiment.bestMetric">
        {{ experiment.bestMetric.toFixed(4) }}
      </Descriptions.Item>
      <Descriptions.Item label="创建时间">
        {{ formatDateTime(experiment.createdAt) }}
      </Descriptions.Item>
      <Descriptions.Item label="实验描述" :span="2">
        {{ experiment.description || '暂无描述' }}
      </Descriptions.Item>
    </Descriptions>

    <!-- 运行历史 -->
    <Card title="运行历史" style="margin-top: 16px;">
      <template #extra>
        <Space>
          <Button type="primary" @click="$emit('run-experiment', experiment.id)">
            <PlayIcon />
            新建运行
          </Button>
          <Button @click="compareSelected">
            <GitCompareIcon />
            比较选中
          </Button>
        </Space>
      </template>
      
      <Table
        v-model:selectedRowKeys="selectedRunKeys"
        :columns="runsColumns"
        :data-source="experimentRuns"
        :pagination="false"
        :row-selection="{ type: 'checkbox' }"
        row-key="id"
        size="small"
      >
        <!-- 运行状态 -->
        <template #status="{ record }">
          <Badge
            :status="getRunStatusBadge(record.status)"
            :text="getRunStatusText(record.status)"
          />
        </template>

        <!-- 参数列 -->
        <template #parameters="{ record }">
          <div class="parameters">
            <Tag
              v-for="(value, key) in record.parameters"
              :key="key"
              size="small"
            >
              {{ key }}: {{ value }}
            </Tag>
          </div>
        </template>

        <!-- 指标列 -->
        <template #metrics="{ record }">
          <div class="metrics">
            <div
              v-for="(value, key) in record.metrics"
              :key="key"
              class="metric-item"
            >
              <span class="metric-name">{{ key }}:</span>
              <span class="metric-value">{{ typeof value === 'number' ? value.toFixed(4) : value }}</span>
            </div>
          </div>
        </template>

        <!-- 操作列 -->
        <template #action="{ record }">
          <Space>
            <Button size="small" @click="viewRunDetail(record.id)">
              <EyeIcon />
            </Button>
            <Button size="small" @click="downloadArtifacts(record.id)">
              <DownloadIcon />
            </Button>
          </Space>
        </template>
      </Table>
    </Card>

    <!-- 指标趋势图表 -->
    <Card title="指标趋势" style="margin-top: 16px;">
      <div ref="chartRef" class="chart-container" />
    </Card>

    <!-- 操作按钮 -->
    <div class="detail-actions" style="margin-top: 24px; text-align: center;">
      <Space>
        <Button @click="$emit('close')">关闭</Button>
        <Button type="primary" @click="$emit('run-experiment', experiment.id)">
          新建运行
        </Button>
      </Space>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import {
  Descriptions,
  Badge,
  Card,
  Table,
  Button,
  Space,
  Tag,
  message,
} from 'ant-design-vue';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import * as echarts from 'echarts';

// 图标
import {
  Play as PlayIcon,
  GitCompare as GitCompareIcon,
  Eye as EyeIcon,
  Download as DownloadIcon,
} from 'lucide-vue-next';

interface Props {
  experiment: any;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  close: [];
  'run-experiment': [id: string];
}>();

// 响应式数据
const selectedRunKeys = ref<string[]>([]);
const chartRef = ref<HTMLElement>();

// 模拟运行历史数据
const experimentRuns = ref([
  {
    id: 'run_1',
    status: 'completed',
    parameters: {
      'learning_rate': 0.001,
      'batch_size': 32,
      'epochs': 100,
    },
    metrics: {
      'accuracy': 0.9234,
      'loss': 0.1567,
      'f1_score': 0.9145,
    },
    startTime: '2024-01-20T10:00:00Z',
    endTime: '2024-01-20T12:30:00Z',
    duration: 9000,
  },
  {
    id: 'run_2',
    status: 'failed',
    parameters: {
      'learning_rate': 0.01,
      'batch_size': 64,
      'epochs': 100,
    },
    metrics: {
      'accuracy': 0.8756,
      'loss': 0.2345,
    },
    startTime: '2024-01-19T14:00:00Z',
    endTime: '2024-01-19T15:45:00Z',
    duration: 6300,
  },
]);

// 运行列表列配置
const runsColumns = [
  {
    title: '运行ID',
    dataIndex: 'id',
    key: 'id',
    width: 120,
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100,
    slots: { customRender: 'status' },
  },
  {
    title: '参数',
    key: 'parameters',
    width: 200,
    slots: { customRender: 'parameters' },
  },
  {
    title: '指标',
    key: 'metrics',
    width: 200,
    slots: { customRender: 'metrics' },
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    key: 'startTime',
    width: 150,
    customRender: ({ text }: { text: string }) => formatDateTime(text),
  },
  {
    title: '持续时间',
    dataIndex: 'duration',
    key: 'duration',
    width: 100,
    customRender: ({ text }: { text: number }) => formatDuration(text),
  },
  {
    title: '操作',
    key: 'action',
    width: 100,
    slots: { customRender: 'action' },
  },
];

// 状态相关方法
const getStatusBadgeType = (status: string) => {
  const statusMap: Record<string, string> = {
    active: 'processing',
    completed: 'success',
    archived: 'default',
  };
  return statusMap[status] || 'default';
};

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '活跃',
    completed: '已完成',
    archived: '已归档',
  };
  return statusMap[status] || status;
};

const getRunStatusBadge = (status: string) => {
  const statusMap: Record<string, string> = {
    running: 'processing',
    completed: 'success',
    failed: 'error',
    cancelled: 'warning',
  };
  return statusMap[status] || 'default';
};

const getRunStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    running: '运行中',
    completed: '已完成',
    failed: '已失败',
    cancelled: '已取消',
  };
  return statusMap[status] || status;
};

// 时间格式化
const formatDateTime = (dateTime: string) => {
  return format(new Date(dateTime), 'yyyy-MM-dd HH:mm', { locale: zhCN });
};

const formatDuration = (seconds: number) => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  return `${hours}h ${minutes}m`;
};

// 操作方法
const compareSelected = () => {
  if (selectedRunKeys.value.length < 2) {
    message.warning('请选择至少两个运行进行比较');
    return;
  }
  message.success(`比较 ${selectedRunKeys.value.length} 个运行`);
};

const viewRunDetail = (runId: string) => {
  message.info(`查看运行详情: ${runId}`);
};

const downloadArtifacts = (runId: string) => {
  message.success(`下载产物: ${runId}`);
};

// 初始化图表
const initChart = () => {
  if (!chartRef.value) return;

  const chart = echarts.init(chartRef.value);
  
  const option = {
    tooltip: {
      trigger: 'axis',
    },
    legend: {
      data: ['Accuracy', 'Loss'],
    },
    xAxis: {
      type: 'category',
      data: ['Run 1', 'Run 2', 'Run 3', 'Run 4', 'Run 5'],
    },
    yAxis: [
      {
        type: 'value',
        name: 'Accuracy',
        position: 'left',
        axisLabel: {
          formatter: '{value}',
        },
      },
      {
        type: 'value',
        name: 'Loss',
        position: 'right',
        axisLabel: {
          formatter: '{value}',
        },
      },
    ],
    series: [
      {
        name: 'Accuracy',
        type: 'line',
        yAxisIndex: 0,
        data: [0.85, 0.89, 0.92, 0.91, 0.93],
        smooth: true,
        lineStyle: { width: 2 },
      },
      {
        name: 'Loss',
        type: 'line',
        yAxisIndex: 1,
        data: [0.35, 0.28, 0.18, 0.20, 0.16],
        smooth: true,
        lineStyle: { width: 2 },
      },
    ],
  };

  chart.setOption(option);
};

onMounted(() => {
  initChart();
});
</script>

<style lang="scss" scoped>
.experiment-detail {
  .parameters {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }

  .metrics {
    .metric-item {
      display: flex;
      justify-content: space-between;
      margin-bottom: 2px;
      font-size: 12px;

      .metric-name {
        color: #666;
      }

      .metric-value {
        font-weight: 500;
      }
    }
  }

  .chart-container {
    height: 300px;
    width: 100%;
  }
}
</style>