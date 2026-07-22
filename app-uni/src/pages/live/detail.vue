<template>
  <view class="page">
    <view v-if="room" class="head">
      <text class="name">{{ room.name }}</text>
      <text class="badge" :class="{ on: room.live_status === 101 }">
        {{ liveLabel(room.live_status) }}
      </text>
      <text class="meta">{{ room.mer_name }} · 主播 {{ room.anchor_name || "-" }}</text>
    </view>
    <view class="section">
      <text class="section-title">直播间商品</text>
      <view v-if="!(room?.goods || []).length" class="empty">暂无挂货</view>
      <view
        v-for="g in room?.goods || []"
        :key="g.product_id"
        class="goods"
        @click="goGoods(g.product_id)"
      >
        <view class="g-meta">
          <text class="g-name">{{ g.store_name || `商品#${g.product_id}` }}</text>
          <text class="g-price">¥{{ g.price ?? "-" }}</text>
        </view>
        <text class="g-go">去看看</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { fetchLiveRoom, type LiveRoom } from "@/api/live";

const room = ref<LiveRoom | null>(null);

onLoad(async (q) => {
  const id = Number(q?.id || 0);
  if (!id) {
    uni.showToast({ title: "参数错误", icon: "none" });
    return;
  }
  try {
    room.value = await fetchLiveRoom(id);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function liveLabel(s: number) {
  if (s === 101) return "直播中";
  if (s === 103) return "已结束";
  return "未开始";
}

function goGoods(id: number) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 24rpx;
  background: var(--qx-bg, #f5f5f5);
}
.head {
  background: #fff;
  border-radius: 12rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
}
.name {
  font-size: 34rpx;
  font-weight: 700;
}
.badge {
  margin-left: 12rpx;
  font-size: 22rpx;
  color: #999;
  &.on {
    color: #e23030;
  }
}
.meta {
  display: block;
  margin-top: 12rpx;
  font-size: 24rpx;
  color: #888;
}
.section {
  background: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
}
.section-title {
  font-size: 28rpx;
  font-weight: 600;
}
.empty {
  text-align: center;
  color: #999;
  padding: 40rpx 0;
}
.goods {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 0;
  border-top: 1px solid #f0f0f0;
  margin-top: 8rpx;
}
.g-name {
  display: block;
  font-size: 28rpx;
}
.g-price {
  display: block;
  margin-top: 6rpx;
  color: #e23030;
  font-weight: 600;
}
.g-go {
  color: #1677ff;
  font-size: 24rpx;
}
</style>
