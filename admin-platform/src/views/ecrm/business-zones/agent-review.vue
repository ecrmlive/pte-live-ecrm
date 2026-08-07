<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage, ElMessageBox, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  auditBusinessZoneAgent,
  fetchBusinessZoneAgents,
  type BusinessZoneAgentRow,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { LIST_KEYWORD_FIELD, listFormOptionsDefaults } from '#/utils/list-form-defaults';

const statusText = (value: number) =>
  ({ '-1': '已驳回', '0': '待审核', '1': '已通过' })[String(value)] ||
  '未知';

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('姓名 / 手机号 / 商户名'),
  {
    component: 'Select',
    componentProps: {
      options: [
        { label: '待审核', value: 0 },
        { label: '已通过', value: 1 },
        { label: '已驳回', value: -1 },
      ],
      placeholder: '审核状态',
    },
    defaultValue: 0,
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<BusinessZoneAgentRow> = {
  columns: [
    { field: 'circle_agent_id', title: 'ID', width: 72 },
    { field: 'name', minWidth: 110, title: '代理人' },
    { field: 'phone', title: '手机号', width: 140 },
    { field: 'business_name', minWidth: 140, title: '关联商户' },
    {
      field: 'qualification',
      minWidth: 180,
      showOverflow: false,
      title: '资质说明',
    },
    {
      field: 'audit_reason',
      minWidth: 160,
      showOverflow: false,
      title: '审核说明',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const statusRaw = formValues?.status;
        const result = await fetchBusinessZoneAgents({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status:
            statusRaw === 0 || statusRaw === 1 || statusRaw === -1
              ? Number(statusRaw)
              : 0,
        });
        return { items: result.list, total: result.total };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'circle_agent_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function approve(row: BusinessZoneAgentRow) {
  try {
    await ElMessageBox.confirm(
      `确认通过“${row.name}”的代理申请？`,
      '代理审核',
    );
    await auditBusinessZoneAgent(row.circle_agent_id, 1);
    ElMessage.success('审核已通过');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

async function reject(row: BusinessZoneAgentRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      '请填写驳回原因',
      '驳回代理申请',
      { inputPattern: /.+/, inputErrorMessage: '驳回原因必填' },
    );
    await auditBusinessZoneAgent(row.circle_agent_id, -1, value.trim());
    ElMessage.success('已驳回');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #status="{ row }">
        <ElTag
          :type="
            row.status === 1
              ? 'success'
              : row.status === -1
                ? 'danger'
                : 'warning'
          "
        >
          {{ statusText(row.status) }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="row.status === 0">
          <ElButton link type="success" @click="approve(row)">通过</ElButton>
          <ElButton link type="danger" @click="reject(row)">驳回</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>
  </Page>
</template>
