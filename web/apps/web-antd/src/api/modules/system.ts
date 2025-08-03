/**
 * 系统管理相关API
 */
import { defHttp } from '#/api/core/http';
import type { BaseListResponse, BaseResponse } from '#/api/types';

// 用户角色
export enum UserRole {
  ADMIN = 'admin',
  USER = 'user',
  VIEWER = 'viewer',
  OPERATOR = 'operator',
}

// 用户状态
export enum UserStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SUSPENDED = 'suspended',
  DELETED = 'deleted',
}

// 系统配置类型
export enum ConfigType {
  SYSTEM = 'system',
  SECURITY = 'security',
  STORAGE = 'storage',
  NETWORK = 'network',
  MONITORING = 'monitoring',
  NOTIFICATION = 'notification',
  CUSTOM = 'custom',
}

// 日志级别
export enum LogLevel {
  DEBUG = 'debug',
  INFO = 'info',
  WARN = 'warn',
  ERROR = 'error',
  FATAL = 'fatal',
}

// 操作类型
export enum AuditAction {
  CREATE = 'create',
  UPDATE = 'update',
  DELETE = 'delete',
  LOGIN = 'login',
  LOGOUT = 'logout',
  UPLOAD = 'upload',
  DOWNLOAD = 'download',
  DEPLOY = 'deploy',
  STOP = 'stop',
  START = 'start',
  VIEW = 'view',
  EXPORT = 'export',
  IMPORT = 'import',
}

// 用户信息
export interface User {
  id: string;
  username: string;
  email: string;
  fullName: string;
  avatar?: string;
  phone?: string;
  department?: string;
  position?: string;
  role: UserRole;
  status: UserStatus;
  permissions: string[];
  lastLoginTime?: string;
  lastLoginIp?: string;
  loginCount: number;
  storageUsed: number;
  storageQuota: number;
  workspaces: string[];
  preferences: {
    language: string;
    timezone: string;
    theme: 'light' | 'dark' | 'auto';
    notifications: {
      email: boolean;
      system: boolean;
      job: boolean;
      alert: boolean;
    };
  };
  createTime: string;
  updateTime: string;
  createdBy: string;
}

// 角色信息
export interface Role {
  id: string;
  name: string;
  description?: string;
  permissions: Permission[];
  userCount: number;
  isSystem: boolean;
  createTime: string;
  updateTime: string;
}

// 权限信息
export interface Permission {
  id: string;
  name: string;
  resource: string;
  action: string;
  description?: string;
  module: string;
}

// 系统配置项
export interface SystemConfig {
  id: string;
  key: string;
  value: string;
  type: ConfigType;
  description?: string;
  isEncrypted: boolean;
  isRequired: boolean;
  defaultValue?: string;
  validation?: {
    type: 'string' | 'number' | 'boolean' | 'json' | 'email' | 'url';
    min?: number;
    max?: number;
    pattern?: string;
    enum?: string[];
  };
  category: string;
  updateTime: string;
  updatedBy: string;
}

// 审计日志
export interface AuditLog {
  id: string;
  userId: string;
  username: string;
  action: AuditAction;
  resource: string;
  resourceId?: string;
  resourceName?: string;
  details: Record<string, any>;
  ipAddress: string;
  userAgent: string;
  timestamp: string;
  duration?: number; // 毫秒
  status: 'success' | 'failed' | 'partial';
  errorMessage?: string;
  sessionId?: string;
}

// 系统日志
export interface SystemLog {
  id: string;
  level: LogLevel;
  service: string;
  component: string;
  message: string;
  details?: Record<string, any>;
  timestamp: string;
  traceId?: string;
  spanId?: string;
  tags: Record<string, string>;
  stackTrace?: string;
}

// 系统信息
export interface SystemInfo {
  version: string;
  buildTime: string;
  gitCommit: string;
  goVersion: string;
  os: string;
  arch: string;
  uptime: number; // 秒
  startTime: string;
  timezone: string;
  hostname: string;
  pid: number;
  components: {
    name: string;
    version: string;
    status: 'healthy' | 'unhealthy' | 'unknown';
    lastCheck: string;
  }[];
}

// 资源使用统计
export interface ResourceUsage {
  cpu: {
    usage: number; // 百分比
    cores: number;
    frequency: number; // MHz
  };
  memory: {
    total: number; // 字节
    used: number;
    free: number;
    usage: number; // 百分比
  };
  disk: {
    total: number; // 字节
    used: number;
    free: number;
    usage: number; // 百分比
    iops: {
      read: number;
      write: number;
    };
  };
  network: {
    bytesIn: number;
    bytesOut: number;
    packetsIn: number;
    packetsOut: number;
    connections: number;
  };
  gpu?: {
    count: number;
    usage: number; // 百分比
    memory: {
      total: number;
      used: number;
    };
    temperature: number;
    powerUsage: number; // 瓦特
  }[];
}

// 系统统计
export interface SystemStatistics {
  users: {
    total: number;
    active: number;
    online: number;
    newToday: number;
  };
  workspaces: {
    total: number;
    active: number;
  };
  jobs: {
    total: number;
    running: number;
    completed: number;
    failed: number;
    todaySubmitted: number;
  };
  models: {
    total: number;
    public: number;
    private: number;
    totalSize: number;
  };
  datasets: {
    total: number;
    public: number;
    private: number;
    totalSize: number;
  };
  storage: {
    total: number;
    used: number;
    available: number;
    userQuotaUsed: number;
    userQuotaTotal: number;
  };
  alerts: {
    total: number;
    active: number;
    resolved: number;
    critical: number;
  };
}

// 备份信息
export interface Backup {
  id: string;
  name: string;
  type: 'full' | 'incremental' | 'differential';
  size: number;
  status: 'pending' | 'running' | 'completed' | 'failed';
  progress?: number;
  filePath?: string;
  downloadUrl?: string;
  includes: string[];
  excludes?: string[];
  createdBy: string;
  createdByName: string;
  createTime: string;
  completeTime?: string;
  errorMessage?: string;
  retention: number; // 保留天数
}

// API查询参数
export interface UsersQueryParams {
  username?: string;
  email?: string;
  role?: UserRole;
  status?: UserStatus;
  department?: string;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface AuditLogsQueryParams {
  userId?: string;
  action?: AuditAction;
  resource?: string;
  startTime?: string;
  endTime?: string;
  ipAddress?: string;
  status?: 'success' | 'failed' | 'partial';
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface SystemLogsQueryParams {
  level?: LogLevel;
  service?: string;
  component?: string;
  message?: string;
  startTime?: string;
  endTime?: string;
  traceId?: string;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

// 用户创建请求
export interface UserCreateRequest {
  username: string;
  email: string;
  password: string;
  fullName: string;
  phone?: string;
  department?: string;
  position?: string;
  role: UserRole;
  permissions?: string[];
  storageQuota?: number;
  workspaces?: string[];
}

// 用户更新请求
export interface UserUpdateRequest {
  email?: string;
  fullName?: string;
  phone?: string;
  department?: string;
  position?: string;
  role?: UserRole;
  status?: UserStatus;
  permissions?: string[];
  storageQuota?: number;
  workspaces?: string[];
  preferences?: User['preferences'];
}

// 角色创建请求
export interface RoleCreateRequest {
  name: string;
  description?: string;
  permissions: string[];
}

// 备份创建请求
export interface BackupCreateRequest {
  name: string;
  type: Backup['type'];
  includes: string[];
  excludes?: string[];
  retention: number;
  description?: string;
}

// ==================== 用户管理API ====================

// 获取用户列表
export const getUsers = (params?: UsersQueryParams) => {
  return defHttp.get<BaseListResponse<User>>('/api/v1/system/users', { params });
};

// 获取用户详情
export const getUserDetail = (id: string) => {
  return defHttp.get<BaseResponse<User>>(`/api/v1/system/users/${id}`);
};

// 创建用户
export const createUser = (data: UserCreateRequest) => {
  return defHttp.post<BaseResponse<User>>('/api/v1/system/users', data);
};

// 更新用户
export const updateUser = (id: string, data: UserUpdateRequest) => {
  return defHttp.put<BaseResponse<User>>(`/api/v1/system/users/${id}`, data);
};

// 删除用户
export const deleteUser = (id: string) => {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/system/users/${id}`);
};

// 批量删除用户
export const batchDeleteUsers = (ids: string[]) => {
  return defHttp.delete<BaseResponse<void>>('/api/v1/system/users/batch', { data: { ids } });
};

// 重置用户密码
export const resetUserPassword = (id: string, newPassword: string) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/system/users/${id}/reset-password`, { password: newPassword });
};

// 启用/禁用用户
export const toggleUserStatus = (id: string, status: UserStatus) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/system/users/${id}/status`, { status });
};

// 获取用户权限
export const getUserPermissions = (id: string) => {
  return defHttp.get<BaseResponse<string[]>>(`/api/v1/system/users/${id}/permissions`);
};

// 更新用户权限
export const updateUserPermissions = (id: string, permissions: string[]) => {
  return defHttp.put<BaseResponse<void>>(`/api/v1/system/users/${id}/permissions`, { permissions });
};

// ==================== 角色管理API ====================

// 获取角色列表
export const getRoles = () => {
  return defHttp.get<BaseListResponse<Role>>('/api/v1/system/roles');
};

// 获取角色详情
export const getRoleDetail = (id: string) => {
  return defHttp.get<BaseResponse<Role>>(`/api/v1/system/roles/${id}`);
};

// 创建角色
export const createRole = (data: RoleCreateRequest) => {
  return defHttp.post<BaseResponse<Role>>('/api/v1/system/roles', data);
};

// 更新角色
export const updateRole = (id: string, data: Partial<RoleCreateRequest>) => {
  return defHttp.put<BaseResponse<Role>>(`/api/v1/system/roles/${id}`, data);
};

// 删除角色
export const deleteRole = (id: string) => {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/system/roles/${id}`);
};

// 获取所有权限
export const getPermissions = () => {
  return defHttp.get<BaseListResponse<Permission>>('/api/v1/system/permissions');
};

// ==================== 系统配置API ====================

// 获取系统配置
export const getSystemConfigs = (type?: ConfigType) => {
  return defHttp.get<BaseListResponse<SystemConfig>>('/api/v1/system/configs', { params: { type } });
};

// 获取配置项
export const getSystemConfig = (key: string) => {
  return defHttp.get<BaseResponse<SystemConfig>>(`/api/v1/system/configs/${key}`);
};

// 更新配置项
export const updateSystemConfig = (key: string, value: string) => {
  return defHttp.put<BaseResponse<void>>(`/api/v1/system/configs/${key}`, { value });
};

// 批量更新配置
export const batchUpdateSystemConfigs = (configs: { key: string; value: string }[]) => {
  return defHttp.put<BaseResponse<void>>('/api/v1/system/configs/batch', { configs });
};

// 重置配置到默认值
export const resetSystemConfig = (key: string) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/system/configs/${key}/reset`);
};

// ==================== 日志管理API ====================

// 获取审计日志
export const getAuditLogs = (params?: AuditLogsQueryParams) => {
  return defHttp.get<BaseListResponse<AuditLog>>('/api/v1/system/audit-logs', { params });
};

// 获取系统日志
export const getSystemLogs = (params?: SystemLogsQueryParams) => {
  return defHttp.get<BaseListResponse<SystemLog>>('/api/v1/system/logs', { params });
};

// 清理日志
export const cleanupLogs = (data: {
  type: 'audit' | 'system';
  beforeDate: string;
  level?: LogLevel;
}) => {
  return defHttp.post<BaseResponse<{
    deletedCount: number;
  }>>('/api/v1/system/logs/cleanup', data);
};

// 导出日志
export const exportLogs = (data: {
  type: 'audit' | 'system';
  startTime?: string;
  endTime?: string;
  filters?: Record<string, any>;
  format: 'csv' | 'json' | 'txt';
}) => {
  return defHttp.post<Blob>('/api/v1/system/logs/export', data, {
    responseType: 'blob',
  });
};

// ==================== 系统信息API ====================

// 获取系统信息
export const getSystemInfo = () => {
  return defHttp.get<BaseResponse<SystemInfo>>('/api/v1/system/info');
};

// 获取资源使用情况
export const getResourceUsage = () => {
  return defHttp.get<BaseResponse<ResourceUsage>>('/api/v1/system/resources');
};

// 获取系统统计
export const getSystemStatistics = () => {
  return defHttp.get<BaseResponse<SystemStatistics>>('/api/v1/system/statistics');
};

// 获取历史资源使用情况
export const getResourceHistory = (params: {
  startTime: string;
  endTime: string;
  interval: '1m' | '5m' | '15m' | '1h' | '6h' | '1d';
  metrics: string[];
}) => {
  return defHttp.get<BaseResponse<{
    timestamps: number[];
    series: {
      metric: string;
      data: number[];
    }[];
  }>>('/api/v1/system/resources/history', { params });
};

// ==================== 备份恢复API ====================

// 获取备份列表
export const getBackups = () => {
  return defHttp.get<BaseListResponse<Backup>>('/api/v1/system/backups');
};

// 创建备份
export const createBackup = (data: BackupCreateRequest) => {
  return defHttp.post<BaseResponse<Backup>>('/api/v1/system/backups', data);
};

// 下载备份
export const downloadBackup = (id: string) => {
  return defHttp.get<Blob>(`/api/v1/system/backups/${id}/download`, {
    responseType: 'blob',
  });
};

// 删除备份
export const deleteBackup = (id: string) => {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/system/backups/${id}`);
};

// 从备份恢复
export const restoreFromBackup = (id: string, options?: {
  includeUsers?: boolean;
  includeConfigs?: boolean;
  includeData?: boolean;
}) => {
  return defHttp.post<BaseResponse<{
    jobId: string;
  }>>(`/api/v1/system/backups/${id}/restore`, options);
};

// 获取备份进度
export const getBackupProgress = (id: string) => {
  return defHttp.get<BaseResponse<{
    status: Backup['status'];
    progress: number;
    message?: string;
  }>>(`/api/v1/system/backups/${id}/progress`);
};

// ==================== 系统维护API ====================

// 重启系统服务
export const restartService = (serviceName: string) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/system/services/${serviceName}/restart`);
};

// 清理缓存
export const clearCache = (types?: string[]) => {
  return defHttp.post<BaseResponse<{
    clearedTypes: string[];
    totalSize: number;
  }>>('/api/v1/system/cache/clear', { types });
};

// 清理临时文件
export const cleanupTempFiles = (olderThan?: string) => {
  return defHttp.post<BaseResponse<{
    deletedFiles: number;
    freedSpace: number;
  }>>('/api/v1/system/cleanup/temp', { olderThan });
};

// 优化数据库
export const optimizeDatabase = () => {
  return defHttp.post<BaseResponse<{
    tablesOptimized: number;
    spaceSaved: number;
    duration: number;
  }>>('/api/v1/system/database/optimize');
};

// 检查系统更新
export const checkSystemUpdate = () => {
  return defHttp.get<BaseResponse<{
    hasUpdate: boolean;
    currentVersion: string;
    latestVersion?: string;
    releaseNotes?: string;
    downloadUrl?: string;
    size?: number;
  }>>('/api/v1/system/update/check');
};

// 应用系统更新
export const applySystemUpdate = () => {
  return defHttp.post<BaseResponse<{
    jobId: string;
  }>>('/api/v1/system/update/apply');
};

// 获取更新进度
export const getUpdateProgress = () => {
  return defHttp.get<BaseResponse<{
    status: 'downloading' | 'installing' | 'restarting' | 'completed' | 'failed';
    progress: number;
    message?: string;
    error?: string;
  }>>('/api/v1/system/update/progress');
};

// 导出系统配置
export const exportSystemConfig = () => {
  return defHttp.get<Blob>('/api/v1/system/config/export', {
    responseType: 'blob',
  });
};

// 导入系统配置
export const importSystemConfig = (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  
  return defHttp.post<BaseResponse<{
    imported: number;
    skipped: number;
    errors: string[];
  }>>('/api/v1/system/config/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};