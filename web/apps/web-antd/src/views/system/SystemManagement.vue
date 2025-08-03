<template>
  <div class="system-management-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1 class="page-title">
            <SettingOutlined class="title-icon" />
            <span class="title-text">系统管理</span>
            <div class="title-glow"></div>
          </h1>
          <p class="page-description">
            <span class="description-text">管理系统配置、用户权限和系统维护</span>
          </p>
        </div>
        <div class="action-section">
          <Space>
            <Button @click="refreshData" :loading="loading">
              <ReloadOutlined />
              刷新
            </Button>
            <Button type="primary" @click="showBackupModal">
              <CloudUploadOutlined />
              创建备份
            </Button>
            <Dropdown>
              <Button>
                <MoreOutlined />
                更多操作
              </Button>
              <template #overlay>
                <Menu>
                  <Menu.Item @click="showMaintenanceModal">
                    <ToolOutlined />
                    系统维护
                  </Menu.Item>
                  <Menu.Item @click="exportSystemConfig">
                    <ExportOutlined />
                    导出配置
                  </Menu.Item>
                  <Menu.Item @click="checkUpdate">
                    <CloudDownloadOutlined />
                    检查更新
                  </Menu.Item>
                </Menu>
              </template>
            </Dropdown>
          </Space>
        </div>
      </div>
    </div>

    <!-- 系统概览 -->
    <div class="overview-section">
      <Row :gutter="16">
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="overview-card">
            <div class="overview-item">
              <div class="overview-icon">
                <UserOutlined />
              </div>
              <div class="overview-info">
                <div class="overview-title">用户总数</div>
                <div class="overview-value">{{ statistics.users.total }}</div>
                <div class="overview-detail">
                  在线: {{ statistics.users.online }} | 活跃: {{ statistics.users.active }}
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="overview-card">
            <div class="overview-item">
              <div class="overview-icon">
                <DesktopOutlined />
              </div>
              <div class="overview-info">
                <div class="overview-title">系统运行时间</div>
                <div class="overview-value">{{ formatUptime(systemInfo.uptime) }}</div>
                <div class="overview-detail">
                  版本: {{ systemInfo.version }}
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="overview-card">
            <div class="overview-item">
              <div class="overview-icon">
                <DatabaseOutlined />
              </div>
              <div class="overview-info">
                <div class="overview-title">存储使用</div>
                <div class="overview-value">{{ Math.round(statistics.storage.used / statistics.storage.total * 100) }}%</div>
                <div class="overview-detail">
                  {{ formatFileSize(statistics.storage.used) }} / {{ formatFileSize(statistics.storage.total) }}
                </div>
              </div>
            </div>
          </Card>
        </Col>
        <Col :xs="24" :sm="12" :md="6" :lg="6">
          <Card class="overview-card">
            <div class="overview-item">
              <div class="overview-icon">
                <WarningOutlined />
              </div>
              <div class="overview-info">
                <div class="overview-title">活跃告警</div>
                <div class="overview-value" :style="{ color: statistics.alerts.active > 0 ? '#ff4d4f' : '#52c41a' }">
                  {{ statistics.alerts.active }}
                </div>
                <div class="overview-detail">
                  严重: {{ statistics.alerts.critical }}
                </div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 资源使用监控 -->
    <div class="resource-section">
      <Card title="资源使用监控" class="resource-card">
        <Row :gutter="16">
          <Col :xs="24" :sm="12" :md="6" :lg="6">
            <div class="resource-item">
              <div class="resource-header">
                <span class="resource-title">CPU使用率</span>
                <span class="resource-value">{{ resourceUsage.cpu.usage }}%</span>
              </div>
              <Progress
                :percent="resourceUsage.cpu.usage"
                :stroke-color="getProgressColor(resourceUsage.cpu.usage)"
                :show-info="false"
              />
              <div class="resource-detail">
                {{ resourceUsage.cpu.cores }}核 @ {{ resourceUsage.cpu.frequency }}MHz
              </div>
            </div>
          </Col>
          <Col :xs="24" :sm="12" :md="6" :lg="6">
            <div class="resource-item">
              <div class="resource-header">
                <span class="resource-title">内存使用率</span>
                <span class="resource-value">{{ resourceUsage.memory.usage }}%</span>
              </div>
              <Progress
                :percent="resourceUsage.memory.usage"
                :stroke-color="getProgressColor(resourceUsage.memory.usage)"
                :show-info="false"
              />
              <div class="resource-detail">
                {{ formatFileSize(resourceUsage.memory.used) }} / {{ formatFileSize(resourceUsage.memory.total) }}
              </div>
            </div>
          </Col>
          <Col :xs="24" :sm="12" :md="6" :lg="6">
            <div class="resource-item">
              <div class="resource-header">
                <span class="resource-title">磁盘使用率</span>
                <span class="resource-value">{{ resourceUsage.disk.usage }}%</span>
              </div>
              <Progress
                :percent="resourceUsage.disk.usage"
                :stroke-color="getProgressColor(resourceUsage.disk.usage)"
                :show-info="false"
              />
              <div class="resource-detail">
                {{ formatFileSize(resourceUsage.disk.used) }} / {{ formatFileSize(resourceUsage.disk.total) }}
              </div>
            </div>
          </Col>
          <Col :xs="24" :sm="12" :md="6" :lg="6">
            <div class="resource-item">
              <div class="resource-header">
                <span class="resource-title">网络连接</span>
                <span class="resource-value">{{ resourceUsage.network.connections }}</span>
              </div>
              <Progress
                :percent="Math.min(resourceUsage.network.connections / 1000 * 100, 100)"
                :stroke-color="getProgressColor(Math.min(resourceUsage.network.connections / 1000 * 100, 100))"
                :show-info="false"
              />
              <div class="resource-detail">
                入: {{ formatFileSize(resourceUsage.network.bytesIn) }} / 出: {{ formatFileSize(resourceUsage.network.bytesOut) }}
              </div>
            </div>
          </Col>
        </Row>
      </Card>
    </div>

    <!-- 主要功能模块 -->
    <div class="modules-section">
      <Row :gutter="16">
        <!-- 用户管理 -->
        <Col :xs="24" :lg="12">
          <Card title="用户管理" class="module-card">
            <template #extra>
              <Space>
                <Button size="small" @click="showUserModal">
                  <PlusOutlined />
                  添加用户
                </Button>
                <Button size="small" @click="showUsersManagement">
                  <SettingOutlined />
                  管理
                </Button>
              </Space>
            </template>

            <div class="user-stats">
              <Row :gutter="8">
                <Col :span="6">
                  <Statistic title="总用户" :value="statistics.users.total" />
                </Col>
                <Col :span="6">
                  <Statistic title="活跃" :value="statistics.users.active" />
                </Col>
                <Col :span="6">
                  <Statistic title="在线" :value="statistics.users.online" />
                </Col>
                <Col :span="6">
                  <Statistic title="今日新增" :value="statistics.users.newToday" />
                </Col>
              </Row>
            </div>

            <Divider style="margin: 16px 0" />

            <div class="recent-users">
              <h4>最近登录用户</h4>
              <List
                :data-source="recentUsers"
                size="small"
                :pagination="false"
              >
                <template #renderItem="{ item }">
                  <List.Item>
                    <List.Item.Meta>
                      <template #avatar>
                        <Avatar>{{ item.fullName?.[0] }}</Avatar>
                      </template>
                      <template #title>
                        {{ item.fullName }}
                      </template>
                      <template #description>
                        <Space size="small">
                          <Tag :color="getRoleColor(item.role)">
                            {{ getRoleText(item.role) }}
                          </Tag>
                          <span class="login-time">{{ formatRelativeTime(item.lastLoginTime) }}</span>
                        </Space>
                      </template>
                    </List.Item.Meta>
                    <div class="user-status">
                      <Badge
                        :status="getUserStatusBadge(item.status)"
                        :text="getUserStatusText(item.status)"
                      />
                    </div>
                  </List.Item>
                </template>
              </List>
            </div>
          </Card>
        </Col>

        <!-- 系统配置 -->
        <Col :xs="24" :lg="12">
          <Card title="系统配置" class="module-card">
            <template #extra>
              <Space>
                <Button size="small" @click="showConfigModal">
                  <EditOutlined />
                  编辑配置
                </Button>
                <Button size="small" @click="exportConfig">
                  <ExportOutlined />
                  导出
                </Button>
              </Space>
            </template>

            <div class="config-categories">
              <div
                v-for="category in configCategories"
                :key="category.key"
                class="config-category"
                @click="editConfigCategory(category)"
              >
                <div class="category-icon">
                  <component :is="category.icon" />
                </div>
                <div class="category-info">
                  <div class="category-name">{{ category.name }}</div>
                  <div class="category-desc">{{ category.description }}</div>
                  <div class="category-count">{{ category.count }} 项配置</div>
                </div>
                <div class="category-arrow">
                  <RightOutlined />
                </div>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 日志和备份 -->
    <div class="logs-backup-section">
      <Row :gutter="16">
        <!-- 系统日志 -->
        <Col :xs="24" :lg="12">
          <Card title="系统日志" class="logs-card">
            <template #extra>
              <Space>
                <Select
                  v-model:value="logFilter.level"
                  style="width: 80px"
                  size="small"
                  placeholder="级别"
                  allow-clear
                  @change="filterLogs"
                >
                  <Select.Option value="error">错误</Select.Option>
                  <Select.Option value="warn">警告</Select.Option>
                  <Select.Option value="info">信息</Select.Option>
                  <Select.Option value="debug">调试</Select.Option>
                </Select>
                <Button size="small" @click="showLogsManagement">
                  <FileTextOutlined />
                  查看全部
                </Button>
              </Space>
            </template>

            <div class="logs-list">
              <div
                v-for="log in filteredLogs"
                :key="log.id"
                class="log-item"
                :class="`log-${log.level}`"
              >
                <div class="log-header">
                  <Tag :color="getLogLevelColor(log.level)" size="small">
                    {{ getLogLevelText(log.level) }}
                  </Tag>
                  <span class="log-service">{{ log.service }}</span>
                  <span class="log-time">{{ formatRelativeTime(log.timestamp) }}</span>
                </div>
                <div class="log-message">{{ log.message }}</div>
                <div v-if="log.component" class="log-component">{{ log.component }}</div>
              </div>
            </div>
          </Card>
        </Col>

        <!-- 备份管理 -->
        <Col :xs="24" :lg="12">
          <Card title="备份管理" class="backup-card">
            <template #extra>
              <Space>
                <Button size="small" @click="showBackupModal">
                  <PlusOutlined />
                  创建备份
                </Button>
                <Button size="small" @click="showBackupsManagement">
                  <SettingOutlined />
                  管理
                </Button>
              </Space>
            </template>

            <div class="backup-stats">
              <Row :gutter="8">
                <Col :span="8">
                  <Statistic title="总备份" :value="backups.length" />
                </Col>
                <Col :span="8">
                  <Statistic title="今日备份" :value="backups.filter(b => isToday(b.createTime)).length" />
                </Col>
                <Col :span="8">
                  <Statistic
                    title="总大小"
                    :value="formatFileSize(backups.reduce((sum, b) => sum + b.size, 0))"
                  />
                </Col>
              </Row>
            </div>

            <Divider style="margin: 16px 0" />

            <div class="recent-backups">
              <h4>最近备份</h4>
              <List
                :data-source="backups.slice(0, 5)"
                size="small"
                :pagination="false"
              >
                <template #renderItem="{ item }">
                  <List.Item>
                    <List.Item.Meta>
                      <template #title>
                        <Space>
                          {{ item.name }}
                          <Tag :color="getBackupTypeColor(item.type)">
                            {{ getBackupTypeText(item.type) }}
                          </Tag>
                        </Space>
                      </template>
                      <template #description>
                        <Space size="small">
                          <span>{{ formatFileSize(item.size) }}</span>
                          <span>{{ formatRelativeTime(item.createTime) }}</span>
                          <span>{{ item.createdByName }}</span>
                        </Space>
                      </template>
                    </List.Item.Meta>
                    <div class="backup-actions">
                      <Space size="small">
                        <Button type="text" size="small" @click="downloadBackup(item.id)">
                          <DownloadOutlined />
                        </Button>
                        <Dropdown>
                          <Button type="text" size="small">
                            <MoreOutlined />
                          </Button>
                          <template #overlay>
                            <Menu>
                              <Menu.Item @click="restoreBackup(item.id)">
                                <UndoOutlined />
                                恢复
                              </Menu.Item>
                              <Menu.Item @click="deleteBackup(item.id)" danger>
                                <DeleteOutlined />
                                删除
                              </Menu.Item>
                            </Menu>
                          </template>
                        </Dropdown>
                      </Space>
                    </div>
                  </List.Item>
                </template>
              </List>
            </div>
          </Card>
        </Col>
      </Row>
    </div>

    <!-- 创建用户模态框 -->
    <Modal
      v-model:open="userModalVisible"
      title="添加用户"
      width="600px"
      @ok="handleCreateUser"
      @cancel="handleUserCancel"
      :confirm-loading="userLoading"
    >
      <Form
        ref="userFormRef"
        :model="userForm"
        :rules="userFormRules"
        layout="vertical"
      >
        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="用户名" name="username">
              <Input v-model:value="userForm.username" placeholder="输入用户名" />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="邮箱" name="email">
              <Input v-model:value="userForm.email" placeholder="输入邮箱" />
            </Form.Item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="姓名" name="fullName">
              <Input v-model:value="userForm.fullName" placeholder="输入姓名" />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="角色" name="role">
              <Select v-model:value="userForm.role" placeholder="选择角色">
                <Select.Option value="admin">管理员</Select.Option>
                <Select.Option value="user">普通用户</Select.Option>
                <Select.Option value="viewer">查看者</Select.Option>
                <Select.Option value="operator">操作员</Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="部门" name="department">
              <Input v-model:value="userForm.department" placeholder="输入部门" />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="职位" name="position">
              <Input v-model:value="userForm.position" placeholder="输入职位" />
            </Form.Item>
          </Col>
        </Row>

        <Form.Item label="初始密码" name="password">
          <Input.Password v-model:value="userForm.password" placeholder="输入密码" />
        </Form.Item>

        <Form.Item label="存储配额(GB)" name="storageQuota">
          <InputNumber
            v-model:value="userForm.storageQuota"
            :min="1"
            :max="1000"
            style="width: 100%"
          />
        </Form.Item>
      </Form>
    </Modal>

    <!-- 创建备份模态框 -->
    <Modal
      v-model:open="backupModalVisible"
      title="创建备份"
      width="600px"
      @ok="handleCreateBackup"
      @cancel="handleBackupCancel"
      :confirm-loading="backupLoading"
    >
      <Form
        ref="backupFormRef"
        :model="backupForm"
        :rules="backupFormRules"
        layout="vertical"
      >
        <Row :gutter="16">
          <Col :span="12">
            <Form.Item label="备份名称" name="name">
              <Input v-model:value="backupForm.name" placeholder="输入备份名称" />
            </Form.Item>
          </Col>
          <Col :span="12">
            <Form.Item label="备份类型" name="type">
              <Select v-model:value="backupForm.type" placeholder="选择备份类型">
                <Select.Option value="full">完整备份</Select.Option>
                <Select.Option value="incremental">增量备份</Select.Option>
                <Select.Option value="differential">差异备份</Select.Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Form.Item label="包含内容" name="includes">
          <Checkbox.Group v-model:value="backupForm.includes">
            <Row>
              <Col :span="8">
                <Checkbox value="users">用户数据</Checkbox>
              </Col>
              <Col :span="8">
                <Checkbox value="configs">系统配置</Checkbox>
              </Col>
              <Col :span="8">
                <Checkbox value="models">模型文件</Checkbox>
              </Col>
              <Col :span="8">
                <Checkbox value="datasets">数据集</Checkbox>
              </Col>
              <Col :span="8">
                <Checkbox value="logs">系统日志</Checkbox>
              </Col>
              <Col :span="8">
                <Checkbox value="jobs">训练任务</Checkbox>
              </Col>
            </Row>
          </Checkbox.Group>
        </Form.Item>

        <Form.Item label="保留天数" name="retention">
          <InputNumber
            v-model:value="backupForm.retention"
            :min="1"
            :max="365"
            style="width: 100%"
            placeholder="输入保留天数"
          />
        </Form.Item>

        <Form.Item label="描述" name="description">
          <Input.TextArea
            v-model:value="backupForm.description"
            placeholder="输入备份描述"
            :rows="3"
          />
        </Form.Item>
      </Form>
    </Modal>

    <!-- 系统维护模态框 -->
    <Modal
      v-model:open="maintenanceModalVisible"
      title="系统维护"
      width="500px"
      :footer="null"
    >
      <div class="maintenance-actions">
        <div class="maintenance-item" @click="clearCache">
          <div class="action-icon">
            <DeleteOutlined />
          </div>
          <div class="action-info">
            <div class="action-title">清理缓存</div>
            <div class="action-desc">清理系统缓存文件，释放内存空间</div>
          </div>
        </div>

        <div class="maintenance-item" @click="cleanupTempFiles">
          <div class="action-icon">
            <FolderOpenOutlined />
          </div>
          <div class="action-info">
            <div class="action-title">清理临时文件</div>
            <div class="action-desc">删除过期的临时文件，释放磁盘空间</div>
          </div>
        </div>

        <div class="maintenance-item" @click="optimizeDatabase">
          <div class="action-icon">
            <DatabaseOutlined />
          </div>
          <div class="action-info">
            <div class="action-title">优化数据库</div>
            <div class="action-desc">优化数据库表结构，提升查询性能</div>
          </div>
        </div>

        <div class="maintenance-item" @click="restartService">
          <div class="action-icon">
            <ReloadOutlined />
          </div>
          <div class="action-info">
            <div class="action-title">重启服务</div>
            <div class="action-desc">重启系统服务，应用最新配置</div>
          </div>
        </div>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import type {
  User,
  UserRole,
  UserStatus,
  UserCreateRequest,
  SystemInfo,
  ResourceUsage,
  SystemStatistics,
  SystemLog,
  Backup,
  BackupCreateRequest,
  LogLevel,
  ConfigType,
} from '#/api/types';
import {
  getUsers,
  createUser,
  getSystemInfo,
  getResourceUsage,
  getSystemStatistics,
  getSystemLogs,
  getBackups,
  createBackup,
  downloadBackup,
  deleteBackup,
  restoreFromBackup,
  clearCache,
  cleanupTempFiles,
  optimizeDatabase,
  restartService,
  checkSystemUpdate,
  exportSystemConfig,
} from '#/api';
import { formatDateTime, formatFileSize } from '#/utils/date';

import {
  Card,
  Row,
  Col,
  Button,
  Space,
  Statistic,
  Progress,
  List,
  Avatar,
  Tag,
  Badge,
  Divider,
  Select,
  Modal,
  Form,
  Input,
  InputNumber,
  Checkbox,
  Menu,
  Dropdown,
} from 'ant-design-vue';
import {
  SettingOutlined,
  ReloadOutlined,
  CloudUploadOutlined,
  MoreOutlined,
  ToolOutlined,
  ExportOutlined,
  CloudDownloadOutlined,
  UserOutlined,
  DesktopOutlined,
  DatabaseOutlined,
  WarningOutlined,
  PlusOutlined,
  EditOutlined,
  RightOutlined,
  FileTextOutlined,
  DownloadOutlined,
  UndoOutlined,
  DeleteOutlined,
  FolderOpenOutlined,
} from '@ant-design/icons-vue';

defineOptions({ name: 'SystemManagement' });

// 响应式数据
const loading = ref(false);
const userLoading = ref(false);
const backupLoading = ref(false);
const userModalVisible = ref(false);
const backupModalVisible = ref(false);
const maintenanceModalVisible = ref(false);

const recentUsers = ref<User[]>([]);
const systemLogs = ref<SystemLog[]>([]);
const backups = ref<Backup[]>([]);

const systemInfo = ref<SystemInfo>({
  version: '1.0.0',
  buildTime: '2024-01-20 10:00:00',
  gitCommit: 'abc123',
  goVersion: '1.21.0',
  os: 'linux',
  arch: 'amd64',
  uptime: 86400 * 7,
  startTime: '2024-01-13 10:00:00',
  timezone: 'Asia/Shanghai',
  hostname: 'volctrain-master',
  pid: 1234,
  components: [],
});

const resourceUsage = ref<ResourceUsage>({
  cpu: { usage: 35, cores: 8, frequency: 2400 },
  memory: { total: 17179869184, used: 7516192768, free: 9663676416, usage: 44 },
  disk: { total: 1099511627776, used: 549755813888, free: 549755813888, usage: 50, iops: { read: 100, write: 80 } },
  network: { bytesIn: 1073741824, bytesOut: 536870912, packetsIn: 1000000, packetsOut: 800000, connections: 256 },
});

const statistics = ref<SystemStatistics>({
  users: { total: 125, active: 98, online: 15, newToday: 3 },
  workspaces: { total: 24, active: 18 },
  jobs: { total: 1234, running: 8, completed: 1200, failed: 26, todaySubmitted: 15 },
  models: { total: 456, public: 123, private: 333, totalSize: 5368709120 },
  datasets: { total: 789, public: 234, private: 555, totalSize: 10737418240 },
  storage: { total: 1099511627776, used: 549755813888, available: 549755813888, userQuotaUsed: 268435456000, userQuotaTotal: 536870912000 },
  alerts: { total: 45, active: 3, resolved: 42, critical: 1 },
});

// 筛选器
const logFilter = reactive({
  level: undefined as LogLevel | undefined,
});

// 用户表单
const userForm = reactive<UserCreateRequest>({
  username: '',
  email: '',
  password: '',
  fullName: '',
  department: '',
  position: '',
  role: 'user' as UserRole,
  storageQuota: 10,
});

const userFormRef = ref();

// 备份表单
const backupForm = reactive<BackupCreateRequest>({
  name: '',
  type: 'full',
  includes: ['users', 'configs'],
  retention: 30,
  description: '',
});

const backupFormRef = ref();

// 配置分类
const configCategories = ref([
  {
    key: 'system',
    name: '系统设置',
    description: '基础系统配置',
    icon: 'SettingOutlined',
    count: 12,
  },
  {
    key: 'security',
    name: '安全配置',
    description: '认证和权限设置',
    icon: 'SafetyCertificateOutlined',
    count: 8,
  },
  {
    key: 'storage',
    name: '存储配置',
    description: '文件存储和备份设置',
    icon: 'DatabaseOutlined',
    count: 6,
  },
  {
    key: 'notification',
    name: '通知配置',
    description: '邮件和消息通知设置',
    icon: 'BellOutlined',
    count: 4,
  },
]);

// 模拟数据
const mockUsers: User[] = [
  {
    id: 'user-001',
    username: 'admin',
    email: 'admin@example.com',
    fullName: '系统管理员',
    role: 'admin' as UserRole,
    status: 'active' as UserStatus,
    permissions: [],
    lastLoginTime: '2024-01-20 14:30:00',
    loginCount: 156,
    storageUsed: 1073741824,
    storageQuota: 10737418240,
    workspaces: [],
    preferences: {
      language: 'zh-CN',
      timezone: 'Asia/Shanghai',
      theme: 'light',
      notifications: { email: true, system: true, job: true, alert: true },
    },
    createTime: '2024-01-01 00:00:00',
    updateTime: '2024-01-20 14:30:00',
    createdBy: 'system',
  },
  // 添加更多模拟用户...
];

const mockLogs: SystemLog[] = [
  {
    id: 'log-001',
    level: 'info' as LogLevel,
    service: 'api-server',
    component: 'auth',
    message: '用户登录成功',
    timestamp: '2024-01-20 14:30:00',
    tags: { userId: 'user-001', ip: '192.168.1.100' },
  },
  {
    id: 'log-002',
    level: 'warn' as LogLevel,
    service: 'scheduler',
    component: 'job-manager',
    message: 'GPU资源不足，任务排队等待',
    timestamp: '2024-01-20 14:25:00',
    tags: { jobId: 'job-123', gpuRequired: '2' },
  },
  // 添加更多模拟日志...
];

const mockBackups: Backup[] = [
  {
    id: 'backup-001',
    name: 'daily-backup-20240120',
    type: 'full',
    size: 1073741824,
    status: 'completed',
    includes: ['users', 'configs', 'models'],
    createdBy: 'user-001',
    createdByName: '系统管理员',
    createTime: '2024-01-20 02:00:00',
    completeTime: '2024-01-20 02:30:00',
    retention: 30,
  },
  // 添加更多模拟备份...
];

// 计算属性
const filteredLogs = computed(() => {
  let filtered = systemLogs.value;

  if (logFilter.level) {
    filtered = filtered.filter(log => log.level === logFilter.level);
  }

  return filtered.slice(0, 10);
});

// 工具方法
const formatUptime = (seconds: number): string => {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);

  if (days > 0) {
    return `${days}天${hours}小时`;
  } else if (hours > 0) {
    return `${hours}小时${minutes}分钟`;
  } else {
    return `${minutes}分钟`;
  }
};

const getProgressColor = (value: number) => {
  if (value >= 90) return '#ff4d4f';
  if (value >= 70) return '#faad14';
  return '#52c41a';
};

const getRoleText = (role: UserRole) => {
  const texts = {
    admin: '管理员',
    user: '用户',
    viewer: '查看者',
    operator: '操作员',
  };
  return texts[role];
};

const getRoleColor = (role: UserRole) => {
  const colors = {
    admin: 'red',
    user: 'blue',
    viewer: 'green',
    operator: 'orange',
  };
  return colors[role];
};

const getUserStatusText = (status: UserStatus) => {
  const texts = {
    active: '活跃',
    inactive: '非活跃',
    suspended: '已暂停',
    deleted: '已删除',
  };
  return texts[status];
};

const getUserStatusBadge = (status: UserStatus) => {
  const badges = {
    active: 'success' as const,
    inactive: 'default' as const,
    suspended: 'warning' as const,
    deleted: 'error' as const,
  };
  return badges[status];
};

const getLogLevelText = (level: LogLevel) => {
  const texts = {
    debug: '调试',
    info: '信息',
    warn: '警告',
    error: '错误',
    fatal: '致命',
  };
  return texts[level];
};

const getLogLevelColor = (level: LogLevel) => {
  const colors = {
    debug: 'default',
    info: 'blue',
    warn: 'orange',
    error: 'red',
    fatal: 'volcano',
  };
  return colors[level];
};

const getBackupTypeText = (type: string) => {
  const texts = {
    full: '完整',
    incremental: '增量',
    differential: '差异',
  };
  return texts[type as keyof typeof texts] || type;
};

const getBackupTypeColor = (type: string) => {
  const colors = {
    full: 'blue',
    incremental: 'green',
    differential: 'orange',
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

const isToday = (time: string): boolean => {
  const today = new Date();
  const target = new Date(time);
  return today.toDateString() === target.toDateString();
};

// 数据加载
const loadUsers = async () => {
  try {
    // const response = await getUsers({ pageSize: 10 });
    // recentUsers.value = response.data.items;
    
    // 模拟数据
    recentUsers.value = mockUsers;
  } catch (error) {
    message.error('加载用户列表失败');
  }
};

const loadSystemInfo = async () => {
  try {
    // const response = await getSystemInfo();
    // systemInfo.value = response.data;
    
    // 模拟数据已设置
  } catch (error) {
    message.error('加载系统信息失败');
  }
};

const loadResourceUsage = async () => {
  try {
    // const response = await getResourceUsage();
    // resourceUsage.value = response.data;
    
    // 模拟数据已设置
  } catch (error) {
    message.error('加载资源使用情况失败');
  }
};

const loadStatistics = async () => {
  try {
    // const response = await getSystemStatistics();
    // statistics.value = response.data;
    
    // 模拟数据已设置
  } catch (error) {
    message.error('加载统计信息失败');
  }
};

const loadSystemLogs = async () => {
  try {
    // const response = await getSystemLogs({ pageSize: 20 });
    // systemLogs.value = response.data.items;
    
    // 模拟数据
    systemLogs.value = mockLogs;
  } catch (error) {
    message.error('加载系统日志失败');
  }
};

const loadBackups = async () => {
  try {
    // const response = await getBackups();
    // backups.value = response.data.items;
    
    // 模拟数据
    backups.value = mockBackups;
  } catch (error) {
    message.error('加载备份列表失败');
  }
};

const refreshData = async () => {
  loading.value = true;
  try {
    await Promise.all([
      loadUsers(),
      loadSystemInfo(),
      loadResourceUsage(),
      loadStatistics(),
      loadSystemLogs(),
      loadBackups(),
    ]);
  } finally {
    loading.value = false;
  }
};

// 事件处理
const filterLogs = () => {
  // 筛选逻辑已在computed中实现
};

const showUserModal = () => {
  userModalVisible.value = true;
  resetUserForm();
};

const resetUserForm = () => {
  Object.assign(userForm, {
    username: '',
    email: '',
    password: '',
    fullName: '',
    department: '',
    position: '',
    role: 'user',
    storageQuota: 10,
  });
};

const handleCreateUser = async () => {
  try {
    await userFormRef.value?.validate();
    
    userLoading.value = true;
    // const response = await createUser(userForm);
    
    // 模拟创建成功
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    message.success('用户创建成功');
    userModalVisible.value = false;
    loadUsers();
  } catch (error) {
    message.error('创建失败');
  } finally {
    userLoading.value = false;
  }
};

const handleUserCancel = () => {
  userModalVisible.value = false;
};

const showBackupModal = () => {
  backupModalVisible.value = true;
  resetBackupForm();
};

const resetBackupForm = () => {
  Object.assign(backupForm, {
    name: `backup-${new Date().toISOString().slice(0, 10)}`,
    type: 'full',
    includes: ['users', 'configs'],
    retention: 30,
    description: '',
  });
};

const handleCreateBackup = async () => {
  try {
    await backupFormRef.value?.validate();
    
    backupLoading.value = true;
    // const response = await createBackup(backupForm);
    
    // 模拟创建成功
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    message.success('备份创建成功');
    backupModalVisible.value = false;
    loadBackups();
  } catch (error) {
    message.error('创建失败');
  } finally {
    backupLoading.value = false;
  }
};

const handleBackupCancel = () => {
  backupModalVisible.value = false;
};

const downloadBackup = async (backupId: string) => {
  try {
    // const blob = await downloadBackup(backupId);
    // const url = window.URL.createObjectURL(blob);
    // const link = document.createElement('a');
    // link.href = url;
    // link.download = 'backup.tar.gz';
    // link.click();
    // window.URL.revokeObjectURL(url);
    
    message.success('备份下载中...');
  } catch (error) {
    message.error('下载失败');
  }
};

const restoreBackup = async (backupId: string) => {
  Modal.confirm({
    title: '确认恢复',
    content: '确定要从此备份恢复吗？此操作将覆盖当前数据。',
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await restoreFromBackup(backupId);
        message.success('恢复任务已启动');
      } catch (error) {
        message.error('恢复失败');
      }
    },
  });
};

const deleteBackup = async (backupId: string) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除此备份吗？此操作不可恢复。',
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await deleteBackup(backupId);
        message.success('备份删除成功');
        loadBackups();
      } catch (error) {
        message.error('删除失败');
      }
    },
  });
};

const showMaintenanceModal = () => {
  maintenanceModalVisible.value = true;
};

const clearCache = async () => {
  try {
    // await clearCache();
    message.success('缓存清理成功');
    maintenanceModalVisible.value = false;
  } catch (error) {
    message.error('清理失败');
  }
};

const cleanupTempFiles = async () => {
  try {
    // await cleanupTempFiles();
    message.success('临时文件清理成功');
    maintenanceModalVisible.value = false;
  } catch (error) {
    message.error('清理失败');
  }
};

const optimizeDatabase = async () => {
  try {
    // await optimizeDatabase();
    message.success('数据库优化成功');
    maintenanceModalVisible.value = false;
  } catch (error) {
    message.error('优化失败');
  }
};

const restartService = async () => {
  Modal.confirm({
    title: '确认重启',
    content: '确定要重启系统服务吗？这将暂时中断服务。',
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      try {
        // await restartService('all');
        message.success('服务重启成功');
        maintenanceModalVisible.value = false;
      } catch (error) {
        message.error('重启失败');
      }
    },
  });
};

const checkUpdate = async () => {
  try {
    // const response = await checkSystemUpdate();
    // if (response.data.hasUpdate) {
    //   message.info(`发现新版本: ${response.data.latestVersion}`);
    // } else {
      message.info('系统已是最新版本');
    // }
  } catch (error) {
    message.error('检查更新失败');
  }
};

const exportSystemConfig = async () => {
  try {
    // const blob = await exportSystemConfig();
    // const url = window.URL.createObjectURL(blob);
    // const link = document.createElement('a');
    // link.href = url;
    // link.download = 'system-config.json';
    // link.click();
    // window.URL.revokeObjectURL(url);
    
    message.success('配置导出成功');
  } catch (error) {
    message.error('导出失败');
  }
};

const showUsersManagement = () => {
  message.info('用户管理功能开发中');
};

const showConfigModal = () => {
  message.info('配置编辑功能开发中');
};

const exportConfig = () => {
  exportSystemConfig();
};

const editConfigCategory = (category: any) => {
  message.info(`${category.name}配置功能开发中`);
};

const showLogsManagement = () => {
  message.info('日志管理功能开发中');
};

const showBackupsManagement = () => {
  message.info('备份管理功能开发中');
};

// 表单验证规则
const userFormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' },
  ],
  fullName: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少 6 个字符', trigger: 'blur' },
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' },
  ],
};

const backupFormRules = {
  name: [
    { required: true, message: '请输入备份名称', trigger: 'blur' },
  ],
  type: [
    { required: true, message: '请选择备份类型', trigger: 'change' },
  ],
  includes: [
    { required: true, message: '请选择包含内容', trigger: 'change' },
  ],
  retention: [
    { required: true, message: '请输入保留天数', trigger: 'blur' },
  ],
};

// 初始化
onMounted(() => {
  refreshData();
});
</script>

<style scoped lang="scss">
.system-management-container {
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

// 系统概览
.overview-section {
  margin-bottom: 24px;
}

.overview-card {
  border-radius: 8px;
  overflow: hidden;
  
  :deep(.ant-card-body) {
    padding: 20px;
  }
}

.overview-item {
  display: flex;
  align-items: center;
  gap: 16px;
  
  .overview-icon {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    background: linear-gradient(135deg, #1890ff, #722ed1);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 24px;
  }
  
  .overview-info {
    flex: 1;
    
    .overview-title {
      font-size: 14px;
      color: #666;
      margin-bottom: 4px;
    }
    
    .overview-value {
      font-size: 24px;
      font-weight: 600;
      color: #262626;
      margin-bottom: 2px;
    }
    
    .overview-detail {
      font-size: 12px;
      color: #999;
    }
  }
}

// 资源监控
.resource-section {
  margin-bottom: 24px;
}

.resource-card {
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
}

.resource-item {
  .resource-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    
    .resource-title {
      font-size: 14px;
      font-weight: 500;
    }
    
    .resource-value {
      font-size: 18px;
      font-weight: 600;
      color: #1890ff;
    }
  }
  
  .resource-detail {
    font-size: 12px;
    color: #999;
    margin-top: 8px;
  }
}

// 功能模块
.modules-section {
  margin-bottom: 24px;
}

.module-card {
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
}

.user-stats {
  margin-bottom: 16px;
  
  :deep(.ant-statistic) {
    text-align: center;
    
    .ant-statistic-title {
      font-size: 12px;
    }
    
    .ant-statistic-content {
      font-size: 18px;
    }
  }
}

.recent-users {
  h4 {
    margin-bottom: 12px;
    font-size: 14px;
    font-weight: 600;
  }
  
  .login-time {
    font-size: 12px;
    color: #999;
  }
  
  .user-status {
    text-align: right;
  }
}

.config-categories {
  .config-category {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    border: 1px solid #f0f0f0;
    border-radius: 8px;
    margin-bottom: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:hover {
      background: #fafafa;
      border-color: #1890ff;
    }
    
    .category-icon {
      width: 40px;
      height: 40px;
      border-radius: 6px;
      background: #f0f0f0;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #666;
      font-size: 20px;
    }
    
    .category-info {
      flex: 1;
      
      .category-name {
        font-size: 14px;
        font-weight: 600;
        margin-bottom: 4px;
      }
      
      .category-desc {
        font-size: 12px;
        color: #666;
        margin-bottom: 2px;
      }
      
      .category-count {
        font-size: 11px;
        color: #999;
      }
    }
    
    .category-arrow {
      color: #bbb;
    }
  }
}

// 日志和备份
.logs-backup-section {
  margin-bottom: 24px;
}

.logs-card,
.backup-card {
  border-radius: 8px;
  
  :deep(.ant-card-head) {
    background: #fafafa;
  }
}

.logs-list {
  max-height: 400px;
  overflow-y: auto;
  
  .log-item {
    padding: 12px;
    border: 1px solid #f0f0f0;
    border-radius: 6px;
    margin-bottom: 8px;
    
    &.log-error {
      border-left: 3px solid #ff4d4f;
    }
    
    &.log-warn {
      border-left: 3px solid #faad14;
    }
    
    &.log-info {
      border-left: 3px solid #1890ff;
    }
    
    .log-header {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 4px;
      
      .log-service {
        font-size: 12px;
        color: #666;
      }
      
      .log-time {
        font-size: 11px;
        color: #999;
        margin-left: auto;
      }
    }
    
    .log-message {
      font-size: 13px;
      margin-bottom: 2px;
    }
    
    .log-component {
      font-size: 11px;
      color: #999;
    }
  }
}

.backup-stats {
  margin-bottom: 16px;
  
  :deep(.ant-statistic) {
    text-align: center;
    
    .ant-statistic-title {
      font-size: 12px;
    }
    
    .ant-statistic-content {
      font-size: 16px;
    }
  }
}

.recent-backups {
  h4 {
    margin-bottom: 12px;
    font-size: 14px;
    font-weight: 600;
  }
  
  .backup-actions {
    text-align: right;
  }
}

// 维护模态框
.maintenance-actions {
  .maintenance-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    border: 1px solid #f0f0f0;
    border-radius: 8px;
    margin-bottom: 12px;
    cursor: pointer;
    transition: all 0.3s ease;
    
    &:hover {
      background: #fafafa;
      border-color: #1890ff;
    }
    
    .action-icon {
      width: 40px;
      height: 40px;
      border-radius: 6px;
      background: #f0f0f0;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #666;
      font-size: 18px;
    }
    
    .action-info {
      flex: 1;
      
      .action-title {
        font-size: 14px;
        font-weight: 600;
        margin-bottom: 4px;
      }
      
      .action-desc {
        font-size: 12px;
        color: #666;
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .system-management-container {
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
  
  .overview-item {
    gap: 12px;
    
    .overview-icon {
      width: 40px;
      height: 40px;
      font-size: 20px;
    }
    
    .overview-value {
      font-size: 20px;
    }
  }
  
  .resource-header {
    .resource-value {
      font-size: 16px;
    }
  }
}
</style>