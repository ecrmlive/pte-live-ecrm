/** 腾讯云 LVB stream_status（wx_live_stream） */
export const STREAM_STATUS = {
	IDLE: 0,
	PUSHING: 1,
	INTERRUPT: 2,
	FORBIDDEN: 3,
	ERROR: 4,
};

/** 主播暂时离开：三端统一文案 */
export const ANCHOR_AWAY_COPY = {
	title: '主播暂时离开',
	sub: '稍等片刻，精彩马上回来',
};

/** 观众端拉流失败且本机网络差时的主标题；副文案与离开态一致 */
export const NETWORK_UNSTABLE_COPY = {
	title: '网络状态不稳定',
	sub: ANCHOR_AWAY_COPY.sub,
};

/** 推流中拉流自愈重试间隔（不在 UI 展示倒计时） */
export const PULL_RETRY_DELAY_MS = 3000;

/** IM 11021 status（与 live-api CalcIMStatus 一致） */
export const IM_LIVE_STATUS = {
	NOT_STARTED: 0,
	LIVING: 1,
	ENDED: 2,
	AWAY: 3,
};

export function parseImPayload(data) {
	if (data == null) return {};
	if (typeof data === 'string') {
		try {
			return JSON.parse(data);
		} catch {
			return {};
		}
	}
	if (typeof data === 'object') return data;
	return {};
}

/** 断流/拉流失败遮罩文案：anchor_away | network_unstable */
export function getStreamInterruptOverlayCopy(kind) {
	if (kind === 'network_unstable') {
		return NETWORK_UNSTABLE_COPY;
	}
	if (kind === 'anchor_away') {
		return ANCHOR_AWAY_COPY;
	}
	return { title: '', sub: '' };
}

/** 检测观众端网络是否不可用或极弱（none / 2g / H5 offline） */
export function detectViewerNetworkWeak() {
	return new Promise((resolve) => {
		if (typeof navigator !== 'undefined' && navigator.onLine === false) {
			resolve(true);
			return;
		}
		if (typeof uni !== 'undefined' && typeof uni.getNetworkType === 'function') {
			uni.getNetworkType({
				success(res) {
					const type = String(res.networkType || '').toLowerCase();
					resolve(type === 'none' || type === '2g');
				},
				fail() {
					resolve(false);
				},
			});
			return;
		}
		resolve(false);
	});
}

/** 与 api-live/pkg/tencentlive/callback.go CalcIMStatus 对齐 */
export function calcImLiveStatus(roomType, sessionStatus, streamStatus) {
	if (roomType === 2) {
		if (sessionStatus === 3) return IM_LIVE_STATUS.ENDED;
		if (sessionStatus === 1) return IM_LIVE_STATUS.LIVING;
		return IM_LIVE_STATUS.NOT_STARTED;
	}
	if (sessionStatus === 3) return IM_LIVE_STATUS.ENDED;
	if (sessionStatus === 0) return IM_LIVE_STATUS.NOT_STARTED;
	if (sessionStatus === 2) return IM_LIVE_STATUS.AWAY;
	if (sessionStatus === 1) {
		if (streamStatus === STREAM_STATUS.INTERRUPT) return IM_LIVE_STATUS.AWAY;
		return IM_LIVE_STATUS.LIVING;
	}
	return IM_LIVE_STATUS.NOT_STARTED;
}

/**  Unix 秒或毫秒 → 毫秒时间戳（场次 actual_start_time） */
export function normalizeSessionStartMs(actualStartTime) {
	const ts = Number(actualStartTime) || 0;
	if (ts <= 0) return 0;
	// 秒级 Unix 时间约 ≤ 1e10；毫秒级 ≥ 1e12
	return ts >= 1e12 ? ts : ts * 1000;
}

/** 场次是否进行中（含主播暂时离开 session_status=2） */
export function isLiveSessionOngoing(sessionStatus) {
	return sessionStatus === 1 || sessionStatus === 2;
}

/** 是否处于「主播暂时离开」（应隐藏拉流画面、展示离开 UI） */
export function isAnchorAwayState({
	sessionStatus = 0,
	streamStatus = 0,
	forceAway = false,
	roomType = 1,
} = {}) {
	if (roomType === 2 || sessionStatus === 3) return false;
	if (forceAway) return true;
	return sessionStatus === 2 || streamStatus === STREAM_STATUS.INTERRUPT;
}

/** 是否应在中控预览区拉流播放（直播房 + stream_status=推流中且未结束） */
export function shouldPullLiveStream(sessionStatus = 0, streamStatus = 0, roomType = 1) {
	if (roomType === 2) return false;
	if (sessionStatus === 3) return false;
	return streamStatus === STREAM_STATUS.PUSHING;
}

/** 录播房是否应播放录播视频 */
export function shouldPlayRecordVideo(roomType = 1, sessionStatus = 0, recordVideoPath = '', recordVodMediaUrl = '') {
	if (roomType !== 2) return false;
	if (sessionStatus !== 1) return false;
	return Boolean(String(recordVodMediaUrl || recordVideoPath || '').trim());
}

/**
 * 中控推流标识：正在推流 / 暂停推流 / 等待推流
 */
export function resolveAnchorStreamDisplay({
	sessionStatus = 0,
	streamStatus = 0,
	imLiveStatus,
	roomType = 1,
} = {}) {
	if (roomType === 2) {
		if (sessionStatus === 1) {
			return {
				key: 'record',
				text: '录播播放中',
				tip: '正在播放录播视频',
				visible: true,
				viewerVisible: true,
			};
		}
		if (sessionStatus === 3) {
			return {
				key: 'idle',
				text: '未开始',
				tip: '本场已结束',
				visible: true,
				viewerVisible: false,
			};
		}
		return {
			key: 'idle',
			text: '未开始',
			tip: '录播尚未开始',
			visible: true,
			viewerVisible: false,
		};
	}

	const imStatus =
		imLiveStatus ?? calcImLiveStatus(roomType, sessionStatus, streamStatus);
	const viewerVisible = imStatus === IM_LIVE_STATUS.LIVING && streamStatus === STREAM_STATUS.PUSHING;

	if (sessionStatus === 0) {
		const idleMap = {
			[STREAM_STATUS.IDLE]: { text: '未推流', tip: '主播尚未推流' },
			[STREAM_STATUS.PUSHING]: { text: '正在推流', tip: '主播推流中' },
			[STREAM_STATUS.INTERRUPT]: { text: '暂停推流', tip: '推流已中断' },
			[STREAM_STATUS.FORBIDDEN]: { text: '禁推', tip: '推流已被禁止' },
			[STREAM_STATUS.ERROR]: { text: '推流异常', tip: '推流状态异常' },
		};
		const item = idleMap[streamStatus] || idleMap[STREAM_STATUS.IDLE];
		return {
			key: streamStatus === STREAM_STATUS.PUSHING ? 'pushing' : streamStatus === STREAM_STATUS.INTERRUPT ? 'paused' : 'idle',
			text: item.text,
			tip: item.tip,
			visible: true,
			viewerVisible: streamStatus === STREAM_STATUS.PUSHING,
		};
	}

	if (sessionStatus === 3) {
		return {
			key: 'idle',
			text: '未推流',
			tip: '本场已结束',
			visible: true,
			viewerVisible: false,
		};
	}

	if (streamStatus === STREAM_STATUS.PUSHING && sessionStatus === 1) {
		return {
			key: 'pushing',
			text: '正在推流',
			tip: '主播推流中，观众可见直播画面',
			visible: true,
			viewerVisible: true,
		};
	}

	if (
		sessionStatus === 2 ||
		streamStatus === STREAM_STATUS.INTERRUPT ||
		imStatus === IM_LIVE_STATUS.AWAY
	) {
		return {
			key: 'paused',
			text: '暂停推流',
			tip: '主播已断流，观众暂不可见画面',
			visible: true,
			viewerVisible: false,
		};
	}

	if (sessionStatus === 1 && streamStatus === STREAM_STATUS.IDLE) {
		return {
			key: 'waiting',
			text: '等待推流',
			tip: '已开播，等待主播开始推流',
			visible: true,
			viewerVisible: false,
		};
	}

	return {
		key: 'unknown',
		text: '未推流',
		tip: '推流状态未知',
		visible: sessionStatus === 1 || sessionStatus === 2,
		viewerVisible: false,
	};
}

/**
 * 解析 IM 直播相关消息。
 * 场次 session_status 仅由 11017/11018（开始/结束直播）变更；11019/11020/11021 只更新推流态与展示态。
 */
export function applyImLiveEvent(state, code, rawData) {
	const next = {
		sessionStatus: state.sessionStatus,
		streamStatus: state.streamStatus,
		imLiveStatus: state.imLiveStatus,
		roomType: state.roomType,
	};
	const data = parseImPayload(rawData);

	switch (code) {
		case 11017:
			next.sessionStatus = 1;
			break;
		case 11018:
			next.sessionStatus = 3;
			next.streamStatus = STREAM_STATUS.IDLE;
			break;
		case 11019:
			next.streamStatus = STREAM_STATUS.PUSHING;
			break;
		case 11020:
			next.streamStatus = STREAM_STATUS.INTERRUPT;
			break;
		case 11021:
			if (data.stream_status != null) {
				next.streamStatus = Number(data.stream_status) || 0;
			}
			if (data.status != null) {
				next.imLiveStatus = Number(data.status);
			}
			break;
		default:
			return null;
	}

	if (next.imLiveStatus == null) {
		next.imLiveStatus = calcImLiveStatus(next.roomType, next.sessionStatus, next.streamStatus);
	}
	return next;
}
