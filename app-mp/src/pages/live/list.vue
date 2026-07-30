<template>
  <view class="page">
    <view class="tip">直播带货演示（无微信推流）</view>
    <view v-if="!list.length" class="empty">暂无直播间</view>
    <view
      v-for="r in list"
      :key="r.broadcast_room_id"
      class="card"
      @click="goDetail(r.broadcast_room_id)"
    >
      <view class="row">
        <text class="name">{{ r.name }}</text>
        <text class="badge" :class="{ on: r.live_status === 101 }">
          {{ liveLabel(r.live_status) }}
        </text>
      </view>
      <text class="mer">{{ r.mer_name || "店铺" }} · {{ r.anchor_name || "主播" }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchLiveRooms, type LiveRoom } from "@/api/live";

const list = ref<LiveRoom[]>([]);

onShow(async () => {
  try {
    const data = await fetchLiveRooms();
    list.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function liveLabel(s: number) {
  if (s === 101) return "直播中";
  if (s === 103) return "已结束";
  return "未开始";
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/live/detail?id=${id}` });
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
</style>
