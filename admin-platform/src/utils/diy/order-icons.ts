/** DIY 订单组件图标（src/assets/diy/img/order） */
function orderIconUrl(style: number, index: number) {
  return new URL(
    `../../assets/diy/img/order/${style}-${index}.png`,
    import.meta.url,
  ).href;
}

export const DIY_ORDER_ICON_LIST: string[][] = Array.from({ length: 10 }, (_, style) =>
  Array.from({ length: 5 }, (_, index) => orderIconUrl(style + 1, index)),
);

export const DIY_CENTER_BG_URL = new URL(
  '../../assets/diy/img/center-bg.png',
  import.meta.url,
).href;

export const DIY_LOGIN_NAME_URL = new URL(
  '../../assets/diy/img/login-name.png',
  import.meta.url,
).href;

export const DIY_COUPON_BG_URL = new URL(
  '../../assets/diy/img/diy-coupon.png',
  import.meta.url,
).href;

export const DIY_SECKILL_ICON_URL = new URL(
  '../../assets/diy/img/seckill.png',
  import.meta.url,
).href;
