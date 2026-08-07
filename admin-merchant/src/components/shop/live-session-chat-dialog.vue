<script setup lang="ts">
import type { LiveStreamLogItem } from '#/api/core/live';

import { useVbenDrawer } from '@vben/common-ui';
import { computed, watch } from 'vue';

import LiveSessionChatPanel from '#/components/shop/live-session-chat-panel.vue';

defineOptions({ name: 'LiveSessionChatDialog' });

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  row?: LiveStreamLogItem | Record<string, never>;
}>();

const liveId = computed(() => {
  const id = Number.parseInt(String(props.row?.live_id ?? ''), 10);
  return Number.isFinite(id) && id > 0 ? id : 0;
});

const sessionId = computed(() => String(props.row?.session_id ?? '').trim());

const canLoad = computed(() => liveId.value > 0 && sessionId.value !== '');

const headerText = computed(() => {
  if (!canLoad.value) return '';
  const name = String(props.row?.name ?? '').trim() || `直播间 ${liveId.value}`;
  const start = props.row?.start_time_text || '-';
  const end = props.row?.end_time_text || '-';
  return `${name} · ${start} ~ ${end}`;
});

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
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
    class="session-chat-dialog w-[720px]"
    content-class="session-chat-dialog__content"
    title="弹幕消息"
  >
    <div v-if="headerText" class="session-chat-dialog__meta">{{ headerText }}</div>
    <LiveSessionChatPanel
      v-if="open && canLoad"
      embedded
      :live-id="liveId"
      :session-id="sessionId"
    />
    <el-empty v-else-if="open" :image-size="80" description="该记录缺少场次信息，无法加载弹幕" />
  </Modal>
</template>

<style scoped lang="scss">
.session-chat-dialog__meta {
  margin: -4px 0 12px;
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
}
</style>

<style lang="scss">
.session-chat-dialog__content {
  display: flex;
  flex-direction: column;
  padding-top: 8px;
  max-height: calc(88vh - 100px);
  overflow: hidden;
}
</style>
