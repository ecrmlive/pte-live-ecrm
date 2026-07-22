<template>
  <view class="page">
    <view v-for="c in list" :key="c.coupon_id" class="card">
      <view class="meta">
        <text class="title">{{ c.title }}</text>
        <text class="sub">¥{{ Number(c.coupon_price).toFixed(2) }} · 满{{ c.use_min_price }}可用</text>
      </view>
      <view class="qx-btn qx-btn-primary mini" :class="{ disabled: c.received }" @click="onReceive(c)">
        {{ c.received ? "已领" : "领取" }}
      </view>
    </view>
    <text v-if="!list.length" class="empty">暂无可领优惠券</text>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchCouponCenter, receiveCoupon, type Coupon } from "@/api/coupon";

const list = ref<Coupon[]>([]);

async function load() {
  try {
    const res = await fetchCouponCenter();
    list.value = res.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

onShow(() => {
  void load();
});

async function onReceive(c: Coupon) {
  if (c.received) return;
  try {
    await receiveCoupon(c.coupon_id);
    uni.showToast({ title: "领取成功", icon: "none" });
    await load();
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "领取失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 20rpx 24rpx;
}
.card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 16rpx;
}
.title {
  font-weight: 700;
  display: block;
}
.sub {
  display: block;
  margin-top: 8rpx;
  color: var(--qx-text-secondary);
  font-size: 24rpx;
}
.mini {
  min-width: 140rpx;
  height: 64rpx;
  line-height: 64rpx;
  font-size: 26rpx;
  padding: 0 20rpx;
}
.disabled {
  opacity: 0.5;
}
.empty {
  display: block;
  text-align: center;
  color: var(--qx-text-secondary);
  margin-top: 80rpx;
}
</style>
