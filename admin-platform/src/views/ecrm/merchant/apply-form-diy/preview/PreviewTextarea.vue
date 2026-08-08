<script setup lang="ts">
import { computed } from 'vue';

import { previewDisplayValue } from '../registry';
import type { ApplyFormField } from '../types';

const props = defineProps<{ field: ApplyFormField }>();

const value = computed(() => previewDisplayValue(props.field));
</script>

<template>
  <div class="preview-textarea">
    <div class="preview-textarea__title">
      <span v-if="field.required" class="req">*</span>
      {{ field.title || '未命名' }}
    </div>
    <div class="preview-textarea__box">
      <span class="place">{{ value }}</span>
    </div>
  </div>
</template>

<style scoped lang="scss">
.preview-textarea {
  padding: 0 12px 12px;
  border-bottom: 1px solid #eee;
  background: #fff;
}

.preview-textarea__title {
  width: 100%;
  padding-top: 14px;
  padding-bottom: 8px;
  color: #282828;
  font-size: 15px;
}

.preview-textarea__box {
  min-height: 100px;
  padding: 8px 10px;
  border: 1px solid #eee;
  border-radius: 4px;
  background: #fafafa;
  line-height: 1.5;
  font-size: 14px;
}

.place {
  color: #ccc;
  font-weight: 400;
  white-space: pre-wrap;
  word-break: break-word;
}

.req {
  margin-right: 2px;
  color: #f56c6c;
}
</style>
