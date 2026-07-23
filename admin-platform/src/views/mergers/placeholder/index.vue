<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';

const route = useRoute();
const title = computed(() => (route.meta.title as string) || '功能建设中');
</script>

<template>
  <Page :title="title" description="该菜单已接入 Vben 壳层，业务页尚未迁移。">
    <div class="text-muted-foreground p-4 text-sm">path: {{ route.path }}</div>
  </Page>
</template>
