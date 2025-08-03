<template>
  <Drawer
    v-model:open="visible"
    title="GPU实时监控"
    width="1000"
    placement="right"
    class="gpu-metrics-drawer"
  >
    <div v-if="gpu" class="drawer-content">
      <!-- 监控头部 -->
      <div class="metrics-header">
        <div class="gpu-info">
          <h3>{{ gpu.name }}</h3>
          <Tag :color="getGPUStatusColor(gpu.status)">
            {{ getGPUStatusLabel(gpu.status) }}
          </Tag>
          <span class="gpu-id">ID: {{ gpu.id }}</span>
        </div>
        
        <div class="controls">
          <Space>
            <Select
              v-model:value="timeRange"
              style="width: 120px"
              @change="handleTimeRangeChange"
            >
              <Select.Option value="5m">5分钟</Select.Option>
              <Select.Option value="15m">15分钟</Select.Option>
              <Select.Option value="1h">1小时</Select.Option>
              <Select.Option value="6h">6小时</Select.Option>
              <Select.Option value="1d">1天</Select.Option>
            </Select>
            
            <Button @click="refreshMetrics" :loading="loading">
              <ReloadOutlined />
              刷新
            </Button>
            
            <Button @click="toggleAutoRefresh" :type="autoRefresh ? 'primary' : 'default'">
              <ClockCircleOutlined />
              {{ autoRefresh ? '停止' : '自动' }}
            </Button>
          </Space>
        </div>
      </div>

      <Divider style="margin: 16px 0" />

      <!-- 实时状态卡片 -->
      <Row :gutter="16" class="realtime-stats">
        <Col :span="6">
          <Card>
            <Statistic
              title="GPU使用率"
              :value="gpu.gpuUtilization"
              suffix="%"
              :value-style="{ color: getUtilizationColor(gpu.gpuUtilization) }"
            />
            <Progress
              :percent="gpu.gpuUtilization"
              :stroke-color="getUtilizationColor(gpu.gpuUtilization)"
              size="small"
              style="margin-top: 8px"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card>
            <Statistic
              title="显存使用率"
              :value="gpu.memoryUsage"
              suffix="%"
              :value-style="{ color: '#722ed1' }"
            />
            <Progress
              :percent="gpu.memoryUsage"
              stroke-color="#722ed1"
              size="small"
              style="margin-top: 8px"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card>
            <Statistic
              title="温度"
              :value="gpu.temperature"
              suffix="°C"
              :value-style="{ color: getTemperatureColor(gpu.temperature) }"
            />
            <div class="temperature-indicator" :style="{ color: getTemperatureColor(gpu.temperature) }">
              {{ getTemperatureStatus(gpu.temperature) }}
            </div>
          </Card>
        </Col>
        <Col :span="6">
          <Card>
            <Statistic
              title="功耗"
              :value="gpu.powerUsage"
              suffix="W"
              :value-style="{ color: '#faad14' }"
            />
            <div class="power-ratio">
              {{ gpu.powerUsage }}W / {{ gpu.maxPower }}W
            </div>
          </Card>
        </Col>
      </Row>

      <!-- 图表区域 -->
      <div class="charts-container">
        <Row :gutter="16">
          <Col :span="12">
            <Card title="GPU使用率趋势" class="chart-card">
              <div ref="gpuUtilChart" class="chart" style="height: 300px"></div>
            </Card>
          </Col>
          <Col :span="12">
            <Card title="显存使用趋势" class="chart-card">
              <div ref="memoryChart" class="chart" style="height: 300px"></div>
            </Card>
          </Col>
        </Row>
        
        <Row :gutter="16" style="margin-top: 16px">
          <Col :span="12">
            <Card title="温度变化" class="chart-card">
              <div ref="temperatureChart" class="chart" style="height: 300px"></div>
            </Card>
          </Col>
          <Col :span="12">
            <Card title="功耗变化" class="chart-card">
              <div ref="powerChart" class="chart" style="height: 300px"></div>
            </Card>
          </Col>
        </Row>
      </div>

      <!-- 历史数据表格 -->
      <Card title="历史数据" class="data-table" style="margin-top: 16px">
        <Table
          :columns="dataColumns"
          :data-source="metricsData"
          :pagination="{ pageSize: 10, size: 'small' }"
          size="small"
          :scroll="{ y: 300 }"
        >
          <template #utilization="{ record }">
            <Progress
              :percent="record.gpuUtilization"
              size="small"
              :stroke-color="getUtilizationColor(record.gpuUtilization)"
              style="min-width: 100px"
            />
          </template>
          <template #memory="{ record }">
            <Progress
              :percent="record.memoryUsage"
              size="small"
              stroke-color="#722ed1"
              style="min-width: 100px"
            />
          </template>
          <template #temperature="{ record }">
            <Tag :color="getTemperatureColor(record.temperature)">
              {{ record.temperature }}°C
            </Tag>
          </template>
        </Table>
      </Card>
    </div>
  </Drawer>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue';
import {
  Drawer,
  Tag,
  Space,
  Select,
  Button,
  Divider,
  Row,
  Col,
  Card,
  Statistic,
  Progress,
  Table,
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons-vue';
import type { GPUDevice } from '#/api/types';
import { getGPUMetrics } from '#/api';
import * as echarts from 'echarts';

const props = defineProps<{
  visible: boolean;
  gpu: GPUDevice | null;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
}>();

// 响应式数据
const loading = ref(false);
const timeRange = ref('1h');
const autoRefresh = ref(false);
const metricsData = ref<any[]>([]);

// 图表引用
const gpuUtilChart = ref<HTMLElement>();
const memoryChart = ref<HTMLElement>();
const temperatureChart = ref<HTMLElement>();
const powerChart = ref<HTMLElement>();

// 图表实例
let gpuUtilChartInstance: echarts.ECharts | null = null;
let memoryChartInstance: echarts.ECharts | null = null;
let temperatureChartInstance: echarts.ECharts | null = null;
let powerChartInstance: echarts.ECharts | null = null;

// 自动刷新定时器
let autoRefreshTimer: NodeJS.Timeout | null = null;

// 表格列定义
const dataColumns = [
  {
    title: '时间',
    dataIndex: 'timestamp',
    key: 'timestamp',
    width: 150,
  },
  {
    title: 'GPU使用率',
    key: 'utilization',
    slots: { customRender: 'utilization' },
    width: 120,
  },
  {
    title: '显存使用率',
    key: 'memory',
    slots: { customRender: 'memory' },
    width: 120,
  },
  {
    title: '温度',
    key: 'temperature',
    slots: { customRender: 'temperature' },
    width: 80,
  },
  {
    title: '功耗(W)',
    dataIndex: 'powerUsage',
    key: 'powerUsage',
    width: 80,
  },
];

// 模拟监控数据
const generateMockData = () => {
  const data = [];
  const now = new Date();
  const timeRangeMinutes = {
    '5m': 5,
    '15m': 15,
    '1h': 60,
    '6h': 360,
    '1d': 1440,
  }[timeRange.value] || 60;
  
  const interval = Math.max(1, timeRangeMinutes / 100);
  
  for (let i = timeRangeMinutes; i >= 0; i -= interval) {
    const time = new Date(now.getTime() - i * 60 * 1000);
    data.push({
      timestamp: time.toLocaleString(),
      gpuUtilization: Math.random() * 100,
      memoryUsage: Math.random() * 100,
      temperature: 45 + Math.random() * 40,
      powerUsage: 100 + Math.random() * 200,
    });
  }
  
  return data;
};

// 工具方法
const getGPUStatusColor = (status: string) => {
  const colors = {
    available: 'success',
    allocated: 'processing',
    busy: 'warning',
    maintenance: 'default',
    offline: 'default',
    error: 'error',
  };
  return colors[status as keyof typeof colors] || 'default';
};

const getGPUStatusLabel = (status: string) => {
  const labels = {
    available: '可用',
    allocated: '已分配',
    busy: '忙碌',
    maintenance: '维护中',
    offline: '离线',
    error: '故障',
  };
  return labels[status as keyof typeof labels] || status;
};

const getUtilizationColor = (util: number) => {
  if (util >= 90) return '#f5222d';
  if (util >= 70) return '#fa8c16';
  if (util >= 40) return '#52c41a';
  return '#1890ff';
};

const getTemperatureColor = (temp: number) => {
  if (temp >= 85) return '#f5222d';
  if (temp >= 75) return '#fa8c16';
  if (temp >= 65) return '#faad14';
  return '#52c41a';
};

const getTemperatureStatus = (temp: number) => {
  if (temp >= 85) return '过热';
  if (temp >= 75) return '偏高';
  if (temp >= 65) return '正常偏高';
  return '正常';
};

// 图表初始化
const initCharts = async () => {
  await nextTick();
  
  if (gpuUtilChart.value) {
    gpuUtilChartInstance = echarts.init(gpuUtilChart.value);
  }
  if (memoryChart.value) {
    memoryChartInstance = echarts.init(memoryChart.value);
  }
  if (temperatureChart.value) {
    temperatureChartInstance = echarts.init(temperatureChart.value);
  }
  if (powerChart.value) {
    powerChartInstance = echarts.init(powerChart.value);
  }
  
  updateCharts();
};

const updateCharts = () => {
  const data = metricsData.value;
  const timestamps = data.map(item => item.timestamp);
  
  // GPU使用率图表
  if (gpuUtilChartInstance) {
    gpuUtilChartInstance.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: timestamps },
      yAxis: { type: 'value', max: 100, name: '使用率(%)' },
      series: [{
        data: data.map(item => item.gpuUtilization),
        type: 'line',
        smooth: true,
        lineStyle: { color: '#1890ff' },
        areaStyle: { color: 'rgba(24, 144, 255, 0.1)' },
      }],
      grid: { left: 50, right: 20, top: 20, bottom: 50 },
    });
  }
  
  // 显存使用率图表
  if (memoryChartInstance) {
    memoryChartInstance.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: timestamps },
      yAxis: { type: 'value', max: 100, name: '使用率(%)' },
      series: [{
        data: data.map(item => item.memoryUsage),
        type: 'line',
        smooth: true,
        lineStyle: { color: '#722ed1' },
        areaStyle: { color: 'rgba(114, 46, 209, 0.1)' },
      }],
      grid: { left: 50, right: 20, top: 20, bottom: 50 },
    });
  }
  
  // 温度图表
  if (temperatureChartInstance) {
    temperatureChartInstance.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: timestamps },
      yAxis: { type: 'value', name: '温度(°C)' },
      series: [{
        data: data.map(item => item.temperature),
        type: 'line',
        smooth: true,
        lineStyle: { color: '#faad14' },
        areaStyle: { color: 'rgba(250, 173, 20, 0.1)' },
      }],
      grid: { left: 50, right: 20, top: 20, bottom: 50 },
    });
  }
  
  // 功耗图表
  if (powerChartInstance) {
    powerChartInstance.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: timestamps },
      yAxis: { type: 'value', name: '功耗(W)' },
      series: [{
        data: data.map(item => item.powerUsage),
        type: 'line',
        smooth: true,
        lineStyle: { color: '#52c41a' },
        areaStyle: { color: 'rgba(82, 196, 26, 0.1)' },
      }],
      grid: { left: 50, right: 20, top: 20, bottom: 50 },
    });
  }
};

// 数据加载
const loadMetrics = async () => {
  if (!props.gpu) return;
  
  try {
    loading.value = true;
    // const response = await getGPUMetrics(props.gpu.id, {
    //   startTime: getStartTime(),
    //   endTime: new Date().toISOString(),
    //   interval: getInterval(),
    // });
    // metricsData.value = response.data;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    metricsData.value = generateMockData();
    
    updateCharts();
  } catch (error) {
    message.error('加载监控数据失败');
  } finally {
    loading.value = false;
  }
};

const refreshMetrics = () => {
  loadMetrics();
};

const handleTimeRangeChange = () => {
  loadMetrics();
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
    loadMetrics();
  }, 30000); // 30秒刷新一次
};

const stopAutoRefresh = () => {
  if (autoRefreshTimer) {
    clearInterval(autoRefreshTimer);
    autoRefreshTimer = null;
  }
};

// 监听抽屉打开
const handleDrawerOpen = () => {
  if (props.visible && props.gpu) {
    loadMetrics();
    setTimeout(() => {
      initCharts();
    }, 100);
  } else {
    stopAutoRefresh();
  }
};

// 生命周期
onMounted(() => {
  // 监听visible变化
  const unwatch = () => {};
  // 实际应该使用watch监听props.visible变化
  
  if (props.visible && props.gpu) {
    handleDrawerOpen();
  }
});

onUnmounted(() => {
  stopAutoRefresh();
  
  // 销毁图表实例
  if (gpuUtilChartInstance) {
    gpuUtilChartInstance.dispose();
  }
  if (memoryChartInstance) {
    memoryChartInstance.dispose();
  }
  if (temperatureChartInstance) {
    temperatureChartInstance.dispose();
  }
  if (powerChartInstance) {
    powerChartInstance.dispose();
  }
});

// 监听visible变化
import { watch } from 'vue';
watch(() => props.visible, (newVal) => {
  if (newVal) {
    handleDrawerOpen();
  } else {
    stopAutoRefresh();
  }
});
</script>

<style scoped lang="scss">
.gpu-metrics-drawer {
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
  overflow-y: auto;
}

.metrics-header {
  flex-shrink: 0;
  
  .gpu-info {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
      color: #1890ff;
    }
    
    .gpu-id {
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

.realtime-stats {
  margin-bottom: 24px;
  
  .temperature-indicator {
    font-size: 12px;
    font-weight: 500;
    margin-top: 4px;
  }
  
  .power-ratio {
    font-size: 12px;
    color: #999;
    margin-top: 4px;
  }
}

.charts-container {
  flex: 1;
  min-height: 0;
}

.chart-card {
  height: 100%;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
  
  :deep(.ant-card-head-title) {
    font-weight: 600;
    color: #1890ff;
    font-size: 14px;
  }
  
  :deep(.ant-card-body) {
    padding: 16px;
  }
}

.chart {
  width: 100%;
}

.data-table {
  flex-shrink: 0;
  margin-top: 16px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
  
  :deep(.ant-card-head-title) {
    font-weight: 600;
    color: #1890ff;
    font-size: 14px;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .drawer-content {
    padding: 16px;
  }
  
  .metrics-header {
    .gpu-info {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;
    }
    
    .controls {
      justify-content: flex-start;
      margin-top: 12px;
      
      :deep(.ant-space) {
        flex-wrap: wrap;
      }
    }
  }
  
  .realtime-stats {
    :deep(.ant-col) {
      margin-bottom: 12px;
    }
  }
  
  .charts-container {
    :deep(.ant-col) {
      margin-bottom: 16px;
    }
  }
  
  .chart {
    height: 250px !important;
  }
}
</style>