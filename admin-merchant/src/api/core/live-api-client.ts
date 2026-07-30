/**
 * 商户直播运行时 API（中控台/直播间）— 统一走 api-platform，路径 /api/v1/*
 */
import type {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from 'axios';

import axios from 'axios';
import qs from 'qs';
import { ElMessage } from 'element-plus';
import { useAccessStore } from '@vben/stores';

import { router } from '#/router';
import useUserStore from '#/adapters/control-user-store';
import { attachShopAppId } from '#/utils/qixi-live-shop-app-id';
import { getDecryptedToken, setEncryptedToken } from '#/utils/qixi-live-token';
import { resolveLiveApiBaseUrl } from '#/utils/qixi-live-api';
import { formatUserFacingApiError } from '#/utils/api-error';
import {
  attachAPIEncryption,
  decryptAPIResponse,
} from '../../../../admin-platform/src/utils/api-crypto';

export interface LiveApiEnvelope<T = unknown> {
  code: number;
  data?: T;
  msg: string;
}

declare module 'axios' {
  interface AxiosRequestConfig {
    silent?: boolean;
  }
}

type LiveRequestConfig = AxiosRequestConfig & { silent?: boolean };
type LivePostData = Record<string, unknown> | undefined;

const liveApiBase = resolveLiveApiBaseUrl();

const LIVE_API_ERROR_TOAST_DEDUP_MS = 4000;
let lastLiveApiErrorToast = {
  message: '',
  timestamp: 0,
};

/**
 * 中控台会在初始化时并发读取房间、商品、统计等数据。
 * API 短暂不可用时只展示一次相同错误，避免多条 Toast 遮挡操作界面。
 */
function notifyLiveApiError(message: string) {
  const now = Date.now();
  if (
    lastLiveApiErrorToast.message === message &&
    now - lastLiveApiErrorToast.timestamp < LIVE_API_ERROR_TOAST_DEDUP_MS
  ) {
    return;
  }
  lastLiveApiErrorToast = { message, timestamp: now };
  ElMessage({
    showClose: true,
    message,
    type: 'error',
  });
}

const liveProductListInFlight = new Map<string, Promise<LiveApiEnvelope>>();

function normalizeProductListValue(value: unknown): unknown {
	if (value === undefined || value === null) {
		return undefined;
	}
	if (Array.isArray(value)) {
		return value.map(normalizeProductListValue);
	}
	if (value instanceof Date) {
		return value.getTime();
	}
	if (typeof value === 'number' && Number.isFinite(value)) {
		return value;
	}
	if (typeof value === 'boolean') {
		return value ? 1 : 0;
	}
	if (typeof value === 'string') {
		const trimmed = value.trim();
		if (trimmed === '') {
			return '';
		}
		const num = Number(trimmed);
		return Number.isFinite(num) && Number.isInteger(num) ? num : trimmed;
	}
	return value;
}

function toSortedRecord(data: Record<string, unknown>) {
	const out: Record<string, unknown> = {};
	const keys = Object.keys(data || {}).sort();
	for (const key of keys) {
		const value = normalizeProductListValue(data[key]);
		if (value === undefined || value === null) {
			continue;
		}
		out[key] = value;
	}
	return out;
}

function cloneProductListPayload(data: LivePostData = {}) {
	if (!data || typeof data !== 'object') {
		return {};
	}
	return toSortedRecord(data as Record<string, unknown>);
}

function buildLiveProductListKey(data: LivePostData = {}) {
	return JSON.stringify(cloneProductListPayload(data));
}

const client: AxiosInstance = axios.create({
  baseURL: liveApiBase,
  timeout: 30_000,
  // 中控台通过 authori-zation Bearer Token 鉴权，不依赖跨域 Cookie。
  // 测试 API 保持 Allow-Origin: * 时，浏览器禁止携带 credentials 的请求，
  // 否则会将正常的接口响应误判为 CORS 错误。
  withCredentials: false,
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
  },
});

function getBearerToken(token?: null | string) {
  const raw = token || getDecryptedToken() || '';
  return raw ? `Bearer ${raw}` : '';
}

client.interceptors.request.use(async (config: InternalAxiosRequestConfig) => {
	const userStore = useUserStore();
	const accessStore = useAccessStore();
	const token =
		accessStore.accessToken || getDecryptedToken() || '';
	config.headers = config.headers || {};
	const bearerToken = getBearerToken(token);
	if (bearerToken) {
		config.headers['Authori-zation'] = bearerToken;
	}
	attachShopAppId(config, userStore.userInfo, token);
	if (
		config.method === 'post' &&
		!config.headers.uploadImg &&
		!(config.data instanceof FormData) &&
    typeof config.data !== 'string'
  ) {
    config.data = qs.stringify(config.data || {});
  }
  return attachAPIEncryption(config, liveApiBase);
});

client.interceptors.response.use(
  async (res: AxiosResponse<LiveApiEnvelope>) => {
    await decryptAPIResponse(res);
    const userStore = useUserStore();
    const authorization =
      res.headers['authori-zation'] || res.headers.authorize || '';
      if (typeof authorization === 'string' && authorization.indexOf('Bearer ') === 0) {
      const nextToken = authorization.slice(7).trim();
      const accessStore = useAccessStore();
      userStore.setToken(nextToken);
      accessStore.setAccessToken(nextToken);
      setEncryptedToken(nextToken);
    }
    if (res.data.code === 1) {
      return res.data;
    }
    if (res.data.code === -1) {
      userStore.afterLogout();
      void router.push({ path: '/auth/login' });
      return Promise.reject(res.data);
    }
    if (res.config?.silent) {
      return Promise.reject(res.data);
    }
    notifyLiveApiError(res.data.msg || '请求失败，请稍后再试');
    return Promise.reject(res.data);
  },
  (error: { config?: LiveRequestConfig; response?: { status?: number } }) => {
	const status = error.response?.status;
	if (status === 401) {
		const userStore = useUserStore();
		const accessStore = useAccessStore();
		userStore.afterLogout();
		accessStore.setAccessToken(null);
		void router.push({ path: '/auth/login' });
	}
    if (error?.config?.silent) {
      return Promise.reject(error);
    }
    let message = formatUserFacingApiError(error, '接口服务无法连接');
    if (status) {
      message =
        status >= 500 ? '接口服务异常，请稍后再试' : `接口请求失败 (HTTP ${status})`;
    }
    notifyLiveApiError(message);
    return Promise.reject(error);
  },
);

function post<T = unknown>(
  path: string,
  data: LivePostData = {},
  config: LiveRequestConfig = {},
) {
  return client.post(path, data, config) as Promise<LiveApiEnvelope<T>>;
}

function get<T = unknown>(path: string, params: Record<string, unknown> = {}) {
  return client.get(path, { params }) as Promise<LiveApiEnvelope<T>>;
}

const LiveApiClient = {
  baseURL: liveApiBase,

  get,

  post,

  /** 房间详情（C 端 user JWT） */
  roomDetail(data: LivePostData = {}) {
    return post('/api/v1/room/detail', data);
  },

  /** 房间详情（商户 shop JWT，中控/后台） */
  roomDetailShop(data: LivePostData = {}) {
    return post('/api/v1/shop/live/detail', data);
  },

  /** 在线/累计人数（高频轮询，silent 失败不弹 toast） */
  roomCounts(data: LivePostData = {}) {
    return post('/api/v1/room/counts', data, { silent: true });
  },

  /** 直播设置读 */
  roomSettingsGet(data: LivePostData = {}) {
    return post('/api/v1/room/settings', data);
  },

  /** 直播设置写（shop JWT） */
  roomSettingsUpdate(data: LivePostData = {}) {
    return post('/api/v1/room/settings/update', data);
  },

  ping() {
    return get('/ping');
  },

  anchorList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/anchor/list', data);
  },

  anchorDetail(data: LivePostData = {}) {
    return post('/api/v1/shop/live/anchor/detail', data);
  },

  anchorCreate(data: LivePostData = {}) {
    return post('/api/v1/shop/live/anchor/create', data);
  },

  anchorUpdate(data: LivePostData = {}) {
    return post('/api/v1/shop/live/anchor/update', data);
  },

  anchorDelete(data: LivePostData = {}) {
    return post('/api/v1/shop/live/anchor/delete', data);
  },

  liveList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/list', data);
  },

  liveCreate(data: LivePostData = {}) {
    return post('/api/v1/shop/live/create', data);
  },

  liveUpdate(data: LivePostData = {}) {
    return post('/api/v1/shop/live/update', data);
  },

  liveDelete(data: LivePostData = {}) {
    return post('/api/v1/shop/live/delete', data);
  },

  liveSetTop(data: LivePostData = {}) {
    return post('/api/v1/shop/live/set-top', data);
  },

  /** TCPlayer License（商户中控 / 与 C 端 H5 共用） */
  livePlayerConfig(data: LivePostData = {}, config: LiveRequestConfig = {}) {
    return post('/api/v1/player/config', data, config);
  },

  liveStream(data: LivePostData = {}, config: LiveRequestConfig = {}) {
    return post('/api/v1/shop/live/stream', data, config);
  },

  /** 中控预览拉流地址刷新（失败由预览区 UI 呈现，不弹 toast） */
  liveStreamRefreshURL(data: LivePostData = {}, config: LiveRequestConfig = {}) {
    return post('/api/v1/shop/live/stream/refresh-url', data, config);
  },

  liveStreamLogList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/stream-log/list', data);
  },

  liveStreamLogStats(data: LivePostData = {}) {
    return post('/api/v1/shop/live/stream-log/stats', data);
  },

  liveAudience(data: LivePostData = {}) {
    return post('/api/v1/shop/live/audience', data);
  },

  liveAudienceList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/audience/list', data);
  },

  liveAudienceMute(data: LivePostData = {}) {
    return post('/api/v1/shop/live/audience/mute', data);
  },

  liveAudienceKick(data: LivePostData = {}) {
    return post('/api/v1/shop/live/audience/kick', data);
  },

  liveAudienceMuteList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/audience/mute/list', data);
  },

  liveAudienceMuteBatchUnmute(data: LivePostData = {}) {
    return post('/api/v1/shop/live/audience/mute/batch-unmute', data);
  },

  liveAudienceKickList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/audience/kick/list', data);
  },

  liveAudienceKickBatchUnkick(data: LivePostData = {}) {
    return post('/api/v1/shop/live/audience/kick/batch-unkick', data);
  },

  liveInviteStats(data: LivePostData = {}) {
    return post('/api/v1/shop/live/invite/stats', data);
  },

  liveBandwidthSummary(data: LivePostData = {}) {
    return post('/api/v1/shop/live/bandwidth/summary', data);
  },

  liveBandwidthDaily(data: LivePostData = {}) {
    return post('/api/v1/shop/live/bandwidth/daily', data);
  },

  liveBandwidthSessionList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/bandwidth/session/list', data);
  },

  liveBandwidthSync(data: LivePostData = {}) {
    return post('/api/v1/shop/live/bandwidth/sync', data);
  },

  trafficSummary(data: LivePostData = {}) {
    return post('/api/v1/shop/traffic/summary', data);
  },

  trafficDailyLvb(data: LivePostData = {}) {
    return post('/api/v1/shop/traffic/daily/lvb', data);
  },

  trafficDailyVod(data: LivePostData = {}) {
    return post('/api/v1/shop/traffic/daily/vod', data);
  },

  trafficSessionList(data: LivePostData = {}) {
    return post('/api/v1/shop/traffic/session/list', data);
  },

  trafficRechargeList(data: LivePostData = {}) {
    return post('/api/v1/shop/traffic/recharge/list', data);
  },

  liveComplaintList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/complaint/list', data);
  },

  liveComplaintDetail(data: LivePostData = {}) {
    return post('/api/v1/shop/live/complaint/detail', data);
  },

  liveComplaintHandle(data: LivePostData = {}) {
    return post('/api/v1/shop/live/complaint/handle', data);
  },

  liveEnd(data: LivePostData = {}) {
    return post('/api/v1/shop/live/end', data);
  },

  liveStart(data: LivePostData = {}) {
    return post('/api/v1/shop/live/start', data);
  },

  liveSessionDetail(data: LivePostData = {}) {
    return post('/api/v1/shop/live/session/detail', data, { silent: true });
  },

  /** 直播间实时聚合数据（C 端 / H5 / 小程序） */
  roomLiveStats(data: LivePostData = {}) {
    return post('/api/v1/room/live-stats', data, { silent: true });
  },

  /** 商户中控实时数据（shop JWT） */
  liveRoomStats(data: LivePostData = {}) {
    return post('/api/v1/shop/live/live-stats', data, { silent: true });
  },

  liveProductList(data: LivePostData = {}, config: LiveRequestConfig = {}) {
	const payload = cloneProductListPayload(data);
	const key = `${buildLiveProductListKey(payload)}`;
	const existing = liveProductListInFlight.get(key);
	if (existing) {
		return existing;
	}

	const request = post('/api/v1/shop/live/product/list', payload, config);
	liveProductListInFlight.set(key, request);
	return request.finally(() => {
		liveProductListInFlight.delete(key);
	});
  },

  liveProductAdd(data: LivePostData = {}) {
    return post('/api/v1/shop/live/product/add', data);
  },

  liveProductUpdate(data: LivePostData = {}) {
    return post('/api/v1/shop/live/product/update', data);
  },

  liveProductOnSale(data: LivePostData = {}) {
    return post('/api/v1/shop/live/product/onsale', data);
  },

  liveProductExplain(data: LivePostData = {}) {
    return post('/api/v1/shop/live/product/explain', data);
  },

  liveProductDelete(data: LivePostData = {}) {
    return post('/api/v1/shop/live/product/delete', data);
  },

  liveProductShelfBatches(data: LivePostData = {}) {
    return post('/api/v1/shop/live/product/shelf-batches', data);
  },

  liveProductScheduleSave(data: LivePostData = {}) {
    return post('/api/v1/shop/live/product/schedule/save', data);
  },

  liveProductScheduleList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/product/schedule/list', data);
  },

  liveProductScheduleCancel(data: LivePostData = {}) {
    return post('/api/v1/shop/live/product/schedule/cancel', data);
  },

  liveRedpackSend(data: LivePostData = {}) {
    return post('/api/v1/shop/live/redpack/send', data);
  },

  liveRedpackList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/redpack/list', data);
  },

  liveRedpackReceiveList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/redpack/receive-list', data);
  },

  liveDanmakuList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku/list', data);
  },

  liveDanmakuRecent(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku/recent', data);
  },

  liveDanmakuSend(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku/send', data);
  },

  liveSessionChatList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/session/chat/list', data);
  },

  liveDanmakuAudit(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku/audit', data);
  },

  liveDanmakuPendingCount(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku/pending-count', data);
  },

  liveSensitiveWordList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/sensitive-word/list', data);
  },

  liveSensitiveWordAdd(data: LivePostData = {}) {
    return post('/api/v1/shop/live/sensitive-word/add', data);
  },

  liveSensitiveWordDelete(data: LivePostData = {}) {
    return post('/api/v1/shop/live/sensitive-word/delete', data);
  },

  liveDanmakuBotRuntime(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/runtime', data);
  },

  liveDanmakuBotScriptTemplateList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/script-template/list', data);
  },

  liveDanmakuBotScriptTemplateSave(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/script-template/save', data);
  },

  liveDanmakuBotScriptTemplateDelete(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/script-template/delete', data);
  },

  liveDanmakuBotRobotTemplateList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/robot-template/list', data);
  },

  liveDanmakuBotRobotTemplateSave(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/robot-template/save', data);
  },

  liveDanmakuBotRobotTemplateDelete(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/robot-template/delete', data);
  },

  liveDanmakuBotPlatformPick(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/user/platform/pick', data);
  },

  liveDanmakuBotCustomSave(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/user/custom/save', data);
  },

  liveDanmakuBotCustomDelete(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/user/custom/delete', data);
  },

  liveDanmakuBotTaskStart(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/task/start', data);
  },

  liveDanmakuBotTaskList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/task/list', data, { silent: true });
  },

  liveDanmakuBotTaskPause(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/task/pause', data);
  },

  liveDanmakuBotTaskResume(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/task/resume', data);
  },

  liveDanmakuBotTaskStop(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/task/stop', data);
  },

  liveDanmakuBotTaskCancelSchedule(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/task/cancel-schedule', data);
  },

  liveDanmakuBotTaskStatus(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/task/status', data, { silent: true });
  },

  vodApplyUpload(data: LivePostData = {}) {
    return post('/api/v1/shop/live/vod/apply-upload', data);
  },

  vodCommitUpload(data: LivePostData = {}) {
    return post('/api/v1/shop/live/vod/commit-upload', data);
  },

  vodPullUpload(data: LivePostData = {}) {
    return post('/api/v1/shop/live/vod/pull-upload', data);
  },

  vodList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/vod/list', data);
  },

  vodDelete(data: LivePostData = {}) {
    return post('/api/v1/shop/live/vod/delete', data);
  },

  liveDanmakuBotLogList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/danmaku-bot/log/list', data);
  },

  liveHeatBoostTaskCreate(data: LivePostData = {}) {
    return post('/api/v1/shop/live/heat-boost/task/create', data);
  },

  liveHeatBoostTaskList(data: LivePostData = {}) {
    return post('/api/v1/shop/live/heat-boost/task/list', data, { silent: true });
  },

  liveHeatBoostTaskStop(data: LivePostData = {}) {
    return post('/api/v1/shop/live/heat-boost/task/stop', data);
  },

  liveHeatBoostTaskPause(data: LivePostData = {}) {
    return post('/api/v1/shop/live/heat-boost/task/pause', data);
  },

  liveHeatBoostTaskResume(data: LivePostData = {}) {
    return post('/api/v1/shop/live/heat-boost/task/resume', data);
  },

  liveHeatBoostTaskCancel(data: LivePostData = {}) {
    return post('/api/v1/shop/live/heat-boost/task/cancel', data);
  },
} as const;

export default LiveApiClient;
