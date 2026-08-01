<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, watch } from 'vue';

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

defineOptions({ name: 'DiyParamsQixiLive' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    params?: Record<string, unknown>;
    style?: Record<string, unknown>;
  };
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

function ensureParams() {
  const params = (props.curItem.params ?? {}) as Record<string, unknown>;
  if (!params.title) params.title = '热门直播';
  if (!params.moreTitle) params.moreTitle = '更多';
  if (params.showMore === undefined || params.showMore === null || params.showMore === '') {
    params.showMore = 1;
  }
  props.curItem.params = params;
}

const showMoreLink = computed(
  () =>
    Number(props.curItem.params?.showMore) === 1 ||
    props.curItem.params?.showMore === '1',
);

const headerSchema = computed((): VbenFormSchema[] => [
  diySection('头部设置'),
  diyInput('params.title', '标题文字：', { maxlength: 8, showWordLimit: true }),
  diyRadioGroup('params.showMore', '显示更多：', [
    { label: '显示', value: 1 },
    { label: '隐藏', value: 0 },
  ]),
  ...(showMoreLink.value
    ? [
        diyInput('params.moreTitle', '按钮文字：', {
          maxlength: 6,
          showWordLimit: true,
        }),
      ]
    : []),
]);

const restSchema = computed((): VbenFormSchema[] => [
  diySection('边距设置'),
  ...diyBgColors('#ffffff', '#f2f2f2'),
  ...DIY_PADDING_FIELDS,
  ...DIY_RADIUS_FIELDS,
  diySection('组件设置'),
  diyInput('params.showNum', '显示数量：'),
]);

const { Form: HeaderForm } = useDiyCurItemForm(() => props.curItem, headerSchema, {
  onInit: ensureParams,
  fieldPaths: ['params.title', 'params.moreTitle', 'params.showMore'],
});

const { Form: RestForm } = useDiyCurItemForm(() => props.curItem, restSchema, {
  fieldPaths: [
    'params.showNum',
    'style.background',
    'style.bgcolor',
    'style.paddingTop',
    'style.paddingBottom',
    'style.paddingLeft',
    'style.topRadio',
    'style.bottomRadio',
    'style.background_image',
  ],
});

watch(
  () => props.curItem,
  () => ensureParams(),
  { deep: true },
);
</script>

<template>
  <div>
    <HeaderForm />
    <div v-if="showMoreLink" class="diy-field-hint">
      标题和按钮文字颜色跟随主题色
    </div>
    <RestForm />
    <div class="form-chink" />
    <div class="form-subtitle">组件样式</div>
    <div class="form-item form-item--stacked">
      <div class="form-label">背景图片：</div>
      <div class="form-item__body">
        <div class="diy-notice-icon">
          <img
            v-img-url="curItem.style?.background_image"
            alt=""
            @click="editor.onEditorSelectImage(curItem.style!, 'background_image')"
          />
        </div>
        <div class="diy-field-hint diy-field-hint--inline">建议尺寸 700px × 90px</div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.diy-notice-icon {
  width: 100%;
  max-width: 320px;
}

.diy-notice-icon img {
  display: block;
  width: 100%;
  height: 45px;
  object-fit: cover;
  border-radius: 4px;
  cursor: pointer;
}
</style>
