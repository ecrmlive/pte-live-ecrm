/**
 * 腾讯云 TCPlayer 5.1（中控预览；直播用 HLS/FLV，录播仍 URL/VOD）
 */
const TCPLAYER_VERSION = '5.1.0';
const TCPLAYER_CSS = `https://web.sdk.qcloud.com/player/tcplayer/release/v${TCPLAYER_VERSION}/tcplayer.min.css`;
const TCPLAYER_JS = `https://web.sdk.qcloud.com/player/tcplayer/release/v${TCPLAYER_VERSION}/tcplayer.v${TCPLAYER_VERSION}.min.js`;

export const TCPLAYER_CONTROL_CONTAINER_ID = 'tcplayer-control-preview';

let sdkLoadPromise = null;

function loadStylesheet(href) {
	if (typeof document === 'undefined') return Promise.resolve();
	if (document.querySelector(`link[href="${href}"]`)) return Promise.resolve();
	return new Promise((resolve, reject) => {
		const link = document.createElement('link');
		link.rel = 'stylesheet';
		link.href = href;
		link.onload = () => resolve();
		link.onerror = () => reject(new Error('tcplayer css load failed'));
		document.head.appendChild(link);
	});
}

function loadScript(src) {
	if (typeof window !== 'undefined' && typeof window.TCPlayer === 'function') {
		return Promise.resolve();
	}
	if (document.querySelector(`script[src="${src}"]`)) {
		return new Promise((resolve, reject) => {
			const timer = setInterval(() => {
				if (typeof window.TCPlayer === 'function') {
					clearInterval(timer);
					resolve();
				}
			}, 50);
			setTimeout(() => {
				clearInterval(timer);
				if (typeof window.TCPlayer === 'function') resolve();
				else reject(new Error('tcplayer js timeout'));
			}, 15000);
		});
	}
	return new Promise((resolve, reject) => {
		const script = document.createElement('script');
		script.src = src;
		script.async = true;
		script.onload = () => resolve();
		script.onerror = () => reject(new Error('tcplayer js load failed'));
		document.head.appendChild(script);
	});
}

export function loadTcPlayerSDK() {
	if (sdkLoadPromise) return sdkLoadPromise;
	sdkLoadPromise = loadStylesheet(TCPLAYER_CSS)
		.then(() => loadScript(TCPLAYER_JS))
		.then(() => {
			if (typeof window.TCPlayer !== 'function') {
				throw new Error('TCPlayer unavailable');
			}
		});
	return sdkLoadPromise;
}

export function resolveTcPlayerContainer(id = TCPLAYER_CONTROL_CONTAINER_ID) {
	if (typeof document === 'undefined') return null;
	return document.getElementById(id);
}

/** TCPlayer dispose 后原 video 常被移除，挂载前在 slot 内重建 */
export function resetTcPlayerMountRoot(mountRoot, videoId, options = {}) {
	if (!mountRoot || !videoId) return null;
	mountRoot.innerHTML = '';
	const video = document.createElement('video');
	video.id = videoId;
	if (options.className) video.className = options.className;
	video.setAttribute('playsinline', '');
	video.setAttribute('webkit-playsinline', '');
	video.setAttribute('preload', 'auto');
	mountRoot.appendChild(video);
	return video;
}

export function ensureTcPlayerVideoElement(mountRoot, videoId, options = {}) {
	if (!mountRoot || !videoId) return null;
	let video = mountRoot.querySelector(`video#${videoId}`);
	if (video) return video;
	const byId = document.getElementById(videoId);
	if (byId?.tagName === 'VIDEO' && mountRoot.contains(byId)) return byId;
	return resetTcPlayerMountRoot(mountRoot, videoId, options);
}

export function guessSourceType(url) {
	if (!url) return 'webrtc';
	const lower = String(url).toLowerCase();
	if (lower.startsWith('webrtc://')) return 'webrtc';
	if (lower.includes('.flv') || lower.includes('format=flv')) return 'video/x-flv';
	if (lower.includes('.mp4') || lower.includes('format=mp4')) return 'video/mp4';
	if (lower.includes('.m3u8')) return 'application/x-mpegURL';
	return 'application/x-mpegURL';
}

/** Web 中控直播拉流地址：优先 HLS，其次 FLV；不使用 RTMP/WebRTC 兜底 */
export function resolveControlLiveSrc({ pullUrlHls = '', pullUrlFlv = '' } = {}) {
	const hls = String(pullUrlHls || '').trim();
	if (hls) return hls;
	return String(pullUrlFlv || '').trim();
}

function bindPlayerEvents(player, handlers = {}) {
	if (!player || typeof player.on !== 'function') return;
	const notifyPlaying = () => handlers.onPlaying?.();
	player.on('playing', notifyPlaying);
	player.on('play', notifyPlaying);
	player.on('loadeddata', notifyPlaying);
	player.on('canplay', notifyPlaying);
	player.on('loadedmetadata', () => {
		try {
			const el = player.el_?.querySelector?.('video') || player.el_;
			if (el && el.readyState >= 2) notifyPlaying();
		} catch {
			notifyPlaying();
		}
	});
	player.on('error', () => {
		suppressTcPlayerErrorUI(player, handlers.mountRoot);
		handlers.onError?.();
	});
	player.on('blocked', () => handlers.onAutoplayBlocked?.());
}

/** 隐藏 TCPlayer / video.js 内置错误层（如 License 域名校验 52） */
export function suppressTcPlayerErrorUI(player, mountRoot) {
	try {
		if (player && typeof player.error === 'function') {
			player.error(null);
		}
	} catch {
		// ignore
	}
	const roots = [];
	if (mountRoot) roots.push(mountRoot);
	if (typeof document !== 'undefined') {
		document.querySelectorAll('.tcplayer-wrap, .tcp-skin').forEach((node) => roots.push(node));
	}
	roots.forEach((root) => {
		root.querySelectorAll('.vjs-error-display, .vjs-modal-dialog, .vjs-error .vjs-big-play-button').forEach((el) => {
			el.style.display = 'none';
		});
	});
}

/**
 * @param {string} containerId DOM 容器 id（须为 <video> 元素，TCPlayer 5.x 要求）
 * @param {{ licenseUrl, licenseKey, src?, vodPlayback?, muted? }} options
 */
export async function createTcPlayer(containerId, options = {}, handlers = {}) {
	await loadTcPlayerSDK();
	const container =
		typeof containerId === 'object' && containerId?.tagName
			? containerId
			: resolveTcPlayerContainer(containerId);
	if (!container) throw new Error('tcplayer container missing');
	if (container.tagName !== 'VIDEO') {
		throw new Error('TCPlayer 容器须为 <video> 元素');
	}
	const playerId = container.id || (typeof containerId === 'string' ? containerId : '');
	if (!playerId) throw new Error('tcplayer container id missing');

	const licenseUrl = String(options.licenseUrl || '').trim();
	const licenseKey = String(options.licenseKey || '').trim();
	if (!licenseUrl || !licenseKey) {
		throw new Error('TCPlayer license 未配置');
	}

	const config = {
		licenseUrl,
		licenseKey,
		autoplay: true,
		muted: options.muted !== false,
		controls: false,
		playsinline: true,
		preload: 'auto',
		width: '100%',
		height: '100%',
		language: 'zh-CN',
		errorDisplay: false,
	};

	const vod = options.vodPlayback;
	if (vod && vod.file_id) {
		config.fileID = vod.file_id;
		config.appID = vod.app_id;
		if (vod.psign) config.psign = vod.psign;
	} else if (options.src) {
		const src = String(options.src).trim();
		const sourceType = guessSourceType(src);
		config.sources = [{ src, type: sourceType }];
		if (sourceType === 'webrtc') {
			config.webrtcConfig = { enableAbr: false };
		}
	} else {
		throw new Error('tcplayer source missing');
	}

	const player = window.TCPlayer(playerId, config);
	bindPlayerEvents(player, { ...handlers, mountRoot: container.parentElement || container });
	suppressTcPlayerErrorUI(player, container.parentElement || container);

	if (config.muted && typeof player.muted === 'function') {
		player.muted(true);
	}
	const playResult = typeof player.play === 'function' ? player.play() : null;
	if (playResult?.catch) {
		playResult.catch(() => handlers.onAutoplayBlocked?.());
	}
	return player;
}

export function disposeTcPlayer(player, mountRoot) {
	if (player) {
		try {
			if (typeof player.dispose === 'function') player.dispose();
			else if (typeof player.destroy === 'function') player.destroy();
		} catch {
			// ignore
		}
	}
	if (!mountRoot) return;
	mountRoot.querySelectorAll('.tcplayer-wrap, .tcp-skin, .vjs-tech').forEach((node) => {
		try {
			node.remove();
		} catch {
			// ignore
		}
	});
}
