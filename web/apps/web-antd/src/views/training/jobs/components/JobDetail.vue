<template>
  <div class="job-detail">
    <!-- 任务基本信息 -->
    <Card class="detail-card" title="任务信息">
      <Descriptions :column="2" bordered>
        <Descriptions.Item label="任务名称">
          {{ job.name }}
        </Descriptions.Item>
        <Descriptions.Item label="任务状态">
          <Badge
            :status="getStatusBadgeType(job.status)"
            :text="getStatusText(job.status)"
          />
        </Descriptions.Item>
        <Descriptions.Item label="训练框架">
          <Tag :color="getFrameworkColor(job.framework)">
            {{ getFrameworkText(job.framework) }}
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="优先级">
          <Tag :color="getPriorityColor(job.priority)">
            {{ getPriorityText(job.priority) }}
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="训练队列">
          {{ job.queueId }}
        </Descriptions.Item>
        <Descriptions.Item label="用户ID">
          {{ job.userId }}
        </Descriptions.Item>
        <Descriptions.Item label="工作空间">
          {{ job.workspaceId }}
        </Descriptions.Item>
        <Descriptions.Item label="进度">
          <Progress 
            :percent="job.progress" 
            :status="job.status === 'failed' ? 'exception' : 
                    job.status === 'completed' ? 'success' : 'normal'"
          />
        </Descriptions.Item>
      </Descriptions>
    </Card>

    <!-- 资源配置 -->
    <Card class="detail-card" title="资源配置" style="margin-top: 16px;">
      <Descriptions :column="2" bordered>
        <Descriptions.Item label="CPU">
          <div class="resource-item">
            <CpuIcon class="resource-icon" />
            {{ job.resourceConfig.cpu }}
          </div>
        </Descriptions.Item>
        <Descriptions.Item label="内存">
          <div class="resource-item">
            <MemoryStickIcon class="resource-icon" />
            {{ job.resourceConfig.memory }}
          </div>
        </Descriptions.Item>
        <Descriptions.Item label="GPU 类型">
          <div class="resource-item">
            <ZapIcon class="resource-icon" />
            {{ job.resourceConfig.gpu.type }}
          </div>
        </Descriptions.Item>
        <Descriptions.Item label="GPU 数量">
          {{ job.resourceConfig.gpu.count }}
        </Descriptions.Item>
        <Descriptions.Item label="副本数">
          {{ job.resourceConfig.replicas }}
        </Descriptions.Item>
        <Descriptions.Item label="最小可用实例">
          {{ job.resourceConfig.minAvailable || '-' }}
        </Descriptions.Item>
      </Descriptions>
    </Card>

    <!-- 运行配置 -->
    <Card class="detail-card" title="运行配置" style="margin-top: 16px;">
      <Descriptions :column="1" bordered>
        <Descriptions.Item label="镜像地址">
          <code>{{ job.image }}</code>
        </Descriptions.Item>
        <Descriptions.Item label="启动命令">
          <code>{{ job.command?.join(' ') || '-' }}</code>
        </Descriptions.Item>
        <Descriptions.Item label="工作目录">
          <code>{{ job.workingDir || '-' }}</code>
        </Descriptions.Item>
        <Descriptions.Item label="环境变量">
          <div v-if="job.env && Object.keys(job.env).length > 0" class="env-vars">
            <div 
              v-for="(value, key) in job.env" 
              :key="key" 
              class="env-var"
            >
              <code>{{ key }}={{ value }}</code>
            </div>
          </div>
          <span v-else>-</span>
        </Descriptions.Item>
      </Descriptions>
    </Card>

    <!-- 时间信息 -->
    <Card class="detail-card" title="时间信息" style="margin-top: 16px;">
      <Descriptions :column="2" bordered>
        <Descriptions.Item label="创建时间">
          <div class="time-item">
            <CalendarIcon class="time-icon" />
            {{ formatDateTime(job.createdAt) }}
          </div>
        </Descriptions.Item>
        <Descriptions.Item label="开始时间">
          <div v-if="job.startedAt" class="time-item">
            <PlayIcon class="time-icon" />
            {{ formatDateTime(job.startedAt) }}
          </div>
          <span v-else>-</span>
        </Descriptions.Item>
        <Descriptions.Item label="完成时间">
          <div v-if="job.completedAt" class="time-item">
            <CheckCircleIcon class="time-icon" />
            {{ formatDateTime(job.completedAt) }}
          </div>
          <span v-else>-</span>
        </Descriptions.Item>
        <Descriptions.Item label="预计时长">
          <div v-if="job.estimatedDuration" class="time-item">
            <ClockIcon class="time-icon" />
            {{ formatDuration(job.estimatedDuration) }}
          </div>
          <span v-else>-</span>
        </Descriptions.Item>
      </Descriptions>
    </Card>

    <!-- 训练指标 -->
    <Card 
      v-if="job.metrics" 
      class="detail-card" 
      title="训练指标" 
      style="margin-top: 16px;"
    >
      <Descriptions :column="2" bordered>
        <Descriptions.Item label="当前 Epoch">
          {{ job.metrics.epoch }}
        </Descriptions.Item>
        <Descriptions.Item label="损失值">
          {{ job.metrics.loss.toFixed(6) }}
        </Descriptions.Item>
        <Descriptions.Item label="准确率" v-if="job.metrics.accuracy">
          {{ (job.metrics.accuracy * 100).toFixed(2) }}%
        </Descriptions.Item>
        <Descriptions.Item label="学习率">
          {{ job.metrics.learningRate.toExponential(2) }}
        </Descriptions.Item>
      </Descriptions>

      <!-- 自定义指标 -->
      <div v-if="job.metrics.customMetrics" class="custom-metrics">
        <h4>自定义指标</h4>
        <Descriptions :column="2" bordered>
          <Descriptions.Item 
            v-for="(value, key) in job.metrics.customMetrics"
            :key="key"
            :label="key"
          >
            {{ typeof value === 'number' ? value.toFixed(4) : value }}
          </Descriptions.Item>
        </Descriptions>
      </div>
    </Card>

    <!-- Volcano 任务规范 -->
    <Card 
      v-if="job.volcanoSpec" 
      class="detail-card" 
      title="Volcano 配置" 
      style="margin-top: 16px;"
    >
      <Descriptions :column="2" bordered>
        <Descriptions.Item label="命名空间">
          {{ job.volcanoSpec.namespace }}
        </Descriptions.Item>
        <Descriptions.Item label="队列">
          {{ job.volcanoSpec.queue }}
        </Descriptions.Item>
        <Descriptions.Item label="调度器">
          {{ job.volcanoSpec.schedulerName }}
        </Descriptions.Item>
        <Descriptions.Item label="优先级类">
          {{ job.volcanoSpec.priorityClassName || '-' }}
        </Descriptions.Item>
        <Descriptions.Item label="最小可用实例">
          {{ job.volcanoSpec.minAvailable }}
        </Descriptions.Item>
        <Descriptions.Item label="最大重试">
          {{ job.volcanoSpec.maxRetry || '-' }}
        </Descriptions.Item>
        <Descriptions.Item label="任务类型">
          {{ job.volcanoSpec.jobType }}
        </Descriptions.Item>
      </Descriptions>

      <!-- 任务列表 -->
      <div v-if="job.volcanoSpec.tasks" class="volcano-tasks">
        <h4>任务列表</h4>
        <Table 
          :columns="taskColumns"
          :data-source="job.volcanoSpec.tasks"
          :pagination="false"
          size="small"
        >
          <template #resources="{ record }">
            <div class="task-resources">
              <div v-if="record.template?.spec?.containers?.[0]?.resources">
                <div class="resource-item">
                  CPU: {{ record.template.spec.containers[0].resources.requests?.cpu || '-' }}
                </div>
                <div class="resource-item">
                  内存: {{ record.template.spec.containers[0].resources.requests?.memory || '-' }}
                </div>
              </div>
            </div>
          </template>
        </Table>
      </div>
    </Card>

    <!-- 操作按钮 -->
    <div class="action-buttons" style="margin-top: 24px;">
      <Space>
        <Button @click="$emit('close')">关闭</Button>
        <Button type="primary" @click="viewLogs">查看日志</Button>
        <Button @click="viewMetrics">查看监控</Button>
        <Button 
          v-if="['running', 'pending'].includes(job.status)"
          danger
          @click="stopJob"
        >
          停止任务
        </Button>
        <Button 
          v-if="job.status === 'suspended'"
          @click="resumeJob"
        >
          恢复任务
        </Button>
      </Space>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import {
  Card,
  Descriptions,
  Badge,
  Tag,
  Progress,
  Table,
  Button,
  Space,
  message,
} from 'ant-design-vue';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

// 图标
import {
  Cpu as CpuIcon,
  MemoryStick as MemoryStickIcon,
  Zap as ZapIcon,
  Calendar as CalendarIcon,
  Play as PlayIcon,
  CheckCircle as CheckCircleIcon,
  Clock as ClockIcon,
} from 'lucide-vue-next';

// 类型和服务
import type { TrainingApi } from '#/types/api';
import { trainingService } from '#/api/services';

interface Props {
  job: TrainingApi.TrainingJob;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  close: [];
}>();

// 任务列表列配置
const taskColumns = [
  { title: '任务名称', dataIndex: 'name', key: 'name' },
  { title: '副本数', dataIndex: 'replicas', key: 'replicas' },
  { title: '最小可用', dataIndex: 'minAvailable', key: 'minAvailable' },
  { title: '资源配置', key: 'resources', slots: { customRender: 'resources' } },
];

// 状态相关方法
const getStatusBadgeType = (status: TrainingApi.JobStatus) => {
  const statusMap: Record<TrainingApi.JobStatus, string> = {
    pending: 'default',
    running: 'processing',
    completed: 'success',
    failed: 'error',
    cancelled: 'warning',
    suspended: 'warning',
  };
  return statusMap[status] || 'default';
};

const getStatusText = (status: TrainingApi.JobStatus) => {
  const statusMap: Record<TrainingApi.JobStatus, string> = {
    pending: '等待中',
    running: '运行中',
    completed: '已完成',
    failed: '已失败',
    cancelled: '已取消',
    suspended: '已暂停',
  };
  return statusMap[status] || status;
};

const getFrameworkColor = (framework: TrainingApi.FrameworkType) => {
  const colorMap: Record<TrainingApi.FrameworkType, string> = {
    tensorflow: 'orange',
    pytorch: 'red',
    paddlepaddle: 'blue',
    mindspore: 'green',
    mpi: 'purple',
  };
  return colorMap[framework] || 'default';
};

const getFrameworkText = (framework: TrainingApi.FrameworkType) => {
  const textMap: Record<TrainingApi.FrameworkType, string> = {
    tensorflow: 'TensorFlow',
    pytorch: 'PyTorch',
    paddlepaddle: 'PaddlePaddle',
    mindspore: 'MindSpore',
    mpi: 'MPI',
  };
  return textMap[framework] || framework;
};

const getPriorityColor = (priority: TrainingApi.Priority) => {
  const colorMap: Record<TrainingApi.Priority, string> = {
    low: 'default',
    normal: 'blue',
    high: 'orange',
    urgent: 'red',
  };
  return colorMap[priority] || 'default';
};

const getPriorityText = (priority: TrainingApi.Priority) => {
  const textMap: Record<TrainingApi.Priority, string> = {
    low: '低',
    normal: '正常',
    high: '高',
    urgent: '紧急',
  };
  return textMap[priority] || priority;
};

// 时间格式化
const formatDateTime = (dateTime: string) => {
  return format(new Date(dateTime), 'yyyy-MM-dd HH:mm:ss', { locale: zhCN });
};

const formatDuration = (seconds: number) => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const remainingSeconds = seconds % 60;
  return `${hours}h ${minutes}m ${remainingSeconds}s`;
};

// 操作方法
const viewLogs = () => {
  // 发送事件给父组件处理
};

const viewMetrics = () => {
  // 发送事件给父组件处理
};

const stopJob = async () => {
  try {
    await trainingService.controlJob({
      id: props.job.id,
      action: 'stop',
    });
    message.success('任务停止成功');
  } catch (error) {
    message.error('任务停止失败');
  }
};

const resumeJob = async () => {
  try {
    await trainingService.controlJob({
      id: props.job.id,
      action: 'resume',
    });
    message.success('任务恢复成功');
  } catch (error) {
    message.error('任务恢复失败');
  }
};
</script>

<style lang="scss" scoped>
.job-detail {
  .detail-card {
    border-radius: 8px;
    border: none;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

    .resource-item {
      display: flex;
      align-items: center;
      
      .resource-icon {
        width: 16px;
        height: 16px;
        margin-right: 8px;
        color: #666;
      }
    }

    .time-item {
      display: flex;
      align-items: center;
      
      .time-icon {
        width: 16px;
        height: 16px;
        margin-right: 8px;
        color: #666;
      }
    }

    .env-vars {
      .env-var {
        margin-bottom: 4px;
        
        &:last-child {
          margin-bottom: 0;
        }
        
        code {
          background: #f6f6f6;
          padding: 2px 6px;
          border-radius: 4px;
          font-size: 12px;
        }
      }
    }

    .custom-metrics {
      margin-top: 16px;
      
      h4 {
        margin-bottom: 12px;
        color: #333;
      }
    }

    .volcano-tasks {
      margin-top: 16px;
      
      h4 {
        margin-bottom: 12px;
        color: #333;
      }

      .task-resources {
        .resource-item {
          font-size: 12px;
          margin-bottom: 2px;
        }
      }
    }
  }

  .action-buttons {
    text-align: center;
    padding: 16px 0;
    border-top: 1px solid #e8e8e8;
  }

  code {
    background: #f6f6f6;
    padding: 4px 8px;
    border-radius: 4px;
    font-family: 'JetBrains Mono', 'Consolas', monospace;
    font-size: 13px;
  }
}
</style>