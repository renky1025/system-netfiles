/**
 * 统一的API响应类型定义
 */

/**
 * 标准API响应结构
 */
export interface ApiResponse<T = any> {
  code: number;
  msg: string;
  data: T;
  error?: string;
}

/**
 * 分页响应数据
 */
export interface PaginatedData<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

/**
 * 认证相关类型
 */
export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user?: {
    id: number;
    username: string;
    email: string;
    roles?: RoleInfo[];
  };
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export interface ResetPasswordRequestRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  new_password: string;
}

/**
 * 文件相关类型
 */
export interface FileInfo {
  ID?: number;
  id?: number;
  Name?: string;
  name?: string;
  Size?: number;
  size?: number;
  UpdatedAt?: string;
  updated_at?: string;
  CreatedAt?: string;
  created_at?: string;
  FolderID?: number;
  folder_id?: number;
  MD5?: string;
  md5?: string;
}

export interface FileListResponse {
  files: FileInfo[];
}

export interface FileCheckRequest {
  md5: string;
}

export interface FileCheckResponse {
  exists: boolean;
}

export interface FileInstantUploadRequest {
  md5: string;
  file_name: string;
  file_size: number;
  folder_id: number | null;
}

export interface FileUploadChunkRequest {
  upload_id: string;
  index: number;
  file: Blob;
}

export interface FileUploadMergeRequest {
  upload_id: string;
  file_name: string;
  total_chunks: number;
  folder_id: number | null;
}

export interface FileRenameRequest {
  new_name: string;
}

export interface FileMoveRequest {
  folder_id: number;
}

export interface FileCopyRequest {
  folder_id: number;
  new_name?: string;
}

export interface FileBatchDeleteRequest {
  file_ids: number[];
}

export interface FileBatchMoveRequest {
  file_ids: number[];
  folder_id: number;
}

export interface FileBatchCopyRequest {
  file_ids: number[];
  folder_id: number;
}

/**
 * 文件夹相关类型
 */
export interface FolderInfo {
  ID?: number;
  id?: number;
  Name?: string;
  name?: string;
  ParentID?: number;
  parent_id?: number;
  UpdatedAt?: string;
  updated_at?: string;
  CreatedAt?: string;
  created_at?: string;
  children?: FolderInfo[];
}

export interface FolderListResponse {
  folders: FolderInfo[];
}

export interface FolderTreeResponse {
  tree: FolderInfo[];
}

export interface FolderCreateRequest {
  name: string;
  parent_id: number | null;
}

export interface FolderUpdateRequest {
  name: string;
}

export interface FolderMoveRequest {
  parent_id: number;
}

export interface BreadcrumbItem {
  id: number;
  name: string;
  path: string;
}

export interface BreadcrumbResponse {
  breadcrumbs: BreadcrumbItem[];
}

/**
 * 分享相关类型
 */
export interface ShareInfo {
  ID?: number;
  id?: number;
  Code?: string;
  code?: string;
  Type?: number;
  type?: number;
  ClickCount?: number;
  click_count?: number;
  ExpiredAt?: string;
  expired_at?: string;
  CreatedAt?: string;
  created_at?: string;
  file_name?: string;
  file_size?: number;
}

export interface ShareCreateRequest {
  file_id: number;
  password?: string;
  expires_at?: string;
  max_downloads?: number;
}

export interface ShareCreateResponse {
  share_code: string;
}

export interface ShareListResponse extends PaginatedData<ShareInfo> { }

export interface ShareValidateRequest {
  share_id: number;
  password: string;
}

export interface ShareAccessResponse {
  share: ShareInfo;
}

/**
 * 回收站相关类型
 */
export interface RecycleItem {
  ID?: number;
  id?: number;
  Name?: string;
  name?: string;
  Size?: number;
  size?: number;
  DeletedAt?: string;
  deleted_at?: string;
}

export interface RecycleListResponse extends PaginatedData<RecycleItem> { }

/**
 * 管理后台相关类型
 */

// 用户管理
export interface AdminUserListResponse extends PaginatedData<UserInfo> { }

export interface UserInfo {
  id: number;
  username: string;
  email: string;
  phone?: string;
  status: number; // 1: active, 2: frozen
  created_at: string;
  last_login_at?: string;
  roles?: RoleInfo[];
}

export interface RoleInfo {
  id: number;
  name: string;
  description?: string;
}

// 系统统计
export interface SystemStatsResponse {
    active_shares: number;
    active_users: number;
    recycle_bin_files: number;
    today_downloads: number;
    today_uploads: number;
    total_files: number;
    total_storage: number;
    total_users: number;

}

export interface StorageStatsResponse {
  total_space: number;
  used_space: number;
  free_space: number;
  usage_percent: number;
}

// 审计日志
export interface FileOpLog {
  id: number;
  user_id: number;
  username: string;
  file_id: number;
  file_name: string;
  operation: string;
  ip: string;
  created_at: string;
  details?: string;
}

export interface LoginLog {
  id: number;
  user_id: number;
  username: string;
  ip: string;
  file: Blob;
}

export interface FileUploadMergeRequest {
  upload_id: string;
  file_name: string;
  total_chunks: number;
  folder_id: number | null;
}

export interface FileRenameRequest {
  new_name: string;
}

export interface FileMoveRequest {
  folder_id: number;
}

export interface FileCopyRequest {
  folder_id: number;
  new_name?: string;
}

export interface FileBatchDeleteRequest {
  file_ids: number[];
}

export interface FileBatchMoveRequest {
  file_ids: number[];
  folder_id: number;
}

/**
 * 文件夹相关类型
 */
export interface FolderInfo {
  ID?: number;
  id?: number;
  Name?: string;
  name?: string;
  ParentID?: number;
  parent_id?: number;
  UpdatedAt?: string;
  updated_at?: string;
  CreatedAt?: string;
  created_at?: string;
  children?: FolderInfo[];
}

export interface FolderListResponse {
  folders: FolderInfo[];
}

export interface FolderTreeResponse {
  tree: FolderInfo[];
}

export interface FolderCreateRequest {
  name: string;
  parent_id: number | null;
}

export interface FolderUpdateRequest {
  name: string;
}

export interface FolderMoveRequest {
  parent_id: number;
}

export interface BreadcrumbItem {
  id: number;
  name: string;
  path: string;
}

export interface BreadcrumbResponse {
  breadcrumbs: BreadcrumbItem[];
}

/**
 * 分享相关类型
 */
export interface ShareInfo {
  ID?: number;
  id?: number;
  Code?: string;
  code?: string;
  Type?: number;
  type?: number;
  ClickCount?: number;
  click_count?: number;
  ExpiredAt?: string;
  expired_at?: string;
  CreatedAt?: string;
  created_at?: string;
  file_name?: string;
  file_size?: number;
}

export interface ShareCreateRequest {
  file_id: number;
  password?: string;
  expires_at?: string;
  max_downloads?: number;
}

export interface ShareCreateResponse {
  share_code: string;
}

export interface ShareListResponse extends PaginatedData<ShareInfo> { }

export interface ShareValidateRequest {
  share_id: number;
  password: string;
}

export interface ShareAccessResponse {
  share: ShareInfo;
}

/**
 * 回收站相关类型
 */
export interface RecycleItem {
  ID?: number;
  id?: number;
  Name?: string;
  name?: string;
  Size?: number;
  size?: number;
  DeletedAt?: string;
  deleted_at?: string;
}

export interface RecycleListResponse extends PaginatedData<RecycleItem> { }

/**
 * 管理后台相关类型
 */

// 用户管理
export interface AdminUserListResponse extends PaginatedData<UserInfo> { }

export interface UserInfo {
  id: number;
  username: string;
  email: string;
  phone?: string;
  status: number; // 1: active, 2: frozen
  created_at: string;
  last_login_at?: string;
  roles?: RoleInfo[];
}

export interface RoleInfo {
  id: number;
  name: string;
  description?: string;
}

// 系统统计
export interface SystemStatsResponse {
  user_count: number;
  file_count: number;
  total_storage: number; // bytes
  active_users_today: number;
}

export interface StorageStatsResponse {
  total_space: number;
  used_space: number;
  free_space: number;
  usage_percent: number;
}

// 审计日志
export interface FileOpLog {
  id: number;
  user_id: number;
  username: string;
  file_id: number;
  file_name: string;
  operation: string;
  ip: string;
  created_at: string;
  details?: string;
}

export interface LoginLog {
  id: number;
  user_id: number;
  username: string;
  ip: string;
  status: number; // 1: success, 0: failed
  created_at: string;
  user_agent?: string;
}

export interface FileOpLogListResponse extends PaginatedData<FileOpLog> { }
export interface LoginLogListResponse extends PaginatedData<LoginLog> { }

/**
 * 版本管理相关类型
 */
export interface FileVersion {
  id: number;
  file_id: number;
  version: number;
  size: number;
  md5: string;
  created_at: string;
  creator_name?: string;
}

export interface VersionListResponse {
  versions: FileVersion[];
}
