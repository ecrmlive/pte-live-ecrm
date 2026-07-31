<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import {
  fetchPresellFinals,
  fetchPresells,
  payPresellFinal,
  presellCreate,
  type PresellActive,
  type PresellFinal,
} from "@/api/presell";
import { fetchAddresses } from "@/api/trade";

const router = useRouter();
const loading = ref(false);
const booking = ref(false);
const list = ref<PresellActive[]>([]);
const finals = ref<PresellFinal[]>([]);

async function loadFinals() {
  try {
    const data = await fetchPresellFinals();
    finals.value = data.list || [];
  } catch {
    finals.value = [];
  }
}

onMounted(async () => {
  loading.value = true;
  try {
    const data = await fetchPresells();
    list.value = data.list || [];
    await loadFinals();
  } finally {
    loading.value = false;
  }
});

async function payFinal(row: PresellFinal) {
  try {
    await payPresellFinal(row.presell_order_id, "balance");
    alert("余额支付成功");
    await loadFinals();
  } catch (e) {
    alert((e as Error).message || "支付失败");
  }
}

async function buy(row: PresellActive) {
  booking.value = true;
  try {
    const addr = await fetchAddresses();
    const arr = addr.list || [];
    if (!arr.length) {
      alert("请先在「我的」添加收货地址");
      router.push("/user");
      return;
    }
    const a = arr.find((x) => x.is_default === 1) || arr[0];
    const g = await presellCreate({
      product_presell_id: row.product_presell_id,
      cart_num: 1,
      address_id: a.address_id,
    });
    router.push(`/pay/${g.group_order_id}`);
  } catch (e) {
    alert((e as Error).message || "下单失败");
  } finally {
    booking.value = false;
  }
}
</script>

<template>
  <div class="page">
    <header class="head">
      <h1>预售专区</h1>
      <p>全款一次付清；定金活动先付定金，再付尾款后发货</p>
    </header>

    <section v-if="finals.length" class="finals">
      <h2>待付尾款</h2>
      <article v-for="f in finals" :key="f.presell_order_id" class="card final">
        <div class="row">
          <strong>{{ f.store_name || "预售尾款" }}</strong>
          <em>¥{{ f.pay_price }}</em>
        </div>
        <p class="mer">{{ f.presell_order_sn }}</p>
        <button type="button" class="btn" @click="payFinal(f)">余额支付尾款</button>
      </article>
    </section>

    <p v-if="loading" class="hint">加载中…</p>
    <p v-else-if="!list.length" class="hint">暂无预售活动</p>
    <div v-else class="grid">
      <article v-for="g in list" :key="g.product_presell_id" class="card">
        <div class="row">
          <strong>{{ g.store_name || "预售商品" }}</strong>
          <span class="badge" :class="{ on: g.in_window }">
            {{ g.in_window ? "进行中" : "未开始" }}
          </span>
        </div>
        <p class="mer">{{ g.mer_name }} · 库存 {{ g.stock }}</p>
        <p class="price">
          <em>¥{{ g.price }}</em>
          <span v-if="g.ot_price">¥{{ g.ot_price }}</span>
        </p>
        <p v-if="g.presell_type === 2" class="mer">
          定金 ¥{{ g.down_price }} · 尾款 ¥{{ g.final_price }}
        </p>
        <button type="button" class="btn" :disabled="booking || !g.in_window" @click="buy(g)">
          {{ g.presell_type === 2 ? "支付定金" : "立即预订" }}
        </button>
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
.finals {
  margin: 24px 0;
}
.finals h2 {
  margin: 0 0 12px;
  font-size: 18px;
}
.finals .final {
  margin-bottom: 12px;
}
.finals em {
  font-style: normal;
  color: #c45c26;
  font-weight: 700;
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
}
.row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
.badge {
  font-size: 12px;
  color: #999;
}
.badge.on {
  color: #c45c26;
}
.mer {
  margin: 8px 0 0;
  font-size: 13px;
  color: #999;
}
.price {
  margin: 16px 0;
  display: flex;
  align-items: baseline;
  gap: 10px;
}
.price em {
  font-style: normal;
  font-size: 22px;
  color: #c45c26;
  font-weight: 700;
}
.price span {
  color: #bbb;
  text-decoration: line-through;
  font-size: 13px;
}
.btn {
  width: 100%;
  height: 40px;
  border: 0;
  border-radius: 8px;
  background: #1a1a1a;
  color: #fff;
  cursor: pointer;
}
.btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>
