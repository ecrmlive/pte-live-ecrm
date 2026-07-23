import { requestClient } from '#/api/request';

export type LiveStatRange = '30d' | '7d' | 'today' | 'yesterday';

export interface LiveStreamLogOverview {
  avg_playback_duration_text?: string;
  avg_watch_duration_text?: string;
  conversion_rate?: string;
  deal_amount?: string;
  deal_order_count?: number;
  deal_user_count?: number;
  live_duration_text?: string;
  order_user_count?: number;
  pay_user_count?: number;
  peak_viewer_count?: number;
  playback_duration_text?: string;
  redpack_sent_amount?: string;
  refund_amount?: string;
  refund_order_count?: number;
  total_revenue?: string;
  traffic_text?: string;
  viewer_count?: number;
}

export interface LiveStreamLogSeries {
  bandwidth_mbps?: number[];
  flux_mb?: number[];
  times?: string[];
  values?: number[];
}

export interface LiveStreamLogDurationBuckets {
  labels?: string[];
  live?: number[];
  playback?: number[];
}

export interface LiveStreamLogStatsResult {
  cover_img?: string;
  duration_buckets?: LiveStreamLogDurationBuckets;
  end_time_text?: string;
  live_id?: number;
  log_id?: number;
  mode?: string;
  name?: string;
  overview?: LiveStreamLogOverview;
  play_series?: LiveStreamLogSeries;
  push_series?: LiveStreamLogSeries;
  range?: LiveStatRange;
  session_id?: string;
  start_time_text?: string;
  viewer_series?: LiveStreamLogSeries;
}

export interface LiveStreamLogStatsParams {
  log_id?: number;
  range?: LiveStatRange;
}

export async function getLiveStreamLogStatsApi(
  params: LiveStreamLogStatsParams = {},
) {
  return requestClient.post<LiveStreamLogStatsResult>(
    '/api/v1/shop/live/stream-log/stats',
    params,
  );
}
