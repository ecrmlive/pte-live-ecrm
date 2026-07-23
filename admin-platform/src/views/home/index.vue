<script setup lang="ts">
import type { Component } from 'vue';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type {
  PlatformDashboardData,
  ProductRankRow,
} from '#/api/core/platform-dashboard';

import { computed, onMounted, ref, watch } from 'vue';

import { Page } from '@vben/common-ui';

import { ElButton, ElCard, ElDatePicker, ElIcon, ElImage } from 'element-plus';
import {
  Coin,
  DataLine,
  Film,
  Money,
  Shop,
  VideoCamera,
  Wallet,
} from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  formatCount,
  formatMoney,
  formatTrafficGB,
  getPlatformDashboardApi,
  getPlatformTencentTrafficApi,
  type PlatformTencentTrafficData,
} from '#/api/core/platform-dashboard';
import { formatLiveApiDateTime } from '#/utils/live-api-time';

import './home.scss';

defineOptions({ name: 'PlatformHomeDashboard' });

interface StatCard {
  icon?: Component;
  label: string;
  suffix?: string;
  tone?: string;
  value: string;
}

type DashboardRangeKey = '30d' | '7d' | 'custom' | 'today';

const dayMs = 24 * 60 * 60 * 1000;

function formatDate(date: Date) {
  const year = date.getFullYear();
  const month = `${date.getMonth() + 1}`.padStart(2, '0');
  const day = `${date.getDate()}`.padStart(2, '0');
  return `${year}-${month}-${day}`;
}

function startOfToday() {
  const now = new Date();
  return new Date(now.getFullYear(), now.getMonth(), now.getDate());
}

function defaultCustomRange(): [string, string] {
  const today = startOfToday();
  return [formatDate(new Date(today.getTime() - 29 * dayMs)), formatDate(today)];
}

const loading = ref(false);
const tencentLoading = ref(false);
const cacheTime = ref('');
const tencentCacheTime = ref('');
const tencentConfigured = ref(false);
const activeRange = ref<DashboardRangeKey>('30d');
const customDateRange = ref<[string, string] | null>(defaultCustomRange());
const tencentTraffic = ref<PlatformTencentTrafficData>({
  configured: false,
  lvb: {
    month_play_gb: 0,
    package_remain_gb: 0,
    package_total_gb: 0,
    package_used_gb: 0,
    today_play_gb: 0,
  },
  update_time: '',
  vod: {
    month_play_gb: 0,
    today_play_gb: 0,
  },
});
const overview = ref<PlatformDashboardData['overview']>({
  consumed_traffic_gb: 0,
  merchant_count: 0,
  recharge_amount_yuan: 0,
  recharge_traffic_gb: 0,
});
const trafficSummary = ref<PlatformDashboardData['traffic_summary']>({
  lvb_play_used_gb: 0,
  remain_gb: 0,
  total_gb: 0,
  vod_play_used_gb: 0,
});
const salesSummary = ref<PlatformDashboardData['sales_summary']>({
  paid_amount: 0,
  paid_order_count: 0,
  refund_amount: 0,
  refund_order_count: 0,
  user_count: 0,
});
const productRank = ref<ProductRankRow[]>([]);

const dashboardRangeOptions: Array<{ label: string; value: DashboardRangeKey }> = [
  { label: '今日', value: 'today' },
  { label: '7日', value: '7d' },
  { label: '30天', value: '30d' },
  { label: '自定义', value: 'custom' },
];

const rankGridOptions: VxeGridProps<ProductRankRow> = {
  border: true,
  columns: [
    { type: 'seq', title: '排名', width: 70 },
    { field: 'app_name', minWidth: 150, showOverflow: true, title: '商城名称' },
    {
      field: 'product_image',
      slots: { default: 'product_image' },
      title: '商品图片',
      width: 96,
    },
    {
      field: 'product_name',
      minWidth: 180,
      showOverflow: true,
      title: '商品名称',
    },
    {
      field: 'product_price',
      slots: { default: 'product_price' },
      title: '单价',
      width: 110,
    },
    { field: 'total_num', title: '销量', width: 90 },
    {
      field: 'total_price',
      minWidth: 120,
      slots: { default: 'total_price' },
      title: '总金额',
    },
  ],
  height: 400,
  maxHeight: 400,
  pagerConfig: { enabled: false },
  rowConfig: { isHover: true, keyField: 'rank_key' },
};

const [RankGrid, rankGridApi] = useVbenVxeGrid({
  gridOptions: rankGridOptions,
});

watch(
  productRank,
  (rows) => {
    rankGridApi.setGridOptions({ data: rows });
  },
  { immediate: true },
);

const tencentCacheTimeText = computed(() =>
  formatLiveApiDateTime(tencentCacheTime.value),
);
const dashboardCacheTimeText = computed(() => formatLiveApiDateTime(cacheTime.value));

const dashboardTimeHint = computed(() => {
  if (dashboardCacheTimeText.value === '-') {
    return '';
  }
  return `数据统计于 · ${dashboardCacheTimeText.value}`;
});

const shopStatsCards = computed<StatCard[]>(() => [
  {
    icon: Shop,
    label: '商城数量',
    tone: 'shop',
    value: formatCount(overview.value.merchant_count),
  },
  {
    icon: DataLine,
    label: '充值流量',
    suffix: 'GB',
    tone: 'traffic',
    value: formatTrafficGB(overview.value.recharge_traffic_gb),
  },
  {
    icon: Money,
    label: '充值金额',
    suffix: '元',
    tone: 'money',
    value: formatMoney(overview.value.recharge_amount_yuan),
  },
  {
    icon: VideoCamera,
    label: '云直播消耗',
    suffix: 'GB',
    tone: 'live',
    value: formatTrafficGB(trafficSummary.value.lvb_play_used_gb),
  },
  {
    icon: Film,
    label: '云点播消耗',
    suffix: 'GB',
    tone: 'vod',
    value: formatTrafficGB(trafficSummary.value.vod_play_used_gb),
  },
  {
    icon: Wallet,
    label: '剩余流量',
    suffix: 'GB',
    tone: 'remain',
    value: formatTrafficGB(trafficSummary.value.remain_gb),
  },
]);

const tencentCards = computed<StatCard[]>(() => [
  {
    icon: Wallet,
    label: '云直播·流量包剩余',
    suffix: 'GB',
    value: formatTrafficGB(tencentTraffic.value.lvb.package_remain_gb),
  },
  {
    icon: VideoCamera,
    label: '云直播·今日播放',
    suffix: 'GB',
    value: formatTrafficGB(tencentTraffic.value.lvb.today_play_gb),
  },
  {
    icon: DataLine,
    label: '云直播·本月播放',
    suffix: 'GB',
    value: formatTrafficGB(tencentTraffic.value.lvb.month_play_gb),
  },
  {
    icon: Film,
    label: '云点播·今日播放',
    suffix: 'GB',
    value: formatTrafficGB(tencentTraffic.value.vod.today_play_gb),
  },
  {
    icon: DataLine,
    label: '云点播·本月播放',
    suffix: 'GB',
    value: formatTrafficGB(tencentTraffic.value.vod.month_play_gb),
  },
]);

const salesCards = computed<StatCard[]>(() => [
  {
    icon: Shop,
    label: '用户数',
    value: formatCount(salesSummary.value.user_count),
  },
  {
    icon: Coin,
    label: '成交订单',
    value: formatCount(salesSummary.value.paid_order_count),
  },
  {
    icon: Money,
    label: '成交金额',
    suffix: '元',
    value: formatMoney(salesSummary.value.paid_amount),
  },
  {
    icon: DataLine,
    label: '退单数',
    value: formatCount(salesSummary.value.refund_order_count),
  },
  {
    icon: Wallet,
    label: '退单金额',
    suffix: '元',
    value: formatMoney(salesSummary.value.refund_amount),
  },
]);

async function fetchTencentTraffic(refresh = false) {
  tencentLoading.value = true;
  try {
    const res = await getPlatformTencentTrafficApi(refresh);
    tencentTraffic.value = res.data;
    tencentConfigured.value = res.data.configured;
    tencentCacheTime.value = res.data.cache_time || res.data.update_time;
  } finally {
    tencentLoading.value = false;
  }
}

async function fetchDashboard() {
  const dateRange =
    activeRange.value === 'custom' && customDateRange.value
      ? {
          end_date: customDateRange.value[1],
          start_date: customDateRange.value[0],
        }
      : undefined;
  loading.value = true;
  try {
    const res = await getPlatformDashboardApi(10, activeRange.value, dateRange);
    overview.value = res.data.overview;
    trafficSummary.value = res.data.traffic_summary;
    salesSummary.value = res.data.sales_summary;
    productRank.value = (res.data.product_rank ?? [])
      .filter(
        (row) =>
          Number(row.product_id) > 0 &&
          (Number(row.total_num) > 0 || Number(row.total_price) > 0),
      )
      .map((row) => ({
        ...row,
        rank_key: `${row.app_id}_${row.product_id}`,
      }));
    cacheTime.value = res.data.cache_time || res.data.update_time;
  } finally {
    loading.value = false;
  }
}

function disabledCustomDate(date: Date) {
  const today = startOfToday();
  const minDate = new Date(today);
  minDate.setFullYear(today.getFullYear() - 3);
  const target = new Date(date.getFullYear(), date.getMonth(), date.getDate());
  return target.getTime() < minDate.getTime() || target.getTime() > today.getTime();
}

function switchDashboardRange(range: DashboardRangeKey) {
  if (activeRange.value === range || loading.value) {
    return;
  }
  activeRange.value = range;
  if (range === 'custom' && !customDateRange.value) {
    customDateRange.value = defaultCustomRange();
  }
  void fetchDashboard();
}

function handleCustomDateChange() {
  if (activeRange.value !== 'custom' || loading.value || !customDateRange.value) {
    return;
  }
  void fetchDashboard();
}

onMounted(() => {
  void fetchTencentTraffic();
  void fetchDashboard();
});
</script>

<template>
  <Page auto-content-height>
    <div class="platform-home">
      <section v-loading="tencentLoading" class="platform-home__block">
        <div
          v-if="tencentCacheTimeText !== '-'"
          class="platform-home__time platform-home__time--source"
        >
          数据来自腾讯云 API · {{ tencentCacheTimeText }}（每 1 小时更新）
        </div>
        <div class="platform-home__section-head">
          <h3 class="platform-home__section-title platform-home__section-title--inline">
            腾讯云流量
          </h3>
          <ElButton
            :loading="tencentLoading"
            size="small"
            type="primary"
            link
            @click="fetchTencentTraffic(true)"
          >
            刷新
          </ElButton>
        </div>
        <p v-if="!tencentConfigured" class="platform-home__hint">
          未配置腾讯云 API 密钥或点播子应用，请在系统设置中完善后刷新。
        </p>
        <div class="platform-home__stats platform-home__stats--5">
          <div
            v-for="card in tencentCards"
            :key="card.label"
            class="platform-stat-card platform-stat-card--tencent"
          >
            <div class="platform-stat-card__top">
              <div class="platform-stat-card__label">{{ card.label }}</div>
              <ElIcon v-if="card.icon" class="platform-stat-card__icon">
                <component :is="card.icon" />
              </ElIcon>
            </div>
            <div class="platform-stat-card__value">
              {{ card.value }}
              <span v-if="card.suffix" class="platform-stat-card__suffix">
                {{ card.suffix }}
              </span>
            </div>
          </div>
        </div>
      </section>

      <div v-if="dashboardTimeHint" class="platform-home__time">
        {{ dashboardTimeHint }}
      </div>

      <div v-loading="loading">
        <section class="platform-home__block">
          <div class="platform-home__section-head platform-home__section-head--stats">
            <h3 class="platform-home__section-title platform-home__section-title--inline">
              商城流量统计
            </h3>
            <div class="platform-home__range">
              <ElButton
                v-for="option in dashboardRangeOptions"
                :key="option.value"
                :plain="activeRange !== option.value"
                :type="activeRange === option.value ? 'primary' : 'default'"
                size="small"
                @click="switchDashboardRange(option.value)"
              >
                {{ option.label }}
              </ElButton>
              <ElDatePicker
                v-if="activeRange === 'custom'"
                v-model="customDateRange"
                class="platform-home__custom-range"
                :clearable="false"
                :disabled-date="disabledCustomDate"
                end-placeholder="结束日期"
                range-separator="至"
                size="small"
                start-placeholder="开始日期"
                type="daterange"
                value-format="YYYY-MM-DD"
                @change="handleCustomDateChange"
              />
            </div>
          </div>
          <div class="platform-home__stats platform-home__stats--6">
            <div
              v-for="card in shopStatsCards"
              :key="card.label"
              class="platform-stat-card"
              :class="card.tone ? `platform-stat-card--${card.tone}` : undefined"
            >
              <div class="platform-stat-card__top">
                <div class="platform-stat-card__label">{{ card.label }}</div>
                <ElIcon v-if="card.icon" class="platform-stat-card__icon">
                  <component :is="card.icon" />
                </ElIcon>
              </div>
              <div class="platform-stat-card__value">
                {{ card.value }}
                <span v-if="card.suffix" class="platform-stat-card__suffix">
                  {{ card.suffix }}
                </span>
              </div>
            </div>
          </div>
        </section>

        <section class="platform-home__block">
          <h3 class="platform-home__section-title">商城销售统计</h3>
          <div class="platform-home__stats platform-home__stats--5">
            <div
              v-for="card in salesCards"
              :key="card.label"
              class="platform-stat-card"
            >
              <div class="platform-stat-card__top">
                <div class="platform-stat-card__label">{{ card.label }}</div>
                <ElIcon v-if="card.icon" class="platform-stat-card__icon">
                  <component :is="card.icon" />
                </ElIcon>
              </div>
              <div class="platform-stat-card__value">
                {{ card.value }}
                <span v-if="card.suffix" class="platform-stat-card__suffix">
                  {{ card.suffix }}
                </span>
              </div>
            </div>
          </div>
        </section>

        <ElCard class="platform-home__panel platform-home__block" shadow="never">
          <template #header>
            <span class="platform-home__card-title">商品排行榜</span>
          </template>
          <RankGrid>
            <template #product_image="{ row }">
              <ElImage
                v-if="row.product_image"
                class="platform-rank-product-image"
                fit="cover"
                :src="row.product_image"
              />
              <div
                v-else
                class="platform-rank-product-image platform-rank-product-image--empty"
              >
                -
              </div>
            </template>
            <template #product_price="{ row }">
              {{ formatMoney(row.product_price ?? 0) }}
            </template>
            <template #total_price="{ row }">
              {{ formatMoney(row.total_price ?? 0) }}
            </template>
          </RankGrid>
        </ElCard>
      </div>
    </div>
  </Page>
</template>
