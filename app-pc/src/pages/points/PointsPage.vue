<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { http } from "@/utils/request";
import { addCart } from "@/api/trade";
import { ApiError } from "@/utils/request";
import { useUserStore } from "@/stores/user";

interface PointsProduct {
  id: number;
  store_name: string;
  price: string;
  integral: number;
}

const user = useUserStore();
const router = useRouter();
const list = ref<PointsProduct[]>([]);
const hint = ref("");

async function load() {
  try {
    const res = await http.get<{ list: PointsProduct[] }>("/catalog/points/products");
    list.value = res.list || [];
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  }
}

onMounted(() => void load());

async function buy(p: PointsProduct) {
  if (!user.isLogin) {
    router.push({ name: "login", query: { redirect: "/points" } });
    return;
  }
  try {
    const item = await addCart({ product_id: p.id, cart_num: 1 });
    router.push({ name: "points-checkout", query: { cart_ids: String(item.cart_id) } });
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加购失败";
  }
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>积分商城</h1>
      <p class="muted">积分商品单独下单（`/order/v3`），与普通 v2 入口隔离。</p>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <ul class="list">
        <li v-for="p in list" :key="p.id">
          <div>
            <strong>{{ p.store_name }}</strong>
            <p>{{ p.integral }} 积分 + ¥{{ p.price }}</p>
          </div>
          <button class="pc-btn" type="button" @click="buy(p)">兑换</button>
        </li>
        <li v-if="!list.length" class="muted">暂无积分商品</li>
      </ul>
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
.muted { color: var(--pc-muted); }
.hint { color: #c0392b; }
.list { list-style: none; padding: 0; margin: 1rem 0 0; }
.list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding: 0.9rem 0;
  border-bottom: 1px solid var(--pc-line);
}
.list p { margin: 0.3rem 0 0; color: var(--pc-muted); font-size: 0.9rem; }
</style>
