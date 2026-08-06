<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';
import { ElButton, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listCommunityCategoriesApi,
  type CommunityCategory,
} from '#/api/core/platform-community';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { LIST_KEYWORD_FIELD, listFormOptionsDefaults } from '#/utils/list-form-defaults';

const router = useRouter();

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('分类名称'),
]);

const gridOptions: VxeGridProps<CommunityCategory> = {
  columns: [
    { field: 'category_id', title: 'ID', width: 90 },
    {
      field: 'cate_name',
      minWidth: 160,
      showOverflow: false,
      title: '分类名称',
    },
    { field: 'pid', title: '上级 ID', width: 100 },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '展示',
      width: 90,
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
        let list = (await listCommunityCategoriesApi()).list || [];
        if (keyword) {
          list = list.filter((row) =>
            row.cate_name.toLowerCase().includes(keyword),
          );
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'category_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });

function openPosts(row: CommunityCategory) {
  void router.push({
    path: '/community/list',
    query: { category_id: String(row.category_id) },
  });
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #is_show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '显示' : '隐藏' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openPosts(row)">查看帖子</ElButton>
      </template>
    </Grid>
  </Page>
</template>
