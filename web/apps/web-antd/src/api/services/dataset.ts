/**
 * 数据集管理服务 API
 * 处理数据集上传、管理、预览等相关的接口调用
 */

import type {
  DatasetApi,
  PageResponse,
  CrudResponse,
  ApiResponse,
} from '#/types/api';
import type { QueryParams } from './base';
import { BaseApiService } from './base';

export class DatasetService extends BaseApiService {
  constructor() {
    super('/api/v1/datasets');
  }

  /**
   * 获取数据集列表
   * @param params 查询参数
   */
  async getDatasets(
    params?: QueryParams & {
      type?: DatasetApi.DatasetType;
      status?: DatasetApi.DatasetStatus;
      userId?: string;
      isPublic?: boolean;
      tags?: string[];
    },
  ): Promise<PageResponse<DatasetApi.Dataset>> {
    return this.getList('', params);
  }

  /**
   * 获取指定数据集详情
   * @param id 数据集ID
   */
  async getDataset(id: string): Promise<DatasetApi.Dataset> {
    return this.getById('', id);
  }

  /**
   * 创建数据集
   * @param datasetData 数据集数据
   */
  async createDataset(
    datasetData: DatasetApi.CreateDatasetRequest,
  ): Promise<CrudResponse> {
    const response = await this.create<CrudResponse>('', datasetData);
    return response.data;
  }

  /**
   * 更新数据集信息
   * @param id 数据集ID
   * @param datasetData 数据集更新数据
   */
  async updateDataset(
    id: string,
    datasetData: Partial<DatasetApi.Dataset>,
  ): Promise<CrudResponse> {
    const response = await this.update<CrudResponse>('', id, datasetData);
    return response.data;
  }

  /**
   * 删除数据集
   * @param id 数据集ID
   */
  async deleteDataset(id: string): Promise<CrudResponse> {
    const response = await this.remove('', id);
    return response.data;
  }

  /**
   * 上传数据集文件
   * @param file 文件
   * @param onProgress 上传进度回调
   */
  async uploadDataset(
    file: File,
    onProgress?: (progressEvent: any) => void,
  ): Promise<{ uploadId: string; filename: string }> {
    const response = await this.upload<{ uploadId: string; filename: string }>(
      '/upload',
      file,
      onProgress,
    );
    return response.data;
  }

  /**
   * 批量上传数据集文件
   * @param files 文件列表
   * @param onProgress 上传进度回调
   */
  async uploadDatasetBatch(
    files: File[],
    onProgress?: (progressEvent: any) => void,
  ): Promise<{ uploadIds: string[]; filenames: string[] }> {
    const formData = new FormData();
    files.forEach((file) => {
      formData.append('files', file);
    });

    const response = await this.upload<{ uploadIds: string[]; filenames: string[] }>(
      '/upload/batch',
      formData,
      onProgress,
    );
    return response.data;
  }

  /**
   * 从 URL 导入数据集
   * @param url 数据集 URL
   * @param metadata 元数据
   */
  async importFromUrl(
    url: string,
    metadata: Partial<DatasetApi.Dataset>,
  ): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      '/import/url',
      { url, metadata },
    );
    return response.data;
  }

  /**
   * 从云存储导入数据集
   * @param path 云存储路径
   * @param metadata 元数据
   */
  async importFromStorage(
    path: string,
    metadata: Partial<DatasetApi.Dataset>,
  ): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      '/import/storage',
      { path, metadata },
    );
    return response.data;
  }

  /**
   * 预览数据集内容
   * @param id 数据集ID
   * @param params 预览参数
   */
  async previewDataset(
    id: string,
    params?: {
      limit?: number;
      offset?: number;
      sampleType?: 'random' | 'first' | 'last';
    },
  ): Promise<{
    samples: any[];
    total: number;
    metadata: DatasetApi.DatasetMetadata;
  }> {
    return this.get(`/${id}/preview`, params);
  }

  /**
   * 获取数据集统计信息
   * @param id 数据集ID
   */
  async getDatasetStatistics(id: string): Promise<{
    fileCount: number;
    totalSize: number;
    distribution: Record<string, number>;
    metadata: DatasetApi.DatasetMetadata;
  }> {
    return this.get(`/${id}/statistics`);
  }

  /**
   * 获取数据集版本列表
   * @param id 数据集ID
   */
  async getDatasetVersions(id: string): Promise<{
    versions: Array<{
      version: string;
      createdAt: string;
      description?: string;
      changes: string[];
    }>;
  }> {
    return this.get(`/${id}/versions`);
  }

  /**
   * 创建数据集版本
   * @param id 数据集ID
   * @param versionData 版本数据
   */
  async createDatasetVersion(
    id: string,
    versionData: {
      version: string;
      description?: string;
      changes?: string[];
    },
  ): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      `/${id}/versions`,
      versionData,
    );
    return response.data;
  }

  /**
   * 下载数据集
   * @param id 数据集ID
   * @param version 版本号（可选）
   */
  async downloadDataset(id: string, version?: string): Promise<Blob> {
    return this.download(`/${id}/download`, undefined, { version });
  }

  /**
   * 分享数据集
   * @param id 数据集ID
   * @param shareData 分享配置
   */
  async shareDataset(
    id: string,
    shareData: {
      isPublic: boolean;
      shareWithUsers?: string[];
      permissions?: ('read' | 'write' | 'download')[];
      expiresAt?: string;
    },
  ): Promise<{ shareLink?: string }> {
    const response = await this.post<ApiResponse<{ shareLink?: string }>>(
      `/${id}/share`,
      shareData,
    );
    return response.data;
  }

  /**
   * 数据集标注
   * @param id 数据集ID
   * @param annotations 标注数据
   */
  async annotateDataset(
    id: string,
    annotations: Array<{
      sampleId: string;
      labels: string[];
      boundingBoxes?: any[];
      metadata?: Record<string, any>;
    }>,
  ): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      `/${id}/annotations`,
      { annotations },
    );
    return response.data;
  }

  /**
   * 获取数据集标注
   * @param id 数据集ID
   * @param params 查询参数
   */
  async getDatasetAnnotations(
    id: string,
    params?: {
      sampleIds?: string[];
      labelFilter?: string[];
    },
  ): Promise<{
    annotations: Array<{
      sampleId: string;
      labels: string[];
      boundingBoxes?: any[];
      metadata?: Record<string, any>;
      createdAt: string;
      updatedAt: string;
    }>;
  }> {
    return this.get(`/${id}/annotations`, params);
  }

  /**
   * 数据集格式转换
   * @param id 数据集ID
   * @param targetFormat 目标格式
   */
  async convertDatasetFormat(
    id: string,
    targetFormat: string,
  ): Promise<{ taskId: string }> {
    const response = await this.post<ApiResponse<{ taskId: string }>>(
      `/${id}/convert`,
      { targetFormat },
    );
    return response.data;
  }

  /**
   * 获取转换任务状态
   * @param taskId 任务ID
   */
  async getConversionStatus(taskId: string): Promise<{
    status: 'pending' | 'running' | 'completed' | 'failed';
    progress: number;
    result?: {
      datasetId: string;
      downloadUrl: string;
    };
    error?: string;
  }> {
    return this.get(`/conversion/${taskId}`);
  }

  /**
   * 搜索公共数据集
   * @param params 搜索参数
   */
  async searchPublicDatasets(
    params?: QueryParams & {
      category?: string;
      domain?: string;
      minSize?: number;
      maxSize?: number;
      format?: string;
    },
  ): Promise<PageResponse<DatasetApi.Dataset>> {
    return this.getList('/public', params);
  }

  /**
   * 获取推荐数据集
   * @param params 推荐参数
   */
  async getRecommendedDatasets(params?: {
    basedOnDataset?: string;
    basedOnTask?: string;
    limit?: number;
  }): Promise<DatasetApi.Dataset[]> {
    return this.get('/recommendations', params);
  }

  /**
   * 数据集质量检查
   * @param id 数据集ID
   */
  async checkDatasetQuality(id: string): Promise<{
    score: number;
    issues: Array<{
      type: 'warning' | 'error';
      message: string;
      suggestion?: string;
    }>;
    recommendations: string[];
  }> {
    return this.get(`/${id}/quality-check`);
  }
}

// 导出服务实例
export const datasetService = new DatasetService();