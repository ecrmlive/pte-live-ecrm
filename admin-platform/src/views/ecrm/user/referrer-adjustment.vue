<script setup lang="ts">
import { reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import { changePlatformUserReferrer } from '#/api/core/ecrm';

const submitting = ref(false);
const form = reactive({ user_id: undefined as number | undefined, parent_user_id: undefined as number | undefined, reason: '' });

function reset() { form.user_id = undefined; form.parent_user_id = undefined; form.reason = ''; }

async function submit() {
  const reason = form.reason.trim();
  if (!form.user_id || reason.length < 2 || reason.length > 500) {
    ElMessage.warning('请填写用户 ID 与 2 至 500 字的调整原因');
    return;
  }
  if (form.parent_user_id === form.user_id) { ElMessage.warning('上级用户不能是本人'); return; }
  submitting.value = true;
  try {
    await changePlatformUserReferrer(form.user_id, { parent_user_id: form.parent_user_id || 0, reason, idempotency_key: `referrer-${form.user_id}-${crypto.randomUUID()}` });
    ElMessage.success(form.parent_user_id ? '推荐关系已调整并写入审计记录' : '上级关系已清除并写入审计记录');
    reset();
  } finally { submitting.value = false; }
}
</script>

<template>
  <Page title="用户推荐关系调整" description="仅平台角色可执行。系统拒绝自荐、无效上级和循环推荐关系；清除上级仅置空当前绑定，不删除佣金、订单或历史审计事实。">
    <el-alert class="mb-4" type="warning" :closable="false" title="请根据已核验的中文工单操作。该操作不会重算已产生的佣金，也不会改变既有订单归属。" />
    <el-card shadow="never" class="max-w-3xl"><el-form label-width="118px" @submit.prevent="submit"><el-form-item label="用户 ID" required><el-input-number v-model="form.user_id" :min="1" controls-position="right" /></el-form-item><el-form-item label="上级用户 ID"><el-input-number v-model="form.parent_user_id" :min="1" controls-position="right" /><span class="ml-3 text-sm text-gray-500">留空即清除上级绑定</span></el-form-item><el-form-item label="调整原因" required><el-input v-model="form.reason" type="textarea" :rows="4" maxlength="500" show-word-limit placeholder="例如：虚构中文工单演示，修正推荐关系。" /></el-form-item><el-form-item><el-button type="primary" :loading="submitting" @click="submit">确认调整</el-button><el-button :disabled="submitting" @click="reset">重置</el-button></el-form-item></el-form></el-card>
  </Page>
</template>
