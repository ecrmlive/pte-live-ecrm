/**
 * api-live 时间字段 — 客户端展示（勿原样显示 API 字符串）
 * 时间展示与服务端返回保持一致。
 */

function pad2(n: number) {
  return String(n).padStart(2, '0');
}

/** 解析 Unix 秒/毫秒、ISO 字符串、YYYY-MM-DD */
export function parseLiveApiTime(value: unknown): Date | null {
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

const DAY_CHART_RANGES = new Set(['7d', '30d', 'custom']);

export function formatLiveApiChartAxis(
  value: unknown,
  range: string,
): string {
  const d = parseLiveApiTime(value);
  if (!d) {
    return value == null ? '' : String(value);
  }
  if (DAY_CHART_RANGES.has(range)) {
    return `${pad2(d.getMonth() + 1)}-${pad2(d.getDate())}`;
  }
  return `${pad2(d.getHours())}:${pad2(d.getMinutes())}`;
}

export function formatLiveApiDateTime(value: unknown): string {
  const d = parseLiveApiTime(value);
  if (!d) {
    return '-';
  }
  return `${d.getFullYear()}-${pad2(d.getMonth() + 1)}-${pad2(d.getDate())} ${pad2(d.getHours())}:${pad2(d.getMinutes())}:${pad2(d.getSeconds())}`;
}

export function formatLiveApiChartAxisList(
  times: unknown,
  range: string,
): string[] {
  if (!Array.isArray(times)) {
    return [];
  }
  return times.map((t) => formatLiveApiChartAxis(t, range));
}

export function resolveLiveStatChartRange(
  stats: { range?: string } | null | undefined,
  activeRange: string,
  isSingleMode: boolean,
): string {
  if (isSingleMode) {
    return 'today';
  }
  return stats?.range || activeRange || '7d';
}
