<template>
  <view class="page">
    <view class="banner">
      <text class="bal">我的积分 {{ integral }}</text>
      <text class="rule">独立入口 /order/v3 · 与普通单不混用</text>
    </view>
    <view v-for="p in list" :key="p.id" class="card">
      <view class="meta">
        <text class="title">{{ p.store_name }}</text>
        <text class="sub">{{ p.mer_name || "" }} · {{ pts(p) }} 积分</text>
      </view>
      <view class="qx-btn qx-btn-primary mini" @click="buy(p)">兑换</view>
    </view>
    <text v-if="!list.length" class="empty">暂无积分商品</text>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchPointsProducts, type ProductItem } from "@/api/catalog";
import { fetchIntegral } from "@/api/order";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const list = ref<ProductItem[]>([]);
const integral = ref(0);

function pts(p: ProductItem) {
  return p.integral || Number(p.ot_price || 0);
}

async function load() {
  try {
    if (user.isLogin) {
      const bal = await fetchIntegral();
      integral.value = bal.integral || 0;
    }
    const res = await fetchPointsProducts();
    list.value = (res.list || []).map((x: ProductItem & { product_id?: number }) => ({
      ...x,
      id: x.id || x.product_id || 0,
    }));
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

onShow(() => {
  void load();
});

function buy(p: ProductItem) {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  const id = Number(p.id);
  uni.navigateTo({ url: `/pages/points/checkout?product_id=${id}` });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 20rpx 24rpx;
}
.banner {
  background: linear-gradient(135deg, #3d2b1f, #8b5a2b);
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  color: #fff;
}
.bal {
  display: block;
  font-size: 36rpx;
  font-weight: 700;
}
.rule {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  opacity: 0.85;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.title {
  font-weight: 700;
  display: block;
}
.sub {
  color: #999;
  font-size: 24rpx;
  margin-top: 8rpx;
  display: block;
}
.mini {
  width: 160rpx;
  height: 64rpx;
  line-height: 64rpx;
  font-size: 26rpx;
}
.empty {
  display: block;
  text-align: center;
  color: #999;
  padding: 80rpx;
}
</style>
