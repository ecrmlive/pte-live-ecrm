<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformCombinationsApi,
  getPlatformCombinationApi,
  deletePlatformCombinationApi,
  updatePlatformCombinationApi,
  type PlatformCombination,
  type PlatformCombinationInput,
} from '#/api/core/platform-combination';

const loading = ref(false);
const rows = ref<PlatformCombination[]>([]);
const total = ref(0);
const canManage = ref(false);
const query = reactive({ limit: 20, mer_id: undefined as number | undefined, page: 1 });
const editOpen = ref(false);
const saving = ref(false);
const editingID = ref<number>();
const form = reactive<Required<PlatformCombinationInput>>({ price: 0, buying_count_num: 2, time: 24, start_time: '', end_time: '', is_show: 1, status: 1 });

async function load() {
  loading.value = true;
  try {
    const data = await listPlatformCombinationsApi(query);
    rows.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
}

async function setStatus(row: PlatformCombination, status: number) {
  const action = status === 1 ? '启用' : '停用';
  try {
    await ElMessageBox.confirm(`确认${action}商品 #${row.product_id} 的拼团活动？`, `${action}确认`, {
      cancelButtonText: '取消',
      confirmButtonText: `确认${action}`,
      type: 'warning',
    });
    await updatePlatformCombinationApi(row.product_group_id, { status });
    ElMessage.success(`活动已${action}`);
    await load();
  } catch {
    // 用户取消或接口错误由统一请求层提示。
  }
}

function dateTime(value: string) {
  return value ? value.replace('T', ' ').slice(0, 19) : '';
}

async function edit(row: PlatformCombination) {
  const detail = await getPlatformCombinationApi(row.product_group_id);
  editingID.value = row.product_group_id;
  Object.assign(form, { price: Number(detail.price), buying_count_num: detail.buying_count_num, time: detail.time || 24, start_time: dateTime(detail.start_time), end_time: dateTime(detail.end_time), is_show: detail.is_show, status: detail.status });
  editOpen.value = true;
}

async function save() {
  if (!editingID.value || form.price <= 0 || form.buying_count_num < 2 || form.time < 1 || !form.start_time || !form.end_time || new Date(form.end_time).valueOf() < new Date(form.start_time).valueOf()) {
    ElMessage.warning('请填写正数拼团价、至少 2 人、有效时长和正确的活动时间');
    return;
  }
  saving.value = true;
  try {
    await updatePlatformCombinationApi(editingID.value, { ...form });
    ElMessage.success('拼团活动已更新');
    editOpen.value = false;
    await load();
  } finally { saving.value = false; }
}

async function remove(row: PlatformCombination) {
  try {
    await ElMessageBox.confirm(`删除商品 #${row.product_id} 的拼团配置只会软删除活动，进行中或已完成团单不会被改写。是否继续？`, '删除拼团活动', { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' });
    await deletePlatformCombinationApi(row.product_group_id);
    ElMessage.success('拼团活动已软删除');
    await load();
  } catch {
    // 取消不提示；统一请求层处理权限和业务错误。
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi(), load()]);
  canManage.value = profile.roles.some((role) => role === 'platform' || role === 'operations') && permissions.includes('marketing.combination.manage');
});
</script>

<template>
  <Page title="拼团监管" description="监管各商户拼团活动；运营可维护活动价格、人数与时间配置或软删除，商品归属与已产生团单不可在此变更。">
    <el-card shadow="never">
      <el-form inline @submit.prevent="query.page = 1; load()">
        <el-form-item label="商户 ID"><el-input-number v-model="query.mer_id" :min="1" /></el-form-item>
        <el-button type="primary" @click="query.page = 1; load()">查询</el-button>
      </el-form>
      <el-table v-loading="loading" :data="rows">
        <el-table-column label="商品" min-width="180"><template #default="{ row }">{{ row.store_name || `商品 #${row.product_id}` }}</template></el-table-column>
        <el-table-column label="商户" min-width="130"><template #default="{ row }">{{ row.mer_name || `商户 #${row.mer_id}` }}</template></el-table-column>
        <el-table-column label="拼团价" width="110"><template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template></el-table-column>
        <el-table-column label="成团人数" prop="buying_count_num" width="100" />
        <el-table-column label="活动时间" min-width="220"><template #default="{ row }">{{ row.start_time }} 至 {{ row.end_time }}</template></el-table-column>
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template></el-table-column>
        <el-table-column v-if="canManage" fixed="right" label="操作" width="172"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link :type="row.status === 1 ? 'danger' : 'success'" @click="setStatus(row, row.status === 1 ? 0 : 1)">{{ row.status === 1 ? '停用' : '启用' }}</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" /></div>
    </el-card>
    <el-dialog v-model="editOpen" title="编辑拼团活动" width="620px" destroy-on-close><el-alert class="mb-4" type="warning" :closable="false" title="仅维护拼团配置；商品、商户与已产生团单的成员、价格快照不可在此修改。" /><el-form label-width="118px"><el-form-item label="拼团价" required><el-input-number v-model="form.price" :min="0.01" :precision="2" :step="1" /></el-form-item><el-form-item label="成团人数" required><el-input-number v-model="form.buying_count_num" :min="2" :max="9999" /></el-form-item><el-form-item label="成团时限（小时）" required><el-input-number v-model="form.time" :min="1" :max="720" /></el-form-item><el-form-item label="活动时间" required><el-date-picker v-model="form.start_time" value-format="YYYY-MM-DD HH:mm:ss" type="datetime" /><span class="mx-2">至</span><el-date-picker v-model="form.end_time" value-format="YYYY-MM-DD HH:mm:ss" type="datetime" /></el-form-item><el-form-item label="前台展示"><el-radio-group v-model="form.is_show"><el-radio :value="1">上架</el-radio><el-radio :value="0">下架</el-radio></el-radio-group></el-form-item><el-form-item label="活动状态"><el-radio-group v-model="form.status"><el-radio :value="1">启用</el-radio><el-radio :value="0">停用</el-radio></el-radio-group></el-form-item></el-form><template #footer><el-button @click="editOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
