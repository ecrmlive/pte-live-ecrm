<script setup lang="ts">
import { onMounted, ref } from "vue";
import { fetchHome, type ProductItem } from "@/api/catalog";
import ProductCard from "@/components/ProductCard.vue";

const banners = ref<{ id: number; title: string; image: string }[]>([]);
const hot = ref<ProductItem[]>([]);
const loading = ref(true);
const usingDemo = ref(false);
const errorMsg = ref("");

const demoBanners = [
  { id: 1, title: "多商户精选 · 本周上新", image: "" },
  { id: 2, title: "品牌店街 · 入驻好店", image: "" },
  { id: 3, title: "限时秒杀 · 场次预告", image: "" },
  { id: 4, title: "优惠券专区 · 领券下单", image: "" },
];

const demoHot: ProductItem[] = [
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
    ot_price: "219.00",
    sales: 180,
    stock: 40,
  },
  {
    id: 1003,
    mer_id: 1,
    mer_name: "栖息优选店",
    store_name: "家居香氛套装",
    image: "",
    price: "59.90",
    sales: 512,
    stock: 200,
  },
  {
    id: 1004,
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

onMounted(async () => {
  loading.value = true;
  errorMsg.value = "";
  try {
    const data = await fetchHome();
    banners.value = data.banners?.length ? data.banners : demoBanners;
    hot.value = data.hot?.length ? data.hot : demoHot;
    usingDemo.value = !data.hot?.length;
  } catch (e) {
    banners.value = demoBanners;
    hot.value = demoHot;
    usingDemo.value = true;
    errorMsg.value = (e as Error).message || "首页接口暂不可用";
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="pc-container home">
    <section class="hero">
      <div class="hero-copy">
        <p class="eyebrow">PC 商城 · 功能表 4</p>
        <h1>搜索好物，逛品牌店，多店结算</h1>
        <p class="sub">
          首页承载轮播、秒杀/预售入口、推荐位与分类广场。当前为 PC-0/1 骨架，接口就绪后接真数据。
        </p>
        <div class="hero-actions">
          <RouterLink class="pc-btn" to="/goods">逛商品</RouterLink>
          <RouterLink class="pc-btn ghost" to="/category">分类广场</RouterLink>
        </div>
      </div>
      <div class="banner-panel">
        <article v-for="b in banners.slice(0, 4)" :key="b.id" class="banner-card">
          <strong>{{ b.title }}</strong>
        </article>
      </div>
    </section>

    <section class="ops">
      <RouterLink class="op" to="/seckill">秒杀</RouterLink>
      <RouterLink class="op" to="/combination">拼团</RouterLink>
      <RouterLink class="op" to="/assist">助力</RouterLink>
      <RouterLink class="op" to="/live">直播</RouterLink>
      <RouterLink class="op" to="/community">社区</RouterLink>
      <RouterLink class="op" to="/reservation">预约</RouterLink>
      <RouterLink class="op" to="/presell">预售</RouterLink>
      <RouterLink class="op" to="/points">积分</RouterLink>
      <RouterLink class="op" to="/coupons">领券中心</RouterLink>
      <RouterLink class="op" to="/goods?tag=hot">热销</RouterLink>
      <RouterLink class="op" to="/goods?tag=brand">品牌店</RouterLink>
      <RouterLink class="op" to="/category">分类广场</RouterLink>
    </section>

    <section class="section">
      <div class="section-head">
        <h2>热销好物</h2>
        <RouterLink to="/goods">全部</RouterLink>
      </div>
      <p v-if="loading" class="hint">加载中…</p>
      <p v-else-if="usingDemo" class="hint">
        {{ errorMsg || "暂无接口数据" }} · 已展示演示商品（后端阶段 2 接真）
      </p>
      <div class="grid">
        <ProductCard v-for="p in hot" :key="p.id" :product="p" />
      </div>
    </section>
  </div>
</template>

<style scoped>
.home {
  display: grid;
  gap: 1.5rem;
}

.hero {
  display: grid;
  grid-template-columns: 1.05fr 1fr;
  gap: 1.25rem;
  align-items: stretch;
}

.hero-copy {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: calc(var(--pc-radius) + 4px);
  padding: 2rem 2.1rem;
  box-shadow: var(--pc-shadow);
}

.eyebrow {
  margin: 0 0 0.6rem;
  color: var(--pc-brand);
  font-weight: 600;
  letter-spacing: 0.04em;
  text-transform: uppercase;
  font-size: 0.8rem;
}

h1 {
  margin: 0;
  font-size: clamp(1.7rem, 2.4vw, 2.3rem);
  line-height: 1.25;
}

.sub {
  margin: 0.9rem 0 1.4rem;
  color: var(--pc-muted);
  line-height: 1.7;
  max-width: 34rem;
}

.hero-actions {
  display: flex;
  gap: 0.7rem;
}

.banner-panel {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

.banner-card {
  min-height: 120px;
  border-radius: var(--pc-radius);
  padding: 1.1rem;
  display: flex;
  align-items: end;
  color: #fff;
  background:
    linear-gradient(160deg, rgba(18, 70, 62, 0.2), rgba(12, 40, 36, 0.55)),
    linear-gradient(135deg, #1a8a76, #244a62);
}

.banner-card:nth-child(2) {
  background:
    linear-gradient(160deg, rgba(80, 40, 20, 0.15), rgba(40, 20, 10, 0.45)),
    linear-gradient(135deg, #c45c26, #8a3d18);
}

.banner-card:nth-child(3) {
  background:
    linear-gradient(160deg, rgba(20, 40, 60, 0.2), rgba(10, 20, 30, 0.5)),
    linear-gradient(135deg, #2f6f8f, #1d3f55);
}

.banner-card:nth-child(4) {
  background:
    linear-gradient(160deg, rgba(40, 50, 30, 0.2), rgba(20, 25, 15, 0.5)),
    linear-gradient(135deg, #5f7a46, #334528);
}

.ops {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 0.65rem;
}

.op {
  text-align: center;
  padding: 0.9rem 0.5rem;
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: 8px;
  color: var(--pc-ink);
  font-weight: 500;
}

.op:hover {
  border-color: rgba(15, 107, 92, 0.4);
  color: var(--pc-brand);
}

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 0.85rem;
}

.section-head h2 {
  margin: 0;
  font-size: 1.25rem;
}

.section-head a {
  color: var(--pc-brand);
}

.hint {
  margin: 0 0 0.85rem;
  color: var(--pc-muted);
  font-size: 0.92rem;
}

.grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
}

@media (max-width: 980px) {
  .hero {
    grid-template-columns: 1fr;
  }

  .ops {
    grid-template-columns: repeat(3, 1fr);
  }

  .grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
