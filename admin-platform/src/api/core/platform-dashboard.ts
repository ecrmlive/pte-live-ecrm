import { liveAdminPost } from '#/api/live-request';

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
