<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { computed, ref } from 'vue';

import { ArrowRight, CloseBold } from '@element-plus/icons-vue';
import { ElIcon, ElSwitch } from 'element-plus';
import draggable from 'vuedraggable';

import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import DiyInputField from './shared/diy-input-field.vue';
import { type FormRecord, getByPath, parseIntFields, setByPath } from './shared/path-utils';
import {
  DIY_PADDING_FIELDS,
  DIY_RADIUS_FIELDS,
  diyBgColors,
  diyRadioGroup,
  diySection,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyColorField from './shared/diy-color-field.vue';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import DiySliderField from './shared/diy-slider-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsNavBar' });

const { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown> & {
    data?: Array<Record<string, unknown>>;
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const isLinkset = ref(false);
const linkIndex = ref(0);
const linkData = ref<null | Record<string, unknown>>(null);

const schema = computed((): VbenFormSchema[] => [
  diySection('样式设置'),
  ...diyBgColors('#ffffff', '#f2f2f2'),
  ...DIY_PADDING_FIELDS,
  ...DIY_RADIUS_FIELDS,
  diySection('导航模式'),
  diyRadioGroup('style.rowsNum', '每行数量：', [
    { label: '3个', value: 3 },
    { label: '4个', value: 4 },
    { label: '5个', value: 5 },
  ]),
  diySection('图片设置', '图片建议宽度88*88px；鼠标拖拽左侧圆点可调整导航顺序'),
]);

function ensureNavBarStyle(item: FormRecord) {
  if (!item.style || typeof item.style !== 'object') {
    item.style = {};
  }
  parseIntFields(item, ['style.rowsNum']);
  const raw = getByPath(item, 'style.rowsNum');
  const num = Number(raw);
  if (!Number.isFinite(num) || num <= 0) {
    setByPath(item, 'style.rowsNum', 5);
  }
}

const { Form } = useDiyCurItemForm(() => props.curItem, schema, {
  fieldPaths: [
    'style.background',
    'style.bgcolor',
    'style.paddingTop',
    'style.paddingBottom',
    'style.paddingLeft',
    'style.topRadio',
    'style.bottomRadio',
    'style.rowsNum',
  ],
  onInit: ensureNavBarStyle,
});

function changeLink(index: number) {
  isLinkset.value = true;
  linkIndex.value = index;
  linkData.value = props.curItem.data?.[index] ?? null;
}

function closeLinkset(e: null | { name?: string; type?: string; url?: string }) {
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
                <DiyInputField v-model="element.text" :maxlength="6" show-word-limit />
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
              <div class="url-box mb16 flex-1 d-s-c ww100">
                <div class="form-label">隐藏：</div>
                <ElSwitch v-model="element.hide" :active-value="true" :inactive-value="false" />
              </div>
            </div>
          </div>
        </template>
      </draggable>
    </template>
    <div class="d-c-c pb16">
      <component :is="PrimaryButton" plain @click="editor.onEditorAddData()">+新增图文导航</component>
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
