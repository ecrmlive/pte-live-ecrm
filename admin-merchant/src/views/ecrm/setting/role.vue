<script setup lang="ts">
import type { MerchantMenuNode, MerchantRoleOption } from '#/api/core/merchant-admin';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElForm, ElFormItem, ElInput, ElMessage, ElSwitch, ElTable, ElTableColumn, ElTag, ElTree } from 'element-plus';
import { nextTick, ref } from 'vue';

import { createMerchantRoleApi, getMerchantMenuTreeApi, listMerchantRoleOptionsApi, updateMerchantRoleApi } from '#/api/core/merchant-admin';

const rows = ref<MerchantRoleOption[]>([]);
const tree = ref<MerchantMenuNode[]>([]);
const treeRef = ref<InstanceType<typeof ElTree>>();
const loading = ref(false);
const saving = ref(false);
const editing = ref<MerchantRoleOption | null>(null);
const roleName = ref('');
const roleStatus = ref(1);


const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

async function load() {
  loading.value = true;
  try {
    const [roles, menus] = await Promise.all([listMerchantRoleOptionsApi(), getMerchantMenuTreeApi()]);
    rows.value = roles.list ?? [];
    tree.value = menus.list ?? [];
  } finally { loading.value = false; }
}

function openCreate() {
  editing.value = null; roleName.value = ''; roleStatus.value = 1; formDrawerApi.setState({ title: '新增角色' }).open();
  void nextTick(() => treeRef.value?.setCheckedKeys([]));
}
function openEdit(row: MerchantRoleOption) {
  if (row.mer_id === 0) { ElMessage.warning('平台预置角色不可修改'); return; }
  editing.value = row; roleName.value = row.role_name; roleStatus.value = row.status; formDrawerApi.setState({ title: '编辑角色权限' }).open();
  const ids = row.rules.split(',').map(Number).filter(Boolean);
  void nextTick(() => treeRef.value?.setCheckedKeys(ids));
}
async function save() {
  if (!roleName.value.trim()) { ElMessage.warning('请填写角色名称'); return; }
  const menuIDs = (treeRef.value?.getCheckedKeys(false) ?? []).map(Number).filter(Boolean);
  saving.value = true;
  formDrawerApi.lock();
  try {
    const data = { role_name: roleName.value.trim(), status: roleStatus.value, menu_ids: menuIDs };
    if (editing.value == null) await createMerchantRoleApi(data);
    else await updateMerchantRoleApi(editing.value.role_id, data);
    formDrawerApi.close(); await load(); ElMessage.success(editing.value == null ? '角色已创建' : '角色权限已保存');
  } finally { saving.value = false; formDrawerApi.unlock(); }
}
async function toggle(row: MerchantRoleOption) {
  if (row.mer_id === 0) return;
  await updateMerchantRoleApi(row.role_id, { role_name: row.role_name, status: row.status, menu_ids: row.rules.split(',').map(Number).filter(Boolean) });
  ElMessage.success('角色状态已保存');
}
void load();
</script>

<template>
  <Page auto-content-height>
    <div class="mb-4"><ElButton type="primary" @click="openCreate">新增角色</ElButton></div>
    <ElTable v-loading="loading" :data="rows" border><ElTableColumn prop="role_id" label="ID" width="80" /><ElTableColumn prop="role_name" label="角色名称" min-width="160" /><ElTableColumn label="来源" width="120"><template #default="{ row }"><ElTag :type="row.mer_id === 0 ? 'info' : 'success'">{{ row.mer_id === 0 ? '平台预置' : '本店自建' }}</ElTag></template></ElTableColumn><ElTableColumn label="状态" width="110"><template #default="{ row }"><ElSwitch v-model="row.status" :active-value="1" :inactive-value="0" :disabled="row.mer_id === 0" @change="toggle(row)" /></template></ElTableColumn><ElTableColumn label="操作" width="90"><template #default="{ row }"><ElButton link type="primary" :disabled="row.mer_id === 0" @click="openEdit(row)">编辑</ElButton></template></ElTableColumn></ElTable>
    <FormDrawer><ElForm label-position="top"><ElFormItem label="角色名称" required><ElInput v-model="roleName" maxlength="32" /></ElFormItem><ElFormItem label="启用角色"><ElSwitch v-model="roleStatus" :active-value="1" :inactive-value="0" /></ElFormItem><ElFormItem label="菜单与按钮权限"><ElTree ref="treeRef" :data="tree" node-key="menu_id" :props="{ children: 'children', label: 'menu_name' }" show-checkbox check-on-click-node default-expand-all /></ElFormItem></ElForm></FormDrawer>
  </Page>
</template>
