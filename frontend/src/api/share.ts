/**
 * 分享相关API
 */
import api from './axios';
import type {
  ApiResponse,
  ShareInfo,
  ShareCreateRequest,
  ShareCreateResponse,
  ShareListResponse,
  ShareValidateRequest,
  ShareAccessResponse,
} from './types';

/**
 * 创建分享
 */
export const createShare = async (
  data: ShareCreateRequest
): Promise<ApiResponse<ShareCreateResponse>> => {
  const response = await api.post<ApiResponse<ShareCreateResponse>>('/api/share/create', data);
  return response.data;
};

/**
 * 获取分享列表
 */
export const getShareList = async (
  page: number = 1,
  pageSize: number = 20
): Promise<ApiResponse<ShareListResponse>> => {
  const response = await api.get<ApiResponse<ShareListResponse>>('/api/share/list', {
    params: {
      page,
      page_size: pageSize,
    },
  });
  return response.data;
};

/**
 * 删除分享
 */
export const deleteShare = async (shareId: number): Promise<ApiResponse<void>> => {
  const response = await api.delete<ApiResponse<void>>(`/api/share/${shareId}`);
  return response.data;
};

/**
 * 验证分享密码
 */
export const validateShare = async (
  data: ShareValidateRequest
): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/api/share/validate', data);
  return response.data;
};

/**
 * 获取分享信息（公开访问）
 */
export const getShareInfo = async (code: string): Promise<ApiResponse<ShareAccessResponse>> => {
  const response = await api.get<ApiResponse<ShareAccessResponse>>(`/share/${code}`);
  return response.data;
};

/**
 * 下载分享文件（公开访问）
 */
export const downloadShareFile = async (code: string, password?: string): Promise<Blob> => {
  const response = await api.get(`/share/${code}/download`, {
    params: password ? { password } : undefined,
    responseType: 'blob',
  });
  return response.data;
};

