/**
 * DIY 通用层（平台后台 + 店铺后台共用）
 *
 * 职责边界：
 * - 本包：预览壳、尺寸常量等与业务端无关的 DIY 通用能力
 * - 各端 `views/ecrm/diy/**`：组件库条目、属性面板、接口与权限（可逐步下沉到本包）
 */
export {
  ANDROID_PREVIEW,
  DEVICE_PREVIEW_METRICS,
  HARMONY_PREVIEW,
  IPHONE_17_PRO,
  PREVIEW_DEVICE_OPTIONS,
  PREVIEW_DEVICE_STORAGE_KEY,
  isPreviewDevice,
  readStoredPreviewDevice,
  writeStoredPreviewDevice,
} from './constants';
export type {
  DevicePreviewMetrics,
  Iphone17ProMetrics,
  PreviewDevice,
} from './constants';

export { default as DevicePreviewFrame } from './components/device-preview-frame.vue';
/** @deprecated 请使用 DevicePreviewFrame；行为相同（含设备切换） */
export { default as IphonePreviewFrame } from './components/device-preview-frame.vue';
