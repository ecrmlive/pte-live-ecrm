<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { fetchProducts, type ProductItem } from "@/api/catalog";
import ProductCard from "@/components/ProductCard.vue";

const route = useRoute();
const list = ref<ProductItem[]>([]);
const hint = ref("");
const loading = ref(false);

const keyword = computed(() =>
  typeof route.query.keyword === "string" ? route.query.keyword : "",
);
const cateId = computed(() => {
  const v = route.query.cate_id;
  return typeof v === "string" && v ? Number(v) : undefined;
});

const demo: ProductItem[] = [
  {
    id: 1001,
    mer_id: 1,
    mer_name: "栖息优选店",
    store_name: "有机燕麦礼盒",
    image: "",
    price: "89.00",
    ot_price: "129.00",
    sales: 320,
    stock: 99,
  },
  {
    id: 1002,
    mer_id: 2,
    mer_name: "山海生鲜",
    store_name: "冷冻海鲜组合装",
    image: "",
    price: "168.00",
    sales: 180,
    stock: 40,
  },
  {
    id: 1003,
    mer_id: 3,
    mer_name: "都市数码",
    store_name: "蓝牙降噪耳机",
    image: "",
    price: "299.00",
    ot_price: "399.00",
    sales: 96,
    stock: 30,
  },
];

async function load() {
  loading.value = true;
  hint.value = "";
  try {
    const data = await fetchProducts({
      keyword: keyword.value || undefined,
      cate_id: cateId.value,
      page: 1,
    });
    list.value = data.list?.length ? data.list : demo;
    hint.value = data.list?.length
      ? `共 ${data.total} 件`
      : "暂无接口数据 · 展示演示列表（阶段 2）";
  } catch {
    list.value = demo;
    hint.value = "商品列表接口暂不可用 · 展示演示列表（阶段 2）";
  } finally {
    loading.value = false;
  }
}

watch(
  () => [route.query.keyword, route.query.cate_id],
  () => {
    void load();
  },
  { immediate: true },
);
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <div class="head">
        <h1>商品列表</h1>
        <p>
          筛选：综合 / 销量 / 价格 / 新品 / 品牌 / 价格区间（阶段 2 接真）。
          <template v-if="keyword">当前关键词：{{ keyword }}</template>
        </p>
      </div>
      <p class="hint">{{ loading ? "加载中…" : hint }}</p>
      <div class="grid">
        <ProductCard v-for="p in list" :key="p.id" :product="p" />
      </div>
    </section>
  </div>
</template>

<style scoped>
.panel {
  background: transparent;
}

.head {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: var(--pc-radius);
  padding: 1.2rem 1.4rem;
  margin-bottom: 1rem;
  box-shadow: var(--pc-shadow);
}

h1 {
  margin: 0;
  font-size: 1.4rem;
}

.head p,
.hint {
  color: var(--pc-muted);
}

.head p {
  margin: 0.4rem 0 0;
}

.hint {
  margin: 0 0 0.9rem;
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
