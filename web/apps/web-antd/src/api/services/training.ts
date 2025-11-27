/**
 * 训练任务服务 API
 * 处理训练任务、队列、实验等相关的接口调用
 */

import type {
  TrainingApi,
  PageResponse,
  CrudResponse,
  ApiResponse,
} from '#/types/api';
import type { QueryParams } from './base';
import { BaseApiService } from './base';

export class TrainingService extends BaseApiService {
  constructor() {
    super('/api/v1/training');
  }

  // ============ 训练任务相关 ============

  /**
   * 获取训练任务列表
   * @param params 查询参数
   */
  async getJobs(
    params?: QueryParams & {
      status?: TrainingApi.JobStatus;
      userId?: string;
      queueId?: string;
      framework?: TrainingApi.FrameworkType;
    },
  ): Promise<PageResponse<TrainingApi.TrainingJob>> {
    return this.getList('/jobs', params);
  }

  /**
   * 获取指定训练任务详情
   * @param id 任务ID
   */
  async getJob(id: string): Promise<TrainingApi.TrainingJob> {
    return this.getById('/jobs', id);
  }

  /**
   * 创建训练任务
   * @param jobData 任务数据
   */
  async createJob(
    jobData: TrainingApi.CreateTrainingJobRequest,
  ): Promise<CrudResponse> {
    const response = await this.create<CrudResponse>('/jobs', jobData);
    return response.data;
  }

  /**
   * 更新训练任务
   * @param jobData 任务更新数据
   */
  async updateJob(
    jobData: TrainingApi.UpdateTrainingJobRequest,
  ): Promise<CrudResponse> {
    const response = await this.update<CrudResponse>(
      '/jobs',
      jobData.id,
      jobData,
    );
    return response.data;
  }

  /**
   * 删除训练任务
   * @param id 任务ID
   */
  async deleteJob(id: string): Promise<CrudResponse> {
    const response = await this.remove('/jobs', id);
    return response.data;
  }

  /**
   * 控制训练任务（启动、停止、重启等）
   * @param controlData 控制数据
   */
  async controlJob(
    controlData: TrainingApi.TrainingJobControlRequest,
  ): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      `/jobs/${controlData.id}/${controlData.action}`,
      { reason: controlData.reason },
    );
    return response.data;
  }

  /**
   * 获取训练任务日志
   * @param id 任务ID
   * @param params 查询参数
   */
  async getJobLogs(
    id: string,
    params?: {
      startTime?: string;
      endTime?: string;
      level?: string;
      lines?: number;
    },
  ): Promise<TrainingApi.TrainingLog[]> {
    return this.get(`/jobs/${id}/logs`, params);
  }

  /**
   * 获取训练任务指标
   * @param id 任务ID
   * @param params 查询参数
   */
  async getJobMetrics(
    id: string,
    params?: {
      startTime?: string;
      endTime?: string;
      step?: number;
    },
  ): Promise<TrainingApi.TrainingMetrics[]> {
    return this.get(`/jobs/${id}/metrics`, params);
  }

  /**
   * 获取训练任务检查点
   * @param id 任务ID
   */
  async getJobCheckpoints(id: string): Promise<TrainingApi.Checkpoint[]> {
    return this.get(`/jobs/${id}/checkpoints`);
  }

  /**
   * 获取训练任务实例信息
   * @param id 任务ID
   */
  async getJobInstances(id: string): Promise<any[]> {
    return this.get(`/jobs/${id}/instances`);
  }

  /**
   * 获取训练任务关联信息
   * @param id 任务ID
   */
  async getJobRelations(id: string): Promise<any> {
    return this.get(`/jobs/${id}/relations`);
  }

  /**
   * 获取任务创建选项（框架、资源等）
   */
  async getJobOptions(): Promise<{
    frameworks: TrainingApi.FrameworkType[];
    gpuTypes: string[];
    queues: TrainingApi.TrainingQueue[];
  }> {
    return this.get('/jobs/options');
  }

  /**
   * 批量控制训练任务
   * @param jobIds 任务ID列表
   * @param action 操作类型
   */
  async batchControlJobs(
    jobIds: string[],
    action: string,
  ): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      '/jobs/batch/control',
      { jobIds, action },
    );
    return response.data;
  }

  /**
   * 批量删除训练任务
   * @param jobIds 任务ID列表
   */
  async batchDeleteJobs(jobIds: string[]): Promise<CrudResponse> {
    const response = await this.batchRemove('/jobs', jobIds);
    return response.data;
  }

  /**
   * 克隆训练任务
   * @param id 原任务ID
   * @param cloneData 克隆配置
   */
  async cloneJob(
    id: string,
    cloneData?: { name?: string; description?: string },
  ): Promise<CrudResponse> {
    const response = await this.post<ApiResponse<CrudResponse>>(
      `/jobs/${id}/clone`,
      cloneData,
    );
    return response.data;
  }

  // ============ 训练队列相关 ============

  /**
   * 获取训练队列列表
   * @param params 查询参数
   */
  async getQueues(params?: QueryParams): Promise<PageResponse<TrainingApi.TrainingQueue>> {
    return this.getList('/queues', params);
  }

  /**
   * 获取指定训练队列详情
   * @param id 队列ID
   */
  async getQueue(id: string): Promise<TrainingApi.TrainingQueue> {
    return this.getById('/queues', id);
  }

  /**
   * 创建训练队列
   * @param queueData 队列数据
   */
  async createQueue(
    queueData: Partial<TrainingApi.TrainingQueue>,
  ): Promise<CrudResponse> {
    const response = await this.create<CrudResponse>('/queues', queueData);
    return response.data;
  }

  /**
   * 更新训练队列
   * @param id 队列ID
   * @param queueData 队列更新数据
   */
  async updateQueue(
    id: string,
    queueData: Partial<TrainingApi.TrainingQueue>,
  ): Promise<CrudResponse> {
    const response = await this.update<CrudResponse>('/queues', id, queueData);
    return response.data;
  }

  /**
   * 删除训练队列
   * @param id 队列ID
   */
  async deleteQueue(id: string): Promise<CrudResponse> {
    const response = await this.remove('/queues', id);
    return response.data;
  }

  /**
   * 获取队列统计信息
   * @param id 队列ID
   */
  async getQueueStatistics(id: string): Promise<{
    runningJobs: number;
    pendingJobs: number;
    completedJobs: number;
    failedJobs: number;
    resourceUsage: TrainingApi.ResourceQuota;
  }> {
    return this.get(`/queues/${id}/statistics`);
  }

  /**
   * 获取队列选项
   */
  async getQueueOptions(): Promise<{
    priorities: number[];
    policies: string[];
  }> {
    return this.get('/queues/options');
  }

  // ============ 统计和报告 ============

  /**
   * 获取训练统计信息
   * @param params 查询参数
   */
  async getTrainingStatistics(params?: {
    startDate?: string;
    endDate?: string;
    groupBy?: 'day' | 'week' | 'month';
  }): Promise<{
    totalJobs: number;
    completedJobs: number;
    failedJobs: number;
    runningJobs: number;
    averageDuration: number;
    resourceUsage: any;
  }> {
    return this.get('/statistics', params);
  }

  /**
   * 获取用户的训练任务列表
   * @param params 查询参数
   */
  async getMyJobs(params?: QueryParams): Promise<PageResponse<TrainingApi.TrainingJob>> {
    return this.getList('/jobs/my', params);
  }

  /**
   * 获取可用队列列表
   */
  async getAvailableQueues(): Promise<TrainingApi.TrainingQueue[]> {
    return this.get('/queues/available');
  }
}

// 导出服务实例
export const trainingService = new TrainingService();