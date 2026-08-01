<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';
import { createPlatformBrandApi, deletePlatformBrandApi, listPlatformBrandsApi, updatePlatformBrandApi, type PlatformBrand } from '#/api/core/platform-catalog';
const rows = ref<PlatformBrand[]>([]); const loading = ref(false); const open = ref(false); const editing = ref<PlatformBrand>(); const form = reactive({ brand_name: '', is_show: 1, sort: 0 });
async function load() { loading.value = true; try { rows.value = (await listPlatformBrandsApi()).list || []; } finally { loading.value = false; } }
function add() { editing.value = undefined; Object.assign(form, { brand_name: '', is_show: 1, sort: 0 }); open.value = true; }
function edit(row: PlatformBrand) { editing.value = row; Object.assign(form, { brand_name: row.brand_name, is_show: row.is_show, sort: row.sort }); open.value = true; }
async function save() { if (!form.brand_name.trim()) { ElMessage.warning('请填写品牌名称'); return; } if (editing.value) await updatePlatformBrandApi(editing.value.brand_id, { brand_name: form.brand_name.trim(), is_show: form.is_show, sort: form.sort }); else await createPlatformBrandApi({ brand_name: form.brand_name.trim(), is_show: form.is_show, sort: form.sort }); open.value = false; ElMessage.success('品牌已保存'); await load(); }
async function remove(row: PlatformBrand) { try { await ElMessageBox.confirm(`删除品牌“${row.brand_name}”后不可恢复，是否继续？`, '删除品牌', { type: 'warning' }); await deletePlatformBrandApi(row.brand_id); ElMessage.success('品牌已删除'); await load(); } catch { /* 取消或请求错误统一提示 */ } }
onMounted(() => void load());
</script>
<template>
  <Page title="品牌管理" description="维护平台品牌；新增、编辑、删除均受 product/brand/manage 按钮权限控制。"><template #extra><el-button type="primary" @click="add">新增品牌</el-button></template><el-card shadow="never"><el-table v-loading="loading" :data="rows" row-key="brand_id"><el-table-column label="ID" prop="brand_id" width="90" /><el-table-column label="品牌名称" min-width="240" prop="brand_name" /><el-table-column label="排序" prop="sort" width="90" /><el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '显示' : '隐藏' }}</el-tag></template></el-table-column><el-table-column label="操作" width="150"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table></el-card><el-dialog v-model="open" :title="editing ? '编辑品牌' : '新增品牌'" width="480px" destroy-on-close><el-form label-width="88px"><el-form-item label="品牌名称" required><el-input v-model="form.brand_name" /></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" class="w-full" /></el-form-item><el-form-item label="显示"><el-switch v-model="form.is_show" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="open = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template></el-dialog></Page>
</template>
