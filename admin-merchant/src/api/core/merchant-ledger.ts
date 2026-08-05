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

export function listFinanceLedgerApi(params: { limit: number; page: number }) {
  return requestClient.get<FinanceLedgerPage>('/finance/ledger', { params });
}

export function listFinanceStatementsApi(params: { limit: number; page: number }) {
  return requestClient.get<FinanceStatementPage>('/finance/statements', { params });
}
