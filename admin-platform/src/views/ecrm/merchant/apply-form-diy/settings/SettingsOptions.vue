<script setup lang="ts">
import { ElFormItem, ElInput, ElSwitch } from 'element-plus';

import type { ApplyFormField } from '../types';

const props = defineProps<{
  field: ApplyFormField;
  disabled?: boolean;
}>();

function optionsText() {
  return (props.field.options || []).join('\n');
}

function updateOptionsText(value: string) {
  props.field.options = value
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean);
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
  <ElFormItem v-if="field.type === 'select'" label="提示语">
    <ElInput
      v-model="field.placeholder"
      maxlength="128"
      placeholder="请输入提示语"
      :disabled="disabled"
    />
  </ElFormItem>
  <ElFormItem label="选项（每行一个）">
    <ElInput
      :model-value="optionsText()"
      type="textarea"
      :rows="4"
      placeholder="选项一&#10;选项二"
      :disabled="disabled"
      @update:model-value="updateOptionsText"
    />
  </ElFormItem>
  <ElFormItem label="是否必填">
    <ElSwitch v-model="field.required" :disabled="disabled" />
  </ElFormItem>
</template>
