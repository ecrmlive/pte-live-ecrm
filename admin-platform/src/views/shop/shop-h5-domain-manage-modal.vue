<script lang="ts" setup>
import { computed, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import ShopH5DomainPanel from '../shop-h5-domain/shop-h5-domain-panel.vue';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  appId: number;
  appName?: string;
}>();

const title = computed(() => {
  const name = props.appName?.trim();
  if (name) {
    return `域名管理 · ${name}（${props.appId}）`;
  }
  return `域名管理 · app_id ${props.appId}`;
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
    class="shop-h5-domain-manage-modal w-[1020px] max-w-[96vw]"
    content-class="shop-h5-domain-manage-modal__content"
    :title="title"
  >
    <ShopH5DomainPanel v-if="open && appId > 0" :app-id="appId" :app-name="appName" />
  </Modal>
</template>

<style scoped>
:global(.shop-h5-domain-manage-modal) {
  max-height: min(92vh, calc(100dvh - 40px)) !important;
}

:global(.platform-overlay-message-box) {
  z-index: 10000 !important;
}

:global(.shop-h5-domain-manage-modal__content) {
  min-height: 520px;
  max-height: calc(100dvh - 132px);
  overflow-y: auto;
}
</style>
