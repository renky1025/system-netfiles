/**
 * 回收站相关API
 */
import api from './axios';
import type {
  ApiResponse,
  RecycleItem,
  RecycleListResponse,
} from './types';

/**
 * 获取回收站列表
 */
export const getRecycleList = async (
  page: number = 1,
  pageSize: number = 20
): Promise<ApiResponse<RecycleListResponse>> => {
  const response = await api.get<ApiResponse<RecycleListResponse>>('/api/recycle/list', {
    params: {
      page,
      page_size: pageSize,
    },
  });
  return response.data;
};

/**
 * 恢复文件
 */
export const restoreFile = async (fileId: number): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>(`/api/recycle/${fileId}/restore`);
  return response.data;
};

/**
 * 永久删除文件
 */
export const permanentDeleteFile = async (fileId: number): Promise<ApiResponse<void>> => {
  const response = await api.delete<ApiResponse<void>>(`/api/recycle/${fileId}`);
  return response.data;
};

/**
 * 清空回收站
 */
export const clearRecycleBin = async (): Promise<ApiResponse<void>> => {
  const response = await api.delete<ApiResponse<void>>('/api/recycle/clear');
  return response.data;
};

