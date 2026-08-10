import { requestClient } from '#/api/request';

export interface PlatformIntegralLogRow {
  balance: number;
  bill_id: number;
  category: string;
  create_time: string;
  mark: string;
  nickname: string;
  number: number;
  pm: number;
  title: string;
  type: string;
  uid: number;
}

export interface PlatformIntegralLogPage {
  limit: number;
  list: PlatformIntegralLogRow[];
  page: number;
  total: number;
}

export interface PlatformIntegralLogTitle {
  freeze_integral: number;
  order_integral: number;
  sign_count: number;
  sign_integral: number;
  total_integral: number;
  used_integral: number;
}

export interface PlatformIntegralLogExport {
  content: string;
  file_name: string;
  row_count: number;
  truncated: boolean;
}

export type PlatformIntegralLogQuery = {
  date_from?: string;
  date_to?: string;
  keyword?: string;
  limit?: number;
  page?: number;
};

export function listPlatformIntegralLogsApi(params: PlatformIntegralLogQuery) {
  return requestClient.get<PlatformIntegralLogPage>('/integral/logs', { params });
}

export function getPlatformIntegralLogTitleApi(params?: {
  date_from?: string;
  date_to?: string;
  keyword?: string;
}) {
  return requestClient.get<PlatformIntegralLogTitle>('/integral/logs/title', {
    params,
  });
}

export function exportPlatformIntegralLogsApi(input: {
  date_from?: string;
  date_to?: string;
  keyword?: string;
}) {
  return requestClient.post<PlatformIntegralLogExport>(
    '/integral/logs/export',
    input,
  );
}
