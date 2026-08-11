<script setup lang="ts">
import type { CustomerServiceSettings } from '#/api/core/customer-service';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';

import { getUserInfoApi } from '#/api/core/auth';
import {
  fetchCustomerServiceSettings,
  updateCustomerServiceSettings,
} from '#/api/core/customer-service';

const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const settings = reactive<CustomerServiceSettings>({
  auto_reply_enabled: false,
  auto_reply_text: '',
  enterprise_wechat_corp_id: '',
  enterprise_wechat_url: '',
  max_sessions_per_agent: 20,
  queue_mode: 'manual',
  redirect_url: '',
  service_phone: '',
  service_type: 'system',
});

const description = computed(() => {
  switch (settings.service_type) {
    case 'disabled':
      return '关闭后，商城端不展示平台客服入口。';
    case 'system':
      return '由平台客服在线接待用户咨询。';
    case 'phone':
      return '商城端点击客服入口后发起电话呼叫。';
    case 'enterprise_wechat':
      return '用户将跳转至企业微信客服。';
    case 'link':
      return '商城端跳转至已配置的第三方客服页面。';
    case 'mini_program':
      return '用户可在小程序内发起客服咨询。';
  }
});

function validBeforeSave() {
  if (settings.service_type === 'phone' && !settings.service_phone.trim()) {
    ElMessage.warning('请填写客服热线');
    return false;
  }
  if (settings.service_type === 'enterprise_wechat' && (!settings.enterprise_wechat_url.trim() || !settings.enterprise_wechat_corp_id.trim())) {
    ElMessage.warning('请填写企业微信客服链接与企业 ID');
    return false;
  }
  if (settings.service_type === 'link' && !settings.redirect_url.trim()) {
    ElMessage.warning('请填写跳转链接');
    return false;
  }
  return true;
}

async function load() {
  loading.value = true;
  try {
    const result = await fetchCustomerServiceSettings();
    Object.assign(settings, result.settings);
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!validBeforeSave()) return;
  saving.value = true;
  try {
    const result = await updateCustomerServiceSettings({
      ...settings,
      auto_reply_text: settings.auto_reply_text.trim(),
      enterprise_wechat_corp_id: settings.enterprise_wechat_corp_id.trim(),
      enterprise_wechat_url: settings.enterprise_wechat_url.trim(),
      redirect_url: settings.redirect_url.trim(),
      service_phone: settings.service_phone.trim(),
    });
    Object.assign(settings, result.settings);
    ElMessage.success('客服设置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const profile = await getUserInfoApi();
  canManage.value = profile.roles.includes('platform');
  await load();
});
</script>

<template>
  <Page auto-content-height>
    <div v-loading="loading" class="customer-service-settings">
      <el-form label-width="132px" class="customer-service-settings__form" :disabled="!canManage">
        <el-form-item label="平台客服类型">
          <el-radio-group v-model="settings.service_type">
            <el-radio value="disabled">关闭</el-radio>
            <el-radio value="system">系统客服</el-radio>
            <el-radio value="phone">拨打电话</el-radio>
            <el-radio value="enterprise_wechat">企业微信</el-radio>
            <el-radio value="link">跳转链接</el-radio>
            <el-radio value="mini_program">小程序客服</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item class="customer-service-settings__description" label="">
          {{ description }}
        </el-form-item>

        <el-form-item v-if="settings.service_type === 'phone'" label="客服热线" required>
          <el-input v-model="settings.service_phone" maxlength="32" placeholder="例如：400-888-8888" />
        </el-form-item>
        <template v-else-if="settings.service_type === 'enterprise_wechat'">
          <el-form-item label="客服链接" required>
            <el-input v-model="settings.enterprise_wechat_url" maxlength="2048" placeholder="请输入企业微信客服链接" />
          </el-form-item>
          <el-form-item label="企业 ID" required>
            <el-input v-model="settings.enterprise_wechat_corp_id" maxlength="128" placeholder="请输入企业 ID" />
          </el-form-item>
        </template>
        <el-form-item v-else-if="settings.service_type === 'link'" label="跳转链接" required>
          <el-input v-model="settings.redirect_url" maxlength="2048" placeholder="请输入客服跳转链接" />
        </el-form-item>
      </el-form>
      <div class="customer-service-settings__footer">
        <ElButton v-if="canManage" type="primary" :loading="saving" @click="save">保存</ElButton>
      </div>
    </div>
  </Page>
</template>

<style scoped>
.customer-service-settings {
  padding: 40px 48px 24px;
  background: hsl(var(--card));
  border-radius: 8px;
}

.customer-service-settings__form {
  width: min(100%, 1800px);
}

.customer-service-settings__description :deep(.el-form-item__content) {
  color: var(--el-text-color-secondary);
  line-height: 1.75;
}

.customer-service-settings__footer {
  display: flex;
  align-items: center;
  justify-content: center;
  width: min(100%, 1800px);
  margin-top: 30px;
  padding-top: 28px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
  border-top: 1px solid var(--el-border-color-lighter);
}
</style>
