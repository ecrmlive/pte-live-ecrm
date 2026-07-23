import { requestClient } from '#/api/request';

import type { PaginatedList, ProductCategoryOption } from './product';

export interface ShopLinkPickerItem {
  id?: number;
  name: string;
  type: string;
  url: string;
}

export interface ShopLinkMarketingLists {
  invitationList: ShopLinkPickerItem[];
  packageList: ShopLinkPickerItem[];
  tableList: ShopLinkPickerItem[];
}

export interface ShopLinkDiyPageItem {
  page_id: number;
  page_name: string;
}

export interface ShopLinkProductRow {
  image: Array<{ file_path: string }>;
  product_id: number;
  product_name: string;
  product_price: number | string;
}

export interface ShopLinkArticleCategory {
  category_id: number;
  child?: ShopLinkArticleCategory[];
  name: string;
}

export interface ShopLinkArticleRow {
  article_id: number;
  article_title: string;
  category?: { name: string };
}

export async function getShopLinkMarketingListsApi() {
  return requestClient.post<ShopLinkMarketingLists>('/shop/data.link/index', {});
}

export async function getShopLinkDiyPageListApi() {
  return requestClient.post<{ list: ShopLinkDiyPageItem[] }>(
    '/shop/data.link/getPageList',
    {},
  );
}

export async function getShopLinkProductCategoryApi() {
  return requestClient.post<{ list: ProductCategoryOption[] }>(
    '/shop/data.product/category',
    {},
  );
}

export async function getShopLinkProductListsApi(params: {
  list_rows?: number;
  page?: number;
  product_name?: string;
}) {
  return requestClient.post<{ list: PaginatedList<ShopLinkProductRow> }>(
    '/shop/data.product/lists',
    params,
  );
}

export async function getShopLinkArticleCategoryApi() {
  return requestClient.post<{ list: ShopLinkArticleCategory[] }>(
    '/shop/data.article/category',
    {},
  );
}

export async function getShopLinkArticleListApi(params: {
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<{ list: PaginatedList<ShopLinkArticleRow> }>(
    '/shop/data.article/index',
    params,
  );
}
