/**
 * API类型定义统一导出
 */

// 通用类型
export * from './common';

// 工作空间管理
export * from './workspace';

// GPU资源管理
export * from './gpu';

// 训练任务管理
export * from './training';

// 监控告警
export * from './monitoring';

// 用户和认证相关类型
export interface User {
  id: string;
  username: string;
  realName: string;
  email: string;
  avatar?: string;
  department?: string;
  role: UserRole;
  status: UserStatus;
  lastLoginTime?: string;
  createTime: string;
  updateTime: string;
}

export enum UserRole {
  SUPER_ADMIN = 'super_admin',
  ADMIN = 'admin',
  USER = 'user',
  GUEST = 'guest',
}

export enum UserStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  LOCKED = 'locked',
  PENDING = 'pending',
}

export interface LoginRequest {
  username: string;
  password: string;
  captcha?: string;
  rememberMe?: boolean;
}

export interface LoginResponse {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
  user: User;
}

// 模型管理相关类型
export interface Model {
  id: string;
  name: string;
  description?: string;
  version: string;
  type: ModelType;
  framework: string;
  size: number;              // 模型大小(字节)
  accuracy?: number;         // 准确率
  filePath: string;
  downloadUrl?: string;
  
  // 创建信息
  creatorId: string;
  creatorName: string;
  workspaceId: string;
  workspaceName: string;
  projectId?: string;
  projectName?: string;
  
  // 训练信息
  trainedJobId?: string;     // 训练任务ID
  trainingDataset?: string;  // 训练数据集
  hyperParameters?: { [key: string]: any };
  
  // 状态信息
  status: ModelStatus;
  isPublic: boolean;
  downloadCount: number;
  useCount: number;
  
  // 元数据
  tags: string[];
  labels?: { [key: string]: string };
  createTime: string;
  updateTime: string;
}

export enum ModelType {
  CLASSIFICATION = 'classification',
  REGRESSION = 'regression',
  OBJECT_DETECTION = 'object_detection',
  SEMANTIC_SEGMENTATION = 'semantic_segmentation',
  NLP = 'nlp',
  RECOMMENDATION = 'recommendation',
  GENERATIVE = 'generative',
  CUSTOM = 'custom',
}

export enum ModelStatus {
  TRAINING = 'training',
  AVAILABLE = 'available',
  DEPRECATED = 'deprecated',
  DELETED = 'deleted',
}

// 数据集管理相关类型
export interface Dataset {
  id: string;
  name: string;
  description?: string;
  type: DatasetType;
  format: DatasetFormat;
  size: number;              // 数据集大小(字节)
  sampleCount: number;       // 样本数量
  filePath: string;
  downloadUrl?: string;
  
  // 创建信息
  creatorId: string;
  creatorName: string;
  workspaceId: string;
  workspaceName: string;
  
  // 数据信息
  schema?: any;              // 数据模式
  statistics?: any;          // 统计信息
  previewData?: any[];       // 预览数据
  
  // 标注信息
  isLabeled: boolean;
  labelFormat?: string;
  labelCount?: number;
  
  // 状态信息
  status: DatasetStatus;
  isPublic: boolean;
  downloadCount: number;
  useCount: number;
  
  // 版本信息
  version: string;
  parentDatasetId?: string;  // 父数据集ID
  versionHistory: DatasetVersion[];
  
  // 元数据
  tags: string[];
  labels?: { [key: string]: string };
  createTime: string;
  updateTime: string;
}

export enum DatasetType {
  IMAGE = 'image',
  TEXT = 'text',
  AUDIO = 'audio',
  VIDEO = 'video',
  TABULAR = 'tabular',
  TIME_SERIES = 'time_series',
  GRAPH = 'graph',
  CUSTOM = 'custom',
}

export enum DatasetFormat {
  CSV = 'csv',
  JSON = 'json',
  PARQUET = 'parquet',
  COCO = 'coco',
  YOLO = 'yolo',
  PASCAL_VOC = 'pascal_voc',
  IMAGENET = 'imagenet',
  CUSTOM = 'custom',
}

export enum DatasetStatus {
  UPLOADING = 'uploading',
  PROCESSING = 'processing',
  AVAILABLE = 'available',
  ERROR = 'error',
  DELETED = 'deleted',
}

export interface DatasetVersion {
  version: string;
  description?: string;
  changeLog: string;
  size: number;
  sampleCount: number;
  createTime: string;
  creatorId: string;
  creatorName: string;
}

// 系统管理相关类型
export interface SystemRole {
  id: string;
  name: string;
  description?: string;
  permissions: string[];     // 权限列表
  isSystem: boolean;         // 是否系统角色
  userCount: number;         // 用户数量
  createTime: string;
  updateTime: string;
}

export interface Permission {
  id: string;
  name: string;
  code: string;
  description?: string;
  module: string;            // 所属模块
  type: PermissionType;
  parentId?: string;         // 父权限ID
  children?: Permission[];   // 子权限
  sort: number;              // 排序
}

export enum PermissionType {
  MODULE = 'module',         // 模块
  MENU = 'menu',            // 菜单
  BUTTON = 'button',        // 按钮
  API = 'api',              // API接口
}

export interface ApiEndpoint {
  id: string;
  path: string;
  method: string;
  description?: string;
  module: string;
  requireAuth: boolean;
  permissions: string[];     // 需要的权限
  rateLimit?: number;        // 速率限制(请求/分钟)
  deprecated: boolean;
  version: string;
  responseTime: number;      // 平均响应时间(ms)
  callCount: number;         // 调用次数
  errorRate: number;         // 错误率(%)
  createTime: string;
  updateTime: string;
}