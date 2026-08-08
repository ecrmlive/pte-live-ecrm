<script setup lang="ts">
import { ElFormItem, ElInput, ElSwitch } from 'element-plus';

import { CITY_LEVEL_OPTIONS } from '../registry';
import type { ApplyCityLevel, ApplyFormField } from '../types';

const props = defineProps<{
  field: ApplyFormField;
  disabled?: boolean;
}>();

function setCityLevel(value: ApplyCityLevel) {
  if (props.disabled) return;
  if (props.field.city_level === value) return;
  props.field.city_level = value;
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
  <ElFormItem label="默认值">
    <!-- 同款视觉：device-preview__switcher / __switcher-btn（灰底轨道 + pill） -->
    <div
      class="city-level-switcher"
      role="radiogroup"
      aria-label="城市默认值级别"
    >
      <button
        v-for="item in CITY_LEVEL_OPTIONS"
        :key="item.value"
        type="button"
        class="city-level-switcher__btn"
        role="radio"
        :aria-checked="field.city_level === item.value"
        :class="{ 'is-active': field.city_level === item.value }"
        :disabled="disabled"
        @click="setCityLevel(item.value)"
      >
        {{ item.label }}
      </button>
    </div>
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

<style scoped lang="scss">
/* 对齐 @pte-live/diy .device-preview__switcher / __switcher-btn */
.city-level-switcher {
  box-sizing: border-box;
  display: flex;
  flex-shrink: 0;
  align-items: center;
  width: 100%;
  height: 32px;
  padding: 2px;
  border: 1px solid hsl(var(--border, 220 13% 91%));
  border-radius: 8px;
  background: hsl(var(--muted, 210 40% 96%));
  box-shadow: 0 1px 2px rgb(15 23 42 / 6%);
}

.city-level-switcher__btn {
  box-sizing: border-box;
  display: inline-flex;
  flex: 1 1 0;
  align-items: center;
  justify-content: center;
  min-width: 0;
  height: 28px;
  margin: 0;
  padding: 0 8px;
  border: 1px solid transparent;
  border-radius: 6px;
  background: transparent;
  color: hsl(var(--muted-foreground, 215 16% 47%));
  font-size: 13px;
  font-weight: 500;
  line-height: 1;
  white-space: nowrap;
  vertical-align: middle;
  appearance: none;
  -webkit-appearance: none;
  transform: none;
  cursor: pointer;
  transition:
    background 0.15s ease,
    color 0.15s ease;
}

.city-level-switcher__btn:hover:not(.is-active):not(:disabled) {
  color: hsl(var(--foreground, 222 47% 11%));
  background: hsl(var(--background, 0 0% 100%) / 70%);
}

.city-level-switcher__btn.is-active {
  margin: 0;
  border-color: transparent;
  color: hsl(var(--primary-foreground, 0 0% 100%));
  background: hsl(var(--primary, 221 83% 53%));
  transform: none;
}

.city-level-switcher__btn:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}
</style>
