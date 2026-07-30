<script setup lang="ts">
import { ElButton } from 'element-plus';
import { ref } from 'vue';

import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';

const model = defineModel<string>({ default: '' });
const props = withDefaults(defineProps<{ hint?: string; previewSize?: number }>(), { hint: '', previewSize: 100 });
const open = ref(false);

function select(items: Array<{ file_path: string }>) {
  model.value = items[0]?.file_path ?? '';
}
</script>

<template>
  <div class="image-field">
    <ElButton type="primary" @click="open = true">选择图片</ElButton>
    <button class="image-field__preview" type="button" :style="{ width: `${props.previewSize}px`, height: `${props.previewSize}px` }" @click="open = true">
      <img v-if="model" :src="model" alt="图片预览" />
      <span v-else>暂无图片</span>
    </button>
    <p v-if="props.hint" class="image-field__hint">{{ props.hint }}</p>
    <ImagePickerDialog v-model:open="open" @select="select" />
  </div>
</template>

<style scoped lang="scss">
.image-field { display: flex; flex-direction: column; align-items: flex-start; gap: 10px; }.image-field__preview { overflow: hidden; padding: 0; border: 1px dashed hsl(var(--border)); border-radius: 8px; color: hsl(var(--muted-foreground)); background: hsl(var(--muted)); }.image-field__preview img { display: block; width: 100%; height: 100%; object-fit: cover; }.image-field__hint { margin: 0; color: hsl(var(--muted-foreground)); font-size: 12px; }
</style>
