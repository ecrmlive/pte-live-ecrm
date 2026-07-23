import * as XLSX from 'xlsx';

import { parseDanmakuScriptLines } from './danmaku-bot-script-lines';

const HEADER_LABELS = new Set(['话术', '内容', '弹幕', '话术内容', 'script', 'content']);

function cellText(row: unknown) {
  if (!Array.isArray(row)) return '';
  const value = row[0];
  if (value == null) return '';
  return String(value).trim();
}

function isHeaderRow(text: string) {
  if (!text) return false;
  const normalized = text.toLowerCase();
  return HEADER_LABELS.has(text) || HEADER_LABELS.has(normalized);
}

export async function parseDanmakuScriptExcelFile(file: File) {
  if (!file) {
    throw new Error('请选择 Excel 文件');
  }
  const ext = String(file.name || '').split('.').pop()?.toLowerCase();
  if (ext !== 'xlsx' && ext !== 'xls') {
    throw new Error('仅支持 .xlsx / .xls 文件');
  }
  if (file.size / 1024 / 1024 > 5) {
    throw new Error('文件大小不能超过 5MB');
  }

  const buffer = await file.arrayBuffer();
  const workbook = XLSX.read(buffer, { type: 'array' });
  const sheetName = workbook.SheetNames?.[0];
  if (!sheetName) {
    throw new Error('Excel 中没有工作表');
  }
  const rows = XLSX.utils.sheet_to_json<(string | number)[]>(workbook.Sheets[sheetName]!, {
    blankrows: false,
    defval: '',
    header: 1,
  });

  const lines: string[] = [];
  rows.forEach((row, index) => {
    const text = cellText(row);
    if (!text) return;
    if (index === 0 && isHeaderRow(text)) return;
    lines.push(text);
  });
  return lines;
}

export function mergeDanmakuScriptLines(
  existingContent: unknown,
  importedLines: string[],
) {
  const existing = parseDanmakuScriptLines(existingContent);
  const seen = new Set(existing);
  const merged = [...existing];
  importedLines.forEach((line) => {
    const t = String(line || '').trim();
    if (!t || seen.has(t)) return;
    seen.add(t);
    merged.push(t);
  });
  return merged.join('\n');
}

export function downloadDanmakuScriptExcelTemplate() {
  const sheet = XLSX.utils.aoa_to_sheet([
    ['话术'],
    ['欢迎光临直播间'],
    ['这个商品今天特价'],
    ['库存不多抓紧下单'],
  ]);
  sheet['!cols'] = [{ wch: 40 }];
  const book = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(book, sheet, '话术');
  XLSX.writeFile(book, '话术模版导入模板.xlsx');
}
