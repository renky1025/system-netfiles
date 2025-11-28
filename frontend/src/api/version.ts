/**
 * 版本管理相关API
 */
import api from './axios';
import type {
    ApiResponse,
    VersionListResponse,
} from './types';

/**
 * 获取文件版本列表
 */
export const getFileVersions = async (fileId: number): Promise<ApiResponse<VersionListResponse>> => {
    const response = await api.get<ApiResponse<VersionListResponse>>(`/api/file/${fileId}/versions`);
    return response.data;
};

/**
 * 回滚版本
 */
export const rollbackVersion = async (fileId: number, versionId: number): Promise<ApiResponse<void>> => {
    const response = await api.post<ApiResponse<void>>(`/api/file/${fileId}/versions/${versionId}/rollback`);
    return response.data;
};

/**
 * 删除版本
 */
export const deleteVersion = async (fileId: number, versionId: number): Promise<ApiResponse<void>> => {
    const response = await api.delete<ApiResponse<void>>(`/api/file/${fileId}/versions/${versionId}`);
    return response.data;
};
