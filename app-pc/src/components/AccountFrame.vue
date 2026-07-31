<script setup lang="ts">
import { computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";

const route = useRoute();
const router = useRouter();
const user = useUserStore();

const phoneText = computed(() => {
  const value = user.profile?.phone || user.profile?.account || "";
  if (value.length < 7) return value || "未登录";
  return `${value.slice(0, 3)}****${value.slice(-4)}`;
});

const menus = [
  { label: "账户管理", to: "/user", match: (path: string) => path === "/user" },
  { label: "我的订单", to: "/orders", match: (path: string) => path.startsWith("/orders") },
  { label: "预售订单", to: "/presell", match: (path: string) => path.startsWith("/presell") },
  { label: "我的积分", to: "/points", match: (path: string) => path.startsWith("/points") },
  { label: "我的余额", to: "/user/balance", match: (path: string) => path.startsWith("/user/balance") },
  { label: "我的优惠券", to: "/coupons", match: (path: string) => path.startsWith("/coupons") },
  { label: "我的发票", to: "/user/invoices", match: (path: string) => path.startsWith("/user/invoices") },
  { label: "地址管理", to: "/user/addresses", match: (path: string) => path.startsWith("/user/addresses") },
];

onMounted(() => {
  if (user.isLogin) void user.refreshMe();
});

function goLogin() {
  router.push({ name: "login", query: { redirect: route.fullPath } });
}
</script>

<template>
  <div class="account-shell pc-container">
    <div class="account-crumb">首页 <span>›</span> 个人中心 <slot name="crumb" /></div>
    <div class="account-grid">
      <aside class="account-sidebar">
        <div class="account-profile">
          <div class="account-avatar">{{ user.displayName.slice(0, 1) }}</div>
          <strong>{{ phoneText }}</strong>
        </div>
        <div class="account-nav">
          <RouterLink
            v-for="item in menus"
            :key="item.to"
            :to="item.to"
            :class="{ active: item.match(route.path) }"
          >
            {{ item.label }}
          </RouterLink>
        </div>
      </aside>
      <section v-if="user.isLogin" class="account-content"><slot /></section>
      <section v-else class="account-content account-login">
        <h1>个人中心</h1>
        <p>登录后即可管理订单、优惠券和收货地址。</p>
        <button class="pc-btn" type="button" @click="goLogin">去登录</button>
      </section>
    </div>
  </div>
</template>

<style scoped>
.account-shell { padding-top: 28px; padding-bottom: 52px; }
.account-crumb { margin-bottom: 22px; color: #555; font-size: 14px; }
.account-crumb span { padding: 0 5px; color: #aaa; }
.account-grid { display: grid; grid-template-columns: 188px minmax(0, 1fr); gap: 10px; align-items: stretch; }
.account-sidebar, .account-content { background: #fff; }
.account-sidebar { min-height: 690px; }
.account-profile { display: grid; justify-items: center; gap: 14px; padding: 46px 16px 34px; border-bottom: 10px solid #f8f8f8; font-size: 14px; }
.account-avatar { width: 72px; height: 72px; display: grid; place-items: center; border-radius: 50%; background: linear-gradient(145deg, #ffd2be, #f5f5f5 48%, #66b354 49%); color: #674b3d; font-size: 28px; font-weight: 700; }
.account-nav { padding: 24px 0 38px; }
.account-nav a { position: relative; display: block; padding: 15px 0 15px 62px; color: #666; font-size: 15px; }
.account-nav a:hover, .account-nav a.active { color: #f13728; }
.account-nav a.active::before { position: absolute; left: 0; top: 13px; width: 4px; height: 26px; background: #f13728; content: ""; }
.account-content { min-height: 690px; padding: 24px 30px 42px; }
.account-login { display: grid; align-content: center; justify-items: center; min-height: 460px; color: #777; }
.account-login h1 { margin: 0; color: #2d2d2d; }.account-login p { margin: 12px 0 22px; }
@media (max-width: 760px) { .account-grid { grid-template-columns: 1fr; }.account-sidebar { min-height: 0; }.account-profile { display: none; }.account-nav { display: flex; overflow-x: auto; padding: 0; }.account-nav a { white-space: nowrap; padding: 15px 16px; }.account-nav a.active::before { left: 50%; top: auto; bottom: 0; width: 24px; height: 2px; transform: translateX(-50%); }.account-content { padding: 22px 18px; } }
</style>
