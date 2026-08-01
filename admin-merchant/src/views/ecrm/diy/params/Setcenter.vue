<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed } from 'vue';

import ImageField from '#/components/shop/image-field.vue';

import { diyInput, diySection } from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';

defineOptions({ name: 'DiyParamsSetcenter' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    params?: Record<string, unknown>;
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const shareImg = computed({
  get: () => String(props.curItem.params?.share_img ?? ''),
  set: (value: string) => {
    if (props.curItem.params) {
      props.curItem.params.share_img = value;
    }
  },
});

const schema = computed((): VbenFormSchema[] => [
  diySection('样式设置'),
  diyInput('params.name', '页面名称：'),
  diyInput('params.share_title', '分享标题：'),
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema);
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <Form />
    <p class="gray ml-[60px] mb-4">页面名称仅用于后台查找</p>
    <p class="gray ml-[60px] mb-4">小程序端转发时显示的标题</p>
    <div class="form-item ml-[60px]">
      <div class="form-label mb-2">分享logo：</div>
      <ImageField
        v-model="shareImg"
        hint="公众号分享logo，建议尺寸80×80"
        :preview-size="80"
      />
    </div>
  </div>
</template>
