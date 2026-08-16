<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';
import type { PlatformProduct } from '#/api/core/platform-catalog';

import { Picture } from '@element-plus/icons-vue';
import { ElIcon, ElSwitch } from 'element-plus';
import { computed, ref } from 'vue';

import ProductPickerDialog from '#/components/shop/product-picker-dialog.vue';

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

type DiyRecord = Record<string, any>;

interface DiscountItem {
  enabled: boolean;
  image: string;
  price: string;
  productId: number;
  productName: string;
  title: string;
}

const props = defineProps<{
  curItem: DiyRecord & {
    params?: DiyRecord;
    style?: DiyRecord;
  };
}>();

const { PrimaryButton } = useDiyAdapterComponents();
const editor = useDiyEditor();
const styleType = ref<'content' | 'style'>('content');
const productPickerOpen = ref(false);
const productTarget = ref<number | null>(null);

const defaultTitles = ['全球美妆', '大牌鞋包', '数码产品', '精品腕表'];

function defaultItem(index: number): DiscountItem {
  return {
    enabled: true,
    image: '',
    price: '',
    productId: 0,
    productName: '',
    title: defaultTitles[index] ?? `折扣分类${index + 1}`,
  };
}

function ensureDiscountDefaults(item: DiyRecord) {
  const params = (item.params ??= {}) as DiyRecord;
  const style = (item.style ??= {}) as DiyRecord;

  if (!params.title) params.title = '心动购物季';
  if (!params.promotion) params.promotion = '券后低至7.3折';
  if (!params.slogan) params.slogan = '真低价 放心买';
  if (params.iconImage === undefined) params.iconImage = '';

  const source = Array.isArray(params.items) ? params.items : [];
  params.items = Array.from({ length: 4 }, (_, index) => ({
    ...defaultItem(index),
    ...(source[index] ?? {}),
  }));

  if (style.float === undefined) style.float = 0;
  if (!style.background) style.background = 'rgba(255, 255, 255, 0)';
  if (!style.bgcolor) style.bgcolor = '#f5f5f5';
  if (style.paddingTop === undefined) style.paddingTop = 9;
  if (style.paddingBottom === undefined) style.paddingBottom = 0;
  if (style.paddingLeft === undefined) style.paddingLeft = 10;
  if (style.paddingRight === undefined) style.paddingRight = style.paddingLeft;
  if (style.marginTop === undefined) style.marginTop = 0;
  if (!style.radiusMode) style.radiusMode = 'all';
  if (style.cardRadius === undefined) style.cardRadius = 0;
  if (style.topLeftRadio === undefined) style.topLeftRadio = 0;
  if (style.topRightRadio === undefined) style.topRightRadio = 0;
  if (style.bottomLeftRadio === undefined) style.bottomLeftRadio = 0;
  if (style.bottomRightRadio === undefined) style.bottomRightRadio = 0;
  if (!style.cardShadow) style.cardShadow = 'off';
}

ensureDiscountDefaults(props.curItem);

const params = computed(() => {
  ensureDiscountDefaults(props.curItem);
  return props.curItem.params as DiyRecord;
});

const items = computed<DiscountItem[]>(() => params.value.items as DiscountItem[]);

const contentSchema = computed((): VbenFormSchema[] => [
  diySection(
    '内容设置',
    '顶部文案、图标与四个折扣分类可独立维护；商品只能通过选择器关联。',
  ),
]);

const styleSchema = computed((): VbenFormSchema[] => {
  const style = props.curItem.style as DiyRecord;
  return [
    diySection('卡片样式'),
    diySlider('style.float', '组件上浮：', { max: 48, min: 0 }),
    diyColor('style.background', '底部背景：', 'rgba(255, 255, 255, 0)', '透明'),
    diyColor('style.bgcolor', '组件背景：', '#f5f5f5', '透明'),
    diySlider('style.paddingTop', '上边距：', { max: 48, min: 0 }),
    diySlider('style.paddingBottom', '下边距：', { max: 48, min: 0 }),
    diySlider('style.paddingLeft', '左右边距：', { max: 48, min: 0 }),
    diySlider('style.marginTop', '页面上间距：', { max: 48, min: 0 }),
    diyRadioGroup(
      'style.radiusMode',
      '背景圆角：',
      [
        { label: '全部', value: 'all' },
        { label: '分别设置', value: 'individual' },
      ],
      true,
    ),
    ...(style.radiusMode === 'individual'
      ? [
          diySlider('style.topLeftRadio', '左上圆角：', { max: 32, min: 0 }),
          diySlider('style.topRightRadio', '右上圆角：', { max: 32, min: 0 }),
          diySlider('style.bottomLeftRadio', '左下圆角：', { max: 32, min: 0 }),
          diySlider('style.bottomRightRadio', '右下圆角：', { max: 32, min: 0 }),
        ]
      : [diySlider('style.cardRadius', '圆角值：', { max: 32, min: 0 })]),
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

const { Form: ContentForm } = useDiyCurItemForm(
  () => props.curItem,
  contentSchema,
  { onInit: ensureDiscountDefaults },
);
const { Form: StyleForm } = useDiyCurItemForm(
  () => props.curItem,
  styleSchema,
  { onInit: ensureDiscountDefaults },
);

function selectHeaderIcon() {
  editor.onEditorSelectImage(params.value, 'iconImage');
}

function openProductPicker(index: number) {
  productTarget.value = index;
  productPickerOpen.value = true;
}

function onProductPicked(product: PlatformProduct) {
  const index = productTarget.value;
  if (index === null || !items.value[index]) return;
  Object.assign(items.value[index], {
    image: String(product.image ?? ''),
    price: String(product.price ?? ''),
    productId: Number(product.product_id ?? 0),
    productName: String(product.store_name ?? '未命名商品'),
  });
  productTarget.value = null;
}

function clearProduct(index: number) {
  const item = items.value[index];
  if (!item) return;
  Object.assign(item, {
    image: '',
    price: '',
    productId: 0,
    productName: '',
  });
}
</script>

<template>
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

  <template v-if="styleType === 'content'">
    <ContentForm />
    <section class="discount-header-card">
      <h4>顶部展示</h4>
      <div class="discount-header-grid">
        <label class="discount-input-field">
          <span>主标题</span>
          <DiyInputField v-model="params.title" :maxlength="14" placeholder="如：心动购物季" />
        </label>
        <label class="discount-input-field">
          <span>优惠文案</span>
          <DiyInputField v-model="params.promotion" :maxlength="18" placeholder="如：券后低至7.3折" />
        </label>
        <label class="discount-input-field discount-input-field--full">
          <span>卖点文案</span>
          <DiyInputField v-model="params.slogan" :maxlength="16" placeholder="如：真低价 放心买" />
        </label>
      </div>
      <button class="discount-icon-picker" type="button" @click="selectHeaderIcon">
        <img v-if="params.iconImage" v-img-url="params.iconImage" alt="折扣组图标" />
        <ElIcon v-else><Picture /></ElIcon>
        <span>{{ params.iconImage ? '更换图标' : '选择图标' }}</span>
      </button>
    </section>

    <section class="discount-module-list">
      <div class="discount-list-heading">
        <div>
          <h4>折扣分类</h4>
          <p>每个商品位独立配置，关联商品后自动带入商品图与价格。</p>
        </div>
      </div>
      <article v-for="(item, index) in items" :key="index" class="discount-module-card">
        <div class="discount-module-card__header">
          <h4>分类 {{ index + 1 }}</h4>
          <ElSwitch v-model="item.enabled" />
        </div>
        <label class="discount-input-field">
          <span>分类标题</span>
          <DiyInputField v-model="item.title" :maxlength="8" placeholder="请输入分类标题" />
        </label>
        <div class="discount-product-row">
          <button class="discount-product-cover" type="button" @click="openProductPicker(index)">
            <img v-if="item.image" v-img-url="item.image" :alt="item.title" />
            <ElIcon v-else><Picture /></ElIcon>
          </button>
          <div class="discount-product-info">
            <strong>{{ item.productName || '未关联商品' }}</strong>
            <span>{{ item.productId ? `商品 ID：${item.productId}` : '请从商品库选择关联商品' }}</span>
            <span v-if="item.price">售价：¥{{ item.price }}</span>
          </div>
        </div>
        <div class="discount-actions">
          <PrimaryButton size="small" @click="openProductPicker(index)">
            {{ item.productId ? '更换商品' : '选择商品' }}
          </PrimaryButton>
          <button v-if="item.productId" class="discount-text-button" type="button" @click="clearProduct(index)">
            取消关联
          </button>
        </div>
      </article>
    </section>
  </template>
  <StyleForm v-else />
  <ProductPickerDialog v-model:open="productPickerOpen" @select="onProductPicked" />
</template>

<style scoped>
.discount-header-card,
.discount-module-list {
  padding: 0 16px 16px;
}

.discount-header-card {
  border-bottom: 8px solid var(--vben-bg-color);
}

.discount-header-card h4,
.discount-list-heading h4,
.discount-module-card h4 {
  margin: 0;
  color: var(--el-text-color-primary);
  font-size: 14px;
  font-weight: 600;
}

.discount-header-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-top: 12px;
}

.discount-input-field {
  display: grid;
  gap: 6px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.discount-input-field--full {
  grid-column: 1 / -1;
}

.discount-icon-picker {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  min-height: 54px;
  margin-top: 12px;
  overflow: hidden;
  color: var(--el-text-color-secondary);
  background: var(--el-fill-color-lighter);
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
}

.discount-icon-picker img {
  width: 36px;
  height: 36px;
  object-fit: cover;
  border-radius: 6px;
}

.discount-icon-picker .el-icon {
  font-size: 20px;
}

.discount-module-list {
  display: grid;
  gap: 12px;
  padding-top: 16px;
}

.discount-list-heading p {
  margin: 6px 0 0;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 18px;
}

.discount-module-card {
  display: grid;
  gap: 12px;
  padding: 12px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
}

.discount-module-card__header,
.discount-product-row,
.discount-actions {
  display: flex;
  align-items: center;
}

.discount-module-card__header {
  justify-content: space-between;
}

.discount-product-row {
  gap: 10px;
}

.discount-product-cover {
  display: grid;
  flex: 0 0 66px;
  width: 66px;
  height: 66px;
  padding: 0;
  overflow: hidden;
  color: var(--el-text-color-placeholder);
  background: var(--el-bg-color);
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  place-items: center;
}

.discount-product-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.discount-product-cover .el-icon {
  font-size: 20px;
}

.discount-product-info {
  display: grid;
  min-width: 0;
  gap: 3px;
}

.discount-product-info strong,
.discount-product-info span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.discount-product-info strong {
  color: var(--el-text-color-primary);
  font-size: 13px;
}

.discount-product-info span {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.discount-actions {
  justify-content: flex-end;
  gap: 10px;
}

.discount-text-button {
  padding: 0;
  color: var(--el-color-primary);
  font-size: 12px;
  background: transparent;
  border: 0;
  cursor: pointer;
}
</style>
