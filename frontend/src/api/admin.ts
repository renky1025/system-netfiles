/**
 * 管理后台相关API
 */
import api from './axios';
import type {
    ApiResponse,
    AdminUserListResponse,
    UserInfo,
    SystemStatsResponse,
    StorageStatsResponse,
    FileOpLogListResponse,
    LoginLogListResponse,
    FileListResponse,
    ShareListResponse,
} from './types';

/**
 * 获取用户列表
 */
export const getUserList = async (page: number, pageSize: number, search?: string): Promise<ApiResponse<AdminUserListResponse>> => {
    const response = await api.get<ApiResponse<AdminUserListResponse>>('/api/admin/users', {
        params: { page, page_size: pageSize, search },
    });
    return response.data;
};

/**
 * 创建用户
 */
export const createUser = async (data: { username: string; email: string; password: string }): Promise<ApiResponse<{ user: UserInfo }>> => {
    const response = await api.post<ApiResponse<{ user: UserInfo }>>('/api/admin/users', data);
    return response.data;
};

/**
 * 获取单个用户详情
 */
export const getUser = async (id: number): Promise<ApiResponse<UserInfo>> => {
    const response = await api.get<ApiResponse<UserInfo>>(`/api/admin/users/${id}`);
    return response.data;
};

/**
 * 更新用户信息
 */
export const updateUser = async (id: number, data: { username?: string; email?: string }): Promise<ApiResponse<void>> => {
    const response = await api.put<ApiResponse<void>>(`/api/admin/users/${id}`, data);
    return response.data;
};

/**
 * 冻结用户
 */
export const freezeUser = async (id: number): Promise<ApiResponse<void>> => {
    const response = await api.post<ApiResponse<void>>(`/api/admin/users/${id}/freeze`);
    return response.data;
};

/**
 * 解冻用户
 */
export const unfreezeUser = async (id: number): Promise<ApiResponse<void>> => {
    const response = await api.post<ApiResponse<void>>(`/api/admin/users/${id}/unfreeze`);
    return response.data;
};

/**
 * 管理员重置用户密码
 */
export const adminResetPassword = async (id: number, newPassword?: string): Promise<ApiResponse<{ new_password: string }>> => {
    const password = newPassword || Math.random().toString(36).slice(-8) + 'A1!';
    const response = await api.post<ApiResponse<{ new_password: string }>>(`/api/admin/users/${id}/reset-password`, {
        new_password: password,
    });
    // Return the password we sent
    return { ...response.data, data: { new_password: password } };
};

/**
 * 删除用户
 */
export const deleteUser = async (id: number): Promise<ApiResponse<void>> => {
    const response = await api.delete<ApiResponse<void>>(`/api/admin/users/${id}`);
    return response.data;
};

/**
 * 获取全局文件列表
 */
export const getAllFiles = async (page: number, pageSize: number): Promise<ApiResponse<FileListResponse>> => {
    const response = await api.get<ApiResponse<FileListResponse>>('/api/admin/files', {
        params: { page, page_size: pageSize },
    });
    return response.data;
};

/**
 * 强制删除文件
 */
export const forceDeleteFile = async (id: number): Promise<ApiResponse<void>> => {
    const response = await api.delete<ApiResponse<void>>(`/api/admin/files/${id}/force`);
    return response.data;
};

/**
 * 获取全局分享列表
 */
export const getAllShares = async (page: number, pageSize: number): Promise<ApiResponse<ShareListResponse>> => {
    const response = await api.get<ApiResponse<ShareListResponse>>('/api/admin/shares', {
        params: { page, page_size: pageSize },
    });
    return response.data;
};

/**
 * 禁用分享
 */
export const disableShare = async (id: number): Promise<ApiResponse<void>> => {
    const response = await api.post<ApiResponse<void>>(`/api/admin/shares/${id}/disable`);
    return response.data;
};

/**
 * 获取系统统计
 */
export const getSystemStats = async (): Promise<ApiResponse<SystemStatsResponse>> => {
    const response = await api.get<ApiResponse<SystemStatsResponse>>('/api/admin/stats/system');
    return response.data;
};

/**
 * 获取存储统计
 */
export const getStorageStats = async (): Promise<ApiResponse<StorageStatsResponse>> => {
    const response = await api.get<ApiResponse<StorageStatsResponse>>('/api/admin/stats/storage');
    return response.data;
};

/**
 * 获取文件操作日志
 */
export const getFileOpLogs = async (page: number, pageSize: number): Promise<ApiResponse<FileOpLogListResponse>> => {
    const response = await api.get<ApiResponse<FileOpLogListResponse>>('/api/admin/logs/file-ops', {
        params: { page, page_size: pageSize },
    });
    return response.data;
};

/**
 * 获取登录日志
 */
export const getLoginLogs = async (page: number, pageSize: number): Promise<ApiResponse<LoginLogListResponse>> => {
    const response = await api.get<ApiResponse<LoginLogListResponse>>('/api/admin/logs/login', {
        params: { page, page_size: pageSize },
    });
    return response.data;
};
