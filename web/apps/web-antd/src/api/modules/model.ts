/**
 * 模型管理相关API
 */
import { defHttp } from '#/api/core/http';
import type { 
  Model, 
  ModelType, 
  ModelStatus,
  BaseListResponse,
  BaseResponse 
} from '#/api/types';

// 模型列表查询参数
export interface ModelListParams {
  name?: string;
  type?: ModelType;
  status?: ModelStatus;
  framework?: string;
  creatorId?: string;
  workspaceId?: string;
  tags?: string[];
  isPublic?: boolean;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

// 模型创建请求
export interface ModelCreateRequest {
  name: string;
  description?: string;
  type: ModelType;
  framework: string;
  workspaceId: string;
  projectId?: string;
  trainedJobId?: string;
  trainingDataset?: string;
  hyperParameters?: { [key: string]: any };
  tags: string[];
  labels?: { [key: string]: string };
  isPublic: boolean;
}

// 模型更新请求
export interface ModelUpdateRequest {
  name?: string;
  description?: string;
  type?: ModelType;
  framework?: string;
  tags?: string[];
  labels?: { [key: string]: string };
  isPublic?: boolean;
  status?: ModelStatus;
}

// 模型上传请求
export interface ModelUploadRequest {
  name: string;
  description?: string;
  type: ModelType;
  framework: string;
  workspaceId: string;
  projectId?: string;
  tags: string[];
  isPublic: boolean;
  file: File;
}

// 模型版本创建请求
export interface ModelVersionRequest {
  version: string;
  description?: string;
  changeLog: string;
  file: File;
  accuracy?: number;
  hyperParameters?: { [key: string]: any };
}

// 模型部署请求
export interface ModelDeployRequest {
  modelId: string;
  version: string;
  targetEnvironment: 'development' | 'staging' | 'production';
  replicas: number;
  resources: {
    cpu: number;
    memory: number;
    gpu?: number;
  };
  configuration?: { [key: string]: any };
}

// 模型统计信息
export interface ModelStatistics {
  totalModels: number;
  publicModels: number;
  privateModels: number;
  modelsByType: Record<ModelType, number>;
  modelsByFramework: Record<string, number>;
  topDownloaded: Model[];
  recentUploaded: Model[];
  storageUsage: {
    total: number;
    used: number;
    available: number;
  };
}

// 获取模型列表
export const getModelList = (params?: ModelListParams) => {
  return defHttp.get<BaseListResponse<Model>>('/api/v1/models', { params });
};

// 获取模型详情
export const getModelDetail = (id: string) => {
  return defHttp.get<BaseResponse<Model>>(`/api/v1/models/${id}`);
};

// 创建模型
export const createModel = (data: ModelCreateRequest) => {
  return defHttp.post<BaseResponse<Model>>('/api/v1/models', data);
};

// 上传模型文件
export const uploadModel = (data: ModelUploadRequest) => {
  const formData = new FormData();
  Object.entries(data).forEach(([key, value]) => {
    if (key === 'file') {
      formData.append(key, value);
    } else if (key === 'tags' && Array.isArray(value)) {
      formData.append(key, JSON.stringify(value));
    } else if (typeof value === 'object') {
      formData.append(key, JSON.stringify(value));
    } else {
      formData.append(key, String(value));
    }
  });

  return defHttp.post<BaseResponse<Model>>('/api/v1/models/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 更新模型
export const updateModel = (id: string, data: ModelUpdateRequest) => {
  return defHttp.put<BaseResponse<Model>>(`/api/v1/models/${id}`, data);
};

// 删除模型
export const deleteModel = (id: string) => {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/models/${id}`);
};

// 批量删除模型
export const batchDeleteModels = (ids: string[]) => {
  return defHttp.delete<BaseResponse<void>>('/api/v1/models/batch', { 
    data: { ids } 
  });
};

// 创建模型版本
export const createModelVersion = (modelId: string, data: ModelVersionRequest) => {
  const formData = new FormData();
  Object.entries(data).forEach(([key, value]) => {
    if (key === 'file') {
      formData.append(key, value);
    } else if (typeof value === 'object') {
      formData.append(key, JSON.stringify(value));
    } else {
      formData.append(key, String(value));
    }
  });

  return defHttp.post<BaseResponse<Model>>(`/api/v1/models/${modelId}/versions`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 获取模型版本历史
export const getModelVersions = (modelId: string) => {
  return defHttp.get<BaseListResponse<any>>(`/api/v1/models/${modelId}/versions`);
};

// 下载模型
export const downloadModel = (id: string, version?: string) => {
  const params = version ? { version } : undefined;
  return defHttp.get<Blob>(`/api/v1/models/${id}/download`, {
    params,
    responseType: 'blob',
  });
};

// 部署模型
export const deployModel = (data: ModelDeployRequest) => {
  return defHttp.post<BaseResponse<any>>('/api/v1/models/deploy', data);
};

// 获取模型部署状态
export const getModelDeployments = (modelId: string) => {
  return defHttp.get<BaseListResponse<any>>(`/api/v1/models/${modelId}/deployments`);
};

// 停止模型部署
export const stopModelDeployment = (deploymentId: string) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/deployments/${deploymentId}/stop`);
};

// 获取模型统计信息
export const getModelStatistics = (workspaceId?: string) => {
  const params = workspaceId ? { workspaceId } : undefined;
  return defHttp.get<BaseResponse<ModelStatistics>>('/api/v1/models/statistics', { params });
};

// 获取模型使用记录
export const getModelUsageHistory = (modelId: string, params?: {
  startTime?: string;
  endTime?: string;
  page?: number;
  pageSize?: number;
}) => {
  return defHttp.get<BaseListResponse<any>>(`/api/v1/models/${modelId}/usage`, { params });
};

// 收藏/取消收藏模型
export const toggleModelFavorite = (modelId: string, favorite: boolean) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/models/${modelId}/favorite`, { favorite });
};

// 获取推荐模型
export const getRecommendedModels = (params?: {
  type?: ModelType;
  framework?: string;
  limit?: number;
}) => {
  return defHttp.get<BaseListResponse<Model>>('/api/v1/models/recommended', { params });
};

// 搜索公开模型
export const searchPublicModels = (query: string, params?: {
  type?: ModelType;
  framework?: string;
  page?: number;
  pageSize?: number;
}) => {
  return defHttp.get<BaseListResponse<Model>>('/api/v1/models/search', {
    params: { query, ...params },
  });
};

// 获取支持的框架列表
export const getSupportedFrameworks = () => {
  return defHttp.get<BaseResponse<string[]>>('/api/v1/models/frameworks');
};

// 验证模型文件
export const validateModelFile = (file: File) => {
  const formData = new FormData();
  formData.append('file', file);
  
  return defHttp.post<BaseResponse<{
    valid: boolean;
    format: string;
    size: number;
    metadata?: any;
    errors?: string[];
  }>>('/api/v1/models/validate', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};