/** 解析话术模版：一行一条；空行忽略；旧 @ 分隔数据兼容 */
export function parseDanmakuScriptLines(content: unknown) {
  let text = String(content || '')
    .replace(/\r\n/g, '\n')
    .replace(/\r/g, '\n');
  if (!text.includes('\n') && text.includes('@')) {
    text = text.split('@').join('\n');
  }
  return text
    .split('\n')
    .map((s) => s.trim())
    .filter(Boolean);
}

export function formatDanmakuScriptForEdit(content: unknown) {
  return parseDanmakuScriptLines(content).join('\n');
}

export function normalizeDanmakuScriptForSave(content: unknown) {
  return parseDanmakuScriptLines(content).join('\n');
}
