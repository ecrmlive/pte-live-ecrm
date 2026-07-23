import type { VbenFormSchema } from '#/adapter/form';

import { markRaw } from 'vue';

import NativeSectionTitle from '#/components/shop/native-section-title.vue';

type SectionTitleOptions = {
  dependencies?: VbenFormSchema['dependencies'];
  formItemClass?: string;
  hint?: string;
};

/** In-form section header: single accent bar + title, no Divider horizontal lines. */
export function nativeSectionTitle(
  fieldName: string,
  title: string,
  options?: SectionTitleOptions,
): VbenFormSchema {
  return {
    component: markRaw(NativeSectionTitle),
    componentProps: { hint: options?.hint, title },
    dependencies: options?.dependencies,
    fieldName,
    formItemClass: options?.formItemClass ?? 'col-span-full native-section-title-item',
    hideLabel: true,
  };
}
