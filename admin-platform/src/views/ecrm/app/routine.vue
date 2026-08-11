<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElCard,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElSwitch,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  getRoutineConfigApi,
  saveRoutineConfigApi,
} from '#/api/core/cloud-config';

const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const secretConfigured = ref(false);
const form = reactive({
  appId: '',
  appSecret: '',
  enabled: false,
});

function toEnabled(value?: string) {
  return value === '1' || value === 'true';
}

async function load() {
  loading.value = true;
  try {
    const data = await getRoutineConfigApi();
    form.enabled = toEnabled(data.values.enabled);
    form.appId = data.values.app_id || '';
    // 接口从不返回密钥明文；清空输入框后保存不会覆盖已配置的 AppSecret。
    form.appSecret = '';
    secretConfigured.value = Boolean(data.values.app_secret);
  } finally {
    loading.value = false;
  }
}

async function save() {
  const appId = form.appId.trim();
  const appSecret = form.appSecret.trim();
  if (!appId) {
    ElMessage.warning('请填写 AppID');
    return;
  }
  if (!appSecret && !secretConfigured.value) {
    ElMessage.warning('请填写 AppSecret');
    return;
  }
  saving.value = true;
  try {
    const data = await saveRoutineConfigApi({
      app_id: appId,
      app_secret: appSecret,
      enabled: form.enabled ? 'true' : 'false',
    });
    form.appId = data.values.app_id || appId;
    form.appSecret = '';
    secretConfigured.value = Boolean(data.values.app_secret);
    ElMessage.success('小程序配置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canManage.value = codes.includes('app.routine.manage');
  if (canManage.value) {
    await load();
  }
});
</script>

<template>
  <Page auto-content-height>
    <ElCard v-loading="loading" shadow="never" class="routine-config">
      <ElAlert
        title="AppSecret 采用加密存储，保存后不会回显明文；如不需要更换，请保持 AppSecret 输入框为空。"
        type="warning"
        :closable="false"
        class="mb-5"
      />

      <ElForm
        label-width="120px"
        class="max-w-3xl"
        :disabled="!canManage"
        @submit.prevent
      >
        <ElFormItem label="启用小程序">
          <ElSwitch v-model="form.enabled" />
        </ElFormItem>
        <ElFormItem label="AppID" required>
          <ElInput
            v-model="form.appId"
            maxlength="64"
            placeholder="请输入微信小程序 AppID"
          />
        </ElFormItem>
        <ElFormItem label="AppSecret" required>
          <ElInput
            v-model="form.appSecret"
            maxlength="128"
            show-password
            type="password"
            :placeholder="
              secretConfigured
                ? '已配置；如需更换请输入新的 AppSecret'
                : '请输入微信小程序 AppSecret'
            "
          />
          <p class="routine-config__hint">
            {{ secretConfigured ? 'AppSecret 已配置，留空即可保留原值。' : '保存后仅显示已配置状态。' }}
          </p>
        </ElFormItem>
      </ElForm>

      <div class="routine-config__actions">
        <ElButton @click="load">重置</ElButton>
        <ElButton
          v-if="canManage"
          type="primary"
          :loading="saving"
          @click="save"
        >
          保存
        </ElButton>
      </div>
    </ElCard>
  </Page>
</template>

<style scoped>
.routine-config {
  max-width: 960px;
}

.routine-config__hint {
  width: 100%;
  margin: 6px 0 0;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 1.5;
}

.routine-config__actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  max-width: 768px;
  margin-top: 28px;
}
</style>
