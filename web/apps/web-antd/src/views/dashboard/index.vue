<template>
  <div class="dashboard-container">
    <!-- 页面头部 -->
    <div class="dashboard-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <MonitorDashboardIcon class="title-icon" />
            <span class="title-text">VolcTrain 控制台</span>
          </h1>
          <p class="page-description">
            <span class="description-text">火山引擎训练平台 - 实时监控与管理中心</span>
          </p>
        </div>
        <div class="action-section">
          <Space>
            <Badge :count="systemAlerts" :offset="[10, 0]">
              <Button type="primary" @click="showAlertsModal">
                <BellIcon />
                告警中心
              </Button>
            </Badge>
            <Dropdown>
              <Button>
                <RefreshCwIcon :class="{ 'animate-spin': refreshing }" />
                自动刷新
                <DownIcon />
              </Button>
              <template #overlay>
                <Menu @click="handleRefreshInterval">
                  <Menu.Item key="10">10秒</Menu.Item>
                  <Menu.Item key="30">30秒</Menu.Item>
                  <Menu.Item key="60">1分钟</Menu.Item>
                  <Menu.Item key="300">5分钟</Menu.Item>
                  <Menu.Item key="0">关闭</Menu.Item>
                </Menu>
              </template>
            </Dropdown>
            <Button type="primary" @click="refreshData">
              <RefreshCwIcon :class="{ 'animate-spin': loading }" />
              刷新数据
            </Button>
          </Space>
        </div>
      </div>
    </div>

    <!-- 系统状态卡片 -->
    <div class="status-cards">
      <Row :gutter="[16, 16]">
        <Col :span="6">
          <Card class="status-card status-card-primary">
            <div class="status-card-content">
              <div class="status-icon">
                <PlayCircleIcon />
              </div>
              <div class="status-info">
                <div class="status-title">运行中的训练任务</div>
                <div class="status-value">{{ dashboardData.runningJobs }}</div>
                <div class="status-trend">
                  <TrendingUpIcon class="trend-icon trend-up" />
                  <span class="trend-text">较昨日 +12%</span>
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :span="6">
          <Card class="status-card status-card-success">
            <div class="status-card-content">
              <div class="status-icon">
                <CpuIcon />
              </div>
              <div class="status-info">
                <div class="status-title">GPU 设备总数</div>
                <div class="status-value">{{ dashboardData.totalGpus }}</div>
                <div class="status-trend">
                  <ActivityIcon class="trend-icon trend-stable" />
                  <span class="trend-text">使用率 {{ dashboardData.gpuUtilization }}%</span>
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :span="6">
          <Card class="status-card status-card-warning">
            <div class="status-card-content">
              <div class="status-icon">
                <DatabaseIcon />
              </div>
              <div class="status-info">
                <div class="status-title">数据集数量</div>
                <div class="status-value">{{ dashboardData.totalDatasets }}</div>
                <div class="status-trend">
                  <TrendingUpIcon class="trend-icon trend-up" />
                  <span class="trend-text">本月新增 {{ dashboardData.newDatasets }}</span>
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :span="6">
          <Card class="status-card status-card-info">
            <div class="status-card-content">
              <div class="status-icon">
                <PackageIcon />
              </div>
              <div class="status-info">
                <div class="status-title">训练模型</div>
                <div class="status-value">{{ dashboardData.totalModels }}</div>
                <div class="status-trend">
                  <CheckCircleIcon class="trend-icon trend-up" />
                  <span class="trend-text">已部署 {{ dashboardData.deployedModels }}</span>
                </div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 主要内容区域 -->
    <div class="dashboard-content">
      <Row :gutter="[16, 16]">
        <!-- 左侧区域 -->
        <Col :span="16">
          <!-- 训练任务统计图表 -->
          <Card class="chart-card" title="训练任务趋势" :loading="loading">
            <template #extra>
              <Space>
                <Radio.Group v-model:value="trainingChartTimeRange" @change="loadTrainingChart">
                  <Radio.Button value="7d">7天</Radio.Button>
                  <Radio.Button value="30d">30天</Radio.Button>
                  <Radio.Button value="90d">90天</Radio.Button>
                </Radio.Group>
              </Space>
            </template>
            <div ref="trainingChartRef" class="chart-container" />
          </Card>

          <!-- GPU 使用情况 -->
          <Card class="chart-card" title="GPU 集群状态" style="margin-top: 16px">
            <div class="gpu-cluster-list">
              <div
                v-for="cluster in dashboardData.clusters"
                :key="cluster.id"
                class="cluster-item"
              >
                <div class="cluster-info">
                  <div class="cluster-name">{{ cluster.name }}</div>
                  <div class="cluster-stats">
                    {{ cluster.availableGpus }}/{{ cluster.totalGpus }} GPU 可用
                  </div>
                </div>
                <div class="cluster-progress">
                  <Progress
                    :percent="Math.round((cluster.allocatedGpus / cluster.totalGpus) * 100)"
                    :status="cluster.status === 'healthy' ? 'normal' : 'exception'"
                    :stroke-color="{
                      from: '#87d068',
                      to: '#52c41a',
                    }"
                  />
                </div>
                <div class="cluster-status">
                  <Badge
                    :status="cluster.status === 'healthy' ? 'success' : 'error'"
                    :text="cluster.status === 'healthy' ? '健康' : '异常'"
                  />
                </div>
              </div>
            </div>
          </Card>
        </Col>

        <!-- 右侧区域 -->
        <Col :span="8">
          <!-- 最近活动 -->
          <Card class="activity-card" title="最近活动">
            <template #extra>
              <a @click="$router.push('/training/jobs')">查看全部</a>
            </template>
            <div class="activity-list">
              <div
                v-for="activity in dashboardData.recentActivities"
                :key="activity.id"
                class="activity-item"
              >
                <div class="activity-icon">
                  <PlayCircleIcon v-if="activity.type === 'job_started'" class="icon-success" />
                  <CheckCircleIcon v-else-if="activity.type === 'job_completed'" class="icon-success" />
                  <XCircleIcon v-else-if="activity.type === 'job_failed'" class="icon-error" />
                  <UploadIcon v-else-if="activity.type === 'dataset_uploaded'" class="icon-info" />
                </div>
                <div class="activity-content">
                  <div class="activity-title">{{ activity.title }}</div>
                  <div class="activity-desc">{{ activity.description }}</div>
                  <div class="activity-time">{{ formatTime(activity.timestamp) }}</div>
                </div>
              </div>
            </div>
          </Card>

          <!-- 快速操作 -->
          <Card class="quick-actions-card" title="快速操作" style="margin-top: 16px">
            <div class="quick-actions">
              <Button
                type="primary"
                size="large"
                block
                class="quick-action-btn"
                @click="$router.push('/training/jobs')"
              >
                <PlusCircleIcon />
                创建训练任务
              </Button>
              <Button
                size="large"
                block
                class="quick-action-btn"
                @click="$router.push('/data/upload')"
              >
                <UploadIcon />
                上传数据集
              </Button>
              <Button
                size="large"
                block
                class="quick-action-btn"
                @click="$router.push('/gpu/clusters')"
              >
                <ServerIcon />
                管理 GPU 集群
              </Button>
              <Button
                size="large"
                block
                class="quick-action-btn"
                @click="$router.push('/monitoring/system')"
              >
                <BarChart3Icon />
                查看监控
              </Button>
            </div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 告警弹窗 -->
    <Modal
      v-model:open="alertsModalVisible"
      title="系统告警"
      :footer="null"
      width="800px"
    >
      <div class="alerts-list">
        <div
          v-for="alert in systemAlertsList"
          :key="alert.id"
          class="alert-item"
        >
          <div class="alert-severity">
            <Badge
              :status="alert.severity === 'critical' ? 'error' : 
                     alert.severity === 'warning' ? 'warning' : 'default'"
              :text="alert.severity"
            />
          </div>
          <div class="alert-content">
            <div class="alert-title">{{ alert.ruleName }}</div>
            <div class="alert-message">{{ alert.message }}</div>
            <div class="alert-time">{{ formatTime(alert.triggeredAt) }}</div>
          </div>
          <div class="alert-actions">
            <Button size="small" @click="resolveAlert(alert.id)">
              标记已解决
            </Button>
          </div>
        </div>
      </div>
    </Modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import {
  Card,
  Row,
  Col,
  Button,
  Space,
  Badge,
  Dropdown,
  Menu,
  Progress,
  Modal,
  Radio,
} from 'ant-design-vue';
import { message } from 'ant-design-vue';
import * as echarts from 'echarts';
import { formatDistanceToNow } from 'date-fns';
import { zhCN } from 'date-fns/locale';

// 图标组件
import {
  MonitorDashboard as MonitorDashboardIcon,
  Bell as BellIcon,
  RefreshCw as RefreshCwIcon,
  ChevronDown as DownIcon,
  PlayCircle as PlayCircleIcon,
  Cpu as CpuIcon,
  Database as DatabaseIcon,
  Package as PackageIcon,
  TrendingUp as TrendingUpIcon,
  Activity as ActivityIcon,
  CheckCircle as CheckCircleIcon,
  XCircle as XCircleIcon,
  Upload as UploadIcon,
  PlusCircle as PlusCircleIcon,
  Server as ServerIcon,
  BarChart3 as BarChart3Icon,
} from 'lucide-vue-next';

// API 服务
import { trainingService, gpuService } from '#/api/services';
import type { TrainingApi, GpuApi, MonitorApi } from '#/types/api';

// 定义数据类型
interface DashboardData {
  runningJobs: number;
  totalGpus: number;
  gpuUtilization: number;
  totalDatasets: number;
  newDatasets: number;
  totalModels: number;
  deployedModels: number;
  clusters: Array<{
    id: string;
    name: string;
    totalGpus: number;
    availableGpus: number;
    allocatedGpus: number;
    status: 'healthy' | 'degraded' | 'unhealthy';
  }>;
  recentActivities: Array<{
    id: string;
    type: string;
    title: string;
    description: string;
    timestamp: string;
  }>;
}

defineOptions({ name: 'Dashboard' });

const router = useRouter();

// 响应式数据
const loading = ref(false);
const refreshing = ref(false);
const systemAlerts = ref(3);
const alertsModalVisible = ref(false);
const trainingChartTimeRange = ref('7d');
const trainingChartRef = ref<HTMLDivElement>();

// 仪表板数据
const dashboardData = ref<DashboardData>({
  runningJobs: 24,
  totalGpus: 128,
  gpuUtilization: 78,
  totalDatasets: 156,
  newDatasets: 12,
  totalModels: 89,
  deployedModels: 23,
  clusters: [
    {
      id: '1',
      name: 'GPU-Cluster-01',
      totalGpus: 64,
      availableGpus: 18,
      allocatedGpus: 46,
      status: 'healthy',
    },
    {
      id: '2',
      name: 'GPU-Cluster-02',
      totalGpus: 64,
      availableGpus: 12,
      allocatedGpus: 52,
      status: 'healthy',
    },
  ],
  recentActivities: [
    {
      id: '1',
      type: 'job_started',
      title: 'PyTorch 图像分类任务已启动',
      description: 'ResNet-50 模型训练开始',
      timestamp: '2024-01-20T10:30:00Z',
    },
    {
      id: '2',
      type: 'job_completed',
      title: 'TensorFlow NLP 任务完成',
      description: 'BERT 模型训练成功完成',
      timestamp: '2024-01-20T09:45:00Z',
    },
    {
      id: '3',
      type: 'dataset_uploaded',
      title: '新数据集上传完成',
      description: 'ImageNet-Mini 数据集已可用',
      timestamp: '2024-01-20T09:15:00Z',
    },
  ],
});

// 系统告警列表
const systemAlertsList = ref<MonitorApi.AlertRecord[]>([
  {
    id: '1',
    ruleId: 'gpu-temp',
    ruleName: 'GPU 温度过高',
    severity: 'warning',
    status: 'firing',
    message: 'GPU-01 温度达到 85°C，超过阈值',
    value: 85,
    triggeredAt: '2024-01-20T10:25:00Z',
  },
  {
    id: '2',
    ruleId: 'disk-usage',
    ruleName: '磁盘空间不足',
    severity: 'critical',
    status: 'firing',
    message: '存储节点 /data 使用率达到 95%',
    value: 95,
    triggeredAt: '2024-01-20T10:20:00Z',
  },
]);

// 自动刷新定时器
let refreshTimer: NodeJS.Timeout | null = null;

/**
 * 初始化训练任务趋势图表
 */
const initTrainingChart = () => {
  if (!trainingChartRef.value) return;

  const chart = echarts.init(trainingChartRef.value);
  
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
        stack: 'Total',
        smooth: true,
        lineStyle: {
          width: 2,
        },
        emphasis: {
          focus: 'series',
        },
        data: [12, 15, 18, 22, 20, 24, 26],
      },
      {
        name: '已完成',
        type: 'line',
        stack: 'Total',
        smooth: true,
        lineStyle: {
          width: 2,
        },
        emphasis: {
          focus: 'series',
        },
        data: [8, 12, 16, 18, 22, 25, 28],
      },
      {
        name: '已失败',
        type: 'line',
        stack: 'Total',
        smooth: true,
        lineStyle: {
          width: 2,
        },
        emphasis: {
          focus: 'series',
        },
        data: [2, 3, 2, 4, 3, 2, 1],
      },
    ],
  };

  chart.setOption(option);

  // 响应式图表
  window.addEventListener('resize', () => {
    chart.resize();
  });
};

/**
 * 加载仪表板数据
 */
const loadDashboardData = async () => {
  try {
    loading.value = true;
    
    // 这里可以并行加载多个API
    // const [trainingStats, gpuStats, datasetStats] = await Promise.all([
    //   trainingService.getTrainingStatistics(),
    //   gpuService.getResourceStatistics(),
    //   datasetService.getStatistics(),
    // ]);
    
    // 模拟数据加载延迟
    await new Promise(resolve => setTimeout(resolve, 500));
    
    message.success('数据加载完成');
  } catch (error) {
    message.error('数据加载失败');
  } finally {
    loading.value = false;
  }
};

/**
 * 刷新数据
 */
const refreshData = async () => {
  await loadDashboardData();
};

/**
 * 加载训练任务趋势图表
 */
const loadTrainingChart = () => {
  // 根据时间范围重新加载图表数据
  initTrainingChart();
};

/**
 * 处理刷新间隔设置
 */
const handleRefreshInterval = ({ key }: { key: string }) => {
  const interval = parseInt(key);
  
  if (refreshTimer) {
    clearInterval(refreshTimer);
    refreshTimer = null;
  }
  
  if (interval > 0) {
    refreshTimer = setInterval(() => {
      refreshing.value = true;
      loadDashboardData().finally(() => {
        refreshing.value = false;
      });
    }, interval * 1000);
    
    message.success(`已设置 ${interval} 秒自动刷新`);
  } else {
    message.info('已关闭自动刷新');
  }
};

/**
 * 显示告警弹窗
 */
const showAlertsModal = () => {
  alertsModalVisible.value = true;
};

/**
 * 解决告警
 */
const resolveAlert = (alertId: string) => {
  systemAlertsList.value = systemAlertsList.value.filter(alert => alert.id !== alertId);
  systemAlerts.value = systemAlertsList.value.length;
  message.success('告警已标记为解决');
};

/**
 * 格式化时间
 */
const formatTime = (timestamp: string) => {
  return formatDistanceToNow(new Date(timestamp), { 
    addSuffix: true,
    locale: zhCN,
  });
};

// 生命周期
onMounted(() => {
  loadDashboardData();
  initTrainingChart();
});

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer);
  }
});
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 0;
  background: #f5f7fa;
  min-height: 100vh;
}

.dashboard-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 24px 32px;
  margin-bottom: 24px;

  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .title-section {
    .page-title {
      display: flex;
      align-items: center;
      font-size: 28px;
      font-weight: 600;
      margin: 0;
      
      .title-icon {
        margin-right: 12px;
        font-size: 32px;
      }
    }

    .page-description {
      margin: 8px 0 0 44px;
      font-size: 16px;
      opacity: 0.9;
    }
  }
}

.status-cards {
  margin-bottom: 24px;
  padding: 0 24px;
}

.status-card {
  border: none;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  }

  &.status-card-primary .status-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.status-card-success .status-icon {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  }

  &.status-card-warning .status-icon {
    background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  }

  &.status-card-info .status-icon {
    background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  }
}

.status-card-content {
  display: flex;
  align-items: center;
  padding: 8px;

  .status-icon {
    width: 64px;
    height: 64px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 24px;
    margin-right: 16px;
  }

  .status-info {
    flex: 1;

    .status-title {
      font-size: 14px;
      color: #666;
      margin-bottom: 4px;
    }

    .status-value {
      font-size: 24px;
      font-weight: 600;
      color: #333;
      margin-bottom: 4px;
    }

    .status-trend {
      display: flex;
      align-items: center;
      font-size: 12px;

      .trend-icon {
        margin-right: 4px;
        
        &.trend-up {
          color: #52c41a;
        }
        
        &.trend-stable {
          color: #1890ff;
        }
      }

      .trend-text {
        color: #666;
      }
    }
  }
}

.dashboard-content {
  padding: 0 24px;
}

.chart-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  .chart-container {
    height: 300px;
    width: 100%;
  }
}

.gpu-cluster-list {
  .cluster-item {
    display: flex;
    align-items: center;
    padding: 16px 0;
    border-bottom: 1px solid #f0f0f0;

    &:last-child {
      border-bottom: none;
    }

    .cluster-info {
      flex: 1;
      margin-right: 16px;

      .cluster-name {
        font-weight: 600;
        margin-bottom: 4px;
      }

      .cluster-stats {
        font-size: 12px;
        color: #666;
      }
    }

    .cluster-progress {
      flex: 2;
      margin-right: 16px;
    }

    .cluster-status {
      width: 80px;
    }
  }
}

.activity-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.activity-list {
  .activity-item {
    display: flex;
    align-items: flex-start;
    padding: 12px 0;
    border-bottom: 1px solid #f0f0f0;

    &:last-child {
      border-bottom: none;
    }

    .activity-icon {
      margin-right: 12px;
      margin-top: 2px;
      
      .icon-success {
        color: #52c41a;
      }
      
      .icon-error {
        color: #ff4d4f;
      }
      
      .icon-info {
        color: #1890ff;
      }
    }

    .activity-content {
      flex: 1;

      .activity-title {
        font-weight: 500;
        margin-bottom: 4px;
        color: #333;
      }

      .activity-desc {
        font-size: 12px;
        color: #666;
        margin-bottom: 4px;
      }

      .activity-time {
        font-size: 11px;
        color: #999;
      }
    }
  }
}

.quick-actions-card {
  border-radius: 12px;
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.quick-actions {
  .quick-action-btn {
    margin-bottom: 12px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 500;
    border-radius: 8px;

    &:last-child {
      margin-bottom: 0;
    }

    :deep(.anticon) {
      margin-right: 8px;
    }
  }
}

.alerts-list {
  .alert-item {
    display: flex;
    align-items: flex-start;
    padding: 16px 0;
    border-bottom: 1px solid #f0f0f0;

    &:last-child {
      border-bottom: none;
    }

    .alert-severity {
      margin-right: 16px;
      margin-top: 2px;
    }

    .alert-content {
      flex: 1;
      margin-right: 16px;

      .alert-title {
        font-weight: 500;
        margin-bottom: 4px;
      }

      .alert-message {
        color: #666;
        margin-bottom: 4px;
      }

      .alert-time {
        font-size: 12px;
        color: #999;
      }
    }

    .alert-actions {
      margin-top: 2px;
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