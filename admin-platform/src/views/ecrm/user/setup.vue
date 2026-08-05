<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformUserSetupConfigApi,
  savePlatformUserSetupConfigApi,
  type PlatformUserSetupConfig,
} from '#/api/core/platform-setting-ext';

const note = ref('');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<PlatformUserSetupConfig>({
  register_enabled: true,
  mobile_required: true,
  invite_required: false,
  remark: '',
});

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformUserSetupConfigApi();
    note.value = data.note;
    form.value = data.config;
  } finally {
    loading.value = false;
  }
}

async function save() {
  saving.value = true;
  try {
    form.value = await savePlatformUserSetupConfigApi({ ...form.value, remark: form.value.remark.trim() });
    ElMessage.success('用户注册设置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes('user.setup.manage');
});
</script>

<template>
  <Page title="用户注册设置" description="维护注册开关与校验规则；不含短信或第三方登录密钥。">
    <el-card v-loading="loading" shadow="never">
      <el-alert :title="note" type="warning" :closable="false" class="mb-4" />
      <el-alert v-if="!canManage" class="mb-4" title="当前账号无写入权限；配置只读展示。" type="info" :closable="false" />
      <el-form :disabled="!canManage" label-width="160px" class="max-w-3xl">
        <el-form-item label="开放注册"><el-switch v-model="form.register_enabled" /></el-form-item>
        <el-form-item label="手机号必填"><el-switch v-model="form.mobile_required" /></el-form-item>
        <el-form-item label="邀请码必填"><el-switch v-model="form.invite_required" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" :rows="4" maxlength="500" show-word-limit type="textarea" /></el-form-item>
      </el-form>
      <div class="mt-4 flex justify-end"><el-button @click="load">重置</el-button><el-button v-if="canManage" :loading="saving" type="primary" @click="save">保存设置</el-button></div>
    </el-card>
  </Page>
</template>
