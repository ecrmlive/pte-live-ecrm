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
    /** 是否展示「从素材库选择」文字按钮；缩略图本身仍可点击打开 */
    showButton?: boolean;
    defaultLibrary?: 'merchant' | 'system';
  }>(),
  {
    disabled: false,
    hint: '',
    previewSize: 60,
    buttonText: '从素材库选择',
    showButton: true,
    defaultLibrary: 'merchant',
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

function clearImage(event: Event) {
  event.stopPropagation();
  if (props.disabled) return;
  model.value = '';
}
</script>

<template>
  <div class="image-field">
    <div class="image-field__row">
      <button
        class="image-field__preview"
        type="button"
        :class="{ 'is-empty': !model }"
        :style="{
          width: `${props.previewSize}px`,
          height: `${props.previewSize}px`,
        }"
        :disabled="props.disabled"
        :title="model ? '更换图片' : '选择图片'"
        @click="openPicker"
      >
        <img v-if="model" :src="model" alt="图片预览" />
        <span v-else class="image-field__plus" aria-hidden="true">+</span>
        <span
          v-if="model && !props.disabled"
          class="image-field__clear"
          @click="clearImage"
        >
          清除
        </span>
      </button>
      <ElButton
        v-if="props.showButton && !props.disabled"
        type="primary"
        link
        @click="openPicker"
      >
        {{ props.buttonText }}
      </ElButton>
    </div>
    <p v-if="props.hint" class="image-field__hint">{{ props.hint }}</p>
    <ImagePickerDialog
      v-if="!props.disabled"
      v-model:open="open"
      :default-library="props.defaultLibrary"
      @select="select"
    />
  </div>
</template>

<style scoped lang="scss">
.image-field {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
}

.image-field__row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.image-field__preview {
  position: relative;
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  padding: 0;
  border: 1px dashed hsl(var(--border));
  border-radius: 4px;
  color: hsl(var(--muted-foreground));
  background: #fff;
  cursor: pointer;
}

.image-field__preview.is-empty {
  background: hsl(var(--background));
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

.image-field__plus {
  font-size: 22px;
  font-weight: 300;
  line-height: 1;
  color: #c0c4cc;
}

.image-field__clear {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgb(0 0 0 / 45%);
  color: #fff;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}

.image-field__preview:hover .image-field__clear {
  opacity: 1;
}

.image-field__hint {
  margin: 0;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 1.5;
}
</style>
