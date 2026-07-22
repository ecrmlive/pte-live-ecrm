<template>
  <view class="page">
    <view
      v-for="g in list"
      :key="g.group_order_id"
      class="card"
      @click="goDetail(g.group_order_id)"
    >
      <text class="sn">{{ g.group_order_sn }}</text>
      <text class="line">¥{{ g.pay_price }} · {{ g.paid === 1 ? "已支付" : "待付款" }} · {{ g.total_num || 0 }} 件</text>
    </view>
    <view v-if="!list.length" class="empty">暂无订单</view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchOrders, type GroupOrder } from "@/api/order";

const list = ref<GroupOrder[]>([]);

onShow(async () => {
  try {
    const res = await fetchOrders();
    list.value = res.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/order/detail?id=${id}` });
}
</script>

<style lang="scss" scoped>
.page { min-height: 100vh; background: var(--qx-bg); padding: 20rpx 0; }
.card { margin: 0 24rpx 20rpx; background: #fff; border-radius: 16rpx; padding: 28rpx; }
.sn { font-weight: 700; display: block; }
.line { margin-top: 8rpx; color: var(--qx-text-secondary); font-size: 24rpx; display: block; }
.empty { text-align: center; padding: 80rpx; color: #999; }
</style>
