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
  DIY_RADIUS_FIELDS,
} from './shared/schema-helpers';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyInputField from './shared/diy-input-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsBanner' });

const { PrimaryButton } = useDiyAdapterComponents();

const props = defineProps<{
  curItem: Record<string, unknown> & {
    data?: Array<Record<string, unknown>>;
  };
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const isLinkset = ref(false);
const linkIndex = ref(0);
const linkData = ref<Record<string, unknown> | null>(null);

const schema = computed((): VbenFormSchema[] => [
  diySection('颜色设置'),
  diyColor('style.background', '底部背景：', '#ff4c01'),
  diySection('指示器设置'),
  diyColor('style.btnColor', '指示点颜色：', '#ffffff'),
  diyRadioGroup(
    'style.imgShape',
    '指示点形状：',
    [
      { label: '圆形', value: 'round' },
      { label: '正方形', value: 'square' },
      { label: '长方形', value: 'rectangle' },
    ],
    true,
  ),
  diySection('边距设置'),
  ...DIY_PADDING_FIELDS,
  diySection('圆角设置'),
  ...DIY_RADIUS_FIELDS,
  diySection('图片设置', '建议上传尺寸相同的图片，建议尺寸750px*340px'),
  diySlider('style.height', '图片高度：', { max: 1900, min: 100 }),
]);

const { Form } = useDiyCurItemForm(() => props.curItem, schema);

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
                  alt=""
                  :style="{ height: Number(curItem.style?.height ?? 0) * 0.5 + 'px' }"
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
      <component :is="PrimaryButton" plain @click="editor.onEditorAddData()">+新增一个</component>
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
  height: 202px;
  background: #eeeeee;
  margin-top: 0;
  margin-bottom: 10px;
}
</style>
