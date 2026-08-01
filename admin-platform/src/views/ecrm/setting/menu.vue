<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import { fetchPlatformMenus, updatePlatformMenu, type PlatformMenuRow } from '#/api/core/ecrm';

const rows = ref<PlatformMenuRow[]>([]);
const loading = ref(false);
const query = reactive({ keyword: '', type: undefined as number | undefined });
const dialogOpen = ref(false);
const editing = ref<PlatformMenuRow>();
const form = reactive({ menu_name: '', sort: 0, is_show: 1 });
const filtered = computed(() => rows.value.filter((row) => (!query.keyword || `${row.menu_name}${row.path}${row.route}`.toLowerCase().includes(query.keyword.toLowerCase())) && (query.type === undefined || row.is_menu === query.type)));
const typeText = (value: number) => value === 2 ? '按钮权限' : '菜单页面';

async function load() { loading.value = true; try { rows.value = await fetchPlatformMenus(); } finally { loading.value = false; } }
function edit(row: PlatformMenuRow) { editing.value = row; Object.assign(form, { menu_name: row.menu_name, sort: row.sort, is_show: row.is_show }); dialogOpen.value = true; }
async function save() { if (!editing.value || !form.menu_name.trim()) { ElMessage.warning('菜单名称不能为空'); return; } await updatePlatformMenu(editing.value.menu_id, { menu_name: form.menu_name.trim(), sort: form.sort, is_show: form.is_show }); ElMessage.success('菜单已保存'); dialogOpen.value = false; await load(); }
onMounted(load);
</script>

<template>
  <Page title="菜单管理" description="维护平台后台菜单与按钮权限的名称、排序和显示状态；路由由功能模块注册，不能在此页任意改写。">
    <el-card shadow="never"><el-form class="grid gap-x-4 md:grid-cols-3" label-width="72px"><el-form-item label="搜索"><el-input v-model="query.keyword" clearable placeholder="名称 / 路径 / 路由" /></el-form-item><el-form-item label="类型"><el-select v-model="query.type" clearable class="w-full" placeholder="全部"><el-option label="菜单页面" :value="1" /><el-option label="按钮权限" :value="2" /></el-select></el-form-item><el-form-item><el-button @click="load">刷新</el-button></el-form-item></el-form>
      <el-table v-loading="loading" :data="filtered" border><el-table-column prop="menu_id" label="ID" width="76" /><el-table-column prop="menu_name" label="名称" min-width="160" /><el-table-column prop="path" label="路径" min-width="210" show-overflow-tooltip /><el-table-column prop="route" label="路由" min-width="150" show-overflow-tooltip /><el-table-column label="类型" width="100"><template #default="{ row }"><el-tag :type="row.is_menu === 2 ? 'warning' : 'success'">{{ typeText(row.is_menu) }}</el-tag></template></el-table-column><el-table-column prop="sort" label="排序" width="80" /><el-table-column label="显示" width="80"><template #default="{ row }"><el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '显示' : '隐藏' }}</el-tag></template></el-table-column><el-table-column label="操作" width="90" fixed="right"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button></template></el-table-column></el-table>
    </el-card>
    <el-dialog v-model="dialogOpen" title="编辑菜单" width="520px"><el-form label-width="100px"><el-form-item label="菜单 ID"><el-input :model-value="String(editing?.menu_id || '')" disabled /></el-form-item><el-form-item label="菜单名称"><el-input v-model="form.menu_name" /></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" /></el-form-item><el-form-item label="显示状态"><el-switch v-model="form.is_show" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
