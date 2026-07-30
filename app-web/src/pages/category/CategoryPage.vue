<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { fetchCategories, type CategoryItem } from "@/api/catalog";

const list = ref<CategoryItem[]>([]);
const hint = ref("加载中…");
const selectedID = ref<number>();

const roots = computed(() => list.value.filter((item) => !item.pid));
const groups = computed(() => {
  if (!roots.value.length) {
    return list.value.map((item) => ({ parent: item, children: [] as CategoryItem[] }));
  }
  return roots.value.map((parent) => ({
    parent,
    children: list.value.filter((item) => item.pid === parent.id),
  }));
});
const selected = computed(
  () => groups.value.find((group) => group.parent.id === selectedID.value) || groups.value[0],
);

onMounted(async () => {
  try {
    const data = await fetchCategories();
    list.value = data || [];
    selectedID.value = list.value.find((item) => !item.pid)?.id || list.value[0]?.id;
    hint.value = list.value.length ? "" : "暂无商品分类";
  } catch (error) {
    list.value = [];
    hint.value = (error as Error).message || "分类加载失败";
  }
});
</script>

<template>
  <div class="category-page">
    <div class="pc-container">
      <nav class="crumb" aria-label="当前位置"><RouterLink to="/">首页</RouterLink><span>›</span><span>商品分类</span></nav>
      <section class="category-panel">
        <aside class="category-nav" aria-label="一级分类">
          <button
            v-for="group in groups"
            :key="group.parent.id"
            type="button"
            :class="{ active: selected?.parent.id === group.parent.id }"
            @click="selectedID = group.parent.id"
          >
            <img v-if="group.parent.pic" :src="group.parent.pic" alt="" />
            <span>{{ group.parent.name }}</span>
          </button>
        </aside>
        <div class="category-content">
          <template v-if="selected">
            <div class="title-row">
              <h1>{{ selected.parent.name }}</h1>
              <RouterLink :to="{ name: 'goods-list', query: { cate_id: String(selected.parent.id) } }">查看全部 ›</RouterLink>
            </div>
            <div class="children">
              <RouterLink
                v-for="child in selected.children.length ? selected.children : [selected.parent]"
                :key="child.id"
                class="item"
                :to="{ name: 'goods-list', query: { cate_id: String(child.id) } }"
              >
                <div class="pic"><img v-if="child.pic" :src="child.pic" :alt="child.name" /><span v-else>商品分类</span></div>
                <span>{{ child.name }}</span>
              </RouterLink>
            </div>
          </template>
          <p v-else class="hint">{{ hint }}</p>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.category-page { min-height: 560px; padding-top: 18px; background: #f6f6f6; }
.crumb { display: flex; gap: .55rem; align-items: center; height: 42px; color: #888; font-size: .88rem; }
.crumb a:hover { color: #ef3727; }
.category-panel { display: grid; grid-template-columns: 218px 1fr; min-height: 440px; border: 1px solid #ececec; background: #fff; }
.category-nav { padding: 10px 0; background: #fafafa; border-right: 1px solid #eee; }
.category-nav button { display: flex; align-items: center; gap: .7rem; width: 100%; min-height: 51px; padding: 0 22px; border: 0; border-left: 3px solid transparent; background: transparent; color: #555; text-align: left; cursor: pointer; }
.category-nav button:hover, .category-nav button.active { border-left-color: #ef3727; background: #fff; color: #ef3727; }
.category-nav img { width: 26px; height: 26px; object-fit: cover; border-radius: 50%; }
.category-content { padding: 29px 36px 38px; }
.title-row { display: flex; align-items: baseline; justify-content: space-between; padding-bottom: 17px; border-bottom: 1px solid #eee; }
.title-row h1 { margin: 0; color: #333; font-size: 1.25rem; }
.title-row a { color: #888; font-size: .88rem; }
.title-row a:hover { color: #ef3727; }
.children { display: grid; grid-template-columns: repeat(5, minmax(0, 1fr)); gap: 24px 18px; padding-top: 30px; }
.item { display: grid; justify-items: center; gap: .65rem; color: #555; font-size: .92rem; text-align: center; }
.item:hover { color: #ef3727; }
.pic { display: grid; width: 92px; height: 92px; place-items: center; overflow: hidden; background: #f6f6f6; color: #aaa; font-size: .76rem; }
.pic img { width: 100%; height: 100%; object-fit: cover; }
.hint { margin: 9rem 0; color: #999; text-align: center; }
@media (max-width: 860px) { .category-panel { grid-template-columns: 150px 1fr; } .category-content { padding-inline: 20px; } .children { grid-template-columns: repeat(3, minmax(0, 1fr)); } }
@media (max-width: 600px) { .category-panel { grid-template-columns: 1fr; } .category-nav { display: flex; overflow-x: auto; border-right: 0; border-bottom: 1px solid #eee; } .category-nav button { width: auto; min-width: 108px; border-left: 0; border-bottom: 3px solid transparent; } .category-nav button:hover, .category-nav button.active { border-left-color: transparent; border-bottom-color: #ef3727; } .children { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
</style>
