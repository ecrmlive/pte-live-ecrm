<script setup lang="ts">
import { ElForm } from 'element-plus';

import { fieldTypeLabel, isDateLikeType, isOptionsType } from '../registry';
import type { ApplyFormField } from '../types';
import SettingsCity from './SettingsCity.vue';
import SettingsDateLike from './SettingsDateLike.vue';
import SettingsImage from './SettingsImage.vue';
import SettingsOptions from './SettingsOptions.vue';
import SettingsText from './SettingsText.vue';
import SettingsTextarea from './SettingsTextarea.vue';

defineProps<{
  field: ApplyFormField | null;
  disabled?: boolean;
}>();
</script>

<template>
  <template v-if="field">
    <div class="diy-right__title">{{ fieldTypeLabel(field.type) }}</div>
    <div class="diy-right__body">
      <ElForm label-position="top" :disabled="disabled" size="default">
        <SettingsText
          v-if="field.type === 'text'"
          :field="field"
          :disabled="disabled"
        />
        <SettingsTextarea
          v-else-if="field.type === 'textarea'"
          :field="field"
          :disabled="disabled"
        />
        <SettingsImage
          v-else-if="field.type === 'image'"
          :field="field"
          :disabled="disabled"
        />
        <SettingsCity
          v-else-if="field.type === 'city'"
          :field="field"
          :disabled="disabled"
        />
        <SettingsDateLike
          v-else-if="isDateLikeType(field.type)"
          :field="field"
          :disabled="disabled"
        />
        <SettingsOptions
          v-else-if="isOptionsType(field.type)"
          :field="field"
          :disabled="disabled"
        />
      </ElForm>
    </div>
  </template>
  <div v-else class="diy-right__empty">请从中间预览选择自定义字段</div>
</template>

<style scoped lang="scss">
.diy-right__title {
  height: 45px;
  padding-left: 24px;
  line-height: 45px;
  border-bottom: 1px solid #eee;
  color: #000;
  font-size: 14px;
}

.diy-right__body {
  padding: 12px 20px 24px;
}

.diy-right__body :deep(.el-form-item) {
  margin-bottom: 14px;
}

.diy-right__empty {
  padding: 48px 24px;
  color: #999;
  text-align: center;
  font-size: 13px;
}
</style>
