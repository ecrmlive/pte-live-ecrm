<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformPayConfigApi,
  savePlatformPayConfigApi,
  type PlatformPayConfig,
} from '#/api/core/platform-mall-setting';

const note = ref('');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<PlatformPayConfig>({
  wechat_enabled: false,
  alipay_enabled: false,
  balance_enabled: true,
  remark: '',
});

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformPayConfigApi();
    note.value = data.note;
    form.value = data.config;
  } finally {
    loading.value = false;
  }
}

async function save() {
  saving.value = true;
  try {
    form.value = await savePlatformPayConfigApi({
      ...form.value,
      remark: form.value.remark.trim(),
    });
    ElMessage.success('支付设置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes('setting.pay.manage');
});
</script>

<template>
  <Page title="支付设置" description="仅维护支付方式开关；微信/支付宝密钥与证书请通过云服务配置维护。">
    <el-card v-loading="loading" shadow="never">
      <el-alert :title="note" type="warning" :closable="false" class="mb-4" />
      <el-form :disabled="!canManage" label-width="160px" class="max-w-3xl">
        <el-form-item label="启用微信支付">
          <el-switch v-model="form.wechat_enabled" />
        </el-form-item>
        <el-form-item label="启用支付宝">
          <el-switch v-model="form.alipay_enabled" />
        </el-form-item>
        <el-form-item label="启用余额支付">
          <el-switch v-model="form.balance_enabled" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" :rows="4" maxlength="500" show-word-limit type="textarea" />
        </el-form-item>
      </el-form>
      <div class="mt-4 flex justify-end">
        <el-button @click="load">重置</el-button>
        <el-button v-if="canManage" :loading="saving" type="primary" @click="save">保存</el-button>
      </div>
    </el-card>
  </Page>
</template>
