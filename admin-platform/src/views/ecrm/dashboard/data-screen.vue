<script setup lang="ts">
import type { EchartsUIType } from '@vben/plugins/echarts';

import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue';

import { EchartsUI, useEcharts } from '@vben/plugins/echarts';
import * as echarts from 'echarts';

import {
  formatCount,
  formatMoney,
  getPlatformDataScreenApi,
  type DataScreenHourPoint,
  type PlatformDataScreen,
} from '#/api/core/platform-dashboard';
import { ADMIN_TIMEZONE } from '#/utils/date-time';

const ASSET = `${import.meta.env.BASE_URL.replace(/\/?$/, '/')}demo/data-screen`;
const CHINA_GEO_URL = `${import.meta.env.BASE_URL.replace(/\/?$/, '/')}geo/china.json`;
/** 列表滚动像素速度（实时订单 / 店铺排行 / 单品排行共用，避免行高差异导致快慢不一） */
const SCROLL_PX_PER_SEC = 28;
const RANK_ROW_APPROX_PX = 40;
/** 与 .order-feed li 实际行高对齐，避免滚动克隆缝出现大空隙 */
const ORDER_ROW_APPROX_PX = 78;
/** CRMEB 科技风双层圆环素材（外虚线刻度 + 内实线双轨） */
const GAUGE_RINGS = [
  `${ASSET}/gauge-ring-cyan.png`,
  `${ASSET}/gauge-ring-green.png`,
  `${ASSET}/gauge-ring-orange.png`,
  `${ASSET}/gauge-ring-magenta.png`,
] as const;
/** 海南主岛纬度下限：低于此的南海诸岛/南沙多边形从地图剔除 */
const HAINAN_MAINLAND_MIN_LAT = 17.8;

const PANEL_TITLES = {
  amount: '今日订单支付金额(元)',
  hour: '订单支付情况',
  merchant: '店铺销售排行',
  month: '本月销售情况统计',
  newOld: '新老客户占比',
  orders: '实时订单',
  product: '单品销售排行',
  today: '今日数据',
} as const;

const MAP_PIECES = [
  { gte: 1000, label: '1000单以上', color: '#1781b5' },
  { gte: 600, label: '600-999单', lte: 999, color: '#2f90b9' },
  { gte: 200, label: '200-599单', lte: 599, color: '#66a9c9' },
  { gte: 50, label: '50-199单', lte: 199, color: '#8abcd1' },
  { gte: 10, label: '10-49单', lte: 49, color: '#5cb3cc' },
  { gte: 1, label: '1-9单', lte: 9, color: '#c3d7df' },
] as const;

const loading = ref(false);
const nowText = ref('');
const isFullscreen = ref(false);
const data = ref<PlatformDataScreen>(emptyScreen());

const title = computed(
  () => data.value.config?.data_screen_title || '数据大屏',
);
const todayNums = computed(() => data.value.today_pay_count_number);
const newOld = computed(() => data.value.today_pay_new_old);
const newOldTotal = computed(
  () => Math.max(newOld.value.new_count + newOld.value.old_count, 0),
);
const newOldTotalSafe = computed(() => Math.max(newOldTotal.value, 1));

/** 店铺销售排行 ← today_pay_merchant_rank */
const merchantRanks = computed(
  () => data.value.today_pay_merchant_rank?.data || [],
);
const merchantRankRows = computed(() => duplicateForScroll(merchantRanks.value));

/** 单品销售排行 ← pay_product_rank */
const productRanks = computed(() => data.value.pay_product_rank || []);
const productRankRows = computed(() => duplicateForScroll(productRanks.value));

const realtimeOrders = computed(() =>
  duplicateForScroll(data.value.today_pay_info || []),
);
const orderSourceLen = computed(() =>
  Math.max(data.value.today_pay_info?.length || 0, 1),
);

function scrollDurationFor(count: number, rowPx: number) {
  const distance = Math.max(count, 1) * rowPx;
  return `${Math.max(12, Math.round(distance / SCROLL_PX_PER_SEC))}s`;
}

const merchantScrollDuration = computed(() =>
  scrollDurationFor(merchantRanks.value.length, RANK_ROW_APPROX_PX),
);
const productScrollDuration = computed(() =>
  scrollDurationFor(productRanks.value.length, RANK_ROW_APPROX_PX),
);
const orderScrollDuration = computed(() =>
  scrollDurationFor(data.value.today_pay_info?.length || 0, ORDER_ROW_APPROX_PX),
);

const newOldChart = ref<EchartsUIType>();
const monthChart = ref<EchartsUIType>();
const hourChart = ref<EchartsUIType>();
const mapRef = ref<HTMLDivElement>();

const { renderEcharts: renderNewOld } = useEcharts(newOldChart);
const { renderEcharts: renderMonth } = useEcharts(monthChart);
const { renderEcharts: renderHour } = useEcharts(hourChart);

let mapChart: echarts.ECharts | null = null;
let clockTimer: ReturnType<typeof setInterval> | undefined;
let refreshTimer: ReturnType<typeof setInterval> | undefined;
let mapReady = false;
let provinceCenters = new Map<string, { center: number[]; name: string }>();

function emptyScreen(): PlatformDataScreen {
  return {
    city_ranking: [],
    config: { data_screen_title: '数据大屏' },
    month_pay_count: [],
    pay_product_rank: [],
    today_pay_count: [],
    today_pay_count_number: {
      today_pay_number: 0,
      today_pay_user_first: 0,
      visit_num: 0,
      visit_user_num: 0,
    },
    today_pay_info: [],
    today_pay_merchant_rank: { data: [], type: '元' },
    today_pay_new_old: { new_count: 0, old_count: 0 },
    today_pay_number: { count: 0, number: 0, order_id: 0, paid: 1 },
  };
}

function duplicateForScroll<T>(rows: T[]): T[] {
  return rows.length > 1 ? [...rows, ...rows] : rows;
}

function placeholderImg(label: string, hue: number) {
  const safe = (label.slice(0, 2) || '示')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(
    `<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64">
      <rect width="64" height="64" rx="6" fill="hsl(${hue} 55% 28%)"/>
      <text x="32" y="38" text-anchor="middle" fill="#d9f8ff" font-size="18" font-family="sans-serif">${safe}</text>
    </svg>`,
  )}`;
}

function normalizeMapName(name: string) {
  return String(name || '')
    .replace(/壮族自治区|维吾尔自治区|回族自治区|自治区|特别行政区|省|市/g, '')
    .trim();
}

function resolveProvinceName(name: string) {
  const direct = provinceCenters.get(name);
  if (direct) return direct.name;
  const short = provinceCenters.get(normalizeMapName(name));
  return short?.name || name;
}

function formatScreenClock(value: Date) {
  const parts = new Intl.DateTimeFormat('zh-CN', {
    day: '2-digit',
    hour: '2-digit',
    hour12: false,
    minute: '2-digit',
    month: '2-digit',
    second: '2-digit',
    timeZone: ADMIN_TIMEZONE,
    weekday: 'short',
    year: 'numeric',
  })
    .formatToParts(value)
    .reduce<Record<string, string>>((result, part) => {
      result[part.type] = part.value;
      return result;
    }, {});
  const weekday = (parts.weekday || '').replace(/^周/, '周');
  return `${parts.year}-${parts.month}-${parts.day} ${weekday} ${parts.hour}:${parts.minute}:${parts.second}`;
}

function polygonMaxLat(poly: number[][][]): number {
  let max = Number.NEGATIVE_INFINITY;
  for (const ring of poly) {
    for (const point of ring) {
      const lat = Number(point[1]);
      if (Number.isFinite(lat) && lat > max) max = lat;
    }
  }
  return max;
}

/** 剔除南海诸岛/南沙等低纬多边形与无名插图 feature，避免右下角南沙小图块 */
function stripSouthChinaSeaIslands(geoJson: {
  features?: Array<{
    geometry?: { coordinates?: unknown; type?: string };
    properties?: { center?: number[]; name?: string };
  }>;
}) {
  const features = geoJson.features || [];
  geoJson.features = features.filter((feature) => {
    const name = String(feature.properties?.name || '').trim();
    if (!name) return false;
    const geom = feature.geometry;
    if (!geom || name !== '海南省') return true;
    if (geom.type === 'MultiPolygon' && Array.isArray(geom.coordinates)) {
      const kept = (geom.coordinates as number[][][][]).filter(
        (poly) => polygonMaxLat(poly) >= HAINAN_MAINLAND_MIN_LAT,
      );
      if (!kept.length) return false;
      if (kept.length === 1) {
        geom.type = 'Polygon';
        geom.coordinates = kept[0];
      } else {
        geom.coordinates = kept;
      }
      return true;
    }
    if (geom.type === 'Polygon' && Array.isArray(geom.coordinates)) {
      return polygonMaxLat(geom.coordinates as number[][][]) >= HAINAN_MAINLAND_MIN_LAT;
    }
    return true;
  });
  return geoJson;
}

async function ensureChinaMap() {
  if (mapReady) return;
  const res = await fetch(CHINA_GEO_URL);
  const geoJson = stripSouthChinaSeaIslands(await res.json());
  provinceCenters.clear();
  for (const feature of geoJson.features as Array<{
    properties?: { center?: number[]; name?: string };
  }>) {
    const { center, name } = feature.properties || {};
    if (!Array.isArray(center) || !name) continue;
    const point = { center, name };
    provinceCenters.set(name, point);
    provinceCenters.set(normalizeMapName(name), point);
  }
  echarts.registerMap('china', geoJson as never);
  mapReady = true;
}

async function renderMap() {
  if (!mapRef.value) return;
  await ensureChinaMap();
  if (!mapChart) {
    mapChart = echarts.init(mapRef.value);
  }
  const list = (data.value.city_ranking || []).map((item) => ({
    name: resolveProvinceName(item.name),
    value: item.value,
  }));
  const mapPoints = list
    .map((item) => {
      const province =
        provinceCenters.get(item.name) ||
        provinceCenters.get(normalizeMapName(item.name));
      return province
        ? {
            name: province.name,
            value: [...province.center, item.value],
          }
        : undefined;
    })
    .filter(
      (point): point is { name: string; value: number[] } => point !== undefined,
    );

  mapChart.setOption(
    {
      backgroundColor: 'transparent',
      geo: {
        itemStyle: {
          areaColor: 'rgba(17, 129, 188, 0.16)',
          borderColor: 'rgba(104, 230, 255, 0.28)',
          borderWidth: 1,
        },
        map: 'china',
        roam: false,
        selectedMode: false,
        show: true,
        top: 28,
        zoom: 1.08,
      },
      series: [
        {
          data: list,
          emphasis: {
            itemStyle: {
              areaColor: 'rgba(80, 210, 255, 0.45)',
            },
            label: { show: false },
          },
          geoIndex: 0,
          itemStyle: {
            areaColor: 'rgba(17, 129, 188, 0.16)',
            borderColor: 'rgba(147, 235, 248, 0.55)',
            borderWidth: 1,
          },
          label: { show: false },
          map: 'china',
          name: '订单分布',
          roam: false,
          top: 28,
          type: 'map',
        },
        {
          coordinateSystem: 'geo',
          data: mapPoints,
          effectType: 'ripple',
          itemStyle: {
            color: '#fff',
            shadowBlur: 8,
            shadowColor: '#8beeff',
          },
          label: {
            color: '#fff',
            fontSize: 11,
            formatter: (params: { name: string }) => params.name.slice(0, 2),
            position: 'bottom',
            show: true,
            textShadowBlur: 4,
            textShadowColor: '#00124c',
          },
          name: '订单城市',
          rippleEffect: { brushType: 'stroke', scale: 5 },
          showEffectOn: 'render',
          symbolSize: 4,
          type: 'effectScatter',
        },
      ],
      tooltip: {
        formatter: (params: { name?: string; value?: number | number[] }) => {
          const value = Array.isArray(params.value)
            ? params.value[params.value.length - 1]
            : params.value;
          return `${params.name || ''}：${formatCount(Number(value || 0))} 单`;
        },
        trigger: 'item',
      },
      visualMap: {
        bottom: 18,
        left: 16,
        pieces: MAP_PIECES.map((piece) => ({ ...piece })),
        show: true,
        textStyle: { color: '#d9f8ff', fontSize: 10 },
        type: 'piecewise',
      },
    },
    true,
  );
}

function hourAxisLabels(points: DataScreenHourPoint[]) {
  if (points.length) {
    return points.map((item) => String(item.hours || '').replace('~', '-'));
  }
  const labels: string[] = [];
  for (let hour = 0; hour < 24; hour += 2) {
    labels.push(
      `${String(hour).padStart(2, '0')}-${String(hour + 1).padStart(2, '0')}`,
    );
  }
  return labels;
}

function renderAllCharts() {
  const oldPct = Math.round(
    (newOld.value.old_count / newOldTotalSafe.value) * 100,
  );
  const newPct = Math.round(
    (newOld.value.new_count / newOldTotalSafe.value) * 100,
  );
  renderNewOld({
    color: ['#f5a623', '#09caf3'],
    series: [
      {
        center: ['50%', '50%'],
        data: [
          { name: '新用户', value: newOld.value.new_count },
          { name: '老用户', value: newOld.value.old_count },
        ],
        itemStyle: {
          borderColor: '#041028',
          borderWidth: 2,
        },
        label: { show: false },
        radius: ['40%', '58%'],
        type: 'pie',
      },
    ],
    title: {
      left: 'center',
      text: `{n|${formatCount(newOldTotal.value)}}\n{t|总数}`,
      textStyle: {
        rich: {
          n: {
            color: '#e2e8f0',
            fontSize: 18,
            fontWeight: 700,
            lineHeight: 24,
          },
          t: { color: '#94a3b8', fontSize: 11, lineHeight: 16 },
        },
      },
      top: '40%',
    },
    tooltip: {
      formatter: () =>
        `老用户 ${formatCount(newOld.value.old_count)}（${oldPct}%）<br/>新用户 ${formatCount(newOld.value.new_count)}（${newPct}%）`,
      trigger: 'item',
    },
  });

  const monthPoints = data.value.month_pay_count || [];
  renderMonth({
    grid: {
      bottom: 44,
      containLabel: true,
      left: 12,
      right: 12,
      top: 14,
    },
    series: [
      {
        barMaxWidth: 12,
        clip: false,
        data: monthPoints.map((item) => item.total_sum),
        itemStyle: {
          borderRadius: [5, 5, 0, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { color: 'rgba(156, 107, 211, 0.86)', offset: 0 },
            { color: 'rgba(156, 107, 211, 0.12)', offset: 1 },
          ]),
        },
        type: 'bar',
      },
    ],
    tooltip: { trigger: 'axis' },
    xAxis: {
      axisLabel: {
        color: '#7dd3fc',
        fontSize: 11,
        hideOverlap: true,
        margin: 14,
      },
      axisLine: { lineStyle: { color: '#33d3ff' } },
      axisTick: { show: false },
      boundaryGap: true,
      data: monthPoints.map((item) => item.day),
      type: 'category',
    },
    yAxis: {
      axisLabel: { color: '#7dd3fc', fontSize: 11, margin: 8 },
      minInterval: 1,
      splitLine: { show: false },
      type: 'value',
    },
  });

  const hourPoints = data.value.today_pay_count || [];
  const hourLabels = hourAxisLabels(hourPoints);
  const orderSeries = hourPoints.length
    ? hourPoints.map((item) => item.order_count)
    : hourLabels.map(() => 0);
  const userSeries = hourPoints.length
    ? hourPoints.map((item) => item.user_count)
    : hourLabels.map(() => 0);

  renderHour({
    grid: {
      borderColor: 'rgba(31, 99, 163, 0.72)',
      bottom: 44,
      containLabel: true,
      left: 12,
      right: 12,
      show: true,
      top: 14,
    },
    legend: { show: false },
    series: [
      {
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { color: 'rgba(252, 144, 16, 0.72)', offset: 0 },
            { color: 'rgba(252, 144, 16, 0)', offset: 1 },
          ]),
        },
        clip: false,
        data: orderSeries,
        itemStyle: { color: '#fc9010' },
        lineStyle: { width: 2 },
        name: '订单数',
        showSymbol: false,
        smooth: true,
        type: 'line',
      },
      {
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { color: 'rgba(9, 202, 243, 0.6)', offset: 0 },
            { color: 'rgba(9, 202, 243, 0)', offset: 1 },
          ]),
        },
        clip: false,
        data: userSeries,
        itemStyle: { color: '#09caf3' },
        lineStyle: { width: 2 },
        name: '购买人数',
        showSymbol: false,
        smooth: true,
        type: 'line',
      },
    ],
    tooltip: { trigger: 'axis' },
    xAxis: {
      axisLabel: {
        color: '#7dd3fc',
        fontSize: 10,
        hideOverlap: true,
        margin: 14,
      },
      axisLine: { lineStyle: { color: 'rgba(31, 99, 163, 0.45)' } },
      data: hourLabels,
      splitLine: {
        lineStyle: { color: 'rgba(31, 99, 163, 0.22)' },
        show: true,
      },
      type: 'category',
    },
    yAxis: {
      axisLabel: { color: '#7dd3fc', fontSize: 11, margin: 8 },
      min: 0,
      minInterval: 1,
      splitLine: { lineStyle: { color: 'rgba(31, 99, 163, 0.22)' } },
      type: 'value',
    },
  });

  void renderMap();
}

async function load() {
  loading.value = true;
  try {
    const raw = await getPlatformDataScreenApi();
    data.value = {
      ...(raw || emptyScreen()),
      pay_product_rank: raw?.pay_product_rank || [],
      today_pay_merchant_rank: raw?.today_pay_merchant_rank || {
        data: [],
        type: '元',
      },
    };
    await nextTick();
    renderAllCharts();
  } catch {
    data.value = emptyScreen();
    await nextTick();
    renderAllCharts();
  } finally {
    loading.value = false;
  }
}

function tickClock() {
  nowText.value = formatScreenClock(new Date());
}

function syncFullscreenFlag() {
  isFullscreen.value = Boolean(document.fullscreenElement);
}

function toggleFullscreen() {
  if (!document.fullscreenElement) {
    void document.documentElement.requestFullscreen?.();
  } else {
    void document.exitFullscreen?.();
  }
}

function onResize() {
  mapChart?.resize();
}

function isUsableImg(url?: string) {
  const s = String(url || '').trim();
  if (!s) return false;
  if (s.startsWith('data:image/')) return true;
  if (/^https?:\/\//i.test(s) || s.startsWith('/')) return true;
  return false;
}

function storeImg(row: { image?: string; store_name?: string }, idx: number) {
  if (isUsableImg(row.image)) return row.image!;
  return placeholderImg(row.store_name || '店', 190 + (idx % 10) * 9);
}

function productImg(
  row: { image?: string; product_name?: string },
  idx: number,
) {
  if (isUsableImg(row.image)) return row.image!;
  return placeholderImg(row.product_name || '品', 30 + (idx % 10) * 11);
}

watch(data, () => {
  nextTick(() => renderAllCharts());
});

onMounted(async () => {
  tickClock();
  clockTimer = setInterval(tickClock, 1000);
  refreshTimer = setInterval(() => {
    void load();
  }, 60_000);
  syncFullscreenFlag();
  document.addEventListener('fullscreenchange', syncFullscreenFlag);
  window.addEventListener('resize', onResize);
  await load();
});

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer);
  if (refreshTimer) clearInterval(refreshTimer);
  document.removeEventListener('fullscreenchange', syncFullscreenFlag);
  window.removeEventListener('resize', onResize);
  mapChart?.dispose();
  mapChart = null;
});
</script>

<template>
  <div class="ds" :class="{ 'is-loading': loading }">
    <div class="ds__side-rail ds__side-rail--left" aria-hidden="true">
      <span class="ds__side-rail__spine" />
      <span class="ds__side-rail__glow" />
      <span class="ds__side-rail__ticks" />
    </div>
    <div class="ds__side-rail ds__side-rail--right" aria-hidden="true">
      <span class="ds__side-rail__spine" />
      <span class="ds__side-rail__glow" />
      <span class="ds__side-rail__ticks" />
    </div>

    <header class="ds__header">
      <div class="ds__bar-ornament ds__bar-ornament--l" />
      <div class="ds__bar-ornament ds__bar-ornament--r" />
      <div class="ds__clock">{{ nowText }}</div>
      <div class="ds__title-wrap">
        <h1 class="ds__title">{{ title }}</h1>
      </div>
      <button class="ds__fullscreen" type="button" @click="toggleFullscreen">
        <span>{{ isFullscreen ? '退出全屏' : '全屏' }}</span>
        <span
          class="ds__fullscreen-icon"
          :class="{ 'is-exit': isFullscreen }"
          aria-hidden="true"
        />
      </button>
    </header>

    <div class="ds__body">
      <aside class="ds__col">
        <section class="panel panel--today">
          <div class="panel__frame" aria-hidden="true">
            <i class="c c--tl" /><i class="c c--tr" /><i class="c c--bl" /><i class="c c--br" />
            <i class="e e--t" /><i class="e e--b" /><i class="e e--l" /><i class="e e--r" />
          </div>
          <h3 class="panel__title">
            <img class="panel__chevron" :src="`${ASSET}/title-chevron.png`" alt="" />
            <span class="panel__title-text">{{ PANEL_TITLES.today }}</span>
            <img
              class="panel__chevron panel__chevron--r"
              :src="`${ASSET}/title-chevron.png`"
              alt=""
            />
          </h3>
          <div class="gauge-grid">
            <div class="gauge-item" style="color: #09caf3">
              <div
                class="gauge-item__ring"
                :style="{ backgroundImage: `url(${GAUGE_RINGS[0]})` }"
                aria-hidden="true"
              />
              <b class="gauge-item__value">{{ formatCount(todayNums.visit_num) }}</b>
              <p>浏览量</p>
            </div>
            <div class="gauge-item" style="color: #33ffc3">
              <div
                class="gauge-item__ring"
                :style="{ backgroundImage: `url(${GAUGE_RINGS[1]})` }"
                aria-hidden="true"
              />
              <b class="gauge-item__value">{{ formatCount(todayNums.visit_user_num) }}</b>
              <p>访客数</p>
            </div>
            <div class="gauge-item" style="color: #f5a623">
              <div
                class="gauge-item__ring"
                :style="{ backgroundImage: `url(${GAUGE_RINGS[2]})` }"
                aria-hidden="true"
              />
              <b class="gauge-item__value">{{
                formatCount(todayNums.today_pay_user_first)
              }}</b>
              <p>新增用户数</p>
            </div>
            <div class="gauge-item" style="color: #a78bfa">
              <div
                class="gauge-item__ring"
                :style="{ backgroundImage: `url(${GAUGE_RINGS[3]})` }"
                aria-hidden="true"
              />
              <b class="gauge-item__value">{{
                formatCount(todayNums.today_pay_number)
              }}</b>
              <p>订单数</p>
            </div>
          </div>
        </section>

        <section class="panel panel--new-old">
          <div class="panel__frame" aria-hidden="true">
            <i class="c c--tl" /><i class="c c--tr" /><i class="c c--bl" /><i class="c c--br" />
            <i class="e e--t" /><i class="e e--b" /><i class="e e--l" /><i class="e e--r" />
          </div>
          <h3 class="panel__title">
            <img class="panel__chevron" :src="`${ASSET}/title-chevron.png`" alt="" />
            <span class="panel__title-text">{{ PANEL_TITLES.newOld }}</span>
            <img
              class="panel__chevron panel__chevron--r"
              :src="`${ASSET}/title-chevron.png`"
              alt=""
            />
          </h3>
          <div class="new-old">
            <div class="new-old__side">
              <small>老用户数量</small>
              <b>{{ formatCount(newOld.old_count) }}个</b>
              <span class="is-old">
                {{ Math.round((newOld.old_count / newOldTotalSafe) * 100) }}%
              </span>
            </div>
            <EchartsUI ref="newOldChart" height="100%" />
            <div class="new-old__side">
              <small>新用户数量</small>
              <b>{{ formatCount(newOld.new_count) }}个</b>
              <span class="is-new">
                {{ Math.round((newOld.new_count / newOldTotalSafe) * 100) }}%
              </span>
            </div>
          </div>
        </section>

        <section class="panel panel--orders">
          <div class="panel__frame" aria-hidden="true">
            <i class="c c--tl" /><i class="c c--tr" /><i class="c c--bl" /><i class="c c--br" />
            <i class="e e--t" /><i class="e e--b" /><i class="e e--l" /><i class="e e--r" />
          </div>
          <h3 class="panel__title">
            <img class="panel__chevron" :src="`${ASSET}/title-chevron.png`" alt="" />
            <span class="panel__title-text">{{ PANEL_TITLES.orders }}</span>
            <img
              class="panel__chevron panel__chevron--r"
              :src="`${ASSET}/title-chevron.png`"
              alt=""
            />
          </h3>
          <div class="order-feed-viewport">
            <ul
              v-if="data.today_pay_info?.length"
              class="order-feed"
              :class="{ 'is-scrolling': data.today_pay_info.length > 1 }"
              :style="{ '--scroll-duration': orderScrollDuration }"
            >
              <li
                v-for="(row, idx) in realtimeOrders"
                :key="`${row.paytime}-${idx}`"
              >
                <span class="no">No{{ (idx % orderSourceLen) + 1 }}</span>
                <div class="order-feed__body">
                  <div class="order-feed__identity">
                    <div class="order-feed__entity">
                      <img
                        :src="storeImg(row.store, idx)"
                        alt=""
                      />
                      <span>{{ row.store.store_name }}</span>
                    </div>
                    <div class="order-feed__entity">
                      <img
                        :src="productImg(row.product, idx)"
                        alt=""
                      />
                      <strong>{{ row.product.product_name }}</strong>
                    </div>
                  </div>
                  <div class="order-feed__meta">
                    <p>支付时间：{{ row.paytime }}</p>
                    <p>支付方式：{{ row.payment_method }}</p>
                    <p>支付金额：{{ formatMoney(row.number) }}</p>
                  </div>
                </div>
              </li>
            </ul>
            <div v-else class="empty empty--panel">暂无实时订单</div>
          </div>
        </section>
      </aside>

      <main class="ds__col ds__col--center">
        <section class="panel amount-panel">
          <h3 class="panel__title">
            <img class="panel__chevron" :src="`${ASSET}/title-chevron.png`" alt="" />
            <span class="panel__title-text">{{ PANEL_TITLES.amount }}</span>
            <img
              class="panel__chevron panel__chevron--r"
              :src="`${ASSET}/title-chevron.png`"
              alt=""
            />
          </h3>
          <div class="amount">{{ formatMoney(data.today_pay_number.number) }}</div>
        </section>
        <section class="panel map-panel">
          <div ref="mapRef" class="map-host" />
        </section>
        <section class="panel panel--month">
          <div class="panel__frame" aria-hidden="true">
            <i class="c c--tl" /><i class="c c--tr" /><i class="c c--bl" /><i class="c c--br" />
            <i class="e e--t" /><i class="e e--b" /><i class="e e--l" /><i class="e e--r" />
          </div>
          <h3 class="panel__title">
            <img class="panel__chevron" :src="`${ASSET}/title-chevron.png`" alt="" />
            <span class="panel__title-text">{{ PANEL_TITLES.month }}</span>
            <img
              class="panel__chevron panel__chevron--r"
              :src="`${ASSET}/title-chevron.png`"
              alt=""
            />
          </h3>
          <EchartsUI ref="monthChart" height="100%" />
        </section>
      </main>

      <aside class="ds__col">
        <section class="panel panel--payment">
          <div class="panel__frame" aria-hidden="true">
            <i class="c c--tl" /><i class="c c--tr" /><i class="c c--bl" /><i class="c c--br" />
            <i class="e e--t" /><i class="e e--b" /><i class="e e--l" /><i class="e e--r" />
          </div>
          <h3 class="panel__title">
            <img class="panel__chevron" :src="`${ASSET}/title-chevron.png`" alt="" />
            <span class="panel__title-text">{{ PANEL_TITLES.hour }}</span>
            <img
              class="panel__chevron panel__chevron--r"
              :src="`${ASSET}/title-chevron.png`"
              alt=""
            />
          </h3>
          <EchartsUI ref="hourChart" height="100%" />
        </section>

        <section class="panel panel--rank panel--merchant">
          <div class="panel__frame" aria-hidden="true">
            <i class="c c--tl" /><i class="c c--tr" /><i class="c c--bl" /><i class="c c--br" />
            <i class="e e--t" /><i class="e e--b" /><i class="e e--l" /><i class="e e--r" />
          </div>
          <h3 class="panel__title">
            <img class="panel__chevron" :src="`${ASSET}/title-chevron.png`" alt="" />
            <span class="panel__title-text">{{ PANEL_TITLES.merchant }}</span>
            <img
              class="panel__chevron panel__chevron--r"
              :src="`${ASSET}/title-chevron.png`"
              alt=""
            />
          </h3>
          <div class="rank-table-viewport">
            <table class="rank-table rank-table--head" aria-hidden="false">
              <colgroup>
                <col class="col-rank" />
                <col class="col-info" />
                <col class="col-count" />
                <col class="col-amount" />
              </colgroup>
              <thead>
                <tr>
                  <th>排名</th>
                  <th>店铺信息</th>
                  <th>销量</th>
                  <th>销售金额</th>
                </tr>
              </thead>
            </table>
            <div class="rank-body-viewport">
              <table
                v-if="merchantRanks.length"
                class="rank-table rank-table--body"
                :class="{ 'is-scrolling': merchantRanks.length > 1 }"
                :style="{ '--scroll-duration': merchantScrollDuration }"
              >
                <colgroup>
                  <col class="col-rank" />
                  <col class="col-info" />
                  <col class="col-count" />
                  <col class="col-amount" />
                </colgroup>
                <tbody>
                  <tr
                    v-for="(row, idx) in merchantRankRows"
                    :key="`m-${row.store.store_name}-${idx}`"
                  >
                    <td>
                      <span
                        class="rank-no"
                        :class="{ top: (idx % merchantRanks.length) < 3 }"
                      >
                        No{{ (idx % merchantRanks.length) + 1 }}
                      </span>
                    </td>
                    <td>
                      <div class="product-cell">
                        <img
                          :src="storeImg(row.store, idx)"
                          alt=""
                          class="product-cell__img"
                        />
                        <strong class="product-cell__text">{{
                          row.store.store_name || '—'
                        }}</strong>
                      </div>
                    </td>
                    <td>{{ formatCount(row.count) }}</td>
                    <td>{{ formatMoney(row.number) }}</td>
                  </tr>
                </tbody>
              </table>
              <div v-else class="empty empty--panel">暂无排行数据</div>
            </div>
          </div>
        </section>

        <section class="panel panel--grow">
          <div class="panel__frame" aria-hidden="true">
            <i class="c c--tl" /><i class="c c--tr" /><i class="c c--bl" /><i class="c c--br" />
            <i class="e e--t" /><i class="e e--b" /><i class="e e--l" /><i class="e e--r" />
          </div>
          <h3 class="panel__title">
            <img class="panel__chevron" :src="`${ASSET}/title-chevron.png`" alt="" />
            <span class="panel__title-text">{{ PANEL_TITLES.product }}</span>
            <img
              class="panel__chevron panel__chevron--r"
              :src="`${ASSET}/title-chevron.png`"
              alt=""
            />
          </h3>
          <div class="rank-table-viewport">
            <table class="rank-table rank-table--head">
              <colgroup>
                <col class="col-rank" />
                <col class="col-info" />
                <col class="col-count" />
                <col class="col-amount" />
              </colgroup>
              <thead>
                <tr>
                  <th>排名</th>
                  <th>商品信息</th>
                  <th>销量</th>
                  <th>销售金额</th>
                </tr>
              </thead>
            </table>
            <div class="rank-body-viewport">
              <table
                v-if="productRanks.length"
                class="rank-table rank-table--body"
                :class="{ 'is-scrolling': productRanks.length > 1 }"
                :style="{ '--scroll-duration': productScrollDuration }"
              >
                <colgroup>
                  <col class="col-rank" />
                  <col class="col-info" />
                  <col class="col-count" />
                  <col class="col-amount" />
                </colgroup>
                <tbody>
                  <tr
                    v-for="(row, idx) in productRankRows"
                    :key="`p-${row.product.product_name}-${idx}`"
                  >
                    <td>
                      <span
                        class="rank-no"
                        :class="{ top: (idx % productRanks.length) < 3 }"
                      >
                        No{{ (idx % productRanks.length) + 1 }}
                      </span>
                    </td>
                    <td>
                      <div class="product-cell">
                        <img
                          :src="productImg(row.product, idx)"
                          alt=""
                          class="product-cell__img"
                        />
                        <strong class="product-cell__text">{{
                          row.product.product_name || '—'
                        }}</strong>
                      </div>
                    </td>
                    <td>{{ formatCount(row.count) }}</td>
                    <td>{{ formatMoney(row.number) }}</td>
                  </tr>
                </tbody>
              </table>
              <div v-else class="empty empty--panel">暂无排行数据</div>
            </div>
          </div>
        </section>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.ds {
  --ds-side-pad: 15px;
  /* 左列三等分时两段 gap(6px*2)，本月销售高度与实时订单同高，只动中列 */
  --ds-left-third: calc((100% - 12px) / 3);
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  width: 100%;
  height: 100vh;
  /* 整页左右 15px：内容内缩，侧边装饰条落在边距带内不压内容 */
  padding: 8px var(--ds-side-pad) 0;
  overflow: hidden;
  color: #d7e9f7;
  background:
    radial-gradient(ellipse 62% 68% at 50% 46%, rgb(17 112 193 / 28%), transparent 62%),
    radial-gradient(ellipse 72% 38% at 50% -6%, rgb(24 105 204 / 36%), transparent 70%),
    url('/demo/data-screen/page-bg.png') center / cover no-repeat,
    linear-gradient(180deg, #020d62 0%, #031365 50%, #020653 100%);
}

.ds.is-loading {
  opacity: 0.94;
}

/* —— 左右整屏修饰条（对齐 CRMEB 侧边机械感线条 + 流光） —— */
.ds__side-rail {
  position: absolute;
  z-index: 3;
  top: 78px;
  bottom: 14px;
  width: 10px;
  pointer-events: none;
}

.ds__side-rail--left {
  /* absolute 相对 padding edge，落在 15px 边距带内 */
  left: 0;
}

.ds__side-rail--right {
  right: 0;
  transform: scaleX(-1);
}

.ds__side-rail__spine {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 3px;
  width: 2px;
  background: linear-gradient(
    180deg,
    transparent 0%,
    rgb(43 205 247 / 55%) 8%,
    rgb(43 205 247 / 85%) 50%,
    rgb(43 205 247 / 55%) 92%,
    transparent 100%
  );
  box-shadow: 0 0 8px rgb(80 230 255 / 35%);
}

.ds__side-rail__spine::before,
.ds__side-rail__spine::after {
  position: absolute;
  left: -3px;
  width: 8px;
  height: 2px;
  content: '';
  background: #2bcdf7;
  box-shadow: 0 0 6px #54dffb;
}

.ds__side-rail__spine::before {
  top: 0;
}

.ds__side-rail__spine::after {
  bottom: 0;
}

.ds__side-rail__ticks {
  position: absolute;
  top: 24px;
  bottom: 24px;
  left: 0;
  width: 8px;
  background:
    repeating-linear-gradient(
      180deg,
      transparent 0 18px,
      rgb(80 230 255 / 55%) 18px 19px,
      transparent 19px 36px
    );
  mask-image: linear-gradient(180deg, transparent, #000 12%, #000 88%, transparent);
  opacity: 0.55;
}

.ds__side-rail__glow {
  position: absolute;
  top: 0;
  left: 1px;
  width: 6px;
  height: 64px;
  background: linear-gradient(
    180deg,
    transparent,
    rgb(128 239 255 / 95%),
    transparent
  );
  filter: blur(0.2px) drop-shadow(0 0 6px #54dffb);
  animation: data-screen-rail-flow 6.5s cubic-bezier(0.4, 0, 0.2, 1) infinite;
}

.ds__header {
  position: relative;
  display: grid;
  flex: 0 0 auto;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  height: 70px;
  margin-bottom: 6px;
  padding: 0 28px;
  background: url('/demo/data-screen/title-bar.png') center / 100% 100% no-repeat;
}

.ds__bar-ornament {
  position: absolute;
  top: -2px;
  width: 140px;
  height: 6px;
  background: url('/demo/data-screen/zuojuxing.png') center / 100% 100% no-repeat;
}

.ds__bar-ornament--l {
  left: 11%;
}

.ds__bar-ornament--r {
  right: 11%;
  transform: rotate(180deg);
}

.ds__clock {
  padding-left: 2px;
  font-size: 22px;
  font-weight: 500;
  letter-spacing: 0.04em;
  color: #33d3ff;
}

.ds__title-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}

.ds__title {
  margin: 0;
  font-size: 38px;
  font-weight: 900;
  letter-spacing: 0.18em;
  background: linear-gradient(92deg, #0072ff, #00eaff 48.85%, #01aaff);
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.ds__fullscreen {
  display: inline-flex;
  gap: 10px;
  align-items: center;
  justify-self: end;
  padding: 4px 2px;
  font-size: 22px;
  line-height: 1;
  color: #33d3ff;
  cursor: pointer;
  background: transparent;
  border: 0;
}

.ds__fullscreen:hover {
  color: #9fe8ff;
}

.ds__fullscreen-icon {
  display: inline-block;
  width: 18px;
  height: 18px;
  background:
    linear-gradient(currentcolor, currentcolor) left top / 7px 2px no-repeat,
    linear-gradient(currentcolor, currentcolor) left top / 2px 7px no-repeat,
    linear-gradient(currentcolor, currentcolor) right top / 7px 2px no-repeat,
    linear-gradient(currentcolor, currentcolor) right top / 2px 7px no-repeat,
    linear-gradient(currentcolor, currentcolor) left bottom / 7px 2px no-repeat,
    linear-gradient(currentcolor, currentcolor) left bottom / 2px 7px no-repeat,
    linear-gradient(currentcolor, currentcolor) right bottom / 7px 2px no-repeat,
    linear-gradient(currentcolor, currentcolor) right bottom / 2px 7px no-repeat;
}

.ds__fullscreen-icon.is-exit {
  background:
    linear-gradient(currentcolor, currentcolor) 3px 3px / 5px 2px no-repeat,
    linear-gradient(currentcolor, currentcolor) 3px 3px / 2px 5px no-repeat,
    linear-gradient(currentcolor, currentcolor) calc(100% - 3px) 3px / 5px 2px no-repeat,
    linear-gradient(currentcolor, currentcolor) calc(100% - 3px) 3px / 2px 5px no-repeat,
    linear-gradient(currentcolor, currentcolor) 3px calc(100% - 3px) / 5px 2px no-repeat,
    linear-gradient(currentcolor, currentcolor) 3px calc(100% - 3px) / 2px 5px no-repeat,
    linear-gradient(currentcolor, currentcolor) calc(100% - 3px) calc(100% - 3px) / 5px 2px
      no-repeat,
    linear-gradient(currentcolor, currentcolor) calc(100% - 3px) calc(100% - 3px) / 2px 5px
      no-repeat;
}

.ds__body {
  display: grid;
  flex: 1;
  /* fr 比例替代 vw；左右边距已由 .ds 的 15px 承担 */
  grid-template-columns: minmax(0, 30fr) minmax(0, 36fr) minmax(0, 30fr);
  gap: 6px;
  min-height: 0;
  padding: 0;
  overflow: hidden;
}

.ds__col {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
  min-height: 0;
  height: 100%;
  overflow: hidden;
}

.ds__col--center {
  /* 与左列同 gap，便于本月销售顶边与实时订单对齐 */
  gap: 6px;
}

.panel {
  position: relative;
  box-sizing: border-box;
  padding: 0 10px;
  background: linear-gradient(
    90deg,
    rgb(15 78 169 / 10%),
    transparent 14% 86%,
    rgb(15 78 169 / 10%)
  );
}

.panel__frame {
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
}

.panel__frame .c,
.panel__frame .e {
  position: absolute;
  display: block;
  pointer-events: none;
}

/* 四角括号（对齐 DataV border-box-1 气质） */
.panel__frame .c {
  width: 22px;
  height: 22px;
  border: 2px solid rgb(43 207 244 / 0%);
  filter: drop-shadow(0 0 0 rgb(69 220 255 / 0%));
  animation: data-screen-corner-breathe 3.6s ease-in-out infinite;
}

.panel__frame .c--tl {
  top: 0;
  left: 0;
  border-top-color: #2ccff4;
  border-left-color: #2ccff4;
}

.panel__frame .c--tr {
  top: 0;
  right: 0;
  border-top-color: #2ccff4;
  border-right-color: #2ccff4;
  animation-delay: 0.4s;
}

.panel__frame .c--bl {
  bottom: 0;
  left: 0;
  border-bottom-color: #2ccff4;
  border-left-color: #2ccff4;
  animation-delay: 0.8s;
}

.panel__frame .c--br {
  right: 0;
  bottom: 0;
  border-right-color: #2ccff4;
  border-bottom-color: #2ccff4;
  animation-delay: 1.2s;
}

.panel__frame .e {
  background: #2ccff4;
  opacity: 0.72;
  box-shadow: 0 0 6px rgb(80 230 255 / 55%);
  animation: data-screen-edge-pulse 4.2s ease-in-out infinite;
}

.panel__frame .e--t,
.panel__frame .e--b {
  left: 28px;
  width: calc(100% - 56px);
  height: 1px;
}

.panel__frame .e--t {
  top: 5px;
}

.panel__frame .e--b {
  bottom: 5px;
  animation-delay: 1s;
}

.panel__frame .e--l,
.panel__frame .e--r {
  top: 28px;
  width: 1px;
  height: calc(100% - 56px);
}

.panel__frame .e--l {
  left: 5px;
  animation-delay: 0.5s;
}

.panel__frame .e--r {
  right: 5px;
  animation-delay: 1.5s;
}

.panel--today,
.panel--new-old,
.panel--orders,
.panel--merchant,
.panel--grow {
  flex: 1 1 0;
  min-height: 0;
  height: auto;
  overflow: hidden;
}

.panel--payment {
  flex: 1 1 0;
  min-height: 0;
  height: auto;
  /* 图表 x 轴标签需可见 */
  overflow: visible;
}

.panel--rank {
  overflow: hidden;
}

.panel__title {
  position: relative;
  z-index: 1;
  display: flex;
  gap: 10px;
  align-items: center;
  justify-content: center;
  height: 38px;
  margin: 0;
  font-size: 16px;
  font-weight: 900;
  letter-spacing: 2px;
  color: #31abe3;
  text-align: center;
}

.panel__title-text {
  font-weight: 900;
  letter-spacing: 2px;
  background: linear-gradient(180deg, #e8fbff, #31abe3);
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.panel__chevron {
  width: 58px;
  height: 14px;
  object-fit: contain;
  animation: data-screen-chevron-flow 2.4s linear infinite;
}

.panel__chevron--r {
  transform: rotate(180deg);
  animation-direction: reverse;
}

.gauge-grid {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  height: calc(100% - 38px);
  /* 标题下方多留空隙（需求 10） */
  padding: 36px 0 0;
}

.gauge-item {
  position: relative;
}

.gauge-item__ring {
  width: 100px;
  height: 100px;
  margin: 0 auto;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  animation: data-screen-gauge-rotate 14s linear infinite;
}

.gauge-item__value {
  position: absolute;
  z-index: 2;
  top: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100px;
  font-family: 'DIN Alternate', Arial, sans-serif;
  font-size: 22px;
  font-weight: 700;
  color: currentcolor;
  text-align: center;
  text-shadow: 0 0 10px currentcolor;
  pointer-events: none;
}

.gauge-item p {
  margin: 10px 0 0;
  font-size: 14px;
  text-align: center;
}

.new-old {
  position: relative;
  z-index: 1;
  display: grid;
  height: calc(100% - 38px);
  grid-template-columns: 116px 1fr 116px;
  gap: 0;
  align-items: center;
}

.new-old__side {
  text-align: center;
}

.new-old__side small {
  color: #fff;
  font-size: 12px;
}

.new-old__side b {
  display: block;
  margin: 6px 0;
  font-size: 14px;
  color: #e2e8f0;
}

.new-old__side span {
  font-size: 14px;
}

.new-old__side .is-old {
  color: #f5a623;
}

.new-old__side .is-new {
  color: #09caf3;
}

.order-feed-viewport {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  height: calc(100% - 38px);
  min-height: 0;
  overflow: hidden;
}

.order-feed {
  margin: 0;
  padding: 0;
  list-style: none;
}

.order-feed.is-scrolling {
  animation: data-screen-rank-scroll var(--scroll-duration, 18s) linear infinite;
}

.order-feed-viewport:hover .order-feed.is-scrolling {
  animation-play-state: paused;
}

.order-feed li {
  display: flex;
  flex-shrink: 0;
  align-items: flex-start;
  gap: 10px;
  box-sizing: border-box;
  min-height: 0;
  padding: 8px 6px;
  border-bottom: 1px dashed rgb(36 149 218 / 65%);
}

.order-feed .no {
  flex: 0 0 auto;
  padding-top: 4px;
  font-size: 14px;
  font-weight: 700;
  color: #f5d76e;
}

.order-feed__body {
  min-width: 0;
  flex: 1;
}

.order-feed__identity {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-bottom: 4px;
}

.order-feed__entity {
  display: flex;
  gap: 6px;
  align-items: center;
  min-width: 0;
  font-size: 12px;
  color: #a3d9f4;
}

.order-feed__entity img {
  flex: 0 0 auto;
  width: 28px;
  height: 28px;
  object-fit: cover;
  border: 1px solid rgb(71 213 255 / 45%);
  border-radius: 2px;
}

.order-feed__entity span,
.order-feed__entity strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.order-feed__entity strong {
  color: #33d3ff;
}

.order-feed__meta {
  display: grid;
  grid-template-columns: 1.4fr 0.9fr 0.9fr;
  gap: 6px;
}

.order-feed p {
  margin: 0;
  font-size: 12px;
  line-height: 1.5;
  color: #94a3b8;
}

.empty {
  color: #64748b;
  border: 0 !important;
}

.empty--panel {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  min-height: 0;
  margin: 0;
  padding: 0;
  font-size: 12px;
  line-height: 1.5;
  white-space: nowrap;
  text-align: center;
}

.amount-panel {
  flex: 0 0 12vh;
  height: 12vh;
  background: transparent;
}

.amount-panel .amount {
  position: relative;
  z-index: 1;
  padding: 8px 0 0;
  font-family: 'DIN Alternate', Arial, sans-serif;
  font-size: 80px;
  font-weight: 700;
  letter-spacing: 0.06em;
  color: #9fe8ff;
  text-align: center;
  text-shadow: 0 0 22px rgb(9 202 243 / 45%);
}

.map-panel {
  flex: 1 1 0;
  min-height: 0;
  height: auto;
  padding: 0;
  overflow: hidden;
  background: transparent;
}

.map-panel .panel__frame,
.amount-panel .panel__frame {
  display: none;
}

.map-host {
  width: 100%;
  height: 100%;
  min-height: 0;
}

/* 仅改中列底栏：高度对齐左列三等分之一（=实时订单高度），多出的空间留给地图 */
.panel--month {
  flex: 0 0 var(--ds-left-third);
  height: var(--ds-left-third);
  min-height: 0;
  overflow: visible;
}

.panel--payment :deep(.vben-echarts),
.panel--month :deep(.vben-echarts) {
  position: relative;
  z-index: 1;
  box-sizing: border-box;
  width: 100%;
  height: calc(100% - 38px) !important;
  min-height: 0;
  /* 底部 padding 给 x 轴数字完整空间，系列整体上移 */
  padding: 4px 8px 18px;
  overflow: visible;
}

.panel--payment :deep(.vben-echarts) > div,
.panel--month :deep(.vben-echarts) > div {
  overflow: visible !important;
}

.panel--new-old :deep(.vben-echarts) {
  position: relative;
  z-index: 1;
  box-sizing: border-box;
  height: 100% !important;
  min-height: 0;
  padding: 0;
  overflow: hidden;
}

.rank-table-viewport {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  height: calc(100% - 38px);
  overflow: hidden;
}

.rank-body-viewport {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.rank-table-viewport:hover .rank-table--body.is-scrolling {
  animation-play-state: paused;
}

.rank-no {
  font-weight: 700;
  color: #8ec8ef;
}

.rank-no.top {
  color: #f5d76e;
}

.rank-table {
  width: 100%;
  font-size: 12px;
  table-layout: fixed;
  border-collapse: collapse;
}

.rank-table .col-rank {
  width: 44px;
}

.rank-table .col-info {
  width: auto;
}

.rank-table .col-count {
  width: 56px;
}

.rank-table .col-amount {
  width: 78px;
}

.rank-table th,
.rank-table td {
  padding: 6px 4px;
  text-align: left;
  vertical-align: middle;
  border-bottom: 1px solid rgb(9 202 243 / 14%);
}

.rank-table th:first-child,
.rank-table td:first-child {
  padding-left: 6px;
  padding-right: 2px;
  text-align: center;
}

.rank-table--head {
  flex: 0 0 auto;
}

.rank-table--head thead tr {
  background: transparent;
}

.rank-table th {
  height: 30px;
  font-weight: 500;
  color: #93c5fd;
  white-space: nowrap;
  background: transparent;
}

.rank-table--body.is-scrolling {
  animation: data-screen-rank-scroll var(--scroll-duration, 18s) linear infinite;
}

.product-cell {
  display: flex;
  gap: 8px;
  align-items: center;
  min-width: 0;
}

.product-cell__img {
  flex: 0 0 auto;
  width: 28px;
  height: 28px;
  object-fit: cover;
  background: rgb(15 23 42 / 80%);
  border: 1px solid rgb(9 202 243 / 25%);
  border-radius: 2px;
}

.product-cell__text {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@keyframes data-screen-rank-scroll {
  from {
    transform: translateY(0);
  }

  to {
    transform: translateY(-50%);
  }
}

@keyframes data-screen-rail-flow {
  from {
    transform: translateY(-10%);
    opacity: 0.35;
  }

  15% {
    opacity: 1;
  }

  85% {
    opacity: 1;
  }

  to {
    transform: translateY(calc(100vh - 160px));
    opacity: 0.35;
  }
}

@keyframes data-screen-corner-breathe {
  0%,
  100% {
    filter: drop-shadow(0 0 0 rgb(69 220 255 / 0%));
    opacity: 0.55;
  }

  50% {
    filter: drop-shadow(0 0 7px rgb(69 220 255 / 85%));
    opacity: 1;
  }
}

@keyframes data-screen-edge-pulse {
  0%,
  100% {
    opacity: 0.28;
  }

  50% {
    opacity: 0.85;
  }
}

@keyframes data-screen-chevron-flow {
  0%,
  100% {
    filter: brightness(0.75);
    opacity: 0.55;
  }

  50% {
    filter: brightness(1.25) drop-shadow(0 0 4px rgb(80 230 255 / 65%));
    opacity: 1;
  }
}

@keyframes data-screen-gauge-rotate {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 960px) {
  .ds {
    position: relative;
    height: auto;
    min-height: 100vh;
    overflow: auto;
  }

  .ds__body {
    grid-template-columns: 1fr;
  }

  .ds__header {
    grid-template-columns: 1fr;
    gap: 8px;
    height: auto;
    text-align: center;
  }

  .ds__fullscreen {
    justify-self: center;
  }

  .ds__side-rail {
    display: none;
  }
}
</style>
