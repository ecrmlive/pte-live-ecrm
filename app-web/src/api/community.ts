import { http } from "@/utils/request";

export interface CommunityTopic {
  topic_id: number;
  topic_name: string;
  is_hot?: number;
}

export interface CommunityPost {
  community_id: number;
  title: string;
  content: string;
  nickname?: string;
  topic_name?: string;
  product_id?: number;
  product_name?: string;
  product_price?: number;
  count_reply?: number;
  pv?: number;
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
