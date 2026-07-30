<template>
  <view class="page">
    <view class="tip">选择服务 → 选日期时段 → 下单后到店核销</view>
    <view v-if="!list.length" class="empty">暂无可预约服务</view>
    <view
      v-for="p in list"
      :key="p.product_id"
      class="card"
      @click="goBook(p.product_id)"
    >
      <text class="name">{{ p.store_name }}</text>
      <text class="mer">{{ p.mer_name || "" }}</text>
      <text class="price">¥{{ p.price }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchReservationProducts, type ReservationProduct } from "@/api/reservation";

const list = ref<ReservationProduct[]>([]);

onShow(async () => {
  try {
    const data = await fetchReservationProducts();
    list.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goBook(id: number) {
  uni.navigateTo({ url: `/pages/reservation/detail?id=${id}` });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 24rpx;
  background: #f5f5f5;
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
.name {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
}
.mer {
  display: block;
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #999;
}
.price {
  display: block;
  margin-top: 12rpx;
  color: #e23030;
  font-size: 32rpx;
}
</style>
