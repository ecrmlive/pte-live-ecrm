<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { ref, computed } from 'vue';

import {
  diyBgColors,
  diyColor,
  diyInput,
  diyRadioGroup,
  diySection,
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsSearch' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    params?: Record<string, unknown>;
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const styleType = ref<'content' | 'style'>('content');

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('展示设置'),
  diyRadioGroup(
    'params.title_type',
    '选择风格：',
    [
      { label: 'logo+搜索', value: 'image' },
      { label: '标题+搜索', value: 'text' },
      { label: '搜索', value: 'search' },
    ],
    true,
  ),
  ...(props.curItem.params?.title_type === 'text'
    ? [diyInput('params.title', '标题：')]
    : []),
  diySection('搜索内容'),
  diyInput('params.searchText', '提示文字：'),
]);

const styleSchema = computed((): VbenFormSchema[] => {
  const fields: VbenFormSchema[] = [];
  if (props.curItem.params?.title_type === 'text') {
    fields.push(
      diySection('标题文字颜色'),
      diyColor('style.titleTextColor', '标题文字颜色：', '#333'),
    );
  }
  fields.push(
    diySection('搜索框样式'),
    diyColor('style.searchBackGround', '搜索框：', '#ffffff'),
    diyColor('style.searchColor', '提示文字：', '#f2f2f2'),
    diySection('组件样式'),
    ...diyBgColors('#ffffff', '#f2f2f2'),
    ...DIY_PADDING_FIELDS,
    diySection('圆角设置'),
    ...DIY_RADIUS_FIELDS,
  );
  return fields;
});

const { Form: ContentForm } = useDiyCurItemForm(() => props.curItem, contentSchema);
const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema);
</script>

<template>
  <div>
    <div class="common-form common-form-new">
      <span>{{ curItem.name }}</span>
      <div class="diy-changes">
        <div
          class="diy-change"
          :class="{ active: styleType == 'content' }"
          @click="styleType = 'content'"
        >
          内容
        </div>
        <div
          class="diy-change"
          :class="{ active: styleType == 'style' }"
          @click="styleType = 'style'"
        >
          样式
        </div>
      </div>
    </div>
    <div v-show="styleType == 'content'">
      <ContentForm />
      <div v-if="curItem.params?.title_type == 'image'" class="form-item ml-[100px]">
        <div class="form-label mb-2">logo图：</div>
        <div class="diy-setpages-cover">
          <img
            v-img-url="curItem.params?.toplogo"
            :width="120"
            alt=""
            @click="editor.onEditorSelectImage(curItem.params!, 'toplogo')"
          />
          <div>建议尺寸78*64</div>
        </div>
      </div>
    </div>
    <div v-show="styleType == 'style'">
      <StyleForm />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.diy-setpages-cover > img {
  width: 60px;
  height: 60px;
}
</style>
