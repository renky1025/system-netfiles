import api from './axios';

export interface SearchRequest {
    query: string;
    type?: 'file' | 'folder' | 'all';
    page?: number;
    page_size?: number;
}

export interface SearchResult {
    type: string;
    id: number;
    name: string;
    path: string;
    size?: number;
    created_at: string;
}

export interface SearchResponse {
    results: SearchResult[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export const search = (data: SearchRequest) => {
    return api.post<SearchResponse>('/api/search', data);
};

export const getSearchSuggestions = (query: string, limit: number = 5) => {
    return api.get<string[]>('/api/search/suggestions', {
        params: { q: query, limit }
    });
};
