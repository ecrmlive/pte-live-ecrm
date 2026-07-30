<template>
  <view class="page">
    <text class="icon">{{ paid ? "✓" : "!" }}</text>
    <text class="title">{{ paid ? "支付成功" : "下单完成" }}</text>
    <text class="desc">主单 #{{ id }}</text>
    <view class="qx-btn qx-btn-primary btn" @click="goDetail">查看订单</view>
    <view class="qx-btn qx-btn-ghost btn" @click="goHome">回首页</view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";

const id = ref(0);
const paid = ref(0);

onLoad((q) => {
  id.value = Number(q?.id || 0);
  paid.value = Number(q?.paid || 0);
});

function goDetail() {
  uni.redirectTo({ url: `/pages/order/detail?id=${id.value}` });
}
function goHome() {
  uni.switchTab({ url: "/pages/index/index" });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: var(--qx-bg);
  padding: 40rpx;
}
.icon {
  width: 120rpx;
  height: 120rpx;
  border-radius: 60rpx;
  background: #e8f8ef;
  color: #1a9f5b;
  text-align: center;
  line-height: 120rpx;
  font-size: 64rpx;
  font-weight: 700;
}
.title {
  margin-top: 32rpx;
  font-size: 40rpx;
  font-weight: 700;
}
.desc {
  margin-top: 12rpx;
  color: #999;
}
.btn {
  margin-top: 28rpx;
  width: 360rpx;
}
</style>
