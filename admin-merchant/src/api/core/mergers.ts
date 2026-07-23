import { requestClient } from '#/api/request';

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export interface InvoiceRow {
  invoice_id: number;
  uid: number;
  order_id: number;
  header: string;
  tax_no: string;
  email: string;
  status: number;
  mark: string;
}

export interface ShippingTemplate {
  template_id: number;
  name: string;
  type: number;
  sort: number;
  regions?: Array<{
    first?: number;
    first_price?: number;
    continue?: number;
    continue_price?: number;
  }>;
}

export function fetchInvoices(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<InvoiceRow>>('/invoices', { params });
}

export function auditInvoice(id: number, data: { status: number; mark?: string }) {
  return requestClient.put<InvoiceRow>(`/invoices/${id}/audit`, data);
}

export function fetchShippingTemplates(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<ShippingTemplate>>('/shipping/templates', {
    params,
  });
}

export function createShippingTemplate(data: Record<string, unknown>) {
  return requestClient.post<ShippingTemplate>('/shipping/templates', data);
}

export function updateShippingTemplate(id: number, data: Record<string, unknown>) {
  return requestClient.put<ShippingTemplate>(`/shipping/templates/${id}`, data);
}

export function deleteShippingTemplate(id: number) {
  return requestClient.delete(`/shipping/templates/${id}`);
}

export function getShippingTemplate(id: number) {
  return requestClient.get<ShippingTemplate>(`/shipping/templates/${id}`);
}
