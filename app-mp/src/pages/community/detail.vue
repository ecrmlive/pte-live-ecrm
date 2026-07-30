<template>
  <view class="page" v-if="post">
    <view class="card">
      <text class="title">{{ post.title }}</text>
      <text class="meta">{{ post.nickname }} · {{ post.topic_name }} · 浏览 {{ post.pv || 0 }}</text>
      <text class="body">{{ post.content }}</text>
      <view v-if="post.product_id" class="goods" @click="goGoods(post.product_id)">
        {{ post.product_name || `商品 #${post.product_id}` }} · ¥{{ Number(post.product_price || 0).toFixed(2) }} ›
      </view>
    </view>
    <view class="section">
      <text class="sec-title">评论 {{ replies.length }}</text>
      <view v-if="!replies.length" class="empty">暂无评论</view>
      <view v-for="r in replies" :key="r.reply_id" class="reply">
        <text class="r-nick">{{ r.nickname || "用户" }}</text>
        <text class="r-body">{{ r.content }}</text>
      </view>
    </view>
    <view class="composer">
      <input v-model="draft" class="input" placeholder="说点什么…" />
      <view class="qx-btn qx-btn-primary send" @click="send">发送</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import {
  createReply,
  fetchCommunityPost,
  fetchReplies,
  type CommunityPost,
  type CommunityReply,
} from "@/api/community";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const id = ref(0);
const post = ref<CommunityPost | null>(null);
const replies = ref<CommunityReply[]>([]);
const draft = ref("");

onLoad(async (q) => {
  id.value = Number(q?.id || 0);
  if (!id.value) return;
  try {
    post.value = await fetchCommunityPost(id.value);
    const data = await fetchReplies(id.value);
    replies.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function goGoods(pid: number) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${pid}` });
}

async function send() {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  const content = draft.value.trim();
  if (!content) return;
  try {
    const row = await createReply(id.value, content);
    replies.value = [...replies.value, row];
    draft.value = "";
    if (post.value) {
      post.value.count_reply = (post.value.count_reply || 0) + 1;
    }
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "发送失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 24rpx 24rpx 160rpx;
  background: #f5f5f5;
}
.card,
.reply {
  background: #fff;
  border-radius: 12rpx;
  padding: 28rpx;
}
.title {
  font-size: 34rpx;
  font-weight: 700;
}
.meta {
  display: block;
  margin-top: 10rpx;
  color: #999;
  font-size: 22rpx;
}
.body {
  display: block;
  margin-top: 16rpx;
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
}
.goods {
  margin-top: 20rpx;
  padding: 20rpx;
  background: #fafafa;
  border-radius: 8rpx;
  color: #e23030;
  font-size: 26rpx;
}
.section {
  margin-top: 24rpx;
}
.sec-title {
  display: block;
  margin-bottom: 12rpx;
  font-weight: 600;
}
.empty {
  color: #999;
  padding: 40rpx 0;
  text-align: center;
}
.reply {
  margin-bottom: 12rpx;
}
.r-nick {
  font-size: 24rpx;
  color: #666;
  font-weight: 600;
}
.r-body {
  display: block;
  margin-top: 8rpx;
  font-size: 26rpx;
}
.composer {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  gap: 16rpx;
  padding: 16rpx 24rpx calc(16rpx + env(safe-area-inset-bottom));
  background: #fff;
  border-top: 1px solid #eee;
}
.input {
  flex: 1;
  height: 72rpx;
  background: #f5f5f5;
  border-radius: 8rpx;
  padding: 0 20rpx;
}
.send {
  width: 160rpx;
}
</style>
