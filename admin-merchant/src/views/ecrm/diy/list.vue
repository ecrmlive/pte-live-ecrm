<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage, ElMessageBox, ElRadioButton, ElRadioGroup, ElTag } from 'element-plus';
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  activeDiyPageApi,
  applyDiyDefaultApi,
  copyDiyPageApi,
  deleteDiyPageApi,
  listDiyDefaultsApi,
  listDiyPagesApi,
  recoveryDiyPageApi,
  type DiyPageRow,
} from '#/api/core/diy';

const router = useRouter();
const route = useRoute();
const editorPath = '/devise/diy/index';
const pageKind = ref<'home' | 'micro'>(route.query.kind === 'micro' ? 'micro' : 'home');

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
        width: 380,
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
            is_diy: pageKind.value === 'home' ? 1 : 0,
          });
          return { items: res.list || [], total: res.total || 0 };
        },
      },
    },
    toolbarConfig: { search: false, refresh: true },
  } as VxeGridProps<DiyPageRow>,
});

function openEditor(id?: number) {
  router.push({
    path: editorPath,
    query: id
      ? { id: String(id), types: pageKind.value === 'home' ? '1' : '0' }
      : { types: pageKind.value === 'home' ? '1' : '0' },
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

async function onRecovery(row: DiyPageRow) {
  await ElMessageBox.confirm(
    `恢复会覆盖「${row.name}」当前的组件和页面设置，确认继续？`,
    '恢复默认页面',
    { cancelButtonText: '取消', confirmButtonText: '确认恢复', type: 'warning' },
  );
  await recoveryDiyPageApi(row.id);
  ElMessage.success('已恢复默认页面');
  gridApi.reload();
}

async function onDelete(row: DiyPageRow) {
  if (pageKind.value === 'home' && row.status === 1) {
    ElMessage.warning('请先取消启用再删除');
    return;
  }
  await ElMessageBox.confirm(`确认删除「${row.name}」？`, '提示');
  await deleteDiyPageApi(row.id);
  ElMessage.success('已删除');
  gridApi.reload();
}

function onKindChange() {
  gridApi.reload();
}

async function onApplyDefault() {
  const res = await listDiyDefaultsApi({ page: 1, limit: 1 });
  const first = res.list?.[0];
  if (!first) {
    ElMessage.warning('暂无平台默认模板');
    return;
  }
  await applyDiyDefaultApi(first.id);
  ElMessage.success('已套用平台模板');
  gridApi.reload();
}
</script>

<template>
  <Page :title="pageKind === 'home' ? '页面装修' : '微页面'" :description="pageKind === 'home' ? '商户首页 DIY 可视化装修' : '可被店铺首页组件链接引用的独立页面'">
    <template #extra>
      <ElRadioGroup v-model="pageKind" class="mr-3" @change="onKindChange"><ElRadioButton value="home">首页</ElRadioButton><ElRadioButton value="micro">微页面</ElRadioButton></ElRadioGroup>
      <ElButton v-if="pageKind === 'home'" @click="onApplyDefault">套用平台模板</ElButton>
      <ElButton type="primary" @click="openEditor()">{{ pageKind === 'home' ? '新建首页' : '新建微页面' }}</ElButton>
    </template>
    <Grid>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? (pageKind === 'home' ? '使用中' : '已发布') : '未启用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEditor(row.id)">装修</ElButton>
        <ElButton
          v-if="pageKind === 'home' && row.status !== 1"
          link
          type="success"
          @click="onActive(row)"
        >
          设为首页
        </ElButton>
        <ElButton link @click="onCopy(row)">复制</ElButton>
        <ElButton link type="warning" @click="onRecovery(row)">恢复</ElButton>
        <ElButton link type="danger" @click="onDelete(row)">删除</ElButton>
      </template>
    </Grid>
  </Page>
</template>
