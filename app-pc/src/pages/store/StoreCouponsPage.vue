<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchCouponCenter, receiveCoupon, type Coupon } from "@/api/coupon";
import { fetchStoreHome } from "@/api/catalog";
import { ApiError } from "@/utils/request";
import { useUserStore } from "@/stores/user";

const route = useRoute();
const router = useRouter();
const user = useUserStore();
const merId = computed(() => Number(route.params.id));
const coupons = ref<Coupon[]>([]);
const storeName = ref("店铺");
const hint = ref("");

async function load() {
  if (!merId.value) return;
  try {
    const [couponRes, storeRes] = await Promise.all([
      fetchCouponCenter(merId.value),
      fetchStoreHome(merId.value),
    ]);
    coupons.value = couponRes.list || [];
    storeName.value = storeRes.mer_name || "店铺";
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "店铺优惠券加载失败";
  }
}

async function receive(coupon: Coupon) {
  if (!user.isLogin) {
    router.push({ name: "login", query: { redirect: route.fullPath } });
    return;
  }
  try {
    await receiveCoupon(coupon.coupon_id);
    await load();
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "领取失败";
  }
}

onMounted(() => void load());

function couponAmount(coupon: Coupon) {
  return coupon.discount_type === "rate" ? `${Number(coupon.discount_value / 10).toFixed(1)}折` : `¥${Number(coupon.discount_value).toFixed(0)}`;
}
</script>

<template>
  <div class="store-coupon-page">
    <nav class="store-nav">
      <div class="pc-container store-nav__inner">
        <RouterLink :to="`/store/${merId}`">店铺首页</RouterLink>
        <RouterLink :to="`/store/${merId}`">全部分类</RouterLink>
        <RouterLink :to="`/store/${merId}/coupons`" class="active">领优惠券</RouterLink>
        <div class="store-search"><input placeholder="店内商品搜索" /><button type="button">搜索</button></div>
      </div>
    </nav>
    <div class="pc-container coupon-layout">
      <aside class="store-profile">
        <img src="/brand/qixi-logo.png" :alt="storeName" />
        <span class="badge">旗舰店</span>
        <h1>{{ storeName }}</h1>
        <p>店铺评分 <b>★★★★★</b></p>
        <p>关注人数 <strong>热度持续上升</strong></p>
        <RouterLink :to="`/store/${merId}`">进入店铺</RouterLink>
      </aside>
      <main class="coupon-main">
        <header><h1>店铺优惠券</h1><span>优惠券仅限本店商品使用</span></header>
        <p v-if="hint" class="hint">{{ hint }}</p>
        <div v-else-if="!coupons.length" class="empty">该店铺暂未发放优惠券</div>
        <div v-else class="coupon-grid">
          <article v-for="coupon in coupons" :key="coupon.coupon_id" class="coupon-card" :class="{ received: coupon.received }">
            <div class="amount" :class="{ rate: coupon.discount_type === 'rate' }">{{ couponAmount(coupon) }}</div>
            <div class="coupon-detail"><b>{{ coupon.title }}</b><span>满{{ Number(coupon.use_min_price).toFixed(0) }}元可用</span><small>剩余 {{ coupon.remain_count }} 张</small></div>
            <button type="button" :disabled="coupon.received || coupon.remain_count <= 0" @click="receive(coupon)">{{ coupon.received ? "已领取" : coupon.remain_count <= 0 ? "已领完" : "立即领取" }}</button>
          </article>
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
.store-coupon-page { min-height: 680px; padding-bottom: 76px; background: #f7f7f7; }.store-nav { background: #e4e4e4; }.store-nav__inner { display: flex; align-items: center; gap: 46px; height: 52px; }.store-nav a { color: #333; font-size: 15px; }.store-nav a.active { padding: 7px 16px; border-radius: 18px; color: #fff; background: #282828; }.store-search { display: flex; width: 238px; margin-left: auto; }.store-search input { min-width: 0; flex: 1; height: 28px; padding: 0 10px; border: 0; background: #fff; }.store-search button { width: 46px; border: 0; color: #fff; background: #282828; }.coupon-layout { display: grid; grid-template-columns: 220px minmax(0, 1fr); gap: 20px; padding-top: 26px; }.store-profile, .coupon-main { background: #fff; }.store-profile { padding: 35px 26px; text-align: center; }.store-profile img { width: 66px; height: 66px; border-radius: 50%; object-fit: cover; }.badge { display: inline-block; margin: 14px 0 4px; padding: 3px 6px; color: #fff; background: #f13728; font-size: 12px; }.store-profile h1 { margin: 0 0 18px; font-size: 16px; }.store-profile p { display: flex; justify-content: space-between; margin: 0; padding: 12px 0; border-top: 1px solid #eee; color: #888; font-size: 13px; }.store-profile p b { color: #f13728; letter-spacing: 2px; }.store-profile p strong { color: #555; font-weight: 400; }.store-profile > a { display: block; margin-top: 18px; padding: 10px; color: #fff; background: #f13728; }.coupon-main { min-height: 500px; padding: 20px 24px; }.coupon-main > header { display: flex; align-items: baseline; gap: 18px; padding-bottom: 20px; border-bottom: 1px solid #eee; }.coupon-main h1 { margin: 0; font-size: 20px; }.coupon-main header span { color: #aaa; font-size: 13px; }.coupon-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 24px; padding-top: 30px; }.coupon-card { display: grid; grid-template-columns: auto minmax(0, 1fr) 66px; min-height: 142px; overflow: hidden; border: 1px solid #ffd6d2; background: linear-gradient(135deg, #fff, #fff8f7); }.coupon-card.received { filter: grayscale(1); opacity: .58; }.amount { align-self: center; padding-left: 18px; color: #f13728; font-size: 36px; font-weight: 700; }.amount.rate { font-size: 28px; }.coupon-detail { display: grid; align-content: center; gap: 7px; padding: 12px; }.coupon-detail b { overflow: hidden; color: #444; text-overflow: ellipsis; white-space: nowrap; }.coupon-detail span, .coupon-detail small { color: #999; font-size: 12px; }.coupon-card button { border: 0; color: #fff; background: #f13728; font-size: 16px; line-height: 1.45; cursor: pointer; writing-mode: vertical-rl; }.coupon-card button:disabled { cursor: default; background: #aaa; }.hint, .empty { padding:48px 0;color:#d9362b;text-align:center; }.empty { color:#aaa; }@media (max-width:1000px){.coupon-grid{grid-template-columns:repeat(2,minmax(0,1fr));}.coupon-layout{grid-template-columns:1fr;}.store-profile{display:none;}}@media(max-width:620px){.coupon-grid{grid-template-columns:1fr;}.store-nav__inner{gap:22px;overflow-x:auto;}.store-search{display:none;}}
</style>
