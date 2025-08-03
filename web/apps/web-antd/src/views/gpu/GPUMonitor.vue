<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import {
  Card,
  Row,
  Col,
  Tabs,
  Select,
  Button,
  Space,
  Statistic,
  Divider,
} from 'ant-design-vue';

defineOptions({ name: 'GPUMonitor' });

// 模拟图表数据
const chartData = ref({
  usage: [65, 59, 80, 81, 56, 55, 72, 78, 81, 85, 90, 92],
  memory: [70, 63, 82, 79, 60, 55, 70, 75, 82, 87, 92, 95],
  temperature: [50, 55, 60, 65, 70, 68, 72, 75, 70, 72, 74, 75],
  power: [200, 210, 230, 240, 250, 245, 255, 260, 250, 255, 265, 270],
});

// 模拟GPU详情数据
const gpuDetails = ref({
  id: 'GPU-001',
  name: 'Tesla A100',
  architecture: 'Ampere',
  cudaCores: 6912,
  memory: '80GB HBM2',
  memoryBandwidth: '1.6 TB/s',
  baseClockRate: '1.41 GHz',
  boostClockRate: '1.73 GHz',
  maxPower: '300W',
  driver: 'NVIDIA 515.65.01',
  cuda: 'CUDA 11.7',
});

// 选择监控的GPU
const selectedGpu = ref('GPU-001');
const selectedMetric = ref('usage');
const timeRange = ref('1hour');

const gpuOptions = [
  { label: 'GPU-001 (NVIDIA A100)', value: 'GPU-001' },
  { label: 'GPU-002 (NVIDIA A100)', value: 'GPU-002' },
  { label: 'GPU-003 (NVIDIA A100)', value: 'GPU-003' },
  { label: 'GPU-004 (NVIDIA A100)', value: 'GPU-004' },
  { label: 'GPU-005 (NVIDIA V100)', value: 'GPU-005' },
];

const metricOptions = [
  { label: 'GPU使用率', value: 'usage' },
  { label: '显存使用率', value: 'memory' },
  { label: '温度', value: 'temperature' },
  { label: '功耗', value: 'power' },
];

const timeRangeOptions = [
  { label: '最近1小时', value: '1hour' },
  { label: '最近6小时', value: '6hours' },
  { label: '最近24小时', value: '24hours' },
  { label: '最近7天', value: '7days' },
];

// 当前显示的监控指标
const currentMetrics = ref({
  usage: 92,
  memory: 95,
  temperature: 75,
  power: 270,
});

const activeKey = ref('1');

// 初始化图表
onMounted(() => {
  // 使用实际数据初始化图表代码
  console.log('组件已加载，初始化图表');
});

const handleChangeGpu = (value: string | undefined) => {
  if (value) {
    selectedGpu.value = value;
    // 在实际应用中，这里应该获取选中GPU的监控数据
    console.log('切换到GPU:', value);
  }
};

const handleChangeMetric = (value: string | undefined) => {
  if (value) {
    selectedMetric.value = value;
    // 切换监控指标后重新渲染图表
    console.log('切换监控指标:', value);
  }
};

const handleChangeTimeRange = (value: string | undefined) => {
  if (value) {
    timeRange.value = value;
    // 切换时间范围后重新获取数据
    console.log('切换时间范围:', value);
  }
};

const refreshData = () => {
  // 实际应用中，这里应该重新获取最新的监控数据
  console.log('刷新数据');
};

// 获取当前选中指标的单位
const getMetricUnit = (metric: string) => {
  switch (metric) {
    case 'usage':
    case 'memory':
      return '%';
    case 'temperature':
      return '°C';
    case 'power':
      return 'W';
    default:
      return '';
  }
};

// 获取当前选中指标的名称
const getMetricName = (metric: string) => {
  switch (metric) {
    case 'usage':
      return 'GPU使用率';
    case 'memory':
      return '显存使用率';
    case 'temperature':
      return '温度';
    case 'power':
      return '功耗';
    default:
      return '';
  }
};
</script>

<template>
  <div class="gpu-monitor-container">
    <Card>
      <div class="monitor-header">
        <Space size="large">
          <div class="select-group">
            <span>GPU:</span>
            <Select
              v-model:value="selectedGpu"
              style="width: 200px"
              @change="handleChangeGpu"
            >
              <Select.Option
                v-for="option in gpuOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </Select.Option>
            </Select>
          </div>
          <div class="select-group">
            <span>监控指标:</span>
            <Select
              v-model:value="selectedMetric"
              style="width: 150px"
              @change="handleChangeMetric"
            >
              <Select.Option
                v-for="option in metricOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </Select.Option>
            </Select>
          </div>
          <div class="select-group">
            <span>时间范围:</span>
            <Select
              v-model:value="timeRange"
              style="width: 150px"
              @change="handleChangeTimeRange"
            >
              <Select.Option
                v-for="option in timeRangeOptions"
                :key="option.value"
                :value="option.value"
              >
                {{ option.label }}
              </Select.Option>
            </Select>
          </div>
          <Button @click="refreshData">刷新</Button>
        </Space>
      </div>

      <Tabs v-model:activeKey="activeKey">
        <Tabs.TabPane key="1" tab="实时监控">
          <div class="metrics-overview">
            <Row :gutter="16">
              <Col :span="6">
                <Card>
                  <Statistic
                    title="GPU使用率"
                    :value="currentMetrics.usage"
                    suffix="%"
                    :valueStyle="{ color: currentMetrics.usage > 90 ? '#cf1322' : '#3f8600' }"
                  />
                </Card>
              </Col>
              <Col :span="6">
                <Card>
                  <Statistic
                    title="显存使用率"
                    :value="currentMetrics.memory"
                    suffix="%"
                    :valueStyle="{ color: currentMetrics.memory > 90 ? '#cf1322' : '#3f8600' }"
                  />
                </Card>
              </Col>
              <Col :span="6">
                <Card>
                  <Statistic
                    title="温度"
                    :value="currentMetrics.temperature"
                    suffix="°C"
                    :valueStyle="{ color: currentMetrics.temperature > 80 ? '#cf1322' : '#3f8600' }"
                  />
                </Card>
              </Col>
              <Col :span="6">
                <Card>
                  <Statistic
                    title="功耗"
                    :value="currentMetrics.power"
                    suffix="W"
                    :valueStyle="{ color: currentMetrics.power > 280 ? '#cf1322' : '#3f8600' }"
                  />
                </Card>
              </Col>
            </Row>
          </div>

          <div class="chart-container" style="margin-top: 20px; height: 350px;">
            <Card title="GPU监控图表">
              <div class="chart-placeholder">
                <p>这里应该是一个动态图表，显示{{ getMetricName(selectedMetric) }}随时间变化的趋势</p>
                <p>时间范围: {{ timeRangeOptions.find(t => t.value === timeRange)?.label }}</p>
                <p>数据点: {{ chartData[selectedMetric as keyof typeof chartData].join(', ') }} {{ getMetricUnit(selectedMetric) }}</p>
              </div>
            </Card>
          </div>
        </Tabs.TabPane>
        <Tabs.TabPane key="2" tab="GPU详情">
          <div class="gpu-details">
            <Row :gutter="16">
              <Col :span="12">
                <Card title="基本信息">
                  <p><strong>GPU ID:</strong> {{ gpuDetails.id }}</p>
                  <p><strong>名称:</strong> {{ gpuDetails.name }}</p>
                  <p><strong>架构:</strong> {{ gpuDetails.architecture }}</p>
                  <p><strong>CUDA核心数:</strong> {{ gpuDetails.cudaCores }}</p>
                  <p><strong>显存容量:</strong> {{ gpuDetails.memory }}</p>
                  <p><strong>显存带宽:</strong> {{ gpuDetails.memoryBandwidth }}</p>
                </Card>
              </Col>
              <Col :span="12">
                <Card title="性能信息">
                  <p><strong>基础时钟频率:</strong> {{ gpuDetails.baseClockRate }}</p>
                  <p><strong>加速时钟频率:</strong> {{ gpuDetails.boostClockRate }}</p>
                  <p><strong>最大功耗:</strong> {{ gpuDetails.maxPower }}</p>
                  <p><strong>驱动版本:</strong> {{ gpuDetails.driver }}</p>
                  <p><strong>CUDA版本:</strong> {{ gpuDetails.cuda }}</p>
                </Card>
              </Col>
            </Row>

            <Divider />

            <Card title="历史性能图表" style="margin-top: 16px">
              <div class="chart-placeholder">
                <p>这里应该是一个历史性能图表，显示GPU过去一段时间的性能指标</p>
              </div>
            </Card>
          </div>
        </Tabs.TabPane>
      </Tabs>
    </Card>
  </div>
</template>

<style scoped>
.gpu-monitor-container {
  padding: 0;
}
.monitor-header {
  margin-bottom: 16px;
}
.select-group {
  display: flex;
  align-items: center;
}
.select-group span {
  margin-right: 8px;
}
.chart-placeholder {
  height: 300px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #fafafa;
  border: 1px dashed #d9d9d9;
  border-radius: 2px;
}
.metrics-overview {
  margin-bottom: 16px;
}
</style>
