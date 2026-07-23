import { http } from '@/api/http';

export interface Invoice {
  invoice_id: number;
  uid: number;
  order_id: number;
  invoice_type: number;
  header_type: number;
  header: string;
  tax_no: string;
  email: string;
  status: number;
  mark: string;
  create_time?: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchInvoices(params: Record<string, unknown>) {
  return http.get<PageResult<Invoice>>('/invoices', { params });
}

export function auditInvoice(id: number, data: { status: number; mark?: string }) {
  return http.put<Invoice>(`/invoices/${id}/audit`, data);
}
