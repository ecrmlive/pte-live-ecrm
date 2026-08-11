import { requestClient } from '#/api/request';

export interface PlatformInvoice {
  id: number;
  order_id: number;
  order_sn: string;
  order_no: string;
  mer_id: number;
  merchant_id: number;
  mer_name: string;
  merchant_name: string;
  store_id: number;
  store_name: string;
  uid: number;
  nickname: string;
  user_phone_mask?: string;
  pay_price: number;
  order_amount: number;
  order_status: string;
  order_status_label: string;
  invoice_amount: number;
  receipt_price: number;
  receipt_sn: string;
  invoice_no: string;
  invoice_type: number;
  invoice_type_label: string;
  detail_title: string;
  profile_type: 'personal' | 'enterprise';
  title_type_label: string;
  title: string;
  contact_name: string;
  contact_info: string;
  tax_no?: string;
  tax_no_masked?: string;
  email?: string;
  email_masked?: string;
  status: 'requested' | 'issued' | 'rejected' | 'voided';
  status_label: string;
  issued: boolean;
  mark?: string;
  rejection_reason?: string;
  create_time: string;
  requested_at: string;
  issued_at?: string;
}

export interface PlatformInvoiceQuery {
  page: number;
  limit: number;
  date_from?: string;
  date_to?: string;
  mer_id?: number;
  status?: string;
  order_type?: string;
  keyword?: string;
  user_type?: string;
  user_keyword?: string;
}

export function listPlatformInvoices(params: PlatformInvoiceQuery) {
  return requestClient.get<{ list: PlatformInvoice[]; total: number }>(
    '/finance/invoices',
    { params },
  );
}

export function getPlatformInvoiceApi(id: number) {
  return requestClient.get<PlatformInvoice>(`/finance/invoices/${id}`);
}
