import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { tokenStore } from '@/api/http';
import { useAuthStore } from '@/stores/auth';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('@/layouts/BasicLayout.vue'),
    redirect: '/dashboard',
    children: [
      { path: 'marketing', redirect: '/marketing/coupon' },
      { path: 'product', redirect: '/product/list' },
      { path: 'order', redirect: '/order/list' },
      { path: 'finance', redirect: '/finance/balance' },
      { path: 'setting', redirect: '/setting/shop' },
      {
        path: 'dashboard',
        name: 'MerDashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '工作台' },
      },
      {
        path: 'product/list',
        name: 'MerProductList',
        component: () => import('@/views/product/list.vue'),
        meta: { title: '商品列表' },
      },
      {
        path: 'product/edit',
        name: 'MerProductEdit',
        component: () => import('@/views/product/edit.vue'),
        meta: { title: '发布商品' },
      },
      {
        path: 'order/list',
        name: 'MerOrderList',
        component: () => import('@/views/order/list.vue'),
        meta: { title: '订单列表' },
      },
      {
        path: 'order/delivery',
        name: 'MerOrderDelivery',
        component: () => import('@/views/order/list.vue'),
        meta: { title: '待发货' },
      },
      {
        path: 'order/refund',
        name: 'MerOrderRefund',
        component: () => import('@/views/order/refund.vue'),
        meta: { title: '售后处理' },
      },
      {
        path: 'finance/balance',
        name: 'MerFinanceBalance',
        component: () => import('@/views/finance/balance.vue'),
        meta: { title: '店铺余额' },
      },
      {
        path: 'finance/withdraw',
        name: 'MerFinanceWithdraw',
        component: () => import('@/views/finance/withdraw.vue'),
        meta: { title: '提现申请' },
      },
      {
        path: 'marketing/coupon',
        name: 'MerMarketingCoupon',
        component: () => import('@/views/marketing/coupon.vue'),
        meta: { title: '店铺优惠券' },
      },
      {
        path: 'marketing/seckill',
        name: 'MerMarketingSeckill',
        component: () => import('@/views/marketing/seckill.vue'),
        meta: { title: '秒杀活动' },
      },
      {
        path: 'marketing/combination',
        name: 'MerMarketingCombination',
        component: () => import('@/views/marketing/combination.vue'),
        meta: { title: '拼团活动' },
      },
      {
        path: 'marketing/reservation',
        name: 'MerMarketingReservation',
        component: () => import('@/views/marketing/reservation.vue'),
        meta: { title: '预约服务' },
      },
      {
        path: 'marketing/presell',
        name: 'MerMarketingPresell',
        component: () => import('@/views/marketing/presell.vue'),
        meta: { title: '预售活动' },
      },
      {
        path: 'marketing/broadcast',
        name: 'MerMarketingBroadcast',
        component: () => import('@/views/marketing/broadcast.vue'),
        meta: { title: '直播间' },
      },
      {
        path: 'marketing/assist',
        name: 'MerMarketingAssist',
        component: () => import('@/views/marketing/assist.vue'),
        meta: { title: '助力活动' },
      },
      {
        path: 'community/list',
        redirect: '/marketing/community',
      },
      {
        path: 'marketing/community',
        name: 'MerCommunityList',
        component: () => import('@/views/community/list.vue'),
        meta: { title: '逛逛社区' },
      },
      {
        path: 'setting/svip',
        name: 'MerSettingSvip',
        component: () => import('@/views/setting/svip.vue'),
        meta: { title: '付费会员配置' },
      },
      {
        path: 'setting/shop',
        name: 'MerSettingShop',
        component: () => import('@/views/setting/shop.vue'),
        meta: { title: '店铺资料' },
      },
      {
        path: 'setting/staff',
        name: 'MerSettingStaff',
        component: () => import('@/views/setting/staff.vue'),
        meta: { title: '店员核销' },
      },
      {
        path: 'setting/admins',
        name: 'MerSettingAdmins',
        component: () => import('@/views/setting/admin.vue'),
        meta: { title: '子账号' },
      },
      {
        path: 'setting/admin',
        redirect: '/setting/admins',
      },
      {
        path: 'setting/attachment',
        name: 'MerSettingAttachment',
        component: () => import('@/views/setting/attachment.vue'),
        meta: { title: '素材库' },
      },
      {
        path: 'setting/replies',
        name: 'MerSettingReplies',
        component: () => import('@/views/setting/replies.vue'),
        meta: { title: '快捷回复' },
      },
      {
        path: 'setting/role',
        name: 'MerSettingRole',
        component: () => import('@/views/setting/role.vue'),
        meta: { title: '角色权限' },
      },
      {
        path: ':pathMatch(.*)*',
        name: 'Placeholder',
        component: () => import('@/views/placeholder/index.vue'),
        meta: { title: '功能建设中' },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to) => {
  const auth = useAuthStore();
  if (!auth.booted) {
    await auth.bootstrap();
  }
  if (to.meta.public) {
    if (tokenStore.getAccess() && to.path === '/login') {
      return '/dashboard';
    }
    return true;
  }
  if (!tokenStore.getAccess()) {
    return { path: '/login', query: { redirect: to.fullPath } };
  }
  return true;
});

export default router;
