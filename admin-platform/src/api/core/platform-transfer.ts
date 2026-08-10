import { requestClient } from '#/api/request';

export interface PlatformTransferTitle {
  applying_amount: number;
  applying_merchant_count: number;
  freeze_amount: number;
  payable_amount: number;
  pending_amount: number;
  withdrawable_amount: number;
}

export interface PlatformTransferRow {
  admin_id?: number | null;
  admin_mark?: string;
  admin_name?: string;
  create_time?: string;
  extract_money: number;
  financial_account: string;
  financial_id: number;
  financial_sn: string;
  financial_status: number;
  financial_type: number;
  image?: string;
  is_trader: number;
  mark?: string;
  mer_id: number;
  mer_money: number;
  mer_name: string;
  refusal?: string;
  status: number;
  status_time?: string;
  update_time?: string;
}

export interface PlatformTransferPage {
  limit: number;
  list: PlatformTransferRow[];
  page: number;
  total: number;
}

export interface PlatformTransferExport {
  content: string;
  file_name: string;
  row_count: number;
  truncated: boolean;
}

export type PlatformTransferQuery = {
  admin_keyword?: string;
  date_from?: string;
  date_to?: string;
  financial_status?: number;
  financial_type?: number;
  is_trader?: number;
  limit?: number;
  mer_name?: string;
  page?: number;
  status?: number;
};

export function getPlatformTransferTitleApi(params?: PlatformTransferQuery) {
  return requestClient.get<PlatformTransferTitle>('/finance/transfers/title', {
    params,
  });
}

export function listPlatformTransfersApi(params: PlatformTransferQuery) {
  return requestClient.get<PlatformTransferPage>('/finance/transfers', {
    params,
  });
}

export function getPlatformTransferApi(id: number) {
  return requestClient.get<PlatformTransferRow>(`/finance/transfers/${id}`);
}

export function auditPlatformTransferApi(
  id: number,
  data: { refusal?: string; status: number },
) {
  return requestClient.post(`/finance/transfers/${id}/status`, data);
}

export function markPlatformTransferApi(id: number, admin_mark: string) {
  return requestClient.post(`/finance/transfers/${id}/mark`, { admin_mark });
}

export function payPlatformTransferApi(id: number, image: string) {
  return requestClient.post(`/finance/transfers/${id}/pay`, { image });
}

export function exportPlatformTransfersApi(params?: PlatformTransferQuery) {
  return requestClient.post<PlatformTransferExport>(
    '/finance/transfers/export',
    params ?? {},
  );
}
