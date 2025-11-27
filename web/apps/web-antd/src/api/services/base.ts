/**
 * 基础 API 服务类
 * 提供通用的 HTTP 请求方法和错误处理
 */

import type { RequestClientOptions } from '@vben/request';
import { requestClient } from '#/api/request';
import type { ApiResponse, PageResponse } from '#/types/api';

// 通用查询参数
export interface QueryParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  [key: string]: any;
}

// 基础 API 服务类
export class BaseApiService {
  protected baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  /**
   * 构建完整的 API URL
   */
  protected buildUrl(path: string): string {
    const cleanBaseUrl = this.baseUrl.replace(/\/$/, '');
    const cleanPath = path.replace(/^\//, '');
    return `${cleanBaseUrl}/${cleanPath}`;
  }

  /**
   * GET 请求
   */
  protected async get<T>(
    path: string,
    params?: QueryParams,
    options?: RequestClientOptions,
  ): Promise<T> {
    return requestClient.get<T>(this.buildUrl(path), {
      params,
      ...options,
    });
  }

  /**
   * POST 请求
   */
  protected async post<T>(
    path: string,
    data?: any,
    options?: RequestClientOptions,
  ): Promise<T> {
    return requestClient.post<T>(this.buildUrl(path), data, options);
  }

  /**
   * PUT 请求
   */
  protected async put<T>(
    path: string,
    data?: any,
    options?: RequestClientOptions,
  ): Promise<T> {
    return requestClient.put<T>(this.buildUrl(path), data, options);
  }

  /**
   * DELETE 请求
   */
  protected async delete<T>(
    path: string,
    data?: any,
    options?: RequestClientOptions,
  ): Promise<T> {
    return requestClient.delete<T>(this.buildUrl(path), { data, ...options });
  }

  /**
   * PATCH 请求
   */
  protected async patch<T>(
    path: string,
    data?: any,
    options?: RequestClientOptions,
  ): Promise<T> {
    return requestClient.patch<T>(this.buildUrl(path), data, options);
  }

  /**
   * 分页查询
   */
  protected async getList<T>(
    path: string,
    params?: QueryParams,
  ): Promise<PageResponse<T>> {
    return this.get<PageResponse<T>>(path, params);
  }

  /**
   * 根据 ID 获取单个资源
   */
  protected async getById<T>(path: string, id: string): Promise<T> {
    return this.get<T>(`${path}/${id}`);
  }

  /**
   * 创建资源
   */
  protected async create<T>(
    path: string,
    data: any,
  ): Promise<ApiResponse<T>> {
    return this.post<ApiResponse<T>>(path, data);
  }

  /**
   * 更新资源
   */
  protected async update<T>(
    path: string,
    id: string,
    data: any,
  ): Promise<ApiResponse<T>> {
    return this.put<ApiResponse<T>>(`${path}/${id}`, data);
  }

  /**
   * 删除资源
   */
  protected async remove(path: string, id: string): Promise<ApiResponse<void>> {
    return this.delete<ApiResponse<void>>(`${path}/${id}`);
  }

  /**
   * 批量删除
   */
  protected async batchRemove(
    path: string,
    ids: string[],
  ): Promise<ApiResponse<void>> {
    return this.delete<ApiResponse<void>>(`${path}/batch`, { ids });
  }

  /**
   * 文件上传
   */
  protected async upload<T>(
    path: string,
    file: File | FormData,
    onProgress?: (progressEvent: any) => void,
  ): Promise<ApiResponse<T>> {
    const formData = file instanceof FormData ? file : new FormData();
    if (file instanceof File) {
      formData.append('file', file);
    }

    return requestClient.post<ApiResponse<T>>(this.buildUrl(path), formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: onProgress,
    });
  }

  /**
   * 文件下载
   */
  protected async download(
    path: string,
    filename?: string,
    params?: QueryParams,
  ): Promise<Blob> {
    const response = await requestClient.get(this.buildUrl(path), {
      params,
      responseType: 'blob',
    });

    // 如果指定了文件名，创建下载链接
    if (filename) {
      const url = window.URL.createObjectURL(response as Blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    }

    return response as Blob;
  }

  /**
   * 健康检查
   */
  protected async healthCheck(): Promise<{ status: string; timestamp: number }> {
    return this.get('/health');
  }
}

// 错误处理工具函数
export const handleApiError = (error: any): string => {
  if (error?.response?.data?.message) {
    return error.response.data.message;
  }
  
  if (error?.message) {
    return error.message;
  }
  
  return '未知错误，请稍后重试';
};

// API 响应工具函数
export const isApiSuccess = <T>(response: ApiResponse<T>): boolean => {
  return response.code === 0 || response.success === true;
};

export const getApiData = <T>(response: ApiResponse<T>): T => {
  return response.data;
};

export const getApiMessage = <T>(response: ApiResponse<T>): string => {
  return response.message || '';
};