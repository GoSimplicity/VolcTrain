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
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  SearchOutlined,
  SettingOutlined,
  EyeOutlined,
  ToolOutlined,
  StopOutlined,
  MoreOutlined,
  ThunderboltOutlined,
  DatabaseOutlined,
} from '@ant-design/icons-vue';
import type { GPUDevice, GPUQuery, GPUStatistics } from '#/api/types';
import { 
  getGPUList, 
  getGPUStatistics,
  releaseGPU,
  maintainGPU
} from '#/api';
import { formatDateTime, formatDuration } from '#/utils/date';
import GPUDetailDrawer from './components/GPUDetailDrawer.vue';
import GPUMetricsDrawer from './components/GPUMetricsDrawer.vue';

defineOptions({ name: 'GPUList' });

// 响应式数据
const loading = ref(false);
const gpuList = ref<GPUDevice[]>([]);
const selectedGPU = ref<GPUDevice | null>(null);
const detailDrawerVisible = ref(false);
const metricsDrawerVisible = ref(false);

// 搜索参数
const searchParams = reactive<GPUQuery>({
  page: 1,
  pageSize: 10,
  keyword: '',
  status: undefined,
  brand: undefined,
  clusterId: '',
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
const statistics = ref<GPUStatistics>({
  totalCount: 0,
  byBrand: {
    nvidia: 0,
    amd: 0,
    intel: 0,
  },
  byStatus: {
    available: 0,
    allocated: 0,
    busy: 0,
    maintenance: 0,
    offline: 0,
    error: 0,
  },
  byModel: {},
  utilizationStats: {
    avgGpuUtilization: 0,
    avgMemoryUsage: 0,
    avgTemperature: 0,
    avgPowerUsage: 0,
  },
});

// 模拟数据
const mockGPUs: GPUDevice[] = [
  {
    id: 'gpu-001',
    name: 'Tesla A100',
    nodeId: 'node-01',
    nodeName: 'gpu-node-01',
    clusterId: 'cluster-01',
    clusterName: 'Main Cluster',
    brand: 'nvidia',
    model: 'Tesla A100',
    architecture: 'Ampere',
    cudaCores: 6912,
    memorySize: 80,
    memoryBandwidth: 1935,
    baseClockRate: 765,
    boostClockRate: 1410,
    maxPower: 400,
    status: 'allocated',
    temperature: 68,
    powerUsage: 320,
    memoryUsage: 85,
    gpuUtilization: 95,
    allocatedTo: 'job-001',
    allocatedUser: '张三',
    allocatedTime: '2024-01-20 09:30:00',
    driverVersion: '525.60.13',
    cudaVersion: '12.0',
    createTime: '2024-01-01 00:00:00',
    updateTime: '2024-01-20 09:30:00',
    lastHeartbeat: '2024-01-20 15:30:00',
  },
  {
    id: 'gpu-002',
    name: 'Tesla A100',
    nodeId: 'node-01',
    nodeName: 'gpu-node-01',
    clusterId: 'cluster-01',
    clusterName: 'Main Cluster',
    brand: 'nvidia',
    model: 'Tesla A100',
    architecture: 'Ampere',
    cudaCores: 6912,
    memorySize: 80,
    memoryBandwidth: 1935,
    baseClockRate: 765,
    boostClockRate: 1410,
    maxPower: 400,
    status: 'available',
    temperature: 45,
    powerUsage: 50,
    memoryUsage: 5,
    gpuUtilization: 0,
    driverVersion: '525.60.13',
    cudaVersion: '12.0',
    createTime: '2024-01-01 00:00:00',
    updateTime: '2024-01-20 15:30:00',
    lastHeartbeat: '2024-01-20 15:30:00',
  },
  {
    id: 'gpu-003',
    name: 'Tesla V100',
    nodeId: 'node-02',
    nodeName: 'gpu-node-02',
    clusterId: 'cluster-01',
    clusterName: 'Main Cluster',
    brand: 'nvidia',
    model: 'Tesla V100',
    architecture: 'Volta',
    cudaCores: 5120,
    memorySize: 32,
    memoryBandwidth: 900,
    baseClockRate: 1230,
    boostClockRate: 1530,
    maxPower: 300,
    status: 'error',
    temperature: 0,
    powerUsage: 0,
    memoryUsage: 0,
    gpuUtilization: 0,
    driverVersion: '525.60.13',
    cudaVersion: '12.0',
    createTime: '2024-01-01 00:00:00',
    updateTime: '2024-01-20 12:00:00',
    lastHeartbeat: '2024-01-20 12:00:00',
  },
];

// 表格列定义
const columns = [
  {
    title: 'GPU名称',
    key: 'name',
    slots: { customRender: 'name' },
    width: 180,
  },
  {
    title: '状态',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100,
  },
  {
    title: '使用率',
    key: 'utilization',
    slots: { customRender: 'utilization' },
    width: 120,
  },
  {
    title: '显存',
    key: 'memory',
    slots: { customRender: 'memory' },
    width: 100,
  },
  {
    title: '温度',
    key: 'temperature',
    slots: { customRender: 'temperature' },
    width: 80,
  },
  {
    title: '功耗',
    key: 'power',
    slots: { customRender: 'power' },
    width: 100,
  },
  {
    title: '节点',
    key: 'node',
    slots: { customRender: 'node' },
    width: 120,
  },
  {
    title: '分配信息',
    key: 'allocation',
    slots: { customRender: 'allocation' },
    width: 150,
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

// 数据加载
const loadGPUs = async () => {
  try {
    loading.value = true;
    // const response = await getGPUList(searchParams);
    // gpuList.value = response.data;
    // pagination.total = response.total;
    // pagination.current = response.page;
    // pagination.pageSize = response.pageSize;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    gpuList.value = mockGPUs;
    pagination.total = mockGPUs.length;
    
    // 更新统计数据
    updateStatistics();
  } catch (error) {
    message.error('加载GPU列表失败');
  } finally {
    loading.value = false;
  }
};

const updateStatistics = () => {
  const stats = {
    totalCount: gpuList.value.length,
    byBrand: { nvidia: 0, amd: 0, intel: 0 },
    byStatus: { available: 0, allocated: 0, busy: 0, maintenance: 0, offline: 0, error: 0 },
    byModel: {},
    utilizationStats: {
      avgGpuUtilization: 0,
      avgMemoryUsage: 0,
      avgTemperature: 0,
      avgPowerUsage: 0,
    },
  };
  
  gpuList.value.forEach(gpu => {
    stats.byBrand[gpu.brand]++;
    stats.byStatus[gpu.status]++;
    stats.utilizationStats.avgGpuUtilization += gpu.gpuUtilization;
    stats.utilizationStats.avgMemoryUsage += gpu.memoryUsage;
    stats.utilizationStats.avgTemperature += gpu.temperature;
    stats.utilizationStats.avgPowerUsage += gpu.powerUsage;
  });
  
  const count = gpuList.value.length;
  if (count > 0) {
    stats.utilizationStats.avgGpuUtilization /= count;
    stats.utilizationStats.avgMemoryUsage /= count;
    stats.utilizationStats.avgTemperature /= count;
    stats.utilizationStats.avgPowerUsage /= count;
  }
  
  statistics.value = stats;
};

const refreshData = () => {
  loadGPUs();
};

// 事件处理
const handleSearch = () => {
  searchParams.page = 1;
  pagination.current = 1;
  loadGPUs();
};

const resetSearch = () => {
  Object.assign(searchParams, {
    page: 1,
    pageSize: 10,
    keyword: '',
    status: undefined,
    brand: undefined,
    clusterId: '',
  });
  handleSearch();
};

const handleTableChange = (pag: any) => {
  searchParams.page = pag.current;
  searchParams.pageSize = pag.pageSize;
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  loadGPUs();
};

// GPU操作
const viewGPUDetail = (gpu: GPUDevice) => {
  selectedGPU.value = gpu;
  detailDrawerVisible.value = true;
};

const viewGPUMetrics = (gpu: GPUDevice) => {
  selectedGPU.value = gpu;
  metricsDrawerVisible.value = true;
};

const releaseGPUResource = async (gpu: GPUDevice) => {
  Modal.confirm({
    title: '确认释放GPU',
    content: `确定要释放GPU "${gpu.name}" 吗？这将终止当前任务。`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        await releaseGPU(gpu.id);
        message.success('GPU释放成功');
        loadGPUs();
      } catch (error) {
        message.error('GPU释放失败');
      }
    },
  });
};

const toggleMaintenance = async (gpu: GPUDevice) => {
  const isMaintenance = gpu.status === 'maintenance';
  const action = isMaintenance ? '退出维护' : '进入维护';
  
  Modal.confirm({
    title: `确认${action}`,
    content: `确定要让GPU "${gpu.name}" ${action}模式吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        await maintainGPU(gpu.id, !isMaintenance);
        message.success(`${action}成功`);
        loadGPUs();
      } catch (error) {
        message.error(`${action}失败`);
      }
    },
  });
};

// 初始化
onMounted(() => {
  loadGPUs();
});
</script>

<template>
  <div class="gpu-list-container">
    <!-- 页面头部 -->
    <Card>
      <div class="page-header">
        <div class="header-left">
          <h2>GPU资源管理</h2>
          <p>管理和监控GPU集群资源的使用情况</p>
        </div>
        <div class="header-right">
          <Space>
            <Button @click="refreshData">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary">
              <SettingOutlined />
              集群配置
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
            title="GPU总数"
            :value="statistics.totalCount"
            :value-style="{ color: '#3f8600' }"
            prefix="🔥"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="可用GPU"
            :value="statistics.byStatus.available"
            :value-style="{ color: '#52c41a' }"
            prefix="✅"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="平均使用率"
            :value="statistics.utilizationStats.avgGpuUtilization"
            precision="1"
            suffix="%"
            :value-style="{ color: '#1890ff' }"
            prefix="📊"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="故障GPU"
            :value="statistics.byStatus.error"
            :value-style="{ color: '#cf1322' }"
            prefix="⚠️"
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
            placeholder="搜索GPU名称"
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
            placeholder="状态"
            style="width: 100%"
            @change="handleSearch"
            allow-clear
          >
            <Select.Option value="">全部状态</Select.Option>
            <Select.Option value="available">可用</Select.Option>
            <Select.Option value="allocated">已分配</Select.Option>
            <Select.Option value="busy">忙碌</Select.Option>
            <Select.Option value="maintenance">维护中</Select.Option>
            <Select.Option value="error">故障</Select.Option>
          </Select>
        </Col>
        <Col :span="3">
          <Select
            v-model:value="searchParams.brand"
            placeholder="品牌"
            style="width: 100%"
            @change="handleSearch"
            allow-clear
          >
            <Select.Option value="">全部品牌</Select.Option>
            <Select.Option value="nvidia">NVIDIA</Select.Option>
            <Select.Option value="amd">AMD</Select.Option>
            <Select.Option value="intel">Intel</Select.Option>
          </Select>
        </Col>
        <Col :span="6">
          <Space>
            <Button @click="resetSearch">重置</Button>
          </Space>
        </Col>
      </Row>
    </Card>

    <!-- GPU列表 -->
    <Card>
      <Table
        :columns="columns"
        :data-source="gpuList"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
      >
        <!-- GPU名称 -->
        <template #name="{ record }">
          <div class="gpu-name">
            <div class="name-main">
              <Button type="link" @click="viewGPUDetail(record)">
                {{ record.name }}
              </Button>
              <Tag color="blue" size="small" style="margin-left: 8px">
                {{ record.model }}
              </Tag>
            </div>
            <div class="name-desc">{{ record.architecture }} • {{ record.memorySize }}GB</div>
          </div>
        </template>

        <!-- 状态 -->
        <template #status="{ record }">
          <Tag :color="getGPUStatusColor(record.status)">
            {{ getGPUStatusLabel(record.status) }}
          </Tag>
        </template>

        <!-- 使用率 -->
        <template #utilization="{ record }">
          <div class="utilization-info">
            <Progress
              :percent="record.gpuUtilization"
              size="small"
              :stroke-color="getUtilizationColor(record.gpuUtilization)"
            />
            <div class="utilization-text">{{ record.gpuUtilization }}%</div>
          </div>
        </template>

        <!-- 显存 -->
        <template #memory="{ record }">
          <div class="memory-info">
            <Progress
              :percent="record.memoryUsage"
              size="small"
              stroke-color="#722ed1"
            />
            <div class="memory-text">{{ record.memoryUsage }}%</div>
          </div>
        </template>

        <!-- 温度 -->
        <template #temperature="{ record }">
          <Tooltip :title="`${record.temperature}°C`">
            <Tag :color="getTemperatureColor(record.temperature)">
              {{ record.temperature }}°C
            </Tag>
          </Tooltip>
        </template>

        <!-- 功耗 -->
        <template #power="{ record }">
          <div class="power-info">
            <span class="power-value">{{ record.powerUsage }}W</span>
            <div class="power-detail">/ {{ record.maxPower }}W</div>
          </div>
        </template>

        <!-- 节点 -->
        <template #node="{ record }">
          <div class="node-info">
            <div class="node-name">{{ record.nodeName }}</div>
            <div class="cluster-name">{{ record.clusterName }}</div>
          </div>
        </template>

        <!-- 分配信息 -->
        <template #allocation="{ record }">
          <div v-if="record.allocatedTo" class="allocation-info">
            <div class="allocated-user">
              <Avatar size="small">{{ record.allocatedUser?.[0] }}</Avatar>
              <span style="margin-left: 8px">{{ record.allocatedUser }}</span>
            </div>
            <div class="allocated-time">{{ formatDateTime(record.allocatedTime, 'MM-DD HH:mm') }}</div>
          </div>
          <div v-else class="allocation-info">
            <span class="unallocated">未分配</span>
          </div>
        </template>

        <!-- 操作 -->
        <template #action="{ record }">
          <Space size="small">
            <Button type="link" size="small" @click="viewGPUDetail(record)">
              <EyeOutlined />
            </Button>
            <Button type="link" size="small" @click="viewGPUMetrics(record)">
              <ThunderboltOutlined />
            </Button>
            <Dropdown>
              <Button type="link" size="small">
                <MoreOutlined />
              </Button>
              <template #overlay>
                <Menu>
                  <Menu.Item 
                    key="release" 
                    @click="releaseGPUResource(record)"
                    :disabled="record.status !== 'allocated'"
                  >
                    <StopOutlined />
                    释放GPU
                  </Menu.Item>
                  <Menu.Item 
                    key="maintenance" 
                    @click="toggleMaintenance(record)"
                  >
                    <ToolOutlined />
                    {{ record.status === 'maintenance' ? '退出维护' : '进入维护' }}
                  </Menu.Item>
                </Menu>
              </template>
            </Dropdown>
          </Space>
        </template>
      </Table>
    </Card>

    <!-- GPU详情抽屉 -->
    <GPUDetailDrawer
      v-model:visible="detailDrawerVisible"
      :gpu="selectedGPU"
      @refresh="loadGPUs"
    />

    <!-- GPU监控抽屉 -->
    <GPUMetricsDrawer
      v-model:visible="metricsDrawerVisible"
      :gpu="selectedGPU"
    />
  </div>
</template>

<style scoped lang="scss">
.gpu-list-container {
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

.gpu-name {
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

.utilization-info {
  .utilization-text {
    font-size: 12px;
    text-align: center;
    margin-top: 4px;
  }
}

.memory-info {
  .memory-text {
    font-size: 12px;
    text-align: center;
    margin-top: 4px;
  }
}

.power-info {
  .power-value {
    font-weight: 500;
  }
  
  .power-detail {
    font-size: 12px;
    color: #999;
    margin-top: 2px;
  }
}

.node-info {
  .node-name {
    font-weight: 500;
  }
  
  .cluster-name {
    font-size: 12px;
    color: #999;
    margin-top: 2px;
  }
}

.allocation-info {
  .allocated-user {
    display: flex;
    align-items: center;
    margin-bottom: 4px;
  }
  
  .allocated-time {
    font-size: 12px;
    color: #999;
  }
  
  .unallocated {
    color: #999;
    font-style: italic;
  }
}
</style>
