<script lang="ts" setup>
import { computed, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import ShopUserPanel from '../shop-user/shop-user-panel.vue';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  appId: number;
  appName?: string;
}>();

const title = computed(() => {
  const name = props.appName?.trim();
  if (name) {
    return `账号管理 · ${name}（${props.appId}）`;
  }
  return `账号管理 · app_id ${props.appId}`;
});

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(
  open,
  (visible) => {
    if (visible) {
      drawerApi.setState({ title: title.value }).open();
      return;
    }
    drawerApi.close();
  },
  { immediate: true },
);

watch(title, (value) => {
  if (open.value) {
    drawerApi.setState({ title: value });
  }
});
</script>

<template>
  <Drawer
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :footer="false"
    class="shop-manage-drawer"
    content-class="shop-manage-drawer__content"
    :title="title"
  >
    <ShopUserPanel v-if="open && appId > 0" :app-id="appId" :app-name="appName" />
  </Drawer>
</template>
