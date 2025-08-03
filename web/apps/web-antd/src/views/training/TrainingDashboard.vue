<template>
  <div class="training-dashboard-container">
    <!-- 页面头部 -->
    <Card>
      <div class="page-header">
        <div class="header-left">
          <h2>训练概览</h2>
          <p>查看训练任务的整体状态和性能指标</p>
        </div>
        <div class="header-right">
          <Space>
            <Button @click="refreshData">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary" @click="createTrainingJob">
              <PlusOutlined />
              创建训练任务
            </Button>
          </Space>
        </div>
      </div>
    </Card>

    <!-- 统计卡片 -->
    <Row :gutter="16" style="margin: 16px 0">
      <Col :span="6">
        <Card>
          <Statistic
            title="总任务数"
            :value="statistics.total"
            :value-style="{ color: '#3f8600' }"
            prefix="📊"
          />
          <div style="margin-top: 8px; font-size: 12px; color: #666">
            <span>今日新增: {{ todayNewJobs }}</span>
          </div>
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="运行中"
            :value="statistics.byStatus.running"
            :value-style="{ color: '#1890ff' }"
            prefix="🚀"
          />
          <div style="margin-top: 8px">
            <Progress
              :percent="getStatusPercent('running')"
              size="small"
              status="active"
            />
          </div>
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="队列中"
            :value="statistics.byStatus.pending + statistics.byStatus.queued"
            :value-style="{ color: '#faad14' }"
            prefix="⏳"
          />
          <div style="margin-top: 8px">
            <Progress
              :percent="getStatusPercent('queued')"
              size="small"
              stroke-color="#faad14"
            />
          </div>
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="成功率"
            :value="statistics.successRate"
            precision="1"
            suffix="%"
            :value-style="{ color: '#52c41a' }"
            prefix="✅"
          />
          <div style="margin-top: 8px; font-size: 12px; color: #666">
            <span>失败: {{ statistics.byStatus.failed }} 个</span>
          </div>
        </Card>
      </Col>
    </Row>

    <!-- 主要内容区域 -->
    <Row :gutter="16">
      <!-- 左侧：任务状态图表 -->
      <Col :span="12">
        <Card title="任务状态分布" :bordered="false">
          <div ref="statusChartRef" style="height: 300px"></div>
        </Card>
      </Col>

      <!-- 右侧：框架使用统计 -->
      <Col :span="12">
        <Card title="训练框架分布" :bordered="false">
          <div ref="frameworkChartRef" style="height: 300px"></div>
        </Card>
      </Col>
    </Row>

    <!-- 下方：最近任务和资源使用 -->
    <Row :gutter="16" style="margin-top: 16px">
      <!-- 最近训练任务 -->
      <Col :span="16">
        <Card title="最近训练任务" :bordered="false">
          <Table
            :columns="recentJobsColumns"
            :data-source="recentJobs"
            :pagination="false"
            size="small"
            :loading="loadingRecentJobs"
          >
            <!-- 任务名称 -->
            <template #name="{ record }">
              <div>
                <Button type="link" @click="viewJobDetail(record)">
                  {{ record.name }}
                </Button>
                <div style="font-size: 12px; color: #999">
                  {{ record.framework }}
                </div>
              </div>
            </template>

            <!-- 状态 -->
            <template #status="{ record }">
              <Tag :color="getJobStatusColor(record.status)">
                {{ getJobStatusLabel(record.status) }}
              </Tag>
            </template>

            <!-- 进度 -->
            <template #progress="{ record }">
              <Progress
                :percent="record.progress"
                size="small"
                :status="record.status === 'failed' ? 'exception' : 'active'"
              />
            </template>

            <!-- 时长 -->
            <template #duration="{ record }">
              <span>{{ formatDuration(record.duration || 0) }}</span>
            </template>

            <!-- 创建时间 -->
            <template #createTime="{ record }">
              <span>{{ formatDateTime(record.createTime, 'MM-DD HH:mm') }}</span>
            </template>

            <!-- 操作 -->
            <template #action="{ record }">
              <Space size="small">
                <Button type="link" size="small" @click="viewJobDetail(record)">
                  详情
                </Button>
                <Button 
                  type="link" 
                  size="small" 
                  @click="controlJob(record, 'stop')"
                  :disabled="!canControl(record.status)"
                  danger
                >
                  停止
                </Button>
              </Space>
            </template>
          </Table>
          
          <div style="text-align: center; margin-top: 16px">
            <Button @click="viewAllJobs">查看全部任务</Button>
          </div>
        </Card>
      </Col>

      <!-- 资源使用情况 -->
      <Col :span="8">
        <Card title="资源使用情况" :bordered="false">
          <div class="resource-stats">
            <div class="resource-item">
              <div class="resource-label">GPU 利用率</div>
              <Progress
                :percent="Math.round(statistics.resourceUtilization.avgGpuUtilization)"
                stroke-color="#52c41a"
              />
              <div class="resource-detail">
                平均: {{ statistics.resourceUtilization.avgGpuUtilization.toFixed(1) }}%
              </div>
            </div>

            <div class="resource-item">
              <div class="resource-label">CPU 利用率</div>
              <Progress
                :percent="Math.round(statistics.resourceUtilization.avgCpuUtilization)"
                stroke-color="#1890ff"
              />
              <div class="resource-detail">
                平均: {{ statistics.resourceUtilization.avgCpuUtilization.toFixed(1) }}%
              </div>
            </div>

            <div class="resource-item">
              <div class="resource-label">内存使用</div>
              <Progress
                :percent="Math.round(statistics.resourceUtilization.avgMemoryUsage)"
                stroke-color="#faad14"
              />
              <div class="resource-detail">
                平均: {{ statistics.resourceUtilization.avgMemoryUsage.toFixed(1) }}%
              </div>
            </div>
          </div>

          <Divider />

          <div class="training-metrics">
            <div class="metric-item">
              <span class="metric-label">总GPU小时数</span>
              <span class="metric-value">{{ statistics.totalGpuHours.toFixed(1) }}h</span>
            </div>
            <div class="metric-item">
              <span class="metric-label">平均任务时长</span>
              <span class="metric-value">{{ statistics.avgJobDuration.toFixed(1) }}h</span>
            </div>
          </div>
        </Card>
      </Col>
    </Row>

    <!-- 快速操作 -->
    <Card title="快速操作" style="margin-top: 16px">
      <Row :gutter="16">
        <Col :span="6">
          <div class="quick-action-item" @click="createFromTemplate">
            <FileTextOutlined class="action-icon" />
            <div class="action-content">
              <div class="action-title">从模板创建</div>
              <div class="action-desc">使用预设模板快速创建训练任务</div>
            </div>
          </div>
        </Col>
        <Col :span="6">
          <div class="quick-action-item" @click="viewTemplates">
            <DatabaseOutlined class="action-icon" />
            <div class="action-content">
              <div class="action-title">模板管理</div>
              <div class="action-desc">管理和编辑训练模板</div>
            </div>
          </div>
        </Col>
        <Col :span="6">
          <div class="quick-action-item" @click="viewExperiments">
            <ExperimentOutlined class="action-icon" />
            <div class="action-content">
              <div class="action-title">实验跟踪</div>
              <div class="action-desc">查看和比较实验结果</div>
            </div>
          </div>
        </Col>
        <Col :span="6">
          <div class="quick-action-item" @click="viewQueues">
            <ClusterOutlined class="action-icon" />
            <div class="action-content">
              <div class="action-title">队列管理</div>
              <div class="action-desc">管理训练队列和资源分配</div>
            </div>
          </div>
        </Col>
      </Row>
    </Card>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, nextTick } from 'vue';
import {
  Card,
  Row,
  Col,
  Button,
  Space,
  Statistic,
  Progress,
  Table,
  Tag,
  Divider,
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  PlusOutlined,
  FileTextOutlined,
  DatabaseOutlined,
  ExperimentOutlined,
  ClusterOutlined,
} from '@ant-design/icons-vue';
import { useRouter } from 'vue-router';
import type { TrainingStatistics, TrainingJob } from '#/api/types';
import { getTrainingStatistics, getMyTrainingJobs, controlTrainingJob } from '#/api';
import { formatDateTime, formatDuration } from '#/utils/date';

// 使用ECharts进行图表渲染
const router = useRouter();
const statusChartRef = ref();
const frameworkChartRef = ref();
const loadingRecentJobs = ref(false);

// 统计数据
const statistics = ref<TrainingStatistics>({
  total: 0,
  active: 0,
  inactive: 0,
  lastUpdated: '',
  byStatus: {
    pending: 0,
    queued: 0,
    running: 0,
    completed: 0,
    failed: 0,
    cancelled: 0,
    paused: 0,
    stopped: 0,
  },
  byFramework: {
    tensorflow: 0,
    pytorch: 0,
    keras: 0,
    paddlepaddle: 0,
    mindspore: 0,
    custom: 0,
  },
  byPriority: {
    low: 0,
    medium: 0,
    high: 0,
    urgent: 0,
  },
  totalGpuHours: 0,
  avgJobDuration: 0,
  successRate: 0,
  resourceUtilization: {
    avgGpuUtilization: 0,
    avgCpuUtilization: 0,
    avgMemoryUsage: 0,
  },
});

const recentJobs = ref<TrainingJob[]>([]);
const todayNewJobs = ref(0);

// 最近任务表格列定义
const recentJobsColumns = [
  {
    title: '任务名称',
    key: 'name',
    slots: { customRender: 'name' },
    width: 200,
  },
  {
    title: '状态',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100,
  },
  {
    title: '进度',
    key: 'progress',
    slots: { customRender: 'progress' },
    width: 120,
  },
  {
    title: '时长',
    key: 'duration',
    slots: { customRender: 'duration' },
    width: 80,
  },
  {
    title: '创建时间',
    key: 'createTime',
    slots: { customRender: 'createTime' },
    width: 100,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 120,
  },
];

defineOptions({ name: 'TrainingDashboard' });

// 工具方法
const getStatusPercent = (status: string) => {
  const total = statistics.value.total || 1;
  switch (status) {
    case 'running':
      return Math.round((statistics.value.byStatus.running / total) * 100);
    case 'queued':
      return Math.round(((statistics.value.byStatus.pending + statistics.value.byStatus.queued) / total) * 100);
    default:
      return 0;
  }
};

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

const canControl = (status: string) => {
  return ['running', 'queued', 'pending'].includes(status);
};

// 数据加载
const loadStatistics = async () => {
  try {
    const response = await getTrainingStatistics();
    statistics.value = response;
  } catch (error) {
    console.error('加载统计数据失败:', error);
  }
};

const loadRecentJobs = async () => {
  try {
    loadingRecentJobs.value = true;
    const response = await getMyTrainingJobs({ pageSize: 10, sortBy: 'createTime', sortOrder: 'desc' });
    recentJobs.value = response.data;
    
    // 计算今日新增任务
    const today = new Date().toDateString();
    todayNewJobs.value = response.data.filter(job => 
      new Date(job.createTime).toDateString() === today
    ).length;
  } catch (error) {
    message.error('加载最近任务失败');
  } finally {
    loadingRecentJobs.value = false;
  }
};

const refreshData = () => {
  loadStatistics();
  loadRecentJobs();
};

// 图表初始化
const initCharts = async () => {
  await nextTick();
  
  // 这里应该使用实际的图表库（如ECharts）来渲染图表
  // 为了简化，这里只是占位符
  if (statusChartRef.value) {
    statusChartRef.value.innerHTML = '<div style="text-align: center; line-height: 300px; color: #999;">状态分布图表</div>';
  }
  
  if (frameworkChartRef.value) {
    frameworkChartRef.value.innerHTML = '<div style="text-align: center; line-height: 300px; color: #999;">框架分布图表</div>';
  }
};

// 事件处理
const createTrainingJob = () => {
  router.push('/training/queue');
};

const createFromTemplate = () => {
  router.push('/training/template');
};

const viewTemplates = () => {
  router.push('/training/template');
};

const viewExperiments = () => {
  router.push('/training/experiments');
};

const viewQueues = () => {
  router.push('/training/queue');
};

const viewAllJobs = () => {
  router.push('/training/history');
};

const viewJobDetail = (job: TrainingJob) => {
  router.push(`/training/jobs/${job.id}`);
};

const controlJob = async (job: TrainingJob, action: string) => {
  try {
    await controlTrainingJob({ id: job.id, action });
    message.success(`${action === 'stop' ? '停止' : action}任务成功`);
    refreshData();
  } catch (error) {
    message.error('操作失败');
  }
};

// 初始化
onMounted(() => {
  refreshData();
  initCharts();
});
</script>

<style scoped lang="scss">
.training-dashboard-container {
  padding: 24px;
  background: #f5f5f5;
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .header-left {
    h2 {
      margin: 0;
      color: #1890ff;
    }
    
    p {
      margin: 8px 0 0 0;
      color: #666;
    }
  }
}

.resource-stats {
  .resource-item {
    margin-bottom: 16px;
    
    .resource-label {
      font-size: 14px;
      margin-bottom: 8px;
      color: #333;
    }
    
    .resource-detail {
      font-size: 12px;
      color: #999;
      margin-top: 4px;
    }
  }
}

.training-metrics {
  .metric-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
    
    .metric-label {
      color: #666;
    }
    
    .metric-value {
      font-weight: 500;
      color: #333;
    }
  }
}

.quick-action-item {
  display: flex;
  align-items: center;
  padding: 16px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  
  &:hover {
    border-color: #1890ff;
    background-color: #f0f8ff;
  }
  
  .action-icon {
    font-size: 24px;
    color: #1890ff;
    margin-right: 12px;
  }
  
  .action-content {
    .action-title {
      font-size: 16px;
      font-weight: 500;
      margin-bottom: 4px;
    }
    
    .action-desc {
      font-size: 12px;
      color: #999;
    }
  }
}
</style>
