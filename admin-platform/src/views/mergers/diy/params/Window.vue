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
  DIY_PADDING_FIELDS,
} from './shared/schema-helpers';
import { parseIntFields } from './shared/path-utils';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyColorField from './shared/diy-color-field.vue';
import DiyInputField from './shared/diy-input-field.vue';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import DiySliderField from './shared/diy-slider-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsWindow' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown> & {
    data?: Array<Record<string, unknown>>;
    style?: Record<string, unknown>;
  };
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const isLinkset = ref(false);
const linkIndex = ref(0);
const linkData = ref<Record<string, unknown> | null>(null);

const schema = computed((): VbenFormSchema[] => [
  diySection('样式设置'),
  diyColor('style.background', '底部背景：', '#F2F2F2'),
  ...DIY_PADDING_FIELDS.map((field) =>
    field.fieldName === 'style.paddingLeft'
      ? { ...field, label: '图片边距：' }
      : field,
  ),
  diySection('布局方式'),
  diyRadioGroup(
    'style.layout',
    '',
    [
      { label: '堆积两列', value: 2 },
      { label: '堆积三列', value: 3 },
      { label: '堆积四列', value: 4 },
      { label: '橱窗样式', value: -1 },
    ],
    true,
  ),
]);

const { Form } = useDiyCurItemForm(
  () => props.curItem,
  schema,
  {
    onInit(item) {
      parseIntFields(item, ['style.paddingTop', 'style.paddingLeft']);
    },
  },
);

function changeLink(index: number) {
  isLinkset.value = true;
  linkIndex.value = index;
  linkData.value = props.curItem.data?.[index] ?? null;
}

function closeLinkset(e: { name?: string; type?: string; url?: string } | null) {
  isLinkset.value = false;
  if (e && props.curItem.data?.[linkIndex.value]) {
    const row = props.curItem.data[linkIndex.value]!;
    row.linkeType = e.type;
    row.linkUrl = e.url;
    row.name = e.name;
  }
}
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <Form />
    <div v-if="curItem.style?.layout === -1" class="red px-4">橱窗样式最多显示四张图片，超出隐藏</div>
    <div class="gray px-4">请确保所有图片的尺寸/比例相同。</div>
    <div class="form-chink" />
    <div class="f16 gray3 form-subtitle mb-2 px-4">
      图片设置
      <span class="gray f12">请确保所有图片的尺寸/比例相同；鼠标拖拽左侧圆点可调整导航顺序</span>
    </div>
    <template v-if="curItem.data && curItem.data.length > 0">
      <draggable v-model="curItem.data" class="draggable-list" group="people" item-key="index">
        <template #item="{ element, index }">
          <div class="d-c-c param-img-item navbar">
            <div class="d-c d-c-c" style="margin-right: 28px">
              <div class="icon">
                <img
                  v-img-url="element.imgUrl"
                  alt=""
                  @click="editor.onEditorSelectImage(element, 'imgUrl')"
                />
              </div>
            </div>
            <div class="right">
              <ElIcon
                class="el-icon-DeleteFilled"
                @click="editor.onEditorDeleleData(index, selectedIndex ?? 0)"
              >
                <CloseBold />
              </ElIcon>
              <div class="url-box mb16 flex-1 d-s-c ww100">
                <span class="key-name">名称</span>
                <DiyInputField disabled maxlength="6" show-word-limit :value="`图${index + 1}`" />
              </div>
              <div class="d-s-c ww100">
                <div class="url-box flex-1 d-s-c">
                  <span class="key-name">链接</span>
                  <DiyLinkInputField v-model="element.linkUrl" @click="changeLink(index)">
<template #suffix>
<ElIcon color="#333" size="16px"><ArrowRight /></ElIcon>
</template>
</DiyLinkInputField>
                </div>
              </div>
            </div>
          </div>
        </template>
      </draggable>
    </template>
    <div class="d-c-c pb16">
      <component :is="PrimaryButton" plain @click="editor.onEditorAddData()">+添加一个</component>
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
