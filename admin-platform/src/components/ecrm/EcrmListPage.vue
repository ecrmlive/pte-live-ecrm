<script setup lang="ts">
/**
 * CRMEB 列表页骨架（结构对齐 mer.crmeb.net，颜色跟 Vben 主题）：
 * 筛选区 → 状态 Tab → 主操作（如「添加」）→ 表格（操作列可 fixed）→ 分页。
 */
defineProps<{
  /** @deprecated 菜单/面包屑已够用；仅特殊页需要时再传 */
  description?: string;
  /** @deprecated 菜单/面包屑已够用；仅特殊页需要时再传 */
  title?: string;
}>();
</script>

<template>
  <div class="ecrm-list-page flex min-w-0 flex-col gap-3">
    <div v-if="title || description" class="ecrm-list-page__header">
      <h2 v-if="title" class="text-base font-semibold leading-6">{{ title }}</h2>
      <p v-if="description" class="mt-1 text-sm text-muted-foreground">
        {{ description }}
      </p>
    </div>

    <!-- 1. 筛选区：多列表单项网格，搜索/重置跟在最后一行 -->
    <div
      v-if="$slots.filters"
      class="ecrm-list-page__filters rounded-md border border-border bg-card px-4 py-4"
    >
      <slot name="filters" />
    </div>

    <!-- 2. 状态子 Tab（下划线），单独一行 -->
    <div v-if="$slots.tabs" class="ecrm-list-page__tabs min-w-0">
      <slot name="tabs" />
    </div>

    <!-- 3. 主操作（添加等）在 Tab 下方、表格上方，左对齐 -->
    <div
      v-if="$slots.actions"
      class="ecrm-list-page__actions flex flex-wrap items-center gap-2"
    >
      <slot name="actions" />
    </div>

    <!-- 4. 表格 + 分页；横向滚动时操作列可 fixed="right" -->
    <div class="ecrm-list-page__table min-w-0 rounded-md border border-border bg-card p-4">
      <div class="ecrm-list-page__table-inner min-w-0">
        <slot />
      </div>
      <div
        v-if="$slots.pager"
        class="ecrm-list-page__pager mt-4 flex justify-end overflow-x-auto"
      >
        <slot name="pager" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.ecrm-list-page__filters :deep(.el-form-item) {
  margin-bottom: 12px;
  margin-right: 0;
}

.ecrm-list-page__filters :deep(.el-form--inline .el-form-item) {
  margin-right: 16px;
}

.ecrm-list-page__pager :deep(.el-pagination) {
  display: inline-flex;
  flex-wrap: nowrap;
  align-items: center;
  justify-content: flex-end;
  max-width: 100%;
}

.ecrm-list-page__pager :deep(.el-pagination > *) {
  flex-shrink: 0;
}

.ecrm-list-page__table-inner :deep(.el-table) {
  width: 100%;
}
</style>
