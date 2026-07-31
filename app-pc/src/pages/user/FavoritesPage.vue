<script setup lang="ts">
import { ref, watch } from "vue";
import AccountFrame from "@/components/AccountFrame.vue";
import ProductCard from "@/components/ProductCard.vue";
import type { ProductItem } from "@/api/catalog";
import {
  fetchProductFavorites,
  fetchStoreFavorites,
  removeProductFavorite,
  removeStoreFavorite,
  type FavoriteStore,
} from "@/api/favorite";
import { ApiError } from "@/utils/request";

const tab = ref<"product" | "store">("product");
const products = ref<ProductItem[]>([]);
const stores = ref<FavoriteStore[]>([]);
const loading = ref(false);
const hint = ref("");

async function load() {
  loading.value = true;
  try {
    if (tab.value === "product") products.value = (await fetchProductFavorites()).list || [];
    else stores.value = (await fetchStoreFavorites()).list || [];
    hint.value = "";
  } catch (error) {
    hint.value = error instanceof ApiError ? error.message : "收藏列表加载失败";
  } finally {
    loading.value = false;
  }
}

async function removeProduct(product: ProductItem) {
  try {
    await removeProductFavorite(product.id);
    products.value = products.value.filter((item) => item.id !== product.id);
  } catch (error) { hint.value = error instanceof ApiError ? error.message : "取消收藏失败"; }
}
async function removeStore(store: FavoriteStore) {
  try {
    await removeStoreFavorite(store.store_id);
    stores.value = stores.value.filter((item) => item.store_id !== store.store_id);
  } catch (error) { hint.value = error instanceof ApiError ? error.message : "取消收藏失败"; }
}

watch(tab, () => void load(), { immediate: true });
</script>

<template>
  <AccountFrame>
    <template #crumb><span>›</span> 我的收藏</template>
    <div class="tabs">
      <button :class="{ active: tab === 'product' }" type="button" @click="tab = 'product'">商品收藏</button>
      <button :class="{ active: tab === 'store' }" type="button" @click="tab = 'store'">店铺收藏</button>
    </div>
    <p v-if="hint" class="hint">{{ hint }}</p>
    <p v-else-if="loading" class="empty">正在加载收藏…</p>
    <p v-else-if="tab === 'product' && !products.length" class="empty">暂无商品收藏，去逛逛喜欢的商品吧。</p>
    <p v-else-if="tab === 'store' && !stores.length" class="empty">暂无店铺收藏，去发现优质店铺吧。</p>
    <section v-else-if="tab === 'product'" class="product-grid">
      <div v-for="product in products" :key="product.id" class="favorite-product">
        <ProductCard :product="product" />
        <button type="button" class="remove" @click="removeProduct(product)">取消收藏</button>
      </div>
    </section>
    <section v-else class="store-grid">
      <article v-for="store in stores" :key="store.store_id" class="store-card">
        <img src="/brand/qixi-logo.png" :alt="store.store_name" />
        <div><h2>{{ store.store_name }}</h2><p>{{ store.follower_count }}人关注</p></div>
        <div class="store-actions"><RouterLink :to="`/store/${store.mer_id}`">进店逛逛</RouterLink><button type="button" @click="removeStore(store)">取消关注</button></div>
      </article>
    </section>
  </AccountFrame>
</template>

<style scoped>
.tabs { display: flex; gap: 58px; height: 74px; border-bottom: 1px solid #eee; }.tabs button { position: relative; border: 0; padding: 0 16px; color: #999; background: transparent; font-size: 20px; cursor: pointer; }.tabs button.active { color: #f13728; }.tabs button.active::after { position: absolute; bottom: -1px; left: 50%; width: 140px; height: 3px; background: #f13728; content: ""; transform: translateX(-50%); }.hint { color: #d9362b; }.empty { padding: 76px 0; color: #aaa; text-align: center; }.product-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 18px; padding-top: 34px; }.favorite-product { position: relative; }.favorite-product :deep(.card) { height: 100%; }.remove { position: absolute; top: 8px; right: 8px; border: 1px solid rgb(255 255 255 / 85%); padding: 5px 8px; color: #777; background: rgb(255 255 255 / 92%); font-size: 12px; cursor: pointer; }.remove:hover { color: #f13728; border-color: #f13728; }.store-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 26px; padding-top: 34px; }.store-card { display: grid; min-height: 474px; grid-template-rows: 140px auto 1fr auto; justify-items: center; align-items: center; padding: 30px 28px; border: 1px solid #eee; text-align: center; }.store-card img { width: 116px; height: 116px; border-radius: 50%; object-fit: cover; }.store-card h2 { margin: 0 0 18px; color: #333; font-size: 20px; }.store-card p { margin: 0; color: #777; font-size: 15px; line-height: 1.7; }.store-actions { display: grid; width: 100%; grid-template-columns: 1fr 1fr; gap: 0; padding-top: 24px; border-top: 1px dashed #e8e8e8; }.store-actions a, .store-actions button { border: 0; padding: 0; color: #333; background: transparent; font-size: 15px; cursor: pointer; white-space: nowrap; }.store-actions a { border-right: 1px solid #eee; }.store-actions button:hover, .store-actions a:hover { color: #f13728; }@media (max-width: 800px) { .product-grid, .store-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }.tabs { gap: 18px; }.tabs button { font-size: 16px; }.tabs button.active::after { width: 90px; } } @media (max-width: 540px) { .product-grid, .store-grid { grid-template-columns: 1fr; } }
</style>
