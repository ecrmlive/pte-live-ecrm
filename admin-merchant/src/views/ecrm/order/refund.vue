<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  approveMerchantRefundApi,
  addMerchantRefundRemarkApi,
  confirmMerchantRefundReturnApi,
  exportMerchantRefundsApi,
  getMerchantRefundApi,
  getMerchantRefundExpressApi,
  hideMerchantRefundApi,
  listMerchantRefundEventsApi,
  listMerchantRefundsApi,
  rejectMerchantRefundApi,
  type MerchantRefundEvent,
  type MerchantRefundOrder,
} from '#/api/core/merchant-aftersale';
import { getAccessCodesApi } from '#/api/core/auth';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<MerchantRefundOrder[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const rejectOpen = ref(false);
const logOpen = ref(false);
const current = ref<MerchantRefundOrder>();
const events = ref<MerchantRefundEvent[]>([]);
const shipment = ref<Awaited<ReturnType<typeof getMerchantRefundExpressApi>>>();
const rejecting = ref(false);
const canLog = ref(false);
const canExport = ref(false);
const canExpress = ref(false);
const canRemark = ref(false);
const canDelete = ref(false);
const canApprove = ref(false);
const canReject = ref(false);
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
  return row.status === 0;
}

function rowCanApprove(row: MerchantRefundOrder) {
  return canApprove.value && canAudit(row);
}

function rowCanReject(row: MerchantRefundOrder) {
  return canReject.value && canAudit(row);
}

function rowCanConfirmReturn(row: MerchantRefundOrder) {
  return canApprove.value && row.status === 2 && row.refund_type === 2;
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

async function openLog(row: MerchantRefundOrder) { current.value = row; events.value = []; logOpen.value = true; try { const result = await listMerchantRefundEventsApi(row.refund_order_id); events.value = result.list || []; } catch { logOpen.value = false; } }
async function loadExpress(row: MerchantRefundOrder) { try { shipment.value = await getMerchantRefundExpressApi(row.refund_order_id); ElMessage.success(`已获取退货物流：${shipment.value.carrier_name}`); } catch { /* 由请求层提示 */ } }
async function exportRows() { try { const result = await exportMerchantRefundsApi({ status: query.status }); const blob = new Blob([result.content], { type: 'text/csv;charset=utf-8' }); const url = URL.createObjectURL(blob); const link = document.createElement('a'); link.href = url; link.download = result.file_name; link.click(); URL.revokeObjectURL(url); ElMessage.success(`已导出 ${result.row_count} 条店铺退款记录${result.truncated ? '（已按 5000 条上限截断）' : ''}`); } catch { /* 由请求层提示 */ } }
async function addRemark(row: MerchantRefundOrder) { try { const { value } = await ElMessageBox.prompt('备注仅记录在店铺操作审计中，不会更改退款原因、状态或资金处理。', '添加售后备注', { inputPattern: /\S/, inputErrorMessage: '备注不能为空', inputType: 'textarea', confirmButtonText: '保存备注', cancelButtonText: '取消' }); await addMerchantRefundRemarkApi(row.refund_order_id, { note: value.trim(), idempotency_key: crypto.randomUUID() }); ElMessage.success('售后备注已记录'); } catch { /* 用户取消或接口错误 */ } }
async function hideRefund(row: MerchantRefundOrder) { try { const { value } = await ElMessageBox.prompt('隐藏后仅本店后台不再显示该退款单，平台监管、退款状态和资金记录不会删除。请填写原因。', '隐藏退款单', { inputPattern: /.{2,}/, inputErrorMessage: '隐藏原因至少 2 个字符', inputType: 'textarea', confirmButtonText: '确认隐藏', cancelButtonText: '取消', type: 'warning' }); await hideMerchantRefundApi(row.refund_order_id, { reason: value.trim(), idempotency_key: crypto.randomUUID() }); ElMessage.success('退款单已从本店列表隐藏'); await load(); } catch { /* 用户取消或接口错误 */ } }

async function approve(row: MerchantRefundOrder) {
	const returnAndRefund = row.refund_type === 2;
  try {
    await ElMessageBox.confirm(
      returnAndRefund
        ? '确认同意退货退款申请？用户将进入待退货状态。'
        : '确认同意仅退款申请？售后单将进入支付渠道退款阶段；已验证的支付回调才会标记为已退款。',
      '同意售后确认',
      { confirmButtonText: '确认同意', cancelButtonText: '取消', type: 'warning' },
    );
    await approveMerchantRefundApi(row.refund_order_id);
    ElMessage.success(returnAndRefund ? '已同意退货退款，等待用户寄回商品' : '已受理仅退款，等待支付渠道退款回调');
    await load();
  } catch {
    // 用户取消或接口已返回错误时，requestClient 统一处理提示。
  }
}

async function confirmReturn(row: MerchantRefundOrder) {
	try {
		await ElMessageBox.confirm('确认已收到用户寄回的商品？确认后售后单进入支付渠道退款阶段，无法撤销。', '确认收货退款', { confirmButtonText: '确认收货', cancelButtonText: '取消', type: 'warning' });
		await confirmMerchantRefundReturnApi(row.refund_order_id);
		ElMessage.success('已确认收货，等待支付渠道退款回调');
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

onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), load()]);
  canApprove.value = permissions.includes('refund.approve');
  canReject.value = permissions.includes('refund.reject');
  canLog.value = permissions.includes('refund.log');
  canExport.value = permissions.includes('refund.export');
  canExpress.value = permissions.includes('refund.express');
  canRemark.value = permissions.includes('refund.remark');
  canDelete.value = permissions.includes('refund.delete');
});
</script>

<template>
  <Page auto-content-height>
    <EcrmListPage
      description="处理本店售后：仅退款进入支付渠道退款阶段；退货退款依次完成同意退货、用户寄回、商户确认收货。资金终态仅由已验证支付回调写入。"
      title="退款订单"
    >
      <template #filters>
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
      </template>

      <template #actions>
        <el-button v-if="canExport" type="success" @click="exportRows">导出</el-button>
      </template>

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
        <el-table-column fixed="right" label="操作" width="210">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
            <el-button v-if="rowCanApprove(row)" link type="success" @click="approve(row)">同意</el-button>
            <el-button v-if="rowCanReject(row)" link type="danger" @click="openReject(row)">拒绝</el-button>
            <el-button v-if="rowCanConfirmReturn(row)" link type="success" @click="confirmReturn(row)">确认收货</el-button>
            <el-button v-if="canLog" link @click="openLog(row)">日志</el-button>
            <el-button v-if="canExpress && row.refund_type === 2" link @click="loadExpress(row)">退货物流</el-button>
            <el-button v-if="canRemark" link @click="addRemark(row)">备注</el-button>
            <el-button v-if="canDelete" link type="danger" @click="hideRefund(row)">隐藏</el-button>
          </template>
        </el-table-column>
      </el-table>

      <template #pager>
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
      </template>
    </EcrmListPage>

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
          <el-descriptions-item v-if="current.return_shipment" :span="2" label="退货物流">{{ current.return_shipment.carrier_name }} / {{ current.return_shipment.tracking_no }}（{{ current.return_shipment.submitted_at }}）</el-descriptions-item>
          <el-descriptions-item v-if="current.return_shipment?.remark" :span="2" label="物流备注">{{ current.return_shipment.remark }}</el-descriptions-item>
          <el-descriptions-item v-if="shipment && shipment.tracking_no === current.return_shipment?.tracking_no" :span="2" label="物流查询快照">{{ shipment.carrier_name }} / {{ shipment.tracking_no }}（{{ shipment.submitted_at }}）</el-descriptions-item>
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
    <el-dialog v-model="logOpen" title="退款操作日志" width="760px" destroy-on-close>
      <el-empty v-if="events.length === 0" description="暂无状态流转日志" />
      <el-table v-else :data="events" border>
        <el-table-column label="时间" min-width="166" prop="created_at" />
        <el-table-column label="操作方" width="108"><template #default="{ row }">{{ ({ user: '用户', merchant: '商户', platform: '平台', system: '系统' } as Record<string, string>)[row.actor_type] || row.actor_type }}</template></el-table-column>
        <el-table-column label="状态流转" min-width="164"><template #default="{ row }">{{ row.from_status || '—' }} → {{ row.to_status }}</template></el-table-column>
        <el-table-column label="说明" min-width="220" prop="reason" show-overflow-tooltip />
      </el-table>
    </el-dialog>
  </Page>
</template>
