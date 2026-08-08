<script setup lang="ts">
import { computed } from 'vue';

import { previewDisplayValue } from '../registry';
import type { ApplyFormField } from '../types';
import PreviewRow from './PreviewRow.vue';

const props = defineProps<{ field: ApplyFormField }>();

const value = computed(() => previewDisplayValue(props.field));
</script>

<template>
  <PreviewRow
    :title="field.title"
    :required="field.required"
    :value="value"
  />
</template>
