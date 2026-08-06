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

export interface MerchantExpressListParams {
  keyword?: string;
  is_show?: 0 | 1;
  limit: number;
  page: number;
}

export function listMerchantExpressApi(params: MerchantExpressListParams) {
  return requestClient.get<MerchantExpressPage>('/express', { params });
}
