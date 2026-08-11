import { requestClient } from '#/api/request';

/** 微信自定义菜单按钮（对齐微信公众平台 button / sub_button） */
export interface WechatMenuButton {
  type?: 'click' | 'view' | 'miniprogram' | string;
  name: string;
  key?: string;
  url?: string;
  appid?: string;
  pagepath?: string;
  sub_button?: WechatMenuButton[];
}

export function getWechatMenusApi() {
  return requestClient.get<{
    wechat_menus: WechatMenuButton[];
    note?: string;
  }>('/setting/wechat-menus');
}

export function saveWechatMenusApi(button: WechatMenuButton[]) {
  return requestClient.put<{
    wechat_menus: WechatMenuButton[];
    published?: boolean;
    note?: string;
  }>('/setting/wechat-menus', { button });
}
