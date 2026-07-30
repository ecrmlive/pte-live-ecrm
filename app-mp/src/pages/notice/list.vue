<template>
  <view class="page">
    <view v-for="n in list" :key="n.notice_id" class="card" @click="goDetail(n.notice_id)">
      <text class="title">{{ n.title }}</text>
      <text class="sub">{{ n.create_time || "" }}</text>
    </view>
    <text v-if="!list.length && !loading" class="empty">暂无公告</text>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchNotices, type Notice } from "@/api/content";

const list = ref<Notice[]>([]);
const loading = ref(false);

async function load() {
  loading.value = true;
  try {
    const res = await fetchNotices(1, 50);
    list.value = res.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  } finally {
    loading.value = false;
  }
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/notice/detail?id=${id}` });
}

onShow(() => {
  void load();
});
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 20rpx 24rpx;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 16rpx;
}
.title {
  font-weight: 700;
  display: block;
}
.sub {
  display: block;
  margin-top: 8rpx;
  color: var(--qx-text-secondary);
  font-size: 24rpx;
}
.empty {
  display: block;
  text-align: center;
  color: var(--qx-text-secondary);
  margin-top: 80rpx;
}
</style>
