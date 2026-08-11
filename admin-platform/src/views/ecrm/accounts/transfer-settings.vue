<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformTransferSettingsConfigApi,
  savePlatformTransferSettingsConfigApi,
  type PlatformTransferSettingsConfig,
} from '#/api/core/platform-setting-ext';

const note = ref('');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<PlatformTransferSettingsConfig>({ enabled: false, min_amount: 1, remark: '' });

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformTransferSettingsConfigApi();
    note.value = data.note;
    form.value = data.config;
  } finally {
    loading.value = false;
  }
}

async function save() {
  saving.value = true;
  try {
    form.value = await savePlatformTransferSettingsConfigApi({ ...form.value, remark: form.value.remark.trim() });
    ElMessage.success('转账设置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes('accounts.transfer_settings.manage');
});
</script>

<template>
  <Page title="转账设置" description="维护转账监管开关与最低金额；真实打款凭据请通过云服务配置维护。">
    <el-card v-loading="loading" shadow="never">
      <el-alert :title="note" type="warning" :closable="false" class="mb-4" />
      <el-form :disabled="!canManage" label-width="160px" class="max-w-3xl">
        <el-form-item label="启用转账监管"><el-switch v-model="form.enabled" /></el-form-item>
        <el-form-item label="最低转账金额"><el-input-number v-model="form.min_amount" :min="0" :precision="2" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" :rows="4" maxlength="500" show-word-limit type="textarea" /></el-form-item>
      </el-form>
      <div class="mt-4 flex justify-center gap-3"><el-button @click="load">重置</el-button><el-button v-if="canManage" :loading="saving" type="primary" @click="save">保存</el-button></div>
    </el-card>
  </Page>
</template>
