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
  getUserStatDealApi,
  getUserStatLineApi,
  getUserStatTopApi,
  type UserDealPoint,
  type UserLinePoint,
  type UserTopCard,
} from '#/api/core/platform-user-stat';

type QuickKey = 'lately7' | 'lately30' | 'month' | 'year';

const quickTabs: Array<{ label: string; value: QuickKey }> = [
  { label: '最近7天', value: 'lately7' },
  { label: '最近30天', value: 'lately30' },
  { label: '本月', value: 'month' },
  { label: '本年', value: 'year' },
];

const LINE_BLUE = '#409EFF';
const AREA_BLUE = 'rgba(64, 158, 255, 0.18)';
const BAR_OLD = '#409EFF';
const BAR_NEW = '#67C23A';

const loading = ref(false);
const failed = ref(false);
const timeType = ref<QuickKey | ''>('lately7');
const timeVal = ref<[string, string] | []>([]);
const baseInfo = ref<UserTopCard[]>([]);
const newUserLine = ref<UserLinePoint[]>([]);
const activeUserLine = ref<UserLinePoint[]>([]);
const dealPoints = ref<UserDealPoint[]>([]);
const svipLine = ref<UserLinePoint[]>([]);

const newChartRef = ref<EchartsUIType>();
const activeChartRef = ref<EchartsUIType>();
const dealChartRef = ref<EchartsUIType>();
const svipChartRef = ref<EchartsUIType>();
const { renderEcharts: renderNew } = useEcharts(newChartRef);
const { renderEcharts: renderActive } = useEcharts(activeChartRef);
const { renderEcharts: renderDeal } = useEcharts(dealChartRef);
const { renderEcharts: renderSvip } = useEcharts(svipChartRef);

const emptyTop: UserTopCard[] = [
  { count: 0, mom: 0, statistic: 0, title: '用户数量' },
  { count: 0, mom: 0, statistic: 0, title: '新增用户' },
  { count: 0, mom: 0, statistic: 0, title: '下单用户' },
  { count: 0, mom: 0, statistic: 0, title: '活跃用户' },
  { count: 0, mom: 0, statistic: 0, title: '付费会员' },
  { count: 0, mom: 0, statistic: 0, title: '新增付费会员' },
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

function areaLineOption(title: string, points: UserLinePoint[]) {
  return {
    color: [LINE_BLUE],
    grid: { bottom: 36, left: 48, right: 24, top: 56 },
    series: [
      {
        areaStyle: { color: AREA_BLUE },
        data: points.map((p) => p.count),
        itemStyle: { color: LINE_BLUE },
        lineStyle: { color: LINE_BLUE, width: 2 },
        name: title,
        showSymbol: true,
        smooth: true,
        symbol: 'circle',
        symbolSize: 7,
        type: 'line',
      },
    ],
    title: {
      left: 12,
      text: title,
      textStyle: { color: '#303133', fontSize: 15, fontWeight: 600 },
      top: 8,
    },
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
      minInterval: 1,
      name: '数量',
      nameTextStyle: { align: 'right', color: '#909399', padding: [0, 7, 0, 0] },
      splitLine: { lineStyle: { color: '#F0F0F0' } },
      type: 'value',
    },
  };
}

function dealBarOption(points: UserDealPoint[]) {
  return {
    color: [BAR_OLD, BAR_NEW],
    grid: { bottom: 36, left: 48, right: 24, top: 56 },
    legend: {
      data: ['老用户', '新用户'],
      icon: 'rect',
      itemHeight: 10,
      itemWidth: 14,
      left: 'center',
      top: 10,
    },
    series: [
      {
        barGap: '30%',
        barMaxWidth: 18,
        data: points.map((p) => p.old),
        itemStyle: { color: BAR_OLD },
        name: '老用户',
        type: 'bar',
      },
      {
        barMaxWidth: 18,
        data: points.map((p) => p.new),
        itemStyle: { color: BAR_NEW },
        name: '新用户',
        type: 'bar',
      },
    ],
    title: {
      left: 12,
      text: '成交用户数量',
      textStyle: { color: '#303133', fontSize: 15, fontWeight: 600 },
      top: 8,
    },
    tooltip: { axisPointer: { type: 'shadow' }, trigger: 'axis' },
    xAxis: {
      axisLabel: { color: '#666' },
      axisLine: { lineStyle: { color: '#E5E8EF' } },
      data: points.map((p) => p.xaxis),
      type: 'category',
    },
    yAxis: {
      axisLabel: { color: '#666' },
      axisLine: { show: false },
      min: 0,
      minInterval: 1,
      name: '数量',
      nameTextStyle: { align: 'right', color: '#909399', padding: [0, 7, 0, 0] },
      splitLine: { lineStyle: { color: '#F0F0F0' } },
      type: 'value',
    },
  };
}

function plainLineOption(title: string, points: UserLinePoint[]) {
  return {
    color: [LINE_BLUE],
    grid: { bottom: 36, left: 48, right: 24, top: 56 },
    series: [
      {
        data: points.map((p) => p.count),
        itemStyle: { color: LINE_BLUE },
        lineStyle: { color: LINE_BLUE, width: 2 },
        name: title,
        showSymbol: true,
        smooth: false,
        symbol: 'circle',
        symbolSize: 7,
        type: 'line',
      },
    ],
    title: {
      left: 12,
      text: title,
      textStyle: { color: '#303133', fontSize: 15, fontWeight: 600 },
      top: 8,
    },
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
      minInterval: 1,
      name: '数量',
      nameTextStyle: { align: 'right', color: '#909399', padding: [0, 7, 0, 0] },
      splitLine: { lineStyle: { color: '#F0F0F0' } },
      type: 'value',
    },
  };
}

async function paintCharts() {
  await nextTick();
  await renderNew(areaLineOption('新增用户数量', newUserLine.value));
  await renderActive(areaLineOption('活跃用户数量', activeUserLine.value));
  await renderDeal(dealBarOption(dealPoints.value));
  await renderSvip(plainLineOption('新增付费会员数量', svipLine.value));
}

async function loadAll() {
  loading.value = true;
  failed.value = false;
  const date = resolveDateParam();
  try {
    const [top, neu, active, deal, svip] = await Promise.all([
      getUserStatTopApi(date),
      getUserStatLineApi(date, 0),
      getUserStatLineApi(date, 1),
      getUserStatDealApi(date),
      getUserStatLineApi(date, 2),
    ]);
    baseInfo.value = Array.isArray(top) ? top : [];
    newUserLine.value = Array.isArray(neu) ? neu : [];
    activeUserLine.value = Array.isArray(active) ? active : [];
    dealPoints.value = Array.isArray(deal) ? deal : [];
    svipLine.value = Array.isArray(svip) ? svip : [];
    await paintCharts();
  } catch {
    failed.value = true;
    baseInfo.value = [];
    newUserLine.value = [];
    activeUserLine.value = [];
    dealPoints.value = [];
    svipLine.value = [];
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
    <div v-loading="loading" class="user-stat">
      <el-alert
        v-if="failed"
        class="mb-3"
        title="统计数据暂不可用，请刷新重试。"
        type="warning"
        :closable="false"
      />

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

      <div class="chart-row">
        <div class="chart-card">
          <EchartsUI ref="newChartRef" height="300px" />
        </div>
        <div class="chart-card">
          <EchartsUI ref="activeChartRef" height="300px" />
        </div>
      </div>

      <div class="chart-row">
        <div class="chart-card">
          <EchartsUI ref="dealChartRef" height="300px" />
        </div>
        <div class="chart-card">
          <EchartsUI ref="svipChartRef" height="300px" />
        </div>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.user-stat {
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

.chart-row {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  margin-bottom: 12px;
}

.chart-card {
  background: #fff;
  border-radius: 4px;
  min-width: 0;
  padding: 4px 8px 8px;
}

@media (max-width: 640px) {
  .kpi-row {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .chart-row {
    grid-template-columns: 1fr;
  }
}
</style>
