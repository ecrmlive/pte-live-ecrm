<script setup lang="ts">
import type { ShopLinkValue } from '#/types/shop-link';

import LinkPickerDialog from '#/components/shop/link-picker-dialog.vue';

defineOptions({ name: 'ShopTabLinkPickerDialog' });

const open = defineModel<boolean>('open', { default: false });

defineProps<{
  linkData?: ShopLinkValue | null;
}>();

const emit = defineEmits<{
  confirm: [ShopLinkValue | undefined];
}>();
</script>

<template>
  <LinkPickerDialog v-model:open="open" :link-data="linkData" @confirm="emit('confirm', $event)" />
</template>
