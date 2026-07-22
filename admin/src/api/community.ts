import { http } from '@/api/http';

export interface CommunityPost {
  community_id: number;
  title: string;
  content: string;
  nickname?: string;
  product_name?: string;
  product_id?: number;
  status: number;
  is_show: number;
  is_hot: number;
  count_reply?: number;
  mer_id?: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchPosts(params: Record<string, unknown>) {
  return http.get<PageResult<CommunityPost>>('/community/posts', { params });
}

export function auditPost(id: number, data: { status: number; refusal?: string; is_show?: number; is_hot?: number }) {
  return http.post<CommunityPost>(`/community/posts/${id}/audit`, data);
}

export function deletePost(id: number) {
  return http.delete<{ ok: boolean }>(`/community/posts/${id}`);
}
