import { requestClient } from '#/api/request';

export type MarketingDecorType = 'application' | 'atmosphere' | 'border' | 'topic';

export interface MarketingDecor {
  code: string;
  cover_url: string;
  decor_type: MarketingDecorType;
  ends_at: string;
  id: number;
  name: string;
  payload: Record<string, unknown>;
  remark: string;
  sort: number;
  starts_at: string;
  status: number;
  updated_at: string;
}

export interface MarketingDecorPage {
  limit: number;
  list: MarketingDecor[];
  page: number;
  total: number;
}

export interface MarketingDecorInput {
  code?: string;
  cover_url?: string;
  ends_at?: string;
  name: string;
  payload?: Record<string, unknown>;
  remark?: string;
  sort?: number;
  starts_at?: string;
  status?: number;
}

const pathByType: Record<MarketingDecorType, string> = {
  application: '/marketing/applications',
  atmosphere: '/marketing/atmosphere',
  border: '/marketing/border',
  topic: '/marketing/topic',
};

export function listMarketingDecorApi(
  type: MarketingDecorType,
  params: { keyword?: string; limit: number; page: number; status?: number },
) {
  return requestClient.get<MarketingDecorPage>(pathByType[type], { params });
}

export function getMarketingDecorApi(type: MarketingDecorType, id: number) {
  return requestClient.get<MarketingDecor>(`${pathByType[type]}/${id}`);
}

export function createMarketingDecorApi(type: MarketingDecorType, body: MarketingDecorInput) {
  return requestClient.post<MarketingDecor>(pathByType[type], body);
}

export function updateMarketingDecorApi(
  type: MarketingDecorType,
  id: number,
  body: MarketingDecorInput,
) {
  return requestClient.put<MarketingDecor>(`${pathByType[type]}/${id}`, body);
}

export function setMarketingDecorStatusApi(type: MarketingDecorType, id: number, status: 0 | 1) {
  return requestClient.put<{ id: number; status: number }>(`${pathByType[type]}/${id}/status`, {
    status,
  });
}

export function deleteMarketingDecorApi(type: MarketingDecorType, id: number) {
  return requestClient.delete<{ ok: boolean }>(`${pathByType[type]}/${id}`);
}
