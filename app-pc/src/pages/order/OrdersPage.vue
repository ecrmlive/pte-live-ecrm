<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { cancelOrder, fetchOrderDetail, fetchOrders, type GroupOrder } from "@/api/trade";
import { ApiError } from "@/utils/request";

const list = ref<GroupOrder[]>([]);
const hint = ref("");
const loading = ref(false);
const statusFilter = ref<"all" | "pending" | "paid" | "aftersale">("all");
const filters: { key: "all" | "pending" | "paid" | "aftersale"; label: string }[] = [
  { key: "all", label: "全部订单" }, { key: "pending", label: "待付款" },
  { key: "paid", label: "待收货/已支付" }, { key: "aftersale", label: "售后中" },
];

function setStatusFilter(value: "all" | "pending" | "paid" | "aftersale") {
  statusFilter.value = value;
}

const filtered = computed(() => list.value.filter((group) => {
  if (statusFilter.value === "all") return true;
  if (statusFilter.value === "pending") return group.pay_status === "pending";
  if (statusFilter.value === "paid") return group.paid === 1 && !group.orders?.some((order) => order.status === "aftersale");
  return group.orders?.some((order) => order.status === "aftersale");
}));

async function load() {
  loading.value = true;
  try {
    const res = await fetchOrders(1);
    const groups = res.list || [];
    list.value = await Promise.all(groups.map(async (group) => {
      try { return await fetchOrderDetail(group.group_order_id); }
      catch { return group; }
    }));
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());

function statusText(g: GroupOrder) {
  if (g.pay_status === "closed") return "已取消";
  if (g.pay_status === "refunded") return "已退款";
  if (g.paid !== 1) return "待付款";
  return "已支付";
}

async function cancel(g: GroupOrder) {
  if (!window.confirm("确定取消该订单吗？已锁定的优惠券会恢复为可用状态。")) return;
  try {
    await cancelOrder(g.group_order_id);
    await load();
  } catch (error) {
    hint.value = error instanceof ApiError ? error.message : "取消订单失败";
  }
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>我的订单</h1>
      <nav class="filters"><button v-for="item in filters" :key="item.key" type="button" :class="{ active: statusFilter === item.key }" @click="setStatusFilter(item.key)">{{ item.label }}</button></nav>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div v-if="loading" class="empty">正在加载订单…</div>
      <div v-else-if="!filtered.length && !hint" class="empty">暂无该状态订单</div>
      <article v-for="g in filtered" :key="g.group_order_id" class="card">
        <header class="group-header"><strong>订单号：{{ g.group_order_sn }}</strong><span>{{ statusText(g) }}</span></header>
        <section v-for="order in g.orders" :key="order.order_id" class="store-order">
          <div class="store-name">{{ order.mer_name || '店铺订单' }}</div>
          <div v-for="product in order.products" :key="`${order.order_id}-${product.product_id}-${product.product_attr_unique}`" class="product-row"><img v-if="product.image" :src="product.image" :alt="product.store_name" /><div class="product-fallback" v-else>暂无图片</div><p>{{ product.store_name }}<small>规格：{{ product.product_attr_unique || '默认' }}</small></p><span>¥{{ Number(product.product_price).toFixed(2) }} × {{ product.product_num }}</span></div>
          <footer><span>店铺小计 ¥{{ Number(order.pay_price).toFixed(2) }}</span><RouterLink v-if="g.paid === 1 && ['paid', 'fulfilling', 'shipped'].includes(order.status)" :to="`/refunds?order_id=${order.order_id}`">申请售后</RouterLink></footer>
        </section>
        <footer class="group-footer"><span>共 {{ g.total_num }} 件商品，应付 <b>¥{{ Number(g.pay_price).toFixed(2) }}</b></span><button v-if="g.pay_status === 'pending'" class="cancel" type="button" @click="cancel(g)">取消订单</button><RouterLink v-if="g.pay_status === 'pending' || g.paid === 1" class="link" :to="`/pay/${g.group_order_id}`">{{ g.paid === 1 ? '订单详情' : '去支付' }}</RouterLink></footer>
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
.filters { display: flex; gap: 34px; margin: 18px 0 0; border-bottom: 1px solid #eee; }.filters button { position: relative; border: 0; padding: 0 2px 16px; color: #666; background: transparent; cursor: pointer; }.filters button.active { color: #f13728; }.filters button.active::after { position: absolute; bottom: -1px; left: 0; width: 100%; height: 2px; background: #f13728; content: ""; }
.card { margin-top: 18px; border: 1px solid #eee; }
.group-header {
  display: flex;
  justify-content: space-between;
  padding: 12px 18px;
  color: #555;
  background: #fafafa;
  font-size: 14px;
}
.group-header span { color: #f13728; }.store-order { padding: 0 18px; }.store-name { padding: 14px 0; border-bottom: 1px dashed #eee; color: #555; font-size: 14px; }.product-row { display: grid; grid-template-columns: 72px minmax(0, 1fr) auto; gap: 14px; align-items: center; padding: 14px 0; border-bottom: 1px solid #eee; }.product-row img, .product-fallback { width: 60px; height: 60px; object-fit: cover; background: #f5f5f5; }.product-fallback { display: grid; place-items: center; color: #aaa; font-size: 10px; }.product-row p { margin: 0; color: #555; }.product-row small { display: block; margin-top: 7px; color: #999; }.product-row > span { color: #666; font-size: 14px; }.store-order footer, .group-footer { display: flex; align-items: center; justify-content: flex-end; gap: 12px; padding: 14px 0; color: #777; font-size: 14px; }.store-order footer a, .link { color: #f13728; }.group-footer { padding: 15px 18px; border-top: 1px solid #eee; }.group-footer b { color: #f13728; font-size: 16px; }.group-footer .link, .group-footer .cancel { border: 1px solid #f13728; padding: 7px 14px; color: #f13728; background: #fff; cursor: pointer; }.group-footer .link { color: #fff; background: #f13728; }@media (max-width: 650px) { .filters { gap: 16px; overflow-x: auto; }.product-row { grid-template-columns: 62px minmax(0, 1fr); }.product-row > span { grid-column: 2; }.group-footer { align-items: flex-end; } }
</style>
