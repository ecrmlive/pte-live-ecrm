<script lang="ts" setup>
import { computed, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';

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

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(
  open,
  (visible) => {
    if (visible) {
      modalApi.open();
      return;
    }
    modalApi.close();
  },
  { immediate: true },
);
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :footer="false"
    class="shop-manage-modal w-[980px] max-w-[96vw]"
    content-class="shop-manage-modal__content"
    :title="title"
  >
    <ShopUserPanel v-if="open && appId > 0" :app-id="appId" :app-name="appName" />
  </Modal>
</template>
