<script setup lang="ts">
import { useRouter } from "vue-router";
import type { ProductItem } from "@/api/catalog";

const props = defineProps<{
  product: ProductItem;
}>();

const router = useRouter();

function goDetail() {
  void router.push({ name: "goods-detail", params: { id: props.product.id } });
}

function goStore(event: MouseEvent) {
  event.stopPropagation();
  if (!props.product.mer_id) return;
  void router.push({ name: "store", params: { id: props.product.mer_id } });
}
</script>

<template>
  <article
    class="card"
    tabindex="0"
    role="link"
    :aria-label="`查看商品：${product.title || product.store_name}`"
    @click="goDetail"
    @keydown.enter="goDetail"
    @keydown.space.prevent="goDetail"
  >
    <div class="cover">
      <img v-if="product.image" :src="product.image" :alt="product.title || product.store_name" />
      <div v-else class="cover-fallback" aria-hidden="true">
        <span>七禧精选</span><strong>{{ (product.title || product.store_name || '好物').slice(0, 4) }}</strong>
      </div>
      <span v-if="product.stock <= 0" class="sold-out">暂时缺货</span>
    </div>
    <div class="meta">
      <h3>{{ product.title || product.store_name }}</h3>
      <button v-if="product.mer_id" class="store" type="button" @click="goStore">
        {{ product.shop_name || product.mer_name || "查看店铺" }}
      </button>
      <p v-else class="store platform">平台自营</p>
      <div class="price-row">
        <strong>¥{{ product.price }}</strong>
        <span v-if="product.ot_price" class="ot">¥{{ product.ot_price }}</span>
        <span class="sales">已售 {{ product.sales }}</span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.card {
  display: block;
  background: var(--pc-surface);
  border: 1px solid #efefef;
  border-radius: 0;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.18s ease, border-color 0.18s ease, box-shadow 0.18s ease;
}

.card:hover {
  transform: translateY(-3px);
  border-color: #ef3727;
  box-shadow: 0 10px 24px rgb(0 0 0 / 8%);
}

.card:focus-visible {
  outline: 2px solid #ef3727;
  outline-offset: 3px;
}

.cover {
  position: relative;
  aspect-ratio: 1;
  display: grid;
  place-items: center;
  background: #f5f5f5; color: #999; font-size: .9rem;
}
.cover img { width: 100%; height: 100%; object-fit: cover; display: block; }
.cover-fallback { width: 100%; height: 100%; display: grid; align-content: center; gap: .45rem; padding: 1rem; color: #7c231c; background: linear-gradient(135deg, #fff3ef 0%, #f4c5b8 50%, #d78875 100%); }
.cover-fallback span { color: rgba(124, 35, 28, .72); font-size: .72rem; letter-spacing: .12em; }
.cover-fallback strong { color: #7c231c; font-size: 1.1rem; line-height: 1.35; word-break: break-all; }

.sold-out {
  position: absolute;
  right: .55rem;
  top: .55rem;
  padding: .2rem .42rem;
  color: #fff;
  background: rgb(0 0 0 / 56%);
  font-size: .75rem;
}

.meta {
  padding: 0.9rem 1rem 1.05rem;
}

h3 {
  margin: 0;
  font-size: 1rem;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.store {
  display: block;
  width: 100%;
  overflow: hidden;
  padding: 0;
  margin: .38rem 0 .58rem;
  border: 0;
  color: var(--pc-muted);
  background: transparent;
  font-size: .86rem;
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
}

.store:hover {
  color: #ef3727;
}

.platform {
  cursor: default;
}

.platform:hover {
  color: var(--pc-muted);
}

.sales {
  margin-left: auto;
  color: var(--pc-muted);
  font-size: .8rem;
}

.price-row {
  display: flex;
  align-items: baseline;
  gap: 0.45rem;
}

.price-row strong {
  color: var(--pc-accent);
  font-size: 1.1rem;
}

.ot {
  color: var(--pc-muted);
  text-decoration: line-through;
  font-size: 0.85rem;
}
</style>
