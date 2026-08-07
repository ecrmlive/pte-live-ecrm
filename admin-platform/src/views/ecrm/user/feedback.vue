<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage, ElMessageBox, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  closeUserFeedback,
  deleteUserFeedback,
  fetchUserFeedback,
  replyUserFeedback,
  type UserFeedbackRow,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const statusLabels: Record<UserFeedbackRow['status'], string> = {
  pending: '待处理',
  replied: '已回复',
  closed: '已关闭',
};

const statusTypes: Record<
  UserFeedbackRow['status'],
  'info' | 'success' | 'warning'
> = {
  pending: 'warning',
  replied: 'success',
  closed: 'info',
};

function idempotencyKey(action: string, id: number) {
  return `feedback-${action}-${id}-${crypto.randomUUID()}`;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '反馈内容 / 用户 ID',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待处理', value: 'pending' },
        { label: '已回复', value: 'replied' },
        { label: '已关闭', value: 'closed' },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<UserFeedbackRow> = {
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    { field: 'user_id', title: '用户 ID', width: 100 },
    { field: 'type', title: '类型', width: 120 },
    {
      field: 'content',
      minWidth: 280,
      showOverflow: false,
      title: '反馈内容',
    },
    {
      field: 'reply',
      minWidth: 220,
      showOverflow: false,
      title: '平台回复',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      title: '提交时间',
      width: 170,
    },
    platformListActionColumn({ width: 180 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = String(formValues?.status ?? '').trim() || undefined;
        const result = await fetchUserFeedback({
          page: page.currentPage,
          limit: page.pageSize,
          status,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          date_from: range[0] as string | undefined,
          date_to: range[1] as string | undefined,
        });
        return {
          items: result.list || [],
          total: result.total || 0,
        };
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

async function reply(row: UserFeedbackRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      '回复会同步给反馈用户，最多 1000 字。',
      '回复用户反馈',
      { inputValidator: (v) => (v.trim() ? true : '请填写回复') },
    );
    await replyUserFeedback(row.id, {
      reply: value.trim(),
      idempotency_key: idempotencyKey('reply', row.id),
    });
    ElMessage.success('已回复');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

async function close(row: UserFeedbackRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      '可填写关闭说明。',
      '关闭反馈',
      { inputPlaceholder: '关闭说明（可选）' },
    );
    await closeUserFeedback(row.id, {
      reply: value.trim(),
      idempotency_key: idempotencyKey('close', row.id),
    });
    ElMessage.success('已关闭');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

async function remove(row: UserFeedbackRow) {
  try {
    await ElMessageBox.confirm(
      `将软删除反馈 #${row.id}，该操作保留审计记录。`,
      '删除反馈',
      { type: 'warning' },
    );
    await deleteUserFeedback(row.id, {
      idempotency_key: idempotencyKey('delete', row.id),
    });
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #status="{ row }">
        <ElTag :type="statusTypes[row.status]">
          {{ statusLabels[row.status] || row.status }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="row.status === 'pending'"
          link
          type="primary"
          @click="reply(row)"
        >
          回复
        </ElButton>
        <ElButton
          v-if="row.status !== 'closed'"
          link
          type="warning"
          @click="close(row)"
        >
          关闭
        </ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>
  </Page>
</template>
