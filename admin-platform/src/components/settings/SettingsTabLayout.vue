<script setup lang="ts">
/**
 * 多 Tab 设置类页面壳（金标准分区）：
 * 顶栏 Tabs → 中间可滚动内容区（居中卡片）→ 底栏固定操作条。
 * 用于协议/说明/配置类编辑页；列表页仍走店铺列表金标准。
 */
defineSlots<{
  /** 顶栏：通常放 ElTabs（仅 header，无 pane 内嵌） */
  tabs?: () => unknown;
  /** 中间主内容：表单 / 富文本等 */
  default?: () => unknown;
  /** 底栏操作：预览 / 提交等 */
  actions?: () => unknown;
}>();
</script>

<template>
  <div class="settings-tab-layout">
    <div class="settings-tab-layout__tabs">
      <slot name="tabs" />
    </div>

    <div class="settings-tab-layout__body">
      <div class="settings-tab-layout__card">
        <slot />
      </div>
    </div>

    <div v-if="$slots.actions" class="settings-tab-layout__footer">
      <div class="settings-tab-layout__footer-inner">
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-tab-layout {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: hsl(var(--background));
}

.settings-tab-layout__tabs {
  flex-shrink: 0;
  padding: 0 16px;
  background: hsl(var(--card));
  border-bottom: 1px solid hsl(var(--border));
}

.settings-tab-layout__tabs :deep(.el-tabs__header) {
  margin: 0;
}

.settings-tab-layout__tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.settings-tab-layout__tabs :deep(.el-tabs__item) {
  height: 48px;
  line-height: 48px;
}

.settings-tab-layout__body {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 16px 16px 24px;
  background: hsl(var(--background-deep, var(--background)));
}

.settings-tab-layout__card {
  width: 100%;
  max-width: 1080px;
  margin: 0 auto;
  padding: 16px;
  background: hsl(var(--card));
  border: 1px solid hsl(var(--border));
  border-radius: 8px;
}

.settings-tab-layout__footer {
  flex-shrink: 0;
  border-top: 1px solid hsl(var(--border));
  background: hsl(var(--card));
}

.settings-tab-layout__footer-inner {
  display: flex;
  gap: 12px;
  align-items: center;
  justify-content: center;
  min-height: 56px;
  padding: 10px 16px;
}
</style>
