/**
 * 通用类型定义
 */

// 基础响应结构
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
  timestamp?: string;
}

// 分页请求参数
export interface PageQuery {
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  keyword?: string;
}

// 分页响应结构
export interface PageResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

// 通用状态枚举
export enum CommonStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  DELETED = 'deleted',
  PENDING = 'pending',
}

// 资源配额定义
export interface ResourceQuota {
  cpu: number;        // CPU核心数
  memory: number;     // 内存 (GB)
  gpu?: number;       // GPU数量
  storage?: number;   // 存储空间 (GB)
  bandwidth?: number; // 网络带宽 (Mbps)
}

// 时间范围查询
export interface DateRange {
  startTime?: string;
  endTime?: string;
}

// 标签和注解
export interface Labels {
  [key: string]: string;
}

export interface Annotations {
  [key: string]: string;
}

// 通用的CRUD操作响应
export interface CrudResponse {
  success: boolean;
  message: string;
  id?: string;
}

// 统计信息基础结构
export interface BaseStatistics {
  total: number;
  active: number;
  inactive: number;
  lastUpdated: string;
}

// 错误响应结构
export interface ErrorResponse {
  code: number;
  message: string;
  details?: string;
  field?: string;
}

// 基础响应结构
export interface BaseResponse<T = any> {
  code: number;
  message: string;
  data: T;
  success: boolean;
  timestamp?: string;
}

// 基础列表响应结构
export interface BaseListResponse<T = any> {
  code: number;
  message: string;
  data: {
    items: T[];
    total: number;
    page: number;
    pageSize: number;
    totalPages: number;
    hasNext: boolean;
    hasPrev: boolean;
  };
  success: boolean;
  timestamp?: string;
}