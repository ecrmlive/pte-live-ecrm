<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import {
  createPlatformAdmin,
  fetchPlatformAdmins,
  updatePlatformAdmin,
  type PlatformAdminRow,
} from '#/api/core/mergers';

const rows = ref<PlatformAdminRow[]>([]);
const total = ref(0);
const loading = ref(false);
const dialogOpen = ref(false);
const editing = ref<PlatformAdminRow>();
const query = reactive({ page: 1, limit: 20 });
const form = reactive({ account: '', password: '', real_name: '', phone: '', roles: '', status: 1, is_agent: 0, region_ids: '', circle_agent_id: 0 });
const title = computed(() => editing.value ? '编辑管理员' : '新增管理员');

async function load() { loading.value = true; try { const result = await fetchPlatformAdmins(query); rows.value = result.list; total.value = result.total; } finally { loading.value = false; } }
function reset() { Object.assign(form, { account: '', password: '', real_name: '', phone: '', roles: '', status: 1, is_agent: 0, region_ids: '', circle_agent_id: 0 }); }
function add() { editing.value = undefined; reset(); dialogOpen.value = true; }
function edit(row: PlatformAdminRow) { editing.value = row; Object.assign(form, { account: row.account, password: '', real_name: row.real_name, phone: row.phone, roles: row.roles, status: row.status, is_agent: row.is_agent, region_ids: row.region_ids, circle_agent_id: row.circle_agent_id }); dialogOpen.value = true; }
async function save() {
  if (!editing.value && (!form.account.trim() || form.password.length < 6)) { ElMessage.warning('账号与至少 6 位的初始密码必填'); return; }
  if (form.is_agent === 1 && !form.region_ids.trim()) { ElMessage.warning('区域管理员必须填写可管理区域 ID，多个以逗号分隔'); return; }
  const payload = { ...form, account: form.account.trim(), real_name: form.real_name.trim(), phone: form.phone.trim(), roles: form.roles.trim(), region_ids: form.region_ids.trim() };
  if (editing.value) await updatePlatformAdmin(editing.value.admin_id, payload); else await createPlatformAdmin(payload);
  ElMessage.success('管理员已保存'); dialogOpen.value = false; await load();
}
onMounted(load);
</script>

<template>
  <Page title="管理员管理" description="平台管理员拥有全量数据；区域管理员必须绑定区域 ID，商户列表、详情和启停接口会按该范围隔离。">
    <el-card shadow="never"><div class="mb-4"><el-button type="primary" @click="add">新增管理员</el-button></div>
      <el-table v-loading="loading" :data="rows" border><el-table-column prop="admin_id" label="ID" width="72" /><el-table-column prop="account" label="账号" min-width="130" /><el-table-column prop="real_name" label="姓名" min-width="110" /><el-table-column prop="phone" label="手机号" width="140" /><el-table-column label="账号类型" width="110"><template #default="{ row }"><el-tag :type="row.is_agent === 1 ? 'warning' : 'success'">{{ row.is_agent === 1 ? '区域管理员' : '平台管理员' }}</el-tag></template></el-table-column><el-table-column prop="region_ids" label="可管理区域" min-width="150"><template #default="{ row }">{{ row.is_agent === 1 ? row.region_ids || '未绑定' : '全平台' }}</template></el-table-column><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'danger'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag></template></el-table-column><el-table-column label="操作" width="90" fixed="right"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button></template></el-table-column></el-table>
      <div class="mt-4 flex justify-end"><el-pagination v-model:current-page="query.page" v-model:page-size="query.limit" layout="total, sizes, prev, pager, next" :total="total" @current-change="load" @size-change="() => { query.page = 1; load(); }" /></div>
    </el-card>
    <el-dialog v-model="dialogOpen" :title="title" width="620px"><el-form label-width="120px"><el-form-item label="登录账号" required><el-input v-model="form.account" :disabled="!!editing" /></el-form-item><el-form-item :label="editing ? '重置密码' : '初始密码'" :required="!editing"><el-input v-model="form.password" show-password type="password" :placeholder="editing ? '留空则不修改' : '至少 6 位'" /></el-form-item><el-form-item label="姓名"><el-input v-model="form.real_name" /></el-form-item><el-form-item label="手机号"><el-input v-model="form.phone" /></el-form-item><el-form-item label="角色 ID"><el-input v-model="form.roles" placeholder="多个角色以逗号分隔" /></el-form-item><el-form-item label="账号类型"><el-radio-group v-model="form.is_agent"><el-radio :value="0">平台管理员</el-radio><el-radio :value="1">区域管理员</el-radio></el-radio-group></el-form-item><el-form-item v-if="form.is_agent === 1" label="可管理区域" required><el-input v-model="form.region_ids" placeholder="例如 1,2,3（商户 region_id）" /></el-form-item><el-form-item v-if="form.is_agent === 1" label="关联代理ID"><el-input-number v-model="form.circle_agent_id" :min="0" /></el-form-item><el-form-item label="启用状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
