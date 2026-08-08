<script setup lang="ts">
import {
  ElFormItem,
  ElInput,
  ElOption,
  ElSelect,
  ElSwitch,
} from 'element-plus';

import { CONTENT_TYPE_OPTIONS } from '../registry';
import type { ApplyFormField } from '../types';

defineProps<{
  field: ApplyFormField;
  disabled?: boolean;
}>();
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
  <ElFormItem label="内容">
    <ElSelect v-model="field.content_type" class="w-full" :disabled="disabled">
      <ElOption
        v-for="item in CONTENT_TYPE_OPTIONS"
        :key="item.value"
        :label="item.label"
        :value="item.value"
      />
    </ElSelect>
  </ElFormItem>
  <ElFormItem label="默认值">
    <ElInput
      v-model="field.default_value"
      maxlength="256"
      placeholder="请输入默认值"
      :disabled="disabled"
    />
  </ElFormItem>
  <ElFormItem label="提示语">
    <ElInput
      v-model="field.placeholder"
      maxlength="128"
      placeholder="请输入提示语"
      :disabled="disabled"
    />
  </ElFormItem>
  <ElFormItem label="是否必填">
    <ElSwitch v-model="field.required" :disabled="disabled" />
  </ElFormItem>
</template>

<style scoped>
.w-full {
  width: 100%;
}
</style>
