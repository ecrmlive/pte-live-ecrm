import { requestClient } from '#/api/request';

export interface MerchantOperationLog {
  account_id: number;
  action: string;
  created_at: string;
  id: number;
  request_id: string;
  resource_id: string;
  resource_type: string;
}

export interface MerchantOperationLogPage {
  limit: number;
  list: MerchantOperationLog[];
  page: number;
  total: number;
}

export function listMerchantOperationLogsApi(params: { limit: number; page: number }) {
  return requestClient.get<MerchantOperationLogPage>('/setting/operation-logs', { params });
}
