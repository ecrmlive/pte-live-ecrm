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
  <Page title="订单统计" description="基于平台 dashboard 汇总接口的订单监管视图；不含支付密钥或用户隐私字段。">
    <el-alert v-if="failed" class="mb-4" title="统计数据暂不可用，请刷新重试。" type="warning" :closable="false" />
    <el-row v-loading="loading" :gutter="16">
      <el-col :md="8" :xs="24"><el-card shadow="never"><div class="text-sm text-gray-500">当日订单数</div><div class="mt-2 text-2xl font-semibold">{{ formatCount(dashboard.today_order_count) }}</div></el-card></el-col>
      <el-col :md="8" :xs="24"><el-card shadow="never"><div class="text-sm text-gray-500">当日支付金额</div><div class="mt-2 text-2xl font-semibold">¥{{ formatMoney(dashboard.today_paid_amount) }}</div></el-card></el-col>
      <el-col :md="8" :xs="24"><el-card shadow="never"><div class="text-sm text-gray-500">当日支付人数</div><div class="mt-2 text-2xl font-semibold">{{ formatCount(dashboard.today_payer_count) }}</div></el-card></el-col>
      <el-col class="mt-4" :md="8" :xs="24"><el-card shadow="never"><div class="text-sm text-gray-500">累计已支付订单</div><div class="mt-2 text-2xl font-semibold">{{ formatCount(dashboard.paid_order) }}</div></el-card></el-col>
      <el-col class="mt-4" :md="8" :xs="24"><el-card shadow="never"><div class="text-sm text-gray-500">待发货订单</div><div class="mt-2 text-2xl font-semibold">{{ formatCount(dashboard.pending_delivery) }}</div></el-card></el-col>
      <el-col class="mt-4" :md="8" :xs="24"><el-card shadow="never"><div class="text-sm text-gray-500">待处理退款</div><div class="mt-2 text-2xl font-semibold">{{ formatCount(dashboard.pending_refund) }}</div></el-card></el-col>
    </el-row>
    <div class="mt-4 flex justify-end text-sm text-gray-500">
      <span class="mr-3">更新时间：{{ updatedAt || '—' }}</span>
      <el-button :loading="loading" text type="primary" @click="load">刷新</el-button>
    </div>
  </Page>
</template>
