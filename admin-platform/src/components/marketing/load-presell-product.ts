import {
  getPlatformProductEditApi,
  type PlatformProductEditDetail,
} from '#/api/core/platform-catalog';
import {
  getPlatformPresellApi,
  type PlatformPresell,
} from '#/api/core/platform-presell';

export type PresellProductBundle = {
  presell: PlatformPresell;
  product?: PlatformProductEditDetail;
  productMissing: boolean;
};

/** 加载预售详情 + 关联商品编辑信息（商品缺失时仍返回预售）。 */
export async function loadPresellProductBundle(
  id: number,
): Promise<PresellProductBundle> {
  const presell = await getPlatformPresellApi(id);
  const productId = Number(presell.product_id || 0);
  if (productId <= 0) {
    return { presell, product: undefined, productMissing: true };
  }
  try {
    const product = await getPlatformProductEditApi(productId);
    return { presell, product, productMissing: false };
  } catch {
    return { presell, product: undefined, productMissing: true };
  }
}
