import DOMPurify from 'dompurify';

/** 富文本 XSS 过滤（DIY 预览等） */
export function sanitizeHtml(html: string) {
  if (!html) return '';
  return DOMPurify.sanitize(html, {
    ALLOWED_ATTR: [
      'href',
      'src',
      'alt',
      'title',
      'class',
      'style',
      'target',
      'width',
      'height',
    ],
    ALLOWED_TAGS: [
      'p',
      'br',
      'strong',
      'em',
      'u',
      'h1',
      'h2',
      'h3',
      'h4',
      'h5',
      'h6',
      'ul',
      'ol',
      'li',
      'a',
      'img',
      'table',
      'thead',
      'tbody',
      'tr',
      'th',
      'td',
      'div',
      'span',
      'blockquote',
      'code',
      'pre',
    ],
  });
}

export function sanitizeText(text: string) {
  if (!text) return '';
  return DOMPurify.sanitize(text, {
    ALLOWED_ATTR: [],
    ALLOWED_TAGS: [],
  });
}
