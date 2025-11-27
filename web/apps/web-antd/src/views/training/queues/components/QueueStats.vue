<template>
  <div class="queue-stats">
    <!-- 统计概览 -->
    <Row :gutter="[16, 16]">
      <Col :span="6">
        <Card class="stat-card">
          <Statistic
            title="总任务数"
            :value="stats.totalJobs"
            :value-style="{ color: '#1890ff' }"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card class="stat-card">
          <Statistic
            title="运行中"
            :value="stats.runningJobs"
            :value-style="{ color: '#52c41a' }"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card class="stat-card">
          <Statistic
            title="已完成"
            :value="stats.completedJobs"
            :value-style="{ color: '#13c2c2' }"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card class="stat-card">
          <Statistic
            title="失败数"
            :value="stats.failedJobs"
            :value-style="{ color: '#ff4d4f' }"
          />
        </Card>
      </Col>
    </Row>

    <!-- 资源使用情况 -->
    <Card title="资源使用情况" style="margin-top: 16px;">
      <Row :gutter="[16, 16]">
        <Col :span="8">
          <div class="resource-usage">
            <div class="resource-title">CPU使用率</div>
            <Progress
              :percent="Math.round((stats.resourceUsage.cpu.used / stats.resourceUsage.cpu.total) * 100)"
              :stroke-color="{ from: '#108ee9', to: '#87d068' }"
            />
            <div class="resource-detail">
              {{ stats.resourceUsage.cpu.used }} / {{ stats.resourceUsage.cpu.total }} 核
            </div>
          </div>
        </Col>
        <Col :span="8">
          <div class="resource-usage">
            <div class="resource-title">内存使用率</div>
            <Progress
              :percent="Math.round((parseFloat(stats.resourceUsage.memory.used) / parseFloat(stats.resourceUsage.memory.total)) * 100)"
              :stroke-color="{ from: '#108ee9', to: '#87d068' }"
            />
            <div class="resource-detail">
              {{ stats.resourceUsage.memory.used }} / {{ stats.resourceUsage.memory.total }}
            </div>
          </div>
        </Col>
        <Col :span="8">
          <div class="resource-usage">
            <div class="resource-title">GPU使用率</div>
            <Progress
              :percent="Math.round((stats.resourceUsage.gpu.used / stats.resourceUsage.gpu.total) * 100)"
              :stroke-color="{ from: '#108ee9', to: '#87d068' }"
            />
            <div class="resource-detail">
              {{ stats.resourceUsage.gpu.used }} / {{ stats.resourceUsage.gpu.total }} 个
            </div>
          </div>
        </Col>
      </Row>
    </Card>

    <!-- 任务趋势图表 -->
    <Card title="任务执行趋势" style="margin-top: 16px;">
      <div ref="chartRef" class="chart-container" />
    </Card>

    <div class="stats-actions" style="margin-top: 24px; text-align: center;">
      <Button @click="$emit('close')">关闭</Button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import {
  Row,
  Col,
  Card,
  Statistic,
  Progress,
  Button,
} from 'ant-design-vue';
import * as echarts from 'echarts';

interface Props {
  queueId: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  close: [];
}>();

// 图表引用
const chartRef = ref<HTMLElement>();

// 统计数据
const stats = ref({
  totalJobs: 45,
  runningJobs: 8,
  completedJobs: 32,
  failedJobs: 5,
  resourceUsage: {
    cpu: { used: 35, total: 50 },
    memory: { used: '70Gi', total: '100Gi' },
    gpu: { used: 7, total: 10 },
  },
});

// 初始化图表
const initChart = () => {
  if (!chartRef.value) return;

  const chart = echarts.init(chartRef.value);
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
      },
    },
    legend: {
      data: ['运行中', '已完成', '已失败'],
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: [
      {
        type: 'category',
        boundaryGap: false,
        data: ['1/14', '1/15', '1/16', '1/17', '1/18', '1/19', '1/20'],
      },
    ],
    yAxis: [
      {
        type: 'value',
      },
    ],
    series: [
      {
        name: '运行中',
        type: 'line',
        smooth: true,
        lineStyle: { width: 2 },
        emphasis: { focus: 'series' },
        data: [5, 7, 8, 10, 9, 8, 8],
      },
      {
        name: '已完成',
        type: 'line',
        smooth: true,
        lineStyle: { width: 2 },
        emphasis: { focus: 'series' },
        data: [3, 5, 6, 8, 12, 15, 18],
      },
      {
        name: '已失败',
        type: 'line',
        smooth: true,
        lineStyle: { width: 2 },
        emphasis: { focus: 'series' },
        data: [1, 1, 2, 1, 2, 3, 2],
      },
    ],
  };

  chart.setOption(option);

  // 响应式图表
  window.addEventListener('resize', () => {
    chart.resize();
  });
};

onMounted(() => {
  initChart();
});
</script>

<style lang="scss" scoped>
.queue-stats {
  .stat-card {
    text-align: center;
    border-radius: 8px;
    border: none;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .resource-usage {
    text-align: center;

    .resource-title {
      margin-bottom: 12px;
      font-weight: 500;
      color: #333;
    }

    .resource-detail {
      margin-top: 8px;
      font-size: 12px;
      color: #666;
    }
  }

  .chart-container {
    height: 300px;
    width: 100%;
  }
}
</style>