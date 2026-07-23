<script lang="ts" setup>
import type { ShopLinkValue } from '#/types/shop-link';

import { ref, watch } from 'vue';

import LinkPickerDialog from '#/components/shop/link-picker-dialog.vue';

defineOptions({ name: 'DiyLinkPickerDialog' });

const props = defineProps<{
  is_linkset?: boolean;
  linkData?: ShopLinkValue | null;
}>();

const emit = defineEmits<{
  closeDialog: [payload?: ShopLinkValue];
}>();

const open = ref(false);

watch(
  () => props.is_linkset,
  (value) => {
    open.value = !!value;
  },
  { immediate: true },
);

function onConfirm(payload?: ShopLinkValue) {
  open.value = false;
  emit('closeDialog', payload);
}
</script>

<template>
  <LinkPickerDialog v-model:open="open" :link-data="linkData" @confirm="onConfirm" />
</template>
