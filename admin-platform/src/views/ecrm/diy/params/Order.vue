<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed } from 'vue';

import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import {
  diyBgColors,
  diySection,
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { setDiyStyleType } from './shared/center-style-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';

defineOptions({ name: 'DiyParamsOrder' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

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
  ...diyBgColors('#ffffff', '#f2f2f2'),
  ...DIY_PADDING_FIELDS,
  ...DIY_RADIUS_FIELDS,
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema, {
  fieldPaths: [
    'style.type',
    'style.background',
    'style.bgcolor',
    'style.paddingTop',
    'style.paddingBottom',
    'style.paddingLeft',
    'style.topRadio',
    'style.bottomRadio',
  ],
});

function changeType(type: number) {
  setDiyStyleType(props.curItem, type, true);
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
          :type="Number(curItem.style?.type) === n ? 'primary' : ''"
          @click="changeType(n)"
        >
          风格{{ ['一', '二', '三', '四'][n - 1] }}
        </component>
      </div>
      <div class="mb-4">
        <component
          :is="PrimaryButton"
          v-for="n in [5, 6, 7, 8]"
          :key="n"
          size="small"
          :type="Number(curItem.style?.type) === n ? 'primary' : ''"
          @click="changeType(n)"
        >
          风格{{ ['五', '六', '七', '八'][n - 5] }}
        </component>
      </div>
      <div>
        <component
          :is="PrimaryButton"
          v-for="n in [9, 10]"
          :key="n"
          size="small"
          :type="Number(curItem.style?.type) === n ? 'primary' : ''"
          @click="changeType(n)"
        >
          风格{{ n === 9 ? '九' : '十' }}
        </component>
      </div>
    </div>
    <div class="form-chink" />
    <Form />
  </div>
</template>
