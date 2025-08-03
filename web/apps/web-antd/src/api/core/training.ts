import { requestClient } from '#/api/request';
import type {
  TrainingJob,
  TrainingJobQuery,
  CreateTrainingJobRequest,
  UpdateTrainingJobRequest,
  TrainingJobControlRequest,
  TrainingQueue,
  TrainingQueueQuery,
  TrainingStatistics,
  TrainingTemplate,
  Experiment,
  ExperimentRun,
  TrainingLog,
  TrainingMetrics,
  Checkpoint,
  PageResponse,
  CrudResponse,
} from '#/api/types';

/**
 * 训练任务管理API接口
 */

// 训练任务管理
export const getTrainingJobList = (params: TrainingJobQuery) => {
  return requestClient.get<PageResponse<TrainingJob>>('/api/v1/training/jobs', { params });
};

export const getTrainingJobById = (id: string) => {
  return requestClient.get<TrainingJob>(`/api/v1/training/jobs/${id}`);
};

export const createTrainingJob = (data: CreateTrainingJobRequest) => {
  return requestClient.post<CrudResponse>('/api/v1/training/jobs', data);
};

export const updateTrainingJob = (data: UpdateTrainingJobRequest) => {
  return requestClient.put<CrudResponse>(`/api/v1/training/jobs/${data.id}`, data);
};

export const deleteTrainingJob = (id: string) => {
  return requestClient.delete<CrudResponse>(`/api/v1/training/jobs/${id}`);
};

export const controlTrainingJob = (data: TrainingJobControlRequest) => {
  return requestClient.post<CrudResponse>(`/api/v1/training/jobs/${data.id}/control`, data);
};

export const getTrainingJobLogs = (id: string, params?: { lines?: number; follow?: boolean }) => {
  return requestClient.get<TrainingLog[]>(`/api/v1/training/jobs/${id}/logs`, { params });
};

export const getTrainingJobMetrics = (id: string, params?: { timeRange?: string }) => {
  return requestClient.get<TrainingMetrics[]>(`/api/v1/training/jobs/${id}/metrics`, { params });
};

export const getTrainingJobCheckpoints = (id: string) => {
  return requestClient.get<Checkpoint[]>(`/api/v1/training/jobs/${id}/checkpoints`);
};

// 训练队列管理
export const getTrainingQueueList = (params: TrainingQueueQuery) => {
  return requestClient.get<PageResponse<TrainingQueue>>('/api/v1/training/queues', { params });
};

export const getTrainingQueueById = (id: string) => {
  return requestClient.get<TrainingQueue>(`/api/v1/training/queues/${id}`);
};

export const createTrainingQueue = (data: any) => {
  return requestClient.post<CrudResponse>('/api/v1/training/queues', data);
};

export const updateTrainingQueue = (id: string, data: any) => {
  return requestClient.put<CrudResponse>(`/api/v1/training/queues/${id}`, data);
};

export const deleteTrainingQueue = (id: string) => {
  return requestClient.delete<CrudResponse>(`/api/v1/training/queues/${id}`);
};

export const getQueueStatistics = (id: string) => {
  return requestClient.get<any>(`/api/v1/training/queues/${id}/statistics`);
};

// 训练统计
export const getTrainingStatistics = (params?: { timeRange?: string; workspaceId?: string }) => {
  return requestClient.get<TrainingStatistics>('/api/v1/training/statistics', { params });
};

// 训练模板管理
export const getTrainingTemplateList = (params?: { category?: string; framework?: string; page?: number; pageSize?: number }) => {
  return requestClient.get<PageResponse<TrainingTemplate>>('/api/v1/training/templates', { params });
};

export const getTrainingTemplateById = (id: string) => {
  return requestClient.get<TrainingTemplate>(`/api/v1/training/templates/${id}`);
};

export const createTrainingTemplate = (data: any) => {
  return requestClient.post<CrudResponse>('/api/v1/training/templates', data);
};

export const updateTrainingTemplate = (id: string, data: any) => {
  return requestClient.put<CrudResponse>(`/api/v1/training/templates/${id}`, data);
};

export const deleteTrainingTemplate = (id: string) => {
  return requestClient.delete<CrudResponse>(`/api/v1/training/templates/${id}`);
};

export const useTrainingTemplate = (id: string, data: any) => {
  return requestClient.post<CreateTrainingJobRequest>(`/api/v1/training/templates/${id}/use`, data);
};

// 实验跟踪
export const getExperimentList = (params?: { projectId?: string; page?: number; pageSize?: number }) => {
  return requestClient.get<PageResponse<Experiment>>('/api/v1/training/experiments', { params });
};

export const getExperimentById = (id: string) => {
  return requestClient.get<Experiment>(`/api/v1/training/experiments/${id}`);
};

export const createExperiment = (data: any) => {
  return requestClient.post<CrudResponse>('/api/v1/training/experiments', data);
};

export const updateExperiment = (id: string, data: any) => {
  return requestClient.put<CrudResponse>(`/api/v1/training/experiments/${id}`, data);
};

export const deleteExperiment = (id: string) => {
  return requestClient.delete<CrudResponse>(`/api/v1/training/experiments/${id}`);
};

export const getExperimentRuns = (experimentId: string) => {
  return requestClient.get<ExperimentRun[]>(`/api/v1/training/experiments/${experimentId}/runs`);
};

export const createExperimentRun = (experimentId: string, data: any) => {
  return requestClient.post<CrudResponse>(`/api/v1/training/experiments/${experimentId}/runs`, data);
};

export const compareExperimentRuns = (experimentId: string, runIds: string[]) => {
  return requestClient.post<any>(`/api/v1/training/experiments/${experimentId}/compare`, { runIds });
};

// 我的训练任务
export const getMyTrainingJobs = (params?: TrainingJobQuery) => {
  return requestClient.get<PageResponse<TrainingJob>>('/api/v1/training/jobs/my', { params });
};

// 获取可用的训练队列
export const getAvailableQueues = (workspaceId?: string) => {
  return requestClient.get<TrainingQueue[]>('/api/v1/training/queues/available', { 
    params: { workspaceId } 
  });
};

// 批量操作
export const batchControlJobs = (jobIds: string[], action: string) => {
  return requestClient.post<CrudResponse>('/api/v1/training/jobs/batch/control', { jobIds, action });
};

export const batchDeleteJobs = (jobIds: string[]) => {
  return requestClient.delete<CrudResponse>('/api/v1/training/jobs/batch', { data: { jobIds } });
};

// 任务克隆
export const cloneTrainingJob = (id: string, data: { name: string; workspaceId?: string }) => {
  return requestClient.post<CrudResponse>(`/api/v1/training/jobs/${id}/clone`, data);
};