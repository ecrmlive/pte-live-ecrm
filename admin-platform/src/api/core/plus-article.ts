import { requestClient } from '#/api/request';

export interface ArticleCategoryOption {
  cid: number;
  status: number;
  title: string;
}

export async function getArticleCategoryListApi() {
  return requestClient.get<{ list: ArticleCategoryOption[] }>('/article/categories');
}

export interface ArticleCategoryInput {
  sort?: number;
  status?: number;
  title: string;
}

export function createArticleCategoryApi(data: ArticleCategoryInput) {
  return requestClient.post<ArticleCategoryOption>('/article/categories', data);
}

export function updateArticleCategoryApi(id: number, data: ArticleCategoryInput) {
  return requestClient.put<ArticleCategoryOption>(`/article/categories/${id}`, data);
}

export function deleteArticleCategoryApi(id: number) {
  return requestClient.delete(`/article/categories/${id}`);
}
