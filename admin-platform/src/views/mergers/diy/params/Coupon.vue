<script lang="ts" setup>
import type { VbenFormSchema } from '#/adapter/form';

import { ArrowRight } from '@element-plus/icons-vue';
import { ElButton, ElIcon, ElInput } from 'element-plus';
import { computed, ref } from 'vue';

import DiyLinkPickerDialog from '#/components/shop/diy-link-picker-dialog.vue';

import {
  diyBgColors,
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

const props = defineProps<{
  curItem: Record<string, unknown> & {
    style?: Record<string, unknown>;
  };
  opts?: unknown;
  selectedIndex?: number;
}>();

const editor = useDiyEditor();

const schema = computed((): VbenFormSchema[] => {
  const bgtype = props.curItem.style?.bgtype;
  const fields: VbenFormSchema[] = [
    diySection('风格设置'),
    diyRadioGroup('style.type', '风格：', [{ label: '风格一', value: 1 }], true),
    diySection('优惠券数据'),
    diySlider('params.limit', '优惠券数量：', { max: 30, min: 1 }),
    diySection('按钮内容'),
    diyInput('params.btntext', '按钮文字：'),
    diySection('优惠券样式'),
    diyColor('style.descolor', '描述颜色：', '#666666', '透明'),
    diyColor('style.pricecolor', '面额颜色：', '#ff4c01', '透明'),
    diyColor('style.cillcolor', '门槛颜色：', '#ff4c01', '透明'),
    diySection('按钮样式'),
    diyColor('style.btncolor', '背景颜色：', '#ff4c01', '透明'),
    diyColor('style.btnTxtcolor', '文字颜色：', '#FFFFFF', '透明'),
    diySlider('style.btnRadio', '按钮圆角：', { max: 24, min: 0 }),
    diySection('背景设置'),
    diyRadioGroup(
      'style.bgtype',
      '组件背景：',
      [
        { label: '背景色', value: 1 },
        { label: '背景图片', value: 2 },
      ],
      true,
    ),
  ];
  if (bgtype == 1) {
    fields.push(diyColor('style.background', '背景颜色：', '#ff4c01', '透明'));
  }
  fields.push(
    diySection('组件样式'),
    diyColor('style.bgcolor', '底部背景：', '#f2f2f2', '透明'),
    ...DIY_PADDING_FIELDS,
    ...DIY_RADIUS_FIELDS,
  );
  return fields;
});

const { Form } = useDiyCurItemForm(
  () => props.curItem,
  schema,
  {
    onInit(item) {
      parseIntFields(item, ['style.paddingTop', 'params.limit']);
    },
  },
);
</script>

<template>
  <div>
    <div class="common-form">
      <span>{{ curItem.name }}</span>
    </div>
    <Form />
    <div v-if="curItem.style?.bgtype == 2" class="form-item ml-[100px]">
      <div class="form-label mb-2">背景图片：</div>
      <div class="diy-special-cover">
        <img
          v-img-url="curItem.style?.bgimage"
          alt=""
          @click="editor.onEditorSelectImage(curItem.style!, 'bgimage')"
        />
        <div class="gray">建议尺寸706px*288px</div>
      </div>
    </div>
  </div>
</template>
