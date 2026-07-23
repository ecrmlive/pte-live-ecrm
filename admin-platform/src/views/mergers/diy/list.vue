<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage, ElMessageBox, ElTag } from 'element-plus';
import { useRouter } from 'vue-router';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  activeDiyPageApi,
  copyDiyPageApi,
  deleteDiyPageApi,
  listDiyPagesApi,
  type DiyPageRow,
} from '#/api/core/diy';

const router = useRouter();
const editorPath = '/setting/diy/index';

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions: {
    columns: [
      { field: 'id', title: 'ID', width: 80 },
      { field: 'name', title: '页面名称', minWidth: 160 },
      { field: 'title', title: '标题', minWidth: 140 },
      {
        field: 'status',
        title: '状态',
        width: 100,
        slots: { default: 'status' },
      },
      { field: 'update_time', title: '更新时间', minWidth: 170 },
      {
        field: 'action',
        title: '操作',
        width: 320,
        fixed: 'right',
        slots: { default: 'action' },
      },
    ],
    height: 'auto',
    keepSource: true,
    proxyConfig: {
      ajax: {
        query: async ({ page }) => {
          const res = await listDiyPagesApi({
            page: page.currentPage,
            limit: page.pageSize,
            is_diy: 1,
          });
          return { items: res.list || [], total: res.total || 0 };
        },
      },
    },
    toolbarConfig: { search: true, refresh: true },
  } as VxeGridProps<DiyPageRow>,
});

function openEditor(id?: number) {
  router.push({
    path: editorPath,
    query: id ? { id: String(id), types: '1' } : { types: '1' },
  });
}

async function onActive(row: DiyPageRow) {
  await activeDiyPageApi(row.id);
  ElMessage.success('已设为首页');
  gridApi.reload();
}

async function onCopy(row: DiyPageRow) {
  await copyDiyPageApi(row.id);
  ElMessage.success('已复制');
  gridApi.reload();
}

async function onDelete(row: DiyPageRow) {
  if (row.status === 1) {
    ElMessage.warning('请先取消启用再删除');
    return;
  }
  await ElMessageBox.confirm(`确认删除「${row.name}」？`, '提示');
  await deleteDiyPageApi(row.id);
  ElMessage.success('已删除');
  gridApi.reload();
}
</script>

<template>
  <Page title="页面装修" description="平台首页 DIY 可视化装修">
    <template #extra>
      <ElButton type="primary" @click="openEditor()">新建首页</ElButton>
    </template>
    <Grid>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '使用中' : '未启用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEditor(row.id)">装修</ElButton>
        <ElButton
          v-if="row.status !== 1"
          link
          type="success"
          @click="onActive(row)"
        >
          设为首页
        </ElButton>
        <ElButton link @click="onCopy(row)">复制</ElButton>
        <ElButton link type="danger" @click="onDelete(row)">删除</ElButton>
      </template>
    </Grid>
  </Page>
</template>
