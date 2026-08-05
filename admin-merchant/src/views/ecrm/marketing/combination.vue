<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createMerchantCombinationGroupApi,
  deleteMerchantCombinationGroupApi,
  listMerchantCombinationGroupsApi,
  setMerchantCombinationShowApi,
  updateMerchantCombinationGroupApi,
  type MerchantCombinationGroup,
  type MerchantCombinationSaveInput,
} from '#/api/core/merchant-combination';
import {
  listMerchantProductsApi,
  type MerchantProduct,
} from '#/api/core/merchant-catalog';

const loading = ref(false);
const saving = ref(false);
const dialogOpen = ref(false);
const editingID = ref<number>();
const rows = ref<MerchantCombinationGroup[]>([]);
const products = ref<MerchantProduct[]>([]);
const total = ref(0);
const query = reactive({ keyword: '', limit: 20, page: 1 });
const form = reactive({
  buying_count_num: 2,
  dates: [] as string[],
  is_show: 1,
  price: 0,
  product_id: 0,
  time: 24,
});

const filteredRows = computed(() => {
  const keyword = query.keyword.trim().toLowerCase();
  if (!keyword) return rows.value;
  return rows.value.filter((row) => `${row.store_name || ''} ${row.product_id}`.toLowerCase().includes(keyword));
});

function formatTime(value: string) {
  return value ? value.replace('T', ' ').slice(0, 19) : '—';
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, { buying_count_num: 2, dates: [], is_show: 1, price: 0, product_id: 0, time: 24 });
}

function openCreate() {
  resetForm();
  dialogOpen.value = true;
}

function openEdit(row: MerchantCombinationGroup) {
  editingID.value = row.product_group_id;
  Object.assign(form, {
    buying_count_num: row.buying_count_num,
    dates: [formatTime(row.start_time), formatTime(row.end_time)],
    is_show: row.is_show,
    price: row.price,
    product_id: row.product_id,
    time: row.time || 24,
  });
  dialogOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const data = await listMerchantCombinationGroupsApi({ limit: query.limit, page: query.page });
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
  if (!form.product_id || form.price <= 0 || form.buying_count_num < 2 || form.dates.length !== 2) {
    ElMessage.warning('请选择商品、填写拼团价、成团人数与活动时间');
    return;
  }
  const body: MerchantCombinationSaveInput = {
    buying_count_num: form.buying_count_num,
    end_time: form.dates[1]!,
    is_show: form.is_show,
    price: form.price,
    product_id: form.product_id,
    start_time: form.dates[0]!,
    time: form.time,
  };
  saving.value = true;
  try {
    if (editingID.value) await updateMerchantCombinationGroupApi(editingID.value, body);
    else await createMerchantCombinationGroupApi(body);
    ElMessage.success(editingID.value ? '拼团活动已更新' : '拼团活动已创建');
    dialogOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function toggleShow(row: MerchantCombinationGroup) {
  const isShow = row.is_show === 1 ? 0 : 1;
  await setMerchantCombinationShowApi(row.product_group_id, isShow);
  row.is_show = isShow;
  ElMessage.success(isShow ? '活动已上架' : '活动已下架');
}

async function remove(row: MerchantCombinationGroup) {
  try {
    await ElMessageBox.confirm(`删除商品 #${row.product_id} 的拼团活动？进行中团单不会被改写。`, '删除确认', { type: 'warning' });
  } catch {
    return;
  }
  await deleteMerchantCombinationGroupApi(row.product_group_id);
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
  <Page title="拼团活动" description="配置本店拼团商品与成团规则；上架状态控制前台展示，已产生团单不受后续改价影响。">
    <template #extra><el-button type="primary" @click="openCreate">新建拼团活动</el-button></template>
    <el-card shadow="never">
      <el-form inline @submit.prevent="search">
        <el-form-item label="搜索"><el-input v-model="query.keyword" clearable placeholder="商品名称 / 商品 ID" @keyup.enter="search" /></el-form-item>
        <el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="query.keyword = ''; search()">重置</el-button></el-form-item>
      </el-form>
    </el-card>
    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="filteredRows" row-key="product_group_id">
        <el-table-column label="商品" min-width="170"><template #default="{ row }">{{ row.store_name || `商品 #${row.product_id}` }}</template></el-table-column>
        <el-table-column label="拼团价" width="105"><template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template></el-table-column>
        <el-table-column label="成团人数" prop="buying_count_num" width="100" />
        <el-table-column label="成团时限" width="100"><template #default="{ row }">{{ row.time }} 小时</template></el-table-column>
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
    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑拼团活动' : '新建拼团活动'" width="640px" destroy-on-close>
      <el-form class="grid grid-cols-2 gap-x-4" label-width="108px">
        <el-form-item class="col-span-2" label="参与商品" required><el-select v-model="form.product_id" :disabled="!!editingID" filterable class="w-full" placeholder="选择本店已审核商品"><el-option v-for="item in products" :key="item.product_id" :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`" :value="item.product_id" /></el-select></el-form-item>
        <el-form-item label="拼团价" required><el-input-number v-model="form.price" :min="0.01" :precision="2" class="w-full" /></el-form-item>
        <el-form-item label="成团人数" required><el-input-number v-model="form.buying_count_num" :min="2" :precision="0" class="w-full" /></el-form-item>
        <el-form-item label="成团时限" required><el-input-number v-model="form.time" :max="720" :min="1" :precision="0" class="w-full" /><span class="ml-2 text-sm text-muted-foreground">小时</span></el-form-item>
        <el-form-item label="前台展示"><el-switch v-model="form.is_show" :active-value="1" :inactive-value="0" /></el-form-item>
        <el-form-item class="col-span-2" label="活动时间" required><el-date-picker v-model="form.dates" class="w-full" end-placeholder="结束时间" start-placeholder="开始时间" type="datetimerange" value-format="YYYY-MM-DD HH:mm:ss" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template>
    </el-dialog>
  </Page>
</template>
