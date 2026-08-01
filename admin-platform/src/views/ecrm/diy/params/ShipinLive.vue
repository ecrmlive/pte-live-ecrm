<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed } from 'vue';

import {
  diyInput,
  diySection,
  diySlider,
} from './shared/schema-helpers';
import { parseIntFields } from './shared/path-utils';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsShipinLive' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    params?: Record<string, unknown>;
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

const schema = computed((): VbenFormSchema[] => [
  diySection('样式设置'),
  diySlider('style.bottom', '底边距：'),
  diySlider('style.right', '右边距：'),
  diySlider('style.opacity', '不透明度：'),
  diyInput('params.finderUserName', '视频号ID：'),
]);

const { Form } = useDiyCurItemForm(
  () => props.curItem,
  schema,
  {
    onInit(item) {
      parseIntFields(item, ['style.bottom', 'style.right', 'style.opacity']);
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
    <div class="gray ml-[100px]">视频号ID,在登录视频号助手首页中获取</div>
    <div class="form-item ml-[100px] mt-2">
      <div class="form-label mb-2">视频号图标：</div>
      <div class="diy-service-icon">
        <img
          v-img-url="curItem.params?.image"
          alt=""
          @click="editor.onEditorSelectImage(curItem.params!, 'image')"
        />
      </div>
      <div class="gray">建议尺寸90×90</div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.diy-service-icon,
.diy-service-icon img {
  width: 40px;
  height: 40px;
}
</style>
