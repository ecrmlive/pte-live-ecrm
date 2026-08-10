<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  auditCommunityReplyApi,
  deleteCommunityReplyApi,
  listAllCommunityRepliesApi,
  type CommunityReply,
} from '#/api/core/platform-community';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);
const canManage = ref(false);
const auditTarget = ref<CommunityReply>();
const auditForm = reactive({
  status: 1 as -1 | 1,
  refusal: '',
});

function statusInfo(status: number) {
  if (status === 1) return { label: '已通过', type: 'success' as const };
  if (status === -1) return { label: '已拒绝', type: 'danger' as const };
  return { label: '待审核', type: 'warning' as const };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入用户名称',
    },
    fieldName: 'username',
    label: '用户名称',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入评论关键字',
    },
    fieldName: 'keyword',
    label: '关键字',
  },
]);

const gridOptions: VxeGridProps<CommunityReply> = {
  columns: [
    { field: 'reply_id', title: 'ID', width: 80 },
    {
      field: 'nickname',
      formatter: ({ row }) =>
        `${row.nickname || '用户'} | ${row.uid}`,
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '用户名/ID',
    },
    {
      field: 'post_title',
      minWidth: 160,
      showOverflow: 'tooltip',
      title: '文章标题',
    },
    {
      field: 'content',
      minWidth: 220,
      showOverflow: 'tooltip',
      title: '评论内容',
    },
    { field: 'count_reply', title: '评论条数', width: 90 },
    { field: 'count_start', title: '评论点赞数', width: 100 },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 170,
      title: '评论时间',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '审核状态',
      width: 100,
    },
    platformListActionColumn({ width: 120 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const result = await listAllCommunityRepliesApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          username: String(formValues?.username ?? '').trim() || undefined,
          date_from: range[0] as string | undefined,
          date_to: range[1] as string | undefined,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'reply_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [AuditDrawer, auditDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => submitAudit(),
});

function openAudit(row: CommunityReply) {
  auditTarget.value = row;
  auditForm.status = 1;
  auditForm.refusal = '';
  auditDrawerApi.setState({ title: '审核评论' }).open();
}

async function submitAudit() {
  if (!auditTarget.value) return false;
  if (auditForm.status === -1 && !auditForm.refusal.trim()) {
    ElMessage.warning('请填写拒绝理由');
    return false;
  }
  await auditCommunityReplyApi(auditTarget.value.reply_id, {
    status: auditForm.status,
    refusal: auditForm.refusal.trim() || undefined,
  });
  ElMessage.success('审核成功');
  auditDrawerApi.close();
  gridApi.reload();
  return true;
}

async function deleteReply(row: CommunityReply) {
  try {
    await confirm({
      title: '提示',
      content: '确认删除该评论？删除后不可恢复。',
      icon: 'warning',
    });
    await deleteCommunityReplyApi(row.reply_id);
    ElMessage.success('评论已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const roleOK = profile.roles.some(
    (role) => role === 'platform' || role === 'operations',
  );
  canRead.value =
    roleOK &&
    (codes.includes('content.community_reply.read') ||
      codes.includes('content.community_reply.manage') ||
      codes.includes('content.community_list.read') ||
      codes.includes('content.community_list.manage') ||
      codes.includes('content.community.delete'));
  canManage.value =
    roleOK &&
    (codes.includes('content.community_reply.manage') ||
      codes.includes('content.community_list.manage') ||
      codes.includes('content.community.delete'));
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status).type">
          {{ statusInfo(row.status).label }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="canManage && row.status === 0"
          link
          type="primary"
          @click="openAudit(row)"
        >
          审核
        </ElButton>
        <ElButton
          v-if="canManage"
          link
          type="danger"
          @click="deleteReply(row)"
        >
          删除
        </ElButton>
        <span v-if="!canManage">—</span>
      </template>
    </Grid>

    <AuditDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="审核状态" required>
          <ElRadioGroup v-model="auditForm.status">
            <ElRadio :label="1">通过</ElRadio>
            <ElRadio :label="-1">拒绝</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="auditForm.status === -1" label="拒绝理由" required>
          <ElInput
            v-model="auditForm.refusal"
            :rows="4"
            maxlength="200"
            placeholder="请填写拒绝理由"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </AuditDrawer>
  </Page>
</template>
