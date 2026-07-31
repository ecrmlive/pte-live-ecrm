<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import { fetchCategories, fetchHome, type CategoryItem, type ProductItem } from "@/api/catalog";
import ProductCard from "@/components/ProductCard.vue";

const banners = ref<{ id: number; title: string; image: string; url?: string }[]>([]);
const hot = ref<ProductItem[]>([]);
const categories = ref<CategoryItem[]>([]);
const loading = ref(true);
const errorMsg = ref("");
const bannerIndex = ref(0);
let bannerTimer: ReturnType<typeof setInterval> | undefined;

const heroImage = computed(() => banners.value[bannerIndex.value]?.image || "");
const visibleCategories = computed(() => categories.value.filter((item) => item.pid === 0).slice(0, 8));
const categoryNav = computed(() => (visibleCategories.value.length ? visibleCategories.value : categories.value.slice(0, 8)));
const featuredProducts = computed(() => hot.value.slice(0, 3));
const seasonProducts = computed(() => hot.value.slice(0, 5));
const recommendationProducts = computed(() => hot.value.slice(2, 5));
const rankingProducts = computed(() => hot.value.slice(0, 3));
const floorProducts = computed(() => hot.value.slice(0, 6));
const primaryStore = computed(() => hot.value.find((item) => item.mer_id > 0));
const primaryStoreProductCount = computed(() => {
  const merchantId = primaryStore.value?.mer_id;
  return merchantId ? hot.value.filter((item) => item.mer_id === merchantId).length : 0;
});

onMounted(async () => {
  loading.value = true;
  errorMsg.value = "";
  try {
    const [home, categoryList] = await Promise.all([fetchHome(), fetchCategories()]);
    banners.value = home.banners || [];
    hot.value = home.hot || [];
    categories.value = categoryList;
    startBannerTimer();
  } catch (e) {
    banners.value = [];
    hot.value = [];
    errorMsg.value = (e as Error).message || "首页数据暂不可用";
  } finally {
    loading.value = false;
  }
});

onUnmounted(() => stopBannerTimer());

function switchBanner(direction: number) {
  if (!banners.value.length) return;
  bannerIndex.value = (bannerIndex.value + direction + banners.value.length) % banners.value.length;
  restartBannerTimer();
}

function selectBanner(index: number) {
  bannerIndex.value = index;
  restartBannerTimer();
}

function startBannerTimer() {
  stopBannerTimer();
  if (banners.value.length > 1) {
    bannerTimer = setInterval(() => {
      bannerIndex.value = (bannerIndex.value + 1) % banners.value.length;
    }, 5000);
  }
}

function stopBannerTimer() {
  if (bannerTimer) window.clearInterval(bannerTimer);
  bannerTimer = undefined;
}

function restartBannerTimer() {
  startBannerTimer();
}

function productName(product: ProductItem) {
  return product.title || product.store_name;
}
</script>

<template>
  <div class="home">
    <section class="pc-container storefront">
      <div class="hero-slot" aria-label="首页活动轮播">
        <RouterLink v-if="heroImage && banners[bannerIndex]?.url" class="hero-link" :to="banners[bannerIndex]?.url || '/goods'">
          <img :src="heroImage" :alt="banners[bannerIndex]?.title || '商城活动'" />
        </RouterLink>
        <img v-else-if="heroImage" :src="heroImage" :alt="banners[bannerIndex]?.title || '商城活动'" />
        <div v-else class="hero-empty" role="status">首页装修内容加载中</div>
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
              @click="selectBanner(index)"
            />
          </div>
        </template>
      </div>
      <aside class="category-rail" aria-label="商品分类">
        <div class="category-rail__body">
          <RouterLink v-for="item in categoryNav" :key="item.id" :to="`/goods?cate_id=${item.id}`">
            <span>{{ item.name }}</span><b>›</b>
          </RouterLink>
          <p v-if="!categories.length" class="category-empty">商品分类加载中</p>
        </div>
      </aside>
    </section>

    <section class="pc-container section first-section">
      <div class="section-head">
        <div><h2>精品推荐</h2><span>诚意推荐 品质商品</span></div>
        <RouterLink to="/goods">更多 ›</RouterLink>
      </div>
      <p v-if="loading" class="hint">加载中…</p>
      <p v-else-if="errorMsg" class="hint">{{ errorMsg }}</p>
      <p v-else-if="!hot.length" class="hint">暂无在售商品</p>
      <div v-else class="feature-grid">
        <RouterLink v-for="product in featuredProducts" :key="product.id" class="feature-card" :to="`/goods/${product.id}`">
          <div class="feature-card__copy">
            <h3>{{ productName(product) }}</h3>
            <p>{{ product.shop_name || product.mer_name || '平台精选好物' }}</p>
            <strong>¥{{ product.price }}</strong>
          </div>
          <img :src="product.image" :alt="productName(product)" />
        </RouterLink>
      </div>
    </section>

    <section class="pc-container section">
      <div class="section-head">
        <div><h2>当季新品</h2><span>本季好物 新鲜上架</span></div>
        <RouterLink to="/goods?sort=sales">更多 ›</RouterLink>
      </div>
      <div class="season-grid">
        <ProductCard v-for="product in seasonProducts" :key="`season-${product.id}`" :product="product" />
      </div>
    </section>

    <section class="pc-container discovery-row" aria-label="推荐、店铺与榜单">
      <article class="discovery-card recommended-card">
        <header><h2>推荐单品</h2><span>精挑细选 天天低价</span></header>
        <RouterLink v-for="product in recommendationProducts" :key="`recommend-${product.id}`" class="compact-product" :to="`/goods/${product.id}`">
          <img :src="product.image" :alt="productName(product)" />
          <div><h3>{{ productName(product) }}</h3><strong>¥{{ product.price }}</strong></div>
        </RouterLink>
      </article>

      <article class="discovery-card store-card">
        <header><h2>品牌好店</h2><span>优质好店值得收藏</span></header>
        <RouterLink v-if="primaryStore" class="store-spotlight" :to="`/store/${primaryStore?.mer_id || 0}`">
          <img :src="primaryStore.image" :alt="primaryStore.shop_name || primaryStore.mer_name" />
          <div><h3>{{ primaryStore.shop_name || primaryStore.mer_name }}</h3><p>在售 {{ primaryStoreProductCount }} 件商品</p><b>进店逛逛 ›</b></div>
        </RouterLink>
        <p v-else class="empty-note">暂无可展示店铺</p>
      </article>

      <article class="discovery-card ranking-card">
        <header><h2>热门榜单</h2><span>大家都在买的好物</span></header>
        <RouterLink v-for="(product, index) in rankingProducts" :key="`rank-${product.id}`" class="rank-item" :to="`/goods/${product.id}`">
          <i>{{ String(index + 1).padStart(2, '0') }}</i>
          <img :src="product.image" :alt="productName(product)" />
          <div><h3>{{ productName(product) }}</h3><p>已售 {{ product.sales }}</p></div>
        </RouterLink>
      </article>
    </section>

    <section class="pc-container section category-floor">
      <div class="section-head floor-title">
        <div><h2>{{ categoryNav[0]?.name || '上门服务' }}</h2><span>品质服务 安心到家</span></div>
        <RouterLink :to="categoryNav[0] ? `/goods?cate_id=${categoryNav[0].id}` : '/goods'">更多 ›</RouterLink>
      </div>
      <div class="floor-layout">
        <RouterLink class="floor-poster" :to="categoryNav[0] ? `/goods?cate_id=${categoryNav[0].id}` : '/goods'">
          <img src="/demo/home-service-vertical-v1.png" alt="上门服务" />
        </RouterLink>
        <div class="floor-products">
          <ProductCard v-for="product in floorProducts.slice(0, 4)" :key="`service-${product.id}`" :product="product" />
        </div>
      </div>
      <RouterLink class="floor-wide" to="/reservation"><img src="/demo/home-service-wide-v1.png" alt="七禧到家服务" /></RouterLink>
    </section>

    <section class="pc-container section category-floor">
      <div class="section-head floor-title">
        <div><h2>{{ categoryNav[1]?.name || '家居百货' }}</h2><span>日常所需 一站购齐</span></div>
        <RouterLink :to="categoryNav[1] ? `/goods?cate_id=${categoryNav[1].id}` : '/goods'">更多 ›</RouterLink>
      </div>
      <div class="floor-layout reverse">
        <RouterLink class="floor-poster" :to="categoryNav[1] ? `/goods?cate_id=${categoryNav[1].id}` : '/goods'">
          <img src="/demo/home-beauty-vertical-v1.png" alt="品质生活" />
        </RouterLink>
        <div class="floor-products">
          <ProductCard v-for="product in floorProducts.slice(2, 6)" :key="`living-${product.id}`" :product="product" />
        </div>
      </div>
      <RouterLink class="floor-wide" to="/goods"><img src="/demo/home-tech-wide-v1.png" alt="品质生活精选" /></RouterLink>
    </section>

    <section class="category-square">
      <div class="pc-container">
        <div class="section-head"><div><h2>分类广场</h2><span>按分类快速找到好商品</span></div></div>
        <div class="category-grid">
          <RouterLink v-for="item in categoryNav" :key="`square-${item.id}`" :to="`/goods?cate_id=${item.id}`">
            <strong>{{ item.name }}</strong><span>精选商品</span>
          </RouterLink>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.home { background: #fff; }
.storefront { position: relative; min-height: 480px; }
.category-rail { position: absolute; z-index: 3; inset: 0 auto 0 0; width: 238px; background: rgb(41 47 45 / 78%); }
.category-rail__body { padding: .5rem 0 1rem; }
.category-rail a { display: flex; align-items: center; justify-content: space-between; padding: .83rem 1.55rem; color: #fff; font-size: 1rem; }
.category-rail b { color: rgba(255, 255, 255, .75); font-size: 1.45rem; font-weight: 300; line-height: 1; }
.category-rail a:hover { background: rgb(0 0 0 / 18%); }
.category-empty { color: #ddd; padding: .78rem 1.55rem; margin: 0; }
.hero-slot { position: relative; min-width: 0; background: #f6f6f6; }
.hero-link { display: block; }
.hero-slot img { display: block; width: 100%; height: 480px; object-fit: cover; }
.hero-empty { display: grid; height: 480px; place-content: center; color: #777; background: #f4f4f4; text-align: center; }
.hero-control { position: absolute; top: 50%; z-index: 1; width: 42px; height: 64px; border: 0; color: #fff; background: rgba(0, 0, 0, .22); font-size: 2.3rem; line-height: 1; transform: translateY(-50%); }
.hero-control:hover { background: rgba(0, 0, 0, .45); }.hero-control.previous { left: 0; }.hero-control.next { right: 0; }
.hero-dots { position: absolute; right: 26px; bottom: 18px; display: flex; gap: .45rem; }
.hero-dots button { width: 25px; height: 3px; border: 0; background: rgba(255, 255, 255, .58); }.hero-dots button.active { background: #fff; }

.section { margin-top: 2.25rem; }
.first-section { margin-top: 2.8rem; }
.section-head { display: flex; align-items: baseline; justify-content: space-between; margin: 0 0 1rem; }
.section-head h2 { display: inline-block; margin: 0; padding-bottom: .25rem; border-bottom: 3px solid #f13728; font-size: 1.5rem; }
.section-head a { padding: .3rem .7rem; border: 1px solid #ddd; color: #666; font-size: .85rem; }.section-head a:hover { border-color: #f13728; color: #f13728; }
.section-head span { margin-left: .75rem; color: #999; font-size: .87rem; }
.hint { margin: 0 0 .85rem; color: var(--pc-muted); font-size: .92rem; }

.feature-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
.feature-card { display: grid; grid-template-columns: minmax(0, 1fr) 45%; min-height: 172px; overflow: hidden; border: 1px solid #eee; background: #fff; }
.feature-card:hover { border-color: #f13728; box-shadow: 0 10px 20px rgb(0 0 0 / 7%); }
.feature-card__copy { padding: 1.3rem 0 1.1rem 1.35rem; }.feature-card h3 { display: -webkit-box; margin: 0; overflow: hidden; color: #333; font-size: 1.03rem; line-height: 1.45; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.feature-card p { margin: .55rem 0 1.05rem; overflow: hidden; color: #999; font-size: .83rem; text-overflow: ellipsis; white-space: nowrap; }.feature-card strong { color: #f13728; font-size: 1.15rem; }
.feature-card img { width: 100%; height: 100%; min-height: 172px; object-fit: cover; }

.season-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 12px; }
.season-grid :deep(.meta) { padding: .72rem .75rem .85rem; }.season-grid :deep(h3) { font-size: .88rem; }.season-grid :deep(.store) { margin: .25rem 0 .45rem; font-size: .76rem; }.season-grid :deep(.price-row strong) { font-size: .98rem; }.season-grid :deep(.sales) { display: none; }

.discovery-row { display: grid; grid-template-columns: 1.04fr 1.08fr 1.18fr; gap: 12px; margin-top: 2.6rem; }
.discovery-card { min-height: 285px; border: 1px solid #eee; background: #fff; }.discovery-card header { padding: 1rem 1.15rem .8rem; border-bottom: 1px solid #f2f2f2; }.discovery-card h2 { display: inline; margin: 0; color: #333; font-size: 1.12rem; }.discovery-card header span { margin-left: .5rem; color: #999; font-size: .75rem; }
.compact-product { display: flex; gap: .75rem; padding: .72rem 1rem 0; }.compact-product img { width: 66px; height: 66px; object-fit: cover; }.compact-product div { min-width: 0; }.compact-product h3, .rank-item h3 { display: -webkit-box; margin: .08rem 0 .35rem; overflow: hidden; color: #444; font-size: .86rem; line-height: 1.35; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }.compact-product strong { color: #f13728; font-size: .95rem; }
.store-spotlight { display: grid; grid-template-columns: 46% 1fr; gap: 1rem; padding: 1.1rem; }.store-spotlight img { width: 100%; aspect-ratio: 1; object-fit: cover; }.store-spotlight h3 { margin: .25rem 0 .65rem; color: #333; font-size: 1.05rem; }.store-spotlight p { margin: 0 0 1rem; color: #999; font-size: .82rem; line-height: 1.55; }.store-spotlight b { color: #f13728; font-size: .84rem; font-weight: 500; }.empty-note { padding: 1.2rem; color: #999; }
.rank-item { display: grid; grid-template-columns: 28px 54px minmax(0, 1fr); gap: .6rem; align-items: center; padding: .65rem 1rem 0; }.rank-item i { display: grid; width: 24px; height: 24px; place-content: center; color: #fff; background: #f0b23c; font-size: .72rem; font-style: normal; }.rank-item:nth-of-type(3) i { background: #b7c0d4; }.rank-item:nth-of-type(n+4) i { background: #d3d3d3; }.rank-item img { width: 54px; height: 54px; object-fit: cover; }.rank-item h3 { margin: 0 0 .18rem; }.rank-item p { margin: 0; color: #999; font-size: .75rem; }

.category-floor { margin-top: 2.8rem; }.floor-layout { display: grid; grid-template-columns: 248px 1fr; gap: 12px; }.floor-poster { display: block; overflow: hidden; background: #f2f2f2; }.floor-poster img { display: block; width: 100%; height: 100%; min-height: 440px; object-fit: cover; }.floor-products { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }.floor-products :deep(.meta) { padding: .72rem .75rem .85rem; }.floor-products :deep(h3) { font-size: .88rem; }.floor-products :deep(.store) { margin: .25rem 0 .45rem; font-size: .76rem; }.floor-products :deep(.sales) { display: none; }.floor-wide { display: block; margin-top: 12px; overflow: hidden; background: #f2f2f2; }.floor-wide img { display: block; width: 100%; height: 154px; object-fit: cover; }

.category-square { margin-top: 2.8rem; padding: 2.2rem 0 3rem; background: #f7f7f7; }.category-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1px; background: #e6e6e6; }.category-grid a { display: grid; min-height: 96px; align-content: center; gap: .4rem; padding: 1.3rem; background: #fff; transition: background .15s ease; }.category-grid a:hover { background: #fff4f2; }.category-grid strong { color: #333; font-size: 1rem; }.category-grid span { color: #999; font-size: .82rem; }

@media (max-width: 980px) { .category-rail { display: none; }.hero-slot img, .hero-empty { height: 300px; }.feature-grid, .season-grid { grid-template-columns: repeat(2, 1fr); }.discovery-row { grid-template-columns: 1fr; }.floor-layout { grid-template-columns: 1fr; }.floor-poster { display: none; }.floor-products { grid-template-columns: repeat(2, 1fr); }.category-grid { grid-template-columns: repeat(2, 1fr); } }
</style>
