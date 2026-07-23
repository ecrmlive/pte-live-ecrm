import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface PointsProductListItem {
  limit_num?: number;
  point_product_id: number;
  product?: {
    image?: Array<{ file_path: string }>;
    product_name?: string;
    spec_type?: number;
  };
  sku?: Array<{ point_money?: number | string; point_num?: number }>;
  sort?: number;
  status?: number;
  stock?: number;
}

export async function getPointsProductListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{
    exclude_ids?: number[];
    list: PaginatedList<PointsProductListItem>;
  }>('/shop/plus.points.product/index', params);
}

export async function deletePointsProductApi(pointProductId: number) {
  return requestClient.post('/shop/plus.points.product/delete', { id: pointProductId });
}

export async function getPointsProductAddMetaApi(productId: number) {
  return requestClient.get<{ model: Record<string, unknown> }>(
    '/shop/plus.points.product/add',
    { params: { product_id: productId } },
  );
}

export async function getPointsProductEditMetaApi(pointProductId: number) {
  return requestClient.get<{ model: Record<string, unknown> }>(
    '/shop/plus.points.product/edit',
    { params: { point_product_id: pointProductId } },
  );
}

export async function addPointsProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.points.product/add', payload);
}

export async function editPointsProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.points.product/edit', payload);
}
