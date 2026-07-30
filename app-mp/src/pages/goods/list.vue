<template>
  <view class="page">
    <view class="search-bar">
      <input
        v-model="keyword"
        class="input"
        placeholder="搜索商品"
        confirm-type="search"
        @confirm="load"
      />
      <view class="btn" @click="load">搜索</view>
    </view>
    <scroll-view scroll-y class="list">
      <view v-if="!items.length" class="empty">暂无商品</view>
      <view v-for="p in items" :key="p.id" class="row" @click="goDetail(p.id)">
        <view class="cover">
          <text>{{ p.store_name.slice(0, 1) }}</text>
        </view>
        <view class="meta">
          <text class="name">{{ p.store_name }}</text>
          <text class="mer" @click.stop="goStore(p.mer_id)">{{ p.mer_name }}</text>
          <view class="price-row">
            <text class="price">¥{{ p.price }}</text>
            <text class="sales">已售 {{ p.sales }}</text>
          </view>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { fetchProducts, type ProductItem } from "@/api/catalog";

const keyword = ref("");
const cateId = ref(0);
const items = ref<ProductItem[]>([]);

onLoad((q) => {
  cateId.value = Number(q?.cate_id || 0);
  load();
});

async function load() {
  try {
    const data = await fetchProducts({
      cate_id: cateId.value || undefined,
      keyword: keyword.value.trim() || undefined,
      page: 1,
    });
    items.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` });
}
function goStore(merId: number) {
  uni.navigateTo({ url: `/pages/store/home?mer_id=${merId}` });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
}
.search-bar {
  display: flex;
  gap: 16rpx;
  padding: 16rpx 24rpx;
  background: #fff;
}
.input {
  flex: 1;
  height: 72rpx;
  background: #f5f5f5;
  border-radius: 12rpx;
  padding: 0 24rpx;
  font-size: 28rpx;
}
.btn {
  height: 72rpx;
  line-height: 72rpx;
  padding: 0 28rpx;
  background: var(--qx-brand);
  color: #fff;
  border-radius: 12rpx;
  font-size: 28rpx;
}
.list {
  height: calc(100vh - 104rpx);
  padding: 16rpx 24rpx;
  box-sizing: border-box;
}
.row {
  display: flex;
  gap: 20rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 20rpx;
  margin-bottom: 16rpx;
}
.cover {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  background: #f3e9e6;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c9a39a;
  font-size: 48rpx;
  font-weight: 700;
}
.meta {
  flex: 1;
}
.name {
  font-size: 28rpx;
  font-weight: 600;
}
.mer {
  display: block;
  margin-top: 8rpx;
  color: var(--qx-text-secondary);
  font-size: 22rpx;
}
.price-row {
  margin-top: 16rpx;
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}
.price {
  color: var(--qx-price);
  font-size: 32rpx;
  font-weight: 700;
}
.sales {
  color: #bbb;
  font-size: 22rpx;
}
.empty {
  text-align: center;
  color: var(--qx-text-secondary);
  padding: 80rpx 0;
}
</style>
