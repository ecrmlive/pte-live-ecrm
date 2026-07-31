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
const allSelected = computed(() => flatIds.value.length > 0 && flatIds.value.every((id) => selected.value.has(id)));
const selectedCount = computed(() => [...selected.value].filter((id) => flatIds.value.includes(id)).length);

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

function toggleAll() {
  selected.value = allSelected.value ? new Set() : new Set(flatIds.value);
}

async function changeNum(id: number, num: number) {
  try {
    if (num < 1) {
      await removeCart(id);
    } else {
      await updateCartNum(id, num);
    }
    await load();
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "购物车更新失败";
  }
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
        <div class="table-head" role="row">
          <label class="select-all"><input type="checkbox" :checked="allSelected" @change="toggleAll" />全选</label>
          <span>商品信息</span><span>单价</span><span>数量</span><span>小计</span><span>操作</span>
        </div>
        <div v-for="b in buckets" :key="b.mer_id" class="bucket">
          <h2><span class="store-dot" />{{ b.mer_name || `商户 #${b.mer_id}` }}</h2>
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
              <span class="muted">库存 {{ it.stock }}</span>
              <span v-if="it.is_fail" class="fail">不可售</span>
            </div>
            <span class="unit-price">¥{{ Number(it.price).toFixed(2) }}</span>
            <div class="qty">
              <button type="button" @click="changeNum(it.cart_id, it.cart_num - 1)">-</button>
              <span>{{ it.cart_num }}</span>
              <button type="button" @click="changeNum(it.cart_id, it.cart_num + 1)">+</button>
            </div>
            <strong class="subtotal">¥{{ ((it.price || 0) * it.cart_num).toFixed(2) }}</strong>
            <button class="remove" type="button" @click="changeNum(it.cart_id, 0)">删除</button>
          </div>
        </div>
        <footer class="bar">
          <span>已选 <b>{{ selectedCount }}</b> 件商品</span>
          <span class="bar-total">合计 <strong>¥{{ selectedPrice }}</strong></span>
          <button class="pc-btn" type="button" @click="checkout">去结算</button>
        </footer>
      </template>
    </section>
  </div>
</template>

<style scoped>
.panel {
  margin-top: 22px;
  background: var(--pc-surface);
  border: 0;
  border-radius: 0;
  padding: 1.5rem;
  box-shadow: var(--pc-shadow);
}
h1 { margin: 0 0 1.15rem; color: #252525; font-size: 1.4rem; }
.hint, .muted { color: var(--pc-muted); }
.empty {
  border: 1px dashed var(--pc-line);
  border-radius: 8px;
  padding: 2rem;
  display: grid;
  justify-items: center;
  gap: 1rem;
}
.table-head {
  display: grid;
  grid-template-columns: 138px minmax(310px, 1fr) 110px 132px 120px 72px;
  align-items: center;
  min-height: 46px;
  padding: 0 20px;
  color: #666;
  background: #f5f5f5;
  font-size: .9rem;
}
.select-all { display: flex; align-items: center; gap: 8px; }
.bucket { margin-top: 16px; border: 1px solid #ececec; }
.bucket h2 { display: flex; align-items: center; gap: 8px; min-height: 42px; margin: 0; padding: 0 18px; border-bottom: 1px solid #eee; color: #444; font-size: .95rem; font-weight: 600; }
.store-dot { width: 14px; height: 14px; border: 1px solid #aaa; border-radius: 2px; }
.row {
  display: grid;
  grid-template-columns: 34px 84px minmax(220px, 1fr) 110px 132px 120px 72px;
  gap: 0.8rem;
  align-items: center;
  min-height: 118px;
  padding: 14px 18px;
  border-top: 1px solid #f0f0f0;
}
.row:first-of-type { border-top: 0; }
.thumb { width: 72px; height: 72px; background: #f6f6f6; display: grid; place-items: center; color: #999; font-size: .72rem; }
.thumb img { width: 100%; height: 100%; object-fit: cover; }
.info { display: grid; gap: 0.2rem; }
.info strong { line-height: 1.45; }
.fail { color: #c0392b; font-size: 0.85rem; }
.unit-price { color: #555; font-size: .92rem; }
.qty { display: flex; align-items: center; gap: 0.5rem; }
.qty button {
  width: 28px; height: 28px; border: 1px solid var(--pc-line); background: #fff; cursor: pointer;
}
.subtotal { color: #e64335; font-size: .96rem; }
.remove { border: 0; padding: 0; color: #666; background: transparent; font-size: .88rem; }
.remove:hover { color: #e64335; }
.bar {
  margin-top: 1.2rem;
  display: flex;
  justify-content: flex-end;
  gap: 26px;
  align-items: center;
  border-top: 1px solid var(--pc-line);
  padding-top: 1rem;
}
.bar-total { color: #555; }
.bar-total strong { margin-left: 6px; color: #e64335; font-size: 1.35rem; }
.bar .pc-btn { border-radius: 0; min-width: 130px; }
@media (max-width: 900px) { .table-head { display: none; }.row { grid-template-columns: 26px 72px 1fr auto; }.unit-price,.subtotal,.remove { display: none; } }
</style>
