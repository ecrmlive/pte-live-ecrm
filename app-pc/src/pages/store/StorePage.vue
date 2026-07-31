<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { activateMerchantContext } from "@/api/auth";
import { fetchCategories, fetchStoreHome, type CategoryItem, type ProductItem } from "@/api/catalog";
import ProductCard from "@/components/ProductCard.vue";
import { useUserStore } from "@/stores/user";

const route = useRoute();
const user = useUserStore();
const merId = computed(() => Number(route.params.id));
const merName = ref("");
const products = ref<ProductItem[]>([]);
const categories = ref<CategoryItem[]>([]);
const hint = ref("");
const selectedCategory = ref<number>();
const followed = ref(false);
const sort = ref<"default" | "sales" | "price">("default");

const rootCategories = computed(() => categories.value.filter((item) => item.pid === 0));
const storeProducts = computed(() => {
  const selected = selectedCategory.value;
  const categoryIds = selected
    ? new Set([selected, ...categories.value.filter((item) => item.pid === selected).map((item) => item.id)])
    : undefined;
  const filtered = categoryIds ? products.value.filter((item) => categoryIds.has(item.category_id || 0)) : products.value;
  return [...filtered].sort((a, b) => {
    if (sort.value === "sales") return b.sales - a.sales;
    if (sort.value === "price") return Number(a.price) - Number(b.price);
    return b.sales - a.sales;
  });
});

async function load() {
  if (!merId.value) return;
  hint.value = "加载中…";
  try {
    const [data, categoryList] = await Promise.all([fetchStoreHome(merId.value), fetchCategories()]);
    merName.value = data.mer_name || `店铺 #${merId.value}`;
    products.value = data.products || [];
    categories.value = categoryList || [];
    if (user.isLogin && data.merchant_app_id) {
      try {
        const context = await activateMerchantContext(data.merchant_app_id);
        user.replaceTokenPair(context.token);
      } catch {
        hint.value = "店铺已加载，重新登录后可购买";
        return;
      }
    }
    hint.value = products.value.length ? `共 ${data.total} 件在售商品` : "暂无在售商品";
  } catch (e) {
    merName.value = `店铺 #${merId.value}`;
    products.value = [];
    hint.value = (e as Error).message || "店铺加载失败";
  }
}

watch(merId, () => void load(), { immediate: true });
</script>

<template>
  <div class="store-page">
    <nav class="store-nav">
      <div class="pc-container store-nav__inner">
        <RouterLink :to="`/store/${merId}`">店铺首页</RouterLink>
        <button type="button" :class="{ active: !selectedCategory }" @click="selectedCategory = undefined">全部分类</button>
        <RouterLink to="/coupons">领优惠券</RouterLink>
        <div class="store-search"><input placeholder="店内商品搜索" /><button type="button">搜索</button></div>
      </div>
    </nav>

    <div class="pc-container store-content">
      <aside class="store-side">
        <section class="store-profile">
          <img class="store-avatar" src="/brand/qixi-logo.png" :alt="merName" />
          <div class="store-badge">旗舰店</div>
          <h1>{{ merName || '店铺加载中' }}</h1>
          <p>店铺评分 <strong>★★★★★</strong></p>
          <p>关注人数 <b>{{ products.reduce((sum, item) => sum + item.sales, 0) || 0 }}</b></p>
          <p>店铺资质 <span class="quality">已认证</span></p>
          <button type="button" class="follow-btn" :class="{ followed }" @click="followed = !followed">
            {{ followed ? '已收藏' : '收藏店铺' }}
          </button>
        </section>
        <section class="store-categories">
          <h2>店内分类</h2>
          <button type="button" :class="{ active: !selectedCategory }" @click="selectedCategory = undefined">全部商品</button>
          <button
            v-for="category in rootCategories"
            :key="category.id"
            type="button"
            :class="{ active: selectedCategory === category.id }"
            @click="selectedCategory = category.id"
          >{{ category.name }}</button>
        </section>
      </aside>

      <main class="store-main">
        <div class="sort-bar">
          <span>排序：</span>
          <button type="button" :class="{ active: sort === 'default' }" @click="sort = 'default'">默认</button>
          <button type="button" :class="{ active: sort === 'sales' }" @click="sort = 'sales'">销量</button>
          <button type="button" :class="{ active: sort === 'price' }" @click="sort = 'price'">价格</button>
          <span class="store-total">{{ hint }}</span>
        </div>
        <p v-if="!storeProducts.length && hint !== '加载中…'" class="empty">该分类暂无在售商品</p>
        <div v-else class="grid">
          <ProductCard v-for="product in storeProducts" :key="product.id" :product="product" />
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
.store-page { min-height: 680px; padding-bottom: 2.8rem; background: #f7f7f7; }
.store-nav { background: #e5e5e5; }.store-nav__inner { display: flex; align-items: center; gap: 2.35rem; height: 52px; }.store-nav a, .store-nav button { border: 0; color: #333; background: transparent; font-size: .94rem; cursor: pointer; }.store-nav button.active, .store-nav a.router-link-active { padding: .38rem .82rem; border-radius: 18px; color: #fff; background: #333; }.store-search { display: flex; width: 232px; margin-left: auto; }.store-search input { min-width: 0; height: 28px; padding: 0 .7rem; border: 0; background: #fff; font-size: .76rem; }.store-search button { width: 39px; color: #fff; background: #333; font-size: .74rem; }
.store-content { display: grid; grid-template-columns: 220px minmax(0, 1fr); gap: 20px; padding-top: 22px; }.store-side { display: grid; align-content: start; gap: 12px; }.store-profile, .store-categories, .store-main { background: #fff; }.store-profile { padding: 1.55rem 1.35rem 1.45rem; text-align: center; }.store-avatar { width: 62px; height: 62px; border-radius: 50%; object-fit: cover; }.store-badge { display: inline-block; margin: .85rem 0 .15rem; padding: .18rem .35rem; color: #fff; background: #ef3727; font-size: .72rem; }.store-profile h1 { margin: 0 0 1rem; color: #444; font-size: 1rem; font-weight: 500; }.store-profile p { display: flex; justify-content: space-between; margin: .72rem 0; padding-bottom: .65rem; border-bottom: 1px solid #eee; color: #888; font-size: .8rem; text-align: left; }.store-profile strong { color: #ef3727; letter-spacing: .12rem; }.quality { color: #d18a18; }.follow-btn { width: 100%; margin-top: .2rem; padding: .6rem; border: 0; color: #fff; background: #ef3727; cursor: pointer; }.follow-btn.followed { background: #999; }
.store-categories h2 { margin: 0; padding: .9rem 1.15rem; border-bottom: 1px solid #eee; color: #555; font-size: .9rem; }.store-categories button { display: block; width: 100%; padding: .82rem 1.15rem; border: 0; border-bottom: 1px solid #f1f1f1; color: #777; background: #fff; text-align: left; cursor: pointer; }.store-categories button:hover, .store-categories button.active { color: #ef3727; background: #fff8f7; }
.store-main { padding: 16px; }.sort-bar { display: flex; align-items: center; gap: 1.4rem; min-height: 48px; padding: 0 12px; border-bottom: 1px solid #eee; color: #777; font-size: .86rem; }.sort-bar button { border: 0; padding: .45rem 0; color: #444; background: transparent; cursor: pointer; }.sort-bar button.active { color: #ef3727; }.store-total { margin-left: auto; color: #999; font-size: .78rem; }.grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 16px; padding-top: 16px; }.grid :deep(.meta) { padding: .75rem .8rem .9rem; }.grid :deep(h3) { font-size: .9rem; }.grid :deep(.store) { display: none; }.grid :deep(.sales) { font-size: .72rem; }.empty { padding: 5rem 0; color: #999; text-align: center; }
@media (max-width: 980px) { .store-content { grid-template-columns: 1fr; }.store-side { grid-template-columns: 1fr 1fr; }.grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }.store-nav__inner { overflow-x: auto; } }
</style>
