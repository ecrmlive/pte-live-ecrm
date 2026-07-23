import { requestClient } from '#/api/request';

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export interface ExpressRow {
  express_id: number;
  name: string;
  code: string;
  sort: number;
  is_show: number;
}

export interface ArticleRow {
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

export interface UserLabelRow {
  label_id: number;
  label_name: string;
  sort: number;
}

export function fetchExpressList(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<ExpressRow>>('/express', { params });
}

export function createExpress(data: {
  name: string;
  code?: string;
  sort?: number;
  is_show?: number;
}) {
  return requestClient.post<ExpressRow>('/express', data);
}

export function updateExpress(
  id: number,
  data: { name: string; code?: string; sort?: number; is_show?: number },
) {
  return requestClient.put<ExpressRow>(`/express/${id}`, data);
}

export function deleteExpress(id: number) {
  return requestClient.delete(`/express/${id}`);
}

export function fetchArticles(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<ArticleRow>>('/articles', { params });
}

export function createArticle(data: Partial<ArticleRow>) {
  return requestClient.post<ArticleRow>('/articles', data);
}

export function updateArticle(id: number, data: Partial<ArticleRow>) {
  return requestClient.put<ArticleRow>(`/articles/${id}`, data);
}

export function deleteArticle(id: number) {
  return requestClient.delete(`/articles/${id}`);
}

export function fetchUserLabels(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<UserLabelRow>>('/user/labels', {
    params,
  });
}

export function createUserLabel(data: { label_name: string; sort?: number }) {
  return requestClient.post<UserLabelRow>('/user/labels', data);
}

export function updateUserLabel(
  id: number,
  data: { label_name: string; sort?: number },
) {
  return requestClient.put<UserLabelRow>(`/user/labels/${id}`, data);
}

export function deleteUserLabel(id: number) {
  return requestClient.delete(`/user/labels/${id}`);
}

export function markUserLabels(uid: number, label_ids: number[]) {
  return requestClient.put(`/user/${uid}/labels`, { label_ids });
}
