import { requestClient } from '#/api/request';

export interface PlatformSeckillActive {
  seckill_active_id: number;
  name: string;
  seckill_time_ids?: string;
  start_day: string;
  end_day: string;
  mer_id: number;
  mer_name?: string;
  is_trader?: number;
  trader_name?: string;
  product_id: number;
  store_name?: string;
  image?: string;
  price?: number;
  seckill_price: number;
  once_pay_count?: number;
  all_pay_count?: number;
  active_status: number;
  status: number;
  is_show: number;
  product_status: number;
  product_status_name?: string;
  activity_status?: number;
  activity_status_text?: string;
  star: number;
  sort: number;
  stock: number;
  sales: number;
  sys_labels?: string;
  refusal?: string;
  time_titles?: string;
}

export interface PlatformSeckillPage {
  limit: number;
  list: PlatformSeckillActive[];
  page: number;
  total: number;
}

export interface PlatformSeckillInput {
  name?: string;
  seckill_time_ids?: string;
  start_day?: string;
  end_day?: string;
  seckill_price?: number;
  once_pay_count?: number;
  all_pay_count?: number;
  status?: number;
  is_show?: number;
  product_status?: number;
  star?: number;
  sort?: number;
  stock?: number;
  sys_labels?: string;
  refusal?: string;
}

export interface PlatformSeckillListParams {
  page: number;
  limit: number;
  type?: number;
  mer_id?: number;
  active_name?: string;
  keyword?: string;
  is_trader?: number;
  star?: number;
  us_status?: number;
  sys_labels?: string;
  active_status?: number;
}

export interface PlatformSeckillStatusFilter {
  type: number;
  name: string;
  count: number;
}

export function listPlatformSeckillApi(params: PlatformSeckillListParams) {
  return requestClient.get<PlatformSeckillPage>('/seckill/actives', { params });
}

export function getPlatformSeckillStatusFilterApi(
  params: Omit<PlatformSeckillListParams, 'page' | 'limit' | 'type'>,
) {
  return requestClient.get<{ list: PlatformSeckillStatusFilter[] }>(
    '/seckill/actives/filter',
    { params },
  );
}

export function getPlatformSeckillApi(id: number) {
  return requestClient.get<PlatformSeckillActive>('/seckill/actives/' + id);
}

export function updatePlatformSeckillApi(id: number, payload: PlatformSeckillInput) {
  return requestClient.put(`/seckill/actives/${id}`, payload);
}

export function setPlatformSeckillShowApi(id: number, isShow: number) {
  return requestClient.put(`/seckill/actives/${id}/show`, { is_show: isShow });
}

export function setPlatformSeckillStarApi(id: number, star: number) {
  return requestClient.put(`/seckill/actives/${id}/star`, { star });
}

export function setPlatformSeckillLabelsApi(id: number, sysLabels: string) {
  return requestClient.put(`/seckill/actives/${id}/labels`, {
    sys_labels: sysLabels,
  });
}

export function forceOffPlatformSeckillApi(ids: number[], reason: string) {
  return requestClient.post<{ ok: true }>('/seckill/actives/force-off', {
    ids,
    reason,
  });
}

export function deletePlatformSeckillApi(id: number) {
  return requestClient.delete<{ ok: true }>('/seckill/actives/' + id);
}

/** 秒杀配置（场次） */
export interface PlatformSeckillTime {
  seckill_time_id: number;
  title: string;
  start_time: number;
  end_time: number;
  status: number;
  pic: string;
  create_time?: string;
}

export interface PlatformSeckillTimePage {
  limit: number;
  list: PlatformSeckillTime[];
  page: number;
  total: number;
}

export interface PlatformSeckillTimeInput {
  title: string;
  start_time: number;
  end_time: number;
  status?: number;
  pic?: string;
}

export function listPlatformSeckillTimesApi(params: {
  page: number;
  limit: number;
  status?: number;
}) {
  return requestClient.get<PlatformSeckillTimePage>('/seckill/times', {
    params,
  });
}

export function createPlatformSeckillTimeApi(payload: PlatformSeckillTimeInput) {
  return requestClient.post<PlatformSeckillTime>('/seckill/times', payload);
}

export function updatePlatformSeckillTimeApi(
  id: number,
  payload: PlatformSeckillTimeInput,
) {
  return requestClient.put<PlatformSeckillTime>(`/seckill/times/${id}`, payload);
}

export function setPlatformSeckillTimeStatusApi(id: number, status: number) {
  return requestClient.put<PlatformSeckillTime>(`/seckill/times/${id}/status`, {
    status,
  });
}

export function deletePlatformSeckillTimeApi(id: number) {
  return requestClient.delete<{ ok: true }>(`/seckill/times/${id}`);
}

/** 启用中的场次（活动表单选择器用，无分页） */
export function listPlatformSeckillTimeOptionsApi() {
  return requestClient.get<{ list: PlatformSeckillTime[] }>('/seckill/times');
}

/** 秒杀活动（活动场，对齐 CRMEB store_seckill_active） */
export interface PlatformSeckillActivity {
  seckill_activity_id: number;
  name: string;
  seckill_time_ids: string;
  start_day: string;
  end_day: string;
  once_pay_count: number;
  all_pay_count: number;
  product_category_ids?: string;
  border_pic?: string;
  status: number;
  active_status: number;
  status_text?: string;
  product_count: number;
  merchant_count: number;
  create_time?: string;
  update_time?: string;
  seckill_time_texts?: string[];
}

export interface PlatformSeckillActivityPage {
  limit: number;
  list: PlatformSeckillActivity[];
  page: number;
  total: number;
}

export interface PlatformSeckillActivityInput {
  name: string;
  seckill_time_ids: string;
  start_day: string;
  end_day: string;
  once_pay_count?: number;
  all_pay_count?: number;
  product_category_ids?: string;
  border_pic?: string;
  status?: number;
}

export interface PlatformSeckillActivityListParams {
  page: number;
  limit: number;
  name?: string;
  date_from?: string;
  date_to?: string;
  active_status?: number;
  status?: number;
}

export function listPlatformSeckillActivitiesApi(
  params: PlatformSeckillActivityListParams,
) {
  return requestClient.get<PlatformSeckillActivityPage>('/seckill/activities', {
    params,
  });
}

export function getPlatformSeckillActivityApi(id: number) {
  return requestClient.get<PlatformSeckillActivity>(
    `/seckill/activities/${id}`,
  );
}

export function createPlatformSeckillActivityApi(
  payload: PlatformSeckillActivityInput,
) {
  return requestClient.post<PlatformSeckillActivity>(
    '/seckill/activities',
    payload,
  );
}

export function updatePlatformSeckillActivityApi(
  id: number,
  payload: PlatformSeckillActivityInput,
) {
  return requestClient.put<PlatformSeckillActivity>(
    `/seckill/activities/${id}`,
    payload,
  );
}

export function setPlatformSeckillActivityStatusApi(
  id: number,
  status: number,
) {
  return requestClient.put<PlatformSeckillActivity>(
    `/seckill/activities/${id}/status`,
    { status },
  );
}

export function clonePlatformSeckillActivityApi(id: number) {
  return requestClient.post<PlatformSeckillActivity>(
    `/seckill/activities/${id}/clone`,
  );
}

export function deletePlatformSeckillActivityApi(id: number) {
  return requestClient.delete<{ ok: true }>(`/seckill/activities/${id}`);
}

export interface PlatformSeckillActivityProduct {
  seckill_active_id: number;
  product_id: number;
  name: string;
  image?: string;
  category_name?: string;
  mer_id: number;
  mer_name?: string;
  price: number;
  seckill_price: number;
  stock: number;
  sales: number;
  seckill_time_texts?: string[];
}

export interface PlatformSeckillActivityStats {
  seckill_activity_id: number;
  name: string;
  orders_people_count: number;
  pay_order_money: number;
  pay_order_people_count: number;
  pay_order_count: number;
}

export interface PlatformSeckillActivityStatPeople {
  uid: number;
  nickname: string;
  phone?: string;
  mer_id: number;
  sum_total_num: number;
  order_count: number;
  sum_pay_price: number;
  last_join_time: string;
}

export interface PlatformSeckillActivityStatOrder {
  order_sn: string;
  uid: number;
  nickname: string;
  mer_id: number;
  status: number;
  status_text: string;
  pay_price: number;
  total_num: number;
  paid: number;
  create_time: string;
  pay_time?: string | null;
}

export interface PlatformSeckillActivityStatPage<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export interface PlatformSeckillActivityStatParams {
  page: number;
  limit: number;
  keyword?: string;
  date_from?: string;
  date_to?: string;
  mer_id?: number;
  status?: number;
}

export function getPlatformSeckillActivityStatsApi(
  id: number,
  params?: { mer_id?: number },
) {
  return requestClient.get<PlatformSeckillActivityStats>(
    `/seckill/activities/${id}/stats`,
    { params },
  );
}

export function listPlatformSeckillActivityStatPeopleApi(
  id: number,
  params: PlatformSeckillActivityStatParams,
) {
  return requestClient.get<
    PlatformSeckillActivityStatPage<PlatformSeckillActivityStatPeople>
  >(`/seckill/activities/${id}/stats/people`, { params });
}

export function listPlatformSeckillActivityStatOrdersApi(
  id: number,
  params: PlatformSeckillActivityStatParams,
) {
  return requestClient.get<
    PlatformSeckillActivityStatPage<PlatformSeckillActivityStatOrder>
  >(`/seckill/activities/${id}/stats/orders`, { params });
}

export function listPlatformSeckillActivityStatProductsApi(
  id: number,
  params: PlatformSeckillActivityStatParams,
) {
  return requestClient.get<
    PlatformSeckillActivityStatPage<PlatformSeckillActivityProduct>
  >(`/seckill/activities/${id}/stats/products`, { params });
}

/** 活动「已加商品」行（含降级 SKU children） */
export interface PlatformSeckillActivityGoodsSKU {
  sku: string;
  image?: string;
  price?: number;
  seckill_price: number;
  stock: number;
  limit_stock: number;
}

export interface PlatformSeckillActivityGoods {
  seckill_active_id: number;
  product_id: number;
  name: string;
  store_name: string;
  image?: string;
  cate_name?: string;
  mer_id: number;
  mer_name?: string;
  price: number;
  seckill_price: number;
  product_stock: number;
  stock: number;
  sort: number;
  product_status: number;
  product_status_name?: string;
  children?: PlatformSeckillActivityGoodsSKU[];
}

export interface PlatformSeckillActivityGoodsPage {
  list: PlatformSeckillActivityGoods[];
  total: number;
  page: number;
  limit: number;
}

export interface PlatformSeckillActivityGoodsParams {
  page: number;
  limit: number;
  keyword?: string;
  product_status?: number;
}

export function listPlatformSeckillActivityGoodsApi(
  id: number,
  params: PlatformSeckillActivityGoodsParams,
) {
  return requestClient.get<PlatformSeckillActivityGoodsPage>(
    `/seckill/activities/${id}/products`,
    { params },
  );
}

/** 新增/编辑 Drawer「秒杀商品」草稿落库 */
export interface PlatformSeckillActivityProductSaveItem {
  product_id: number;
  seckill_price: number;
  stock: number;
  status?: number;
  sort?: number;
}

export function savePlatformSeckillActivityProductsApi(
  id: number,
  products: PlatformSeckillActivityProductSaveItem[],
) {
  return requestClient.post<{ ok: true }>(`/seckill/activities/${id}/products`, {
    products,
  });
}
