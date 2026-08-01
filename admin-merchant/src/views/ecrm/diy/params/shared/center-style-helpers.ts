/** Helpers for center DIY style.type buttons and preview reactivity. */

export function ensureDiyStyle(curItem: Record<string, unknown>): Record<string, unknown> {
  if (!curItem.style || typeof curItem.style !== 'object') {
    curItem.style = {};
  }
  return curItem.style as Record<string, unknown>;
}

export function diyStyleType(curItem: Record<string, unknown>): number {
  const raw = ensureDiyStyle(curItem).type;
  const type = Number(raw);
  return Number.isFinite(type) && type > 0 ? type : 1;
}

export function setDiyStyleType(
  curItem: Record<string, unknown>,
  type: number,
  resetOnType1 = false,
): void {
  const style = ensureDiyStyle(curItem);
  style.type = type;
  if (resetOnType1 && type === 1) {
    style.paddingTop = 0;
    style.paddingBottom = 0;
    style.paddingLeft = 0;
    style.topRadio = 0;
  }
}
