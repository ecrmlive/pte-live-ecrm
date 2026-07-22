import { http } from "@/utils/request";

export interface Notice {
  notice_id: number;
  title: string;
  content: string;
  sort?: number;
  create_time?: string;
}

export function fetchNotices(page = 1, limit = 20) {
  return http.get<{ list: Notice[]; total: number }>(
    `/notices?page=${page}&limit=${limit}`,
    false,
  );
}

export function fetchNotice(id: number) {
  return http.get<Notice>(`/notices/${id}`, false);
}

export interface Agreement {
  key: string;
  label: string;
  content: string;
}

export function fetchAgreement(key: string) {
  return http.get<Agreement>(`/agreements/${encodeURIComponent(key)}`, false);
}
