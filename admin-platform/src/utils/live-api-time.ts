/**
 * api-live 时间字段 — 客户端展示
 * 时间展示与服务端返回保持一致。
 */

function pad2(n: number) {
  return String(n).padStart(2, '0');
}

export function parseLiveApiTime(value: null | number | string | undefined) {
  if (value == null || value === '') {
    return null;
  }
  if (typeof value === 'number' && Number.isFinite(value)) {
    const ms = value > 1e12 ? value : value * 1000;
    const d = new Date(ms);
    return Number.isNaN(d.getTime()) ? null : d;
  }
  const str = String(value).trim();
  if (!str || str.startsWith('0000-00-00')) {
    return null;
  }
  if (/^\d+$/.test(str)) {
    const n = Number(str);
    const ms = str.length > 10 ? n : n * 1000;
    const d = new Date(ms);
    return Number.isNaN(d.getTime()) ? null : d;
  }
  const d = new Date(str);
  return Number.isNaN(d.getTime()) ? null : d;
}

/** 列表、详情等完整时间 */
export function formatLiveApiDateTime(value: null | number | string | undefined) {
  const d = parseLiveApiTime(value);
  if (!d) {
    return '-';
  }
  return `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())} ${pad2(d.getHours())}:${pad2(d.getMinutes())}:${pad2(d.getSeconds())}`;
}

const SETTLEMENT_STATUS_LABELS: Record<string, string> = {
  pending: '待结算',
  done: '已完成',
  failed: '失败',
  unknown: '未知',
};

export function formatSettlementStatus(value: null | string | undefined) {
  if (!value) {
    return '-';
  }
  return SETTLEMENT_STATUS_LABELS[value] ?? SETTLEMENT_STATUS_LABELS.unknown!;
}
