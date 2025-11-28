import api from './axios';
import type { ApiResponse } from './types';

export interface OrgInfo {
    ID: number;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt?: string | null;
    Name: string;
    Type: string;
    ParentID?: number | null;
    ManagerID?: number | null;
    children?: OrgInfo[];
}

export interface OrgTreeResponse {
    roles: OrgInfo[];
}

export interface CreateOrgRequest {
    name: string;
    type: string;
    parent_id?: number | null;
    manager_id?: number | null;
}

export interface UpdateOrgRequest {
    name: string;
    type: string;
    manager_id?: number | null;
}

export const getOrgTree = async (): Promise<ApiResponse<OrgTreeResponse>> => {
    const response = await api.get<ApiResponse<OrgTreeResponse>>('/api/org/tree');
    return response.data;
};

export const createOrg = async (data: CreateOrgRequest): Promise<ApiResponse<OrgInfo>> => {
    const response = await api.post<ApiResponse<OrgInfo>>('/api/org', data);
    return response.data;
};

export const updateOrg = async (id: number, data: UpdateOrgRequest): Promise<ApiResponse<void>> => {
    const response = await api.put<ApiResponse<void>>(`/api/org/${id}`, data);
    return response.data;
};

export const deleteOrg = async (id: number): Promise<ApiResponse<void>> => {
    const response = await api.delete<ApiResponse<void>>(`/api/org/${id}`);
    return response.data;
};
