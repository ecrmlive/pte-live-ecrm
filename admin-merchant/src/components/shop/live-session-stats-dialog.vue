<script setup lang="ts">
import type { LiveStreamLogItem } from '#/api/core/live';

import { useVbenModal } from '@vben/common-ui';
import { computed, watch } from 'vue';

import LiveSessionStatsPanel from '#/views/native/live/session/components/live-session-stats-panel.vue';

defineOptions({ name: 'LiveSessionStatsDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  row?: LiveStreamLogItem | Record<string, never>;
}>();

const statsLogId = computed(() => {
  const id = Number.parseInt(String(props.row?.log_id ?? ''), 10);
  return Number.isFinite(id) && id > 0 ? id : 0;
});

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) modalApi.open();
  else modalApi.close();
});
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="session-stats-dialog w-[1280px]"
    content-class="session-stats-dialog__content"
    title="场次数据统计"
  >
    <LiveSessionStatsPanel v-if="open && statsLogId > 0" :log-id="statsLogId" />
  </Modal>
</template>

<style lang="scss">
.session-stats-dialog__content {
  padding-top: 8px;
  max-height: calc(92vh - 100px);
  overflow-y: auto;
}
</style>
