import { requestClient } from '#/api/request';

export enum JobStatus {
  PENDING = 'pending',
  RUNNING = 'running',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'cancelled',
  PAUSED = 'paused',
}

export enum JobPriority {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  URGENT = 'urgent',
}

export enum PodStatus {
  PENDING = 'Pending',
  RUNNING = 'Running',
  SUCCEEDED = 'Succeeded',
  FAILED = 'Failed',
  UNKNOWN = 'Unknown',
}

export enum EventType {
  NORMAL = 'Normal',
  WARNING = 'Warning',
  ERROR = 'Error',
}

export interface JobResources {
  cpu: number;
  memory: number;
  gpu?: number;
  storage?: number;
}

export interface EnvVar {
  key: string;
  value: string;
}

export interface VolumeMount {
  name: string;
  mountPath: string;
  readOnly?: boolean;
  subPath?: string;
}

export interface TrainingJob {
  id: string;
  name: string;
  namespace: string;
  queue: string;
  status: JobStatus;
  priority: JobPriority;
  creator: string;
  image: string;
  createTime: string;
  startTime?: string;
  completionTime?: string;
  duration: number;
  progress: number;
  resources: JobResources;
  replicas: number;
  maxRetries: number;
  currentRetries: number;
  script: string;
  envVars?: EnvVar[];
  dataPath?: string;
  outputPath?: string;
  volumeMounts?: VolumeMount[];
  description?: string;
  labels?: Record<string, string>;
  annotations?: Record<string, string>;
}

export interface PodInfo {
  name: string;
  status: PodStatus;
  nodeName: string;
  podIP?: string;
  hostIP?: string;
  startTime: string;
  resources: JobResources;
  restartCount: number;
  message?: string;
  reason?: string;
}

export interface EventInfo {
  id: string;
  type: EventType;
  time: string;
  message: string;
  source?: string;
  object?: string;
  count?: number;
}

export interface TrainingMetrics {
  loss?: number;
  accuracy?: number;
  learningRate?: number;
  epoch?: number;
  step?: number;
  timestamp: string;
}

export interface CreateJobRequest {
  name: string;
  namespace: string;
  queue: string;
  priority: JobPriority;
  image: string;
  resources: JobResources;
  replicas: number;
  maxRetries: number;
  script: string;
  envVars?: EnvVar[];
  dataPath?: string;
  outputPath?: string;
  volumeMounts?: VolumeMount[];
  description?: string;
  labels?: Record<string, string>;
}

export interface UpdateJobRequest {
  id: string;
  priority?: JobPriority;
  replicas?: number;
  description?: string;
  labels?: Record<string, string>;
}

export interface CloneJobRequest {
  sourceJobId: string;
  name: string;
  queue?: string;
  priority?: JobPriority;
}

export interface QueryJobsRequest {
  page?: number;
  pageSize?: number;
  namespace?: string;
  queue?: string;
  status?: JobStatus;
  priority?: JobPriority;
  creator?: string;
  keyword?: string;
  createTimeStart?: string;
  createTimeEnd?: string;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface PageResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface JobStatistics {
  total: number;
  running: number;
  pending: number;
  completed: number;
  failed: number;
  cancelled: number;
  paused: number;
}
