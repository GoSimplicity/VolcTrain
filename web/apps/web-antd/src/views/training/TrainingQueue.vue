<template>
  <div class="training-queue-container">
    <!-- 页面头部 -->
    <Card>
      <div class="page-header">
        <div class="header-left">
          <h2>任务队列</h2>
          <p>管理训练任务队列和任务提交</p>
        </div>
        <div class="header-right">
          <Space>
            <Button @click="refreshData">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary" @click="showCreateJobModal">
              <PlusOutlined />
              提交任务
            </Button>
          </Space>
        </div>
      </div>
    </Card>

    <!-- 队列选择和统计 -->
    <Row :gutter="16" style="margin: 16px 0">
      <Col :span="6">
        <Card>
          <div class="queue-selector">
            <div class="selector-label">选择队列</div>
            <Select
              v-model:value="selectedQueueId"
              style="width: 100%"
              placeholder="选择训练队列"
              @change="handleQueueChange"
            >
              <Select.Option
                v-for="queue in availableQueues"
                :key="queue.id"
                :value="queue.id"
              >
                {{ queue.name }}
              </Select.Option>
            </Select>
          </div>
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="队列中任务"
            :value="queueStats.queuedCount"
            :value-style="{ color: '#faad14' }"
            prefix="⏳"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="运行中任务"
            :value="queueStats.runningCount"
            :value-style="{ color: '#1890ff' }"
            prefix="🚀"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="资源利用率"
            :value="queueStats.resourceUtilization"
            suffix="%"
            :value-style="{ color: '#52c41a' }"
            prefix="📊"
          />
        </Card>
      </Col>
    </Row>

    <!-- 任务筛选 -->
    <Card style="margin-bottom: 16px">
      <Row :gutter="16">
        <Col :span="6">
          <Input
            v-model:value="searchParams.keyword"
            placeholder="搜索任务名称"
            @change="handleSearch"
          >
            <template #prefix>
              <SearchOutlined />
            </template>
          </Input>
        </Col>
        <Col :span="4">
          <Select
            v-model:value="searchParams.status"
            placeholder="任务状态"
            style="width: 100%"
            @change="handleSearch"
          >
            <Select.Option value="">全部状态</Select.Option>
            <Select.Option value="pending">等待中</Select.Option>
            <Select.Option value="queued">队列中</Select.Option>
            <Select.Option value="running">运行中</Select.Option>
            <Select.Option value="completed">已完成</Select.Option>
            <Select.Option value="failed">失败</Select.Option>
            <Select.Option value="cancelled">已取消</Select.Option>
          </Select>
        </Col>
        <Col :span="4">
          <Select
            v-model:value="searchParams.framework"
            placeholder="训练框架"
            style="width: 100%"
            @change="handleSearch"
          >
            <Select.Option value="">全部框架</Select.Option>
            <Select.Option value="pytorch">PyTorch</Select.Option>
            <Select.Option value="tensorflow">TensorFlow</Select.Option>
            <Select.Option value="keras">Keras</Select.Option>
            <Select.Option value="paddlepaddle">PaddlePaddle</Select.Option>
            <Select.Option value="mindspore">MindSpore</Select.Option>
          </Select>
        </Col>
        <Col :span="4">
          <Select
            v-model:value="searchParams.priority"
            placeholder="优先级"
            style="width: 100%"
            @change="handleSearch"
          >
            <Select.Option value="">全部优先级</Select.Option>
            <Select.Option value="urgent">紧急</Select.Option>
            <Select.Option value="high">高</Select.Option>
            <Select.Option value="medium">中</Select.Option>
            <Select.Option value="low">低</Select.Option>
          </Select>
        </Col>
        <Col :span="6">
          <Space>
            <Button @click="resetSearch">重置</Button>
            <Button @click="batchOperation" :disabled="!hasSelectedJobs">
              批量操作
            </Button>
          </Space>
        </Col>
      </Row>
    </Card>

    <!-- 任务列表 -->
    <Card>
      <Table
        :columns="columns"
        :data-source="jobList"
        :loading="loading"
        :pagination="pagination"
        :row-selection="{ selectedRowKeys: selectedJobIds, onChange: onSelectChange }"
        @change="handleTableChange"
        row-key="id"
      >
        <!-- 任务名称 -->
        <template #name="{ record }">
          <div class="job-name">
            <div class="name-main">
              <Button type="link" @click="viewJobDetail(record)">
                {{ record.name }}
              </Button>
              <Tag :color="getFrameworkColor(record.framework)" size="small" style="margin-left: 8px">
                {{ record.framework }}
              </Tag>
            </div>
            <div class="name-desc">{{ record.description || '暂无描述' }}</div>
          </div>
        </template>

        <!-- 状态 -->
        <template #status="{ record }">
          <div class="status-info">
            <Tag :color="getJobStatusColor(record.status)">
              {{ getJobStatusLabel(record.status) }}
            </Tag>
            <div v-if="record.status === 'running'" class="status-detail">
              运行时长: {{ formatDuration(record.duration || 0) }}
            </div>
          </div>
        </template>

        <!-- 优先级 -->
        <template #priority="{ record }">
          <Tag :color="getPriorityColor(record.priority)">
            {{ getPriorityLabel(record.priority) }}
          </Tag>
        </template>

        <!-- 资源需求 -->
        <template #resources="{ record }">
          <div class="resource-info">
            <div v-if="record.resourceRequirements.gpu">
              GPU: {{ record.resourceRequirements.gpu }}
            </div>
            <div>
              CPU: {{ record.resourceRequirements.cpu }}核
            </div>
            <div>
              内存: {{ record.resourceRequirements.memory }}GB
            </div>
          </div>
        </template>

        <!-- 进度 -->
        <template #progress="{ record }">
          <div class="progress-info">
            <Progress
              :percent="record.progress"
              size="small"
              :status="record.status === 'failed' ? 'exception' : 'active'"
            />
            <div class="progress-text">{{ record.progress }}%</div>
          </div>
        </template>

        <!-- 创建者 -->
        <template #creator="{ record }">
          <div class="creator-info">
            <Avatar size="small">{{ record.creatorName?.[0] }}</Avatar>
            <span style="margin-left: 8px">{{ record.creatorName }}</span>
          </div>
        </template>

        <!-- 提交时间 -->
        <template #submitTime="{ record }">
          <div>
            <div>{{ formatDateTime(record.submitTime, 'MM-DD HH:mm') }}</div>
            <div style="font-size: 12px; color: #999">
              {{ formatRelativeTime(record.submitTime) }}
            </div>
          </div>
        </template>

        <!-- 操作 -->
        <template #action="{ record }">
          <Space size="small">
            <Button type="link" size="small" @click="viewJobDetail(record)">
              <EyeOutlined />
            </Button>
            <Button 
              type="link" 
              size="small" 
              @click="viewJobLogs(record)"
              :disabled="!canViewLogs(record.status)"
            >
              <FileTextOutlined />
            </Button>
            <Dropdown>
              <Button type="link" size="small">
                <MoreOutlined />
              </Button>
              <template #overlay>
                <Menu>
                  <Menu.Item 
                    key="start" 
                    @click="controlJob(record, 'start')"
                    :disabled="!canStart(record.status)"
                  >
                    <PlayCircleOutlined />
                    启动
                  </Menu.Item>
                  <Menu.Item 
                    key="pause" 
                    @click="controlJob(record, 'pause')"
                    :disabled="!canPause(record.status)"
                  >
                    <PauseCircleOutlined />
                    暂停
                  </Menu.Item>
                  <Menu.Item 
                    key="stop" 
                    @click="controlJob(record, 'stop')"
                    :disabled="!canStop(record.status)"
                    danger
                  >
                    <StopOutlined />
                    停止
                  </Menu.Item>
                  <Menu.Divider />
                  <Menu.Item key="clone" @click="cloneJob(record)">
                    <CopyOutlined />
                    克隆
                  </Menu.Item>
                  <Menu.Item key="edit" @click="editJob(record)">
                    <EditOutlined />
                    编辑
                  </Menu.Item>
                  <Menu.Divider />
                  <Menu.Item 
                    key="delete" 
                    @click="deleteJob(record)"
                    danger
                    :disabled="!canDelete(record.status)"
                  >
                    <DeleteOutlined />
                    删除
                  </Menu.Item>
                </Menu>
              </template>
            </Dropdown>
          </Space>
        </template>
      </Table>
    </Card>

    <!-- 创建任务模态框 -->
    <CreateJobModal
      v-model:visible="createJobModalVisible"
      :available-queues="availableQueues"
      :default-queue-id="selectedQueueId"
      @success="handleCreateSuccess"
    />

    <!-- 任务详情抽屉 -->
    <JobDetailDrawer
      v-model:visible="detailDrawerVisible"
      :job="selectedJob"
      @view-logs="handleViewLogsFromDetail"
      @refresh="loadJobs"
    />

    <!-- 日志查看抽屉 -->
    <JobLogsDrawer
      v-model:visible="logsDrawerVisible"
      :job="selectedJob"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue';
import {
  Card,
  Row,
  Col,
  Button,
  Space,
  Statistic,
  Select,
  Input,
  Table,
  Tag,
  Progress,
  Avatar,
  Dropdown,
  Menu,
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  PlusOutlined,
  SearchOutlined,
  EyeOutlined,
  FileTextOutlined,
  MoreOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  StopOutlined,
  CopyOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons-vue';
import { useRouter } from 'vue-router';
import type { 
  TrainingJob, 
  TrainingJobQuery, 
  TrainingQueue,
  TrainingJobControlRequest 
} from '#/api/types';
import { 
  getTrainingJobList, 
  getAvailableQueues,
  controlTrainingJob,
  deleteTrainingJob as deleteTrainingJobApi,
  cloneTrainingJob
} from '#/api';
import { formatDateTime, formatDuration, formatRelativeTime } from '#/utils/date';
import CreateJobModal from './components/CreateJobModal.vue';
import JobDetailDrawer from './components/JobDetailDrawer.vue';
import JobLogsDrawer from './components/JobLogsDrawer.vue';

const router = useRouter();

defineOptions({ name: 'TrainingQueue' });

// 响应式数据
const loading = ref(false);
const jobList = ref<TrainingJob[]>([]);
const availableQueues = ref<TrainingQueue[]>([]);
const selectedQueueId = ref<string>('');
const selectedJob = ref<TrainingJob | null>(null);
const selectedJobIds = ref<string[]>([]);

// 模态框和抽屉状态
const createJobModalVisible = ref(false);
const detailDrawerVisible = ref(false);
const logsDrawerVisible = ref(false);

// 搜索参数
const searchParams = reactive<TrainingJobQuery>({
  page: 1,
  pageSize: 10,
  keyword: '',
  status: undefined,
  framework: undefined,
  priority: undefined,
  queueId: '',
});

// 分页配置
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total: number) => `共 ${total} 条记录`,
});

// 队列统计
const queueStats = ref({
  queuedCount: 0,
  runningCount: 0,
  resourceUtilization: 0,
});

// 计算属性
const hasSelectedJobs = computed(() => selectedJobIds.value.length > 0);

// 表格列定义
const columns = [
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
    width: 120,
  },
  {
    title: '优先级',
    key: 'priority',
    slots: { customRender: 'priority' },
    width: 80,
  },
  {
    title: '资源需求',
    key: 'resources',
    slots: { customRender: 'resources' },
    width: 120,
  },
  {
    title: '进度',
    key: 'progress',
    slots: { customRender: 'progress' },
    width: 100,
  },
  {
    title: '创建者',
    key: 'creator',
    slots: { customRender: 'creator' },
    width: 120,
  },
  {
    title: '提交时间',
    key: 'submitTime',
    slots: { customRender: 'submitTime' },
    width: 100,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 120,
    fixed: 'right' as const,
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

const getFrameworkColor = (framework: string) => {
  const colors = {
    pytorch: 'orange',
    tensorflow: 'blue',
    keras: 'red',
    paddlepaddle: 'green',
    mindspore: 'purple',
    custom: 'default',
  };
  return colors[framework as keyof typeof colors] || 'default';
};

const canViewLogs = (status: string) => {
  return ['running', 'completed', 'failed', 'stopped'].includes(status);
};

const canStart = (status: string) => {
  return ['pending', 'paused', 'stopped'].includes(status);
};

const canPause = (status: string) => {
  return status === 'running';
};

const canStop = (status: string) => {
  return ['running', 'queued', 'pending'].includes(status);
};

const canDelete = (status: string) => {
  return ['completed', 'failed', 'cancelled', 'stopped'].includes(status);
};

// 数据加载
const loadJobs = async () => {
  try {
    loading.value = true;
    const params = {
      ...searchParams,
      queueId: selectedQueueId.value,
    };
    
    const response = await getTrainingJobList(params);
    jobList.value = response.data;
    pagination.total = response.total;
    pagination.current = response.page;
    pagination.pageSize = response.pageSize;
    
    // 更新队列统计
    updateQueueStats(response.data);
  } catch (error) {
    message.error('加载任务列表失败');
  } finally {
    loading.value = false;
  }
};

const loadAvailableQueues = async () => {
  try {
    const response = await getAvailableQueues();
    availableQueues.value = response;
    if (response.length > 0 && !selectedQueueId.value) {
      selectedQueueId.value = response[0].id;
    }
  } catch (error) {
    message.error('加载可用队列失败');
  }
};

const updateQueueStats = (jobs: TrainingJob[]) => {
  queueStats.value.queuedCount = jobs.filter(job => 
    ['pending', 'queued'].includes(job.status)
  ).length;
  
  queueStats.value.runningCount = jobs.filter(job => 
    job.status === 'running'
  ).length;
  
  // 简单的资源利用率计算
  queueStats.value.resourceUtilization = Math.min(
    (queueStats.value.runningCount / Math.max(jobs.length, 1)) * 100,
    100
  );
};

// 事件处理
const handleQueueChange = (queueId: string) => {
  selectedQueueId.value = queueId;
  searchParams.queueId = queueId;
  searchParams.page = 1;
  pagination.current = 1;
  loadJobs();
};

const handleSearch = () => {
  searchParams.page = 1;
  pagination.current = 1;
  loadJobs();
};

const resetSearch = () => {
  searchParams.keyword = '';
  searchParams.status = undefined;
  searchParams.framework = undefined;
  searchParams.priority = undefined;
  handleSearch();
};

const handleTableChange = (pag: any) => {
  searchParams.page = pag.current;
  searchParams.pageSize = pag.pageSize;
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  loadJobs();
};

const onSelectChange = (keys: string[]) => {
  selectedJobIds.value = keys;
};

const refreshData = () => {
  loadJobs();
  loadAvailableQueues();
};

// 任务操作
const showCreateJobModal = () => {
  createJobModalVisible.value = true;
};

const viewJobDetail = (job: TrainingJob) => {
  selectedJob.value = job;
  detailDrawerVisible.value = true;
};

const viewJobLogs = (job: TrainingJob) => {
  selectedJob.value = job;
  logsDrawerVisible.value = true;
};

const controlJob = async (job: TrainingJob, action: string) => {
  try {
    const request: TrainingJobControlRequest = {
      id: job.id,
      action: action as any,
    };
    await controlTrainingJob(request);
    message.success(`${action === 'start' ? '启动' : action === 'pause' ? '暂停' : '停止'}任务成功`);
    loadJobs();
  } catch (error) {
    message.error('操作失败');
  }
};

const cloneJob = async (job: TrainingJob) => {
  try {
    await cloneTrainingJob(job.id, { name: `${job.name}_copy` });
    message.success('克隆任务成功');
    loadJobs();
  } catch (error) {
    message.error('克隆失败');
  }
};

const editJob = (job: TrainingJob) => {
  router.push(`/training/jobs/${job.id}/edit`);
};

const deleteJob = async (job: TrainingJob) => {
  try {
    await deleteTrainingJobApi(job.id);
    message.success('删除任务成功');
    loadJobs();
  } catch (error) {
    message.error('删除失败');
  }
};

const batchOperation = () => {
  console.log('批量操作选中的任务:', selectedJobIds.value);
  // 实现批量操作逻辑
};

const handleCreateSuccess = () => {
  createJobModalVisible.value = false;
  loadJobs();
};

const handleViewLogsFromDetail = (job: TrainingJob) => {
  selectedJob.value = job;
  detailDrawerVisible.value = false;
  logsDrawerVisible.value = true;
};

// 初始化
onMounted(() => {
  loadAvailableQueues().then(() => {
    if (selectedQueueId.value) {
      loadJobs();
    }
  });
});
</script>

<style scoped lang="scss">
.training-queue-container {
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

.queue-selector {
  .selector-label {
    font-size: 14px;
    margin-bottom: 8px;
    color: #333;
    font-weight: 500;
  }
}

.job-name {
  .name-main {
    display: flex;
    align-items: center;
  }
  
  .name-desc {
    font-size: 12px;
    color: #999;
    margin-top: 4px;
  }
}

.status-info {
  .status-detail {
    font-size: 12px;
    color: #999;
    margin-top: 4px;
  }
}

.resource-info {
  font-size: 12px;
  
  div {
    margin-bottom: 2px;
  }
}

.progress-info {
  .progress-text {
    font-size: 12px;
    text-align: center;
    margin-top: 4px;
  }
}

.creator-info {
  display: flex;
  align-items: center;
}
</style>
