import { createApp, watchEffect } from 'vue';

import { registerAccessDirective } from '@vben/access';
import { registerLoadingDirective, setDefaultDrawerProps } from '@vben/common-ui';
import { preferences } from '@vben/preferences';
import { initStores, useAccessStore } from '@vben/stores';
import '@vben/styles';
import '@vben/styles/ele';

import { useTitle } from '@vueuse/core';
import { ElLoading } from 'element-plus';

import { $t, setupI18n } from '#/locales';

import { initComponentAdapter } from './adapter/component';
import { initSetupVbenForm } from './adapter/form';
import App from './app.vue';
import { registerMerchantDirectives } from './directives/merchant';
import { router } from './router';
import { registerMerchantShell } from './adapters/merchant-shell';
import { restoreShopSessionBootstrappedFromStorage } from '#/utils/jwt-session';
import {
  getDecryptedToken,
  resolveMerchantAccessToken,
  setEncryptedToken,
} from '#/utils/pte-live-token';

/** Vben 壳层兼容样式（列表间距、Drawer/Modal 等） */
import './styles/vben-legacy-compat.css';
import './styles/admin-form-actions.scss';
import './styles/merchant-list-page.scss';
import './styles/native-list-page.scss';
import './styles/native-form-page.scss';
import './styles/native-form-dialog.scss';
import './styles/native-merchant-ui.scss';

async function bootstrap(namespace: string) {
  setDefaultDrawerProps({
    cancelText: '取消',
    confirmText: '保存',
  });

  // 初始化组件适配器
  await initComponentAdapter();

  // 初始化表单组件
  await initSetupVbenForm();

  // // 设置弹窗的默认配置
  // setDefaultModalProps({
  //   fullscreenButton: false,
  // });
  const app = createApp(App);

  // 注册Element Plus提供的v-loading指令
  app.directive('loading', ElLoading.directive);

  // 注册Vben提供的v-loading和v-spinning指令
  registerLoadingDirective(app, {
    loading: false, // Vben提供的v-loading指令和Element Plus提供的v-loading指令二选一即可，此处false表示不注册Vben提供的v-loading指令
    spinning: 'spinning',
  });

  // 国际化 i18n 配置
  await setupI18n(app);

  // 配置 pinia-tore
  await initStores(app, { namespace });

  const accessStore = useAccessStore();
  resolveMerchantAccessToken(
    () => accessStore.accessToken,
    (token) => accessStore.setAccessToken(token),
  );
  if (accessStore.accessToken && !getDecryptedToken()) {
    setEncryptedToken(accessStore.accessToken);
  }
  restoreShopSessionBootstrappedFromStorage();

  // 安装权限指令
  registerAccessDirective(app);
  registerMerchantDirectives(app);
  registerMerchantShell(app);

  // 初始化 tippy
  const { initTippy } = await import('@vben/common-ui/es/tippy');
  initTippy(app);

  // 配置路由及路由守卫
  app.use(router);

  // 配置Motion插件
  const { MotionPlugin } = await import('@vben/plugins/motion');
  app.use(MotionPlugin);

  // 动态更新标题
  watchEffect(() => {
    if (preferences.app.dynamicTitle) {
      const routeTitle = router.currentRoute.value.meta?.title;
      const pageTitle =
        (routeTitle ? `${$t(routeTitle)} - ` : '') + preferences.app.name;
      useTitle(pageTitle);
    }
  });

  app.mount('#app');
}

export { bootstrap };
