<script lang="ts" setup>
import { ElButton } from 'element-plus';

/** @deprecated 请使用 Vben 查询表单（withPlatformListSearchLayout） */
defineEmits<{
  reset: [];
  search: [];
}>();
</script>

<template>
  <ElButton type="primary" @click="$emit('search')">查询</ElButton>
  <ElButton @click="$emit('reset')">重置</ElButton>
</template>
