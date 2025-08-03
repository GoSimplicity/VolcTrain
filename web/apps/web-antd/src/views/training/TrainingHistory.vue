<template>
  <div class="training-history-container">
    <!-- 页面头部 -->
    <Card>
      <div class="page-header">
        <div class="header-left">
          <h2>训练历史</h2>
          <p>查看所有历史训练任务的详细记录</p>
        </div>
        <div class="header-right">
          <Space>
            <Button @click="exportData" :loading="exporting">
              <ExportOutlined />
              导出数据
            </Button>
            <Button @click="refreshData">
              <ReloadOutlined />
              刷新
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
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="平均耗时"
            :value="statistics.avgDuration"
            precision="1"
            suffix="h"
            :value-style="{ color: '#1890ff' }"
            prefix="⏱️"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="GPU总时长"
            :value="statistics.totalGpuHours"
            precision="1"
            suffix="h"
            :value-style="{ color: '#722ed1' }"
            prefix="🔥"
          />
        </Card>
      </Col>
    </Row>

    <!-- 筛选器 -->
    <Card style="margin-bottom: 16px">
      <Row :gutter="16">
        <Col :span="4">
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
        <Col :span="3">
          <Select
            v-model:value="searchParams.status"
            placeholder="任务状态"
            style="width: 100%"
            @change="handleSearch"
            allow-clear
          >
            <Select.Option value="">全部状态</Select.Option>
            <Select.Option value="completed">已完成</Select.Option>
            <Select.Option value="failed">失败</Select.Option>
            <Select.Option value="cancelled">已取消</Select.Option>
            <Select.Option value="stopped">已停止</Select.Option>
          </Select>
        </Col>
        <Col :span="3">
          <Select
            v-model:value="searchParams.framework"
            placeholder="训练框架"
            style="width: 100%"
            @change="handleSearch"
            allow-clear
          >
            <Select.Option value="">全部框架</Select.Option>
            <Select.Option value="pytorch">PyTorch</Select.Option>
            <Select.Option value="tensorflow">TensorFlow</Select.Option>
            <Select.Option value="keras">Keras</Select.Option>
            <Select.Option value="paddlepaddle">PaddlePaddle</Select.Option>
            <Select.Option value="mindspore">MindSpore</Select.Option>
          </Select>
        </Col>
        <Col :span="10">
          <Space>
            <Button @click="resetSearch">重置</Button>
            <Button type="primary" @click="advancedSearch">
              高级搜索
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
        @change="handleTableChange"
        row-key="id"
        :row-selection="{ selectedRowKeys: selectedJobIds, onChange: onSelectChange }"
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
            <div class="status-detail">
              <div v-if="record.endTime">
                结束时间: {{ formatDateTime(record.endTime, 'MM-DD HH:mm') }}
              </div>
            </div>
          </div>
        </template>

        <!-- 优先级 -->
        <template #priority="{ record }">
          <Tag :color="getPriorityColor(record.priority)">
            {{ getPriorityLabel(record.priority) }}
          </Tag>
        </template>

        <!-- 资源使用 -->
        <template #resources="{ record }">
          <div class="resource-info">
            <div class="resource-summary">
              <Tooltip title="CPU使用">
                <Tag color="blue">{{ record.resourceRequirements.cpu }}C</Tag>
              </Tooltip>
              <Tooltip title="内存使用">
                <Tag color="green">{{ record.resourceRequirements.memory }}G</Tag>
              </Tooltip>
              <Tooltip title="GPU使用" v-if="record.resourceRequirements.gpu">
                <Tag color="purple">{{ record.resourceRequirements.gpu }}GPU</Tag>
              </Tooltip>
            </div>
          </div>
        </template>

        <!-- 进度 -->
        <template #progress="{ record }">
          <div class="progress-info">
            <Progress
              :percent="record.progress"
              size="small"
              :status="getProgressStatus(record.status)"
            />
            <div class="progress-text">{{ record.progress }}%</div>
          </div>
        </template>

        <!-- 运行时长 -->
        <template #duration="{ record }">
          <div class="duration-info">
            <div class="duration-main">{{ formatDuration(record.duration || 0) }}</div>
            <div class="duration-detail" v-if="record.startTime && record.endTime">
              {{ formatDateTime(record.startTime, 'MM-DD HH:mm') }} - 
              {{ formatDateTime(record.endTime, 'MM-DD HH:mm') }}
            </div>
          </div>
        </template>

        <!-- 创建者 -->
        <template #creator="{ record }">
          <div class="creator-info">
            <Avatar size="small">{{ record.creatorName?.[0] }}</Avatar>
            <span style="margin-left: 8px">{{ record.creatorName }}</span>
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
            >
              <FileTextOutlined />
            </Button>
            <Dropdown>
              <Button type="link" size="small">
                <MoreOutlined />
              </Button>
              <template #overlay>
                <Menu>
                  <Menu.Item key="clone" @click="cloneJob(record)">
                    <CopyOutlined />
                    克隆
                  </Menu.Item>
                  <Menu.Item key="export" @click="exportJob(record)">
                    <ExportOutlined />
                    导出
                  </Menu.Item>
                  <Menu.Divider />
                  <Menu.Item 
                    key="delete" 
                    @click="deleteJob(record)"
                    danger
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

      <!-- 批量操作 -->
      <div v-if="selectedJobIds.length > 0" class="batch-actions">
        <Space>
          <span>已选择 {{ selectedJobIds.length }} 项</span>
          <Button @click="batchExport" :loading="batchLoading">
            <ExportOutlined />
            批量导出
          </Button>
          <Button @click="batchDelete" :loading="batchLoading" danger>
            <DeleteOutlined />
            批量删除
          </Button>
          <Button @click="clearSelection">
            清空选择
          </Button>
        </Space>
      </div>
    </Card>

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

    <!-- 高级搜索模态框 -->
    <Modal
      v-model:open="advancedSearchVisible"
      title="高级搜索"
      width="600px"
      @ok="handleAdvancedSearch"
      @cancel="handleAdvancedSearchCancel"
    >
      <Form
        :model="advancedSearchForm"
        layout="vertical"
      >
        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="GPU使用范围">
              <Space>
                <InputNumber v-model:value="advancedSearchForm.minGpu" placeholder="最小" :min="0" />
                <span>-</span>
                <InputNumber v-model:value="advancedSearchForm.maxGpu" placeholder="最大" :min="0" />
              </Space>
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="运行时长范围(小时)">
              <Space>
                <InputNumber v-model:value="advancedSearchForm.minDuration" placeholder="最小" :min="0" />
                <span>-</span>
                <InputNumber v-model:value="advancedSearchForm.maxDuration" placeholder="最大" :min="0" />
              </Space>
            </Form.Item>
          </Col>
        </Row>
      </Form>
    </Modal>
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
  Input,
  Select,
  Table,
  Tag,
  Progress,
  Avatar,
  Dropdown,
  Menu,
  Tooltip,
  Modal,
  Form,
  InputNumber,
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  SearchOutlined,
  EyeOutlined,
  FileTextOutlined,
  MoreOutlined,
  CopyOutlined,
  ExportOutlined,
  DeleteOutlined,
} from '@ant-design/icons-vue';
import { useRouter } from 'vue-router';
import type { TrainingJob, TrainingJobQuery } from '#/api/types';
import { 
  getTrainingJobList, 
  deleteTrainingJob,
  cloneTrainingJob,
  batchDeleteJobs
} from '#/api';
import { formatDateTime, formatDuration, formatRelativeTime } from '#/utils/date';
import JobDetailDrawer from './components/JobDetailDrawer.vue';
import JobLogsDrawer from './components/JobLogsDrawer.vue';

const router = useRouter();

defineOptions({ name: 'TrainingHistory' });

// 响应式数据
const loading = ref(false);
const exporting = ref(false);
const batchLoading = ref(false);
const jobList = ref<TrainingJob[]>([]);
const selectedJob = ref<TrainingJob | null>(null);
const selectedJobIds = ref<string[]>([]);
const detailDrawerVisible = ref(false);
const logsDrawerVisible = ref(false);
const advancedSearchVisible = ref(false);

// 搜索参数
const searchParams = reactive<TrainingJobQuery>({
  page: 1,
  pageSize: 10,
  keyword: '',
  status: undefined,
  framework: undefined,
});

// 高级搜索表单
const advancedSearchForm = reactive({
  minGpu: undefined,
  maxGpu: undefined,
  minDuration: undefined,
  maxDuration: undefined,
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

// 统计数据
const statistics = ref({
  total: 0,
  successRate: 0,
  avgDuration: 0,
  totalGpuHours: 0,
});

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
    title: '资源使用',
    key: 'resources',
    slots: { customRender: 'resources' },
    width: 140,
  },
  {
    title: '进度',
    key: 'progress',
    slots: { customRender: 'progress' },
    width: 100,
  },
  {
    title: '运行时长',
    key: 'duration',
    slots: { customRender: 'duration' },
    width: 120,
  },
  {
    title: '创建者',
    key: 'creator',
    slots: { customRender: 'creator' },
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
    completed: 'success',
    failed: 'error',
    cancelled: 'warning',
    stopped: 'default',
  };
  return colors[status as keyof typeof colors] || 'default';
};

const getJobStatusLabel = (status: string) => {
  const labels = {
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消',
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

const getProgressStatus = (status: string) => {
  if (status === 'failed') return 'exception';
  if (status === 'completed') return 'success';
  return 'normal';
};

// 数据加载
const loadJobs = async () => {
  try {
    loading.value = true;
    const params = {
      ...searchParams,
      // 只查询历史任务
      status: searchParams.status || 'completed,failed,cancelled,stopped',
    };
    
    const response = await getTrainingJobList(params);
    jobList.value = response.data;
    pagination.total = response.total;
    pagination.current = response.page;
    pagination.pageSize = response.pageSize;
    
    // 更新统计数据
    updateStatistics(response.data);
  } catch (error) {
    message.error('加载任务列表失败');
  } finally {
    loading.value = false;
  }
};

const updateStatistics = (jobs: TrainingJob[]) => {
  statistics.value.total = jobs.length;
  
  const completedJobs = jobs.filter(job => job.status === 'completed');
  statistics.value.successRate = jobs.length > 0 
    ? (completedJobs.length / jobs.length) * 100 
    : 0;
  
  const totalDuration = jobs.reduce((sum, job) => sum + (job.duration || 0), 0);
  statistics.value.avgDuration = jobs.length > 0 ? totalDuration / jobs.length / 3600 : 0;
  
  const totalGpuHours = jobs.reduce((sum, job) => {
    const gpuCount = job.resourceRequirements.gpu || 0;
    const hours = (job.duration || 0) / 3600;
    return sum + (gpuCount * hours);
  }, 0);
  statistics.value.totalGpuHours = totalGpuHours;
};

const refreshData = () => {
  loadJobs();
};

// 事件处理
const handleSearch = () => {
  searchParams.page = 1;
  pagination.current = 1;
  loadJobs();
};

const resetSearch = () => {
  Object.assign(searchParams, {
    page: 1,
    pageSize: 10,
    keyword: '',
    status: undefined,
    framework: undefined,
  });
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

const clearSelection = () => {
  selectedJobIds.value = [];
};

// 任务操作
const viewJobDetail = (job: TrainingJob) => {
  selectedJob.value = job;
  detailDrawerVisible.value = true;
};

const viewJobLogs = (job: TrainingJob) => {
  selectedJob.value = job;
  logsDrawerVisible.value = true;
};

const handleViewLogsFromDetail = (job: TrainingJob) => {
  selectedJob.value = job;
  detailDrawerVisible.value = false;
  logsDrawerVisible.value = true;
};

const cloneJob = async (job: TrainingJob) => {
  try {
    await cloneTrainingJob(job.id, { name: `${job.name}_copy` });
    message.success('克隆任务成功');
    router.push('/training/queue');
  } catch (error) {
    message.error('克隆失败');
  }
};

const exportJob = (job: TrainingJob) => {
  // 实现单个任务导出
  const data = JSON.stringify(job, null, 2);
  const blob = new Blob([data], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${job.name}-export.json`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
  message.success('导出成功');
};

const deleteJob = async (job: TrainingJob) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除训练任务 "${job.name}" 吗？此操作不可恢复。`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        await deleteTrainingJob(job.id);
        message.success('删除任务成功');
        loadJobs();
      } catch (error) {
        message.error('删除失败');
      }
    },
  });
};

// 批量操作
const batchExport = async () => {
  try {
    batchLoading.value = true;
    const selectedJobs = jobList.value.filter(job => 
      selectedJobIds.value.includes(job.id)
    );
    
    const data = JSON.stringify(selectedJobs, null, 2);
    const blob = new Blob([data], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `training-jobs-export-${new Date().getTime()}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    message.success(`成功导出 ${selectedJobs.length} 个任务`);
    clearSelection();
  } catch (error) {
    message.error('批量导出失败');
  } finally {
    batchLoading.value = false;
  }
};

const batchDelete = async () => {
  Modal.confirm({
    title: '确认批量删除',
    content: `确定要删除选中的 ${selectedJobIds.value.length} 个训练任务吗？此操作不可恢复。`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        batchLoading.value = true;
        await batchDeleteJobs(selectedJobIds.value);
        message.success(`成功删除 ${selectedJobIds.value.length} 个任务`);
        clearSelection();
        loadJobs();
      } catch (error) {
        message.error('批量删除失败');
      } finally {
        batchLoading.value = false;
      }
    },
  });
};

const exportData = async () => {
  try {
    exporting.value = true;
    // 导出当前筛选结果的所有数据
    const allData = JSON.stringify(jobList.value, null, 2);
    const blob = new Blob([allData], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `training-history-${new Date().getTime()}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    
    message.success('数据导出成功');
  } catch (error) {
    message.error('导出失败');
  } finally {
    exporting.value = false;
  }
};

// 高级搜索
const advancedSearch = () => {
  advancedSearchVisible.value = true;
};

const handleAdvancedSearch = () => {
  advancedSearchVisible.value = false;
  handleSearch();
  message.success('高级搜索已应用');
};

const handleAdvancedSearchCancel = () => {
  advancedSearchVisible.value = false;
};

// 初始化
onMounted(() => {
  loadJobs();
});
</script>

<style scoped lang="scss">
.training-history-container {
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
  .resource-summary {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-bottom: 4px;
  }
}

.progress-info {
  .progress-text {
    font-size: 12px;
    text-align: center;
    margin-top: 4px;
  }
}

.duration-info {
  .duration-main {
    font-weight: 500;
  }
  
  .duration-detail {
    font-size: 12px;
    color: #999;
    margin-top: 4px;
  }
}

.creator-info {
  display: flex;
  align-items: center;
}

.batch-actions {
  margin-top: 16px;
  padding: 12px 16px;
  background: #f5f5f5;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>