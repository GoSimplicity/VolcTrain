import { defHttp } from '#/api/core/http';
import type {
  GPUDevice,
  GPUNode,
  GPUCluster,
  GPUQuery,
  NodeQuery,
  ClusterQuery,
  GPUStatistics,
  GPUUsageRecord,
  ResourceScheduleRequest,
  ResourceScheduleResponse,
  BillingConfig,
  CostStatistics,
} from '../types/gpu';
import type { BaseResponse, BaseListResponse } from '../types/common';

/**
 * GPU设备相关API
 */
export function getGPUList(params?: GPUQuery): Promise<BaseListResponse<GPUDevice>> {
  return defHttp.get<BaseListResponse<GPUDevice>>('/api/v1/gpu/devices', { params });
}

export function getGPUDetail(id: string): Promise<BaseResponse<GPUDevice>> {
  return defHttp.get<BaseResponse<GPUDevice>>(`/api/v1/gpu/devices/${id}`);
}

export function updateGPUStatus(id: string, status: 'available' | 'occupied' | 'maintenance'): Promise<BaseResponse<void>> {
  return defHttp.put<BaseResponse<void>>(`/api/v1/gpu/devices/${id}/status`, { status });
}

// 释放GPU设备
export function releaseGPU(id: string): Promise<BaseResponse<void>> {
  return defHttp.post<BaseResponse<void>>(`/api/v1/gpu/devices/${id}/release`);
}

// 设置GPU维护状态
export function maintainGPU(id: string, reason?: string): Promise<BaseResponse<void>> {
  return defHttp.post<BaseResponse<void>>(`/api/v1/gpu/devices/${id}/maintain`, { reason });
}

export function getGPUMetrics(id: string, timeRange?: string): Promise<BaseResponse<any>> {
  return defHttp.get<BaseResponse<any>>(`/api/v1/gpu/devices/${id}/metrics`, { 
    params: timeRange ? { timeRange } : undefined 
  });
}

/**
 * GPU节点相关API
 */
export function getNodeList(params?: NodeQuery): Promise<BaseListResponse<GPUNode>> {
  return defHttp.get<BaseListResponse<GPUNode>>('/api/v1/gpu/nodes', { params });
}

export function getNodeDetail(id: string): Promise<BaseResponse<GPUNode>> {
  return defHttp.get<BaseResponse<GPUNode>>(`/api/v1/gpu/nodes/${id}`);
}

export function addNode(data: Partial<GPUNode>): Promise<BaseResponse<GPUNode>> {
  return defHttp.post<BaseResponse<GPUNode>>('/api/v1/gpu/nodes', data);
}

export function updateNode(id: string, data: Partial<GPUNode>): Promise<BaseResponse<GPUNode>> {
  return defHttp.put<BaseResponse<GPUNode>>(`/api/v1/gpu/nodes/${id}`, data);
}

export function removeNode(id: string): Promise<BaseResponse<void>> {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/gpu/nodes/${id}`);
}

/**
 * GPU集群相关API
 */
export function getClusterList(params?: ClusterQuery): Promise<BaseListResponse<GPUCluster>> {
  return defHttp.get<BaseListResponse<GPUCluster>>('/api/v1/gpu/clusters', { params });
}

export function getClusterDetail(id: string): Promise<BaseResponse<GPUCluster>> {
  return defHttp.get<BaseResponse<GPUCluster>>(`/api/v1/gpu/clusters/${id}`);
}

export function createCluster(data: Partial<GPUCluster>): Promise<BaseResponse<GPUCluster>> {
  return defHttp.post<BaseResponse<GPUCluster>>('/api/v1/gpu/clusters', data);
}

export function updateCluster(id: string, data: Partial<GPUCluster>): Promise<BaseResponse<GPUCluster>> {
  return defHttp.put<BaseResponse<GPUCluster>>(`/api/v1/gpu/clusters/${id}`, data);
}

export function deleteCluster(id: string): Promise<BaseResponse<void>> {
  return defHttp.delete<BaseResponse<void>>(`/api/v1/gpu/clusters/${id}`);
}

/**
 * 资源调度相关API
 */
export function scheduleResource(data: ResourceScheduleRequest): Promise<BaseResponse<ResourceScheduleResponse>> {
  return defHttp.post<BaseResponse<ResourceScheduleResponse>>('/api/v1/gpu/schedule', data);
}

// 批量资源调度
export function scheduleResources(requests: ResourceScheduleRequest[]): Promise<BaseResponse<ResourceScheduleResponse[]>> {
  return defHttp.post<BaseResponse<ResourceScheduleResponse[]>>('/api/v1/gpu/schedule/batch', { requests });
}

export function getScheduleHistory(params?: {
  startTime?: string;
  endTime?: string;
  status?: string;
  page?: number;
  pageSize?: number;
}): Promise<BaseListResponse<any>> {
  return defHttp.get<BaseListResponse<any>>('/api/v1/gpu/schedule/history', { params });
}

/**
 * 统计信息相关API
 */
export function getGPUStatistics(clusterId?: string): Promise<BaseResponse<GPUStatistics>> {
  return defHttp.get<BaseResponse<GPUStatistics>>('/api/v1/gpu/statistics', { 
    params: clusterId ? { clusterId } : undefined 
  });
}

export function getUsageHistory(params?: {
  deviceId?: string;
  startTime?: string;
  endTime?: string;
  page?: number;
  pageSize?: number;
}): Promise<BaseListResponse<GPUUsageRecord>> {
  return defHttp.get<BaseListResponse<GPUUsageRecord>>('/api/v1/gpu/usage', { params });
}

export function exportUsageReport(params?: {
  format: 'csv' | 'excel';
  startTime?: string;
  endTime?: string;
  deviceIds?: string[];
}): Promise<Blob> {
  return defHttp.get<Blob>('/api/v1/gpu/usage/export', { 
    params,
    responseType: 'blob' 
  });
}

/**
 * 计费相关API
 */
export function getBillingConfig(): Promise<BaseResponse<BillingConfig>> {
  return defHttp.get<BaseResponse<BillingConfig>>('/api/v1/gpu/billing/config');
}

export function updateBillingConfig(data: BillingConfig): Promise<BaseResponse<BillingConfig>> {
  return defHttp.put<BaseResponse<BillingConfig>>('/api/v1/gpu/billing/config', data);
}

export function getCostStatistics(params?: {
  workspaceId?: string;
  userId?: string;
  startTime?: string;
  endTime?: string;
}): Promise<BaseResponse<CostStatistics>> {
  return defHttp.get<BaseResponse<CostStatistics>>('/api/v1/gpu/billing/statistics', { params });
}

export function generateBill(params: {
  workspaceId?: string;
  userId?: string;
  startTime: string;
  endTime: string;
}): Promise<BaseResponse<any>> {
  return defHttp.post<BaseResponse<any>>('/api/v1/gpu/billing/generate', params);
}

/**
 * 实时监控相关API
 */
export function getRealtimeMetrics(deviceIds?: string[]): Promise<BaseResponse<any>> {
  return defHttp.get<BaseResponse<any>>('/api/v1/gpu/realtime/metrics', { 
    params: deviceIds ? { deviceIds: deviceIds.join(',') } : undefined 
  });
}

export function getNodeHealth(nodeId: string): Promise<BaseResponse<any>> {
  return defHttp.get<BaseResponse<any>>(`/api/v1/gpu/nodes/${nodeId}/health`);
}