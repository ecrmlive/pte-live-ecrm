<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import {
  assignPlatformUserGroups,
  fetchPlatformUserGroupOptions,
  type PlatformUserGroupOption,
} from '#/api/core/ecrm';
import UserRelationSelect from '#/components/ecrm/UserRelationSelect.vue';

const groups = ref<PlatformUserGroupOption[]>([]);
const loading = ref(false);
const submitting = ref(false);
const form = reactive({ user_ids: [] as number[], group_id: 0, reason: '' });

async function loadGroups() {
  loading.value = true;
  try { groups.value = (await fetchPlatformUserGroupOptions()).list || []; } finally { loading.value = false; }
}

function reset() { Object.assign(form, { user_ids: [], group_id: 0, reason: '' }); }

async function submit() {
  const userIDs = [...new Set(form.user_ids)];
  const reason = form.reason.trim();
  if (!userIDs || reason.length < 2 || reason.length > 500) {
    ElMessage.warning('请选择 1 至 100 个用户，并填写 2 至 500 字的调整原因');
    return;
  }
  submitting.value = true;
  try {
    await assignPlatformUserGroups({ user_ids: userIDs, group_id: form.group_id, reason, idempotency_key: `user-group-${crypto.randomUUID()}` });
    ElMessage.success(form.group_id ? '用户分组已调整并写入审计记录' : '用户已移出运营分组并写入审计记录');
    reset();
  } finally { submitting.value = false; }
}

onMounted(() => void loadGroups());
</script>

<template>
  <Page title="用户分组归属" description="仅平台角色可执行。可一次调整 1 至 100 位用户的运营分组；不改变登录、会员、订单、资金或佣金事实。">
    <el-alert class="mb-4" type="warning" :closable="false" title="按虚构中文工单核验后操作。留在“未分组”代表移出运营分组，已删除分组不可再分配。" />
    <el-card v-loading="loading" shadow="never" class="max-w-3xl">
      <el-form label-width="116px" @submit.prevent="submit">
        <el-form-item label="关联用户" required>
          <UserRelationSelect v-model="form.user_ids" multiple placeholder="请选择用户，最多 100 位" />
        </el-form-item>
        <el-form-item label="目标分组" required>
          <el-select v-model="form.group_id" class="w-96">
            <el-option :value="0" label="未分组（移出运营分组）" />
            <el-option v-for="group in groups" :key="group.group_id" :value="group.group_id" :label="`${group.group_name}（#${group.group_id}）`" />
          </el-select>
        </el-form-item>
        <el-form-item label="调整原因" required>
          <el-input v-model="form.reason" type="textarea" :rows="4" maxlength="500" show-word-limit placeholder="例如：虚构中文工单演示，将活动报名用户纳入新品体验分组。" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submit">确认调整</el-button>
          <el-button :disabled="submitting" @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </Page>
</template>
