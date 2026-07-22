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
    component: () => import('@/layouts/Shell.vue'),
    redirect: '/verify',
    children: [
      {
        path: 'verify',
        name: 'VerifyList',
        component: () => import('@/views/verify/list.vue'),
        meta: { title: '待核销' },
      },
      {
        path: 'verify/:id',
        name: 'VerifyDetail',
        component: () => import('@/views/verify/detail.vue'),
        meta: { title: '核销详情' },
      },
      {
        path: 'refund',
        name: 'RefundList',
        component: () => import('@/views/refund/list.vue'),
        meta: { title: '代退处理' },
      },
      {
        path: 'delivery',
        name: 'DeliveryList',
        component: () => import('@/views/delivery/list.vue'),
        meta: { title: '待发货' },
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
      return '/verify';
    }
    return true;
  }
  if (!tokenStore.getAccess()) {
    return { path: '/login', query: { redirect: to.fullPath } };
  }
  return true;
});

export default router;
