<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { getShopLinkArticleCategoryApi } from '#/api/core/shop-link';

import {
  diyBgColors,
  diyColor,
  diyRadioGroup,
  diySection,
  diySlider,
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
  diySection('展示设置'),
  diyRadioGroup('params.layout', '选择风格：', [
    { label: '大图展示', value: 'list' },
    { label: '两列展示（纵向）', value: 'grid' },
    { label: '两列展示（横向）', value: 'scroll' },
  ]),
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
  diyRadioGroup('params.showDate', '时间日期：', [{ label: '显示', value: true }, { label: '隐藏', value: false }]),
  diyRadioGroup('params.showViews', '浏览量：', [{ label: '显示', value: true }, { label: '隐藏', value: false }]),
  diySection('列表样式'),
  diySlider('style.imageRadius', '图片圆角：', { min: 0, max: 32 }),
  diyRadioGroup('style.titleWeight', '文章标题：', [{ label: '加粗', value: 'bold' }, { label: '常规', value: 'normal' }]),
  diyColor('style.titleColor', '文章标题：', '#333333'),
  diyColor('style.metaColor', '时间日期：', '#999999'),
  diyColor('style.viewColor', '浏览元素：', '#999999'),
  diyRadioGroup('style.shadow', '开启阴影：', [{ label: '关闭', value: 'off' }, { label: '开启', value: 'on' }]),
  diySection('卡片样式'),
  ...diyBgColors('#ffffff', '#F5F5F5'),
  ...DIY_PADDING_FIELDS,
  ...DIY_RADIUS_FIELDS,
]);

const { Form } = useDiyCurItemForm(
  () => props.curItem,
  schema,
  {
    onInit(item) {
      item.params = item.params ?? {};
      item.style = item.style ?? {};
      Object.assign(item.params as Record<string, unknown>, { layout: 'list', showNum: 3, showDate: true, showViews: true, ...(item.params as Record<string, unknown>) });
      Object.assign(item.style as Record<string, unknown>, { bgcolor: '#F5F5F5', background: '#ffffff', paddingTop: 0, paddingBottom: 0, paddingLeft: 10, marginTop: 10, radius: 0, imageRadius: 0, titleColor: '#333333', metaColor: '#999999', viewColor: '#999999', titleWeight: 'normal', shadow: 'off', ...(item.style as Record<string, unknown>) });
      parseIntFields(item, ['style.imageRadius']);
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
  void router.push('/cms/article');
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
