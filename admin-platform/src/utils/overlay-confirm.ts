import type { ElMessageBoxOptions } from 'element-plus';

import { ElMessageBox } from 'element-plus';

/** 取当前页面最高 overlay z-index，保证确认框在 Vben Modal 之上 */
function resolveOverlayZIndex() {
  if (typeof document === 'undefined') {
    return 3000;
  }
  let max = 2000;
  for (const el of document.querySelectorAll<HTMLElement>('[data-dismissable-modal], .el-overlay')) {
    const z = Number.parseInt(getComputedStyle(el).zIndex || '0', 10);
    if (z > max) {
      max = z;
    }
  }
  return max + 100;
}

/** 嵌套在 Modal 内时使用，避免 ElMessageBox 被遮挡导致「点击无反应」 */
export function overlayConfirm(
  message: string,
  title: string,
  options: ElMessageBoxOptions = {},
) {
  return ElMessageBox.confirm(message, title, {
    ...options,
    appendTo: document.body,
    customClass: ['platform-overlay-message-box', options.customClass]
      .filter(Boolean)
      .join(' '),
    zIndex: options.zIndex ?? resolveOverlayZIndex(),
  });
}
