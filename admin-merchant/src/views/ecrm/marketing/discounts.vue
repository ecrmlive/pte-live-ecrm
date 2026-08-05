<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  createMerchantDiscountApi,
  deleteMerchantDiscountApi,
  listMerchantDiscountsApi,
  setMerchantDiscountStatusApi,
  updateMerchantDiscountApi,
  type MerchantDiscount,
  type MerchantDiscountStatus,
} from '#/api/core/merchant-discount';
import { EcrmFormDialog, EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const saving = ref(false);
const rows = ref<MerchantDiscount[]>([]);
const total = ref(0);
const canManage = ref(false);
const dialog = ref(false);
const editingID = ref<number>();
const query = reactive({
  keyword: '',
  limit: 20,
  page: 1,
  status: undefined as MerchantDiscountStatus | undefined,
});
const form = reactive({
  ends_at: '',
  free_shipping: false,
  name: '',
  package_price: 0,
  product_ids_text: '',
  remark: '',
  starts_at: '',
  status: 'draft' as MerchantDiscountStatus,
});

const statusLabels: Record<MerchantDiscountStatus, string> = {
  active: '进行中',
  closed: '已关闭',
  draft: '草稿',
  pending: '待审核',
  rejected: '已拒绝',
};

async function load() {
  loading.value = true;
  try {
    const page = await listMerchantDiscountsApi({
      keyword: query.keyword.trim() || undefined,
      limit: query.limit,
      page: query.page,
      status: query.status,
    });
    rows.value = page.list || [];
    total.value = page.total || 0;
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.keyword = '';
  query.status = undefined;
  query.page = 1;
  void load();
}

function openCreate() {
  editingID.value = undefined;
  Object.assign(form, {
    ends_at: '',
    free_shipping: false,
    name: '',
    package_price: 0,
    product_ids_text: '',
    remark: '',
    starts_at: '',
    status: 'draft',
  });
  dialog.value = true;
}

function openEdit(row: MerchantDiscount) {
  editingID.value = row.activity_id;
  Object.assign(form, {
    ends_at: row.ends_at,
    free_shipping: row.free_shipping,
    name: row.name,
    package_price: row.package_price,
    product_ids_text: (row.product_ids || []).join(','),
    remark: row.remark,
    starts_at: row.starts_at,
    status: row.status,
  });
  dialog.value = true;
}

function parseProductIDs() {
  return form.product_ids_text
    .split(/[,，\s]+/)
    .map((part) => Number(part.trim()))
    .filter((id) => Number.isFinite(id) && id > 0);
}

async function save() {
  const productIDs = parseProductIDs();
  if (!form.name.trim() || form.package_price <= 0 || productIDs.length === 0) {
    ElMessage.warning('请填写名称、正数套餐价，并至少填写一个商品 ID');
    return;
  }
  saving.value = true;
  try {
    const body = {
      ends_at: form.ends_at || undefined,
      free_shipping: form.free_shipping,
      name: form.name.trim(),
      package_price: form.package_price,
      product_ids: productIDs,
      remark: form.remark.trim(),
      starts_at: form.starts_at || undefined,
      status: form.status,
    };
    if (editingID.value) {
      await updateMerchantDiscountApi(editingID.value, body);
      ElMessage.success('优惠套餐已更新');
    } else {
      await createMerchantDiscountApi(body);
      ElMessage.success('优惠套餐已创建');
    }
    dialog.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function toggleStatus(row: MerchantDiscount) {
  const next: MerchantDiscountStatus = row.status === 'active' ? 'closed' : 'active';
  const action = next === 'active' ? '上架' : '关闭';
  try {
    await ElMessageBox.confirm(`确认${action}「${row.name}」？`, `${action}确认`, { type: 'warning' });
    await setMerchantDiscountStatusApi(row.activity_id, next);
    ElMessage.success(`已${action}`);
    await load();
  } catch {
    /* cancel or request layer toast */
  }
}

async function remove(row: MerchantDiscount) {
  try {
    await ElMessageBox.confirm(`删除「${row.name}」后不可恢复，确认继续？`, '删除优惠套餐', {
      type: 'warning',
      confirmButtonText: '删除',
    });
    await deleteMerchantDiscountApi(row.activity_id);
    ElMessage.success('已删除');
    await load();
  } catch {
    /* cancel or request layer toast */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('marketing.discounts.manage');
  await load();
});
</script>

<template>
  <EcrmListPage
    title="优惠套餐"
    description="店铺优惠套餐写入 qixi_crm_m_marketing_activity（activity_type=discount），并同步业务投影供 C 端/平台监管。"
  >
    <template #filters>
      <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" clearable maxlength="64" placeholder="套餐名称" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" clearable class="w-32" placeholder="全部">
            <el-option
              v-for="(label, status) in statusLabels"
              :key="status"
              :label="label"
              :value="status"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">查询</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </template>

    <template #actions>
      <el-button v-if="canManage" type="primary" @click="openCreate">新增套餐</el-button>
    </template>

    <el-table v-loading="loading" :data="rows" row-key="activity_id">
      <el-table-column prop="activity_id" label="ID" width="90" />
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column label="套餐价" width="110">
        <template #default="{ row }">¥{{ Number(row.package_price).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="商品数" width="90">
        <template #default="{ row }">{{ (row.product_ids || []).length }}</template>
      </el-table-column>
      <el-table-column label="包邮" width="80">
        <template #default="{ row }">{{ row.free_shipping ? '是' : '否' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag>{{ statusLabels[row.status] || row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="有效期" min-width="220">
        <template #default="{ row }">
          {{ row.starts_at || '—' }} ~ {{ row.ends_at || '—' }}
        </template>
      </el-table-column>
      <el-table-column v-if="canManage" fixed="right" label="操作" width="200">
        <template #default="{ row }">
          <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
          <el-button link type="warning" @click="toggleStatus(row)">
            {{ row.status === 'active' ? '关闭' : '上架' }}
          </el-button>
          <el-button link type="danger" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <template #pager>
      <el-pagination
        :current-page="query.page"
        :page-size="query.limit"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="(page: number) => { query.page = page; load(); }"
      />
    </template>
  </EcrmListPage>

  <EcrmFormDialog
    v-model="dialog"
    :title="editingID ? '编辑优惠套餐' : '新增优惠套餐'"
    width="560px"
  >
    <el-form label-width="96px">
      <el-form-item label="名称" required>
        <el-input v-model="form.name" maxlength="128" show-word-limit />
      </el-form-item>
      <el-form-item label="套餐价" required>
        <el-input-number v-model="form.package_price" :min="0.01" :precision="2" :step="1" />
      </el-form-item>
      <el-form-item label="商品 ID" required>
        <el-input
          v-model="form.product_ids_text"
          placeholder="多个 ID 用逗号分隔，如 1001,1006"
        />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="form.status" class="w-40">
          <el-option
            v-for="(label, status) in statusLabels"
            :key="status"
            :label="label"
            :value="status"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="开始时间">
        <el-input v-model="form.starts_at" placeholder="YYYY-MM-DD HH:mm:ss" />
      </el-form-item>
      <el-form-item label="结束时间">
        <el-input v-model="form.ends_at" placeholder="YYYY-MM-DD HH:mm:ss" />
      </el-form-item>
      <el-form-item label="包邮">
        <el-switch v-model="form.free_shipping" />
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="form.remark" type="textarea" :rows="2" maxlength="255" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialog = false">取消</el-button>
      <el-button :loading="saving" type="primary" @click="save">保存</el-button>
    </template>
  </EcrmFormDialog>
</template>
