import type { ShopLinkValue } from '#/types/shop-link';

export function finalizeShopLinkValue(link: ShopLinkValue | null | undefined) {
  if (!link) return null;
  const result = { ...link };
  const url = String(result.url ?? '');
  if (url.includes('giftpackage') && result.name) {
    result.name = `礼包购-${result.name}`;
  }
  if (url.includes('invite') && result.name) {
    result.name = `邀请有礼-${result.name}`;
  }
  return result;
}
