<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createMerchantCouponApi,
  deleteMerchantCouponApi,
  listMerchantCouponsApi,
  sendMerchantCouponApi,
  setMerchantCouponStatusApi,
  updateMerchantCouponApi,
  type MerchantCoupon,
  type MerchantCouponSaveInput,
} from '#/api/core/merchant-promotion';

const loading = ref(false);
const saving = ref(false);
const rows = ref<MerchantCoupon[]>([]);
const total = ref(0);
const dialogOpen = ref(false);
const sendOpen = ref(false);
const editingID = ref<number>();
const sendingCoupon = ref<MerchantCoupon>();
const query = reactive({ limit: 20, page: 1 });
const form = reactive<MerchantCouponSaveInput>({
  coupon_price: 0,
  coupon_time: 30,
  is_limited: 0,
  sort: 1,
  status: 1,
  title: '',
  total_count: 0,
  use_min_price: 0,
});
const sendForm = reactive({ mark: '', uidsText: '' });

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, { coupon_price: 0, coupon_time: 30, is_limited: 0, sort: 1, status: 1, title: '', total_count: 0, use_min_price: 0 });
}

function openCreate() {
  resetForm();
  dialogOpen.value = true;
}

function openEdit(row: MerchantCoupon) {
  editingID.value = row.coupon_id;
  Object.assign(form, {
    coupon_price: row.coupon_price,
    coupon_time: row.coupon_time,
    is_limited: row.is_limited,
    sort: row.sort,
    status: row.status,
    title: row.title,
    total_count: row.total_count,
    use_min_price: row.use_min_price,
  });
  dialogOpen.value = true;
}

function openSend(row: MerchantCoupon) {
  sendingCoupon.value = row;
  Object.assign(sendForm, { mark: '', uidsText: '' });
  sendOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const result = await listMerchantCouponsApi(query);
    rows.value = result.list;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!form.title.trim()) {
    ElMessage.warning('请填写优惠券名称');
    return;
  }
  if (form.coupon_price <= 0 || form.use_min_price < 0 || form.coupon_time <= 0) {
    ElMessage.warning('请检查优惠金额、使用门槛和有效天数');
    return;
  }
  if (form.is_limited === 1 && form.total_count <= 0) {
    ElMessage.warning('限量发放时必须填写发放总数');
    return;
  }
  saving.value = true;
  try {
    const body = { ...form, title: form.title.trim(), total_count: form.is_limited === 1 ? form.total_count : 0 };
    if (editingID.value) await updateMerchantCouponApi(editingID.value, body);
    else await createMerchantCouponApi(body);
    dialogOpen.value = false;
    ElMessage.success(editingID.value ? '优惠券已更新' : '优惠券已创建');
    await load();
  } finally {
    saving.value = false;
  }
}

async function toggle(row: MerchantCoupon) {
  const next = row.status === 1 ? 0 : 1;
  await setMerchantCouponStatusApi(row.coupon_id, next);
  row.status = next;
  ElMessage.success(next === 1 ? '优惠券已启用' : '优惠券已停用');
}

async function remove(row: MerchantCoupon) {
  try {
    await ElMessageBox.confirm(`删除优惠券“${row.title}”后不可恢复，是否继续？`, '删除确认', { type: 'warning' });
    await deleteMerchantCouponApi(row.coupon_id);
    ElMessage.success('优惠券已删除');
    await load();
  } catch {
    // 用户取消或接口错误时由统一请求层处理。
  }
}

async function send() {
  if (!sendingCoupon.value) return;
  const uids = [...new Set(sendForm.uidsText.split(/[，,\s]+/).map((value) => Number(value)).filter((value) => Number.isSafeInteger(value) && value > 0))];
  if (uids.length === 0 || uids.length > 100) {
    ElMessage.warning('请填写 1 至 100 个用户 ID，并用逗号或换行分隔');
    return;
  }
  saving.value = true;
  try {
    await sendMerchantCouponApi(sendingCoupon.value.coupon_id, { mark: sendForm.mark.trim(), uids });
    sendOpen.value = false;
    ElMessage.success('优惠券已发送');
    await load();
  } finally {
    saving.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="店铺优惠券" description="管理本店券模板，并可向已有支付订单的本店用户定向发送优惠券。">
    <template #extra><el-button type="primary" @click="openCreate">新增优惠券</el-button></template>
    <el-card shadow="never"><el-table v-loading="loading" :data="rows" row-key="coupon_id"><el-table-column label="ID" prop="coupon_id" width="80" /><el-table-column label="优惠券名称" min-width="180" prop="title" /><el-table-column label="面额" width="104"><template #default="{ row }">¥{{ Number(row.coupon_price).toFixed(2) }}</template></el-table-column><el-table-column label="使用门槛" width="120"><template #default="{ row }">满 ¥{{ Number(row.use_min_price).toFixed(2) }} 可用</template></el-table-column><el-table-column label="有效期" width="108"><template #default="{ row }">领取后 {{ row.coupon_time }} 天</template></el-table-column><el-table-column label="发放数量" width="112"><template #default="{ row }">{{ row.is_limited === 1 ? `${row.remain_count}/${row.total_count}` : '不限量' }}</template></el-table-column><el-table-column label="状态" width="88"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '已启用' : '已停用' }}</el-tag></template></el-table-column><el-table-column label="创建时间" min-width="170" prop="create_time" /><el-table-column fixed="right" label="操作" width="222"><template #default="{ row }"><el-button link type="success" :disabled="row.status !== 1" @click="openSend(row)">发券</el-button><el-button link type="primary" @click="openEdit(row)">编辑</el-button><el-button link type="warning" @click="toggle(row)">{{ row.status === 1 ? '停用' : '启用' }}</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑优惠券' : '新增优惠券'" width="580px" destroy-on-close><el-form class="grid grid-cols-2 gap-x-4" label-width="90px"><el-form-item class="col-span-2" label="优惠券名称" required><el-input v-model="form.title" maxlength="40" show-word-limit /></el-form-item><el-form-item label="优惠金额" required><el-input-number v-model="form.coupon_price" :min="0.01" :precision="2" class="w-full" /></el-form-item><el-form-item label="使用门槛"><el-input-number v-model="form.use_min_price" :min="0" :precision="2" class="w-full" /></el-form-item><el-form-item label="有效天数" required><el-input-number v-model="form.coupon_time" :min="1" class="w-full" /></el-form-item><el-form-item label="排序"><el-input-number v-model="form.sort" :min="0" class="w-full" /></el-form-item><el-form-item label="限量发放"><el-switch v-model="form.is_limited" :active-value="1" :inactive-value="0" /></el-form-item><el-form-item v-if="form.is_limited === 1" label="发放总数" required><el-input-number v-model="form.total_count" :min="1" class="w-full" /></el-form-item><el-form-item label="初始状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" /></el-form-item></el-form><template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">保存</el-button></template></el-dialog>
    <el-dialog v-model="sendOpen" title="定向发送优惠券" width="540px" destroy-on-close><el-alert :closable="false" class="mb-4" title="仅允许向本商户存在已支付订单的用户发送；同一用户不可重复领取同一张券。" type="warning" /><el-form label-width="88px"><el-form-item label="优惠券"><span>{{ sendingCoupon?.title }}</span></el-form-item><el-form-item label="用户 ID" required><el-input v-model="sendForm.uidsText" :rows="4" placeholder="输入用户 ID，多个 ID 用逗号或换行分隔，最多 100 个" type="textarea" /></el-form-item><el-form-item label="发送说明"><el-input v-model="sendForm.mark" :rows="3" maxlength="200" show-word-limit type="textarea" /></el-form-item></el-form><template #footer><el-button @click="sendOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="send">确认发送</el-button></template></el-dialog>
  </Page>
</template>
