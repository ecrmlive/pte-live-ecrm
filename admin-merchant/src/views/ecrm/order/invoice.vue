<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { auditInvoice, fetchInvoices, type InvoiceRow } from '#/api/core/ecrm';

const gridOptions: VxeGridProps<InvoiceRow> = {
  border: true,
  columns: [
    { field: 'invoice_id', title: 'ID', width: 80 },
    { field: 'order_id', title: '订单', width: 100 },
    { field: 'uid', title: '用户', width: 90 },
    { field: 'header', minWidth: 160, title: '抬头' },
    { field: 'tax_no', title: '税号', width: 160 },
    { field: 'email', title: '邮箱', width: 160 },
    {
      field: 'status',
      title: '状态',
      width: 90,
      formatter: ({ cellValue }) => {
        if (cellValue === 1) return '已开';
        if (cellValue === -1) return '驳回';
        return '待审';
      },
    },
    { fixed: 'right', slots: { default: 'action' }, title: '操作', width: 140 },
  ],
  height: 'auto',
  pagerConfig: { enabled: true, pageSize: 20 },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const data = await fetchInvoices({
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'invoice_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

async function onAudit(row: InvoiceRow, status: number) {
  await auditInvoice(row.invoice_id, {
    status,
    mark: status === 1 ? '已开具' : '驳回',
  });
  ElMessage.success(status === 1 ? '已开票' : '已驳回');
  gridApi.reload();
}
</script>

<template>
  <Page auto-content-height title="发票管理">
    <Grid>
      <template #action="{ row }">
        <template v-if="row.status === 0">
          <ElButton link type="primary" @click="onAudit(row, 1)">开票</ElButton>
          <ElButton link type="danger" @click="onAudit(row, -1)">驳回</ElButton>
        </template>
      </template>
    </Grid>
  </Page>
</template>
