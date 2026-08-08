<script setup lang="ts">
import type { ApplyFormField } from '../types';
import PreviewImage from './PreviewImage.vue';
import PreviewOptions from './PreviewOptions.vue';
import PreviewSelectLike from './PreviewSelectLike.vue';
import PreviewText from './PreviewText.vue';
import PreviewTextarea from './PreviewTextarea.vue';

defineProps<{ field: ApplyFormField }>();
</script>

<template>
  <PreviewImage v-if="field.type === 'image'" :field="field" />
  <PreviewTextarea v-else-if="field.type === 'textarea'" :field="field" />
  <PreviewOptions
    v-else-if="field.type === 'radio' || field.type === 'checkbox'"
    :field="field"
  />
  <PreviewSelectLike
    v-else-if="
      field.type === 'select' ||
      field.type === 'city' ||
      field.type === 'date' ||
      field.type === 'daterange' ||
      field.type === 'time' ||
      field.type === 'timerange'
    "
    :field="field"
  />
  <PreviewText v-else :field="field" />
</template>
