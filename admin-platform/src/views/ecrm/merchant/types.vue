<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { deleteMerchantType, fetchMerchantType, fetchMerchantTypes, saveMerchantType, setMerchantTypeRemark, setMerchantTypeStatus, type MerchantTypeRow } from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';

const rows = ref<MerchantTypeRow[]>([]);
const loading = ref(false);
const open = ref(false);
const canManage = ref(false);
const editing = ref<MerchantTypeRow>();
const form = reactive({ name: '', type_info: '', is_margin: false, margin: 0, description: '', remark: '', status: true, menu_codes: [] as string[] });

async function load() { loading.value = true; try { rows.value = (await fetchMerchantTypes()).list || []; } finally { loading.value = false; } }
function reset(row?: MerchantTypeRow) { editing.value = row; Object.assign(form, { name: row?.name || '', type_info: row?.type_info || '', is_margin: row?.is_margin === 1, margin: Number(row?.margin || 0), description: row?.description || '', remark: row?.remark || '', status: row ? row.status === 1 : true, menu_codes: row?.menu_codes || [] }); }
function add() { reset(); open.value = true; }
async function edit(row: MerchantTypeRow) { reset(await fetchMerchantType(row.id)); open.value = true; }
async function save() { if (!form.name.trim() || !form.description.trim() || (form.is_margin && form.margin <= 0)) { ElMessage.warning('请填写类型名称、说明；启用保证金时金额必须大于 0'); return; } await saveMerchantType(editing.value?.id, { ...form, name: form.name.trim(), type_info: form.type_info.trim(), description: form.description.trim(), remark: form.remark.trim() }); open.value = false; ElMessage.success('店铺类型已保存'); await load(); }
async function toggle(row: MerchantTypeRow) { await setMerchantTypeStatus(row.id, row.status !== 1); ElMessage.success('店铺类型状态已更新'); await load(); }
async function mark(row: MerchantTypeRow) { try { const { value } = await ElMessageBox.prompt('填写不超过 500 字的内部备注。', `备注：${row.name}`, { inputValue: row.remark, inputPattern: /^[\s\S]{0,500}$/, inputErrorMessage: '备注不超过 500 字' }); await setMerchantTypeRemark(row.id, value.trim()); ElMessage.success('备注已更新'); await load(); } catch { /* 取消或统一请求层提示 */ } }
async function remove(row: MerchantTypeRow) { try { await ElMessageBox.confirm(`删除“${row.name}”将移除它的店铺菜单授权，是否继续？`, '删除店铺类型', { type: 'warning' }); await deleteMerchantType(row.id); ElMessage.success('店铺类型已删除'); await load(); } catch { /* 取消或统一请求层提示 */ } }
onMounted(async () => { const [codes] = await Promise.all([getAccessCodesApi(), load()]); canManage.value = codes.includes('merchant.type.manage'); });
</script>
<template>
  <Page title="店铺类型" description="维护商户入驻类型、保证金规则、类型说明与店铺菜单授权；仅平台角色可操作。">
    <template #extra><el-button v-if="canManage" type="primary" @click="add">新增类型</el-button></template>
    <el-card shadow="never"><el-table v-loading="loading" :data="rows" row-key="id"><el-table-column label="类型名称" prop="name" min-width="150" /><el-table-column label="类型简介" prop="type_info" min-width="220" /><el-table-column label="保证金" width="130"><template #default="{ row }">{{ row.is_margin ? `¥${Number(row.margin).toFixed(2)}` : '不要求' }}</template></el-table-column><el-table-column label="授权菜单" min-width="220"><template #default="{ row }">{{ row.menu_codes?.join('、') || '未配置' }}</template></el-table-column><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column label="操作" fixed="right" width="240"><template #default="{ row }"><template v-if="canManage"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link type="primary" @click="mark(row)">备注</el-button><el-button link type="primary" @click="toggle(row)">{{ row.status ? '停用' : '启用' }}</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></template></el-table-column></el-table></el-card>
    <el-dialog v-model="open" :title="editing ? '编辑店铺类型' : '新增店铺类型'" width="640px" destroy-on-close><el-form label-width="106px"><el-form-item label="类型名称" required><el-input v-model="form.name" maxlength="128" /></el-form-item><el-form-item label="类型简介"><el-input v-model="form.type_info" maxlength="500" /></el-form-item><el-form-item label="要求保证金"><el-switch v-model="form.is_margin" /></el-form-item><el-form-item v-if="form.is_margin" label="保证金金额" required><el-input-number v-model="form.margin" :min="0.01" :precision="2" /></el-form-item><el-form-item label="类型说明" required><el-input v-model="form.description" :rows="4" maxlength="65535" show-word-limit type="textarea" /></el-form-item><el-form-item label="店铺菜单授权"><el-select v-model="form.menu_codes" multiple allow-create filterable default-first-option class="w-full" placeholder="输入统一菜单代码，如 merchant.catalog"><el-option label="商户控制台" value="merchant.dashboard" /><el-option label="商品管理" value="merchant.catalog" /><el-option label="订单管理" value="merchant.order" /></el-select></el-form-item><el-form-item label="内部备注"><el-input v-model="form.remark" :rows="2" maxlength="500" type="textarea" /></el-form-item><el-form-item label="启用状态"><el-switch v-model="form.status" /></el-form-item></el-form><template #footer><el-button @click="open = false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
