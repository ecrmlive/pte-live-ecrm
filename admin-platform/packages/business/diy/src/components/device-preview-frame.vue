<script setup lang="ts">
/**
 * DIY 多端预览壳：苹果 / 安卓 / 鸿蒙
 *
 * 结构：状态栏 → 导航栏 → 内容区 → 底部安全区
 * 白底内容对齐屏宽内侧；sideGutter 只占外层 padding + 内容伸出区（工具条）。
 * 内容区可滚轮/触控板滚动，不展示原生或代理滚动条。
 *
 * 切换器吸附在预览工作区（灰底 host）右上角，非机身右缘；偏好持久化到 localStorage。
 */
import { computed, ref } from 'vue';

import {
  DEVICE_PREVIEW_METRICS,
  PREVIEW_DEVICE_OPTIONS,
  type PreviewDevice,
  readStoredPreviewDevice,
  writeStoredPreviewDevice,
} from '../constants';

const props = withDefaults(
  defineProps<{
    /** 受控设备；不传则用内部状态 + localStorage */
    device?: PreviewDevice;
    /** 是否显示外壳切换器 */
    showDeviceSwitcher?: boolean;
    /** 导航栏标题 */
    title?: string;
    /** 是否显示返回箭头 */
    showBack?: boolean;
    /** 隐藏导航栏（H5 iframe 全页预览时由页面自带导航） */
    hideNav?: boolean;
    /** 导航栏是否高亮（如 DIY「页面设置」选中） */
    navActive?: boolean;
    /** 右侧留给悬浮工具条的外边距（避免机身 overflow 裁切） */
    sideGutter?: number;
    /** 左侧留给组件名标签等的外边距 */
    sideGutterLeft?: number;
    /** 内容区背景 */
    contentBg?: string;
    /** 底部「提交」条占位（入驻表单等） */
    showSubmitBar?: boolean;
    submitText?: string;
  }>(),
  {
    showDeviceSwitcher: true,
    title: '页面标题',
    showBack: true,
    hideNav: false,
    navActive: false,
    sideGutter: 52,
    sideGutterLeft: 0,
    contentBg: '#ffffff',
    showSubmitBar: false,
    submitText: '提交',
  },
);

const emit = defineEmits<{
  'nav-click': [];
  'back-click': [];
  'update:device': [PreviewDevice];
}>();

const internalDevice = ref<PreviewDevice>(readStoredPreviewDevice());

const activeDevice = computed<PreviewDevice>(
  () => props.device ?? internalDevice.value,
);

const metrics = computed(() => DEVICE_PREVIEW_METRICS[activeDevice.value]);

function setDevice(next: PreviewDevice) {
  if (activeDevice.value === next) return;
  if (props.device === undefined) {
    internalDevice.value = next;
  }
  writeStoredPreviewDevice(next);
  emit('update:device', next);
}

const statusTime = computed(() => {
  const d = new Date();
  const h = d.getHours();
  const m = String(d.getMinutes()).padStart(2, '0');
  return `${h}:${m}`;
});

/** 有侧边 gutter 时加宽非滚动壳，供工具条/组件名标签伸出 */
const needsGutterChrome = computed(
  () => props.sideGutter > 0 || props.sideGutterLeft > 0,
);

const outerStyle = computed(() => {
  const m = metrics.value;
  return {
    width: `${m.screenWidth + props.sideGutter + props.sideGutterLeft}px`,
    ['--device-screen-w' as string]: `${m.screenWidth}px`,
    ['--device-screen-h' as string]: `${m.screenHeight}px`,
    ['--device-safe-top' as string]: `${m.safeTop}px`,
    ['--device-safe-bottom' as string]: `${m.safeBottom}px`,
    ['--device-status-h' as string]: `${m.statusBarHeight}px`,
    ['--device-nav-h' as string]: `${m.navBarHeight}px`,
    ['--device-side-gutter' as string]: `${props.sideGutter}px`,
    ['--device-side-gutter-left' as string]: `${props.sideGutterLeft}px`,
  };
});
</script>

<template>
  <!-- 铺满灰底预览区：切换器相对此 host 右上吸附；机身仍水平居中 -->
  <div
    class="device-preview-host"
    :class="{ 'device-preview-host--with-switcher': showDeviceSwitcher }"
  >
    <div
      v-if="showDeviceSwitcher"
      class="device-preview__switcher-bar"
    >
      <div
        class="device-preview__switcher"
        role="radiogroup"
        aria-label="预览设备"
      >
        <button
          v-for="opt in PREVIEW_DEVICE_OPTIONS"
          :key="opt.value"
          type="button"
          class="device-preview__switcher-btn"
          role="radio"
          :aria-checked="activeDevice === opt.value"
          :class="{ 'is-active': activeDevice === opt.value }"
          @click="setDevice(opt.value)"
        >
          {{ opt.label }}
        </button>
      </div>
    </div>

    <div
      class="device-preview"
      :class="[`device-preview--${activeDevice}`]"
      :style="outerStyle"
    >
      <div class="device-preview__device">
        <div class="device-preview__screen">
          <!-- 1. 状态栏 -->
          <div class="device-preview__status" aria-hidden="true">
            <!-- iOS -->
            <template v-if="activeDevice === 'ios'">
              <span class="device-preview__time">{{ statusTime }}</span>
              <div class="device-preview__island" />
              <div class="device-preview__status-right">
                <svg class="device-preview__icon" viewBox="0 0 18 12" width="18" height="12">
                  <rect x="0" y="8" width="3" height="4" rx="0.5" fill="currentColor" />
                  <rect x="4.5" y="5.5" width="3" height="6.5" rx="0.5" fill="currentColor" />
                  <rect x="9" y="3" width="3" height="9" rx="0.5" fill="currentColor" />
                  <rect x="13.5" y="0.5" width="3" height="11.5" rx="0.5" fill="currentColor" />
                </svg>
                <svg class="device-preview__icon" viewBox="0 0 16 12" width="16" height="12">
                  <path
                    fill="currentColor"
                    d="M8 3.2c1.7 0 3.2.6 4.4 1.7l-1.1 1.1A4.6 4.6 0 0 0 8 4.8c-1.2 0-2.3.5-3.2 1.2L3.7 4.9A6.3 6.3 0 0 1 8 3.2zm0 2.7c1 0 1.9.4 2.6 1l-1.1 1.1A2.2 2.2 0 0 0 8 7.2c-.6 0-1.1.2-1.5.6L5.4 6.7A3.7 3.7 0 0 1 8 5.9zM8 10.2a1.1 1.1 0 1 0 0-2.2 1.1 1.1 0 0 0 0 2.2z"
                  />
                </svg>
                <div class="device-preview__battery device-preview__battery--ios">
                  <div class="device-preview__battery-body">
                    <div class="device-preview__battery-level" />
                  </div>
                  <div class="device-preview__battery-cap" />
                </div>
              </div>
            </template>

          <!-- Android -->
          <template v-else-if="activeDevice === 'android'">
            <span class="device-preview__time device-preview__time--android">
              {{ statusTime }}
            </span>
            <div class="device-preview__punch" />
            <div class="device-preview__status-right device-preview__status-right--android">
              <svg class="device-preview__icon" viewBox="0 0 16 12" width="15" height="11">
                <path
                  fill="currentColor"
                  d="M8 2.4c1.9 0 3.6.7 4.9 1.9l-1.2 1.2A5 5 0 0 0 8 4.1a5 5 0 0 0-3.7 1.4L3.1 4.3A6.8 6.8 0 0 1 8 2.4zm0 3c1.1 0 2.1.4 2.9 1.1L9.7 7.7A2.4 2.4 0 0 0 8 7a2.4 2.4 0 0 0-1.7.7L5.1 6.5A3.9 3.9 0 0 1 8 5.4zM8 10.5a1.2 1.2 0 1 0 0-2.4 1.2 1.2 0 0 0 0 2.4z"
                />
              </svg>
              <svg class="device-preview__icon" viewBox="0 0 18 12" width="16" height="11">
                <rect x="0" y="7.5" width="2.8" height="3.5" rx="0.4" fill="currentColor" />
                <rect x="4.2" y="5" width="2.8" height="6" rx="0.4" fill="currentColor" />
                <rect x="8.4" y="2.5" width="2.8" height="8.5" rx="0.4" fill="currentColor" />
                <rect x="12.6" y="0.5" width="2.8" height="10.5" rx="0.4" fill="currentColor" opacity="0.35" />
              </svg>
              <div class="device-preview__battery device-preview__battery--android">
                <div class="device-preview__battery-body">
                  <div class="device-preview__battery-level" />
                </div>
              </div>
            </div>
          </template>

          <!-- HarmonyOS -->
          <template v-else>
            <span class="device-preview__time device-preview__time--harmony">
              {{ statusTime }}
            </span>
            <div class="device-preview__capsule" />
            <div class="device-preview__status-right device-preview__status-right--harmony">
              <svg class="device-preview__icon" viewBox="0 0 18 12" width="16" height="11">
                <rect x="0" y="8" width="3" height="4" rx="1" fill="currentColor" />
                <rect x="4.5" y="5.5" width="3" height="6.5" rx="1" fill="currentColor" />
                <rect x="9" y="3" width="3" height="9" rx="1" fill="currentColor" />
                <rect x="13.5" y="0.5" width="3" height="11.5" rx="1" fill="currentColor" />
              </svg>
              <svg class="device-preview__icon" viewBox="0 0 14 12" width="13" height="11">
                <path
                  fill="currentColor"
                  d="M7 1.8c1.6 0 3 .6 4.1 1.6l-1 1A4.2 4.2 0 0 0 7 3.4a4.2 4.2 0 0 0-3.1 1.3l-1-1A5.8 5.8 0 0 1 7 1.8zm0 2.6c.9 0 1.7.3 2.3.9L8.3 6.3A1.9 1.9 0 0 0 7 5.7c-.5 0-1 .2-1.3.6L4.7 5.3A3.3 3.3 0 0 1 7 4.4zM7 9.6a1 1 0 1 0 0-2 1 1 0 0 0 0 2z"
                />
              </svg>
              <div class="device-preview__battery device-preview__battery--harmony">
                <div class="device-preview__battery-body">
                  <div class="device-preview__battery-level" />
                </div>
                <div class="device-preview__battery-cap" />
              </div>
            </div>
          </template>
        </div>

        <!-- 2. 导航栏 -->
        <div
          v-if="!hideNav"
          class="device-preview__nav"
          :class="{ 'is-active': navActive }"
          @click="emit('nav-click')"
        >
          <button
            v-if="showBack"
            type="button"
            class="device-preview__back"
            aria-label="返回"
            @click.stop="emit('back-click')"
          >
            <svg viewBox="0 0 12 20" width="12" height="20">
              <path
                d="M10.5 1.5 2 10l8.5 8.5"
                fill="none"
                stroke="currentColor"
                stroke-width="2.2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>
          <div v-else class="device-preview__back-spacer" />
          <div class="device-preview__title">{{ title }}</div>
          <div class="device-preview__nav-right" />
        </div>

        <!-- 3. 内容区：白底对齐屏宽；gutter 仅供工具条伸出；滚动条隐藏 -->
        <div
          class="device-preview__body"
          :class="{ 'device-preview__body--gutter': needsGutterChrome }"
        >
          <div class="device-preview__scroll">
            <div
              class="device-preview__body-inner"
              :style="{ background: contentBg }"
            >
              <slot />
            </div>
          </div>
        </div>

        <!-- 可选提交条占位 -->
        <div v-if="showSubmitBar" class="device-preview__submit">
          <div class="device-preview__submit-btn">{{ submitText }}</div>
        </div>

        <!-- 4. 底部安全区 -->
        <div class="device-preview__home" aria-hidden="true">
          <div class="device-preview__home-bar" />
        </div>
      </div>
    </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
/* 铺满父级灰底工作区（.diy-center / .phone-frame-host），供切换器定位 */
.device-preview-host {
  position: relative;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  align-self: stretch;
  width: 100%;
  min-width: 0;
  min-height: 100%;
}

/*
 * sticky 全宽零高条：滚动灰底预览列时仍吸附在可视区右上；
 * 按钮溢出可见，不占机身上方布局空间。
 */
.device-preview__switcher-bar {
  position: sticky;
  top: 12px;
  z-index: 20;
  box-sizing: border-box;
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  width: 100%;
  height: 0;
  padding-right: 12px;
  overflow: visible;
  pointer-events: none;
}

.device-preview {
  position: relative;
  box-sizing: border-box;
  display: flex;
  flex-shrink: 0;
  justify-content: flex-start;
  padding-right: var(--device-side-gutter, 52px);
  padding-left: var(--device-side-gutter-left, 0px);
  overflow: visible;
}

/* 灰底工作区右上角分段切换：轨道 + 等高胶囊，激活仅改色不改布局 */
.device-preview__switcher {
  box-sizing: border-box;
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  height: 32px;
  padding: 2px;
  border: 1px solid hsl(var(--border, 220 13% 91%));
  border-radius: 8px;
  background: hsl(var(--muted, 210 40% 96%));
  box-shadow: 0 1px 2px rgb(15 23 42 / 6%);
  pointer-events: auto;
}

.device-preview__switcher-btn {
  box-sizing: border-box;
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  min-width: 48px;
  height: 28px;
  margin: 0;
  padding: 0 12px;
  border: 1px solid transparent;
  border-radius: 6px;
  background: transparent;
  color: hsl(var(--muted-foreground, 215 16% 47%));
  font-size: 13px;
  font-weight: 500;
  line-height: 1;
  white-space: nowrap;
  vertical-align: middle;
  appearance: none;
  -webkit-appearance: none;
  transform: none;
  cursor: pointer;
  transition:
    background 0.15s ease,
    color 0.15s ease;
}

.device-preview__switcher-btn:hover:not(.is-active) {
  color: hsl(var(--foreground, 222 47% 11%));
  background: hsl(var(--background, 0 0% 100%) / 70%);
}

.device-preview__switcher-btn.is-active {
  margin: 0;
  border-color: transparent;
  color: hsl(var(--primary-foreground, 0 0% 100%));
  background: hsl(var(--primary, 221 83% 53%));
  transform: none;
}

.device-preview__device {
  box-sizing: border-box;
  width: calc(var(--device-screen-w) + 16px);
  padding: 8px;
  overflow: visible;
  transition:
    border-radius 0.2s ease,
    background 0.2s ease,
    box-shadow 0.2s ease;
}

/* —— iOS 壳 —— */
.device-preview--ios .device-preview__device {
  border: 1px solid #1c1c1e;
  border-radius: 44px;
  background: linear-gradient(160deg, #3a3a3c 0%, #1c1c1e 55%, #0a0a0a 100%);
  box-shadow:
    0 12px 40px rgb(0 0 0 / 18%),
    0 2px 8px rgb(0 0 0 / 8%),
    inset 0 1px 0 rgb(255 255 255 / 12%);
}

.device-preview--ios .device-preview__screen {
  border-radius: 36px;
}

.device-preview--ios .device-preview__status {
  padding: 0 22px 10px;
  border-radius: 36px 36px 0 0;
}

.device-preview--ios .device-preview__home {
  border-radius: 0 0 36px 36px;
}

.device-preview--ios .device-preview__home-bar {
  width: 128px;
  height: 5px;
  border-radius: 3px;
  background: #111;
}

/* —— Android 壳 —— */
.device-preview--android .device-preview__device {
  border: 2px solid #2a2a2e;
  border-radius: 28px;
  background: linear-gradient(180deg, #2e2e32 0%, #141416 100%);
  box-shadow:
    0 10px 32px rgb(0 0 0 / 16%),
    inset 0 0 0 1px rgb(255 255 255 / 6%);
}

.device-preview--android .device-preview__screen {
  border-radius: 20px;
}

.device-preview--android .device-preview__status {
  align-items: center;
  padding: 0 16px;
  border-radius: 20px 20px 0 0;
}

.device-preview--android .device-preview__home {
  border-radius: 0 0 20px 20px;
}

.device-preview--android .device-preview__home-bar {
  width: 108px;
  height: 4px;
  border-radius: 999px;
  background: #3c4043;
}

/* —— HarmonyOS 壳 —— */
.device-preview--harmony .device-preview__device {
  border: 1px solid #1a2230;
  border-radius: 36px;
  background: linear-gradient(155deg, #2a3548 0%, #121820 50%, #0b1018 100%);
  box-shadow:
    0 12px 36px rgb(15 30 60 / 20%),
    0 2px 8px rgb(0 0 0 / 10%),
    inset 0 1px 0 rgb(120 160 220 / 18%);
}

.device-preview--harmony .device-preview__screen {
  border-radius: 28px;
}

.device-preview--harmony .device-preview__status {
  align-items: center;
  padding: 0 18px;
  border-radius: 28px 28px 0 0;
}

.device-preview--harmony .device-preview__home {
  border-radius: 0 0 28px 28px;
}

.device-preview--harmony .device-preview__home-bar {
  width: 96px;
  height: 6px;
  border-radius: 8px 8px 4px 4px;
  background: linear-gradient(180deg, #1a2332 0%, #0d121a 100%);
}

/* overflow visible：选中工具条可伸出到 sideGutter，不被机身裁切 */
.device-preview__screen {
  display: flex;
  flex-direction: column;
  width: var(--device-screen-w);
  height: var(--device-screen-h);
  overflow: visible;
  background: #fff;
  transition: border-radius 0.2s ease;
}

.device-preview__status {
  position: relative;
  z-index: 1;
  display: flex;
  flex-shrink: 0;
  align-items: flex-end;
  justify-content: space-between;
  height: var(--device-status-h);
  color: #111;
  background: #fff;
}

.device-preview__time {
  min-width: 54px;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 0.2px;
  line-height: 1;
}

.device-preview__time--android {
  min-width: 40px;
  font-size: 13px;
  font-weight: 500;
  font-variant-numeric: tabular-nums;
}

.device-preview__time--harmony {
  min-width: 48px;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.4px;
}

/* iOS Dynamic Island */
.device-preview__island {
  position: absolute;
  top: 11px;
  left: 50%;
  width: 118px;
  height: 34px;
  border-radius: 20px;
  background: #0a0a0a;
  transform: translateX(-50%);
}

/* Android 居中挖孔 */
.device-preview__punch {
  position: absolute;
  top: 8px;
  left: 50%;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #0a0a0a;
  box-shadow: inset 0 0 0 1.5px #1c1c1e;
  transform: translateX(-50%);
}

/* Harmony 顶部胶囊 */
.device-preview__capsule {
  position: absolute;
  top: 8px;
  left: 50%;
  width: 86px;
  height: 22px;
  border-radius: 12px;
  background: #0b1018;
  box-shadow: inset 0 0 0 1px rgb(80 120 180 / 25%);
  transform: translateX(-50%);
}

.device-preview__status-right {
  display: flex;
  align-items: center;
  gap: 5px;
  min-width: 54px;
  justify-content: flex-end;
  color: #111;
}

.device-preview__status-right--android {
  gap: 6px;
  min-width: 62px;
}

.device-preview__status-right--harmony {
  gap: 4px;
  min-width: 58px;
}

.device-preview__icon {
  display: block;
}

.device-preview__battery {
  display: flex;
  align-items: center;
}

.device-preview__battery-body {
  box-sizing: border-box;
  width: 22px;
  height: 11px;
  padding: 1px;
  border: 1px solid rgb(0 0 0 / 45%);
  border-radius: 3px;
}

.device-preview__battery--android .device-preview__battery-body {
  width: 18px;
  height: 10px;
  border-radius: 2px;
  border-color: rgb(0 0 0 / 55%);
}

.device-preview__battery--harmony .device-preview__battery-body {
  width: 20px;
  height: 10px;
  border-radius: 4px;
}

.device-preview__battery-level {
  width: 80%;
  height: 100%;
  border-radius: 1.5px;
  background: #111;
}

.device-preview__battery--android .device-preview__battery-level {
  width: 70%;
  border-radius: 1px;
  background: #1a73e8;
}

.device-preview__battery--harmony .device-preview__battery-level {
  width: 75%;
  border-radius: 2px;
  background: #0a59f7;
}

.device-preview__battery-cap {
  width: 1.5px;
  height: 4px;
  margin-left: 1px;
  border-radius: 0 1px 1px 0;
  background: rgb(0 0 0 / 45%);
}

.device-preview__nav {
  position: relative;
  display: flex;
  flex-shrink: 0;
  align-items: center;
  height: var(--device-nav-h);
  padding: 0 6px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
}

.device-preview__nav.is-active {
  box-shadow: inset 0 0 0 2px hsl(var(--primary));
}

.device-preview__back,
.device-preview__back-spacer,
.device-preview__nav-right {
  flex-shrink: 0;
  width: 44px;
  height: 44px;
}

.device-preview__back {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  border: 0;
  background: transparent;
  color: #111;
  cursor: pointer;
}

.device-preview__title {
  flex: 1;
  overflow: hidden;
  color: #111;
  font-size: 17px;
  font-weight: 600;
  line-height: 44px;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.device-preview--android .device-preview__title,
.device-preview--harmony .device-preview__title {
  font-size: 16px;
  font-weight: 500;
}

/*
 * 滚动壳默认 = 屏宽。
 * 有 gutter 时壳加宽仅用于绘制伸出的工具条/标签；
 * 滚动面隐藏原生条（滚轮/触控板仍可滚动）。
 */
.device-preview__body {
  position: relative;
  flex: 1;
  min-height: 0;
  box-sizing: border-box;
  width: var(--device-screen-w);
  overflow: visible;
  background: transparent;
}

.device-preview__body--gutter {
  width: calc(
    var(--device-screen-w) + var(--device-side-gutter-left, 0px) +
      var(--device-side-gutter, 0px)
  );
  margin-left: calc(-1 * var(--device-side-gutter-left, 0px));
}

.device-preview__scroll {
  box-sizing: border-box;
  width: 100%;
  height: 100%;
  overflow-x: hidden;
  overflow-y: auto;
  overscroll-behavior: contain;
  background: transparent;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.device-preview__scroll::-webkit-scrollbar {
  display: none;
  width: 0;
  height: 0;
}

.device-preview__body-inner {
  position: relative;
  box-sizing: border-box;
  width: var(--device-screen-w);
  min-height: 100%;
  margin-left: var(--device-side-gutter-left, 0px);
  overflow: visible;
  background: #fff;
}

/* 无 gutter：壳即屏宽，inner 跟随 100% */
.device-preview__body:not(.device-preview__body--gutter) .device-preview__body-inner {
  width: 100%;
  margin-left: 0;
}

.device-preview__submit {
  flex-shrink: 0;
  padding: 8px 16px 0;
  background: #fff;
  border-top: 1px solid #f0f0f0;
}

.device-preview__submit-btn {
  height: 40px;
  border-radius: 20px;
  background: hsl(var(--primary));
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  line-height: 40px;
  text-align: center;
}

.device-preview--android .device-preview__submit-btn {
  border-radius: 8px;
}

.device-preview--harmony .device-preview__submit-btn {
  border-radius: 12px;
}

.device-preview__home {
  z-index: 1;
  display: flex;
  flex-shrink: 0;
  align-items: flex-end;
  justify-content: center;
  height: var(--device-safe-bottom);
  padding-bottom: 8px;
  background: #fff;
}
</style>
