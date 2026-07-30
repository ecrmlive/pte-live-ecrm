<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import {
  fetchDaySlots,
  fetchReserveProducts,
  reservationCreate,
  type ReserveProduct,
  type SlotDay,
} from "@/api/reservation";
import { fetchAddresses } from "@/api/trade";

const router = useRouter();
const loading = ref(false);
const list = ref<ReserveProduct[]>([]);
const active = ref<ReserveProduct | null>(null);
const days = ref<string[]>([]);
const date = ref("");
const slots = ref<SlotDay[]>([]);
const slotId = ref(0);
const booking = ref(false);

function buildDays(n: number) {
  const out: string[] = [];
  const base = new Date();
  for (let i = 0; i < n; i++) {
    const d = new Date(base);
    d.setDate(base.getDate() + i);
    out.push(
      `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`,
    );
  }
  return out;
}

onMounted(async () => {
  loading.value = true;
  try {
    const data = await fetchReserveProducts();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
});

async function openBook(p: ReserveProduct) {
  active.value = p;
  days.value = buildDays(Math.max(1, p.show_reservation_days || 7));
  date.value = days.value[0] || "";
  slotId.value = 0;
  await loadSlots();
}

async function loadSlots() {
  if (!active.value || !date.value) return;
  const data = await fetchDaySlots(active.value.product_id, date.value);
  slots.value = data.list || [];
  slotId.value = 0;
}

async function pickDate(d: string) {
  date.value = d;
  await loadSlots();
}

async function submit() {
  if (!active.value || !slotId.value || !date.value) return;
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
    const g = await reservationCreate({
      product_id: active.value.product_id,
      slot_id: slotId.value,
      date: date.value,
      address_id: a.address_id,
    });
    router.push(`/pay/${g.group_order_id}`);
  } catch (e) {
    alert((e as Error).message || "预约失败");
  } finally {
    booking.value = false;
  }
}
</script>

<template>
  <div class="page">
    <header class="head">
      <h1>预约服务</h1>
      <p>选日期时段下单，支付后凭核销码到店（店员端可核销）</p>
    </header>
    <p v-if="loading" class="hint">加载中…</p>
    <p v-else-if="!list.length" class="hint">暂无预约服务</p>
    <div v-else class="grid">
      <article
        v-for="p in list"
        :key="p.product_id"
        class="card"
        @click="openBook(p)"
      >
        <strong>{{ p.store_name }}</strong>
        <p class="mer">{{ p.mer_name }}</p>
        <p class="price"><em>¥{{ p.price }}</em></p>
      </article>
    </div>

    <div v-if="active" class="panel">
      <h2>预约：{{ active.store_name }}</h2>
      <div class="days">
        <button
          v-for="d in days"
          :key="d"
          type="button"
          :class="{ on: d === date }"
          @click="pickDate(d)"
        >
          {{ d.slice(5) }}
        </button>
      </div>
      <div class="slots">
        <button
          v-for="s in slots"
          :key="s.attr_reservation_id"
          type="button"
          :disabled="s.remain <= 0"
          :class="{ on: s.attr_reservation_id === slotId }"
          @click="slotId = s.attr_reservation_id"
        >
          {{ s.label || `${s.start_time}-${s.end_time}` }} · 余{{ s.remain }}
        </button>
      </div>
      <button class="cta" type="button" :disabled="booking || !slotId" @click="submit">
        确认预约并支付
      </button>
      <button class="ghost" type="button" @click="active = null">取消</button>
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
.hint,
.mer {
  color: #888;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
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
.price em {
  font-style: normal;
  color: #e23030;
  font-size: 22px;
  font-weight: 700;
}
.panel {
  margin-top: 28px;
  padding: 24px;
  border: 1px solid #eee;
  border-radius: 12px;
  background: #fff;
}
.days,
.slots {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin: 12px 0 16px;
}
.days button,
.slots button {
  border: 1px solid #ddd;
  background: #fafafa;
  border-radius: 8px;
  padding: 8px 12px;
  cursor: pointer;
}
.days button.on,
.slots button.on {
  border-color: #e23030;
  color: #e23030;
  background: #fff1f0;
}
.cta {
  margin-right: 12px;
  background: #e23030;
  color: #fff;
  border: 0;
  border-radius: 8px;
  padding: 10px 18px;
  cursor: pointer;
}
.ghost {
  border: 0;
  background: transparent;
  color: #888;
  cursor: pointer;
}
</style>
