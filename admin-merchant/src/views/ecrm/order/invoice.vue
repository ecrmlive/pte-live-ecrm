<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage, ElMessageBox, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import { auditInvoice, fetchInvoices, type InvoiceRow } from '#/api/core/ecrm';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canAudit = ref(false);

const statusMap: Record<
  number,
  { label: string; type: 'danger' | 'info' | 'success' | 'warning' }
> = {
  [-1]: { label: '已驳回', type: 'danger' },
  0: { label: '待审核', type: 'warning' },
  1: { label: '已开票', type: 'success' },
};

function statusInfo(status: number) {
  return statusMap[status] || { label: '未知', type: 'info' as const };
}

function rowCanAudit(row: InvoiceRow) {
  return canAudit.value && row.status === 0;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '订单 ID' },
    fieldName: 'order_id',
    label: '订单 ID',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '抬头 / 税号 / 邮箱',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 'requested' },
        { label: '已开票', value: 'issued' },
        { label: '已驳回', value: 'rejected' },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<InvoiceRow> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'invoice_id', title: 'ID', width: 80 },
    { field: 'order_id', title: '订单', width: 100 },
    { field: 'uid', title: '用户', width: 90 },
    {
      field: 'header',
      minWidth: 160,
      showOverflow: false,
      title: '抬头',
    },
    {
      field: 'tax_no',
      minWidth: 160,
      showOverflow: false,
      title: '税号',
    },
    {
      field: 'email',
      minWidth: 160,
      showOverflow: false,
      title: '邮箱',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    {
      field: 'mark',
      minWidth: 140,
      showOverflow: false,
      title: '备注/票号',
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '申请时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 140 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const orderIdRaw = String(formValues?.order_id ?? '').trim();
        const orderId = orderIdRaw ? Number(orderIdRaw) : undefined;
        const status = String(formValues?.status ?? '').trim();
        const result = await fetchInvoices({
          page: page.currentPage,
          limit: page.pageSize,
          status: status || undefined,
          order_id:
            orderId && Number.isFinite(orderId) && orderId > 0
              ? orderId
              : undefined,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'invoice_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function onAudit(row: InvoiceRow, status: 1 | -1) {
  try {
    const { value } = await ElMessageBox.prompt(
      status === 1 ? '请填写发票号码（可留空自动生成）' : '请填写驳回原因',
      status === 1 ? '确认开票' : '驳回发票',
      {
        confirmButtonText: status === 1 ? '开票' : '驳回',
        cancelButtonText: '取消',
        inputPlaceholder: status === 1 ? '发票号码' : '驳回原因',
        inputValue: status === 1 ? '' : '商户驳回',
        type: status === 1 ? 'info' : 'warning',
      },
    );
    await auditInvoice(row.invoice_id, {
      status,
      mark: value?.trim() || undefined,
    });
    ElMessage.success(status === 1 ? '已开票' : '已驳回');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canAudit.value = codes.includes('invoice.audit');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status).type" size="small">
          {{ statusInfo(row.status).label }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="rowCanAudit(row)">
          <ElButton link type="primary" @click="onAudit(row, 1)">开票</ElButton>
          <ElButton link type="danger" @click="onAudit(row, -1)">驳回</ElButton>
        </template>
        <span v-else class="text-muted-foreground">—</span>
      </template>
    </Grid>
  </Page>
</template>
