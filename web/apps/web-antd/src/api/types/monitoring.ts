import type { CommonStatus, PageQuery, BaseStatistics } from './common';

/**
 * 监控告警类型定义
 */

// 告警级别枚举
export enum AlertSeverity {
  CRITICAL = 'critical',     // 严重
  HIGH = 'high',            // 高
  MEDIUM = 'medium',        // 中等
  LOW = 'low',              // 低
  INFO = 'info',            // 信息
}

// 告警状态枚举
export enum AlertStatus {
  FIRING = 'firing',        // 触发中
  RESOLVED = 'resolved',    // 已解决
  ACKNOWLEDGED = 'acknowledged', // 已确认
  SUPPRESSED = 'suppressed', // 已抑制
  SILENCED = 'silenced',    // 已静默
}

// 告警规则类型枚举
export enum AlertRuleType {
  THRESHOLD = 'threshold',   // 阈值检测
  TREND = 'trend',          // 趋势分析
  ANOMALY = 'anomaly',      // 异常检测
}

// 通知渠道类型枚举
export enum NotificationChannelType {
  EMAIL = 'email',
  SMS = 'sms',
  WEBHOOK = 'webhook',
  DINGTALK = 'dingtalk',
  SLACK = 'slack',
  WECHAT = 'wechat',
}

// 指标类型枚举
export enum MetricType {
  COUNTER = 'counter',      // 计数器
  GAUGE = 'gauge',          // 仪表盘
  HISTOGRAM = 'histogram',  // 直方图
  SUMMARY = 'summary',      // 摘要
}

// 监控指标定义
export interface Metric {
  name: string;
  type: MetricType;
  value: number;
  timestamp: string;
  labels?: { [key: string]: string };
  unit?: string;            // 单位
  description?: string;     // 描述
}

// 系统监控指标
export interface SystemMetrics {
  // CPU指标
  cpuUsagePercent: number;
  cpuCoresUsed: number;
  cpuLoadAvg1m: number;
  cpuLoadAvg5m: number;
  cpuLoadAvg15m: number;
  
  // 内存指标
  memoryUsagePercent: number;
  memoryUsedGB: number;
  memoryTotalGB: number;
  memoryAvailableGB: number;
  
  // 磁盘指标
  diskUsagePercent: number;
  diskUsedGB: number;
  diskTotalGB: number;
  diskReadBytesPerSec: number;
  diskWriteBytesPerSec: number;
  diskIOPS: number;
  
  // 网络指标
  networkRxBytesPerSec: number;
  networkTxBytesPerSec: number;
  networkRxPacketsPerSec: number;
  networkTxPacketsPerSec: number;
  
  // GPU指标(如果有)
  gpuCount?: number;
  gpuUsagePercent?: number;
  gpuMemoryUsagePercent?: number;
  gpuTemperature?: number;
  gpuPowerUsage?: number;
  
  timestamp: string;
}

// 业务监控指标
export interface BusinessMetrics {
  // 训练任务指标
  totalJobs: number;
  runningJobs: number;
  queuedJobs: number;
  completedJobs: number;
  failedJobs: number;
  
  // 用户活跃度
  activeUsers: number;
  newUsers: number;
  
  // 资源利用率
  gpuUtilization: number;
  cpuUtilization: number;
  memoryUtilization: number;
  
  // 队列指标
  avgQueueTime: number;     // 平均排队时间(分钟)
  maxQueueTime: number;     // 最大排队时间(分钟)
  queueLength: number;      // 队列长度
  
  // 集群指标
  healthyClusters: number;
  totalClusters: number;
  healthyNodes: number;
  totalNodes: number;
  
  timestamp: string;
}

// 告警规则条件
export interface AlertCondition {
  metric: string;           // 监控指标名称
  operator: '>' | '<' | '>=' | '<=' | '==' | '!='; // 比较操作符
  threshold: number;        // 阈值
  duration: number;         // 持续时间(秒)
  labels?: { [key: string]: string }; // 标签过滤器
}

// 告警规则
export interface AlertRule {
  id: string;
  name: string;
  description?: string;
  type: AlertRuleType;
  severity: AlertSeverity;
  
  // 规则配置
  conditions: AlertCondition[];
  expression?: string;       // 自定义表达式(高级用法)
  
  // 触发配置
  evaluationInterval: number; // 评估间隔(秒)
  maxEvaluations?: number;   // 最大评估次数
  
  // 抑制配置
  suppressDuration?: number; // 抑制时长(秒)
  groupBy?: string[];       // 分组字段
  
  // 通知配置
  notificationChannels: string[]; // 通知渠道ID列表
  notificationTemplate?: string;   // 通知模板
  
  // 状态信息
  enabled: boolean;
  status: CommonStatus;
  lastTriggered?: string;   // 最后触发时间
  triggerCount: number;     // 触发次数
  
  // 创建信息
  creatorId: string;
  creatorName: string;
  
  createTime: string;
  updateTime: string;
}

// 告警实例
export interface Alert {
  id: string;
  ruleId: string;
  ruleName: string;
  
  // 告警信息
  severity: AlertSeverity;
  status: AlertStatus;
  message: string;
  summary?: string;
  
  // 触发信息
  triggerTime: string;
  resolveTime?: string;
  acknowledgeTime?: string;
  acknowledgedBy?: string;
  
  // 相关信息
  labels: { [key: string]: string };
  annotations: { [key: string]: string };
  source: string;           // 告警源(节点、服务等)
  
  // 处理信息
  assignedTo?: string;      // 分配给谁处理
  comments: AlertComment[]; // 处理评论
  
  // 抑制信息
  suppressedBy?: string[];  // 被哪些规则抑制
  silencedBy?: string[];    // 被哪些静默规则静默
  
  createTime: string;
  updateTime: string;
}

// 告警评论
export interface AlertComment {
  id: string;
  alertId: string;
  userId: string;
  userName: string;
  content: string;
  createTime: string;
}

// 通知渠道配置
export interface NotificationChannel {
  id: string;
  name: string;
  type: NotificationChannelType;
  description?: string;
  
  // 通用配置
  enabled: boolean;
  
  // 邮件配置
  emailConfig?: {
    smtpHost: string;
    smtpPort: number;
    username: string;
    password: string;
    fromAddress: string;
    toAddresses: string[];
    subject?: string;
    template?: string;
  };
  
  // 短信配置
  smsConfig?: {
    provider: string;       // 服务提供商
    apiKey: string;
    secretKey: string;
    phoneNumbers: string[];
    template?: string;
  };
  
  // Webhook配置
  webhookConfig?: {
    url: string;
    method: 'POST' | 'GET' | 'PUT';
    headers?: { [key: string]: string };
    timeout: number;
    retryCount: number;
  };
  
  // 钉钉配置
  dingtalkConfig?: {
    accessToken: string;
    secret?: string;
    atMobiles?: string[];
    isAtAll?: boolean;
  };
  
  // 测试配置
  lastTestTime?: string;
  lastTestResult?: boolean;
  
  createTime: string;
  updateTime: string;
}

// 静默规则
export interface SilenceRule {
  id: string;
  name: string;
  description?: string;
  
  // 匹配条件
  matchers: {
    name: string;           // 标签名
    value: string;          // 标签值
    isRegex: boolean;       // 是否正则表达式
  }[];
  
  // 时间配置
  startTime: string;
  endTime: string;
  
  // 创建信息
  creatorId: string;
  creatorName: string;
  
  // 状态
  active: boolean;
  
  createTime: string;
  updateTime: string;
}

// 监控面板配置
export interface Dashboard {
  id: string;
  name: string;
  description?: string;
  category?: string;
  
  // 面板配置
  panels: DashboardPanel[];
  layout: {
    rows: number;
    cols: number;
  };
  
  // 时间配置
  timeRange: {
    from: string;
    to: string;
    refreshInterval?: number; // 刷新间隔(秒)
  };
  
  // 变量配置
  variables?: {
    name: string;
    label: string;
    query: string;
    multi: boolean;         // 是否多选
    options: Array<{ text: string; value: string }>;
  }[];
  
  // 权限配置
  isPublic: boolean;
  allowedUsers?: string[];
  allowedRoles?: string[];
  
  // 创建信息
  creatorId: string;
  creatorName: string;
  
  createTime: string;
  updateTime: string;
}

// 面板配置
export interface DashboardPanel {
  id: string;
  title: string;
  type: 'graph' | 'stat' | 'table' | 'heatmap' | 'gauge' | 'logs';
  
  // 位置配置
  gridPos: {
    x: number;
    y: number;
    w: number;             // 宽度
    h: number;             // 高度
  };
  
  // 数据查询
  targets: {
    expr: string;          // 查询表达式
    legendFormat?: string; // 图例格式
    interval?: string;     // 查询间隔
  }[];
  
  // 显示配置
  options: {
    displayMode?: string;
    colorMode?: string;
    unit?: string;
    decimals?: number;
    thresholds?: Array<{
      color: string;
      value: number;
    }>;
  };
  
  // 告警配置
  alert?: {
    ruleId?: string;
    conditions: AlertCondition[];
    frequency: string;
  };
}

// 查询参数
export interface AlertQuery extends PageQuery {
  ruleId?: string;
  severity?: AlertSeverity;
  status?: AlertStatus;
  source?: string;
  assignedTo?: string;
  triggerTimeStart?: string;
  triggerTimeEnd?: string;
}

export interface AlertRuleQuery extends PageQuery {
  type?: AlertRuleType;
  severity?: AlertSeverity;
  enabled?: boolean;
  creatorId?: string;
}

export interface NotificationChannelQuery extends PageQuery {
  type?: NotificationChannelType;
  enabled?: boolean;
}

// 统计信息
export interface AlertStatistics extends BaseStatistics {
  bySeverity: {
    [key in AlertSeverity]: number;
  };
  byStatus: {
    [key in AlertStatus]: number;
  };
  avgResolutionTime: number;  // 平均解决时间(分钟)
  mttr: number;              // 平均修复时间(分钟)
  mtbf: number;              // 平均故障间隔时间(小时)
}

// 系统概览
export interface SystemOverview {
  // 集群状态
  totalClusters: number;
  healthyClusters: number;
  totalNodes: number;
  healthyNodes: number;
  totalGPUs: number;
  availableGPUs: number;
  
  // 任务状态
  totalJobs: number;
  runningJobs: number;
  queuedJobs: number;
  completedJobsToday: number;
  failedJobsToday: number;
  
  // 用户状态
  totalUsers: number;
  activeUsersToday: number;
  totalWorkspaces: number;
  
  // 告警状态
  criticalAlerts: number;
  highAlerts: number;
  totalActiveAlerts: number;
  
  // 资源使用情况
  avgCpuUsage: number;
  avgMemoryUsage: number;
  avgGpuUsage: number;
  
  timestamp: string;
}