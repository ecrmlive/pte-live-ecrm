<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import {
  formatCount,
  formatMoney,
  getPlatformDashboardSummaryApi,
  type PlatformDashboardSummary,
} from '#/api/core/platform-dashboard';
import { formatShanghaiDateTime } from '#/utils/date-time';

const loading = ref(false);
const failed = ref(false);
const updatedAt = ref('');
const dashboard = ref<PlatformDashboardSummary>({
  new_users: { month: 0, today: 0, yesterday: 0 },
  on_sale_product: 0,
  page_views: { month: 0, today: 0, yesterday: 0 },
  paid_order: 0,
  scope: 'all',
  pending_delivery: 0,
  pending_product_audit: 0,
  pending_refund: 0,
  pending_service: 0,
  pending_store_audit: 0,
  store_count: 0,
  store_sales_rank: [],
  store_total: 0,
  today_order_count: 0,
  today_paid_amount: 0,
  today_payer_count: 0,
  visitors: { month: 0, today: 0, yesterday: 0 },
});

async function load() {
  loading.value = true;
  failed.value = false;
  try {
    dashboard.value = await getPlatformDashboardSummaryApi();
  } catch {
    failed.value = true;
  } finally {
    updatedAt.value = formatShanghaiDateTime(new Date());
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="商品统计" description="基于平台 dashboard 汇总接口的商品监管视图。">
    <el-alert v-if="failed" class="mb-4" title="统计数据暂不可用，请刷新重试。" type="warning" :closable="false" />
    <el-row v-loading="loading" :gutter="16">
      <el-col :md="8" :xs="24"><el-card shadow="never"><div class="text-sm text-gray-500">在售商品</div><div class="mt-2 text-2xl font-semibold">{{ formatCount(dashboard.on_sale_product) }}</div></el-card></el-col>
      <el-col :md="8" :xs="24"><el-card shadow="never"><div class="text-sm text-gray-500">待审核商品</div><div class="mt-2 text-2xl font-semibold">{{ formatCount(dashboard.pending_product_audit) }}</div></el-card></el-col>
      <el-col :md="8" :xs="24"><el-card shadow="never"><div class="text-sm text-gray-500">店铺总数</div><div class="mt-2 text-2xl font-semibold">{{ formatCount(dashboard.store_count) }}</div></el-card></el-col>
    </el-row>
    <el-card v-loading="loading" class="mt-4" shadow="never">
      <template #header><strong>店铺销售排行 TOP 10</strong></template>
      <el-table :data="dashboard.store_sales_rank || []" size="small">
        <el-table-column label="店铺" min-width="160" prop="store_name" />
        <el-table-column align="center" label="销量" width="100"><template #default="{ row }">{{ formatCount(row.sale_count) }}</template></el-table-column>
        <el-table-column align="right" label="销售额" width="140"><template #default="{ row }">¥{{ formatMoney(row.sale_amount) }}</template></el-table-column>
      </el-table>
    </el-card>
    <div class="mt-4 flex justify-end text-sm text-gray-500">
      <span class="mr-3">更新时间：{{ updatedAt || '—' }}</span>
      <el-button :loading="loading" text type="primary" @click="load">刷新</el-button>
    </div>
  </Page>
</template>
