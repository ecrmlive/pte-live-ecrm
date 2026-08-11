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

export interface CustomerServiceEvent {
  id: number;
  binding_id: number;
  sender_role: 'service' | 'system' | 'user';
  sender_id: number;
  msg_type: 'order' | 'system';
  content: string;
  created_at: string;
}

export interface CustomerServiceEventPage {
  conversation_id: string;
  list: CustomerServiceEvent[];
  total: number;
  page: number;
  limit: number;
  note: string;
}

export interface CustomerServiceAssignmentLog {
  id: number;
  binding_id: number;
  from_admin_id?: number | null;
  target_admin_id: number;
  operator_admin_id: number;
  reason: string;
  created_at: string;
}

export interface CustomerServiceAssignmentLogPage {
  list: CustomerServiceAssignmentLog[];
  total: number;
  page: number;
  limit: number;
}

export interface CustomerServiceQuickReply {
  id: number;
  store_id: number;
  title: string;
  content: string;
  message_type: 'image' | 'text';
  status: 'disabled' | 'enabled';
  created_by: number;
  updated_by: number;
  created_at: string;
  updated_at: string;
}

export interface CustomerServiceQuickReplyPage {
  list: CustomerServiceQuickReply[];
  total: number;
  page: number;
  limit: number;
}

export interface CustomerServiceQuickReplyInput {
  store_id: number;
  title: string;
  content: string;
  message_type?: 'image' | 'text';
  status?: 'disabled' | 'enabled';
}

export interface CustomerServiceAgent {
  account: string;
  avatar_url?: string;
  created_at: string;
  id: number;
  display_name: string;
  linked_user_id?: number;
  phone: string;
  roles: string;
  status: 0 | 1;
  service_store_ids: number[];
  wechat_username?: string;
}

export interface CustomerServiceAgentPage {
  list: CustomerServiceAgent[];
  total: number;
  page: number;
  limit: number;
}

export interface CustomerServiceAgentUser {
  binding_id: number;
  user_id: number;
  nickname: string;
  mobile: string;
  store_id: number;
  store_name: string;
  status: 'closed' | 'open';
  updated_at: string;
}

export interface CustomerServiceAgentUserPage {
  list: CustomerServiceAgentUser[];
  total: number;
  page: number;
  limit: number;
}

export interface CustomerServiceSettings {
  auto_reply_enabled: boolean;
  auto_reply_text: string;
  enterprise_wechat_corp_id: string;
  enterprise_wechat_url: string;
  queue_mode: 'manual' | 'round_robin';
  max_sessions_per_agent: number;
  redirect_url: string;
  service_phone: string;
  service_type:
    | 'disabled'
    | 'enterprise_wechat'
    | 'link'
    | 'mini_program'
    | 'phone'
    | 'system';
}

export interface CustomerServiceSettingsResult {
  settings: CustomerServiceSettings;
  updated_at?: string | null;
}

export function fetchCustomerServiceAgents(params?: {
  keyword?: string;
  status?: number;
  page?: number;
  limit?: number;
}) {
  return requestClient.get<CustomerServiceAgentPage>('/customer-service/agents', { params });
}

export function fetchCustomerServiceAgentUsers(agentID: number, params: { page: number; limit: number }) {
  return requestClient.get<CustomerServiceAgentUserPage>(`/customer-service/agents/${agentID}/users`, { params });
}

export function fetchCustomerServiceSettings() {
  return requestClient.get<CustomerServiceSettingsResult>('/customer-service/settings');
}

export function updateCustomerServiceSettings(data: CustomerServiceSettings) {
  return requestClient.put<CustomerServiceSettingsResult>('/customer-service/settings', data);
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

export function fetchCustomerServiceThread(id: number) {
  return requestClient.get<CustomerServiceThread>(`/customer-service/threads/${id}`);
}

export function claimCustomerServiceThread(id: number) {
  return requestClient.post<CustomerServiceThread>(
    `/customer-service/threads/${id}/claim`,
  );
}

export function transferCustomerServiceThread(
  id: number,
  data: { reason: string; target_admin_id: number },
  idempotencyKey: string,
) {
  return requestClient.post<CustomerServiceThread>(
    `/customer-service/threads/${id}/transfer`,
    data,
    { headers: { 'X-Idempotency-Key': idempotencyKey } },
  );
}

export function fetchCustomerServiceEvents(id: number, params: { page: number; limit: number }) {
  return requestClient.get<CustomerServiceEventPage>(`/customer-service/threads/${id}/messages`, { params });
}

export function fetchCustomerServiceAssignmentLogs(id: number, params: { page: number; limit: number }) {
  return requestClient.get<CustomerServiceAssignmentLogPage>(`/customer-service/threads/${id}/assignment-logs`, { params });
}

export function fetchCustomerServiceOrder(id: number) {
  return requestClient.get<Record<string, unknown>>(`/customer-service/threads/${id}/order`);
}

export function fetchCustomerServiceDelivery(id: number) {
  return requestClient.get<Record<string, unknown>>(`/customer-service/threads/${id}/delivery`);
}

export function fetchCustomerServiceProducts(id: number) {
  return requestClient.get<{ list: Record<string, unknown>[] }>(`/customer-service/threads/${id}/products`);
}

export function fetchCustomerServiceRefunds(id: number) {
  return requestClient.get<{ list: Record<string, unknown>[] }>(`/customer-service/threads/${id}/refunds`);
}

export function fetchCustomerServiceUser(id: number) {
  return requestClient.get<Record<string, unknown>>(`/customer-service/threads/${id}/user`);
}

export function updateCustomerServiceUserNote(id: number, content: string) {
  return requestClient.put<{ user_id: number; store_id: number; content: string }>(`/customer-service/threads/${id}/user-note`, { content });
}

export function fetchCustomerServiceQuickReplies(params: {
  page: number;
  limit: number;
  keyword?: string;
  store_id?: number;
}) {
  return requestClient.get<CustomerServiceQuickReplyPage>('/customer-service/quick-replies', { params });
}

export function createCustomerServiceQuickReply(data: CustomerServiceQuickReplyInput) {
  return requestClient.post<CustomerServiceQuickReply>('/customer-service/quick-replies', data);
}

export function updateCustomerServiceQuickReply(
  id: number,
  data: Omit<CustomerServiceQuickReplyInput, 'store_id'> & { store_id?: number },
) {
  return requestClient.put<CustomerServiceQuickReply>(`/customer-service/quick-replies/${id}`, data);
}

export function deleteCustomerServiceQuickReply(id: number) {
  return requestClient.delete(`/customer-service/quick-replies/${id}`);
}
