<script setup lang="ts">
import { ElButton } from 'element-plus';
import { ref } from 'vue';

import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';

const model = defineModel<string>({ default: '' });
const props = withDefaults(
  defineProps<{
    disabled?: boolean;
    hint?: string;
    previewSize?: number;
    buttonText?: string;
  }>(),
  {
    disabled: false,
    hint: '',
    previewSize: 100,
    buttonText: '从素材库选择',
  },
);
const open = ref(false);

function select(items: Array<{ file_path: string }>) {
  model.value = items[0]?.file_path ?? '';
}

function openPicker() {
  if (props.disabled) return;
  open.value = true;
}
</script>

<template>
  <div class="image-field">
    <ElButton v-if="!props.disabled" type="primary" @click="openPicker">
      {{ props.buttonText }}
    </ElButton>
    <button
      class="image-field__preview"
      type="button"
      :style="{ width: `${props.previewSize}px`, height: `${props.previewSize}px` }"
      :disabled="props.disabled"
      @click="openPicker"
    >
      <img v-if="model" :src="model" alt="图片预览" />
      <span v-else>暂无图片</span>
    </button>
    <p v-if="props.hint" class="image-field__hint">{{ props.hint }}</p>
    <ImagePickerDialog v-if="!props.disabled" v-model:open="open" @select="select" />
  </div>
</template>

<style scoped lang="scss">
.image-field {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
}
.image-field__preview {
  overflow: hidden;
  padding: 0;
  border: 1px dashed hsl(var(--border));
  border-radius: 8px;
  color: hsl(var(--muted-foreground));
  background: hsl(var(--muted));
  cursor: pointer;
}
.image-field__preview:disabled {
  cursor: default;
}
.image-field__preview img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.image-field__hint {
  margin: 0;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
}
</style>
