<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createProductGuaranteeApi,
  deleteProductGuaranteeApi,
  listProductGuaranteesApi,
  updateProductGuaranteeApi,
  type ProductGuarantee,
  type ProductGuaranteeInput,
} from '#/api/core/product-meta';

const loading = ref(false);
const saving = ref(false);
const rows = ref<ProductGuarantee[]>([]);
const total = ref(0);
const dialogOpen = ref(false);
const editingID = ref<number>();
const query = reactive({ page: 1, limit: 20 });
const form = reactive<ProductGuaranteeInput>({ name: '', content: '', sort: 0, status: 1 });

function resetForm() {
  editingID.value = undefined;
  form.name = '';
  form.content = '';
  form.sort = 0;
  form.status = 1;
}

function openCreate() { resetForm(); dialogOpen.value = true; }
function openEdit(row: ProductGuarantee) {
  editingID.value = row.guarantee_id;
  form.name = row.name;
  form.content = row.content;
  form.sort = row.sort;
  form.status = row.status;
  dialogOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const data = await listProductGuaranteesApi(query);
    rows.value = data.list;
    total.value = data.total;
  } finally { loading.value = false; }
}

async function save() {
  if (!form.name.trim()) { ElMessage.warning('请填写保障名称'); return; }
  saving.value = true;
  try {
    if (editingID.value) await updateProductGuaranteeApi(editingID.value, form);
    else await createProductGuaranteeApi(form);
    dialogOpen.value = false;
    ElMessage.success('保存成功');
    await load();
  } finally { saving.value = false; }
}

async function remove(row: ProductGuarantee) {
  try {
    await ElMessageBox.confirm(`删除保障服务“${row.name}”后不可恢复，是否继续？`, '删除确认', { type: 'warning' });
    await deleteProductGuaranteeApi(row.guarantee_id);
    ElMessage.success('已删除');
    await load();
  } catch {
    // 用户取消不提示错误。
  }
}

async function changeStatus(row: ProductGuarantee) {
  try {
    await updateProductGuaranteeApi(row.guarantee_id, { name: row.name, content: row.content, sort: row.sort, status: row.status });
    ElMessage.success('状态已更新');
  } catch { row.status = row.status === 1 ? 0 : 1; }
}

onMounted(load);
</script>

<template>
  <Page title="保障服务" description="维护商品售后保障说明；仅启用的保障服务可在本店商品中选择。">
    <template #extra><el-button type="primary" @click="openCreate">新增保障</el-button></template>
    <el-card shadow="never"><el-table v-loading="loading" :data="rows" row-key="guarantee_id"><el-table-column label="ID" prop="guarantee_id" width="88" /><el-table-column label="保障名称" min-width="160" prop="name" /><el-table-column label="服务说明" min-width="260" prop="content" show-overflow-tooltip /><el-table-column label="排序" prop="sort" width="90" /><el-table-column label="启用" width="100"><template #default="{ row }"><el-switch v-model="row.status" :active-value="1" :inactive-value="0" @change="changeStatus(row)" /></template></el-table-column><el-table-column label="创建时间" min-width="170" prop="create_time" /><el-table-column fixed="right" label="操作" width="128"><template #default="{ row }"><el-button link type="primary" @click="openEdit(row)">编辑</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑保障服务' : '新增保障服务'" width="560px" destroy-on-close><el-form label-width="88px"><el-form-item label="保障名称" required><el-input v-model="form.name" maxlength="64" /></el-form-item><el-form-item label="服务说明"><el-input v-model="form.content" :rows="5" maxlength="1000" show-word-limit type="textarea" /></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item><el-form-item label="启用"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
