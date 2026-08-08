<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  addMerchantRefundRemarkApi,
  approveMerchantRefundApi,
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
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

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

const statusMap: Record<
  number,
  { label: string; type: 'danger' | 'info' | 'success' | 'warning' }
> = {
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

function buildRefundQuery(
  formValues: Record<string, unknown> | undefined,
  page: { currentPage: number; pageSize: number },
) {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  const status = formValues?.status;
  const refundTypeValue = formValues?.refund_type;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    order_sn: String(formValues?.order_sn ?? '').trim() || undefined,
    phone: String(formValues?.phone ?? '').trim() || undefined,
    date_from: range[0],
    date_to: range[1],
    status:
      typeof status === 'number' &&
      [-2, -1, 0, 1, 2, 3, 4].includes(Number(status))
        ? Number(status)
        : undefined,
    refund_type:
      refundTypeValue === 1 || refundTypeValue === 2
        ? Number(refundTypeValue)
        : undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '退款单号 / 订单号 / 关键词',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '订单号' },
    fieldName: 'order_sn',
    label: '订单号',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 0 },
        { label: '待退货', value: 1 },
        { label: '待收货', value: 2 },
        { label: '已退款', value: 3 },
        { label: '平台介入', value: 4 },
        { label: '审核拒绝', value: -1 },
        { label: '用户已取消', value: -2 },
      ],
      placeholder: '请选择',
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
    componentProps: { clearable: true, placeholder: '用户手机号' },
    fieldName: 'phone',
    label: '手机号',
  },
]);

const gridOptions: VxeGridProps<MerchantRefundOrder> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    {
      field: 'refund_order_sn',
      minWidth: 170,
      showOverflow: false,
      title: '退款单号',
    },
    { field: 'order_id', title: '订单 ID', width: 100 },
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
    { field: 'refund_num', title: '退款件数', width: 96 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 110,
    },
    {
      field: 'refund_message',
      minWidth: 180,
      showOverflow: true,
      title: '申请原因',
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '申请时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 240 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listMerchantRefundsApi(buildRefundQuery(formValues, page));
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

const [RejectDrawer, rejectDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  title: '拒绝退款',
  confirmText: '确认拒绝',
  onConfirm: async () => {
    const message = rejectForm.failMessage.trim();
    if (!message) {
      ElMessage.warning('请填写拒绝原因');
      return;
    }
    if (!current.value) return;
    rejecting.value = true;
    rejectDrawerApi.lock();
    try {
      await rejectMerchantRefundApi(current.value.refund_order_id, message);
      ElMessage.success('退款申请已拒绝');
      rejectDrawerApi.close();
      gridApi.reload();
    } finally {
      rejecting.value = false;
      rejectDrawerApi.unlock();
    }
  },
});

const [LogDrawer, logDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  title: '退款操作日志',
  showConfirmButton: false,
  cancelText: '关闭',
});

async function openDetail(row: MerchantRefundOrder) {
  detailDrawerApi.setState({ title: '退款详情', loading: true }).open();
  try {
    current.value = await getMerchantRefundApi(row.refund_order_id);
  } finally {
    detailDrawerApi.setState({ loading: false });
  }
}

async function openLog(row: MerchantRefundOrder) {
  current.value = row;
  events.value = [];
  logDrawerApi.open();
  try {
    const result = await listMerchantRefundEventsApi(row.refund_order_id);
    events.value = result.list || [];
  } catch {
    logDrawerApi.close();
  }
}

async function loadExpress(row: MerchantRefundOrder) {
  try {
    shipment.value = await getMerchantRefundExpressApi(row.refund_order_id);
    ElMessage.success(`已获取退货物流：${shipment.value.carrier_name}`);
  } catch {
    // 由请求层提示
  }
}

async function exportRows() {
  try {
    const formValues = (await gridApi.formApi?.getValues()) || {};
    const result = await exportMerchantRefundsApi(
      buildRefundQuery(formValues, { currentPage: 1, pageSize: 10 }),
    );
    const blob = new Blob([result.content], { type: 'text/csv;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      `已导出 ${result.row_count} 条店铺退款记录${result.truncated ? '（已按 5000 条上限截断）' : ''}`,
    );
  } catch {
    // 由请求层提示
  }
}

async function addRemark(row: MerchantRefundOrder) {
  try {
    const { value } = await ElMessageBox.prompt(
      '备注仅记录在店铺操作审计中，不会更改退款原因、状态或资金处理。',
      '新增售后备注',
      {
        inputPattern: /\S/,
        inputErrorMessage: '备注不能为空',
        inputType: 'textarea',
        confirmButtonText: '保存备注',
        cancelButtonText: '取消',
      },
    );
    await addMerchantRefundRemarkApi(row.refund_order_id, {
      note: value.trim(),
      idempotency_key: crypto.randomUUID(),
    });
    ElMessage.success('售后备注已记录');
  } catch {
    // 用户取消或接口错误
  }
}

async function hideRefund(row: MerchantRefundOrder) {
  try {
    const { value } = await ElMessageBox.prompt(
      '隐藏后仅本店后台不再显示该退款单，平台监管、退款状态和资金记录不会删除。请填写原因。',
      '隐藏退款单',
      {
        inputPattern: /.{2,}/,
        inputErrorMessage: '隐藏原因至少 2 个字符',
        inputType: 'textarea',
        confirmButtonText: '确认隐藏',
        cancelButtonText: '取消',
        type: 'warning',
      },
    );
    await hideMerchantRefundApi(row.refund_order_id, {
      reason: value.trim(),
      idempotency_key: crypto.randomUUID(),
    });
    ElMessage.success('退款单已从本店列表隐藏');
    gridApi.reload();
  } catch {
    // 用户取消或接口错误
  }
}

async function approve(row: MerchantRefundOrder) {
  const returnAndRefund = row.refund_type === 2;
  try {
    await confirm({
      content: returnAndRefund
        ? '确认同意退货退款申请？用户将进入待退货状态。'
        : '确认同意仅退款申请？售后单将进入支付渠道退款阶段；已验证的支付回调才会标记为已退款。',
      icon: 'warning',
      title: '同意售后确认',
    });
    await approveMerchantRefundApi(row.refund_order_id);
    ElMessage.success(
      returnAndRefund
        ? '已同意退货退款，等待用户寄回商品'
        : '已受理仅退款，等待支付渠道退款回调',
    );
    gridApi.reload();
  } catch {
    // cancelled
  }
}

async function confirmReturn(row: MerchantRefundOrder) {
  try {
    await confirm({
      content:
        '确认已收到用户寄回的商品？确认后售后单进入支付渠道退款阶段，无法撤销。',
      icon: 'warning',
      title: '确认收货退款',
    });
    await confirmMerchantRefundReturnApi(row.refund_order_id);
    ElMessage.success('已确认收货，等待支付渠道退款回调');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

function openReject(row: MerchantRefundOrder) {
  current.value = row;
  rejectForm.failMessage = '';
  rejectDrawerApi.open();
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
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
    <Grid>
      <template #toolbar-actions>
        <ElButton v-if="canExport" type="success" @click="exportRows">
          导出
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status).type">
          {{ statusInfo(row.status).label }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          v-if="rowCanApprove(row)"
          link
          type="success"
          @click="approve(row)"
        >
          同意
        </ElButton>
        <ElButton
          v-if="rowCanReject(row)"
          link
          type="danger"
          @click="openReject(row)"
        >
          拒绝
        </ElButton>
        <ElButton
          v-if="rowCanConfirmReturn(row)"
          link
          type="success"
          @click="confirmReturn(row)"
        >
          确认收货
        </ElButton>
        <ElButton v-if="canLog" link @click="openLog(row)">日志</ElButton>
        <ElButton
          v-if="canExpress && row.refund_type === 2"
          link
          @click="loadExpress(row)"
        >
          退货物流
        </ElButton>
        <ElButton v-if="canRemark" link @click="addRemark(row)">备注</ElButton>
        <ElButton v-if="canDelete" link type="danger" @click="hideRefund(row)">
          隐藏
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="current">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="退款单号">
            {{ current.refund_order_sn }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="订单 ID">
            {{ current.order_id }}
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
            <ElTag :type="statusInfo(current.status).type">
              {{ statusInfo(current.status).label }}
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
            v-if="
              shipment &&
              shipment.tracking_no === current.return_shipment?.tracking_no
            "
            :span="2"
            label="物流查询快照"
          >
            {{ shipment.carrier_name }} / {{ shipment.tracking_no }}（{{
              formatShanghaiDateTime(shipment.submitted_at)
            }}）
          </ElDescriptionsItem>
          <ElDescriptionsItem
            v-if="current.fail_message"
            :span="2"
            label="拒绝原因"
          >
            {{ current.fail_message }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="申请时间">
            {{ formatShanghaiDateTime(current.create_time) }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <div class="mb-3 mt-6 text-base font-medium">退款商品</div>
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
    </DetailDrawer>

    <RejectDrawer>
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
    </RejectDrawer>

    <LogDrawer class="w-[760px] max-w-[96vw]">
      <ElTable v-if="events.length" :data="events" border>
        <ElTableColumn
          label="时间"
          min-width="166"
          prop="created_at"
          :formatter="(_row, _col, val) => formatShanghaiDateTime(val)"
        />
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
      <div v-else class="py-8 text-center text-muted-foreground">
        暂无状态流转日志
      </div>
    </LogDrawer>
  </Page>
</template>
