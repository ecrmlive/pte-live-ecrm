import { requestClient } from '#/api/request';

export type NoticeScopeType =
  | 'all'
  | 'store_category'
  | 'store_name'
  | 'store_type';

export interface PlatformNoticeScope {
  id: number;
  name: string;
  scope_type: NoticeScopeType;
}

export interface PlatformNotice {
  content: string;
  create_time: string;
  is_show: number;
  notice_id: number;
  scope_ids: number[];
  scope_items: PlatformNoticeScope[];
  scope_type: NoticeScopeType;
  title: string;
}
export interface PlatformNoticePage { limit: number; list: PlatformNotice[]; page: number; total: number; }
export interface NoticeSaveInput {
  content: string;
  is_show: number;
  scope_ids: number[];
  scope_type: NoticeScopeType;
  title: string;
}
export interface PlatformNoticeQuery {
  date_from?: string;
  date_to?: string;
  is_show?: number;
  keyword?: string;
  limit: number;
  page: number;
}
export function listPlatformNoticesApi(params: PlatformNoticeQuery) { return requestClient.get<PlatformNoticePage>('/notices', { params }); }
export function getPlatformNoticeApi(id: number) { return requestClient.get<PlatformNotice>(`/notices/${id}`); }
export function createPlatformNoticeApi(data: NoticeSaveInput) { return requestClient.post<PlatformNotice>('/notices', data); }
export function updatePlatformNoticeApi(id: number, data: NoticeSaveInput) { return requestClient.put<PlatformNotice>(`/notices/${id}`, data); }
export function updatePlatformNoticeStatusApi(id: number, is_show: number) { return requestClient.put(`/notices/${id}/status`, { is_show }); }
export function deletePlatformNoticeApi(id: number) { return requestClient.delete(`/notices/${id}`); }

export interface PlatformAgreement { content: string; key: string; label: string; }
export function listPlatformAgreementsApi() { return requestClient.get<{ list: PlatformAgreement[] }>('/agreements'); }
export function getPlatformAgreementApi(key: string) {
  return requestClient.get<PlatformAgreement>(`/agreements/${encodeURIComponent(key)}`);
}
export function savePlatformAgreementApi(key: string, content: string) { return requestClient.put<PlatformAgreement>(`/agreements/${key}`, { content }); }
