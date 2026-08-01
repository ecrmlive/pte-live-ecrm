import { http } from "@/utils/request";

export interface CustomerServiceThread {
  thread_id: number;
  mer_id: number;
  uid: number;
  service_id: number;
  im_conversation_id: number;
  last_msg: string;
  last_time?: string;
  user_unread: number;
  service_unread: number;
  status: number;
  create_time: string;
  mer_name?: string;
}

export interface CustomerServiceThreadPage {
  list: CustomerServiceThread[];
  total: number;
  page: number;
  limit: number;
}

/** 此凭证仅用于内存中的 PTE IM Web SDK，严禁持久化。 */
export interface CustomerServiceIMCredential {
  mode: string;
  app_id: string;
  sdk_app_id: string;
  im_user_id: string;
  identifier: string;
  user_sig: string;
  expire_at: number;
  api_url: string;
  ws_url: string;
  im_conversation_id: number;
}

export function openCustomerServiceThread(merID: number) {
  return http.post<CustomerServiceThread>("/cs/threads", { mer_id: merID });
}

export function fetchCustomerServiceThreads(page = 1, limit = 50) {
  return http.get<CustomerServiceThreadPage>(`/cs/threads?page=${page}&limit=${limit}`);
}

export function fetchCustomerServiceCredential(threadID: number) {
  return http.get<CustomerServiceIMCredential>(`/cs/im/credential?thread_id=${threadID}`);
}
