<script setup lang="ts">
import type { ShopLinkPickerItem } from '#/api/core/shop-link';

import { onMounted, ref } from 'vue';

import { USER_CENTER_MENU_LINKS } from '../constants';

const emit = defineEmits<{
  change: [ShopLinkPickerItem];
}>();

const activePage = ref<ShopLinkPickerItem | null>(null);

function changeFunc(item: ShopLinkPickerItem) {
  activePage.value = item;
  emit('change', item);
}

onMounted(() => {
  if (USER_CENTER_MENU_LINKS[0]) changeFunc(USER_CENTER_MENU_LINKS[0]);
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
      v-for="item in USER_CENTER_MENU_LINKS"
      :key="item.url"
      :label="item.name"
      :value="item"
    />
  </el-select>
</template>
