<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  fetchCategories,
  fetchProducts,
  fetchStores,
  type CategoryItem,
  type ProductItem,
  type ProductSort,
  type StoreDirectoryItem,
} from "@/api/catalog";
import ProductCard from "@/components/ProductCard.vue";

const route = useRoute();
const router = useRouter();
const list = ref<ProductItem[]>([]);
const categories = ref<CategoryItem[]>([]);
const stores = ref<StoreDirectoryItem[]>([]);
const hint = ref("");
const loading = ref(false);
const total = ref(0);
const limit = 20;

const keyword = computed(() => typeof route.query.keyword === "string" ? route.query.keyword : "");
const cateId = computed(() => toPositiveInteger(route.query.cate_id));
const merId = computed(() => toPositiveInteger(route.query.mer_id));
const page = computed(() => Math.max(1, Number(route.query.page || 1) || 1));
const sort = computed<ProductSort>(() => {
  const value = route.query.sort;
  return value === "sales" || value === "price" ? value : "default";
});
const order = computed<"asc" | "desc">(() => route.query.order === "asc" ? "asc" : "desc");
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / limit)));
const selectedCategory = computed(() => categories.value.find((item) => item.id === cateId.value));
const selectedStore = computed(() => stores.value.find((item) => item.mer_id === merId.value));
const pageTitle = computed(() => {
  if (keyword.value) return `“${keyword.value}” 的搜索结果`;
  if (selectedStore.value) return `${selectedStore.value.name} 商品`;
  return selectedCategory.value?.name || "全部商品";
});

function toPositiveInteger(value: unknown) {
  const id = typeof value === "string" ? Number(value) : Number.NaN;
  return Number.isSafeInteger(id) && id > 0 ? id : undefined;
}

async function load() {
  loading.value = true;
  hint.value = "";
  try {
    const [data, cateData, storeData] = await Promise.all([
      fetchProducts({
        keyword: keyword.value || undefined,
        cate_id: cateId.value,
        mer_id: merId.value,
        page: page.value,
        limit,
        sort: sort.value,
        order: order.value,
      }),
      categories.value.length ? Promise.resolve(categories.value) : fetchCategories(),
      stores.value.length ? Promise.resolve({ list: stores.value, total: stores.value.length }) : fetchStores(),
    ]);
    list.value = data.list || [];
    total.value = data.total || 0;
    categories.value = cateData || [];
    stores.value = storeData.list || [];
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

function selectStore(id?: number) {
  updateQuery({ mer_id: id ? String(id) : undefined, page: undefined });
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

function goStore(item: StoreDirectoryItem) {
  void router.push({ name: "store", params: { id: item.mer_id } });
}

watch(
  () => [route.query.keyword, route.query.cate_id, route.query.mer_id, route.query.page, route.query.sort, route.query.order],
  () => void load(),
  { immediate: true },
);
</script>

<template>
  <div class="goods-page">
    <div class="pc-container">
      <nav class="crumb" aria-label="当前位置">
        <RouterLink to="/">首页</RouterLink><span>›</span><span>全部商品</span>
        <template v-if="keyword"><span>›</span><span>搜索：{{ keyword }}</span></template>
      </nav>

      <section v-if="stores.length" class="store-directory" aria-label="推荐店铺">
        <div class="section-title"><h1>推荐店铺</h1><button type="button" @click="selectStore()">查看全部商品</button></div>
        <div class="store-grid">
          <article v-for="item in stores.slice(0, 3)" :key="item.store_id" class="store-card" @click="goStore(item)">
            <img v-if="item.cover_url" :src="item.cover_url" :alt="item.name" />
            <div v-else class="store-image-empty">店铺</div>
            <div class="store-card__info">
              <h2>{{ item.name }}</h2>
              <p>在售 {{ item.product_count }} 件 · 已售 {{ item.sales_count }}</p>
              <button type="button" @click.stop="selectStore(item.mer_id)">查看商品</button>
            </div>
          </article>
        </div>
      </section>

      <section class="filter-panel" aria-label="商品筛选">
        <div class="filter-row">
          <b>商品分类</b>
          <div class="filter-options">
            <button type="button" :class="{ active: !cateId }" @click="selectCategory()">全部</button>
            <button v-for="category in categories" :key="category.id" type="button" :class="{ active: cateId === category.id }" @click="selectCategory(category.id)">{{ category.name }}</button>
          </div>
        </div>
        <div class="filter-row">
          <b>店铺筛选</b>
          <div class="filter-options">
            <button type="button" :class="{ active: !merId }" @click="selectStore()">全部店铺</button>
            <button v-for="store in stores" :key="store.store_id" type="button" :class="{ active: merId === store.mer_id }" @click="selectStore(store.mer_id)">{{ store.name }}</button>
          </div>
        </div>
        <div class="filter-row sort-row">
          <b>排序方式</b>
          <div class="filter-options sort-options">
            <button type="button" :class="{ active: sort === 'default' }" @click="selectSort('default')">默认排序</button>
            <button type="button" :class="{ active: sort === 'sales' }" @click="selectSort('sales')">销量优先</button>
            <button type="button" :class="{ active: sort === 'price' }" @click="selectSort('price')">价格 <span v-if="sort === 'price'">{{ order === 'asc' ? '↑' : '↓' }}</span></button>
          </div>
        </div>
      </section>

      <section class="result-head">
        <div><h1>{{ pageTitle }}</h1><p v-if="!loading">共 <em>{{ total }}</em> 件商品</p></div>
        <div v-if="totalPages > 1" class="head-pages"><span>{{ page }} / {{ totalPages }}</span><button type="button" :disabled="page <= 1" @click="goPage(page - 1)">‹</button><button type="button" :disabled="page >= totalPages" @click="goPage(page + 1)">›</button></div>
      </section>

      <p v-if="loading" class="hint">商品加载中…</p>
      <p v-else-if="hint" class="hint empty">{{ hint }}</p>
      <div v-else class="grid"><ProductCard v-for="product in list" :key="product.id" :product="product" /></div>

      <nav v-if="totalPages > 1" class="pagination" aria-label="商品分页">
        <button type="button" :disabled="page <= 1" @click="goPage(page - 1)">上一页</button>
        <button v-for="number in totalPages" :key="number" type="button" :class="{ active: number === page }" @click="goPage(number)">{{ number }}</button>
        <button type="button" :disabled="page >= totalPages" @click="goPage(page + 1)">下一页</button>
      </nav>
    </div>
  </div>
</template>

<style scoped>
.goods-page { min-height: 680px; padding: 16px 0 50px; background: #f5f5f5; }
.crumb { display: flex; gap: 9px; align-items: center; height: 42px; color: #777; font-size: 13px; }.crumb a:hover { color: #f13728; }
.store-directory { overflow: hidden; margin-bottom: 16px; border: 1px solid #ececec; background: #fff; }.section-title { display: flex; align-items: center; justify-content: space-between; height: 56px; padding: 0 20px; border-bottom: 1px solid #eee; }.section-title h1 { margin: 0; color: #333; font-size: 17px; }.section-title button { border: 0; color: #888; background: transparent; font-size: 13px; }.section-title button:hover { color: #f13728; }
.store-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 20px; padding: 18px 20px; }.store-card { display: grid; grid-template-columns: 92px minmax(0, 1fr); gap: 14px; align-items: center; min-width: 0; cursor: pointer; }.store-card > img, .store-image-empty { width: 92px; height: 92px; object-fit: cover; background: #f6f6f6; }.store-image-empty { display: grid; place-items: center; color: #aaa; font-size: 13px; }.store-card__info { min-width: 0; }.store-card h2 { overflow: hidden; margin: 0; color: #333; font-size: 15px; font-weight: 600; text-overflow: ellipsis; white-space: nowrap; }.store-card p { overflow: hidden; margin: 8px 0 11px; color: #999; font-size: 12px; text-overflow: ellipsis; white-space: nowrap; }.store-card button { border: 1px solid #f13728; padding: 5px 10px; color: #f13728; background: #fff; font-size: 12px; }.store-card button:hover { color: #fff; background: #f13728; }
.filter-panel { border: 1px solid #ededed; background: #fff; }.filter-row { display: grid; grid-template-columns: 116px 1fr; min-height: 52px; border-bottom: 1px solid #f0f0f0; }.filter-row:last-child { border-bottom: 0; }.filter-row > b { display: flex; align-items: center; padding-left: 22px; color: #555; background: #fafafa; font-size: 13px; font-weight: 500; }.filter-options { display: flex; flex-wrap: wrap; align-items: center; gap: 5px 22px; padding: 8px 18px; }.filter-options button { border: 0; padding: 4px 0; color: #555; background: transparent; font-size: 13px; }.filter-options button:hover, .filter-options button.active { color: #f13728; }.sort-options { gap: 0; padding-block: 0; }.sort-options button { min-width: 96px; min-height: 34px; padding: 0 13px; border-right: 1px solid #eee; }.sort-options button:first-child { border-left: 1px solid #eee; }.sort-options button.active { color: #fff; background: #f13728; }
.result-head { display: flex; align-items: end; justify-content: space-between; margin: 22px 0 13px; }.result-head h1 { margin: 0; color: #333; font-size: 18px; font-weight: 600; }.result-head p { margin: 5px 0 0; color: #888; font-size: 13px; }.result-head em { color: #f13728; font-style: normal; }.head-pages { display: flex; align-items: center; gap: 8px; color: #999; font-size: 12px; }.head-pages button { width: 28px; height: 25px; border: 1px solid #ddd; background: #fff; color: #666; font-size: 18px; line-height: 1; }.head-pages button:disabled { color: #ccc; }
.hint { margin: 40px 0; color: #999; text-align: center; }.hint.empty { padding: 52px 0; border: 1px solid #eee; background: #fff; }.grid { display: grid; grid-template-columns: repeat(5, minmax(0, 1fr)); gap: 14px; }.pagination { display: flex; justify-content: center; gap: 8px; margin-top: 32px; }.pagination button { min-width: 36px; height: 34px; padding: 0 10px; border: 1px solid #ddd; color: #555; background: #fff; }.pagination button:hover:not(:disabled), .pagination button.active { border-color: #f13728; color: #fff; background: #f13728; }.pagination button:disabled { color: #bbb; }
@media (max-width: 980px) { .store-grid { grid-template-columns: 1fr 1fr; }.grid { grid-template-columns: repeat(3, minmax(0, 1fr)); } }
@media (max-width: 680px) { .store-grid { grid-template-columns: 1fr; }.filter-row { grid-template-columns: 84px 1fr; }.filter-row > b { padding-left: 12px; }.filter-options { gap: 5px 12px; padding-inline: 12px; }.grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
</style>
