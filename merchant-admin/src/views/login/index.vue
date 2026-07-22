<template>
  <div class="login-page">
    <div class="panel">
      <div class="hero">
        <p class="eyebrow">MERCHANT CONSOLE</p>
        <h1>商户经营后台</h1>
        <p class="desc">本店商品 · 订单履约 · 营销财务 · 强制 mer_id 隔离</p>
      </div>
      <a-form class="form" layout="vertical" @finish="onSubmit">
        <h2>商户登录</h2>
        <a-form-item label="账号" name="account" :rules="[{ required: true, message: '请输入账号' }]">
          <a-input v-model:value="form.account" size="large" placeholder="meradmin" autocomplete="username" />
        </a-form-item>
        <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password
            v-model:value="form.password"
            size="large"
            placeholder="密码"
            autocomplete="current-password"
          />
        </a-form-item>
        <a-button type="primary" html-type="submit" size="large" block :loading="loading">进入店铺</a-button>
        <p class="hint">本地种子：meradmin / admin123（演示商户）</p>
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
const form = reactive({ account: 'meradmin', password: '' });

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
    radial-gradient(circle at 20% 15%, rgba(31, 111, 139, 0.2), transparent 42%),
    radial-gradient(circle at 80% 80%, rgba(20, 33, 43, 0.1), transparent 40%),
    linear-gradient(160deg, #eef3f6 0%, #f8f9fa 50%, #e7eef2 100%);
}
.panel {
  width: min(920px, 100%);
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(20, 33, 43, 0.08);
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 24px 60px rgba(20, 33, 43, 0.08);
}
.hero {
  padding: 48px 40px;
  background: linear-gradient(145deg, #1f6f8b, #124556);
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
  font-size: 36px;
  font-weight: 700;
}
.desc {
  margin: 18px 0 0;
  max-width: 280px;
  line-height: 1.7;
  opacity: 0.92;
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
}
</style>
