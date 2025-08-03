<template>
  <Drawer
    v-model:open="visible"
    title="模型部署"
    width="1000"
    placement="right"
    class="model-deploy-drawer"
  >
    <div v-if="model" class="drawer-content">
      <!-- 部署头部 -->
      <div class="deploy-header">
        <div class="model-info">
          <h3>{{ model.name }}</h3>
          <div class="model-meta">
            <Tag color="blue">{{ model.version }}</Tag>
            <Tag :color="getModelTypeColor(model.type)">
              {{ getModelTypeText(model.type) }}
            </Tag>
            <span class="model-size">{{ formatFileSize(model.size) }}</span>
          </div>
        </div>
        
        <div class="header-actions">
          <Space>
            <Button @click="refreshDeployments" :loading="loading">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary" @click="showDeployModal">
              <RocketOutlined />
              新建部署
            </Button>
          </Space>
        </div>
      </div>

      <Divider style="margin: 16px 0" />

      <!-- 部署统计 -->
      <Row :gutter="16" class="deploy-stats">
        <Col :span="6">
          <Card>
            <Statistic
              title="总部署数"
              :value="deploymentList.length"
              :value-style="{ color: '#3f8600' }"
              prefix="🚀"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card>
            <Statistic
              title="运行中"
              :value="runningDeployments"
              :value-style="{ color: '#52c41a' }"
              prefix="✅"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card>
            <Statistic
              title="总请求数"
              :value="totalRequests"
              :value-style="{ color: '#1890ff' }"
              prefix="📊"
            />
          </Card>
        </Col>
        <Col :span="6">
          <Card>
            <Statistic
              title="资源使用"
              :value="totalResourceUsage"
              suffix="核"
              :value-style="{ color: '#722ed1' }"
              prefix="💻"
            />
          </Card>
        </Col>
      </Row>

      <!-- 部署列表 -->
      <Card title="部署列表" class="deployment-list-card">
        <Table
          :columns="deploymentColumns"
          :data-source="deploymentList"
          :loading="loading"
          row-key="id"
          :pagination="{ pageSize: 10, size: 'small' }"
        >
          <!-- 部署信息 -->
          <template #deploymentInfo="{ record }">
            <div class="deployment-info">
              <div class="deployment-header">
                <span class="deployment-name">{{ record.name }}</span>
                <Tag :color="getEnvironmentColor(record.environment)" size="small">
                  {{ getEnvironmentText(record.environment) }}
                </Tag>
              </div>
              <div class="deployment-desc">{{ record.description || '无描述' }}</div>
              <div class="deployment-endpoint" v-if="record.endpoint">
                <Text code style="font-size: 11px">{{ record.endpoint }}</Text>
              </div>
            </div>
          </template>

          <!-- 状态 -->
          <template #status="{ record }">
            <Badge 
              :status="getDeploymentStatusColor(record.status) as any" 
              :text="getDeploymentStatusText(record.status)"
            />
            <div class="status-detail" v-if="record.status === 'running'">
              <Progress
                :percent="record.health || 100"
                size="small"
                :stroke-color="record.health >= 90 ? '#52c41a' : record.health >= 70 ? '#faad14' : '#ff4d4f'"
                style="margin-top: 4px"
              />
              <div class="health-text">健康度: {{ record.health || 100 }}%</div>
            </div>
          </template>

          <!-- 资源配置 -->
          <template #resources="{ record }">
            <div class="resources-info">
              <div class="resource-item">
                <span class="resource-label">CPU:</span>
                <span class="resource-value">{{ record.resources.cpu }}核</span>
              </div>
              <div class="resource-item">
                <span class="resource-label">内存:</span>
                <span class="resource-value">{{ record.resources.memory }}GB</span>
              </div>
              <div class="resource-item" v-if="record.resources.gpu">
                <span class="resource-label">GPU:</span>
                <span class="resource-value">{{ record.resources.gpu }}块</span>
              </div>
              <div class="resource-item">
                <span class="resource-label">副本:</span>
                <span class="resource-value">{{ record.replicas }}</span>
              </div>
            </div>
          </template>

          <!-- 性能指标 -->
          <template #metrics="{ record }">
            <div class="metrics-info">
              <div class="metric-item">
                <span class="metric-label">请求数:</span>
                <span class="metric-value">{{ record.requestCount || 0 }}</span>
              </div>
              <div class="metric-item">
                <span class="metric-label">平均延迟:</span>
                <span class="metric-value">{{ record.avgLatency || 0 }}ms</span>
              </div>
              <div class="metric-item">
                <span class="metric-label">错误率:</span>
                <span class="metric-value" :style="{ color: record.errorRate > 5 ? '#ff4d4f' : '#52c41a' }">
                  {{ record.errorRate || 0 }}%
                </span>
              </div>
            </div>
          </template>

          <!-- 创建时间 -->
          <template #createTime="{ record }">
            <div class="time-info">
              <div>{{ formatDateTime(record.createTime, 'MM-DD') }}</div>
              <div class="time-detail">{{ formatDateTime(record.createTime, 'HH:mm') }}</div>
            </div>
          </template>

          <!-- 操作 -->
          <template #action="{ record }">
            <Space size="small">
              <Tooltip title="查看详情">
                <Button type="text" size="small" @click="viewDeploymentDetail(record)">
                  <EyeOutlined />
                </Button>
              </Tooltip>
              <Tooltip title="监控面板">
                <Button type="text" size="small" @click="openMonitoringDashboard(record)">
                  <BarChartOutlined />
                </Button>
              </Tooltip>
              <Tooltip title="测试接口">
                <Button 
                  type="text" 
                  size="small" 
                  @click="testDeployment(record)"
                  :disabled="record.status !== 'running'"
                >
                  <ApiOutlined />
                </Button>
              </Tooltip>
              <Dropdown>
                <Button type="text" size="small">
                  <MoreOutlined />
                </Button>
                <template #overlay>
                  <Menu>
                    <Menu.Item 
                      @click="startDeployment(record)"
                      :disabled="record.status === 'running'"
                    >
                      <PlayCircleOutlined />
                      启动
                    </Menu.Item>
                    <Menu.Item 
                      @click="stopDeployment(record)"
                      :disabled="record.status !== 'running'"
                    >
                      <PauseCircleOutlined />
                      停止
                    </Menu.Item>
                    <Menu.Item @click="scaleDeployment(record)">
                      <ExpandOutlined />
                      扩缩容
                    </Menu.Item>
                    <Menu.Item @click="updateDeployment(record)">
                      <EditOutlined />
                      更新配置
                    </Menu.Item>
                    <Menu.Divider />
                    <Menu.Item @click="deleteDeployment(record)" danger>
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
    </div>

    <!-- 新建部署模态框 -->
    <Modal
      v-model:open="deployModalVisible"
      title="新建部署"
      width="800px"
      @ok="handleDeploySubmit"
      @cancel="handleDeployCancel"
      :confirm-loading="deployLoading"
    >
      <Form
        ref="deployFormRef"
        :model="deployForm"
        :rules="deployFormRules"
        layout="vertical"
      >
        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="部署名称" name="name">
              <Input v-model:value="deployForm.name" placeholder="输入部署名称" />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="目标环境" name="targetEnvironment">
              <Select v-model:value="deployForm.targetEnvironment" placeholder="选择环境">
                <Select.Option value="development">开发环境</Select.Option>
                <Select.Option value="staging">测试环境</Select.Option>
                <Select.Option value="production">生产环境</Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="8">
            <Form.Item label="CPU(核)" name="cpu">
              <InputNumber
                v-model:value="deployForm.resources.cpu"
                :min="0.5"
                :max="32"
                :step="0.5"
                style="width: 100%"
              />
            </Form.Item>
          </Col>
          <Col :span="8">
            <Form.Item label="内存(GB)" name="memory">
              <InputNumber
                v-model:value="deployForm.resources.memory"
                :min="1"
                :max="128"
                style="width: 100%"
              />
            </Form.Item>
          </Col>
          <Col :span="8">
            <Form.Item label="GPU(块)" name="gpu">
              <InputNumber
                v-model:value="deployForm.resources.gpu"
                :min="0"
                :max="8"
                style="width: 100%"
              />
            </Form.Item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="副本数量" name="replicas">
              <InputNumber
                v-model:value="deployForm.replicas"
                :min="1"
                :max="10"
                style="width: 100%"
              />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="自动扩缩容" name="autoScaling">
              <Select v-model:value="deployForm.autoScaling" placeholder="选择策略">
                <Select.Option value="disabled">禁用</Select.Option>
                <Select.Option value="cpu">基于CPU</Select.Option>
                <Select.Option value="memory">基于内存</Select.Option>
                <Select.Option value="requests">基于请求数</Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Form.Item label="描述" name="description">
          <Input.TextArea
            v-model:value="deployForm.description"
            placeholder="输入部署描述"
            :rows="3"
          />
        </Form.Item>

        <Form.Item label="环境变量" name="envVars">
          <div class="env-vars-editor">
            <div v-for="(envVar, index) in deployForm.envVars" :key="index" class="env-var-row">
              <Input
                v-model:value="envVar.key"
                placeholder="变量名"
                style="width: 40%; margin-right: 8px"
              />
              <Input
                v-model:value="envVar.value"
                placeholder="变量值"
                style="width: 40%; margin-right: 8px"
              />
              <Button @click="removeEnvVar(index)" size="small" danger>
                <DeleteOutlined />
              </Button>
            </div>
            <Button @click="addEnvVar" type="dashed" style="width: 100%; margin-top: 8px">
              <PlusOutlined />
              添加环境变量
            </Button>
          </div>
        </Form.Item>

        <Form.Item label="配置信息" name="configurationText">
          <Input.TextArea
            v-model:value="deployForm.configurationText"
            placeholder='JSON格式的配置信息，例如: {"timeout": 30, "batch_size": 1}'
            :rows="4"
          />
        </Form.Item>
      </Form>
    </Modal>

    <!-- 扩缩容模态框 -->
    <Modal
      v-model:open="scaleModalVisible"
      title="扩缩容"
      width="400px"
      @ok="handleScaleSubmit"
      @cancel="handleScaleCancel"
      :confirm-loading="scaleLoading"
    >
      <Form
        :model="scaleForm"
        layout="vertical"
      >
        <Form.Item label="副本数量">
          <InputNumber
            v-model:value="scaleForm.replicas"
            :min="1"
            :max="10"
            style="width: 100%"
          />
          <div style="margin-top: 8px; color: #999; font-size: 12px">
            当前副本数: {{ selectedDeployment?.replicas || 0 }}
          </div>
        </Form.Item>
      </Form>
    </Modal>

    <!-- 部署详情模态框 -->
    <Modal
      v-model:open="deployDetailModalVisible"
      title="部署详情"
      width="900px"
      :footer="null"
    >
      <div v-if="selectedDeployment" class="deployment-detail">
        <Tabs>
          <Tabs.TabPane key="overview" tab="概览">
            <Descriptions :column="2" bordered>
              <Descriptions.Item label="部署名称">
                {{ selectedDeployment.name }}
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                <Badge 
                  :status="getDeploymentStatusColor(selectedDeployment.status) as any" 
                  :text="getDeploymentStatusText(selectedDeployment.status)"
                />
              </Descriptions.Item>
              <Descriptions.Item label="环境">
                <Tag :color="getEnvironmentColor(selectedDeployment.environment)">
                  {{ getEnvironmentText(selectedDeployment.environment) }}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="副本数">
                {{ selectedDeployment.replicas }}
              </Descriptions.Item>
              <Descriptions.Item label="服务端点" :span="2">
                <Text code v-if="selectedDeployment.endpoint">
                  {{ selectedDeployment.endpoint }}
                </Text>
                <Text type="secondary" v-else>未分配</Text>
              </Descriptions.Item>
              <Descriptions.Item label="创建时间" :span="2">
                {{ formatDateTime(selectedDeployment.createTime) }}
              </Descriptions.Item>
            </Descriptions>
          </Tabs.TabPane>
          
          <Tabs.TabPane key="resources" tab="资源配置">
            <Card title="资源分配">
              <Row :gutter="16">
                <Col :span="8">
                  <Statistic
                    title="CPU"
                    :value="selectedDeployment.resources.cpu"
                    suffix="核"
                  />
                </Col>
                <Col :span="8">
                  <Statistic
                    title="内存"
                    :value="selectedDeployment.resources.memory"
                    suffix="GB"
                  />
                </Col>
                <Col :span="8">
                  <Statistic
                    title="GPU"
                    :value="selectedDeployment.resources.gpu || 0"
                    suffix="块"
                  />
                </Col>
              </Row>
            </Card>
          </Tabs.TabPane>
          
          <Tabs.TabPane key="monitoring" tab="监控">
            <Alert
              message="监控面板"
              description="这里可以展示部署的详细监控信息，包括CPU使用率、内存使用率、请求量、响应时间等指标的实时图表。"
              type="info"
              show-icon
            />
          </Tabs.TabPane>
        </Tabs>
      </div>
    </Modal>
  </Drawer>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, onMounted } from 'vue';
import {
  Drawer,
  Space,
  Button,
  Divider,
  Row,
  Col,
  Card,
  Statistic,
  Table,
  Tag,
  Badge,
  Progress,
  Tooltip,
  Dropdown,
  Menu,
  Modal,
  Form,
  Input,
  InputNumber,
  Select,
  Descriptions,
  Typography,
  Tabs,
  Alert,
  message,
} from 'ant-design-vue';
import {
  ReloadOutlined,
  RocketOutlined,
  EyeOutlined,
  BarChartOutlined,
  ApiOutlined,
  MoreOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  ExpandOutlined,
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue';
import type { Model, ModelType } from '#/api/types';
import { formatDateTime, formatFileSize } from '#/utils/date';

const { Text } = Typography;

const props = defineProps<{
  visible: boolean;
  model: Model | null;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
}>();

// 响应式数据
const loading = ref(false);
const deployModalVisible = ref(false);
const scaleModalVisible = ref(false);
const deployDetailModalVisible = ref(false);
const deployLoading = ref(false);
const scaleLoading = ref(false);
const deploymentList = ref<any[]>([]);
const selectedDeployment = ref<any>(null);

// 部署表单
interface DeployForm {
  name: string;
  targetEnvironment: string;
  replicas: number;
  resources: {
    cpu: number;
    memory: number;
    gpu: number;
  };
  autoScaling: string;
  description: string;
  envVars: Array<{ key: string; value: string }>;
  configurationText: string;
}

const deployForm = reactive<DeployForm>({
  name: '',
  targetEnvironment: 'development',
  replicas: 1,
  resources: {
    cpu: 2,
    memory: 4,
    gpu: 0,
  },
  autoScaling: 'disabled',
  description: '',
  envVars: [],
  configurationText: '',
});

const scaleForm = reactive({
  replicas: 1,
});

const deployFormRef = ref();

// 模拟部署数据
const mockDeployments = [
  {
    id: 'deploy-001',
    name: 'bert-chinese-prod',
    modelId: 'model-001',
    modelVersion: 'v1.2.0',
    environment: 'production',
    status: 'running',
    health: 98,
    replicas: 3,
    resources: {
      cpu: 4,
      memory: 8,
      gpu: 1,
    },
    endpoint: 'https://api.example.com/models/bert-chinese',
    requestCount: 15420,
    avgLatency: 45,
    errorRate: 0.2,
    autoScaling: 'cpu',
    description: '生产环境BERT模型部署',
    createTime: '2024-01-20 10:00:00',
  },
  {
    id: 'deploy-002',
    name: 'bert-chinese-staging',
    modelId: 'model-001',
    modelVersion: 'v1.2.0',
    environment: 'staging',
    status: 'running',
    health: 95,
    replicas: 1,
    resources: {
      cpu: 2,
      memory: 4,
      gpu: 0,
    },
    endpoint: 'https://staging-api.example.com/models/bert-chinese',
    requestCount: 234,
    avgLatency: 52,
    errorRate: 1.2,
    autoScaling: 'disabled',
    description: '测试环境BERT模型部署',
    createTime: '2024-01-19 14:30:00',
  },
  {
    id: 'deploy-003',
    name: 'bert-chinese-dev',
    modelId: 'model-001',
    modelVersion: 'v1.1.0',
    environment: 'development',
    status: 'stopped',
    health: 0,
    replicas: 1,
    resources: {
      cpu: 1,
      memory: 2,
      gpu: 0,
    },
    endpoint: null,
    requestCount: 45,
    avgLatency: 68,
    errorRate: 2.1,
    autoScaling: 'disabled',
    description: '开发环境BERT模型部署',
    createTime: '2024-01-18 09:15:00',
  },
];

// 表格列定义
const deploymentColumns = [
  {
    title: '部署信息',
    key: 'deploymentInfo',
    slots: { customRender: 'deploymentInfo' },
    width: 200,
  },
  {
    title: '状态',
    key: 'status',
    slots: { customRender: 'status' },
    width: 120,
  },
  {
    title: '资源配置',
    key: 'resources',
    slots: { customRender: 'resources' },
    width: 150,
  },
  {
    title: '性能指标',
    key: 'metrics',
    slots: { customRender: 'metrics' },
    width: 150,
  },
  {
    title: '创建时间',
    key: 'createTime',
    slots: { customRender: 'createTime' },
    width: 120,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 180,
    fixed: 'right' as const,
  },
];

// 计算属性
const runningDeployments = computed(() => {
  return deploymentList.value.filter(d => d.status === 'running').length;
});

const totalRequests = computed(() => {
  return deploymentList.value.reduce((sum, d) => sum + (d.requestCount || 0), 0);
});

const totalResourceUsage = computed(() => {
  return deploymentList.value.reduce((sum, d) => sum + d.resources.cpu, 0);
});

// 工具方法
const getModelTypeText = (type: ModelType) => {
  const types = {
    [ModelType.CLASSIFICATION]: '分类',
    [ModelType.REGRESSION]: '回归',
    [ModelType.OBJECT_DETECTION]: '目标检测',
    [ModelType.SEMANTIC_SEGMENTATION]: '语义分割',
    [ModelType.NLP]: '自然语言处理',
    [ModelType.RECOMMENDATION]: '推荐系统',
    [ModelType.GENERATIVE]: '生成模型',
    [ModelType.CUSTOM]: '自定义',
  };
  return types[type] || type;
};

const getModelTypeColor = (type: ModelType) => {
  const colors = {
    [ModelType.CLASSIFICATION]: 'blue',
    [ModelType.REGRESSION]: 'green',
    [ModelType.OBJECT_DETECTION]: 'orange',
    [ModelType.SEMANTIC_SEGMENTATION]: 'purple',
    [ModelType.NLP]: 'cyan',
    [ModelType.RECOMMENDATION]: 'magenta',
    [ModelType.GENERATIVE]: 'red',
    [ModelType.CUSTOM]: 'default',
  };
  return colors[type] || 'default';
};

const getEnvironmentText = (env: string) => {
  const envs = {
    development: '开发环境',
    staging: '测试环境',
    production: '生产环境',
  };
  return envs[env as keyof typeof envs] || env;
};

const getEnvironmentColor = (env: string) => {
  const colors = {
    development: 'default',
    staging: 'orange',
    production: 'red',
  };
  return colors[env as keyof typeof colors] || 'default';
};

const getDeploymentStatusText = (status: string) => {
  const statuses = {
    running: '运行中',
    stopped: '已停止',
    starting: '启动中',
    stopping: '停止中',
    failed: '失败',
    updating: '更新中',
  };
  return statuses[status as keyof typeof statuses] || status;
};

const getDeploymentStatusColor = (status: string) => {
  const colors = {
    running: 'success',
    stopped: 'default',
    starting: 'processing',
    stopping: 'warning',
    failed: 'error',
    updating: 'processing',
  };
  return colors[status as keyof typeof colors] || 'default';
};

// 数据加载
const loadDeployments = async () => {
  if (!props.model) return;
  
  try {
    loading.value = true;
    // const response = await getModelDeployments(props.model.id);
    // deploymentList.value = response.data.items;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    deploymentList.value = mockDeployments;
  } catch (error) {
    message.error('加载部署列表失败');
  } finally {
    loading.value = false;
  }
};

const refreshDeployments = () => {
  loadDeployments();
};

// 事件处理
const showDeployModal = () => {
  deployModalVisible.value = true;
  resetDeployForm();
};

const resetDeployForm = () => {
  Object.assign(deployForm, {
    name: `${props.model?.name}-deploy-${Date.now()}`,
    targetEnvironment: 'development',
    replicas: 1,
    resources: {
      cpu: 2,
      memory: 4,
      gpu: 0,
    },
    autoScaling: 'disabled',
    description: '',
    envVars: [],
    configurationText: '',
  });
};

const addEnvVar = () => {
  deployForm.envVars.push({ key: '', value: '' });
};

const removeEnvVar = (index: number) => {
  deployForm.envVars.splice(index, 1);
};

const handleDeploySubmit = async () => {
  try {
    await deployFormRef.value?.validate();
    
    deployLoading.value = true;
    
    // 解析配置信息
    let configuration = null;
    if (deployForm.configurationText.trim()) {
      try {
        configuration = JSON.parse(deployForm.configurationText);
      } catch (error) {
        message.error('配置信息格式错误，请使用有效的JSON格式');
        return;
      }
    }
    
    // const request = {
    //   modelId: props.model!.id,
    //   version: props.model!.version,
    //   targetEnvironment: deployForm.targetEnvironment,
    //   replicas: deployForm.replicas,
    //   resources: deployForm.resources,
    //   configuration,
    // };
    
    // const response = await deployModel(request);
    
    // 模拟部署成功
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    message.success('模型部署成功');
    deployModalVisible.value = false;
    loadDeployments();
  } catch (error) {
    message.error('部署失败');
  } finally {
    deployLoading.value = false;
  }
};

const handleDeployCancel = () => {
  deployModalVisible.value = false;
};

const viewDeploymentDetail = (deployment: any) => {
  selectedDeployment.value = deployment;
  deployDetailModalVisible.value = true;
};

const openMonitoringDashboard = (deployment: any) => {
  message.info('打开监控面板功能开发中');
};

const testDeployment = (deployment: any) => {
  message.info('接口测试功能开发中');
};

const startDeployment = async (deployment: any) => {
  try {
    // await startModelDeployment(deployment.id);
    message.success('部署启动成功');
    loadDeployments();
  } catch (error) {
    message.error('启动失败');
  }
};

const stopDeployment = async (deployment: any) => {
  Modal.confirm({
    title: '确认停止',
    content: `确定要停止部署 "${deployment.name}" 吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await stopModelDeployment(deployment.id);
        message.success('部署停止成功');
        loadDeployments();
      } catch (error) {
        message.error('停止失败');
      }
    },
  });
};

const scaleDeployment = (deployment: any) => {
  selectedDeployment.value = deployment;
  scaleForm.replicas = deployment.replicas;
  scaleModalVisible.value = true;
};

const handleScaleSubmit = async () => {
  try {
    scaleLoading.value = true;
    // await scaleModelDeployment(selectedDeployment.value.id, scaleForm.replicas);
    
    // 模拟扩缩容成功
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    message.success('扩缩容操作成功');
    scaleModalVisible.value = false;
    loadDeployments();
  } catch (error) {
    message.error('扩缩容失败');
  } finally {
    scaleLoading.value = false;
  }
};

const handleScaleCancel = () => {
  scaleModalVisible.value = false;
};

const updateDeployment = (deployment: any) => {
  message.info('更新配置功能开发中');
};

const deleteDeployment = async (deployment: any) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除部署 "${deployment.name}" 吗？此操作不可恢复。`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await deleteModelDeployment(deployment.id);
        message.success('部署删除成功');
        loadDeployments();
      } catch (error) {
        message.error('删除失败');
      }
    },
  });
};

// 表单验证规则
const deployFormRules = {
  name: [
    { required: true, message: '请输入部署名称', trigger: 'blur' },
  ],
  targetEnvironment: [
    { required: true, message: '请选择目标环境', trigger: 'change' },
  ],
  replicas: [
    { required: true, message: '请输入副本数量', trigger: 'blur' },
  ],
};

// 监听模型变化
import { watch } from 'vue';
watch(() => props.model, (newModel) => {
  if (newModel && props.visible) {
    loadDeployments();
  }
});

// 监听visible变化
watch(() => props.visible, (newVal) => {
  if (newVal && props.model) {
    loadDeployments();
  }
});

// 初始化
onMounted(() => {
  if (props.visible && props.model) {
    loadDeployments();
  }
});
</script>

<style scoped lang="scss">
.model-deploy-drawer {
  :deep(.ant-drawer-body) {
    padding: 0;
    display: flex;
    flex-direction: column;
  }
}

.drawer-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 24px;
  overflow-y: auto;
}

.deploy-header {
  flex-shrink: 0;
  
  .model-info {
    margin-bottom: 16px;
    
    h3 {
      margin: 0 0 8px 0;
      color: #1890ff;
    }
    
    .model-meta {
      display: flex;
      align-items: center;
      gap: 8px;
      
      .model-size {
        font-size: 12px;
        color: #999;
      }
    }
  }
  
  .header-actions {
    display: flex;
    justify-content: flex-end;
  }
}

.deploy-stats {
  margin-bottom: 24px;
  
  .ant-statistic {
    text-align: center;
  }
}

.deployment-list-card {
  flex: 1;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
  
  :deep(.ant-card-head-title) {
    font-weight: 600;
    color: #1890ff;
    font-size: 14px;
  }
}

.deployment-info {
  .deployment-header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 4px;
    
    .deployment-name {
      font-weight: 600;
      color: #1890ff;
    }
  }
  
  .deployment-desc {
    color: #666;
    font-size: 12px;
    margin-bottom: 4px;
  }
  
  .deployment-endpoint {
    font-size: 11px;
  }
}

.status-detail {
  margin-top: 4px;
  
  .health-text {
    font-size: 11px;
    color: #666;
    text-align: center;
    margin-top: 2px;
  }
}

.resources-info {
  .resource-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 2px;
    
    .resource-label {
      font-size: 12px;
      color: #666;
    }
    
    .resource-value {
      font-size: 12px;
      font-weight: 500;
    }
  }
}

.metrics-info {
  .metric-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 2px;
    
    .metric-label {
      font-size: 12px;
      color: #666;
    }
    
    .metric-value {
      font-size: 12px;
      font-weight: 500;
    }
  }
}

.time-info {
  .time-detail {
    font-size: 12px;
    color: #999;
    margin-top: 2px;
  }
}

.env-vars-editor {
  .env-var-row {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
  }
}

.deployment-detail {
  .ant-descriptions {
    margin-top: 16px;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .drawer-content {
    padding: 16px;
  }
  
  .deploy-header {
    .model-meta {
      flex-direction: column;
      align-items: flex-start;
      gap: 4px;
    }
    
    .header-actions {
      justify-content: flex-start;
      margin-top: 12px;
    }
  }
  
  .deploy-stats {
    :deep(.ant-col) {
      margin-bottom: 12px;
    }
  }
}
</style>