import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

/** 直播商品封面：优先 API 全 URL，相对键拼 COS。 */
export function resolveLiveProductImage(path) {
	return resolveCosMediaUrl(path);
}
