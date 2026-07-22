<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchProductDetail, type ProductDetail } from "@/api/catalog";
import { addCart } from "@/api/trade";
import { useUserStore } from "@/stores/user";
import { ApiError } from "@/utils/request";

const route = useRoute();
const router = useRouter();
const user = useUserStore();
const detail = ref<ProductDetail | null>(null);
const hint = ref("");
const adding = ref(false);

const id = computed(() => Number(route.params.id));

async function load() {
  if (!id.value) return;
  try {
    detail.value = await fetchProductDetail(id.value);
    hint.value = "";
  } catch {
    detail.value = {
      id: id.value,
      mer_id: 1,
      mer_name: "栖息优选店",
      store_name: "演示商品详情",
      image: "",
      price: "99.00",
      ot_price: "129.00",
      sales: 12,
      stock: 50,
      unit_name: "件",
      store_info: "商品详情将在阶段 2 对接目录接口。当前为 PC 骨架占位。",
      slider_image: [],
      spec_type: 0,
      delivery_way: "快递",
    };
    hint.value = "详情接口暂不可用 · 展示演示数据（阶段 2）";
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
    <section class="panel">
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div class="layout">
        <div class="gallery">
          <div class="cover">{{ detail.store_name.slice(0, 1) }}</div>
        </div>
        <div class="info">
          <p class="mer">
            <RouterLink :to="`/store/${detail.mer_id}`">{{ detail.mer_name }}</RouterLink>
          </p>
          <h1>{{ detail.store_name }}</h1>
          <p class="desc">{{ detail.store_info }}</p>
          <div class="price">
            <strong>¥{{ detail.price }}</strong>
            <span v-if="detail.ot_price">¥{{ detail.ot_price }}</span>
          </div>
          <ul>
            <li>销量 {{ detail.sales }} · 库存 {{ detail.stock }} {{ detail.unit_name }}</li>
            <li>配送：{{ detail.delivery_way }}</li>
            <li>规格：{{ detail.spec_type === 0 ? "单规格" : "多规格（阶段 2）" }}</li>
          </ul>
          <div class="actions">
            <button class="pc-btn" type="button" :disabled="adding" @click="onAddCart">加入购物车</button>
            <RouterLink class="pc-btn ghost" to="/cart">去购物车</RouterLink>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.panel {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: calc(var(--pc-radius) + 4px);
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
  border-radius: var(--pc-radius);
  display: grid;
  place-items: center;
  font-size: 4rem;
  font-weight: 700;
  color: var(--pc-brand);
  background: linear-gradient(160deg, #dceeea, #f7faf9);
}

.mer {
  margin: 0;
  color: var(--pc-brand);
  font-weight: 600;
}

h1 {
  margin: 0.4rem 0 0.7rem;
  font-size: 1.7rem;
}

.desc {
  margin: 0;
  color: var(--pc-muted);
  line-height: 1.7;
}

.price {
  margin: 1rem 0;
  display: flex;
  gap: 0.6rem;
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

@media (max-width: 900px) {
  .layout {
    grid-template-columns: 1fr;
  }
}
</style>
