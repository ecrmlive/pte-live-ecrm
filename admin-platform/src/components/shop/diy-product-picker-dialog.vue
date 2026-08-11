<script setup lang="ts">
import type { PlatformProduct } from '#/api/core/platform-catalog';

import { ref } from 'vue';

import ProductPickerDialog from '#/components/shop/product-picker-dialog.vue';

const props = withDefaults(
  defineProps<{
    excludeIds?: number[];
    islist?: boolean;
    isproduct?: boolean;
  }>(),
  { excludeIds: () => [], islist: false, isproduct: false },
);

const emit = defineEmits<{
  closeDialog: [payload: { openDialog: boolean; params?: unknown; type: string }];
}>();

const completed = ref(false);

function close(payload: { openDialog: boolean; params?: unknown; type: string }) {
  emit('closeDialog', payload);
}

function select(row: PlatformProduct) {
  if (props.excludeIds.includes(row.product_id)) return;
  completed.value = true;
  close({ openDialog: false, params: row, type: 'success' });
}

function confirm(rows: PlatformProduct[]) {
  const picked = rows.filter((row) => !props.excludeIds.includes(row.product_id));
  completed.value = true;
  close({ openDialog: false, params: picked, type: 'success' });
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
  <ProductPickerDialog
    :open="props.isproduct"
    :multiple="props.islist"
    @update:open="onOpenChange"
    @confirm="confirm"
    @select="select"
  />
</template>
