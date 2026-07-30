<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createMerchantPresellActiveApi,
  deleteMerchantPresellActiveApi,
  listMerchantPresellActivesApi,
  setMerchantPresellShowApi,
  updateMerchantPresellActiveApi,
  type MerchantPresellActive,
  type MerchantPresellSaveInput,
} from '#/api/core/merchant-presell';
import {
  listMerchantProductsApi,
  type MerchantProduct,
} from '#/api/core/merchant-catalog';

const loading = ref(false);
const saving = ref(false);
const dialogOpen = ref(false);
const editingID = ref<number>();
const rows = ref<MerchantPresellActive[]>([]);
const products = ref<MerchantProduct[]>([]);
const total = ref(0);
const query = reactive({ keyword: '', limit: 20, page: 1 });
const form = reactive({
  activity_dates: [] as string[],
  down_price: 0,
  final_dates: [] as string[],
  final_price: 0,
  is_show: 1,
  presell_type: 1,
  price: 0,
  product_id: 0,
  stock: 100,
  store_info: '',
  store_name: '',
});

const selectedProduct = computed(() => products.value.find((item) => item.product_id === form.product_id));
const filteredRows = computed(() => {
  const keyword = query.keyword.trim().toLowerCase();
  if (!keyword) return rows.value;
  return rows.value.filter((row) => row.store_name.toLowerCase().includes(keyword));
});

function dateText(value: unknown) {
  return String(value || '').replace('T', ' ').slice(0, 19);
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, { activity_dates: [], down_price: 0, final_dates: [], final_price: 0, is_show: 1, presell_type: 1, price: 0, product_id: 0, stock: 100, store_info: '', store_name: '' });
}

function openCreate() {
  resetForm();
  dialogOpen.value = true;
}

function openEdit(row: MerchantPresellActive) {
  editingID.value = row.product_presell_id;
  Object.assign(form, {
    activity_dates: [dateText(row.start_time), dateText(row.end_time)],
    down_price: row.down_price,
    final_dates: [dateText(row.final_start_time), dateText(row.final_end_time)].filter(Boolean),
    final_price: row.final_price,
    is_show: row.is_show,
    presell_type: row.presell_type,
    price: row.price,
    product_id: row.product_id,
    stock: row.stock,
    store_info: row.store_info,
    store_name: row.store_name,
  });
  dialogOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const data = await listMerchantPresellActivesApi({ limit: query.limit, page: query.page });
    rows.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
}

async function loadProducts() {
  const data = await listMerchantProductsApi({ limit: 100, page: 1, status: 1 });
  products.value = data.list || [];
}

function selectedProductChanged(productID: number) {
  const product = products.value.find((item) => item.product_id === productID);
  if (!product || editingID.value) return;
  form.store_name = `${product.store_name} · 预售`;
  form.price = product.price;
  form.stock = product.stock;
}

function validDeposit() {
  return form.down_price > 0 && form.final_price > 0 && Math.abs(form.down_price + form.final_price - form.price) < 0.01;
}

async function save() {
  if (!form.product_id || !form.store_name.trim() || form.price <= 0 || form.stock <= 0 || form.activity_dates.length !== 2) {
    ElMessage.warning('请填写活动名称、商品、预售价、库存和活动时间');
    return;
  }
  if (selectedProduct.value && form.price > selectedProduct.value.price) {
    ElMessage.warning('预售价不能高于商品销售价');
    return;
  }
  if (form.presell_type === 2 && (!validDeposit() || form.final_dates.length !== 2)) {
    ElMessage.warning('定金预售需填写尾款支付期，且定金加尾款必须等于预售价');
    return;
  }
  const body: MerchantPresellSaveInput = {
    down_price: form.presell_type === 2 ? form.down_price : 0,
    end_time: form.activity_dates[1]!,
    final_end_time: form.presell_type === 2 ? form.final_dates[1]! : '',
    final_price: form.presell_type === 2 ? form.final_price : 0,
    final_start_time: form.presell_type === 2 ? form.final_dates[0]! : '',
    is_show: form.is_show,
    presell_type: form.presell_type,
    price: form.price,
    product_id: form.product_id,
    stock: form.stock,
    store_info: form.store_info.trim(),
    store_name: form.store_name.trim(),
    start_time: form.activity_dates[0]!,
  };
  saving.value = true;
  try {
    if (editingID.value) await updateMerchantPresellActiveApi(editingID.value, body);
    else await createMerchantPresellActiveApi(body);
    ElMessage.success(editingID.value ? '预售活动已更新' : '预售活动已创建');
    dialogOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function toggle(row: MerchantPresellActive) {
  const isShow = row.is_show === 1 ? 0 : 1;
  await setMerchantPresellShowApi(row.product_presell_id, isShow);
  row.is_show = isShow;
  ElMessage.success(isShow ? '活动已上架' : '活动已下架');
}

async function remove(row: MerchantPresellActive) {
  try {
    await ElMessageBox.confirm(`删除预售活动「${row.store_name}」？删除后不可恢复。`, '删除确认', { type: 'warning' });
  } catch {
    return;
  }
  await deleteMerchantPresellActiveApi(row.product_presell_id);
  ElMessage.success('已删除');
  await load();
}

watch(() => form.presell_type, (value) => {
  if (value === 1) {
    form.down_price = 0;
    form.final_price = 0;
    form.final_dates = [];
  }
});

onMounted(async () => {
  await Promise.all([load(), loadProducts()]);
});
</script>

<template>
  <Page title="预售活动" description="全款预售可按预售规则计价；定金预售在定金支付后生成尾款单，尾款必须在配置的支付期内完成。">
    <template #extra><el-button type="primary" @click="openCreate">新建预售活动</el-button></template>
    <el-card shadow="never"><el-form inline @submit.prevent="query.page = 1"><el-form-item label="搜索"><el-input v-model="query.keyword" clearable placeholder="活动名称" @keyup.enter="query.page = 1" /></el-form-item><el-form-item><el-button type="primary" @click="query.page = 1">查询</el-button><el-button @click="query.keyword = ''; query.page = 1">重置</el-button></el-form-item></el-form></el-card>
    <el-card class="mt-4" shadow="never"><el-table v-loading="loading" :data="filteredRows" row-key="product_presell_id"><el-table-column label="预售名称" min-width="180" prop="store_name" show-overflow-tooltip /><el-table-column label="类型" width="90"><template #default="{ row }"><el-tag :type="row.presell_type === 2 ? 'warning' : 'success'">{{ row.presell_type === 2 ? '定金预售' : '全款预售' }}</el-tag></template></el-table-column><el-table-column label="预售价" width="100"><template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template></el-table-column><el-table-column label="定金/尾款" min-width="130"><template #default="{ row }"><span v-if="row.presell_type === 2">¥{{ Number(row.down_price).toFixed(2) }} / ¥{{ Number(row.final_price).toFixed(2) }}</span><span v-else>—</span></template></el-table-column><el-table-column label="库存/销量" width="110"><template #default="{ row }">{{ row.stock }} / {{ row.seles }}</template></el-table-column><el-table-column label="活动时间" min-width="210"><template #default="{ row }">{{ dateText(row.start_time) }} 至 {{ dateText(row.end_time) }}</template></el-table-column><el-table-column label="上架" width="80"><template #default="{ row }"><el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '上架' : '下架' }}</el-tag></template></el-table-column><el-table-column fixed="right" label="操作" width="190"><template #default="{ row }"><el-button link type="primary" @click="openEdit(row)">编辑</el-button><el-button link type="warning" @click="toggle(row)">{{ row.is_show === 1 ? '下架' : '上架' }}</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑预售活动' : '新建预售活动'" width="720px" destroy-on-close><el-form class="grid grid-cols-2 gap-x-4" label-width="96px"><el-form-item class="col-span-2" label="预售名称" required><el-input v-model="form.store_name" maxlength="120" show-word-limit /></el-form-item><el-form-item label="参与商品" required><el-select v-model="form.product_id" :disabled="!!editingID" filterable class="w-full" placeholder="选择本店已审核商品" @change="selectedProductChanged"><el-option v-for="item in products" :key="item.product_id" :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`" :value="item.product_id" /></el-select></el-form-item><el-form-item label="预售类型"><el-radio-group v-model="form.presell_type"><el-radio :value="1">全款</el-radio><el-radio :value="2">定金</el-radio></el-radio-group></el-form-item><el-form-item label="预售价" required><el-input-number v-model="form.price" :min="0.01" :precision="2" class="w-full" /></el-form-item><el-form-item label="预售库存" required><el-input-number v-model="form.stock" :min="1" :precision="0" class="w-full" /></el-form-item><el-form-item class="col-span-2" label="活动时间" required><el-date-picker v-model="form.activity_dates" class="w-full" end-placeholder="活动结束" start-placeholder="活动开始" type="datetimerange" value-format="YYYY-MM-DD HH:mm:ss" /></el-form-item><template v-if="form.presell_type === 2"><el-form-item label="定金" required><el-input-number v-model="form.down_price" :min="0.01" :precision="2" class="w-full" /></el-form-item><el-form-item label="尾款" required><el-input-number v-model="form.final_price" :min="0.01" :precision="2" class="w-full" /></el-form-item><el-form-item class="col-span-2" label="尾款支付期" required><el-date-picker v-model="form.final_dates" class="w-full" end-placeholder="尾款截止" start-placeholder="尾款开始" type="datetimerange" value-format="YYYY-MM-DD HH:mm:ss" /></el-form-item></template><el-form-item class="col-span-2" label="活动说明"><el-input v-model="form.store_info" :rows="3" maxlength="500" show-word-limit type="textarea" /></el-form-item><el-form-item label="初始上架"><el-switch v-model="form.is_show" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
