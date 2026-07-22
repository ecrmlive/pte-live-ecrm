<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { fetchStore, type ProductItem } from "@/api/catalog";
import ProductCard from "@/components/ProductCard.vue";

const route = useRoute();
const merId = computed(() => Number(route.params.id));
const merName = ref("");
const products = ref<ProductItem[]>([]);
const hint = ref("");

async function load() {
  if (!merId.value) return;
  hint.value = "加载中…";
  try {
    const data = await fetchStore(merId.value);
    merName.value = data.mer_name || `店铺 #${merId.value}`;
    products.value = data.products || [];
    hint.value = products.value.length ? `共 ${data.total} 件在售` : "暂无在售商品";
  } catch (e) {
    merName.value = `店铺 #${merId.value}`;
    products.value = [];
    hint.value = (e as Error).message || "店铺加载失败";
  }
}

watch(merId, () => void load(), { immediate: true });
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>{{ merName }}</h1>
      <p class="sub">店铺首页 · 店内商品（功能表 4）</p>
      <p class="hint">{{ hint }}</p>
      <div class="grid">
        <ProductCard v-for="p in products" :key="p.id" :product="p" />
      </div>
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

h1 {
  margin: 0;
}

.sub,
.hint {
  color: var(--pc-muted);
}

.sub {
  margin: 0.4rem 0 0.8rem;
}

.hint {
  margin: 0 0 1rem;
  font-size: 0.92rem;
}

.grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
}

@media (max-width: 980px) {
  .grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
