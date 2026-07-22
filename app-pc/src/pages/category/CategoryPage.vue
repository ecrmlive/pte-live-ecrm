<script setup lang="ts">
import { onMounted, ref } from "vue";
import { fetchCategories, type CategoryItem } from "@/api/catalog";

const list = ref<CategoryItem[]>([]);
const hint = ref("加载中…");

const demo: CategoryItem[] = [
  { id: 1, name: "食品生鲜", pid: 0 },
  { id: 2, name: "家居生活", pid: 0 },
  { id: 3, name: "数码电器", pid: 0 },
  { id: 4, name: "美妆个护", pid: 0 },
  { id: 5, name: "服饰鞋包", pid: 0 },
  { id: 6, name: "母婴亲子", pid: 0 },
];

onMounted(async () => {
  try {
    const data = await fetchCategories();
    list.value = data?.length ? data : demo;
    hint.value = data?.length ? "" : "分类接口暂无数据 · 展示演示类目（阶段 2）";
  } catch {
    list.value = demo;
    hint.value = "分类接口暂不可用 · 展示演示类目（阶段 2）";
  }
});
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>分类广场</h1>
      <p class="sub">平台三级分类入口；商户店内二级分类在店铺页（阶段 2）。</p>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div class="grid">
        <RouterLink
          v-for="c in list"
          :key="c.id"
          class="item"
          :to="{ name: 'goods-list', query: { cate_id: String(c.id) } }"
        >
          {{ c.name }}
        </RouterLink>
      </div>
    </section>
  </div>
</template>

<style scoped>
.panel {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: calc(var(--pc-radius) + 4px);
  padding: 1.5rem 1.7rem 1.8rem;
  box-shadow: var(--pc-shadow);
}

h1 {
  margin: 0;
  font-size: 1.45rem;
}

.sub,
.hint {
  color: var(--pc-muted);
}

.sub {
  margin: 0.4rem 0 0.9rem;
}

.hint {
  margin: 0 0 1rem;
  font-size: 0.92rem;
}

.grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.75rem;
}

.item {
  padding: 1.1rem 1rem;
  border: 1px solid var(--pc-line);
  border-radius: 8px;
  font-weight: 600;
}

.item:hover {
  border-color: rgba(15, 107, 92, 0.45);
  color: var(--pc-brand);
  background: var(--pc-brand-soft);
}

@media (max-width: 720px) {
  .grid {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
