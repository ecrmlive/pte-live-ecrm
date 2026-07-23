<script lang="ts" setup>
import { ref, watch } from 'vue';

import ProductMultiPickerDialog from '#/components/shop/product-multi-picker-dialog.vue';

defineOptions({ name: 'DiyProductPickerDialog' });

const props = defineProps<{
  excludeIds?: number[];
  islist?: boolean;
  isproduct?: boolean;
}>();

const emit = defineEmits<{
  closeDialog: [payload: { openDialog: boolean; params?: unknown; type: string }];
}>();

const open = ref(false);

watch(
  () => props.isproduct,
  (value) => {
    open.value = !!value;
  },
  { immediate: true },
);

function onCloseDialog(payload: { openDialog: boolean; params?: unknown; type: string }) {
  emit('closeDialog', payload);
}
</script>

<template>
    <ProductMultiPickerDialog
      v-model:open="open"
      :exclude-ids="excludeIds"
      @close-dialog="onCloseDialog"
    />
</template>
