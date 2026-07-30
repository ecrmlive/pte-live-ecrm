import { requestClient } from '#/api/request';

export interface CustomerServiceThread {
  id: number;
  user_id: number;
  store_id: number;
  store_name: string;
  merchant_id: number;
  im_sdk_app_id: string;
  order_id?: number | null;
  im_conversation_id: string;
  status: 'open' | 'closed';
  assigned_admin_id?: number | null;
  assigned_at?: string | null;
  created_at: string;
  updated_at: string;
}

export interface CustomerServiceThreadPage {
  list: CustomerServiceThread[];
  total: number;
  page: number;
  limit: number;
}

export function fetchCustomerServiceThreads(params: {
  page: number;
  limit: number;
  mine?: boolean;
  status?: 'closed' | 'open';
}) {
  return requestClient.get<CustomerServiceThreadPage>('/customer-service/threads', {
    params: { ...params, mine: params.mine ? 1 : undefined },
  });
}

export function claimCustomerServiceThread(id: number) {
  return requestClient.post<CustomerServiceThread>(
    `/customer-service/threads/${id}/claim`,
  );
}
