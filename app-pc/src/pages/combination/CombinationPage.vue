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
</style>
