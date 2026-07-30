import { requestClient } from '#/api/request';

export interface ArticleCategoryOption {
  cid: number;
  status: number;
  title: string;
}

export async function getArticleCategoryListApi() {
  return requestClient.get<{ list: ArticleCategoryOption[] }>('/article/categories');
}
