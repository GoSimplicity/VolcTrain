<template>
  <div class="training-experiments-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <FlaskConicalIcon class="title-icon" />
            <span class="title-text">实验管理</span>
          </h1>
          <p class="page-description">
            <span class="description-text">管理机器学习实验，跟踪模型版本和性能指标</span>
          </p>
        </div>
        <div class="action-section">
          <Space>
            <Button @click="refreshExperiments">
              <RefreshCwIcon :class="{ 'animate-spin': loading }" />
              刷新
            </Button>
            <Button type="primary" @click="showCreateModal">
              <PlusCircleIcon />
              创建实验
            </Button>
          </Space>
        </div>
      </div>
    </div>

    <!-- 实验列表 -->
    <Card class="experiments-table-card">
      <Table
        :columns="tableColumns"
        :data-source="experimentsList"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <!-- 实验名称列 -->
        <template #name="{ record }">
          <div class="experiment-name-cell">
            <div class="experiment-name">
              <a @click="viewExperimentDetail(record.id)">{{ record.name }}</a>
            </div>
            <div class="experiment-description">{{ record.description }}</div>
          </div>
        </template>

        <!-- 状态列 -->
        <template #status="{ record }">
          <Badge
            :status="getStatusBadgeType(record.status)"
            :text="getStatusText(record.status)"
          />
        </template>

        <!-- 运行数量列 -->
        <template #runs="{ record }">
          <div class="runs-info">
            <div class="runs-count">{{ record.totalRuns || 0 }} 次运行</div>
            <div class="best-metric" v-if="record.bestMetric">
              最佳: {{ record.bestMetric.toFixed(4) }}
            </div>
          </div>
        </template>

        <!-- 时间列 -->
        <template #time="{ record }">
          <div class="time-info">
            <div class="time-item">创建: {{ formatDateTime(record.createdAt) }}</div>
            <div v-if="record.lastRunAt" class="time-item">
              最后运行: {{ formatDateTime(record.lastRunAt) }}
            </div>
          </div>
        </template>

        <!-- 操作列 -->
        <template #action="{ record }">
          <Space>
            <Tooltip title="查看详情">
              <Button size="small" @click="viewExperimentDetail(record.id)">
                <EyeIcon />
              </Button>
            </Tooltip>
            <Tooltip title="运行实验">
              <Button size="small" type="primary" @click="runExperiment(record.id)">
                <PlayIcon />
              </Button>
            </Tooltip>
            <Tooltip title="比较运行">
              <Button size="small" @click="compareRuns(record.id)">
                <GitCompareIcon />
              </Button>
            </Tooltip>
            <Dropdown>
              <Button size="small">
                <MoreHorizontalIcon />
              </Button>
              <template #overlay>
                <Menu @click="(e) => handleExperimentAction(record, e.key)">
                  <Menu.Item key="edit">
                    <EditIcon />
                    编辑
                  </Menu.Item>
                  <Menu.Item key="clone">
                    <CopyIcon />
                    克隆
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

    <!-- 创建实验弹窗 -->
    <Modal
      v-model:open="createModalVisible"
      title="创建实验"
      width="600px"
      :confirm-loading="createLoading"
      @ok="handleCreateExperiment"
      @cancel="cancelCreate"
    >
      <CreateExperimentForm ref="createExperimentFormRef" />
    </Modal>

    <!-- 实验详情弹窗 -->
    <Modal
      v-model:open="detailModalVisible"
      title="实验详情"
      width="1200px"
      :footer="null"
    >
      <ExperimentDetail
        v-if="selectedExperiment"
        :experiment="selectedExperiment"
        @close="detailModalVisible = false"
        @run-experiment="runExperiment"
      />
    </Modal>

    <!-- 运行比较弹窗 -->
    <Modal
      v-model:open="compareModalVisible"
      title="运行比较"
      width="1400px"
      :footer="null"
    >
      <RunComparison
        v-if="selectedExperimentId"
        :experiment-id="selectedExperimentId"
        @close="compareModalVisible = false"
      />
    </Modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, reactive } from 'vue';
import {
  Card,
  Table,
  Button,
  Space,
  Badge,
  Tooltip,
  Dropdown,
  Menu,
  Modal,
  message,
} from 'ant-design-vue';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

// 图标组件
import {
  FlaskConical as FlaskConicalIcon,
  RefreshCw as RefreshCwIcon,
  PlusCircle as PlusCircleIcon,
  Eye as EyeIcon,
  Play as PlayIcon,
  GitCompare as GitCompareIcon,
  MoreHorizontal as MoreHorizontalIcon,
  Edit as EditIcon,
  Copy as CopyIcon,
  Trash as TrashIcon,
} from 'lucide-vue-next';

// 子组件
import CreateExperimentForm from './components/CreateExperimentForm.vue';
import ExperimentDetail from './components/ExperimentDetail.vue';
import RunComparison from './components/RunComparison.vue';

defineOptions({ name: 'TrainingExperiments' });

// 响应式数据
const loading = ref(false);
const createLoading = ref(false);

// 选择相关
const selectedExperiment = ref<any>(null);
const selectedExperimentId = ref<string | null>(null);

// 弹窗状态
const createModalVisible = ref(false);
const detailModalVisible = ref(false);
const compareModalVisible = ref(false);

// 表单引用
const createExperimentFormRef = ref();

// 实验列表数据（模拟数据）
const experimentsList = ref([
  {
    id: '1',
    name: 'ResNet-50-ImageNet',
    description: '使用ImageNet数据集训练ResNet-50模型',
    status: 'active',
    totalRuns: 15,
    bestMetric: 0.9234,
    createdAt: '2024-01-15T10:00:00Z',
    lastRunAt: '2024-01-20T14:30:00Z',
    createdBy: 'user1',
  },
  {
    id: '2',
    name: 'BERT-Classification',
    description: '基于BERT的文本分类实验',
    status: 'completed',
    totalRuns: 8,
    bestMetric: 0.8967,
    createdAt: '2024-01-16T09:00:00Z',
    lastRunAt: '2024-01-19T11:20:00Z',
    createdBy: 'user2',
  },
]);

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
    title: '实验名称',
    dataIndex: 'name',
    key: 'name',
    width: 250,
    slots: { customRender: 'name' },
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100,
    slots: { customRender: 'status' },
  },
  {
    title: '运行统计',
    key: 'runs',
    width: 150,
    slots: { customRender: 'runs' },
  },
  {
    title: '创建者',
    dataIndex: 'createdBy',
    key: 'createdBy',
    width: 120,
  },
  {
    title: '时间信息',
    key: 'time',
    width: 200,
    slots: { customRender: 'time' },
  },
  {
    title: '操作',
    key: 'action',
    width: 200,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

/**
 * 获取状态徽章类型
 */
const getStatusBadgeType = (status: string) => {
  const statusMap: Record<string, string> = {
    active: 'processing',
    completed: 'success',
    archived: 'default',
  };
  return statusMap[status] || 'default';
};

/**
 * 获取状态文本
 */
const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    active: '活跃',
    completed: '已完成',
    archived: '已归档',
  };
  return statusMap[status] || status;
};

/**
 * 格式化日期时间
 */
const formatDateTime = (dateTime: string) => {
  return format(new Date(dateTime), 'MM-dd HH:mm', { locale: zhCN });
};

/**
 * 加载实验列表
 */
const loadExperiments = async () => {
  try {
    loading.value = true;
    
    // 模拟数据加载
    await new Promise(resolve => setTimeout(resolve, 500));
    pagination.total = experimentsList.value.length;
  } catch (error) {
    message.error('加载实验列表失败');
  } finally {
    loading.value = false;
  }
};

/**
 * 刷新实验列表
 */
const refreshExperiments = () => {
  loadExperiments();
};

/**
 * 处理表格变化
 */
const handleTableChange = (paginationConfig: any) => {
  pagination.current = paginationConfig.current;
  pagination.pageSize = paginationConfig.pageSize;
  loadExperiments();
};

/**
 * 显示创建实验弹窗
 */
const showCreateModal = () => {
  createModalVisible.value = true;
};

/**
 * 处理创建实验
 */
const handleCreateExperiment = async () => {
  try {
    const formData = await createExperimentFormRef.value.validate();
    
    createLoading.value = true;
    // await experimentService.createExperiment(formData);
    
    message.success('实验创建成功');
    createModalVisible.value = false;
    loadExperiments();
  } catch (error) {
    message.error('创建实验失败');
  } finally {
    createLoading.value = false;
  }
};

/**
 * 取消创建
 */
const cancelCreate = () => {
  createModalVisible.value = false;
  createExperimentFormRef.value?.resetFields();
};

/**
 * 查看实验详情
 */
const viewExperimentDetail = (experimentId: string) => {
  selectedExperiment.value = experimentsList.value.find(e => e.id === experimentId);
  detailModalVisible.value = true;
};

/**
 * 运行实验
 */
const runExperiment = (experimentId: string) => {
  message.success('实验运行已启动');
};

/**
 * 比较运行
 */
const compareRuns = (experimentId: string) => {
  selectedExperimentId.value = experimentId;
  compareModalVisible.value = true;
};

/**
 * 处理实验操作
 */
const handleExperimentAction = async (experiment: any, action: string) => {
  try {
    switch (action) {
      case 'edit':
        // 编辑实验
        break;
      case 'clone':
        message.success(`实验 "${experiment.name}" 已克隆`);
        loadExperiments();
        break;
      case 'delete':
        Modal.confirm({
          title: '确认删除',
          content: `确定要删除实验 "${experiment.name}" 吗？`,
          onOk: async () => {
            message.success('删除成功');
            loadExperiments();
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
  loadExperiments();
});
</script>

<style lang="scss" scoped>
.training-experiments-container {
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

.experiments-table-card {
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

.experiment-name-cell {
  .experiment-name {
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

  .experiment-description {
    font-size: 12px;
    color: #666;
  }
}

.runs-info {
  .runs-count {
    font-weight: 500;
    margin-bottom: 2px;
  }

  .best-metric {
    font-size: 12px;
    color: #52c41a;
  }
}

.time-info {
  .time-item {
    margin-bottom: 2px;
    font-size: 12px;
    color: #666;
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