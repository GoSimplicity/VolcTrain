/**
 * API 响应数据类型定义
 * 统一定义后端 API 返回的数据结构
 */

// 通用响应结构
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
  success: boolean;
}

// 分页响应结构
export interface PageResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// CRUD 响应结构
export interface CrudResponse {
  id: string;
  success: boolean;
  message: string;
}

// 用户信息
export interface UserInfo {
  id: string;
  username: string;
  email: string;
  fullName: string;
  avatar?: string;
  roles: string[];
  permissions: string[];
  createdAt: string;
  lastLoginAt?: string;
}

// 认证相关类型
export namespace AuthApi {
  export interface LoginRequest {
    username: string;
    password: string;
  }

  export interface LoginResult {
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    tokenType: string;
    userInfo: UserInfo;
  }

  export interface RefreshTokenRequest {
    refreshToken: string;
  }

  export interface RefreshTokenResult {
    accessToken: string;
    refreshToken: string;
    expiresIn: number;
    tokenType: string;
  }
}

// 训练任务相关类型
export namespace TrainingApi {
  // 训练任务状态
  export type JobStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled' | 'suspended';
  
  // 优先级
  export type Priority = 'low' | 'normal' | 'high' | 'urgent';

  // 训练框架类型
  export type FrameworkType = 'tensorflow' | 'pytorch' | 'mpi' | 'paddlepaddle' | 'mindspore';

  // 基础训练任务信息
  export interface TrainingJob {
    id: string;
    name: string;
    description?: string;
    status: JobStatus;
    priority: Priority;
    framework: FrameworkType;
    queueId: string;
    userId: string;
    workspaceId: string;
    
    // 资源配置
    resourceConfig: ResourceConfig;
    
    // 运行配置
    command: string[];
    image: string;
    workingDir?: string;
    env?: Record<string, string>;
    
    // 时间信息
    createdAt: string;
    startedAt?: string;
    completedAt?: string;
    estimatedDuration?: number;
    
    // 统计信息
    progress: number;
    metrics?: TrainingMetrics;
    
    // Volcano 任务配置
    volcanoSpec?: VolcanoJobSpec;
  }

  // 资源配置
  export interface ResourceConfig {
    cpu: string;
    memory: string;
    gpu: {
      type: string;
      count: number;
    };
    replicas: number;
    minAvailable?: number;
  }

  // 训练指标
  export interface TrainingMetrics {
    epoch: number;
    loss: number;
    accuracy?: number;
    learningRate: number;
    customMetrics?: Record<string, number>;
  }

  // Volcano 任务规范
  export interface VolcanoJobSpec {
    name: string;
    namespace: string;
    minAvailable: number;
    queue: string;
    priorityClassName?: string;
    maxRetry?: number;
    schedulerName: string;
    tasks: VolcanoTaskSpec[];
    plugins?: Record<string, any>;
    jobType: string;
  }

  // Volcano 任务规范
  export interface VolcanoTaskSpec {
    name: string;
    replicas: number;
    template: PodTemplateSpec;
    policies?: TaskPolicy[];
    minAvailable?: number;
    maxRetry?: number;
  }

  // Pod 模板规范
  export interface PodTemplateSpec {
    metadata?: ObjectMeta;
    spec: PodSpec;
    affinity?: any;
    nodeSelector?: Record<string, string>;
    tolerations?: any[];
    priorityClass?: string;
  }

  // K8s 对象元数据
  export interface ObjectMeta {
    name?: string;
    labels?: Record<string, string>;
    annotations?: Record<string, string>;
  }

  // Pod 规范
  export interface PodSpec {
    containers: Container[];
    initContainers?: Container[];
    volumes?: Volume[];
    serviceAccount?: string;
    restartPolicy?: string;
    hostNetwork?: boolean;
    securityContext?: any;
  }

  // 容器规范
  export interface Container {
    name: string;
    image: string;
    command?: string[];
    args?: string[];
    workingDir?: string;
    ports?: ContainerPort[];
    env?: EnvVar[];
    resources?: ResourceRequirements;
    volumeMounts?: VolumeMount[];
    livenessProbe?: any;
    readinessProbe?: any;
    imagePullPolicy?: string;
    securityContext?: any;
  }

  // 其他辅助类型
  export interface TaskPolicy {
    event: string;
    action: string;
  }

  export interface ContainerPort {
    name?: string;
    containerPort: number;
    protocol?: string;
  }

  export interface EnvVar {
    name: string;
    value?: string;
    valueFrom?: any;
  }

  export interface ResourceRequirements {
    limits?: Record<string, string>;
    requests?: Record<string, string>;
  }

  export interface VolumeMount {
    name: string;
    mountPath: string;
    subPath?: string;
    readOnly?: boolean;
  }

  export interface Volume {
    name: string;
    hostPath?: any;
    emptyDir?: any;
    configMap?: any;
    secret?: any;
    persistentVolumeClaim?: any;
  }

  // 创建训练任务请求
  export interface CreateTrainingJobRequest {
    name: string;
    description?: string;
    queueId: string;
    priority: Priority;
    framework: FrameworkType;
    
    // 资源配置
    resourceConfig: ResourceConfig;
    
    // 运行配置
    command: string[];
    image: string;
    workingDir?: string;
    env?: Record<string, string>;
    
    // 数据和输出配置
    datasetIds?: string[];
    outputPath?: string;
    
    // Volcano 配置
    volcanoSpec?: Partial<VolcanoJobSpec>;
  }

  // 更新训练任务请求
  export interface UpdateTrainingJobRequest {
    id: string;
    name?: string;
    description?: string;
    priority?: Priority;
    resourceConfig?: Partial<ResourceConfig>;
  }

  // 训练任务控制请求
  export interface TrainingJobControlRequest {
    id: string;
    action: 'start' | 'stop' | 'restart' | 'suspend' | 'resume' | 'cancel';
    reason?: string;
  }

  // 训练队列
  export interface TrainingQueue {
    id: string;
    name: string;
    description?: string;
    priority: number;
    maxRunningJobs: number;
    resourceQuota: ResourceQuota;
    policies: QueuePolicy[];
    createdAt: string;
    updatedAt: string;
  }

  // 资源配额
  export interface ResourceQuota {
    cpu: string;
    memory: string;
    gpu: number;
  }

  // 队列策略
  export interface QueuePolicy {
    type: string;
    config: Record<string, any>;
  }

  // 训练日志
  export interface TrainingLog {
    timestamp: string;
    level: 'info' | 'warn' | 'error' | 'debug';
    message: string;
    source: string;
  }

  // 检查点
  export interface Checkpoint {
    id: string;
    jobId: string;
    epoch: number;
    path: string;
    size: number;
    metrics: Record<string, number>;
    createdAt: string;
  }
}

// GPU 资源相关类型
export namespace GpuApi {
  // GPU 设备类型
  export type GpuType = 'V100' | 'A100' | 'T4' | 'RTX3080' | 'RTX3090' | 'H100';
  
  // 设备状态
  export type DeviceStatus = 'online' | 'offline' | 'maintenance' | 'error';
  
  // 集群状态
  export type ClusterStatus = 'healthy' | 'degraded' | 'unhealthy';

  // GPU 集群
  export interface GpuCluster {
    id: string;
    name: string;
    description?: string;
    status: ClusterStatus;
    nodeCount: number;
    totalGpus: number;
    availableGpus: number;
    allocatedGpus: number;
    createdAt: string;
    updatedAt: string;
  }

  // GPU 节点
  export interface GpuNode {
    id: string;
    clusterId: string;
    name: string;
    hostname: string;
    ip: string;
    status: DeviceStatus;
    os: string;
    kernelVersion: string;
    driverVersion: string;
    cudaVersion: string;
    
    // 硬件信息
    cpuCores: number;
    memory: string;
    storage: string;
    
    // GPU 设备
    gpuDevices: GpuDevice[];
    
    // 资源使用情况
    cpuUsage: number;
    memoryUsage: number;
    
    createdAt: string;
    updatedAt: string;
  }

  // GPU 设备
  export interface GpuDevice {
    id: string;
    nodeId: string;
    index: number;
    name: string;
    type: GpuType;
    memory: string;
    status: DeviceStatus;
    temperature: number;
    powerDraw: number;
    utilization: number;
    memoryUsage: number;
    
    // 分配信息
    isAllocated: boolean;
    allocatedTo?: string; // 任务ID
    allocatedAt?: string;
    
    createdAt: string;
    updatedAt: string;
  }

  // GPU 使用记录
  export interface GpuUsageRecord {
    id: string;
    deviceId: string;
    jobId: string;
    userId: string;
    startTime: string;
    endTime?: string;
    duration?: number;
    utilization: number;
    memoryUsage: number;
    energyConsumed?: number;
  }
}

// 数据集相关类型
export namespace DatasetApi {
  // 数据集类型
  export type DatasetType = 'image' | 'text' | 'audio' | 'video' | 'tabular' | 'other';
  
  // 数据集状态
  export type DatasetStatus = 'uploading' | 'processing' | 'ready' | 'error';

  // 数据集
  export interface Dataset {
    id: string;
    name: string;
    description?: string;
    type: DatasetType;
    status: DatasetStatus;
    size: number;
    fileCount: number;
    format: string;
    tags: string[];
    
    // 元数据
    metadata: DatasetMetadata;
    
    // 路径信息
    path: string;
    downloadUrl?: string;
    
    // 用户信息
    userId: string;
    workspaceId: string;
    isPublic: boolean;
    
    createdAt: string;
    updatedAt: string;
  }

  // 数据集元数据
  export interface DatasetMetadata {
    columns?: string[];
    classCount?: number;
    sampleCount: number;
    features?: Record<string, any>;
    statistics?: Record<string, any>;
  }

  // 创建数据集请求
  export interface CreateDatasetRequest {
    name: string;
    description?: string;
    type: DatasetType;
    tags?: string[];
    isPublic?: boolean;
    
    // 上传方式：'upload' | 'url' | 'path'
    uploadType: 'upload' | 'url' | 'path';
    source?: string; // URL 或路径
  }
}

// 模型相关类型
export namespace ModelApi {
  // 模型状态
  export type ModelStatus = 'training' | 'trained' | 'deployed' | 'archived';
  
  // 模型类型
  export type ModelType = 'classification' | 'regression' | 'detection' | 'segmentation' | 'nlp' | 'other';

  // 模型
  export interface Model {
    id: string;
    name: string;
    description?: string;
    type: ModelType;
    status: ModelStatus;
    framework: TrainingApi.FrameworkType;
    version: string;
    
    // 文件信息
    filePath: string;
    fileSize: number;
    checksum: string;
    
    // 性能指标
    metrics: Record<string, number>;
    
    // 关联信息
    datasetId?: string;
    trainingJobId?: string;
    
    // 用户信息
    userId: string;
    workspaceId: string;
    isPublic: boolean;
    
    // 部署信息
    deploymentInfo?: DeploymentInfo;
    
    createdAt: string;
    updatedAt: string;
  }

  // 部署信息
  export interface DeploymentInfo {
    endpoint?: string;
    replicas: number;
    resources: TrainingApi.ResourceConfig;
    status: 'deploying' | 'running' | 'stopped' | 'error';
  }
}

// 工作空间相关类型
export namespace WorkspaceApi {
  // 工作空间
  export interface Workspace {
    id: string;
    name: string;
    description?: string;
    ownerId: string;
    members: WorkspaceMember[];
    resourceQuota: TrainingApi.ResourceQuota;
    
    createdAt: string;
    updatedAt: string;
  }

  // 工作空间成员
  export interface WorkspaceMember {
    userId: string;
    username: string;
    role: 'owner' | 'admin' | 'member' | 'viewer';
    joinedAt: string;
  }
}

// 监控相关类型
export namespace MonitorApi {
  // 系统指标
  export interface SystemMetrics {
    timestamp: number;
    cpu: {
      usage: number;
      cores: number;
    };
    memory: {
      used: number;
      total: number;
      usage: number;
    };
    disk: {
      used: number;
      total: number;
      usage: number;
    };
    network: {
      rxBytes: number;
      txBytes: number;
    };
  }

  // GPU 指标
  export interface GpuMetrics {
    timestamp: number;
    deviceId: string;
    utilization: number;
    memoryUsage: number;
    temperature: number;
    powerDraw: number;
  }

  // 告警规则
  export interface AlertRule {
    id: string;
    name: string;
    description?: string;
    metric: string;
    condition: string;
    threshold: number;
    duration: number;
    severity: 'info' | 'warning' | 'critical';
    enabled: boolean;
    
    createdAt: string;
    updatedAt: string;
  }

  // 告警记录
  export interface AlertRecord {
    id: string;
    ruleId: string;
    ruleName: string;
    severity: 'info' | 'warning' | 'critical';
    status: 'firing' | 'resolved';
    message: string;
    value: number;
    
    triggeredAt: string;
    resolvedAt?: string;
  }
}