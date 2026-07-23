/**
 * 直播间弹幕区配色（商户中控 / 主播端 / C 端统一常量）
 *
 * 角色标签/ID/昵称色见 DANMAKU_ROLE_COLORS；正文 rgba(255,255,255,0.9)
 * 系统：进房欢迎 / 讲解 / 平台提示 见 resolveSystemDanmakuStyle
 */

export const DANMAKU_CONTENT_COLOR = 'rgba(255, 255, 255, 0.9)';
export const DANMAKU_SYSTEM_BODY_COLOR = 'rgba(255, 255, 255, 0.6)';

/** @deprecated 使用 DANMAKU_CONTENT_COLOR */
export const DANMAKU_TEXT_COLOR = DANMAKU_CONTENT_COLOR;

export const DANMAKU_ROLE_COLORS = {
	system: '#00D4AA',
	anchor: '#FBBF24',
	admin: '#1A5CFF',
	viewer: '#A78BFA',
	bot: '#EC4899',
};

export const DANMAKU_SYSTEM_LABEL_COLORS = {
	welcome: '#F472B6',
	explain: '#FBBF24',
};

const ROLE_LABELS = {
	0: '观众',
	1: '管理员',
	2: '主播',
	3: '机器人',
};

const ENTER_WELCOME_SUFFIX = '进入直播间🎉';

function normalizeRole(role, source) {
	if (Number(source) === 2) return 3;
	const n = Number(role);
	if (n === 1 || n === 2 || n === 3) return n;
	if (typeof role === 'string') {
		const key = role.trim().toLowerCase();
		if (key === 'admin' || key === 'manager') return 1;
		if (key === 'anchor' || key === 'host') return 2;
		if (key === 'bot' || key === 'robot') return 3;
	}
	return 0;
}

export function getDanmakuRoleColor(role, source) {
	const r = normalizeRole(role, source);
	if (r === 2) return DANMAKU_ROLE_COLORS.anchor;
	if (r === 1) return DANMAKU_ROLE_COLORS.admin;
	if (r === 3) return DANMAKU_ROLE_COLORS.bot;
	return DANMAKU_ROLE_COLORS.viewer;
}

export function getDanmakuRoleLabel(role, source, roleText = '') {
	const text = String(roleText || '').trim();
	if (text) return text;
	return ROLE_LABELS[normalizeRole(role, source)] || '观众';
}

export function resolveSystemDanmakuStyle(type = 'default') {
	if (type === 'welcome') {
		return {
			systemLabel: '系统',
			systemLabelColor: DANMAKU_SYSTEM_LABEL_COLORS.welcome,
			nickColor: DANMAKU_SYSTEM_LABEL_COLORS.welcome,
			textColor: DANMAKU_SYSTEM_BODY_COLOR,
			bodyColor: DANMAKU_SYSTEM_BODY_COLOR,
			showBadge: false,
			system: true,
			systemType: 'welcome',
		};
	}
	if (type === 'explain') {
		return {
			systemLabel: '系统',
			systemLabelColor: DANMAKU_SYSTEM_LABEL_COLORS.explain,
			nickColor: DANMAKU_SYSTEM_LABEL_COLORS.explain,
			textColor: DANMAKU_SYSTEM_BODY_COLOR,
			bodyColor: DANMAKU_SYSTEM_BODY_COLOR,
			showBadge: false,
			system: true,
			systemType: 'explain',
		};
	}
	return {
		systemLabel: type === 'notice' ? '' : '系统',
		systemLabelColor: DANMAKU_ROLE_COLORS.system,
		nickColor: DANMAKU_ROLE_COLORS.system,
		textColor: DANMAKU_ROLE_COLORS.system,
		bodyColor: DANMAKU_ROLE_COLORS.system,
		showBadge: false,
		system: true,
		systemType: type,
	};
}

export function buildEnterWelcomeSystemMessage(data = {}, sendUserId = '') {
	const userId = String(data.userId || data.user_id || sendUserId || '').trim();
	const nickName =
		String(data.nickName || data.nick_name || data.nick || '').trim() ||
		(userId ? `用户${userId}` : '用户');
	return {
		id: `sys-enter-${userId || Date.now()}`,
		type: 'system',
		systemType: 'welcome',
		systemLabel: '系统',
		systemLabelColor: DANMAKU_SYSTEM_LABEL_COLORS.welcome,
		text: `欢迎 ${nickName} ${ENTER_WELCOME_SUFFIX}`,
		bodyColor: DANMAKU_SYSTEM_BODY_COLOR,
		textColor: DANMAKU_SYSTEM_BODY_COLOR,
	};
}

export function buildProductExplainSystemMessage(data = {}) {
	const name = String(data.name || data.product_name || data.productName || '').trim() || '商品';
	return {
		id: `sys-explain-${Date.now()}`,
		type: 'system',
		systemType: 'explain',
		systemLabel: '系统',
		systemLabelColor: DANMAKU_SYSTEM_LABEL_COLORS.explain,
		text: `主播正在讲解商品：${name}`,
		bodyColor: DANMAKU_SYSTEM_BODY_COLOR,
		textColor: DANMAKU_SYSTEM_BODY_COLOR,
	};
}

/** IM 11028 → 进房欢迎系统弹幕（中控聊天室） */
export function parseEnterWelcomeMessage(msg) {
	if (!msg || Number(msg.code) !== 11028) return null;
	const data =
		typeof msg.data === 'string'
			? (() => {
					try {
						return JSON.parse(msg.data);
					} catch {
						return {};
					}
				})()
			: msg.data || {};
	return buildEnterWelcomeSystemMessage(data, msg.sendUserId);
}
