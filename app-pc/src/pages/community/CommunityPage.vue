<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import {
  fetchCommunityPosts,
  fetchCommunityTopics,
  type CommunityPost,
  type CommunityTopic,
} from "@/api/community";

const router = useRouter();
const loading = ref(false);
const list = ref<CommunityPost[]>([]);
const topics = ref<CommunityTopic[]>([]);
const topicId = ref(0);

async function load() {
  loading.value = true;
  try {
    const [t, data] = await Promise.all([
      fetchCommunityTopics().catch(() => ({ list: [] as CommunityTopic[] })),
      fetchCommunityPosts(1, 20, topicId.value),
    ]);
    topics.value = t.list || [];
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  void load();
});

function selectTopic(id: number) {
  topicId.value = id;
  void load();
}

function open(id: number) {
  router.push(`/community/${id}`);
}
</script>

<template>
  <div class="page">
    <header class="head">
      <div>
        <h1>社区种草</h1>
        <p>用户分享好物 · 审核通过后展示</p>
      </div>
      <button type="button" class="cta" @click="router.push('/community/create')">发帖</button>
    </header>
    <div class="topics">
      <button type="button" class="chip" :class="{ on: topicId === 0 }" @click="selectTopic(0)">
        全部
      </button>
      <button
        v-for="t in topics"
        :key="t.topic_id"
        type="button"
        class="chip"
        :class="{ on: topicId === t.topic_id }"
        @click="selectTopic(t.topic_id)"
      >
        {{ t.topic_name }}
      </button>
    </div>
    <p v-if="loading" class="hint">加载中…</p>
    <p v-else-if="!list.length" class="hint">暂无内容</p>
    <div v-else class="list">
      <article v-for="p in list" :key="p.community_id" class="card" @click="open(p.community_id)">
        <strong>{{ p.title }}</strong>
        <p class="meta">{{ p.nickname }} · {{ p.topic_name }}</p>
        <p class="body">{{ p.content }}</p>
        <p v-if="p.product_name" class="goods">
          {{ p.product_name }} · ¥{{ Number(p.product_price || 0).toFixed(2) }}
        </p>
      </article>
    </div>
  </div>
</template>

<style scoped>
.page {
  max-width: 800px;
  margin: 0 auto;
  padding: 32px 24px 64px;
}
.head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}
.head h1 {
  margin: 0 0 8px;
}
.head p,
.hint {
  color: #888;
}
.cta {
  border: 0;
  border-radius: 10px;
  padding: 10px 18px;
  background: #c45c26;
  color: #fff;
  cursor: pointer;
  white-space: nowrap;
}
.topics {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 20px 0 8px;
}
.chip {
  border: 1px solid #eee;
  background: #fff;
  border-radius: 999px;
  padding: 6px 14px;
  cursor: pointer;
  font-size: 13px;
  color: #666;
}
.chip.on {
  background: #c45c26;
  border-color: #c45c26;
  color: #fff;
}
.list {
  margin-top: 24px;
  display: grid;
  gap: 12px;
}
.card {
  background: #fff;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
}
.meta {
  margin: 8px 0;
  color: #999;
  font-size: 13px;
}
.body {
  margin: 0;
  color: #444;
}
.goods {
  margin: 12px 0 0;
  color: #c45c26;
  font-size: 14px;
}
</style>
