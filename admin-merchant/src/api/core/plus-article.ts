import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface ArticleCategoryItem {
  category_id: number;
  create_time?: string;
  name: string;
  sort?: number;
}

export interface ArticleListItem {
  article_id: number;
  article_sort: number;
  article_status: number;
  article_title: string;
  category?: { name: string };
  create_time: string;
  image?: { file_path: string };
  update_time: string;
  virtual_views: number;
}

export interface ArticleListQuery {
  article_title?: string;
  category_id?: number;
  list_rows?: number;
  page?: number;
}

export interface ArticleFormValues {
  article_content: string;
  article_id?: number;
  article_sort: number | string;
  article_status: number;
  article_title: string;
  category_id: number | string;
  dec: string;
  image?: { file_path?: string };
  image_id: number | string;
  virtual_views: number | string;
}

export async function getArticleListApi(params: ArticleListQuery) {
  return requestClient.post<{
    categoryList: ArticleCategoryItem[];
    list: PaginatedList<ArticleListItem>;
  }>('/shop/plus.article.article/index', params);
}

export async function deleteArticleApi(articleId: number) {
  return requestClient.post('/shop/plus.article.article/delete', { article_id: articleId });
}

export async function getArticleCategoryListApi(params?: { page_id?: number }) {
  return requestClient.post<{ category: ArticleCategoryItem[] }>(
    '/shop/plus.article.category/index',
    params ?? {},
  );
}

export async function addArticleCategoryApi(payload: { name: string; sort: number | string }) {
  return requestClient.post('/shop/plus.article.category/add', payload);
}

export async function editArticleCategoryApi(payload: ArticleCategoryItem) {
  return requestClient.post('/shop/plus.article.category/edit', payload);
}

export async function deleteArticleCategoryApi(categoryId: number) {
  return requestClient.post('/shop/plus.article.category/delete', { category_id: categoryId });
}

export async function getArticleAddMetaApi() {
  return requestClient.get<{ catgory: ArticleCategoryItem[] }>('/shop/plus.article.article/add');
}

export async function addArticleApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.article.article/add', payload);
}

export async function getArticleEditMetaApi(articleId: number) {
  return requestClient.get<{ catgory: ArticleCategoryItem[]; model: ArticleFormValues }>(
    '/shop/plus.article.article/edit',
    { params: { article_id: articleId } },
  );
}

export async function editArticleApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.article.article/edit', payload);
}
