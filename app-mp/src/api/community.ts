import { http } from "@/utils/request";

export interface CommunityTopic {
  topic_id: number;
  topic_name: string;
  is_hot?: number;
  count_use?: number;
}

export interface CommunityPost {
  community_id: number;
  title: string;
  content: string;
  nickname?: string;
  topic_id?: number;
  topic_name?: string;
  product_id?: number;
  product_name?: string;
  product_price?: number;
  count_reply?: number;
  count_start?: number;
  pv?: number;
  create_time?: string;
}

export interface CommunityReply {
  reply_id: number;
  content: string;
  nickname?: string;
  create_time?: string;
}

export function fetchCommunityTopics() {
  return http.get<{ list: CommunityTopic[] }>("/community/topics", false);
}

export function fetchCommunityPosts(page = 1, limit = 20, topicId = 0) {
  const q = topicId ? `&topic_id=${topicId}` : "";
  return http.get<{ list: CommunityPost[]; total: number }>(
    `/community/posts?page=${page}&limit=${limit}${q}`,
    false,
  );
}

export function fetchCommunityPost(id: number) {
  return http.get<CommunityPost>(`/community/posts/${id}`, false);
}

export function createCommunityPost(body: {
  title: string;
  content: string;
  product_id?: number;
  topic_id?: number;
  category_id?: number;
}) {
  return http.post<CommunityPost>("/community/posts", body);
}

export function fetchReplies(postId: number) {
  return http.get<{ list: CommunityReply[]; total: number }>(
    `/community/posts/${postId}/replies`,
    false,
  );
}

export function createReply(postId: number, content: string) {
  return http.post<CommunityReply>(`/community/posts/${postId}/replies`, { content });
}
