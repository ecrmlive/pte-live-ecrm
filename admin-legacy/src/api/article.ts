import { http } from '@/api/http';

export interface Article {
  article_id: number;
  cid: number;
  title: string;
  author: string;
  synopsis: string;
  content: string;
  sort: number;
  status: number;
  visit?: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchArticles(params: Record<string, unknown>) {
  return http.get<PageResult<Article>>('/articles', { params });
}

export function createArticle(data: Partial<Article>) {
  return http.post<Article>('/articles', data);
}

export function updateArticle(id: number, data: Partial<Article>) {
  return http.put<Article>(`/articles/${id}`, data);
}

export function deleteArticle(id: number) {
  return http.delete<{ ok: boolean }>(`/articles/${id}`);
}
