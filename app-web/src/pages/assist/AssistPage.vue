<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchAssists, type ProductAssist } from "@/api/assist";

const router = useRouter();
const loading = ref(false);
const list = ref<ProductAssist[]>([]);

onMounted(async () => {
  loading.value = true;
  try {
    const data = await fetchAssists();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
});

function go(id: number) {
  router.push(`/assist/${id}`);
}
</script>

<template>
  <div class="page">
    <header class="head">
      <h1>好友助力</h1>
      <p>邀请好友助力满员后按助力价下单；演示活动需 1 人（种子 set 已满员）</p>
    </header>
    <p v-if="loading" class="hint">加载中…</p>
    <p v-else-if="!list.length" class="hint">暂无助力活动</p>
    <div v-else class="grid">
      <article
        v-for="a in list"
        :key="a.product_assist_id"
        class="card"
        @click="go(a.product_assist_id)"
      >
        <div class="row">
          <strong>{{ a.store_name || "助力商品" }}</strong>
          <span class="badge">需{{ a.assist_count }}人</span>
        </div>
        <p class="mer">{{ a.mer_name }}</p>
        <p class="price">
          <em>¥{{ a.assist_price }}</em>
          <span v-if="a.ot_price">¥{{ a.ot_price }}</span>
        </p>
      </article>
    </div>
  </div>
</template>

<style scoped>
.page {
  max-width: 1080px;
  margin: 0 auto;
  padding: 32px 24px 64px;
}
.head h1 {
  margin: 0 0 8px;
  font-size: 28px;
}
.head p,
.hint {
  color: #888;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  margin-top: 24px;
}
.card {
  background: #fff;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
}
.row {
  display: flex;
  justify-content: space-between;
  gap: 8px;
}
.badge {
  font-size: 12px;
  color: #e23030;
  background: #fff1f0;
  padding: 2px 8px;
  border-radius: 999px;
}
.mer {
  color: #999;
  font-size: 13px;
  margin: 8px 0;
}
.price em {
  color: #e23030;
  font-style: normal;
  font-size: 22px;
  font-weight: 700;
  margin-right: 8px;
}
.price span {
  color: #bbb;
  text-decoration: line-through;
  font-size: 13px;
}
</style>
