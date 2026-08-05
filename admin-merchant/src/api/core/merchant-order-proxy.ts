import { requestClient } from '#/api/request';

function newIdempotencyKey() {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return `proxy:${crypto.randomUUID()}`;
  }
  return `proxy:${Date.now()}-${Math.random().toString(16).slice(2)}`;
}

export function createProxyOrderApi(body: {
  product_id: number;
  quantity: number;
  remark?: string;
  user_id: number;
  idempotency_key?: string;
}) {
  return requestClient.post<{
    group_order_id: number;
    order_id: number;
    order_sn: string;
    pay_amount: number;
    paid: number;
    replayed: boolean;
  }>('/orders/proxy', {
    ...body,
    idempotency_key: body.idempotency_key || newIdempotencyKey(),
  });
}
