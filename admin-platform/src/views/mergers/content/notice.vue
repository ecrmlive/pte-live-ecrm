<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { Page } from '@vben/common-ui';
const route = useRoute();
const title = computed(() => (route.meta.title as string) || '功能建设中');
</script>
<template>
  <Page :title="title" description="Vben 壳已就绪，业务页待迁移。">
    <div class="text-muted-foreground p-4 text-sm">path: {{ route.path }}</div>
  </Page>
</template>
