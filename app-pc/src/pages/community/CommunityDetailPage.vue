<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchCommunityPost, type CommunityPost } from "@/api/community";

const route = useRoute();
const router = useRouter();
const post = ref<CommunityPost | null>(null);

onMounted(async () => {
  const id = Number(route.params.id || 0);
  if (!id) return;
  try {
    post.value = await fetchCommunityPost(id);
  } catch (e) {
    alert((e as Error).message || "加载失败");
  }
});

function goGoods(id: number) {
  router.push(`/goods/${id}`);
}
</script>

<template>
  <div class="page" v-if="post">
    <h1>{{ post.title }}</h1>
    <p class="meta">{{ post.nickname }} · {{ post.topic_name }} · 浏览 {{ post.pv || 0 }}</p>
    <p class="body">{{ post.content }}</p>
    <button v-if="post.product_id" type="button" class="goods" @click="goGoods(post.product_id!)">
      {{ post.product_name || `商品 #${post.product_id}` }} · ¥{{ Number(post.product_price || 0).toFixed(2) }}
    </button>
  </div>
</template>

<style scoped>
.page {
  max-width: 720px;
  margin: 0 auto;
  padding: 32px 24px 64px;
}
.meta {
  color: #999;
}
.body {
  margin: 24px 0;
  line-height: 1.7;
  color: #333;
}
.goods {
  border: 1px solid #eee;
  background: #fff;
  border-radius: 8px;
  padding: 14px 16px;
  cursor: pointer;
  color: #c45c26;
}
</style>
