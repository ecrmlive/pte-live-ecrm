/**
 * 商户直播 API — 统一走 api-platform (:11503)。
 * - 运行时：`/api/v1/shop/live/*`（api-platform 进程内 Go 路由）
 * - 微信带货库：`plus-live-wx-product.ts` → `/shop/plus.live.product/*`（api-platform Go）
 * - 下列 legacy 路径亦为 api-platform Go（非 PHP）：auto_syn / 手动同步 / 分享二维码
 */
import { requestClient } from '#/api/request';
import LiveApiClient from '#/api/core/live-api-client';

import type { PaginatedList } from './product';

export type LiveTrafficRange = '30d' | '7d' | 'today' | 'yesterday';

export interface LiveTrafficSummary {
  app_id: number;
  app_name: string;
  block_threshold_gb: number;
  lvb_play_used_gb: number;
  push_flux_gb_total: number;
  remain_gb: number;
  total_gb: number;
  used_gb_total: number;
  vod_play_used_gb: number;
  warn_threshold_gb: number;
}

export interface LiveTrafficDailyItem {
  play_flux_gb?: number;
  play_gb?: number;
  push_flux_gb?: number;
  push_gb?: number;
  stat_date: string;
}

export interface LiveTrafficRechargeItem {
  amount_yuan: number | string;
  create_time_text: string;
  delta_gb: number;
  invoice_no?: string;
  operator_name?: string;
  recharge_type: number;
  remark?: string;
  total_gb_after: number;
}

export interface LiveRoomListItem {
  allow_chat?: number;
  anchor_id?: number;
  anchor_name?: string;
  anchor_wechat?: string;
  background_img?: string;
  cover_img?: string;
  create_time_text?: string;
  enable_gift?: number;
  enable_linkmic?: number;
  enable_replay?: number;
  end_time_text?: string;
  feeds_img?: string;
  fire_value?: number;
  is_top?: number;
  is_visible?: number;
  live_id: number;
  live_status?: number;
  name: string;
  record_video_path?: string;
  record_vod_file_id?: string;
  record_vod_media_url?: string;
  room_type?: number;
  roomid?: number | string;
  session_status?: number;
  share_img?: string;
  share_intro?: string;
  show_heat?: number;
  show_home?: number;
  show_online_count?: number;
  show_total_count?: number;
  start_time_text?: string;
  stream_orientation?: number;
  stream_status?: number;
  system_notice?: string;
  watch_password?: string;
  [key: string]: unknown;
}

export interface LiveRoomListQuery {
  anchor_name?: string;
  list_rows?: number;
  name?: string;
  page?: number;
  room_type?: number;
  search?: string;
  session_status?: number;
}

export interface LiveRoomListResult {
  list: PaginatedList<LiveRoomListItem>;
}

export interface LiveImUserSigResult {
  app_id: string;
  auth_mode: string;
  device_id?: string;
  expire_at?: number;
  identifier: string;
  im_token: string;
  platform?: string;
  room_id?: string;
  sdk_app_id: string;
  token: string;
  user_id: string;
  user_sig: string;
  ws_url: string;
}

export async function getLiveTrafficSummaryApi() {
  return requestClient.post<LiveTrafficSummary>(
    '/api/v1/shop/traffic/summary',
    {},
  );
}

export async function getLiveTrafficDailyLvbApi(params: { range: LiveTrafficRange }) {
  return requestClient.post<{ list: LiveTrafficDailyItem[] }>(
    '/api/v1/shop/traffic/daily/lvb',
    params,
  );
}

export async function getLiveTrafficDailyVodApi(params: { range: LiveTrafficRange }) {
  return requestClient.post<{ list: LiveTrafficDailyItem[] }>(
    '/api/v1/shop/traffic/daily/vod',
    params,
  );
}

export async function getLiveTrafficRechargeListApi(params: {
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<{ list: PaginatedList<LiveTrafficRechargeItem> }>(
    '/api/v1/shop/traffic/recharge/list',
    params,
  );
}

export async function getLiveRoomListApi(params: LiveRoomListQuery) {
  return requestClient.post<LiveRoomListResult>(
    '/api/v1/shop/live/list',
    params,
  );
}

export async function getLiveImUserSigApi(params: {
  device_id?: string;
  live_id?: number | string;
  platform?: string;
  room_id?: number | string;
}) {
  return requestClient.post<LiveImUserSigResult>(
    '/api/v1/shop/live/im/usersig',
    params,
  );
}

export async function createLiveRoomApi(payload: Record<string, unknown>) {
  return requestClient.post<{ data?: { push_url?: string }; msg?: string }>(
    '/api/v1/shop/live/create',
    payload,
  );
}

export async function updateLiveRoomApi(payload: Record<string, unknown>) {
  return requestClient.post<{ msg?: string }>(
    '/api/v1/shop/live/update',
    payload,
  );
}

export async function deleteLiveRoomApi(liveId: number) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/delete', {
    live_id: liveId,
  });
}

export async function endLiveRoomApi(liveId: number) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/end', {
    live_id: liveId,
  });
}

export async function setLiveRoomTopApi(payload: { is_top: number; live_id: number }) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/set-top', payload);
}

export async function syncLiveRoomsApi() {
  return requestClient.post<{ msg?: string }>('/shop/plus.live.wx/syn', {});
}

export async function setLiveAutoSynApi(payload: { auto_syn: boolean }) {
  return requestClient.post<{ msg?: string }>('/shop/plus.live.wx/setSyn', payload);
}

export async function getLiveAutoSynSettingApi() {
  return requestClient.post<{ auto_syn: boolean }>('/shop/plus.live.wx/index', {});
}

export interface LiveAnchorListItem {
  account: string;
  anchor_id: number;
  avatar?: string;
  create_time_text?: string;
  intro?: string;
  nick_name: string;
  phone: string;
  sort: number;
  status: number;
  wechat?: string;
}

export interface LiveAnchorListQuery {
  list_rows?: number;
  page?: number;
  search?: string;
  status?: number | string;
}

export async function getLiveAnchorListApi(params: LiveAnchorListQuery) {
  return requestClient.post<{ list: PaginatedList<LiveAnchorListItem> }>(
    '/api/v1/shop/live/anchor/list',
    params,
  );
}

export async function deleteLiveAnchorApi(anchorId: number) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/anchor/delete', {
    anchor_id: anchorId,
  });
}

export interface LiveAnchorForm {
  account?: string;
  anchor_id?: number;
  avatar?: string;
  intro?: string;
  nick_name: string;
  password?: string;
  phone: string;
  sort?: number;
  status?: number;
  wechat?: string;
}

export async function createLiveAnchorApi(payload: LiveAnchorForm) {
  return requestClient.post<{ data?: LiveAnchorListItem; msg?: string }>(
    '/api/v1/shop/live/anchor/create',
    payload,
  );
}

export async function updateLiveAnchorApi(payload: LiveAnchorForm & { anchor_id: number }) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/anchor/update', payload);
}

export interface LiveRoomForm {
  allow_chat?: number;
  anchor_id?: number;
  anchor_name: string;
  anchor_wechat: string;
  background_img?: string;
  close_comment?: number;
  close_goods?: number;
  close_like?: number;
  close_replay?: number;
  cover_img: string;
  enable_gift?: number;
  enable_linkmic?: number;
  enable_replay?: number;
  end_time: string;
  feeds_img: string;
  fire_value?: number;
  is_visible?: number;
  live_id?: number;
  name: string;
  record_video_path?: string;
  record_vod_file_id?: string;
  record_vod_media_url?: string;
  room_type: number;
  share_img: string;
  share_intro?: string;
  show_heat?: number;
  show_home?: number;
  show_online_count?: number;
  show_total_count?: number;
  start_time: string;
  stream_orientation?: number;
  system_notice?: string;
  watch_password?: string;
}

export interface LiveStreamInfo {
  pull_url_flv?: string;
  pull_url_hls?: string;
  pull_url_rtmp?: string;
  pull_url_webrtc?: string;
  push_auth_query?: string;
  push_stream_key?: string;
  push_server_url?: string;
  push_url?: string;
  stream_status?: number;
}

function parseShopProductListPayload<T>(raw: unknown): { list: T[]; total: number } {
	const result = { list: [], total: 0 } as { list: T[]; total: number };
	if (!raw || typeof raw !== 'object') {
		return result;
	}
	const payload = raw as { list?: unknown; total?: unknown };
	const total = Number(payload.total);
	let hasTopTotal = Number.isFinite(total);
	if (hasTopTotal && total >= 0) {
		result.total = total;
	}

	const listValue = payload.list;
	if (Array.isArray(listValue)) {
		result.list = listValue as T[];
		return result;
	}
	if (listValue && typeof listValue === 'object') {
		const nested = listValue as { data?: unknown; total?: unknown };
		if (Array.isArray(nested.data)) {
			result.list = nested.data as T[];
			const nestedTotal = Number(nested.total);
			if (!hasTopTotal && Number.isFinite(nestedTotal) && nestedTotal >= 0) {
				result.total = Number(nested.total);
			}
		}
	}
	if (!hasTopTotal && !result.total && result.list.length > 0) {
		result.total = result.list.length;
	}
	return result;
}

async function fetchShopProductList<T>(params: {
	list_rows?: number;
	live_id: number;
	page?: number;
	room_id?: number | string;
}) {
	const res = await LiveApiClient.liveProductList(params as Record<string, unknown>);
	return parseShopProductListPayload<T>(res.data);
}

export async function getLiveStreamApi(liveId: number) {
  return requestClient.post<LiveStreamInfo>('/api/v1/shop/live/stream', {
    live_id: liveId,
  });
}

export async function refreshLiveStreamUrlApi(liveId: number) {
  return requestClient.post<LiveStreamInfo & { msg?: string }>(
    '/api/v1/shop/live/stream/refresh-url',
    { live_id: liveId },
  );
}

export interface LiveRoomProductItem {
  isPush?: number;
  live_product_id: number;
  name?: string;
  on_sale?: number;
  product?: { cover_img?: string; name?: string };
  [key: string]: unknown;
}

export async function getLiveRoomProductListApi(params: {
  list_rows?: number;
  live_id: number;
  page?: number;
  room_id: number | string;
}) {
  return fetchShopProductList<LiveRoomProductItem>({
    list_rows: params.list_rows,
    live_id: params.live_id,
    page: params.page,
    room_id: params.room_id,
  });
}

export async function addLiveRoomProductsApi(payload: {
  live_id: number;
  productIds: Array<number | string>;
  room_id: number | string;
}) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/product/add', payload);
}

export async function setLiveRoomProductOnSaleApi(payload: {
  live_product_id: number;
  on_sale: number;
  room_id?: number | string;
}) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/product/onsale', payload);
}

export async function pushLiveRoomProductApi(payload: {
  live_product_id: number;
  room_id?: number | string;
}) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/product/explain', payload);
}

export async function deleteLiveRoomProductApi(liveProductId: number) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/product/delete', {
    live_product_id: liveProductId,
  });
}

export interface LiveMallProductRow {
  activity_price?: number;
  cover_img?: string;
  limit_per_user?: number;
  live_product_id: number;
  name?: string;
  on_sale?: number;
  original_price?: number;
  price?: number;
  price2?: number;
  product_id?: number;
  product_stock?: number;
}

export async function listLiveMallProductsApi(params: {
  live_id: number;
  list_rows?: number;
  page?: number;
}) {
  return fetchShopProductList<LiveMallProductRow>({
    list_rows: params.list_rows,
    live_id: params.live_id,
    page: params.page,
  });
}

export async function addLiveMallProductsApi(payload: {
  live_id: number;
  product_ids: string;
}) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/product/add', payload);
}

export async function updateLiveMallProductApi(payload: {
  limit_per_user?: number;
  live_product_id: number;
}) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/product/update', payload);
}

export async function setLiveMallProductOnSaleApi(payload: {
  live_product_id: number;
  on_sale: number;
  shelf_qty?: number;
}) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/product/onsale', payload);
}

export async function deleteLiveMallProductApi(liveProductId: number) {
  return requestClient.post<{ msg?: string }>('/api/v1/shop/live/product/delete', {
    live_product_id: liveProductId,
  });
}

export interface LiveRoomShareH5Item {
  domain?: string;
  h5_qrcode?: string;
  h5_url?: string;
}

export interface LiveRoomShareQrcodeResult {
  h5_list?: LiveRoomShareH5Item[];
  h5_qrcode?: string;
  h5_url?: string;
  live_id?: number;
  mp_entry?: string;
  mp_path?: string;
  mp_scene?: string;
  room_id?: number;
  wx_error?: string;
  wx_qrcode?: string;
}

export async function getLiveRoomShareQrcodeApi(payload: {
  live_id: number;
  room_id: number | string;
}) {
  return requestClient.post<LiveRoomShareQrcodeResult>('/shop/plus.live.room/qrcode', payload);
}

export interface LiveAudienceInviter {
  avatar?: string;
  nick_name?: string;
  user_id?: number;
}

export interface LiveAudienceListItem {
  avatar?: string;
  enter_method_text?: string;
  first_enter_time_text?: string;
  first_inviter?: LiveAudienceInviter;
  last_active_time_text?: string;
  nick_name?: string;
  session_inviter?: LiveAudienceInviter;
  user_id?: number;
  watch_duration_text?: string;
}

export interface LiveAudienceListQuery {
  enter_time_end?: string;
  enter_time_start?: string;
  inviter_id?: number | string;
  inviter_nick_name?: string;
  last_active_end?: string;
  last_active_start?: string;
  list_rows?: number;
  live_id: number;
  nick_name?: string;
  page?: number;
  user_id?: number | string;
}

export interface LiveAudienceListResult {
  list: PaginatedList<LiveAudienceListItem>;
  summary: {
    anonymous_count: number;
    logged_in_count: number;
  };
}

export interface LiveInviteStatsItem {
  inviter_avatar?: string;
  inviter_id?: number;
  inviter_nick_name?: string;
  invite_user_count?: number;
  invite_watch_duration_text?: string;
  order_count?: number;
  order_pay_amount_text?: string;
}

export interface LiveInviteStatsQuery {
  list_rows?: number;
  live_id: number;
  page?: number;
}

export async function getLiveAudienceListApi(params: LiveAudienceListQuery) {
  return requestClient.post<LiveAudienceListResult>(
    '/api/v1/shop/live/audience/list',
    params,
  );
}

export async function getLiveInviteStatsApi(params: LiveInviteStatsQuery) {
  return requestClient.post<{ list: PaginatedList<LiveInviteStatsItem> }>(
    '/api/v1/shop/live/invite/stats',
    params,
  );
}

export interface LiveSessionChatItem {
  content?: string;
  message_id?: number | string;
  nick_name?: string;
  role?: number;
  role_text?: string;
  send_time_text?: string;
  source?: number;
  user_id?: number;
}

export interface LiveSessionChatListQuery {
  list_rows?: number;
  live_id: number;
  page?: number;
  session_id: string;
}

export interface LiveSessionChatListResult {
  list: LiveSessionChatItem[];
  total: number;
}

export async function getLiveSessionChatListApi(params: LiveSessionChatListQuery) {
  return requestClient.post<LiveSessionChatListResult>(
    '/api/v1/shop/live/session/chat/list',
    params,
  );
}

export function formatTrafficGB(value: null | number | undefined) {
  const n = Number(value ?? 0);
  return n.toFixed(n >= 100 ? 1 : 2);
}

export interface LiveStreamLogItem {
  cover_img?: string;
  disconnect_reason?: string;
  duration_text?: string;
  end_time?: number;
  end_time_text?: string;
  live_id: number;
  log_id: number;
  name?: string;
  push_client_ip?: string;
  session_id?: string;
  start_time?: number;
  start_time_text?: string;
}

export interface LiveStreamLogListQuery {
  list_rows?: number;
  live_id?: number;
  order?: 'asc' | 'desc';
  order_field?: string;
  page?: number;
  search?: string;
}

export interface LiveTrafficSessionItem {
  live_id: number;
  lvb_play_gb?: number;
  push_flux_gb?: number;
  quota_gb_total?: number;
  room_name?: string;
  session_end?: string;
  session_id?: string;
  session_start?: string;
  settlement_status?: string;
  vod_play_gb?: number;
}

export type LiveTrafficSettlementStatus = 'done' | 'failed' | 'pending';

export interface LiveTrafficSessionListQuery {
  list_rows?: number;
  live_id?: number;
  page?: number;
  settlement_status?: LiveTrafficSettlementStatus;
}

export async function getLiveStreamLogListApi(params: LiveStreamLogListQuery) {
  return requestClient.post<{ list: PaginatedList<LiveStreamLogItem> }>(
    '/api/v1/shop/live/stream-log/list',
    params,
  );
}

export async function getLiveTrafficSessionListApi(
  params: LiveTrafficSessionListQuery,
) {
  return requestClient.post<{ list: PaginatedList<LiveTrafficSessionItem> }>(
    '/api/v1/shop/traffic/session/list',
    params,
  );
}
