<script setup lang="ts">
/**
 * CRMEB 列表页骨架：顶筛 → 工具条 → 表格区 → 底部分页。
 * 字段顺序由调用方 slot 控制，本组件只统一间距与分区。
 */
defineProps<{
  description?: string;
  title: string;
}>();
</script>

<template>
  <div class="ecrm-list-page flex flex-col gap-4">
    <div class="ecrm-list-page__header">
      <h2 class="text-lg font-semibold leading-7">{{ title }}</h2>
      <p v-if="description" class="mt-1 text-sm text-muted-foreground">
        {{ description }}
      </p>
    </div>

    <div
      v-if="$slots.filters || $slots.actions"
      class="ecrm-list-page__toolbar rounded-md border border-border bg-card p-4"
    >
      <div
        v-if="$slots.filters"
        class="ecrm-list-page__filters flex flex-wrap items-start gap-2"
      >
        <slot name="filters" />
      </div>
      <div
        v-if="$slots.actions"
        class="ecrm-list-page__actions mt-3 flex flex-wrap items-center gap-2"
      >
        <slot name="actions" />
      </div>
    </div>

    <div class="ecrm-list-page__table rounded-md border border-border bg-card p-4">
      <slot />
      <div
        v-if="$slots.pager"
        class="ecrm-list-page__pager mt-4 flex justify-end"
      >
        <slot name="pager" />
      </div>
    </div>
  </div>
</template>
