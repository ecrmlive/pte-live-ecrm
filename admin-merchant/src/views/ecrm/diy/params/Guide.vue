<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed } from 'vue';

import {
  diyColor,
  diyRadioGroup,
  diySection,
  diySlider,
  DIY_PADDING_FIELDS,
} from './shared/schema-helpers';
import { parseIntFields } from './shared/path-utils';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';

defineOptions({ name: 'DiyParamsGuide' });

const props = defineProps<{
  curItem: Record<string, unknown>;
  opts?: unknown;
  selectedIndex?: number;
}>();

const schema = computed((): VbenFormSchema[] => [
  diySection('样式设置'),
  diyColor('style.background', '背景颜色：', '#f2f2f2'),
  diyRadioGroup(
    'style.lineStyle',
    '线条样式：',
    [
      { label: '实线', value: 'solid' },
      { label: '虚线', value: 'dashed' },
      { label: '点状', value: 'dotted' },
    ],
    true,
  ),
  diyColor('style.lineColor', '线条颜色：', '#eeeeee'),
  diySlider('style.lineHeight', '线条高度：'),
  ...DIY_PADDING_FIELDS,
]);

const { Form } = useDiyCurItemForm(
  () => props.curItem,
  schema,
  {
    onInit(item) {
      parseIntFields(item, ['style.lineHeight', 'style.paddingTop']);
    },
  },
);
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <Form />
  </div>
</template>
