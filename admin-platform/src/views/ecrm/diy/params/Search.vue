<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, ref } from 'vue';
import { Close, Rank } from '@element-plus/icons-vue';
import draggable from 'vuedraggable';

import {
  diyBgColors,
  diyColor,
  diyInput,
  diyRadioGroup,
  diySection,
  diySlider,
  DIY_PADDING_FIELDS,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsSearch' });

type DiyRecord = Record<string, unknown>;

const props = defineProps<{
  curItem: DiyRecord & { params?: DiyRecord; style?: DiyRecord };
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const { PrimaryButton } = useDiyAdapterComponents();
const styleType = ref<'content' | 'style'>('content');

function ensureSearchDefaults(item: DiyRecord) {
  const params = (item.params ??= {}) as DiyRecord;
  const style = (item.style ??= {}) as DiyRecord;

  // 兼容旧版 image/text/search 配置，同时让已有页面无需手工重建组件即可使用新选项。
  if (params.title_type === 'image') params.title_type = 'location';
  if (params.title_type === 'text') params.title_type = 'title';
  if (!params.title_type) params.title_type = 'search';
  if (!params.style_type) {
    params.style_type = params.toplogo ? 'logo' : params.title_type === 'title' ? 'title' : 'location';
  }
  if (!params.display_mode) params.display_mode = 'normal';
  if (!params.title) params.title = '首页';
  if (!params.locationText) params.locationText = '金湾大厦';
  if (!params.searchText) params.searchText = '搜索商品';
  if (!Array.isArray(params.hotWords)) params.hotWords = ['雅诗兰黛', 'Only澳白瓶'];
  if (!params.hotWordInterval) params.hotWordInterval = 3;

  if (!style.searchBackGround) style.searchBackGround = '#ffffff';
  if (!style.searchColor) style.searchColor = '#cccccc';
  if (!style.hotWordColor) style.hotWordColor = '#666666';
  if (!style.titleTextColor) style.titleTextColor = '#333333';
  if (!style.titleTextSize) style.titleTextSize = 14;
  if (!style.titleAlign) style.titleAlign = 'left';
  if (!style.titleTextStyle) style.titleTextStyle = 'normal';
  if (style.float === undefined) style.float = 0;
  if (!style.background || style.background === '#fd642a') style.background = '#fff1f0';
  if (!style.bgcolor || style.bgcolor === '#fd642a') style.bgcolor = '#fff1f0';
  if (style.paddingTop === undefined) style.paddingTop = 0;
  if (style.paddingBottom === undefined) style.paddingBottom = 0;
  if (style.paddingLeft === undefined) style.paddingLeft = 0;
  if (style.paddingRight === undefined) style.paddingRight = style.paddingLeft;
  if (!style.radiusMode) style.radiusMode = 'all';
  if (style.topRadio === undefined) style.topRadio = 0;
  if (style.bottomRadio === undefined) style.bottomRadio = 0;
  if (!style.shadow) style.shadow = 'off';
}

function addHotWord() {
  hotWords.value = [...hotWords.value, `热词${hotWords.value.length + 1}`];
}

function removeHotWord(index: number) {
  if (hotWords.value.length <= 1) return;
  hotWords.value = hotWords.value.filter((_, currentIndex) => currentIndex !== index);
}

const hotWords = computed<string[]>({
  get() {
    const words = props.curItem.params?.hotWords;
    return Array.isArray(words) ? [...(words as string[])] : [];
  },
  set(words) {
    const params = (props.curItem.params ??= {}) as DiyRecord;
    params.hotWords = words;
  },
});

const hotWordInterval = computed<number>({
  get() {
    return Number(props.curItem.params?.hotWordInterval ?? 3);
  },
  set(value) {
    const params = (props.curItem.params ??= {}) as DiyRecord;
    params.hotWordInterval = Math.min(10, Math.max(1, Number(value) || 3));
  },
});

function updateHotWord(index: number, value: string) {
  const words = [...hotWords.value];
  words[index] = value;
  hotWords.value = words;
}

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('展示设置'),
  diyRadioGroup(
    'params.display_mode',
    '搜索设置：',
    [
      { label: '正常显示', value: 'normal' },
      { label: '滚动至顶部固定', value: 'fixed' },
    ],
  ),
  diyRadioGroup(
    'params.title_type',
    '选择风格：',
    [
      { label: '标题', value: 'title' },
      { label: '定位', value: 'location' },
      { label: '搜索', value: 'search' },
    ],
  ),
  diyRadioGroup(
    'params.style_type',
    '样式类型：',
    [
      { label: '标题', value: 'title' },
      { label: '定位', value: 'location' },
      { label: 'logo', value: 'logo' },
    ],
  ),
  ...(props.curItem.params?.style_type === 'title' ? [diyInput('params.title', '标题文字：')] : []),
  ...(props.curItem.params?.style_type === 'location'
    ? [diyInput('params.locationText', '定位文字：')]
    : []),
  diySection('搜索内容'),
  diyInput('params.searchText', '提示文字：'),
]);

const styleSchema = computed((): VbenFormSchema[] => {
  const style = props.curItem.style ?? {};
  return [
    diySection('搜索框'),
    diyColor('style.searchBackGround', '搜索框：', '#ffffff'),
    diyColor('style.searchColor', '提示文字：', '#cccccc'),
    diyColor('style.hotWordColor', '热词文字：', '#666666'),
    diySection('文字设置'),
    diyColor('style.titleTextColor', '文字颜色：', '#333333'),
    diySlider('style.titleTextSize', '文字大小：', { max: 36, min: 12 }),
    diyRadioGroup(
      'style.titleAlign',
      '文字位置：',
      [
        { label: '左对齐', value: 'left' },
        { label: '居中对齐', value: 'center' },
        { label: '右对齐', value: 'right' },
      ],
    ),
    diyRadioGroup(
      'style.titleTextStyle',
      '文字样式：',
      [
        { label: '正常', value: 'normal' },
        { label: '倾斜', value: 'italic' },
        { label: '加粗', value: 'bold' },
      ],
    ),
    diySection('卡片样式'),
    diySlider('style.float', '组件上浮：', { max: 48, min: 0 }),
    ...diyBgColors('rgba(255, 255, 255, 0)', 'rgba(255, 255, 255, 0)'),
    ...DIY_PADDING_FIELDS,
    diySlider('style.paddingRight', '右边距：', { max: 48, min: 0 }),
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
          diySlider('style.topLeftRadio', '左上圆角：', { max: 48, min: 0 }),
          diySlider('style.topRightRadio', '右上圆角：', { max: 48, min: 0 }),
          diySlider('style.bottomLeftRadio', '左下圆角：', { max: 48, min: 0 }),
          diySlider('style.bottomRightRadio', '右下圆角：', { max: 48, min: 0 }),
        ]
      : [diySlider('style.topRadio', '圆角值：', { max: 48, min: 0 })]),
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
  onInit: ensureSearchDefaults,
});
const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema, {
  onInit: ensureSearchDefaults,
});
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
      <div v-if="curItem.params?.style_type === 'logo'" class="form-item ml-[100px]">
        <div class="form-label mb-2">logo图：</div>
        <div class="diy-setpages-cover">
          <img
            v-img-url="curItem.params?.toplogo"
            :width="120"
            alt=""
            @click="editor.onEditorSelectImage(curItem.params!, 'toplogo')"
          />
          <div>建议尺寸78 × 64</div>
        </div>
      </div>
      <div class="form-chink"></div>
      <div class="f16 gray3 form-subtitle">搜索热词</div>
      <div class="search-hotword-panel">
        <draggable v-model="hotWords" item-key="index" handle=".hotword-grip" animation="180">
          <template #item="{ index }">
            <div class="search-hotword-row">
              <el-icon class="hotword-grip"><Rank /></el-icon>
              <el-input
                :model-value="hotWords[index]"
                maxlength="20"
                @update:model-value="updateHotWord(index, $event)"
              />
              <el-icon
                class="hotword-remove"
                :class="{ 'is-disabled': hotWords.length <= 1 }"
                @click="removeHotWord(index)"
              >
                <Close />
              </el-icon>
            </div>
          </template>
        </draggable>
        <component :is="PrimaryButton" class="search-hotword-add" plain @click="addHotWord">+ 添加</component>
        <div class="hotword-duration">
          <span>显示时间</span>
          <el-input-number v-model="hotWordInterval" :max="10" :min="1" controls-position="right" />
          <em>秒</em>
        </div>
      </div>
    </div>
    <div v-show="styleType === 'style'">
      <StyleForm />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.diy-setpages-cover > img {
  width: 60px;
  height: 60px;
}

.search-hotword-panel {
  padding: 0 20px 24px;
}

.search-hotword-row {
  display: flex;
  align-items: center;
  gap: 14px;
  margin: 0 0 14px;

  :deep(.el-input) {
    flex: 1;
  }
}

.hotword-grip,
.hotword-remove {
  flex: 0 0 auto;
  color: var(--el-text-color-placeholder);
  cursor: pointer;
  font-size: 22px;
}

.hotword-grip {
  cursor: grab;
}

.hotword-remove.is-disabled {
  cursor: not-allowed;
  opacity: 0.35;
}

.search-hotword-add {
  width: 100%;
  min-height: 42px;
  margin: 2px 0 20px;
}

.hotword-duration {
  display: flex;
  align-items: center;
  gap: 14px;
  padding-top: 20px;
  border-top: 1px dashed var(--el-border-color-lighter);
  color: var(--el-text-color-secondary);

  :deep(.el-input-number) {
    flex: 1;
  }

  em {
    font-style: normal;
  }
}
</style>
