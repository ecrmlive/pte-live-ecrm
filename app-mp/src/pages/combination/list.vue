<template>
  <view class="page">
    <view class="tip">2 人成团 · 支付后计入人数 · 满员自动待发货</view>
    <view v-if="!list.length" class="empty">暂无拼团活动</view>
    <view v-for="g in list" :key="g.product_group_id" class="card" @click="goDetail(g)">
      <text class="name">{{ g.store_name || "拼团商品" }}</text>
      <text class="mer">{{ g.mer_name }} · {{ g.buying_count_num }}人团</text>
      <view class="price-row">
        <text class="price">¥{{ g.price }}</text>
        <text v-if="g.ot_price" class="ot">¥{{ g.ot_price }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchGroups, type ProductGroup } from "@/api/combination";

const list = ref<ProductGroup[]>([]);

onShow(async () => {
  try {
    const data = await fetchGroups();
    list.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goDetail(g: ProductGroup) {
  uni.navigateTo({ url: `/pages/combination/detail?id=${g.product_group_id}` });
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
  font-size: 30rpx;
  font-weight: 600;
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
  gap: 12rpx;
  align-items: baseline;
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
