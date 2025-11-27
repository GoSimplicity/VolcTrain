<template>
  <div class="queue-detail">
    <Descriptions :column="2" bordered>
      <Descriptions.Item label="队列名称">
        {{ queue.name }}
      </Descriptions.Item>
      <Descriptions.Item label="优先级">
        <Tag :color="getPriorityColor(queue.priority)">
          优先级 {{ queue.priority }}
        </Tag>
      </Descriptions.Item>
      <Descriptions.Item label="最大运行任务">
        {{ queue.maxRunningJobs }}
      </Descriptions.Item>
      <Descriptions.Item label="创建时间">
        {{ formatDateTime(queue.createdAt) }}
      </Descriptions.Item>
      <Descriptions.Item label="更新时间">
        {{ formatDateTime(queue.updatedAt) }}
      </Descriptions.Item>
      <Descriptions.Item label="队列描述" :span="2">
        {{ queue.description || '暂无描述' }}
      </Descriptions.Item>
    </Descriptions>

    <!-- 资源配额 -->
    <Card title="资源配额" style="margin-top: 16px;">
      <Descriptions :column="3" bordered>
        <Descriptions.Item label="CPU">
          <div class="resource-item">
            <CpuIcon class="resource-icon" />
            {{ queue.resourceQuota.cpu }} 核
          </div>
        </Descriptions.Item>
        <Descriptions.Item label="内存">
          <div class="resource-item">
            <MemoryStickIcon class="resource-icon" />
            {{ queue.resourceQuota.memory }}
          </div>
        </Descriptions.Item>
        <Descriptions.Item label="GPU">
          <div class="resource-item">
            <ZapIcon class="resource-icon" />
            {{ queue.resourceQuota.gpu }} 个
          </div>
        </Descriptions.Item>
      </Descriptions>
    </Card>

    <!-- 调度策略 -->
    <Card title="调度策略" style="margin-top: 16px;">
      <div v-if="queue.policies.length > 0">
        <Table 
          :columns="policyColumns"
          :data-source="queue.policies"
          :pagination="false"
          size="small"
        >
          <template #type="{ record }">
            <Tag>{{ getPolicyTypeName(record.type) }}</Tag>
          </template>
        </Table>
      </div>
      <Empty v-else description="暂无调度策略" />
    </Card>

    <div class="detail-actions" style="margin-top: 24px; text-align: center;">
      <Space>
        <Button @click="$emit('close')">关闭</Button>
        <Button type="primary" @click="editQueue">编辑队列</Button>
      </Space>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  Descriptions,
  Tag,
  Card,
  Table,
  Empty,
  Button,
  Space,
} from 'ant-design-vue';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

// 图标
import {
  Cpu as CpuIcon,
  MemoryStick as MemoryStickIcon,
  Zap as ZapIcon,
} from 'lucide-vue-next';

// 类型
import type { TrainingApi } from '#/types/api';

interface Props {
  queue: TrainingApi.TrainingQueue;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  close: [];
}>();

// 策略表格列
const policyColumns = [
  { title: '策略类型', key: 'type', slots: { customRender: 'type' } },
  { title: '配置', dataIndex: ['config', 'description'], key: 'config' },
];

const getPriorityColor = (priority: number) => {
  if (priority >= 3) return 'red';
  if (priority >= 2) return 'orange';
  return 'blue';
};

const getPolicyTypeName = (type: string) => {
  const typeMap: Record<string, string> = {
    fairshare: '公平共享',
    priority: '优先级',
    drf: '主导资源公平',
    gang: 'Gang调度',
  };
  return typeMap[type] || type;
};

const formatDateTime = (dateTime: string) => {
  return format(new Date(dateTime), 'yyyy-MM-dd HH:mm:ss', { locale: zhCN });
};

const editQueue = () => {
  // 编辑队列逻辑
};
</script>

<style lang="scss" scoped>
.queue-detail {
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
}
</style>