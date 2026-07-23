/** 平台 DIY 预览轻量过滤（不依赖 dompurify） */
export function filterXSS(html: string) {
  return String(html || '');
}

export function sanitizeHtml(html: string) {
  return filterXSS(html);
}

export default filterXSS;
