import { requestClient } from '#/api/request';

export interface PlatformCapitalFlowRow {
  create_time?: string;
  financial_pm?: number;
  financial_record_id: number;
  financial_record_sn?: string;
  financial_type?: string;
  financial_type_name?: string;
  mer_id?: number;
  number?: number;
  order_id?: number;
  order_sn?: string;
  pay_type?: number;
  pay_type_name?: string;
  signed_number?: number;
  transaction_id?: string;
  type?: number;
  user_id?: number;
  user_info?: string;
}

export interface PlatformCapitalFlowPage {
  limit: number;
  list: PlatformCapitalFlowRow[];
  page: number;
  total: number;
}

export type PlatformCapitalFlowQuery = {
  date_from?: string;
  date_to?: string;
  limit?: number;
  order_sn?: string;
  page?: number;
  pay_type?: string;
  user_keyword?: string;
  user_type?: string;
};

export interface PlatformCapitalFlowExport {
  content: string;
  file_name: string;
  row_count: number;
  truncated: boolean;
}

export function listPlatformCapitalFlowsApi(params: PlatformCapitalFlowQuery) {
  return requestClient.get<PlatformCapitalFlowPage>('/finance/capital-flows', {
    params,
  });
}

export function getPlatformCapitalFlowApi(id: number) {
  return requestClient.get<PlatformCapitalFlowRow>(
    `/finance/capital-flows/${id}`,
  );
}

export function exportPlatformCapitalFlowsApi(data: PlatformCapitalFlowQuery) {
  return requestClient.post<PlatformCapitalFlowExport>(
    '/finance/capital-flows/export',
    data,
  );
}
