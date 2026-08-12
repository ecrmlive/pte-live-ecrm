<script setup lang="ts">
import type { Component } from 'vue';
import type { EchartsUIType } from '@vben/plugins/echarts';

import { computed, nextTick, onMounted, ref, watch } from 'vue';

import { Page } from '@vben/common-ui';
import { EchartsUI, useEcharts } from '@vben/plugins/echarts';
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
  Van,
  Wallet,
} from '@element-plus/icons-vue';
import { ElButton, ElCard, ElEmpty, ElIcon, ElRadioButton, ElRadioGroup, ElSkeleton, ElTag } from 'element-plus';

import {
  emptyDashboardSummary,
  getPlatformDashboardSummaryApi,
  getPlatformDealApi,
  getPlatformMerchantTopApi,
  getPlatformUserTrendApi,
  type DashboardMetric,
  type DashboardPeriod,
  type PlatformDashboardSummary,
  type SparkStat,
  type StoreSalesRankRow,
} from '#/api/core/platform-dashboard';
import { getUserInfoApi } from '#/api/core/auth';
import { formatShanghaiDateTime } from '#/utils/date-time';

const loading = ref(false);
const failed = ref(false);
const updatedAt = ref('');
const roles = ref<string[]>([]);
const dashboard = ref<PlatformDashboardSummary>(emptyDashboardSummary());

const rankPeriod = ref<DashboardPeriod>('month');
const rankLoading = ref(false);
const rankRows = ref<StoreSalesRankRow[]>([]);

const userPeriod = ref<'7d' | '30d' | 'month'>('30d');
const dealPeriod = ref<DashboardPeriod>('month');
const ratioPeriod = ref<DashboardPeriod>('month');
const ratioMode = ref<'amount' | 'users'>('amount');

const amountChartRef = ref<EchartsUIType>();
const spark0Ref = ref<EchartsUIType>();
const spark1Ref = ref<EchartsUIType>();
const userChartRef = ref<EchartsUIType>();
const ratioChartRef = ref<EchartsUIType>();

const { renderEcharts: renderAmountChart } = useEcharts(amountChartRef);
const { renderEcharts: renderSpark0 } = useEcharts(spark0Ref);
const { renderEcharts: renderSpark1 } = useEcharts(spark1Ref);
const { renderEcharts: renderUserChart } = useEcharts(userChartRef);
const { renderEcharts: renderRatioChart } = useEcharts(ratioChartRef);
const sparkRenderers = [renderSpark0, renderSpark1];

const isPlatform = computed(() => roles.value.includes('platform'));
const isMerchantOrRegion = computed(
  () => roles.value.includes('merchant') || roles.value.includes('region'),
);
const isCustomerService = computed(() => roles.value.includes('customer_service'));

type MetricCard = {
  footerLabel: string;
  footerValue: number;
  label: string;
  metric: DashboardMetric;
};

const metricCards = computed<MetricCard[]>(() => {
  const d = dashboard.value;
  return [
    {
      footerLabel: '本月新增用户',
      footerValue: d.new_users.month,
      label: d.scope === 'all' ? '新增用户' : '下单用户',
      metric: d.new_users,
    },
    {
      footerLabel: '本月浏览量',
      footerValue: d.page_views.month,
      label: '浏览量',
      metric: d.page_views,
    },
    {
      footerLabel: '本月访客',
      footerValue: d.visitors.month,
      label: '访客数',
      metric: d.visitors,
    },
    {
      footerLabel: '当前店铺',
      footerValue: d.store_count,
      label: '店铺数',
      metric: d.stores?.today !== undefined ? d.stores : { month: d.store_count, today: 0, yesterday: 0, week_ratio: 0 },
    },
  ];
});

const platformQuickLinks: Array<{ icon: Component; label: string; route: string; tone: string }> = [
  { icon: Goods, label: '商品管理', route: '/product/audit', tone: 'mint' },
  { icon: User, label: '用户管理', route: '/user/list', tone: 'blue' },
  { icon: Document, label: '订单管理', route: '/order/list', tone: 'orange' },
  { icon: Connection, label: '分销管理', route: '/marketing/spread', tone: 'purple' },
  { icon: Van, label: '服务配置', route: '/service/settings', tone: 'gold' },
  { icon: Document, label: '文章管理', route: '/cms/article', tone: 'cyan' },
  { icon: Ticket, label: '优惠券', route: '/marketing/coupon', tone: 'peach' },
  { icon: Setting, label: '系统设置', route: '/setting/admin', tone: 'amber' },
];

const quickLinks = computed(() => {
  if (isPlatform.value) return platformQuickLinks;
  if (isCustomerService.value) return platformQuickLinks.filter((item) => item.route === '/user/list');
  if (isMerchantOrRegion.value) {
    return platformQuickLinks.filter((item) => ['/product/audit', '/order/list'].includes(item.route));
  }
  return [];
});

const allTodos = computed(() => {
  const d = dashboard.value;
  return [
    { icon: Goods, label: '待审核商品', route: '/product/audit', tone: 'blue', value: d.pending_product_audit },
    { icon: Shop, label: '待审核分销礼包', route: '/promoter/gift', tone: 'cyan', value: d.pending_spread_gift },
    { icon: Shop, label: '待审核商户入驻', route: '/merchant/audit', tone: 'orange', value: d.pending_store_audit },
    { icon: Wallet, label: '待审核提现', route: '/accounts/withdraw', tone: 'green', value: d.pending_withdraw },
    { icon: Money, label: '待审核转账', route: '/accounts/transferRecord', tone: 'purple', value: d.pending_transfer },
    { icon: ChatDotRound, label: '待审核社区内容', route: '/content/community', tone: 'teal', value: d.pending_community },
    { icon: Document, label: '待退款订单', route: '/order/refund', tone: 'red', value: d.pending_refund },
    { icon: ChatDotRound, label: '待处理用户反馈', route: '/user/feedback/list', tone: 'lime', value: d.pending_feedback },
    { icon: Van, label: '待发货积分订单', route: '/marketing/points', tone: 'amber', value: d.pending_integral_ship },
  ];
});

const todos = computed(() => {
  if (isPlatform.value) return allTodos.value;
  if (isCustomerService.value) {
    return allTodos.value.filter((item) =>
      ['/user/feedback/list', '/content/community'].includes(item.route),
    );
  }
  if (isMerchantOrRegion.value) {
    return allTodos.value.filter((item) =>
      ['/product/audit', '/order/refund', '/order/list'].includes(item.route),
    );
  }
  return [];
});

const sparkStats = computed(() => {
  const s = dashboard.value.order_stats;
  return [
    { hasChart: true, key: 'today_order', label: '当日订单数', ratioLabel: '日同比', stat: s?.today_order_count },
    { hasChart: true, key: 'today_payer', label: '当日支付人数', ratioLabel: '日同比', stat: s?.today_payer_count },
    { hasChart: false, key: 'month_order', label: '当月订单数', ratioLabel: '月同比', stat: s?.month_order_count },
    { hasChart: false, key: 'month_payer', label: '当月支付人数', ratioLabel: '月同比', stat: s?.month_payer_count },
  ];
});

const dealFunnel = computed(() => dashboard.value.deal_funnel);

const periodTabs: Array<{ label: string; value: DashboardPeriod }> = [
  { label: '近7天', value: '7d' },
  { label: '近30天', value: '30d' },
  { label: '本月', value: 'month' },
  { label: '本年', value: 'year' },
];

const userPeriodTabs: Array<{ label: string; value: '7d' | '30d' | 'month' }> = [
  { label: '近7天', value: '7d' },
  { label: '近30天', value: '30d' },
  { label: '本月', value: 'month' },
];

const formatCount = (value: number | null | undefined) => Number(value || 0).toLocaleString('zh-CN');
const formatMoney = (value: number | null | undefined) =>
  Number(value || 0).toLocaleString('zh-CN', { maximumFractionDigits: 2, minimumFractionDigits: 2 });

function formatWeekRatio(ratio: number | null | undefined) {
  const n = Number(ratio || 0);
  const sign = n > 0 ? '+' : '';
  return `${sign}${n.toFixed(2)}%`;
}

function formatSparkRatio(ratio: number | null | undefined) {
  return `${Number(ratio || 0).toFixed(0)}%`;
}

function ratioClass(ratio: number | null | undefined) {
  const n = Number(ratio || 0);
  if (n === 0) return 'delta-flat';
  return n > 0 ? 'delta-up' : 'delta-down';
}

/** 数据统计：0% 也按上涨红▲，下跌绿▼ */
function sparkRatioClass(ratio: number | null | undefined) {
  return Number(ratio || 0) >= 0 ? 'delta-up' : 'delta-down';
}

function sparkOption(stat?: SparkStat) {
  const raw = Array.isArray(stat?.spark) ? stat!.spark : [];
  const data = Array.from({ length: 24 }, (_, i) => Number(raw[i] ?? 0));
  const hours = Array.from({ length: 24 }, (_, i) => `${String(i).padStart(2, '0')}:00`);
  const maxVal = Math.max(...data, 0);
  return {
    animation: false,
    color: ['#409eff'],
    grid: { bottom: 18, containLabel: false, left: 28, right: 8, top: 6 },
    series: [
      {
        data,
        itemStyle: { color: '#409eff' },
        lineStyle: { color: '#409eff', width: 1.5 },
        name: '今天',
        showSymbol: false,
        type: 'line',
      },
    ],
    tooltip: {
      formatter: (params: any) => {
        const p = Array.isArray(params) ? params[0] : params;
        return `${p?.axisValue ?? ''} ${p?.data ?? 0}`;
      },
      trigger: 'axis',
    },
    xAxis: {
      axisLabel: {
        color: '#c0c4cc',
        fontSize: 11,
        interval: 5,
      },
      axisLine: { lineStyle: { color: '#ebeef5' } },
      axisTick: { show: false },
      boundaryGap: false,
      data: hours,
      type: 'category',
    },
    yAxis: {
      axisLabel: { color: '#c0c4cc', fontSize: 11, margin: 6 },
      axisLine: { show: false },
      axisTick: { show: false },
      min: 0,
      minInterval: maxVal === 0 ? 0.2 : undefined,
      splitLine: { lineStyle: { color: '#f0f2f5', type: 'solid' } },
      splitNumber: 4,
      type: 'value',
    },
  };
}

function renderAmount() {
  const hours = dashboard.value.today_order_hours || [];
  const x = hours.length ? hours.map((h) => h.hour) : Array.from({ length: 24 }, (_, i) => `${String(i).padStart(2, '0')}:00`);
  const today = hours.length ? hours.map((h) => h.today_amount) : Array.from({ length: 24 }, () => 0);
  const yesterday = hours.length ? hours.map((h) => h.yesterday_amount) : Array.from({ length: 24 }, () => 0);
  renderAmountChart({
    color: ['#409eff', '#67c23a'],
    grid: { bottom: 28, containLabel: true, left: 12, right: 16, top: 36 },
    legend: { data: ['今天', '昨天'], left: 0, top: 0 },
    series: [
      { data: today, name: '今天', showSymbol: false, smooth: true, type: 'line' },
      { data: yesterday, name: '昨天', showSymbol: false, smooth: true, type: 'line' },
    ],
    tooltip: { trigger: 'axis' },
    xAxis: { boundaryGap: false, data: x, type: 'category' },
    yAxis: { min: 0, splitNumber: 5, type: 'value' },
  });
}

function renderSparks() {
  sparkStats.value
    .filter((item) => item.hasChart)
    .forEach((item, index) => {
      sparkRenderers[index]?.(sparkOption(item.stat));
    });
}

function renderUserTrend() {
  const list = dashboard.value.user_trend || [];
  const days = list.map((i) => i.day);
  renderUserChart({
    color: ['#409eff', '#e6a23c', '#9b73e7'],
    grid: { bottom: 28, containLabel: true, left: 16, right: 24, top: 48 },
    legend: { data: ['新用户', '访问用户', '累计用户'], top: 8 },
    series: [
      { data: list.map((i) => i.new_users), name: '新用户', showSymbol: false, smooth: true, type: 'line' },
      { data: list.map((i) => i.visit_users), name: '访问用户', showSymbol: false, smooth: true, type: 'line' },
      {
        data: list.map((i) => i.total_users),
        name: '累计用户',
        showSymbol: false,
        smooth: true,
        type: 'line',
        yAxisIndex: 1,
      },
    ],
    tooltip: { trigger: 'axis' },
    xAxis: { boundaryGap: false, data: days.length ? days : ['—'], type: 'category' },
    yAxis: [
      { name: '新/访', splitNumber: 4, type: 'value' },
      { name: '累计用户', splitNumber: 4, type: 'value' },
    ],
  });
}

function renderRatio() {
  const r = dashboard.value.deal_ratio;
  const newVal = ratioMode.value === 'amount' ? r?.new_amount || 0 : r?.new_users || 0;
  const oldVal = ratioMode.value === 'amount' ? r?.old_amount || 0 : r?.old_users || 0;
  renderRatioChart({
    color: ['#409eff', '#e6a23c'],
    legend: {
      data: ['新用户', '老用户'],
      icon: 'rect',
      itemGap: 16,
      itemHeight: 10,
      itemWidth: 10,
      left: 8,
      orient: 'vertical',
      textStyle: { color: '#606266', fontSize: 13 },
      top: 'middle',
    },
    series: [
      {
        center: ['58%', '52%'],
        data: [
          { name: '新用户', value: newVal },
          { name: '老用户', value: oldVal },
        ],
        label: { show: false },
        radius: ['46%', '68%'],
        type: 'pie',
      },
    ],
    tooltip: {
      formatter: (p: any) => `${p.name}: ${ratioMode.value === 'amount' ? `¥${formatMoney(p.value)}` : formatCount(p.value)}`,
    },
  });
}

async function renderAllCharts() {
  await nextTick();
  renderAmount();
  renderSparks();
  renderUserTrend();
  renderRatio();
}

async function loadDashboard() {
  loading.value = true;
  failed.value = false;
  try {
    const data = await getPlatformDashboardSummaryApi();
    dashboard.value = {
      ...emptyDashboardSummary(),
      ...data,
      deal_funnel: { ...emptyDashboardSummary().deal_funnel, ...(data.deal_funnel || {}) },
      deal_ratio: { ...emptyDashboardSummary().deal_ratio, ...(data.deal_ratio || {}) },
      order_stats: {
        month_order_count: { ...emptySpark(), ...(data.order_stats?.month_order_count || {}) },
        month_payer_count: { ...emptySpark(), ...(data.order_stats?.month_payer_count || {}) },
        today_order_count: { ...emptySpark(), ...(data.order_stats?.today_order_count || {}) },
        today_payer_count: { ...emptySpark(), ...(data.order_stats?.today_payer_count || {}) },
      },
      store_sales_rank: data.store_sales_rank || [],
      stores: { ...emptyDashboardSummary().stores, ...(data.stores || {}) },
      today_order_hours: data.today_order_hours || [],
      user_trend: data.user_trend || [],
    };
    rankRows.value = dashboard.value.store_sales_rank || [];
    rankPeriod.value = 'month';
    await renderAllCharts();
  } catch {
    failed.value = true;
  } finally {
    updatedAt.value = formatShanghaiDateTime(new Date());
    loading.value = false;
  }
}

function emptySpark(): SparkStat {
  return { ratio: 0, spark: [], value: 0 };
}

async function loadRank(period: DashboardPeriod) {
  rankPeriod.value = period;
  rankLoading.value = true;
  try {
    const res = await getPlatformMerchantTopApi(period);
    rankRows.value = res.list || [];
  } catch {
    rankRows.value = [];
  } finally {
    rankLoading.value = false;
  }
}

async function loadUserTrend(period: '7d' | '30d' | 'month') {
  userPeriod.value = period;
  try {
    const res = await getPlatformUserTrendApi(period);
    dashboard.value.user_trend = res.list || [];
    await nextTick();
    renderUserTrend();
  } catch {
    /* keep previous */
  }
}

async function loadDeal(period: DashboardPeriod, target: 'deal' | 'ratio' = 'deal') {
  if (target === 'deal') dealPeriod.value = period;
  else ratioPeriod.value = period;
  try {
    const res = await getPlatformDealApi(period);
    if (target === 'deal') {
      dashboard.value.deal_funnel = { ...emptyDashboardSummary().deal_funnel, ...(res.funnel || {}) };
    } else {
      dashboard.value.deal_ratio = { ...emptyDashboardSummary().deal_ratio, ...(res.ratio || {}) };
      await nextTick();
      renderRatio();
    }
  } catch {
    /* keep previous */
  }
}

watch(ratioMode, () => {
  renderRatio();
});

onMounted(async () => {
  try {
    roles.value = (await getUserInfoApi()).roles || [];
  } finally {
    await loadDashboard();
  }
});
</script>

<template>
  <Page auto-content-height content-class="platform-dashboard-page">
    <div class="platform-dashboard" :class="{ 'is-loading': loading }">
      <p v-if="failed" class="load-error">经营数据暂不可用，请刷新重试；不会使用演示数据替代真实统计。</p>

      <!-- A1 指标卡 -->
      <section class="metric-grid" aria-label="今日经营指标">
        <ElCard v-for="card in metricCards" :key="card.label" shadow="never" class="metric-card">
          <template v-if="loading"><ElSkeleton :rows="3" animated /></template>
          <template v-else>
            <div class="metric-card__header">
              <span>{{ card.label }}</span>
              <ElTag size="small" effect="plain" type="primary">今日</ElTag>
            </div>
            <strong>{{ formatCount(card.metric.today) }}</strong>
            <div class="metric-card__compare">
              <span>昨日 {{ formatCount(card.metric.yesterday) }}</span>
              <span :class="ratioClass(card.metric.week_ratio)">
                周环比 {{ formatWeekRatio(card.metric.week_ratio) }}
                <i>{{ Number(card.metric.week_ratio || 0) >= 0 ? '↑' : '↓' }}</i>
              </span>
            </div>
            <div class="metric-card__footer">
              <span>{{ card.footerLabel }}</span>
              <b>{{ formatCount(card.footerValue) }}</b>
            </div>
          </template>
        </ElCard>
      </section>

      <!-- A2 快捷入口 -->
      <section class="quick-panel" aria-label="常用功能">
        <ElCard shadow="never" class="dashboard-panel quick-card-wrap">
          <div class="quick-grid">
            <RouterLink v-for="item in quickLinks" :key="item.label" :to="item.route" class="quick-card">
              <span class="quick-card__icon" :class="`quick-card__icon--${item.tone}`">
                <ElIcon><component :is="item.icon" /></ElIcon>
              </span>
              <span>{{ item.label }}</span>
            </RouterLink>
          </div>
        </ElCard>
      </section>

      <!-- A3 待办 | 店铺排行 -->
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
              <strong>店铺销售情况排行 TOP10</strong>
              <ElRadioGroup :model-value="rankPeriod" size="small" @change="(v: any) => loadRank(v)">
                <ElRadioButton v-for="tab in periodTabs" :key="tab.value" :value="tab.value">
                  {{ tab.label }}
                </ElRadioButton>
              </ElRadioGroup>
            </div>
          </template>
          <ElEmpty
            v-if="!rankLoading && !(rankRows?.length)"
            description="所选周期暂无已支付订单"
            :image-size="76"
          />
          <div v-else v-loading="rankLoading" class="rank-list">
            <div class="rank-row rank-row--head" role="row">
              <span>排名</span>
              <span>店铺图</span>
              <span>店铺名称</span>
              <span>关注人数</span>
              <span>销量</span>
              <span>销售金额</span>
            </div>
            <div
              v-for="(row, index) in rankRows || []"
              :key="row.store_id || index"
              class="rank-row"
              role="row"
            >
              <span class="rank-cell rank-cell--center">
                <span class="rank-badge" :class="{ 'rank-badge--top': index < 3 }">{{ index + 1 }}</span>
              </span>
              <span class="rank-cell rank-cell--center">
                <img
                  v-if="row.store_image"
                  class="store-avatar-img"
                  :src="row.store_image"
                  :alt="row.store_name || '店铺'"
                />
                <span v-else class="store-avatar">{{ (row.store_name || '店').slice(0, 1) }}</span>
              </span>
              <span class="rank-cell rank-cell--name" :title="row.store_name">{{ row.store_name || '—' }}</span>
              <span class="rank-cell rank-cell--center">{{ formatCount(row.follower_count) }}</span>
              <span class="rank-cell rank-cell--center">{{ formatCount(row.sale_count) }}</span>
              <span class="rank-cell rank-cell--amount">¥{{ formatMoney(row.sale_amount) }}</span>
            </div>
          </div>
        </ElCard>
      </section>

      <!-- B 当日金额 | 数据统计 -->
      <section class="bottom-grid">
        <ElCard shadow="never" class="dashboard-panel amount-panel">
          <template #header><strong>当日订单金额</strong></template>
          <div class="amount-main">¥ {{ formatMoney(dashboard.today_paid_amount) }}</div>
          <EchartsUI ref="amountChartRef" height="260px" />
        </ElCard>

        <ElCard shadow="never" class="dashboard-panel order-panel">
          <template #header><strong>数据统计</strong></template>
          <div class="order-stats">
            <div
              v-for="(item, index) in sparkStats"
              :key="item.key"
              class="order-stat-cell"
              :class="{ 'order-stat-cell--chart': item.hasChart }"
            >
              <span class="order-stat-cell__label">{{ item.label }}</span>
              <b>{{ formatCount(item.stat?.value) }}</b>
              <em :class="sparkRatioClass(item.stat?.ratio)">
                {{ item.ratioLabel }}: {{ formatSparkRatio(item.stat?.ratio) }}
                {{ Number(item.stat?.ratio || 0) >= 0 ? '▲' : '▼' }}
              </em>
              <div v-if="item.hasChart" class="order-stat-cell__chart">
                <div class="order-stat-cell__legend">
                  <i class="order-stat-cell__legend-mark" />
                  <span>今天</span>
                </div>
                <EchartsUI v-if="index === 0" ref="spark0Ref" height="96px" />
                <EchartsUI v-else ref="spark1Ref" height="96px" />
              </div>
            </div>
          </div>
        </ElCard>
      </section>

      <!-- B 用户数据（唯一三折线） -->
      <section>
        <ElCard shadow="never" class="dashboard-panel user-panel">
          <template #header>
            <div class="panel-header">
              <strong>用户数据</strong>
              <ElRadioGroup :model-value="userPeriod" size="small" @change="(v: any) => loadUserTrend(v)">
                <ElRadioButton v-for="tab in userPeriodTabs" :key="tab.value" :value="tab.value">
                  {{ tab.label }}
                </ElRadioButton>
              </ElRadioGroup>
            </div>
          </template>
          <EchartsUI ref="userChartRef" height="320px" />
        </ElCard>
      </section>

      <!-- C 成交用户 | 成交用户占比 -->
      <section class="deal-grid">
        <ElCard shadow="never" class="dashboard-panel deal-panel">
          <template #header>
            <div class="panel-header">
              <strong>成交用户</strong>
              <ElRadioGroup :model-value="dealPeriod" size="small" @change="(v: any) => loadDeal(v, 'deal')">
                <ElRadioButton v-for="tab in periodTabs" :key="tab.value" :value="tab.value">
                  {{ tab.label }}
                </ElRadioButton>
              </ElRadioGroup>
            </div>
          </template>
          <div class="deal-body">
            <div class="deal-metrics">
              <div class="deal-band deal-band--visit">
                <div class="deal-cell">
                  <b>{{ formatCount(dealFunnel?.visit_users) }}</b>
                  <span>访客人数</span>
                </div>
                <div class="deal-cell deal-cell--empty" />
                <div class="deal-cell deal-cell--empty" />
              </div>
              <div class="deal-band deal-band--order">
                <div class="deal-cell">
                  <b>{{ formatCount(dealFunnel?.order_users) }}</b>
                  <span>下单人数</span>
                </div>
                <div class="deal-cell">
                  <b>{{ formatMoney(dealFunnel?.order_amount) }}</b>
                  <span>下单金额</span>
                </div>
                <div class="deal-cell deal-cell--empty" />
              </div>
              <div class="deal-band deal-band--pay">
                <div class="deal-cell">
                  <b>{{ formatCount(dealFunnel?.pay_users) }}</b>
                  <span>支付人数</span>
                </div>
                <div class="deal-cell">
                  <b>{{ formatMoney(dealFunnel?.pay_amount) }}</b>
                  <span>支付金额</span>
                </div>
                <div class="deal-cell">
                  <b>{{ formatMoney(dealFunnel?.avg_order_amount) }}</b>
                  <span>客单价</span>
                </div>
              </div>
            </div>
            <div class="deal-funnel">
              <div class="funnel-stack" aria-hidden="true">
                <div class="funnel-seg funnel-seg--visit">访客</div>
                <div class="funnel-seg funnel-seg--order">下单</div>
                <div class="funnel-seg funnel-seg--pay">支付</div>
              </div>
              <div class="funnel-rates">
                <div class="funnel-rate funnel-rate--mid1">
                  <i class="funnel-rate__line"></i>
                  <span>访客-下单转化率：{{ Number(dealFunnel?.visit_order_rate || 0).toFixed(0) }} %</span>
                </div>
                <div class="funnel-rate funnel-rate--mid2">
                  <i class="funnel-rate__line"></i>
                  <span>下单-支付转化率：{{ Number(dealFunnel?.order_pay_rate || 0).toFixed(0) }} %</span>
                </div>
              </div>
            </div>
          </div>
        </ElCard>

        <ElCard shadow="never" class="dashboard-panel ratio-panel">
          <template #header>
            <div class="panel-header">
              <strong>成交用户占比</strong>
              <ElRadioGroup :model-value="ratioPeriod" size="small" @change="(v: any) => loadDeal(v, 'ratio')">
                <ElRadioButton v-for="tab in periodTabs" :key="tab.value" :value="tab.value">
                  {{ tab.label }}
                </ElRadioButton>
              </ElRadioGroup>
            </div>
          </template>
          <div class="ratio-modes">
            <button
              type="button"
              class="ratio-mode"
              :class="{ 'is-active': ratioMode === 'amount' }"
              @click="ratioMode = 'amount'"
            >
              金额
            </button>
            <button
              type="button"
              class="ratio-mode"
              :class="{ 'is-active': ratioMode === 'users' }"
              @click="ratioMode = 'users'"
            >
              客户数
            </button>
          </div>
          <EchartsUI ref="ratioChartRef" height="260px" />
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
.platform-dashboard { display: grid; gap: 16px; min-width: 0; padding-bottom: 12px; }
.load-error {
  margin: 0; padding: 10px 14px; border: 1px solid var(--el-color-warning-light-5);
  border-radius: 8px; color: var(--el-color-warning-dark-2); background: var(--el-color-warning-light-9); font-size: 13px;
}
.metric-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 16px; }
.metric-card { min-height: 186px; overflow: hidden; border-color: var(--el-border-color-lighter); }
.metric-card :deep(.el-card__body) { display: flex; min-height: 186px; padding: 22px 26px 0; flex-direction: column; }
.metric-card__header {
  display: flex; align-items: center; justify-content: space-between;
  color: var(--el-text-color-primary); font-size: 16px; font-weight: 600;
}
.metric-card strong { margin-top: 20px; color: var(--el-text-color-primary); font-size: 36px; font-weight: 600; line-height: 1; }
.metric-card__compare {
  display: flex; flex-wrap: wrap; gap: 16px; margin-top: 14px;
  color: var(--el-text-color-secondary); font-size: 13px;
}
.metric-card__compare i { font-style: normal; margin-left: 2px; }
.delta-up { color: #ef4444; }
.delta-down { color: #22c55e; }
.delta-flat { color: var(--el-text-color-secondary); }
.metric-card__footer {
  display: flex; align-items: center; justify-content: space-between;
  margin: auto -26px 0; padding: 14px 26px; border-top: 1px solid var(--el-border-color-lighter);
  color: var(--el-text-color-secondary); font-size: 13px;
}
.metric-card__footer b { color: var(--el-text-color-primary); font-weight: 600; }

.quick-card-wrap :deep(.el-card__body) { padding: 18px 12px; }
.quick-grid { display: grid; grid-template-columns: repeat(8, minmax(0, 1fr)); gap: 8px; }
.quick-card {
  display: grid; min-height: 108px; place-content: center; gap: 12px;
  color: var(--el-text-color-primary); text-align: center; text-decoration: none;
  transition: color .18s ease, transform .18s ease;
}
.quick-card:hover { color: var(--el-color-primary); transform: translateY(-2px); }
.quick-card__icon {
  display: grid; width: 44px; height: 44px; margin: auto; place-items: center;
  border-radius: 12px; font-size: 22px; color: #fff;
}
.quick-card__icon--mint { background: #32bea2; }
.quick-card__icon--blue { background: #4f87ff; }
.quick-card__icon--orange { background: #f3975d; }
.quick-card__icon--purple { background: #9b73e7; }
.quick-card__icon--gold { background: #df9d0a; }
.quick-card__icon--cyan { background: #32bdbd; }
.quick-card__icon--peach { background: #ee9a7d; }
.quick-card__icon--amber { background: #e8a80d; }

.content-grid { display: grid; grid-template-columns: minmax(360px, 1fr) minmax(520px, 1.04fr); gap: 16px; }
.dashboard-panel { border-color: var(--el-border-color-lighter); }
.dashboard-panel :deep(.el-card__header) {
  min-height: 56px; padding: 14px 22px; border-bottom-color: var(--el-border-color-lighter); font-size: 16px;
}
.dashboard-panel :deep(.el-card__body) { padding: 20px 22px; }
.panel-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; flex-wrap: wrap; }

.todo-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 28px 18px; padding: 8px 0; }
.todo-item {
  display: flex; align-items: center; gap: 12px; min-width: 0;
  color: var(--el-text-color-primary); text-decoration: none;
}
.todo-item:hover small { color: var(--el-color-primary); }
.todo-item__icon {
  display: grid; flex: 0 0 auto; width: 48px; height: 48px; place-items: center;
  border-radius: 4px; font-size: 24px;
}
.todo-item__icon--blue { color: #4688ef; background: #edf5ff; }
.todo-item__icon--cyan { color: #20bbb9; background: #e9fbfa; }
.todo-item__icon--orange { color: #ef851a; background: #fff5e8; }
.todo-item__icon--green { color: #65ba38; background: #f1fae9; }
.todo-item__icon--purple { color: #9871ed; background: #f6efff; }
.todo-item__icon--teal { color: #14b8a6; background: #ecfdf8; }
.todo-item__icon--red { color: #f56c6c; background: #fef0f0; }
.todo-item__icon--lime { color: #84cc16; background: #f7fee7; }
.todo-item__icon--amber { color: #d97706; background: #fffbeb; }
.todo-item__text { display: grid; gap: 5px; min-width: 0; }
.todo-item__text small {
  overflow: hidden; color: var(--el-text-color-secondary); font-size: 13px;
  text-overflow: ellipsis; white-space: nowrap;
}
.todo-item__text b { font-size: 22px; font-weight: 600; line-height: 1; }

.rank-list { width: 100%; min-height: 120px; }
.rank-row {
  display: grid;
  grid-template-columns: 56px 56px minmax(100px, 1.5fr) 88px 64px 112px;
  align-items: center;
  column-gap: 8px;
  min-height: 48px;
  padding: 0 4px;
  border-bottom: 1px solid var(--el-border-color-extra-light);
  color: var(--el-text-color-primary);
  font-size: 13px;
}
.rank-row--head {
  min-height: 40px;
  border-bottom-color: var(--el-border-color-lighter);
  color: var(--el-text-color-primary);
  font-weight: 600;
}
.rank-row--head > span:nth-child(1),
.rank-row--head > span:nth-child(2),
.rank-row--head > span:nth-child(4),
.rank-row--head > span:nth-child(5),
.rank-cell--center { text-align: center; }
.rank-row--head > span:nth-child(6),
.rank-cell--amount { text-align: right; font-variant-numeric: tabular-nums; }
.rank-cell--name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.rank-badge {
  display: inline-grid; width: 24px; height: 24px; place-items: center; border-radius: 50%;
  color: var(--el-text-color-secondary); background: var(--el-fill-color-light); font-size: 12px;
}
.rank-badge--top { color: #fff; background: #d9ac4c; }
.store-avatar,
.store-avatar-img {
  display: inline-grid; width: 32px; height: 32px; place-items: center; border-radius: 6px;
  background: #edf5ff; color: #4688ef; font-size: 13px; font-weight: 600;
  object-fit: cover;
}

.bottom-grid { display: grid; grid-template-columns: minmax(300px, .97fr) minmax(500px, 1.03fr); gap: 16px; }
.amount-main { margin-bottom: 8px; color: var(--el-text-color-primary); font-size: 32px; font-weight: 600; line-height: 1.1; }
.order-panel :deep(.el-card__body) { padding: 18px 24px 22px; }
.order-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 28px 48px;
  min-height: 320px;
}
.order-stat-cell { display: grid; gap: 10px; min-width: 0; align-content: start; }
.order-stat-cell--chart { min-height: 200px; }
.order-stat-cell__label { color: var(--el-text-color-secondary); font-size: 13px; line-height: 1.2; }
.order-stat-cell b { color: var(--el-text-color-primary); font-size: 30px; font-weight: 700; line-height: 1.15; }
.order-stat-cell em { font-style: normal; font-size: 12px; line-height: 1.2; }
.order-stat-cell__chart { display: grid; gap: 2px; margin-top: 2px; }
.order-stat-cell__legend {
  display: inline-flex; align-items: center; gap: 8px;
  color: var(--el-text-color-secondary); font-size: 12px;
}
.order-stat-cell__legend-mark {
  position: relative;
  display: inline-block;
  width: 18px;
  height: 8px;
  box-sizing: border-box;
}
.order-stat-cell__legend-mark::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  top: 50%;
  height: 1.5px;
  margin-top: -0.75px;
  background: #409eff;
}
.order-stat-cell__legend-mark::after {
  content: '';
  position: absolute;
  left: 50%;
  top: 50%;
  width: 7px;
  height: 7px;
  margin: -3.5px 0 0 -3.5px;
  border: 1.5px solid #409eff;
  border-radius: 50%;
  background: #fff;
  box-sizing: border-box;
}

.user-panel :deep(.el-card__body) { padding-top: 8px; }

.deal-grid { display: grid; grid-template-columns: 1.2fr .8fr; gap: 16px; }
.deal-body {
  display: grid;
  grid-template-columns: minmax(260px, 1.1fr) minmax(250px, .95fr);
  gap: 16px 12px;
  align-items: stretch;
  min-height: 280px;
}
.deal-metrics { display: grid; grid-template-rows: repeat(3, minmax(78px, 1fr)); gap: 8px; }
.deal-band {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  align-items: center;
  column-gap: 8px;
  padding: 14px 20px;
  border-radius: 2px;
}
.deal-band--visit { background: #edf5ff; }
.deal-band--order { background: #f1fae9; }
.deal-band--pay { background: #f2f3f5; }
.deal-cell { display: grid; gap: 6px; min-width: 0; }
.deal-cell--empty { visibility: hidden; pointer-events: none; }
.deal-cell b {
  color: var(--el-text-color-primary); font-size: 26px; font-weight: 700; line-height: 1.1;
  font-variant-numeric: tabular-nums;
}
.deal-cell span { color: #909399; font-size: 13px; line-height: 1.2; }

.deal-funnel {
  position: relative;
  display: grid;
  grid-template-columns: minmax(130px, 1fr) minmax(148px, 1.15fr);
  gap: 0;
  align-items: stretch;
  min-height: 260px;
  padding-left: 4px;
}
.funnel-stack {
  display: grid;
  grid-template-rows: 1fr 1fr 1fr;
  gap: 5px;
  align-items: stretch;
  justify-items: center;
  padding: 2px 0;
}
.funnel-seg {
  display: grid;
  place-items: center;
  width: 100%;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: 0.02em;
}
.funnel-seg--visit {
  width: 100%;
  background: #409eff;
  clip-path: polygon(0 0, 100% 0, 88% 100%, 12% 100%);
}
.funnel-seg--order {
  width: 78%;
  background: #67c23a;
  clip-path: polygon(2% 0, 98% 0, 84% 100%, 16% 100%);
}
.funnel-seg--pay {
  width: 58%;
  background: #5a6a85;
  clip-path: polygon(4% 0, 96% 0, 78% 100%, 22% 100%);
}

.funnel-rates { position: relative; }
.funnel-rate {
  position: absolute;
  left: -6px;
  right: 0;
  display: flex;
  align-items: center;
  gap: 0;
  color: #606266;
  font-size: 13px;
  white-space: nowrap;
}
.funnel-rate--mid1 { top: calc(33.333% - 2px); transform: translateY(-50%); }
.funnel-rate--mid2 { top: calc(66.666% + 2px); transform: translateY(-50%); }
.funnel-rate__line {
  flex: 0 0 48px;
  height: 0;
  margin-right: 8px;
  border-top: 1px solid #c0c4cc;
}
.funnel-rate span { line-height: 1.35; white-space: normal; }

.ratio-panel :deep(.el-card__body) { padding-top: 12px; }
.ratio-modes { display: flex; gap: 22px; margin-bottom: 8px; padding-left: 2px; }
.ratio-mode {
  border: 0; background: transparent; padding: 0 0 8px; cursor: pointer;
  color: var(--el-text-color-secondary); font-size: 14px; border-bottom: 2px solid transparent;
}
.ratio-mode.is-active { color: var(--el-color-primary); border-bottom-color: var(--el-color-primary); font-weight: 500; }

.dashboard-footer {
  display: flex; align-items: center; justify-content: flex-end; gap: 12px;
  color: var(--el-text-color-secondary); font-size: 12px;
}

@media (max-width: 1440px) {
  .quick-grid { grid-template-columns: repeat(4, minmax(0, 1fr)); }
  .content-grid, .deal-grid, .bottom-grid { grid-template-columns: 1fr; }
  .metric-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 860px) {
  .todo-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .deal-body { grid-template-columns: 1fr; }
  .deal-funnel { grid-template-columns: 1fr 1fr; min-height: 200px; }
  .order-stats { gap: 24px 28px; min-height: 0; }
  .rank-row {
    grid-template-columns: 48px 48px minmax(80px, 1.2fr) 72px 56px 96px;
    font-size: 12px;
  }
}
@media (max-width: 600px) {
  .metric-grid, .quick-grid, .todo-grid, .order-stats { grid-template-columns: 1fr; }
  .deal-band { grid-template-columns: 1fr; }
  .deal-cell--empty { display: none; }
  .rank-row {
    grid-template-columns: 40px 40px minmax(0, 1fr);
    row-gap: 4px;
  }
  .rank-row > :nth-child(n + 4) { grid-column: span 1; }
}
</style>
