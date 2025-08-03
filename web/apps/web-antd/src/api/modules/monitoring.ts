/**
 * 监控告警相关API
 */
import { defHttp } from '#/api/core/http';
import type { BaseListResponse, BaseResponse } from '#/api/types';

// 监控指标类型
export enum MetricType {
  CPU = 'cpu',
  MEMORY = 'memory',
  GPU = 'gpu',
  DISK = 'disk',
  NETWORK = 'network',
  TEMPERATURE = 'temperature',
  POWER = 'power',
  CUSTOM = 'custom',
}

// 告警级别
export enum AlertLevel {
  INFO = 'info',
  WARNING = 'warning',
  ERROR = 'error',
  CRITICAL = 'critical',
}

// 告警状态
export enum AlertStatus {
  PENDING = 'pending',
  FIRING = 'firing',
  RESOLVED = 'resolved',
  SILENCED = 'silenced',
}

// 监控规则类型
export enum RuleType {
  THRESHOLD = 'threshold',
  ANOMALY = 'anomaly',
  COMPOSITE = 'composite',
  CUSTOM = 'custom',
}

// 监控指标数据点
export interface MetricDataPoint {
  timestamp: number;
  value: number;
  labels?: Record<string, string>;
}

// 监控指标
export interface Metric {
  id: string;
  name: string;
  type: MetricType;
  description?: string;
  unit: string;
  labels: Record<string, string>;
  dataPoints: MetricDataPoint[];
  currentValue: number;
  status: 'normal' | 'warning' | 'error';
  lastUpdateTime: string;
}

// 告警规则
export interface AlertRule {
  id: string;
  name: string;
  description?: string;
  type: RuleType;
  metricType: MetricType;
  query: string;
  conditions: {
    operator: 'gt' | 'lt' | 'eq' | 'ne' | 'gte' | 'lte';
    threshold: number;
    duration: number; // 秒
  }[];
  level: AlertLevel;
  enabled: boolean;
  labels: Record<string, string>;
  annotations: Record<string, string>;
  notificationChannels: string[];
  evaluationInterval: number; // 秒
  creatorId: string;
  creatorName: string;
  createTime: string;
  updateTime: string;
}

// 告警事件
export interface Alert {
  id: string;
  ruleId: string;
  ruleName: string;
  level: AlertLevel;
  status: AlertStatus;
  summary: string;
  description: string;
  metricType: MetricType;
  currentValue: number;
  threshold: number;
  labels: Record<string, string>;
  annotations: Record<string, string>;
  startsAt: string;
  endsAt?: string;
  duration: number; // 秒
  fingerprint: string;
  notificationsSent: number;
  relatedResources: {
    type: 'gpu' | 'node' | 'job' | 'model' | 'dataset';
    id: string;
    name: string;
  }[];
  actions: {
    type: 'acknowledge' | 'silence' | 'resolve' | 'escalate';
    userId: string;
    userName: string;
    timestamp: string;
    comment?: string;
  }[];
}

// 监控仪表板
export interface Dashboard {
  id: string;
  name: string;
  description?: string;
  isDefault: boolean;
  isPublic: boolean;
  panels: DashboardPanel[];
  layout: {
    rows: number;
    cols: number;
  };
  refreshInterval: number; // 秒
  timeRange: {
    from: string;
    to: string;
  };
  variables: {
    name: string;
    type: 'query' | 'custom' | 'constant';
    values: string[];
    selected: string;
  }[];
  creatorId: string;
  creatorName: string;
  createTime: string;
  updateTime: string;
}

// 仪表板面板
export interface DashboardPanel {
  id: string;
  title: string;
  type: 'graph' | 'stat' | 'table' | 'heatmap' | 'gauge' | 'progress';
  position: {
    x: number;
    y: number;
    w: number;
    h: number;
  };
  queries: {
    metric: string;
    query: string;
    legend?: string;
  }[];
  options: {
    unit?: string;
    decimals?: number;
    min?: number;
    max?: number;
    thresholds?: {
      value: number;
      color: string;
    }[];
    colors?: string[];
  };
  alerts?: {
    ruleId: string;
    enabled: boolean;
  }[];
}

// 通知渠道
export interface NotificationChannel {
  id: string;
  name: string;
  type: 'email' | 'webhook' | 'slack' | 'dingtalk' | 'wechat' | 'sms';
  enabled: boolean;
  config: {
    // Email配置
    to?: string[];
    cc?: string[];
    subject?: string;
    
    // Webhook配置
    url?: string;
    method?: 'POST' | 'PUT' | 'PATCH';
    headers?: Record<string, string>;
    body?: string;
    
    // Slack配置
    webhook_url?: string;
    channel?: string;
    username?: string;
    
    // 钉钉配置
    webhook_url_dingtalk?: string;
    secret?: string;
    
    // 微信配置
    corp_id?: string;
    agent_id?: string;
    secret_wechat?: string;
    
    // 短信配置
    provider?: string;
    api_key?: string;
    template_id?: string;
  };
  testSettings?: {
    lastTest: string;
    success: boolean;
    error?: string;
  };
  createTime: string;
  updateTime: string;
}

// 监控统计信息
export interface MonitoringStatistics {
  totalMetrics: number;
  totalAlerts: number;
  activeAlerts: number;
  resolvedAlerts: number;
  criticalAlerts: number;
  warningAlerts: number;
  totalRules: number;
  enabledRules: number;
  totalDashboards: number;
  notificationChannels: number;
  systemHealth: {
    overall: 'healthy' | 'warning' | 'critical';
    cpu: number;
    memory: number;
    disk: number;
    network: number;
  };
  alertsByLevel: Record<AlertLevel, number>;
  alertsByType: Record<MetricType, number>;
  recentAlerts: Alert[];
  topAlertRules: {
    ruleId: string;
    ruleName: string;
    alertCount: number;
  }[];
}

// API查询参数
export interface MetricsQueryParams {
  type?: MetricType;
  labels?: Record<string, string>;
  startTime?: string;
  endTime?: string;
  step?: number; // 秒
  aggregation?: 'avg' | 'sum' | 'min' | 'max' | 'count';
}

export interface AlertsQueryParams {
  level?: AlertLevel;
  status?: AlertStatus;
  ruleId?: string;
  metricType?: MetricType;
  startTime?: string;
  endTime?: string;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface RulesQueryParams {
  type?: RuleType;
  metricType?: MetricType;
  enabled?: boolean;
  level?: AlertLevel;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

// 告警规则创建请求
export interface AlertRuleCreateRequest {
  name: string;
  description?: string;
  type: RuleType;
  metricType: MetricType;
  query: string;
  conditions: AlertRule['conditions'];
  level: AlertLevel;
  labels?: Record<string, string>;
  annotations?: Record<string, string>;
  notificationChannels: string[];
  evaluationInterval: number;
}

// 告警规则更新请求
export interface AlertRuleUpdateRequest {
  name?: string;
  description?: string;
  query?: string;
  conditions?: AlertRule['conditions'];
  level?: AlertLevel;
  enabled?: boolean;
  labels?: Record<string, string>;
  annotations?: Record<string, string>;
  notificationChannels?: string[];
  evaluationInterval?: number;
}

// 仪表板创建请求
export interface DashboardCreateRequest {
  name: string;
  description?: string;
  isPublic: boolean;
  panels: Omit<DashboardPanel, 'id'>[];
  layout: Dashboard['layout'];
  refreshInterval: number;
  variables?: Dashboard['variables'];
}

// 通知渠道创建请求
export interface NotificationChannelCreateRequest {
  name: string;
  type: NotificationChannel['type'];
  config: NotificationChannel['config'];
}

// ==================== 监控指标API ====================

// 获取监控指标列表
export const getMetrics = (params?: MetricsQueryParams) => {
  return defHttp.get<BaseListResponse<Metric>>('/api/v1/monitoring/metrics', { params });
};

// 获取指定指标数据
export const getMetricData = (metricName: string, params?: MetricsQueryParams) => {
  return defHttp.get<BaseResponse<MetricDataPoint[]>>(`/api/v1/monitoring/metrics/${metricName}/data`, { params });
};

// 查询自定义指标
export const queryMetrics = (query: string, params?: {
  startTime?: string;
  endTime?: string;
  step?: number;
}) => {
  return defHttp.post<BaseResponse<{
    metric: Record<string, string>;
    values: [number, string][];
  }[]>>('/api/v1/monitoring/metrics/query', { query, ...params });
};

// ==================== 告警API ====================

// 获取告警列表
export const getAlerts = (params?: AlertsQueryParams) => {
  return defHttp.get<BaseListResponse<Alert>>('/api/v1/monitoring/alerts', { params });
};

// 获取告警详情
export const getAlertDetail = (id: string) => {
  return defHttp.get<BaseResponse<Alert>>(`/api/v1/monitoring/alerts/${id}`);
};

// 确认告警
export const acknowledgeAlert = (id: string, comment?: string) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/monitoring/alerts/${id}/acknowledge`, { comment });
};

// 静默告警
export const silenceAlert = (id: string, duration: number, comment?: string) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/monitoring/alerts/${id}/silence`, { duration, comment });
};

// 解决告警
export const resolveAlert = (id: string, comment?: string) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/monitoring/alerts/${id}/resolve`, { comment });
};

// 批量操作告警
export const batchOperateAlerts = (data: {
  alertIds: string[];
  action: 'acknowledge' | 'silence' | 'resolve';
  duration?: number; // 静默时间（秒）
  comment?: string;
}) => {
  return defHttp.post<BaseResponse<void>>('/api/v1/monitoring/alerts/batch', data);
};

// ==================== 告警规则API ====================

// 获取告警规则列表
export const getAlertRules = (params?: RulesQueryParams) => {
  return defHttp.get<BaseListResponse<AlertRule>>('/api/v1/monitoring/rules', { params });
};

// 获取告警规则详情
export const getAlertRuleDetail = (id: string) => {
  return defHttp.get<BaseResponse<AlertRule>>(`/api/v1/monitoring/rules/${id}`);
};

// 创建告警规则
export const createAlertRule = (data: AlertRuleCreateRequest) => {
  return defHttp.post<BaseResponse<AlertRule>>('/api/v1/monitoring/rules', data);
};

// 更新告警规则
export const updateAlertRule = (id: string, data: AlertRuleUpdateRequest) => {
  return defHttp.put<BaseResponse<AlertRule>>(`/api/v1/monitoring/rules/${id}`, data);
};

// 删除告警规则
export const deleteAlertRule = (id: string) => {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/monitoring/rules/${id}`);
};

// 批量删除告警规则
export const batchDeleteAlertRules = (ids: string[]) => {
  return defHttp.delete<BaseResponse<void>>('/api/v1/monitoring/rules/batch', { data: { ids } });
};

// 启用/禁用告警规则
export const toggleAlertRule = (id: string, enabled: boolean) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/monitoring/rules/${id}/toggle`, { enabled });
};

// 测试告警规则
export const testAlertRule = (data: Omit<AlertRuleCreateRequest, 'notificationChannels'>) => {
  return defHttp.post<BaseResponse<{
    valid: boolean;
    matchCount: number;
    sampleData: MetricDataPoint[];
    errors?: string[];
  }>>('/api/v1/monitoring/rules/test', data);
};

// ==================== 仪表板API ====================

// 获取仪表板列表
export const getDashboards = (params?: {
  isPublic?: boolean;
  creatorId?: string;
  page?: number;
  pageSize?: number;
}) => {
  return defHttp.get<BaseListResponse<Dashboard>>('/api/v1/monitoring/dashboards', { params });
};

// 获取仪表板详情
export const getDashboardDetail = (id: string) => {
  return defHttp.get<BaseResponse<Dashboard>>(`/api/v1/monitoring/dashboards/${id}`);
};

// 创建仪表板
export const createDashboard = (data: DashboardCreateRequest) => {
  return defHttp.post<BaseResponse<Dashboard>>('/api/v1/monitoring/dashboards', data);
};

// 更新仪表板
export const updateDashboard = (id: string, data: Partial<DashboardCreateRequest>) => {
  return defHttp.put<BaseResponse<Dashboard>>(`/api/v1/monitoring/dashboards/${id}`, data);
};

// 删除仪表板
export const deleteDashboard = (id: string) => {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/monitoring/dashboards/${id}`);
};

// 复制仪表板
export const cloneDashboard = (id: string, name: string) => {
  return defHttp.post<BaseResponse<Dashboard>>(`/api/v1/monitoring/dashboards/${id}/clone`, { name });
};

// ==================== 通知渠道API ====================

// 获取通知渠道列表
export const getNotificationChannels = () => {
  return defHttp.get<BaseListResponse<NotificationChannel>>('/api/v1/monitoring/notification-channels');
};

// 获取通知渠道详情
export const getNotificationChannelDetail = (id: string) => {
  return defHttp.get<BaseResponse<NotificationChannel>>(`/api/v1/monitoring/notification-channels/${id}`);
};

// 创建通知渠道
export const createNotificationChannel = (data: NotificationChannelCreateRequest) => {
  return defHttp.post<BaseResponse<NotificationChannel>>('/api/v1/monitoring/notification-channels', data);
};

// 更新通知渠道
export const updateNotificationChannel = (id: string, data: Partial<NotificationChannelCreateRequest>) => {
  return defHttp.put<BaseResponse<NotificationChannel>>(`/api/v1/monitoring/notification-channels/${id}`, data);
};

// 删除通知渠道
export const deleteNotificationChannel = (id: string) => {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/monitoring/notification-channels/${id}`);
};

// 测试通知渠道
export const testNotificationChannel = (id: string, message?: string) => {
  return defHttp.post<BaseResponse<{
    success: boolean;
    error?: string;
    response?: any;
  }>>(`/api/v1/monitoring/notification-channels/${id}/test`, { message: message || '这是一条测试消息' });
};

// ==================== 统计API ====================

// 获取监控统计信息
export const getMonitoringStatistics = () => {
  return defHttp.get<BaseResponse<MonitoringStatistics>>('/api/v1/monitoring/statistics');
};

// 获取系统健康状况
export const getSystemHealth = () => {
  return defHttp.get<BaseResponse<{
    status: 'healthy' | 'warning' | 'critical';
    components: {
      name: string;
      status: 'healthy' | 'warning' | 'critical';
      message?: string;
      lastCheck: string;
    }[];
    uptime: number;
    version: string;
  }>>('/api/v1/monitoring/health');
};

// 获取告警趋势
export const getAlertTrends = (params: {
  startTime: string;
  endTime: string;
  interval: '1h' | '6h' | '1d' | '1w';
}) => {
  return defHttp.get<BaseResponse<{
    timestamps: number[];
    series: {
      level: AlertLevel;
      data: number[];
    }[];
  }>>('/api/v1/monitoring/alerts/trends', { params });
};

// 获取指标概览
export const getMetricsOverview = () => {
  return defHttp.get<BaseResponse<{
    totalMetrics: number;
    metricsByType: Record<MetricType, number>;
    healthyMetrics: number;
    warningMetrics: number;
    errorMetrics: number;
    lastUpdateTime: string;
  }>>('/api/v1/monitoring/metrics/overview');
};

// 导出监控配置
export const exportMonitoringConfig = () => {
  return defHttp.get<Blob>('/api/v1/monitoring/export', {
    responseType: 'blob',
  });
};

// 导入监控配置
export const importMonitoringConfig = (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  
  return defHttp.post<BaseResponse<{
    imported: {
      rules: number;
      dashboards: number;
      channels: number;
    };
    errors?: string[];
  }>>('/api/v1/monitoring/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};