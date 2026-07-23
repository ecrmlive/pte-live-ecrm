import type { LiveRoomForm } from '#/api/core/live';

import { toCosRelativePath } from '#/utils/live/cosMediaUrl.js';

export const DEFAULT_LIVE_ROOM_SYSTEM_NOTICE =
  '📢 温馨提示：本直播间所有商品均支持平台售后服务，商品价格及优惠活动以直播页面实际显示为准。如遇物流或售后问题，请联系在线客服，我们将竭诚为您服务。';

export function buildLiveRoomPayload(
  form: LiveRoomForm & {
    close_comment?: number;
    close_goods?: number;
    close_like?: number;
    close_replay?: number;
  },
  extra?: { product_ids?: string },
) {
  const coverKey = toCosRelativePath(form.cover_img || '');
  const backgroundKey = toCosRelativePath(form.background_img || '');
  const enableReplay = Number(form.enable_replay ?? 1) === 1;
  const allowChat = Number(form.allow_chat ?? 1) === 1;

  return {
    ...form,
    allow_chat: allowChat ? 1 : 0,
    anchor_id: form.anchor_id || 0,
    background_img: backgroundKey,
    close_comment: allowChat ? 0 : 1,
    close_replay: enableReplay ? 0 : 1,
    cover_img: coverKey,
    enable_gift: Number(form.enable_gift ?? 1),
    enable_linkmic: Number(form.enable_linkmic ?? 0),
    enable_replay: enableReplay ? 1 : 0,
    feeds_img: toCosRelativePath(form.feeds_img || coverKey),
    fire_value: Number(form.fire_value ?? 0),
    is_visible: Number(form.is_visible ?? 1),
    product_ids: extra?.product_ids,
    record_vod_media_url: toCosRelativePath(form.record_vod_media_url || ''),
    share_img: toCosRelativePath(form.share_img || coverKey),
    show_heat: Number(form.show_heat ?? 1),
    show_home: Number(form.show_home ?? 1),
    show_online_count: Number(form.show_online_count ?? 1),
    show_total_count: Number(form.show_total_count ?? 1),
    stream_orientation: Number(form.stream_orientation ?? 2),
    system_notice: form.system_notice || '',
    share_intro: form.share_intro || '',
    watch_password: form.watch_password || '',
  };
}

export function defaultLiveRoomForm(): LiveRoomForm & {
  close_comment: number;
  close_goods: number;
  close_like: number;
  close_replay: number;
} {
  return {
    allow_chat: 1,
    anchor_id: 0,
    anchor_name: '',
    anchor_wechat: '',
    background_img: '',
    close_comment: 0,
    close_goods: 0,
    close_like: 0,
    close_replay: 0,
    cover_img: '',
    enable_gift: 1,
    enable_linkmic: 0,
    enable_replay: 1,
    end_time: '',
    feeds_img: '',
    fire_value: 0,
    is_visible: 1,
    name: '',
    record_video_path: '',
    record_vod_file_id: '',
    record_vod_media_url: '',
    room_type: 1,
    share_img: '',
    share_intro: '',
    show_heat: 1,
    show_home: 1,
    show_online_count: 1,
    show_total_count: 1,
    start_time: '',
    stream_orientation: 2,
    system_notice: DEFAULT_LIVE_ROOM_SYSTEM_NOTICE,
    watch_password: '',
  };
}
