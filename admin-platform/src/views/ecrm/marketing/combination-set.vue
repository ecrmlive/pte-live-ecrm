<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformGroupBuyingConfigApi,
  savePlatformGroupBuyingConfigApi,
  type PlatformGroupBuyingConfig,
} from '#/api/core/platform-mall-setting';
import SettingsTabLayout from '#/components/settings/SettingsTabLayout.vue';

const loading = ref(false);
const saving = ref(false);
const canManage = ref(false);

const form = reactive<PlatformGroupBuyingConfig>({
  ficti_status: 1,
  group_buying_rate: 30,
});

const fictiFillMaxRate = computed(() => {
  const rate = Number(form.group_buying_rate || 0);
  if (!Number.isFinite(rate)) return 0;
  return Math.max(0, Math.min(100, 100 - Math.trunc(rate)));
});

function validate(): string | null {
  if (form.ficti_status !== 0 && form.ficti_status !== 1) {
    return '请选择虚拟成团启用状态';
  }
  const rate = form.group_buying_rate;
  if (!Number.isInteger(rate) || rate < 0 || rate > 100) {
    return '真实成团最小比例请填写 0～100 之间的整数';
  }
  return null;
}

async function load() {
  loading.value = true;
  try {
    const data = await getPlatformGroupBuyingConfigApi();
    form.ficti_status = data.config.ficti_status === 1 ? 1 : 0;
    form.group_buying_rate = Number(data.config.group_buying_rate ?? 30);
  } finally {
    loading.value = false;
  }
}

async function save() {
  const err = validate();
  if (err) {
    ElMessage.warning(err);
    return;
  }
  saving.value = true;
  try {
    const saved = await savePlatformGroupBuyingConfigApi({
      ficti_status: form.ficti_status === 1 ? 1 : 0,
      group_buying_rate: Math.trunc(Number(form.group_buying_rate)),
    });
    form.ficti_status = saved.ficti_status === 1 ? 1 : 0;
    form.group_buying_rate = saved.group_buying_rate;
    ElMessage.success('保存成功');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  canManage.value =
    codes.includes('marketing.combination.manage') ||
    codes.includes('marketing.combination.set');
  await load();
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <SettingsTabLayout v-loading="loading">
      <ElForm
        :disabled="!canManage"
        class="combination-set-form"
        label-width="180px"
      >
        <ElFormItem label="虚拟成团启用：" required>
          <ElRadioGroup v-model="form.ficti_status">
            <ElRadio :label="1">启用</ElRadio>
            <ElRadio :label="0">关闭</ElRadio>
          </ElRadioGroup>
        </ElFormItem>

        <ElFormItem label="真实成团最小比例：" required>
          <div class="combination-set-form__rate">
            <ElInputNumber
              v-model="form.group_buying_rate"
              :max="100"
              :min="0"
              :precision="0"
              :step="1"
              controls-position="right"
            />
            <span class="combination-set-form__unit">%</span>
          </div>
        </ElFormItem>

        <ElFormItem label="虚拟成团补齐最大比例：">
          <span class="combination-set-form__readonly">
            {{ fictiFillMaxRate }}%
          </span>
        </ElFormItem>
      </ElForm>

      <template #actions>
        <ElButton
          type="primary"
          :disabled="!canManage"
          :loading="saving"
          @click="save"
        >
          确认
        </ElButton>
      </template>
    </SettingsTabLayout>
  </Page>
</template>

<style scoped>
.combination-set-form {
  width: 100%;
  max-width: none;
  padding: 8px 0 0;
}

.combination-set-form :deep(.el-form-item) {
  margin-bottom: 22px;
}

.combination-set-form :deep(.el-form-item__label) {
  align-items: center;
  color: hsl(var(--foreground));
  font-weight: 400;
}

.combination-set-form__rate {
  display: inline-flex;
  gap: 8px;
  align-items: center;
}

.combination-set-form__unit,
.combination-set-form__readonly {
  color: hsl(var(--foreground));
  font-size: 14px;
  line-height: 32px;
}
</style>
