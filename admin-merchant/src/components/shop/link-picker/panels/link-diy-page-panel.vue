<script setup lang="ts">
import type { ShopLinkPickerItem } from '#/api/core/shop-link';

import { onMounted, ref } from 'vue';

import { getShopLinkDiyPageListApi } from '#/api/core/shop-link';

const emit = defineEmits<{
  change: [ShopLinkPickerItem];
}>();

const loading = ref(false);
const pages = ref<ShopLinkPickerItem[]>([]);
const activePage = ref<ShopLinkPickerItem | null>(null);

function changeFunc(item: ShopLinkPickerItem) {
  activePage.value = item;
  emit('change', item);
}

onMounted(async () => {
  loading.value = true;
  try {
    const res = await getShopLinkDiyPageListApi();
    pages.value = (res.list ?? []).map((item) => ({
      name: item.page_name,
      type: '自定义页面',
      url: `pages/main/diy-page/diy-page?page_id=${item.page_id}`,
    }));
    if (pages.value[0]) changeFunc(pages.value[0]);
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <el-select
    v-loading="loading"
    class="w-full"
    :model-value="activePage"
    placeholder="请选择"
    value-key="url"
    @change="changeFunc"
  >
    <el-option v-for="item in pages" :key="item.url" :label="item.name" :value="item" />
  </el-select>
</template>
