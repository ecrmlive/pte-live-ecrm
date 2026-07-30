<template>
  <view class="page">
    <view class="tip">场次内下单自动按秒杀价结算</view>
    <view v-if="!list.length" class="empty">暂无秒杀活动</view>
    <view
      v-for="a in list"
      :key="a.seckill_active_id"
      class="card"
      @click="goDetail(a.product_id)"
    >
      <view class="row">
        <text class="name">{{ a.store_name || a.name }}</text>
        <text class="badge" :class="{ on: a.in_window }">
          {{ a.in_window ? "抢购中" : "未开场" }}
        </text>
      </view>
      <text class="mer">{{ a.mer_name }} · {{ a.start_day }} ~ {{ a.end_day }}</text>
      <view class="price-row">
        <text class="price">¥{{ a.seckill_price }}</text>
        <text v-if="a.price" class="ot">¥{{ a.price }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchSeckillList, type SeckillActive } from "@/api/seckill";

const list = ref<SeckillActive[]>([]);

onShow(async () => {
  try {
    const data = await fetchSeckillList();
    list.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 24rpx;
  background: var(--qx-bg, #f5f5f5);
}
.tip {
  font-size: 24rpx;
  color: #888;
  margin-bottom: 20rpx;
}
.empty {
  text-align: center;
  color: #999;
  padding: 80rpx 0;
}
.card {
  background: #fff;
  padding: 28rpx;
  margin-bottom: 20rpx;
  border-radius: 12rpx;
}
.row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.name {
  font-size: 30rpx;
  font-weight: 600;
}
.badge {
  font-size: 22rpx;
  color: #999;
  &.on {
    color: #e23030;
  }
}
.mer {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: #999;
}
.price-row {
  margin-top: 16rpx;
  display: flex;
  align-items: baseline;
  gap: 12rpx;
}
.price {
  color: #e23030;
  font-size: 36rpx;
  font-weight: 700;
}
.ot {
  color: #bbb;
  text-decoration: line-through;
  font-size: 24rpx;
}
</style>
