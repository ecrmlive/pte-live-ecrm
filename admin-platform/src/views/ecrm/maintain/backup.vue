<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElDescriptions, ElDescriptionsItem, ElMessage, ElTable, ElTableColumn, ElTabs, ElTabPane, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createDatabaseBackupApi,
  deleteDatabaseBackupApi,
  downloadDatabaseBackupApi,
  getDatabaseTableDetailApi,
  listDatabaseBackupsApi,
  listDatabaseTablesApi,
  maintainDatabaseTablesApi,
  type DatabaseBackupRecord,
  type DatabaseScope,
  type DatabaseTableColumn,
  type DatabaseTableRow,
} from '#/api/core/platform-database-backup';
import { platformListActionColumn, platformListPagerConfig } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

const READ_CODE = 'maintain.backup';
const MANAGE_CODE = 'maintain.backup.manage';

const activeTab = ref<'backups' | 'tables'>('tables');
const canManage = ref(false);
const canRead = ref(false);
const selectedTables = ref<DatabaseTableRow[]>([]);
const detailColumns = ref<DatabaseTableColumn[]>([]);
const detailLoading = ref(false);
const detailTable = ref<DatabaseTableRow>();
const operating = ref(false);

const selectedCount = computed(() => selectedTables.value.length);

const databaseGridOptions: VxeGridProps<DatabaseTableRow> = {
  border: true,
  checkboxConfig: { highlight: true, reserve: true, trigger: 'row' },
  columns: [
    { type: 'checkbox', width: 48 },
    { field: 'table_name', minWidth: 230, showOverflow: true, title: '表名称' },
    { field: 'table_comment', minWidth: 180, showOverflow: true, title: '备注' },
    { field: 'engine', title: '类型', width: 120 },
    { field: 'size_bytes', slots: { default: 'sizeBytes' }, title: '大小', width: 140 },
    {
      field: 'updated_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '更新时间',
    },
    { field: 'row_count', title: '行数', width: 120 },
    platformListActionColumn({ width: 100 }),
  ],
  emptyText: '暂无可备份的数据表',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const result = await listDatabaseTablesApi({ limit: page.pageSize, page: page.currentPage });
        selectedTables.value = [];
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'table_name' },
  toolbarConfig: { custom: false, export: false, refresh: false, zoom: false },
};

const backupGridOptions: VxeGridProps<DatabaseBackupRecord> = {
  columns: [
    { field: 'file_name', minWidth: 280, showOverflow: true, title: '备份文件' },
    { field: 'database_scope', slots: { default: 'scope' }, title: '数据库', width: 120 },
    { field: 'table_count', title: '表数量', width: 110 },
    { field: 'size_bytes', slots: { default: 'backupSize' }, title: '文件大小', width: 140 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '创建时间',
    },
    { field: 'status', slots: { default: 'status' }, title: '状态', width: 100 },
    platformListActionColumn({ width: 160 }),
  ],
  emptyText: '暂无数据库备份',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const result = await listDatabaseBackupsApi({ limit: page.pageSize, page: page.currentPage });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: { custom: false, export: false, refresh: false, zoom: false },
};

const [DatabaseGrid, databaseGridApi] = useVbenVxeGrid({
  gridEvents: {
    checkboxAll: syncSelectedTables,
    checkboxChange: syncSelectedTables,
  },
  gridOptions: databaseGridOptions,
});
const [BackupGrid, backupGridApi] = useVbenVxeGrid({ gridOptions: backupGridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  cancelText: '关闭',
  class: 'w-[min(860px,96vw)]',
  confirmText: '关闭',
  placement: 'right',
  onConfirm: () => detailDrawerApi.close(),
});

function syncSelectedTables() {
  selectedTables.value = (databaseGridApi.grid?.getCheckboxRecords?.() ?? []) as DatabaseTableRow[];
}

function formatBytes(value: number) {
  if (!value) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const unit = Math.min(Math.floor(Math.log(value) / Math.log(1024)), units.length - 1);
  return `${(value / 1024 ** unit).toFixed(unit === 0 ? 0 : 2)} ${units[unit]}`;
}

function groupSelectedTables() {
  return selectedTables.value.reduce<Record<DatabaseScope, string[]>>(
    (groups, row) => {
      groups[row.database_scope].push(row.table_name);
      return groups;
    },
    { admin: [], business: [] },
  );
}

async function runBackup() {
  if (!selectedCount.value) {
    ElMessage.warning('请先勾选需要备份的数据表');
    return;
  }
  try {
    await confirm({
      title: '确认备份',
      content: `将为选中的 ${selectedCount.value} 张数据表生成本地 SQL 备份文件，是否继续？`,
      icon: 'warning',
    });
  } catch {
    return;
  }
  operating.value = true;
  try {
    const groups = groupSelectedTables();
    const results = await Promise.all(
      (Object.entries(groups) as Array<[DatabaseScope, string[]]>).filter(([, tables]) => tables.length).map(([scope, tables]) =>
        createDatabaseBackupApi({ scope, tables }),
      ),
    );
    ElMessage.success(`已生成 ${results.length} 个本地备份文件`);
    selectedTables.value = [];
    databaseGridApi.grid?.clearCheckboxRow?.();
    backupGridApi.reload();
  } finally {
    operating.value = false;
  }
}

async function runMaintenance(action: 'optimize' | 'repair') {
  if (!selectedCount.value) {
    ElMessage.warning(`请先勾选需要${action === 'optimize' ? '优化' : '修复'}的数据表`);
    return;
  }
  const actionText = action === 'optimize' ? '优化表' : '修复表';
  try {
    await confirm({
      title: `确认${actionText}`,
      content: action === 'optimize'
        ? `将对选中的 ${selectedCount.value} 张表执行优化操作，期间可能占用数据库资源，是否继续？`
        : `将对选中的 ${selectedCount.value} 张表执行修复操作；InnoDB 表会执行安全的统计分析，是否继续？`,
      icon: 'warning',
    });
  } catch {
    return;
  }
  operating.value = true;
  try {
    const groups = groupSelectedTables();
    await Promise.all(
      (Object.entries(groups) as Array<[DatabaseScope, string[]]>).filter(([, tables]) => tables.length).map(([scope, tables]) =>
        maintainDatabaseTablesApi(action, { scope, tables }),
      ),
    );
    ElMessage.success(`${actionText}已完成`);
    selectedTables.value = [];
    databaseGridApi.grid?.clearCheckboxRow?.();
    databaseGridApi.reload();
  } finally {
    operating.value = false;
  }
}

async function openTableDetail(row: DatabaseTableRow) {
  detailTable.value = row;
  detailColumns.value = [];
  detailLoading.value = true;
  detailDrawerApi.setState({ title: `数据表详情 · ${row.table_name}` }).open();
  try {
    const result = await getDatabaseTableDetailApi(row.database_scope, row.table_name);
    detailColumns.value = result.columns || [];
  } catch {
    detailDrawerApi.close();
  } finally {
    detailLoading.value = false;
  }
}

async function downloadBackup(row: DatabaseBackupRecord) {
  try {
    const result = await downloadDatabaseBackupApi(row.id);
    const url = URL.createObjectURL(result.blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.fileName || row.file_name;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success('备份文件已下载');
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '下载备份失败');
  }
}

async function removeBackup(row: DatabaseBackupRecord) {
  try {
    await confirm({
      title: '删除备份',
      content: `确定删除本地备份文件“${row.file_name}”吗？删除后无法恢复。`,
      icon: 'warning',
    });
  } catch {
    return;
  }
  await deleteDatabaseBackupApi(row.id);
  ElMessage.success('备份文件已删除');
  backupGridApi.reload();
}

function onTabChange(name: string | number) {
  activeTab.value = name === 'backups' ? 'backups' : 'tables';
  if (activeTab.value === 'tables') databaseGridApi.reload();
  else backupGridApi.reload();
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  const platform = profile.roles.includes('platform');
  canRead.value = platform && (codes.includes(READ_CODE) || codes.includes(MANAGE_CODE));
  canManage.value = platform && codes.includes(MANAGE_CODE);
  if (canRead.value) databaseGridApi.reload();
});
</script>

<template>
  <Page>
    <ElTabs v-model="activeTab" class="database-backup-tabs" @tab-change="onTabChange">
      <ElTabPane label="数据库备份" name="tables">
        <DatabaseGrid>
          <template #toolbar-actions>
            <ElButton :disabled="!canManage || !selectedCount" :loading="operating" type="primary" @click="runBackup">
              备份
            </ElButton>
            <ElButton :disabled="!canManage || !selectedCount" :loading="operating" @click="runMaintenance('optimize')">
              优化表
            </ElButton>
            <ElButton :disabled="!canManage || !selectedCount" :loading="operating" @click="runMaintenance('repair')">
              修复表
            </ElButton>
          </template>
          <template #sizeBytes="{ row }">{{ formatBytes(row.size_bytes) }}</template>
          <template #action="{ row }">
            <ElButton link type="primary" @click="openTableDetail(row)">详情</ElButton>
          </template>
        </DatabaseGrid>
      </ElTabPane>

      <ElTabPane label="数据库列表备份" name="backups">
        <BackupGrid>
          <template #scope="{ row }">
            <ElTag size="small" type="info">
              {{ row.database_scope === 'admin' ? '管理库' : '业务库' }}
            </ElTag>
          </template>
          <template #backupSize="{ row }">{{ formatBytes(row.size_bytes) }}</template>
          <template #status="{ row }">
            <ElTag :type="row.status === 'ready' ? 'success' : 'info'" size="small">
              {{ row.status === 'ready' ? '可下载' : row.status }}
            </ElTag>
          </template>
          <template #action="{ row }">
            <ElButton link type="primary" @click="downloadBackup(row)">下载</ElButton>
            <ElButton v-if="canManage" link type="danger" @click="removeBackup(row)">删除</ElButton>
          </template>
        </BackupGrid>
      </ElTabPane>
    </ElTabs>

    <DetailDrawer>
      <ElDescriptions v-if="detailTable" :column="2" border class="mb-4">
        <ElDescriptionsItem label="表名称">{{ detailTable.table_name }}</ElDescriptionsItem>
        <ElDescriptionsItem label="数据库">{{ detailTable.database_scope === 'admin' ? '管理库' : '业务库' }}</ElDescriptionsItem>
        <ElDescriptionsItem label="存储引擎">{{ detailTable.engine }}</ElDescriptionsItem>
        <ElDescriptionsItem label="数据行数">{{ detailTable.row_count }}</ElDescriptionsItem>
      </ElDescriptions>
      <ElTable :data="detailColumns" :loading="detailLoading" border max-height="calc(100dvh - 280px)">
        <ElTableColumn min-width="150" prop="column_name" label="字段名称" />
        <ElTableColumn min-width="160" prop="column_type" label="字段类型" />
        <ElTableColumn width="90" prop="is_nullable" label="可为空" />
        <ElTableColumn min-width="140" prop="column_default" label="默认值" />
        <ElTableColumn min-width="180" prop="column_comment" label="说明" />
      </ElTable>
    </DetailDrawer>
  </Page>
</template>

<style scoped>
.database-backup-tabs :deep(.el-tabs__header) {
  margin-bottom: 16px;
}
</style>
