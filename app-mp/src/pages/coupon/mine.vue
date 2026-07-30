<template>
  <view class="page">
    <view class="tabs">
      <text :class="{ active: status === 0 }" @click="setStatus(0)">未使用</text>
      <text :class="{ active: status === 1 }" @click="setStatus(1)">已使用</text>
      <text :class="{ active: status === 2 }" @click="setStatus(2)">已过期</text>
    </view>
    <view v-for="c in list" :key="c.coupon_user_id" class="card">
      <text class="title">{{ c.coupon_title }}</text>
      <text class="sub">¥{{ Number(c.coupon_price).toFixed(2) }} · 满{{ c.use_min_price }} · #{{ c.coupon_user_id }}</text>
    </view>
    <text v-if="!list.length" class="empty">暂无优惠券</text>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchMyCoupons, type CouponUser } from "@/api/coupon";

const status = ref(0);
const list = ref<CouponUser[]>([]);

async function load() {
  try {
    const res = await fetchMyCoupons(status.value);
    list.value = res.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

function setStatus(s: number) {
  status.value = s;
  void load();
}

onShow(() => {
  void load();
});
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 20rpx 24rpx;
}
.tabs {
  display: flex;
  gap: 24rpx;
  margin-bottom: 20rpx;
  text {
    color: var(--qx-text-secondary);
    font-size: 28rpx;
  }
  .active {
    color: var(--qx-primary, #e23030);
    font-weight: 700;
  }
}
.card {
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
.empty {
  display: block;
  text-align: center;
  color: var(--qx-text-secondary);
  margin-top: 80rpx;
}
</style>
