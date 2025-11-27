<template>
  <div class="training-queues-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <ListOrderedIcon class="title-icon" />
            <span class="title-text">训练队列管理</span>
          </h1>
          <p class="page-description">
            <span class="description-text">管理和配置训练任务队列，控制资源分配策略</span>
          </p>
        </div>
        <div class="action-section">
          <Space>
            <Button @click="refreshQueues">
              <RefreshCwIcon :class="{ 'animate-spin': loading }" />
              刷新
            </Button>
            <Button type="primary" @click="showCreateModal">
              <PlusCircleIcon />
              创建队列
            </Button>
          </Space>
        </div>
      </div>
    </div>

    <!-- 队列统计 -->
    <div class="queues-stats">
      <Row :gutter="[16, 16]">
        <Col :span="6">
          <Card class="stat-card">
            <Statistic
              title="总队列数"
              :value="queueStats.total"
              :value-style="{ color: '#3f8600' }"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card class="stat-card">
            <Statistic
              title="活跃队列"
              :value="queueStats.active"
              :value-style="{ color: '#1890ff' }"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card class="stat-card">
            <Statistic
              title="运行任务"
              :value="queueStats.runningJobs"
              :value-style="{ color: '#722ed1' }"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card class="stat-card">
            <Statistic
              title="排队任务"
              :value="queueStats.pendingJobs"
              :value-style="{ color: '#faad14' }"
            />
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 队列列表 -->
    <Card class="queues-table-card">
      <Table
        :columns="tableColumns"
        :data-source="queuesList"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <!-- 队列名称列 -->
        <template #name="{ record }">
          <div class="queue-name-cell">
            <div class="queue-name">
              <a @click="viewQueueDetail(record.id)">{{ record.name }}</a>
            </div>
            <div class="queue-description">{{ record.description }}</div>
          </div>
        </template>

        <!-- 优先级列 -->
        <template #priority="{ record }">
          <Tag :color="getPriorityColor(record.priority)">
            优先级 {{ record.priority }}
          </Tag>
        </template>

        <!-- 资源配额列 -->
        <template #quota="{ record }">
          <div class="resource-quota">
            <div class="quota-item">
              <CpuIcon class="quota-icon" />
              {{ record.resourceQuota.cpu }}
            </div>
            <div class="quota-item">
              <MemoryStickIcon class="quota-icon" />
              {{ record.resourceQuota.memory }}
            </div>
            <div class="quota-item">
              <ZapIcon class="quota-icon" />
              {{ record.resourceQuota.gpu }} GPU
            </div>
          </div>
        </template>

        <!-- 任务统计列 -->
        <template #jobs="{ record }">
          <div class="jobs-stats">
            <div class="jobs-item">
              <span class="jobs-label">运行中:</span>
              <span class="jobs-value running">{{ record.runningJobs || 0 }}</span>
            </div>
            <div class="jobs-item">
              <span class="jobs-label">排队:</span>
              <span class="jobs-value pending">{{ record.pendingJobs || 0 }}</span>
            </div>
            <div class="jobs-item">
              <span class="jobs-label">最大:</span>
              <span class="jobs-value max">{{ record.maxRunningJobs }}</span>
            </div>
          </div>
        </template>

        <!-- 状态列 -->
        <template #status="{ record }">
          <Badge
            :status="getQueueStatusBadge(record)"
            :text="getQueueStatusText(record)"
          />
        </template>

        <!-- 操作列 -->
        <template #action="{ record }">
          <Space>
            <Tooltip title="查看详情">
              <Button size="small" @click="viewQueueDetail(record.id)">
                <EyeIcon />
              </Button>
            </Tooltip>
            <Tooltip title="队列统计">
              <Button size="small" @click="viewQueueStats(record.id)">
                <BarChart3Icon />
              </Button>
            </Tooltip>
            <Dropdown>
              <Button size="small">
                <MoreHorizontalIcon />
              </Button>
              <template #overlay>
                <Menu @click="(e) => handleQueueAction(record, e.key)">
                  <Menu.Item key="edit">
                    <EditIcon />
                    编辑
                  </Menu.Item>
                  <Menu.Item key="pause">
                    <PauseIcon />
                    暂停
                  </Menu.Item>
                  <Menu.Item key="resume">
                    <PlayIcon />
                    恢复
                  </Menu.Item>
                  <Menu.Divider />
                  <Menu.Item key="delete" danger>
                    <TrashIcon />
                    删除
                  </Menu.Item>
                </Menu>
              </template>
            </Dropdown>
          </Space>
        </template>
      </Table>
    </Card>

    <!-- 创建队列弹窗 -->
    <Modal
      v-model:open="createModalVisible"
      title="创建训练队列"
      width="600px"
      :confirm-loading="createLoading"
      @ok="handleCreateQueue"
      @cancel="cancelCreate"
    >
      <CreateQueueForm ref="createQueueFormRef" />
    </Modal>

    <!-- 队列详情弹窗 -->
    <Modal
      v-model:open="detailModalVisible"
      title="队列详情"
      width="800px"
      :footer="null"
    >
      <QueueDetail
        v-if="selectedQueue"
        :queue="selectedQueue"
        @close="detailModalVisible = false"
      />
    </Modal>

    <!-- 队列统计弹窗 -->
    <Modal
      v-model:open="statsModalVisible"
      title="队列统计"
      width="1000px"
      :footer="null"
    >
      <QueueStats
        v-if="selectedQueueId"
        :queue-id="selectedQueueId"
        @close="statsModalVisible = false"
      />
    </Modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, reactive } from 'vue';
import {
  Card,
  Table,
  Button,
  Space,
  Tag,
  Badge,
  Tooltip,
  Dropdown,
  Menu,
  Modal,
  Row,
  Col,
  Statistic,
  message,
} from 'ant-design-vue';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

// 图标组件
import {
  ListOrdered as ListOrderedIcon,
  RefreshCw as RefreshCwIcon,
  PlusCircle as PlusCircleIcon,
  Cpu as CpuIcon,
  MemoryStick as MemoryStickIcon,
  Zap as ZapIcon,
  Eye as EyeIcon,
  BarChart3 as BarChart3Icon,
  MoreHorizontal as MoreHorizontalIcon,
  Edit as EditIcon,
  Pause as PauseIcon,
  Play as PlayIcon,
  Trash as TrashIcon,
} from 'lucide-vue-next';

// API 服务和类型
import { trainingService } from '#/api/services';
import type { TrainingApi } from '#/types/api';

// 子组件
import CreateQueueForm from './components/CreateQueueForm.vue';
import QueueDetail from './components/QueueDetail.vue';
import QueueStats from './components/QueueStats.vue';

defineOptions({ name: 'TrainingQueues' });

// 响应式数据
const loading = ref(false);
const createLoading = ref(false);

// 选择相关
const selectedQueue = ref<TrainingApi.TrainingQueue | null>(null);
const selectedQueueId = ref<string | null>(null);

// 弹窗状态
const createModalVisible = ref(false);
const detailModalVisible = ref(false);
const statsModalVisible = ref(false);

// 表单引用
const createQueueFormRef = ref();

// 队列列表数据
const queuesList = ref<TrainingApi.TrainingQueue[]>([
  {
    id: '1',
    name: 'default',
    description: '默认训练队列',
    priority: 1,
    maxRunningJobs: 10,
    resourceQuota: {
      cpu: '100',
      memory: '200Gi',
      gpu: 20,
    },
    policies: [],
    createdAt: '2024-01-15T10:00:00Z',
    updatedAt: '2024-01-20T14:30:00Z',
    runningJobs: 5,
    pendingJobs: 3,
  },
  {
    id: '2',
    name: 'high-priority',
    description: '高优先级队列',
    priority: 2,
    maxRunningJobs: 5,
    resourceQuota: {
      cpu: '50',
      memory: '100Gi',
      gpu: 10,
    },
    policies: [],
    createdAt: '2024-01-16T09:00:00Z',
    updatedAt: '2024-01-20T12:15:00Z',
    runningJobs: 2,
    pendingJobs: 1,
  },
]);

// 队列统计
const queueStats = computed(() => ({
  total: queuesList.value.length,
  active: queuesList.value.filter(q => (q.runningJobs || 0) > 0).length,
  runningJobs: queuesList.value.reduce((sum, q) => sum + (q.runningJobs || 0), 0),
  pendingJobs: queuesList.value.reduce((sum, q) => sum + (q.pendingJobs || 0), 0),
}));

// 分页配置
const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number, range: [number, number]) =>
    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
});

// 表格列配置
const tableColumns = [
  {
    title: '队列名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
    slots: { customRender: 'name' },
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
    width: 100,
    slots: { customRender: 'priority' },
  },
  {
    title: '资源配额',
    key: 'quota',
    width: 200,
    slots: { customRender: 'quota' },
  },
  {
    title: '任务统计',
    key: 'jobs',
    width: 150,
    slots: { customRender: 'jobs' },
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    slots: { customRender: 'status' },
  },
  {
    title: '创建时间',
    dataIndex: 'createdAt',
    key: 'createdAt',
    width: 150,
    customRender: ({ text }: { text: string }) => formatDateTime(text),
  },
  {
    title: '操作',
    key: 'action',
    width: 150,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

/**
 * 获取优先级颜色
 */
const getPriorityColor = (priority: number) => {
  if (priority >= 3) return 'red';
  if (priority >= 2) return 'orange';
  return 'blue';
};

/**
 * 获取队列状态徽章
 */
const getQueueStatusBadge = (queue: TrainingApi.TrainingQueue & { runningJobs?: number }) => {
  if ((queue.runningJobs || 0) >= queue.maxRunningJobs) {
    return 'error'; // 满载
  }
  if ((queue.runningJobs || 0) > 0) {
    return 'processing'; // 运行中
  }
  return 'default'; // 空闲
};

/**
 * 获取队列状态文本
 */
const getQueueStatusText = (queue: TrainingApi.TrainingQueue & { runningJobs?: number }) => {
  if ((queue.runningJobs || 0) >= queue.maxRunningJobs) {
    return '满载';
  }
  if ((queue.runningJobs || 0) > 0) {
    return '运行中';
  }
  return '空闲';
};

/**
 * 格式化日期时间
 */
const formatDateTime = (dateTime: string) => {
  return format(new Date(dateTime), 'yyyy-MM-dd HH:mm', { locale: zhCN });
};

/**
 * 加载队列列表
 */
const loadQueues = async () => {
  try {
    loading.value = true;
    
    const params = {
      page: pagination.current,
      pageSize: pagination.pageSize,
    };

    // const response = await trainingService.getQueues(params);
    // queuesList.value = response.items;
    // pagination.total = response.total;
    
    // 模拟数据加载
    await new Promise(resolve => setTimeout(resolve, 500));
    pagination.total = queuesList.value.length;
  } catch (error) {
    message.error('加载队列列表失败');
  } finally {
    loading.value = false;
  }
};

/**
 * 刷新队列列表
 */
const refreshQueues = () => {
  loadQueues();
};

/**
 * 处理表格变化
 */
const handleTableChange = (paginationConfig: any) => {
  pagination.current = paginationConfig.current;
  pagination.pageSize = paginationConfig.pageSize;
  loadQueues();
};

/**
 * 显示创建队列弹窗
 */
const showCreateModal = () => {
  createModalVisible.value = true;
};

/**
 * 处理创建队列
 */
const handleCreateQueue = async () => {
  try {
    const formData = await createQueueFormRef.value.validate();
    
    createLoading.value = true;
    // await trainingService.createQueue(formData);
    
    message.success('训练队列创建成功');
    createModalVisible.value = false;
    loadQueues();
  } catch (error) {
    message.error('创建训练队列失败');
  } finally {
    createLoading.value = false;
  }
};

/**
 * 取消创建
 */
const cancelCreate = () => {
  createModalVisible.value = false;
  createQueueFormRef.value?.resetFields();
};

/**
 * 查看队列详情
 */
const viewQueueDetail = async (queueId: string) => {
  try {
    // const queue = await trainingService.getQueue(queueId);
    selectedQueue.value = queuesList.value.find(q => q.id === queueId) || null;
    detailModalVisible.value = true;
  } catch (error) {
    message.error('获取队列详情失败');
  }
};

/**
 * 查看队列统计
 */
const viewQueueStats = (queueId: string) => {
  selectedQueueId.value = queueId;
  statsModalVisible.value = true;
};

/**
 * 处理队列操作
 */
const handleQueueAction = async (queue: TrainingApi.TrainingQueue, action: string) => {
  try {
    switch (action) {
      case 'edit':
        // 打开编辑弹窗
        break;
      case 'pause':
        message.success(`队列 "${queue.name}" 已暂停`);
        break;
      case 'resume':
        message.success(`队列 "${queue.name}" 已恢复`);
        break;
      case 'delete':
        Modal.confirm({
          title: '确认删除',
          content: `确定要删除队列 "${queue.name}" 吗？`,
          onOk: async () => {
            // await trainingService.deleteQueue(queue.id);
            message.success('删除成功');
            loadQueues();
          },
        });
        break;
    }
  } catch (error) {
    message.error('操作失败');
  }
};

// 初始化
onMounted(() => {
  loadQueues();
});
</script>

<style lang="scss" scoped>
.training-queues-container {
  padding: 0;
  background: #f5f7fa;
  min-height: 100vh;
}

.page-header {
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
      font-size: 24px;
      font-weight: 600;
      margin: 0;
      
      .title-icon {
        margin-right: 12px;
        font-size: 28px;
      }
    }

    .page-description {
      margin: 8px 0 0 40px;
      font-size: 14px;
      opacity: 0.9;
    }
  }
}

.queues-stats {
  margin: 0 24px 24px 24px;

  .stat-card {
    text-align: center;
    border-radius: 8px;
    border: none;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
}

.queues-table-card {
  margin: 0 24px;
  border-radius: 8px;
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  :deep(.ant-table) {
    .ant-table-thead > tr > th {
      background: #fafafa;
      border-bottom: 1px solid #e8e8e8;
      font-weight: 600;
    }
  }
}

.queue-name-cell {
  .queue-name {
    font-weight: 500;
    margin-bottom: 4px;
    
    a {
      color: #1890ff;
      text-decoration: none;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }

  .queue-description {
    font-size: 12px;
    color: #666;
  }
}

.resource-quota {
  .quota-item {
    display: flex;
    align-items: center;
    margin-bottom: 2px;
    font-size: 12px;
    color: #666;

    .quota-icon {
      width: 14px;
      height: 14px;
      margin-right: 4px;
    }
  }
}

.jobs-stats {
  .jobs-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 2px;
    font-size: 12px;

    .jobs-label {
      color: #666;
    }

    .jobs-value {
      font-weight: 500;
      
      &.running {
        color: #52c41a;
      }
      
      &.pending {
        color: #faad14;
      }
      
      &.max {
        color: #1890ff;
      }
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