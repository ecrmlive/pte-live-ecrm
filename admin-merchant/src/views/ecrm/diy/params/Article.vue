<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { getShopLinkArticleCategoryApi } from '#/api/core/shop-link';

import {
  diyBgColors,
  diyInput,
  diyRadioGroup,
  diySection,
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { parseIntFields } from './shared/path-utils';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';

defineOptions({ name: 'DiyParamsArticle' });

const props = defineProps<{
  curItem: Record<string, unknown>;
  opts?: unknown;
  selectedIndex?: number;
}>();

const router = useRouter();
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
    help: '文章数据请到内容管理 - 文章列表中管理',
  },
  diyRadioGroup('style.display', '显示类型：', [
    { label: '有图模式', value: 10 },
    { label: '无图模式', value: 20 },
  ]),
]);

const { Form } = useDiyCurItemForm(
  () => props.curItem,
  schema,
  {
    onInit(item) {
      parseIntFields(item, ['style.display']);
    },
    onValuesChange(values, fieldsChanged) {
      const item = props.curItem;
      for (const path of fieldsChanged) {
        if (path === 'params.auto.showNum') {
          const num = Number(values[path]);
          item.params = item.params ?? {};
          (item.params as Record<string, unknown>).auto = (
            (item.params as Record<string, unknown>).auto ?? {}
          ) as Record<string, unknown>;
          ((item.params as Record<string, unknown>).auto as Record<string, unknown>).showNum =
            num <= 1 ? 1 : num;
        } else {
          const keys = path.split('.');
          let current: Record<string, unknown> = item;
          for (let i = 0; i < keys.length - 1; i++) {
            current = current[keys[i]!] as Record<string, unknown>;
          }
          current[keys[keys.length - 1]!] = values[path];
        }
      }
    },
  },
);

function gotoArticle() {
  void router.push('/plus/article/index');
}

onMounted(async () => {
  try {
    const res = await getShopLinkArticleCategoryApi();
    category.value = res.list ?? [];
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
    <div class="gray ml-[100px]">
      文章数据请到
      <span class="cursor-pointer text-primary" @click="gotoArticle">内容管理 - 文章列表</span>
      中管理
    </div>
  </div>
</template>
