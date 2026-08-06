import { requestClient } from '#/api/request';

export interface FinanceLedgerEntry {
  amount: number;
  created_at: string;
  entry_type: string;
  id: number;
  reference_id: string;
  reference_type: string;
}

export interface FinanceLedgerPage {
  limit: number;
  list: FinanceLedgerEntry[];
  page: number;
  total: number;
}

export interface FinanceStatement {
  amount: number;
  period_end: string;
  period_start: string;
  statement_id: number;
  status: string;
  updated_at: string;
}

export interface FinanceStatementPage {
  limit: number;
  list: FinanceStatement[];
  page: number;
  total: number;
}

export interface FinanceLedgerListParams {
  date_from?: string;
  date_to?: string;
  entry_type?: string;
  limit: number;
  page: number;
}

export function listFinanceLedgerApi(params: FinanceLedgerListParams) {
  return requestClient.get<FinanceLedgerPage>('/finance/ledger', { params });
}

export type FinanceStatementStatus =
  | 'approved'
  | 'bill_frozen'
  | 'bill_pending'
  | 'paid'
  | 'rejected'
  | 'withdraw_applied';

export interface FinanceStatementListParams {
  date_from?: string;
  date_to?: string;
  limit: number;
  page: number;
  status?: FinanceStatementStatus;
}

export function listFinanceStatementsApi(params: FinanceStatementListParams) {
  return requestClient.get<FinanceStatementPage>('/finance/statements', { params });
}
