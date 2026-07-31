<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";

const router = useRouter();
const route = useRoute();
const user = useUserStore();

const keyword = ref(typeof route.query.keyword === "string" ? route.query.keyword : "");
const popularKeywords = ["海鲜", "预制菜", "国潮", "箱包", "家居", "口红", "iPhone"];

type NavItem = { to: string; label: string; activePath?: string };

const navs: NavItem[] = [
  { to: "/goods", label: "全部商品" },
  { to: "/coupons?view=center", activePath: "/coupons", label: "领券中心" },
  { to: "/notices", label: "新闻中心" },
  { to: "/seckill", label: "秒杀列表" },
  { to: "/merchant/apply", label: "商户入驻" },
  { to: "/user", label: "我的信息" },
];

const serviceLinks = [
  { label: "精品推荐", to: "/goods?type=best" },
  { label: "火爆新品", to: "/goods?type=new" },
  { label: "分类广场", to: "/category" },
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
            <RouterLink to="/user/balance">我的余额</RouterLink>
            <RouterLink to="/merchant/apply">店铺入驻</RouterLink>
            <RouterLink to="/notices">资讯信息</RouterLink>
          </div>
        </div>
      </div>
      <div class="pc-container masthead">
        <RouterLink class="brand" to="/" aria-label="七禧商城首页">
          <img class="brand-mark" src="/brand/qixi-logo.png" alt="七禧商城" />
          <span class="brand-name"><strong>七禧</strong><em>商城</em></span>
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
            :class="{ active: item.to === '/' ? activePath === '/' : activePath.startsWith(item.activePath || item.to) }"
          >
            {{ item.label }}
          </RouterLink>
        </div>
      </nav>
    </header>

    <main class="main">
      <RouterView />
    </main>

    <aside class="right-rail" aria-label="快捷服务">
      <RouterLink v-for="item in serviceLinks" :key="item.label" :to="item.to">{{ item.label }}</RouterLink>
      <RouterLink class="right-rail__service" to="/live">在线客服</RouterLink>
    </aside>

    <footer class="footer">
      <div class="pc-container footer-inner">
        <div class="assurance-list" aria-label="商城服务保障">
          <div><b>品质商品</b><span>严选商品，购物放心</span></div>
          <div><b>多仓直发</b><span>覆盖全国，快速送达</span></div>
          <div><b>正品保障</b><span>正规渠道，售后无忧</span></div>
          <div><b>全天客服</b><span>在线服务，及时响应</span></div>
        </div>
        <p>Copyright © 七禧商城 All Rights Reserved</p>
        <p class="muted">七禧多商户商城</p>
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

.site-header {
  position: fixed;
  right: 0;
  left: 0;
  top: 0;
  z-index: 30;
  background: #fff;
  box-shadow: 0 2px 10px rgb(0 0 0 / 6%);
}

.utility-bar { background: #2e2e2e; color: #c9c9c9; font-size: .78rem; }

.utility-inner { height: 40px; display: flex; align-items: center; justify-content: space-between; }

.utility-links { display: flex; gap: 1rem; align-items: center; }

.utility-links a, .utility-links button { color: inherit; border: 0; padding: 0; background: transparent; font: inherit; cursor: pointer; }

.utility-links a:hover, .utility-links button:hover { color: #fff; }

.masthead {
  display: grid; grid-template-columns: 250px minmax(360px, 730px) 150px;
  min-height: 160px; align-items: center; justify-content: space-between; gap: 1.55rem;
}

.brand { display: flex; align-items: center; gap: 10px; color: #151515; white-space: nowrap; }
.brand-mark { width: 58px; height: 58px; border-radius: 13px; object-fit: cover; box-shadow: 0 5px 12px rgb(10 54 103 / 14%); }
.brand-name { display: flex; align-items: baseline; gap: 5px; letter-spacing: -.07em; }
.brand-name strong { color: #242424; font-size: 2.25rem; font-weight: 800; }
.brand-name em { color: #f13728; font-size: .88rem; font-style: normal; font-weight: 700; letter-spacing: 0; }

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
  height: 48px;
  padding: 0.85rem 1rem;
  background: #fff;
  outline: none;
}

.search input:focus {
  border-color: #f13728;
}

.search > button { width: 104px; border: 0; background: #f13728; color: #fff; font-weight: 700; cursor: pointer; }
.popular-keywords { position: absolute; top: 56px; left: .75rem; display: flex; gap: 1.1rem; white-space: nowrap; }
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
  padding: 256px 0 2.5rem;
}

.footer {
  border-top: 1px solid var(--pc-line);
  background: rgba(255, 255, 255, 0.7);
}

.footer-inner {
  padding: 1.75rem 0 1.5rem;
  text-align: center;
}

.footer p {
  margin: 0;
}

.footer .muted {
  margin-top: 0.35rem;
  color: var(--pc-muted);
  font-size: 0.9rem;
}

.assurance-list { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1rem; padding: 0 0 1.65rem; margin-bottom: 1.25rem; border-bottom: 1px solid #eee; }
.assurance-list div { display: grid; gap: .2rem; text-align: left; }
.assurance-list b { color: #555; font-size: .9rem; }.assurance-list span { margin-top: .2rem; color: #999; font-size: .75rem; }

.right-rail { position: fixed; z-index: 8; right: max(12px, calc((100vw - var(--pc-max)) / 2 - 78px)); top: 52%; display: grid; width: 58px; overflow: hidden; border: 1px solid #eee; background: #fff; box-shadow: 0 4px 18px rgb(0 0 0 / 5%); transform: translateY(-50%); }
.right-rail a { display: grid; min-height: 62px; place-items: center; padding: 8px 7px; border-bottom: 1px solid #f0f0f0; color: #666; font-size: 12px; line-height: 18px; text-align: center; }
.right-rail a:last-child { border-bottom: 0; }.right-rail a:hover { color: #f13728; background: #fff7f6; }.right-rail__service { color: #f13728 !important; font-weight: 600; }

@media (max-width: 860px) {
  .masthead {
    grid-template-columns: 1fr;
    padding: 1.5rem 0;
    gap: 1rem;
  }

  .main { padding-top: 358px; }

  .search {
    max-width: none;
  }

  .nav-inner {
    overflow-x: auto;
  }

  .right-rail { display: none; }
  .assurance-list { grid-template-columns: repeat(2, 1fr); }
}
</style>
