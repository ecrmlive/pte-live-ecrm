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

export interface MerchantOperationLogListParams {
  action?: string;
  date_from?: string;
  date_to?: string;
  keyword?: string;
  limit: number;
  page: number;
}

export function listMerchantOperationLogsApi(params: MerchantOperationLogListParams) {
  return requestClient.get<MerchantOperationLogPage>('/setting/operation-logs', { params });
}
