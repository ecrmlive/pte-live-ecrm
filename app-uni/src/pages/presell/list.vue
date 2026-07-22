<template>
  <view class="page">
    <view class="tip">全款一次付清 · 定金先付定金再付尾款</view>
    <view class="tip link" @click="goFinals">查看待付尾款 ›</view>
    <view v-if="!list.length" class="empty">暂无预售活动</view>
    <view v-for="g in list" :key="g.product_presell_id" class="card" @click="goDetail(g)">
      <text class="name">{{ g.store_name || "预售商品" }}</text>
      <text class="mer">
        {{ g.presell_type === 2 ? "定金预售" : "全款预售" }} · {{ g.mer_name }} · 库存 {{ g.stock }}
      </text>
      <view class="price-row">
        <text class="price">¥{{ g.price }}</text>
        <text v-if="g.presell_type === 2" class="ot">定金¥{{ g.down_price }}</text>
        <text v-else-if="g.ot_price" class="ot">¥{{ g.ot_price }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchPresells, type PresellActive } from "@/api/presell";

const list = ref<PresellActive[]>([]);

onShow(async () => {
  try {
    const data = await fetchPresells();
    list.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goDetail(g: PresellActive) {
  uni.navigateTo({ url: `/pages/presell/detail?id=${g.product_presell_id}` });
}
function goFinals() {
  uni.navigateTo({ url: "/pages/presell/finals" });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 24rpx;
  background: #f5f5f5;
}
.tip {
  margin-bottom: 12rpx;
  color: #666;
  font-size: 24rpx;
}
.tip.link {
  color: #c45c26;
  margin-bottom: 20rpx;
}
.empty {
  margin-top: 80rpx;
  text-align: center;
  color: #999;
}
.card {
  background: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}
.name {
  font-size: 30rpx;
  font-weight: 600;
}
.mer {
  display: block;
  margin-top: 8rpx;
  color: #888;
  font-size: 24rpx;
}
.price-row {
  margin-top: 16rpx;
  display: flex;
  align-items: baseline;
  gap: 12rpx;
}
.price {
  color: #e23030;
  font-size: 34rpx;
  font-weight: 700;
}
.ot {
  color: #999;
  font-size: 24rpx;
  text-decoration: line-through;
}
</style>
