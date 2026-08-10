import { requestClient } from '#/api/request';

export interface CommunityCategory {
  category_id: number;
  cate_name: string;
  pid: number;
  is_show: number;
  sort: number;
}

export interface CommunityTopic {
  topic_id: number;
  topic_name: string;
  pic?: string;
  status: number;
  is_hot: number;
  category_id: number;
  cate_name?: string;
  count_use: number;
  sort: number;
  create_time: string;
}

export interface CommunityTopicInput {
  topic_name: string;
  pic?: string;
  category_id: number;
  sort?: number;
  status?: number;
  is_hot?: number;
}

export interface CommunityPost {
  cate_name?: string;
  category_id: number;
  community_id: number;
  content: string;
  count_reply: number;
  count_start: number;
  create_time: string;
  image?: string;
  is_hot: number;
  is_show: number;
  is_type: number;
  mer_id: number;
  nickname?: string;
  product_id: number;
  product_image?: string;
  product_name?: string;
  product_price?: number;
  pv: number;
  refusal?: string;
  start: number;
  status: number;
  status_time?: string;
  title: string;
  topic_id?: number;
  topic_name?: string;
  uid?: number;
  video_link?: string;
}

export interface CommunityReply {
  community_id: number;
  content: string;
  count_reply: number;
  count_start: number;
  create_time: string;
  nickname?: string;
  post_title?: string;
  refusal?: string;
  reply_id: number;
  status: number;
  uid: number;
}

export interface CommunityPage<T> {
  limit: number;
  list: T[];
  page: number;
  total: number;
}

export interface CommunityPostPage extends CommunityPage<CommunityPost> {
  image_count: number;
  video_count: number;
}

export interface CommunityAuditPayload {
  is_hot?: 0 | 1;
  is_show?: 0 | 1;
  refusal?: string;
  status: -2 | -1 | 0 | 1;
}

export interface CommunityReplyAuditPayload {
  refusal?: string;
  status: -1 | 1;
}

export interface CommunityPostListParams {
  author?: string;
  author_type?: string;
  category_id?: number;
  is_show?: number;
  is_type?: number;
  keyword?: string;
  limit: number;
  page: number;
  status?: number;
  topic_id?: number;
}

export interface CommunityReplyListParams {
  date_from?: string;
  date_to?: string;
  keyword?: string;
  limit: number;
  page: number;
  username?: string;
}

export interface CommunityCategoryInput {
  cate_name: string;
  is_show?: number;
  sort?: number;
}

export function listCommunityCategoriesApi() {
  return requestClient.get<{ list: CommunityCategory[] }>('/community/categories');
}

export function createCommunityCategoryApi(data: CommunityCategoryInput) {
  return requestClient.post<CommunityCategory>('/community/categories', data);
}

export function updateCommunityCategoryApi(id: number, data: CommunityCategoryInput) {
  return requestClient.put<CommunityCategory>(`/community/categories/${id}`, data);
}

export function updateCommunityCategoryStatusApi(id: number, isShow: number) {
  return requestClient.put(`/community/categories/${id}/status`, { is_show: isShow });
}

export function deleteCommunityCategoryApi(id: number) {
  return requestClient.delete(`/community/categories/${id}`);
}

export function listCommunityTopicsApi() {
  return requestClient.get<{ list: CommunityTopic[] }>('/community/topics');
}

export function createCommunityTopicApi(data: CommunityTopicInput) {
  return requestClient.post<CommunityTopic>('/community/topics', data);
}

export function updateCommunityTopicApi(id: number, data: CommunityTopicInput) {
  return requestClient.put<CommunityTopic>(`/community/topics/${id}`, data);
}

export function updateCommunityTopicStatusApi(id: number, status: number) {
  return requestClient.put(`/community/topics/${id}/status`, { status });
}

export function updateCommunityTopicHotApi(id: number, isHot: number) {
  return requestClient.put(`/community/topics/${id}/hot`, { is_hot: isHot });
}

export function deleteCommunityTopicApi(id: number) {
  return requestClient.delete(`/community/topics/${id}`);
}

export function listCommunityPostsApi(params: CommunityPostListParams) {
  return requestClient.get<CommunityPostPage>('/community/posts', { params });
}

export function getCommunityPostApi(id: number) {
  return requestClient.get<CommunityPost>(`/community/posts/${id}`);
}

export function listCommunityRepliesApi(id: number, params: { limit: number; page: number }) {
  return requestClient.get<CommunityPage<CommunityReply>>(`/community/posts/${id}/replies`, { params });
}

export function listAllCommunityRepliesApi(params: CommunityReplyListParams) {
  return requestClient.get<CommunityPage<CommunityReply>>('/community/replies', { params });
}

export function auditCommunityPostApi(id: number, data: CommunityAuditPayload) {
  return requestClient.post<CommunityPost>(`/community/posts/${id}/audit`, data);
}

export function updateCommunityPostStarApi(id: number, start: number) {
  return requestClient.put<CommunityPost>(`/community/posts/${id}/star`, { start });
}

export function updateCommunityPostShowApi(id: number, isShow: number) {
  return requestClient.put<CommunityPost>(`/community/posts/${id}/show`, { is_show: isShow });
}

export function deleteCommunityPostApi(id: number) {
  return requestClient.delete(`/community/posts/${id}`);
}

export function auditCommunityReplyApi(id: number, data: CommunityReplyAuditPayload) {
  return requestClient.post<CommunityReply>(`/community/replies/${id}/audit`, data);
}

export function deleteCommunityReplyApi(id: number) {
  return requestClient.delete(`/community/replies/${id}`);
}
