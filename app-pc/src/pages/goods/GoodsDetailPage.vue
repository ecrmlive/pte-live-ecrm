<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { activateMerchantContext } from "@/api/auth";
import { fetchProductDetail, type ProductDetail } from "@/api/catalog";
import { addCart } from "@/api/trade";
import { fetchProductFavoriteState, removeProductFavorite, saveProductFavorite } from "@/api/favorite";
import { useUserStore } from "@/stores/user";
import { ApiError } from "@/utils/request";

const route = useRoute();
const router = useRouter();
const user = useUserStore();
const detail = ref<ProductDetail | null>(null);
const hint = ref("");
const adding = ref(false);
const followed = ref(false);
const following = ref(false);

const id = computed(() => Number(route.params.id));

async function load() {
  if (!id.value) return;
  try {
    detail.value = await fetchProductDetail(id.value);
    followed.value = false;
    hint.value = "";
    if (user.isLogin && detail.value.merchant_app_id) {
      try {
        const context = await activateMerchantContext(detail.value.merchant_app_id);
        user.replaceTokenPair(context.token);
      } catch {
        hint.value = "商品已加载；登录上下文初始化失败，请重新登录后再购买";
      }
    }
    if (user.isLogin) {
      const state = await fetchProductFavoriteState(detail.value.id);
      followed.value = state.followed;
    }
  } catch (e) {
    detail.value = null;
    hint.value = (e as Error).message || "商品详情加载失败";
  }
}

async function toggleFavorite() {
  if (!detail.value || following.value) return;
  if (!user.isLogin) {
    await router.push({ name: "login", query: { redirect: route.fullPath } });
    return;
  }
  following.value = true;
  try {
    const state = followed.value ? await removeProductFavorite(detail.value.id) : await saveProductFavorite(detail.value.id);
    followed.value = state.followed;
    hint.value = state.followed ? "商品已收藏" : "已取消商品收藏";
  } catch (error) {
    hint.value = error instanceof ApiError ? error.message : "收藏操作失败";
  } finally {
    following.value = false;
  }
}

watch(id, () => void load(), { immediate: true });

async function onAddCart() {
  if (!detail.value) return;
  if (!user.isLogin) {
    router.push({ name: "login", query: { redirect: route.fullPath } });
    return;
  }
  adding.value = true;
  try {
    await addCart({ product_id: detail.value.id, cart_num: 1 });
    hint.value = "已加入购物车";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加购失败";
  } finally {
    adding.value = false;
  }
}
</script>

<template>
  <div class="pc-container" v-if="detail">
    <p class="crumb"><RouterLink to="/">首页</RouterLink> / <RouterLink to="/goods">全部商品</RouterLink> / 商品详情</p>
    <section class="panel">
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div class="layout">
        <div class="gallery">
          <div class="cover">
            <img v-if="detail.image" :src="detail.image" :alt="detail.store_name" />
            <span v-else>暂无商品图片</span>
          </div>
        </div>
        <div class="info">
          <p class="mer">
            <RouterLink :to="`/store/${detail.mer_id}`">{{ detail.mer_name }}</RouterLink>
          </p>
          <h1>{{ detail.title || detail.store_name }}</h1>
          <p class="desc">{{ detail.store_info }}</p>
          <div class="price-box">
            <span>售价</span><div class="price"><strong>¥{{ detail.price }}</strong><span v-if="detail.ot_price">¥{{ detail.ot_price }}</span></div>
          </div>
          <ul>
            <li>销量 {{ detail.sales }} · 库存 {{ detail.stock }} {{ detail.unit_name }}</li>
            <li>配送：{{ detail.delivery_way }}</li>
            <li>规格：{{ detail.spec_type === 0 ? "单规格" : "多规格" }}</li>
          </ul>
          <div class="actions">
            <button class="pc-btn" type="button" :disabled="adding" @click="onAddCart">加入购物车</button>
            <button class="pc-btn ghost" type="button" :disabled="following" @click="toggleFavorite">{{ following ? "处理中…" : followed ? "已收藏" : "收藏商品" }}</button>
            <RouterLink class="pc-btn ghost" to="/cart">去购物车</RouterLink>
          </div>
        </div>
      </div>
    </section>
  </div>
  <div v-else class="pc-container"><p class="hint">{{ hint || "商品不存在" }}</p></div>
</template>

<style scoped>
.crumb { margin: 1.35rem 0 .8rem; font-size: .85rem; color: #888; }
.crumb a { color: #555; }
.panel {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: 0;
  padding: 1.4rem;
  box-shadow: var(--pc-shadow);
}

.hint {
  margin: 0 0 1rem;
  color: var(--pc-muted);
}

.layout {
  display: grid;
  grid-template-columns: 420px 1fr;
  gap: 1.5rem;
}

.cover {
  aspect-ratio: 1;
  display: grid;
  place-items: center;
  color: #999; background: #f7f7f7; font-size: .9rem;
}
.cover img { width: 100%; height: 100%; object-fit: cover; }

.mer {
  margin: 0;
  color: var(--pc-brand);
  font-weight: 600;
}

h1 {
  margin: 0.4rem 0 0.7rem;
  font-size: 1.45rem;
}

.desc {
  margin: 0;
  color: var(--pc-muted);
  line-height: 1.7;
}

.price-box { margin: 1rem 0; padding: 1rem 1.2rem; background: #f7f7f7; display: flex; align-items: baseline; gap: 1.2rem; color: #777; }
.price {
  display: flex;
  gap: 0.8rem;
  align-items: baseline;
}

.price strong {
  color: var(--pc-accent);
  font-size: 1.8rem;
}

.price span {
  color: var(--pc-muted);
  text-decoration: line-through;
}

ul {
  margin: 0 0 1.2rem;
  padding-left: 1.1rem;
  color: var(--pc-muted);
  line-height: 1.8;
}

.actions {
  display: flex;
  gap: 0.7rem;
}
.actions .pc-btn { border-radius: 0; min-width: 155px; }

@media (max-width: 900px) {
  .layout {
    grid-template-columns: 1fr;
  }
}
</style>
