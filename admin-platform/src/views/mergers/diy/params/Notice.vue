<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed } from 'vue';

import {
  diyBgColors,
  diyColor,
  diyInput,
  diySection,
  diySlider,
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { parseIntFields } from './shared/path-utils';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsNotice' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    params?: Record<string, unknown>;
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

const schema = computed((): VbenFormSchema[] => [
  diySection('边距设置'),
  ...diyBgColors('#ffffff', '#f2f2f2'),
  ...DIY_PADDING_FIELDS,
  diySlider('style.padding', '左右内边距：'),
  ...DIY_RADIUS_FIELDS,
  diyColor('style.textColor', '文字颜色：', '#000000'),
  diyInput('params.text', '公告内容：'),
]);

const { Form } = useDiyCurItemForm(
  () => props.curItem,
  schema,
  {
    onInit(item) {
      parseIntFields(item, ['style.paddingTop']);
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
    <div class="form-item ml-[100px] mt-2">
      <div class="form-label mb-2">公告图标：</div>
      <div class="diy-notice-icon">
        <img
          v-img-url="curItem.params?.icon"
          alt=""
          style="width: 100%; height: auto"
          @click="editor.onEditorSelectImage(curItem.params!, 'icon')"
        />
      </div>
      <div class="ww100 gray">建议尺寸32×32</div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.diy-notice-icon,
.diy-notice-icon img {
  width: 32px;
  height: 32px;
}
</style>
