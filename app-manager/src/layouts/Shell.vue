<template>
  <div class="shell">
    <header class="top">
      <div>
        <strong>店员端</strong>
        <small>{{ auth.user?.mer_name || `mer ${auth.user?.mer_id || '-'}` }} · {{ auth.user?.nickname }}</small>
      </div>
      <a-button type="link" danger @click="onLogout">退出</a-button>
    </header>
    <nav class="nav">
      <router-link v-if="auth.user?.is_verify === 1" to="/verify">待核销</router-link>
      <router-link v-if="auth.user?.is_goods === 1" to="/delivery">待发货</router-link>
      <router-link to="/refund">代退</router-link>
    </nav>
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const router = useRouter();

function onLogout() {
  auth.logout();
  router.replace('/login');
}
</script>

<style scoped>
.shell {
  min-height: 100vh;
  background: #f5f5f5;
}
.top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fff;
  border-bottom: 1px solid #eee;
  min-height: 44px;
}
.top small {
  display: block;
  color: #888;
  margin-top: 2px;
}
.nav {
  display: flex;
  gap: 16px;
  padding: 10px 16px;
  background: #fff;
  border-bottom: 1px solid #f0f0f0;
}
.nav a {
  color: #666;
  text-decoration: none;
}
.nav a.router-link-active {
  color: #1677ff;
  font-weight: 600;
}
.main {
  padding: 12px;
}
</style>
