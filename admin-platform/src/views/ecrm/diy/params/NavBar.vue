<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, ref } from 'vue';

import { ArrowRight, CloseBold, Rank } from '@element-plus/icons-vue';
import { ElIcon, ElMessage, ElSwitch } from 'element-plus';
import draggable from 'vuedraggable';

import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import DiyInputField from './shared/diy-input-field.vue';
import {
  diyColor,
  diyRadioGroup,
  diySection,
  diySlider,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsNavBar' });

type DiyRecord = Record<string, any>;

const { PrimaryButton } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: DiyRecord & { data?: DiyRecord[]; params?: DiyRecord; style?: DiyRecord };
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const styleType = ref<'content' | 'style'>('content');
const isLinkset = ref(false);
const linkIndex = ref(0);
const linkData = ref<DiyRecord | null>(null);

const navigationItems = computed<DiyRecord[]>({
  get: () => props.curItem.data ?? [],
  set: (value) => {
    props.curItem.data = value;
  },
});

function ensureNavBarDefaults(item: DiyRecord) {
  const style = (item.style ??= {}) as DiyRecord;

  if (!Array.isArray(item.data)) item.data = [];
  if (!style.navigationType) style.navigationType = 'icon-text';
  if (!style.displayMode) style.displayMode = 'fixed';
  const rowsNum = Number(style.rowsNum);
  if (!Number.isFinite(rowsNum) || rowsNum < 3 || rowsNum > 5) style.rowsNum = 5;

  if (style.iconRadiusMode === undefined) style.iconRadiusMode = 'all';
  if (style.iconRadius === undefined) style.iconRadius = style.topRadio ?? 8;
  if (style.iconTopLeftRadius === undefined) style.iconTopLeftRadius = style.iconRadius;
  if (style.iconTopRightRadius === undefined) style.iconTopRightRadius = style.iconRadius;
  if (style.iconBottomRightRadius === undefined) style.iconBottomRightRadius = style.iconRadius;
  if (style.iconBottomLeftRadius === undefined) style.iconBottomLeftRadius = style.iconRadius;
  if (!style.iconShadow) style.iconShadow = 'off';
  if (!style.iconShadowColor) style.iconShadowColor = 'rgba(20, 37, 63, 0.16)';
  if (style.iconShadowX === undefined) style.iconShadowX = 0;
  if (style.iconShadowY === undefined) style.iconShadowY = 4;
  if (style.iconShadowBlur === undefined) style.iconShadowBlur = 10;
  if (style.iconShadowSpread === undefined) style.iconShadowSpread = 0;

  if (style.float === undefined) style.float = 0;
  if (style.bgcolor === undefined) style.bgcolor = 'rgba(255, 255, 255, 0)';
  if (style.bgcolorEnd === undefined) style.bgcolorEnd = 'rgba(255, 255, 255, 0)';
  if (style.background === undefined) style.background = 'rgba(255, 0, 0, 0)';
  if (!style.textColor) style.textColor = '#333333';
  if (style.cardBorderWidth === undefined) style.cardBorderWidth = 0;
  if (!style.cardBorderColor) style.cardBorderColor = 'rgba(0, 0, 0, 0)';
  if (style.paddingTop === undefined) style.paddingTop = 0;
  if (style.paddingBottom === undefined) style.paddingBottom = 10;
  if (style.paddingLeft === undefined) style.paddingLeft = 0;
  if (style.paddingRight === undefined) style.paddingRight = style.paddingLeft;
  if (style.marginTop === undefined) style.marginTop = 0;
  if (!style.cardRadiusMode) style.cardRadiusMode = 'all';
  if (style.cardRadius === undefined) style.cardRadius = 0;
  if (style.cardTopLeftRadius === undefined) style.cardTopLeftRadius = style.cardRadius;
  if (style.cardTopRightRadius === undefined) style.cardTopRightRadius = style.cardRadius;
  if (style.cardBottomRightRadius === undefined) style.cardBottomRightRadius = style.cardRadius;
  if (style.cardBottomLeftRadius === undefined) style.cardBottomLeftRadius = style.cardRadius;
  if (!style.cardShadow) style.cardShadow = 'off';
  if (!style.cardShadowColor) style.cardShadowColor = 'rgba(20, 37, 63, 0.12)';
  if (style.cardShadowX === undefined) style.cardShadowX = 0;
  if (style.cardShadowY === undefined) style.cardShadowY = 8;
  if (style.cardShadowBlur === undefined) style.cardShadowBlur = 18;
  if (style.cardShadowSpread === undefined) style.cardShadowSpread = 0;

  item.data.forEach((entry) => {
    if (!entry._navUid) entry._navUid = `nav-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
    if (!entry.text) entry.text = '导航名称';
    if (entry.hide === undefined) entry.hide = false;
  });
}

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('展示设置'),
  diyRadioGroup('style.navigationType', '导航样式：', [
    { label: '图片加文字', value: 'icon-text' },
    { label: '图片', value: 'icon' },
    { label: '文字', value: 'text' },
  ]),
  diyRadioGroup('style.rowsNum', '单行显示：', [
    { label: '3个', value: 3 },
    { label: '4个', value: 4 },
    { label: '5个', value: 5 },
  ]),
  diyRadioGroup('style.displayMode', '展示样式：', [
    { label: '固定显示', value: 'fixed' },
    { label: '分页滑动', value: 'page' },
  ]),
  diySection('内容设置', '最多可添加13张图片，建议宽度90 × 90px'),
]);

const styleSchema = computed((): VbenFormSchema[] => {
  const style = props.curItem.style ?? {};
  return [
    diySection('图标样式'),
    diyRadioGroup(
      'style.iconRadiusMode',
      '背景圆角：',
      [
        { label: '全部', value: 'all' },
        { label: '单个', value: 'individual' },
      ],
      true,
    ),
    ...(style.iconRadiusMode === 'individual'
      ? [
          diySlider('style.iconTopLeftRadius', '左上圆角：', { max: 48, min: 0 }),
          diySlider('style.iconTopRightRadius', '右上圆角：', { max: 48, min: 0 }),
          diySlider('style.iconBottomRightRadius', '右下圆角：', { max: 48, min: 0 }),
          diySlider('style.iconBottomLeftRadius', '左下圆角：', { max: 48, min: 0 }),
        ]
      : [diySlider('style.iconRadius', '圆角值：', { max: 48, min: 0 })]),
    diyRadioGroup('style.iconShadow', '开启阴影：', [
      { label: '关闭', value: 'off' },
      { label: '开启', value: 'on' },
    ]),
    ...(style.iconShadow === 'on'
      ? [
          diyColor('style.iconShadowColor', '阴影颜色：', 'rgba(20, 37, 63, 0.16)'),
          diySlider('style.iconShadowX', '横轴：', { max: 24, min: -24 }),
          diySlider('style.iconShadowY', '纵轴：', { max: 24, min: -24 }),
          diySlider('style.iconShadowBlur', '宽度：', { max: 48, min: 0 }),
          diySlider('style.iconShadowSpread', '扩散：', { max: 48, min: -24 }),
        ]
      : []),
    diySection('卡片样式'),
    diySlider('style.float', '组件上浮：', { max: 48, min: 0 }),
    diyColor('style.bgcolor', '组件背景：', 'rgba(255, 255, 255, 0)', '透明'),
    diyColor('style.bgcolorEnd', '组件背景渐变：', 'rgba(255, 255, 255, 0)', '透明'),
    diyColor('style.background', '底部背景：', 'rgba(255, 0, 0, 0)', '透明'),
    diyColor('style.textColor', '文字颜色：', '#333333'),
    diySlider('style.cardBorderWidth', '边框宽度：', { max: 8, min: 0 }),
    diyColor('style.cardBorderColor', '边框颜色：', 'rgba(0, 0, 0, 0)', '透明'),
    diySlider('style.paddingTop', '上边距：', { max: 48, min: 0 }),
    diySlider('style.paddingBottom', '下边距：', { max: 48, min: 0 }),
    diySlider('style.paddingLeft', '左右边距：', { max: 48, min: 0 }),
    diySlider('style.marginTop', '页面上间距：', { max: 96, min: 0 }),
    diyRadioGroup(
      'style.cardRadiusMode',
      '背景圆角：',
      [
        { label: '全部', value: 'all' },
        { label: '单个', value: 'individual' },
      ],
      true,
    ),
    ...(style.cardRadiusMode === 'individual'
      ? [
          diySlider('style.cardTopLeftRadius', '左上圆角：', { max: 48, min: 0 }),
          diySlider('style.cardTopRightRadius', '右上圆角：', { max: 48, min: 0 }),
          diySlider('style.cardBottomRightRadius', '右下圆角：', { max: 48, min: 0 }),
          diySlider('style.cardBottomLeftRadius', '左下圆角：', { max: 48, min: 0 }),
        ]
      : [diySlider('style.cardRadius', '圆角值：', { max: 48, min: 0 })]),
    diyRadioGroup('style.cardShadow', '开启阴影：', [
      { label: '关闭', value: 'off' },
      { label: '开启', value: 'on' },
    ]),
    ...(style.cardShadow === 'on'
      ? [
          diyColor('style.cardShadowColor', '阴影颜色：', 'rgba(20, 37, 63, 0.12)'),
          diySlider('style.cardShadowX', '横轴：', { max: 24, min: -24 }),
          diySlider('style.cardShadowY', '纵轴：', { max: 24, min: -24 }),
          diySlider('style.cardShadowBlur', '宽度：', { max: 48, min: 0 }),
          diySlider('style.cardShadowSpread', '扩散：', { max: 48, min: -24 }),
        ]
      : []),
  ];
});

const { Form: ContentForm } = useDiyCurItemForm(() => props.curItem, contentSchema, {
  onInit: ensureNavBarDefaults,
});
const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema, {
  onInit: ensureNavBarDefaults,
});

function changeLink(index: number) {
  isLinkset.value = true;
  linkIndex.value = index;
  linkData.value = navigationItems.value[index] ?? null;
}

function closeLinkset(value: { name?: string; type?: string; url?: string } | null) {
  isLinkset.value = false;
  const item = navigationItems.value[linkIndex.value];
  if (!value || !item) return;
  item.linkeType = value.type;
  item.linkUrl = value.url;
  item.name = value.name;
  item.text = value.name || item.text;
}

function addNavigationItem() {
  if (navigationItems.value.length >= 13) {
    ElMessage.warning('导航组最多添加13项');
    return;
  }
  navigationItems.value.push({
    _navUid: `nav-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`,
    hide: false,
    imgUrl: '',
    linkUrl: '',
    text: '导航名称',
  });
}

function removeNavigationItem(index: number) {
  if (navigationItems.value.length <= 1) {
    ElMessage.warning('导航组至少保留 1 项');
    return;
  }
  navigationItems.value.splice(index, 1);
}

function setItemEnabled(item: DiyRecord, enabled: boolean) {
  item.hide = !enabled;
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
    <draggable
      v-model="navigationItems"
      class="draggable-list"
      group="navigation-items"
      item-key="_navUid"
      :animation="180"
      handle=".nav-item-card__drag"
    >
        <template #item="{ element, index }">
          <div class="d-c-c param-img-item navbar nav-item-card">
            <div class="nav-item-card__drag" title="拖拽调整顺序">
              <ElIcon><Rank /></ElIcon>
            </div>
            <div class="d-c d-c-c nav-item-card__image">
              <div class="icon">
                <img
                  v-img-url="element.imgUrl"
                  alt=""
                  @click="editor.onEditorSelectImage(element, 'imgUrl')"
                />
              </div>
            </div>
            <div class="right nav-item-card__fields">
              <ElIcon
                class="el-icon-DeleteFilled"
                @click.stop="removeNavigationItem(index)"
              >
                <CloseBold />
              </ElIcon>
              <div class="url-box mb16 flex-1 d-s-c ww100">
                <span class="key-name">标题</span>
                <DiyInputField v-model="element.text" :maxlength="6" show-word-limit />
              </div>
              <div class="url-box mb16 flex-1 d-s-c ww100">
                <span class="key-name">链接</span>
                <DiyLinkInputField v-model="element.linkUrl" @click="changeLink(index)">
                  <template #suffix>
                    <ElIcon color="#333" size="16px"><ArrowRight /></ElIcon>
                  </template>
                </DiyLinkInputField>
              </div>
              <div class="url-box flex-1 d-s-c ww100">
                <span class="key-name">状态</span>
                <ElSwitch
                  :model-value="!element.hide"
                  @update:model-value="setItemEnabled(element, $event as boolean)"
                />
              </div>
            </div>
          </div>
        </template>
      </draggable>
      <div class="d-c-c pb16">
        <component :is="PrimaryButton" plain @click="addNavigationItem">+ 添加</component>
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
.nav-item-card {
  align-items: center;
  min-height: 132px;
}

.nav-item-card__drag {
  display: flex;
  width: 24px;
  align-items: center;
  justify-content: center;
  color: #b7c0cf;
  cursor: move;
}

.nav-item-card__image {
  width: 108px;
  margin-right: 16px;
}

.nav-item-card__image .icon img {
  display: block;
  width: 90px;
  height: 90px;
  margin: 0;
  cursor: pointer;
  border: 1px dashed #d9dfe9;
  border-radius: 8px;
  object-fit: cover;
}

.nav-item-card__fields {
  flex: 1;
}
</style>
