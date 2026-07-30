<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { fetchCategories, fetchHome, type CategoryItem, type ProductItem } from "@/api/catalog";
import ProductCard from "@/components/ProductCard.vue";

const banners = ref<{ id: number; title: string; image: string }[]>([]);
const hot = ref<ProductItem[]>([]);
const categories = ref<CategoryItem[]>([]);
const loading = ref(true);
const errorMsg = ref("");
const bannerIndex = ref(0);

const floors = [
  { title: "精品推荐", subtitle: "诚意推荐 品质商品", type: "best" },
  { title: "火爆新品", subtitle: "爆款尖货 每日上新", type: "new" },
  { title: "推荐单品", subtitle: "精选爆款 天天低价", type: "good" },
];

onMounted(async () => {
  loading.value = true;
  errorMsg.value = "";
  try {
    const data = await fetchHome();
    banners.value = data.banners || [];
    hot.value = data.hot || [];
    categories.value = await fetchCategories();
  } catch (e) {
    banners.value = [];
    hot.value = [];
    errorMsg.value = (e as Error).message || "首页接口暂不可用";
  } finally {
    loading.value = false;
  }
});

const heroImage = computed(() => banners.value[bannerIndex.value]?.image || "");
const visibleCategories = computed(() => categories.value.slice(0, 8));
const showProducts = computed(() => hot.value.slice(0, 8));

function switchBanner(direction: number) {
  if (!banners.value.length) return;
  bannerIndex.value = (bannerIndex.value + direction + banners.value.length) % banners.value.length;
}
</script>

<template>
  <div class="home">
    <section class="pc-container storefront">
      <aside class="category-rail" aria-label="商品分类">
        <RouterLink v-for="item in visibleCategories" :key="item.id" :to="`/goods?cate_id=${item.id}`">
          <span>{{ item.name }}</span><b>›</b>
        </RouterLink>
        <p v-if="!categories.length" class="category-empty">商品分类加载中</p>
      </aside>
      <div class="hero-slot">
        <img v-if="heroImage" :src="heroImage" :alt="banners[0]?.title || '商城活动'" />
        <div v-else class="hero-empty">
          <strong>七禧精选</strong><span>运营 Banner 将由后台装修发布</span>
        </div>
        <template v-if="banners.length > 1">
          <button class="hero-control previous" type="button" aria-label="上一张 Banner" @click="switchBanner(-1)">‹</button>
          <button class="hero-control next" type="button" aria-label="下一张 Banner" @click="switchBanner(1)">›</button>
          <div class="hero-dots" aria-label="Banner 页码">
            <button
              v-for="(_, index) in banners"
              :key="index"
              type="button"
              :class="{ active: bannerIndex === index }"
              :aria-label="`第 ${index + 1} 张 Banner`"
              @click="bannerIndex = index"
            />
          </div>
        </template>
      </div>
    </section>

    <section v-for="floor in floors" :key="floor.type" class="pc-container section">
      <div class="section-head">
        <div><h2>{{ floor.title }}</h2><span>{{ floor.subtitle }}</span></div>
        <RouterLink :to="`/goods?type=${floor.type}`">更多 ›</RouterLink>
      </div>
      <p v-if="loading" class="hint">加载中…</p>
      <p v-else-if="errorMsg" class="hint">{{ errorMsg }}</p>
      <p v-else-if="!hot.length" class="hint">暂无在售商品</p>
      <div class="grid">
        <ProductCard v-for="p in showProducts" :key="`${floor.type}-${p.id}`" :product="p" />
      </div>
    </section>

    <section class="category-square">
      <div class="pc-container">
        <div class="section-head"><div><h2>分类广场</h2><span>按分类快速找到好商品</span></div></div>
        <div class="category-grid">
          <RouterLink v-for="item in visibleCategories" :key="`square-${item.id}`" :to="`/goods?cate_id=${item.id}`">
            <strong>{{ item.name }}</strong><span>精选商品</span>
          </RouterLink>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.home { background: #fff; }
.storefront { display: grid; grid-template-columns: 238px 1fr; min-height: 480px; }
.category-rail { background: #646b69; padding: 1rem 0; }
.category-rail a { display: flex; align-items: center; justify-content: space-between; padding: .78rem 1.55rem; color: #fff; font-size: 1rem; }
.category-rail b { color: rgba(255, 255, 255, .75); font-size: 1.45rem; font-weight: 300; line-height: 1; }
.category-rail a:hover { background: #515756; }
.category-empty { color: #ddd; padding: .78rem 1.55rem; margin: 0; }
.hero-slot { position: relative; background: #f6f6f6; min-width: 0; }
.hero-slot img { width: 100%; height: 480px; object-fit: cover; display: block; }
.hero-empty { height: 480px; display: grid; place-content: center; gap: .85rem; color: #999; background: linear-gradient(115deg, #ececeb, #c8d0cf 58%, #5b6664); text-align: center; }
.hero-empty strong { color: #292929; font-size: 3rem; letter-spacing: .25rem; }
.hero-control { position: absolute; top: 50%; z-index: 1; width: 42px; height: 64px; border: 0; color: #fff; background: rgba(0, 0, 0, .22); font-size: 2.3rem; line-height: 1; transform: translateY(-50%); }
.hero-control:hover { background: rgba(0, 0, 0, .45); }
.hero-control.previous { left: 0; }
.hero-control.next { right: 0; }
.hero-dots { position: absolute; bottom: 18px; right: 26px; display: flex; gap: .45rem; }
.hero-dots button { width: 25px; height: 3px; border: 0; background: rgba(255,255,255,.58); }
.hero-dots button.active { background: #fff; }

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin: 2.25rem 0 1rem;
}

.section-head h2 {
  margin: 0;
  font-size: 1.6rem;
  display: inline-block;
  border-bottom: 3px solid #f13728;
  padding-bottom: .25rem;
}

.section-head a {
  color: #666; border: 1px solid #ddd; padding: .35rem .7rem;
}
.section-head span { margin-left: .75rem; color: #999; }

.hint {
  margin: 0 0 0.85rem;
  color: var(--pc-muted);
  font-size: 0.92rem;
}

.grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.25rem;
}

.category-square { margin-top: 2.25rem; padding: 1px 0 3rem; background: #f7f7f7; }
.category-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1px; background: #e6e6e6; }
.category-grid a { min-height: 112px; display: grid; align-content: center; gap: .45rem; padding: 1.5rem; background: #fff; transition: background .15s ease; }
.category-grid a:hover { background: #fff4f2; }
.category-grid strong { color: #333; font-size: 1.05rem; }
.category-grid span { color: #999; font-size: .85rem; }

@media (max-width: 980px) {
  .storefront { grid-template-columns: 1fr; }
  .category-rail { display: grid; grid-template-columns: repeat(2, 1fr); }
  .hero-slot img, .hero-empty { height: 300px; }

  .grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .category-grid { grid-template-columns: repeat(2, 1fr); }
}
</style>
