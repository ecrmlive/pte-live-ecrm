import { requestClient } from '#/api/request';

export interface CommunityPost {
  cate_name?: string;
  category_id: number;
  community_id: number;
  content: string;
  count_reply: number;
  create_time: string;
  image?: string;
  is_hot: number;
  is_show: number;
  mer_id: number;
  nickname?: string;
  product_id: number;
  product_name?: string;
  product_price?: number;
  pv: number;
  refusal?: string;
  status: number;
  status_time?: string;
  title: string;
  topic_name?: string;
}

export interface CommunityReply {
  community_id: number;
  content: string;
  create_time: string;
  nickname?: string;
  reply_id: number;
  uid: number;
}

export interface CommunityPage<T> {
  limit: number;
  list: T[];
  page: number;
  total: number;
}

export interface CommunityAuditPayload {
  is_hot?: 0 | 1;
  is_show?: 0 | 1;
  refusal?: string;
  status: -1 | 0 | 1;
}

export function listCommunityPostsApi(params: { keyword?: string; limit: number; page: number; status?: number }) {
  return requestClient.get<CommunityPage<CommunityPost>>('/community/posts', { params });
}

export function getCommunityPostApi(id: number) {
  return requestClient.get<CommunityPost>(`/community/posts/${id}`);
}

export function listCommunityRepliesApi(id: number, params: { limit: number; page: number }) {
  return requestClient.get<CommunityPage<CommunityReply>>(`/community/posts/${id}/replies`, { params });
}

export function auditCommunityPostApi(id: number, data: CommunityAuditPayload) {
  return requestClient.post<CommunityPost>(`/community/posts/${id}/audit`, data);
}

export function deleteCommunityPostApi(id: number) {
  return requestClient.delete(`/community/posts/${id}`);
}

export function deleteCommunityReplyApi(id: number) {
  return requestClient.delete(`/community/replies/${id}`);
}
