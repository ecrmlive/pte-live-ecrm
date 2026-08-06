<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage, ElMessageBox } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  fetchProfitsharingApplications,
  reviewProfitsharingApplication,
  saveProfitsharingApplicationNote,
  type ProfitsharingApplication,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const can = ref(false);

const statusLabels: Record<string, string> = {
  applied: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
};

function formatTime(value?: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function validateNote(value: string) {
  const note = value.trim();
  return note && [...note].length <= 500
    ? true
    : '审核说明不能为空，且不能超过 500 个字符。';
}

function isPromptDismissed(error: unknown) {
  return error === 'cancel' || error === 'close' || error === 'escape';
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const merIdRaw = String(formValues?.mer_id ?? '').trim();
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    merchant_id: merIdRaw ? Number(merIdRaw) : undefined,
    status: String(formValues?.status ?? '').trim() || undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '申请编号 / 说明关键词',
    },
    fieldName: 'keyword',
    label: '申请搜索',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '商户 ID' },
    fieldName: 'mer_id',
    label: '商户 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 'applied' },
        { label: '已通过', value: 'approved' },
        { label: '已拒绝', value: 'rejected' },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<ProfitsharingApplication> = {
  columns: [
    { field: 'application_no', minWidth: 160, title: '申请编号' },
    { field: 'merchant_id', title: '商户 ID', width: 100 },
    {
      field: 'description',
      minWidth: 260,
      showOverflow: false,
      title: '申请说明',
    },
    {
      field: 'status',
      formatter: ({ cellValue }) => statusLabels[cellValue] || cellValue,
      title: '状态',
      width: 100,
    },
    {
      field: 'review_note',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 180,
      showOverflow: false,
      title: '审核备注',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 170,
      title: '申请时间',
    },
    platformListActionColumn({ width: 210 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchProfitsharingApplications(
          buildListParams(page, formValues),
        );
        const list = data.list || [];
        return { items: list, total: Number(data.total ?? 0) };
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function review(row: ProfitsharingApplication, approved: boolean) {
  try {
    const { value } = await ElMessageBox.prompt(
      '填写审核说明。',
      approved ? '同意分账申请' : '拒绝分账申请',
      { inputValidator: validateNote },
    );
    await reviewProfitsharingApplication(row.id, approved, value.trim());
    ElMessage.success('审核已保存');
    gridApi.reload();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

async function note(row: ProfitsharingApplication) {
  try {
    const { value } = await ElMessageBox.prompt(
      '填写内部审核备注。',
      '分账申请备注',
      { inputValue: row.review_note, inputValidator: validateNote },
    );
    await saveProfitsharingApplicationNote(row.id, value.trim());
    ElMessage.success('备注已保存');
    gridApi.reload();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  can.value = codes.includes('merchant.profitsharing.review');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #action="{ row }">
        <template v-if="can">
          <template v-if="row.status === 'applied'">
            <ElButton link type="success" @click="review(row, true)">同意</ElButton>
            <ElButton link type="danger" @click="review(row, false)">拒绝</ElButton>
          </template>
          <ElButton link type="primary" @click="note(row)">备注</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>
  </Page>
</template>
