<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { ElButton, ElDialog } from 'element-plus';
import { computed, ref } from 'vue';

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

defineOptions({ name: 'DiyParamsCoupon' });

type DiyItem = Record<string, any>;

const props = defineProps<{
  curItem: DiyItem;
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const activePanel = ref<'content' | 'style'>('content');
const stylePickerVisible = ref(false);
const pendingStyle = ref(1);

const couponStyles = [
  { value: 1, label: '风格1' },
  { value: 2, label: '风格2' },
  { value: 3, label: '风格3' },
  { value: 4, label: '风格4' },
  { value: 5, label: '风格5' },
];

function ensureCouponDefaults(item: DiyItem) {
  item.params ??= {};
  item.style ??= {};
  const { params, style } = item;

  params.limit ??= 5;
  params.btntext ??= '立即领取';
  style.type ??= 1;
  style.colorMode ??= 'theme';
  style.couponSpacing ??= 6;
  style.shadow ??= 'off';
  style.descolor ??= '#666666';
  style.pricecolor ??= '#ff4c01';
  style.cillcolor ??= '#999999';
  style.btncolor ??= '#ff4c01';
  style.btnTxtcolor ??= '#ffffff';
  style.btnRadio ??= 12;
  style.bgtype ??= 1;
  style.background ??= '#ffffff';
  style.bgcolor ??= '#f5f5f5';
  style.paddingTop ??= 0;
  style.paddingBottom ??= 0;
  style.paddingLeft ??= 10;
  style.topRadio ??= 8;
  style.bottomRadio ??= 8;
}

const selectedStyleLabel = computed(() => {
  const value = Number(props.curItem.style?.type ?? 1);
  return couponStyles.find((item) => item.value === value)?.label ?? '风格1';
});

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('展示设置'),
  diySlider('params.limit', '展示数量：', { max: 30, min: 1 }),
  diyInput('params.btntext', '领取按钮：', { maxlength: 8 }),
]);

const styleSchema = computed((): VbenFormSchema[] => {
  const style = props.curItem.style ?? {};
  const fields: VbenFormSchema[] = [
    diySection('优惠券样式'),
    diyRadioGroup(
      'style.colorMode',
      '色调：',
      [
        { label: '跟随主题风格', value: 'theme' },
        { label: '自定义', value: 'custom' },
      ],
      true,
    ),
  ];
  if (style.colorMode === 'custom') {
    fields.push(
      diyColor('style.pricecolor', '面额颜色：', '#ff4c01'),
      diyColor('style.btncolor', '按钮颜色：', '#ff4c01'),
      diyColor('style.btnTxtcolor', '按钮文字：', '#ffffff'),
    );
  }
  fields.push(
    diySlider('style.couponSpacing', '优惠券间距：', { max: 24, min: 0 }),
    diyRadioGroup(
      'style.shadow',
      '开启阴影：',
      [
        { label: '关闭', value: 'off' },
        { label: '开启', value: 'on' },
      ],
      true,
    ),
    diySection('卡片样式'),
    diyRadioGroup(
      'style.bgtype',
      '组件背景：',
      [
        { label: '背景色', value: 1 },
        { label: '背景图片', value: 2 },
      ],
      true,
    ),
  );
  if (Number(style.bgtype) === 1) {
    fields.push(diyColor('style.background', '组件背景：', '#ffffff'));
  }
  fields.push(
    diyColor('style.bgcolor', '底部背景：', '#f5f5f5'),
    ...DIY_PADDING_FIELDS,
    ...DIY_RADIUS_FIELDS,
  );
  return fields;
});

const { Form: ContentForm } = useDiyCurItemForm(() => props.curItem, contentSchema, {
  onInit(item) {
    ensureCouponDefaults(item as DiyItem);
    parseIntFields(item, ['params.limit']);
  },
});

const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema, {
  onInit(item) {
    ensureCouponDefaults(item as DiyItem);
    parseIntFields(item, [
      'style.btnRadio',
      'style.couponSpacing',
      'style.paddingBottom',
      'style.paddingLeft',
      'style.paddingTop',
      'style.topRadio',
      'style.bottomRadio',
    ]);
  },
});

function openStylePicker() {
  ensureCouponDefaults(props.curItem);
  pendingStyle.value = Number(props.curItem.style.type ?? 1);
  stylePickerVisible.value = true;
}

function confirmStylePicker() {
  props.curItem.style.type = pendingStyle.value;
  stylePickerVisible.value = false;
}
</script>

<template>
  <div class="coupon-params">
    <div class="common-form coupon-title-row">
      <span>{{ curItem.name }}</span>
      <div class="coupon-tabs">
        <button :class="{ active: activePanel === 'content' }" type="button" @click="activePanel = 'content'">
          内容
        </button>
        <button :class="{ active: activePanel === 'style' }" type="button" @click="activePanel = 'style'">
          样式
        </button>
      </div>
    </div>

    <template v-if="activePanel === 'content'">
      <div class="coupon-style-entry">
        <span>选择风格</span>
        <ElButton type="primary" @click="openStylePicker">修改风格</ElButton>
        <span class="coupon-style-current">当前：{{ selectedStyleLabel }}</span>
      </div>
      <ContentForm />
    </template>

    <template v-else>
      <StyleForm />
      <div v-if="Number(curItem.style?.bgtype) === 2" class="form-item ml-[100px]">
        <div class="form-label mb-2">背景图片：</div>
        <div class="diy-special-cover">
          <img
            v-img-url="curItem.style?.bgimage"
            alt=""
            @click="editor.onEditorSelectImage(curItem.style, 'bgimage')"
          />
          <div class="gray">建议尺寸 710px × 180px</div>
        </div>
      </div>
    </template>

    <ElDialog v-model="stylePickerVisible" class="coupon-style-dialog" title="风格选择器" width="min(94vw, 1120px)" append-to-body>
      <div class="coupon-style-grid">
        <button
          v-for="style in couponStyles"
          :key="style.value"
          class="coupon-style-card"
          :class="{ selected: pendingStyle === style.value }"
          type="button"
          @click="pendingStyle = style.value"
        >
          <div class="coupon-style-preview" :class="`preview-${style.value}`">
            <template v-if="style.value === 2">
              <span>品类券<br /><b>¥80</b><small>去领取</small></span>
              <span>通用券<br /><b>¥80</b><small>去领取</small></span>
              <span>商品券<br /><b>¥80</b><small>去领取</small></span>
            </template>
            <template v-else>
              <span v-for="index in 3" :key="index"><b>¥80</b><small>满500元可用</small><em>去领取</em></span>
            </template>
          </div>
          <strong>{{ style.label }}</strong>
        </button>
      </div>
      <template #footer>
        <div class="coupon-dialog-footer">
          <ElButton @click="stylePickerVisible = false">取消</ElButton>
          <ElButton type="primary" @click="confirmStylePicker">确定</ElButton>
        </div>
      </template>
    </ElDialog>
  </div>
</template>

<style lang="scss" scoped>
.coupon-title-row { display: flex; align-items: center; justify-content: space-between; }
.coupon-tabs { display: flex; overflow: hidden; border-radius: 18px; background: #f7f7f7; }
.coupon-tabs button { min-width: 66px; padding: 7px 16px; color: #666; background: transparent; border: 0; cursor: pointer; }
.coupon-tabs button.active { color: #fff; background: var(--el-color-primary); border-radius: 18px; }
.coupon-style-entry { display: flex; align-items: center; gap: 12px; padding: 14px 20px 8px 100px; color: #555; }
.coupon-style-current { color: #999; font-size: 13px; }
.coupon-style-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 20px; }
.coupon-style-card { padding: 14px; text-align: center; background: #fff; border: 1px solid #e2e5ec; border-radius: 8px; cursor: pointer; }
.coupon-style-card.selected { border-color: var(--el-color-primary); box-shadow: 0 0 0 1px var(--el-color-primary); }
.coupon-style-card strong { display: block; margin-top: 10px; font-weight: 500; color: #333; }
.coupon-style-preview { display: flex; align-items: stretch; justify-content: center; min-height: 116px; padding: 20px 12px; overflow: hidden; background: #f6f6f6; gap: 6px; }
.coupon-style-preview span { display: flex; flex: 1; flex-direction: column; align-items: center; justify-content: center; min-width: 0; color: #ff4c01; background: #fff; border-radius: 8px; }
.coupon-style-preview b { font-size: 20px; }.coupon-style-preview small { margin: 4px 0; color: #999; font-size: 10px; }.coupon-style-preview em { padding: 3px 10px; color: #fff; font-size: 10px; font-style: normal; background: #ff4c01; border-radius: 12px; }
.preview-2 { background: linear-gradient(135deg, #ef3122, #ff7a21); }.preview-2 span { padding: 4px; }.preview-3 span { border: 1px solid #ff8aa4; border-radius: 6px; }.preview-4 { background: #f63c20; }.preview-4 span { border-radius: 6px; }.preview-5 span { border-radius: 0; }
.coupon-dialog-footer { display: flex; justify-content: center; gap: 12px; }
.diy-special-cover img { width: 220px; cursor: pointer; }
@media (max-width: 960px) { .coupon-style-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
</style>
