<template>
  <view class="page">
    <view class="tip">种草分享 · 审核通过后展示 · 可挂商品</view>
    <scroll-view scroll-x class="topics">
      <view
        class="chip"
        :class="{ on: topicId === 0 }"
        @click="selectTopic(0)"
      >
        全部
      </view>
      <view
        v-for="t in topics"
        :key="t.topic_id"
        class="chip"
        :class="{ on: topicId === t.topic_id }"
        @click="selectTopic(t.topic_id)"
      >
        {{ t.topic_name }}
      </view>
    </scroll-view>
    <view class="actions">
      <view class="qx-btn qx-btn-primary" @click="goCreate">发帖</view>
    </view>
    <view v-if="!list.length" class="empty">暂无内容</view>
    <view
      v-for="p in list"
      :key="p.community_id"
      class="card"
      @click="goDetail(p.community_id)"
    >
      <text class="name">{{ p.title }}</text>
      <text class="meta">{{ p.nickname || "用户" }} · {{ p.topic_name || "话题" }}</text>
      <text class="body">{{ p.content }}</text>
      <text v-if="p.product_name" class="goods">
        挂货 {{ p.product_name }} · ¥{{ Number(p.product_price || 0).toFixed(2) }}
      </text>
      <text class="stat">评论 {{ p.count_reply || 0 }} · 浏览 {{ p.pv || 0 }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import {
  fetchCommunityPosts,
  fetchCommunityTopics,
  type CommunityPost,
  type CommunityTopic,
} from "@/api/community";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const topics = ref<CommunityTopic[]>([]);
const list = ref<CommunityPost[]>([]);
const topicId = ref(0);

async function load() {
  try {
    const [t, data] = await Promise.all([
      fetchCommunityTopics().catch(() => ({ list: [] as CommunityTopic[] })),
      fetchCommunityPosts(1, 20, topicId.value),
    ]);
    topics.value = t.list || [];
    list.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

onShow(() => {
  void load();
});

function selectTopic(id: number) {
  topicId.value = id;
  void load();
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/community/detail?id=${id}` });
}

function goCreate() {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  uni.navigateTo({ url: "/pages/community/create" });
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
  margin-bottom: 16rpx;
}
.topics {
  white-space: nowrap;
  margin-bottom: 16rpx;
}
.chip {
  display: inline-block;
  margin-right: 16rpx;
  padding: 10rpx 28rpx;
  border-radius: 999rpx;
  background: #fff;
  font-size: 24rpx;
  color: #666;
  &.on {
    background: #e23030;
    color: #fff;
  }
}
.actions {
  margin-bottom: 20rpx;
}
.empty {
  text-align: center;
  color: #999;
  padding: 80rpx 0;
}
.card {
  background: #fff;
  border-radius: 12rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
}
.name {
  font-size: 30rpx;
  font-weight: 600;
}
.meta {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: #999;
}
.body {
  margin-top: 12rpx;
  color: #444;
  font-size: 26rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
.goods {
  display: block;
  margin-top: 12rpx;
  color: #e23030;
  font-size: 24rpx;
}
.stat {
  display: block;
  margin-top: 10rpx;
  font-size: 22rpx;
  color: #aaa;
}
</style>
