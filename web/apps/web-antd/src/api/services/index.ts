/**
 * API 服务统一导出文件
 * 提供所有 API 服务的统一入口
 */

// 导出服务类
export { BaseApiService } from './base';
export { AuthService, authService } from './auth';
export { TrainingService, trainingService } from './training';
export { GpuService, gpuService } from './gpu';
export { DatasetService, datasetService } from './dataset';

// 导出工具函数
export {
  handleApiError,
  isApiSuccess,
  getApiData,
  getApiMessage,
} from './base';

// 重新导出类型
export type { QueryParams } from './base';

// 服务注册表 - 便于统一管理和依赖注入
export const apiServices = {
  auth: authService,
  training: trainingService,
  gpu: gpuService,
  dataset: datasetService,
} as const;

// 服务类型定义
export type ApiServiceRegistry = typeof apiServices;
export type ApiServiceName = keyof ApiServiceRegistry;

/**
 * 获取指定的 API 服务实例
 * @param serviceName 服务名称
 */
export function getApiService<T extends ApiServiceName>(
  serviceName: T,
): ApiServiceRegistry[T] {
  return apiServices[serviceName];
}

/**
 * 批量健康检查所有服务
 */
export async function healthCheckAllServices(): Promise<{
  [K in ApiServiceName]: {
    status: 'healthy' | 'unhealthy';
    timestamp: number;
    error?: string;
  };
}> {
  const results = {} as any;
  
  for (const [name, service] of Object.entries(apiServices)) {
    try {
      await service.healthCheck();
      results[name] = {
        status: 'healthy',
        timestamp: Date.now(),
      };
    } catch (error: any) {
      results[name] = {
        status: 'unhealthy',
        timestamp: Date.now(),
        error: error.message || 'Unknown error',
      };
    }
  }
  
  return results;
}

// 默认导出主要服务
export default apiServices;