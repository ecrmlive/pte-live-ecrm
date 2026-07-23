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
      { path: 'content', redirect: '/content/notice' },
      { path: 'product', redirect: '/product/audit' },
      { path: 'order', redirect: '/order/list' },
      { path: 'merchant', redirect: '/merchant/list' },
      { path: 'accounts', redirect: '/accounts/withdraw' },
      { path: 'setting', redirect: '/setting/admin' },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '工作台' },
      },
      {
        path: 'merchant/list',
        name: 'MerchantList',
        component: () => import('@/views/merchant/list.vue'),
        meta: { title: '商户列表' },
      },
      {
        path: 'merchant/audit',
        name: 'MerchantAudit',
        component: () => import('@/views/merchant/audit.vue'),
        meta: { title: '入驻审核' },
      },
      {
        path: 'product/category',
        name: 'ProductCategory',
        component: () => import('@/views/product/category.vue'),
        meta: { title: '平台分类' },
      },
      {
        path: 'product/brand',
        name: 'ProductBrand',
        component: () => import('@/views/product/brand.vue'),
        meta: { title: '品牌管理' },
      },
      {
        path: 'product/audit',
        name: 'ProductAudit',
        component: () => import('@/views/product/audit.vue'),
        meta: { title: '商品审核' },
      },
      {
        path: 'order/list',
        name: 'OrderList',
        component: () => import('@/views/order/list.vue'),
        meta: { title: '订单监管' },
      },
      {
        path: 'order/refund',
        name: 'OrderRefund',
        component: () => import('@/views/order/refund.vue'),
        meta: { title: '退款监管' },
      },
      {
        path: 'accounts/withdraw',
        name: 'AccountsWithdraw',
        component: () => import('@/views/accounts/withdraw.vue'),
        meta: { title: '提现审核' },
      },
      {
        path: 'marketing/coupon',
        name: 'MarketingCoupon',
        component: () => import('@/views/marketing/coupon.vue'),
        meta: { title: '平台优惠券' },
      },
      {
        path: 'marketing/spread',
        name: 'MarketingSpread',
        component: () => import('@/views/marketing/spread.vue'),
        meta: { title: '分销监管' },
      },
      {
        path: 'content/notice',
        name: 'ContentNotice',
        component: () => import('@/views/content/notice.vue'),
        meta: { title: '平台公告' },
      },
      {
        path: 'content/diy',
        name: 'ContentDiy',
        component: () => import('@/views/content/diy.vue'),
        meta: { title: '页面装修' },
      },
      {
        path: 'content/attachment',
        name: 'ContentAttachment',
        component: () => import('@/views/content/attachment.vue'),
        meta: { title: '素材库' },
      },
      {
        path: 'marketing/seckill',
        name: 'MarketingSeckill',
        component: () => import('@/views/marketing/seckill.vue'),
        meta: { title: '秒杀监管' },
      },
      {
        path: 'marketing/combination',
        name: 'MarketingCombination',
        component: () => import('@/views/marketing/combination.vue'),
        meta: { title: '拼团监管' },
      },
      {
        path: 'marketing/presell',
        name: 'MarketingPresell',
        component: () => import('@/views/marketing/presell.vue'),
        meta: { title: '预售监管' },
      },
      {
        path: 'marketing/broadcast',
        name: 'MarketingBroadcast',
        component: () => import('@/views/marketing/broadcast.vue'),
        meta: { title: '直播监管' },
      },
      {
        path: 'marketing/assist',
        name: 'MarketingAssist',
        component: () => import('@/views/marketing/assist.vue'),
        meta: { title: '助力监管' },
      },
      {
        path: 'content/community',
        name: 'ContentCommunity',
        component: () => import('@/views/content/community.vue'),
        meta: { title: '社区种草' },
      },
      {
        path: 'freight/express',
        name: 'FreightExpress',
        component: () => import('@/views/freight/express.vue'),
        meta: { title: '物流公司' },
      },
      {
        path: 'cms/article',
        name: 'CmsArticle',
        component: () => import('@/views/cms/article.vue'),
        meta: { title: '文章管理' },
      },
      {
        path: 'user/label',
        name: 'UserLabel',
        component: () => import('@/views/user/label.vue'),
        meta: { title: '用户标签' },
      },
      {
        path: 'user/svip',
        name: 'UserSvip',
        component: () => import('@/views/user/svip.vue'),
        meta: { title: '付费会员' },
      },
      {
        path: 'setting/admin',
        name: 'SettingAdmin',
        component: () => import('@/views/setting/admin.vue'),
        meta: { title: '管理员' },
      },
      {
        path: 'setting/role',
        name: 'SettingRole',
        component: () => import('@/views/setting/role.vue'),
        meta: { title: '角色管理' },
      },
      {
        path: 'setting/menu',
        name: 'SettingMenu',
        component: () => import('@/views/setting/menu.vue'),
        meta: { title: '菜单管理' },
      },
      {
        path: 'setting/agreements',
        name: 'SettingAgreements',
        component: () => import('@/views/setting/agreements.vue'),
        meta: { title: '协议规则' },
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
