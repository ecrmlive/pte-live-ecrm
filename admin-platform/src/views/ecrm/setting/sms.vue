<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElCard, ElForm, ElFormItem, ElInput, ElMessage, ElSwitch } from 'element-plus';

import {
  getTencentSMSConfigApi,
  saveTencentSMSConfigApi,
} from '#/api/core/cloud-config';

const loading = ref(false);
const saving = ref(false);
const form = reactive({
  app_key: '',
  enabled: true,
  sdk_app_id: '',
  sign_content: '',
  sign_id: '',
  template_id: '',
});

async function load() {
  loading.value = true;
  try {
    const config = await getTencentSMSConfigApi();
    form.enabled = config.values.enabled !== 'false' && config.values.enabled !== '0';
    form.sdk_app_id = config.values.sdk_app_id || '';
    form.sign_id = config.values.sign_id || '';
    form.sign_content = config.values.sign_content || '';
    form.template_id = config.values.template_id || '';
    form.app_key = config.values.app_key === '********' ? '********' : '';
  } catch {
    ElMessage.error('加载短信配置失败');
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!form.sdk_app_id || !form.sign_id || !form.sign_content.trim() || !form.template_id) {
    ElMessage.warning('请填写完整的短信配置');
    return;
  }
  if (!form.app_key.trim()) {
    ElMessage.warning('请填写 App Key');
    return;
  }
  saving.value = true;
  try {
    await saveTencentSMSConfigApi({
      app_key: form.app_key.trim(),
      enabled: String(form.enabled),
      sdk_app_id: form.sdk_app_id.trim(),
      sign_content: form.sign_content.trim(),
      sign_id: form.sign_id.trim(),
      template_id: form.template_id.trim(),
    });
    ElMessage.success('短信配置已保存');
    await load();
  } catch {
    // 请求层会展示服务端返回的具体错误，避免覆盖成笼统提示。
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>

<template>
  <Page auto-content-height>
    <ElCard v-loading="loading" class="max-w-3xl" shadow="never">
      <ElForm label-width="130px">
        <ElFormItem label="启用短信验证码">
          <ElSwitch v-model="form.enabled" />
        </ElFormItem>
        <ElFormItem label="SDKAppID" required>
          <ElInput v-model="form.sdk_app_id" inputmode="numeric" placeholder="请输入 SDKAppID" />
        </ElFormItem>
        <ElFormItem label="App Key" required>
          <ElInput
            v-model="form.app_key"
            show-password
            type="password"
            placeholder="请输入 App Key"
          />
        </ElFormItem>
        <ElFormItem label="签名管理 ID" required>
          <ElInput v-model="form.sign_id" inputmode="numeric" placeholder="请输入签名管理 ID" />
        </ElFormItem>
        <ElFormItem label="签名内容" required>
          <ElInput v-model="form.sign_content" placeholder="请输入签名内容" />
        </ElFormItem>
        <ElFormItem label="模板 ID" required>
          <ElInput v-model="form.template_id" inputmode="numeric" placeholder="请输入模板 ID" />
        </ElFormItem>
        <div class="flex justify-center gap-3 pt-4">
          <ElButton @click="load">重置</ElButton>
          <ElButton :loading="saving" type="primary" @click="save">保存</ElButton>
        </div>
      </ElForm>
    </ElCard>
  </Page>
</template>
