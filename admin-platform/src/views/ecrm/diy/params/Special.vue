<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, ref } from 'vue';

import { getArticleCategoryListApi } from '#/api/core/plus-article';

import {
  diyBgColors,
  diyInput,
  diyRadioGroup,
  diySection,
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsSpecial' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    style?: Record<string, unknown>;
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const category = ref<Array<{ category_id: number; name: string }>>([]);

const schema = computed((): VbenFormSchema[] => [
  diySection('边距设置'),
  ...diyBgColors('#ffffff', '#F2F2F2'),
  ...DIY_PADDING_FIELDS,
  ...DIY_RADIUS_FIELDS,
  diySection('文章设置'),
  {
    component: 'Select',
    componentProps: {
      options: [
        { label: '全部分类', value: 0 },
        ...category.value.map((item) => ({
          label: item.name,
          value: item.category_id,
        })),
      ],
      style: { width: '100%' },
    },
    fieldName: 'params.auto.category',
    label: '文章分类：',
  },
  {
    component: 'Input',
    componentProps: { min: 1, style: { width: '100%' }, type: 'number' },
    fieldName: 'params.auto.showNum',
    label: '显示数量：',
  },
  diyRadioGroup('style.display', '显示类型：', [
    { label: '单行', value: 1 },
    { label: '两行', value: 2 },
  ]),
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema);

onMounted(async () => {
  try {
    const res = await getArticleCategoryListApi({});
    category.value = res.category ?? [];
  } catch {
    category.value = [];
  }
});
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <Form />
    <div class="form-item ml-[100px]">
      <div class="form-label mb-2">图片：</div>
      <div class="diy-special-cover">
        <img
          v-img-url="curItem.style?.image"
          alt=""
          @click="editor.onEditorSelectImage(curItem.style!, 'image')"
        />
        <div class="gray">建议尺寸140×38</div>
      </div>
    </div>
    <div class="gray ml-[100px]">
      文章数据请到
      <a href="#/plus/article/index" target="_blank">内容管理 - 文章列表</a>
      中管理
    </div>
  </div>
</template>

<style lang="scss" scoped>
.diy-special-cover img {
  width: 140px;
  height: 38px;
}
</style>
