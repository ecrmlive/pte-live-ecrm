import { http } from '@/api/http';

export interface CommunityPost {
  community_id: number;
  title: string;
  content: string;
  image?: string;
  category_id?: number;
  topic_id?: number;
  product_id?: number;
  nickname?: string;
  product_name?: string;
  status: number;
  count_reply?: number;
  refusal?: string;
}

export interface CommunityReply {
  reply_id: number;
  content: string;
  nickname?: string;
  create_time?: string;
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

export function createPost(data: Record<string, unknown>) {
  return http.post<CommunityPost>('/community/posts', data);
}

export function updatePost(id: number, data: Record<string, unknown>) {
  return http.put<CommunityPost>(`/community/posts/${id}`, data);
}

export function deletePost(id: number) {
  return http.delete<{ ok: boolean }>(`/community/posts/${id}`);
}

export function fetchReplies(id: number, params: Record<string, unknown> = {}) {
  return http.get<PageResult<CommunityReply>>(`/community/posts/${id}/replies`, { params });
}
