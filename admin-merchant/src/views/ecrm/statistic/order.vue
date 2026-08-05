<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { getMerchantOrderStatsApi, type MerchantOrderStats } from '#/api/core/merchant-statistic';

const loading = ref(false);
const stats = ref<MerchantOrderStats>();

async function load() {
  loading.value = true;
  try {
    stats.value = await getMerchantOrderStatsApi();
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="订单统计" description="按当前店铺统计有效订单与售后待办。">
    <el-card v-loading="loading" shadow="never">
      <div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-5">
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">已支付订单</div><div class="mt-2 text-2xl font-semibold">{{ stats?.paid_order ?? '—' }}</div></div>
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">待发货</div><div class="mt-2 text-2xl font-semibold">{{ stats?.pending_ship ?? '—' }}</div></div>
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">已发货</div><div class="mt-2 text-2xl font-semibold">{{ stats?.shipped ?? '—' }}</div></div>
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">已完成</div><div class="mt-2 text-2xl font-semibold">{{ stats?.completed ?? '—' }}</div></div>
        <div class="rounded border p-4"><div class="text-sm text-muted-foreground">待处理退款</div><div class="mt-2 text-2xl font-semibold text-danger">{{ stats?.pending_refund ?? '—' }}</div></div>
      </div>
      <div class="mt-4"><el-button @click="load">刷新</el-button></div>
    </el-card>
  </Page>
</template>
