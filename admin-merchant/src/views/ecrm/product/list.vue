<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createMerchantProductApi,
  deleteMerchantProductApi,
  getMerchantProductCategoriesApi,
  listMerchantProductsApi,
  listMerchantProductRecycleBinApi,
  restoreMerchantProductApi,
  setMerchantProductShowApi,
  setMerchantProductStockApi,
  updateMerchantProductApi,
  type MerchantCategoryNode,
  type MerchantProduct,
  type MerchantRecycleProduct,
  type MerchantProductSaveInput,
} from '#/api/core/merchant-catalog';
import { getAccessCodesApi } from '#/api/core/auth';

const loading = ref(false);
const saving = ref(false);
const rows = ref<MerchantProduct[]>([]);
const categories = ref<MerchantCategoryNode[]>([]);
const total = ref(0);
const dialogOpen = ref(false);
const stockOpen = ref(false);
const recycleOpen = ref(false);
const canCreate = ref(false);
const canUpdate = ref(false);
const canDelete = ref(false);
const canShow = ref(false);
const canStock = ref(false);
const canRestore = ref(false);
const editingID = ref<number>();
const stockRow = ref<MerchantProduct>();
const stockValue = ref(0);
const recycleRows = ref<MerchantRecycleProduct[]>([]);
const recycleTotal = ref(0);
const recycleQuery = reactive({ page: 1, limit: 20 });
const query = reactive({ keyword: '', page: 1, limit: 20, status: undefined as number | undefined });
const form = reactive<MerchantProductSaveInput>({
  cate_id: 0, cost: 0, image: '', is_show: 1, keyword: '', ot_price: 0, price: 0,
  slider_image: '', spec_type: 0, stock: 0, store_info: '', store_name: '', type: 0, unit_name: '件',
});

const categoryOptions = computed(() => flattenCategories(categories.value));
const statusText = (status: number) => ({ '-2': '已下架', '-1': '审核驳回', 0: '待审核', 1: '已通过' }[status] || '未知');
const statusType = (status: number) => ({ '-2': 'info', '-1': 'danger', 0: 'warning', 1: 'success' }[status] || 'info');

function flattenCategories(nodes: MerchantCategoryNode[], prefix = ''): Array<{ label: string; value: number }> {
  return nodes.flatMap((node) => [
    { label: `${prefix}${node.cate_name}`, value: node.store_category_id },
    ...flattenCategories(node.children || [], `${prefix}— `),
  ]);
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, { cate_id: 0, cost: 0, image: '', is_show: 1, keyword: '', ot_price: 0, price: 0, slider_image: '', spec_type: 0, stock: 0, store_info: '', store_name: '', type: 0, unit_name: '件' });
}

function openCreate() { resetForm(); dialogOpen.value = true; }
function openEdit(row: MerchantProduct) {
  editingID.value = row.product_id;
  Object.assign(form, { cate_id: row.cate_id, cost: -1, image: row.image, is_show: row.is_show, keyword: row.keyword, ot_price: row.ot_price, price: row.price, slider_image: row.image, spec_type: 0, stock: row.stock, store_info: row.store_info, store_name: row.store_name, type: 0, unit_name: row.unit_name || '件' });
  dialogOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const result = await listMerchantProductsApi({ ...query, keyword: query.keyword.trim() || undefined });
    rows.value = result.list;
    total.value = result.total;
  } finally { loading.value = false; }
}

async function loadCategories() {
  const result = await getMerchantProductCategoriesApi();
  categories.value = result.list;
}

function search() { query.page = 1; void load(); }
function reset() { query.keyword = ''; query.status = undefined; query.page = 1; void load(); }

async function save() {
  if (!form.store_name.trim()) { ElMessage.warning('请填写商品名称'); return; }
  if (!form.cate_id) { ElMessage.warning('请选择平台分类'); return; }
  if (form.price < 0 || form.stock < 0) { ElMessage.warning('价格和库存不能小于 0'); return; }
  saving.value = true;
  try {
    if (editingID.value) await updateMerchantProductApi(editingID.value, form);
    else await createMerchantProductApi(form);
    dialogOpen.value = false;
    ElMessage.success(editingID.value ? '商品已更新，可能需重新审核' : '商品已创建');
    await load();
  } finally { saving.value = false; }
}

async function changeShow(row: MerchantProduct) {
  const next = row.is_show === 1 ? 0 : 1;
  try {
    await setMerchantProductShowApi(row.product_id, next);
    row.is_show = next;
    ElMessage.success(next ? '商品已上架' : '商品已下架');
  } catch {
    ElMessage.error('商品状态更新失败');
  }
}

function openStock(row: MerchantProduct) {
  stockRow.value = row;
  stockValue.value = row.stock;
  stockOpen.value = true;
}

async function saveStock() {
  if (!stockRow.value || stockValue.value < 0) return;
  await setMerchantProductStockApi(stockRow.value.product_id, stockValue.value);
  stockRow.value.stock = stockValue.value;
  stockOpen.value = false;
  ElMessage.success('库存已更新');
}

async function loadRecycleBin() {
  const result = await listMerchantProductRecycleBinApi(recycleQuery);
  recycleRows.value = result.list || [];
  recycleTotal.value = result.total || 0;
}

async function openRecycleBin() {
  recycleQuery.page = 1;
  await loadRecycleBin();
  recycleOpen.value = true;
}

async function restore(row: MerchantRecycleProduct) {
  await restoreMerchantProductApi(row.product_id);
  ElMessage.success('商品已恢复，等待平台重新审核');
  await Promise.all([loadRecycleBin(), load()]);
}

async function remove(row: MerchantProduct) {
  try {
    await ElMessageBox.confirm(`将商品“${row.store_name}”移入回收站，30 天内可恢复，是否继续？`, '移入回收站', { type: 'warning' });
    await deleteMerchantProductApi(row.product_id);
    ElMessage.success('商品已移入回收站');
    await load();
  } catch {
    // 用户取消不提示错误。
  }
}

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load(), loadCategories()]);
  canCreate.value = permissions.includes('product.create');
  canUpdate.value = permissions.includes('product.update');
  canDelete.value = permissions.includes('product.delete');
  canShow.value = permissions.includes('product.show');
  canStock.value = permissions.includes('product.stock');
  canRestore.value = permissions.includes('product.restore');
});
</script>

<template>
  <Page title="商品列表" description="管理本店商品、库存与上架状态；商品编辑会按店铺审核规则重新进入审核流程。">
    <template #extra>
      <el-button v-if="canRestore" @click="openRecycleBin">回收站</el-button>
      <el-button v-if="canCreate" type="primary" @click="openCreate">新增商品</el-button>
    </template>
    <el-card shadow="never"><el-form class="grid gap-x-4 md:grid-cols-3" label-width="72px" @submit.prevent="search"><el-form-item label="商品搜索"><el-input v-model="query.keyword" clearable placeholder="商品名称 / 关键词" @keyup.enter="search" /></el-form-item><el-form-item label="审核状态"><el-select v-model="query.status" clearable class="w-full" placeholder="全部"><el-option label="待审核" :value="0" /><el-option label="已通过" :value="1" /><el-option label="已驳回" :value="-1" /><el-option label="已下架" :value="-2" /></el-select></el-form-item><el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button></el-form-item></el-form></el-card>
    <el-card class="mt-4" shadow="never"><el-table v-loading="loading" :data="rows" row-key="product_id"><el-table-column label="ID" prop="product_id" width="80" /><el-table-column label="商品名称" min-width="190" prop="store_name" show-overflow-tooltip /><el-table-column label="分类" min-width="120" prop="cate_name" /><el-table-column label="售价" width="110"><template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template></el-table-column><el-table-column label="库存" prop="stock" width="90" /><el-table-column label="审核" width="96"><template #default="{ row }"><el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag></template></el-table-column><el-table-column label="上架" width="88"><template #default="{ row }"><el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '上架中' : '已下架' }}</el-tag></template></el-table-column><el-table-column label="驳回原因" min-width="160" prop="refusal" show-overflow-tooltip /><el-table-column label="创建时间" min-width="170" prop="create_time" /><el-table-column fixed="right" label="操作" width="248"><template #default="{ row }"><el-button v-if="canUpdate" link type="primary" @click="openEdit(row)">编辑</el-button><el-button v-if="canStock" link type="primary" @click="openStock(row)">库存</el-button><el-button v-if="canShow" link type="warning" @click="changeShow(row)">{{ row.is_show === 1 ? '下架' : '上架' }}</el-button><el-button v-if="canDelete" link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑商品' : '新增商品'" width="760px" destroy-on-close><el-form class="grid grid-cols-2 gap-x-4" label-width="88px"><el-form-item label="商品名称" required><el-input v-model="form.store_name" /></el-form-item><el-form-item label="平台分类" required><el-select v-model="form.cate_id" filterable class="w-full"><el-option v-for="item in categoryOptions" :key="item.value" :label="item.label" :value="item.value" /></el-select></el-form-item><el-form-item label="销售价"><el-input-number v-model="form.price" :min="0" :precision="2" class="w-full" /></el-form-item><el-form-item label="划线价"><el-input-number v-model="form.ot_price" :min="0" :precision="2" class="w-full" /></el-form-item><el-form-item label="库存"><el-input-number v-model="form.stock" :min="0" class="w-full" /></el-form-item><el-form-item label="单位"><el-input v-model="form.unit_name" /></el-form-item><el-form-item class="col-span-2" label="关键词"><el-input v-model="form.keyword" /></el-form-item><el-form-item class="col-span-2" label="主图"><el-input v-model="form.image" placeholder="图片 URL（素材库接入后可直接选择）" /></el-form-item><el-form-item class="col-span-2" label="商品简介"><el-input v-model="form.store_info" :rows="4" type="textarea" /></el-form-item><el-form-item label="初始上架"><el-switch v-model="form.is_show" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
    <el-dialog v-model="stockOpen" title="调整库存" width="380px" destroy-on-close><el-form label-width="72px"><el-form-item label="商品"><span>{{ stockRow?.store_name }}</span></el-form-item><el-form-item label="库存"><el-input-number v-model="stockValue" :min="0" class="w-full" /></el-form-item></el-form><template #footer><el-button @click="stockOpen = false">取消</el-button><el-button type="primary" @click="saveStock">保存</el-button></template></el-dialog>
    <el-dialog v-model="recycleOpen" title="商品回收站" width="760px" destroy-on-close><el-alert class="mb-4" :closable="false" show-icon title="恢复后的商品会重新进入平台审核，超过恢复期限后不可恢复。" type="info" /><el-table :data="recycleRows" row-key="product_id"><el-table-column label="商品" min-width="240" prop="store_name" /><el-table-column label="移入时间" min-width="170" prop="deleted_at" /><el-table-column label="可恢复至" min-width="170" prop="restore_until" /><el-table-column fixed="right" label="操作" width="100"><template #default="{ row }"><el-button link type="primary" @click="restore(row)">恢复</el-button></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="recycleQuery.page" :page-size="recycleQuery.limit" :total="recycleTotal" background layout="total, prev, pager, next" @current-change="(page: number) => { recycleQuery.page = page; loadRecycleBin(); }" /></div></el-dialog>
  </Page>
</template>
