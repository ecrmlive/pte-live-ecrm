import type { VbenFormSchema } from '#/adapter/form';

import { markRaw } from 'vue';

import NativeSectionTitle from '#/components/shop/native-section-title.vue';

import DiyColorField from './diy-color-field.vue';
import DiySliderField from './diy-slider-field.vue';

export function diySection(title: string, hint?: string, fieldName?: string): VbenFormSchema {
  return {
    component: markRaw(NativeSectionTitle),
    componentProps: { hint, title },
    fieldName: fieldName ?? `_section_${title.replace(/\s/g, '_')}`,
    formItemClass: 'col-span-full native-section-title-item diy-section-title-item',
    hideLabel: true,
  };
}

export function diyColor(
  path: string,
  label: string,
  defaultColor: string,
  placeholder?: string,
): VbenFormSchema {
  return {
    component: markRaw(DiyColorField),
    componentProps: { defaultColor, placeholder },
    fieldName: path,
    label,
  };
}

export function diySlider(
  path: string,
  label: string,
  opts?: { max?: number; min?: number },
): VbenFormSchema {
  return {
    component: markRaw(DiySliderField),
    componentProps: { max: opts?.max, min: opts?.min },
    fieldName: path,
    label,
  };
}

export function diyInput(
  path: string,
  label: string,
  componentProps?: Record<string, unknown>,
): VbenFormSchema {
  return {
    component: 'Input',
    componentProps,
    fieldName: path,
    label,
  };
}

export function diyRadioGroup(
  path: string,
  label: string,
  options: Array<{ label: string; value: number | string }>,
  isButton = false,
): VbenFormSchema {
  return {
    component: 'RadioGroup',
    componentProps: { isButton, options },
    fieldName: path,
    label,
  };
}

export const DIY_PADDING_FIELDS: VbenFormSchema[] = [
  diySlider('style.paddingTop', '上边距：'),
  diySlider('style.paddingBottom', '下边距：'),
  diySlider('style.paddingLeft', '左右边距：'),
];

export const DIY_RADIUS_FIELDS: VbenFormSchema[] = [
  diySlider('style.topRadio', '上圆角：'),
  diySlider('style.bottomRadio', '下圆角：'),
];

export function diyBgColors(
  backgroundDefault = '#ffffff',
  bgcolorDefault = '#f2f2f2',
): VbenFormSchema[] {
  return [
    diyColor('style.background', '底部背景：', backgroundDefault, '透明'),
    diyColor('style.bgcolor', '组件背景：', bgcolorDefault, '透明'),
  ];
}
