<script setup lang="ts">
import { onMounted, ref } from "vue";
import { fetchOrders, type GroupOrder } from "@/api/trade";
import { ApiError } from "@/utils/request";

const list = ref<GroupOrder[]>([]);
const hint = ref("");

onMounted(async () => {
  try {
    const res = await fetchOrders(1);
    list.value = res.list || [];
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  }
});

function statusText(g: GroupOrder) {
  if (g.paid !== 1) return "待付款";
  return "已支付";
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>我的订单</h1>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div v-if="!list.length && !hint" class="empty">暂无订单</div>
      <article v-for="g in list" :key="g.group_order_id" class="card">
        <header>
          <strong>{{ g.group_order_sn }}</strong>
          <span>{{ statusText(g) }}</span>
        </header>
        <p>¥{{ Number(g.pay_price).toFixed(2) }} · {{ g.total_num }} 件</p>
        <RouterLink class="link" :to="`/pay/${g.group_order_id}`">详情 / 支付</RouterLink>
      </article>
    </section>
  </div>
</template>

<style scoped>
.panel {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: calc(var(--pc-radius) + 4px);
  padding: 1.5rem;
  box-shadow: var(--pc-shadow);
}
.hint { color: #c0392b; }
.empty { color: var(--pc-muted); padding: 2rem 0; }
.card {
  border-top: 1px solid var(--pc-line);
  padding: 1rem 0;
}
.card header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.4rem;
}
.link { color: var(--pc-accent, #1f6f8b); }
</style>
