<script setup lang="ts">
import type { Component } from 'vue';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ChatDotRound,
  Connection,
  Document,
  Goods,
  Money,
  Setting,
  Shop,
  Ticket,
  User,
  UserFilled,
  View,
} from '@element-plus/icons-vue';
import { ElButton, ElCard, ElEmpty, ElIcon, ElSkeleton, ElTag } from 'element-plus';

import {
  getPlatformDashboardSummaryApi,
  type DashboardMetric,
  type PlatformDashboardSummary,
} from '#/api/core/platform-dashboard';

type MetricCard = {
  icon: Component;
  label: string;
  metric: DashboardMetric;
  tone: 'blue' | 'cyan' | 'purple';
};

const loading = ref(false);
const failed = ref(false);
const updatedAt = ref('');
const dashboard = ref<PlatformDashboardSummary>({
  new_users: { month: 0, today: 0, yesterday: 0 },
  on_sale_product: 0,
  page_views: { month: 0, today: 0, yesterday: 0 },
  paid_order: 0,
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

const metricCards = computed<MetricCard[]>(() => [
  { icon: User, label: '新增用户', metric: dashboard.value.new_users, tone: 'blue' },
  { icon: View, label: '浏览量', metric: dashboard.value.page_views, tone: 'cyan' },
  { icon: UserFilled, label: '访客数', metric: dashboard.value.visitors, tone: 'purple' },
]);

const quickLinks: Array<{ icon: Component; label: string; route: string; tone: string }> = [
  { icon: Goods, label: '商品管理', route: '/product/audit', tone: 'mint' },
  { icon: User, label: '用户管理', route: '/user/group', tone: 'blue' },
  { icon: Document, label: '订单管理', route: '/order/list', tone: 'orange' },
  { icon: Connection, label: '分销管理', route: '/marketing/spread', tone: 'purple' },
  { icon: ChatDotRound, label: '客服管理', route: '/service', tone: 'gold' },
  { icon: Document, label: '文章管理', route: '/content/article', tone: 'cyan' },
  { icon: Ticket, label: '平台优惠券', route: '/marketing/coupon', tone: 'peach' },
  { icon: Setting, label: '系统设置', route: '/setting/admin', tone: 'amber' },
];

const todos = computed(() => [
  { icon: Goods, label: '待审核商品', route: '/product/audit', tone: 'blue', value: dashboard.value.pending_product_audit },
  { icon: Shop, label: '待审核商户入驻', route: '/merchant/audit', tone: 'orange', value: dashboard.value.pending_store_audit },
  { icon: Document, label: '待发货订单', route: '/order/list', tone: 'purple', value: dashboard.value.pending_delivery },
  { icon: Money, label: '待处理退款', route: '/order/refund', tone: 'cyan', value: dashboard.value.pending_refund },
  { icon: ChatDotRound, label: '待处理用户咨询', route: '/service', tone: 'green', value: dashboard.value.pending_service },
]);

const formatCount = (value: number | null | undefined) => Number(value || 0).toLocaleString('zh-CN');
const formatMoney = (value: number | null | undefined) => Number(value || 0).toLocaleString('zh-CN', {
  maximumFractionDigits: 2,
  minimumFractionDigits: 2,
});

async function loadDashboard() {
  loading.value = true;
  failed.value = false;
  try {
    dashboard.value = await getPlatformDashboardSummaryApi();
  } catch {
    failed.value = true;
  } finally {
    updatedAt.value = new Date().toLocaleString('zh-CN', { hour12: false });
    loading.value = false;
  }
}

onMounted(() => void loadDashboard());
</script>

<template>
  <Page title="控制台" description="当前账号的数据范围与待办事项概览" content-class="platform-dashboard-page">
    <div class="platform-dashboard" :class="{ 'is-loading': loading }">
      <p v-if="failed" class="load-error">经营数据暂不可用，请刷新重试；不会使用演示数据替代真实统计。</p>

      <section class="metric-grid" aria-label="今日经营指标">
        <ElCard v-for="card in metricCards" :key="card.label" shadow="never" class="metric-card">
          <template v-if="loading"><ElSkeleton :rows="3" animated /></template>
          <template v-else>
            <div class="metric-card__header">
              <span>{{ card.label }}</span>
              <ElTag size="small" effect="plain">今日</ElTag>
            </div>
            <strong>{{ formatCount(card.metric.today) }}</strong>
            <div class="metric-card__compare">
              <span>昨日 <b>{{ formatCount(card.metric.yesterday) }}</b></span>
              <span :class="`metric-card__month metric-card__month--${card.tone}`">本月 {{ formatCount(card.metric.month) }}</span>
            </div>
            <div class="metric-card__footer">本月累计 {{ formatCount(card.metric.month) }}</div>
          </template>
        </ElCard>

        <ElCard shadow="never" class="metric-card metric-card--store">
          <template v-if="loading"><ElSkeleton :rows="3" animated /></template>
          <template v-else>
            <div class="metric-card__header">
              <span>店铺数</span>
              <ElTag size="small" effect="plain">当前</ElTag>
            </div>
            <strong>{{ formatCount(dashboard.store_count) }}</strong>
            <div class="metric-card__compare"><span>已入驻店铺</span></div>
            <div class="metric-card__footer">当前营业 {{ formatCount(dashboard.store_count) }}</div>
          </template>
        </ElCard>
      </section>

      <section class="quick-grid" aria-label="常用功能">
        <RouterLink v-for="item in quickLinks" :key="item.route" :to="item.route" class="quick-card">
          <span class="quick-card__icon" :class="`quick-card__icon--${item.tone}`">
            <ElIcon><component :is="item.icon" /></ElIcon>
          </span>
          <span>{{ item.label }}</span>
        </RouterLink>
      </section>

      <section class="content-grid">
        <ElCard shadow="never" class="dashboard-panel todo-panel">
          <template #header><strong>待办事项</strong></template>
          <div class="todo-grid">
            <RouterLink v-for="item in todos" :key="item.label" :to="item.route" class="todo-item">
              <span class="todo-item__icon" :class="`todo-item__icon--${item.tone}`">
                <ElIcon><component :is="item.icon" /></ElIcon>
              </span>
              <span class="todo-item__text">
                <small>{{ item.label }}</small>
                <b>{{ formatCount(item.value) }}</b>
              </span>
            </RouterLink>
          </div>
        </ElCard>

        <ElCard shadow="never" class="dashboard-panel rank-panel">
          <template #header>
            <div class="panel-header">
              <strong>店铺销售情况排行 TOP 10</strong>
              <ElTag effect="plain" type="primary">本月</ElTag>
            </div>
          </template>
          <ElEmpty v-if="!loading && dashboard.store_sales_rank.length === 0" description="本月暂无已支付订单" :image-size="76" />
          <el-table v-else v-loading="loading" class="rank-table" :data="dashboard.store_sales_rank" size="small">
            <el-table-column align="center" label="排名" width="68">
              <template #default="{ $index }"><span class="rank-badge" :class="{ 'rank-badge--top': $index < 3 }">{{ $index + 1 }}</span></template>
            </el-table-column>
            <el-table-column label="店铺名称" min-width="160" prop="store_name" show-overflow-tooltip />
            <el-table-column align="center" label="关注人数" min-width="100">
              <template #default="{ row }">{{ formatCount(row.follower_count) }}</template>
            </el-table-column>
            <el-table-column align="center" label="销量" min-width="82">
              <template #default="{ row }">{{ formatCount(row.sale_count) }}</template>
            </el-table-column>
            <el-table-column align="right" label="销售金额" min-width="128">
              <template #default="{ row }">¥{{ formatMoney(row.sale_amount) }}</template>
            </el-table-column>
          </el-table>
        </ElCard>
      </section>

      <section class="bottom-grid">
        <ElCard shadow="never" class="dashboard-panel amount-panel">
          <template #header><strong>当日订单金额</strong></template>
          <div class="amount-main">¥{{ formatMoney(dashboard.today_paid_amount) }}</div>
          <p>只统计已支付、履约中、已发货和已完成的订单。</p>
        </ElCard>

        <ElCard shadow="never" class="dashboard-panel order-panel">
          <template #header><strong>数据统计</strong></template>
          <div class="order-stats">
            <div><span>当日订单数</span><b>{{ formatCount(dashboard.today_order_count) }}</b></div>
            <div><span>当日支付人数</span><b>{{ formatCount(dashboard.today_payer_count) }}</b></div>
            <div><span>在售商品</span><b>{{ formatCount(dashboard.on_sale_product) }}</b></div>
            <div><span>累计已支付订单</span><b>{{ formatCount(dashboard.paid_order) }}</b></div>
          </div>
        </ElCard>
      </section>

      <footer class="dashboard-footer">
        <span>统计更新时间：{{ updatedAt || '—' }}</span>
        <ElButton text type="primary" :loading="loading" @click="loadDashboard">刷新数据</ElButton>
      </footer>
    </div>
  </Page>
</template>

<style scoped>
.platform-dashboard { display: grid; gap: 16px; min-width: 0; }
.load-error { margin: 0; padding: 10px 14px; border: 1px solid var(--el-color-warning-light-5); border-radius: 8px; color: var(--el-color-warning-dark-2); background: var(--el-color-warning-light-9); font-size: 13px; }
.metric-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 16px; }
.metric-card { min-height: 186px; overflow: hidden; border-color: var(--el-border-color-lighter); }
.metric-card :deep(.el-card__body) { display: flex; min-height: 186px; padding: 22px 26px 0; flex-direction: column; }
.metric-card__header { display: flex; align-items: center; justify-content: space-between; color: var(--el-text-color-primary); font-size: 16px; font-weight: 600; }
.metric-card__header :deep(.el-tag) { border-color: var(--el-color-primary-light-5); color: var(--el-color-primary); }
.metric-card strong { margin-top: 20px; color: var(--el-text-color-primary); font-size: 36px; font-weight: 600; line-height: 1; }
.metric-card__compare { display: flex; flex-wrap: wrap; gap: 12px; margin-top: 12px; color: var(--el-text-color-secondary); font-size: 14px; }
.metric-card__compare b { margin-left: 6px; color: var(--el-text-color-primary); font-weight: 500; }
.metric-card__month { font-size: 13px; }
.metric-card__month--blue { color: #4f87ff; }.metric-card__month--cyan { color: #16b8b8; }.metric-card__month--purple { color: #9667e8; }
.metric-card__footer { margin: auto -26px 0; padding: 14px 26px; border-top: 1px solid var(--el-border-color-lighter); color: var(--el-text-color-secondary); font-size: 13px; }
.quick-grid { display: grid; grid-template-columns: repeat(8, minmax(0, 1fr)); gap: 16px; }
.quick-card { display: grid; min-height: 138px; place-content: center; gap: 14px; border: 1px solid var(--el-border-color-lighter); border-radius: 8px; color: var(--el-text-color-primary); background: var(--el-bg-color); text-align: center; text-decoration: none; transition: border-color .18s ease, box-shadow .18s ease, transform .18s ease; }
.quick-card:hover { border-color: var(--el-color-primary-light-5); box-shadow: 0 6px 18px rgb(36 86 180 / 10%); color: var(--el-color-primary); transform: translateY(-2px); }
.quick-card__icon { display: grid; width: 34px; height: 34px; margin: auto; place-items: center; font-size: 28px; }.quick-card__icon--mint { color: #32bea2; }.quick-card__icon--blue { color: #4f87ff; }.quick-card__icon--orange { color: #f3975d; }.quick-card__icon--purple { color: #9b73e7; }.quick-card__icon--gold { color: #df9d0a; }.quick-card__icon--cyan { color: #32bdbd; }.quick-card__icon--peach { color: #ee9a7d; }.quick-card__icon--amber { color: #e8a80d; }
.content-grid { display: grid; grid-template-columns: minmax(360px, 1fr) minmax(520px, 1.04fr); gap: 16px; }
.dashboard-panel { border-color: var(--el-border-color-lighter); }.dashboard-panel :deep(.el-card__header) { min-height: 56px; padding: 17px 22px; border-bottom-color: var(--el-border-color-lighter); font-size: 16px; }.dashboard-panel :deep(.el-card__body) { padding: 20px 22px; }
.todo-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 30px 22px; padding: 8px 0; }.todo-item { display: flex; align-items: center; gap: 12px; min-width: 0; color: var(--el-text-color-primary); text-decoration: none; }.todo-item:hover small { color: var(--el-color-primary); }.todo-item__icon { display: grid; flex: 0 0 auto; width: 52px; height: 52px; place-items: center; border-radius: 2px; font-size: 26px; }.todo-item__icon--blue { color: #4688ef; background: #edf5ff; }.todo-item__icon--orange { color: #ef851a; background: #fff5e8; }.todo-item__icon--purple { color: #9871ed; background: #f6efff; }.todo-item__icon--cyan { color: #20bbb9; background: #e9fbfa; }.todo-item__icon--green { color: #65ba38; background: #f1fae9; }.todo-item__text { display: grid; gap: 5px; min-width: 0; }.todo-item__text small { overflow: hidden; color: var(--el-text-color-secondary); font-size: 14px; text-overflow: ellipsis; white-space: nowrap; }.todo-item__text b { font-size: 24px; font-weight: 600; line-height: 1; }
.panel-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; }.rank-table { width: 100%; }.rank-table :deep(.el-table__cell) { height: 48px; padding: 7px 0; }.rank-badge { display: inline-grid; width: 24px; height: 24px; place-items: center; border-radius: 50%; color: var(--el-text-color-secondary); background: var(--el-fill-color-light); font-size: 12px; }.rank-badge--top { color: #fff; background: #d9ac4c; }
.bottom-grid { display: grid; grid-template-columns: minmax(300px, .97fr) minmax(500px, 1.03fr); gap: 16px; }.amount-main { color: var(--el-text-color-primary); font-size: 32px; font-weight: 600; line-height: 1.1; }.amount-panel p { margin: 16px 0 0; color: var(--el-text-color-secondary); font-size: 13px; }.order-stats { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 24px 48px; }.order-stats div { display: grid; gap: 8px; }.order-stats span { color: var(--el-text-color-secondary); font-size: 13px; }.order-stats b { color: var(--el-text-color-primary); font-size: 23px; font-weight: 600; }.dashboard-footer { display: flex; align-items: center; justify-content: flex-end; gap: 12px; color: var(--el-text-color-secondary); font-size: 12px; }
@media (max-width: 1440px) { .quick-grid { grid-template-columns: repeat(4, minmax(0, 1fr)); }.content-grid { grid-template-columns: 1fr; }.metric-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
@media (max-width: 860px) { .todo-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }.bottom-grid { grid-template-columns: 1fr; } }
@media (max-width: 600px) { .metric-grid, .quick-grid, .todo-grid, .order-stats { grid-template-columns: 1fr; }.metric-card { min-height: 172px; }.quick-card { min-height: 112px; }.dashboard-panel :deep(.el-card__body) { padding: 16px; }.dashboard-footer { align-items: flex-end; flex-direction: column; } }
</style>
