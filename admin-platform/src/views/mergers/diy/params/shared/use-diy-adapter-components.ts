import { computed } from 'vue';

import { globalShareState } from '@vben/common-ui';

/** Adapter components for DIY param panels (no direct El* in templates). */
export function useDiyAdapterComponents() {
  const components = globalShareState.getComponents();

  return {
    Checkbox: computed(() => components.Checkbox),
    CheckboxGroup: computed(() => components.CheckboxGroup),
    DefaultButton: computed(() => components.DefaultButton),
    Input: computed(() => components.Input),
    PrimaryButton: computed(() => components.PrimaryButton),
    RadioGroup: computed(() => components.RadioGroup),
  };
}
