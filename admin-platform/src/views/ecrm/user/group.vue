<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createPlatformUserGroupApi,
  deletePlatformUserGroupApi,
  listPlatformUserGroupsApi,
  updatePlatformUserGroupApi,
  type PlatformUserGroup,
} from '#/api/core/platform-user-group';

const loading = ref(false);
const saving = ref(false);
const rows = ref<PlatformUserGroup[]>([]);
const total = ref(0);
const dialogOpen = ref(false);
const editing = ref<PlatformUserGroup>();
const query = reactive({ limit: 20, page: 1 });
const form = reactive({ group_name: '', sort: 0 });

async function load() {
  loading.value = true;
  try {
    const result = await listPlatformUserGroupsApi(query);
    rows.value = result.list;
    total.value = result.total;
  } finally { loading.value = false; }
}

function openCreate() {
  editing.value = undefined;
  Object.assign(form, { group_name: '', sort: 0 });
  dialogOpen.value = true;
}

function openEdit(row: PlatformUserGroup) {
  editing.value = row;
  Object.assign(form, { group_name: row.group_name, sort: row.sort });
  dialogOpen.value = true;
}

async function save() {
  if (!form.group_name.trim()) {
    ElMessage.warning('请填写分组名称');
    return;
  }
  saving.value = true;
  try {
    const body = { group_name: form.group_name.trim(), sort: form.sort };
    if (editing.value) await updatePlatformUserGroupApi(editing.value.group_id, body);
    else await createPlatformUserGroupApi(body);
    dialogOpen.value = false;
    ElMessage.success(editing.value ? '用户分组已更新' : '用户分组已创建');
    await load();
  } finally { saving.value = false; }
}

async function remove(row: PlatformUserGroup) {
  try {
    await ElMessageBox.confirm(`删除用户分组“${row.group_name}”后不可恢复，是否继续？`, '删除确认', { type: 'warning' });
    await deletePlatformUserGroupApi(row.group_id);
    ElMessage.success('用户分组已删除');
    await load();
  } catch {
    // 用户取消或接口错误时由统一请求层处理。
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="用户分组" description="维护平台用户分组；分组仅用于用户运营管理，不改变用户登录或订单归属。">
    <template #extra><el-button type="primary" @click="openCreate">新增分组</el-button></template>
    <el-card shadow="never"><el-table v-loading="loading" :data="rows" row-key="group_id"><el-table-column label="ID" prop="group_id" width="88" /><el-table-column label="分组名称" min-width="220" prop="group_name" /><el-table-column label="排序" prop="sort" width="96" /><el-table-column label="创建时间" min-width="180" prop="create_time" /><el-table-column fixed="right" label="操作" width="130"><template #default="{ row }"><el-button link type="primary" @click="openEdit(row)">编辑</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
    <el-dialog v-model="dialogOpen" :title="editing ? '编辑用户分组' : '新增用户分组'" width="460px" destroy-on-close><el-form label-width="80px"><el-form-item label="分组名称" required><el-input v-model="form.group_name" maxlength="32" show-word-limit /></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" class="w-full" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
