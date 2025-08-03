/**
 * 数据集管理相关API
 */
import { defHttp } from '#/api/core/http';
import type { 
  Dataset, 
  DatasetType, 
  DatasetFormat,
  DatasetStatus,
  DatasetVersion,
  BaseListResponse,
  BaseResponse 
} from '#/api/types';

// 数据集列表查询参数
export interface DatasetListParams {
  name?: string;
  type?: DatasetType;
  format?: DatasetFormat;
  status?: DatasetStatus;
  creatorId?: string;
  workspaceId?: string;
  tags?: string[];
  isPublic?: boolean;
  isLabeled?: boolean;
  minSize?: number;
  maxSize?: number;
  minSamples?: number;
  maxSamples?: number;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

// 数据集创建请求
export interface DatasetCreateRequest {
  name: string;
  description?: string;
  type: DatasetType;
  format: DatasetFormat;
  workspaceId: string;
  tags: string[];
  labels?: { [key: string]: string };
  isPublic: boolean;
  schema?: any;
  labelFormat?: string;
}

// 数据集更新请求
export interface DatasetUpdateRequest {
  name?: string;
  description?: string;
  type?: DatasetType;
  format?: DatasetFormat;
  tags?: string[];
  labels?: { [key: string]: string };
  isPublic?: boolean;
  status?: DatasetStatus;
  schema?: any;
  labelFormat?: string;
}

// 数据集上传请求
export interface DatasetUploadRequest {
  name: string;
  description?: string;
  type: DatasetType;
  format: DatasetFormat;
  workspaceId: string;
  tags: string[];
  isPublic: boolean;
  files: File[];
  labelFiles?: File[];
  schema?: any;
}

// 数据集预处理请求
export interface DatasetPreprocessRequest {
  operations: Array<{
    type: string;
    parameters: { [key: string]: any };
  }>;
  outputName?: string;
  description?: string;
}

// 数据集标注请求
export interface DatasetAnnotationRequest {
  annotationType: 'manual' | 'auto' | 'semi_auto';
  labelFormat: string;
  annotationTasks: Array<{
    taskType: string;
    parameters: { [key: string]: any };
  }>;
  assignees?: string[];
  deadline?: string;
}

// 数据集版本创建请求
export interface DatasetVersionRequest {
  version: string;
  description?: string;
  changeLog: string;
  files: File[];
  labelFiles?: File[];
}

// 数据集统计信息
export interface DatasetStatistics {
  totalDatasets: number;
  publicDatasets: number;
  privateDatasets: number;
  datasetsByType: Record<DatasetType, number>;
  datasetsByFormat: Record<DatasetFormat, number>;
  totalSamples: number;
  totalSize: number;
  labeledDatasets: number;
  storageUsage: {
    total: number;
    used: number;
    available: number;
  };
  recentUploaded: Dataset[];
  topDownloaded: Dataset[];
}

// 数据预览信息
export interface DatasetPreview {
  sampleData: any[];
  schema: any;
  statistics: {
    sampleCount: number;
    columnCount?: number;
    missingValues?: number;
    dataTypes?: Record<string, string>;
  };
  visualizations?: {
    type: string;
    data: any;
  }[];
}

// 获取数据集列表
export const getDatasetList = (params?: DatasetListParams) => {
  return defHttp.get<BaseListResponse<Dataset>>('/api/v1/datasets', { params });
};

// 获取数据集详情
export const getDatasetDetail = (id: string) => {
  return defHttp.get<BaseResponse<Dataset>>(`/api/v1/datasets/${id}`);
};

// 创建数据集
export const createDataset = (data: DatasetCreateRequest) => {
  return defHttp.post<BaseResponse<Dataset>>('/api/v1/datasets', data);
};

// 上传数据集
export const uploadDataset = (data: DatasetUploadRequest) => {
  const formData = new FormData();
  
  // 添加基本信息
  Object.entries(data).forEach(([key, value]) => {
    if (key === 'files' || key === 'labelFiles') {
      // 文件数组单独处理
      return;
    } else if (key === 'tags' && Array.isArray(value)) {
      formData.append(key, JSON.stringify(value));
    } else if (typeof value === 'object') {
      formData.append(key, JSON.stringify(value));
    } else {
      formData.append(key, String(value));
    }
  });
  
  // 添加数据文件
  data.files.forEach((file) => {
    formData.append(`files`, file);
  });
  
  // 添加标注文件（如果有）
  if (data.labelFiles) {
    data.labelFiles.forEach((file) => {
      formData.append(`labelFiles`, file);
    });
  }

  return defHttp.post<BaseResponse<Dataset>>('/api/v1/datasets/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 更新数据集
export const updateDataset = (id: string, data: DatasetUpdateRequest) => {
  return defHttp.put<BaseResponse<Dataset>>(`/api/v1/datasets/${id}`, data);
};

// 删除数据集
export const deleteDataset = (id: string) => {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/datasets/${id}`);
};

// 批量删除数据集
export const batchDeleteDatasets = (ids: string[]) => {
  return defHttp.delete<BaseResponse<void>>('/api/v1/datasets/batch', { 
    data: { ids } 
  });
};

// 预览数据集
export const previewDataset = (id: string, params?: {
  sampleSize?: number;
  offset?: number;
}) => {
  return defHttp.get<BaseResponse<DatasetPreview>>(`/api/v1/datasets/${id}/preview`, { params });
};

// 下载数据集
export const downloadDataset = (id: string, version?: string) => {
  const params = version ? { version } : {};
  return defHttp.get<Blob>(`/api/v1/datasets/${id}/download`, {
    params,
    responseType: 'blob',
  });
};

// 数据集预处理
export const preprocessDataset = (id: string, data: DatasetPreprocessRequest) => {
  return defHttp.post<BaseResponse<Dataset>>(`/api/v1/datasets/${id}/preprocess`, data);
};

// 创建标注任务
export const createAnnotationTask = (id: string, data: DatasetAnnotationRequest) => {
  return defHttp.post<BaseResponse<any>>(`/api/v1/datasets/${id}/annotation`, data);
};

// 获取标注任务列表
export const getAnnotationTasks = (datasetId: string) => {
  return defHttp.get<BaseListResponse<any>>(`/api/v1/datasets/${datasetId}/annotation/tasks`);
};

// 创建数据集版本
export const createDatasetVersion = (datasetId: string, data: DatasetVersionRequest) => {
  const formData = new FormData();
  
  Object.entries(data).forEach(([key, value]) => {
    if (key === 'files' || key === 'labelFiles') {
      return;
    } else {
      formData.append(key, String(value));
    }
  });
  
  // 添加数据文件
  data.files.forEach((file) => {
    formData.append('files', file);
  });
  
  // 添加标注文件（如果有）
  if (data.labelFiles) {
    data.labelFiles.forEach((file) => {
      formData.append('labelFiles', file);
    });
  }

  return defHttp.post<BaseResponse<Dataset>>(`/api/v1/datasets/${datasetId}/versions`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 获取数据集版本历史
export const getDatasetVersions = (datasetId: string) => {
  return defHttp.get<BaseListResponse<DatasetVersion>>(`/api/v1/datasets/${datasetId}/versions`);
};

// 获取数据集统计信息
export const getDatasetStatistics = (workspaceId?: string) => {
  const params = workspaceId ? { workspaceId } : {};
  return defHttp.get<BaseResponse<DatasetStatistics>>('/api/v1/datasets/statistics', { params });
};

// 获取数据集使用记录
export const getDatasetUsageHistory = (datasetId: string, params?: {
  startTime?: string;
  endTime?: string;
  page?: number;
  pageSize?: number;
}) => {
  return defHttp.get<BaseListResponse<any>>(`/api/v1/datasets/${datasetId}/usage`, { params });
};

// 收藏/取消收藏数据集
export const toggleDatasetFavorite = (datasetId: string, favorite: boolean) => {
  return defHttp.post<BaseResponse<void>>(`/api/v1/datasets/${datasetId}/favorite`, { favorite });
};

// 获取推荐数据集
export const getRecommendedDatasets = (params?: {
  type?: DatasetType;
  format?: DatasetFormat;
  limit?: number;
}) => {
  return defHttp.get<BaseListResponse<Dataset>>('/api/v1/datasets/recommended', { params });
};

// 搜索公开数据集
export const searchPublicDatasets = (query: string, params?: {
  type?: DatasetType;
  format?: DatasetFormat;
  page?: number;
  pageSize?: number;
}) => {
  return defHttp.get<BaseListResponse<Dataset>>('/api/v1/datasets/search', {
    params: { query, ...params },
  });
};

// 获取支持的数据格式列表
export const getSupportedFormats = () => {
  return defHttp.get<BaseResponse<string[]>>('/api/v1/datasets/formats');
};

// 验证数据集文件
export const validateDatasetFiles = (files: File[]) => {
  const formData = new FormData();
  files.forEach((file) => {
    formData.append('files', file);
  });
  
  return defHttp.post<BaseResponse<{
    valid: boolean;
    errors?: string[];
    warnings?: string[];
    suggestions?: string[];
    estimatedSize?: number;
    estimatedSamples?: number;
  }>>('/api/v1/datasets/validate', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// 数据集质量分析
export const analyzeDatasetQuality = (id: string) => {
  return defHttp.post<BaseResponse<{
    overallScore: number;
    completeness: number;
    consistency: number;
    accuracy: number;
    issues: Array<{
      type: string;
      severity: 'low' | 'medium' | 'high';
      description: string;
      suggestions: string[];
    }>;
  }>>(`/api/v1/datasets/${id}/quality`);
};

// 数据集转换格式
export const convertDatasetFormat = (id: string, data: {
  targetFormat: DatasetFormat;
  options?: { [key: string]: any };
}) => {
  return defHttp.post<BaseResponse<Dataset>>(`/api/v1/datasets/${id}/convert`, data);
};

// 合并数据集
export const mergeDatasets = (data: {
  name: string;
  description?: string;
  sourceDatasetIds: string[];
  mergeStrategy: 'concat' | 'union' | 'intersect';
  workspaceId: string;
}) => {
  return defHttp.post<BaseResponse<Dataset>>('/api/v1/datasets/merge', data);
};

// 分割数据集
export const splitDataset = (id: string, data: {
  splitRatio: number[];
  splitNames: string[];
  stratify?: boolean;
  randomSeed?: number;
}) => {
  return defHttp.post<BaseResponse<Dataset[]>>(`/api/v1/datasets/${id}/split`, data);
};