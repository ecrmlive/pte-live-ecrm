import { requestClient } from '#/api/request';

export interface MerchantCommunityCategory {
  cate_name: string;
  category_id: number;
}

export interface MerchantCommunityTopic {
  category_id: number;
  is_hot: number;
  topic_id: number;
  topic_name: string;
}

export interface MerchantCommunityPost {
  cate_name?: string;
  category_id: number;
  community_id: number;
  content: string;
  count_reply: number;
  create_time: string;
  image?: string;
  is_show: number;
  mer_id: number;
  nickname?: string;
  product_id: number;
  product_name?: string;
  refusal?: string;
  status: number;
  title: string;
  topic_id: number;
  topic_name?: string;
}

export interface MerchantCommunityReply {
  content: string;
  create_time: string;
  nickname?: string;
  reply_id: number;
}

export interface MerchantCommunityPostInput {
  category_id: number;
  content: string;
  image: string;
  product_id: number;
  title: string;
  topic_id: number;
}

export interface MerchantCommunityPage<T> {
  limit: number;
  list: T[];
  page: number;
  total: number;
}

export function listMerchantCommunityPostsApi(params: { limit: number; page: number }) {
  return requestClient.get<MerchantCommunityPage<MerchantCommunityPost>>('/community/posts', { params });
}

export function getMerchantCommunityPostApi(id: number) {
  return requestClient.get<MerchantCommunityPost>(`/community/posts/${id}`);
}

export function createMerchantCommunityPostApi(data: MerchantCommunityPostInput) {
  return requestClient.post<MerchantCommunityPost>('/community/posts', data);
}

export function updateMerchantCommunityPostApi(id: number, data: MerchantCommunityPostInput) {
  return requestClient.put<MerchantCommunityPost>(`/community/posts/${id}`, data);
}

export function deleteMerchantCommunityPostApi(id: number) {
  return requestClient.delete(`/community/posts/${id}`);
}

export function listMerchantCommunityRepliesApi(id: number, params: { limit: number; page: number }) {
  return requestClient.get<MerchantCommunityPage<MerchantCommunityReply>>(`/community/posts/${id}/replies`, { params });
}

export function listMerchantCommunityCategoriesApi() {
  return requestClient.get<{ list: MerchantCommunityCategory[] }>('/community/categories');
}

export function listMerchantCommunityTopicsApi() {
  return requestClient.get<{ list: MerchantCommunityTopic[] }>('/community/topics');
}
