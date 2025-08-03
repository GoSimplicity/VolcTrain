import type { CommonStatus, ResourceQuota, Labels, Annotations, PageQuery, BaseStatistics } from './common';

/**
 * 工作空间管理类型定义
 */

// 工作空间类型枚举
export enum WorkspaceType {
  PERSONAL = 'personal',      // 个人工作空间
  TEAM = 'team',             // 团队工作空间
  PROJECT = 'project',       // 项目工作空间
  DEPARTMENT = 'department', // 部门工作空间
}

// 工作空间权限枚举
export enum WorkspacePermission {
  OWNER = 'owner',       // 所有者
  ADMIN = 'admin',       // 管理员
  MEMBER = 'member',     // 成员
  VIEWER = 'viewer',     // 查看者
}

// 工作空间成员信息
export interface WorkspaceMember {
  id: string;
  userId: string;
  username: string;
  realName: string;
  email: string;
  avatar?: string;
  permission: WorkspacePermission;
  joinTime: string;
  lastActiveTime?: string;
  status: CommonStatus;
}

// 工作空间配置
export interface WorkspaceConfig {
  allowExternalAccess: boolean;  // 允许外部访问
  enableAutoCleanup: boolean;    // 启用自动清理
  cleanupDays: number;           // 清理天数
  maxMembers: number;            // 最大成员数
  defaultResourceQuota: ResourceQuota; // 默认资源配额
}

// 工作空间信息
export interface Workspace {
  id: string;
  name: string;
  description?: string;
  type: WorkspaceType;
  ownerId: string;
  ownerName: string;
  memberCount: number;
  projectCount: number;
  status: CommonStatus;
  resourceQuota: ResourceQuota;
  resourceUsed: ResourceQuota;
  config: WorkspaceConfig;
  labels?: Labels;
  annotations?: Annotations;
  createTime: string;
  updateTime: string;
  lastAccessTime?: string;
}

// 创建工作空间请求
export interface CreateWorkspaceRequest {
  name: string;
  description?: string;
  type: WorkspaceType;
  resourceQuota?: ResourceQuota;
  config?: Partial<WorkspaceConfig>;
  labels?: Labels;
  memberIds?: string[]; // 初始成员ID列表
}

// 更新工作空间请求
export interface UpdateWorkspaceRequest {
  id: string;
  name?: string;
  description?: string;
  resourceQuota?: ResourceQuota;
  config?: Partial<WorkspaceConfig>;
  labels?: Labels;
}

// 工作空间查询参数
export interface WorkspaceQuery extends PageQuery {
  type?: WorkspaceType;
  status?: CommonStatus;
  ownerId?: string;
  memberUserId?: string; // 查询用户所属的工作空间
  createTimeStart?: string;
  createTimeEnd?: string;
}

// 工作空间统计信息
export interface WorkspaceStatistics extends BaseStatistics {
  byType: {
    [key in WorkspaceType]: number;
  };
  totalMembers: number;
  totalProjects: number;
  resourceUtilization: {
    cpu: number;    // 使用率 %
    memory: number;
    gpu: number;
    storage: number;
  };
}

// 邀请成员请求
export interface InviteMemberRequest {
  workspaceId: string;
  userIds: string[];
  permission: WorkspacePermission;
  message?: string; // 邀请消息
}

// 更新成员权限请求
export interface UpdateMemberPermissionRequest {
  workspaceId: string;
  userId: string;
  permission: WorkspacePermission;
}

// 工作空间项目信息
export interface WorkspaceProject {
  id: string;
  name: string;
  description?: string;
  workspaceId: string;
  ownerId: string;
  ownerName: string;
  status: CommonStatus;
  type: string; // 项目类型：训练、推理、数据处理等
  resourceUsed: ResourceQuota;
  trainingJobCount: number;
  datasetCount: number;
  modelCount: number;
  createTime: string;
  updateTime: string;
  lastActiveTime?: string;
  labels?: Labels;
}

// 创建项目请求
export interface CreateProjectRequest {
  name: string;
  description?: string;
  workspaceId: string;
  type: string;
  labels?: Labels;
}

// 项目查询参数
export interface ProjectQuery extends PageQuery {
  workspaceId?: string;
  type?: string;
  status?: CommonStatus;
  ownerId?: string;
  createTimeStart?: string;
  createTimeEnd?: string;
}