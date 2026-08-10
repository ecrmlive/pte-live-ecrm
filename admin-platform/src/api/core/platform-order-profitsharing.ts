import { requestClient } from '#/api/request';

export interface PlatformOrderProfitsharingRow {
  create_time?: string;
  error_msg?: string;
  is_combine?: number;
  mer_id: number;
  mer_name: string;
  order_id: number;
  order_sn: string;
  profitsharing_id: number;
  profitsharing_mer_price?: number;
  profitsharing_price: number;
  profitsharing_refund?: number;
  profitsharing_sn: string;
  profitsharing_time?: string;
  status: number;
  status_name?: string;
  transaction_id?: string;
  type: string;
  type_name?: string;
}

export interface PlatformOrderProfitsharingPage {
  limit: number;
  list: PlatformOrderProfitsharingRow[];
  page: number;
  total: number;
}

export interface PlatformOrderProfitsharingExport {
  content: string;
  file_name: string;
  row_count: number;
  truncated: boolean;
}

export type PlatformOrderProfitsharingQuery = {
  date_from?: string;
  date_to?: string;
  limit?: number;
  mer_id?: number;
  mer_name?: string;
  page?: number;
  profit_date_from?: string;
  profit_date_to?: string;
  status?: number;
  type?: string;
};

export function listPlatformOrderProfitsharingsApi(
  params: PlatformOrderProfitsharingQuery,
) {
  return requestClient.get<PlatformOrderProfitsharingPage>(
    '/finance/order-profitsharings',
    { params },
  );
}

export function getPlatformOrderProfitsharingApi(id: number) {
  return requestClient.get<PlatformOrderProfitsharingRow>(
    `/finance/order-profitsharings/${id}`,
  );
}

export function againPlatformOrderProfitsharingApi(id: number) {
  return requestClient.post<{ ok: boolean }>(
    `/finance/order-profitsharings/${id}/again`,
  );
}

export function exportPlatformOrderProfitsharingsApi(
  data: PlatformOrderProfitsharingQuery,
) {
  return requestClient.post<PlatformOrderProfitsharingExport>(
    '/finance/order-profitsharings/export',
    data,
  );
}
