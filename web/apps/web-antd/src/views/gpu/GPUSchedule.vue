<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue';
import {
  Card,
  Row,
  Col,
  Button,
  Space,
  Statistic,
  Tabs,
  Table,
  Tag,
  Badge,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  DatePicker,
  TimePicker,
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  PlusOutlined,
  CalendarOutlined,
  BarChartOutlined,
  PlayCircleOutlined,
  StopOutlined,
  EditOutlined,
  DeleteOutlined,
} from '@ant-design/icons-vue';
import type { Dayjs } from 'dayjs';
import type { 
  ResourceScheduleRequest,
  ResourceScheduleResponse,
  GPUCluster
} from '#/api/types';
import { 
  scheduleResources,
  getScheduleHistory,
  getClusterList
} from '#/api';
import { formatDateTime, formatDuration } from '#/utils/date';
import ScheduleCalendar from './components/ScheduleCalendar.vue';

defineOptions({ name: 'GPUSchedule' });

// 响应式数据
const loading = ref(false);
const scheduleList = ref<any[]>([]);
const availableResources = ref<any[]>([]);
const clusters = ref<GPUCluster[]>([]);
const createModalVisible = ref(false);
const activeTab = ref('schedule');

// 统计数据
const scheduleStats = ref({
  totalSchedules: 0,
  runningCount: 0,
  pendingCount: 0,
  completedCount: 0,
  resourceUtilization: 0,
});

// 表单数据
interface ScheduleForm {
  name: string;
  clusterId: string;
  gpuType: string;
  gpuCount: number;
  startDate: Dayjs | null;
  startTime: Dayjs | null;
  duration: number;
  priority: number;
  strategy: string;
  description: string;
}

const scheduleForm = reactive<ScheduleForm>({
  name: '',
  clusterId: '',
  gpuType: '',
  gpuCount: 1,
  startDate: null,
  startTime: null,
  duration: 1,
  priority: 5,
  strategy: 'best_fit',
  description: '',
});

const formRef = ref();

// 模拟数据
const mockSchedules = [
  {
    id: 'schedule-001',
    name: 'BERT模型训练',
    userId: 'user-001',
    userName: '张三',
    clusterId: 'cluster-01',
    clusterName: 'Main Cluster',
    gpuType: 'Tesla A100',
    gpuCount: 4,
    status: 'running',
    priority: 8,
    strategy: 'best_fit',
    startTime: '2024-01-20 09:00:00',
    endTime: '2024-01-20 15:00:00',
    duration: 21600, // 6小时
    progress: 65,
    createTime: '2024-01-19 16:00:00',
  },
  {
    id: 'schedule-002',
    name: 'ResNet图像分类',
    userId: 'user-002',
    userName: '李四',
    clusterId: 'cluster-01',
    clusterName: 'Main Cluster',
    gpuType: 'Tesla V100',
    gpuCount: 2,
    status: 'pending',
    priority: 5,
    strategy: 'balanced',
    startTime: '2024-01-20 16:00:00',
    endTime: '2024-01-21 04:00:00',
    duration: 43200, // 12小时
    progress: 0,
    createTime: '2024-01-20 08:00:00',
  },
  {
    id: 'schedule-003',
    name: 'GPT训练任务',
    userId: 'user-003',
    userName: '王五',
    clusterId: 'cluster-02',
    clusterName: 'Training Cluster',
    gpuType: 'Tesla A100',
    gpuCount: 8,
    status: 'completed',
    priority: 9,
    strategy: 'gpu_optimized',
    startTime: '2024-01-19 10:00:00',
    endTime: '2024-01-20 02:00:00',
    duration: 57600, // 16小时
    progress: 100,
    createTime: '2024-01-19 08:00:00',
  },
];

const mockResources = [
  {
    clusterId: 'cluster-01',
    clusterName: 'Main Cluster',
    gpuType: 'Tesla A100',
    total: 16,
    allocated: 8,
    available: 8,
    utilization: 50,
  },
  {
    clusterId: 'cluster-01',
    clusterName: 'Main Cluster',
    gpuType: 'Tesla V100',
    total: 24,
    allocated: 6,
    available: 18,
    utilization: 25,
  },
  {
    clusterId: 'cluster-02',
    clusterName: 'Training Cluster',
    gpuType: 'Tesla A100',
    total: 32,
    allocated: 24,
    available: 8,
    utilization: 75,
  },
];

// 表格列定义
const scheduleColumns = [
  {
    title: '调度任务',
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
    title: '资源需求',
    key: 'resources',
    slots: { customRender: 'resources' },
    width: 120,
  },
  {
    title: '优先级',
    key: 'priority',
    slots: { customRender: 'priority' },
    width: 80,
  },
  {
    title: '进度',
    key: 'progress',
    slots: { customRender: 'progress' },
    width: 100,
  },
  {
    title: '调度时间',
    key: 'scheduleTime',
    slots: { customRender: 'scheduleTime' },
    width: 180,
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
    width: 150,
    fixed: 'right' as const,
  },
];

const resourceColumns = [
  {
    title: '集群',
    dataIndex: 'clusterName',
    key: 'clusterName',
    width: 150,
  },
  {
    title: 'GPU类型',
    dataIndex: 'gpuType',
    key: 'gpuType',
    width: 120,
  },
  {
    title: '总数',
    dataIndex: 'total',
    key: 'total',
    width: 80,
  },
  {
    title: '已分配',
    dataIndex: 'allocated',
    key: 'allocated',
    width: 80,
  },
  {
    title: '可用',
    dataIndex: 'available',
    key: 'available',
    width: 80,
  },
  {
    title: '使用率',
    key: 'utilization',
    slots: { customRender: 'utilization' },
    width: 120,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'resourceAction' },
    width: 100,
  },
];
// 工具方法
const getScheduleStatusColor = (status: string) => {
  const colors = {
    pending: 'processing',
    running: 'success',
    completed: 'default',
    failed: 'error',
    cancelled: 'warning',
  };
  return colors[status as keyof typeof colors] || 'default';
};

const getScheduleStatusLabel = (status: string) => {
  const labels = {
    pending: '等待中',
    running: '运行中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消',
  };
  return labels[status as keyof typeof labels] || status;
};

const getPriorityColor = (priority: number) => {
  if (priority >= 8) return 'red';
  if (priority >= 6) return 'orange';
  if (priority >= 4) return 'blue';
  return 'default';
};

const getUtilizationColor = (utilization: number) => {
  if (utilization >= 90) return '#f5222d';
  if (utilization >= 70) return '#fa8c16';
  if (utilization >= 40) return '#52c41a';
  return '#1890ff';
};

// 数据加载
const loadSchedules = async () => {
  try {
    loading.value = true;
    // const response = await getScheduleHistory();
    // scheduleList.value = response.data;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    scheduleList.value = mockSchedules;
    
    updateStats();
  } catch (error) {
    message.error('加载调度列表失败');
  } finally {
    loading.value = false;
  }
};

const loadResources = async () => {
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 300));
    availableResources.value = mockResources;
  } catch (error) {
    message.error('加载资源信息失败');
  }
};

const loadClusters = async () => {
  try {
    // const response = await getClusterList();
    // clusters.value = response.data;
    
    // 模拟集群数据
    clusters.value = [
      {
        id: 'cluster-01',
        name: 'Main Cluster',
        description: '主训练集群',
        type: 'kubernetes',
        apiEndpoint: 'https://main-cluster.example.com',
        nodeCount: 10,
        gpuCount: 40,
        status: 'running',
        healthScore: 95,
        totalResources: { cpu: 320, memory: 1280, gpu: 40 },
        usedResources: { cpu: 180, memory: 720, gpu: 24 },
        availableResources: { cpu: 140, memory: 560, gpu: 16 },
        monitoringEnabled: true,
        alertingEnabled: true,
        createTime: '2024-01-01 00:00:00',
        updateTime: '2024-01-20 15:30:00',
      },
    ];
  } catch (error) {
    message.error('加载集群信息失败');
  }
};

const updateStats = () => {
  const stats = {
    totalSchedules: scheduleList.value.length,
    runningCount: scheduleList.value.filter(s => s.status === 'running').length,
    pendingCount: scheduleList.value.filter(s => s.status === 'pending').length,
    completedCount: scheduleList.value.filter(s => s.status === 'completed').length,
    resourceUtilization: 0,
  };
  
  // 计算平均资源利用率
  if (availableResources.value.length > 0) {
    const totalUtil = availableResources.value.reduce((sum, res) => sum + res.utilization, 0);
    stats.resourceUtilization = Math.round(totalUtil / availableResources.value.length);
  }
  
  scheduleStats.value = stats;
};

const refreshData = () => {
  loadSchedules();
  loadResources();
};

// 事件处理
const showCreateModal = () => {
  createModalVisible.value = true;
  resetForm();
};

const resetForm = () => {
  Object.assign(scheduleForm, {
    name: '',
    clusterId: '',
    gpuType: '',
    gpuCount: 1,
    startDate: null,
    startTime: null,
    duration: 1,
    priority: 5,
    strategy: 'best_fit',
    description: '',
  });
};

const handleCreateSubmit = async () => {
  try {
    await formRef.value?.validate();
    
    const request: ResourceScheduleRequest = {
      strategy: scheduleForm.strategy as any,
      requirements: {
        cpu: scheduleForm.gpuCount * 4, // 假设每个GPU需要4核CPU
        memory: scheduleForm.gpuCount * 16, // 假设每个GPU需要16GB内存
        gpu: scheduleForm.gpuCount,
      },
      constraints: {
        clusterId: scheduleForm.clusterId,
        gpuModel: scheduleForm.gpuType,
      },
      priority: scheduleForm.priority,
      timeout: scheduleForm.duration * 3600, // 转换为秒
    };
    
    // const response = await scheduleResources(request);
    
    // 模拟API调用成功
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    message.success('资源调度请求提交成功');
    createModalVisible.value = false;
    loadSchedules();
  } catch (error) {
    message.error('提交失败');
  }
};

const handleCreateCancel = () => {
  createModalVisible.value = false;
};

// 调度操作
const startSchedule = async (schedule: any) => {
  try {
    // 实际应该调用启动API
    message.success('调度任务启动成功');
    loadSchedules();
  } catch (error) {
    message.error('启动失败');
  }
};

const stopSchedule = async (schedule: any) => {
  Modal.confirm({
    title: '确认停止',
    content: `确定要停止调度任务 "${schedule.name}" 吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // 实际应该调用停止API
        message.success('调度任务停止成功');
        loadSchedules();
      } catch (error) {
        message.error('停止失败');
      }
    },
  });
};

const editSchedule = (schedule: any) => {
  // 实现编辑逻辑
  message.info('编辑功能开发中');
};

const deleteSchedule = async (schedule: any) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除调度任务 "${schedule.name}" 吗？此操作不可恢复。`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // 实际应该调用删除API
        message.success('调度任务删除成功');
        loadSchedules();
      } catch (error) {
        message.error('删除失败');
      }
    },
  });
};

const reserveResource = (resource: any) => {
  message.info('资源预约功能开发中');
};

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入调度任务名称', trigger: 'blur' },
  ],
  clusterId: [
    { required: true, message: '请选择集群', trigger: 'change' },
  ],
  gpuType: [
    { required: true, message: '请选择GPU类型', trigger: 'change' },
  ],
  gpuCount: [
    { required: true, message: '请输入GPU数量', trigger: 'blur' },
  ],
  startDate: [
    { required: true, message: '请选择开始日期', trigger: 'change' },
  ],
  startTime: [
    { required: true, message: '请选择开始时间', trigger: 'change' },
  ],
  duration: [
    { required: true, message: '请输入持续时间', trigger: 'blur' },
  ],
};

// 初始化
onMounted(() => {
  loadSchedules();
  loadResources();
  loadClusters();
});
</script>

<template>
  <div class="gpu-schedule-container">
    <!-- 页面头部 -->
    <Card>
      <div class="page-header">
        <div class="header-left">
          <h2>GPU资源调度</h2>
          <p>管理GPU资源的分配和调度策略</p>
        </div>
        <div class="header-right">
          <Space>
            <Button @click="refreshData">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary" @click="showCreateModal">
              <PlusOutlined />
              创建调度
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
            title="总调度数"
            :value="scheduleStats.totalSchedules"
            :value-style="{ color: '#3f8600' }"
            prefix="📋"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="运行中"
            :value="scheduleStats.runningCount"
            :value-style="{ color: '#52c41a' }"
            prefix="🚀"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="等待中"
            :value="scheduleStats.pendingCount"
            :value-style="{ color: '#faad14' }"
            prefix="⏳"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="资源利用率"
            :value="scheduleStats.resourceUtilization"
            suffix="%"
            :value-style="{ color: '#1890ff' }"
            prefix="📊"
          />
        </Card>
      </Col>
    </Row>

    <!-- 标签页 -->
    <Tabs v-model:activeKey="activeTab">
      <Tabs.TabPane key="schedule" tab="调度任务">
        <Card>
          <Table
            :columns="scheduleColumns"
            :data-source="scheduleList"
            :loading="loading"
            row-key="id"
            :pagination="{ pageSize: 10 }"
          >
            <!-- 任务名称 -->
            <template #name="{ record }">
              <div class="schedule-name">
                <div class="name-main">{{ record.name }}</div>
                <div class="name-desc">{{ record.description || '无描述' }}</div>
              </div>
            </template>

            <!-- 状态 -->
            <template #status="{ record }">
              <Badge 
                :status="getScheduleStatusColor(record.status) as any" 
                :text="getScheduleStatusLabel(record.status)"
              />
            </template>

            <!-- 资源需求 -->
            <template #resources="{ record }">
              <div class="resource-info">
                <div>{{ record.gpuType }}</div>
                <div>{{ record.gpuCount }} GPU</div>
              </div>
            </template>

            <!-- 优先级 -->
            <template #priority="{ record }">
              <Tag :color="getPriorityColor(record.priority)">
                {{ record.priority }}
              </Tag>
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

            <!-- 调度时间 -->
            <template #scheduleTime="{ record }">
              <div class="time-info">
                <div>{{ formatDateTime(record.startTime, 'MM-DD HH:mm') }}</div>
                <div class="time-duration">
                  {{ formatDuration(record.duration) }}
                </div>
              </div>
            </template>

            <!-- 创建者 -->
            <template #creator="{ record }">
              <div class="creator-info">
                <Avatar size="small">{{ record.userName?.[0] }}</Avatar>
                <span style="margin-left: 8px">{{ record.userName }}</span>
              </div>
            </template>

            <!-- 操作 -->
            <template #action="{ record }">
              <Space size="small">
                <Button 
                  v-if="record.status === 'pending'" 
                  type="link" 
                  size="small" 
                  @click="startSchedule(record)"
                >
                  <PlayCircleOutlined />
                </Button>
                <Button 
                  v-if="record.status === 'running'" 
                  type="link" 
                  size="small" 
                  @click="stopSchedule(record)"
                  danger
                >
                  <StopOutlined />
                </Button>
                <Button 
                  v-if="record.status !== 'completed'" 
                  type="link" 
                  size="small" 
                  @click="editSchedule(record)"
                >
                  <EditOutlined />
                </Button>
                <Button 
                  type="link" 
                  size="small" 
                  @click="deleteSchedule(record)"
                  danger
                >
                  <DeleteOutlined />
                </Button>
              </Space>
            </template>
          </Table>
        </Card>
      </Tabs.TabPane>

      <Tabs.TabPane key="resources" tab="可用资源">
        <Card>
          <Table
            :columns="resourceColumns"
            :data-source="availableResources"
            row-key="clusterId"
            :pagination="false"
          >
            <!-- 使用率 -->
            <template #utilization="{ record }">
              <div class="utilization-info">
                <Progress
                  :percent="record.utilization"
                  size="small"
                  :stroke-color="getUtilizationColor(record.utilization)"
                />
                <div class="utilization-text">{{ record.utilization }}%</div>
              </div>
            </template>

            <!-- 资源操作 -->
            <template #resourceAction="{ record }">
              <Button 
                type="primary" 
                size="small" 
                @click="reserveResource(record)"
                :disabled="record.available === 0"
              >
                预约
              </Button>
            </template>
          </Table>
        </Card>
      </Tabs.TabPane>

      <Tabs.TabPane key="calendar" tab="调度日历">
        <Card>
          <ScheduleCalendar :schedules="scheduleList" />
        </Card>
      </Tabs.TabPane>
    </Tabs>

    <!-- 创建调度模态框 -->
    <Modal
      v-model:open="createModalVisible"
      title="创建资源调度"
      width="600px"
      @ok="handleCreateSubmit"
      @cancel="handleCreateCancel"
      :confirm-loading="loading"
    >
      <Form
        ref="formRef"
        :model="scheduleForm"
        :rules="formRules"
        layout="vertical"
      >
        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="任务名称" name="name">
              <Input v-model:value="scheduleForm.name" placeholder="请输入任务名称" />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="集群" name="clusterId">
              <Select v-model:value="scheduleForm.clusterId" placeholder="选择集群">
                <Select.Option 
                  v-for="cluster in clusters" 
                  :key="cluster.id" 
                  :value="cluster.id"
                >
                  {{ cluster.name }}
                </Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="GPU类型" name="gpuType">
              <Select v-model:value="scheduleForm.gpuType" placeholder="选择GPU类型">
                <Select.Option value="Tesla A100">Tesla A100</Select.Option>
                <Select.Option value="Tesla V100">Tesla V100</Select.Option>
                <Select.Option value="Tesla T4">Tesla T4</Select.Option>
              </Select>
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="GPU数量" name="gpuCount">
              <InputNumber
                v-model:value="scheduleForm.gpuCount"
                :min="1"
                :max="32"
                style="width: 100%"
              />
            </Form.Item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="开始日期" name="startDate">
              <DatePicker
                v-model:value="scheduleForm.startDate"
                style="width: 100%"
                placeholder="选择日期"
              />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="开始时间" name="startTime">
              <TimePicker
                v-model:value="scheduleForm.startTime"
                style="width: 100%"
                placeholder="选择时间"
                format="HH:mm"
              />
            </Form.Item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="8">
            <Form.Item label="持续时间(小时)" name="duration">
              <InputNumber
                v-model:value="scheduleForm.duration"
                :min="1"
                :max="168"
                style="width: 100%"
              />
            </Form.Item>
          </Col>
          <Col :span="8">
            <Form.Item label="优先级" name="priority">
              <InputNumber
                v-model:value="scheduleForm.priority"
                :min="1"
                :max="10"
                style="width: 100%"
              />
            </Form.Item>
          </Col>
          <Col :span="8">
            <Form.Item label="调度策略" name="strategy">
              <Select v-model:value="scheduleForm.strategy">
                <Select.Option value="best_fit">最佳适配</Select.Option>
                <Select.Option value="balanced">负载均衡</Select.Option>
                <Select.Option value="gpu_optimized">GPU优化</Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Form.Item label="描述" name="description">
          <Input.TextArea
            v-model:value="scheduleForm.description"
            placeholder="请输入任务描述"
            :rows="3"
          />
        </Form.Item>
      </Form>
    </Modal>
  </div>
</template>

<style scoped lang="scss">
.gpu-schedule-container {
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

.schedule-name {
  .name-main {
    font-weight: 500;
    margin-bottom: 4px;
  }
  
  .name-desc {
    font-size: 12px;
    color: #999;
  }
}

.resource-info {
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

.time-info {
  .time-duration {
    font-size: 12px;
    color: #999;
    margin-top: 2px;
  }
}

.creator-info {
  display: flex;
  align-items: center;
}

.utilization-info {
  .utilization-text {
    font-size: 12px;
    text-align: center;
    margin-top: 4px;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .gpu-schedule-container {
    padding: 16px;
  }
}
</style>
