<script setup lang="ts">
import { computed } from 'vue';

import type { ApplyFormField } from '../types';
import PreviewUploadSlots from './PreviewUploadSlots.vue';

const props = defineProps<{ field: ApplyFormField }>();

const slots = computed(() =>
  Math.max(1, Math.min(props.field.max_upload || 1, 8)),
);
</script>

<template>
  <div class="preview-image">
    <div class="preview-image__title">
      <span v-if="field.required" class="req">*</span>
      {{ field.title || '未命名' }}
    </div>
    <PreviewUploadSlots :count="slots" />
  </div>
</template>

<style scoped lang="scss">
.preview-image {
  padding: 0 12px 15px;
  border-bottom: 1px solid #eee;
  background: #fff;
}

.preview-image__title {
  width: 100%;
  padding-top: 14px;
  color: #282828;
  font-size: 15px;
}

.req {
  margin-right: 2px;
  color: #f56c6c;
}
</style>
