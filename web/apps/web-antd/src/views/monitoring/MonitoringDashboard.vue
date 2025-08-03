<template>
  <div class="monitoring-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <DashboardOutlined class="title-icon" />
            <span class="title-text">监控告警</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">实时监控系统状态和管理告警规则</span>
          </p>
        </div>
        <div class="action-section">
          <Space>
            <Button @click="refreshData" :loading="loading">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary" @click="showCreateRuleModal">
              <PlusOutlined />
              创建规则
            </Button>
            <Button @click="showDashboardModal">
              <BarChartOutlined />
              仪表板
            </Button>
          </Space>
        </div>
      </div>
    </div>

    <!-- 系统健康概览 -->
    <div class="health-overview">
      <Row :gutter="16">
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="health-card">
            <div class="health-item">
              <div class="health-icon" :class="`health-${systemHealth.overall}`">
                <CheckCircleOutlined v-if="systemHealth.overall === 'healthy'" />
                <ExclamationCircleOutlined v-else-if="systemHealth.overall === 'warning'" />
                <CloseCircleOutlined v-else />
              </div>
              <div class="health-info">
                <div class="health-title">系统状态</div>
                <div class="health-value">{{ getHealthText(systemHealth.overall) }}</div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="health-card">
            <div class="health-item">
              <div class="health-progress">
                <Progress
                  type="circle"
                  :percent="systemHealth.cpu"
                  :size="60"
                  :stroke-color="getProgressColor(systemHealth.cpu)"
                />
              </div>
              <div class="health-info">
                <div class="health-title">CPU使用率</div>
                <div class="health-value">{{ systemHealth.cpu }}%</div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="health-card">
            <div class="health-item">
              <div class="health-progress">
                <Progress
                  type="circle"
                  :percent="systemHealth.memory"
                  :size="60"
                  :stroke-color="getProgressColor(systemHealth.memory)"
                />
              </div>
              <div class="health-info">
                <div class="health-title">内存使用率</div>
                <div class="health-value">{{ systemHealth.memory }}%</div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="health-card">
            <div class="health-item">
              <div class="health-progress">
                <Progress
                  type="circle"
                  :percent="systemHealth.disk"
                  :size="60"
                  :stroke-color="getProgressColor(systemHealth.disk)"
                />
              </div>
              <div class="health-info">
                <div class="health-title">磁盘使用率</div>
                <div class="health-value">{{ systemHealth.disk }}%</div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <Row :gutter="16">
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="活跃告警"
              :value="statistics.activeAlerts"
              :value-style="{ color: statistics.activeAlerts > 0 ? '#ff4d4f' : '#52c41a' }"
              prefix="🚨"
            />
            <div class="stat-detail">
              <span class="critical-count">严重: {{ statistics.criticalAlerts }}</span>
              <span class="warning-count">警告: {{ statistics.warningAlerts }}</span>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="监控指标"
              :value="statistics.totalMetrics"
              :value-style="{ color: '#1890ff' }"
              prefix="📊"
            />
            <div class="stat-detail">
              <span>已启用规则: {{ statistics.enabledRules }}</span>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="仪表板"
              :value="statistics.totalDashboards"
              :value-style="{ color: '#722ed1' }"
              prefix="📈"
            />
            <div class="stat-detail">
              <span>通知渠道: {{ statistics.notificationChannels }}</span>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="stat-card">
            <Statistic
              title="今日告警"
              :value="statistics.resolvedAlerts"
              :value-style="{ color: '#52c41a' }"
              prefix="✅"
            />
            <div class="stat-detail">
              <span>已处理: {{ statistics.resolvedAlerts }}</span>
            </div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <Row :gutter="16">
        <!-- 告警列表 -->
        <Col :xs="24" :lg="16">
          <Card title="最新告警" class="alerts-card">
            <template #extra>
              <Space>
                <Select
                  v-model:value="alertFilter.level"
                  placeholder="级别"
                  style="width: 100px"
                  allow-clear
                  size="small"
                  @change="filterAlerts"
                >
                  <Select.Option value="critical">严重</Select.Option>
                  <Select.Option value="error">错误</Select.Option>
                  <Select.Option value="warning">警告</Select.Option>
                  <Select.Option value="info">信息</Select.Option>
                </Select>
                <Select
                  v-model:value="alertFilter.status"
                  placeholder="状态"
                  style="width: 100px"
                  allow-clear
                  size="small"
                  @change="filterAlerts"
                >
                  <Select.Option value="firing">触发中</Select.Option>
                  <Select.Option value="pending">待处理</Select.Option>
                  <Select.Option value="resolved">已解决</Select.Option>
                  <Select.Option value="silenced">已静默</Select.Option>
                </Select>
                <Button size="small" @click="loadAlerts">
                  <ReloadOutlined />
                </Button>
              </Space>
            </template>

            <div class="alerts-list">
              <div
                v-for="alert in filteredAlerts"
                :key="alert.id"
                class="alert-item"
                :class="`alert-${alert.level}`"
              >
                <div class="alert-header">
                  <div class="alert-info">
                    <Badge
                      :status="getAlertStatusBadge(alert.status)"
                      :text="alert.summary"
                      class="alert-summary"
                    />
                    <Tag
                      :color="getAlertLevelColor(alert.level)"
                      size="small"
                      class="alert-level"
                    >
                      {{ getAlertLevelText(alert.level) }}
                    </Tag>
                  </div>
                  <div class="alert-time">
                    {{ formatRelativeTime(alert.startsAt) }}
                  </div>
                </div>
                <div class="alert-description">{{ alert.description }}</div>
                <div class="alert-details">
                  <Space size="small">
                    <span class="alert-metric">
                      指标: {{ getMetricTypeText(alert.metricType) }}
                    </span>
                    <span class="alert-value">
                      当前值: {{ alert.currentValue }}
                    </span>
                    <span class="alert-threshold">
                      阈值: {{ alert.threshold }}
                    </span>
                  </Space>
                </div>
                <div class="alert-actions">
                  <Space size="small">
                    <Button size="small" @click="acknowledgeAlert(alert.id)">
                      确认
                    </Button>
                    <Button size="small" @click="silenceAlert(alert.id)">
                      静默
                    </Button>
                    <Button size="small" @click="resolveAlert(alert.id)">
                      解决
                    </Button>
                    <Button size="small" @click="viewAlertDetail(alert)">
                      详情
                    </Button>
                  </Space>
                </div>
              </div>
              
              <div v-if="filteredAlerts.length === 0" class="empty-alerts">
                <Empty description="暂无告警" />
              </div>
            </div>
          </Card>
        </Col>

        <!-- 告警趋势图表 -->
        <Col :xs="24" :lg="8">
          <Card title="告警趋势" class="trend-card">
            <div class="trend-chart" ref="trendChartRef"></div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 告警规则管理 -->
    <div class="rules-section">
      <Card title="告警规则" class="rules-card">
        <template #extra>
          <Space>
            <Button size="small" @click="showCreateRuleModal">
              <PlusOutlined />
              创建规则
            </Button>
            <Button size="small" @click="showRulesManagement">
              <SettingOutlined />
              管理
            </Button>
          </Space>
        </template>

        <Table
          :columns="ruleColumns"
          :data-source="alertRules"
          :loading="rulesLoading"
          :pagination="{ pageSize: 10, size: 'small' }"
          row-key="id"
          size="small"
        >
          <!-- 规则信息 -->
          <template #ruleInfo="{ record }">
            <div class="rule-info">
              <div class="rule-name">{{ record.name }}</div>
              <div class="rule-desc">{{ record.description || '无描述' }}</div>
            </div>
          </template>

          <!-- 规则类型 -->
          <template #type="{ record }">
            <Tag color="blue">
              {{ getRuleTypeText(record.type) }}
            </Tag>
          </template>

          <!-- 监控指标 -->
          <template #metric="{ record }">
            <Tag :color="getMetricTypeColor(record.metricType)">
              {{ getMetricTypeText(record.metricType) }}
            </Tag>
          </template>

          <!-- 告警级别 -->
          <template #level="{ record }">
            <Tag :color="getAlertLevelColor(record.level)">
              {{ getAlertLevelText(record.level) }}
            </Tag>
          </template>

          <!-- 状态 -->
          <template #status="{ record }">
            <Switch
              :checked="record.enabled"
              @change="toggleRule(record.id, $event)"
              size="small"
            />
          </template>

          <!-- 操作 -->
          <template #action="{ record }">
            <Space size="small">
              <Button type="text" size="small" @click="editRule(record)">
                <EditOutlined />
              </Button>
              <Button type="text" size="small" @click="testRule(record)">
                <PlayCircleOutlined />
              </Button>
              <Button type="text" size="small" danger @click="deleteRule(record)">
                <DeleteOutlined />
              </Button>
            </Space>
          </template>
        </Table>
      </Card>
    </div>

    <!-- 创建告警规则模态框 -->
    <Modal
      v-model:open="createRuleModalVisible"
      title="创建告警规则"
      width="800px"
      @ok="handleCreateRule"
      @cancel="handleCreateRuleCancel"
      :confirm-loading="createRuleLoading"
    >
      <Form
        ref="createRuleFormRef"
        :model="createRuleForm"
        :rules="createRuleFormRules"
        layout="vertical"
      >
        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="规则名称" name="name">
              <Input v-model:value="createRuleForm.name" placeholder="输入规则名称" />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="告警级别" name="level">
              <Select v-model:value="createRuleForm.level" placeholder="选择告警级别">
                <Select.Option value="info">信息</Select.Option>
                <Select.Option value="warning">警告</Select.Option>
                <Select.Option value="error">错误</Select.Option>
                <Select.Option value="critical">严重</Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="监控指标" name="metricType">
              <Select v-model:value="createRuleForm.metricType" placeholder="选择监控指标">
                <Select.Option value="cpu">CPU</Select.Option>
                <Select.Option value="memory">内存</Select.Option>
                <Select.Option value="gpu">GPU</Select.Option>
                <Select.Option value="disk">磁盘</Select.Option>
                <Select.Option value="network">网络</Select.Option>
                <Select.Option value="temperature">温度</Select.Option>
                <Select.Option value="power">功耗</Select.Option>
              </Select>
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="规则类型" name="type">
              <Select v-model:value="createRuleForm.type" placeholder="选择规则类型">
                <Select.Option value="threshold">阈值</Select.Option>
                <Select.Option value="anomaly">异常检测</Select.Option>
                <Select.Option value="composite">复合</Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Form.Item label="查询表达式" name="query">
          <Input.TextArea
            v-model:value="createRuleForm.query"
            placeholder="输入监控查询表达式，如: avg(cpu_usage) > 80"
            :rows="3"
          />
        </Form.Item>

        <Form.Item label="触发条件">
          <div class="conditions-editor">
            <div v-for="(condition, index) in createRuleForm.conditions" :key="index" class="condition-row">
              <Select
                v-model:value="condition.operator"
                style="width: 80px; margin-right: 8px"
                size="small"
              >
                <Select.Option value="gt">></Select.Option>
                <Select.Option value="gte">≥</Select.Option>
                <Select.Option value="lt"><</Select.Option>
                <Select.Option value="lte">≤</Select.Option>
                <Select.Option value="eq">=</Select.Option>
                <Select.Option value="ne">≠</Select.Option>
              </Select>
              <InputNumber
                v-model:value="condition.threshold"
                placeholder="阈值"
                style="width: 120px; margin-right: 8px"
                size="small"
              />
              <InputNumber
                v-model:value="condition.duration"
                placeholder="持续时间(秒)"
                style="width: 130px; margin-right: 8px"
                size="small"
              />
              <Button size="small" danger @click="removeCondition(index)">
                <DeleteOutlined />
              </Button>
            </div>
            <Button @click="addCondition" type="dashed" size="small" style="width: 100%; margin-top: 8px">
              <PlusOutlined />
              添加条件
            </Button>
          </div>
        </Form.Item>

        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="评估间隔(秒)" name="evaluationInterval">
              <InputNumber
                v-model:value="createRuleForm.evaluationInterval"
                :min="10"
                :max="3600"
                style="width: 100%"
              />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="通知渠道" name="notificationChannels">
              <Select
                v-model:value="createRuleForm.notificationChannels"
                mode="multiple"
                placeholder="选择通知渠道"
                style="width: 100%"
              >
                <Select.Option
                  v-for="channel in notificationChannels"
                  :key="channel.id"
                  :value="channel.id"
                >
                  {{ channel.name }}
                </Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Form.Item label="描述" name="description">
          <Input.TextArea
            v-model:value="createRuleForm.description"
            placeholder="输入规则描述"
            :rows="2"
          />
        </Form.Item>
      </Form>
    </Modal>

    <!-- 告警详情模态框 -->
    <Modal
      v-model:open="alertDetailModalVisible"
      title="告警详情"
      width="600px"
      :footer="null"
    >
      <div v-if="selectedAlert" class="alert-detail">
        <Descriptions :column="1" bordered>
          <Descriptions.Item label="告警摘要">
            {{ selectedAlert.summary }}
          </Descriptions.Item>
          <Descriptions.Item label="告警级别">
            <Tag :color="getAlertLevelColor(selectedAlert.level)">
              {{ getAlertLevelText(selectedAlert.level) }}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Badge
              :status="getAlertStatusBadge(selectedAlert.status)"
              :text="getAlertStatusText(selectedAlert.status)"
            />
          </Descriptions.Item>
          <Descriptions.Item label="监控指标">
            {{ getMetricTypeText(selectedAlert.metricType) }}
          </Descriptions.Item>
          <Descriptions.Item label="当前值">
            {{ selectedAlert.currentValue }}
          </Descriptions.Item>
          <Descriptions.Item label="阈值">
            {{ selectedAlert.threshold }}
          </Descriptions.Item>
          <Descriptions.Item label="开始时间">
            {{ formatDateTime(selectedAlert.startsAt) }}
          </Descriptions.Item>
          <Descriptions.Item label="持续时间">
            {{ formatDuration(selectedAlert.duration) }}
          </Descriptions.Item>
          <Descriptions.Item label="描述">
            {{ selectedAlert.description }}
          </Descriptions.Item>
        </Descriptions>

        <div style="margin-top: 16px">
          <h4>操作历史</h4>
          <Timeline>
            <Timeline.Item
              v-for="action in selectedAlert.actions"
              :key="action.timestamp"
              :color="getActionColor(action.type)"
            >
              <div class="action-item">
                <div class="action-header">
                  <span class="action-type">{{ getActionText(action.type) }}</span>
                  <span class="action-user">{{ action.userName }}</span>
                  <span class="action-time">{{ formatDateTime(action.timestamp) }}</span>
                </div>
                <div v-if="action.comment" class="action-comment">{{ action.comment }}</div>
              </div>
            </Timeline.Item>
          </Timeline>
        </div>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue';
import { message } from 'ant-design-vue';
import type {
  Alert,
  AlertRule,
  AlertLevel,
  AlertStatus,
  MetricType,
  RuleType,
  NotificationChannel,
  MonitoringStatistics,
  AlertRuleCreateRequest,
} from '#/api/types';
import {
  getAlerts,
  getAlertRules,
  getMonitoringStatistics,
  getSystemHealth,
  createAlertRule,
  toggleAlertRule,
  deleteAlertRule,
  acknowledgeAlert,
  silenceAlert,
  resolveAlert,
  getNotificationChannels,
  testAlertRule,
} from '#/api';
import { formatDateTime, formatFileSize } from '#/utils/date';
import * as echarts from 'echarts';

import {
  Card,
  Row,
  Col,
  Button,
  Space,
  Statistic,
  Progress,
  Table,
  Tag,
  Badge,
  Select,
  Empty,
  Modal,
  Form,
  Input,
  InputNumber,
  Switch,
  Descriptions,
  Timeline,
} from 'ant-design-vue';
import {
  DashboardOutlined,
  ReloadOutlined,
  PlusOutlined,
  BarChartOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  CloseCircleOutlined,
  EditOutlined,
  DeleteOutlined,
  PlayCircleOutlined,
  SettingOutlined,
} from '@ant-design/icons-vue';

defineOptions({ name: 'MonitoringDashboard' });

// 响应式数据
const loading = ref(false);
const rulesLoading = ref(false);
const createRuleLoading = ref(false);
const createRuleModalVisible = ref(false);
const alertDetailModalVisible = ref(false);

const alerts = ref<Alert[]>([]);
const alertRules = ref<AlertRule[]>([]);
const notificationChannels = ref<NotificationChannel[]>([]);
const selectedAlert = ref<Alert | null>(null);

const statistics = ref<MonitoringStatistics>({
  totalMetrics: 0,
  totalAlerts: 0,
  activeAlerts: 0,
  resolvedAlerts: 0,
  criticalAlerts: 0,
  warningAlerts: 0,
  totalRules: 0,
  enabledRules: 0,
  totalDashboards: 0,
  notificationChannels: 0,
  systemHealth: {
    overall: 'healthy',
    cpu: 0,
    memory: 0,
    disk: 0,
    network: 0,
  },
  alertsByLevel: {
    info: 0,
    warning: 0,
    error: 0,
    critical: 0,
  },
  alertsByType: {
    cpu: 0,
    memory: 0,
    gpu: 0,
    disk: 0,
    network: 0,
    temperature: 0,
    power: 0,
    custom: 0,
  },
  recentAlerts: [],
  topAlertRules: [],
});

const systemHealth = ref({
  overall: 'healthy' as 'healthy' | 'warning' | 'critical',
  cpu: 25,
  memory: 45,
  disk: 60,
  network: 10,
});

// 筛选器
const alertFilter = reactive({
  level: undefined as AlertLevel | undefined,
  status: undefined as AlertStatus | undefined,
});

// 创建规则表单
const createRuleForm = reactive<AlertRuleCreateRequest>({
  name: '',
  description: '',
  type: 'threshold' as RuleType,
  metricType: 'cpu' as MetricType,
  query: '',
  conditions: [{
    operator: 'gt' as const,
    threshold: 80,
    duration: 300,
  }],
  level: 'warning' as AlertLevel,
  notificationChannels: [],
  evaluationInterval: 60,
});

const createRuleFormRef = ref();
const trendChartRef = ref();

// 模拟数据
const mockAlerts: Alert[] = [
  {
    id: 'alert-001',
    ruleId: 'rule-001',
    ruleName: 'GPU温度过高',
    level: 'critical' as AlertLevel,
    status: 'firing' as AlertStatus,
    summary: 'GPU-0 温度超过85°C',
    description: 'GPU设备温度持续超过安全阈值，可能导致性能下降或硬件损坏',
    metricType: 'temperature' as MetricType,
    currentValue: 87.5,
    threshold: 85,
    labels: { device: 'gpu-0', node: 'node-001' },
    annotations: {},
    startsAt: '2024-01-20 14:30:00',
    duration: 1800,
    fingerprint: 'fp-001',
    notificationsSent: 2,
    relatedResources: [
      { type: 'gpu', id: 'gpu-001', name: 'NVIDIA RTX 4090' },
      { type: 'node', id: 'node-001', name: 'worker-node-01' },
    ],
    actions: [
      {
        type: 'acknowledge',
        userId: 'user-001',
        userName: '张三',
        timestamp: '2024-01-20 14:35:00',
        comment: '正在检查冷却系统',
      },
    ],
  },
  // 添加更多模拟告警...
];

const mockRules: AlertRule[] = [
  {
    id: 'rule-001',
    name: 'GPU温度监控',
    description: '监控GPU温度，超过85度时告警',
    type: 'threshold' as RuleType,
    metricType: 'temperature' as MetricType,
    query: 'nvidia_gpu_temperature_celsius > 85',
    conditions: [{
      operator: 'gt',
      threshold: 85,
      duration: 300,
    }],
    level: 'critical' as AlertLevel,
    enabled: true,
    labels: { component: 'gpu' },
    annotations: { description: 'GPU温度过高告警' },
    notificationChannels: ['channel-001'],
    evaluationInterval: 30,
    creatorId: 'user-001',
    creatorName: '管理员',
    createTime: '2024-01-15 10:00:00',
    updateTime: '2024-01-20 14:00:00',
  },
  // 添加更多模拟规则...
];

// 计算属性
const filteredAlerts = computed(() => {
  let filtered = alerts.value;

  if (alertFilter.level) {
    filtered = filtered.filter(alert => alert.level === alertFilter.level);
  }

  if (alertFilter.status) {
    filtered = filtered.filter(alert => alert.status === alertFilter.status);
  }

  return filtered.slice(0, 10); // 只显示前10个
});

// 表格列定义
const ruleColumns = [
  {
    title: '规则信息',
    key: 'ruleInfo',
    slots: { customRender: 'ruleInfo' },
    width: 200,
  },
  {
    title: '类型',
    key: 'type',
    slots: { customRender: 'type' },
    width: 80,
  },
  {
    title: '监控指标',
    key: 'metric',
    slots: { customRender: 'metric' },
    width: 100,
  },
  {
    title: '告警级别',
    key: 'level',
    slots: { customRender: 'level' },
    width: 80,
  },
  {
    title: '状态',
    key: 'status',
    slots: { customRender: 'status' },
    width: 60,
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
const getHealthText = (status: string) => {
  const texts = {
    healthy: '健康',
    warning: '警告',
    critical: '严重',
  };
  return texts[status as keyof typeof texts] || status;
};

const getProgressColor = (value: number) => {
  if (value >= 90) return '#ff4d4f';
  if (value >= 70) return '#faad14';
  return '#52c41a';
};

const getAlertLevelText = (level: AlertLevel) => {
  const texts = {
    info: '信息',
    warning: '警告',
    error: '错误',
    critical: '严重',
  };
  return texts[level];
};

const getAlertLevelColor = (level: AlertLevel) => {
  const colors = {
    info: 'blue',
    warning: 'orange',
    error: 'red',
    critical: 'volcano',
  };
  return colors[level];
};

const getAlertStatusText = (status: AlertStatus) => {
  const texts = {
    pending: '待处理',
    firing: '触发中',
    resolved: '已解决',
    silenced: '已静默',
  };
  return texts[status];
};

const getAlertStatusBadge = (status: AlertStatus) => {
  const badges = {
    pending: 'warning' as const,
    firing: 'error' as const,
    resolved: 'success' as const,
    silenced: 'default' as const,
  };
  return badges[status];
};

const getMetricTypeText = (type: MetricType) => {
  const texts = {
    cpu: 'CPU',
    memory: '内存',
    gpu: 'GPU',
    disk: '磁盘',
    network: '网络',
    temperature: '温度',
    power: '功耗',
    custom: '自定义',
  };
  return texts[type];
};

const getMetricTypeColor = (type: MetricType) => {
  const colors = {
    cpu: 'blue',
    memory: 'green',
    gpu: 'purple',
    disk: 'orange',
    network: 'cyan',
    temperature: 'red',
    power: 'magenta',
    custom: 'default',
  };
  return colors[type];
};

const getRuleTypeText = (type: RuleType) => {
  const texts = {
    threshold: '阈值',
    anomaly: '异常',
    composite: '复合',
    custom: '自定义',
  };
  return texts[type];
};

const getActionText = (type: string) => {
  const texts = {
    acknowledge: '确认',
    silence: '静默',
    resolve: '解决',
    escalate: '升级',
  };
  return texts[type as keyof typeof texts] || type;
};

const getActionColor = (type: string) => {
  const colors = {
    acknowledge: 'blue',
    silence: 'orange',
    resolve: 'green',
    escalate: 'red',
  };
  return colors[type as keyof typeof colors] || 'default';
};

const formatRelativeTime = (time: string): string => {
  const now = new Date();
  const target = new Date(time);
  const diffMs = now.getTime() - target.getTime();
  const diffMinutes = Math.floor(diffMs / (1000 * 60));

  if (diffMinutes < 60) {
    return `${diffMinutes}分钟前`;
  } else if (diffMinutes < 1440) {
    const diffHours = Math.floor(diffMinutes / 60);
    return `${diffHours}小时前`;
  } else {
    const diffDays = Math.floor(diffMinutes / 1440);
    return `${diffDays}天前`;
  }
};

const formatDuration = (seconds: number): string => {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  if (hours > 0) {
    return `${hours}小时${minutes}分钟`;
  } else if (minutes > 0) {
    return `${minutes}分钟${secs}秒`;
  } else {
    return `${secs}秒`;
  }
};

// 数据加载
const loadAlerts = async () => {
  try {
    loading.value = true;
    // const response = await getAlerts();
    // alerts.value = response.data.items;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    alerts.value = mockAlerts;
  } catch (error) {
    message.error('加载告警列表失败');
  } finally {
    loading.value = false;
  }
};

const loadAlertRules = async () => {
  try {
    rulesLoading.value = true;
    // const response = await getAlertRules();
    // alertRules.value = response.data.items;
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    alertRules.value = mockRules;
  } catch (error) {
    message.error('加载告警规则失败');
  } finally {
    rulesLoading.value = false;
  }
};

const loadStatistics = async () => {
  try {
    // const response = await getMonitoringStatistics();
    // statistics.value = response.data;
    
    // 模拟统计数据
    statistics.value = {
      totalMetrics: 245,
      totalAlerts: 1234,
      activeAlerts: 3,
      resolvedAlerts: 45,
      criticalAlerts: 1,
      warningAlerts: 2,
      totalRules: 15,
      enabledRules: 12,
      totalDashboards: 8,
      notificationChannels: 5,
      systemHealth: {
        overall: 'warning',
        cpu: 25,
        memory: 45,
        disk: 60,
        network: 10,
      },
      alertsByLevel: {
        info: 10,
        warning: 25,
        error: 8,
        critical: 2,
      },
      alertsByType: {
        cpu: 15,
        memory: 12,
        gpu: 8,
        disk: 5,
        network: 3,
        temperature: 2,
        power: 0,
        custom: 0,
      },
      recentAlerts: [],
      topAlertRules: [],
    };
  } catch (error) {
    message.error('加载统计信息失败');
  }
};

const loadNotificationChannels = async () => {
  try {
    // const response = await getNotificationChannels();
    // notificationChannels.value = response.data.items;
    
    // 模拟通知渠道
    notificationChannels.value = [
      { id: 'channel-001', name: '邮件通知', type: 'email' },
      { id: 'channel-002', name: '钉钉群', type: 'dingtalk' },
      { id: 'channel-003', name: '企业微信', type: 'wechat' },
    ] as NotificationChannel[];
  } catch (error) {
    message.error('加载通知渠道失败');
  }
};

const refreshData = async () => {
  await Promise.all([
    loadAlerts(),
    loadAlertRules(),
    loadStatistics(),
    loadNotificationChannels(),
  ]);
};

// 事件处理
const filterAlerts = () => {
  // 筛选逻辑已在computed中实现
};

const viewAlertDetail = (alert: Alert) => {
  selectedAlert.value = alert;
  alertDetailModalVisible.value = true;
};

const acknowledgeAlert = async (alertId: string) => {
  try {
    // await acknowledgeAlert(alertId, '已确认');
    message.success('告警已确认');
    loadAlerts();
  } catch (error) {
    message.error('确认失败');
  }
};

const silenceAlert = async (alertId: string) => {
  try {
    // await silenceAlert(alertId, 3600, '临时静默1小时');
    message.success('告警已静默');
    loadAlerts();
  } catch (error) {
    message.error('静默失败');
  }
};

const resolveAlert = async (alertId: string) => {
  try {
    // await resolveAlert(alertId, '问题已解决');
    message.success('告警已解决');
    loadAlerts();
  } catch (error) {
    message.error('解决失败');
  }
};

const showCreateRuleModal = () => {
  createRuleModalVisible.value = true;
  resetCreateRuleForm();
};

const resetCreateRuleForm = () => {
  Object.assign(createRuleForm, {
    name: '',
    description: '',
    type: 'threshold',
    metricType: 'cpu',
    query: '',
    conditions: [{
      operator: 'gt',
      threshold: 80,
      duration: 300,
    }],
    level: 'warning',
    notificationChannels: [],
    evaluationInterval: 60,
  });
};

const addCondition = () => {
  createRuleForm.conditions.push({
    operator: 'gt',
    threshold: 0,
    duration: 300,
  });
};

const removeCondition = (index: number) => {
  createRuleForm.conditions.splice(index, 1);
};

const handleCreateRule = async () => {
  try {
    await createRuleFormRef.value?.validate();
    
    createRuleLoading.value = true;
    // const response = await createAlertRule(createRuleForm);
    
    // 模拟创建成功
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    message.success('告警规则创建成功');
    createRuleModalVisible.value = false;
    loadAlertRules();
  } catch (error) {
    message.error('创建失败');
  } finally {
    createRuleLoading.value = false;
  }
};

const handleCreateRuleCancel = () => {
  createRuleModalVisible.value = false;
};

const toggleRule = async (ruleId: string, enabled: boolean) => {
  try {
    // await toggleAlertRule(ruleId, enabled);
    message.success(`规则已${enabled ? '启用' : '禁用'}`);
    loadAlertRules();
  } catch (error) {
    message.error('操作失败');
  }
};

const editRule = (rule: AlertRule) => {
  message.info('编辑功能开发中');
};

const testRule = async (rule: AlertRule) => {
  try {
    // await testAlertRule(rule);
    message.success('规则测试通过');
  } catch (error) {
    message.error('规则测试失败');
  }
};

const deleteRule = async (rule: AlertRule) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除规则 "${rule.name}" 吗？`,
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await deleteAlertRule(rule.id);
        message.success('规则删除成功');
        loadAlertRules();
      } catch (error) {
        message.error('删除失败');
      }
    },
  });
};

const showDashboardModal = () => {
  message.info('仪表板功能开发中');
};

const showRulesManagement = () => {
  message.info('规则管理功能开发中');
};

// 初始化图表
const initTrendChart = () => {
  const chartDom = trendChartRef.value;
  if (!chartDom) return;

  const chart = echarts.init(chartDom);
  const option = {
    title: {
      text: '24小时告警趋势',
      textStyle: {
        fontSize: 12,
        color: '#666',
      },
    },
    tooltip: {
      trigger: 'axis',
    },
    xAxis: {
      type: 'category',
      data: Array.from({ length: 24 }, (_, i) => `${i}:00`),
    },
    yAxis: {
      type: 'value',
    },
    series: [
      {
        name: '严重',
        type: 'line',
        data: [1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        itemStyle: { color: '#ff4d4f' },
      },
      {
        name: '警告',
        type: 'line',
        data: [2, 1, 0, 1, 0, 2, 1, 0, 1, 0, 0, 1, 3, 2, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1],
        itemStyle: { color: '#faad14' },
      },
    ],
  };

  chart.setOption(option);
};

// 表单验证规则
const createRuleFormRules = {
  name: [
    { required: true, message: '请输入规则名称', trigger: 'blur' },
  ],
  metricType: [
    { required: true, message: '请选择监控指标', trigger: 'change' },
  ],
  level: [
    { required: true, message: '请选择告警级别', trigger: 'change' },
  ],
  query: [
    { required: true, message: '请输入查询表达式', trigger: 'blur' },
  ],
  evaluationInterval: [
    { required: true, message: '请输入评估间隔', trigger: 'blur' },
  ],
};

// 初始化
onMounted(() => {
  refreshData();
  nextTick(() => {
    initTrendChart();
  });
});
</script>

<style scoped lang="scss">
.monitoring-container {
  padding: 24px;
  min-height: 100vh;
  background: #f5f5f5;
}

// 页面头部
.page-header {
  margin-bottom: 24px;
  
  .header-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 24px;
  }
  
  .title-section {
    flex: 1;
  }
  
  .page-title {
    font-size: 28px;
    font-weight: 600;
    margin: 0 0 8px 0;
    display: flex;
    align-items: center;
    gap: 12px;
    position: relative;
    
    .title-icon {
      font-size: 32px;
      color: #1890ff;
    }
    
    .title-glow {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(45deg, #1890ff20, transparent);
      border-radius: 8px;
      z-index: -1;
    }
  }
  
  .page-description {
    font-size: 16px;
    margin: 0;
    color: #666;
  }
}

// 系统健康概览
.health-overview {
  margin-bottom: 24px;
}

.health-card {
  border-radius: 8px;
  overflow: hidden;
  
  :deep(.ant-card-body) {
    padding: 16px;
  }
}

.health-item {
  display: flex;
  align-items: center;
  gap: 12px;
  
  .health-icon {
    font-size: 32px;
    
    &.health-healthy {
      color: #52c41a;
    }
    
    &.health-warning {
      color: #faad14;
    }
    
    &.health-critical {
      color: #ff4d4f;
    }
  }
  
  .health-progress {
    display: flex;
    align-items: center;
  }
  
  .health-info {
    .health-title {
      font-size: 12px;
      color: #666;
      margin-bottom: 4px;
    }
    
    .health-value {
      font-size: 16px;
      font-weight: 600;
    }
  }
}

// 统计卡片
.stats-section {
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 8px;
  
  :deep(.ant-statistic) {
    .ant-statistic-title {
      font-size: 12px;
      margin-bottom: 4px;
    }
    
    .ant-statistic-content {
      font-size: 20px;
    }
  }
  
  .stat-detail {
    margin-top: 8px;
    font-size: 12px;
    color: #666;
    display: flex;
    gap: 12px;
    
    .critical-count {
      color: #ff4d4f;
    }
    
    .warning-count {
      color: #faad14;
    }
  }
}

// 主要内容
.main-content {
  margin-bottom: 24px;
}

.alerts-card {
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
}

.alerts-list {
  max-height: 600px;
  overflow-y: auto;
}

.alert-item {
  padding: 16px;
  border: 1px solid #f0f0f0;
  border-radius: 8px;
  margin-bottom: 12px;
  transition: all 0.3s ease;
  
  &:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  &.alert-critical {
    border-left: 4px solid #ff4d4f;
  }
  
  &.alert-error {
    border-left: 4px solid #ff7875;
  }
  
  &.alert-warning {
    border-left: 4px solid #faad14;
  }
  
  &.alert-info {
    border-left: 4px solid #1890ff;
  }
  
  .alert-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 8px;
    
    .alert-info {
      display: flex;
      align-items: center;
      gap: 8px;
      flex: 1;
    }
    
    .alert-time {
      font-size: 12px;
      color: #999;
    }
  }
  
  .alert-description {
    font-size: 14px;
    color: #666;
    margin-bottom: 8px;
  }
  
  .alert-details {
    margin-bottom: 12px;
    
    .alert-metric,
    .alert-value,
    .alert-threshold {
      font-size: 12px;
      color: #999;
    }
  }
  
  .alert-actions {
    display: flex;
    justify-content: flex-end;
  }
}

.empty-alerts {
  text-align: center;
  padding: 40px 0;
}

// 趋势图表
.trend-card {
  border-radius: 8px;
  
  .trend-chart {
    height: 300px;
  }
}

// 规则管理
.rules-section {
  margin-bottom: 24px;
}

.rules-card {
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
}

.rule-info {
  .rule-name {
    font-weight: 600;
    margin-bottom: 4px;
  }
  
  .rule-desc {
    font-size: 12px;
    color: #666;
  }
}

// 表单
.conditions-editor {
  .condition-row {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
  }
}

// 告警详情
.alert-detail {
  .action-item {
    .action-header {
      display: flex;
      gap: 12px;
      margin-bottom: 4px;
      
      .action-type {
        font-weight: 600;
      }
      
      .action-user {
        color: #1890ff;
      }
      
      .action-time {
        font-size: 12px;
        color: #999;
      }
    }
    
    .action-comment {
      font-size: 12px;
      color: #666;
      font-style: italic;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .monitoring-container {
    padding: 16px;
  }
  
  .header-content {
    flex-direction: column;
    gap: 16px;
  }
  
  .page-title {
    font-size: 24px;
    
    .title-icon {
      font-size: 28px;
    }
  }
  
  .alert-header {
    flex-direction: column;
    gap: 8px;
  }
  
  .alert-actions {
    justify-content: flex-start;
  }
  
  .condition-row {
    flex-direction: column;
    gap: 8px;
    align-items: stretch !important;
  }
}
</style>