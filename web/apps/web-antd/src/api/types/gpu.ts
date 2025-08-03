import type { CommonStatus, ResourceQuota, Labels, Annotations, PageQuery, BaseStatistics } from './common';

/**
 * GPU资源管理类型定义
 */

// GPU状态枚举
export enum GPUStatus {
  AVAILABLE = 'available',   // 可用
  ALLOCATED = 'allocated',   // 已分配
  BUSY = 'busy',            // 忙碌
  MAINTENANCE = 'maintenance', // 维护中
  OFFLINE = 'offline',      // 离线
  ERROR = 'error',          // 故障
}

// GPU品牌枚举
export enum GPUBrand {
  NVIDIA = 'nvidia',
  AMD = 'amd',
  INTEL = 'intel',
}

// 集群类型枚举
export enum ClusterType {
  KUBERNETES = 'kubernetes',
  SLURM = 'slurm',
  CUSTOM = 'custom',
}

// GPU设备信息
export interface GPUDevice {
  id: string;
  name: string;              // 设备名称，如 "Tesla A100"
  nodeId: string;            // 所属节点ID
  nodeName: string;          // 所属节点名称
  clusterId: string;         // 所属集群ID
  clusterName: string;       // 所属集群名称
  brand: GPUBrand;           // 品牌
  model: string;             // 型号
  architecture: string;      // 架构，如 "Ampere"
  
  // 硬件规格
  cudaCores: number;         // CUDA核心数
  memorySize: number;        // 显存大小 (GB)
  memoryBandwidth: number;   // 显存带宽 (GB/s)
  baseClockRate: number;     // 基础时钟频率 (MHz)
  boostClockRate: number;    // 加速时钟频率 (MHz)
  maxPower: number;          // 最大功耗 (W)
  
  // 当前状态
  status: GPUStatus;
  temperature: number;       // 当前温度 (°C)
  powerUsage: number;        // 当前功耗 (W)
  memoryUsage: number;       // 显存使用率 (%)
  gpuUtilization: number;    // GPU使用率 (%)
  
  // 分配信息
  allocatedTo?: string;      // 分配给的任务ID
  allocatedUser?: string;    // 分配给的用户
  allocatedTime?: string;    // 分配时间
  
  // 驱动和软件信息
  driverVersion?: string;    // 驱动版本
  cudaVersion?: string;      // CUDA版本
  
  // 元数据
  labels?: Labels;
  annotations?: Annotations;
  createTime: string;
  updateTime: string;
  lastHeartbeat?: string;    // 最后心跳时间
}

// GPU节点信息
export interface GPUNode {
  id: string;
  name: string;
  clusterId: string;
  clusterName: string;
  
  // 节点基本信息
  hostname: string;
  ipAddress: string;
  osType: string;            // 操作系统类型
  osVersion: string;         // 操作系统版本
  kernelVersion: string;     // 内核版本
  
  // 硬件资源
  cpuCores: number;          // CPU核心数
  memoryTotal: number;       // 总内存 (GB)
  storageTotal: number;      // 总存储 (GB)
  gpuCount: number;          // GPU数量
  gpuTotal: ResourceQuota;   // GPU总资源
  gpuUsed: ResourceQuota;    // GPU已使用资源
  
  // 节点状态
  status: CommonStatus;
  cpuUsage: number;          // CPU使用率 (%)
  memoryUsage: number;       // 内存使用率 (%)
  storageUsage: number;      // 存储使用率 (%)
  networkRx: number;         // 网络接收 (MB/s)
  networkTx: number;         // 网络发送 (MB/s)
  
  // Kubernetes相关
  labels?: Labels;
  annotations?: Annotations;
  taints?: any[];            // 污点配置
  
  // 时间信息
  createTime: string;
  updateTime: string;
  lastHeartbeat?: string;
}

// GPU集群信息
export interface GPUCluster {
  id: string;
  name: string;
  description?: string;
  type: ClusterType;
  
  // 集群配置
  apiEndpoint: string;       // API端点
  region?: string;           // 地域
  zone?: string;             // 可用区
  
  // 资源统计
  nodeCount: number;         // 节点数量
  gpuCount: number;          // GPU总数
  totalResources: ResourceQuota; // 总资源
  usedResources: ResourceQuota;  // 已使用资源
  availableResources: ResourceQuota; // 可用资源
  
  // 集群状态
  status: CommonStatus;
  healthScore: number;       // 健康评分 (0-100)
  
  // 监控配置
  monitoringEnabled: boolean;
  alertingEnabled: boolean;
  
  // 元数据
  labels?: Labels;
  annotations?: Annotations;
  createTime: string;
  updateTime: string;
  lastHealthCheck?: string;
}

// 调度策略枚举
export enum SchedulingStrategy {
  ROUND_ROBIN = 'round_robin',      // 轮询
  BEST_FIT = 'best_fit',           // 最佳适配
  BALANCED = 'balanced',            // 负载均衡
  ANTI_AFFINITY = 'anti_affinity', // 反亲和性
  GPU_OPTIMIZED = 'gpu_optimized', // GPU使用率优化
  MEMORY_OPTIMIZED = 'memory_optimized', // 内存优化
}

// 资源调度请求
export interface ResourceScheduleRequest {
  strategy: SchedulingStrategy;
  requirements: ResourceQuota;
  constraints?: {
    nodeAffinity?: string[];      // 节点亲和性
    nodeAntiAffinity?: string[];  // 节点反亲和性
    gpuModel?: string;            // 指定GPU型号
    clusterId?: string;           // 指定集群
  };
  priority: number;              // 优先级 (1-10)
  timeout: number;               // 超时时间 (秒)
}

// 资源调度响应
export interface ResourceScheduleResponse {
  success: boolean;
  message: string;
  allocation?: {
    clusterId: string;
    nodeId: string;
    gpuIds: string[];
    estimatedStartTime: string;
    estimatedDuration: number;
  };
  reason?: string;               // 失败原因
  alternatives?: ResourceScheduleResponse[]; // 备选方案
}

// GPU查询参数
export interface GPUQuery extends PageQuery {
  clusterId?: string;
  nodeId?: string;
  brand?: GPUBrand;
  model?: string;
  status?: GPUStatus;
  allocatedTo?: string;
  minMemory?: number;
  maxMemory?: number;
  available?: boolean;           // 仅查询可用GPU
}

// 节点查询参数
export interface NodeQuery extends PageQuery {
  clusterId?: string;
  status?: CommonStatus;
  minGpuCount?: number;
  maxGpuCount?: number;
  osType?: string;
}

// 集群查询参数
export interface ClusterQuery extends PageQuery {
  type?: ClusterType;
  status?: CommonStatus;
  region?: string;
  minNodeCount?: number;
}

// GPU统计信息
export interface GPUStatistics extends BaseStatistics {
  byBrand: {
    [key in GPUBrand]: number;
  };
  byStatus: {
    [key in GPUStatus]: number;
  };
  byModel: { [model: string]: number };
  utilizationStats: {
    avgGpuUtilization: number;
    avgMemoryUsage: number;
    avgTemperature: number;
    avgPowerUsage: number;
  };
}

// 使用记录
export interface GPUUsageRecord {
  id: string;
  gpuId: string;
  gpuName: string;
  userId: string;
  userName: string;
  taskId: string;
  taskName: string;
  workspaceId: string;
  workspaceName: string;
  
  startTime: string;
  endTime?: string;
  duration: number;              // 使用时长 (分钟)
  
  cost: number;                  // 使用成本
  currency: string;              // 货币单位
  
  avgGpuUtilization: number;     // 平均GPU使用率
  avgMemoryUsage: number;        // 平均内存使用率
  peakMemoryUsage: number;       // 峰值内存使用率
  
  labels?: Labels;
  createTime: string;
}

// 计费配置
export interface BillingConfig {
  enabled: boolean;
  currency: string;              // 货币单位
  rates: {
    [gpuModel: string]: {
      hourlyRate: number;        // 小时费率
      dailyRate?: number;        // 日费率
      monthlyRate?: number;      // 月费率
    };
  };
  freeQuota?: {
    hoursPerMonth: number;       // 每月免费小时数
    applicableModels: string[];  // 适用的GPU型号
  };
}

// 成本统计
export interface CostStatistics {
  totalCost: number;
  currency: string;
  period: string;                // 统计周期
  
  costByUser: Array<{
    userId: string;
    userName: string;
    cost: number;
    duration: number;
  }>;
  
  costByGpuModel: Array<{
    model: string;
    cost: number;
    usage: number;
  }>;
  
  costByWorkspace: Array<{
    workspaceId: string;
    workspaceName: string;
    cost: number;
    duration: number;
  }>;
}