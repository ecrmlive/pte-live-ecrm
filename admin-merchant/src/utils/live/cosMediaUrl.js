/** API 已出库全 URL；相对对象键仍拼 COS 域名（DIY 默认项等）。 */
const COS_BASE = (import.meta.env.VITE_COS_URL || 'https://cos.qxkejiwl.top').replace(/\/$/, '');

export function resolveCosMediaUrl(path = '') {
	const s = String(path || '').trim();
	if (!s) {
		return '';
	}
	if (/^https?:\/\//i.test(s)) {
		return s;
	}
	return `${COS_BASE}/${s.replace(/^\//, '')}`;
}

/** 保存时原样提交完整 URL（live-api / PHP 入库）。 */
export function toCosRelativePath(fileOrPath = '') {
	if (fileOrPath && typeof fileOrPath === 'object') {
		const url = fileOrPath.file_path || fileOrPath.file_name || '';
		return resolveCosMediaUrl(url);
	}
	return resolveCosMediaUrl(fileOrPath);
}
