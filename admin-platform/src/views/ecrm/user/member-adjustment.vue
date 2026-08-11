<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import { adjustPlatformUserMemberLevel, fetchPlatformMemberLevels, type PlatformMemberLevel } from '#/api/core/ecrm';
import UserRelationSelect from '#/components/ecrm/UserRelationSelect.vue';

const levels = ref<PlatformMemberLevel[]>([]);
const loading = ref(false);
const submitting = ref(false);
const form = reactive({ user_id: undefined as number | undefined, level_id: 0, reason: '' });

async function loadLevels() {
  loading.value = true;
  try {
    levels.value = (await fetchPlatformMemberLevels()).list || [];
  } finally { loading.value = false; }
}

function reset() { form.user_id = undefined; form.level_id = 0; form.reason = ''; }

async function submit() {
  const reason = form.reason.trim();
  if (!form.user_id || reason.length < 2 || reason.length > 500) {
    ElMessage.warning('请填写用户 ID 与 2 至 500 字的调整原因');
    return;
  }
  submitting.value = true;
  try {
    await adjustPlatformUserMemberLevel(form.user_id, { level_id: form.level_id, reason, idempotency_key: `member-level-${form.user_id}-${crypto.randomUUID()}` });
    ElMessage.success('会员等级已调整，并已写入不可变变更日志');
    reset();
  } finally { submitting.value = false; }
}

onMounted(() => void loadLevels());
</script>

<template>
  <Page title="用户会员等级调整" description="仅平台角色可执行。调整只更新用户当前会员等级投影和不可变日志；不会直接赠送 SVIP、优惠券、余额或改变既有订单价格。">
    <el-alert class="mb-4" type="warning" :closable="false" title="请依据已核验的中文工单填写原因。选择“普通会员”会清除当前会员等级，不删除历史等级变更记录。" />
    <el-card v-loading="loading" shadow="never" class="max-w-3xl">
      <el-form label-width="118px" @submit.prevent="submit">
        <el-form-item label="关联用户" required><UserRelationSelect v-model="form.user_id" /></el-form-item>
        <el-form-item label="目标会员等级" required><el-select v-model="form.level_id" class="w-80"><el-option :value="0" label="普通会员（清除等级）" /><el-option v-for="level in levels" :key="level.id" :value="level.id" :label="`${level.name}（等级 ${level.rank}）`" /></el-select></el-form-item>
        <el-form-item label="调整原因" required><el-input v-model="form.reason" type="textarea" :rows="4" maxlength="500" show-word-limit placeholder="例如：虚构中文工单演示，核验后修正会员等级。" /></el-form-item>
        <el-form-item><el-button type="primary" :loading="submitting" @click="submit">确认调整</el-button><el-button :disabled="submitting" @click="reset">重置</el-button></el-form-item>
      </el-form>
    </el-card>
  </Page>
</template>
