<script lang="ts" setup>
import { computed, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';

import TrafficPanel from './traffic-panel.vue';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  appId: number;
  appName?: string;
}>();

const title = computed(() => {
  const name = props.appName?.trim();
  if (name) {
    return `流量管理 · ${name}（${props.appId}）`;
  }
  return `流量管理 · app_id ${props.appId}`;
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
    class="shop-manage-modal w-[1080px] max-w-[96vw]"
    content-class="shop-manage-modal__content"
    :title="title"
  >
    <TrafficPanel v-if="open && appId > 0" :app-id="appId" embedded />
  </Modal>
</template>
