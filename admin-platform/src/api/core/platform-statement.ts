import { requestClient } from '#/api/request';

export type PlatformStatementType = 1 | 2;

export interface PlatformStatementTitle {
  brokerage_expense: number;
  coupon_amount: number;
  merchant_count: number;
  order_income: number;
  platform_charge: number;
  recharge_amount: number;
  recharge_consume: number;
  refund_expense: number;
}

export interface PlatformStatementRow {
  charge: number;
  date: string;
  expend: number;
  income: number;
  offline: number;
}

export interface PlatformStatementPage {
  limit: number;
  list: PlatformStatementRow[];
  page: number;
  total: number;
}

export interface PlatformStatementDetailLine {
  amount: string;
  count: string;
  label: string;
}

export interface PlatformStatementDetailBlock {
  count: string;
  data: PlatformStatementDetailLine[];
  number: number;
  title: string;
}

export interface PlatformStatementDetail {
  bill: PlatformStatementDetailBlock;
  charge: PlatformStatementDetailBlock;
  date: string;
  expend: PlatformStatementDetailBlock;
  income: PlatformStatementDetailBlock;
}

export interface PlatformStatementExport {
  content: string;
  file_name: string;
  row_count: number;
  truncated: boolean;
}

export type PlatformStatementQuery = {
  date_from?: string;
  date_to?: string;
  limit?: number;
  page?: number;
  type?: PlatformStatementType;
};

export function getPlatformStatementTitleApi(params?: {
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PlatformStatementTitle>('/finance/statements/title', {
    params,
  });
}

export function listPlatformStatementsApi(params: PlatformStatementQuery) {
  return requestClient.get<PlatformStatementPage>('/finance/statements', {
    params,
  });
}

export function getPlatformStatementDetailApi(params: {
  date: string;
  type: PlatformStatementType;
}) {
  return requestClient.get<PlatformStatementDetail>(
    '/finance/statements/detail',
    { params },
  );
}

export function exportPlatformStatementApi(input: {
  date: string;
  type: PlatformStatementType;
}) {
  return requestClient.post<PlatformStatementExport>(
    '/finance/statements/export',
    input,
  );
}
