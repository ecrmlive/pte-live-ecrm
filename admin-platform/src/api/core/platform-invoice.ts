import { requestClient } from '#/api/request';
export interface PlatformInvoice { id:number; order_id:number; order_no:string; merchant_id:number; merchant_name:string; store_id:number; store_name:string; profile_type:'personal'|'enterprise'; title:string; tax_no_masked:string; email_masked:string; status:'requested'|'issued'|'rejected'|'voided'; invoice_no:string; rejection_reason:string; requested_at:string; issued_at?:string; }
export function listPlatformInvoices(params: {
  page: number;
  limit: number;
  status?: PlatformInvoice['status'];
  order_no?: string;
  keyword?: string;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<{ list: PlatformInvoice[]; total: number }>(
    '/finance/invoices',
    { params },
  );
}
