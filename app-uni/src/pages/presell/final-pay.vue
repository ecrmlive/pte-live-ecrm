<template>
  <view class="page">
    <view v-if="row" class="card">
      <text class="title">{{ row.store_name || "预售尾款" }}</text>
      <text class="line">尾款单号 {{ row.presell_order_sn }}</text>
      <text class="price">¥{{ Number(row.pay_price).toFixed(2) }}</text>
      <text class="tip">支付窗口内完成尾款后进入待发货</text>
    </view>
    <view v-if="row && row.paid !== 1" class="actions">
      <view class="qx-btn qx-btn-primary" @click="pay('mock')">Mock 支付尾款</view>
      <view class="qx-btn qx-btn-ghost mt" @click="pay('balance')">余额支付尾款</view>
    </view>
    <view v-else-if="row" class="done">尾款已支付</view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { fetchPresellFinal, payPresellFinal, type PresellFinal } from "@/api/presell";

const id = ref(0);
const row = ref<PresellFinal | null>(null);

onLoad(async (q) => {
  id.value = Number(q?.id || 0);
  if (!id.value) return;
  try {
    row.value = await fetchPresellFinal(id.value);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

async function pay(type: "mock" | "balance") {
  try {
    row.value = await payPresellFinal(id.value, type);
    uni.showToast({ title: "尾款支付成功", icon: "success" });
    setTimeout(() => {
      uni.redirectTo({ url: "/pages/order/list" });
    }, 500);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "支付失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 40rpx 24rpx;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 32rpx;
}
.title {
  font-size: 32rpx;
  font-weight: 700;
}
.line {
  display: block;
  margin-top: 12rpx;
  color: #666;
}
.price {
  display: block;
  margin-top: 24rpx;
  font-size: 48rpx;
  font-weight: 700;
  color: #e23030;
}
.tip {
  display: block;
  margin-top: 12rpx;
  color: #999;
  font-size: 24rpx;
}
.actions {
  margin-top: 40rpx;
}
.mt {
  margin-top: 20rpx;
}
.done {
  margin-top: 40rpx;
  text-align: center;
  color: #666;
}
</style>
