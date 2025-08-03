<template>
  <div class="workspace-list-container">
    <!-- 页面头部 -->
    <Card>
      <div class="page-header">
        <div class="header-left">
          <h2>工作空间管理</h2>
          <p>管理和查看您的工作空间，创建新的项目空间</p>
        </div>
        <div class="header-right">
          <Space>
            <Button @click="refreshData">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary" @click="showCreateModal">
              <PlusOutlined />
              创建工作空间
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
            title="总工作空间"
            :value="statistics.total"
            :value-style="{ color: '#3f8600' }"
            prefix="📁"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="活跃空间"
            :value="statistics.active"
            :value-style="{ color: '#1890ff' }"
            prefix="🟢"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="总成员数"
            :value="statistics.totalMembers"
            :value-style="{ color: '#722ed1' }"
            prefix="👥"
          />
        </Card>
      </Col>
      <Col :span="6">
        <Card>
          <Statistic
            title="总项目数"
            :value="statistics.totalProjects"
            :value-style="{ color: '#eb2f96' }"
            prefix="🚀"
          />
        </Card>
      </Col>
    </Row>

    <!-- 搜索和筛选 -->
    <Card style="margin-bottom: 16px">
      <Row :gutter="16">
        <Col :span="6">
          <Input
            v-model:value="searchParams.keyword"
            placeholder="搜索工作空间名称"
            @change="handleSearch"
          >
            <template #prefix>
              <SearchOutlined />
            </template>
          </Input>
        </Col>
        <Col :span="4">
          <Select
            v-model:value="searchParams.type"
            placeholder="选择类型"
            style="width: 100%"
            @change="handleSearch"
          >
            <Select.Option value="">全部类型</Select.Option>
            <Select.Option value="personal">个人</Select.Option>
            <Select.Option value="team">团队</Select.Option>
            <Select.Option value="project">项目</Select.Option>
            <Select.Option value="department">部门</Select.Option>
          </Select>
        </Col>
        <Col :span="4">
          <Select
            v-model:value="searchParams.status"
            placeholder="选择状态"
            style="width: 100%"
            @change="handleSearch"
          >
            <Select.Option value="">全部状态</Select.Option>
            <Select.Option value="active">活跃</Select.Option>
            <Select.Option value="inactive">非活跃</Select.Option>
          </Select>
        </Col>
        <Col :span="6">
          <RangePicker
            v-model:value="searchParams.dateRange"
            @change="handleSearch"
            placeholder="['开始时间', '结束时间']"
          />
        </Col>
        <Col :span="4">
          <Button @click="resetSearch">重置</Button>
        </Col>
      </Row>
    </Card>

    <!-- 工作空间列表 -->
    <Card>
      <Table
        :columns="columns"
        :data-source="workspaceList"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
      >
        <!-- 工作空间名称 -->
        <template #name="{ record }">
          <div class="workspace-name">
            <div class="name-main">
              <Button type="link" @click="viewWorkspace(record)">
                {{ record.name }}
              </Button>
              <Tag :color="getTypeColor(record.type)" style="margin-left: 8px">
                {{ getTypeLabel(record.type) }}
              </Tag>
            </div>
            <div class="name-desc">{{ record.description || '暂无描述' }}</div>
          </div>
        </template>

        <!-- 所有者 -->
        <template #owner="{ record }">
          <div class="owner-info">
            <Avatar size="small" :src="record.ownerAvatar">
              {{ record.ownerName?.[0] }}
            </Avatar>
            <span style="margin-left: 8px">{{ record.ownerName }}</span>
          </div>
        </template>

        <!-- 成员数量 -->
        <template #members="{ record }">
          <span>{{ record.memberCount }} 人</span>
        </template>

        <!-- 项目数量 -->
        <template #projects="{ record }">
          <span>{{ record.projectCount }} 个</span>
        </template>

        <!-- 资源使用情况 -->
        <template #resources="{ record }">
          <div class="resource-info">
            <div>
              <Progress
                :percent="getResourceUsagePercent(record, 'cpu')"
                size="small"
                status="active"
              />
              <span style="font-size: 12px">CPU: {{ record.resourceUsed.cpu }}/{{ record.resourceQuota.cpu }}</span>
            </div>
            <div style="margin-top: 4px">
              <Progress
                :percent="getResourceUsagePercent(record, 'memory')"
                size="small"
                status="active"
              />
              <span style="font-size: 12px">内存: {{ record.resourceUsed.memory }}GB/{{ record.resourceQuota.memory }}GB</span>
            </div>
          </div>
        </template>

        <!-- 状态 -->
        <template #status="{ record }">
          <Tag :color="getStatusColor(record.status)">
            {{ getStatusLabel(record.status) }}
          </Tag>
        </template>

        <!-- 最后访问时间 -->
        <template #lastAccess="{ record }">
          <span>{{ formatDateTime(record.lastAccessTime) || '从未访问' }}</span>
        </template>

        <!-- 操作 -->
        <template #action="{ record }">
          <Space size="middle">
            <Button type="link" size="small" @click="viewWorkspace(record)">
              <EyeOutlined />
              查看
            </Button>
            <Button type="link" size="small" @click="editWorkspace(record)">
              <EditOutlined />
              编辑
            </Button>
            <Button type="link" size="small" @click="manageMembers(record)">
              <TeamOutlined />
              成员
            </Button>
            <Dropdown>
              <Button type="link" size="small">
                <MoreOutlined />
              </Button>
              <template #overlay>
                <Menu>
                  <Menu.Item key="clone" @click="cloneWorkspace(record)">
                    <CopyOutlined />
                    克隆
                  </Menu.Item>
                  <Menu.Item key="settings" @click="workspaceSettings(record)">
                    <SettingOutlined />
                    设置
                  </Menu.Item>
                  <Menu.Divider />
                  <Menu.Item key="delete" danger @click="deleteWorkspace(record)">
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

    <!-- 创建工作空间模态框 -->
    <CreateWorkspaceModal
      v-model:visible="createModalVisible"
      @success="handleCreateSuccess"
    />

    <!-- 编辑工作空间模态框 -->
    <EditWorkspaceModal
      v-model:visible="editModalVisible"
      :workspace="selectedWorkspace"
      @success="handleEditSuccess"
    />

    <!-- 成员管理模态框 -->
    <MemberManageModal
      v-model:visible="memberModalVisible"
      :workspace="selectedWorkspace"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Card,
  Row,
  Col,
  Button,
  Space,
  Statistic,
  Input,
  Select,
  DatePicker,
  Table,
  Tag,
  Avatar,
  Progress,
  Dropdown,
  Menu,
  message,
} from 'ant-design-vue';
import {
  PlusOutlined,
  ReloadOutlined,
  SearchOutlined,
  EyeOutlined,
  EditOutlined,
  TeamOutlined,
  MoreOutlined,
  CopyOutlined,
  SettingOutlined,
  DeleteOutlined,
} from '@ant-design/icons-vue';
import type { Workspace, WorkspaceQuery, WorkspaceStatistics } from '#/api/types';
import { getWorkspaceList, getWorkspaceStatistics, deleteWorkspace as deleteWorkspaceApi } from '#/api';
import { formatDateTime } from '#/utils/date';
import CreateWorkspaceModal from './components/CreateWorkspaceModal.vue';
import EditWorkspaceModal from './components/EditWorkspaceModal.vue';
import MemberManageModal from './components/MemberManageModal.vue';

const { RangePicker } = DatePicker;

defineOptions({ name: 'WorkspaceList' });

// 响应式数据
const loading = ref(false);
const workspaceList = ref<Workspace[]>([]);
const selectedWorkspace = ref<Workspace | null>(null);

// 模态框状态
const createModalVisible = ref(false);
const editModalVisible = ref(false);
const memberModalVisible = ref(false);

// 搜索参数
const searchParams = reactive<WorkspaceQuery & { dateRange?: any }>({
  page: 1,
  pageSize: 10,
  keyword: '',
  type: undefined,
  status: undefined,
  dateRange: undefined,
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
const statistics = ref<WorkspaceStatistics>({
  total: 0,
  active: 0,
  inactive: 0,
  lastUpdated: '',
  byType: {
    personal: 0,
    team: 0,
    project: 0,
    department: 0,
  },
  totalMembers: 0,
  totalProjects: 0,
  resourceUtilization: {
    cpu: 0,
    memory: 0,
    gpu: 0,
    storage: 0,
  },
});

// 表格列定义
const columns = [
  {
    title: '工作空间名称',
    dataIndex: 'name',
    key: 'name',
    slots: { customRender: 'name' },
    width: 250,
  },
  {
    title: '所有者',
    dataIndex: 'ownerName',
    key: 'owner',
    slots: { customRender: 'owner' },
    width: 150,
  },
  {
    title: '成员',
    dataIndex: 'memberCount',
    key: 'members',
    slots: { customRender: 'members' },
    width: 80,
  },
  {
    title: '项目',
    dataIndex: 'projectCount',
    key: 'projects',
    slots: { customRender: 'projects' },
    width: 80,
  },
  {
    title: '资源使用',
    key: 'resources',
    slots: { customRender: 'resources' },
    width: 200,
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100,
  },
  {
    title: '最后访问',
    dataIndex: 'lastAccessTime',
    key: 'lastAccess',
    slots: { customRender: 'lastAccess' },
    width: 150,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 200,
    fixed: 'right' as const,
  },
];

// 工具方法
const getTypeColor = (type: string) => {
  const colors = {
    personal: 'blue',
    team: 'green',
    project: 'orange',
    department: 'purple',
  };
  return colors[type as keyof typeof colors] || 'default';
};

const getTypeLabel = (type: string) => {
  const labels = {
    personal: '个人',
    team: '团队',
    project: '项目',
    department: '部门',
  };
  return labels[type as keyof typeof labels] || type;
};

const getStatusColor = (status: string) => {
  return status === 'active' ? 'success' : 'default';
};

const getStatusLabel = (status: string) => {
  return status === 'active' ? '活跃' : '非活跃';
};

const getResourceUsagePercent = (workspace: Workspace, resource: string) => {
  const used = workspace.resourceUsed[resource as keyof typeof workspace.resourceUsed] || 0;
  const quota = workspace.resourceQuota[resource as keyof typeof workspace.resourceQuota] || 1;
  return Math.round((used / quota) * 100);
};

// 数据加载
const loadWorkspaces = async () => {
  try {
    loading.value = true;
    const params = {
      ...searchParams,
      createTimeStart: searchParams.dateRange?.[0],
      createTimeEnd: searchParams.dateRange?.[1],
    };
    delete params.dateRange;
    
    const response = await getWorkspaceList(params);
    workspaceList.value = response.data;
    pagination.total = response.total;
    pagination.current = response.page;
    pagination.pageSize = response.pageSize;
  } catch (error) {
    message.error('加载工作空间列表失败');
  } finally {
    loading.value = false;
  }
};

const loadStatistics = async () => {
  try {
    const response = await getWorkspaceStatistics();
    statistics.value = response;
  } catch (error) {
    console.error('加载统计数据失败:', error);
  }
};

// 事件处理
const handleSearch = () => {
  searchParams.page = 1;
  pagination.current = 1;
  loadWorkspaces();
};

const resetSearch = () => {
  searchParams.keyword = '';
  searchParams.type = undefined;
  searchParams.status = undefined;
  searchParams.dateRange = undefined;
  handleSearch();
};

const handleTableChange = (pag: any) => {
  searchParams.page = pag.current;
  searchParams.pageSize = pag.pageSize;
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  loadWorkspaces();
};

const refreshData = () => {
  loadWorkspaces();
  loadStatistics();
};

// 工作空间操作
const showCreateModal = () => {
  createModalVisible.value = true;
};

const viewWorkspace = (workspace: Workspace) => {
  // 跳转到工作空间详情页
  console.log('查看工作空间:', workspace);
};

const editWorkspace = (workspace: Workspace) => {
  selectedWorkspace.value = workspace;
  editModalVisible.value = true;
};

const manageMembers = (workspace: Workspace) => {
  selectedWorkspace.value = workspace;
  memberModalVisible.value = true;
};

const cloneWorkspace = (workspace: Workspace) => {
  console.log('克隆工作空间:', workspace);
  // 实现克隆逻辑
};

const workspaceSettings = (workspace: Workspace) => {
  console.log('工作空间设置:', workspace);
  // 跳转到设置页面
};

const deleteWorkspace = async (workspace: Workspace) => {
  try {
    await deleteWorkspaceApi(workspace.id);
    message.success('删除成功');
    refreshData();
  } catch (error) {
    message.error('删除失败');
  }
};

// 模态框事件
const handleCreateSuccess = () => {
  createModalVisible.value = false;
  refreshData();
};

const handleEditSuccess = () => {
  editModalVisible.value = false;
  refreshData();
};

// 初始化
onMounted(() => {
  loadWorkspaces();
  loadStatistics();
});
</script>

<style scoped lang="scss">
.workspace-list-container {
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

.workspace-name {
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

.owner-info {
  display: flex;
  align-items: center;
}

.resource-info {
  font-size: 12px;
}
</style>
