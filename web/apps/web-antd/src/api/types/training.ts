import type { CommonStatus, ResourceQuota, Labels, Annotations, PageQuery, BaseStatistics } from './common';

/**
 * 训练任务管理类型定义
 */

// 训练任务状态枚举
export enum TrainingJobStatus {
  PENDING = 'pending',       // 等待中
  QUEUED = 'queued',        // 队列中
  RUNNING = 'running',       // 运行中
  COMPLETED = 'completed',   // 已完成
  FAILED = 'failed',         // 失败
  CANCELLED = 'cancelled',   // 已取消
  PAUSED = 'paused',         // 已暂停
  STOPPED = 'stopped',       // 已停止
}

// 训练优先级枚举
export enum TrainingPriority {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  URGENT = 'urgent',
}

// 训练框架枚举
export enum TrainingFramework {
  TENSORFLOW = 'tensorflow',
  PYTORCH = 'pytorch',
  KERAS = 'keras',
  PADDLE = 'paddlepaddle',
  MINDSPORE = 'mindspore',
  CUSTOM = 'custom',
}

// 分布式训练类型
export enum DistributedType {
  SINGLE = 'single',         // 单机
  DATA_PARALLEL = 'data_parallel',     // 数据并行
  MODEL_PARALLEL = 'model_parallel',   // 模型并行
  PIPELINE_PARALLEL = 'pipeline_parallel', // 流水线并行
  HYBRID = 'hybrid',         // 混合并行
}

// 环境变量
export interface EnvVar {
  key: string;
  value: string;
}

// 存储卷挂载
export interface VolumeMount {
  name: string;
  mountPath: string;
  readOnly?: boolean;
  subPath?: string;
}

// 训练指标
export interface TrainingMetrics {
  // 基础指标
  loss?: number;
  accuracy?: number;
  learningRate?: number;
  epoch?: number;
  step?: number;
  
  // 时间指标
  epochTime?: number;        // 每个epoch的时间(秒)
  samplesPerSecond?: number; // 每秒处理样本数
  
  // 资源指标
  gpuUtilization?: number;   // GPU使用率
  memoryUsage?: number;      // 内存使用量
  
  // 模型指标
  precision?: number;
  recall?: number;
  f1Score?: number;
  
  // 自定义指标
  customMetrics?: { [key: string]: number };
  
  timestamp: string;
}

// 训练日志
export interface TrainingLog {
  timestamp: string;
  level: 'debug' | 'info' | 'warn' | 'error';
  message: string;
  source?: string;           // 日志来源
  podName?: string;          // Pod名称
  containerName?: string;    // 容器名称
}

// 检查点信息
export interface Checkpoint {
  id: string;
  jobId: string;
  epoch: number;
  step: number;
  filePath: string;
  fileSize: number;          // 文件大小(字节)
  metrics: TrainingMetrics;
  isBest: boolean;           // 是否为最佳模型
  createTime: string;
}

// 训练队列
export interface TrainingQueue {
  id: string;
  name: string;
  description?: string;
  workspaceId: string;
  workspaceName: string;
  
  // 资源配额
  resourceQuota: ResourceQuota;
  resourceUsed: ResourceQuota;
  
  // 调度配置
  schedulingPolicy: string;   // 调度策略
  maxRunningJobs: number;     // 最大并发任务数
  priority: number;           // 队列优先级
  
  // 访问控制
  allowedUsers: string[];     // 允许的用户ID列表
  allowedWorkspaces: string[]; // 允许的工作空间ID列表
  
  // 状态信息
  status: CommonStatus;
  jobCount: number;           // 任务总数
  runningJobCount: number;    // 运行中任务数
  queuedJobCount: number;     // 队列中任务数
  
  // 元数据
  labels?: Labels;
  annotations?: Annotations;
  createTime: string;
  updateTime: string;
}

// 训练任务
export interface TrainingJob {
  id: string;
  name: string;
  description?: string;
  
  // 所属信息
  workspaceId: string;
  workspaceName: string;
  projectId?: string;
  projectName?: string;
  queueId: string;
  queueName: string;
  
  // 用户信息
  creatorId: string;
  creatorName: string;
  
  // 任务配置
  framework: TrainingFramework;
  distributedType: DistributedType;
  image: string;              // 容器镜像
  command: string[];          // 启动命令
  args: string[];             // 启动参数
  workingDir?: string;        // 工作目录
  
  // 资源需求
  resourceRequirements: ResourceQuota;
  replicas: number;           // 副本数(分布式训练)
  
  // 数据配置
  datasetIds: string[];       // 数据集ID列表
  inputDataPath?: string;     // 输入数据路径
  outputDataPath?: string;    // 输出数据路径
  
  // 环境配置
  envVars: EnvVar[];
  volumeMounts: VolumeMount[];
  
  // 训练参数
  hyperParameters: { [key: string]: any }; // 超参数
  maxEpochs?: number;         // 最大训练轮数
  maxSteps?: number;          // 最大训练步数
  
  // 重试和恢复配置
  maxRetries: number;         // 最大重试次数
  currentRetries: number;     // 当前重试次数
  restartPolicy: 'Never' | 'OnFailure' | 'Always'; // 重启策略
  checkpointEnabled: boolean; // 是否启用检查点
  checkpointInterval?: number; // 检查点保存间隔(步数)
  
  // 状态信息
  status: TrainingJobStatus;
  priority: TrainingPriority;
  progress: number;           // 进度百分比
  
  // 时间信息
  createTime: string;
  updateTime: string;
  submitTime: string;         // 提交时间
  startTime?: string;         // 开始时间
  endTime?: string;           // 结束时间
  duration?: number;          // 运行时长(秒)
  estimatedDuration?: number; // 预估时长(秒)
  
  // 运行信息
  podStatus?: any;            // Pod状态信息
  nodeName?: string;          // 运行节点
  gpuIds?: string[];          // 分配的GPU ID列表
  
  // 监控信息
  currentMetrics?: TrainingMetrics;
  bestMetrics?: TrainingMetrics;
  checkpointCount: number;    // 检查点数量
  
  // 错误信息
  failureReason?: string;     // 失败原因
  errorMessage?: string;      // 错误消息
  
  // 元数据
  labels?: Labels;
  annotations?: Annotations;
}

// 创建训练任务请求
export interface CreateTrainingJobRequest {
  name: string;
  description?: string;
  workspaceId: string;
  projectId?: string;
  queueId: string;
  
  framework: TrainingFramework;
  distributedType: DistributedType;
  image: string;
  command: string[];
  args?: string[];
  workingDir?: string;
  
  resourceRequirements: ResourceQuota;
  replicas?: number;
  
  datasetIds?: string[];
  inputDataPath?: string;
  outputDataPath?: string;
  
  envVars?: EnvVar[];
  volumeMounts?: VolumeMount[];
  
  hyperParameters?: { [key: string]: any };
  maxEpochs?: number;
  maxSteps?: number;
  
  maxRetries?: number;
  restartPolicy?: string;
  checkpointEnabled?: boolean;
  checkpointInterval?: number;
  
  priority?: TrainingPriority;
  labels?: Labels;
}

// 更新训练任务请求
export interface UpdateTrainingJobRequest {
  id: string;
  description?: string;
  priority?: TrainingPriority;
  replicas?: number;
  labels?: Labels;
  annotations?: Annotations;
}

// 训练任务控制请求
export interface TrainingJobControlRequest {
  id: string;
  action: 'start' | 'stop' | 'pause' | 'resume' | 'cancel' | 'restart';
  reason?: string;
}

// 训练任务查询参数
export interface TrainingJobQuery extends PageQuery {
  workspaceId?: string;
  projectId?: string;
  queueId?: string;
  creatorId?: string;
  status?: TrainingJobStatus;
  priority?: TrainingPriority;
  framework?: TrainingFramework;
  distributedType?: DistributedType;
  createTimeStart?: string;
  createTimeEnd?: string;
  submitTimeStart?: string;
  submitTimeEnd?: string;
}

// 训练队列查询参数
export interface TrainingQueueQuery extends PageQuery {
  workspaceId?: string;
  status?: CommonStatus;
  allowedUserId?: string;     // 查询用户有权限的队列
}

// 训练统计信息
export interface TrainingStatistics extends BaseStatistics {
  byStatus: {
    [key in TrainingJobStatus]: number;
  };
  byFramework: {
    [key in TrainingFramework]: number;
  };
  byPriority: {
    [key in TrainingPriority]: number;
  };
  
  totalGpuHours: number;      // 总GPU小时数
  avgJobDuration: number;     // 平均任务时长(小时)
  successRate: number;        // 成功率(%)
  resourceUtilization: {
    avgGpuUtilization: number;
    avgCpuUtilization: number;
    avgMemoryUsage: number;
  };
}

// 训练模板
export interface TrainingTemplate {
  id: string;
  name: string;
  description?: string;
  category: string;           // 模板分类
  
  // 模板配置
  framework: TrainingFramework;
  distributedType: DistributedType;
  image: string;
  command: string[];
  args?: string[];
  
  // 默认资源需求
  defaultResources: ResourceQuota;
  defaultReplicas: number;
  
  // 默认参数
  defaultHyperParameters: { [key: string]: any };
  defaultEnvVars: EnvVar[];
  
  // 使用统计
  useCount: number;           // 使用次数
  
  // 创建信息
  creatorId: string;
  creatorName: string;
  isPublic: boolean;          // 是否公开
  
  // 元数据
  tags: string[];             // 标签
  labels?: Labels;
  createTime: string;
  updateTime: string;
}

// 实验跟踪
export interface Experiment {
  id: string;
  name: string;
  description?: string;
  projectId: string;
  projectName: string;
  
  // 实验配置
  framework: TrainingFramework;
  modelType: string;          // 模型类型
  datasetIds: string[];       // 使用的数据集
  
  // 超参数
  hyperParameters: { [key: string]: any };
  
  // 运行记录
  runs: ExperimentRun[];
  
  // 最佳结果
  bestRun?: ExperimentRun;
  bestMetrics?: TrainingMetrics;
  
  // 创建信息
  creatorId: string;
  creatorName: string;
  
  // 元数据
  tags: string[];
  labels?: Labels;
  createTime: string;
  updateTime: string;
}

// 实验运行
export interface ExperimentRun {
  id: string;
  experimentId: string;
  jobId?: string;             // 关联的训练任务ID
  runNumber: number;          // 运行编号
  
  // 配置
  hyperParameters: { [key: string]: any };
  
  // 结果
  status: TrainingJobStatus;
  metrics: TrainingMetrics[];  // 指标历史
  finalMetrics?: TrainingMetrics; // 最终指标
  
  // 模型产出
  modelPath?: string;         // 模型文件路径
  checkpoints: Checkpoint[];  // 检查点列表
  
  // 时间信息
  startTime: string;
  endTime?: string;
  duration?: number;
  
  // 运行环境
  image: string;
  gpuType?: string;
  gpuCount: number;
  
  // 错误信息
  errorMessage?: string;
  
  createTime: string;
}