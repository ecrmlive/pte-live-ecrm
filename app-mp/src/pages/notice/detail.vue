<template>
  <view class="page">
    <text class="title">{{ row?.title || "公告" }}</text>
    <text class="body">{{ row?.content || "" }}</text>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { fetchNotice, type Notice } from "@/api/content";

const row = ref<Notice | null>(null);

onLoad(async (q) => {
  const id = Number(q?.id || 0);
  if (!id) return;
  try {
    row.value = await fetchNotice(id);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: #fff;
  padding: 32rpx 28rpx;
}
.title {
  font-size: 36rpx;
  font-weight: 700;
  display: block;
  margin-bottom: 24rpx;
}
.body {
  display: block;
  line-height: 1.7;
  color: #333;
  white-space: pre-wrap;
  font-size: 28rpx;
}
</style>
