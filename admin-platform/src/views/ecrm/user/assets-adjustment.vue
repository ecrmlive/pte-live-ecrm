<script setup lang="ts">
import { reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import { adjustPlatformUserAsset } from '#/api/core/ecrm';

const submitting = ref(false);
const form = reactive({
  user_id: undefined as number | undefined,
  asset_type: 'balance' as 'balance' | 'points',
  amount: undefined as number | undefined,
  reason: '',
});

function reset() {
  form.user_id = undefined;
  form.asset_type = 'balance';
  form.amount = undefined;
  form.reason = '';
}

async function submit() {
  const reason = form.reason.trim();
  if (!form.user_id || !form.amount || !reason || reason.length < 2 || reason.length > 500) {
    ElMessage.warning('请填写用户 ID、非零调整数和 2 至 500 字的调账原因');
    return;
  }
  if (form.asset_type === 'points' && !Number.isInteger(form.amount)) {
    ElMessage.warning('积分调整必须为整数');
    return;
  }
  submitting.value = true;
  try {
    await adjustPlatformUserAsset(form.user_id, {
      asset_type: form.asset_type,
      amount: form.amount,
      reason,
      idempotency_key: `user-${form.asset_type}-${form.user_id}-${crypto.randomUUID()}`,
    });
    ElMessage.success('调账已登记到不可变流水与审计记录');
    reset();
  } finally {
    submitting.value = false;
  }
}
</script>

<template>
  <Page title="用户资产调整" description="仅平台角色可执行。正数增加、负数扣减；系统在同一事务锁定账户、写资产流水和不可变人工审计，不接入外部支付或真实账户信息。">
    <el-alert class="mb-4" type="warning" :closable="false" title="这是高风险财务操作：请基于已核验的工单执行，并填写可审计原因。余额或积分不足时系统拒绝扣减。" />
    <el-card shadow="never" class="max-w-3xl">
      <el-form label-width="108px" @submit.prevent="submit">
        <el-form-item label="用户 ID" required><el-input-number v-model="form.user_id" :min="1" controls-position="right" /></el-form-item>
        <el-form-item label="资产类型" required><el-radio-group v-model="form.asset_type"><el-radio value="balance">余额</el-radio><el-radio value="points">积分</el-radio></el-radio-group></el-form-item>
        <el-form-item label="调整数" required><el-input-number v-model="form.amount" :precision="form.asset_type === 'points' ? 0 : 2" :step="form.asset_type === 'points' ? 1 : 0.01" controls-position="right" /><span class="ml-3 text-sm text-gray-500">正数增加，负数扣减</span></el-form-item>
        <el-form-item label="调账原因" required><el-input v-model="form.reason" type="textarea" :rows="4" maxlength="500" show-word-limit placeholder="例如：虚构中文工单演示，修正重复扣减。" /></el-form-item>
        <el-form-item><el-button type="primary" :loading="submitting" @click="submit">确认登记</el-button><el-button :disabled="submitting" @click="reset">重置</el-button></el-form-item>
      </el-form>
    </el-card>
  </Page>
</template>
