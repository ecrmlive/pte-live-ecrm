<template>
  <view class="page">
    <view class="tip">邀请好友助力 · 满员后按助力价下单支付</view>
    <view v-if="!list.length" class="empty">暂无助力活动</view>
    <view v-for="a in list" :key="a.product_assist_id" class="card" @click="goDetail(a)">
      <text class="name">{{ a.store_name || "助力商品" }}</text>
      <text class="mer">{{ a.mer_name }} · 需 {{ a.assist_count }} 人助力</text>
      <view class="price-row">
        <text class="price">¥{{ a.assist_price }}</text>
        <text v-if="a.ot_price" class="ot">¥{{ a.ot_price }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchAssists, type ProductAssist } from "@/api/assist";

const list = ref<ProductAssist[]>([]);

onShow(async () => {
  try {
    const data = await fetchAssists();
    list.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goDetail(a: ProductAssist) {
  uni.navigateTo({ url: `/pages/assist/detail?id=${a.product_assist_id}` });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 24rpx;
  background: #f5f5f5;
}
.tip {
  font-size: 24rpx;
  color: #888;
  margin-bottom: 20rpx;
}
.empty {
  text-align: center;
  color: #999;
  padding: 80rpx 0;
}
.card {
  background: #fff;
  padding: 28rpx;
  margin-bottom: 20rpx;
  border-radius: 12rpx;
}
.name {
  font-size: 30rpx;
  font-weight: 600;
}
.mer {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: #999;
}
.price-row {
  margin-top: 16rpx;
  display: flex;
  gap: 12rpx;
  align-items: baseline;
}
.price {
  color: #e23030;
  font-size: 36rpx;
  font-weight: 700;
}
.ot {
  color: #bbb;
  text-decoration: line-through;
  font-size: 24rpx;
}
</style>
