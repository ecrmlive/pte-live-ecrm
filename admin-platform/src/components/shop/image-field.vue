<script setup lang="ts">
import { ref } from 'vue';

import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const model = defineModel<string>({ default: '' });
const props = withDefaults(
  defineProps<{
    disabled?: boolean;
    hint?: string;
    previewSize?: number;
    /** 预览块形状：用户头像等场景用 circle */
    shape?: 'square' | 'circle';
    /** @deprecated 已废弃：统一点图开素材库，不再展示旁侧文案按钮 */
    buttonText?: string;
    /** @deprecated 已废弃：默认不展示旁侧按钮 */
    showButton?: boolean;
    defaultLibrary?: 'merchant' | 'system';
  }>(),
  {
    disabled: false,
    hint: '',
    previewSize: 72,
    shape: 'square',
    buttonText: '',
    showButton: false,
    defaultLibrary: 'merchant',
  },
);
const open = ref(false);

function previewSrc() {
  return resolveCosMediaUrl(String(model.value || '').trim());
}

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
    <button
      class="image-field__tile"
      type="button"
      :class="{
        'is-empty': !previewSrc(),
        'is-circle': props.shape === 'circle',
      }"
      :style="{
        width: `${props.previewSize}px`,
        height: `${props.previewSize}px`,
      }"
      :disabled="props.disabled"
      :title="previewSrc() ? '更换图片' : '选择图片'"
      @click="openPicker"
    >
      <img v-if="previewSrc()" :src="previewSrc()" alt="图片预览" />
      <span v-else class="image-field__plus" aria-hidden="true">+</span>
      <span
        v-if="previewSrc() && !props.disabled"
        class="image-field__remove"
        title="删除"
        @click="clearImage"
      >
        ×
      </span>
    </button>
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

.image-field__tile {
  position: relative;
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  padding: 0;
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  color: var(--el-text-color-secondary);
  background: var(--el-fill-color-blank);
  cursor: pointer;
}

.image-field__tile.is-circle {
  border-radius: 50%;
}

.image-field__tile.is-circle .image-field__remove {
  border-radius: 0 50% 0 6px;
}

.image-field__tile.is-empty {
  background: var(--el-fill-color-lighter);
}

.image-field__tile:disabled {
  cursor: default;
}

.image-field__tile img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-field__plus {
  font-size: 22px;
  font-weight: 300;
  line-height: 1;
  color: var(--el-text-color-placeholder);
}

.image-field__remove {
  position: absolute;
  top: 0;
  right: 0;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 0 0 0 4px;
  background: rgb(0 0 0 / 55%);
  color: #fff;
  font-size: 14px;
  line-height: 1;
  opacity: 0;
  transition: opacity 0.15s;
}

.image-field__tile:hover .image-field__remove {
  opacity: 1;
}

.image-field__hint {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.5;
}
</style>
