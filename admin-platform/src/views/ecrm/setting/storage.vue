<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformStorageConfigApi,
  savePlatformStorageConfigApi,
  type PlatformStorageConfig,
} from '#/api/core/platform-setting-ext';

const note = ref('');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<PlatformStorageConfig>({
  provider: 'cos',
  region: '',
  bucket_name: '',
  enabled: false,
  remark: '',
});

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformStorageConfigApi();
    note.value = data.note;
    form.value = data.config;
  } finally {
    loading.value = false;
  }
}

async function save() {
  saving.value = true;
  try {
    form.value = await savePlatformStorageConfigApi({
      ...form.value,
      provider: form.value.provider.trim(),
      region: form.value.region.trim(),
      bucket_name: form.value.bucket_name.trim(),
      remark: form.value.remark.trim(),
    });
    ElMessage.success('存储设置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes('setting.storage.manage');
});
</script>

<template>
  <Page title="存储配置" description="仅维护对象存储开关、区域与桶名展示；SecretId/SecretKey 请通过云服务配置维护。">
    <el-card v-loading="loading" shadow="never">
      <el-alert :title="note" type="warning" :closable="false" class="mb-4" />
      <el-form :disabled="!canManage" label-width="160px" class="max-w-3xl">
        <el-form-item label="存储提供商"><el-select v-model="form.provider" class="w-full"><el-option label="腾讯云 COS" value="cos" /></el-select></el-form-item>
        <el-form-item label="区域"><el-input v-model="form.region" maxlength="64" placeholder="如 ap-guangzhou" /></el-form-item>
        <el-form-item label="桶名称"><el-input v-model="form.bucket_name" maxlength="128" placeholder="仅展示，不含密钥" /></el-form-item>
        <el-form-item label="启用存储"><el-switch v-model="form.enabled" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" :rows="4" maxlength="500" show-word-limit type="textarea" /></el-form-item>
      </el-form>
      <div class="mt-4 flex justify-end"><el-button @click="load">重置</el-button><el-button v-if="canManage" :loading="saving" type="primary" @click="save">保存设置</el-button></div>
    </el-card>
  </Page>
</template>
