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
        <img v-if="a.image" :src="a.image" :alt="a.store_name || '助力商品'" />
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
  max-width: 1200px;
  margin: 0 auto;
  padding: 28px 0 64px;
}
.head { display: flex; align-items: end; justify-content: space-between; gap: 24px; min-height: 104px; padding: 0 30px; color: #3c2c28; background: #fff6f1 url('/demo/seckill-hero-v1.png') center / cover no-repeat; }.head h1 { margin: 0 0 25px; font-size: 28px; white-space: nowrap; }.head p { max-width: 560px; margin: 0 0 28px; color: #7b625c; font-size: 13px; line-height: 1.55; text-align: right; }.hint { margin: 44px 0; color: #999; text-align: center; }
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  margin-top: 24px;
}
.card {
  background: #fff;
  border: 1px solid #eee;
  padding: 0 16px 16px;
  cursor: pointer;
  transition: transform .18s ease, box-shadow .18s ease;
}
.card:hover { transform: translateY(-3px); box-shadow: 0 8px 18px rgb(0 0 0 / 9%); }.card > img { display: block; width: calc(100% + 32px); aspect-ratio: 1; margin: 0 -16px 14px; object-fit: cover; }
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
@media (max-width: 860px) { .page { padding-inline: 16px; }.grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }.head { align-items: flex-start; flex-direction: column; justify-content: center; padding: 16px 20px; }.head h1,.head p { margin: 0; text-align: left; } }
</style>
