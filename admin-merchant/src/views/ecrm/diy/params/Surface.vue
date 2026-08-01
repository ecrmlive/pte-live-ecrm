<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { ArrowRight } from '@element-plus/icons-vue';
import { ElIcon } from 'element-plus';
import { computed, ref } from 'vue';

import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import {
  diyRadioGroup,
  diySection,
  diySlider,
} from './shared/schema-helpers';
import { parseIntFields } from './shared/path-utils';
import { useDiyCurItemForm } from './shared/use-diy-cur-item-form';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import { useDiyEditor } from './shared/use-diy-editor';

defineOptions({ name: 'DiyParamsSurface' });

const props = defineProps<{
  curItem: Record<string, unknown> & {
    params?: Record<string, unknown> & {
      link?: Record<string, unknown>;
    };
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();
const styleType = ref<'content' | 'style'>('content');
const isLinkset = ref(false);
const linkData = ref<Record<string, unknown>>({});

const contentSchema = computed((): VbenFormSchema[] => [
  diySection('展示设置'),
  diyRadioGroup('params.type', '按钮类型：', [
    { label: '返回顶部', value: 1 },
    { label: '页面链接', value: 2 },
  ]),
  diyRadioGroup('params.showType', '显示方式：', [
    { label: '常驻', value: 1 },
    { label: '滚动后', value: 2 },
  ]),
]);

const styleSchema = computed((): VbenFormSchema[] => [
  diySection('样式设置'),
  diySlider('style.bottom', '底边距：'),
  diySlider('style.right', '右边距：'),
  diySlider('style.opacity', '不透明度：'),
]);

const { Form: ContentForm } = useDiyCurItemForm(() => props.curItem, contentSchema, {
  onValuesChange(values, fieldsChanged) {
    if (fieldsChanged.includes('params.type') && values['params.type'] === 2) {
      props.curItem.params = props.curItem.params ?? {};
      (props.curItem.params as Record<string, unknown>).link = {
        linkeType: '',
        linkUrl: '',
        name: '',
      };
    }
  },
});

const { Form: StyleForm } = useDiyCurItemForm(() => props.curItem, styleSchema, {
  onInit(item) {
    parseIntFields(item, ['style.bottom', 'style.right', 'style.opacity']);
  },
});

function changeLink() {
  isLinkset.value = true;
  linkData.value = (props.curItem.params?.link ?? {}) as Record<string, unknown>;
}

function closeLinkset(e: { name?: string; type?: string; url?: string }) {
  isLinkset.value = false;
  if (props.curItem.params) {
    props.curItem.params.link = {
      linkeType: e.type,
      linkUrl: e.url,
      name: e.name,
    };
  }
}
</script>

<template>
  <div>
    <div class="common-form common-form-new">
      <span>{{ curItem.name }}</span>
      <div class="diy-changes">
        <div
          class="diy-change"
          :class="{ active: styleType == 'content' }"
          @click="styleType = 'content'"
        >
          内容
        </div>
        <div
          class="diy-change"
          :class="{ active: styleType == 'style' }"
          @click="styleType = 'style'"
        >
          样式
        </div>
      </div>
    </div>
    <div v-if="styleType == 'content'">
      <ContentForm />
      <div
        v-if="curItem.params?.type == 2 && curItem.params?.link"
        class="form-item ml-[100px]"
      >
        <div class="form-label mb-2">链接：</div>
        <DiyLinkInputField v-model="curItem.params.link.linkUrl" @click="changeLink">
<template #suffix>
<ElIcon color="#333" size="16px"><ArrowRight /></ElIcon>
</template>
</DiyLinkInputField>
      </div>
      <div class="form-item ml-[100px]">
        <div class="form-label mb-2">悬浮图标：</div>
        <div class="diy-service-icon">
          <img
            v-img-url="curItem.params?.image"
            alt=""
            @click="editor.onEditorSelectImage(curItem.params!, 'image')"
          />
        </div>
        <div>建议尺寸120px×120px</div>
      </div>
    </div>
    <div v-if="styleType == 'style'">
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
.diy-service-icon,
.diy-service-icon img {
  width: 40px;
  height: 40px;
}
</style>
