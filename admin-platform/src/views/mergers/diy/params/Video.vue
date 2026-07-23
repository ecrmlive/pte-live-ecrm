<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed } from 'vue';

import {
  diyColor,
  diyInput,
  diyRadioGroup,
  diySection,
  diySlider,
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { parseIntFields } from './shared/path-utils';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsVideo' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    params?: Record<string, unknown>;
  };
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

const schema = computed((): VbenFormSchema[] => [
  diySection('组件样式'),
  ...DIY_PADDING_FIELDS,
  ...DIY_RADIUS_FIELDS,
  diyColor('style.bgcolor', '底部背景：', '#f2f2f2', '透明'),
  diySlider('style.height', '视频高度：', { max: 800 }),
  diySection('视频设置'),
  diyInput('params.videoUrl', '视频地址：'),
  diyRadioGroup('params.autoplay', '自动播放：', [
    { label: '否', value: 0 },
    { label: '是', value: 1 },
  ]),
]);

const { Form } = useDiyCurItemForm(
  () => props.curItem,
  schema,
  {
    onInit(item) {
      parseIntFields(item, ['style.paddingTop', 'style.height']);
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
    <div class="gray ml-[100px]">滑块可用键盘的左右方向键精确调整</div>
    <div class="form-item ml-[100px] mt-2">
      <div class="form-label mb-2">视频封面：</div>
      <div class="diy-video-cover">
        <img
          v-img-url="curItem.params?.poster"
          alt=""
          @click="editor.onEditorSelectImage(curItem.params!, 'poster')"
        />
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.diy-video-cover {
  width: 200px;
}
.diy-video-cover img {
  width: 100%;
}
</style>
