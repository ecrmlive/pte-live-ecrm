<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import {
  approvePlatformRefundApi,
  exportPlatformRefundsApi,
  getPlatformRefundApi,
  listPlatformRefundEventsApi,
  listPlatformRefundsApi,
  rejectPlatformRefundApi,
  type PlatformRefundEvent,
  type PlatformRefundOrder,
} from '#/api/core/platform-aftersale';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<PlatformRefundOrder[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const rejectOpen = ref(false);
const logOpen = ref(false);
const current = ref<PlatformRefundOrder>();
const events = ref<PlatformRefundEvent[]>([]);
const rejecting = ref(false);
const canApprove = ref(false);
const canReject = ref(false);
const canViewLog = ref(false);
const canExport = ref(false);
const rejectForm = reactive({ failMessage: '' });
const query = reactive({ limit: 20, page: 1, status: undefined as string | undefined });

const statusMap: Record<string, { label: string; type: 'danger' | 'info' | 'success' | 'warning' }> = {
  applied: { label: '待审核', type: 'warning' },
  merchant_handling: { label: '商户处理中', type: 'warning' },
  awaiting_return: { label: '待退货', type: 'warning' },
  awaiting_receipt: { label: '待收货', type: 'warning' },
  platform_intervene: { label: '平台介入', type: 'danger' },
  refunding: { label: '退款处理中', type: 'warning' },
  refunded: { label: '已退款', type: 'success' },
  rejected: { label: '审核拒绝', type: 'danger' },
  cancelled: { label: '用户已取消', type: 'info' },
};

function statusInfo(statusCode: string) {
  return statusMap[statusCode] || { label: '未知状态', type: 'info' as const };
}

function refundType(type: number) {
  return type === 1 ? '仅退款' : type === 2 ? '退货退款' : '未知';
}

function canAudit(row: PlatformRefundOrder) {
  return ['applied', 'merchant_handling', 'platform_intervene'].includes(row.status_code);
}

async function load() {
  loading.value = true;
  try {
    const result = await listPlatformRefundsApi(query);
    rows.value = result.list || [];
    total.value = result.total || 0;
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

async function openDetail(row: PlatformRefundOrder) {
  current.value = await getPlatformRefundApi(row.refund_order_id);
  detailOpen.value = true;
}

async function openLog(row: PlatformRefundOrder) {
  current.value = row;
  events.value = [];
  logOpen.value = true;
  try {
    const result = await listPlatformRefundEventsApi(row.refund_order_id);
    events.value = result.list || [];
  } catch {
    logOpen.value = false;
  }
}

async function exportRows() {
  try {
    const { value } = await ElMessageBox.prompt(
      '请填写导出原因。文件不含用户身份、退款原因及退货物流，最多 5000 行。',
      '导出退款监管清单',
      {
        inputPattern: /.{2,}/,
        inputErrorMessage: '导出原因至少 2 个字符',
        confirmButtonText: '生成 CSV',
        cancelButtonText: '取消',
      },
    );
    const result = await exportPlatformRefundsApi({ reason: value.trim(), status: query.status });
    const blob = new Blob([result.content], { type: 'text/csv;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      `已导出 ${result.row_count} 条退款监管记录${result.truncated ? '（已按 5000 条上限截断）' : ''}`,
    );
  } catch {
    // 用户取消或接口错误时由统一请求层提示。
  }
}

async function approve(row: PlatformRefundOrder) {
  const returnAndRefund = row.refund_type === 2;
  try {
    await ElMessageBox.confirm(
      returnAndRefund
        ? '确认同意退货退款申请？用户将进入待退货状态；商户确认收货后才会进入支付渠道退款阶段。'
        : '确认同意仅退款申请？售后单将进入支付渠道退款阶段；已验证的支付回调才会标记为已退款。',
      '同意售后确认',
      { confirmButtonText: '确认同意', cancelButtonText: '取消', type: 'warning' },
    );
    await approvePlatformRefundApi(row.refund_order_id);
    ElMessage.success(
      returnAndRefund ? '已同意退货退款，等待用户寄回商品' : '已受理仅退款，等待支付渠道退款回调',
    );
    await load();
  } catch {
    // 用户取消或接口已返回错误时，requestClient 统一处理提示。
  }
}

function openReject(row: PlatformRefundOrder) {
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
    await rejectPlatformRefundApi(current.value.refund_order_id, message);
    rejectOpen.value = false;
    ElMessage.success('退款申请已拒绝');
    await load();
  } finally {
    rejecting.value = false;
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi(), load()]);
  const isPlatform = profile.roles.includes('platform') && profile.is_agent !== 1;
  canApprove.value = isPlatform && permissions.includes('order.refund.approve');
  canReject.value = isPlatform && permissions.includes('order.refund.reject');
  canViewLog.value = isPlatform && permissions.includes('order.refund.log');
  canExport.value = isPlatform && permissions.includes('order.refund.export');
});
</script>

<template>
  <EcrmListPage
    title="退款监管"
    description="平台账号可处理待审核或平台介入的售后；区域账号仅查看所属区域商户数据。仅退款进入支付渠道退款阶段；退货退款需等待用户寄回与商户确认收货。"
  >
    <template #filters>
      <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
        <el-form-item label="退款状态">
          <el-select v-model="query.status" clearable class="w-44" placeholder="全部状态">
            <el-option label="待审核" value="applied" />
            <el-option label="商户处理中" value="merchant_handling" />
            <el-option label="待退货" value="awaiting_return" />
            <el-option label="待收货" value="awaiting_receipt" />
            <el-option label="退款处理中" value="refunding" />
            <el-option label="已退款" value="refunded" />
            <el-option label="平台介入" value="platform_intervene" />
            <el-option label="审核拒绝" value="rejected" />
            <el-option label="用户已取消" value="cancelled" />
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
      <el-table-column label="商户 ID" width="100" prop="mer_id" />
      <el-table-column label="用户 ID" width="100" prop="uid" />
      <el-table-column label="售后类型" width="104">
        <template #default="{ row }">{{ refundType(row.refund_type) }}</template>
      </el-table-column>
      <el-table-column label="退款金额" width="116">
        <template #default="{ row }">¥{{ Number(row.refund_price).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusInfo(row.status_code).type">{{ statusInfo(row.status_code).label }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="申请原因" min-width="180" prop="refund_message" show-overflow-tooltip />
      <el-table-column label="申请时间" min-width="170" prop="create_time" />
      <el-table-column fixed="right" label="操作" width="172">
        <template #default="{ row }">
          <el-button link type="primary" @click="openDetail(row)">详情</el-button>
          <template v-if="canAudit(row)">
            <el-button v-if="canApprove" link type="success" @click="approve(row)">同意</el-button>
            <el-button v-if="canReject" link type="danger" @click="openReject(row)">拒绝</el-button>
          </template>
          <el-button v-if="canViewLog" link @click="openLog(row)">日志</el-button>
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
        @current-change="(page: number) => { query.page = page; load(); }"
        @size-change="(limit: number) => { query.limit = limit; query.page = 1; load(); }"
      />
    </template>
  </EcrmListPage>

  <el-drawer v-model="detailOpen" :with-header="false" size="640px">
    <template v-if="current">
      <div class="mb-5 text-lg font-medium">退款详情</div>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="退款单号">{{ current.refund_order_sn }}</el-descriptions-item>
        <el-descriptions-item label="订单 ID">{{ current.order_id }}</el-descriptions-item>
        <el-descriptions-item label="商户 ID">{{ current.mer_id }}</el-descriptions-item>
        <el-descriptions-item label="用户 ID">{{ current.uid }}</el-descriptions-item>
        <el-descriptions-item label="售后类型">{{ refundType(current.refund_type) }}</el-descriptions-item>
        <el-descriptions-item label="退款金额">
          ¥{{ Number(current.refund_price).toFixed(2) }}
        </el-descriptions-item>
        <el-descriptions-item label="退款件数">{{ current.refund_num }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusInfo(current.status_code).type">
            {{ statusInfo(current.status_code).label }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态时间">{{ current.status_time }}</el-descriptions-item>
        <el-descriptions-item :span="2" label="申请原因">
          {{ current.refund_message || '—' }}
        </el-descriptions-item>
        <el-descriptions-item v-if="current.return_shipment" :span="2" label="退货物流">
          {{ current.return_shipment.carrier_name }} / {{ current.return_shipment.tracking_no }}（{{
            current.return_shipment.submitted_at
          }}）
        </el-descriptions-item>
        <el-descriptions-item v-if="current.return_shipment?.remark" :span="2" label="物流备注">
          {{ current.return_shipment.remark }}
        </el-descriptions-item>
        <el-descriptions-item v-if="current.fail_message" :span="2" label="拒绝原因">
          {{ current.fail_message }}
        </el-descriptions-item>
      </el-descriptions>
      <div class="mb-3 mt-6 text-base font-medium">退款商品</div>
      <el-table :data="current.products || []" border>
        <el-table-column label="订单商品 ID" min-width="130" prop="order_product_id" />
        <el-table-column label="退款金额" min-width="110">
          <template #default="{ row }">¥{{ Number(row.refund_price).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="退款数量" min-width="100" prop="refund_num" />
      </el-table>
    </template>
  </el-drawer>

  <el-dialog v-model="rejectOpen" destroy-on-close title="拒绝退款" width="480px">
    <el-form label-width="84px">
      <el-form-item label="拒绝原因" required>
        <el-input
          v-model="rejectForm.failMessage"
          :rows="4"
          maxlength="200"
          placeholder="请向用户说明拒绝原因"
          show-word-limit
          type="textarea"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="rejectOpen = false">取消</el-button>
      <el-button :loading="rejecting" type="danger" @click="reject">确认拒绝</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="logOpen" destroy-on-close title="退款操作日志" width="760px">
    <el-empty v-if="events.length === 0" description="暂无状态流转日志" />
    <el-table v-else :data="events" border>
      <el-table-column label="时间" min-width="166" prop="created_at" />
      <el-table-column label="操作方" width="108">
        <template #default="{ row }">
          {{
            ({ user: '用户', merchant: '商户', platform: '平台', system: '系统' } as Record<
              string,
              string
            >)[row.actor_type] || row.actor_type
          }}
        </template>
      </el-table-column>
      <el-table-column label="状态流转" min-width="164">
        <template #default="{ row }">{{ row.from_status || '—' }} → {{ row.to_status }}</template>
      </el-table-column>
      <el-table-column label="说明" min-width="220" prop="reason" show-overflow-tooltip />
    </el-table>
  </el-dialog>
</template>
