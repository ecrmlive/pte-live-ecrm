<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  APP_MANAGE_CODES,
  getAppStubConfigApi,
  saveAppStubConfigApi,
  type AppConfigKey,
  type AppStubConfig,
} from '#/api/core/platform-app-config';

const props = defineProps<{ configKey: AppConfigKey; title: string }>();

const note = ref('');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<AppStubConfig>({ name: '', enabled: false, remark: '' });

async function load() {
  loading.value = true;
  try {
    const data = await getAppStubConfigApi(props.configKey);
    note.value = data.note;
    form.value = data.config;
  } finally {
    loading.value = false;
  }
}

async function save() {
  saving.value = true;
  try {
    form.value = await saveAppStubConfigApi(props.configKey, {
      ...form.value,
      name: form.value.name.trim(),
      remark: form.value.remark.trim(),
    });
    ElMessage.success(`${props.title}已保存`);
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes(APP_MANAGE_CODES[props.configKey]);
});
</script>

<template>
  <Page :title="title" description="仅维护名称与启用开关；AppSecret、Token 等密钥请通过云服务配置维护。">
    <el-card v-loading="loading" shadow="never">
      <el-alert :title="note" type="warning" :closable="false" class="mb-4" />
      <el-form :disabled="!canManage" label-width="160px" class="max-w-3xl">
        <el-form-item label="名称"><el-input v-model="form.name" maxlength="64" show-word-limit /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="form.enabled" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" :rows="4" maxlength="500" show-word-limit type="textarea" /></el-form-item>
      </el-form>
      <div class="mt-4 flex justify-center gap-3"><el-button @click="load">重置</el-button><el-button v-if="canManage" :loading="saving" type="primary" @click="save">保存</el-button></div>
    </el-card>
  </Page>
</template>
