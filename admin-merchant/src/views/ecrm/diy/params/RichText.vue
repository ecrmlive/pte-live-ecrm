<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, ref, watch } from 'vue';

import RichTextField from '#/components/shop/rich-text-field.vue';

import {
  diyColor,
  diySection,
  diySlider,
} from './shared/schema-helpers';
import { parseIntFields } from './shared/path-utils';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';

defineOptions({ name: 'DiyParamsRichText' });

const props = defineProps<{
  curItem: {
    name?: string;
    params: { content?: string };
    style: {
      background?: string;
      paddingLeft?: number | string;
      paddingTop?: number | string;
    };
  };
  selectedIndex?: number;
}>();

const content = ref('');

const schema = computed((): VbenFormSchema[] => [
  diySection('样式设置'),
  diySlider('style.paddingTop', '上下边距：'),
  diySlider('style.paddingLeft', '左右边距：'),
  diyColor('style.background', '背景颜色：', '#ffffff'),
]);

const { Form } = useDiyCurItemForm(
  () => props.curItem as Record<string, unknown>,
  schema,
  {
    onInit(item) {
      parseIntFields(item, ['style.paddingTop', 'style.paddingLeft']);
    },
  },
);

function syncContent() {
  content.value = props.curItem.params.content ?? '';
}

watch(
  () => props.curItem.params.content,
  (value) => {
    content.value = value ?? '';
  },
);

watch(content, (value) => {
  props.curItem.params.content = value;
});

onMounted(() => {
  syncContent();
});
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <Form />
    <div class="mx-auto w-[402px]">
      <RichTextField v-model="content" />
    </div>
  </div>
</template>
