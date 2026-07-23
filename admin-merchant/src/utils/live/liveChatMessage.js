import { parseImPayload } from '#/utils/live/liveStreamStatus.js';

/** wx_live_danmaku / wx_live_audience：0 观众 1 管理 2 主播；3 仅展示用（source=2 机器人） */
export const LIVE_CHAT_ROLE = {
	VIEWER: 0,
	ADMIN: 1,
	ANCHOR: 2,
	BOT: 3,
};

/** wx_live_danmaku.source：0 用户 2 机器人 */
export const LIVE_CHAT_SOURCE = {
	USER: 0,
	BOT: 2,
};

const ROLE_LABELS = {
	[LIVE_CHAT_ROLE.VIEWER]: '观众',
	[LIVE_CHAT_ROLE.ADMIN]: '管理员',
	[LIVE_CHAT_ROLE.ANCHOR]: '主播',
	[LIVE_CHAT_ROLE.BOT]: '机器人',
};

export function normalizeLiveChatRole(role, source) {
	if (Number(source) === LIVE_CHAT_SOURCE.BOT) {
		return LIVE_CHAT_ROLE.BOT;
	}
	const n = Number(role);
	if (n === LIVE_CHAT_ROLE.ADMIN || n === LIVE_CHAT_ROLE.ANCHOR) return n;
	return LIVE_CHAT_ROLE.VIEWER;
}

export function getLiveChatRoleLabel(role, source) {
	return ROLE_LABELS[normalizeLiveChatRole(role, source)] || '观众';
}

export function getLiveChatRoleTone(role, source) {
	const r = normalizeLiveChatRole(role, source);
	if (r === LIVE_CHAT_ROLE.ANCHOR) return 'anchor';
	if (r === LIVE_CHAT_ROLE.ADMIN) return 'admin';
	if (r === LIVE_CHAT_ROLE.BOT) return 'bot';
	return 'viewer';
}

/** 仅观众可禁言/踢出 */
export function canModerateLiveAudience(role, source) {
	return normalizeLiveChatRole(role, source) === LIVE_CHAT_ROLE.VIEWER;
}

/** 在线观众唯一键：角色 + 用户 ID（与 IM Redis u:/a:/s: 及 live-api 去重一致） */
export function audienceMemberKey(row) {
	const role = normalizeLiveChatRole(row?.role, row?.source);
	const uid = Number(row?.user_id) || 0;
	return `${role}:${uid}`;
}

/** 列表按角色+ID 去重（保留后者覆盖前者） */
export function dedupeAudienceList(list) {
	const map = new Map();
	(list || []).forEach((row) => {
		if (!row?.user_id) return;
		const key = audienceMemberKey(row);
		map.set(key, { ...(map.get(key) || {}), ...row });
	});
	return Array.from(map.values());
}

const BOT_PUBLIC_USER_ID_BASE = 900000000;

function resolveLiveChatSource(data, userId) {
	const source = Number(data?.source ?? 0);
	if (source === LIVE_CHAT_SOURCE.BOT) return LIVE_CHAT_SOURCE.BOT;
	const uid = Number(userId);
	if (uid >= BOT_PUBLIC_USER_ID_BASE) return LIVE_CHAT_SOURCE.BOT;
	return LIVE_CHAT_SOURCE.USER;
}

/** 是否机器人弹幕（不计入弹幕次数/人数） */
export function isLiveChatBotMessage(item = {}) {
	return normalizeLiveChatRole(item.role, item.source) === LIVE_CHAT_ROLE.BOT;
}

/** IM 11003 → 聊天室展示结构 */
export function parseLiveDanmakuMessage(msg) {
	if (!msg || msg.code !== 11003) return null;
	const data = parseImPayload(msg.data);
	const text = String(data.text ?? data.content ?? '').trim();
	if (!text) return null;

	const userId = String(msg.sendUserId || data.userId || data.user_id || '').trim();
	const nickName =
		String(data.nickName || data.nick_name || data.nick || '').trim() || '用户';
	const avatar = String(data.avatar || data.avatarUrl || data.avatar_url || '').trim();
	const source = resolveLiveChatSource(data, userId);
	const messageId = Number(data.message_id || data.messageId || 0) || 0;

	return {
		id: messageId ? `d-${messageId}` : msg.messageId || `d-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
		messageId,
		type: 'user',
		userId,
		nickName,
		avatar,
		source,
		role: normalizeLiveChatRole(data.role, source),
		text,
	};
}

/** 中控进场补拉弹幕 API → 聊天室展示结构 */
export function mapDanmakuViewToChatItem(row = {}) {
	const text = String(row.content || row.text || '').trim();
	if (!text) return null;
	const userId = String(row.user_id || row.userId || '').trim();
	const messageId = Number(row.message_id || row.messageId || 0) || 0;
	const source = Number(row.source ?? 0);
	return {
		id: messageId ? `d-${messageId}` : `d-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
		messageId,
		type: 'user',
		userId,
		nickName: String(row.nick_name || row.nickName || '').trim() || '用户',
		avatar: String(row.avatar || '').trim(),
		source,
		role: normalizeLiveChatRole(row.role, source),
		text,
	};
}

/** IM 11027 / API danmaku → 审核列表项 */
export function mapImPayloadToAuditItem(raw) {
	const data = parseImPayload(raw);
	const messageId = Number(data.message_id || data.messageId || 0);
	if (!messageId) return null;
	const userId = Number(data.user_id || data.userId || 0);
	const nickName =
		String(data.nick_name || data.nickName || data.nick || '').trim() ||
		(userId ? `用户${userId}` : '用户');
	const content = String(data.content || data.text || '').trim();
	if (!content) return null;
	const role = normalizeLiveChatRole(data.role, data.source);
	return {
		message_id: messageId,
		user_id: userId,
		nick_name: nickName,
		avatar: String(data.avatar || data.avatarUrl || data.avatar_url || '').trim(),
		role,
		role_text: getLiveChatRoleLabel(role, data.source),
		content,
		audit_status: 0,
		audit_status_text: '待审核',
		block_type: Number(data.block_type || 0),
		block_type_text: '',
		send_time_text: '刚刚',
	};
}
