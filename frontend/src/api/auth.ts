/**
 * 认证相关API
 */
import api from './axios';
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  ChangePasswordRequest,
  ResetPasswordRequestRequest,
  ResetPasswordRequest,
} from './types';

/**
 * 用户登录
 */
export const login = async (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
  const response = await api.post<ApiResponse<LoginResponse>>('/auth/login', data);
  return response.data;
};

/**
 * 用户注册
 */
export const register = async (data: RegisterRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/auth/register', data);
  return response.data;
};

/**
 * 修改密码
 */
export const changePassword = async (data: ChangePasswordRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/auth/password/change', data);
  return response.data;
};

/**
 * 请求重置密码
 */
export const requestResetPassword = async (
  data: ResetPasswordRequestRequest
): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/auth/password/reset-request', data);
  return response.data;
};

/**
 * 重置密码
 */
export const resetPassword = async (data: ResetPasswordRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/auth/password/reset', data);
  return response.data;
};

/**
 * 刷新Token
 */
export const refreshToken = async (): Promise<ApiResponse<{ token: string }>> => {
  const response = await api.post<ApiResponse<{ token: string }>>('/api/auth/refresh', {});
  return response.data;
};

