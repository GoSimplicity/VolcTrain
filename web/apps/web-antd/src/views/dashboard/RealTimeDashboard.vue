<template>
  <div class="real-time-dashboard">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <DashboardOutlined class="title-icon" />
            <span class="title-text">实时仪表板</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">实时监控系统状态和性能指标</span>
          </p>
        </div>
        <div class="action-section">
          <Space>
            <Badge :count="activeAlerts" :offset="[10, 0]">
              <Button @click="showAlertsPanel">
                <BellOutlined />
                告警中心
              </Button>
            </Badge>
            <Dropdown>
              <Button>
                <SettingOutlined />
                刷新设置
              </Button>
              <template #overlay>
                <Menu @click="handleRefreshSetting">
                  <Menu.Item key="5">5秒</Menu.Item>
                  <Menu.Item key="10">10秒</Menu.Item>
                  <Menu.Item key="30">30秒</Menu.Item>
                  <Menu.Item key="60">1分钟</Menu.Item>
                  <Menu.Item key="stop">停止自动刷新</Menu.Item>
                </Menu>
              </template>
            </Dropdown>
            <Button type="primary" @click="refreshAllData">
              <ReloadOutlined :spin="globalLoading" />
              刷新数据
            </Button>
          </Space>
        </div>
      </div>
    </div>

    <!-- 系统状态概览 -->
    <div class="status-overview">
      <Row :gutter="16">
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="status-card">
            <div class="status-item">
              <div class="status-icon" :class="`status-${systemHealth.overall}`">
                <CheckCircleOutlined v-if="systemHealth.overall === 'healthy'" />
                <ExclamationCircleOutlined v-else-if="systemHealth.overall === 'warning'" />
                <CloseCircleOutlined v-else />
              </div>
              <div class="status-info">
                <div class="status-title">系统状态</div>
                <div class="status-value">{{ getHealthText(systemHealth.overall) }}</div>
                <div class="status-detail">运行时间: {{ formatUptime(systemInfo.uptime) }}</div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="status-card">
            <div class="status-item">
              <div class="status-chart">
                <div ref="cpuGaugeRef" class="gauge-chart"></div>
              </div>
              <div class="status-info">
                <div class="status-title">CPU使用率</div>
                <div class="status-value">{{ resourceUsage.cpu.usage }}%</div>
                <div class="status-detail">{{ resourceUsage.cpu.cores }}核心</div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="status-card">
            <div class="status-item">
              <div class="status-chart">
                <div ref="memoryGaugeRef" class="gauge-chart"></div>
              </div>
              <div class="status-info">
                <div class="status-title">内存使用率</div>
                <div class="status-value">{{ resourceUsage.memory.usage }}%</div>
                <div class="status-detail">{{ formatFileSize(resourceUsage.memory.used) }} / {{ formatFileSize(resourceUsage.memory.total) }}</div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="status-card">
            <div class="status-item">
              <div class="status-chart">
                <div ref="diskGaugeRef" class="gauge-chart"></div>
              </div>
              <div class="status-info">
                <div class="status-title">磁盘使用率</div>
                <div class="status-value">{{ resourceUsage.disk.usage }}%</div>
                <div class="status-detail">{{ formatFileSize(resourceUsage.disk.free) }} 可用</div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 性能趋势图表 -->
    <div class="performance-section">
      <Row :gutter="16">
        <Col :xs="24" :lg="16">
          <Card title="性能趋势" class="performance-card">
            <template #extra>
              <Space>
                <Select
                  v-model:value="timeRange"
                  style="width: 120px"
                  size="small"
                  @change="handleTimeRangeChange"
                >
                  <Select.Option value="1h">最近1小时</Select.Option>
                  <Select.Option value="6h">最近6小时</Select.Option>
                  <Select.Option value="1d">最近1天</Select.Option>
                  <Select.Option value="1w">最近1周</Select.Option>
                </Select>
                <Button size="small" @click="refreshPerformanceData">
                  <ReloadOutlined :spin="performanceLoading" />
                </Button>
              </Space>
            </template>
            <div ref="performanceTrendRef" class="performance-chart" v-loading="performanceLoading"></div>
          </Card>
        </Col>
        <Col :xs="24" :lg="8">
          <Card title="资源分布" class="resource-distribution-card">
            <div ref="resourcePieRef" class="pie-chart"></div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- GPU监控和任务状态 -->
    <div class="gpu-job-section">
      <Row :gutter="16">
        <Col :xs="24" :lg="12">
          <Card title="GPU监控" class="gpu-card">
            <template #extra>
              <Tag :color="gpuStatus.overall === 'healthy' ? 'green' : 'red'">
                {{ gpuStatus.totalGPUs }}个GPU {{ gpuStatus.overall === 'healthy' ? '正常' : '异常' }}
              </Tag>
            </template>
            <div class="gpu-grid">
              <div
                v-for="(gpu, index) in gpuStatus.gpus"
                :key="index"
                class="gpu-item"
                :class="`gpu-${getGPUStatusLevel(gpu.utilization)}`"
              >
                <div class="gpu-header">
                  <span class="gpu-name">GPU {{ index }}</span>
                  <span class="gpu-utilization">{{ gpu.utilization }}%</span>
                </div>
                <Progress
                  :percent="gpu.utilization"
                  :stroke-color="getProgressColor(gpu.utilization)"
                  size="small"
                  :show-info="false"
                />
                <div class="gpu-details">
                  <div class="gpu-detail-item">
                    <span class="detail-label">温度:</span>
                    <span class="detail-value" :style="{ color: getTemperatureColor(gpu.temperature) }">
                      {{ gpu.temperature }}°C
                    </span>
                  </div>
                  <div class="gpu-detail-item">
                    <span class="detail-label">显存:</span>
                    <span class="detail-value">{{ Math.round(gpu.memoryUsed / gpu.memoryTotal * 100) }}%</span>
                  </div>
                  <div class="gpu-detail-item">
                    <span class="detail-label">功耗:</span>
                    <span class="detail-value">{{ gpu.powerUsage }}W</span>
                  </div>
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :lg="12">
          <Card title="任务状态" class="job-status-card">
            <template #extra>
              <Space>
                <Statistic
                  title="运行中"
                  :value="jobStats.running"
                  :value-style="{ color: '#52c41a', fontSize: '16px' }"
                />
                <Divider type="vertical" />
                <Statistic
                  title="排队中"
                  :value="jobStats.pending"
                  :value-style="{ color: '#faad14', fontSize: '16px' }"
                />
              </Space>
            </template>
            <div ref="jobStatusPieRef" class="pie-chart"></div>
            <div class="job-queue">
              <h4>任务队列</h4>
              <List
                :data-source="recentJobs"
                size="small"
                :pagination="false"
              >
                <template #renderItem="{ item }">
                  <List.Item>
                    <List.Item.Meta>
                      <template #title>
                        <Space>
                          {{ item.name }}
                          <Tag :color="getJobStatusColor(item.status)">
                            {{ getJobStatusText(item.status) }}
                          </Tag>
                        </Space>
                      </template>
                      <template #description>
                        <Space size="small" style="font-size: 12px">
                          <span>GPU: {{ item.gpuCount }}</span>
                          <span>用时: {{ formatDuration(item.duration) }}</span>
                          <span>{{ formatRelativeTime(item.startTime) }}</span>
                        </Space>
                      </template>
                    </List.Item.Meta>
                    <div class="job-progress">
                      <Progress
                        v-if="item.status === 'running'"
                        :percent="item.progress"
                        size="small"
                        :stroke-color="getProgressColor(item.progress)"
                      />
                    </div>
                  </List.Item>
                </template>
              </List>
            </div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 告警和日志 -->
    <div class="alerts-logs-section">
      <Row :gutter="16">
        <Col :xs="24" :lg="12">
          <Card title="最新告警" class="alerts-card">
            <template #extra>
              <Space>
                <Badge :count="criticalAlertsCount">
                  <Tag color="red">严重</Tag>
                </Badge>
                <Badge :count="warningAlertsCount">
                  <Tag color="orange">警告</Tag>
                </Badge>
                <Button size="small" @click="showAlertsPanel">
                  查看全部
                </Button>
              </Space>
            </template>
            <List
              :data-source="alerts.slice(0, 5)"
              size="small"
              :pagination="false"
            >
              <template #renderItem="{ item }">
                <List.Item>
                  <List.Item.Meta>
                    <template #title>
                      <Space>
                        <Tag :color="getAlertLevelColor(item.level)">
                          {{ getAlertLevelText(item.level) }}
                        </Tag>
                        {{ item.summary }}
                      </Space>
                    </template>
                    <template #description>
                      <div class="alert-description">{{ item.description }}</div>
                      <div class="alert-time">{{ formatRelativeTime(item.startsAt) }}</div>
                    </template>
                  </List.Item.Meta>
                  <div class="alert-actions">
                    <Space size="small">
                      <Button type="text" size="small" @click="handleAcknowledgeAlert(item.id)">
                        确认
                      </Button>
                      <Button type="text" size="small" @click="viewAlertDetail(item)">
                        详情
                      </Button>
                    </Space>
                  </div>
                </List.Item>
              </template>
            </List>
            <div v-if="alerts.length === 0" class="empty-state">
              <Empty description="暂无告警" />
            </div>
          </Card>
        </Col>
        <Col :xs="24" :lg="12">
          <Card title="系统活动" class="activity-card">
            <template #extra>
              <Select
                v-model:value="activityFilter"
                style="width: 100px"
                size="small"
                @change="filterActivity"
              >
                <Select.Option value="all">全部</Select.Option>
                <Select.Option value="job">任务</Select.Option>
                <Select.Option value="user">用户</Select.Option>
                <Select.Option value="system">系统</Select.Option>
              </Select>
            </template>
            <Timeline class="activity-timeline">
              <Timeline.Item
                v-for="activity in filteredActivities"
                :key="activity.id"
                :color="getActivityColor(activity.type)"
              >
                <div class="activity-item">
                  <div class="activity-header">
                    <span class="activity-type">{{ getActivityTypeText(activity.type) }}</span>
                    <span class="activity-time">{{ formatRelativeTime(activity.timestamp) }}</span>
                  </div>
                  <div class="activity-content">{{ activity.message }}</div>
                  <div v-if="activity.user" class="activity-user">{{ activity.user }}</div>
                </div>
              </Timeline.Item>
            </Timeline>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 网络和存储监控 -->
    <div class="network-storage-section">
      <Row :gutter="16">
        <Col :xs="24" :lg="12">
          <Card title="网络流量" class="network-card">
            <div ref="networkChartRef" class="network-chart"></div>
          </Card>
        </Col>
        <Col :xs="24" :lg="12">
          <Card title="存储监控" class="storage-card">
            <div class="storage-stats">
              <div class="storage-item">
                <div class="storage-header">
                  <span class="storage-title">系统存储</span>
                  <span class="storage-usage">{{ Math.round(storageStats.system.used / storageStats.system.total * 100) }}%</span>
                </div>
                <Progress
                  :percent="Math.round(storageStats.system.used / storageStats.system.total * 100)"
                  :stroke-color="getProgressColor(Math.round(storageStats.system.used / storageStats.system.total * 100))"
                />
                <div class="storage-detail">
                  {{ formatFileSize(storageStats.system.used) }} / {{ formatFileSize(storageStats.system.total) }}
                </div>
              </div>
              <div class="storage-item">
                <div class="storage-header">
                  <span class="storage-title">用户数据</span>
                  <span class="storage-usage">{{ Math.round(storageStats.userData.used / storageStats.userData.total * 100) }}%</span>
                </div>
                <Progress
                  :percent="Math.round(storageStats.userData.used / storageStats.userData.total * 100)"
                  :stroke-color="getProgressColor(Math.round(storageStats.userData.used / storageStats.userData.total * 100))"
                />
                <div class="storage-detail">
                  {{ formatFileSize(storageStats.userData.used) }} / {{ formatFileSize(storageStats.userData.total) }}
                </div>
              </div>
              <div class="storage-item">
                <div class="storage-header">
                  <span class="storage-title">模型存储</span>
                  <span class="storage-usage">{{ Math.round(storageStats.models.used / storageStats.models.total * 100) }}%</span>
                </div>
                <Progress
                  :percent="Math.round(storageStats.models.used / storageStats.models.total * 100)"
                  :stroke-color="getProgressColor(Math.round(storageStats.models.used / storageStats.models.total * 100))"
                />
                <div class="storage-detail">
                  {{ formatFileSize(storageStats.models.used) }} / {{ formatFileSize(storageStats.models.total) }}
                </div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { message } from 'ant-design-vue';
import {
  Card,
  Row,
  Col,
  Space,
  Button,
  Badge,
  Dropdown,
  Menu,
  Select,
  Tag,
  Progress,
  Statistic,
  Divider,
  List,
  Empty,
  Timeline,
} from 'ant-design-vue';
import {
  DashboardOutlined,
  BellOutlined,
  SettingOutlined,
  ReloadOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  CloseCircleOutlined,
} from '@ant-design/icons-vue';

import {
  useRealTimeData,
  useGaugeChart,
  useTimeSeriesChart,
  usePieChart,
  colorUtils,
  chartPresets
} from '#/utils/charts';
import {
  getResourceUsage,
  getSystemInfo,
  getAlerts
} from '#/api';
import { formatFileSize } from '#/utils/date';

defineOptions({ name: 'RealTimeDashboard' });

// 响应式数据
const globalLoading = ref(false);
const performanceLoading = ref(false);
const timeRange = ref('1h');
const activityFilter = ref('all');
const refreshInterval = ref(10); // 秒

// 图表引用
const cpuGaugeRef = ref();
const memoryGaugeRef = ref();
const diskGaugeRef = ref();
const performanceTrendRef = ref();
const resourcePieRef = ref();
const jobStatusPieRef = ref();
const networkChartRef = ref();

// 图表实例
const cpuGauge = useGaugeChart(cpuGaugeRef, chartPresets.gpuGauge);
const memoryGauge = useGaugeChart(memoryGaugeRef, chartPresets.gpuGauge);
const diskGauge = useGaugeChart(diskGaugeRef, chartPresets.gpuGauge);
const performanceTrend = useTimeSeriesChart(performanceTrendRef, {
  title: '',
  yAxisName: '使用率',
  unit: '%',
  smooth: true
});
const resourcePie = usePieChart(resourcePieRef, {
  title: '',
  radius: ['40%', '70%']
});
const jobStatusPie = usePieChart(jobStatusPieRef, {
  title: '',
  radius: '70%'
});
const networkChart = useTimeSeriesChart(networkChartRef, {
  title: '',
  yAxisName: '流量',
  unit: 'MB/s',
  smooth: true
});

// 实时数据
const { data: resourceUsage, refresh: refreshResourceUsage } = useRealTimeData(
  getResourceUsage,
  refreshInterval.value * 1000
);

const { data: systemInfo, refresh: refreshSystemInfo } = useRealTimeData(
  getSystemInfo,
  60000 // 系统信息1分钟刷新一次
);

const { data: alerts, refresh: refreshAlerts } = useRealTimeData(
  () => getAlerts({ pageSize: 20 }),
  refreshInterval.value * 1000
);

// 模拟数据
const systemHealth = ref({
  overall: 'healthy' as 'healthy' | 'warning' | 'critical'
});

const gpuStatus = ref({
  overall: 'healthy' as 'healthy' | 'warning' | 'critical',
  totalGPUs: 4,
  gpus: [
    { utilization: 75, temperature: 65, memoryUsed: 6000, memoryTotal: 8000, powerUsage: 220 },
    { utilization: 82, temperature: 70, memoryUsed: 7200, memoryTotal: 8000, powerUsage: 245 },
    { utilization: 45, temperature: 58, memoryUsed: 3500, memoryTotal: 8000, powerUsage: 180 },
    { utilization: 0, temperature: 35, memoryUsed: 0, memoryTotal: 8000, powerUsage: 25 },
  ]
});

const jobStats = ref({
  running: 8,
  pending: 15,
  completed: 1234,
  failed: 12
});

const recentJobs = ref([
  {
    id: 'job-001',
    name: 'BERT训练任务',
    status: 'running',
    progress: 75,
    gpuCount: 2,
    duration: 3600,
    startTime: '2024-01-20 14:30:00'
  },
  {
    id: 'job-002',
    name: 'ResNet-50训练',
    status: 'pending',
    progress: 0,
    gpuCount: 1,
    duration: 0,
    startTime: '2024-01-20 15:00:00'
  },
  // 添加更多任务...
]);

const activities = ref([
  {
    id: 'act-001',
    type: 'job',
    message: '用户张三提交了新的训练任务',
    user: '张三',
    timestamp: '2024-01-20 15:30:00'
  },
  {
    id: 'act-002',
    type: 'system',
    message: 'GPU节点 node-001 重新上线',
    timestamp: '2024-01-20 15:25:00'
  },
  // 添加更多活动...
]);

const storageStats = ref({
  system: { used: 500 * 1024 * 1024 * 1024, total: 1000 * 1024 * 1024 * 1024 },
  userData: { used: 750 * 1024 * 1024 * 1024, total: 2000 * 1024 * 1024 * 1024 },
  models: { used: 300 * 1024 * 1024 * 1024, total: 500 * 1024 * 1024 * 1024 }
});

// 计算属性
const activeAlerts = computed(() => {
  return alerts.value?.items?.filter((alert: { status: string }) => alert.status === 'firing').length || 0;
});

const criticalAlertsCount = computed(() => {
  return alerts.value?.items?.filter((alert: { level: string }) => alert.level === 'critical').length || 0;
});

const warningAlertsCount = computed(() => {
  return alerts.value?.items?.filter((alert: { level: string }) => alert.level === 'warning').length || 0;
});

const filteredActivities = computed(() => {
  if (activityFilter.value === 'all') {
    return activities.value.slice(0, 10);
  }
  return activities.value.filter(activity => activity.type === activityFilter.value).slice(0, 10);
});

// 工具方法
const getHealthText = (status: string) => {
  const texts = {
    healthy: '健康',
    warning: '警告',
    critical: '严重'
  };
  return texts[status as keyof typeof texts] || status;
};

const formatUptime = (seconds: number): string => {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);

  if (days > 0) {
    return `${days}天${hours}小时`;
  } else if (hours > 0) {
    return `${hours}小时${minutes}分钟`;
  } else {
    return `${minutes}分钟`;
  }
};

const getProgressColor = (value: number) => {
  return colorUtils.getColorByValue(value, [
    { value: 0, color: '#52c41a' },
    { value: 70, color: '#faad14' },
    { value: 90, color: '#ff4d4f' }
  ]);
};

const getTemperatureColor = (temp: number) => {
  return colorUtils.getColorByValue(temp, [
    { value: 0, color: '#52c41a' },
    { value: 70, color: '#faad14' },
    { value: 85, color: '#ff4d4f' }
  ]);
};

const getGPUStatusLevel = (utilization: number) => {
  if (utilization >= 90) return 'high';
  if (utilization >= 70) return 'medium';
  if (utilization > 0) return 'low';
  return 'idle';
};

const getJobStatusText = (status: string) => {
  const texts = {
    running: '运行中',
    pending: '排队中',
    completed: '已完成',
    failed: '失败'
  };
  return texts[status as keyof typeof texts] || status;
};

const getJobStatusColor = (status: string) => {
  const colors = {
    running: 'green',
    pending: 'orange',
    completed: 'blue',
    failed: 'red'
  };
  return colors[status as keyof typeof colors] || 'default';
};

const getAlertLevelText = (level: string) => {
  const texts = {
    info: '信息',
    warning: '警告',
    error: '错误',
    critical: '严重'
  };
  return texts[level as keyof typeof texts] || level;
};

const getAlertLevelColor = (level: string) => {
  const colors = {
    info: 'blue',
    warning: 'orange',
    error: 'red',
    critical: 'volcano'
  };
  return colors[level as keyof typeof colors] || 'default';
};

const getActivityTypeText = (type: string) => {
  const texts = {
    job: '任务',
    user: '用户',
    system: '系统'
  };
  return texts[type as keyof typeof texts] || type;
};

const getActivityColor = (type: string) => {
  const colors = {
    job: 'blue',
    user: 'green',
    system: 'orange'
  };
  return colors[type as keyof typeof colors] || 'default';
};

const formatRelativeTime = (time: string): string => {
  const now = new Date();
  const target = new Date(time);
  const diffMs = now.getTime() - target.getTime();
  const diffMinutes = Math.floor(diffMs / (1000 * 60));

  if (diffMinutes < 60) {
    return `${diffMinutes}分钟前`;
  } else if (diffMinutes < 1440) {
    const diffHours = Math.floor(diffMinutes / 60);
    return `${diffHours}小时前`;
  } else {
    const diffDays = Math.floor(diffMinutes / 1440);
    return `${diffDays}天前`;
  }
};

const formatDuration = (seconds: number): string => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  if (hours > 0) {
    return `${hours}小时${minutes}分钟`;
  } else if (minutes > 0) {
    return `${minutes}分钟${secs}秒`;
  } else {
    return `${secs}秒`;
  }
};

// 数据更新
const updateCharts = () => {
  if (!resourceUsage.value) return;

  // 更新仪表盘
  cpuGauge.updateData(resourceUsage.value.cpu.usage, 'CPU');
  memoryGauge.updateData(resourceUsage.value.memory.usage, '内存');
  diskGauge.updateData(resourceUsage.value.disk.usage, '磁盘');

  // 更新资源分布饼图
  const resourceData = [
    { name: 'CPU', value: resourceUsage.value.cpu.usage, color: '#1890ff' },
    { name: '内存', value: resourceUsage.value.memory.usage, color: '#52c41a' },
    { name: '磁盘', value: resourceUsage.value.disk.usage, color: '#faad14' },
    { name: '网络', value: Math.min(resourceUsage.value.network.connections / 1000 * 100, 100), color: '#722ed1' }
  ];
  resourcePie.updateData(resourceData);

  // 更新任务状态饼图
  const jobData = [
    { name: '运行中', value: jobStats.value.running, color: '#52c41a' },
    { name: '排队中', value: jobStats.value.pending, color: '#faad14' },
    { name: '已完成', value: jobStats.value.completed, color: '#1890ff' },
    { name: '失败', value: jobStats.value.failed, color: '#ff4d4f' }
  ];
  jobStatusPie.updateData(jobData);
};

const loadPerformanceData = async () => {
  try {
    performanceLoading.value = true;
    
    // 生成模拟的性能趋势数据
    const now = Date.now();
    const interval = 5 * 60 * 1000; // 5分钟间隔
    const points = 24; // 24个数据点
    
    const timestamps = Array.from({ length: points }, (_, i) => now - (points - 1 - i) * interval);
    
    const performanceData = {
      timestamps,
      series: [
        {
          name: 'CPU使用率',
          data: Array.from({ length: points }, () => Math.random() * 40 + 30),
          color: '#1890ff'
        },
        {
          name: '内存使用率',
          data: Array.from({ length: points }, () => Math.random() * 30 + 40),
          color: '#52c41a'
        },
        {
          name: 'GPU使用率',
          data: Array.from({ length: points }, () => Math.random() * 50 + 25),
          color: '#722ed1'
        }
      ]
    };
    
    performanceTrend.updateData(performanceData);
    
    // 生成网络流量数据
    const networkData = {
      timestamps,
      series: [
        {
          name: '入站流量',
          data: Array.from({ length: points }, () => Math.random() * 100 + 50),
          color: '#1890ff'
        },
        {
          name: '出站流量',
          data: Array.from({ length: points }, () => Math.random() * 80 + 30),
          color: '#52c41a'
        }
      ]
    };
    
    networkChart.updateData(networkData);
    
  } catch (error) {
    message.error('加载性能数据失败');
  } finally {
    performanceLoading.value = false;
  }
};

// 事件处理
const refreshAllData = async () => {
  globalLoading.value = true;
  try {
    await Promise.all([
      refreshResourceUsage(),
      refreshSystemInfo(),
      refreshAlerts(),
      loadPerformanceData()
    ]);
    updateCharts();
  } finally {
    globalLoading.value = false;
  }
};

const handleRefreshSetting = (info: { key: string | number }) => {
  const { key } = info;
  const keyStr = String(key);
  if (keyStr === 'stop') {
    refreshInterval.value = 0;
    message.info('已停止自动刷新');
  } else {
    refreshInterval.value = parseInt(keyStr);
    message.info(`刷新间隔设置为${keyStr}秒`);
  }
};

const handleTimeRangeChange = (value: string | number | undefined) => {
  if (!value) return;
  timeRange.value = String(value);
  timeRange.value = value;
  loadPerformanceData();
};

const refreshPerformanceData = () => {
  loadPerformanceData();
};

const filterActivity = () => {
  // 筛选逻辑已在computed中实现
};

const showAlertsPanel = () => {
  message.info('告警面板功能开发中');
};

const handleAcknowledgeAlert = async (alertId: string) => {
  try {
    // TODO: Implement API call
    // await acknowledgeAlertApi(alertId);
    message.success('告警已确认');
    refreshAlerts();
  } catch (error) {
    message.error('确认失败');
  }
};

const viewAlertDetail = (alert: { id: string; summary: string; level: string; startsAt: string }) => {
  // TODO: Implement alert detail view
  message.info('查看告警详情功能开发中');
};

// 生命周期
onMounted(async () => {
  await nextTick();
  
  // 初始化数据
  await refreshAllData();
  
  // 设置定时器
  const timer = setInterval(() => {
    if (refreshInterval.value > 0) {
      updateCharts();
    }
  }, refreshInterval.value * 1000);
  
  onUnmounted(() => {
    clearInterval(timer);
  });
});

// 监听资源数据变化
import { watch } from 'vue';
watch(resourceUsage, () => {
  if (resourceUsage.value) {
    updateCharts();
  }
}, { deep: true });
</script>

<style scoped lang="scss">
.real-time-dashboard {
  padding: 24px;
  min-height: 100vh;
  background: #f5f5f5;
}

// 页面头部
.page-header {
  margin-bottom: 24px;
  
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
    position: relative;
    
    .title-icon {
      font-size: 32px;
      color: #1890ff;
    }
    
    .title-glow {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(45deg, #1890ff20, transparent);
      border-radius: 8px;
      z-index: -1;
    }
  }
  
  .page-description {
    font-size: 16px;
    margin: 0;
    color: #666;
  }
}

// 状态概览
.status-overview {
  margin-bottom: 24px;
}

.status-card {
  border-radius: 8px;
  overflow: hidden;
  
  :deep(.ant-card-body) {
    padding: 20px;
  }
}

.status-item {
  display: flex;
  align-items: center;
  gap: 16px;
  
  .status-icon {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 24px;
    
    &.status-healthy {
      background: linear-gradient(135deg, #52c41a, #73d13d);
    }
    
    &.status-warning {
      background: linear-gradient(135deg, #faad14, #ffc53d);
    }
    
    &.status-critical {
      background: linear-gradient(135deg, #ff4d4f, #ff7875);
    }
  }
  
  .status-chart {
    width: 48px;
    height: 48px;
    
    .gauge-chart {
      width: 100%;
      height: 100%;
    }
  }
  
  .status-info {
    flex: 1;
    
    .status-title {
      font-size: 14px;
      color: #666;
      margin-bottom: 4px;
    }
    
    .status-value {
      font-size: 20px;
      font-weight: 600;
      color: #262626;
      margin-bottom: 2px;
    }
    
    .status-detail {
      font-size: 12px;
      color: #999;
    }
  }
}

// 性能图表
.performance-section {
  margin-bottom: 24px;
}

.performance-card,
.resource-distribution-card {
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
}

.performance-chart,
.pie-chart,
.network-chart {
  height: 300px;
  width: 100%;
}

// GPU监控
.gpu-job-section {
  margin-bottom: 24px;
}

.gpu-card,
.job-status-card {
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
}

.gpu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.gpu-item {
  padding: 16px;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  background: #fafafa;
  transition: all 0.3s ease;
  
  &.gpu-high {
    border-left: 4px solid #ff4d4f;
  }
  
  &.gpu-medium {
    border-left: 4px solid #faad14;
  }
  
  &.gpu-low {
    border-left: 4px solid #52c41a;
  }
  
  &.gpu-idle {
    border-left: 4px solid #d9d9d9;
  }
  
  .gpu-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    
    .gpu-name {
      font-weight: 600;
      font-size: 14px;
    }
    
    .gpu-utilization {
      font-size: 16px;
      font-weight: 600;
      color: #1890ff;
    }
  }
  
  .gpu-details {
    margin-top: 12px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    
    .gpu-detail-item {
      display: flex;
      justify-content: space-between;
      font-size: 12px;
      
      .detail-label {
        color: #666;
      }
      
      .detail-value {
        font-weight: 500;
      }
    }
  }
}

.job-queue {
  margin-top: 16px;
  
  h4 {
    margin-bottom: 12px;
    font-size: 14px;
    font-weight: 600;
  }
  
  .job-progress {
    width: 100px;
  }
}

// 告警和日志
.alerts-logs-section {
  margin-bottom: 24px;
}

.alerts-card,
.activity-card {
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
}

.alert-description {
  font-size: 14px;
  margin-bottom: 4px;
}

.alert-time {
  font-size: 12px;
  color: #999;
}

.alert-actions {
  text-align: right;
}

.empty-state {
  text-align: center;
  padding: 40px 0;
}

.activity-timeline {
  max-height: 400px;
  overflow-y: auto;
  
  .activity-item {
    .activity-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 4px;
      
      .activity-type {
        font-weight: 600;
        font-size: 12px;
        color: #1890ff;
      }
      
      .activity-time {
        font-size: 11px;
        color: #999;
      }
    }
    
    .activity-content {
      font-size: 14px;
      margin-bottom: 2px;
    }
    
    .activity-user {
      font-size: 12px;
      color: #666;
    }
  }
}

// 网络和存储
.network-storage-section {
  margin-bottom: 24px;
}

.network-card,
.storage-card {
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
}

.storage-stats {
  .storage-item {
    margin-bottom: 20px;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .storage-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 8px;
      
      .storage-title {
        font-weight: 500;
        font-size: 14px;
      }
      
      .storage-usage {
        font-size: 16px;
        font-weight: 600;
        color: #1890ff;
      }
    }
    
    .storage-detail {
      font-size: 12px;
      color: #999;
      margin-top: 4px;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .real-time-dashboard {
    padding: 16px;
  }
  
  .header-content {
    flex-direction: column;
    gap: 16px;
  }
  
  .page-title {
    font-size: 24px;
    
    .title-icon {
      font-size: 28px;
    }
  }
  
  .status-item {
    gap: 12px;
    
    .status-icon,
    .status-chart {
      width: 40px;
      height: 40px;
    }
    
    .status-value {
      font-size: 18px;
    }
  }
  
  .gpu-grid {
    grid-template-columns: 1fr;
  }
  
  .performance-chart,
  .pie-chart,
  .network-chart {
    height: 250px;
  }
}
</style>