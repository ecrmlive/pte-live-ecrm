<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import {
  createPlatformRole,
  fetchPlatformMenuTree,
  fetchPlatformRoles,
  updatePlatformRole,
  type PlatformMenuNode,
  type PlatformRoleRow,
} from '#/api/core/ecrm';

const rows = ref<PlatformRoleRow[]>([]);
const menus = ref<PlatformMenuNode[]>([]);
const total = ref(0);
const loading = ref(false);
const dialogOpen = ref(false);
const editing = ref<PlatformRoleRow>();
const menuTree = ref();
const query = reactive({ page: 1, limit: 50 });
const form = reactive({ code: '', role_name: '', status: 1, is_agent: 0, circle_id: 0 });
const title = computed(() => editing.value ? '编辑角色' : '新增角色');
const roleType = (value: number) => value === 1 ? '区域角色' : '平台角色';

function roleIDs(row?: PlatformRoleRow) { return row?.rules ? row.rules.split(',').map(Number).filter(Boolean) : []; }
async function load() { loading.value = true; try { const result = await fetchPlatformRoles(query); rows.value = result.list; total.value = result.total; } finally { loading.value = false; } }
async function ensureMenus() { if (!menus.value.length) menus.value = await fetchPlatformMenuTree(); }
async function add() { editing.value = undefined; Object.assign(form, { code: '', role_name: '', status: 1, is_agent: 0, circle_id: 0 }); await ensureMenus(); dialogOpen.value = true; requestAnimationFrame(() => menuTree.value?.setCheckedKeys([])); }
async function edit(row: PlatformRoleRow) { editing.value = row; Object.assign(form, { code: row.code, role_name: row.role_name, status: row.status, is_agent: row.is_agent, circle_id: row.circle_id }); await ensureMenus(); dialogOpen.value = true; requestAnimationFrame(() => menuTree.value?.setCheckedKeys(roleIDs(row))); }
async function save() { if (!form.role_name.trim() || (!editing.value && !/^[a-z0-9_]{1,32}$/.test(form.code))) { ElMessage.warning('请填写角色名称和小写角色代码'); return; } const menu_ids = menuTree.value?.getCheckedKeys(false) || []; if (editing.value) await updatePlatformRole(editing.value.role_id, { role_name: form.role_name.trim(), status: form.status, menu_ids }); else await createPlatformRole({ code: form.code.trim(), role_name: form.role_name.trim(), status: form.status, menu_ids }); ElMessage.success('角色权限已保存'); dialogOpen.value = false; await load(); }
onMounted(load);
</script>

<template>
  <Page title="角色权限" description="平台角色管理全平台菜单；区域角色绑定商圈后用于区域管理员授权。菜单树同时包含页面和按钮权限。">
    <el-card shadow="never"><div class="mb-4"><el-button type="primary" @click="add">新增角色</el-button></div><el-table v-loading="loading" :data="rows" border><el-table-column prop="role_id" label="ID" width="72" /><el-table-column prop="code" label="角色代码" min-width="140" /><el-table-column prop="role_name" label="角色名称" min-width="160" /><el-table-column label="类型" width="120"><template #default="{ row }"><el-tag :type="row.is_agent === 1 ? 'warning' : 'success'">{{ roleType(row.is_agent) }}</el-tag></template></el-table-column><el-table-column label="权限数" width="100"><template #default="{ row }">{{ roleIDs(row).length }}</template></el-table-column><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'danger'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag></template></el-table-column><el-table-column label="操作" width="90"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination v-model:current-page="query.page" v-model:page-size="query.limit" layout="total, sizes, prev, pager, next" :total="total" @current-change="load" @size-change="() => { query.page = 1; load(); }" /></div></el-card>
    <el-dialog v-model="dialogOpen" :title="title" width="720px"><el-form label-width="110px"><el-form-item label="角色代码" required><el-input v-model="form.code" :disabled="!!editing" placeholder="如 operations、region_custom" /></el-form-item><el-form-item label="角色名称" required><el-input v-model="form.role_name" /></el-form-item><el-form-item label="状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item><el-form-item label="菜单与按钮"><div class="w-full rounded border p-3"><el-tree ref="menuTree" :data="menus" :default-expand-all="false" :props="{ label: 'menu_name', children: 'children' }" check-on-click-node node-key="menu_id" show-checkbox /></div></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
