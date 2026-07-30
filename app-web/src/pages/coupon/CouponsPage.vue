<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import {
  fetchCouponCenter,
  fetchMyCoupons,
  receiveCoupon,
  type Coupon,
  type CouponUser,
} from "@/api/coupon";
import { ApiError } from "@/utils/request";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const router = useRouter();
const tab = ref<"center" | "mine">("center");
const center = ref<Coupon[]>([]);
const mine = ref<CouponUser[]>([]);
const hint = ref("");
const loading = ref(false);

async function load() {
  if (!user.isLogin) {
    hint.value = "请先登录";
    return;
  }
  loading.value = true;
  try {
    if (tab.value === "center") {
      const res = await fetchCouponCenter();
      center.value = res.list || [];
    } else {
      const res = await fetchMyCoupons(0);
      mine.value = res.list || [];
    }
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());

async function onReceive(id: number) {
  try {
    await receiveCoupon(id);
    await load();
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "领取失败";
  }
}

function switchTab(t: "center" | "mine") {
  tab.value = t;
  void load();
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>优惠券</h1>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div v-if="!user.isLogin" class="guest">
        <button class="pc-btn" type="button" @click="router.push({ name: 'login', query: { redirect: '/coupons' } })">
          去登录
        </button>
      </div>
      <template v-else>
        <div class="tabs">
          <button type="button" :class="{ active: tab === 'center' }" @click="switchTab('center')">领券中心</button>
          <button type="button" :class="{ active: tab === 'mine' }" @click="switchTab('mine')">我的未用券</button>
        </div>
        <div v-if="loading" class="muted">加载中…</div>
        <ul v-else-if="tab === 'center'" class="list">
          <li v-for="c in center" :key="c.coupon_id">
            <div>
              <strong>{{ c.title }}</strong>
              <p>¥{{ Number(c.coupon_price).toFixed(2) }} · 满 {{ c.use_min_price }} 可用</p>
            </div>
            <button
              class="pc-btn"
              type="button"
              :disabled="c.received"
              @click="onReceive(c.coupon_id)"
            >
              {{ c.received ? "已领取" : "领取" }}
            </button>
          </li>
          <li v-if="!center.length" class="muted">暂无可领平台券</li>
        </ul>
        <ul v-else class="list">
          <li v-for="c in mine" :key="c.coupon_user_id">
            <div>
              <strong>{{ c.coupon_title }}</strong>
              <p>¥{{ Number(c.coupon_price).toFixed(2) }} · 满 {{ c.use_min_price }} · #{{ c.coupon_user_id }}</p>
            </div>
          </li>
          <li v-if="!mine.length" class="muted">暂无未使用优惠券</li>
        </ul>
      </template>
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
.tabs { display: flex; gap: 0.6rem; margin: 1rem 0; }
.tabs button {
  border: 1px solid var(--pc-line);
  background: transparent;
  border-radius: 8px;
  padding: 0.4rem 0.9rem;
  cursor: pointer;
}
.tabs button.active {
  border-color: var(--pc-accent, #127a5d);
  color: var(--pc-accent, #127a5d);
}
.list { list-style: none; padding: 0; margin: 0; }
.list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding: 0.9rem 0;
  border-bottom: 1px solid var(--pc-line);
}
.list p { margin: 0.3rem 0 0; color: var(--pc-muted); font-size: 0.9rem; }
.muted { color: var(--pc-muted); }
</style>
