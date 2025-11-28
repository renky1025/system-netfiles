/**
 * 文件夹相关API
 */
import api from './axios';
import type {
  ApiResponse,
  FolderInfo,
  FolderListResponse,
  FolderTreeResponse,
  FolderCreateRequest,
  FolderUpdateRequest,
  FolderMoveRequest,
  BreadcrumbResponse,
} from './types';

/**
 * 获取文件夹列表
 */
export const getFolderList = async (
  parentId: number | null
): Promise<ApiResponse<FolderListResponse>> => {
  const response = await api.get<ApiResponse<FolderListResponse>>('/api/folder/list', {
    params: { parent_id: parentId },
  });
  return response.data;
};

/**
 * 获取文件夹树
 */
export const getFolderTree = async (): Promise<ApiResponse<FolderTreeResponse>> => {
  const response = await api.get<ApiResponse<FolderTreeResponse>>('/api/folder/tree');
  return response.data;
};

/**
 * 创建文件夹
 */
export const createFolder = async (data: FolderCreateRequest): Promise<ApiResponse<FolderInfo>> => {
  const response = await api.post<ApiResponse<FolderInfo>>('/api/folder', data);
  return response.data;
};

/**
 * 更新文件夹（重命名）
 */
export const updateFolder = async (
  folderId: number,
  data: FolderUpdateRequest
): Promise<ApiResponse<void>> => {
  const response = await api.put<ApiResponse<void>>(`/api/folder/${folderId}`, data);
  return response.data;
};

/**
 * 移动文件夹
 */
export const moveFolder = async (
  folderId: number,
  data: FolderMoveRequest
): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>(`/api/folder/${folderId}/move`, data);
  return response.data;
};

/**
 * 删除文件夹
 */
export const deleteFolder = async (folderId: number): Promise<ApiResponse<void>> => {
  const response = await api.delete<ApiResponse<void>>(`/api/folder/${folderId}`);
  return response.data;
};

/**
 * 获取面包屑导航
 */
export const getBreadcrumb = async (
  folderId: number
): Promise<ApiResponse<BreadcrumbResponse>> => {
  const response = await api.get<ApiResponse<BreadcrumbResponse>>(
    `/api/folder/${folderId}/breadcrumb`
  );
  return response.data;
};

