/** 与 C 端 main.js / uni.scss 主题色一致（页面 → 主题设置） */
export const SHOP_THEME_COLORS = [
  '#fd6303',
  '#42ca4d',
  '#ff448f',
  '#409eff',
  '#101012',
  '#e2aa62',
  '#a253ff',
  '#fe0024',
] as const;

export function getThemeColorByIndex(theme = '0') {
  const idx = Number.parseInt(String(theme), 10);
  if (Number.isNaN(idx) || idx < 0 || idx >= SHOP_THEME_COLORS.length) {
    return SHOP_THEME_COLORS[0];
  }
  return SHOP_THEME_COLORS[idx];
}
