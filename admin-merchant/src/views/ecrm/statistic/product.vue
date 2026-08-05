<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { getMerchantProductStatsApi, type MerchantProductStats } from '#/api/core/merchant-statistic';

const loading = ref(false);
const stats = ref<MerchantProductStats>();

async function load() {
  loading.value = true;
  try {
    stats.value = await getMerchantProductStatsApi();
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="商品统计" description="按当前店铺统计各状态商品数量。">
    <el-card v-loading="loading" shadow="never">
      <div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-5">
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">商品总数</div><div class="mt-2 text-2xl font-semibold">{{ stats?.total ?? '—' }}</div></div>
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">在售</div><div class="mt-2 text-2xl font-semibold">{{ stats?.on_sale ?? '—' }}</div></div>
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">待审核</div><div class="mt-2 text-2xl font-semibold">{{ stats?.pending_review ?? '—' }}</div></div>
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">草稿</div><div class="mt-2 text-2xl font-semibold">{{ stats?.draft ?? '—' }}</div></div>
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">已下架</div><div class="mt-2 text-2xl font-semibold">{{ stats?.off_sale ?? '—' }}</div></div>
      </div>
      <div class="mt-4"><el-button @click="load">刷新</el-button></div>
    </el-card>
  </Page>
</template>
