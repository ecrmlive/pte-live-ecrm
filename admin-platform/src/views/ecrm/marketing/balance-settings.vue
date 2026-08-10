<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformBalanceConfigApi,
  savePlatformBalanceConfigApi,
  type PlatformBalanceConfig,
} from '#/api/core/platform-mall-setting';
import SettingsTabLayout from '#/components/settings/SettingsTabLayout.vue';

const loading = ref(false);
const saving = ref(false);
const canRead = ref(false);
const canManage = ref(false);

const form = reactive<PlatformBalanceConfig>({
  balance_func_status: 1,
  recharge_switch: 1,
  store_user_min_recharge: 1,
  recharge_attention:
    '1、账户充值仅限用于购买商城内商品，不可提现\n2、账户充值成功后，一般1～5分钟到账\n3、如有疑问，请联系客服',
});

function applyConfig(data: PlatformBalanceConfig) {
  form.balance_func_status = data.balance_func_status === 1 ? 1 : 0;
  form.recharge_switch = data.recharge_switch === 1 ? 1 : 0;
  form.store_user_min_recharge = Number(data.store_user_min_recharge ?? 1);
  form.recharge_attention = data.recharge_attention || '';
}

function validate(): string | null {
  if (form.balance_func_status !== 0 && form.balance_func_status !== 1) {
    return '请设置余额功能开关';
  }
  if (form.recharge_switch !== 0 && form.recharge_switch !== 1) {
    return '请设置余额充值开关';
  }
  if (
    !Number.isFinite(form.store_user_min_recharge) ||
    form.store_user_min_recharge < 0
  ) {
    return '用户最低充值金额请填写不小于 0 的数字';
  }
  return null;
}

async function load() {
  if (!canRead.value) {
    ElMessage.warning('无权查看余额设置');
    return;
  }
  loading.value = true;
  try {
    const data = await getPlatformBalanceConfigApi();
    applyConfig(data.config);
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!canManage.value) {
    ElMessage.warning('无权保存余额设置');
    return;
  }
  const err = validate();
  if (err) {
    ElMessage.warning(err);
    return;
  }
  saving.value = true;
  try {
    const saved = await savePlatformBalanceConfigApi({
      balance_func_status: form.balance_func_status === 1 ? 1 : 0,
      recharge_switch: form.recharge_switch === 1 ? 1 : 0,
      store_user_min_recharge: Number(form.store_user_min_recharge),
      recharge_attention: form.recharge_attention || '',
    });
    applyConfig(saved);
    ElMessage.success('保存成功');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  canRead.value =
    codes.includes('marketing.balance.settings.read') ||
    codes.includes('marketing.balance.settings.manage');
  canManage.value = codes.includes('marketing.balance.settings.manage');
  await load();
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <SettingsTabLayout v-loading="loading">
      <ElForm
        :disabled="!canManage"
        class="balance-settings-form"
        label-width="180px"
      >
        <ElFormItem label="余额功能：">
          <div>
            <ElSwitch
              v-model="form.balance_func_status"
              :active-value="1"
              :inactive-value="0"
            />
            <div class="form-tip">商城余额功能启用或关闭</div>
          </div>
        </ElFormItem>

        <ElFormItem label="余额充值开关：">
          <div>
            <ElSwitch
              v-model="form.recharge_switch"
              :active-value="1"
              :inactive-value="0"
            />
            <div class="form-tip">商城余额充值启用或关闭</div>
          </div>
        </ElFormItem>

        <ElFormItem label="用户最低充值金额：">
          <div>
            <ElInputNumber
              v-model="form.store_user_min_recharge"
              :min="0"
              :precision="2"
              :step="1"
              controls-position="right"
            />
            <div class="form-tip">单次充值最低金额</div>
          </div>
        </ElFormItem>

        <ElFormItem label="充值注意事项：">
          <div class="balance-settings-form__attention">
            <ElInput
              v-model="form.recharge_attention"
              :autosize="{ minRows: 5, maxRows: 12 }"
              type="textarea"
              placeholder="请输入充值注意事项"
            />
            <div class="form-tip">将在充值页展示</div>
          </div>
        </ElFormItem>
      </ElForm>

      <template #actions>
        <ElButton
          type="primary"
          :disabled="!canManage"
          :loading="saving"
          @click="save"
        >
          提交
        </ElButton>
      </template>
    </SettingsTabLayout>
  </Page>
</template>

<style scoped>
.balance-settings-form {
  width: 100%;
  max-width: none;
  padding: 8px 0 0;
}

.balance-settings-form :deep(.el-form-item) {
  margin-bottom: 22px;
}

.balance-settings-form :deep(.el-form-item__label) {
  align-items: flex-start;
  padding-top: 6px;
  color: hsl(var(--foreground));
  font-weight: 400;
}

.balance-settings-form__attention {
  width: 100%;
  max-width: 560px;
}

.form-tip {
  margin-top: 8px;
  font-size: 12px;
  line-height: 1.5;
  color: hsl(var(--muted-foreground));
}
</style>
