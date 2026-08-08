<script setup lang="ts">
import { ElFormItem, ElInput, ElInputNumber, ElSwitch } from 'element-plus';

import type { ApplyFormField } from '../types';

const props = defineProps<{
  field: ApplyFormField;
  disabled?: boolean;
}>();

function onMaxUploadChange(val: number | undefined) {
  props.field.max_upload = Math.max(1, Math.min(Number(val) || 1, 20));
}
</script>

<template>
  <ElFormItem label="标题">
    <ElInput
      v-model="field.title"
      maxlength="64"
      placeholder="请输入标题"
      :disabled="disabled"
    />
  </ElFormItem>
  <ElFormItem label="最多上传">
    <ElInputNumber
      :model-value="field.max_upload || 1"
      :min="1"
      :max="20"
      :disabled="disabled"
      controls-position="right"
      @update:model-value="onMaxUploadChange"
    />
  </ElFormItem>
  <ElFormItem label="是否必填">
    <ElSwitch v-model="field.required" :disabled="disabled" />
  </ElFormItem>
</template>
