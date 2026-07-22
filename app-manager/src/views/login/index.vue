<template>
  <div class="login">
    <div class="panel">
      <h1>店员核销</h1>
      <p class="desc">登录后核销本店订单 / 代退</p>
      <a-form layout="vertical" @finish="onSubmit">
        <a-form-item label="账号" name="account" :rules="[{ required: true }]">
          <a-input v-model:value="form.account" size="large" placeholder="staff1" />
        </a-form-item>
        <a-form-item label="密码" name="password" :rules="[{ required: true }]">
          <a-input-password v-model:value="form.password" size="large" />
        </a-form-item>
        <a-button type="primary" html-type="submit" size="large" block :loading="loading">登录</a-button>
      </a-form>
      <p class="hint">种子：staff1 / admin123（mer_id=1）</p>
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
const form = reactive({ account: 'staff1', password: '' });

async function onSubmit() {
  loading.value = true;
  try {
    await auth.login(form.account, form.password);
    message.success('登录成功');
    router.replace((route.query.redirect as string) || '/verify');
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  background: linear-gradient(160deg, #1f2a37, #3d4f66);
}
.panel {
  width: min(400px, 100%);
  background: #fff;
  border-radius: 16px;
  padding: 28px;
}
h1 {
  margin: 0 0 8px;
  font-size: 24px;
}
.desc,
.hint {
  color: #888;
  font-size: 13px;
}
.hint {
  margin-top: 16px;
}
</style>
