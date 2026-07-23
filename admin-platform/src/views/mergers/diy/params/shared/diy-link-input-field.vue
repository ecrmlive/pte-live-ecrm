<script lang="ts" setup>
import type { InputProps } from 'element-plus';

import { computed } from 'vue';

import { globalShareState } from '@vben/common-ui';

const props = withDefaults(
  defineProps<{
    disabled?: boolean;
    maxlength?: number;
    readonly?: boolean;
    showWordLimit?: boolean;
    size?: InputProps['size'];
  }>(),
  {
    disabled: false,
    maxlength: undefined,
    readonly: false,
    showWordLimit: false,
    size: undefined,
  },
);

const modelValue = defineModel<string>({ default: '' });

const InputComponent = computed(() => globalShareState.getComponents().Input);
</script>

<template>
  <component
    :is="InputComponent"
    v-model="modelValue"
    class="diy-link-input-field min-w-0 flex-1"
    :disabled="props.disabled"
    :maxlength="props.maxlength"
    :readonly="props.readonly"
    :show-word-limit="props.showWordLimit"
    :size="props.size"
  >
    <template v-if="$slots.suffix" #suffix>
      <slot name="suffix"></slot>
    </template>
  </component>
</template>
