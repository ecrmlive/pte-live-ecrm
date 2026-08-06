<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElButton, ElMessage, ElMessageBox } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  deleteCommunityReplyApi,
  listAllCommunityRepliesApi,
  type CommunityReply,
} from '#/api/core/platform-community';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { LIST_KEYWORD_FIELD, listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canDelete = ref(false);

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('评论内容'),
]);

const gridOptions: VxeGridProps<CommunityReply> = {
  columns: [
    { field: 'reply_id', title: 'ID', width: 90 },
    {
      field: 'post_title',
      minWidth: 180,
      showOverflow: false,
      title: '帖子',
    },
    { field: 'community_id', title: '帖子 ID', width: 100 },
    { field: 'nickname', minWidth: 110, title: '用户' },
    {
      field: 'content',
      minWidth: 260,
      showOverflow: false,
      title: '内容',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '时间',
    },
    platformListActionColumn({ width: 80 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const result = await listAllCommunityRepliesApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
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

async function deleteReply(row: CommunityReply) {
  try {
    await ElMessageBox.confirm('确认删除该评论？删除后不可恢复。', '删除评论', {
      cancelButtonText: '取消',
      confirmButtonText: '确认删除',
      type: 'error',
    });
    await deleteCommunityReplyApi(row.reply_id);
    ElMessage.success('评论已删除');
    gridApi.reload();
  } catch {
    /* 用户取消或 requestClient 已统一提示接口异常。 */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const isContentOperator =
    profile.roles.includes('platform') || profile.roles.includes('operations');
  canDelete.value =
    isContentOperator && permissions.includes('content.community.delete');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #action="{ row }">
        <ElButton
          v-if="canDelete"
          link
          type="danger"
          @click="deleteReply(row)"
        >
          删除
        </ElButton>
        <span v-else>—</span>
      </template>
    </Grid>
    <ElAlert
      class="mt-4"
      title="有删帖权限的账号可删除违规评论。"
      type="info"
      :closable="false"
    />
  </Page>
</template>
