<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import AccountFrame from "@/components/AccountFrame.vue";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const router = useRouter();
const maskedPhone = computed(() => {
  const value = user.profile?.phone || user.profile?.account || "";
  return value.length >= 7 ? `${value.slice(0, 3)}****${value.slice(-4)}` : value || "未绑定";
});

function logout() {
  user.logout();
  router.push({ name: "home" });
}
</script>

<template>
  <AccountFrame>
    <template #crumb><span>›</span> 账户管理</template>
    <h1 class="content-title">账户管理</h1>
    <section class="profile-section">
      <h2>我的信息</h2>
      <dl>
        <div><dt>我的头像：</dt><dd><span class="large-avatar">{{ user.displayName.slice(0, 1) }}</span></dd></div>
        <div><dt>我的昵称：</dt><dd>{{ user.displayName }}</dd></div>
        <div><dt>我的 ID：</dt><dd>{{ user.profile?.uid || "-" }}</dd></div>
        <div><dt>手机号：</dt><dd>{{ maskedPhone }}</dd></div>
        <div><dt>登录账号：</dt><dd>{{ user.profile?.account || "-" }}</dd></div>
      </dl>
      <div class="account-actions">
        <button class="pc-btn" type="button" @click="logout">退出登录</button>
      </div>
    </section>
  </AccountFrame>
</template>

<style scoped>
.content-title { margin: 0; padding-bottom: 20px; border-bottom: 1px solid #eee; font-size: 20px; }.profile-section { padding: 34px 16px 0; }.profile-section h2 { margin: 0 0 20px; font-size: 18px; }.profile-section dl { margin: 0; }.profile-section dl > div { display: grid; grid-template-columns: 112px 1fr; align-items: center; min-height: 73px; border-bottom: 1px dashed #e8e8e8; }.profile-section dt { color: #888; }.profile-section dd { margin: 0; color: #444; }.large-avatar { display: grid; width: 80px; height: 80px; place-items: center; border-radius: 50%; background: linear-gradient(145deg, #ffd2be, #f5f5f5 48%, #66b354 49%); color: #674b3d; font-size: 30px; font-weight: 700; }.profile-section dl > div:first-child { min-height: 145px; }.account-actions { display: flex; justify-content: flex-end; padding-top: 38px; }.account-actions .pc-btn { min-width: 135px; background: #f13728; border-radius: 4px; }
</style>
