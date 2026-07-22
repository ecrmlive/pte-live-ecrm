<template>
  <view class="page">
    <view class="hero">
      <text class="title">{{ merName || `店铺 ${merId}` }}</text>
      <text class="sub">阶段 2 店铺首页基础 · 本店可售商品</text>
    </view>
    <view class="section">
      <view v-if="!items.length" class="empty">暂无本店上架商品</view>
      <view v-for="p in items" :key="p.id" class="row" @click="goDetail(p.id)">
        <text class="name">{{ p.store_name }}</text>
        <text class="price">¥{{ p.price }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { fetchStoreHome, type ProductItem } from "@/api/catalog";

const merId = ref(0);
const merName = ref("");
const items = ref<ProductItem[]>([]);

onLoad(async (q) => {
  merId.value = Number(q?.mer_id || 0);
  if (!merId.value) return;
  try {
    const data = await fetchStoreHome(merId.value);
    merName.value = data.mer_name || "";
    items.value = data.products || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
}
.hero {
  padding: 48rpx 32rpx;
  background: linear-gradient(135deg, #e23030, #f06a3d);
  color: #fff;
}
.title {
  display: block;
  font-size: 40rpx;
  font-weight: 700;
}
.sub {
  display: block;
  margin-top: 12rpx;
  opacity: 0.9;
  font-size: 24rpx;
}
.section {
  padding: 24rpx;
}
.row {
  background: #fff;
  border-radius: 12rpx;
  padding: 28rpx;
  margin-bottom: 16rpx;
  display: flex;
  justify-content: space-between;
}
.name {
  font-size: 28rpx;
  font-weight: 600;
}
.price {
  color: var(--qx-price);
  font-weight: 700;
}
.empty {
  text-align: center;
  color: var(--qx-text-secondary);
  padding: 60rpx 0;
}
</style>
