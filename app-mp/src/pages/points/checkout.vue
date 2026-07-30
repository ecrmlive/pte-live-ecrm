<template>
  <view class="page">
    <view class="card">
      <text class="title">收货地址</text>
      <view v-if="addr" class="addr">
        <text>{{ addr.real_name }} {{ addr.phone }}</text>
        <text class="line">{{ addr.province }}{{ addr.city }}{{ addr.district }} {{ addr.detail }}</text>
      </view>
      <text v-else class="line">暂无地址</text>
      <text class="link" @click="goAddress">管理地址</text>
    </view>
    <view class="card">
      <text class="title">{{ storeName || "积分商品" }}</text>
      <text class="line">需积分 {{ totalIntegral }} · 现金 ¥0.00</text>
      <text class="line">我的积分 {{ userIntegral }}</text>
    </view>
    <view class="foot qx-safe-bottom">
      <view class="qx-btn qx-btn-primary" @click="submit">确认兑换</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad, onShow } from "@dcloudio/uni-app";
import { fetchAddresses, type Address } from "@/api/cart";
import { pointsPay, v3Check, v3Create } from "@/api/order";

const productId = ref(0);
const addr = ref<Address | null>(null);
const totalIntegral = ref(0);
const userIntegral = ref(0);
const storeName = ref("");

async function refreshCheck() {
  if (!productId.value) return;
  try {
    const check = await v3Check({ product_id: productId.value, cart_num: 1 });
    totalIntegral.value = check.integral || 0;
    userIntegral.value = check.user_integral || 0;
    storeName.value = check.store_name || "";
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "核对失败", icon: "none" });
  }
}

onLoad(async (q) => {
  productId.value = Number(q?.product_id || 0);
  await refreshCheck();
});

onShow(async () => {
  try {
    const res = await fetchAddresses();
    addr.value = (res.list || []).find((a) => a.is_default === 1) || res.list?.[0] || null;
  } catch {
    /* ignore */
  }
});

function goAddress() {
  uni.navigateTo({ url: "/pages/address/list" });
}

async function submit() {
  if (!addr.value || !productId.value) {
    uni.showToast({ title: "缺少地址或商品", icon: "none" });
    return;
  }
  try {
    const g = await v3Create({
      product_id: productId.value,
      address_id: addr.value.address_id,
      cart_num: 1,
    });
    await pointsPay(g.group_order_id);
    uni.redirectTo({ url: `/pages/order/result?id=${g.group_order_id}` });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "兑换失败", icon: "none" });
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
  color: var(--qx-text-secondary);
  margin-top: 8rpx;
  font-size: 24rpx;
}
.link {
  display: block;
  margin-top: 16rpx;
  color: var(--qx-primary, #e23030);
  font-size: 24rpx;
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
