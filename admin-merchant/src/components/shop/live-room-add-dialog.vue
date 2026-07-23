<script setup lang="ts">
import { useVbenModal } from '@vben/common-ui';

import { watch } from 'vue';

import LiveRoomAddModal from '#/views/native/live/room/live-room-add-modal.vue';

defineOptions({ name: 'LiveRoomAddDialog' });

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{ success: [] }>();

const [AddModal, addModalApi] = useVbenModal({
  connectedComponent: LiveRoomAddModal,
  onOpenChange(isOpen) {
    if (!isOpen) {
      open.value = false;
    }
  },
});

watch(open, (visible) => {
  if (visible) {
    addModalApi.open();
  } else {
    addModalApi.close();
  }
});

function onSuccess() {
  open.value = false;
  emit('success');
}
</script>

<template>
  <AddModal @success="onSuccess" />
</template>
