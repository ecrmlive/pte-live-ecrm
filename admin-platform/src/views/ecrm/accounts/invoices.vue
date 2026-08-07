<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformInvoices,
  type PlatformInvoice,
} from '#/api/core/platform-invoice';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);

const labels: Record<PlatformInvoice['status'], string> = {
  issued: '已开票',
  rejected: '已拒绝',
  requested: '待开票',
  voided: '已作废',
};

function formatTime(value?: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const keyword = String(formValues?.keyword ?? '').trim();
  const orderNo = String(formValues?.order_no ?? '').trim();
  return {
    page: page.currentPage,
    limit: page.pageSize,
    status: (String(formValues?.status ?? '').trim() ||
      undefined) as PlatformInvoice['status'] | undefined,
    order_no: orderNo || keyword || undefined,
    keyword: keyword || undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(labels).map(([value, label]) => ({ label, value })),
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '开票状态',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      maxlength: 64,
      placeholder: '订单号 / 抬头 / 发票号',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, maxlength: 64, placeholder: '精确订单号' },
    fieldName: 'order_no',
    label: '订单号',
  },
]);

const gridOptions: VxeGridProps<PlatformInvoice> = {
  columns: [
    { field: 'id', title: '发票 ID', width: 110 },
    { field: 'order_no', minWidth: 210, showOverflow: false, title: '订单号' },
    {
      field: 'merchant_name',
      formatter: ({ row }) => `${row.merchant_name} / ${row.store_name}`,
      minWidth: 180,
      showOverflow: false,
      title: '店铺',
    },
    { field: 'title', minWidth: 200, showOverflow: false, title: '抬头' },
    { field: 'tax_no_masked', minWidth: 160, title: '税号（脱敏）' },
    { field: 'email_masked', minWidth: 180, title: '邮箱（脱敏）' },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 110,
    },
    {
      field: 'invoice_no',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 180,
      title: '发票号',
    },
    {
      field: 'requested_at',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 180,
      title: '申请时间',
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listPlatformInvoices(buildListParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canRead.value =
    profile.roles.includes('platform') &&
    permissions.includes('accounts.invoice.read');
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有发票监管权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #status="{ row }">
        <ElTag>{{ labels[row.status] }}</ElTag>
      </template>
    </Grid>
  </Page>
</template>
