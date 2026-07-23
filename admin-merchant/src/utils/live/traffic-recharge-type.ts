/** 流量充值类型（与 live-api recharge_type 一致） */
export const RECHARGE_TYPE_LABELS: Record<string, string> = {
  adjust: '调账',
  gift: '赠送',
  initial: '开户',
  purchase: '购买',
  reduce: '减少',
};

export function rechargeTypeLabel(value: unknown) {
  const key = String(value ?? '').trim();
  return RECHARGE_TYPE_LABELS[key] || key || '—';
}
