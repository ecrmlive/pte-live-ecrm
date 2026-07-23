<script lang="ts" setup>
import { defineAsyncComponent } from 'vue';

const modelValue = defineModel<number>({ default: 0 });

defineProps<{
  max?: number;
  min?: number;
}>();

const SliderComponent = defineAsyncComponent(() =>
  Promise.all([
    import('element-plus/es/components/slider/index'),
    import('element-plus/es/components/slider/style/css'),
  ]).then(([res]) => res.ElSlider),
);
</script>

<template>
  <SliderComponent
    v-model="modelValue"
    class="diy-slider-field"
    input-size="small"
    :max="max"
    :min="min"
    :show-input-controls="false"
    show-input
    size="small"
  />
</template>
