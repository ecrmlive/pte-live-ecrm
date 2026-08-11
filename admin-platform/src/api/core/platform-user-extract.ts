import { requestClient } from '#/api/request';

export interface PlatformUserExtractRow {
  account?: string;
  admin_id?: number;
  alipay_code?: string;
  balance?: number;
  bank_address?: string;
  bank_code?: string;
  bank_name?: string;
  create_time?: string;
  extract_id: number;
  extract_pic?: string;
  extract_price: number;
  extract_sn?: string;
  extract_type: number;
  extract_type_name?: string;
  fail_msg?: string;
  mark?: string;
  nickname?: string;
  real_name?: string;
  status: number;
  status_name?: string;
  status_time?: string;
  uid: number;
  wechat?: string;
}

export interface PlatformUserExtractPage {
  limit: number;
  list: PlatformUserExtractRow[];
  page: number;
  total: number;
}

export interface PlatformUserExtractExport {
  content: string;
  file_name: string;
  row_count: number;
  truncated: boolean;
}

export type PlatformUserExtractQuery = {
  account_keyword?: string;
  date_from?: string;
  date_to?: string;
  extract_type?: number;
  limit?: number;
  page?: number;
  status?: number;
  user_keyword?: string;
  user_type?: string;
};

export function listPlatformUserExtractsApi(params: PlatformUserExtractQuery) {
  return requestClient.get<PlatformUserExtractPage>('/finance/user-extracts', {
    params,
  });
}

export function getPlatformUserExtractApi(id: number) {
  return requestClient.get<PlatformUserExtractRow>(
    `/finance/user-extracts/${id}`,
  );
}

export function switchPlatformUserExtractStatusApi(
  id: number,
  data: { fail_msg?: string; mark?: string; status: number },
) {
  return requestClient.post<{ ok: boolean; status: number }>(
    `/finance/user-extracts/${id}/status`,
    data,
  );
}

export function exportPlatformUserExtractsApi(data: PlatformUserExtractQuery) {
  return requestClient.post<PlatformUserExtractExport>(
    '/finance/user-extracts/export',
    data,
  );
}
