<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElDialog,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElSkeleton,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  approvePlatformRefundApi,
  exportPlatformRefundsApi,
  getPlatformRefundApi,
  listPlatformRefundEventsApi,
  listPlatformRefundsApi,
  rejectPlatformRefundApi,
  type PlatformRefundEvent,
  type PlatformRefundListParams,
  type PlatformRefundOrder,
} from '#/api/core/platform-aftersale';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const REFUND_STATUS_OPTIONS = [
  { label: '待审核', value: 'applied' },
  { label: '商户处理中', value: 'merchant_handling' },
  { label: '待退货', value: 'awaiting_return' },
  { label: '待收货', value: 'awaiting_receipt' },
  { label: '退款处理中', value: 'refunding' },
  { label: '已退款', value: 'refunded' },
  { label: '平台介入', value: 'platform_intervene' },
  { label: '审核拒绝', value: 'rejected' },
  { label: '用户已取消', value: 'cancelled' },
] as const;

const statusMap: Record<
  string,
  { label: string; type: 'danger' | 'info' | 'success' | 'warning' }
> = {
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

const current = ref<PlatformRefundOrder>();
const detailLoading = ref(false);
const events = ref<PlatformRefundEvent[]>([]);
const logOpen = ref(false);
const rejectOpen = ref(false);
const rejecting = ref(false);
const canApprove = ref(false);
const canReject = ref(false);
const canViewLog = ref(false);
const canExport = ref(false);
const rejectForm = reactive({ failMessage: '' });
const lastQueryParams = ref<Partial<PlatformRefundListParams>>({});

function statusInfo(statusCode: string) {
  return statusMap[statusCode] || { label: '未知状态', type: 'info' as const };
}

function refundType(type: number) {
  return type === 1 ? '仅退款' : type === 2 ? '退货退款' : '未知';
}

function canAudit(row: PlatformRefundOrder) {
  return ['applied', 'merchant_handling', 'platform_intervene'].includes(
    row.status_code,
  );
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
): PlatformRefundListParams {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  return {
    page: page.currentPage,
    limit: page.pageSize,
    status: String(formValues?.status ?? '').trim() || undefined,
    refund_order_sn:
      String(formValues?.refund_order_sn ?? '').trim() || undefined,
    refund_type:
      formValues?.refund_type === 1 || formValues?.refund_type === 2
        ? Number(formValues.refund_type)
        : undefined,
    order_sn: String(formValues?.order_sn ?? '').trim() || undefined,
    phone: String(formValues?.phone ?? '').trim() || undefined,
    real_name: String(formValues?.real_name ?? '').trim() || undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '退款单号' },
    fieldName: 'refund_order_sn',
    label: '退款单号',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [...REFUND_STATUS_OPTIONS],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '退款状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '仅退款', value: 1 },
        { label: '退货退款', value: 2 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'refund_type',
    label: '售后类型',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '关联订单号' },
    fieldName: 'order_sn',
    label: '订单号',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '收货人手机号' },
    fieldName: 'phone',
    label: '手机号',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '收货人姓名' },
    fieldName: 'real_name',
    label: '收货人',
  },
]);

const gridOptions: VxeGridProps<PlatformRefundOrder> = {
  columns: [
    { field: 'refund_order_sn', minWidth: 180, showOverflow: false, title: '退款单号' },
    { field: 'order_id', title: '订单 ID', width: 100 },
    { field: 'mer_id', title: '商户 ID', width: 100 },
    { field: 'uid', title: '用户 ID', width: 100 },
    {
      field: 'refund_type',
      title: '售后类型',
      width: 104,
      formatter: ({ cellValue }) => refundType(Number(cellValue)),
    },
    {
      field: 'refund_price',
      title: '退款金额',
      width: 116,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'status_code',
      slots: { default: 'status' },
      title: '状态',
      width: 110,
    },
    {
      field: 'refund_message',
      minWidth: 180,
      showOverflow: false,
      title: '申请原因',
      formatter: ({ cellValue }) => cellValue || '—',
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '申请时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    platformListActionColumn({ width: 220 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const params = buildListParams(page, formValues);
        lastQueryParams.value = params;
        const data = await listPlatformRefundsApi(params);
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'refund_order_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[900px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

async function openDetail(row: PlatformRefundOrder) {
  current.value = undefined;
  detailLoading.value = true;
  detailDrawerApi.setState({ title: '退款详情', loading: true }).open();
  try {
    current.value = await getPlatformRefundApi(row.refund_order_id);
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
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
    const { page: _page, limit: _limit, ...filters } = lastQueryParams.value;
    const result = await exportPlatformRefundsApi({
      ...filters,
      reason: value.trim(),
    });
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
      {
        confirmButtonText: '确认同意',
        cancelButtonText: '取消',
        type: 'warning',
      },
    );
    await approvePlatformRefundApi(row.refund_order_id);
    ElMessage.success(
      returnAndRefund
        ? '已同意退货退款，等待用户寄回商品'
        : '已受理仅退款，等待支付渠道退款回调',
    );
    gridApi.reload();
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
    gridApi.reload();
  } finally {
    rejecting.value = false;
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const isPlatform =
    profile.roles.includes('platform') && profile.is_agent !== 1;
  canApprove.value = isPlatform && permissions.includes('order.refund.approve');
  canReject.value = isPlatform && permissions.includes('order.refund.reject');
  canViewLog.value = isPlatform && permissions.includes('order.refund.log');
  canExport.value = isPlatform && permissions.includes('order.refund.export');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton v-if="canExport" type="success" @click="exportRows">
          导出
        </ElButton>
      </template>

      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status_code).type">
          {{ statusInfo(row.status_code).label }}
        </ElTag>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <template v-if="canAudit(row)">
          <ElButton
            v-if="canApprove"
            link
            type="success"
            @click="approve(row)"
          >
            同意
          </ElButton>
          <ElButton v-if="canReject" link type="danger" @click="openReject(row)">
            拒绝
          </ElButton>
        </template>
        <ElButton v-if="canViewLog" link @click="openLog(row)">日志</ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <template #default>
          <template v-if="current">
            <ElDescriptions :column="2" border>
              <ElDescriptionsItem label="退款单号">
                {{ current.refund_order_sn }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="订单 ID">
                {{ current.order_id }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="商户 ID">
                {{ current.mer_id }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="用户 ID">
                {{ current.uid }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="售后类型">
                {{ refundType(current.refund_type) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="退款金额">
                ¥{{ Number(current.refund_price).toFixed(2) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="退款件数">
                {{ current.refund_num }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="状态">
                <ElTag :type="statusInfo(current.status_code).type">
                  {{ statusInfo(current.status_code).label }}
                </ElTag>
              </ElDescriptionsItem>
              <ElDescriptionsItem label="状态时间">
                {{ formatShanghaiDateTime(current.status_time) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem :span="2" label="申请原因">
                {{ current.refund_message || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem
                v-if="current.return_shipment"
                :span="2"
                label="退货物流"
              >
                {{ current.return_shipment.carrier_name }} /
                {{ current.return_shipment.tracking_no }}（{{
                  formatShanghaiDateTime(current.return_shipment.submitted_at)
                }}）
              </ElDescriptionsItem>
              <ElDescriptionsItem
                v-if="current.return_shipment?.remark"
                :span="2"
                label="物流备注"
              >
                {{ current.return_shipment.remark }}
              </ElDescriptionsItem>
              <ElDescriptionsItem
                v-if="current.fail_message"
                :span="2"
                label="拒绝原因"
              >
                {{ current.fail_message }}
              </ElDescriptionsItem>
            </ElDescriptions>
            <div class="mb-2 mt-4 text-sm font-medium">退款商品</div>
            <ElTable :data="current.products || []" border>
              <ElTableColumn
                label="订单商品 ID"
                min-width="130"
                prop="order_product_id"
              />
              <ElTableColumn label="退款金额" min-width="110">
                <template #default="{ row }">
                  ¥{{ Number(row.refund_price).toFixed(2) }}
                </template>
              </ElTableColumn>
              <ElTableColumn label="退款数量" min-width="100" prop="refund_num" />
            </ElTable>
          </template>
        </template>
      </ElSkeleton>
    </DetailDrawer>

    <ElDialog v-model="rejectOpen" destroy-on-close title="拒绝退款" width="480px">
      <ElForm label-width="84px">
        <ElFormItem label="拒绝原因" required>
          <ElInput
            v-model="rejectForm.failMessage"
            :rows="4"
            maxlength="200"
            placeholder="请向用户说明拒绝原因"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
      <template #footer>
        <ElButton @click="rejectOpen = false">取消</ElButton>
        <ElButton :loading="rejecting" type="danger" @click="reject">
          确认拒绝
        </ElButton>
      </template>
    </ElDialog>

    <ElDialog
      v-model="logOpen"
      destroy-on-close
      title="退款操作日志"
      width="760px"
    >
      <ElEmpty v-if="events.length === 0" description="暂无状态流转日志" />
      <ElTable v-else :data="events" border>
        <ElTableColumn label="时间" min-width="166">
          <template #default="{ row }">
            {{ formatShanghaiDateTime(row.created_at) }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="操作方" width="108">
          <template #default="{ row }">
            {{
              (
                {
                  user: '用户',
                  merchant: '商户',
                  platform: '平台',
                  system: '系统',
                } as Record<string, string>
              )[row.actor_type] || row.actor_type
            }}
          </template>
        </ElTableColumn>
        <ElTableColumn label="状态流转" min-width="164">
          <template #default="{ row }">
            {{ row.from_status || '—' }} → {{ row.to_status }}
          </template>
        </ElTableColumn>
        <ElTableColumn
          label="说明"
          min-width="220"
          prop="reason"
          show-overflow-tooltip
        />
      </ElTable>
    </ElDialog>
  </Page>
</template>
