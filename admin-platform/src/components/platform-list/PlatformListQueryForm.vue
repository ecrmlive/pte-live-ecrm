<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';

import { useVbenForm } from '#/adapter/form';

const props = defineProps<{
  options: VbenFormProps;
}>();

const emit = defineEmits<{
  reset: [values: Record<string, unknown>];
  submit: [values: Record<string, unknown>];
}>();

const [QueryForm] = useVbenForm({
  ...props.options,
  handleReset: async () => {
    emit('reset', {});
  },
  handleSubmit: async (values) => {
    emit('submit', values);
  },
});
</script>

<template>
  <QueryForm />
</template>
