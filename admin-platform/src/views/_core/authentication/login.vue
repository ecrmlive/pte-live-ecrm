<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';

import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
} from 'element-plus';

import { getAdminLoginBaseApi } from '#/api/core/passport';
import { QIXI_PLATFORM_LOGIN_NAME } from '#/preferences';
import { useAuthStore } from '#/store';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();

const formRef = ref<InstanceType<typeof ElForm>>();
const codeImage = ref('');
const codeKey = ref('');
const baseLoading = ref(false);
const baseLoadFailed = ref(false);

const form = reactive({
  username: '',
  password: '',
  code: '',
});

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
};

async function loadLoginBase() {
  baseLoading.value = true;
  baseLoadFailed.value = false;
  try {
    const data = await getAdminLoginBaseApi();
    codeImage.value = data?.codeData?.codeImage || '';
    codeKey.value = data?.codeData?.codeKey || '';
    if (!codeImage.value) {
      form.code = '';
    }
  } catch {
    codeImage.value = '';
    codeKey.value = '';
    baseLoadFailed.value = true;
  } finally {
    baseLoading.value = false;
  }
}

async function handleSubmit() {
  if (!formRef.value || authStore.loginLoading) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    try {
      await authStore.authLogin({
        username: form.username,
        password: form.password,
        code: codeImage.value ? form.code : undefined,
        codeKey: codeImage.value ? codeKey.value : undefined,
      });
    } catch {
      if (codeImage.value) {
        await loadLoginBase();
      }
    }
  });
}

onMounted(() => {
  loadLoginBase();
});
</script>

<template>
  <div class="qixi-live-login-form">
    <h2 class="qixi-live-login-form__title">{{ QIXI_PLATFORM_LOGIN_NAME }}</h2>
    <p class="qixi-live-login-form__subtitle">平台、商户、区域、客服、运营账号统一登录</p>

    <ElAlert
      v-if="baseLoadFailed"
      :closable="false"
      class="qixi-live-login-form__alert"
      show-icon
      title="内部错误，请稍后重试"
      type="error"
    />

    <ElForm
      ref="formRef"
      v-loading="baseLoading"
      :model="form"
      :rules="
        codeImage
          ? rules
          : { username: rules.username, password: rules.password }
      "
      label-position="top"
      size="large"
      @keyup.enter="handleSubmit"
    >
      <ElFormItem label="账号" prop="username">
        <ElInput
          v-model="form.username"
          autocomplete="username"
          placeholder="账号"
        />
      </ElFormItem>
      <ElFormItem label="密码" prop="password">
        <ElInput
          v-model="form.password"
          autocomplete="current-password"
          placeholder="密码"
          show-password
          type="password"
        />
      </ElFormItem>
      <ElFormItem v-if="codeImage" label="验证码" prop="code">
        <div class="qixi-live-login-form__captcha">
          <ElInput v-model="form.code" placeholder="验证码" />
          <img
            :src="codeImage"
            alt="验证码"
            class="qixi-live-login-form__captcha-img"
            title="点击刷新验证码"
            @click="loadLoginBase"
          />
        </div>
      </ElFormItem>
      <ElButton
        :loading="authStore.loginLoading"
        class="qixi-live-login-form__submit"
        type="primary"
        @click="handleSubmit"
      >
        登录
      </ElButton>
    </ElForm>
  </div>
</template>

<style scoped>
.qixi-live-login-form {
  width: 100%;
  max-width: 420px;
}

.qixi-live-login-form__title {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.qixi-live-login-form__subtitle {
  margin: 0 0 24px;
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.qixi-live-login-form__alert {
  margin-bottom: 16px;
}

.qixi-live-login-form__captcha {
  display: flex;
  gap: 12px;
  width: 100%;
}

.qixi-live-login-form__captcha-img {
  height: 40px;
  cursor: pointer;
  border-radius: 4px;
  flex-shrink: 0;
}

.qixi-live-login-form__submit {
  width: 100%;
  margin-top: 8px;
}
</style>
