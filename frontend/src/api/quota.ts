import api from './axios';
import type { ApiResponse } from './types';

/**
 * 配额信息响应
 */
export interface QuotaInfo {
    user_id: number;
    storage_quota: number;
    used_storage: number;
    free_storage: number;
    usage_percent: number;
    quota_source: string;
    rate_limit_source: string;
    download_rate_limit: number;
}

/**
 * 获取当前用户配额
 */
export const getMyQuota = async (): Promise<ApiResponse<QuotaInfo>> => {
    const response = await api.get<ApiResponse<QuotaInfo>>('/api/quota');
    return response.data;
};

/**
 * 获取指定用户配额 (管理员)
 */
export const getUserQuota = async (userId: number): Promise<ApiResponse<QuotaInfo>> => {
    const response = await api.get<ApiResponse<QuotaInfo>>(`/api/admin/users/${userId}/quota`);
    return response.data;
};

/**
 * 设置用户个人配额 (管理员)
 */
export const setUserQuota = async (userId: number, quota: number): Promise<ApiResponse<void>> => {
    const response = await api.put<ApiResponse<void>>(`/api/admin/users/${userId}/quota`, { quota });
    return response.data;
};

/**
 * 重新计算用户存储使用量 (管理员)
 */
export const recalculateUserStorage = async (userId: number): Promise<ApiResponse<{ quota: QuotaInfo }>> => {
    const response = await api.post<ApiResponse<{ quota: QuotaInfo }>>(`/api/admin/users/${userId}/recalculate-storage`);
    return response.data;
};

/**
 * 设置角色配额 (管理员)
 */
export const setRoleQuota = async (roleId: number, quota: number, rateLimit: number): Promise<ApiResponse<void>> => {
    const response = await api.put<ApiResponse<void>>(`/api/admin/roles/${roleId}/quota`, {
        quota,
        rate_limit: rateLimit
    });
    return response.data;
};

/**
 * 设置部门配额 (管理员)
 */
export const setOrgQuota = async (orgId: number, quota: number, rateLimit: number): Promise<ApiResponse<void>> => {
    const response = await api.put<ApiResponse<void>>(`/api/admin/orgs/${orgId}/quota`, {
        quota,
        rate_limit: rateLimit
    });
    return response.data;
};
