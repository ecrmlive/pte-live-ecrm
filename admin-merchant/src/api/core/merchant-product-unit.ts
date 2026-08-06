import { requestClient } from '#/api/request';

export interface ProductUnit {
  name: string;
  sort: number;
  unit_id: number;
}

export interface ProductUnitListParams {
  keyword?: string;
}

export function listProductUnitsApi(params: ProductUnitListParams = {}) {
  return requestClient.get<{ list: ProductUnit[] }>('/product/units', { params });
}

export function saveProductUnitsApi(body: { list: ProductUnit[] }) {
  return requestClient.post<{ list: ProductUnit[] }>('/product/units', body);
}
