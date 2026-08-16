<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, ref, watch } from 'vue';

import { ElButton, ElInputNumber } from 'element-plus';

import ProductPickerDialog from '#/components/shop/product-picker-dialog.vue';
import type { PlatformProduct } from '#/api/core/platform-catalog';

import DiyInputField from './shared/diy-input-field.vue';
import { diyColor, diyRadioGroup, diySection, diySlider } from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyEditor } from './shared/use-diy-editor';

type DiyRecord = Record<string, any>;
const props = defineProps<{ curItem: DiyRecord & { params?: DiyRecord; style?: DiyRecord; data?: DiyRecord[] } }>();
const tab = ref<'content' | 'style'>('content');
const productPickerOpen = ref(false);
const target = ref<{ board: number; product: number } | null>(null);
const editor = useDiyEditor();

function product() { return { name: '幸运美物', image: '', price: '350.00' }; }
function board(index: number) { return { icon: '🔥', title: ['销量榜', '好评榜', '人气榜'][index] || `榜单${index + 1}`, products: [product(), product(), product()] }; }
function ensure(item: DiyRecord) {
  const params = (item.params ??= {}); const style = (item.style ??= {});
  Object.assign(params, { title: '排行榜', more: '更多', titleType: 'text', titleImage: '', boardCount: 2, productCount: 3, ...params });
  Object.assign(style, { background: '#f5f5f5', cardBackground: '#ffffff', boardBackground: '#fceae9', boardTitleColor: '#ff4c8d', moreColor: '#999999', priceColor: '#ff4c8d', paddingTop: 10, paddingBottom: 10, paddingLeft: 10, marginTop: 10, radius: 10, boardRadius: 8, shadow: 'off', ...style });
  if (!Array.isArray(item.data)) item.data = [];
  syncBoards(item);
}
function syncBoards(item = props.curItem) {
  const count = Math.max(1, Math.min(3, Number(item.params?.boardCount || 2)));
  const productCount = Math.max(1, Math.min(5, Number(item.params?.productCount || 3)));
  item.params.boardCount = count; item.params.productCount = productCount;
  while (item.data!.length < count) item.data!.push(board(item.data!.length));
  item.data!.splice(count);
  item.data!.forEach((entry, index) => { if (!entry.title) entry.title = board(index).title; if (!Array.isArray(entry.products)) entry.products = []; while (entry.products.length < productCount) entry.products.push(product()); entry.products.splice(productCount); });
}
ensure(props.curItem);
const params = computed(() => { ensure(props.curItem); return props.curItem.params!; });
const boards = computed(() => props.curItem.data ?? []);
watch(() => [params.value.boardCount, params.value.productCount], () => syncBoards());

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('头部设置'), diyRadioGroup('params.titleType', '标题类型：', [{ label: '文字', value: 'text' }, { label: '图片', value: 'image' }]),
  diySection('展示设置'), diySlider('params.boardCount', '榜单数量：', { min: 1, max: 3 }), diySlider('params.productCount', '每榜商品数：', { min: 1, max: 5 }),
]);
const styleSchema = computed((): VbenFormSchema[] => [
  diySection('头部样式'), diyColor('style.moreColor', '按钮颜色：', '#999999'),
  diySection('榜单样式'), diyColor('style.boardBackground', '榜单背景：', '#fceae9'), diyColor('style.boardTitleColor', '榜单标题：', '#ff4c8d'), diyColor('style.priceColor', '商品价格：', '#ff4c8d'), diySlider('style.boardRadius', '榜单圆角：', { min: 0, max: 24 }),
  diySection('卡片样式'), diyColor('style.background', '底部背景：', '#f5f5f5'), diyColor('style.cardBackground', '组件背景：', '#ffffff'), diySlider('style.paddingTop', '上边距：', { min: 0, max: 48 }), diySlider('style.paddingBottom', '下边距：', { min: 0, max: 48 }), diySlider('style.paddingLeft', '左右边距：', { min: 0, max: 48 }), diySlider('style.marginTop', '页面上间距：', { min: 0, max: 48 }), diySlider('style.radius', '背景圆角：', { min: 0, max: 32 }), diyRadioGroup('style.shadow', '开启阴影：', [{ label: '关闭', value: 'off' }, { label: '开启', value: 'on' }]),
]);
const { Form: ContentForm } = useDiyCurItemForm(() => props.curItem, contentSchema, { onInit: ensure });
const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema, { onInit: ensure });
function openPicker(boardIndex: number, productIndex: number) { target.value = { board: boardIndex, product: productIndex }; productPickerOpen.value = true; }
function pick(value: PlatformProduct) { const point = target.value; if (!point) return; const item = boards.value[point.board]?.products?.[point.product]; if (item) Object.assign(item, { name: value.store_name || '未命名商品', image: value.image || '', price: String(value.price || '') }); target.value = null; }
</script>

<template>
  <div class="ranking-params">
    <div class="common-form common-form-new"><span>{{ curItem.name }}</span><div class="diy-changes"><div class="diy-change" :class="{ active: tab === 'content' }" @click="tab = 'content'">内容</div><div class="diy-change" :class="{ active: tab === 'style' }" @click="tab = 'style'">样式</div></div></div>
    <div v-show="tab === 'content'"><ContentForm />
      <section class="ranking-param-section"><h4>头部内容</h4><label v-if="params.titleType === 'text'"><span>标题文字</span><DiyInputField v-model="params.title" :maxlength="8" /></label><div v-else class="ranking-image"><img v-img-url="params.titleImage" alt="标题图片" @click="editor.onEditorSelectImage(params, 'titleImage')" /><ElButton link type="primary" @click="params.titleImage = ''">重置</ElButton><small>建议尺寸 102px × 32px</small></div><label><span>右侧按钮</span><DiyInputField v-model="params.more" :maxlength="6" /></label></section>
      <section v-for="(entry, boardIndex) in boards" :key="boardIndex" class="ranking-param-section"><h4>榜单 {{ boardIndex + 1 }}</h4><label><span>榜单名称</span><DiyInputField v-model="entry.title" :maxlength="8" /></label><div v-for="(item, productIndex) in entry.products" :key="productIndex" class="ranking-product-editor"><span>商品 {{ productIndex + 1 }}</span><DiyInputField v-model="item.name" :maxlength="16" /><ElInputNumber v-model="item.price" :min="0" :precision="2" controls-position="right" /><ElButton plain size="small" @click="openPicker(boardIndex, productIndex)">选择商品</ElButton></div></section>
    </div>
    <div v-show="tab === 'style'"><StyleForm /></div>
    <ProductPickerDialog v-model:open="productPickerOpen" @select="pick" />
  </div>
</template>

<style scoped>
.ranking-param-section { margin: 12px 0; padding: 12px; border-radius: 8px; background: #f7f8fa; }.ranking-param-section h4 { margin: 0 0 10px; color: #303133; font-size: 14px; }.ranking-param-section label { display: flex; align-items: center; gap: 10px; margin-top: 8px; }.ranking-param-section label > span { width: 68px; flex: none; color: #606266; font-size: 13px; }.ranking-image { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; }.ranking-image img { width: 102px; height: 32px; cursor: pointer; object-fit: contain; background: #fff; }.ranking-image small { width: 100%; color: #909399; }.ranking-product-editor { display: grid; grid-template-columns: 44px minmax(0, 1fr); gap: 7px; align-items: center; margin-top: 8px; }.ranking-product-editor > span { color: #606266; font-size: 12px; }.ranking-product-editor :deep(.el-input-number), .ranking-product-editor :deep(.el-button) { grid-column: 2; width: 100%; }
</style>
