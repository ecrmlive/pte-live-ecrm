<template>
  <view class="page">
    <view v-if="!buckets.length" class="empty">
      <text class="title">购物车还是空的</text>
      <view class="qx-btn qx-btn-primary btn" @click="goHome">去逛逛</view>
    </view>
    <view v-else>
      <view v-for="b in buckets" :key="b.mer_id" class="bucket">
        <text class="mer">{{ b.mer_name || "商户" }}</text>
        <view v-for="it in b.items" :key="it.cart_id" class="row">
          <view class="info">
            <text class="name">{{ it.store_name }}</text>
            <text class="price">¥{{ it.price }} × {{ it.cart_num }}</text>
          </view>
          <text class="rm" @click="onRemove(it.cart_id)">删除</text>
        </view>
      </view>
      <view class="foot qx-safe-bottom">
        <text>合计 ¥{{ total }}</text>
        <view class="qx-btn qx-btn-primary half" @click="goCheckout">去结算</view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchCart, removeCart, type CartBucket } from "@/api/cart";
import { getToken } from "@/utils/storage";

const buckets = ref<CartBucket[]>([]);
const total = computed(() =>
  buckets.value.reduce((s, b) => s + (b.subtotal || 0), 0).toFixed(2),
);

async function load() {
  if (!getToken()) {
    buckets.value = [];
    return;
  }
  try {
    const res = await fetchCart();
    buckets.value = res.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

onShow(() => void load());

function goHome() {
  uni.switchTab({ url: "/pages/index/index" });
}

async function onRemove(id: number) {
  await removeCart(id);
  await load();
}

function goCheckout() {
  const ids = buckets.value.flatMap((b) => b.items.filter((i) => !i.is_fail).map((i) => i.cart_id));
  if (!ids.length) return;
  uni.navigateTo({ url: `/pages/order/checkout?cart_ids=${ids.join(",")}` });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding-bottom: 140rpx;
}
.empty {
  padding: 120rpx 40rpx;
  text-align: center;
}
.title {
  font-size: 34rpx;
  font-weight: 700;
  display: block;
}
.btn {
  margin: 40rpx auto 0;
  width: 320rpx;
}
.bucket {
  margin: 20rpx 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
}
.mer {
  font-weight: 700;
  display: block;
  margin-bottom: 16rpx;
}
.row {
  display: flex;
  justify-content: space-between;
  padding: 16rpx 0;
  border-top: 1px solid #f0f0f0;
}
.name {
  display: block;
  font-size: 28rpx;
}
.price {
  color: var(--qx-price);
  font-size: 24rpx;
}
.rm {
  color: #999;
  font-size: 24rpx;
}
.foot {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  background: #fff;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 32rpx;
  border-top: 1px solid #eee;
}
.half {
  width: 280rpx;
}
</style>
