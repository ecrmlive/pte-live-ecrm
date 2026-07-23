<script lang="ts" setup>
import { computed, defineAsyncComponent } from 'vue';

import { globalShareState } from '@vben/common-ui';

const modelValue = defineModel<string>({ default: '' });

const props = defineProps<{
  defaultColor: string;
  placeholder?: string;
}>();

const ColorPickerComponent = defineAsyncComponent(() =>
  Promise.all([
    import('element-plus/es/components/color-picker/index'),
    import('element-plus/es/components/color-picker/style/css'),
  ]).then(([res]) => res.ElColorPicker),
);

const InputComponent = computed(() => globalShareState.getComponents().Input);
const PrimaryButton = computed(() => globalShareState.getComponents().PrimaryButton);

function resetColor() {
  modelValue.value = props.defaultColor;
}
</script>

<template>
  <div class="flex flex-1 items-center" style="height: 36px">
    <ColorPickerComponent v-model="modelValue" size="default" />
    <component
      :is="InputComponent"
      v-model="modelValue"
      class="ml-2 flex-1"
      :placeholder="placeholder"
    />
    <component :is="PrimaryButton" class="ml-2" link type="primary" @click="resetColor">
      重置
    </component>
  </div>
</template>
