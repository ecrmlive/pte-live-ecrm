<template>
  <a-card :bordered="false" class="card">
    <a-result status="info" :title="title" sub-title="该模块将在后续阶段竖切实现，当前仅菜单入口占位。">
      <template #extra>
        <a-button type="primary" @click="$router.push('/dashboard')">返回工作台</a-button>
      </template>
    </a-result>
  </a-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const title = computed(() => `「${route.path}」建设中`);
</script>

<style scoped>
.card {
  border-radius: 14px;
  min-height: 420px;
}
</style>
