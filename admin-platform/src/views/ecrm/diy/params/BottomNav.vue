<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, ref } from 'vue';

import { ArrowRight, CloseBold, Rank } from '@element-plus/icons-vue';
import { ElIcon, ElMessage, ElSwitch } from 'element-plus';
import draggable from 'vuedraggable';

import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import DiyInputField from './shared/diy-input-field.vue';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import { diyColor, diyRadioGroup, diySection, diySlider } from './shared/schema-helpers';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsBottomNav' });
type DiyRecord = Record<string, any>;

const props = defineProps<{ curItem: DiyRecord & { data?: DiyRecord[]; params?: DiyRecord; style?: DiyRecord }; selectedIndex?: number }>();
const { PrimaryButton } = useDiyAdapterComponents();
const editor = useDiyEditor();
const tab = ref<'content' | 'style'>('content');
const isLinkset = ref(false);
const linkIndex = ref(0);
const linkData = ref<DiyRecord | null>(null);

const navigationItems = computed<DiyRecord[]>({ get: () => props.curItem.data ?? [], set: (value) => { props.curItem.data = value; } });

function ensureDefaults(item: DiyRecord) {
  const style = (item.style ??= {}) as DiyRecord;
  const params = (item.params ??= {}) as DiyRecord;
  if (!Array.isArray(item.data)) item.data = [];
  if (!style.navigationType) style.navigationType = 'icon-text';
  if (!style.positionType) style.positionType = 'fixed';
  if (!style.themeMode) style.themeMode = 'system';
  if (!style.background) style.background = 'rgba(255,255,255,0.96)';
  if (!style.activeColor) style.activeColor = '#f62c2c';
  if (!style.textColor) style.textColor = '#282828';
  for (const key of ['paddingTop', 'paddingBottom', 'pagePadding', 'bottomSpacing', 'radius']) if (style[key] === undefined) style[key] = 0;
  if (params.activeIndex === undefined) params.activeIndex = 0;
  item.data.forEach((entry: DiyRecord) => { if (!entry.text) entry.text = '导航'; if (entry.hide === undefined) entry.hide = false; });
}

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('展示设置'),
  diyRadioGroup('style.positionType', '导航类型：', [{ label: '底部固定', value: 'fixed' }, { label: '底部悬浮', value: 'float' }]),
  diyRadioGroup('style.navigationType', '导航样式：', [{ label: '图片+文字', value: 'icon-text' }, { label: '文字', value: 'text' }, { label: '图片', value: 'icon' }]),
  diySection('导航内容', '拖拽左侧圆点可调整导航顺序；最多 5 项'),
]);
const styleSchema = computed((): VbenFormSchema[] => [
  diySection('颜色设置'),
  diyRadioGroup('style.themeMode', '色调：', [{ label: '跟随主题风格', value: 'system' }, { label: '自定义', value: 'custom' }]),
  ...(props.curItem.style?.themeMode === 'custom' ? [diyColor('style.activeColor', '选中文字颜色：', '#f62c2c'), diyColor('style.textColor', '文字颜色：', '#282828')] : []),
  diySection('卡片样式'),
  diyColor('style.background', '背景颜色：', 'rgba(255,255,255,0.96)'),
  diySlider('style.paddingTop', '上边距：', { min: 0, max: 24 }),
  diySlider('style.paddingBottom', '下边距：', { min: 0, max: 24 }),
  diySlider('style.pagePadding', '左右边距：', { min: 0, max: 32 }),
  diySlider('style.bottomSpacing', '页面下间距：', { min: 0, max: 32 }),
  ...(props.curItem.style?.positionType === 'float' ? [diySlider('style.radius', '圆角值：', { min: 0, max: 36 })] : []),
]);

const { Form: ContentForm } = useDiyCurItemForm(() => props.curItem, contentSchema, { onInit: ensureDefaults });
const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema, { onInit: ensureDefaults });

function changeLink(index: number) { isLinkset.value = true; linkIndex.value = index; linkData.value = navigationItems.value[index] ?? null; }
function closeLinkset(value: { name?: string; type?: string; url?: string } | null) { isLinkset.value = false; const entry = navigationItems.value[linkIndex.value]; if (value && entry) { entry.linkeType = value.type; entry.linkUrl = value.url; entry.name = value.name; } }
function addItem() { if (navigationItems.value.length >= 5) { ElMessage.warning('底部导航最多添加 5 项'); return; } navigationItems.value.push({ text: '导航', linkUrl: '', selectedImgUrl: '', unselectedImgUrl: '', hide: false }); }
function removeItem(index: number) { if (navigationItems.value.length <= 2) { ElMessage.warning('底部导航至少保留 2 项'); return; } navigationItems.value.splice(index, 1); }
</script>

<template>
  <div class="bottom-nav-params">
    <div class="common-form common-form-new"><span>{{ curItem.name }}</span><div class="diy-changes"><div class="diy-change" :class="{ active: tab === 'content' }" @click="tab = 'content'">内容</div><div class="diy-change" :class="{ active: tab === 'style' }" @click="tab = 'style'">样式</div></div></div>
    <div v-show="tab === 'content'"><ContentForm />
      <draggable v-model="navigationItems" class="draggable-list" item-key="text"><template #item="{ element, index }"><div class="d-c-c param-img-item navbar nav-item-card"><div class="nav-item-card__drag"><ElIcon><Rank /></ElIcon></div><div class="nav-images"><div class="icon"><img v-img-url="element.selectedImgUrl" alt="" @click="editor.onEditorSelectImage(element, 'selectedImgUrl')" /><small>选中</small></div><div class="icon"><img v-img-url="element.unselectedImgUrl" alt="" @click="editor.onEditorSelectImage(element, 'unselectedImgUrl')" /><small>未选中</small></div></div><div class="right nav-item-card__fields"><ElIcon class="el-icon-DeleteFilled" @click.stop="removeItem(index)"><CloseBold /></ElIcon><div class="url-box mb16 flex-1 d-s-c ww100"><span class="key-name">名称</span><DiyInputField v-model="element.text" :maxlength="6" show-word-limit /></div><div class="url-box mb16 flex-1 d-s-c ww100"><span class="key-name">链接</span><DiyLinkInputField v-model="element.linkUrl" @click="changeLink(index)"><template #suffix><ElIcon color="#333" size="16px"><ArrowRight /></ElIcon></template></DiyLinkInputField></div><div class="url-box flex-1 d-s-c ww100"><span class="key-name">状态</span><ElSwitch :model-value="!element.hide" @update:model-value="element.hide = !$event" /></div></div></div></template></draggable>
      <div class="d-c-c pb16"><component :is="PrimaryButton" plain @click="addItem">+ 添加导航</component></div>
    </div>
    <div v-show="tab === 'style'"><StyleForm /></div>
    <DiyLinkPickerDialog v-if="isLinkset" :is_linkset="isLinkset" :link-data="linkData" @close-dialog="closeLinkset">选择链接</DiyLinkPickerDialog>
  </div>
</template>

<style scoped>
.nav-images { display: flex; gap: 8px; margin-right: 12px; }.nav-images .icon { width: 48px; text-align: center; }.nav-images img { width: 42px; height: 42px; object-fit: contain; cursor: pointer; }.nav-images small { display: block; color: #86909c; font-size: 11px; }.nav-item-card__drag { margin-right: 10px; color: #86909c; cursor: grab; }
</style>
