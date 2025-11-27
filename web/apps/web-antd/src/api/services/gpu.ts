/**
 * GPU 资源管理服务 API
 * 处理 GPU 集群、节点、设备、使用记录等相关的接口调用
 */

import type {
  GpuApi,
  PageResponse,
  CrudResponse,
  ApiResponse,
} from '#/types/api';
import type { QueryParams } from './base';
import { BaseApiService } from './base';

export class GpuService extends BaseApiService {
  constructor() {
    super('/api/v1');
  }

  // ============ GPU 集群管理 ============

  /**
   * 获取 GPU 集群列表
   * @param params 查询参数
   */
  async getClusters(
    params?: QueryParams & {
      status?: GpuApi.ClusterStatus;
    },
  ): Promise<PageResponse<GpuApi.GpuCluster>> {
    return this.getList('/gpuclusters', params);
  }

  /**
   * 获取指定 GPU 集群详情
   * @param id 集群ID
   */
  async getCluster(id: string): Promise<GpuApi.GpuCluster> {
    return this.getById('/gpuclusters', id);
  }

  /**
   * 创建 GPU 集群
   * @param clusterData 集群数据
   */
  async createCluster(
    clusterData: Partial<GpuApi.GpuCluster>,
  ): Promise<CrudResponse> {
    const response = await this.create<CrudResponse>('/gpuclusters', clusterData);
    return response.data;
  }

  /**
   * 更新 GPU 集群
   * @param id 集群ID
   * @param clusterData 集群更新数据
   */
  async updateCluster(
    id: string,
    clusterData: Partial<GpuApi.GpuCluster>,
  ): Promise<CrudResponse> {
    const response = await this.update<CrudResponse>('/gpuclusters', id, clusterData);
    return response.data;
  }

  /**
   * 删除 GPU 集群
   * @param id 集群ID
   */
  async deleteCluster(id: string): Promise<CrudResponse> {
    const response = await this.remove('/gpuclusters', id);
    return response.data;
  }

  /**
   * 获取集群节点列表
   * @param clusterId 集群ID
   */
  async getClusterNodes(clusterId: string): Promise<GpuApi.GpuNode[]> {
    return this.get(`/gpuclusters/${clusterId}/nodes`);
  }

  /**
   * 向集群添加节点
   * @param nodeData 节点数据
   */
  async addNodeToCluster(nodeData: {
    clusterId: string;
    nodeId: string;
  }): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      '/gpuclusters/nodes',
      nodeData,
    );
    return response.data;
  }

  /**
   * 从集群移除节点
   * @param clusterId 集群ID
   * @param nodeId 节点ID
   */
  async removeNodeFromCluster(
    clusterId: string,
    nodeId: string,
  ): Promise<CrudResponse> {
    const response = await this.delete<ApiResponse<CrudResponse>>(
      `/gpuclusters/${clusterId}/nodes/${nodeId}`,
    );
    return response.data;
  }

  // ============ GPU 节点管理 ============

  /**
   * 获取 GPU 节点列表
   * @param params 查询参数
   */
  async getNodes(
    params?: QueryParams & {
      clusterId?: string;
      status?: GpuApi.DeviceStatus;
    },
  ): Promise<PageResponse<GpuApi.GpuNode>> {
    return this.getList('/gpunodes', params);
  }

  /**
   * 获取指定 GPU 节点详情
   * @param id 节点ID
   */
  async getNode(id: string): Promise<GpuApi.GpuNode> {
    return this.getById('/gpunodes', id);
  }

  /**
   * 创建 GPU 节点
   * @param nodeData 节点数据
   */
  async createNode(
    nodeData: Partial<GpuApi.GpuNode>,
  ): Promise<CrudResponse> {
    const response = await this.create<CrudResponse>('/gpunodes', nodeData);
    return response.data;
  }

  /**
   * 更新 GPU 节点
   * @param id 节点ID
   * @param nodeData 节点更新数据
   */
  async updateNode(
    id: string,
    nodeData: Partial<GpuApi.GpuNode>,
  ): Promise<CrudResponse> {
    const response = await this.update<CrudResponse>('/gpunodes', id, nodeData);
    return response.data;
  }

  /**
   * 删除 GPU 节点
   * @param id 节点ID
   */
  async deleteNode(id: string): Promise<CrudResponse> {
    const response = await this.remove('/gpunodes', id);
    return response.data;
  }

  /**
   * 获取节点设备列表
   * @param nodeId 节点ID
   */
  async getNodeDevices(nodeId: string): Promise<GpuApi.GpuDevice[]> {
    return this.get(`/gpunodes/${nodeId}/devices`);
  }

  /**
   * 向节点添加设备
   * @param deviceData 设备数据
   */
  async addDeviceToNode(deviceData: {
    nodeId: string;
    deviceId: string;
  }): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      '/gpunodes/devices',
      deviceData,
    );
    return response.data;
  }

  /**
   * 从节点移除设备
   * @param nodeId 节点ID
   * @param deviceId 设备ID
   */
  async removeDeviceFromNode(
    nodeId: string,
    deviceId: string,
  ): Promise<CrudResponse> {
    const response = await this.delete<ApiResponse<CrudResponse>>(
      `/gpunodes/${nodeId}/devices/${deviceId}`,
    );
    return response.data;
  }

  // ============ GPU 设备管理 ============

  /**
   * 获取 GPU 设备列表
   * @param params 查询参数
   */
  async getDevices(
    params?: QueryParams & {
      nodeId?: string;
      type?: GpuApi.GpuType;
      status?: GpuApi.DeviceStatus;
      isAllocated?: boolean;
    },
  ): Promise<PageResponse<GpuApi.GpuDevice>> {
    return this.getList('/gpudevices', params);
  }

  /**
   * 获取指定 GPU 设备详情
   * @param id 设备ID
   */
  async getDevice(id: string): Promise<GpuApi.GpuDevice> {
    return this.getById('/gpudevices', id);
  }

  /**
   * 创建 GPU 设备
   * @param deviceData 设备数据
   */
  async createDevice(
    deviceData: Partial<GpuApi.GpuDevice>,
  ): Promise<CrudResponse> {
    const response = await this.create<CrudResponse>('/gpudevices', deviceData);
    return response.data;
  }

  /**
   * 更新 GPU 设备
   * @param id 设备ID
   * @param deviceData 设备更新数据
   */
  async updateDevice(
    id: string,
    deviceData: Partial<GpuApi.GpuDevice>,
  ): Promise<CrudResponse> {
    const response = await this.update<CrudResponse>('/gpudevices', id, deviceData);
    return response.data;
  }

  /**
   * 删除 GPU 设备
   * @param id 设备ID
   */
  async deleteDevice(id: string): Promise<CrudResponse> {
    const response = await this.remove('/gpudevices', id);
    return response.data;
  }

  /**
   * 分配 GPU 设备
   * @param allocationData 分配数据
   */
  async allocateDevice(allocationData: {
    deviceId: string;
    jobId: string;
    userId: string;
  }): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      '/gpudevices/allocations',
      allocationData,
    );
    return response.data;
  }

  /**
   * 释放 GPU 设备
   * @param id 分配记录ID
   */
  async releaseDevice(id: string): Promise<CrudResponse> {
    const response = await this.delete<ApiResponse<CrudResponse>>(
      `/gpudevices/allocations/${id}`,
    );
    return response.data;
  }

  /**
   * 获取设备分配列表
   * @param params 查询参数
   */
  async getDeviceAllocations(params?: QueryParams): Promise<PageResponse<any>> {
    return this.getList('/gpudevices/allocations', params);
  }

  // ============ GPU 使用记录 ============

  /**
   * 获取 GPU 使用记录列表
   * @param params 查询参数
   */
  async getUsageRecords(
    params?: QueryParams & {
      deviceId?: string;
      userId?: string;
      jobId?: string;
      startTime?: string;
      endTime?: string;
    },
  ): Promise<PageResponse<GpuApi.GpuUsageRecord>> {
    return this.getList('/gpuusage', params);
  }

  /**
   * 获取指定使用记录详情
   * @param id 记录ID
   */
  async getUsageRecord(id: string): Promise<GpuApi.GpuUsageRecord> {
    return this.getById('/gpuusage', id);
  }

  /**
   * 创建使用记录
   * @param recordData 记录数据
   */
  async createUsageRecord(
    recordData: Partial<GpuApi.GpuUsageRecord>,
  ): Promise<CrudResponse> {
    const response = await this.create<CrudResponse>('/gpuusage', recordData);
    return response.data;
  }

  /**
   * 更新使用记录
   * @param id 记录ID
   * @param recordData 记录更新数据
   */
  async updateUsageRecord(
    id: string,
    recordData: Partial<GpuApi.GpuUsageRecord>,
  ): Promise<CrudResponse> {
    const response = await this.update<CrudResponse>('/gpuusage', id, recordData);
    return response.data;
  }

  /**
   * 删除使用记录
   * @param id 记录ID
   */
  async deleteUsageRecord(id: string): Promise<CrudResponse> {
    const response = await this.remove('/gpuusage', id);
    return response.data;
  }

  /**
   * 获取使用记录关联信息
   * @param usageRecordId 使用记录ID
   */
  async getUsageRelations(usageRecordId: string): Promise<any[]> {
    return this.get(`/gpuusage/${usageRecordId}/relations`);
  }

  /**
   * 添加使用记录关联
   * @param relationData 关联数据
   */
  async addUsageRelation(relationData: any): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      '/gpuusage/relations',
      relationData,
    );
    return response.data;
  }

  // ============ 统计和监控 ============

  /**
   * 获取 GPU 资源统计
   * @param params 查询参数
   */
  async getResourceStatistics(params?: {
    clusterId?: string;
    nodeId?: string;
    timeRange?: string;
  }): Promise<{
    totalGpus: number;
    availableGpus: number;
    allocatedGpus: number;
    utilizationRate: number;
    avgTemperature: number;
    avgPowerDraw: number;
  }> {
    return this.get('/gpu/statistics', params);
  }

  /**
   * 获取实时 GPU 监控数据
   * @param deviceId 设备ID (可选)
   */
  async getRealtimeMetrics(deviceId?: string): Promise<any[]> {
    return this.get('/gpu/metrics/realtime', { deviceId });
  }

  /**
   * 获取历史 GPU 监控数据
   * @param params 查询参数
   */
  async getHistoricalMetrics(params: {
    deviceIds?: string[];
    startTime: string;
    endTime: string;
    interval?: string;
  }): Promise<any[]> {
    return this.get('/gpu/metrics/historical', params);
  }
}

// 导出服务实例
export const gpuService = new GpuService();