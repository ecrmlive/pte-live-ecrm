<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import { assignPlatformUserLabels, fetchPlatformUserLabelOptions, type PlatformUserLabelOption } from '#/api/core/ecrm';
import UserRelationSelect from '#/components/ecrm/UserRelationSelect.vue';

const labels = ref<PlatformUserLabelOption[]>([]);
const loading = ref(false);
const submitting = ref(false);
const form = reactive({ user_ids: [] as number[], label_ids: [] as number[], reason: '' });
async function loadLabels() { loading.value = true; try { labels.value = (await fetchPlatformUserLabelOptions()).list || []; } finally { loading.value = false; } }
function reset() { Object.assign(form, { user_ids: [], label_ids: [], reason: '' }); }
async function submit() {
  const userIDs = [...new Set(form.user_ids)];
  const reason = form.reason.trim();
  if (!userIDs.length || userIDs.length > 100 || form.label_ids.length > 50 || reason.length < 2 || reason.length > 500) { ElMessage.warning('请选择 1 至 100 个用户、至多 50 个标签并填写调整原因'); return; }
  submitting.value = true;
  try {
    await assignPlatformUserLabels({ user_ids: userIDs, label_ids: form.label_ids, reason, idempotency_key: `user-label-${crypto.randomUUID()}` });
    ElMessage.success(form.label_ids.length ? '用户标签已替换并写入审计记录' : '用户运营标签已清除并写入审计记录'); reset();
  } finally { submitting.value = false; }
}
onMounted(() => void loadLabels());
</script>

<template>
  <Page title="用户标签归属" description="仅平台角色可执行。一次替换 1 至 100 位用户的运营标签；不改变登录、会员、订单、资金或佣金事实。">
    <el-alert class="mb-4" type="warning" :closable="false" title="标签集合按整体替换。清空选择后提交即清除运营标签；已删除标签不能再分配。" />
    <el-card v-loading="loading" shadow="never" class="max-w-3xl"><el-form label-width="116px" @submit.prevent="submit"><el-form-item label="关联用户" required><UserRelationSelect v-model="form.user_ids" multiple placeholder="请选择用户，最多 100 位" /></el-form-item><el-form-item label="目标标签"><el-select v-model="form.label_ids" class="w-96" multiple clearable filterable placeholder="留空代表清除运营标签"><el-option v-for="label in labels" :key="label.label_id" :value="label.label_id" :label="`${label.label_name}（#${label.label_id}）`" /></el-select></el-form-item><el-form-item label="调整原因" required><el-input v-model="form.reason" type="textarea" :rows="4" maxlength="500" show-word-limit placeholder="例如：虚构中文工单演示，标记为新品体验与高频复购用户。" /></el-form-item><el-form-item><el-button type="primary" :loading="submitting" @click="submit">确认替换标签</el-button><el-button :disabled="submitting" @click="reset">重置</el-button></el-form-item></el-form></el-card>
  </Page>
</template>
