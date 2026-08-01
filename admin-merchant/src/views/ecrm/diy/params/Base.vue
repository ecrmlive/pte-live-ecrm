<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed } from 'vue';

import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';

import {
  diyBgColors,
  diySection,
  diySlider,
  DIY_PADDING_FIELDS,
} from './shared/schema-helpers';
import { setDiyStyleType } from './shared/center-style-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';

defineOptions({ name: 'DiyParamsBase' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    style?: Record<string, unknown>;
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const schema = computed((): VbenFormSchema[] => [
  { component: 'Input', fieldName: 'style.type', hideLabel: true, formItemClass: 'hidden' },
  diySection('样式设置'),
  diySlider('style.padding', '间距：'),
  ...diyBgColors('#ffffff', '#f2f2f2'),
  ...DIY_PADDING_FIELDS,
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema, {
  fieldPaths: [
    'style.type',
    'style.padding',
    'style.background',
    'style.bgcolor',
    'style.paddingTop',
    'style.paddingBottom',
    'style.paddingLeft',
  ],
});

const { PrimaryButton } = useDiyAdapterComponents();

function changeType(type: number) {
  setDiyStyleType(props.curItem, type);
}
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <div class="f16 gray3 form-subtitle">风格设置</div>
    <div class="mb-4 ml-4">
      <div class="mb-4">
        <component
          :is="PrimaryButton"
          v-for="n in [1, 2, 3, 4]"
          :key="n"
          size="small"
          :type="curItem.style?.type == n ? 'primary' : ''"
          @click="changeType(n)"
        >
          {{ ['活力橙', '生鲜绿', '魅力粉', '科技蓝'][n - 1] }}
        </component>
      </div>
      <div>
        <component
          :is="PrimaryButton"
          v-for="n in [5, 6, 7, 8]"
          :key="n"
          size="small"
          :type="curItem.style?.type == n ? 'primary' : ''"
          @click="changeType(n)"
        >
          {{ ['亮黑色', '香槟金', '丁香紫', '热情红'][n - 5] }}
        </component>
      </div>
    </div>
    <div class="form-chink" />
    <Form />
  </div>
</template>
