<script setup lang="ts">
import { computed, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";

const router = useRouter();
const route = useRoute();
const user = useUserStore();

const keyword = ref(typeof route.query.keyword === "string" ? route.query.keyword : "");

const navs = [
  { to: "/", label: "首页" },
  { to: "/category", label: "分类" },
  { to: "/goods", label: "商品" },
  { to: "/cart", label: "购物车" },
  { to: "/user", label: "个人中心" },
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
</script>

<template>
  <div class="shell">
    <header class="topbar">
      <div class="pc-container topbar-inner">
        <RouterLink class="brand" to="/">
          <span class="brand-mark">栖</span>
          <span class="brand-text">栖息商城</span>
          <span class="brand-tag">PC</span>
        </RouterLink>

        <form class="search" @submit.prevent="search">
          <input
            v-model="keyword"
            type="search"
            placeholder="搜索商品 / 店铺"
            aria-label="搜索商品或店铺"
          />
          <button class="pc-btn" type="submit">搜索</button>
        </form>

        <div class="account">
          <template v-if="user.isLogin">
            <RouterLink class="account-link" to="/user">{{ user.displayName }}</RouterLink>
            <button class="pc-btn ghost" type="button" @click="logout">退出</button>
          </template>
          <template v-else>
            <button class="pc-btn ghost" type="button" @click="goLogin">登录 / 注册</button>
          </template>
        </div>
      </div>
      <nav class="nav">
        <div class="pc-container nav-inner">
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
        <p>栖息商城 PC · 对齐功能表 4 · API `/api/app/v1`</p>
        <p class="muted">阶段 PC-0/1 骨架 · 与 H5 共用 api-app</p>
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

.topbar {
  position: sticky;
  top: 0;
  z-index: 20;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.92);
  border-bottom: 1px solid var(--pc-line);
}

.topbar-inner {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 1.25rem;
  align-items: center;
  padding: 1rem 0 0.75rem;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 0.55rem;
}

.brand-mark {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  background: linear-gradient(145deg, #127a69, #0b5348);
  color: #fff;
  font-weight: 700;
}

.brand-text {
  font-size: 1.2rem;
  font-weight: 700;
  letter-spacing: 0.02em;
}

.brand-tag {
  font-size: 0.72rem;
  color: var(--pc-brand);
  background: var(--pc-brand-soft);
  padding: 0.15rem 0.45rem;
  border-radius: 999px;
}

.search {
  display: flex;
  gap: 0.5rem;
  max-width: 560px;
  width: 100%;
  justify-self: center;
}

.search input {
  flex: 1;
  border: 1px solid var(--pc-line);
  border-radius: 8px;
  padding: 0.65rem 0.9rem;
  background: #fff;
  outline: none;
}

.search input:focus {
  border-color: rgba(15, 107, 92, 0.55);
  box-shadow: 0 0 0 3px rgba(15, 107, 92, 0.12);
}

.account {
  display: flex;
  align-items: center;
  gap: 0.65rem;
}

.account-link {
  color: var(--pc-brand);
  font-weight: 600;
}

.nav {
  border-top: 1px solid rgba(228, 233, 239, 0.8);
}

.nav-inner {
  display: flex;
  gap: 1.25rem;
  padding: 0.55rem 0 0.7rem;
}

.nav-link {
  color: var(--pc-muted);
  padding: 0.25rem 0;
  border-bottom: 2px solid transparent;
}

.nav-link.active,
.nav-link:hover {
  color: var(--pc-brand);
  border-bottom-color: var(--pc-brand);
}

.main {
  padding: 1.5rem 0 2.5rem;
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
  .topbar-inner {
    grid-template-columns: 1fr;
  }

  .search {
    max-width: none;
  }

  .nav-inner {
    overflow-x: auto;
  }
}
</style>
