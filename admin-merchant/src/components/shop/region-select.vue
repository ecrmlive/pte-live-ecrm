<script setup lang="ts">
import type { RegionTree } from '#/api/core/data';

import { computed, ref, watch } from 'vue';

const provinceId = defineModel<number | string>('provinceId', { default: '' });
const cityId = defineModel<number | string>('cityId', { default: '' });
const regionId = defineModel<number | string>('regionId', { default: '' });

const props = defineProps<{
  regionData: RegionTree;
}>();

const internalProvince = ref<number | string>('');
const internalCity = ref<number | string>('');
const internalRegion = ref<number | string>('');

watch(
  () => [provinceId.value, cityId.value, regionId.value],
  () => {
    internalProvince.value = provinceId.value;
    internalCity.value = cityId.value;
    internalRegion.value = regionId.value;
  },
  { immediate: true },
);

const provinces = computed(() =>
  Object.values(props.regionData).map((item) => ({ id: item.id, name: item.name })),
);

const cities = computed(() => {
  if (!internalProvince.value) return [];
  const province = props.regionData[Number(internalProvince.value)];
  if (!province?.city) return [];
  return Object.values(province.city).map((item) => ({ id: item.id, name: item.name }));
});

const regions = computed(() => {
  if (!internalProvince.value || !internalCity.value) return [];
  const province = props.regionData[Number(internalProvince.value)];
  const city = province?.city?.[Number(internalCity.value)];
  if (!city?.region) return [];
  return Object.values(city.region).map((item) => ({ id: item.id, name: item.name }));
});

function onProvinceChange() {
  internalCity.value = '';
  internalRegion.value = '';
  provinceId.value = internalProvince.value;
  cityId.value = '';
  regionId.value = '';
}

function onCityChange() {
  internalRegion.value = '';
  cityId.value = internalCity.value;
  regionId.value = '';
}

function onRegionChange() {
  regionId.value = internalRegion.value;
}
</script>

<template>
  <div class="flex flex-col gap-3">
    <el-select
      v-model="internalProvince"
      class="max-w-[460px]"
      placeholder="省"
      @change="onProvinceChange"
    >
      <el-option v-for="item in provinces" :key="item.id" :label="item.name" :value="item.id" />
    </el-select>
    <el-select
      v-if="internalProvince"
      v-model="internalCity"
      class="max-w-[460px]"
      placeholder="市"
      @change="onCityChange"
    >
      <el-option v-for="item in cities" :key="item.id" :label="item.name" :value="item.id" />
    </el-select>
    <el-select
      v-if="internalCity"
      v-model="internalRegion"
      class="max-w-[460px]"
      placeholder="区"
      @change="onRegionChange"
    >
      <el-option v-for="item in regions" :key="item.id" :label="item.name" :value="item.id" />
    </el-select>
  </div>
</template>

<style scoped>
.max-w-\[460px\] {
  max-width: 460px;
}
</style>
