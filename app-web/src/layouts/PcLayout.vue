<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";

const router = useRouter();
const route = useRoute();
const user = useUserStore();

const keyword = ref(typeof route.query.keyword === "string" ? route.query.keyword : "");
const popularKeywords = ["海鲜", "预制菜", "国潮", "箱包", "家居", "口红", "iPhone"];

const navs = [
  { to: "/goods", label: "全部商品" },
  { to: "/coupons", label: "领券中心" },
  { to: "/notices", label: "新闻中心" },
  { to: "/seckill", label: "秒杀列表" },
  { to: "/reservation", label: "服务咨询" },
  { to: "/merchant/apply", label: "商户入驻" },
  { to: "/user", label: "我的信息" },
];

const activePath = computed(() => route.path);

function search() {
  const q = keyword.value.trim();
  router.push({ name: "goods-list", query: q ? { keyword: q } : {} });
}

function goLogin() {
  router.push({ name: "login", query: { redirect: route.fullPath } });
}

function logout() {
  user.logout();
  router.push({ name: "home" });
}

onMounted(() => void user.refreshCartCount());
watch(
  () => user.isLogin,
  () => void user.refreshCartCount(),
);
</script>

<template>
  <div class="shell">
    <header class="site-header">
      <div class="utility-bar">
        <div class="pc-container utility-inner">
          <span>欢迎来到七禧商城</span>
          <div class="utility-links">
            <template v-if="user.isLogin">
              <RouterLink to="/user">{{ user.displayName }}</RouterLink>
              <button type="button" @click="logout">退出</button>
            </template>
            <button v-else type="button" @click="goLogin">登录 / 注册</button>
            <RouterLink to="/orders">我的订单</RouterLink>
            <RouterLink to="/merchant/apply">店铺入驻</RouterLink>
            <RouterLink to="/notices">资讯信息</RouterLink>
          </div>
        </div>
      </div>
      <div class="pc-container masthead">
        <RouterLink class="brand" to="/">
          <strong>七禧</strong><span>商城</span>
        </RouterLink>
        <form class="search" @submit.prevent="search">
          <input
            v-model="keyword"
            type="search"
            placeholder="搜索商品 / 店铺"
            aria-label="搜索商品或店铺"
          />
          <button type="submit">搜索</button>
          <div class="popular-keywords" aria-label="热门搜索">
            <button v-for="item in popularKeywords" :key="item" type="button" @click="keyword = item; search()">
              {{ item }}
            </button>
          </div>
        </form>
        <RouterLink class="cart-entry" to="/cart">购物车 <b>({{ user.cartCount }})</b></RouterLink>
      </div>
      <nav class="main-nav">
        <div class="pc-container nav-inner">
          <RouterLink class="category-tab" to="/category">商品分类</RouterLink>
          <RouterLink
            v-for="item in navs"
            :key="item.to"
            :to="item.to"
            class="nav-link"
            :class="{ active: item.to === '/' ? activePath === '/' : activePath.startsWith(item.to) }"
          >
            {{ item.label }}
          </RouterLink>
        </div>
      </nav>
    </header>

    <main class="main">
      <RouterView />
    </main>

    <footer class="footer">
      <div class="pc-container footer-inner">
        <p>七禧商城 PC · 对齐功能表 4 · API `/api/app/v1`</p>
        <p class="muted">与小程序、H5 共享 C 端业务契约</p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.shell {
  min-height: 100vh;
  display: grid;
  grid-template-rows: auto 1fr auto;
}

.site-header { background: #fff; }

.utility-bar { background: #2e2e2e; color: #c9c9c9; font-size: .78rem; }

.utility-inner { height: 40px; display: flex; align-items: center; justify-content: space-between; }

.utility-links { display: flex; gap: 1rem; align-items: center; }

.utility-links a, .utility-links button { color: inherit; border: 0; padding: 0; background: transparent; font: inherit; cursor: pointer; }

.utility-links a:hover, .utility-links button:hover { color: #fff; }

.masthead {
  display: grid; grid-template-columns: 230px minmax(360px, 730px) 150px;
  min-height: 160px; align-items: center; justify-content: space-between; gap: 2.2rem;
}

.brand {
  color: #151515; font-family: Arial, sans-serif; letter-spacing: -.11em; white-space: nowrap;
}
.brand strong { font-size: 3.55rem; font-weight: 800; }
.brand span { color: #f13728; font-size: 1.1rem; letter-spacing: 0; margin-left: .42rem; }

.search {
  position: relative;
  display: flex;
  max-width: 730px;
  width: 100%;
  justify-self: center;
}

.search input {
  flex: 1;
  border: 1px solid #f13728;
  border-right: 0;
  border-radius: 0;
  height: 46px;
  padding: 0.85rem 1rem;
  background: #fff;
  outline: none;
}

.search input:focus {
  border-color: #f13728;
}

.search > button { width: 104px; border: 0; background: #f13728; color: #fff; font-weight: 700; cursor: pointer; }
.popular-keywords { position: absolute; top: 54px; left: .75rem; display: flex; gap: 1.1rem; white-space: nowrap; }
.popular-keywords button { padding: 0; color: #999; border: 0; background: transparent; font-size: .88rem; }
.popular-keywords button:hover { color: #f13728; }
.cart-entry { border: 1px solid #e5e5e5; padding: .8rem 1rem; text-align: center; color: #f13728; }
.cart-entry b { font-weight: 500; }

.main-nav { border-top: 1px solid #f0f0f0; border-bottom: 1px solid #e9e9e9; }

.nav-inner {
  display: flex;
  gap: 2.6rem; padding: 0;
}

.category-tab { width: 238px; padding: 1rem 0; color: #fff; background: #f13728; text-align: center; font-weight: 700; }

.nav-link {
  color: #2c2c2c; padding: 1rem 0; border-bottom: 2px solid transparent; font-weight: 600;
}

.nav-link.active,
.nav-link:hover {
  color: #f13728; border-bottom-color: #f13728;
}

.main {
  padding: 0 0 2.5rem;
}

.footer {
  border-top: 1px solid var(--pc-line);
  background: rgba(255, 255, 255, 0.7);
}

.footer-inner {
  padding: 1.25rem 0 1.5rem;
}

.footer p {
  margin: 0;
}

.footer .muted {
  margin-top: 0.35rem;
  color: var(--pc-muted);
  font-size: 0.9rem;
}

@media (max-width: 860px) {
  .masthead {
    grid-template-columns: 1fr;
    padding: 1.5rem 0;
    gap: 1rem;
  }

  .search {
    max-width: none;
  }

  .nav-inner {
    overflow-x: auto;
  }
}
</style>
