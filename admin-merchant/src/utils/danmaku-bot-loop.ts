export function isDanmakuLoopEnabled(value: unknown) {
  return Number(value) === 1;
}

export function formatDanmakuLoopEnabled(value: unknown) {
  return isDanmakuLoopEnabled(value) ? '循环执行' : '不循环执行';
}

export function normalizeDanmakuLoopEnabled(value: unknown, defaultValue = 0) {
  if (value === undefined || value === null || value === '') {
    return Number(defaultValue) === 1 ? 1 : 0;
  }
  return isDanmakuLoopEnabled(value) ? 1 : 0;
}
