/**
 * 中控 IM WebSocket（qixi-live-im :11510）
 * 握手：/ws?sdkAppID=&identifier=&userSig=（浏览器无法自定义 Header，走 Query）
 * 进房：握手成功后发送 scene.enter，不再使用旧 /ws?roomId=。
 * 断线后自动重连，直至手动 disconnect()。
 */

import { parseImPayload } from '#/utils/live/liveStreamStatus.js';

/** 电商直播 IM 消息码（与 qixi-live-im live-commerce 一致） */
export const LIVE_IM_CODE = {
	PRODUCT_EXPLAIN_START: 11001,
	PRODUCT_EXPLAIN_CANCEL: 11002,
	PRODUCT_ON_SHELF: 11005,
	PRODUCT_OFF_SHELF: 11006,
	MUTE_ALL: 11007,
	UNMUTE_ALL: 11008,
	MUTE_USER: 11009,
	UNMUTE_USER: 11010,
	KICK_USER: 11016,
	ONLINE_COUNT: 11022,
	TOTAL_COUNT: 11023,
	CONFIG_CHANGE: 11024,
	USER_ENTER_WELCOME: 11028,
	STREAM_EXCEPTION: 11029,
};

const PRODUCT_LIST_REFRESH_CODES = new Set([
	LIVE_IM_CODE.PRODUCT_EXPLAIN_START,
	LIVE_IM_CODE.PRODUCT_EXPLAIN_CANCEL,
	LIVE_IM_CODE.PRODUCT_ON_SHELF,
	LIVE_IM_CODE.PRODUCT_OFF_SHELF,
]);

export function isProductListRefreshImCode(code) {
	return PRODUCT_LIST_REFRESH_CODES.has(Number(code));
}

/** 在线观众列表需校准的 IM 码（进房 / 踢人 / 禁言 / 人数变化） */
const AUDIENCE_LIST_SYNC_CODES = new Set([
	LIVE_IM_CODE.USER_ENTER_WELCOME,
	LIVE_IM_CODE.KICK_USER,
	LIVE_IM_CODE.MUTE_USER,
	LIVE_IM_CODE.UNMUTE_USER,
	LIVE_IM_CODE.ONLINE_COUNT,
]);

export function isAudienceListSyncImCode(code) {
	return AUDIENCE_LIST_SYNC_CODES.has(Number(code));
}

/** IM 11028 → 观众列表行（乐观插入） */
export function parseImAudienceEnter(msg) {
	if (!msg || Number(msg.code) !== LIVE_IM_CODE.USER_ENTER_WELCOME) return null;
	const data = parseImPayload(msg.data);
	const userId = Number(msg.sendUserId || data.userId || data.user_id) || 0;
	if (userId <= 0) return null;
	const nickName =
		String(data.nickName || data.nick_name || data.nick || '').trim() ||
		`用户${userId}`;
	let role = Number(data.role);
	if (Number.isNaN(role)) role = 0;
	return {
		user_id: userId,
		nick_name: nickName,
		avatar: String(data.avatar || data.avatarUrl || '').trim(),
		role,
		is_muted: 0,
		is_kicked: 0,
	};
}

/** IM 11016 → 被踢 user_id */
export function parseImKickUser(msg) {
	if (!msg || Number(msg.code) !== LIVE_IM_CODE.KICK_USER) return null;
	const data = parseImPayload(msg.data);
	const userId = Number(data.userId || data.user_id || msg.sendUserId) || 0;
	return userId > 0 ? { user_id: userId } : null;
}
const RECONNECT_MS = 2000;
const HANDSHAKE_TIMEOUT_MS = 5000;
const CONNECTING_STUCK_MS = HANDSHAKE_TIMEOUT_MS + 1500;

export function resolveImWsBase() {
	const fromEnv = import.meta.env.VITE_LIVE_IM_WS_URL;
	if (fromEnv) {
		return String(fromEnv).replace(/\/$/, '');
	}
	if (import.meta.env.DEV) {
		return 'ws://127.0.0.1:11510/ws';
	}
	const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
	return `${protocol}//${window.location.host}/ws`;
}

function buildWsUrl({ sdkAppID, sdk_app_id, identifier, userSig, user_sig, extend }) {
	const params = new URLSearchParams({
		sdkAppID: String(sdkAppID || sdk_app_id || ''),
		identifier: String(identifier || ''),
		userSig: String(userSig || user_sig || ''),
	});
	if (extend) {
		params.set('extend', typeof extend === 'string' ? extend : JSON.stringify(extend));
	}
	return `${resolveImWsBase()}?${params.toString()}`;
}

function createSceneRequestId(action) {
	return `${action}:${Date.now()}:${Math.random().toString(36).slice(2, 8)}`;
}

function parseSceneAck(msg) {
	const data = parseImPayload(msg?.data);
	return data?.type === 'scene.ack' ? data : null;
}

export function createLiveImClient(hooks = {}) {
	let ws = null;
	let reconnectInterval = null;
	let handshakeTimer = null;
	let manualClose = false;
	let shouldReconnect = false;
	let connectParams = null;
	let status = 'disconnected';
	let connectingSince = 0;
	let session = { clientId: '', userId: '', appId: '', roomId: '' };
	let boundOnlineHandler = null;
	let boundVisibilityHandler = null;

	const emitStatus = (next) => {
		status = next;
		if (next === 'connected') {
			clearReconnect();
			connectingSince = 0;
		} else if (next === 'connecting') {
			connectingSince = Date.now();
		} else if (next !== 'connecting') {
			connectingSince = 0;
		}
		hooks.onStatus?.(next);
	};

	const emitError = (payload) => {
		hooks.onError?.(payload);
	};

	const clearSession = () => {
		session = { clientId: '', userId: '', appId: '', roomId: '' };
	};

	const applySessionFromHandshake = (msg, params) => {
		const data = parseImPayload(msg?.data);
		session = {
			clientId: String(data.clientId || '').trim(),
			userId: String(data.userId || '').trim(),
			appId: String(params?.appId || '').trim(),
			roomId: String(params?.roomId || '').trim(),
			sdkAppID: String(data.sdkAppID || params?.sdkAppID || '').trim(),
			identifier: String(data.identifier || params?.identifier || '').trim(),
		};
		hooks.onSession?.({ ...session });
	};

	const sendSceneEnter = (params) => {
		if (!ws || !params?.roomId) return;
		const payload = {
			action: 'scene.enter',
			request_id: createSceneRequestId('scene.enter'),
			scene: 'shop',
			room_id: String(params.roomId),
			extend: typeof params.extend === 'string' ? params.extend : JSON.stringify(params.extend || {}),
		};
		try {
			ws.send(JSON.stringify(payload));
		} catch {}
	};

	const sendDanmaku = async () => {
		if (status !== 'connected' || !session.clientId) {
			throw new Error('IM 未连接，请稍候重试');
		}
		throw new Error('弹幕发送请走 api-platform 接口');
	};

	const clearHandshakeTimer = () => {
		if (handshakeTimer) {
			clearTimeout(handshakeTimer);
			handshakeTimer = null;
		}
	};

	const resolveConnectParams = (params) => {
		try {
			const fromHook = hooks.resolveParams?.();
			if (fromHook?.appId && fromHook?.roomId && fromHook?.sdkAppID && fromHook?.identifier && fromHook?.userSig) {
				return fromHook;
			}
		} catch (err) {
			emitError({ phase: 'resolveParams', message: err?.message || String(err) });
		}
		if (params?.appId && params?.roomId && params?.sdkAppID && params?.identifier && params?.userSig) {
			return params;
		}
		if (connectParams?.appId && connectParams?.roomId && connectParams?.sdkAppID && connectParams?.identifier && connectParams?.userSig) {
			return connectParams;
		}
		return null;
	};

	const clearReconnect = () => {
		if (reconnectInterval) {
			clearInterval(reconnectInterval);
			reconnectInterval = null;
		}
	};

	const tickReconnect = () => {
		if (manualClose || !shouldReconnect) return;
		if (status === 'connected') return;
		if (status === 'connecting') {
			if (connectingSince && Date.now() - connectingSince >= CONNECTING_STUCK_MS) {
				emitError({ phase: 'handshake', message: '连接超时，正在重试…' });
				closeSocket();
				emitStatus('disconnected');
			} else {
				return;
			}
		}
		void doConnect(connectParams);
	};

	const scheduleReconnect = () => {
		if (manualClose || !shouldReconnect) return;
		tickReconnect();
		if (reconnectInterval) return;
		reconnectInterval = setInterval(tickReconnect, RECONNECT_MS);
	};

	const tryReconnectNow = () => {
		if (manualClose || !shouldReconnect || status === 'connected') {
			return;
		}
		if (status === 'connecting') {
			if (!connectingSince || Date.now() - connectingSince < CONNECTING_STUCK_MS) {
				return;
			}
			closeSocket();
			emitStatus('disconnected');
		}
		clearReconnect();
		doConnect(connectParams);
	};

	const bindNetworkListeners = () => {
		if (typeof window === 'undefined' || boundOnlineHandler) return;
		boundOnlineHandler = () => tryReconnectNow();
		boundVisibilityHandler = () => {
			if (document.visibilityState === 'visible') {
				tryReconnectNow();
			}
		};
		window.addEventListener('online', boundOnlineHandler);
		document.addEventListener('visibilitychange', boundVisibilityHandler);
	};

	const unbindNetworkListeners = () => {
		if (typeof window === 'undefined' || !boundOnlineHandler) return;
		window.removeEventListener('online', boundOnlineHandler);
		document.removeEventListener('visibilitychange', boundVisibilityHandler);
		boundOnlineHandler = null;
		boundVisibilityHandler = null;
	};

	const closeSocket = () => {
		clearHandshakeTimer();
		if (!ws) return;
		ws.onopen = null;
		ws.onmessage = null;
		ws.onerror = null;
		ws.onclose = null;
		ws.close();
		ws = null;
	};

	const handleHandshakeFailure = (message, code) => {
		emitError({ phase: 'handshake', message, code });
		emitStatus('disconnected');
		scheduleReconnect();
	};

	const startHandshakeTimer = () => {
		clearHandshakeTimer();
		handshakeTimer = setTimeout(() => {
			if (status !== 'connecting') return;
			emitError({ phase: 'handshake', message: '连接超时，正在重试…' });
			closeSocket();
			emitStatus('disconnected');
			scheduleReconnect();
		}, HANDSHAKE_TIMEOUT_MS);
	};

	const doConnect = async (params) => {
		if (manualClose) return;

		const resolved = resolveConnectParams(params);
		if (!resolved?.appId || !resolved?.roomId || !resolved?.sdkAppID || !resolved?.identifier || !resolved?.userSig) {
			emitError({
				phase: 'params',
				message: '缺少 appId / roomId / sdkAppID / identifier / userSig，无法连接 IM',
			});
			emitStatus('disconnected');
			scheduleReconnect();
			return;
		}

		connectParams = resolved;

		closeSocket();
		emitStatus('connecting');

		try {
			ws = new WebSocket(buildWsUrl(resolved));
		} catch (err) {
			emitError({ phase: 'socket', message: err?.message || 'WebSocket 创建失败' });
			emitStatus('disconnected');
			scheduleReconnect();
			return;
		}

		startHandshakeTimer();

		ws.onmessage = (ev) => {
			let msg;
			try {
				msg = JSON.parse(ev.data);
			} catch {
				return;
			}
			if (msg && typeof msg.code === 'number') {
				const sceneAck = parseSceneAck(msg);
				if (sceneAck) {
					if (sceneAck.ok === false) {
						emitError({ phase: 'scene', message: msg.msg || 'IM 场景订阅失败', code: msg.code });
					}
				} else if (msg.code === 0 && status === 'connecting') {
					clearHandshakeTimer();
					applySessionFromHandshake(msg, connectParams);
					sendSceneEnter(connectParams);
					emitStatus('connected');
				} else if (status === 'connecting') {
					handleHandshakeFailure(msg.msg || 'IM 握手失败', msg.code);
				}
			}
			hooks.onMessage?.(msg);
		};

		ws.onerror = () => {
			if (manualClose) return;
			emitError({
				phase: 'socket',
				message: '无法连接聊天服务，请确认 IM 已启动',
			});
			if (status === 'connecting' || status === 'connected') {
				emitStatus('disconnected');
			}
			scheduleReconnect();
		};

		ws.onclose = () => {
			ws = null;
			clearHandshakeTimer();
			if (manualClose) return;
			if (status !== 'disconnected') {
				emitStatus('disconnected');
			}
			scheduleReconnect();
		};
	};

	const connect = (params) => {
		manualClose = false;
		shouldReconnect = true;
		bindNetworkListeners();
		const resolved = resolveConnectParams(params);
		if (
			resolved &&
			status === 'connected' &&
			session.clientId &&
			session.appId === String(resolved.appId) &&
			session.roomId === String(resolved.roomId)
		) {
			return;
		}
		void doConnect(params);
	};

	const disconnect = (manual = true) => {
		if (manual) {
			manualClose = true;
			shouldReconnect = false;
			connectParams = null;
			unbindNetworkListeners();
		}
		clearReconnect();
		closeSocket();
		clearSession();
		if (manual) {
			emitStatus('disconnected');
		}
	};

	return {
		connect,
		disconnect,
		getStatus: () => status,
		getSession: () => ({ ...session }),
		sendDanmaku,
	};
}
