<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createUserFeedbackCategory,
  deleteUserFeedbackCategory,
  fetchUserFeedbackCategories,
  setUserFeedbackCategoryStatus,
  updateUserFeedbackCategory,
  type UserFeedbackCategory,
} from '#/api/core/ecrm';

const loading = ref(false);
const rows = ref<UserFeedbackCategory[]>([]);
const form = reactive({ id: 0, name: '', sort: 0, status: 1 as 0 | 1 });
const key = (action: string, id = 0) => `feedback-category-${action}-${id}-${crypto.randomUUID()}`;

async function load() {
  loading.value = true;
  try {
    rows.value = (await fetchUserFeedbackCategories()).list || [];
  } finally {
    loading.value = false;
  }
}
function reset() {
  Object.assign(form, { id: 0, name: '', sort: 0, status: 1 });
}
function edit(row: UserFeedbackCategory) {
  Object.assign(form, { id: row.id, name: row.name, sort: row.sort, status: row.status });
}
async function save() {
  const name = form.name.trim();
  if (!name || [...name].length > 32) {
    ElMessage.warning('请填写不超过 32 字的分类名称');
    return;
  }
  const payload = { name, sort: Math.max(0, Math.min(9999, Number(form.sort) || 0)), status: form.status, idempotency_key: key(form.id ? 'update' : 'create', form.id) };
  if (form.id) await updateUserFeedbackCategory(form.id, payload);
  else await createUserFeedbackCategory(payload);
  ElMessage.success('已保存');
  reset();
  await load();
}
async function switchStatus(row: UserFeedbackCategory) {
  await setUserFeedbackCategoryStatus(row.id, { status: row.status === 1 ? 0 : 1, idempotency_key: key('status', row.id) });
  ElMessage.success('状态已更新');
  await load();
}
async function remove(row: UserFeedbackCategory) {
  try {
    await ElMessageBox.confirm(`删除“${row.name}”仅从可选分类中移除，既有反馈仍保留原分类文本。`, '删除反馈分类', { type: 'warning' });
    await deleteUserFeedbackCategory(row.id, { idempotency_key: key('delete', row.id) });
    ElMessage.success('已删除');
    if (form.id === row.id) reset();
    await load();
  } catch {}
}
onMounted(() => void load());
</script>

<template>
  <Page title="反馈分类" description="分类维护通过业务服务命令处理；停用或删除不会改写既有反馈的分类快照。">
    <el-card shadow="never">
      <el-form inline @submit.prevent="save">
        <el-form-item label="分类名称"><el-input v-model="form.name" maxlength="32" show-word-limit /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" :max="9999" /></el-form-item>
        <el-form-item label="状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="停用" /></el-form-item>
        <el-button type="primary" @click="save">{{ form.id ? '保存修改' : '新增分类' }}</el-button>
        <el-button v-if="form.id" @click="reset">取消编辑</el-button>
      </el-form>
    </el-card>
    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="id" label="ID" width="90" />
        <el-table-column prop="name" label="分类名称" min-width="180" />
        <el-table-column prop="sort" label="排序" width="100" />
        <el-table-column label="状态" width="110"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column label="操作" width="220"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link @click="switchStatus(row)">{{ row.status ? '停用' : '启用' }}</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column>
      </el-table>
    </el-card>
  </Page>
</template>
