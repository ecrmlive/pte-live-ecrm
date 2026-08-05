/** 管理后台统一时间展示：Asia/Shanghai · yyyy-MM-dd HH:mm:ss */
export const ADMIN_DATETIME_FORMAT = 'yyyy-MM-dd HH:mm:ss';
export const ADMIN_TIMEZONE = 'Asia/Shanghai';

export function formatShanghaiDateTime(value?: Date | null | string): string {
  if (!value) return '—';
  const date = value instanceof Date ? value : new Date(value);
  if (Number.isNaN(date.valueOf())) return String(value);
  const parts = new Intl.DateTimeFormat('zh-CN', {
    day: '2-digit',
    hour: '2-digit',
    hour12: false,
    minute: '2-digit',
    month: '2-digit',
    second: '2-digit',
    timeZone: ADMIN_TIMEZONE,
    year: 'numeric',
  })
    .formatToParts(date)
    .reduce<Record<string, string>>((result, part) => {
      result[part.type] = part.value;
      return result;
    }, {});
  return `${parts.year}-${parts.month}-${parts.day} ${parts.hour}:${parts.minute}:${parts.second}`;
}
