<script setup lang="ts">
import type { UserRegionRow } from '#/api/core/platform-dashboard';

import { usePreferences } from '@vben/preferences';
import * as echarts from 'echarts';
import { onBeforeUnmount, onMounted, ref, watch } from 'vue';

const CHINA_GEO_URL = `${import.meta.env.BASE_URL.replace(/\/?$/, '/')}geo/china.json`;

const props = defineProps<{
  list: UserRegionRow[];
}>();

const { isDark } = usePreferences();
const chartRef = ref<HTMLDivElement>();
let chart: echarts.ECharts | null = null;
let mapReady = false;

function normalizeProvince(name: string) {
  const raw = String(name || '').trim();
  if (!raw) return raw;
  if (
    raw.endsWith('省') ||
    raw.endsWith('市') ||
    raw.includes('自治区') ||
    raw.includes('特别行政区')
  ) {
    return raw;
  }
  const special: Record<string, string> = {
    内蒙古: '内蒙古自治区',
    广西: '广西壮族自治区',
    西藏: '西藏自治区',
    宁夏: '宁夏回族自治区',
    新疆: '新疆维吾尔自治区',
    香港: '香港特别行政区',
    澳门: '澳门特别行政区',
  };
  if (special[raw]) return special[raw];
  if (['北京', '天津', '上海', '重庆'].includes(raw)) return `${raw}市`;
  return `${raw}省`;
}

async function ensureChinaMap() {
  if (mapReady) return;
  const res = await fetch(CHINA_GEO_URL);
  const geoJson = await res.json();
  echarts.registerMap('china', geoJson);
  mapReady = true;
}

async function renderChart() {
  if (!chart) return;

  if (!props.list.length) {
    chart.setOption(
      {
        title: {
          left: 'center',
          text: '暂无地域数据',
          textStyle: {
            color: isDark.value ? '#94a3b8' : '#64748b',
            fontSize: 14,
          },
          top: 'middle',
        },
      },
      true,
    );
    chart.resize();
    return;
  }

  await ensureChinaMap();

  const maxValue = Math.max(...props.list.map((item) => item.user_count), 1);
  const data = props.list.map((item) => ({
    name: normalizeProvince(item.province),
    value: item.user_count,
  }));

  chart.setOption(
    {
      geo: {
        itemStyle: {
          areaColor: isDark.value ? '#1e293b' : '#f1f5f9',
          borderColor: isDark.value ? '#334155' : '#cbd5e1',
        },
        label: { show: false },
        map: 'china',
        roam: false,
        zoom: 1.15,
      },
      series: [
        {
          data,
          emphasis: {
            itemStyle: { areaColor: '#f97316' },
            label: { color: '#fff', show: true },
          },
          itemStyle: {
            areaColor: isDark.value ? '#334155' : '#e2e8f0',
            borderColor: isDark.value ? '#475569' : '#cbd5e1',
          },
          label: { show: false },
          map: 'china',
          name: '用户数',
          type: 'map',
        },
      ],
      tooltip: {
        formatter(params: unknown) {
          const item = params as { name?: string; value?: number };
          return `${item.name ?? ''}<br/>用户数：${Number(item.value ?? 0).toLocaleString('zh-CN')}`;
        },
        trigger: 'item',
      },
      visualMap: {
        bottom: 12,
        calculable: false,
        inRange: {
          color: ['#fde68a', '#fb923c', '#ea580c', '#c2410c'],
        },
        left: 12,
        max: maxValue,
        min: 0,
        text: ['高', '低'],
        textStyle: { color: isDark.value ? '#cbd5e1' : '#64748b' },
      },
    },
    true,
  );
  chart.resize();
}

function initChart() {
  if (!chartRef.value) return;
  chart?.dispose();
  chart = echarts.init(chartRef.value, isDark.value ? 'dark' : undefined);
  void renderChart();
}

onMounted(initChart);
watch(() => props.list, () => void renderChart(), { deep: true });
watch(isDark, initChart);

onBeforeUnmount(() => {
  chart?.dispose();
  chart = null;
});
</script>

<template>
  <div ref="chartRef" class="user-region-map" />
</template>

<style scoped>
.user-region-map {
  width: 100%;
  height: 360px;
}
</style>
