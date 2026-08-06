import type { RouteRecordRaw } from 'vue-router';

/** 全屏独立页：不走 BasicLayout（无侧栏/顶栏 chrome） */
const standaloneRoutes: RouteRecordRaw[] = [
  {
    name: 'PlatformDiyEditor',
    path: '/setting/diy/index',
    component: () => import('#/views/ecrm/diy/editor.vue'),
    meta: {
      hideInBreadcrumb: true,
      hideInMenu: true,
      hideInTab: true,
      noBasicLayout: true,
      title: '页面装修',
    },
  },
  {
    name: 'PlatformDataScreen',
    path: '/data-screen/index',
    component: () => import('#/views/ecrm/dashboard/data-screen.vue'),
    meta: {
      hideInBreadcrumb: true,
      hideInMenu: true,
      hideInTab: true,
      noBasicLayout: true,
      // 侧栏点击走 useNavigation：新开浏览器标签，不嵌入当前 admin Tab
      openInNewWindow: true,
      title: '数据大屏',
    },
  },
];

/** 菜单动态路由里需排除的路径，避免再挂到 BasicLayout 下出现双层壳 */
export const STANDALONE_ROUTE_PATHS = new Set(
  standaloneRoutes.map((route) => route.path),
);

export default standaloneRoutes;
