<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createProductAttrTemplateApi,
  deleteProductAttrTemplateApi,
  listProductAttrTemplatesApi,
  updateProductAttrTemplateApi,
  type ProductAttrTemplate,
  type ProductAttrTemplateInput,
} from '#/api/core/product-meta';

const loading = ref(false);
const saving = ref(false);
const rows = ref<ProductAttrTemplate[]>([]);
const total = ref(0);
const dialogOpen = ref(false);
const editingID = ref<number>();
const query = reactive({ page: 1, limit: 20 });
const form = reactive<ProductAttrTemplateInput>({ template_name: '', template_value: '[]', sort: 0 });

function resetForm() {
  editingID.value = undefined;
  form.template_name = '';
  form.template_value = '[]';
  form.sort = 0;
}

function openCreate() { resetForm(); dialogOpen.value = true; }
function openEdit(row: ProductAttrTemplate) {
  editingID.value = row.template_id;
  form.template_name = row.template_name;
  form.template_value = row.template_value || '[]';
  form.sort = row.sort;
  dialogOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const data = await listProductAttrTemplatesApi(query);
    rows.value = data.list;
    total.value = data.total;
  } finally { loading.value = false; }
}

function validateTemplate() {
  try {
    const parsed = JSON.parse(form.template_value || '[]');
    return Array.isArray(parsed) || typeof parsed === 'object';
  } catch { return false; }
}

async function save() {
  if (!form.template_name.trim()) { ElMessage.warning('请填写模板名称'); return; }
  if (!validateTemplate()) { ElMessage.warning('参数定义必须是合法 JSON 对象或数组'); return; }
  saving.value = true;
  try {
    if (editingID.value) await updateProductAttrTemplateApi(editingID.value, form);
    else await createProductAttrTemplateApi(form);
    dialogOpen.value = false;
    ElMessage.success('保存成功');
    await load();
  } finally { saving.value = false; }
}

async function remove(row: ProductAttrTemplate) {
  try {
    await ElMessageBox.confirm(`删除参数模板“${row.template_name}”后不可恢复，是否继续？`, '删除确认', { type: 'warning' });
    await deleteProductAttrTemplateApi(row.template_id);
    ElMessage.success('已删除');
    await load();
  } catch {
    // 用户取消不提示错误。
  }
}

onMounted(load);
</script>

<template>
  <Page title="商品参数模板" description="维护商品规格/参数模板，在新增商品时复用；模板仅属于当前商户。">
    <template #extra><el-button type="primary" @click="openCreate">新增模板</el-button></template>
    <el-card shadow="never"><el-table v-loading="loading" :data="rows" row-key="template_id"><el-table-column label="ID" prop="template_id" width="88" /><el-table-column label="模板名称" min-width="180" prop="template_name" /><el-table-column label="参数定义" min-width="360" prop="template_value" show-overflow-tooltip /><el-table-column label="排序" prop="sort" width="90" /><el-table-column label="创建时间" min-width="170" prop="create_time" /><el-table-column fixed="right" label="操作" width="128"><template #default="{ row }"><el-button link type="primary" @click="openEdit(row)">编辑</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑参数模板' : '新增参数模板'" width="680px" destroy-on-close><el-form label-width="92px"><el-form-item label="模板名称" required><el-input v-model="form.template_name" maxlength="64" /></el-form-item><el-form-item label="参数定义" required><el-input v-model="form.template_value" :rows="10" type="textarea" /><div class="text-xs text-gray-500">使用 JSON 对象或数组保存规格定义，例如 [{&quot;name&quot;:&quot;颜色&quot;,&quot;values&quot;:[&quot;黑&quot;,&quot;白&quot;]}]。</div></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
