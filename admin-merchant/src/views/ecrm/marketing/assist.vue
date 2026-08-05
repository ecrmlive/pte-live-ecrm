<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createMerchantAssistActiveApi,
  deleteMerchantAssistActiveApi,
  listMerchantAssistActivesApi,
  setMerchantAssistShowApi,
  updateMerchantAssistActiveApi,
  type MerchantAssistActive,
  type MerchantAssistSaveInput,
} from '#/api/core/merchant-assist';
import {
  listMerchantProductsApi,
  type MerchantProduct,
} from '#/api/core/merchant-catalog';

const loading = ref(false);
const saving = ref(false);
const dialogOpen = ref(false);
const editingID = ref<number>();
const rows = ref<MerchantAssistActive[]>([]);
const products = ref<MerchantProduct[]>([]);
const total = ref(0);
const query = reactive({ keyword: '', limit: 20, page: 1 });
const form = reactive({
  assist_count: 2,
  assist_price: 0,
  assist_user_count: 1,
  dates: [] as string[],
  is_show: 1,
  product_id: 0,
  stock: 100,
  store_info: '',
  store_name: '',
});

const filteredRows = computed(() => {
  const keyword = query.keyword.trim().toLowerCase();
  if (!keyword) return rows.value;
  return rows.value.filter((row) => `${row.store_name} ${row.product_id}`.toLowerCase().includes(keyword));
});

function formatTime(value: string) {
  return value ? value.replace('T', ' ').slice(0, 19) : '—';
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    assist_count: 2,
    assist_price: 0,
    assist_user_count: 1,
    dates: [],
    is_show: 1,
    product_id: 0,
    stock: 100,
    store_info: '',
    store_name: '',
  });
}

function openCreate() {
  resetForm();
  dialogOpen.value = true;
}

function openEdit(row: MerchantAssistActive) {
  editingID.value = row.product_assist_id;
  Object.assign(form, {
    assist_count: row.assist_count,
    assist_price: row.assist_price,
    assist_user_count: row.assist_user_count,
    dates: [formatTime(row.start_time), formatTime(row.end_time)],
    is_show: row.is_show,
    product_id: row.product_id,
    stock: row.stock,
    store_info: row.store_info || '',
    store_name: row.store_name,
  });
  dialogOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const data = await listMerchantAssistActivesApi({ limit: query.limit, page: query.page });
    rows.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
}

async function loadOptions() {
  const productPage = await listMerchantProductsApi({ limit: 100, page: 1, status: 1 });
  products.value = productPage.list || [];
}

async function save() {
  if (!form.product_id || form.assist_price <= 0 || form.dates.length !== 2) {
    ElMessage.warning('请选择商品、填写助力价并设置活动时间');
    return;
  }
  const body: MerchantAssistSaveInput = {
    assist_count: form.assist_count,
    assist_price: form.assist_price,
    assist_user_count: form.assist_user_count,
    end_time: form.dates[1]!,
    is_show: form.is_show,
    product_id: form.product_id,
    start_time: form.dates[0]!,
    stock: form.stock,
    store_info: form.store_info.trim(),
    store_name: form.store_name.trim(),
  };
  saving.value = true;
  try {
    if (editingID.value) await updateMerchantAssistActiveApi(editingID.value, body);
    else await createMerchantAssistActiveApi(body);
    ElMessage.success(editingID.value ? '助力活动已更新' : '助力活动已创建');
    dialogOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function toggleShow(row: MerchantAssistActive) {
  const isShow = row.is_show === 1 ? 0 : 1;
  await setMerchantAssistShowApi(row.product_assist_id, isShow);
  row.is_show = isShow;
  ElMessage.success(isShow ? '活动已上架' : '活动已下架');
}

async function remove(row: MerchantAssistActive) {
  try {
    await ElMessageBox.confirm(`删除助力活动「${row.store_name}」？删除后不可恢复。`, '删除确认', { type: 'warning' });
  } catch {
    return;
  }
  await deleteMerchantAssistActiveApi(row.product_assist_id);
  ElMessage.success('已删除');
  await load();
}

function search() {
  query.page = 1;
}

onMounted(async () => {
  await Promise.all([load(), loadOptions()]);
});
</script>

<template>
  <Page title="好友助力" description="配置本店好友助力活动；上架状态控制前台展示，已发起的助力单不受后续改价影响。">
    <template #extra><el-button type="primary" @click="openCreate">新建助力活动</el-button></template>
    <el-card shadow="never">
      <el-form inline @submit.prevent="search">
        <el-form-item label="搜索"><el-input v-model="query.keyword" clearable placeholder="活动名称 / 商品 ID" @keyup.enter="search" /></el-form-item>
        <el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="query.keyword = ''; search()">重置</el-button></el-form-item>
      </el-form>
    </el-card>
    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="filteredRows" row-key="product_assist_id">
        <el-table-column label="活动 / 商品" min-width="170" prop="store_name" show-overflow-tooltip />
        <el-table-column label="商品 ID" prop="product_id" width="92" />
        <el-table-column label="助力价" width="105"><template #default="{ row }">¥{{ Number(row.assist_price).toFixed(2) }}</template></el-table-column>
        <el-table-column label="助力规则" min-width="150"><template #default="{ row }">{{ row.assist_count }} 人 / 每人 {{ row.assist_user_count }} 次</template></el-table-column>
        <el-table-column label="库存" prop="stock" width="80" />
        <el-table-column label="活动时间" min-width="220"><template #default="{ row }">{{ formatTime(row.start_time) }} 至 {{ formatTime(row.end_time) }}</template></el-table-column>
        <el-table-column label="展示" width="88"><template #default="{ row }"><el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '已上架' : '已下架' }}</el-tag></template></el-table-column>
        <el-table-column fixed="right" label="操作" width="190">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button link type="warning" @click="toggleShow(row)">{{ row.is_show === 1 ? '下架' : '上架' }}</el-button>
            <el-button link type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div>
    </el-card>
    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑助力活动' : '新建助力活动'" width="640px" destroy-on-close>
      <el-form class="grid grid-cols-2 gap-x-4" label-width="96px">
        <el-form-item class="col-span-2" label="活动名称"><el-input v-model="form.store_name" maxlength="60" placeholder="留空则使用商品名" show-word-limit /></el-form-item>
        <el-form-item class="col-span-2" label="活动简介"><el-input v-model="form.store_info" :rows="2" maxlength="200" show-word-limit type="textarea" /></el-form-item>
        <el-form-item label="参与商品" required><el-select v-model="form.product_id" :disabled="!!editingID" filterable class="w-full" placeholder="选择本店已审核商品"><el-option v-for="item in products" :key="item.product_id" :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`" :value="item.product_id" /></el-select></el-form-item>
        <el-form-item label="助力价" required><el-input-number v-model="form.assist_price" :min="0.01" :precision="2" class="w-full" /></el-form-item>
        <el-form-item label="所需人数" required><el-input-number v-model="form.assist_count" :min="1" :precision="0" class="w-full" /></el-form-item>
        <el-form-item label="单人次数" required><el-input-number v-model="form.assist_user_count" :min="1" :precision="0" class="w-full" /></el-form-item>
        <el-form-item label="活动库存" required><el-input-number v-model="form.stock" :min="1" :precision="0" class="w-full" /></el-form-item>
        <el-form-item label="前台展示"><el-switch v-model="form.is_show" :active-value="1" :inactive-value="0" /></el-form-item>
        <el-form-item class="col-span-2" label="活动时间" required><el-date-picker v-model="form.dates" class="w-full" end-placeholder="结束时间" start-placeholder="开始时间" type="datetimerange" value-format="YYYY-MM-DD HH:mm:ss" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template>
    </el-dialog>
  </Page>
</template>
