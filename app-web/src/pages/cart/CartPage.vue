<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchCart, removeCart, updateCartNum, type CartBucket } from "@/api/trade";
import { useUserStore } from "@/stores/user";
import { ApiError } from "@/utils/request";

const router = useRouter();
const user = useUserStore();
const buckets = ref<CartBucket[]>([]);
const loading = ref(false);
const hint = ref("");
const selected = ref<Set<number>>(new Set());

const flatIds = computed(() =>
  buckets.value.flatMap((b) => b.items.filter((i) => !i.is_fail).map((i) => i.cart_id)),
);

const selectedPrice = computed(() => {
  let sum = 0;
  for (const b of buckets.value) {
    for (const it of b.items) {
      if (selected.value.has(it.cart_id) && !it.is_fail) {
        sum += (it.price || 0) * it.cart_num;
      }
    }
  }
  return sum.toFixed(2);
});

async function load() {
  if (!user.isLogin) {
    hint.value = "请先登录";
    buckets.value = [];
    return;
  }
  loading.value = true;
  try {
    const data = await fetchCart();
    buckets.value = data.list || [];
    selected.value = new Set(flatIds.value);
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());

function toggle(id: number) {
  const next = new Set(selected.value);
  if (next.has(id)) next.delete(id);
  else next.add(id);
  selected.value = next;
}

async function changeNum(id: number, num: number) {
  if (num < 1) {
    await removeCart(id);
  } else {
    await updateCartNum(id, num);
  }
  await load();
}

function checkout() {
  const ids = [...selected.value];
  if (!ids.length) {
    hint.value = "请选择商品";
    return;
  }
  router.push({ name: "checkout", query: { cart_ids: ids.join(",") } });
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>我的购物车</h1>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <p v-if="!user.isLogin">
        <RouterLink class="pc-btn" to="/login?redirect=/cart">去登录</RouterLink>
      </p>
      <div v-else-if="loading" class="muted">加载中…</div>
      <div v-else-if="!buckets.length" class="empty">
        <strong>购物车还是空的</strong>
        <RouterLink class="pc-btn" to="/goods">去逛逛</RouterLink>
      </div>
      <template v-else>
        <div v-for="b in buckets" :key="b.mer_id" class="bucket">
          <h2>{{ b.mer_name || `商户 #${b.mer_id}` }}</h2>
          <div v-for="it in b.items" :key="it.cart_id" class="row">
            <label>
              <input
                type="checkbox"
                :checked="selected.has(it.cart_id)"
                :disabled="!!it.is_fail"
                @change="toggle(it.cart_id)"
              />
            </label>
            <div class="thumb"><img v-if="it.image" :src="it.image" :alt="it.title || it.store_name" /><span v-else>暂无图片</span></div>
            <div class="info">
              <strong>{{ it.title || it.store_name }}</strong>
              <span class="muted">¥{{ it.price }} · 库存 {{ it.stock }}</span>
              <span v-if="it.is_fail" class="fail">不可售</span>
            </div>
            <div class="qty">
              <button type="button" @click="changeNum(it.cart_id, it.cart_num - 1)">-</button>
              <span>{{ it.cart_num }}</span>
              <button type="button" @click="changeNum(it.cart_id, it.cart_num + 1)">+</button>
            </div>
          </div>
        </div>
        <footer class="bar">
          <span>合计 <strong>¥{{ selectedPrice }}</strong></span>
          <button class="pc-btn" type="button" @click="checkout">去结算</button>
        </footer>
      </template>
    </section>
  </div>
</template>

<style scoped>
.panel {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: 0;
  padding: 1.5rem;
  box-shadow: var(--pc-shadow);
}
h1 { margin: 0 0 0.5rem; }
.hint, .muted { color: var(--pc-muted); }
.empty {
  border: 1px dashed var(--pc-line);
  border-radius: 8px;
  padding: 2rem;
  display: grid;
  justify-items: center;
  gap: 1rem;
}
.bucket { margin-top: 1.2rem; }
.bucket h2 { font-size: 1rem; margin: 0 0 0.6rem; }
.row {
  display: grid;
  grid-template-columns: auto 84px 1fr auto;
  gap: 0.8rem;
  align-items: center;
  padding: 0.7rem 0;
  border-top: 1px solid var(--pc-line);
}
.thumb { width: 72px; height: 72px; background: #f6f6f6; display: grid; place-items: center; color: #999; font-size: .72rem; }
.thumb img { width: 100%; height: 100%; object-fit: cover; }
.info { display: grid; gap: 0.2rem; }
.fail { color: #c0392b; font-size: 0.85rem; }
.qty { display: flex; align-items: center; gap: 0.5rem; }
.qty button {
  width: 28px; height: 28px; border: 1px solid var(--pc-line); background: #fff; cursor: pointer;
}
.bar {
  margin-top: 1.2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid var(--pc-line);
  padding-top: 1rem;
}
.bar .pc-btn { border-radius: 0; min-width: 130px; }
</style>
