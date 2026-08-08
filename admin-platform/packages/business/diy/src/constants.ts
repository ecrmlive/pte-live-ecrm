/**
 * DIY 预览设备尺寸与壳层常量（逻辑点 / CSS px）
 *
 * 内容区逻辑宽统一 402，避免切换外壳时 DIY 布局跳动；
 * 状态栏 / 底部安全区高度按系统习惯略有差异。
 */

export type PreviewDevice = 'ios' | 'android' | 'harmony';

export const PREVIEW_DEVICE_STORAGE_KEY = 'diy-preview-device';

export const PREVIEW_DEVICE_OPTIONS: ReadonlyArray<{
  value: PreviewDevice;
  label: string;
}> = [
  { value: 'ios', label: '苹果' },
  { value: 'android', label: '安卓' },
  { value: 'harmony', label: '鸿蒙' },
] as const;

/** iPhone 17 Pro（默认外壳） */
export const IPHONE_17_PRO = {
  screenWidth: 402,
  screenHeight: 874,
  safeTop: 62,
  safeBottom: 34,
  /** 含 Dynamic Island 的状态栏高度 */
  statusBarHeight: 54,
  /** 导航栏高度（项目约定，对齐 iOS UINavigationBar） */
  navBarHeight: 44,
} as const;

export type Iphone17ProMetrics = typeof IPHONE_17_PRO;

export type DevicePreviewMetrics = {
  screenWidth: number;
  screenHeight: number;
  safeTop: number;
  safeBottom: number;
  statusBarHeight: number;
  navBarHeight: number;
};

/** Android 典型全面屏（同逻辑宽，扁平底栏） */
export const ANDROID_PREVIEW = {
  screenWidth: 402,
  screenHeight: 874,
  safeTop: 40,
  safeBottom: 24,
  statusBarHeight: 36,
  navBarHeight: 48,
} as const satisfies DevicePreviewMetrics;

/** HarmonyOS 典型全面屏（胶囊状态区 + 弧形底指示） */
export const HARMONY_PREVIEW = {
  screenWidth: 402,
  screenHeight: 874,
  safeTop: 48,
  safeBottom: 30,
  statusBarHeight: 42,
  navBarHeight: 48,
} as const satisfies DevicePreviewMetrics;

export const DEVICE_PREVIEW_METRICS: Record<
  PreviewDevice,
  DevicePreviewMetrics
> = {
  ios: { ...IPHONE_17_PRO },
  android: { ...ANDROID_PREVIEW },
  harmony: { ...HARMONY_PREVIEW },
};

export function isPreviewDevice(value: unknown): value is PreviewDevice {
  return value === 'ios' || value === 'android' || value === 'harmony';
}

export function readStoredPreviewDevice(): PreviewDevice {
  if (typeof localStorage === 'undefined') return 'ios';
  try {
    const raw = localStorage.getItem(PREVIEW_DEVICE_STORAGE_KEY);
    return isPreviewDevice(raw) ? raw : 'ios';
  } catch {
    return 'ios';
  }
}

export function writeStoredPreviewDevice(device: PreviewDevice) {
  if (typeof localStorage === 'undefined') return;
  try {
    localStorage.setItem(PREVIEW_DEVICE_STORAGE_KEY, device);
  } catch {
    /* ignore quota / private mode */
  }
}
