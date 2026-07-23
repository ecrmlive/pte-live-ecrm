<script setup lang="ts">
import type { ShopLinkPickerItem } from '#/api/core/shop-link';

import { onMounted, ref, watch } from 'vue';

import { getShopLinkMarketingListsApi } from '#/api/core/shop-link';

import { MARKETING_STATIC_LINKS } from '../constants';

const emit = defineEmits<{
  change: [ShopLinkPickerItem];
}>();

const activeTab = ref('signin');
const activePage = ref<ShopLinkPickerItem | null>(null);
const pages = ref<ShopLinkPickerItem[]>([]);

const packageList = ref<ShopLinkPickerItem[]>([]);
const invitationList = ref<ShopLinkPickerItem[]>([]);
const tableList = ref<ShopLinkPickerItem[]>([]);

function resolvePages(tab: string) {
  if (tab === 'package') return packageList.value;
  if (tab === 'invitation') return invitationList.value;
  if (tab === 'table') return tableList.value;
  return MARKETING_STATIC_LINKS[tab] ?? [];
}

function autoSend() {
  const list = resolvePages(activeTab.value);
  pages.value = list;
  if (list[0]) {
    activePage.value = list[0];
    emit('change', list[0]);
  }
}

function onTabChange() {
  autoSend();
}

function changeFunc(item: ShopLinkPickerItem) {
  activePage.value = item;
  emit('change', item);
}

watch(activeTab, onTabChange);

onMounted(async () => {
  try {
    const res = await getShopLinkMarketingListsApi();
    packageList.value = res.packageList ?? [];
    invitationList.value = res.invitationList ?? [];
    tableList.value = res.tableList ?? [];
  } finally {
    autoSend();
  }
});
</script>

<template>
  <div class="marketing-box">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="签到" name="signin" />
      <el-tab-pane label="积分商城" name="points" />
      <el-tab-pane label="预售" name="presale" />
      <el-tab-pane label="秒杀" name="seckill" />
      <el-tab-pane label="预告" name="preview" />
      <el-tab-pane label="拼团" name="assemble" />
      <el-tab-pane label="砍价" name="bargain" />
      <el-tab-pane :disabled="packageList.length <= 0" label="礼包购" name="package" />
      <el-tab-pane :disabled="invitationList.length <= 0" label="邀请有礼" name="invitation" />
      <el-tab-pane label="优惠券" name="coupon" />
      <el-tab-pane label="幸运转盘" name="lottery" />
      <el-tab-pane :disabled="tableList.length <= 0" label="万能表单" name="table" />
      <el-tab-pane label="任务中心" name="task" />
      <el-tab-pane label="充值" name="recharge" />
    </el-tabs>
    <el-select
      class="w-full"
      :model-value="activePage"
      placeholder="请选择"
      value-key="id"
      @change="changeFunc"
    >
      <el-option v-for="(item, index) in pages" :key="index" :label="item.name" :value="item" />
    </el-select>
  </div>
</template>

<style scoped>
.marketing-box :deep(.el-tabs__item) {
  font-size: 12px;
}
</style>
