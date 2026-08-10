import {
  getPlatformProductEditApi,
  type PlatformProductEditDetail,
} from '#/api/core/platform-catalog';
import {
  getPlatformAssistApi,
  type PlatformAssistActive,
} from '#/api/core/platform-assist';

export type AssistProductBundle = {
  assist: PlatformAssistActive;
  product?: PlatformProductEditDetail;
  productMissing: boolean;
};

/** 加载助力活动 + 关联商品编辑信息（商品缺失时仍返回活动）。 */
export async function loadAssistProductBundle(
  id: number,
): Promise<AssistProductBundle> {
  const assist = await getPlatformAssistApi(id);
  const productId = Number(assist.product_id || 0);
  if (productId <= 0) {
    return { assist, product: undefined, productMissing: true };
  }
  try {
    const product = await getPlatformProductEditApi(productId);
    return { assist, product, productMissing: false };
  } catch {
    return { assist, product: undefined, productMissing: true };
  }
}
