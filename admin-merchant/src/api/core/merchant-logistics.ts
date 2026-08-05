import { requestClient } from '#/api/request';

export interface MerchantExpressRow {
  code: string;
  create_time: string;
  express_id: number;
  is_show: number;
  name: string;
  sort: number;
}

export interface MerchantExpressPage {
  limit: number;
  list: MerchantExpressRow[];
  page: number;
  total: number;
}

export function listMerchantExpressApi(params: { limit: number; page: number }) {
  return requestClient.get<MerchantExpressPage>('/express', { params });
}
