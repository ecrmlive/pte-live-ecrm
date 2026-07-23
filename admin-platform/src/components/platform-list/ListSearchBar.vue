<script lang="ts" setup>
/** @deprecated 请使用 Vben 查询表单（`withPlatformListSearchLayout` + `useVbenVxeGrid` formOptions，或 `PlatformListQueryForm`） */
import { ElButton, ElForm } from 'element-plus';

defineEmits<{
  reset: [];
  search: [];
}>();
</script>

<template>
  <div class="platform-list-search">
    <div class="flex flex-wrap items-start justify-between gap-3">
      <ElForm :inline="true" class="!mb-0 flex flex-1 flex-wrap gap-2">
        <slot name="filters" />
      </ElForm>
      <div class="flex shrink-0 items-center gap-2">
        <ElButton @click="$emit('reset')">重置</ElButton>
        <ElButton type="primary" @click="$emit('search')">查询</ElButton>
      </div>
    </div>
    <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>
