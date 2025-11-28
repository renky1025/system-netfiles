import api from './axios';
import type { ApiResponse } from './types';

export interface PermissionInfo {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  Name: string;
  Code: string;
  Description?: string;
}

export interface RoleInfo {
  id: number;
  created_at: string;
  updated_at: string;
  name: string;
  description?: string;
  permissions?: PermissionInfo[];
}

export interface RoleListResponse {
  roles: RoleInfo[];
}

export interface PermissionListResponse {
  permissions: PermissionInfo[];
}

export interface CreateRoleRequest {
  name: string;
  description?: string;
  permission_ids?: number[];
}

export interface UpdateRoleRequest {
  name?: string;
  description?: string;
  permission_ids?: number[];
}

export interface AssignRoleRequest {
  user_id: number;
  role_id: number;
}

// Get all roles
export const getRoles = async (): Promise<ApiResponse<RoleListResponse>> => {
  const response = await api.get<ApiResponse<RoleListResponse>>('/api/role/list');
  return response.data;
};

// Get role by ID
export const getRole = async (id: number): Promise<ApiResponse<{ role: RoleInfo }>> => {
  const response = await api.get<ApiResponse<{ role: RoleInfo }>>(`/api/role/${id}`);
  return response.data;
};

// Create role
export const createRole = async (data: CreateRoleRequest): Promise<ApiResponse<{ role: RoleInfo }>> => {
  const response = await api.post<ApiResponse<{ role: RoleInfo }>>('/api/role', data);
  return response.data;
};

// Update role
export const updateRole = async (id: number, data: UpdateRoleRequest): Promise<ApiResponse<void>> => {
  const response = await api.put<ApiResponse<void>>(`/api/role/${id}`, data);
  return response.data;
};

// Delete role
export const deleteRole = async (id: number): Promise<ApiResponse<void>> => {
  const response = await api.delete<ApiResponse<void>>(`/api/role/${id}`);
  return response.data;
};

// Get all permissions
export const getPermissions = async (): Promise<ApiResponse<PermissionListResponse>> => {
  const response = await api.get<ApiResponse<PermissionListResponse>>('/api/role/permissions');
  return response.data;
};

// Assign role to user
export const assignRoleToUser = async (data: AssignRoleRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/api/role/assign', data);
  return response.data;
};

// Remove role from user
export const removeRoleFromUser = async (data: AssignRoleRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/api/role/remove', data);
  return response.data;
};
