import { requestClient } from '#/api/request';
import type {
  Workspace,
  WorkspaceQuery,
  CreateWorkspaceRequest,
  UpdateWorkspaceRequest,
  WorkspaceStatistics,
  WorkspaceMember,
  InviteMemberRequest,
  UpdateMemberPermissionRequest,
  WorkspaceProject,
  CreateProjectRequest,
  ProjectQuery,
  PageResponse,
  CrudResponse,
} from '#/api/types';

/**
 * 工作空间管理API接口
 */

// 工作空间管理
export const getWorkspaceList = (params: WorkspaceQuery) => {
  return requestClient.get<PageResponse<Workspace>>('/api/v1/workspaces', { params });
};

export const getWorkspaceById = (id: string) => {
  return requestClient.get<Workspace>(`/api/v1/workspaces/${id}`);
};

export const createWorkspace = (data: CreateWorkspaceRequest) => {
  return requestClient.post<CrudResponse>('/api/v1/workspaces', data);
};

export const updateWorkspace = (data: UpdateWorkspaceRequest) => {
  return requestClient.put<CrudResponse>(`/api/v1/workspaces/${data.id}`, data);
};

export const deleteWorkspace = (id: string) => {
  return requestClient.delete<CrudResponse>(`/api/v1/workspaces/${id}`);
};

export const getWorkspaceStatistics = () => {
  return requestClient.get<WorkspaceStatistics>('/api/v1/workspaces/statistics');
};

// 工作空间成员管理
export const getWorkspaceMembers = (workspaceId: string, params?: { page?: number; pageSize?: number }) => {
  return requestClient.get<PageResponse<WorkspaceMember>>(`/api/v1/workspaces/${workspaceId}/members`, { params });
};

export const inviteMembers = (data: InviteMemberRequest) => {
  return requestClient.post<CrudResponse>(`/api/v1/workspaces/${data.workspaceId}/members/invite`, data);
};

export const updateMemberPermission = (data: UpdateMemberPermissionRequest) => {
  return requestClient.put<CrudResponse>(
    `/api/v1/workspaces/${data.workspaceId}/members/${data.userId}/permission`,
    { permission: data.permission }
  );
};

export const removeMember = (workspaceId: string, userId: string) => {
  return requestClient.delete<CrudResponse>(`/api/v1/workspaces/${workspaceId}/members/${userId}`);
};

// 工作空间项目管理
export const getWorkspaceProjects = (workspaceId: string, params?: ProjectQuery) => {
  return requestClient.get<PageResponse<WorkspaceProject>>(`/api/v1/workspaces/${workspaceId}/projects`, { params });
};

export const createProject = (data: CreateProjectRequest) => {
  return requestClient.post<CrudResponse>(`/api/v1/workspaces/${data.workspaceId}/projects`, data);
};

export const getProjectById = (workspaceId: string, projectId: string) => {
  return requestClient.get<WorkspaceProject>(`/api/v1/workspaces/${workspaceId}/projects/${projectId}`);
};

export const updateProject = (workspaceId: string, projectId: string, data: Partial<CreateProjectRequest>) => {
  return requestClient.put<CrudResponse>(`/api/v1/workspaces/${workspaceId}/projects/${projectId}`, data);
};

export const deleteProject = (workspaceId: string, projectId: string) => {
  return requestClient.delete<CrudResponse>(`/api/v1/workspaces/${workspaceId}/projects/${projectId}`);
};

// 我参与的工作空间
export const getMyWorkspaces = (params?: { status?: string; type?: string }) => {
  return requestClient.get<PageResponse<Workspace>>('/api/v1/workspaces/my', { params });
};

// 获取工作空间资源使用情况
export const getWorkspaceResourceUsage = (workspaceId: string) => {
  return requestClient.get<any>(`/api/v1/workspaces/${workspaceId}/resources/usage`);
};

// 获取工作空间配额信息
export const getWorkspaceQuota = (workspaceId: string) => {
  return requestClient.get<any>(`/api/v1/workspaces/${workspaceId}/quota`);
};

// 更新工作空间配额
export const updateWorkspaceQuota = (workspaceId: string, quota: any) => {
  return requestClient.put<CrudResponse>(`/api/v1/workspaces/${workspaceId}/quota`, quota);
};