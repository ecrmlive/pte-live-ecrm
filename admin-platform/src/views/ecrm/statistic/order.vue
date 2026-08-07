<script setup lang="ts">
import type { EchartsUIType } from '@vben/plugins/echarts';

import { computed, nextTick, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { EchartsUI, useEcharts } from '@vben/plugins/echarts';
import { CaretBottom, CaretTop } from '@element-plus/icons-vue';
import { ElDatePicker } from 'element-plus';
import 'element-plus/es/components/date-picker/style/css';
import dayjs from 'dayjs';

import {
  getOrderStatLineApi,
  getOrderStatPieApi,
  getOrderStatTopApi,
  type OrderLinePoint,
  type OrderPieSlice,
  type OrderTopCard,
} from '#/api/core/platform-order-stat';

type QuickKey = 'lately7' | 'lately30' | 'month' | 'year';

const quickTabs: Array<{ label: string; value: QuickKey }> = [
  { label: '最近7天', value: 'lately7' },
  { label: '最近30天', value: 'lately30' },
  { label: '本月', value: 'month' },
  { label: '本年', value: 'year' },
];

const moneyTitles = new Set(['订单实付金额', '用券金额', '退款金额']);

const loading = ref(false);
const failed = ref(false);
const timeType = ref<QuickKey | ''>('lately7');
const timeVal = ref<[string, string] | []>([]);
const baseInfo = ref<OrderTopCard[]>([]);
const linePoints = ref<OrderLinePoint[]>([]);
const typePie = ref<OrderPieSlice[]>([]);
const deliveryPie = ref<OrderPieSlice[]>([]);

const lineChartRef = ref<EchartsUIType>();
const typeChartRef = ref<EchartsUIType>();
const deliveryChartRef = ref<EchartsUIType>();
const { renderEcharts: renderLine } = useEcharts(lineChartRef);
const { renderEcharts: renderType } = useEcharts(typeChartRef);
const { renderEcharts: renderDelivery } = useEcharts(deliveryChartRef);

const emptyTop: OrderTopCard[] = [
  { count: 0, mom: 0, statistic: 0, title: '支付订单数' },
  { count: 0, mom: 0, statistic: 0, title: '订单实付金额' },
  { count: 0, mom: 0, statistic: 0, title: '用券金额' },
  { count: 0, mom: 0, statistic: 0, title: '退款金额' },
  { count: 0, mom: 0, statistic: 0, title: '退款订单数' },
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

/** CRMEB：0.00% 也红↑，负值绿↓ */
function isUp(stat: number) {
  return Number(stat) >= 0;
}

function formatCount(card: OrderTopCard) {
  const n = Number(card.count || 0);
  if (moneyTitles.has(card.title)) {
    return n.toLocaleString('zh-CN', {
      maximumFractionDigits: 2,
      minimumFractionDigits: 0,
    });
  }
  return Math.round(n).toLocaleString('zh-CN');
}

function lineOption(points: OrderLinePoint[]) {
  return {
    color: ['#3B82F6', '#34D399', '#F59E0B', '#A78BFA'],
    grid: { bottom: 40, left: 60, right: 60, top: 50 },
    legend: {
      data: ['订单金额', '退款金额', '订单数量', '退款数量'],
      icon: 'circle',
      itemHeight: 8,
      itemWidth: 8,
      left: 'center',
      top: 10,
    },
    series: [
      {
        barWidth: 18,
        data: points.map((p) => Number(p.pay_price) || 0),
        itemStyle: { borderRadius: [4, 4, 0, 0] },
        name: '订单金额',
        type: 'bar',
        yAxisIndex: 0,
      },
      {
        barWidth: 18,
        data: points.map((p) => Number(p.refund_price) || 0),
        itemStyle: { borderRadius: [4, 4, 0, 0] },
        name: '退款金额',
        type: 'bar',
        yAxisIndex: 0,
      },
      {
        data: points.map((p) => Number(p.order_num) || 0),
        name: '订单数量',
        showSymbol: true,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        type: 'line',
        yAxisIndex: 1,
      },
      {
        data: points.map((p) => Number(p.refund_num) || 0),
        name: '退款数量',
        showSymbol: true,
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        type: 'line',
        yAxisIndex: 1,
      },
    ],
    tooltip: { axisPointer: { type: 'cross' }, trigger: 'axis' },
    xAxis: {
      axisLabel: { color: '#666' },
      axisLine: { lineStyle: { color: '#E5E8EF' } },
      axisTick: { alignWithLabel: true },
      boundaryGap: true,
      data: points.map((p) => p.xaxis),
      type: 'category',
    },
    yAxis: [
      {
        axisLabel: { color: '#666', formatter: '{value}' },
        axisLine: { show: false },
        min: 0,
        name: '金额',
        nameTextStyle: { align: 'right', padding: [0, 7, 0, 0] },
        position: 'left',
        splitLine: { lineStyle: { color: '#F0F0F0' }, show: true },
        type: 'value',
      },
      {
        axisLabel: { color: '#666', formatter: '{value}' },
        axisLine: { show: false },
        min: 0,
        name: '数量',
        nameTextStyle: { align: 'left', padding: [0, 0, 0, 7] },
        position: 'right',
        splitLine: { show: false },
        type: 'value',
      },
    ],
  };
}

function pieOption(title: string, data: OrderPieSlice[], colors: string[]) {
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
  await renderType(
    pieOption('订单类型分析', typePie.value, [
      '#3B82F6',
      '#34D399',
      '#60A5FA',
      '#F59E0B',
      '#FB923C',
      '#EF4444',
      '#A78BFA',
      '#F472B6',
    ]),
  );
  await renderDelivery(
    pieOption('发货方式统计', deliveryPie.value, [
      '#3B82F6',
      '#34D399',
      '#60A5FA',
      '#EF4444',
      '#FB923C',
    ]),
  );
}

async function loadAll() {
  loading.value = true;
  failed.value = false;
  const date = resolveDateParam();
  try {
    const [top, line, type, delivery] = await Promise.all([
      getOrderStatTopApi(date),
      getOrderStatLineApi(date),
      getOrderStatPieApi(0, date),
      getOrderStatPieApi(1, date),
    ]);
    baseInfo.value = Array.isArray(top) ? top : [];
    linePoints.value = Array.isArray(line) ? line : [];
    typePie.value = Array.isArray(type) ? type : [];
    deliveryPie.value = Array.isArray(delivery) ? delivery : [];
    await paintCharts();
  } catch {
    failed.value = true;
    baseInfo.value = [];
    linePoints.value = [];
    typePie.value = [];
    deliveryPie.value = [];
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
    <div v-loading="loading" class="order-stat">
      <el-alert
        v-if="failed"
        class="mb-3"
        title="统计数据暂不可用，请刷新重试。"
        type="warning"
        :closable="false"
      />

      <!-- 1. 时间选择（日期宽度/间距对齐商品统计；选中蓝底白字对齐 CRMEB） -->
      <div class="filter-card">
        <div class="filter-row">
          <span class="filter-label">时间选择：</span>
          <div class="filter-date-wrap">
            <ElDatePicker
              :model-value="timeVal.length === 2 ? timeVal : undefined"
              type="daterange"
              format="YYYY/MM/DD"
              value-format="YYYY/MM/DD"
              range-separator="-"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
              teleported
              @update:model-value="onRangeChange"
            />
          </div>
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

      <!-- 2. 一行 5 张等宽 KPI 小卡 -->
      <div class="kpi-row">
        <div v-for="(card, idx) in kpiCards" :key="idx" class="kpi-card">
          <div class="kpi-title">{{ card.title }}</div>
          <div class="kpi-count">{{ formatCount(card) }}</div>
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

      <!-- 3. 全宽双轴趋势图 -->
      <div class="chart-card line-card">
        <EchartsUI ref="lineChartRef" height="380px" />
      </div>

      <!-- 4. 底栏两卡：订单类型 | 发货方式 -->
      <div class="pie-row">
        <div class="chart-card pie-card">
          <EchartsUI ref="typeChartRef" height="360px" />
        </div>
        <div class="chart-card pie-card">
          <EchartsUI ref="deliveryChartRef" height="360px" />
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.order-stat {
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
  gap: 8px;
}

.filter-label {
  color: #606266;
  flex-shrink: 0;
  font-size: 14px;
  line-height: 32px;
}

/* 仅锁输入框宽度；弹层 teleported 到 body，勿用 :deep 碰 .el-picker-panel */
.filter-date-wrap {
  flex-shrink: 0;
  width: 300px;
}

.filter-date-wrap :deep(.el-date-editor.el-input__wrapper),
.filter-date-wrap :deep(.el-date-editor) {
  --el-date-editor-daterange-width: 300px;
  --el-date-editor-width: 300px;
  box-sizing: border-box;
  height: 32px;
  justify-content: flex-start;
  max-width: 300px;
  padding: 0 8px;
  width: 300px !important;
}

.filter-date-wrap :deep(.el-date-editor .el-range__icon) {
  float: none;
  flex-shrink: 0;
  margin-inline-end: 4px;
}

.filter-date-wrap :deep(.el-date-editor .el-range-input) {
  flex: 0 0 auto;
  text-align: left;
  width: 5.6em;
}

.filter-date-wrap :deep(.el-date-editor .el-range-separator) {
  flex: 0 0 auto;
  padding: 0 4px;
  width: auto;
}

.filter-date-wrap :deep(.el-date-editor .el-range__close-icon) {
  flex-shrink: 0;
  margin-inline-start: 2px;
}

.quick-tabs {
  align-items: center;
  display: inline-flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-inline-start: 6px;
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
    background-color 0.15s,
    border-color 0.15s,
    color 0.15s;
}

.quick-tab:hover {
  border-color: var(--el-color-primary, #409eff);
  color: var(--el-color-primary, #409eff);
}

/* CRMEB el-radio-button 选中：蓝底白字 */
.quick-tab.is-active {
  background: var(--el-color-primary, #409eff);
  border-color: var(--el-color-primary, #409eff);
  color: #fff;
}

.quick-tab.is-active:hover {
  background: var(--el-color-primary, #409eff);
  border-color: var(--el-color-primary, #409eff);
  color: #fff;
}

.kpi-row {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(5, minmax(0, 1fr));
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

@media (max-width: 640px) {
  .kpi-row {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .pie-row {
    grid-template-columns: 1fr;
  }
}
</style>
