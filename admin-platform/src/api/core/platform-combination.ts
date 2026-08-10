import { requestClient } from '#/api/request';

export interface PlatformCombination {
  action_status?: number;
  buying_count_num: number;
  create_time: string;
  end_time: string;
  image?: string;
  is_show?: number;
  mer_id: number;
  mer_name?: string;
  ot_price?: number;
  price: number;
  product_group_id: number;
  product_id: number;
  product_status?: number;
  refusal?: string;
  start_time: string;
  status: number;
  store_name?: string;
  success_num?: number;
  time?: number;
}

export interface PlatformCombinationPage {
  limit: number;
  list: PlatformCombination[];
  page: number;
  total: number;
}

export interface PlatformCombinationInput {
  price?: number;
  buying_count_num?: number;
  time?: number;
  start_time?: string;
  end_time?: string;
  is_show?: number;
  status?: number;
  product_status?: number;
  refusal?: string;
}

export function listPlatformCombinationsApi(params: {
  limit: number;
  page: number;
  mer_id?: number;
  keyword?: string;
  status?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PlatformCombinationPage>('/combination/groups', {
    params,
  });
}

export function getPlatformCombinationApi(id: number) {
  return requestClient.get<PlatformCombination & { time: number; is_show: number }>(
    '/combination/groups/' + id,
  );
}

export function updatePlatformCombinationApi(
  id: number,
  payload: PlatformCombinationInput,
) {
  return requestClient.put(`/combination/groups/${id}`, payload);
}

export function deletePlatformCombinationApi(id: number) {
  return requestClient.delete<{ ok: true }>('/combination/groups/' + id);
}

export function auditPlatformCombinationApi(
  id: number,
  status: number,
  refusal?: string,
) {
  return requestClient.post<PlatformCombination>('/combination/groups/audit', {
    id,
    status,
    refusal: refusal || '',
  });
}

export function forceOffPlatformCombinationApi(ids: number[], reason: string) {
  return requestClient.post<{ ok: true }>('/combination/groups/force-off', {
    ids,
    reason,
  });
}

/** 用户开团记录（CRMEB StoreProductGroupBuying / 拼团活动列表） */
export interface PlatformCombinationBuying {
  buying_count_num: number;
  create_time: string;
  end_time: number;
  group_buying_id: number;
  image?: string;
  is_trader?: number;
  members?: PlatformCombinationMember[];
  mer_id: number;
  mer_name?: string;
  nickname?: string;
  price?: number;
  product_group_id: number;
  product_id?: number;
  remain?: number;
  status: number;
  status_text?: string;
  stop_time?: string;
  store_name?: string;
  trader_name?: string;
  uid?: number;
  yet_buying_num: number;
}

export interface PlatformCombinationMember {
  avatar?: string;
  create_time?: string;
  id: number;
  is_initiator?: number;
  is_leader?: number;
  nickname: string;
  order_id?: number;
  status?: number;
  uid: number;
}

export interface PlatformCombinationBuyingPage {
  limit: number;
  list: PlatformCombinationBuying[];
  page: number;
  total: number;
}

export function listPlatformCombinationBuyingsApi(params: {
  limit: number;
  page: number;
  mer_id?: number;
  keyword?: string;
  user_name?: string;
  is_trader?: number;
  status?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PlatformCombinationBuyingPage>(
    '/combination/buyings',
    { params },
  );
}

export function getPlatformCombinationBuyingApi(id: number) {
  return requestClient.get<PlatformCombinationBuying>(
    `/combination/buyings/${id}`,
  );
}
