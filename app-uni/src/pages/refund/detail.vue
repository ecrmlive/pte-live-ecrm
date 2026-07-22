<template>
  <view class="page" v-if="item">
    <view class="card">
      <text class="st">{{ refundStatusText(item.status) }}</text>
      <text class="line">单号 {{ item.refund_order_sn }}</text>
      <text class="line">金额 ¥{{ Number(item.refund_price).toFixed(2) }}</text>
      <text class="line">原因 {{ item.refund_message }}</text>
      <text v-if="item.fail_message" class="line muted">拒绝：{{ item.fail_message }}</text>
    </view>
    <view v-if="item.status === 0" class="foot">
      <view class="qx-btn" @click="onPlatform">申请平台介入</view>
      <view class="qx-btn qx-btn-ghost" @click="onCancel">取消申请</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import {
  cancelRefund,
  fetchRefund,
  refundStatusText,
  requestPlatformRefund,
  type RefundOrder,
} from "@/api/refund";

const item = ref<RefundOrder | null>(null);
const id = ref(0);

onLoad(async (q) => {
  id.value = Number(q?.id || 0);
  await load();
});

async function load() {
  if (!id.value) return;
  try {
    item.value = await fetchRefund(id.value);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

async function onCancel() {
  try {
    await cancelRefund(id.value);
    uni.showToast({ title: "已取消", icon: "success" });
    await load();
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "取消失败", icon: "none" });
  }
}

async function onPlatform() {
  try {
    await requestPlatformRefund(id.value);
    uni.showToast({ title: "已申请平台介入", icon: "success" });
    await load();
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "申请失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 20rpx 24rpx 140rpx;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
}
.st {
  font-size: 34rpx;
  font-weight: 700;
  display: block;
  margin-bottom: 12rpx;
}
.line {
  display: block;
  margin-top: 8rpx;
  font-size: 26rpx;
}
.muted {
  color: #c0392b;
}
.foot {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  gap: 16rpx;
  padding: 16rpx 24rpx;
  background: #fff;
}
.foot .qx-btn {
  flex: 1;
}
</style>
