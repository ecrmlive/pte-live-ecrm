<script setup lang="ts">
import type { ImageUploadOptions } from '@vben/plugins/tiptap';

import { VbenTiptap } from '@vben/plugins/tiptap';
import { ElMessage } from 'element-plus';

import { uploadShopImageDirect } from '#/api/core/file';

const model = defineModel<string>({ default: '' });

withDefaults(
  defineProps<{
    maxHeight?: number | string;
    minHeight?: number;
  }>(),
  { maxHeight: 400, minHeight: 420 },
);

const tiptapImageUpload: ImageUploadOptions = {
  maxSize: 5 * 1024 * 1024,
  upload: async (file) => {
    const uploaded = await uploadShopImageDirect(file);
    return uploaded.file_path;
  },
  onUploadError: (error) => {
    ElMessage.error(error instanceof Error ? error.message : '图片上传失败');
  },
};
</script>

<template>
  <div class="rich-text-field w-full min-w-0">
    <VbenTiptap
      v-model="model"
      :image-upload="tiptapImageUpload"
      :max-height="maxHeight"
      :min-height="minHeight"
    />
  </div>
</template>

<style scoped></style>
