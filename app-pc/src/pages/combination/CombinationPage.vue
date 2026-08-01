<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchGroups, type ProductGroup } from "@/api/combination";

const router = useRouter();
const loading = ref(false);
const list = ref<ProductGroup[]>([]);

onMounted(async () => {
  loading.value = true;
  try {
    const data = await fetchGroups();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
});

function go(id: number) {
  router.push(`/goods/${id}`);
}
</script>

<template>
  <div class="page">
    <header class="head">
      <h1>拼团好物</h1>
      <p>独立开团/参团下单；成团前订单状态为拼团中；店铺券不可用</p>
    </header>
    <p v-if="loading" class="hint">加载中…</p>
    <p v-else-if="!list.length" class="hint">暂无拼团活动</p>
    <div v-else class="grid">
      <article
        v-for="g in list"
        :key="g.product_group_id"
        class="card"
        @click="go(g.product_id)"
      >
        <img v-if="g.image" :src="g.image" :alt="g.store_name || '拼团商品'" />
        <div class="row">
          <strong>{{ g.store_name || "拼团商品" }}</strong>
          <span class="badge">{{ g.buying_count_num }}人团</span>
        </div>
        <p class="mer">{{ g.mer_name }}</p>
        <p class="price">
          <em>¥{{ g.price }}</em>
          <span v-if="g.ot_price">¥{{ g.ot_price }}</span>
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
.head { display: flex; align-items: end; justify-content: space-between; min-height: 104px; padding: 0 30px; color: #3c2c28; background: #fff6f1 url('/demo/seckill-hero-v1.png') center / cover no-repeat; }
.head h1 {
  margin: 0 0 25px;
  font-size: 28px;
}
.head p,
.hint {
  color: #7b625c;
  margin: 0 0 28px;
  font-size: 13px;
}
.hint { margin: 44px 0; color: #999; text-align: center; }
.grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 15px;
  margin-top: 22px;
}
.card {
  background: #fff;
  border: 1px solid #eee;
  padding: 0 16px 16px;
  cursor: pointer;
  transition: transform .18s ease, box-shadow .18s ease;
}
.card:hover { transform: translateY(-3px); box-shadow: 0 8px 18px rgb(0 0 0 / 9%); }
.card > img { display: block; width: calc(100% + 32px); aspect-ratio: 1; margin: 0 -16px 14px; object-fit: cover; }
.row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
.badge {
  font-size: 12px;
  color: #e23030;
}
.mer {
  margin: 8px 0 0;
  font-size: 13px;
  color: #999;
}
.price {
  margin: 16px 0 0;
  display: flex;
  align-items: baseline;
  gap: 10px;
}
.price em {
  font-style: normal;
  color: #e23030;
  font-size: 22px;
  font-weight: 700;
}
.price span {
  color: #bbb;
  text-decoration: line-through;
  font-size: 13px;
}
@media (max-width: 860px) { .page { padding-inline: 16px; }.grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }.head { padding: 0 20px; } }
</style>
