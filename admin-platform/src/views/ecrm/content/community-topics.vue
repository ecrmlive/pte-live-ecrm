<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';
import { ElButton, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listCommunityTopicsApi,
  type CommunityTopic,
} from '#/api/core/platform-community';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const router = useRouter();

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('话题名称'),
  LIST_ENABLE_STATUS_FIELD('状态'),
]);

const gridOptions: VxeGridProps<CommunityTopic> = {
  columns: [
    { field: 'topic_id', title: 'ID', width: 90 },
    {
      field: 'topic_name',
      minWidth: 180,
      showOverflow: false,
      title: '话题名称',
    },
    { field: 'category_id', title: '分类 ID', width: 100 },
    { field: 'count_use', title: '引用次数', width: 100 },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'is_hot',
      slots: { default: 'is_hot' },
      title: '推荐',
      width: 90,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '创建时间',
    },
    platformListActionColumn({ width: 120 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const statusRaw = formValues?.status;
        let list = (await listCommunityTopicsApi()).list || [];
        if (keyword) {
          list = list.filter((row) =>
            row.topic_name.toLowerCase().includes(keyword),
          );
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter((row) => row.status === Number(statusRaw));
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'topic_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });

function openPosts(row: CommunityTopic) {
  void router.push({
    path: '/community/list',
    query: { topic_id: String(row.topic_id) },
  });
}

onMounted(() => {
  /* grid loads on mount via proxy */
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #is_hot="{ row }">
        <ElTag :type="row.is_hot === 1 ? 'warning' : 'info'">
          {{ row.is_hot === 1 ? '是' : '否' }}
        </ElTag>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openPosts(row)">查看帖子</ElButton>
      </template>
    </Grid>
  </Page>
</template>
