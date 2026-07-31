<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import AccountFrame from "@/components/AccountFrame.vue";
import {
  fetchCouponCenter,
  fetchMyCoupons,
  receiveCoupon,
  type Coupon,
  type CouponUser,
} from "@/api/coupon";
import { ApiError } from "@/utils/request";
import { useUserStore } from "@/stores/user";

type View = "mine" | "history" | "center";
const user = useUserStore();
const route = useRoute();
const router = useRouter();
const view = computed<View>(() => route.query.view === "center" ? "center" : route.query.tab === "history" ? "history" : "mine");
const center = ref<Coupon[]>([]);
const mine = ref<CouponUser[]>([]);
const hint = ref("");
const loading = ref(false);

function setView(next: View) {
  const query = next === "center" ? { view: "center" } : next === "history" ? { tab: "history" } : {};
  void router.replace({ name: "coupons", query });
}

function formatDate(value?: string | null) {
  if (!value) return "长期有效";
  return value.slice(0, 10);
}

function discountLabel(c: Pick<Coupon, "discount_type" | "discount_value"> | Pick<CouponUser, "discount_type" | "discount_value">) {
  return c.discount_type === "rate" ? `${Number(c.discount_value / 10).toFixed(1)}折` : `¥${Number(c.discount_value).toFixed(0)}`;
}

function kindLabel(c: Pick<Coupon, "type"> | Pick<CouponUser, "coupon_kind">) {
  return "type" in c ? (c.type ? "店铺券" : "平台券") : (c.coupon_kind ? "店铺券" : "通用券");
}

async function load() {
  if (!user.isLogin) {
    hint.value = "请先登录后查看优惠券";
    center.value = [];
    mine.value = [];
    return;
  }
  loading.value = true;
  try {
    if (view.value === "center") center.value = (await fetchCouponCenter()).list || [];
    else mine.value = (await fetchMyCoupons(view.value === "history" ? "history" : "unused")).list || [];
    hint.value = "";
  } catch (error) {
    hint.value = error instanceof ApiError ? error.message : "优惠券加载失败";
  } finally {
    loading.value = false;
  }
}

async function onReceive(couponID: number) {
  try {
    await receiveCoupon(couponID);
    await load();
  } catch (error) {
    hint.value = error instanceof ApiError ? error.message : "优惠券领取失败";
  }
}

watch(view, () => void load(), { immediate: true });
</script>

<template>
  <AccountFrame>
    <template #crumb><span>›</span> {{ view === "center" ? "领券中心" : "我的优惠券" }}</template>
    <div v-if="!user.isLogin" class="coupon-guest">
      <p>{{ hint }}</p>
      <button class="pc-btn" type="button" @click="router.push({ name: 'login', query: { redirect: route.fullPath } })">去登录</button>
    </div>
    <template v-else>
      <nav class="coupon-tabs" aria-label="优惠券分类">
        <button type="button" :class="{ active: view === 'mine' }" @click="setView('mine')">未使用</button>
        <button type="button" :class="{ active: view === 'history' }" @click="setView('history')">已使用/已过期</button>
        <button type="button" :class="{ active: view === 'center' }" @click="setView('center')">领券中心</button>
      </nav>
      <p v-if="hint" class="coupon-hint">{{ hint }}</p>
      <p v-else-if="loading" class="coupon-empty">正在加载优惠券…</p>
      <section v-else-if="view === 'center'" class="coupon-grid">
        <article v-for="coupon in center" :key="coupon.coupon_id" class="coupon-card">
          <div class="coupon-main">
            <span class="coupon-kind">{{ kindLabel(coupon) }}</span>
            <strong class="coupon-value">{{ discountLabel(coupon) }}</strong>
            <div class="coupon-copy"><b>{{ coupon.title }}</b><span>满{{ Number(coupon.use_min_price).toFixed(0) }}可用</span></div>
          </div>
          <button class="coupon-side" type="button" :disabled="coupon.received" @click="onReceive(coupon.coupon_id)">{{ coupon.received ? "已领取" : "立即领取" }}</button>
        </article>
        <p v-if="!center.length" class="coupon-empty">暂无可领取的优惠券</p>
      </section>
      <section v-else class="coupon-grid">
        <article v-for="coupon in mine" :key="coupon.coupon_user_id" class="coupon-card" :class="{ invalid: coupon.status !== 0 }">
          <div class="coupon-main">
            <span class="coupon-kind">{{ kindLabel(coupon) }}</span>
            <strong class="coupon-value">{{ discountLabel(coupon) }}</strong>
            <div class="coupon-copy"><b>{{ coupon.coupon_title }}</b><span>满{{ Number(coupon.use_min_price).toFixed(0) }}可用</span></div>
            <small>有效时间：{{ formatDate(coupon.starts_at) }}-{{ formatDate(coupon.ends_at) }}</small>
          </div>
          <span class="coupon-side">{{ coupon.status === 0 ? "可使用" : coupon.status === 2 ? "已使用" : "已过期" }}</span>
        </article>
        <p v-if="!mine.length" class="coupon-empty">{{ view === 'history' ? '暂无已使用或已过期优惠券' : '暂无未使用优惠券，去领券中心看看吧。' }}</p>
      </section>
    </template>
  </AccountFrame>
</template>

<style scoped>
.coupon-tabs { display: flex; gap: 70px; height: 74px; border-bottom: 1px solid #eee; }
.coupon-tabs button { position: relative; border: 0; padding: 0 16px; background: transparent; color: #999; font-size: 20px; cursor: pointer; }
.coupon-tabs button.active { color: #f13728; }
.coupon-tabs button.active::after { position: absolute; bottom: -1px; left: 50%; width: 140px; height: 3px; background: #f13728; content: ""; transform: translateX(-50%); }
.coupon-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 20px; padding-top: 38px; }
.coupon-card { display: grid; grid-template-columns: minmax(0, 1fr) 72px; min-height: 182px; overflow: hidden; background: #fff; box-shadow: 0 8px 23px rgb(0 0 0 / 8%); }
.coupon-main { position: relative; display: grid; grid-template-columns: 104px minmax(0, 1fr); align-content: center; gap: 9px 12px; padding: 21px 18px 18px 24px; }
.coupon-kind { position: absolute; top: 12px; left: 24px; padding: 4px 10px; background: #fff0ef; color: #f13728; font-size: 13px; }
.coupon-value { align-self: center; margin-top: 14px; color: #f13728; font-size: 41px; font-weight: 500; letter-spacing: -2px; }
.coupon-copy { display: grid; align-content: center; gap: 10px; margin-top: 14px; color: #333; }.coupon-copy b { overflow: hidden; font-size: 17px; font-weight: 500; text-overflow: ellipsis; white-space: nowrap; }.coupon-copy span, .coupon-main small { color: #999; font-size: 14px; }.coupon-main small { grid-column: 1 / -1; white-space: nowrap; }
.coupon-side { display: grid; width: 72px; place-items: center; border: 0; background: #f13728; color: #fff; font-size: 17px; line-height: 1.55; cursor: pointer; writing-mode: vertical-rl; }.coupon-side:disabled { cursor: default; opacity: .72; }.coupon-card.invalid { filter: grayscale(1); opacity: .67; }.coupon-card.invalid .coupon-side { background: #999; }
.coupon-empty { grid-column: 1 / -1; padding: 92px 0; color: #aaa; text-align: center; }.coupon-hint { margin: 24px 0 0; color: #d9362b; }.coupon-guest { display: grid; min-height: 380px; place-content: center; justify-items: center; gap: 14px; color: #777; }
@media (max-width: 960px) { .coupon-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }.coupon-tabs { gap: 22px; }.coupon-tabs button { font-size: 16px; } } @media (max-width: 620px) { .coupon-grid { grid-template-columns: 1fr; }.coupon-tabs { gap: 0; justify-content: space-between; }.coupon-tabs button { padding: 0 7px; font-size: 14px; }.coupon-tabs button.active::after { width: 70px; } }
</style>
