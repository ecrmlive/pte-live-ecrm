import { createApp, watchEffect } from 'vue';

import { registerAccessDirective } from '@vben/access';
import { registerLoadingDirective } from '@vben/common-ui';
import { preferences } from '@vben/preferences';
import { initStores } from '@vben/stores';
import '@vben/styles';
import '@vben/styles/ele';

import { useTitle } from '@vueuse/core';
import ElementPlus from 'element-plus';
import { $t, setupI18n } from '#/locales';

import { initComponentAdapter } from './adapter/component';
import { initSetupVbenForm } from './adapter/form';
import '#/adapter/vxe-table';
import App from './app.vue';
import { router } from './router';
import { hydratePlatformUserStore } from '#/store/platform-user';

async function bootstrap(namespace: string) {
  await initComponentAdapter();
  await initSetupVbenForm();

  const app = createApp(App);

  // 业务页使用 Element Plus 原生组件。Vben 的表单适配器只服务于
  // schema 表单，不会全局注册 el-table、el-card、el-form 等模板组件。
  // 未注册时 Vue 会把这些组件当成未知标签，最终导致页面渲染中断。
  app.use(ElementPlus);

  registerLoadingDirective(app, {
    loading: false,
    spinning: 'spinning',
  });

  await setupI18n(app);
  await initStores(app, { namespace });

  registerAccessDirective(app);

  hydratePlatformUserStore();

  const { initTippy } = await import('@vben/common-ui/es/tippy');
  initTippy(app);

  app.use(router);

  const { MotionPlugin } = await import('@vben/plugins/motion');
  app.use(MotionPlugin);

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
