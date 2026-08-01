import type { RouteRecordRaw } from 'vue-router';

/** 全屏独立页：DIY 编辑器不走 BasicLayout */
const standaloneRoutes: RouteRecordRaw[] = [
  {
    name: 'MerchantDiyEditor',
    path: '/devise/diy/index',
    component: () => import('#/views/ecrm/diy/editor.vue'),
    meta: {
      hideInBreadcrumb: true,
      hideInMenu: true,
      hideInTab: true,
      noBasicLayout: true,
      title: '页面装修',
    },
  },
];

export default standaloneRoutes;
