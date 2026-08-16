<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { ArrowRight, CloseBold } from '@element-plus/icons-vue';
import { ElIcon } from 'element-plus';
import { computed, ref } from 'vue';
import draggable from 'vuedraggable';

import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import {
  diyColor,
  diyRadioGroup,
  diySection,
  diySlider,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyInputField from './shared/diy-input-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsBanner' });

type DiyRecord = Record<string, unknown>;

const { PrimaryButton } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: DiyRecord & {
    data?: DiyRecord[];
    params?: DiyRecord;
    style?: DiyRecord;
  };
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const styleType = ref<'content' | 'style'>('content');
const isLinkset = ref(false);
const linkIndex = ref(0);
const linkData = ref<DiyRecord | null>(null);

function ensureBannerDefaults(item: DiyRecord) {
  const params = (item.params ??= {}) as DiyRecord;
  const style = (item.style ??= {}) as DiyRecord;

  if (!params.bannerStyle) params.bannerStyle = 'style1';
  if (!style.indicatorStyle) {
    const legacyShape = String(style.imgShape ?? 'round');
    style.indicatorStyle =
      legacyShape === 'square' ? 'style2' : legacyShape === 'rectangle' ? 'style3' : 'style1';
  }
  if (!style.indicatorPosition) style.indicatorPosition = style.btnShape ?? 'center';
  if (!style.indicatorTone) style.indicatorTone = 'theme';
  if (!style.indicatorColor) style.indicatorColor = style.btnColor ?? '#ffffff';
  if (!style.radiusMode) style.radiusMode = 'all';
  if (style.imageRadius === undefined) style.imageRadius = style.topRadio ?? 0;
  if (style.imageShadow === undefined) style.imageShadow = 'off';
  if (style.float === undefined) style.float = 0;
  if (!style.background) style.background = '#ffffff';
  if (style.paddingTop === undefined) style.paddingTop = 0;
  if (style.paddingBottom === undefined) style.paddingBottom = 0;
  if (style.paddingLeft === undefined) style.paddingLeft = 0;
  if (style.paddingRight === undefined) style.paddingRight = style.paddingLeft;
  if (style.marginTop === undefined) style.marginTop = 0;
  if (style.cardShadow === undefined) style.cardShadow = 'off';
  if (!style.height) style.height = 340;
}

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('展示设置'),
  diyRadioGroup(
    'params.bannerStyle',
    '选择风格：',
    [
      { label: '样式一', value: 'style1' },
      { label: '样式二', value: 'style2' },
      { label: '样式三', value: 'style3' },
    ],
  ),
  diySection('图片设置', '建议上传尺寸相同的图片，建议尺寸750px × 340px'),
  diySlider('style.height', '图片高度：', { max: 1900, min: 100 }),
]);

const styleSchema = computed((): VbenFormSchema[] => {
  const style = props.curItem.style ?? {};
  return [
    diySection('指示器设置'),
    diyRadioGroup(
      'style.indicatorStyle',
      '指示器样式：',
      [
        { label: '样式一', value: 'style1' },
        { label: '样式二', value: 'style2' },
        { label: '样式三', value: 'style3' },
        { label: '样式四', value: 'style4' },
      ],
    ),
    diyRadioGroup(
      'style.indicatorPosition',
      '指示器位置：',
      [
        { label: '左对齐', value: 'left' },
        { label: '居中对齐', value: 'center' },
        { label: '右对齐', value: 'right' },
      ],
    ),
    diyRadioGroup(
      'style.indicatorTone',
      '色调：',
      [
        { label: '跟随主题风格', value: 'theme' },
        { label: '自定义', value: 'custom' },
      ],
    ),
    ...(style.indicatorTone === 'custom'
      ? [diyColor('style.indicatorColor', '指示器颜色：', '#ffffff')]
      : []),
    diySection('图片设置'),
    diyRadioGroup(
      'style.radiusMode',
      '图片圆角：',
      [
        { label: '全部', value: 'all' },
        { label: '分别设置', value: 'individual' },
      ],
      true,
    ),
    ...(style.radiusMode === 'individual'
      ? [
          diySlider('style.topLeftRadio', '左上圆角：', { max: 48, min: 0 }),
          diySlider('style.topRightRadio', '右上圆角：', { max: 48, min: 0 }),
          diySlider('style.bottomLeftRadio', '左下圆角：', { max: 48, min: 0 }),
          diySlider('style.bottomRightRadio', '右下圆角：', { max: 48, min: 0 }),
        ]
      : [diySlider('style.imageRadius', '圆角值：', { max: 48, min: 0 })]),
    diyRadioGroup(
      'style.imageShadow',
      '开启阴影：',
      [
        { label: '关闭', value: 'off' },
        { label: '开启', value: 'on' },
      ],
    ),
    diySection('卡片样式'),
    diySlider('style.float', '组件上浮：', { max: 48, min: 0 }),
    diyColor('style.background', '底部背景：', '#f62c2c', '透明'),
    diySlider('style.paddingTop', '上边距：', { max: 48, min: 0 }),
    diySlider('style.paddingBottom', '下边距：', { max: 48, min: 0 }),
    diySlider('style.paddingLeft', '左右边距：', { max: 48, min: 0 }),
    diySlider('style.marginTop', '页面上间距：', { max: 96, min: 0 }),
    diyRadioGroup(
      'style.cardShadow',
      '开启阴影：',
      [
        { label: '关闭', value: 'off' },
        { label: '开启', value: 'on' },
      ],
    ),
  ];
});

const { Form: ContentForm } = useDiyCurItemForm(() => props.curItem, contentSchema, {
  onInit: ensureBannerDefaults,
});
const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema, {
  onInit: ensureBannerDefaults,
});

function changeLink(index: number) {
  isLinkset.value = true;
  linkIndex.value = index;
  linkData.value = props.curItem.data?.[index] ?? null;
}

function closeLinkset(e: { name?: string; type?: string; url?: string } | null) {
  isLinkset.value = false;
  if (e && props.curItem.data?.[linkIndex.value]) {
    const row = props.curItem.data[linkIndex.value];
    row.linkeType = e.type;
    row.linkUrl = e.url;
    row.name = e.name;
  }
}
</script>

<template>
  <div>
    <div class="common-form common-form-new">
      <span>{{ curItem.name }}</span>
      <div class="diy-changes">
        <div class="diy-change" :class="{ active: styleType === 'content' }" @click="styleType = 'content'">
          内容
        </div>
        <div class="diy-change" :class="{ active: styleType === 'style' }" @click="styleType = 'style'">
          样式
        </div>
      </div>
    </div>

    <div v-show="styleType === 'content'">
      <ContentForm />
      <template v-if="curItem.data && curItem.data.length > 0">
        <draggable v-model="curItem.data" class="draggable-list" group="people" item-key="index">
          <template #item="{ element, index }">
            <div class="d-c-c param-img-item navbar">
              <div class="right pr">
                <div class="icon param-img-thumb">
                  <ElIcon
                    class="el-icon-DeleteFilled"
                    @click.stop="editor.onEditorDeleleData(index, selectedIndex ?? 0)"
                  >
                    <CloseBold />
                  </ElIcon>
                  <img
                    v-img-url="element.imgUrl"
                    alt="轮播图片"
                    :style="{ height: Number(curItem.style?.height ?? 340) * 0.5 + 'px' }"
                    @click="editor.onEditorSelectImage(element, 'imgUrl')"
                  />
                </div>
                <div class="d-s-c ww100">
                  <div class="url-box flex-1 d-s-c">
                    <span class="key-name">链接</span>
                    <DiyInputField v-model="element.linkUrl" style="padding-bottom: 10px">
                      <template #suffix>
                        <ElIcon color="#333" size="16px" @click="changeLink(index)">
                          <ArrowRight />
                        </ElIcon>
                      </template>
                    </DiyInputField>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </draggable>
      </template>
      <div class="d-c-c pb16">
        <component :is="PrimaryButton" plain @click="editor.onEditorAddData()">+ 添加图片</component>
      </div>
    </div>

    <div v-show="styleType === 'style'">
      <StyleForm />
    </div>

    <DiyLinkPickerDialog
      v-if="isLinkset"
      :is_linkset="isLinkset"
      :link-data="linkData"
      @close-dialog="closeLinkset"
    >
      选择链接
    </DiyLinkPickerDialog>
  </div>
</template>

<style lang="scss" scoped>
.param-img-item.navbar {
  min-height: 132px;
  height: auto;
}

.param-img-item.navbar .param-img-thumb {
  position: relative;
  display: inline-block;
  line-height: 0;

  .el-icon-DeleteFilled {
    right: 8px;
    top: 8px;
  }
}

.param-img-item.navbar .icon img {
  display: block;
  width: 408px;
  min-height: 110px;
  max-height: 260px;
  background: #eeeeee;
  margin-top: 0;
  margin-bottom: 10px;
  object-fit: cover;
}
</style>
