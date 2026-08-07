<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { ElButton, ElForm, ElFormItem, ElInput } from 'element-plus';

import { getShopLoginBaseApi, saasLoginApi } from '#/api';
import { $t } from '#/locales';
import { useAuthStore } from '#/store';
import { applyLoginPageBranding } from '#/utils/shop-display-name';
import { loadShopBootstrapData } from '#/utils/shop-bootstrap';
import { useAccessStore, useUserStore } from '@vben/stores';
import { preferences } from '@vben/preferences';
import {
  getDecryptedToken,
  hydrateAccessTokenFromLegacy,
  setEncryptedToken,
} from '#/utils/pte-live-token';
import {
  markJwtIssuedFromLogin,
  markShopSessionBootstrapped,
} from '#/utils/jwt-session';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();
const accessStore = useAccessStore();
const userStore = useUserStore();
const route = useRoute();
const router = useRouter();

const formRef = ref<InstanceType<typeof ElForm>>();
const codeImage = ref('');
const codeKey = ref('');

const form = reactive({
  username: '',
  password: '',
  code: '',
});

const rules = {
  username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
};

async function loadLoginBase() {
  try {
    const data = await getShopLoginBaseApi();
    applyLoginPageBranding(data?.settings?.shop_logo_img);
    codeImage.value = data?.codeData?.codeImage || '';
    codeKey.value = data?.codeData?.codeKey || '';
    if (!codeImage.value) {
      form.code = '';
    }
  } catch {
    codeImage.value = '';
    codeKey.value = '';
    applyLoginPageBranding();
  }
}

async function handleSubmit() {
  if (!formRef.value || authStore.loginLoading) return;
  await formRef.value.validate(async (valid) => {
    if (!valid) return;
    try {
      await authStore.authLogin({
        code: codeImage.value ? form.code : undefined,
        codeKey: codeImage.value ? codeKey.value : undefined,
        password: form.password,
        username: form.username,
      });
    } catch {
      if (codeImage.value) {
        await loadLoginBase();
      }
    }
  });
}

onMounted(() => {
  applyLoginPageBranding();
  void loadLoginBase();
  void trySaasAutoLogin();
});

async function trySaasAutoLogin() {
  if (route.query.from !== 'admin') {
    return;
  }
  // 跨域新标签：平台通过 URL ?token= 下发 store_console access JWT。
  const queryToken = String(route.query.token || '').trim();
  if (queryToken) {
    markJwtIssuedFromLogin();
    accessStore.setAccessToken(queryToken);
    setEncryptedToken(queryToken);
    // 立刻去掉 URL 中的 token，避免历史/Referer 泄露。
    const nextQuery = { ...route.query };
    delete nextQuery.token;
    await router.replace({ path: route.path, query: nextQuery });
  } else {
    hydrateAccessTokenFromLegacy((token) => {
      accessStore.setAccessToken(token);
      setEncryptedToken(token);
    });
  }
  const token = accessStore.accessToken || getDecryptedToken();
  if (!token) {
    return;
  }
  try {
    await saasLoginApi();
    const bootstrap = await loadShopBootstrapData();
    userStore.setUserInfo(bootstrap.userInfo);
    accessStore.setAccessCodes(bootstrap.accessCodes);
    markShopSessionBootstrapped();
    accessStore.setIsAccessChecked(false);
    const redirectQuery = route.query.redirect as string | undefined;
    const target = redirectQuery
      ? decodeURIComponent(redirectQuery)
      : bootstrap.userInfo?.homePath || preferences.app.defaultHomePath;
    await router.replace(target);
  } catch {
    // 保留登录表单，由用户手动登录
  }
}
</script>

<template>
  <div class="pte-live-login-form">
    <h2 class="pte-live-login-form__title">{{ $t('authentication.pageTitle') }}</h2>
    <p class="pte-live-login-form__subtitle">{{ $t('authentication.loginSubtitle') }}</p>

    <ElForm
      ref="formRef"
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
        <div class="pte-live-login-form__captcha">
          <ElInput v-model="form.code" placeholder="验证码" />
          <img
            :src="codeImage"
            alt="验证码"
            class="pte-live-login-form__captcha-img"
            title="点击刷新验证码"
            @click="loadLoginBase"
          />
        </div>
      </ElFormItem>
      <ElButton
        :loading="authStore.loginLoading"
        class="pte-live-login-form__submit"
        type="primary"
        @click="handleSubmit"
      >
        登录
      </ElButton>
    </ElForm>
  </div>
</template>

<style scoped>
.pte-live-login-form {
  width: 100%;
  max-width: 420px;
}

.pte-live-login-form__title {
  margin: 0 0 8px;
  font-size: 24px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.pte-live-login-form__subtitle {
  margin: 0 0 24px;
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.pte-live-login-form__captcha {
  display: flex;
  gap: 12px;
  width: 100%;
}

.pte-live-login-form__captcha :deep(.el-input) {
  flex: 1;
  min-width: 0;
}

.pte-live-login-form__captcha-img {
  height: 40px;
  cursor: pointer;
  border-radius: 4px;
  flex-shrink: 0;
}

.pte-live-login-form__submit {
  width: 100%;
  margin-top: 8px;
}
</style>
