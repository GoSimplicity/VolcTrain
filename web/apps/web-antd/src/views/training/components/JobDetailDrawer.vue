<template>
  <Drawer
    v-model:open="visible"
    title="训练任务详情"
    width="800"
    placement="right"
    class="job-detail-drawer"
  >
    <div v-if="job" class="drawer-content">
      <!-- 基本信息卡片 -->
      <Card title="基本信息" class="info-card">
        <Descriptions :column="2" bordered>
          <Descriptions.Item label="任务名称">
            <div class="job-name">
              {{ job.name }}
              <Tag :color="getJobStatusColor(job.status)" class="status-tag">
                {{ getJobStatusLabel(job.status) }}
              </Tag>
            </div>
          </Descriptions.Item>
          <Descriptions.Item label="任务ID">
            <code class="job-id">{{ job.id }}</code>
          </Descriptions.Item>
          <Descriptions.Item label="训练框架">
            <Tag color="blue">{{ job.framework }}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="分布式类型">
            {{ job.distributedType }}
          </Descriptions.Item>
          <Descriptions.Item label="创建者">
            <Avatar size="small">{{ job.creatorName?.[0] }}</Avatar>
            <span style="margin-left: 8px">{{ job.creatorName }}</span>
          </Descriptions.Item>
          <Descriptions.Item label="工作空间">
            {{ job.workspaceId }}
          </Descriptions.Item>
          <Descriptions.Item label="训练队列">
            {{ job.queueId }}
          </Descriptions.Item>
          <Descriptions.Item label="优先级">
            <Tag :color="getPriorityColor(job.priority)">
              {{ getPriorityLabel(job.priority) }}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="创建时间" :span="2">
            {{ formatDateTime(job.createTime, 'YYYY-MM-DD HH:mm:ss') }}
          </Descriptions.Item>
          <Descriptions.Item label="描述" :span="2">
            {{ job.description || '暂无描述' }}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <!-- 资源配置卡片 -->
      <Card title="资源配置" class="info-card">
        <Row :gutter="16">
          <Col :span="12">
            <div class="resource-item">
              <div class="resource-header">
                <span class="resource-label">CPU 需求</span>
                <span class="resource-value">{{ job.resourceRequirements.cpu }} 核</span>
              </div>
              <Progress 
                :percent="75" 
                size="small" 
                stroke-color="#52c41a"
                :show-info="false"
              />
            </div>
          </Col>
          <Col :span="12">
            <div class="resource-item">
              <div class="resource-header">
                <span class="resource-label">内存需求</span>
                <span class="resource-value">{{ job.resourceRequirements.memory }} GB</span>
              </div>
              <Progress 
                :percent="68" 
                size="small" 
                stroke-color="#1890ff"
                :show-info="false"
              />
            </div>
          </Col>
          <Col :span="12" v-if="job.resourceRequirements.gpu">
            <div class="resource-item">
              <div class="resource-header">
                <span class="resource-label">GPU 需求</span>
                <span class="resource-value">{{ job.resourceRequirements.gpu }} 卡</span>
              </div>
              <Progress 
                :percent="92" 
                size="small" 
                stroke-color="#722ed1"
                :show-info="false"
              />
            </div>
          </Col>
          <Col :span="12">
            <div class="resource-item">
              <div class="resource-header">
                <span class="resource-label">存储需求</span>
                <span class="resource-value">{{ job.resourceRequirements.storage }} GB</span>
              </div>
              <Progress 
                :percent="45" 
                size="small" 
                stroke-color="#faad14"
                :show-info="false"
              />
            </div>
          </Col>
        </Row>

        <Divider />

        <Descriptions :column="2" size="small">
          <Descriptions.Item label="副本数量">
            {{ job.replicas }} 个
          </Descriptions.Item>
          <Descriptions.Item label="最大重试次数">
            {{ job.maxRetries }} 次
          </Descriptions.Item>
          <Descriptions.Item label="检查点启用">
            <Tag :color="job.checkpointEnabled ? 'green' : 'default'">
              {{ job.checkpointEnabled ? '是' : '否' }}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="检查点间隔" v-if="job.checkpointEnabled">
            {{ job.checkpointInterval }} 步
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <!-- 训练配置卡片 -->
      <Card title="训练配置" class="info-card">
        <Descriptions :column="1" bordered>
          <Descriptions.Item label="容器镜像">
            <code class="config-value">{{ job.image }}</code>
          </Descriptions.Item>
          <Descriptions.Item label="启动命令">
            <code class="config-value">{{ Array.isArray(job.command) ? job.command.join(' ') : job.command }}</code>
          </Descriptions.Item>
          <Descriptions.Item label="启动参数" v-if="job.args && job.args.length > 0">
            <code class="config-value">{{ Array.isArray(job.args) ? job.args.join(' ') : job.args }}</code>
          </Descriptions.Item>
          <Descriptions.Item label="输入数据路径" v-if="job.inputDataPath">
            <code class="config-value">{{ job.inputDataPath }}</code>
          </Descriptions.Item>
          <Descriptions.Item label="输出数据路径" v-if="job.outputDataPath">
            <code class="config-value">{{ job.outputDataPath }}</code>
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <!-- 运行状态卡片 -->
      <Card title="运行状态" class="info-card">
        <div class="status-info">
          <Row :gutter="24">
            <Col :span="8">
              <Statistic
                title="进度"
                :value="job.progress"
                suffix="%"
                :value-style="{ color: getProgressColor(job.progress) }"
              />
              <Progress
                :percent="job.progress"
                :status="getProgressStatus(job.status)"
                style="margin-top: 8px"
              />
            </Col>
            <Col :span="8">
              <Statistic
                title="运行时长"
                :value="formatDuration(job.duration || 0)"
                :value-style="{ color: '#1890ff' }"
              />
            </Col>
            <Col :span="8">
              <Statistic
                title="提交时间"
                :value="formatRelativeTime(job.submitTime)"
                :value-style="{ color: '#666' }"
              />
            </Col>
          </Row>
        </div>

        <Divider />

        <!-- Pod 状态 -->
        <div class="pod-status">
          <h4>Pod 状态</h4>
          <Table
            :columns="podColumns"
            :data-source="mockPods"
            :pagination="false"
            size="small"
            class="pod-table"
          >
            <template #podStatus="{ record }">
              <Tag :color="getPodStatusColor(record.status)">
                {{ record.status }}
              </Tag>
            </template>
            <template #resources="{ record }">
              <div class="pod-resources">
                <div>CPU: {{ record.cpu }}</div>
                <div>内存: {{ record.memory }}</div>
                <div v-if="record.gpu">GPU: {{ record.gpu }}</div>
              </div>
            </template>
          </Table>
        </div>
      </Card>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <Space>
          <Button 
            type="primary" 
            @click="viewLogs"
            :loading="actionLoading"
          >
            <FileTextOutlined />
            查看日志
          </Button>
          <Button 
            @click="controlJob('pause')"
            :disabled="!canPause(job.status)"
            :loading="actionLoading"
          >
            <PauseCircleOutlined />
            暂停
          </Button>
          <Button 
            @click="controlJob('resume')"
            :disabled="!canResume(job.status)"
            :loading="actionLoading"
          >
            <PlayCircleOutlined />
            恢复
          </Button>
          <Button 
            danger
            @click="controlJob('stop')"
            :disabled="!canStop(job.status)"
            :loading="actionLoading"
          >
            <StopOutlined />
            停止
          </Button>
          <Button @click="cloneJob" :loading="actionLoading">
            <CopyOutlined />
            克隆
          </Button>
        </Space>
      </div>
    </div>
  </Drawer>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue';
import {
  Drawer,
  Card,
  Descriptions,
  Tag,
  Avatar,
  Row,
  Col,
  Progress,
  Divider,
  Statistic,
  Table,
  Button,
  Space,
  message,
} from 'ant-design-vue';
import {
  FileTextOutlined,
  PauseCircleOutlined,
  PlayCircleOutlined,
  StopOutlined,
  CopyOutlined,
} from '@ant-design/icons-vue';
import type { TrainingJob } from '#/api/types';
import { controlTrainingJob, cloneTrainingJob } from '#/api';
import { formatDateTime, formatDuration, formatRelativeTime } from '#/utils/date';

const props = defineProps<{
  visible: boolean;
  job: TrainingJob | null;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
  'view-logs': [job: TrainingJob];
  'refresh': [];
}>();

const actionLoading = ref(false);

// 模拟Pod数据
const mockPods = ref([
  {
    name: 'job-worker-0',
    status: 'Running',
    cpu: '2 cores',
    memory: '4 GB',
    gpu: '1 card',
    node: 'node-1',
    startTime: '2024-01-20 10:30:00',
  },
  {
    name: 'job-worker-1', 
    status: 'Running',
    cpu: '2 cores',
    memory: '4 GB',
    gpu: '1 card',
    node: 'node-2',
    startTime: '2024-01-20 10:30:00',
  },
]);

// Pod表格列定义
const podColumns = [
  {
    title: 'Pod名称',
    dataIndex: 'name',
    key: 'name',
    width: 120,
  },
  {
    title: '状态',
    key: 'podStatus',
    slots: { customRender: 'podStatus' },
    width: 80,
  },
  {
    title: '资源',
    key: 'resources',
    slots: { customRender: 'resources' },
    width: 100,
  },
  {
    title: '节点',
    dataIndex: 'node',
    key: 'node',
    width: 80,
  },
  {
    title: '启动时间',
    dataIndex: 'startTime',
    key: 'startTime',
    width: 120,
  },
];

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

const getPriorityColor = (priority: string) => {
  const colors = {
    urgent: 'red',
    high: 'orange',
    medium: 'blue',
    low: 'default',
  };
  return colors[priority as keyof typeof colors] || 'default';
};

const getPriorityLabel = (priority: string) => {
  const labels = {
    urgent: '紧急',
    high: '高',
    medium: '中',
    low: '低',
  };
  return labels[priority as keyof typeof labels] || priority;
};

const getProgressColor = (progress: number) => {
  if (progress >= 80) return '#52c41a';
  if (progress >= 60) return '#1890ff';
  if (progress >= 40) return '#faad14';
  return '#ff4d4f';
};

const getProgressStatus = (status: string) => {
  if (status === 'failed') return 'exception';
  if (status === 'completed') return 'success';
  return 'active';
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

const canPause = (status: string) => status === 'running';
const canResume = (status: string) => status === 'paused';
const canStop = (status: string) => ['running', 'queued', 'pending'].includes(status);

// 事件处理
const viewLogs = () => {
  if (props.job) {
    emit('view-logs', props.job);
  }
};

const controlJob = async (action: string) => {
  if (!props.job) return;
  
  try {
    actionLoading.value = true;
    await controlTrainingJob({ id: props.job.id, action });
    
    const actionNames = {
      pause: '暂停',
      resume: '恢复',
      stop: '停止',
    };
    
    message.success(`${actionNames[action as keyof typeof actionNames]}任务成功`);
    emit('refresh');
  } catch (error) {
    message.error('操作失败');
  } finally {
    actionLoading.value = false;
  }
};

const cloneJob = async () => {
  if (!props.job) return;
  
  try {
    actionLoading.value = true;
    await cloneTrainingJob(props.job.id, { name: `${props.job.name}_copy` });
    message.success('克隆任务成功');
    emit('refresh');
  } catch (error) {
    message.error('克隆失败');
  } finally {
    actionLoading.value = false;
  }
};
</script>

<style scoped lang="scss">
.job-detail-drawer {
  :deep(.ant-drawer-body) {
    padding: 0;
  }
}

.drawer-content {
  height: 100%;
  overflow-y: auto;
  padding: 24px;
}

.info-card {
  margin-bottom: 16px;
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
    border-radius: 8px 8px 0 0;
  }
  
  :deep(.ant-card-head-title) {
    font-weight: 600;
    color: #1890ff;
  }
}

.job-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.job-id {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 12px;
}

.status-tag {
  border-radius: 6px;
  font-weight: 500;
}

.resource-item {
  margin-bottom: 16px;
  
  .resource-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    
    .resource-label {
      font-size: 14px;
      color: #666;
    }
    
    .resource-value {
      font-weight: 600;
      color: #333;
    }
  }
}

.config-value {
  background: #f5f5f5;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 12px;
  display: inline-block;
  max-width: 100%;
  word-break: break-all;
}

.status-info {
  margin-bottom: 24px;
}

.pod-status {
  h4 {
    margin-bottom: 16px;
    color: #1890ff;
    font-weight: 600;
  }
}

.pod-table {
  :deep(.ant-table-thead > tr > th) {
    background: #fafafa;
    font-weight: 600;
  }
}

.pod-resources {
  font-size: 12px;
  
  div {
    margin-bottom: 2px;
  }
}

.action-buttons {
  position: sticky;
  bottom: 0;
  background: white;
  padding: 16px 0;
  border-top: 1px solid #f0f0f0;
  margin-top: 24px;
}

// 响应式设计
@media (max-width: 768px) {
  .drawer-content {
    padding: 16px;
  }
  
  .action-buttons {
    padding: 12px 0;
    
    :deep(.ant-space) {
      width: 100%;
      justify-content: center;
    }
    
    :deep(.ant-btn) {
      min-width: 80px;
    }
  }
}
</style>