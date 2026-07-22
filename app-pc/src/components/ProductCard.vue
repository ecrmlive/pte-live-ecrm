<script setup lang="ts">
import type { ProductItem } from "@/api/catalog";

defineProps<{
  product: ProductItem;
}>();
</script>

<template>
  <RouterLink class="card" :to="`/goods/${product.id}`">
    <div class="cover">
      <span>{{ product.store_name.slice(0, 1) }}</span>
    </div>
    <div class="meta">
      <h3>{{ product.store_name }}</h3>
      <p class="mer">{{ product.mer_name }}</p>
      <div class="price-row">
        <strong>¥{{ product.price }}</strong>
        <span v-if="product.ot_price" class="ot">¥{{ product.ot_price }}</span>
      </div>
    </div>
  </RouterLink>
</template>

<style scoped>
.card {
  display: block;
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: var(--pc-radius);
  overflow: hidden;
  box-shadow: var(--pc-shadow);
  transition: transform 0.18s ease, border-color 0.18s ease;
}

.card:hover {
  transform: translateY(-3px);
  border-color: rgba(15, 107, 92, 0.35);
}

.cover {
  aspect-ratio: 1;
  display: grid;
  place-items: center;
  background: linear-gradient(160deg, #dceeea, #f7faf9);
  font-size: 2rem;
  font-weight: 700;
  color: var(--pc-brand);
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

.mer {
  margin: 0.35rem 0 0.55rem;
  color: var(--pc-muted);
  font-size: 0.86rem;
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
