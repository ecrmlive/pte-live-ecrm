import { requestClient } from '#/api/request';

export interface PlatformUserBillRow {
  balance?: number;
  bill_id: number;
  category?: string;
  create_time?: string;
  mark?: string;
  nickname?: string;
  number: number;
  pm?: number;
  title?: string;
  type?: string;
  uid: number;
}

export interface PlatformUserBillPage {
  limit: number;
  list: PlatformUserBillRow[];
  page: number;
  total: number;
}

export interface PlatformUserBillTypeOption {
  title: string;
  type: string;
}

export type PlatformUserBillQuery = {
  date_from?: string;
  date_to?: string;
  limit?: number;
  page?: number;
  type?: string;
  user_keyword?: string;
  user_type?: string;
};

export interface PlatformUserBillExport {
  content: string;
  file_name: string;
  row_count: number;
  truncated: boolean;
}

export function listPlatformUserBillTypesApi() {
  return requestClient.get<{ list: PlatformUserBillTypeOption[] }>(
    '/finance/user-bills/types',
  );
}

export function listPlatformUserBillsApi(params: PlatformUserBillQuery) {
  return requestClient.get<PlatformUserBillPage>('/finance/user-bills', {
    params,
  });
}

export function exportPlatformUserBillsApi(data: PlatformUserBillQuery) {
  return requestClient.post<PlatformUserBillExport>(
    '/finance/user-bills/export',
    data,
  );
}
