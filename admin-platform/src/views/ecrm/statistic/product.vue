<script setup lang="ts">
import type { EchartsUIType } from '@vben/plugins/echarts';

import { computed, nextTick, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { EchartsUI, useEcharts } from '@vben/plugins/echarts';
import { CaretBottom, CaretTop } from '@element-plus/icons-vue';
import dayjs from 'dayjs';

import {
  getProductStatLineApi,
  getProductStatPieApi,
  getProductStatTopApi,
  type ProductLinePoint,
  type ProductPieSlice,
  type ProductTopCard,
} from '#/api/core/platform-product-stat';

type QuickKey = 'lately7' | 'lately30' | 'month' | 'year';

const quickTabs: Array<{ label: string; value: QuickKey }> = [
  { label: '最近7天', value: 'lately7' },
  { label: '最近30天', value: 'lately30' },
  { label: '本月', value: 'month' },
  { label: '本年', value: 'year' },
];

const loading = ref(false);
const failed = ref(false);
const timeType = ref<QuickKey | ''>('lately7');
const timeVal = ref<[string, string] | []>([]);
const baseInfo = ref<ProductTopCard[]>([]);
const linePoints = ref<ProductLinePoint[]>([]);
const catePie = ref<ProductPieSlice[]>([]);
const typePie = ref<ProductPieSlice[]>([]);

const lineChartRef = ref<EchartsUIType>();
const cateChartRef = ref<EchartsUIType>();
const typeChartRef = ref<EchartsUIType>();
const { renderEcharts: renderLine } = useEcharts(lineChartRef);
const { renderEcharts: renderCate } = useEcharts(cateChartRef);
const { renderEcharts: renderType } = useEcharts(typeChartRef);

const emptyTop: ProductTopCard[] = [
  { count: 0, mom: 0, statistic: 0, title: '浏览量' },
  { count: 0, mom: 0, statistic: 0, title: '收藏量' },
  { count: 0, mom: 0, statistic: 0, title: '加购数' },
  { count: 0, mom: 0, statistic: 0, title: '下单数' },
  { count: 0, mom: 0, statistic: 0, title: '支付数' },
  { count: 0, mom: 0, statistic: 0, title: '退款数' },
];

const kpiCards = computed(() => (baseInfo.value.length ? baseInfo.value : emptyTop));

function calculateDateRange(key: QuickKey): [string, string] {
  const now = dayjs();
  switch (key) {
    case 'lately7':
      return [now.subtract(6, 'day').format('YYYY/MM/DD'), now.format('YYYY/MM/DD')];
    case 'lately30':
      return [now.subtract(29, 'day').format('YYYY/MM/DD'), now.format('YYYY/MM/DD')];
    case 'month':
      return [now.startOf('month').format('YYYY/MM/DD'), now.endOf('month').format('YYYY/MM/DD')];
    case 'year':
      return [now.startOf('year').format('YYYY/MM/DD'), now.endOf('year').format('YYYY/MM/DD')];
  }
}

function resolveDateParam(): string {
  if (timeType.value) return timeType.value;
  if (timeVal.value.length === 2) return `${timeVal.value[0]}-${timeVal.value[1]}`;
  return 'lately7';
}

function formatStat(n: number | null | undefined) {
  return Number(n || 0).toFixed(2);
}

function isUp(stat: number) {
  return Number(stat) > 0;
}

function lineOption(points: ProductLinePoint[]) {
  return {
    color: ['#5B8FF9', '#5AD8A6', '#F6BD16', '#A78BFA'],
    grid: { bottom: 40, left: 50, right: 30, top: 50 },
    legend: {
      data: ['商品浏览量', '收藏量', '下单量', '支付量'],
      icon: 'circle',
      itemHeight: 8,
      itemWidth: 8,
      left: 'center',
      top: 10,
    },
    series: [
      {
        data: points.map((p) => p.visit),
        lineStyle: { width: 2 },
        name: '商品浏览量',
        showSymbol: false,
        smooth: true,
        type: 'line',
      },
      {
        data: points.map((p) => p.relation),
        lineStyle: { width: 2 },
        name: '收藏量',
        showSymbol: false,
        smooth: true,
        type: 'line',
      },
      {
        data: points.map((p) => p.total_num),
        lineStyle: { width: 2 },
        name: '下单量',
        showSymbol: false,
        smooth: true,
        type: 'line',
      },
      {
        data: points.map((p) => p.paid_num),
        lineStyle: { width: 2 },
        name: '支付量',
        showSymbol: false,
        smooth: true,
        type: 'line',
      },
    ],
    tooltip: { axisPointer: { type: 'line' }, trigger: 'axis' },
    xAxis: {
      axisLabel: { color: '#666' },
      axisLine: { lineStyle: { color: '#E5E8EF' } },
      boundaryGap: false,
      data: points.map((p) => p.xaxis),
      type: 'category',
    },
    yAxis: {
      axisLabel: { color: '#666' },
      axisLine: { show: false },
      min: 0,
      name: '数量',
      nameTextStyle: { align: 'right', padding: [0, 7, 0, 0] },
      splitLine: { lineStyle: { color: '#F0F0F0' } },
      type: 'value',
    },
  };
}

function pieOption(title: string, data: ProductPieSlice[], colors: string[]) {
  const hasData = data.some((d) => Number(d.value) > 0);
  const seriesData = hasData
    ? data.map((d) => ({ name: d.name, value: Number(d.value) || 0 }))
    : [{ itemStyle: { color: '#D1D5DB' }, name: title, value: 1 }];
  return {
    color: colors,
    legend: {
      icon: 'circle',
      itemHeight: 8,
      itemWidth: 8,
      left: 20,
      orient: 'vertical' as const,
      top: 50,
    },
    series: [
      {
        avoidLabelOverlap: true,
        center: ['60%', '55%'],
        data: seriesData,
        label: hasData
          ? { formatter: '{b}', position: 'outside', show: true }
          : { show: false },
        labelLine: hasData
          ? { length: 14, length2: 10, show: true, smooth: true }
          : { show: false },
        name: title,
        radius: ['40%', '70%'],
        silent: !hasData,
        type: 'pie',
      },
    ],
    title: {
      left: 20,
      text: title,
      textStyle: { fontSize: 18, fontWeight: 600 },
      top: 10,
    },
    tooltip: hasData
      ? { formatter: '{b}: {c} ({d}%)', trigger: 'item' }
      : { show: false },
  };
}

async function paintCharts() {
  await nextTick();
  await renderLine(lineOption(linePoints.value));
  await renderCate(
    pieOption('商品分类', catePie.value, [
      '#3B82F6',
      '#34D399',
      '#60A5FA',
      '#F59E0B',
      '#FB923C',
      '#EF4444',
      '#A78BFA',
      '#F472B6',
      '#C084FC',
    ]),
  );
  await renderType(
    pieOption('商品类型', typePie.value, ['#3B82F6', '#34D399', '#60A5FA', '#EF4444', '#FB923C']),
  );
}

async function loadAll() {
  loading.value = true;
  failed.value = false;
  const date = resolveDateParam();
  try {
    const [top, line, cate, type] = await Promise.all([
      getProductStatTopApi(date),
      getProductStatLineApi(date),
      getProductStatPieApi(1),
      getProductStatPieApi(0),
    ]);
    baseInfo.value = Array.isArray(top) ? top : [];
    linePoints.value = Array.isArray(line) ? line : [];
    catePie.value = Array.isArray(cate) ? cate : [];
    typePie.value = Array.isArray(type) ? type : [];
    await paintCharts();
  } catch {
    failed.value = true;
    baseInfo.value = [];
    linePoints.value = [];
    catePie.value = [];
    typePie.value = [];
    await paintCharts();
  } finally {
    loading.value = false;
  }
}

function onQuickChange(key: QuickKey) {
  timeType.value = key;
  timeVal.value = calculateDateRange(key);
  void loadAll();
}

function onRangeChange(val: [string, string] | string[] | null | undefined) {
  if (!val || val.length !== 2 || !val[0] || !val[1]) {
    timeType.value = 'lately7';
    timeVal.value = calculateDateRange('lately7');
    void loadAll();
    return;
  }
  timeVal.value = [val[0], val[1]];
  timeType.value = '';
  void loadAll();
}

onMounted(() => {
  timeVal.value = calculateDateRange('lately7');
  void loadAll();
});
</script>

<template>
  <Page content-class="!bg-[#f0f2f5]" title="">
    <div v-loading="loading" class="product-stat">
      <el-alert
        v-if="failed"
        class="mb-3"
        title="统计数据暂不可用，请刷新重试。"
        type="warning"
        :closable="false"
      />

      <!-- 1. 时间选择（对齐 CRMEB：日期范围 + 快捷蓝边选中） -->
      <div class="filter-card">
        <div class="filter-row">
          <span class="filter-label">时间选择：</span>
          <el-date-picker
            :model-value="timeVal.length === 2 ? timeVal : undefined"
            type="daterange"
            format="YYYY/MM/DD"
            value-format="YYYY/MM/DD"
            range-separator="-"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            class="filter-date"
            @update:model-value="onRangeChange"
          />
          <div class="quick-tabs">
            <button
              v-for="tab in quickTabs"
              :key="tab.value"
              type="button"
              class="quick-tab"
              :class="{ 'is-active': timeType === tab.value }"
              @click="onQuickChange(tab.value)"
            >
              {{ tab.label }}
            </button>
          </div>
        </div>
      </div>

      <!-- 2. 一行 6 张等宽 KPI 小卡（禁止竖排全宽） -->
      <div class="kpi-row">
        <div v-for="(card, idx) in kpiCards" :key="idx" class="kpi-card">
          <div class="kpi-title">{{ card.title }}</div>
          <div class="kpi-count">{{ Number(card.count || 0).toLocaleString('zh-CN') }}</div>
          <div class="kpi-stat">
            环比增长:
            <span :class="isUp(card.statistic) ? 'text-up' : 'text-down'">
              {{ formatStat(card.statistic) }}%
            </span>
            <el-icon :class="isUp(card.statistic) ? 'text-up' : 'text-down'" class="kpi-arrow">
              <CaretTop v-if="isUp(card.statistic)" />
              <CaretBottom v-else />
            </el-icon>
          </div>
        </div>
      </div>

      <!-- 3. 全宽折线图 -->
      <div class="chart-card line-card">
        <EchartsUI ref="lineChartRef" height="280px" />
      </div>

      <!-- 4. 底栏两卡：商品分类 | 商品类型 -->
      <div class="pie-row">
        <div class="chart-card pie-card">
          <EchartsUI ref="cateChartRef" height="360px" />
        </div>
        <div class="chart-card pie-card">
          <EchartsUI ref="typeChartRef" height="360px" />
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.product-stat {
  min-height: 100%;
}

.filter-card {
  background: #fff;
  border-radius: 4px;
  margin-bottom: 12px;
  padding: 12px 16px;
}

.filter-row {
  align-items: center;
  display: flex;
  flex-wrap: wrap;
  gap: 8px 10px;
}

.filter-label {
  color: #606266;
  flex-shrink: 0;
  font-size: 14px;
  line-height: 32px;
}

/* 覆盖 EP daterange：默认 350px + text-align:center + separator flex:1 导致左侧大空 */
.filter-date.el-date-editor {
  --el-date-editor-daterange-width: 246px;
  --el-date-editor-width: 246px;
  box-sizing: border-box;
  flex-shrink: 0;
  height: 32px;
  justify-content: flex-start;
  padding: 0 8px;
  width: 246px;
}

.filter-date :deep(.el-range__icon) {
  flex-shrink: 0;
  margin-inline-end: 6px;
}

.filter-date :deep(.el-range-input) {
  flex: 1 1 0;
  min-width: 0;
  text-align: left;
  width: auto;
}

.filter-date :deep(.el-range-separator) {
  flex: 0 0 auto;
  padding: 0 4px;
  width: auto;
}

.filter-date :deep(.el-range__close-icon) {
  flex-shrink: 0;
  margin-inline-start: 4px;
}

.quick-tabs {
  align-items: center;
  display: inline-flex;
  flex-wrap: wrap;
  gap: 8px;
}

.quick-tab {
  background: #fff;
  border: 1px solid #dcdfe6;
  border-radius: 2px;
  box-sizing: border-box;
  color: #606266;
  cursor: pointer;
  font-size: 13px;
  height: 32px;
  line-height: 30px;
  padding: 0 14px;
  transition:
    border-color 0.15s,
    color 0.15s;
}

.quick-tab:hover {
  border-color: var(--el-color-primary, #409eff);
  color: var(--el-color-primary, #409eff);
}

.quick-tab.is-active {
  background: #fff;
  border-color: var(--el-color-primary, #409eff);
  color: var(--el-color-primary, #409eff);
}

/* 强制一行 6 等宽，对齐 CRMEB 图片2 */
.kpi-row {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  margin-bottom: 12px;
}

.kpi-card {
  background: #fff;
  border-radius: 4px;
  box-sizing: border-box;
  min-width: 0;
  padding: 16px 14px 14px;
}

.kpi-title {
  color: #909399;
  font-size: 13px;
  line-height: 1.2;
}

.kpi-count {
  color: #303133;
  font-size: 26px;
  font-weight: 600;
  line-height: 1.25;
  margin-top: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kpi-stat {
  align-items: center;
  color: #909399;
  display: flex;
  flex-wrap: wrap;
  font-size: 12px;
  gap: 2px;
  margin-top: 10px;
}

.kpi-arrow {
  font-size: 12px;
}

.text-up {
  color: #f56c6c;
}

.text-down {
  color: #67c23a;
}

.chart-card {
  background: #fff;
  border-radius: 4px;
  margin-bottom: 12px;
  padding: 8px 4px 4px;
}

.line-card {
  padding: 8px 12px 4px;
}

.pie-row {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.pie-card {
  margin-bottom: 0;
  min-width: 0;
}

/* 仅极窄屏折行；桌面始终 6 列横排，禁止全宽竖堆 */
@media (max-width: 640px) {
  .kpi-row {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .pie-row {
    grid-template-columns: 1fr;
  }
}
</style>
