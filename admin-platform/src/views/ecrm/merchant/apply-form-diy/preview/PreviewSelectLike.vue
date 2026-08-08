<script setup lang="ts">
import { computed } from 'vue';

import { previewDisplayValue } from '../registry';
import type { ApplyFormField } from '../types';
import PreviewRow from './PreviewRow.vue';

const props = defineProps<{ field: ApplyFormField }>();

/** 下拉框预览固定「请选择」，对齐 CRMEB home_select */
const value = computed(() =>
  props.field.type === 'select'
    ? '请选择'
    : previewDisplayValue(props.field),
);
</script>

<template>
  <PreviewRow
    :title="field.title"
    :required="field.required"
    :value="value"
    show-arrow
  />
</template>
