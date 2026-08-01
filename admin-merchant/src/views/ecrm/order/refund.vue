<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  approveMerchantRefundApi,
  getMerchantRefundApi,
  listMerchantRefundsApi,
  rejectMerchantRefundApi,
  type MerchantRefundOrder,
} from '#/api/core/merchant-aftersale';

const loading = ref(false);
const rows = ref<MerchantRefundOrder[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const rejectOpen = ref(false);
const current = ref<MerchantRefundOrder>();
const rejecting = ref(false);
const rejectForm = reactive({ failMessage: '' });
const query = reactive({ limit: 20, page: 1, status: undefined as number | undefined });

const statusMap: Record<number, { label: string; type: 'danger' | 'info' | 'success' | 'warning' }> = {
  [-2]: { label: '用户已取消', type: 'info' },
  [-1]: { label: '审核拒绝', type: 'danger' },
  0: { label: '待审核', type: 'warning' },
  1: { label: '待退货', type: 'warning' },
  2: { label: '待收货', type: 'warning' },
  3: { label: '已退款', type: 'success' },
  4: { label: '平台介入', type: 'danger' },
};

function statusInfo(status: number) {
  return statusMap[status] || { label: '未知状态', type: 'info' as const };
}

function refundType(type: number) {
  return type === 1 ? '仅退款' : type === 2 ? '退货退款' : '未知';
}

function canAudit(row: MerchantRefundOrder) {
  return row.status === 0 && row.refund_type === 1;
}

async function load() {
  loading.value = true;
  try {
    const result = await listMerchantRefundsApi(query);
    rows.value = result.list;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.status = undefined;
  query.page = 1;
  void load();
}

async function openDetail(row: MerchantRefundOrder) {
  const result = await getMerchantRefundApi(row.refund_order_id);
  current.value = result;
  detailOpen.value = true;
}

async function approve(row: MerchantRefundOrder) {
  try {
    await ElMessageBox.confirm(
      '确认同意该仅退款申请？操作会执行退款、恢复商品库存并更新用户余额，提交后不可撤销。',
      '同意退款确认',
      { confirmButtonText: '确认同意', cancelButtonText: '取消', type: 'warning' },
    );
    await approveMerchantRefundApi(row.refund_order_id);
    ElMessage.success('退款已处理');
    await load();
  } catch {
    // 用户取消或接口已返回错误时，requestClient 统一处理提示。
  }
}

function openReject(row: MerchantRefundOrder) {
  current.value = row;
  rejectForm.failMessage = '';
  rejectOpen.value = true;
}

async function reject() {
  const message = rejectForm.failMessage.trim();
  if (!message) {
    ElMessage.warning('请填写拒绝原因');
    return;
  }
  if (!current.value) return;
  rejecting.value = true;
  try {
    await rejectMerchantRefundApi(current.value.refund_order_id, message);
    rejectOpen.value = false;
    ElMessage.success('退款申请已拒绝');
    await load();
  } finally {
    rejecting.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="退款订单" description="处理本商户仅退款申请；退货退款、平台介入等状态仅展示，不在此页执行越权处理。">
    <el-card shadow="never">
      <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
        <el-form-item label="退款状态">
          <el-select v-model="query.status" clearable class="w-44" placeholder="全部状态">
            <el-option label="待审核" :value="0" />
            <el-option label="待退货" :value="1" />
            <el-option label="待收货" :value="2" />
            <el-option label="已退款" :value="3" />
            <el-option label="平台介入" :value="4" />
            <el-option label="审核拒绝" :value="-1" />
            <el-option label="用户已取消" :value="-2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">查询</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="refund_order_id">
        <el-table-column label="退款单号" min-width="180" prop="refund_order_sn" />
        <el-table-column label="订单 ID" width="100" prop="order_id" />
        <el-table-column label="用户 ID" width="100" prop="uid" />
        <el-table-column label="售后类型" width="104">
          <template #default="{ row }">{{ refundType(row.refund_type) }}</template>
        </el-table-column>
        <el-table-column label="退款金额" width="116">
          <template #default="{ row }">¥{{ Number(row.refund_price).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="退款件数" width="96" prop="refund_num" />
        <el-table-column label="状态" width="110">
          <template #default="{ row }"><el-tag :type="statusInfo(row.status).type">{{ statusInfo(row.status).label }}</el-tag></template>
        </el-table-column>
        <el-table-column label="申请原因" min-width="180" prop="refund_message" show-overflow-tooltip />
        <el-table-column label="申请时间" min-width="170" prop="create_time" />
        <el-table-column fixed="right" label="操作" width="172">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
            <template v-if="canAudit(row)">
              <el-button link type="success" @click="approve(row)">同意</el-button>
              <el-button link type="danger" @click="openReject(row)">拒绝</el-button>
            </template>
          </template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination
          :current-page="query.page"
          :page-size="query.limit"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          background
          layout="total, sizes, prev, pager, next"
          @current-change="(page) => { query.page = page; load(); }"
          @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }"
        />
      </div>
    </el-card>

    <el-drawer v-model="detailOpen" :with-header="false" size="640px">
      <template v-if="current">
        <div class="mb-5 text-lg font-medium">退款详情</div>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="退款单号">{{ current.refund_order_sn }}</el-descriptions-item>
          <el-descriptions-item label="订单 ID">{{ current.order_id }}</el-descriptions-item>
          <el-descriptions-item label="用户 ID">{{ current.uid }}</el-descriptions-item>
          <el-descriptions-item label="售后类型">{{ refundType(current.refund_type) }}</el-descriptions-item>
          <el-descriptions-item label="退款金额">¥{{ Number(current.refund_price).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="退款件数">{{ current.refund_num }}</el-descriptions-item>
          <el-descriptions-item label="状态"><el-tag :type="statusInfo(current.status).type">{{ statusInfo(current.status).label }}</el-tag></el-descriptions-item>
          <el-descriptions-item label="状态时间">{{ current.status_time }}</el-descriptions-item>
          <el-descriptions-item :span="2" label="申请原因">{{ current.refund_message || '—' }}</el-descriptions-item>
          <el-descriptions-item v-if="current.fail_message" :span="2" label="拒绝原因">{{ current.fail_message }}</el-descriptions-item>
        </el-descriptions>
        <div class="mb-3 mt-6 text-base font-medium">退款商品</div>
        <el-table :data="current.products || []" border>
          <el-table-column label="订单商品 ID" prop="order_product_id" min-width="130" />
          <el-table-column label="退款金额" min-width="110"><template #default="{ row }">¥{{ Number(row.refund_price).toFixed(2) }}</template></el-table-column>
          <el-table-column label="退款数量" prop="refund_num" min-width="100" />
        </el-table>
      </template>
    </el-drawer>

    <el-dialog v-model="rejectOpen" title="拒绝退款" width="480px" destroy-on-close>
      <el-form label-width="84px"><el-form-item label="拒绝原因" required><el-input v-model="rejectForm.failMessage" :rows="4" maxlength="200" placeholder="请向用户说明拒绝原因" show-word-limit type="textarea" /></el-form-item></el-form>
      <template #footer><el-button @click="rejectOpen = false">取消</el-button><el-button :loading="rejecting" type="danger" @click="reject">确认拒绝</el-button></template>
    </el-dialog>
  </Page>
</template>
