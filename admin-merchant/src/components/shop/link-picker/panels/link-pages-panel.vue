<script setup lang="ts">
import type { ShopLinkPickerItem } from '#/api/core/shop-link';

import { onMounted, ref } from 'vue';

import { MP_BUILTIN_PAGES } from '../constants';

const emit = defineEmits<{
  change: [ShopLinkPickerItem];
}>();

const activePage = ref<ShopLinkPickerItem | null>(null);

function changeFunc(item: ShopLinkPickerItem) {
  activePage.value = item;
  emit('change', item);
}

onMounted(() => {
  if (MP_BUILTIN_PAGES[0]) changeFunc(MP_BUILTIN_PAGES[0]);
});
</script>

<template>
  <el-select
    class="w-full"
    :model-value="activePage"
    placeholder="请选择"
    value-key="url"
    @change="changeFunc"
  >
    <el-option
      v-for="item in MP_BUILTIN_PAGES"
      :key="item.url"
      :label="item.name"
      :value="item"
    />
  </el-select>
</template>
