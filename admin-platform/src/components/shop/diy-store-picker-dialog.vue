<script setup lang="ts">
import type { PickedStore } from '#/components/ecrm/store-picker-modal.vue';

import { ref } from 'vue';

import StorePickerModal from '#/components/ecrm/store-picker-modal.vue';

const props = withDefaults(
  defineProps<{
    excludeIds?: number[];
    islist?: boolean;
    isstore?: boolean;
  }>(),
  { excludeIds: () => [], islist: false, isstore: false },
);

const emit = defineEmits<{
  closeDialog: [payload: { openDialog: boolean; params?: unknown; type: string }];
}>();

const completed = ref(false);

function close(payload: { openDialog: boolean; params?: unknown; type: string }) {
  emit('closeDialog', payload);
}

function confirm(stores: PickedStore[]) {
  const available = stores.filter((store) => !props.excludeIds.includes(store.mer_id));
  completed.value = true;
  close({
    openDialog: false,
    params: props.islist ? available : available[0],
    type: 'success',
  });
}

function onOpenChange(visible: boolean) {
  if (visible) {
    completed.value = false;
    return;
  }
  if (!completed.value) close({ openDialog: false, type: 'error' });
  completed.value = false;
}
</script>

<template>
  <StorePickerModal
    :open="props.isstore"
    @update:open="onOpenChange"
    @confirm="confirm"
  />
</template>
