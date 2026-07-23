<script setup lang="ts">
import { computed, ref } from 'vue';

import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';
import ImagePickerTrigger from '#/components/shop/image-picker-trigger.vue';

const model = defineModel<string>({ default: '' });

const props = withDefaults(
  defineProps<{
    buttonText?: string;
    /** 打开选择弹窗时的默认素材库 */
    defaultLibrary?: 'merchant' | 'system';
    hint?: string;
    /** @deprecated use buttonText */
    label?: string;
    previewSize?: number;
  }>(),
  {
    buttonText: '',
    defaultLibrary: 'merchant',
    hint: '',
    label: '',
    previewSize: 100,
  },
);

const open = ref(false);

const resolvedButtonText = computed(
  () => props.buttonText || props.label || '选择图片',
);

function onSelect(items: Array<{ file_path: string }>) {
  model.value = items[0]?.file_path ?? '';
}
</script>

<template>
  <ImagePickerTrigger
    :button-text="resolvedButtonText"
    :hint="hint"
    :preview-url="model"
    :size="previewSize"
    @pick="open = true"
  />
  <ImagePickerDialog
    v-model:open="open"
    :default-library="defaultLibrary"
    @select="onSelect"
  />
</template>
