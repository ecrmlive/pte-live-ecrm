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
const centerKind = ref<"all" | "platform" | "store">("all");
const visibleCenter = computed(() => centerKind.value === "all"
  ? center.value
  : center.value.filter((coupon) => centerKind.value === "store" ? coupon.type === 1 : coupon.type === 0),
);

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
  if (view.value !== "center" && !user.isLogin) {
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
  if (!user.isLogin) {
    await router.push({ name: "login", query: { redirect: route.fullPath } });
    return;
  }
  try {
    await receiveCoupon(couponID);
    center.value = center.value.map((coupon) => coupon.coupon_id === couponID ? { ...coupon, received: true } : coupon);
    hint.value = "优惠券已领取，可在“我的优惠券”中查看";
  } catch (error) {
    hint.value = error instanceof ApiError ? error.message : "优惠券领取失败";
  }
}

watch(view, () => void load(), { immediate: true });
</script>

<template>
  <div v-if="view === 'center'" class="coupon-center-page">
    <section class="coupon-center-hero">
      <div class="coupon-center-title"><span>领券中心</span><small>精选优惠，随心领取</small></div>
    </section>
    <main class="pc-container coupon-center-main">
      <nav class="center-filters" aria-label="优惠券分类">
        <button type="button" :class="{ active: centerKind === 'all' }" @click="centerKind = 'all'">全部</button>
        <button type="button" :class="{ active: centerKind === 'platform' }" @click="centerKind = 'platform'">通用券</button>
        <button type="button" :class="{ active: centerKind === 'store' }" @click="centerKind = 'store'">店铺券</button>
        <RouterLink to="/coupons">我的优惠券</RouterLink>
      </nav>
      <p v-if="hint" class="coupon-hint">{{ hint }}</p>
      <p v-else-if="loading" class="coupon-empty">正在加载优惠券…</p>
      <section v-else class="center-coupon-grid">
        <article v-for="coupon in visibleCenter" :key="coupon.coupon_id" class="coupon-card">
          <div class="coupon-main">
            <span class="coupon-kind">{{ kindLabel(coupon) }}</span>
            <strong class="coupon-value">{{ discountLabel(coupon) }}</strong>
            <div class="coupon-copy"><b>{{ coupon.title }}</b><span>满{{ Number(coupon.use_min_price).toFixed(0) }}可用</span><small>{{ coupon.type ? "店铺专享" : "全场通用" }}</small></div>
          </div>
          <button class="coupon-side" type="button" :disabled="coupon.received" @click="onReceive(coupon.coupon_id)">{{ coupon.received ? "已领取" : "立即领取" }}</button>
        </article>
        <p v-if="!visibleCenter.length" class="coupon-empty">暂无可领取的优惠券</p>
      </section>
    </main>
  </div>
  <AccountFrame v-else>
    <template #crumb><span>›</span> 我的优惠券</template>
    <div v-if="!user.isLogin" class="coupon-guest">
      <p>{{ hint }}</p>
      <button class="pc-btn" type="button" @click="router.push({ name: 'login', query: { redirect: route.fullPath } })">去登录</button>
    </div>
    <template v-else>
      <nav class="coupon-tabs" aria-label="优惠券分类">
        <button type="button" :class="{ active: view === 'mine' }" @click="setView('mine')">未使用</button>
        <button type="button" :class="{ active: view === 'history' }" @click="setView('history')">已使用/已过期</button>
        <button type="button" @click="setView('center')">前往领券中心</button>
      </nav>
      <p v-if="hint" class="coupon-hint">{{ hint }}</p>
      <p v-else-if="loading" class="coupon-empty">正在加载优惠券…</p>
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
.coupon-center-page { min-height: 680px; padding-bottom: 64px; background: #f7f7f7; }
.coupon-center-hero { height: 254px; display: grid; place-items: center; background: #fbf5ed url("/demo/coupon-center-hero-v1.png") center / cover no-repeat; }
.coupon-center-title { display: grid; justify-items: center; gap: 8px; color: #2f2c2b; text-align: center; text-shadow: 0 1px 0 rgb(255 255 255 / 72%); }.coupon-center-title span { font-size: 46px; font-weight: 900; letter-spacing: .04em; }.coupon-center-title small { color: #9c8271; font-size: 13px; letter-spacing: .16em; }
.coupon-center-main { margin-top: -42px; }.center-filters { display: flex; align-items: center; gap: 22px; min-height: 92px; padding: 0 36px; background: #fff; box-shadow: 0 8px 22px rgb(119 83 47 / 8%); }.center-filters button, .center-filters a { min-width: 88px; height: 34px; border: 1px solid #ddd; border-radius: 18px; color: #666; background: #fff; font-size: 14px; text-align: center; line-height: 32px; }.center-filters button.active, .center-filters button:hover { border-color: #f13728; color: #fff; background: #f13728; }.center-filters a { margin-left: auto; border-color: transparent; color: #999; }.center-filters a:hover { color: #f13728; }.center-coupon-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 460px)); align-content: start; gap: 24px; min-height: 380px; padding: 36px 0; }
.coupon-tabs { display: flex; gap: 70px; height: 74px; border-bottom: 1px solid #eee; }
.coupon-tabs button { position: relative; border: 0; padding: 0 16px; background: transparent; color: #999; font-size: 20px; cursor: pointer; }
.coupon-tabs button.active { color: #f13728; }
.coupon-tabs button.active::after { position: absolute; bottom: -1px; left: 50%; width: 140px; height: 3px; background: #f13728; content: ""; transform: translateX(-50%); }
.coupon-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 20px; padding-top: 38px; }
.coupon-card { display: grid; grid-template-columns: minmax(0, 1fr) 72px; min-height: 194px; overflow: hidden; background: #fff; box-shadow: 0 9px 26px rgb(0 0 0 / 8%); }
.coupon-main { position: relative; display: grid; grid-template-columns: 122px minmax(0, 1fr); align-content: center; gap: 9px 12px; padding: 24px 22px 20px 28px; border-right: 1px dashed #f1cbc6; }
.coupon-kind { position: absolute; top: 14px; left: 28px; padding: 4px 10px; background: #fff0ef; color: #f13728; font-size: 13px; }
.coupon-value { align-self: center; margin-top: 16px; color: #f13728; font-size: 44px; font-weight: 500; letter-spacing: -2px; }
.coupon-copy { display: grid; align-content: center; gap: 8px; margin-top: 16px; color: #333; }.coupon-copy b { overflow: hidden; font-size: 17px; font-weight: 500; text-overflow: ellipsis; white-space: nowrap; }.coupon-copy span, .coupon-copy small, .coupon-main small { color: #999; font-size: 13px; }.coupon-copy small { color: #c2a28e; }.coupon-main small { grid-column: 1 / -1; white-space: nowrap; }
.coupon-side { display: grid; width: 72px; place-items: center; border: 0; background: #f13728; color: #fff; font-size: 17px; line-height: 1.55; cursor: pointer; writing-mode: vertical-rl; }.coupon-side:disabled { cursor: default; opacity: .72; }.coupon-card.invalid { filter: grayscale(1); opacity: .67; }.coupon-card.invalid .coupon-side { background: #999; }
.coupon-empty { grid-column: 1 / -1; padding: 92px 0; color: #aaa; text-align: center; }.coupon-hint { margin: 24px 0 0; color: #d9362b; }.coupon-guest { display: grid; min-height: 380px; place-content: center; justify-items: center; gap: 14px; color: #777; }
@media (max-width: 960px) { .coupon-grid, .center-coupon-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }.coupon-tabs { gap: 22px; }.coupon-tabs button { font-size: 16px; } } @media (max-width: 620px) { .coupon-center-hero { height: 190px; }.coupon-center-title span { font-size: 30px; }.center-filters { gap: 8px; padding: 0 12px; }.center-filters button { min-width: 66px; }.coupon-grid, .center-coupon-grid { grid-template-columns: 1fr; }.coupon-tabs { gap: 0; justify-content: space-between; }.coupon-tabs button { padding: 0 7px; font-size: 14px; }.coupon-tabs button.active::after { width: 70px; } }
</style>
