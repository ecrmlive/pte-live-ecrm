import type { App, Component } from 'vue';

import {
  ElButton,
  ElCascader,
  ElCheckbox,
  ElCheckboxGroup,
  ElColorPicker,
  ElDatePicker,
  ElDescriptions,
  ElDescriptionsItem,
  ElDialog,
  ElDivider,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElIcon,
  ElImage,
  ElInput,
  ElInputNumber,
  ElLink,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElPagination,
  ElRadio,
  ElRadioButton,
  ElRadioGroup,
  ElSelect,
  ElSlider,
  ElSwitch,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag,
  ElTimeline,
  ElTimelineItem,
  ElTooltip,
  ElUpload,
} from 'element-plus';

import { hydrateMerchantSession } from './merchant-session-bridge';

declare global {
  interface Window {
    ElMessage: typeof ElMessage;
    ElMessageBox: typeof ElMessageBox;
  }
}

/** native Options API / 未显式 import 的页直接使用 `<el-*>`，需全局注册 */
const MERCHANT_GLOBAL_COMPONENTS: Component[] = [
  ElButton,
  ElCascader,
  ElCheckbox,
  ElCheckboxGroup,
  ElColorPicker,
  ElDatePicker,
  ElDescriptions,
  ElDescriptionsItem,
  ElDialog,
  ElDivider,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElIcon,
  ElImage,
  ElInput,
  ElInputNumber,
  ElLink,
  ElOption,
  ElPagination,
  ElRadio,
  ElRadioButton,
  ElRadioGroup,
  ElSelect,
  ElSlider,
  ElSwitch,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag,
  ElTimeline,
  ElTimelineItem,
  ElTooltip,
  ElUpload,
];

/** 商户壳层全局：ElMessage 兼容 + merchantSession 水合 */
export function registerMerchantShell(app: App) {
  for (const comp of MERCHANT_GLOBAL_COMPONENTS) {
    const name = (comp as { name?: string }).name;
    if (name) {
      app.component(name, comp);
    }
  }
  window.ElMessage = ElMessage;
  window.ElMessageBox = ElMessageBox;
  hydrateMerchantSession();
}

/** @deprecated 使用 registerMerchantShell */
export const registerLegacyApp = registerMerchantShell;
