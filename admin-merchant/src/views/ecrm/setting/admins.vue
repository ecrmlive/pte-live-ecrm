<script setup lang="ts">
import type { MerchantAdmin, MerchantAdminSaveInput, MerchantRoleOption } from '#/api/core/merchant-admin';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElForm, ElFormItem, ElInput, ElMessage, ElPagination, ElSelect, ElOption, ElSwitch, ElTable, ElTableColumn, ElTag } from 'element-plus';
import { computed, reactive, ref } from 'vue';

import { createMerchantAdminApi, listMerchantAdminsApi, listMerchantRoleOptionsApi, updateMerchantAdminApi } from '#/api/core/merchant-admin';

const rows = ref<MerchantAdmin[]>([]);
const roles = ref<MerchantRoleOption[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = 20;
const total = ref(0);
const saving = ref(false);
const editingID = ref<number | null>(null);
const selectedRoleIDs = ref<number[]>([]);
const form = reactive<Required<MerchantAdminSaveInput>>({ account: '', password: '', phone: '', real_name: '', roles: '', status: 1 });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

const roleNames = computed(() => new Map(roles.value.map((role) => [role.role_id, role.role_name])));

function displayRoles(value: string) { return value.split(',').map(Number).filter(Boolean).map((id) => roleNames.value.get(id) || `角色#${id}`).join('、'); }
async function load() {
  loading.value = true;
  try {
    const [admins, roleResult] = await Promise.all([listMerchantAdminsApi({ page: page.value, limit: pageSize }), listMerchantRoleOptionsApi()]);
    rows.value = admins.list ?? []; total.value = admins.total ?? 0; roles.value = (roleResult.list ?? []).filter((role) => role.status === 1);
  } finally { loading.value = false; }
}
function resetForm() { form.account = ''; form.password = ''; form.phone = ''; form.real_name = ''; form.roles = ''; form.status = 1; selectedRoleIDs.value = []; }
function openCreate() { editingID.value = null; resetForm(); formDrawerApi.setState({ title: '新增子账号' }).open(); }
function openEdit(row: MerchantAdmin) { editingID.value = row.merchant_admin_id; form.account = row.account; form.password = ''; form.phone = row.phone; form.real_name = row.real_name; form.status = row.status; selectedRoleIDs.value = row.roles.split(',').map(Number).filter(Boolean); formDrawerApi.setState({ title: '编辑子账号' }).open(); }
async function save() {
  if (!form.real_name.trim()) { ElMessage.warning('请填写姓名'); return; }
  if (!selectedRoleIDs.value.length) { ElMessage.warning('请至少选择一个角色'); return; }
  if (editingID.value == null && (!form.account.trim() || form.password.length < 6)) { ElMessage.warning('请填写账号和至少 6 位密码'); return; }
  saving.value = true;
  formDrawerApi.lock();
  try {
    const data: MerchantAdminSaveInput = { ...form, roles: selectedRoleIDs.value.join(',') };
    if (editingID.value == null) await createMerchantAdminApi(data);
    else { delete data.account; if (!data.password) delete data.password; await updateMerchantAdminApi(editingID.value, data); }
    formDrawerApi.close(); await load(); ElMessage.success(editingID.value == null ? '子账号已创建' : '子账号已更新');
  } finally { saving.value = false; formDrawerApi.unlock(); }
}
async function toggle(row: MerchantAdmin) { await updateMerchantAdminApi(row.merchant_admin_id, { status: row.status }); ElMessage.success('状态已保存'); }
void load();
</script>

<template>
  <Page auto-content-height>
    <div class="mb-4"><ElButton type="primary" @click="openCreate">新增子账号</ElButton></div>
    <ElTable v-loading="loading" :data="rows" border><ElTableColumn prop="merchant_admin_id" label="ID" width="80" /><ElTableColumn prop="account" label="账号" min-width="130" /><ElTableColumn prop="real_name" label="姓名" min-width="120" /><ElTableColumn prop="phone" label="手机号" min-width="130" /><ElTableColumn label="角色" min-width="180"><template #default="{ row }">{{ displayRoles(row.roles) || '-' }}</template></ElTableColumn><ElTableColumn label="身份" width="100"><template #default="{ row }"><ElTag :type="row.level === 0 ? 'success' : 'info'">{{ row.level === 0 ? '主账号' : '子账号' }}</ElTag></template></ElTableColumn><ElTableColumn label="状态" width="100"><template #default="{ row }"><ElSwitch v-model="row.status" :active-value="1" :inactive-value="0" :disabled="row.level === 0" @change="toggle(row)" /></template></ElTableColumn><ElTableColumn label="操作" width="90" fixed="right"><template #default="{ row }"><ElButton link type="primary" @click="openEdit(row)">编辑</ElButton></template></ElTableColumn></ElTable>
    <ElPagination v-if="total > pageSize" v-model:current-page="page" class="mt-4" :page-size="pageSize" :total="total" background layout="prev, pager, next" @current-change="load" />
    <FormDrawer><ElForm label-position="top"><ElFormItem label="登录账号" required><ElInput v-model="form.account" :disabled="editingID != null" /></ElFormItem><ElFormItem label="姓名" required><ElInput v-model="form.real_name" /></ElFormItem><ElFormItem label="手机号"><ElInput v-model="form.phone" /></ElFormItem><ElFormItem :label="editingID == null ? '登录密码' : '登录密码（留空不修改）'" :required="editingID == null"><ElInput v-model="form.password" type="password" show-password /></ElFormItem><ElFormItem label="角色" required><ElSelect v-model="selectedRoleIDs" multiple class="w-full" placeholder="请选择角色"><ElOption v-for="role in roles" :key="role.role_id" :label="role.role_name" :value="role.role_id" /></ElSelect></ElFormItem><ElFormItem label="启用账号"><ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" /></ElFormItem></ElForm></FormDrawer>
  </Page>
</template>
