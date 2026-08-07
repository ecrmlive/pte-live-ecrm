<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformMarginConfigApi,
  savePlatformMarginConfigApi,
  type PlatformMarginConfig,
} from '#/api/core/platform-mall-setting';

const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);
const form = ref<PlatformMarginConfig>({
  margin_remind_switch: false,
  margin_remind_day: 30,
});

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformMarginConfigApi();
    form.value = data.config;
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (form.value.margin_remind_day < 0 || form.value.margin_remind_day > 3650) {
    ElMessage.warning('补缴提醒时间需在 0～3650 天内');
    return;
  }
  saving.value = true;
  try {
    form.value = await savePlatformMarginConfigApi({
      margin_remind_switch: !!form.value.margin_remind_switch,
      margin_remind_day: Number(form.value.margin_remind_day) || 0,
    });
    ElMessage.success('保证金配置已保存');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = permissions.includes('store.margin_config.manage');
});
</script>

<template>
  <!-- 对齐 CRMEB systemForm/Basics/margin：无页内标题，整块白底表单，非灰底小卡片 -->
  <Page auto-content-height content-class="!bg-card !p-0">
    <div v-loading="loading" class="margin-config-page">
      <ElForm
        :disabled="!canManage"
        label-width="160px"
        class="margin-config-form"
      >
        <ElFormItem label="保证金补缴提醒">
          <div>
            <ElSwitch
              v-model="form.margin_remind_switch"
              inline-prompt
              active-text="开启"
              inactive-text="关闭"
              :width="56"
            />
            <div class="tips">
              开启后，保证金不足时自动向商户发送补缴提醒
            </div>
          </div>
        </ElFormItem>
        <ElFormItem label="补缴提醒时间（天）">
          <div>
            <ElInputNumber
              v-model="form.margin_remind_day"
              class="margin-day-input"
              :min="0"
              :max="3650"
              :step="1"
              :precision="0"
            />
            <div class="tips">
              店铺保证金补缴提醒开启，并设置天数，比如:填写30天，即自保证金不足日开始计算连续提醒商户补缴保证金30天，如果期间商户还未补足保证金，30天满后自动关闭店铺
            </div>
          </div>
        </ElFormItem>
        <ElFormItem v-if="canManage">
          <ElButton :loading="saving" type="primary" @click="save">
            提交
          </ElButton>
        </ElFormItem>
      </ElForm>
    </div>
  </Page>
</template>

<style scoped>
.margin-config-page {
  box-sizing: border-box;
  min-height: 100%;
  padding: 24px 32px 40px;
  background: hsl(var(--card));
}

.margin-config-form {
  width: 100%;
  max-width: none;
}

.margin-config-form :deep(.el-form-item) {
  margin-bottom: 18px;
}

.margin-config-form :deep(.el-form-item__label) {
  align-items: center;
  color: hsl(var(--foreground));
  font-weight: 400;
}

.margin-config-form :deep(.el-switch) {
  --el-switch-on-color: hsl(var(--primary));
}

.margin-day-input {
  width: 160px;
}

.tips {
  margin-top: 6px;
  max-width: 720px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  line-height: 1.5;
}
</style>
