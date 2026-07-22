<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import {
  createCommunityPost,
  fetchCommunityTopics,
  type CommunityTopic,
} from "@/api/community";

const router = useRouter();
const title = ref("");
const content = ref("");
const productId = ref("");
const topicId = ref(1);
const topics = ref<CommunityTopic[]>([]);
const msg = ref("");
const submitting = ref(false);

onMounted(async () => {
  try {
    const data = await fetchCommunityTopics();
    topics.value = data.list || [];
    if (topics.value.length) topicId.value = topics.value[0].topic_id;
  } catch {
    topics.value = [{ topic_id: 1, topic_name: "好物推荐" }];
  }
});

async function submit() {
  const t = title.value.trim();
  const c = content.value.trim();
  if (!t || !c) {
    msg.value = "请填写标题和内容";
    return;
  }
  submitting.value = true;
  msg.value = "";
  try {
    await createCommunityPost({
      title: t,
      content: c,
      product_id: Number(productId.value) || 0,
      topic_id: topicId.value || 1,
      category_id: 1,
    });
    msg.value = "已提交审核";
    setTimeout(() => router.replace("/community"), 400);
  } catch (e) {
    msg.value = (e as Error).message || "提交失败";
  } finally {
    submitting.value = false;
  }
}
</script>

<template>
  <div class="page">
    <header class="head">
      <h1>发帖</h1>
      <p>发布后需平台审核通过才会展示</p>
    </header>
    <label class="field">
      <span>标题</span>
      <input v-model="title" type="text" maxlength="64" placeholder="种草标题" />
    </label>
    <label class="field">
      <span>内容</span>
      <textarea v-model="content" rows="6" placeholder="分享你的种草心得…" />
    </label>
    <div class="field">
      <span>话题</span>
      <div class="topics">
        <button
          v-for="t in topics"
          :key="t.topic_id"
          type="button"
          class="chip"
          :class="{ on: topicId === t.topic_id }"
          @click="topicId = t.topic_id"
        >
          {{ t.topic_name }}
        </button>
      </div>
    </div>
    <label class="field">
      <span>挂货商品 ID（可选）</span>
      <input v-model="productId" type="number" min="0" placeholder="如 1" />
    </label>
    <p v-if="msg" class="hint">{{ msg }}</p>
    <button class="cta" type="button" :disabled="submitting" @click="submit">
      {{ submitting ? "提交中…" : "提交审核" }}
    </button>
  </div>
</template>

<style scoped>
.page {
  max-width: 640px;
  margin: 0 auto;
  padding: 1.5rem 1.25rem 3rem;
}
.head h1 {
  margin: 0;
  font-size: 1.5rem;
}
.head p {
  margin: 0.4rem 0 1.4rem;
  color: var(--pc-muted, #888);
}
.field {
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
  margin-bottom: 1rem;
}
.field span {
  font-size: 0.9rem;
  color: var(--pc-muted, #666);
}
input,
textarea {
  border: 1px solid var(--pc-line, #e5e5e5);
  border-radius: 10px;
  padding: 0.7rem 0.85rem;
  font: inherit;
}
.topics {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}
.chip {
  border: 1px solid var(--pc-line, #e5e5e5);
  background: #fff;
  border-radius: 999px;
  padding: 0.35rem 0.9rem;
  cursor: pointer;
}
.chip.on {
  border-color: var(--pc-brand, #e23030);
  color: var(--pc-brand, #e23030);
}
.hint {
  color: var(--pc-muted, #888);
  margin: 0 0 0.8rem;
}
.cta {
  width: 100%;
  border: 0;
  border-radius: 10px;
  padding: 0.85rem;
  background: var(--pc-brand, #e23030);
  color: #fff;
  font: inherit;
  cursor: pointer;
}
.cta:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
