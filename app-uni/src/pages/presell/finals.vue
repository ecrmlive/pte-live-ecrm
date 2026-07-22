<template>
  <view class="page">
    <view class="tip">定金已付 · 请在尾款窗口内支付</view>
    <view v-if="!list.length" class="empty">暂无待付尾款</view>
    <view
      v-for="row in list"
      :key="row.presell_order_id"
      class="card"
      @click="goPay(row.presell_order_id)"
    >
      <text class="name">{{ row.store_name || "预售尾款" }}</text>
      <text class="mer">单号 {{ row.presell_order_sn }}</text>
      <view class="price-row">
        <text class="price">¥{{ Number(row.pay_price).toFixed(2) }}</text>
        <text class="link">去支付 ›</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchPresellFinals, type PresellFinal } from "@/api/presell";

const list = ref<PresellFinal[]>([]);

onShow(async () => {
  try {
    const data = await fetchPresellFinals(true);
    list.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goPay(id: number) {
  uni.navigateTo({ url: `/pages/presell/final-pay?id=${id}` });
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
  justify-content: space-between;
  align-items: baseline;
}
.price {
  color: #e23030;
  font-size: 36rpx;
  font-weight: 700;
}
.link {
  color: #666;
  font-size: 26rpx;
}
</style>
