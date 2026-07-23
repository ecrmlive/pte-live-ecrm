import { http } from '@/api/http';

export interface Region {
  city_ids?: string;
  first?: number;
  first_price?: number;
  continue?: number;
  continue_price?: number;
}

export interface ShippingTemplate {
  template_id: number;
  name: string;
  type: number;
  appoint: number;
  sort: number;
  regions?: Region[];
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchShippingTemplates(params: Record<string, unknown>) {
  return http.get<PageResult<ShippingTemplate>>('/shipping/templates', { params });
}

export function getShippingTemplate(id: number) {
  return http.get<ShippingTemplate>(`/shipping/templates/${id}`);
}

export function createShippingTemplate(data: {
  name: string;
  type?: number;
  appoint?: number;
  sort?: number;
  regions?: Region[];
}) {
  return http.post<ShippingTemplate>('/shipping/templates', data);
}

export function updateShippingTemplate(
  id: number,
  data: { name: string; type?: number; appoint?: number; sort?: number; regions?: Region[] },
) {
  return http.put<ShippingTemplate>(`/shipping/templates/${id}`, data);
}

export function deleteShippingTemplate(id: number) {
  return http.delete<{ ok: boolean }>(`/shipping/templates/${id}`);
}
