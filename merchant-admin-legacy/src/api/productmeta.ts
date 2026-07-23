import { http } from '@/api/http';

export interface ProductLabel {
  label_id: number;
  name: string;
  info: string;
  sort: number;
  status: number;
}

export interface Guarantee {
  guarantee_id: number;
  name: string;
  content: string;
  sort: number;
  status: number;
}

export interface AttrTemplate {
  template_id: number;
  template_name: string;
  template_value: string;
  sort: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchLabels(params: Record<string, unknown>) {
  return http.get<PageResult<ProductLabel>>('/product/labels', { params });
}
export function createLabel(data: { name: string; info?: string; sort?: number; status?: number }) {
  return http.post<ProductLabel>('/product/labels', data);
}
export function updateLabel(
  id: number,
  data: { name: string; info?: string; sort?: number; status?: number },
) {
  return http.put<ProductLabel>(`/product/labels/${id}`, data);
}
export function deleteLabel(id: number) {
  return http.delete<{ ok: boolean }>(`/product/labels/${id}`);
}

export function fetchGuarantees(params: Record<string, unknown>) {
  return http.get<PageResult<Guarantee>>('/product/guarantees', { params });
}
export function createGuarantee(data: {
  name: string;
  content?: string;
  sort?: number;
  status?: number;
}) {
  return http.post<Guarantee>('/product/guarantees', data);
}
export function updateGuarantee(
  id: number,
  data: { name: string; content?: string; sort?: number; status?: number },
) {
  return http.put<Guarantee>(`/product/guarantees/${id}`, data);
}
export function deleteGuarantee(id: number) {
  return http.delete<{ ok: boolean }>(`/product/guarantees/${id}`);
}

export function fetchAttrTemplates(params: Record<string, unknown>) {
  return http.get<PageResult<AttrTemplate>>('/product/attr-templates', { params });
}
export function createAttrTemplate(data: {
  template_name: string;
  template_value?: string;
  sort?: number;
}) {
  return http.post<AttrTemplate>('/product/attr-templates', data);
}
export function updateAttrTemplate(
  id: number,
  data: { template_name: string; template_value?: string; sort?: number },
) {
  return http.put<AttrTemplate>(`/product/attr-templates/${id}`, data);
}
export function deleteAttrTemplate(id: number) {
  return http.delete<{ ok: boolean }>(`/product/attr-templates/${id}`);
}
