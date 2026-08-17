<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { ArrowRight, CloseBold, Rank } from '@element-plus/icons-vue';
import { ElIcon } from 'element-plus';
import { computed, ref } from 'vue';
import draggable from 'vuedraggable';

import { getShopLinkProductCategoryApi } from '#/api/core/shop-link';
import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import {
  diyBgColors,
  diyColor,
  diyRadioGroup,
  diySection,
  diySlider,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyInputField from './shared/diy-input-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';

defineOptions({ name: 'DiyParamsOption' });

type DiyRecord = Record<string, unknown>;
type OptionDataType = 'category' | 'page';

const { PrimaryButton, RadioGroup } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: DiyRecord & {
    data?: DiyRecord[];
    params?: DiyRecord;
    style?: DiyRecord;
  };
}>();

const styleType = ref<'content' | 'style'>('content');
const stylePickerVisible = ref(false);
const categoryLoading = ref(false);
const categoryList = ref<DiyRecord[]>([]);
const isLinkset = ref(false);
const linkIndex = ref(0);
const linkData = ref<DiyRecord | null>(null);

const tabs = computed<DiyRecord[]>({
  get: () => props.curItem.data ?? [],
  set: (value) => {
    props.curItem.data = value;
  },
});

let nextTabId = 0;

function createTab(overrides: DiyRecord = {}): DiyRecord {
  nextTabId += 1;
  return {
    _diyTabId: `option-tab-${Date.now()}-${nextTabId}`,
    dataType: 'page',
    linkName: '微页面',
    linkUrl: '/pages/small_page/index',
    text: '选项卡',
    ...overrides,
  };
}

function ensureOptionDefaults(item: DiyRecord) {
  const params = (item.params ??= {}) as DiyRecord;
  const style = (item.style ??= {}) as DiyRecord;

  if (!params.type) params.type = '2';
  if (params.topUp === undefined) params.topUp = '0';
  if (!style.themeType) style.themeType = 'theme';
  if (!style.activeColor) style.activeColor = style.active_color1 ?? '#ff4d7d';
  if (!style.activeText) style.activeText = '#ffffff';
  if (!style.background) style.background = 'rgba(255, 255, 255, 0)';
  if (!style.bgcolor) style.bgcolor = '#fff0f3';
  if (style.float === undefined) style.float = 0;
  if (style.marginTop === undefined) style.marginTop = 0;
  if (style.paddingTop === undefined) style.paddingTop = 0;
  if (style.paddingBottom === undefined) style.paddingBottom = 0;
  if (style.paddingLeft === undefined) style.paddingLeft = 0;
  if (style.paddingRight === undefined) style.paddingRight = style.paddingLeft;
  if (!style.radiusMode) style.radiusMode = 'all';
  if (style.radius === undefined) style.radius = style.topRadio ?? 0;
  if (style.topLeftRadius === undefined) style.topLeftRadius = style.radius;
  if (style.topRightRadius === undefined) style.topRightRadius = style.radius;
  if (style.bottomRightRadius === undefined) style.bottomRightRadius = style.bottomRadio ?? style.radius;
  if (style.bottomLeftRadius === undefined) style.bottomLeftRadius = style.bottomRadio ?? style.radius;
  if (!style.shadow) style.shadow = 'off';

  if (!Array.isArray(item.data) || item.data.length === 0) {
    item.data = [
      createTab({ linkName: '首页', linkUrl: '/pages/index/index', text: '首页' }),
      createTab({ text: '果蔬生鲜' }),
      createTab({ text: '健康医疗' }),
      createTab({ text: '非遗文创' }),
      createTab({ text: '优选茶叶' }),
    ];
  }

  item.data.forEach((tab) => {
    if (!tab._diyTabId) tab._diyTabId = createTab()._diyTabId;
    if (!tab.dataType) tab.dataType = tab.category_id ? 'category' : 'page';
    if (!tab.text) tab.text = '选项卡';
    if (tab.dataType === 'page' && !tab.linkName) tab.linkName = tab.name ?? '选择微页面';
  });
}

async function loadCategories() {
  if (categoryList.value.length > 0 || categoryLoading.value) return;
  categoryLoading.value = true;
  try {
    const response = await getShopLinkProductCategoryApi();
    categoryList.value = response.list ?? [];
  } finally {
    categoryLoading.value = false;
  }
}

function selectedStyleLabel() {
  const type = String(props.curItem.params?.type ?? '2');
  return type === '1' ? '样式一' : type === '3' ? '样式三' : '样式二';
}

function addTab() {
  tabs.value.push(createTab());
}

function removeTab(index: number) {
  if (tabs.value.length <= 1) return;
  tabs.value.splice(index, 1);
}

function changeLink(index: number) {
  isLinkset.value = true;
  linkIndex.value = index;
  linkData.value = tabs.value[index] ?? null;
}

function closeLinkset(value: { name?: string; type?: string; url?: string } | null) {
  isLinkset.value = false;
  const tab = tabs.value[linkIndex.value];
  if (!value || !tab) return;
  tab.linkeType = value.type;
  tab.linkName = value.name;
  tab.linkUrl = value.url;
}

function findCategory(nodes: DiyRecord[], target: unknown): DiyRecord | undefined {
  for (const node of nodes) {
    if (String(node.category_id) === String(target)) return node;
    const nested = Array.isArray(node.child) ? (node.child as DiyRecord[]) : [];
    const match = findCategory(nested, target);
    if (match) return match;
  }
}

function changeCategory(value: unknown, index: number) {
  const tab = tabs.value[index];
  if (!tab) return;
  const category = findCategory(categoryList.value, value);
  tab.category_id = value as string | number;
  tab.categoryName = category?.name ?? '';
}

function changeDataType(tab: DiyRecord, value: OptionDataType) {
  tab.dataType = value;
  if (value === 'page') {
    tab.category_id = undefined;
    tab.categoryName = undefined;
  } else {
    tab.linkUrl = undefined;
    tab.linkName = undefined;
  }
}

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('展示设置'),
  diyRadioGroup(
    'params.topUp',
    '滑动置顶：',
    [
      { label: '启用', value: '1' },
      { label: '不启用', value: '0' },
    ],
  ),
]);

const styleSchema = computed((): VbenFormSchema[] => {
  const style = props.curItem.style ?? {};
  return [
    diySection('选项卡样式'),
    diyRadioGroup(
      'style.themeType',
      '色调：',
      [
        { label: '跟随主题风格', value: 'theme' },
        { label: '自定义', value: 'custom' },
      ],
    ),
    ...(style.themeType === 'custom'
      ? [
          diyColor('style.activeColor', '选中颜色：', '#ff4d7d'),
          diyColor('style.activeText', '选中文字：', '#ffffff'),
        ]
      : []),
    diySection('卡片样式'),
    diySlider('style.float', '组件上浮：', { max: 48, min: 0 }),
    ...diyBgColors('rgba(255, 255, 255, 0)', '#fff0f3'),
    diySlider('style.paddingTop', '上边距：', { max: 48, min: 0 }),
    diySlider('style.paddingBottom', '下边距：', { max: 48, min: 0 }),
    diySlider('style.paddingLeft', '左右边距：', { max: 48, min: 0 }),
    diySlider('style.marginTop', '页面上间距：', { max: 96, min: 0 }),
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
          diySlider('style.topLeftRadius', '左上圆角：', { max: 48, min: 0 }),
          diySlider('style.topRightRadius', '右上圆角：', { max: 48, min: 0 }),
          diySlider('style.bottomRightRadius', '右下圆角：', { max: 48, min: 0 }),
          diySlider('style.bottomLeftRadius', '左下圆角：', { max: 48, min: 0 }),
        ]
      : [diySlider('style.radius', '圆角值：', { max: 48, min: 0 })]),
    diyRadioGroup(
      'style.shadow',
      '开启阴影：',
      [
        { label: '关闭', value: 'off' },
        { label: '开启', value: 'on' },
      ],
    ),
  ];
});

const { Form: ContentForm } = useDiyCurItemForm(() => props.curItem, contentSchema, {
  onInit: ensureOptionDefaults,
});
const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema, {
  onInit: ensureOptionDefaults,
});

void loadCategories();
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
      <div class="form-vben-item option-style-picker" label="选择风格：">
        <component :is="PrimaryButton" plain size="small" @click="stylePickerVisible = !stylePickerVisible">
          修改风格
        </component>
        <span class="option-style-picker__current">当前：{{ selectedStyleLabel() }}</span>
      </div>
      <div v-if="stylePickerVisible" class="form-vben-item" label="风格样式：">
        <component :is="RadioGroup" v-model="curItem.params!.type">
          <el-radio label="1">样式一</el-radio>
          <el-radio label="2">样式二</el-radio>
          <el-radio label="3">样式三</el-radio>
        </component>
      </div>

      <div class="form-chink"></div>
      <div class="f16 gray3 form-subtitle">
        选项卡设置
        <span class="gray f12">鼠标拖拽版块可调整选项卡顺序</span>
      </div>
      <draggable
        v-model="tabs"
        class="draggable-list"
        item-key="_diyTabId"
        handle=".option-tab-item__handle"
        :animation="180"
      >
        <template #item="{ element, index }">
          <div class="d-c-c param-img-item option-tab-item">
            <ElIcon class="option-tab-item__handle" size="20">
              <Rank />
            </ElIcon>
            <div class="right pr">
              <ElIcon
                class="el-icon-DeleteFilled"
                :class="{ 'is-disabled': tabs.length <= 1 }"
                :title="tabs.length <= 1 ? '至少保留一个选项卡' : '删除选项卡'"
                @click.stop="removeTab(index)"
              >
                <CloseBold />
              </ElIcon>
              <div class="url-box mb16 flex-1 d-s-c ww100">
                <span class="key-name">显示文字</span>
                <DiyInputField v-model="element.text" maxlength="6" show-word-limit />
              </div>
              <div class="url-box mb16 flex-1 d-s-c ww100">
                <span class="key-name">数据类型</span>
                <component
                  :is="RadioGroup"
                  :model-value="element.dataType"
                  @update:model-value="changeDataType(element, $event as OptionDataType)"
                >
                  <el-radio label="page">微页面</el-radio>
                  <el-radio label="category">商品分类</el-radio>
                </component>
              </div>
              <div v-if="element.dataType === 'page'" class="url-box flex-1 d-s-c ww100">
                <span class="key-name">微页面</span>
                <DiyInputField :model-value="String(element.linkName ?? element.linkUrl ?? '')" readonly>
                  <template #suffix>
                    <ElIcon color="#98a2b3" size="16px" @click.stop="changeLink(index)">
                      <ArrowRight />
                    </ElIcon>
                  </template>
                </DiyInputField>
              </div>
              <div v-else class="url-box flex-1 d-s-c ww100">
                <span class="key-name">商品分类</span>
                <el-cascader
                  v-model="element.category_id"
                  class="ww100"
                  clearable
                  :loading="categoryLoading"
                  :options="categoryList"
                  :props="{
                    checkStrictly: true,
                    children: 'child',
                    emitPath: false,
                    label: 'name',
                    value: 'category_id',
                  }"
                  @change="changeCategory($event, index)"
                />
              </div>
            </div>
          </div>
        </template>
      </draggable>
      <div class="d-c-c pb16">
        <component :is="PrimaryButton" plain @click="addTab">+ 添加选项卡</component>
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
      选择微页面
    </DiyLinkPickerDialog>
  </div>
</template>

<style lang="scss" scoped>
.option-style-picker {
  align-items: center;
  display: flex;
  gap: 12px;
  margin: 0 0 16px;
}

.option-style-picker__current {
  color: #98a2b3;
  font-size: 13px;
}

.option-tab-item {
  align-items: flex-start;
  background: #f8fafc;
  border: 1px solid #edf0f5;
  border-radius: 10px;
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
  padding: 14px 12px;

  .right {
    flex: 1;
    min-width: 0;
  }

  .url-box {
    margin-bottom: 12px;
  }

  .url-box:last-child {
    margin-bottom: 0;
  }
}

.option-tab-item__handle {
  color: #c7cfdd;
  cursor: grab;
  flex: none;
  margin-top: 11px;

  &:active {
    cursor: grabbing;
  }
}

.el-icon-DeleteFilled.is-disabled {
  cursor: not-allowed;
  opacity: 0.4;
  pointer-events: none;
}
</style>
