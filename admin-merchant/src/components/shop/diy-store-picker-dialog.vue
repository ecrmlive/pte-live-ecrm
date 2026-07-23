<script lang="ts" setup>
import { ref, watch } from 'vue';

import StorePickerDialog from '#/components/shop/store-picker-dialog.vue';

defineOptions({ name: 'DiyStorePickerDialog' });

const props = defineProps<{
  excludeIds?: number[];
  islist?: boolean;
  isstore?: boolean;
}>();

const emit = defineEmits<{
  closeDialog: [payload: { params?: unknown; type: string }];
}>();

const open = ref(false);

watch(
  () => props.isstore,
  (value) => {
    open.value = !!value;
  },
  { immediate: true },
);

function onCloseDialog(payload: { params?: unknown; type: string }) {
  emit('closeDialog', payload);
}
</script>

<template>
  <StorePickerDialog
    v-model:open="open"
    :exclude-ids="excludeIds"
    :multiple="!!islist"
    @close-dialog="onCloseDialog"
  />
</template>
