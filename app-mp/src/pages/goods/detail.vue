<template>
  <view class="page" v-if="detail">
    <view class="hero">
      <text class="hero-letter">{{ detail.store_name.slice(0, 1) }}</text>
    </view>
    <view class="card">
      <view class="price-row">
        <text class="price">¥{{ detail.price }}</text>
        <text v-if="detail.ot_price" class="ot">¥{{ detail.ot_price }}</text>
        <text v-if="svipPriceText" class="svip">会员¥{{ svipPriceText }}</text>
      </view>
      <text class="name">{{ detail.store_name }}</text>
      <text class="info">{{ detail.store_info || "暂无简介" }}</text>
      <view class="meta">
        <text @click="goStore">店铺：{{ detail.mer_name }}</text>
        <text>库存 {{ detail.stock }} · 已售 {{ detail.sales }}</text>
      </view>
    </view>
    <view class="card">
      <text class="block-title">配送 / 规格</text>
      <text class="line">配送：{{ deliveryText }}</text>
      <text class="line">规格：{{ detail.spec_type === 1 ? "多规格" : "单规格" }} · {{ detail.unit_name }}</text>
    </view>
    <view class="foot qx-safe-bottom">
      <view class="qx-btn qx-btn-ghost half" @click="goStore">进店</view>
      <view class="qx-btn qx-btn-primary half" @click="onAddCart">加入购物车</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { addCart } from "@/api/cart";
import { fetchProductDetail, type ProductDetail } from "@/api/catalog";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const detail = ref<ProductDetail | null>(null);

const deliveryText = computed(() => {
  const w = detail.value?.delivery_way || "";
  const map: Record<string, string> = { "1": "自提", "2": "快递", "3": "包邮" };
  return (
    w
      .split(",")
      .map((x) => map[x] || x)
      .filter(Boolean)
      .join(" / ") || "快递"
  );
});

const svipPriceText = computed(() => {
  const d = detail.value;
  if (!d || !d.svip_price_type) return "";
  if (d.svip_price_type === 2 && d.svip_price) return Number(d.svip_price).toFixed(2);
  if (d.svip_price_type === 1) return (Number(d.price) * 0.9).toFixed(2);
  return "";
});

onLoad(async (q) => {
  const id = Number(q?.id || 0);
  if (!id) return;
  try {
    detail.value = await fetchProductDetail(id);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goStore() {
  if (!detail.value) return;
  uni.navigateTo({ url: `/pages/store/home?mer_id=${detail.value.mer_id}` });
}

async function onAddCart() {
  if (!detail.value) return;
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  try {
    await addCart({ product_id: Number(detail.value.id), cart_num: 1 });
    uni.showToast({ title: "已加入购物车", icon: "success" });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加购失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding-bottom: 140rpx;
}
.hero {
  height: 520rpx;
  background: linear-gradient(135deg, #f3e9e6, #ffe8e0);
  display: flex;
  align-items: center;
  justify-content: center;
}
.hero-letter {
  font-size: 120rpx;
  font-weight: 700;
  color: #c9a39a;
}
.card {
  margin: 20rpx 24rpx 0;
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
}
.price-row {
  display: flex;
  align-items: baseline;
  gap: 12rpx;
}
.svip {
  color: #b8860b;
  font-size: 26rpx;
  font-weight: 600;
}
.price {
  color: var(--qx-price);
  font-size: 44rpx;
  font-weight: 700;
}
.ot {
  color: #bbb;
  text-decoration: line-through;
  font-size: 24rpx;
}
.name {
  display: block;
  margin-top: 16rpx;
  font-size: 34rpx;
  font-weight: 700;
}
.info {
  display: block;
  margin-top: 12rpx;
  color: var(--qx-text-secondary);
  font-size: 26rpx;
  line-height: 1.5;
}
.meta {
  margin-top: 20rpx;
  display: flex;
  justify-content: space-between;
  color: #888;
  font-size: 24rpx;
}
.block-title {
  display: block;
  font-weight: 700;
  margin-bottom: 12rpx;
}
.line {
  display: block;
  color: #666;
  font-size: 26rpx;
  margin-top: 8rpx;
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
.half {
  flex: 1;
}
</style>
