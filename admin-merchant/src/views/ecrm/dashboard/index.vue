<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElCard, ElSkeleton, ElTag } from 'element-plus';

import { getMerchantDashboardSummaryApi } from '#/api/core/merchant-dashboard';

type Stat = {
  label: string;
  route: string;
  tone: 'blue' | 'green' | 'orange' | 'red';
  value: number | null;
  formatter?: (value: number) => string;
};

const loading = ref(false);
const failed = ref(false);
const updatedAt = ref('');
const stats = ref<Stat[]>([
  { label: '商品总数', route: '/product/list', tone: 'blue', value: null },
  { label: '已支付订单', route: '/order/list', tone: 'green', value: null },
  { label: '可用余额', route: '/finance/balance', tone: 'orange', value: null, formatter: (value) => `¥${value.toFixed(2)}` },
  { label: '待处理退款', route: '/order/refund', tone: 'red', value: null },
]);

const quickLinks = [
  { label: '发布商品', route: '/product/list' },
  { label: '订单处理', route: '/order/list' },
  { label: '售后处理', route: '/order/refund' },
  { label: '优惠券', route: '/marketing/coupon' },
  { label: '物流模板', route: '/config/freight/shippingTemplates' },
  { label: '店铺装修', route: '/devise/diy/list' },
];

const todos = computed(() => [
  { label: '已支付订单', route: '/order/list', value: stats.value[1]?.value },
  { label: '待处理退款', route: '/order/refund', value: stats.value[3]?.value },
  { label: '已录入商品', route: '/product/list', value: stats.value[0]?.value },
]);

async function loadDashboard() {
  loading.value = true;
  failed.value = false;
  try {
    const summary = await getMerchantDashboardSummaryApi();
    const values = [summary.product_total, summary.paid_order, summary.available_balance, summary.pending_refund];
    stats.value = stats.value.map((item, index) => ({ ...item, value: Number(values[index] || 0) }));
  } catch {
    stats.value = stats.value.map((item) => ({ ...item, value: null }));
    failed.value = true;
  }
  updatedAt.value = new Date().toLocaleString('zh-CN', { hour12: false });
  loading.value = false;
}

function formatStat(stat: Stat) {
  if (stat.value === null) return '—';
  return stat.formatter ? stat.formatter(stat.value) : stat.value.toLocaleString('zh-CN');
}

onMounted(() => void loadDashboard());
</script>

<template>
  <Page title="控制台" description="当前店铺的商品、订单、资金与售后工作台">
    <div class="dashboard">
      <header class="dashboard-head">
        <div>
          <h2>店铺经营概览</h2>
          <p>数据仅在当前店铺 AppId 范围内查询，店员权限仍由服务端校验。</p>
        </div>
        <div class="dashboard-actions">
          <span v-if="updatedAt">更新于 {{ updatedAt }}</span>
          <ElButton :loading="loading" @click="loadDashboard">刷新数据</ElButton>
        </div>
      </header>

      <p v-if="failed" class="load-error">经营数据暂不可用，请检查店铺 API、X-AppId 和当前登录状态。</p>

      <section class="stat-grid" aria-label="店铺经营统计">
        <ElCard v-for="stat in stats" :key="stat.label" shadow="never" class="stat-card" :class="stat.tone">
          <RouterLink :to="stat.route" class="stat-link">
            <span>{{ stat.label }}</span>
            <ElSkeleton v-if="loading" :rows="1" animated />
            <strong v-else>{{ formatStat(stat) }}</strong>
            <small>查看详情 →</small>
          </RouterLink>
        </ElCard>
      </section>

      <section class="panel-grid">
        <ElCard shadow="never" class="panel">
          <template #header><strong>常用功能</strong></template>
          <div class="quick-grid">
            <RouterLink v-for="item in quickLinks" :key="item.route" :to="item.route">{{ item.label }}</RouterLink>
          </div>
        </ElCard>
        <ElCard shadow="never" class="panel">
          <template #header><strong>待办事项</strong></template>
          <RouterLink v-for="item in todos" :key="item.label" :to="item.route" class="todo-row">
            <span>{{ item.label }}</span>
            <ElTag type="danger" effect="plain">{{ item.value === null ? '—' : item.value }}</ElTag>
          </RouterLink>
        </ElCard>
      </section>
    </div>
  </Page>
</template>

<style scoped>
.dashboard { display: grid; gap: 16px; padding: 4px; }
.dashboard-head { display: flex; align-items: center; justify-content: space-between; gap: 24px; padding: 12px 4px; }
.dashboard-head h2 { margin: 0; font-size: 20px; font-weight: 600; }
.dashboard-head p, .dashboard-actions span { margin: 6px 0 0; color: var(--el-text-color-secondary); font-size: 13px; }
.dashboard-actions { display: flex; align-items: center; gap: 12px; }
.load-error { margin: 0; padding: 10px 12px; border: 1px solid var(--el-color-warning-light-5); color: var(--el-color-warning-dark-2); background: var(--el-color-warning-light-9); font-size: 13px; }
.stat-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 16px; }
.stat-card { overflow: hidden; border-left: 4px solid var(--el-color-primary); }
.stat-card.green { border-left-color: var(--el-color-success); }
.stat-card.orange { border-left-color: var(--el-color-warning); }
.stat-card.red { border-left-color: var(--el-color-danger); }
.stat-link { display: grid; gap: 10px; color: inherit; text-decoration: none; }
.stat-link span, .stat-link small { color: var(--el-text-color-secondary); font-size: 13px; }
.stat-link strong { font-size: 30px; line-height: 1; }
.panel-grid { display: grid; grid-template-columns: 1.2fr .8fr; gap: 16px; }
.panel :deep(.el-card__header) { padding: 14px 18px; }
.quick-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; }
.quick-grid a { padding: 12px; border: 1px solid var(--el-border-color-lighter); color: var(--el-text-color-primary); text-align: center; text-decoration: none; }
.quick-grid a:hover { border-color: var(--el-color-primary); color: var(--el-color-primary); }
.todo-row { display: flex; align-items: center; justify-content: space-between; padding: 12px 0; border-bottom: 1px solid var(--el-border-color-lighter); color: var(--el-text-color-primary); text-decoration: none; }
.todo-row:last-child { border-bottom: 0; }
@media (max-width: 900px) { .stat-grid, .panel-grid { grid-template-columns: 1fr 1fr; } .panel-grid { grid-template-columns: 1fr; } }
@media (max-width: 600px) { .dashboard-head { align-items: flex-start; flex-direction: column; } .stat-grid, .quick-grid { grid-template-columns: 1fr; } }
</style>
