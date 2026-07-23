<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed } from 'vue';

import {
  diyBgColors,
  diyInput,
  diySection,
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsWxlive' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    style?: Record<string, unknown>;
  };
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

const schema = computed((): VbenFormSchema[] => [
  diySection('边距设置'),
  ...diyBgColors('#ffffff', '#f2f2f2'),
  ...DIY_PADDING_FIELDS,
  ...DIY_RADIUS_FIELDS,
  diySection('组件设置'),
  diyInput('params.showNum', '显示数量：'),
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema, {
  fieldPaths: ['params.showNum', 'style.background_image'],
});
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <Form />
    <div class="form-chink" />
    <div class="f16 gray3 form-subtitle px-4">组件样式</div>
    <div class="form-item ml-[100px]">
      <div class="form-label mb-2">背景图片：</div>
      <div class="diy-notice-icon">
        <img
          v-img-url="curItem.style?.background_image"
          alt=""
          @click="editor.onEditorSelectImage(curItem.style!, 'background_image')"
        />
      </div>
      <div class="gray">建议尺寸700px*90px</div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.diy-notice-icon,
.diy-notice-icon img {
  width: 350px;
  height: 45px;
}
</style>
