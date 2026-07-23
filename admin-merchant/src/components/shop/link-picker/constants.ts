import type { ShopLinkPickerItem } from '#/api/core/shop-link';

export const MP_BUILTIN_PAGES: ShopLinkPickerItem[] = [
  { url: 'pages/main/index/index', name: '首页', type: '页面' },
  { url: 'pages/mall/product/category', name: '分类', type: '页面' },
  { url: 'pages/content/article/list/list', name: '文章首页', type: '页面' },
  { url: 'pages/mall/cart/cart', name: '购物车', type: '页面' },
  { url: 'pages/order/myorder', name: '订单', type: '页面' },
  { url: 'pages/member/user/index/index', name: '我的', type: '页面' },
  { url: 'pages/live/room/list', name: '直播列表', type: '页面' },
  { url: 'pages/order/code/code', name: '提货页面', type: '页面' },
  { url: 'pages/order/codeorder', name: '提货订单', type: '页面' },
];

export const USER_CENTER_MENU_LINKS: ShopLinkPickerItem[] = [
  { url: '/pages/member/user/address/address', name: '收货地址', type: '菜单' },
  { url: '/pages/member/coupon/coupon', name: '领券中心', type: '菜单' },
  { url: '/pages/member/user/my-coupon/my-coupon', name: '我的优惠券', type: '菜单' },
  { url: '/pages/member/agent/index/index', name: '分销中心', type: '菜单' },
  { url: '/pages/member/user/my-bargain/my-bargain', name: '我的砍价', type: '菜单' },
  { url: '/pages/marketing/giftpackage/giftlist', name: '我的礼包', type: '菜单' },
  { url: '/pages/member/user/favorite/favorite', name: '我的收藏', type: '菜单' },
  { url: '/pages/order/assemble-order', name: '我的拼团', type: '菜单' },
  { url: '/pages/marketing/signin/signin', name: '签到有礼', type: '菜单' },
  { url: '/pages/member/task/index', name: '任务中心', type: '菜单' },
  { url: '/pages/member/user/evaluate/list', name: '我的评价', type: '菜单' },
  { url: '/pages/member/user/set/set', name: '设置', type: '菜单' },
  { url: '/pages/member/card/index', name: '我的等级', type: '菜单' },
];

export const MARKETING_STATIC_LINKS: Record<string, ShopLinkPickerItem[]> = {
  assemble: [{ id: 0, url: 'pages/marketing/assemble/list/list', name: '拼团', type: '营销' }],
  bargain: [{ id: 0, url: 'pages/marketing/bargain/list/list', name: '砍价', type: '营销' }],
  coupon: [{ id: 0, url: 'pages/member/coupon/coupon', name: '优惠券', type: '营销' }],
  lottery: [{ id: 0, url: 'pages/marketing/lottery/lottery', name: '幸运转盘', type: '营销' }],
  points: [{ id: 0, url: 'pages/marketing/points/list/list', name: '积分商城', type: '营销' }],
  presale: [{ id: 0, url: 'pages/marketing/presale/list', name: '预售', type: '营销' }],
  preview: [{ id: 0, url: 'pages/marketing/preview/list', name: '预告', type: '营销' }],
  recharge: [{ id: 0, url: 'pages/order/recharge', name: '充值', type: '营销' }],
  seckill: [{ id: 0, url: 'pages/marketing/seckill/list/list', name: '秒杀', type: '营销' }],
  signin: [{ id: 0, url: 'pages/marketing/signin/signin', name: '签到', type: '营销' }],
  task: [{ id: 0, url: 'pages/member/task/index', name: '任务中心', type: '营销' }],
};
