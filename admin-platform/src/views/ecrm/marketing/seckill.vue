<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformSeckillApi,
  getPlatformSeckillApi,
  deletePlatformSeckillApi,
  updatePlatformSeckillApi,
  type PlatformSeckillInput,
  type PlatformSeckillActive,
} from '#/api/core/platform-seckill';

const loading = ref(false);
const rows = ref<PlatformSeckillActive[]>([]);
const total = ref(0);
const canManage = ref(false);
const query = reactive({ limit: 20, mer_id: undefined as number | undefined, page: 1 });
const editOpen = ref(false);
const saving = ref(false);
const editingID = ref<number>();
const form = reactive<Required<PlatformSeckillInput>>({ name: '', seckill_time_ids: '', start_day: '', end_day: '', seckill_price: 0, once_pay_count: 1, status: 1 });

async function load() {
  loading.value = true;
  try {
    const data = await listPlatformSeckillApi(query);
    rows.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
}

async function setStatus(row: PlatformSeckillActive, status: number) {
  const action = status === 1 ? '启用' : '停用';
  try {
    await ElMessageBox.confirm(`确认${action}秒杀活动“${row.name}”？`, `${action}确认`, {
      cancelButtonText: '取消',
      confirmButtonText: `确认${action}`,
      type: 'warning',
    });
    await updatePlatformSeckillApi(row.seckill_active_id, { status });
    ElMessage.success(`活动已${action}`);
    await load();
  } catch {
    // 用户取消或接口错误由统一请求层提示。
  }
}

async function edit(row: PlatformSeckillActive) {
  const detail = await getPlatformSeckillApi(row.seckill_active_id);
  editingID.value = row.seckill_active_id;
  Object.assign(form, { name: detail.name, seckill_time_ids: detail.seckill_time_ids || '', start_day: detail.start_day, end_day: detail.end_day, seckill_price: Number(detail.seckill_price), once_pay_count: detail.once_pay_count || 1, status: detail.status });
  editOpen.value = true;
}

async function save() {
  if (!editingID.value || !form.name.trim() || !form.start_day || !form.end_day || form.end_day < form.start_day || form.seckill_price <= 0 || form.once_pay_count < 1) {
    ElMessage.warning('请填写活动名称、有效日期、正数秒杀价和限购数量');
    return;
  }
  saving.value = true;
  try {
    await updatePlatformSeckillApi(editingID.value, { ...form, name: form.name.trim(), seckill_time_ids: form.seckill_time_ids.trim() });
    ElMessage.success('秒杀活动已更新');
    editOpen.value = false;
    await load();
  } finally { saving.value = false; }
}

async function remove(row: PlatformSeckillActive) {
  try {
    await ElMessageBox.confirm(`删除“${row.name}”只会软删除活动配置，已产生订单不会被修改。是否继续？`, '删除秒杀活动', { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' });
    await deletePlatformSeckillApi(row.seckill_active_id);
    ElMessage.success('秒杀活动已软删除');
    await load();
  } catch {
    // 取消不提示；统一请求层处理权限和业务错误。
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi(), load()]);
  canManage.value = profile.roles.some((role) => role === 'platform' || role === 'operations') && permissions.includes('marketing.seckill.manage');
});
</script>

<template>
  <Page title="秒杀监管" description="监管各商户秒杀活动；运营可维护活动配置或软删除，已产生订单、商品归属和历史价格快照不会被改写。">
    <el-card shadow="never">
      <el-form inline @submit.prevent="query.page = 1; load()">
        <el-form-item label="商户 ID"><el-input-number v-model="query.mer_id" :min="1" /></el-form-item>
        <el-button type="primary" @click="query.page = 1; load()">查询</el-button>
      </el-form>
      <el-table v-loading="loading" :data="rows">
        <el-table-column label="活动" min-width="160" prop="name" />
        <el-table-column label="商户" min-width="130"><template #default="{ row }">{{ row.mer_name || `商户 #${row.mer_id}` }}</template></el-table-column>
        <el-table-column label="商品" min-width="150"><template #default="{ row }">{{ row.store_name || `商品 #${row.product_id}` }}</template></el-table-column>
        <el-table-column label="秒杀价" width="110"><template #default="{ row }">¥{{ Number(row.seckill_price).toFixed(2) }}</template></el-table-column>
        <el-table-column label="活动日期" min-width="200"><template #default="{ row }">{{ row.start_day }} 至 {{ row.end_day }}</template></el-table-column>
        <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template></el-table-column>
        <el-table-column v-if="canManage" fixed="right" label="操作" width="172"><template #default="{ row }"><el-button link type="primary" @click="edit(row)">编辑</el-button><el-button link :type="row.status === 1 ? 'danger' : 'success'" @click="setStatus(row, row.status === 1 ? 0 : 1)">{{ row.status === 1 ? '停用' : '启用' }}</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" /></div>
    </el-card>
    <el-dialog v-model="editOpen" title="编辑秒杀活动" width="620px" destroy-on-close><el-alert class="mb-4" type="warning" :closable="false" title="仅修改活动配置；商品、商户和已产生订单快照不可在此变更。" /><el-form label-width="118px"><el-form-item label="活动名称" required><el-input v-model="form.name" maxlength="128" show-word-limit /></el-form-item><el-form-item label="秒杀场次 ID"><el-input v-model="form.seckill_time_ids" placeholder="多个场次用英文逗号分隔，例如 1,2" /></el-form-item><el-form-item label="活动日期" required><el-date-picker v-model="form.start_day" value-format="YYYY-MM-DD" type="date" /><span class="mx-2">至</span><el-date-picker v-model="form.end_day" value-format="YYYY-MM-DD" type="date" /></el-form-item><el-form-item label="秒杀价" required><el-input-number v-model="form.seckill_price" :min="0.01" :precision="2" :step="1" /></el-form-item><el-form-item label="单次限购" required><el-input-number v-model="form.once_pay_count" :min="1" :max="9999" /></el-form-item><el-form-item label="活动状态"><el-radio-group v-model="form.status"><el-radio :value="1">启用</el-radio><el-radio :value="0">停用</el-radio></el-radio-group></el-form-item></el-form><template #footer><el-button @click="editOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
  </Page>
</template>
