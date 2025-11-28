/**
 * 文件相关API
 */
import api from './axios';
import type {
  ApiResponse,
  FileInfo,
  FileListResponse,
  FileCheckRequest,
  FileCheckResponse,
  FileInstantUploadRequest,
  FileUploadChunkRequest,
  FileUploadMergeRequest,
  FileRenameRequest,
  FileMoveRequest,
  FileCopyRequest,
  FileBatchDeleteRequest,
  FileBatchMoveRequest,
  FileBatchCopyRequest,
} from './types';

/**
 * 获取文件列表
 */
export const getFileList = async (
  folderId: number | null
): Promise<ApiResponse<FileListResponse>> => {
  const response = await api.get<ApiResponse<FileListResponse>>('/api/file/list', {
    params: { folder_id: folderId },
  });
  return response.data;
};

/**
 * 检查文件是否存在（通过MD5）
 */
export const checkFile = async (data: FileCheckRequest): Promise<ApiResponse<FileCheckResponse>> => {
  const response = await api.post<ApiResponse<FileCheckResponse>>('/api/file/check', data);
  return response.data;
};

/**
 * 秒传文件
 */
export const instantUpload = async (
  data: FileInstantUploadRequest
): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/api/file/instant-upload', data);
  return response.data;
};

/**
 * 上传文件分片
 */
export const uploadChunk = async (data: FileUploadChunkRequest): Promise<ApiResponse<void>> => {
  const formData = new FormData();
  formData.append('upload_id', data.upload_id);
  formData.append('index', data.index.toString());
  formData.append('file', data.file);

  const response = await api.post<ApiResponse<void>>('/api/file/upload/chunk', formData);
  return response.data;
};

/**
 * 合并文件分片
 */
export const mergeChunks = async (data: FileUploadMergeRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/api/file/upload/merge', data);
  return response.data;
};

/**
 * 下载文件
 */
export const downloadFile = async (fileId: number): Promise<Blob> => {
  const response = await api.get(`/api/file/${fileId}/download`, {
    responseType: 'blob',
  });
  return response.data;
};

/**
 * 重命名文件
 */
export const renameFile = async (
  fileId: number,
  data: FileRenameRequest
): Promise<ApiResponse<void>> => {
  const response = await api.put<ApiResponse<void>>(`/api/file/${fileId}/rename`, data);
  return response.data;
};

/**
 * 移动文件
 */
export const moveFile = async (fileId: number, data: FileMoveRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>(`/api/file/${fileId}/move`, data);
  return response.data;
};

/**
 * 复制文件
 */
export const copyFile = async (fileId: number, data: FileCopyRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>(`/api/file/${fileId}/copy`, data);
  return response.data;
};

/**
 * 删除文件
 */
export const deleteFile = async (fileId: number): Promise<ApiResponse<void>> => {
  const response = await api.delete<ApiResponse<void>>(`/api/file/${fileId}`);
  return response.data;
};

/**
 * 批量删除文件
 */
export const batchDeleteFiles = async (
  data: FileBatchDeleteRequest
): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/api/file/batch/delete', data);
  return response.data;
};

/**
 * 批量移动文件
 */
export const batchMoveFiles = async (data: FileBatchMoveRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/api/file/batch/move', data);
  return response.data;
};

/**
 * 批量复制文件
 */
export const batchCopyFiles = async (data: FileBatchCopyRequest): Promise<ApiResponse<void>> => {
  const response = await api.post<ApiResponse<void>>('/api/file/batch/copy', data);
  return response.data;
};

