import { requestClient } from '#/api/request';

export interface ProductUnit {
  name: string;
  sort: number;
  unit_id: number;
}

export function listProductUnitsApi() {
  return requestClient.get<{ list: ProductUnit[] }>('/product/units');
}

export function saveProductUnitsApi(body: { list: ProductUnit[] }) {
  return requestClient.post<{ list: ProductUnit[] }>('/product/units', body);
}
