import { requestClient } from '#/api/request';

export interface PlatformNotice { content: string; create_time: string; is_show: number; notice_id: number; sort: number; title: string; }
export interface PlatformNoticePage { limit: number; list: PlatformNotice[]; page: number; total: number; }
export interface NoticeSaveInput { content: string; is_show: number; sort: number; title: string; }
export function listPlatformNoticesApi(params: { limit: number; page: number }) { return requestClient.get<PlatformNoticePage>('/notices', { params }); }
export function createPlatformNoticeApi(data: NoticeSaveInput) { return requestClient.post<PlatformNotice>('/notices', data); }
export function updatePlatformNoticeApi(id: number, data: NoticeSaveInput) { return requestClient.put<PlatformNotice>(`/notices/${id}`, data); }
export function deletePlatformNoticeApi(id: number) { return requestClient.delete(`/notices/${id}`); }

export interface PlatformAgreement { content: string; key: string; label: string; }
export function listPlatformAgreementsApi() { return requestClient.get<{ list: PlatformAgreement[] }>('/agreements'); }
export function getPlatformAgreementApi(key: string) {
  return requestClient.get<PlatformAgreement>(`/agreements/${encodeURIComponent(key)}`);
}
export function savePlatformAgreementApi(key: string, content: string) { return requestClient.put<PlatformAgreement>(`/agreements/${key}`, { content }); }
