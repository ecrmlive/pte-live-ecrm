import { requestClient } from '#/api/request';
export interface PlatformSeckillActive { active_status: number; end_day: string; mer_id: number; mer_name?: string; name: string; product_id: number; seckill_active_id: number; seckill_price: number; start_day: string; status: number; store_name?: string; }
export interface PlatformSeckillPage { limit: number; list: PlatformSeckillActive[]; page: number; total: number; }
export function listPlatformSeckillApi(params: { limit: number; mer_id?: number; page: number }) { return requestClient.get<PlatformSeckillPage>('/seckill/actives', { params }); }
