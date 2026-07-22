<template>
  <div class="login-page">
    <div class="panel">
      <div class="hero">
        <p class="eyebrow">QIXI MERGERS</p>
        <h1>栖息多商户</h1>
        <p class="desc">平台运营后台 · 商户审核 · 全局监管 · 结算权限</p>
      </div>
      <a-form class="form" layout="vertical" @finish="onSubmit">
        <h2>登录</h2>
        <a-form-item label="账号" name="account" :rules="[{ required: true, message: '请输入账号' }]">
          <a-input v-model:value="form.account" size="large" placeholder="admin" autocomplete="username" />
        </a-form-item>
        <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password
            v-model:value="form.password"
            size="large"
            placeholder="密码"
            autocomplete="current-password"
          />
        </a-form-item>
        <a-button type="primary" html-type="submit" size="large" block :loading="loading">进入后台</a-button>
        <p class="hint">本地种子：admin / admin123</p>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { message } from 'ant-design-vue';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();
const loading = ref(false);
const form = reactive({ account: 'admin', password: '' });

async function onSubmit() {
  loading.value = true;
  try {
    await auth.login(form.account, form.password);
    message.success('登录成功');
    const redirect = (route.query.redirect as string) || '/dashboard';
    router.replace(redirect);
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  background:
    radial-gradient(circle at 15% 20%, rgba(15, 110, 86, 0.18), transparent 40%),
    radial-gradient(circle at 85% 10%, rgba(20, 33, 43, 0.12), transparent 35%),
    linear-gradient(160deg, #edf3f0 0%, #f7f8f6 45%, #e8eef2 100%);
}
.panel {
  width: min(920px, 100%);
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid rgba(20, 33, 43, 0.08);
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 24px 60px rgba(20, 33, 43, 0.08);
}
.hero {
  padding: 48px 40px;
  background:
    linear-gradient(145deg, rgba(15, 110, 86, 0.92), rgba(11, 79, 61, 0.95)),
    url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.06'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  color: #fff;
}
.eyebrow {
  letter-spacing: 0.18em;
  font-size: 12px;
  opacity: 0.8;
  margin: 0 0 16px;
}
.hero h1 {
  margin: 0;
  font-size: 40px;
  font-weight: 700;
  letter-spacing: 0.04em;
}
.desc {
  margin: 18px 0 0;
  max-width: 280px;
  line-height: 1.7;
  opacity: 0.9;
}
.form {
  padding: 48px 40px;
}
.form h2 {
  margin: 0 0 24px;
  font-size: 24px;
}
.hint {
  margin-top: 16px;
  color: #6b7785;
  font-size: 13px;
}
@media (max-width: 800px) {
  .panel {
    grid-template-columns: 1fr;
  }
  .hero {
    padding: 32px 24px;
  }
  .form {
    padding: 32px 24px;
  }
}
</style>
