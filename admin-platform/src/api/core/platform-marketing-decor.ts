import { requestClient } from '#/api/request';

export type MarketingDecorType = 'application' | 'atmosphere' | 'border' | 'topic';

/** 对齐 CRMEB：0 全部 / 1 商品 / 2 分类 / 3 店铺 / 4 商品标签 */
export type MarketingDecorScopeType = 0 | 1 | 2 | 3 | 4;

export interface MarketingDecor {
  activity_status?: number;
  activity_status_text?: string;
  cate_ids?: number[];
  code: string;
  cover_url: string;
  created_at?: string;
  decor_type: MarketingDecorType;
  ends_at: string;
  id: number;
  label_ids?: number[];
  mer_ids?: number[];
  name: string;
  payload: Record<string, unknown>;
  remark: string;
  scope_type?: MarketingDecorScopeType | number;
  sort: number;
  spu_ids?: number[];
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
  cate_ids?: number[];
  code?: string;
  cover_url?: string;
  ends_at?: string;
  label_ids?: number[];
  mer_ids?: number[];
  name: string;
  payload?: Record<string, unknown>;
  remark?: string;
  scope_type?: MarketingDecorScopeType | number;
  sort?: number;
  spu_ids?: number[];
  starts_at?: string;
  status?: number;
}

export interface MarketingDecorListParams {
  activity_status?: number;
  date_from?: string;
  date_to?: string;
  keyword?: string;
  limit: number;
  page: number;
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
  params: MarketingDecorListParams,
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
