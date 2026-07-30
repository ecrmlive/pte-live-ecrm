<template>
  <view class="page">
    <view class="tip">发布后需平台审核通过才会展示</view>
    <input v-model="title" class="input" placeholder="标题" />
    <textarea v-model="content" class="textarea" placeholder="分享你的种草心得…" />
    <view class="label">话题</view>
    <scroll-view scroll-x class="topics">
      <view
        v-for="t in topics"
        :key="t.topic_id"
        class="chip"
        :class="{ on: topicId === t.topic_id }"
        @click="topicId = t.topic_id"
      >
        {{ t.topic_name }}
      </view>
    </scroll-view>
    <input v-model="productId" class="input" type="number" placeholder="挂货商品 ID（可选）" />
    <view class="qx-btn qx-btn-primary" @click="submit">提交审核</view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import {
  createCommunityPost,
  fetchCommunityTopics,
  type CommunityTopic,
} from "@/api/community";

const title = ref("");
const content = ref("");
const productId = ref("");
const topicId = ref(1);
const topics = ref<CommunityTopic[]>([]);

onLoad(async () => {
  try {
    const data = await fetchCommunityTopics();
    topics.value = data.list || [];
    if (topics.value.length && !topics.value.some((t) => t.topic_id === topicId.value)) {
      topicId.value = topics.value[0].topic_id;
    }
  } catch {
    topics.value = [{ topic_id: 1, topic_name: "好物推荐" }];
  }
});

async function submit() {
  const t = title.value.trim();
  const c = content.value.trim();
  if (!t || !c) {
    uni.showToast({ title: "请填写标题和内容", icon: "none" });
    return;
  }
  try {
    await createCommunityPost({
      title: t,
      content: c,
      product_id: Number(productId.value) || 0,
      topic_id: topicId.value || 1,
      category_id: 1,
    });
    uni.showToast({ title: "已提交审核", icon: "success" });
    setTimeout(() => {
      uni.redirectTo({ url: "/pages/community/list" });
    }, 400);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "提交失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 24rpx;
  background: var(--qx-bg, #f5f5f5);
}
.tip {
  color: #888;
  font-size: 24rpx;
  margin-bottom: 20rpx;
}
.label {
  font-size: 24rpx;
  color: #666;
  margin-bottom: 12rpx;
}
.topics {
  white-space: nowrap;
  margin-bottom: 20rpx;
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
.input,
.textarea {
  width: 100%;
  background: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  box-sizing: border-box;
}
.textarea {
  min-height: 280rpx;
}
</style>
