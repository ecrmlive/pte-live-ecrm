<script setup lang="ts">
import { computed, ref, watch } from 'vue';

import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

type PickerItem = {
  attachment_id: number;
  file_path: string;
};

const model = defineModel<string[]>({ default: () => [] });
const attachmentIds = defineModel<number[]>('attachmentIds', {
  default: () => [],
});

const props = withDefaults(
  defineProps<{
    disabled?: boolean;
    hint?: string;
    previewSize?: number;
    limit?: number;
    defaultLibrary?: 'merchant' | 'system';
  }>(),
  {
    disabled: false,
    hint: '',
    previewSize: 72,
    limit: 9,
    defaultLibrary: 'merchant',
  },
);

const open = ref(false);
const replaceIndex = ref<number | null>(null);
/** 与 urls 等长；0 表示未知 id（仅有 url 场景） */
const internalIds = ref<number[]>([]);

const urls = computed(() =>
  (model.value ?? []).map((url) => String(url || '').trim()).filter(Boolean),
);

const remaining = computed(() =>
  Math.max(0, props.limit - urls.value.length),
);

const pickerLimit = computed(() => {
  if (replaceIndex.value !== null) return 1;
  return Math.max(1, remaining.value);
});

const showAdd = computed(
  () => !props.disabled && urls.value.length < props.limit,
);

function previewSrc(url: string) {
  return resolveCosMediaUrl(url);
}

function syncState(nextUrls: string[], nextIds: number[]) {
  const alignedIds = nextUrls.map((_, index) => Number(nextIds[index] || 0));
  model.value = nextUrls;
  internalIds.value = alignedIds;
  attachmentIds.value = alignedIds.filter((id) => id > 0);
}

watch(
  urls,
  (next) => {
    if (internalIds.value.length === next.length) return;
    const prev = internalIds.value;
    internalIds.value = next.map((_, index) => Number(prev[index] || 0));
    attachmentIds.value = internalIds.value.filter((id) => id > 0);
  },
  { immediate: true },
);

function openAdd() {
  if (props.disabled || remaining.value <= 0) return;
  replaceIndex.value = null;
  open.value = true;
}

function openReplace(index: number) {
  if (props.disabled) return;
  replaceIndex.value = index;
  open.value = true;
}

function removeAt(index: number, event: Event) {
  event.stopPropagation();
  if (props.disabled) return;
  const nextUrls = [...urls.value];
  const nextIds = [...internalIds.value];
  while (nextIds.length < nextUrls.length) nextIds.push(0);
  nextUrls.splice(index, 1);
  nextIds.splice(index, 1);
  syncState(nextUrls, nextIds);
}

function onSelect(items: PickerItem[]) {
  if (!items.length) return;
  const nextUrls = [...urls.value];
  const nextIds = [...internalIds.value];
  while (nextIds.length < nextUrls.length) nextIds.push(0);

  if (replaceIndex.value !== null) {
    const idx = replaceIndex.value;
    const item = items[0];
    if (item && idx >= 0 && idx < nextUrls.length) {
      nextUrls[idx] = item.file_path;
      nextIds[idx] = item.attachment_id;
    }
  } else {
    for (const item of items) {
      if (nextUrls.length >= props.limit) break;
      nextUrls.push(item.file_path);
      nextIds.push(item.attachment_id);
    }
  }

  syncState(nextUrls, nextIds);
  replaceIndex.value = null;
}
</script>

<template>
  <div class="images-field">
    <div class="images-field__grid">
      <button
        v-for="(url, index) in urls"
        :key="`${url}-${index}`"
        class="images-field__tile"
        type="button"
        :style="{
          width: `${props.previewSize}px`,
          height: `${props.previewSize}px`,
        }"
        :disabled="props.disabled"
        title="更换图片"
        @click="openReplace(index)"
      >
        <img :src="previewSrc(url)" alt="图片预览" />
        <span
          v-if="!props.disabled"
          class="images-field__remove"
          title="删除"
          @click="removeAt(index, $event)"
        >
          ×
        </span>
      </button>
      <button
        v-if="showAdd"
        class="images-field__tile is-empty"
        type="button"
        :style="{
          width: `${props.previewSize}px`,
          height: `${props.previewSize}px`,
        }"
        title="选择图片"
        @click="openAdd"
      >
        <span class="images-field__plus" aria-hidden="true">+</span>
      </button>
    </div>
    <p v-if="props.hint" class="images-field__hint">{{ props.hint }}</p>
    <ImagePickerDialog
      v-if="!props.disabled"
      v-model:open="open"
      :default-library="props.defaultLibrary"
      kind="image"
      :limit="pickerLimit"
      @select="onSelect"
    />
  </div>
</template>

<style scoped lang="scss">
.images-field {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
}

.images-field__grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.images-field__tile {
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

.images-field__tile.is-empty {
  background: var(--el-fill-color-lighter);
}

.images-field__tile:disabled {
  cursor: default;
}

.images-field__tile img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.images-field__plus {
  font-size: 22px;
  font-weight: 300;
  line-height: 1;
  color: var(--el-text-color-placeholder);
}

.images-field__remove {
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

.images-field__tile:hover .images-field__remove {
  opacity: 1;
}

.images-field__hint {
  margin: 0;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.5;
}
</style>
