<script setup lang="ts">
import { computed, watch } from 'vue';

import {
  ElDatePicker,
  ElFormItem,
  ElInput,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElSwitch,
  ElTimePicker,
} from 'element-plus';

import type { ApplyFormField } from '../types';

const RANGE_SEP = ' ~ ';

const props = defineProps<{
  field: ApplyFormField;
  disabled?: boolean;
}>();

watch(
  () => props.field,
  (field) => {
    if (!field.default_visible) field.default_visible = 'show';
    if (!field.default_mode) field.default_mode = 'current';
  },
  { immediate: true },
);

const isDateFamily = computed(
  () => props.field.type === 'date' || props.field.type === 'daterange',
);

const currentLabel = computed(() =>
  isDateFamily.value ? '当前日期' : '当前时间',
);
const specifyLabel = computed(() =>
  isDateFamily.value ? '指定日期' : '指定时间',
);

const showSpecify = computed(
  () =>
    props.field.default_visible === 'show' &&
    props.field.default_mode === 'specify',
);

/** daterange / timerange：picker 用 [start, end]，落库仍为 "a ~ b" 字符串 */
const specifyRangeValue = computed({
  get(): [string, string] | undefined {
    const raw = (props.field.specify_value || '').trim();
    if (!raw) return undefined;
    const parts = raw.split(RANGE_SEP).map((s) => s.trim());
    if (parts.length === 2 && parts[0] && parts[1]) {
      return [parts[0], parts[1]];
    }
    return undefined;
  },
  set(v: [string, string] | null | undefined) {
    if (!v?.[0] || !v[1]) {
      props.field.specify_value = '';
      return;
    }
    props.field.specify_value = `${v[0]}${RANGE_SEP}${v[1]}`;
  },
});
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
    <div class="default-box">
      <ElRadioGroup v-model="field.default_visible" :disabled="disabled">
        <ElRadio label="show">显示</ElRadio>
        <ElRadio label="hide">隐藏</ElRadio>
      </ElRadioGroup>
      <ElSelect
        v-if="field.default_visible === 'show'"
        v-model="field.default_mode"
        class="w-full"
        :disabled="disabled"
      >
        <ElOption :label="currentLabel" value="current" />
        <ElOption :label="specifyLabel" value="specify" />
      </ElSelect>
      <ElDatePicker
        v-if="showSpecify && field.type === 'date'"
        v-model="field.specify_value"
        class="w-full"
        type="date"
        format="YYYY-MM-DD"
        value-format="YYYY-MM-DD"
        placeholder="请选择日期"
        clearable
        :disabled="disabled"
      />
      <ElDatePicker
        v-else-if="showSpecify && field.type === 'daterange'"
        v-model="specifyRangeValue"
        class="w-full"
        type="daterange"
        format="YYYY-MM-DD"
        value-format="YYYY-MM-DD"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        clearable
        :disabled="disabled"
      />
      <ElTimePicker
        v-else-if="showSpecify && field.type === 'time'"
        v-model="field.specify_value"
        class="w-full"
        format="HH:mm"
        value-format="HH:mm"
        placeholder="请选择时间"
        clearable
        :disabled="disabled"
      />
      <ElTimePicker
        v-else-if="showSpecify && field.type === 'timerange'"
        v-model="specifyRangeValue"
        class="w-full"
        is-range
        format="HH:mm"
        value-format="HH:mm"
        start-placeholder="开始时间"
        end-placeholder="结束时间"
        range-separator="~"
        clearable
        :disabled="disabled"
      />
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

<style scoped>
.w-full {
  width: 100%;
}

.default-box {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}

.default-box :deep(.el-date-editor) {
  width: 100%;
}
</style>
