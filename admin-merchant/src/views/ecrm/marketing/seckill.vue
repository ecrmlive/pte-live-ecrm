<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createMerchantSeckillActiveApi,
  deleteMerchantSeckillActiveApi,
  listMerchantSeckillActivesApi,
  listMerchantSeckillTimesApi,
  setMerchantSeckillStatusApi,
  updateMerchantSeckillActiveApi,
  type MerchantSeckillActive,
  type MerchantSeckillSaveInput,
  type MerchantSeckillTime,
} from '#/api/core/merchant-seckill';
import {
  listMerchantProductsApi,
  type MerchantProduct,
} from '#/api/core/merchant-catalog';

const loading = ref(false);
const saving = ref(false);
const dialogOpen = ref(false);
const editingID = ref<number>();
const rows = ref<MerchantSeckillActive[]>([]);
const products = ref<MerchantProduct[]>([]);
const slots = ref<MerchantSeckillTime[]>([]);
const total = ref(0);
const query = reactive({ keyword: '', limit: 20, page: 1 });
const form = reactive({
  dates: [] as string[],
  name: '',
  once_pay_count: 1,
  product_id: 0,
  seckill_price: 0,
  slot_ids: [] as number[],
  status: 1,
});

const filteredRows = computed(() => {
  const keyword = query.keyword.trim().toLowerCase();
  if (!keyword) return rows.value;
  return rows.value.filter((row) => `${row.name} ${row.store_name || ''}`.toLowerCase().includes(keyword));
});
const selectedProduct = computed(() => products.value.find((item) => item.product_id === form.product_id));

function statusType(status: number) {
  return status === 1 ? 'success' : 'info';
}

function slotsText(ids: string) {
  const selected = new Set(ids.split(',').map((id) => Number(id)));
  return slots.value.filter((slot) => selected.has(slot.seckill_time_id)).map((slot) => slot.title).join('、') || '全部场次';
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, { dates: [], name: '', once_pay_count: 1, product_id: 0, seckill_price: 0, slot_ids: [], status: 1 });
}

function openCreate() {
  resetForm();
  dialogOpen.value = true;
}

function openEdit(row: MerchantSeckillActive) {
  editingID.value = row.seckill_active_id;
  Object.assign(form, {
    dates: [row.start_day, row.end_day],
    name: row.name,
    once_pay_count: row.once_pay_count,
    product_id: row.product_id,
    seckill_price: row.seckill_price,
    slot_ids: row.seckill_time_ids.split(',').map((id) => Number(id)).filter((id) => id > 0),
    status: row.status,
  });
  dialogOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const data = await listMerchantSeckillActivesApi({ limit: query.limit, page: query.page });
    rows.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
}

async function loadOptions() {
  const [times, productPage] = await Promise.all([
    listMerchantSeckillTimesApi(),
    listMerchantProductsApi({ limit: 100, page: 1, status: 1 }),
  ]);
  slots.value = times.list || [];
  products.value = productPage.list || [];
}

async function save() {
  if (!form.name.trim() || !form.product_id || form.seckill_price <= 0 || form.dates.length !== 2) {
    ElMessage.warning('请填写活动名称、商品、秒杀价和活动日期');
    return;
  }
  if (selectedProduct.value && form.seckill_price >= selectedProduct.value.price) {
    ElMessage.warning('秒杀价应低于商品销售价');
    return;
  }
  const body: MerchantSeckillSaveInput = {
    end_day: form.dates[1]!,
    name: form.name.trim(),
    once_pay_count: form.once_pay_count,
    product_id: form.product_id,
    seckill_price: form.seckill_price,
    seckill_time_ids: form.slot_ids.join(','),
    start_day: form.dates[0]!,
    status: form.status,
  };
  saving.value = true;
  try {
    if (editingID.value) await updateMerchantSeckillActiveApi(editingID.value, body);
    else await createMerchantSeckillActiveApi(body);
    ElMessage.success(editingID.value ? '秒杀活动已更新' : '秒杀活动已创建');
    dialogOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function toggle(row: MerchantSeckillActive) {
  const status = row.status === 1 ? 0 : 1;
  await setMerchantSeckillStatusApi(row.seckill_active_id, status);
  row.status = status;
  ElMessage.success(status ? '活动已启用' : '活动已停用');
}

async function remove(row: MerchantSeckillActive) {
  try {
    await ElMessageBox.confirm(`删除秒杀活动「${row.name}」？删除后不可恢复。`, '删除确认', { type: 'warning' });
  } catch {
    return;
  }
  await deleteMerchantSeckillActiveApi(row.seckill_active_id);
  ElMessage.success('已删除');
  await load();
}

function search() { query.page = 1; }

onMounted(async () => {
  await Promise.all([load(), loadOptions()]);
});
</script>

<template>
  <Page title="秒杀活动" description="秒杀商品结算不参与店铺券或平台券；启停仅影响后续下单报价，不修改已创建订单。">
    <template #extra><el-button type="primary" @click="openCreate">新建秒杀活动</el-button></template>
    <el-card shadow="never">
      <el-form inline @submit.prevent="search"><el-form-item label="搜索"><el-input v-model="query.keyword" clearable placeholder="活动名称 / 商品名称" @keyup.enter="search" /></el-form-item><el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="query.keyword = ''; search()">重置</el-button></el-form-item></el-form>
    </el-card>
    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="filteredRows" row-key="seckill_active_id"><el-table-column label="活动名称" min-width="150" prop="name" show-overflow-tooltip /><el-table-column label="商品" min-width="170" prop="store_name" show-overflow-tooltip /><el-table-column label="原价" width="100"><template #default="{ row }">¥{{ Number(row.price || 0).toFixed(2) }}</template></el-table-column><el-table-column label="秒杀价" width="105"><template #default="{ row }">¥{{ Number(row.seckill_price).toFixed(2) }}</template></el-table-column><el-table-column label="限购" prop="once_pay_count" width="88"><template #default="{ row }">{{ row.once_pay_count }} 件</template></el-table-column><el-table-column label="活动日期" min-width="190"><template #default="{ row }">{{ row.start_day }} 至 {{ row.end_day }}</template></el-table-column><el-table-column label="场次" min-width="150"><template #default="{ row }">{{ slotsText(row.seckill_time_ids) }}</template></el-table-column><el-table-column label="状态" width="88"><template #default="{ row }"><el-tag :type="statusType(row.status)">{{ row.status === 1 ? '已启用' : '已停用' }}</el-tag></template></el-table-column><el-table-column fixed="right" label="操作" width="190"><template #default="{ row }"><el-button link type="primary" @click="openEdit(row)">编辑</el-button><el-button link type="warning" @click="toggle(row)">{{ row.status === 1 ? '停用' : '启用' }}</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div>
    </el-card>
    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑秒杀活动' : '新建秒杀活动'" width="640px" destroy-on-close><el-form class="grid grid-cols-2 gap-x-4" label-width="88px"><el-form-item class="col-span-2" label="活动名称" required><el-input v-model="form.name" maxlength="60" show-word-limit /></el-form-item><el-form-item label="参与商品" required><el-select v-model="form.product_id" :disabled="!!editingID" filterable class="w-full" placeholder="选择本店已审核商品"><el-option v-for="item in products" :key="item.product_id" :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`" :value="item.product_id" /></el-select></el-form-item><el-form-item label="秒杀价" required><el-input-number v-model="form.seckill_price" :min="0.01" :precision="2" class="w-full" /></el-form-item><el-form-item class="col-span-2" label="活动日期" required><el-date-picker v-model="form.dates" class="w-full" end-placeholder="结束日期" start-placeholder="开始日期" type="daterange" value-format="YYYY-MM-DD" /></el-form-item><el-form-item class="col-span-2" label="参与场次"><el-select v-model="form.slot_ids" multiple class="w-full" placeholder="不选则使用默认场次"><el-option v-for="slot in slots" :key="slot.seckill_time_id" :label="slot.title" :value="slot.seckill_time_id" /></el-select></el-form-item><el-form-item label="单次限购"><el-input-number v-model="form.once_pay_count" :min="1" :precision="0" class="w-full" /></el-form-item><el-form-item label="初始状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
