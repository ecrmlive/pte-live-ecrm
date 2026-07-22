<template>
  <view class="page">
    <view
      v-for="r in list"
      :key="r.refund_order_id"
      class="card"
      @click="goDetail(r.refund_order_id)"
    >
      <text class="sn">{{ r.refund_order_sn }}</text>
      <text class="line">¥{{ Number(r.refund_price).toFixed(2) }} · {{ refundStatusText(r.status) }}</text>
    </view>
    <view v-if="!list.length" class="empty">暂无售后单</view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchRefunds, refundStatusText, type RefundOrder } from "@/api/refund";

const list = ref<RefundOrder[]>([]);

onShow(async () => {
  try {
    const res = await fetchRefunds();
    list.value = res.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/refund/detail?id=${id}` });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 20rpx 0;
}
.card {
  margin: 0 24rpx 20rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
}
.sn {
  font-weight: 700;
  display: block;
}
.line {
  margin-top: 8rpx;
  color: #999;
  font-size: 24rpx;
  display: block;
}
.empty {
  text-align: center;
  padding: 80rpx;
  color: #999;
}
</style>
