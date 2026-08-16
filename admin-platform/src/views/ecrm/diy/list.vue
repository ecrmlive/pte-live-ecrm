<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';
import { ElButton, ElEmpty, ElMessage, ElMessageBox } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';
import { computed, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  activeDiyPageApi,
  copyDiyPageApi,
  deleteDiyPageApi,
  listDiyPagesApi,
  type DiyPageDoc,
  type DiyPageRow,
} from '#/api/core/diy';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

import DiyPreview from './Preview.vue';

const router = useRouter();
const route = useRoute();
const editorPath = '/setting/diy/index';
const SYSTEM_DEFAULT_HOME_PAGE_ID = 4001;
const pageKind = route.path.includes('/micro/') || route.query.kind === 'micro' ? 'micro' : 'home';
const canManage = ref(false);
const previewPage = ref<DiyPageRow>();
const previewForm = ref({ curItem: {}, selectedIndex: -1 });

const pageTitle = computed(() => (pageKind === 'home' ? '页面装修' : '微页面'));
const createLabel = computed(() => (pageKind === 'home' ? '新增页面' : '新增微页面'));
const previewDoc = computed<DiyPageDoc>(() =>
  previewPage.value?.doc || {
    page: { params: { name: '页面预览', title: '页面预览' } },
    items: [],
  },
);

function isSystemDefaultTemplate(row: DiyPageRow) {
  return pageKind === 'home' && row.id === SYSTEM_DEFAULT_HOME_PAGE_ID;
}

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions: {
    columns: [
      { field: 'id', title: '页面ID', width: 100 },
      { field: 'name', title: '模板名称', minWidth: 180 },
      {
        field: 'add_time',
        formatter: ({ cellValue }) => formatTime(cellValue),
        minWidth: 170,
        title: '添加时间',
      },
      {
        field: 'update_time',
        formatter: ({ cellValue }) => formatTime(cellValue),
        minWidth: 170,
        title: '更新时间',
      },
      platformListActionColumn({ width: 350 }),
    ],
    pagerConfig: platformListPagerConfig(),
    proxyConfig: {
      ajax: {
        query: async ({ page }) => {
          const res = await listDiyPagesApi({
            page: page.currentPage,
            limit: page.pageSize,
            is_diy: pageKind === 'home' ? 1 : 0,
          });
          const rows = res.list || [];
          if (!previewPage.value || !rows.some((row) => row.id === previewPage.value?.id)) {
            previewPage.value =
              rows.find((row) => isSystemDefaultTemplate(row)) ||
              rows.find((row) => row.status === 1) ||
              rows[0];
          }
          return { items: rows, total: res.total || 0 };
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
  } as VxeGridProps<DiyPageRow>,
});

function formatTime(value?: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function openEditor(id?: number) {
  router.push({
    path: editorPath,
    query: id ? { id: String(id), types: pageKind === 'home' ? '1' : '0' } : { types: pageKind === 'home' ? '1' : '0' },
  });
}

function selectPreview(row: DiyPageRow) {
  previewForm.value = { curItem: {}, selectedIndex: -1 };
  previewPage.value = row;
}

async function onActive(row: DiyPageRow) {
  if (isSystemDefaultTemplate(row)) {
    ElMessage.warning('系统默认模板仅支持预览和复制');
    return;
  }
  await activeDiyPageApi(row.id);
  ElMessage.success('已设为首页');
  await gridApi.reload();
}

async function onCopy(row: DiyPageRow) {
  await copyDiyPageApi(row.id);
  ElMessage.success('已复制');
  await gridApi.reload();
}

async function onDelete(row: DiyPageRow) {
  if (isSystemDefaultTemplate(row)) {
    ElMessage.warning('系统默认模板仅支持预览和复制');
    return;
  }
  if (pageKind === 'home' && row.status === 1) {
    ElMessage.warning('请先设置其他首页模板，再删除当前首页');
    return;
  }
  await ElMessageBox.confirm(`确认删除「${row.name}」？`, '删除页面', {
    cancelButtonText: '取消',
    confirmButtonText: '删除',
    type: 'warning',
  });
  await deleteDiyPageApi(row.id);
  if (previewPage.value?.id === row.id) previewPage.value = undefined;
  ElMessage.success('已删除');
  await gridApi.reload();
}

getAccessCodesApi().then((permissions) => {
  canManage.value = permissions.includes('operations.diy.manage');
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <section class="diy-page-list" :aria-label="pageTitle">
      <aside class="diy-page-list__preview" aria-label="移动端预览">
        <div v-if="previewPage" class="diy-page-list__phone">
          <DiyPreview
            :compact="true"
            :default-data="{}"
            :diy-data="previewDoc"
            diy-type=""
            :form="previewForm"
          />
        </div>
        <ElEmpty v-else description="暂无页面模板" />
      </aside>

      <section class="diy-page-list__content">
        <header class="diy-page-list__toolbar">
          <ElButton v-if="canManage" :icon="Plus" type="primary" @click="openEditor()">
            {{ createLabel }}
          </ElButton>
          <p v-if="pageKind === 'home'" class="diy-page-list__tip">
            系统默认模板仅支持预览和复制；复制后可编辑并设为首页。
          </p>
        </header>

        <Grid>
          <template #action="{ row }">
            <ElButton v-if="canManage && !isSystemDefaultTemplate(row)" link type="primary" @click="openEditor(row.id)">编辑</ElButton>
            <ElButton v-if="canManage && !isSystemDefaultTemplate(row)" link type="danger" @click="onDelete(row)">删除</ElButton>
            <ElButton
              v-if="canManage && pageKind === 'home' && row.status !== 1 && !isSystemDefaultTemplate(row)"
              link
              type="primary"
              @click="onActive(row)"
            >
              设为首页
            </ElButton>
            <ElButton link type="primary" @click="selectPreview(row)">预览</ElButton>
            <ElButton v-if="canManage" link type="primary" @click="onCopy(row)">复制</ElButton>
          </template>
        </Grid>
      </section>
    </section>
  </Page>
</template>

<style scoped>
.diy-page-list {
  display: grid;
  grid-template-columns: minmax(380px, 460px) minmax(0, 1fr);
  min-height: calc(100vh - 148px);
  overflow: hidden;
  background: var(--vben-bg-color-overlay, #fff);
}

.diy-page-list__preview {
  display: flex;
  min-height: 0;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  overflow: auto;
  border-right: 1px solid var(--el-border-color-lighter);
  background: #f7f8fa;
  padding: 24px 20px 32px;
}

.diy-page-list__phone {
  display: flex;
  width: 100%;
  justify-content: center;
}

.diy-page-list__content {
  min-width: 0;
  padding: 24px;
}

.diy-page-list__toolbar {
  display: flex;
  min-height: 40px;
  align-items: center;
  gap: 14px;
  margin-bottom: 18px;
}

.diy-page-list__tip {
  margin: 0;
  color: var(--el-color-danger);
  font-size: 14px;
  line-height: 22px;
}

@media (max-width: 1180px) {
  .diy-page-list {
    grid-template-columns: 1fr;
    overflow: visible;
  }

  .diy-page-list__preview {
    max-height: 620px;
    border-right: 0;
    border-bottom: 1px solid var(--el-border-color-lighter);
  }
}

@media (max-width: 720px) {
  .diy-page-list__content {
    padding: 16px;
  }

  .diy-page-list__toolbar {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
