<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformShopConfigApi,
  savePlatformShopConfigApi,
  type PlatformShopConfig,
} from '#/api/core/platform-mall-setting';

const note = ref('');
const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<PlatformShopConfig>({
  site_name: '',
  site_url: '',
  order_auto_cancel_minutes: 30,
  order_auto_receive_days: 7,
  enabled: true,
  remark: '',
});

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformShopConfigApi();
    note.value = data.note;
    form.value = data.config;
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!form.value.site_name.trim()) {
    ElMessage.warning('站点名称不能为空');
    return;
  }
  saving.value = true;
  try {
    form.value = await savePlatformShopConfigApi({
      ...form.value,
      site_name: form.value.site_name.trim(),
      site_url: form.value.site_url.trim(),
      remark: form.value.remark.trim(),
    });
    ElMessage.success('商城设置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes('setting.shop.manage');
});
</script>

<template>
  <Page title="商城设置" description="维护站点名称、访问地址与订单超时规则；不含密钥或支付凭据。">
    <el-card v-loading="loading" shadow="never">
      <el-alert :title="note" type="info" :closable="false" class="mb-4" />
      <el-form :disabled="!canManage" label-width="160px" class="max-w-3xl">
        <el-form-item label="站点名称" required>
          <el-input v-model="form.site_name" maxlength="128" placeholder="例如 七禧商城" />
        </el-form-item>
        <el-form-item label="站点地址">
          <el-input v-model="form.site_url" maxlength="256" placeholder="https://example.test" />
        </el-form-item>
        <el-form-item label="未支付自动取消（分钟）">
          <el-input-number v-model="form.order_auto_cancel_minutes" :min="0" :max="10080" />
        </el-form-item>
        <el-form-item label="自动确认收货（天）">
          <el-input-number v-model="form.order_auto_receive_days" :min="0" :max="365" />
        </el-form-item>
        <el-form-item label="启用商城">
          <el-switch v-model="form.enabled" />
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
