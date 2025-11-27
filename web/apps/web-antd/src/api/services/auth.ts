/**
 * 认证服务 API
 * 处理用户登录、认证、权限相关的接口调用
 */

import type {
  AuthApi,
  UserInfo,
  ApiResponse,
} from '#/types/api';
import { BaseApiService } from './base';
import { baseRequestClient } from '#/api/request';

export class AuthService extends BaseApiService {
  constructor() {
    super('/api/v1/auth');
  }

  /**
   * 用户登录
   * @param loginData 登录信息
   */
  async login(loginData: AuthApi.LoginRequest): Promise<AuthApi.LoginResult> {
    // 使用不需要认证的客户端进行登录请求
    return baseRequestClient.post<AuthApi.LoginResult>(
      '/api/v1/auth/login',
      loginData,
    );
  }

  /**
   * 刷新访问令牌
   * @param refreshToken 刷新令牌
   */
  async refreshToken(
    refreshToken: string,
  ): Promise<AuthApi.RefreshTokenResult> {
    return baseRequestClient.post<AuthApi.RefreshTokenResult>(
      '/api/v1/auth/refresh',
      { refreshToken },
    );
  }

  /**
   * 用户登出
   */
  async logout(): Promise<ApiResponse<void>> {
    return this.post('/logout');
  }

  /**
   * 获取用户权限代码
   */
  async getAccessCodes(): Promise<string[]> {
    const response = await this.get<{ codes: string[] }>('/codes');
    return response.codes;
  }

  /**
   * 获取用户信息
   */
  async getUserInfo(): Promise<UserInfo> {
    const response = await this.get<{ userInfo: UserInfo }>('/api/v1/user/info');
    return response.userInfo;
  }

  /**
   * 修改密码
   * @param oldPassword 旧密码
   * @param newPassword 新密码
   */
  async changePassword(
    oldPassword: string,
    newPassword: string,
  ): Promise<ApiResponse<void>> {
    return this.post('/change-password', {
      oldPassword,
      newPassword,
    });
  }

  /**
   * 更新用户信息
   * @param userInfo 用户信息
   */
  async updateUserInfo(
    userInfo: Partial<UserInfo>,
  ): Promise<ApiResponse<UserInfo>> {
    return this.put('/user/profile', userInfo);
  }

  /**
   * 上传用户头像
   * @param file 头像文件
   */
  async uploadAvatar(file: File): Promise<ApiResponse<{ avatarUrl: string }>> {
    return this.upload('/user/avatar', file);
  }

  /**
   * 验证令牌有效性
   */
  async validateToken(): Promise<boolean> {
    try {
      await this.get('/validate');
      return true;
    } catch {
      return false;
    }
  }

  /**
   * 获取用户权限列表
   */
  async getUserPermissions(): Promise<string[]> {
    const response = await this.get<{ permissions: string[] }>('/permissions');
    return response.permissions;
  }

  /**
   * 检查用户是否有特定权限
   * @param permission 权限代码
   */
  async hasPermission(permission: string): Promise<boolean> {
    const response = await this.get<{ hasPermission: boolean }>(
      `/permissions/check`,
      { permission },
    );
    return response.hasPermission;
  }

  /**
   * 获取用户角色列表
   */
  async getUserRoles(): Promise<string[]> {
    const response = await this.get<{ roles: string[] }>('/roles');
    return response.roles;
  }
}

// 导出服务实例
export const authService = new AuthService();