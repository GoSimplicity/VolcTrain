<template>
  <div class="training-jobs-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <BrainCircuitIcon class="title-icon" />
            <span class="title-text">训练任务管理</span>
          </h1>
          <p class="page-description">
            <span class="description-text">创建、管理和监控深度学习训练任务</span>
          </p>
        </div>
        <div class="action-section">
          <Space>
            <Button @click="refreshJobs">
              <RefreshCwIcon :class="{ 'animate-spin': loading }" />
              刷新
            </Button>
            <Button type="primary" @click="showCreateModal">
              <PlusCircleIcon />
              创建训练任务
            </Button>
          </Space>
        </div>
      </div>
    </div>

    <!-- 筛选和搜索 -->
    <Card class="filter-card">
      <div class="filter-content">
        <Row :gutter="16" align="middle">
          <Col :span="6">
            <Input.Search
              v-model:value="searchKeyword"
              placeholder="搜索任务名称或描述"
              @search="handleSearch"
              allow-clear
            />
          </Col>
          <Col :span="4">
            <Select
              v-model:value="statusFilter"
              placeholder="任务状态"
              style="width: 100%"
              allow-clear
              @change="handleFilter"
            >
              <Select.Option value="pending">等待中</Select.Option>
              <Select.Option value="running">运行中</Select.Option>
              <Select.Option value="completed">已完成</Select.Option>
              <Select.Option value="failed">已失败</Select.Option>
              <Select.Option value="cancelled">已取消</Select.Option>
              <Select.Option value="suspended">已暂停</Select.Option>
            </Select>
          </Col>
          <Col :span="4">
            <Select
              v-model:value="frameworkFilter"
              placeholder="训练框架"
              style="width: 100%"
              allow-clear
              @change="handleFilter"
            >
              <Select.Option value="tensorflow">TensorFlow</Select.Option>
              <Select.Option value="pytorch">PyTorch</Select.Option>
              <Select.Option value="paddlepaddle">PaddlePaddle</Select.Option>
              <Select.Option value="mindspore">MindSpore</Select.Option>
              <Select.Option value="mpi">MPI</Select.Option>
            </Select>
          </Col>
          <Col :span="4">
            <Select
              v-model:value="queueFilter"
              placeholder="训练队列"
              style="width: 100%"
              allow-clear
              @change="handleFilter"
            >
              <Select.Option
                v-for="queue in availableQueues"
                :key="queue.id"
                :value="queue.id"
              >
                {{ queue.name }}
              </Select.Option>
            </Select>
          </Col>
          <Col :span="6">
            <Space>
              <Button @click="handleFilter">
                <FilterIcon />
                筛选
              </Button>
              <Button @click="resetFilters">
                <XIcon />
                重置
              </Button>
              <Dropdown>
                <Button>
                  <MoreHorizontalIcon />
                  批量操作
                  <DownIcon />
                </Button>
                <template #overlay>
                  <Menu @click="handleBatchAction">
                    <Menu.Item key="start">批量启动</Menu.Item>
                    <Menu.Item key="stop">批量停止</Menu.Item>
                    <Menu.Item key="delete" danger>批量删除</Menu.Item>
                  </Menu>
                </template>
              </Dropdown>
            </Space>
          </Col>
        </Row>
      </div>
    </Card>

    <!-- 任务列表 -->
    <Card class="jobs-table-card">
      <Table
        v-model:selectedRowKeys="selectedJobKeys"
        :columns="tableColumns"
        :data-source="jobsList"
        :loading="loading"
        :pagination="pagination"
        :row-selection="{ type: 'checkbox', onChange: handleSelectionChange }"
        row-key="id"
        @change="handleTableChange"
      >
        <!-- 任务名称列 -->
        <template #name="{ record }">
          <div class="job-name-cell">
            <div class="job-name">
              <a @click="viewJobDetail(record.id)">{{ record.name }}</a>
            </div>
            <div class="job-description">{{ record.description }}</div>
          </div>
        </template>

        <!-- 状态列 -->
        <template #status="{ record }">
          <Badge
            :status="getStatusBadgeType(record.status)"
            :text="getStatusText(record.status)"
          />
        </template>

        <!-- 框架列 -->
        <template #framework="{ record }">
          <Tag :color="getFrameworkColor(record.framework)">
            {{ getFrameworkText(record.framework) }}
          </Tag>
        </template>

        <!-- 资源配置列 -->
        <template #resources="{ record }">
          <div class="resource-info">
            <div class="resource-item">
              <CpuIcon class="resource-icon" />
              {{ record.resourceConfig.cpu }}
            </div>
            <div class="resource-item">
              <MemoryStickIcon class="resource-icon" />
              {{ record.resourceConfig.memory }}
            </div>
            <div class="resource-item">
              <ZapIcon class="resource-icon" />
              {{ record.resourceConfig.gpu.count }}x {{ record.resourceConfig.gpu.type }}
            </div>
          </div>
        </template>

        <!-- 进度列 -->
        <template #progress="{ record }">
          <div class="progress-cell">
            <Progress
              :percent="record.progress"
              size="small"
              :status="record.status === 'failed' ? 'exception' : 
                      record.status === 'completed' ? 'success' : 'normal'"
            />
            <div class="progress-text">{{ record.progress }}%</div>
          </div>
        </template>

        <!-- 时间列 -->
        <template #time="{ record }">
          <div class="time-info">
            <div class="time-item">
              <CalendarIcon class="time-icon" />
              创建: {{ formatDateTime(record.createdAt) }}
            </div>
            <div v-if="record.startedAt" class="time-item">
              <PlayIcon class="time-icon" />
              开始: {{ formatDateTime(record.startedAt) }}
            </div>
            <div v-if="record.completedAt" class="time-item">
              <CheckCircleIcon class="time-icon" />
              完成: {{ formatDateTime(record.completedAt) }}
            </div>
            <div v-if="record.estimatedDuration" class="time-item">
              <ClockIcon class="time-icon" />
              预计: {{ formatDuration(record.estimatedDuration) }}
            </div>
          </div>
        </template>

        <!-- 操作列 -->
        <template #action="{ record }">
          <Space>
            <Tooltip title="查看详情">
              <Button size="small" @click="viewJobDetail(record.id)">
                <EyeIcon />
              </Button>
            </Tooltip>
            <Tooltip title="查看日志">
              <Button size="small" @click="viewJobLogs(record.id)">
                <FileTextIcon />
              </Button>
            </Tooltip>
            <Tooltip title="查看监控">
              <Button size="small" @click="viewJobMetrics(record.id)">
                <BarChart3Icon />
              </Button>
            </Tooltip>
            <Dropdown>
              <Button size="small">
                <MoreHorizontalIcon />
              </Button>
              <template #overlay>
                <Menu @click="(e) => handleJobAction(record, e.key)">
                  <Menu.Item
                    v-if="record.status === 'pending'"
                    key="start"
                  >
                    <PlayIcon />
                    启动
                  </Menu.Item>
                  <Menu.Item
                    v-if="['running', 'pending'].includes(record.status)"
                    key="stop"
                  >
                    <StopCircleIcon />
                    停止
                  </Menu.Item>
                  <Menu.Item
                    v-if="record.status === 'suspended'"
                    key="resume"
                  >
                    <PlayIcon />
                    恢复
                  </Menu.Item>
                  <Menu.Item
                    v-if="record.status === 'running'"
                    key="suspend"
                  >
                    <PauseIcon />
                    暂停
                  </Menu.Item>
                  <Menu.Item key="restart">
                    <RotateCcwIcon />
                    重启
                  </Menu.Item>
                  <Menu.Item key="clone">
                    <CopyIcon />
                    克隆
                  </Menu.Item>
                  <Menu.Divider />
                  <Menu.Item key="edit">
                    <EditIcon />
                    编辑
                  </Menu.Item>
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

    <!-- 创建任务弹窗 -->
    <Modal
      v-model:open="createModalVisible"
      title="创建训练任务"
      width="800px"
      :confirm-loading="createLoading"
      @ok="handleCreateJob"
      @cancel="cancelCreate"
    >
      <CreateJobForm
        ref="createJobFormRef"
        :queues="availableQueues"
      />
    </Modal>

    <!-- 任务详情弹窗 -->
    <Modal
      v-model:open="detailModalVisible"
      title="任务详情"
      width="1200px"
      :footer="null"
    >
      <JobDetail
        v-if="selectedJob"
        :job="selectedJob"
        @close="detailModalVisible = false"
      />
    </Modal>

    <!-- 任务日志弹窗 -->
    <Modal
      v-model:open="logsModalVisible"
      title="任务日志"
      width="1000px"
      :footer="null"
    >
      <JobLogs
        v-if="selectedJobId"
        :job-id="selectedJobId"
        @close="logsModalVisible = false"
      />
    </Modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import {
  Card,
  Table,
  Button,
  Space,
  Input,
  Select,
  Badge,
  Tag,
  Progress,
  Tooltip,
  Dropdown,
  Menu,
  Modal,
  Row,
  Col,
  message,
} from 'ant-design-vue';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

// 图标组件
import {
  BrainCircuit as BrainCircuitIcon,
  RefreshCw as RefreshCwIcon,
  PlusCircle as PlusCircleIcon,
  Filter as FilterIcon,
  X as XIcon,
  MoreHorizontal as MoreHorizontalIcon,
  ChevronDown as DownIcon,
  Cpu as CpuIcon,
  MemoryStick as MemoryStickIcon,
  Zap as ZapIcon,
  Calendar as CalendarIcon,
  Play as PlayIcon,
  CheckCircle as CheckCircleIcon,
  Clock as ClockIcon,
  Eye as EyeIcon,
  FileText as FileTextIcon,
  BarChart3 as BarChart3Icon,
  StopCircle as StopCircleIcon,
  Pause as PauseIcon,
  RotateCcw as RotateCcwIcon,
  Copy as CopyIcon,
  Edit as EditIcon,
  Trash as TrashIcon,
} from 'lucide-vue-next';

// API 服务和类型
import { trainingService } from '#/api/services';
import type { TrainingApi } from '#/types/api';

// 子组件
import CreateJobForm from './components/CreateJobForm.vue';
import JobDetail from './components/JobDetail.vue';
import JobLogs from './components/JobLogs.vue';

defineOptions({ name: 'TrainingJobs' });

const router = useRouter();

// 响应式数据
const loading = ref(false);
const createLoading = ref(false);
const searchKeyword = ref('');
const statusFilter = ref<string | undefined>();
const frameworkFilter = ref<string | undefined>();
const queueFilter = ref<string | undefined>();

// 选择相关
const selectedJobKeys = ref<string[]>([]);
const selectedJob = ref<TrainingApi.TrainingJob | null>(null);
const selectedJobId = ref<string | null>(null);

// 弹窗状态
const createModalVisible = ref(false);
const detailModalVisible = ref(false);
const logsModalVisible = ref(false);

// 表单引用
const createJobFormRef = ref();

// 任务列表数据
const jobsList = ref<TrainingApi.TrainingJob[]>([]);
const availableQueues = ref<TrainingApi.TrainingQueue[]>([]);

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
    title: '任务名称',
    dataIndex: 'name',
    key: 'name',
    width: 200,
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
    title: '训练框架',
    dataIndex: 'framework',
    key: 'framework',
    width: 120,
    slots: { customRender: 'framework' },
  },
  {
    title: '资源配置',
    key: 'resources',
    width: 180,
    slots: { customRender: 'resources' },
  },
  {
    title: '进度',
    dataIndex: 'progress',
    key: 'progress',
    width: 120,
    slots: { customRender: 'progress' },
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
    width: 160,
    fixed: 'right',
    slots: { customRender: 'action' },
  },
];

/**
 * 获取状态徽章类型
 */
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

/**
 * 获取状态文本
 */
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

/**
 * 获取框架颜色
 */
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

/**
 * 获取框架文本
 */
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

/**
 * 格式化日期时间
 */
const formatDateTime = (dateTime: string) => {
  return format(new Date(dateTime), 'MM-dd HH:mm', { locale: zhCN });
};

/**
 * 格式化持续时间
 */
const formatDuration = (seconds: number) => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  return `${hours}h${minutes}m`;
};

/**
 * 加载训练任务列表
 */
const loadJobs = async () => {
  try {
    loading.value = true;
    
    const params = {
      page: pagination.current,
      pageSize: pagination.pageSize,
      search: searchKeyword.value,
      status: statusFilter.value,
      framework: frameworkFilter.value,
      queueId: queueFilter.value,
    };

    const response = await trainingService.getJobs(params);
    
    jobsList.value = response.items;
    pagination.total = response.total;
  } catch (error) {
    message.error('加载训练任务失败');
  } finally {
    loading.value = false;
  }
};

/**
 * 加载可用队列
 */
const loadAvailableQueues = async () => {
  try {
    const queues = await trainingService.getAvailableQueues();
    availableQueues.value = queues;
  } catch (error) {
    message.error('加载队列列表失败');
  }
};

/**
 * 刷新任务列表
 */
const refreshJobs = () => {
  loadJobs();
};

/**
 * 处理搜索
 */
const handleSearch = () => {
  pagination.current = 1;
  loadJobs();
};

/**
 * 处理筛选
 */
const handleFilter = () => {
  pagination.current = 1;
  loadJobs();
};

/**
 * 重置筛选条件
 */
const resetFilters = () => {
  searchKeyword.value = '';
  statusFilter.value = undefined;
  frameworkFilter.value = undefined;
  queueFilter.value = undefined;
  handleFilter();
};

/**
 * 处理表格变化
 */
const handleTableChange = (paginationConfig: any) => {
  pagination.current = paginationConfig.current;
  pagination.pageSize = paginationConfig.pageSize;
  loadJobs();
};

/**
 * 处理选择变化
 */
const handleSelectionChange = (selectedRowKeys: string[]) => {
  selectedJobKeys.value = selectedRowKeys;
};

/**
 * 显示创建任务弹窗
 */
const showCreateModal = () => {
  createModalVisible.value = true;
};

/**
 * 处理创建任务
 */
const handleCreateJob = async () => {
  try {
    const formData = await createJobFormRef.value.validate();
    
    createLoading.value = true;
    await trainingService.createJob(formData);
    
    message.success('训练任务创建成功');
    createModalVisible.value = false;
    loadJobs();
  } catch (error) {
    message.error('创建训练任务失败');
  } finally {
    createLoading.value = false;
  }
};

/**
 * 取消创建
 */
const cancelCreate = () => {
  createModalVisible.value = false;
  createJobFormRef.value?.resetFields();
};

/**
 * 查看任务详情
 */
const viewJobDetail = async (jobId: string) => {
  try {
    const job = await trainingService.getJob(jobId);
    selectedJob.value = job;
    detailModalVisible.value = true;
  } catch (error) {
    message.error('获取任务详情失败');
  }
};

/**
 * 查看任务日志
 */
const viewJobLogs = (jobId: string) => {
  selectedJobId.value = jobId;
  logsModalVisible.value = true;
};

/**
 * 查看任务监控
 */
const viewJobMetrics = (jobId: string) => {
  router.push(`/training/jobs/${jobId}/metrics`);
};

/**
 * 处理任务操作
 */
const handleJobAction = async (job: TrainingApi.TrainingJob, action: string) => {
  try {
    switch (action) {
      case 'start':
      case 'stop':
      case 'restart':
      case 'suspend':
      case 'resume':
        await trainingService.controlJob({
          id: job.id,
          action: action as any,
        });
        message.success(`任务${action === 'start' ? '启动' : 
                       action === 'stop' ? '停止' : 
                       action === 'restart' ? '重启' :
                       action === 'suspend' ? '暂停' : '恢复'}成功`);
        loadJobs();
        break;
      case 'clone':
        const clonedJob = await trainingService.cloneJob(job.id, {
          name: `${job.name}-副本`,
        });
        message.success('任务克隆成功');
        loadJobs();
        break;
      case 'edit':
        // 打开编辑弹窗
        break;
      case 'delete':
        Modal.confirm({
          title: '确认删除',
          content: `确定要删除训练任务 "${job.name}" 吗？`,
          onOk: async () => {
            await trainingService.deleteJob(job.id);
            message.success('删除成功');
            loadJobs();
          },
        });
        break;
    }
  } catch (error) {
    message.error('操作失败');
  }
};

/**
 * 处理批量操作
 */
const handleBatchAction = async ({ key }: { key: string }) => {
  if (selectedJobKeys.value.length === 0) {
    message.warning('请选择要操作的任务');
    return;
  }

  try {
    switch (key) {
      case 'start':
      case 'stop':
        await trainingService.batchControlJobs(selectedJobKeys.value, key);
        message.success(`批量${key === 'start' ? '启动' : '停止'}成功`);
        break;
      case 'delete':
        Modal.confirm({
          title: '确认删除',
          content: `确定要删除选中的 ${selectedJobKeys.value.length} 个训练任务吗？`,
          onOk: async () => {
            await trainingService.batchDeleteJobs(selectedJobKeys.value);
            message.success('批量删除成功');
            selectedJobKeys.value = [];
          },
        });
        return;
    }
    
    loadJobs();
    selectedJobKeys.value = [];
  } catch (error) {
    message.error('批量操作失败');
  }
};

// 初始化
onMounted(() => {
  loadJobs();
  loadAvailableQueues();
});
</script>

<style lang="scss" scoped>
.training-jobs-container {
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

.filter-card {
  margin: 0 24px 16px 24px;
  border-radius: 8px;
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  .filter-content {
    padding: 8px 0;
  }
}

.jobs-table-card {
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

.job-name-cell {
  .job-name {
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

  .job-description {
    font-size: 12px;
    color: #666;
  }
}

.resource-info {
  .resource-item {
    display: flex;
    align-items: center;
    margin-bottom: 2px;
    font-size: 12px;
    color: #666;

    .resource-icon {
      width: 14px;
      height: 14px;
      margin-right: 4px;
    }
  }
}

.progress-cell {
  .progress-text {
    font-size: 12px;
    color: #666;
    text-align: center;
    margin-top: 4px;
  }
}

.time-info {
  .time-item {
    display: flex;
    align-items: center;
    margin-bottom: 2px;
    font-size: 11px;
    color: #666;

    .time-icon {
      width: 12px;
      height: 12px;
      margin-right: 4px;
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