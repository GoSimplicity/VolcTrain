<template>
  <Drawer
    v-model:open="visible"
    title="GPU设备详情"
    width="800"
    placement="right"
    class="gpu-detail-drawer"
  >
    <div v-if="gpu" class="drawer-content">
      <!-- 基本信息卡片 -->
      <Card title="基本信息" class="info-card">
        <Descriptions :column="2" bordered>
          <Descriptions.Item label="设备名称">
            <div class="gpu-name">
              {{ gpu.name }}
              <Tag :color="getGPUStatusColor(gpu.status)" class="status-tag">
                {{ getGPUStatusLabel(gpu.status) }}
              </Tag>
            </div>
          </Descriptions.Item>
          <Descriptions.Item label="设备ID">
            <code class="gpu-id">{{ gpu.id }}</code>
          </Descriptions.Item>
          <Descriptions.Item label="品牌型号">
            <Tag color="blue">{{ gpu.brand.toUpperCase() }} {{ gpu.model }}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="架构">
            {{ gpu.architecture }}
          </Descriptions.Item>
          <Descriptions.Item label="所属节点">
            {{ gpu.nodeName }}
          </Descriptions.Item>
          <Descriptions.Item label="所属集群">
            {{ gpu.clusterName }}
          </Descriptions.Item>
          <Descriptions.Item label="CUDA核心数">
            {{ gpu.cudaCores.toLocaleString() }}
          </Descriptions.Item>
          <Descriptions.Item label="显存容量">
            {{ gpu.memorySize }} GB
          </Descriptions.Item>
          <Descriptions.Item label="驱动版本" :span="2">
            {{ gpu.driverVersion }}
          </Descriptions.Item>
          <Descriptions.Item label="CUDA版本" :span="2">
            {{ gpu.cudaVersion }}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <!-- 硬件规格卡片 -->
      <Card title="硬件规格" class="info-card">
        <Row :gutter="16">
          <Col :span="12">
            <div class="spec-item">
              <div class="spec-header">
                <span class="spec-label">显存带宽</span>
                <span class="spec-value">{{ gpu.memoryBandwidth }} GB/s</span>
              </div>
            </div>
          </Col>
          <Col :span="12">
            <div class="spec-item">
              <div class="spec-header">
                <span class="spec-label">基础频率</span>
                <span class="spec-value">{{ gpu.baseClockRate }} MHz</span>
              </div>
            </div>
          </Col>
          <Col :span="12">
            <div class="spec-item">
              <div class="spec-header">
                <span class="spec-label">加速频率</span>
                <span class="spec-value">{{ gpu.boostClockRate }} MHz</span>
              </div>
            </div>
          </Col>
          <Col :span="12">
            <div class="spec-item">
              <div class="spec-header">
                <span class="spec-label">最大功耗</span>
                <span class="spec-value">{{ gpu.maxPower }} W</span>
              </div>
            </div>
          </Col>
        </Row>
      </Card>

      <!-- 实时状态卡片 -->
      <Card title="实时状态" class="info-card">
        <Row :gutter="24">
          <Col :span="12">
            <div class="status-item">
              <div class="status-header">
                <span class="status-label">GPU使用率</span>
                <span class="status-value">{{ gpu.gpuUtilization }}%</span>
              </div>
              <Progress
                :percent="gpu.gpuUtilization"
                :stroke-color="getUtilizationColor(gpu.gpuUtilization)"
                style="margin-top: 8px"
              />
            </div>
          </Col>
          <Col :span="12">
            <div class="status-item">
              <div class="status-header">
                <span class="status-label">显存使用率</span>
                <span class="status-value">{{ gpu.memoryUsage }}%</span>
              </div>
              <Progress
                :percent="gpu.memoryUsage"
                stroke-color="#722ed1"
                style="margin-top: 8px"
              />
            </div>
          </Col>
          <Col :span="12">
            <div class="status-item">
              <div class="status-header">
                <span class="status-label">当前温度</span>
                <Tag :color="getTemperatureColor(gpu.temperature)" class="temperature-tag">
                  {{ gpu.temperature }}°C
                </Tag>
              </div>
            </div>
          </Col>
          <Col :span="12">
            <div class="status-item">
              <div class="status-header">
                <span class="status-label">当前功耗</span>
                <span class="power-display">
                  {{ gpu.powerUsage }}W / {{ gpu.maxPower }}W
                </span>
              </div>
              <Progress
                :percent="(gpu.powerUsage / gpu.maxPower) * 100"
                stroke-color="#faad14"
                style="margin-top: 8px"
              />
            </div>
          </Col>
        </Row>

        <Divider />

        <div class="heartbeat-info">
          <Row :gutter="16">
            <Col :span="8">
              <Statistic
                title="创建时间"
                :value="formatDateTime(gpu.createTime, 'YYYY-MM-DD HH:mm:ss')"
                :value-style="{ fontSize: '14px' }"
              />
            </Col>
            <Col :span="8">
              <Statistic
                title="最后更新"
                :value="formatDateTime(gpu.updateTime, 'YYYY-MM-DD HH:mm:ss')"
                :value-style="{ fontSize: '14px' }"
              />
            </Col>
            <Col :span="8">
              <Statistic
                title="最后心跳"
                :value="gpu.lastHeartbeat ? formatRelativeTime(gpu.lastHeartbeat) : '无'"
                :value-style="{ fontSize: '14px', color: gpu.lastHeartbeat ? '#52c41a' : '#ff4d4f' }"
              />
            </Col>
          </Row>
        </div>
      </Card>

      <!-- 分配信息卡片 -->
      <Card v-if="gpu.allocatedTo" title="分配信息" class="info-card">
        <Descriptions :column="1" bordered>
          <Descriptions.Item label="分配给用户">
            <Avatar size="small">{{ gpu.allocatedUser?.[0] }}</Avatar>
            <span style="margin-left: 8px">{{ gpu.allocatedUser }}</span>
          </Descriptions.Item>
          <Descriptions.Item label="分配时间">
            {{ formatDateTime(gpu.allocatedTime, 'YYYY-MM-DD HH:mm:ss') }}
          </Descriptions.Item>
          <Descriptions.Item label="分配任务">
            <Button type="link">{{ gpu.allocatedTo }}</Button>
          </Descriptions.Item>
          <Descriptions.Item label="使用时长">
            {{ formatDuration(getDurationSinceAllocation(gpu.allocatedTime)) }}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <Space>
          <Button 
            type="primary" 
            @click="viewMetrics"
            :loading="actionLoading"
          >
            <ThunderboltOutlined />
            查看监控
          </Button>
          <Button 
            @click="releaseGPU"
            :disabled="gpu.status !== 'allocated'"
            :loading="actionLoading"
            danger
          >
            <StopOutlined />
            释放GPU
          </Button>
          <Button 
            @click="toggleMaintenance"
            :loading="actionLoading"
          >
            <ToolOutlined />
            {{ gpu.status === 'maintenance' ? '退出维护' : '进入维护' }}
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
  Row,
  Col,
  Progress,
  Divider,
  Statistic,
  Avatar,
  Button,
  Space,
  Modal,
  message,
} from 'ant-design-vue';
import {
  ThunderboltOutlined,
  StopOutlined,
  ToolOutlined,
} from '@ant-design/icons-vue';
import type { GPUDevice } from '#/api/types';
import { releaseGPU as releaseGPUApi, maintainGPU } from '#/api';
import { formatDateTime, formatDuration, formatRelativeTime } from '#/utils/date';

const props = defineProps<{
  visible: boolean;
  gpu: GPUDevice | null;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
  'view-metrics': [gpu: GPUDevice];
  'refresh': [];
}>();

const actionLoading = ref(false);

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

const getTemperatureColor = (temp: number) => {
  if (temp >= 85) return 'red';
  if (temp >= 75) return 'orange';
  if (temp >= 65) return 'gold';
  return 'green';
};

const getUtilizationColor = (util: number) => {
  if (util >= 90) return '#f5222d';
  if (util >= 70) return '#fa8c16';
  if (util >= 40) return '#52c41a';
  return '#1890ff';
};

const getDurationSinceAllocation = (allocatedTime?: string) => {
  if (!allocatedTime) return 0;
  const now = new Date();
  const allocated = new Date(allocatedTime);
  return Math.floor((now.getTime() - allocated.getTime()) / 1000);
};

// 事件处理
const viewMetrics = () => {
  if (props.gpu) {
    emit('view-metrics', props.gpu);
  }
};

const releaseGPU = async () => {
  if (!props.gpu) return;
  
  Modal.confirm({
    title: '确认释放GPU',
    content: `确定要释放GPU "${props.gpu.name}" 吗？这将终止当前任务。`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        actionLoading.value = true;
        await releaseGPUApi(props.gpu!.id);
        message.success('GPU释放成功');
        emit('refresh');
      } catch (error) {
        message.error('GPU释放失败');
      } finally {
        actionLoading.value = false;
      }
    },
  });
};

const toggleMaintenance = async () => {
  if (!props.gpu) return;
  
  const isMaintenance = props.gpu.status === 'maintenance';
  const action = isMaintenance ? '退出维护' : '进入维护';
  
  Modal.confirm({
    title: `确认${action}`,
    content: `确定要让GPU "${props.gpu.name}" ${action}模式吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        actionLoading.value = true;
        await maintainGPU(props.gpu!.id, !isMaintenance);
        message.success(`${action}成功`);
        emit('refresh');
      } catch (error) {
        message.error(`${action}失败`);
      } finally {
        actionLoading.value = false;
      }
    },
  });
};
</script>

<style scoped lang="scss">
.gpu-detail-drawer {
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

.gpu-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.gpu-id {
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

.spec-item {
  margin-bottom: 16px;
  
  .spec-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .spec-label {
      font-size: 14px;
      color: #666;
    }
    
    .spec-value {
      font-weight: 600;
      color: #333;
    }
  }
}

.status-item {
  margin-bottom: 16px;
  
  .status-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    
    .status-label {
      font-size: 14px;
      color: #666;
    }
    
    .status-value {
      font-weight: 600;
      color: #333;
    }
  }
}

.temperature-tag {
  font-weight: 600;
}

.power-display {
  font-weight: 600;
  color: #333;
}

.heartbeat-info {
  margin-top: 16px;
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