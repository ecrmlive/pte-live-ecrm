import { liveAdminPost } from '#/api/live-request';
import { requestClient } from '#/api/request';

export interface DashboardMetric {
  last_week?: number;
  month: number;
  today: number;
  week_ratio?: number;
  yesterday: number;
}

export interface StoreSalesRankRow {
  follower_count: number;
  sale_amount: number;
  sale_count: number;
  store_id: number;
  store_image?: string;
  store_name: string;
}

export interface HourAmountPoint {
  hour: string;
  today_amount: number;
  yesterday_amount: number;
}

export interface SparkStat {
  ratio: number;
  spark: number[];
  spark_label?: string;
  value: number;
}

export interface OrderStatsBlock {
  month_order_count: SparkStat;
  month_payer_count: SparkStat;
  today_order_count: SparkStat;
  today_payer_count: SparkStat;
}

export interface UserTrendPoint {
  day: string;
  new_users: number;
  total_users: number;
  visit_users: number;
}

export interface DealFunnelBlock {
  avg_order_amount: number;
  order_amount: number;
  order_pay_rate: number;
  order_users: number;
  pay_amount: number;
  pay_users: number;
  visit_order_rate: number;
  visit_users: number;
}

export interface DealRatioBlock {
  new_amount: number;
  new_users: number;
  old_amount: number;
  old_users: number;
}

export type DashboardPeriod = '7d' | '30d' | 'month' | 'year';

export interface PlatformDashboardSummary {
  /** all: 平台全量；store: 当前账号获授权店铺范围 */
  scope: 'all' | 'store';
  deal_funnel: DealFunnelBlock;
  deal_ratio: DealRatioBlock;
  new_users: DashboardMetric;
  on_sale_product: number;
  order_stats: OrderStatsBlock;
  page_views: DashboardMetric;
  paid_order: number;
  pending_community: number;
  pending_delivery: number;
  pending_feedback: number;
  pending_integral_ship: number;
  pending_product_audit: number;
  pending_refund: number;
  pending_service: number;
  pending_spread_gift: number;
  pending_store_audit: number;
  pending_transfer: number;
  pending_withdraw: number;
  store_count: number;
  store_sales_rank: StoreSalesRankRow[];
  store_total: number;
  stores: DashboardMetric;
  today_order_count: number;
  today_order_hours: HourAmountPoint[];
  today_paid_amount: number;
  today_payer_count: number;
  user_trend: UserTrendPoint[];
  visitors: DashboardMetric;
}

function emptySpark(): SparkStat {
  return { ratio: 0, spark: [], value: 0 };
}

export function emptyDashboardSummary(): PlatformDashboardSummary {
  const metric = (): DashboardMetric => ({ month: 0, today: 0, week_ratio: 0, yesterday: 0 });
  return {
    deal_funnel: {
      avg_order_amount: 0,
      order_amount: 0,
      order_pay_rate: 0,
      order_users: 0,
      pay_amount: 0,
      pay_users: 0,
      visit_order_rate: 0,
      visit_users: 0,
    },
    deal_ratio: { new_amount: 0, new_users: 0, old_amount: 0, old_users: 0 },
    new_users: metric(),
    on_sale_product: 0,
    order_stats: {
      month_order_count: emptySpark(),
      month_payer_count: emptySpark(),
      today_order_count: emptySpark(),
      today_payer_count: emptySpark(),
    },
    page_views: metric(),
    paid_order: 0,
    pending_community: 0,
    pending_delivery: 0,
    pending_feedback: 0,
    pending_integral_ship: 0,
    pending_product_audit: 0,
    pending_refund: 0,
    pending_service: 0,
    pending_spread_gift: 0,
    pending_store_audit: 0,
    pending_transfer: 0,
    pending_withdraw: 0,
    scope: 'all',
    store_count: 0,
    store_sales_rank: [],
    store_total: 0,
    stores: metric(),
    today_order_count: 0,
    today_order_hours: [],
    today_paid_amount: 0,
    today_payer_count: 0,
    user_trend: [],
    visitors: metric(),
  };
}

export function getPlatformDashboardSummaryApi() {
  return requestClient.get<PlatformDashboardSummary>('/dashboard/summary');
}

export function getPlatformMerchantTopApi(period: DashboardPeriod = 'month') {
  return requestClient.get<{ list: StoreSalesRankRow[]; period: string }>('/dashboard/merchant-top', {
    params: { period },
  });
}

export function getPlatformUserTrendApi(period: Exclude<DashboardPeriod, 'year'> = '30d') {
  return requestClient.get<{ list: UserTrendPoint[]; period: string }>('/dashboard/user-trend', {
    params: { period },
  });
}

export function getPlatformDealApi(period: DashboardPeriod = 'month') {
  return requestClient.get<{ funnel: DealFunnelBlock; period: string; ratio: DealRatioBlock }>(
    '/dashboard/deal',
    { params: { period } },
  );
}

export interface DataScreenTodayNumbers {
  today_pay_number: number;
  today_pay_user_first: number;
  visit_num: number;
  visit_user_num: number;
}

export interface DataScreenNewOld {
  new_count: number;
  old_count: number;
}

export interface DataScreenPaymentAmount {
  count: number;
  number: number;
  order_id: number;
  paid: number;
}

export interface DataScreenCityRank {
  code?: string;
  name: string;
  value: number;
}

export interface DataScreenMonthPoint {
  day: string;
  total_sum: number;
}

export interface DataScreenHourPoint {
  hours: string;
  order_count: number;
  user_count: number;
}

export interface DataScreenOrderInfo {
  number: number;
  payment_method: string;
  paytime: string;
  product: DataScreenProduct;
  store: DataScreenStore;
}

export interface DataScreenProduct {
  image?: string;
  product_name: string;
}

export interface DataScreenStore {
  image?: string;
  store_name: string;
}

export interface DataScreenMerchantRank {
  count: number;
  number: number;
  store: DataScreenStore;
}

export interface DataScreenMerchantRankBoard {
  data: DataScreenMerchantRank[];
  type: string;
}

export interface DataScreenProductRank {
  count: number;
  number: number;
  product: DataScreenProduct;
}

export interface PlatformDataScreen {
  city_ranking: DataScreenCityRank[];
  config: { data_screen_title: string };
  month_pay_count: DataScreenMonthPoint[];
  pay_product_rank: DataScreenProductRank[];
  today_pay_count: DataScreenHourPoint[];
  today_pay_count_number: DataScreenTodayNumbers;
  today_pay_info: DataScreenOrderInfo[];
  today_pay_merchant_rank: DataScreenMerchantRankBoard;
  today_pay_new_old: DataScreenNewOld;
  today_pay_number: DataScreenPaymentAmount;
}

export function getPlatformDataScreenApi() {
  return requestClient.get<PlatformDataScreen>('/dashboard/data-screen');
}

export interface PlatformDashboardOverview {
  merchant_count: number;
  recharge_traffic_gb: number;
  recharge_amount_yuan: number;
  consumed_traffic_gb: number;
}

export interface PlatformTrafficSummary {
  total_gb: number;
  lvb_play_used_gb: number;
  vod_play_used_gb: number;
  remain_gb: number;
}

export interface PlatformSalesSummary {
  user_count: number;
  paid_order_count: number;
  paid_amount: number;
  refund_order_count: number;
  refund_amount: number;
}

export interface ProductRankRow {
  app_id: number;
  app_name: string;
  app_logo?: string;
  product_id: number;
  product_image?: string;
  product_name: string;
  product_price?: number;
  rank_key?: string;
  total_num: number;
  total_price: number;
}

export interface UserRegionRow {
  province: string;
  user_count: number;
}

export interface PlatformDashboardData {
  overview: PlatformDashboardOverview;
  traffic_summary: PlatformTrafficSummary;
  sales_summary: PlatformSalesSummary;
  product_rank: ProductRankRow[];
  range?: string;
  user_region: UserRegionRow[];
  update_time: string;
  cache_time?: string;
}

export interface PlatformMyStatsData {
  merchant_count: number;
  recharge_traffic_gb: number;
  recharge_amount_yuan: number;
  lvb_play_used_gb: number;
  vod_play_used_gb: number;
  update_time: string;
}

export interface PlatformTencentLvbTraffic {
  today_play_gb: number;
  month_play_gb: number;
  package_total_gb: number;
  package_used_gb: number;
  package_remain_gb: number;
}

export interface PlatformTencentVodTraffic {
  today_play_gb: number;
  month_play_gb: number;
}

export interface PlatformTencentTrafficData {
  configured: boolean;
  lvb: PlatformTencentLvbTraffic;
  vod: PlatformTencentVodTraffic;
  update_time: string;
  cache_time?: string;
}

export async function getPlatformDashboardApi(
  productLimit = 10,
  range = '30d',
  dateRange?: { end_date?: string; start_date?: string },
) {
  return liveAdminPost<PlatformDashboardData>(
    '/api/v1/admin/platform/dashboard',
    { product_limit: productLimit, range, ...(dateRange ?? {}) },
  );
}

export async function getPlatformMyStatsApi() {
  return liveAdminPost<PlatformMyStatsData>('/api/v1/admin/platform/my-stats', {});
}

export async function getPlatformTencentTrafficApi(refresh = false) {
  return liveAdminPost<PlatformTencentTrafficData>(
    '/api/v1/admin/platform/tencent-traffic',
    refresh ? { refresh: 1 } : {},
  );
}

export function formatTrafficGB(value: null | number | undefined) {
  const n = Number(value ?? 0);
  return n.toFixed(n >= 100 ? 1 : 2);
}

export function formatMoney(value: null | number | undefined) {
  return Number(value ?? 0).toLocaleString('zh-CN', {
    maximumFractionDigits: 2,
    minimumFractionDigits: 2,
  });
}

export function formatCount(value: null | number | undefined) {
  return Number(value ?? 0).toLocaleString('zh-CN');
}
