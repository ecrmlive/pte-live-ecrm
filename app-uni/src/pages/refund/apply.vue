<template>
  <view class="page">
    <view class="card">
      <text class="title">申请仅退款</text>
      <text class="line">子单 ID {{ orderId }}</text>
      <textarea
        v-model="message"
        class="area"
        placeholder="请填写退款原因"
        maxlength="100"
      />
    </view>
    <view class="foot qx-safe-bottom">
      <view class="qx-btn qx-btn-primary" @click="submit">提交申请</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { applyRefund } from "@/api/refund";

const orderId = ref(0);
const message = ref("不想要了");

onLoad((q) => {
  orderId.value = Number(q?.order_id || 0);
});

async function submit() {
  if (!orderId.value || !message.value.trim()) {
    uni.showToast({ title: "请填写原因", icon: "none" });
    return;
  }
  try {
    const ro = await applyRefund({
      order_id: orderId.value,
      refund_type: 1,
      refund_message: message.value.trim(),
    });
    uni.showToast({ title: "已提交", icon: "success" });
    setTimeout(() => {
      uni.redirectTo({ url: `/pages/refund/detail?id=${ro.refund_order_id}` });
    }, 400);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "申请失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding-bottom: 140rpx;
}
.card {
  margin: 20rpx 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
}
.title {
  font-weight: 700;
  display: block;
  margin-bottom: 12rpx;
}
.line {
  display: block;
  color: #999;
  font-size: 24rpx;
  margin-bottom: 20rpx;
}
.area {
  width: 100%;
  min-height: 200rpx;
  background: #f7f7f7;
  border-radius: 12rpx;
  padding: 20rpx;
  box-sizing: border-box;
}
.foot {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 20rpx 32rpx;
  background: #fff;
}
</style>
