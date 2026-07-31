<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  fetchCategories,
  fetchProducts,
  type CategoryItem,
  type ProductItem,
  type ProductSort,
} from "@/api/catalog";
import ProductCard from "@/components/ProductCard.vue";

const route = useRoute();
const router = useRouter();
const list = ref<ProductItem[]>([]);
const categories = ref<CategoryItem[]>([]);
const hint = ref("");
const loading = ref(false);
const total = ref(0);
const limit = 20;

const keyword = computed(() =>
  typeof route.query.keyword === "string" ? route.query.keyword : "",
);
const cateId = computed(() => {
  const value = route.query.cate_id;
  const id = typeof value === "string" ? Number(value) : Number.NaN;
  return Number.isSafeInteger(id) && id > 0 ? id : undefined;
});
const page = computed(() => {
  const value = Number(route.query.page || 1);
  return Number.isSafeInteger(value) && value > 0 ? value : 1;
});
const sort = computed<ProductSort>(() => {
  const value = route.query.sort;
  return value === "sales" || value === "price" ? value : "default";
});
const order = computed<"asc" | "desc">(() =>
  route.query.order === "asc" ? "asc" : "desc",
);
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit)));
const pageNumbers = computed(() => {
  const start = Math.max(1, Math.min(page.value - 2, totalPages.value - 4));
  const end = Math.min(totalPages.value, start + 4);
  return Array.from({ length: end - start + 1 }, (_, index) => start + index);
});
const selectedCategory = computed(() =>
  categories.value.find((item) => item.id === cateId.value),
);
const pageTitle = computed(() => {
  if (keyword.value) return `“${keyword.value}” 的搜索结果`;
  return selectedCategory.value?.name || "全部商品";
});

async function load() {
  loading.value = true;
  hint.value = "";
  try {
    const [data, cateData] = await Promise.all([
      fetchProducts({
        keyword: keyword.value || undefined,
        cate_id: cateId.value,
        page: page.value,
        limit,
        sort: sort.value,
        order: order.value,
      }),
      categories.value.length ? Promise.resolve(categories.value) : fetchCategories(),
    ]);
    list.value = data.list || [];
    total.value = data.total || 0;
    categories.value = cateData || [];
    hint.value = list.value.length ? "" : "暂无符合条件的在售商品";
  } catch (error) {
    list.value = [];
    total.value = 0;
    hint.value = (error as Error).message || "商品列表加载失败";
  } finally {
    loading.value = false;
  }
}

function updateQuery(next: Record<string, string | undefined>) {
  const query: Record<string, string> = {};
  for (const [key, value] of Object.entries({ ...route.query, ...next })) {
    if (typeof value === "string" && value) query[key] = value;
  }
  void router.push({ name: "goods-list", query });
}

function selectCategory(id?: number) {
  updateQuery({ cate_id: id ? String(id) : undefined, page: undefined });
}

function selectSort(nextSort: ProductSort) {
  if (nextSort === "price" && sort.value === "price") {
    updateQuery({ order: order.value === "asc" ? "desc" : "asc", page: undefined });
    return;
  }
  updateQuery({
    sort: nextSort === "default" ? undefined : nextSort,
    order: nextSort === "price" ? "asc" : undefined,
    page: undefined,
  });
}

function goPage(nextPage: number) {
  if (nextPage < 1 || nextPage > totalPages.value || nextPage === page.value) return;
  updateQuery({ page: nextPage === 1 ? undefined : String(nextPage) });
}

watch(
  () => [route.query.keyword, route.query.cate_id, route.query.page, route.query.sort, route.query.order],
  () => void load(),
  { immediate: true },
);
</script>

<template>
  <div class="goods-page">
    <div class="pc-container">
      <nav class="crumb" aria-label="当前位置">
        <RouterLink to="/">首页</RouterLink><span>›</span><span>商品列表</span>
        <template v-if="keyword"><span>›</span><span>搜索：{{ keyword }}</span></template>
      </nav>

      <section class="filter-panel" aria-label="商品筛选">
        <div class="filter-row">
          <b>商品分类</b>
          <div class="filter-options">
            <button type="button" :class="{ active: !cateId }" @click="selectCategory()">全部</button>
            <button
              v-for="category in categories"
              :key="category.id"
              type="button"
              :class="{ active: cateId === category.id }"
              @click="selectCategory(category.id)"
            >
              {{ category.name }}
            </button>
            <span v-if="!categories.length && !loading" class="no-category">当前暂无可用分类</span>
          </div>
        </div>
        <div class="filter-row">
          <b>排序方式</b>
          <div class="filter-options sort-options">
            <button type="button" :class="{ active: sort === 'default' }" @click="selectSort('default')">综合排序</button>
            <button type="button" :class="{ active: sort === 'sales' }" @click="selectSort('sales')">销量优先</button>
            <button type="button" :class="{ active: sort === 'price' }" @click="selectSort('price')">
              价格 <span v-if="sort === 'price'">{{ order === 'asc' ? '↑' : '↓' }}</span>
            </button>
          </div>
        </div>
      </section>

      <section class="result-head">
        <div>
          <h1>{{ pageTitle }}</h1>
          <p v-if="!loading">共 <em>{{ total }}</em> 件商品</p>
        </div>
        <div v-if="totalPages > 1" class="head-pages">
          <span>{{ page }} / {{ totalPages }}</span>
          <button type="button" :disabled="page <= 1" @click="goPage(page - 1)">‹</button>
          <button type="button" :disabled="page >= totalPages" @click="goPage(page + 1)">›</button>
        </div>
      </section>

      <p v-if="loading" class="hint">商品加载中…</p>
      <p v-else-if="hint" class="hint empty">{{ hint }}</p>
      <div v-else class="grid">
        <ProductCard v-for="product in list" :key="product.id" :product="product" />
      </div>

      <nav v-if="totalPages > 1" class="pagination" aria-label="商品分页">
        <button type="button" :disabled="page <= 1" @click="goPage(page - 1)">上一页</button>
        <button
          v-for="number in pageNumbers"
          :key="number"
          type="button"
          :class="{ active: number === page }"
          @click="goPage(number)"
        >{{ number }}</button>
        <button type="button" :disabled="page >= totalPages" @click="goPage(page + 1)">下一页</button>
      </nav>
    </div>
  </div>
</template>

<style scoped>
.goods-page { min-height: 560px; padding-top: 18px; background: #f6f6f6; }
.crumb { display: flex; gap: .55rem; align-items: center; height: 42px; color: #888; font-size: .88rem; }
.crumb a:hover { color: #ef3727; }
.filter-panel { background: #fff; border: 1px solid #ededed; }
.filter-row { display: grid; grid-template-columns: 116px 1fr; min-height: 56px; border-bottom: 1px solid #f0f0f0; }
.filter-row:last-child { border-bottom: 0; }
.filter-row > b { display: flex; align-items: center; padding-left: 22px; background: #fafafa; color: #555; font-size: .9rem; font-weight: 500; }
.filter-options { display: flex; flex-wrap: wrap; align-items: center; gap: .35rem 1.25rem; padding: .55rem 1.15rem; }
.filter-options button { border: 0; padding: .38rem .45rem; color: #555; background: transparent; font-size: .9rem; cursor: pointer; }
.filter-options button:hover, .filter-options button.active { color: #ef3727; }
.no-category { color: #999; font-size: .88rem; }
.sort-options { gap: .35rem; }
.sort-options button { min-width: 92px; border: 1px solid #e4e4e4; }
.sort-options button.active { border-color: #ef3727; background: #fff7f6; }
.result-head { display: flex; align-items: end; justify-content: space-between; margin: 22px 0 13px; }
.result-head h1 { margin: 0; color: #333; font-size: 1.24rem; font-weight: 600; }
.result-head p { margin: .4rem 0 0; color: #888; font-size: .88rem; }
.result-head em { color: #ef3727; font-style: normal; }
.head-pages { display: flex; align-items: center; gap: .5rem; color: #999; font-size: .84rem; }
.head-pages button { width: 28px; height: 25px; border: 1px solid #ddd; background: #fff; color: #666; font-size: 1.12rem; line-height: 1; }
.head-pages button:disabled { color: #ccc; cursor: not-allowed; }
.hint { margin: 40px 0; color: var(--pc-muted); text-align: center; }
.hint.empty { padding: 3rem 0; background: #fff; border: 1px solid #eee; }
.grid { display: grid; grid-template-columns: repeat(5, minmax(0, 1fr)); gap: 14px; }
.pagination { display: flex; justify-content: center; gap: .5rem; margin: 32px 0 10px; }
.pagination button { min-width: 36px; height: 34px; padding: 0 .7rem; border: 1px solid #ddd; background: #fff; color: #555; }
.pagination button:hover:not(:disabled), .pagination button.active { border-color: #ef3727; color: #fff; background: #ef3727; }
.pagination button:disabled { color: #bbb; cursor: not-allowed; }
@media (max-width: 980px) { .grid { grid-template-columns: repeat(3, minmax(0, 1fr)); } }
@media (max-width: 680px) { .filter-row { grid-template-columns: 84px 1fr; } .filter-row > b { padding-left: 12px; } .filter-options { gap: .2rem .45rem; padding-inline: .55rem; } .grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
</style>
