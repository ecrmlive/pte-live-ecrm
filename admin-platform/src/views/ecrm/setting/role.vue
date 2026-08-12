<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

import { confirm, Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElSwitch,
  ElTabPane,
  ElTabs,
  ElTag,
  ElTree,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformRole,
  deletePlatformRole,
  fetchPlatformMenuTree,
  fetchPlatformRoles,
  updatePlatformRole,
  type PlatformMenuNode,
  type PlatformRoleRow,
  type PlatformRoleType,
} from '#/api/core/ecrm';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

const ROLE_TABS: Array<{ label: string; value: PlatformRoleType }> = [
  { label: '平台', value: 'platform' },
  { label: '商户', value: 'merchant' },
  { label: '区域', value: 'region' },
];

const activeRoleType = ref<PlatformRoleType>('platform');
const menus = ref<PlatformMenuNode[]>([]);
const editing = ref<PlatformRoleRow>();
const menuTree = ref();
const form = reactive({
  role_name: '',
  status: 1,
});

const modalTitle = computed(() =>
  editing.value ? '编辑身份管理' : '新增身份管理',
);
const activeRoleTypeLabel = computed(
  () => ROLE_TABS.find((item) => item.value === activeRoleType.value)?.label || '平台',
);

function roleTypeLabel(value: PlatformRoleType) {
  return ROLE_TABS.find((item) => item.value === value)?.label || '平台';
}

function roleIDs(row?: PlatformRoleRow) {
  return row?.rules ? row.rules.split(',').map(Number).filter(Boolean) : [];
}

const gridOptions: VxeGridProps<PlatformRoleRow> = {
  columns: [
    { field: 'role_id', title: 'ID', width: 84 },
    { field: 'role_name', minWidth: 220, title: '身份名称' },
    {
      align: 'center',
      field: 'role_type',
      slots: { default: 'role_type' },
      title: '身份类型',
      width: 160,
    },
    {
      align: 'center',
      field: 'status',
      slots: { default: 'status' },
      title: '是否开启',
      width: 130,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 190,
      title: '创建时间',
    },
    {
      field: 'updated_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 190,
      title: '更新时间',
    },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const result = await fetchPlatformRoles({
          limit: page.pageSize,
          page: page.currentPage,
          role_type: activeRoleType.value,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'role_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: save,
});

async function ensureMenus() {
  menus.value = await fetchPlatformMenuTree(activeRoleType.value);
}

function onRoleTypeChange(value: string | number) {
  activeRoleType.value = value as PlatformRoleType;
  menus.value = [];
  gridApi.reload({ page: 1 });
}

async function openCreate() {
  editing.value = undefined;
  Object.assign(form, { role_name: '', status: 1 });
  await ensureMenus();
  formDrawerApi.setState({ title: '新增身份管理' }).open();
  requestAnimationFrame(() => menuTree.value?.setCheckedKeys([]));
}

async function openEdit(row: PlatformRoleRow) {
  editing.value = row;
  Object.assign(form, { role_name: row.role_name, status: row.status });
  await ensureMenus();
  formDrawerApi.setState({ title: '编辑身份管理' }).open();
  requestAnimationFrame(() => menuTree.value?.setCheckedKeys(roleIDs(row)));
}

async function save() {
  const roleName = form.role_name.trim();
  if (!roleName) {
    ElMessage.warning('请填写身份名称');
    return;
  }
  formDrawerApi.lock();
  try {
    const menu_ids = menuTree.value?.getCheckedKeys(false) || [];
    if (editing.value) {
      await updatePlatformRole(editing.value.role_id, {
        menu_ids,
        role_name: roleName,
        status: form.status,
      });
    } else {
      await createPlatformRole({
        menu_ids,
        role_name: roleName,
        role_type: activeRoleType.value,
        status: form.status,
      });
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    await gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function toggleStatus(row: PlatformRoleRow, enabled: boolean | string | number) {
  await updatePlatformRole(row.role_id, {
    menu_ids: roleIDs(row),
    role_name: row.role_name,
    status: enabled ? 1 : 0,
  });
  ElMessage.success('已保存');
  await gridApi.reload();
}

async function remove(row: PlatformRoleRow) {
  try {
    await confirm({
      content: `删除身份“${row.role_name}”后不可恢复。`,
      icon: 'warning',
      title: '删除身份管理',
    });
  } catch {
    return;
  }
  await deletePlatformRole(row.role_id);
  ElMessage.success('已删除');
  await gridApi.reload();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="role-toolbar">
          <ElTabs :model-value="activeRoleType" @tab-change="onRoleTypeChange">
            <ElTabPane
              v-for="item in ROLE_TABS"
              :key="item.value"
              :label="item.label"
              :name="item.value"
            />
          </ElTabs>
          <div class="role-toolbar__actions">
            <ElButton :icon="Plus" type="primary" @click="openCreate">
              新增身份管理
            </ElButton>
          </div>
        </div>
      </template>
      <template #role_type="{ row }">
        <ElTag effect="plain" type="info">
          {{ roleTypeLabel(row.role_type) }}
        </ElTag>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          @update:model-value="toggleStatus(row, $event)"
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer :title="modalTitle">
      <ElForm label-width="120px">
        <ElFormItem label="身份类型">
          <ElTag effect="plain" type="info">{{ activeRoleTypeLabel }}</ElTag>
        </ElFormItem>
        <ElFormItem label="身份名称" required>
          <ElInput v-model="form.role_name" maxlength="64" placeholder="请输入身份名称" />
        </ElFormItem>
        <ElFormItem label="是否开启">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem label="权限菜单">
          <div class="w-full rounded border p-3">
            <ElTree
              ref="menuTree"
              :data="menus"
              :default-expand-all="false"
              :props="{ label: 'menu_name', children: 'children' }"
              check-on-click-node
              node-key="menu_id"
              show-checkbox
            />
          </div>
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.role-toolbar {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 620px;
  width: 100%;
}

.role-toolbar__actions {
  display: flex;
  width: fit-content;
}

.role-toolbar :deep(.el-tabs__header) {
  margin: 0;
}

.role-toolbar :deep(.el-tabs__nav-wrap::after) {
  background-color: var(--el-border-color-lighter);
}
</style>
