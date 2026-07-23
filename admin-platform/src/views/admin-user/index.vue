<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page } from '@vben/common-ui';

import { ElButton, ElTag } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import PlatformAdminUserApi from '#/api/core/platform-admin-user';
import type { PlatformAdminUserRow } from '#/api/core/platform-admin-user';
import { parseApiList } from '#/utils/list-response';

import AdminUserAddModal from './admin-user-add-modal.vue';
import AdminUserEditModal from './admin-user-edit-modal.vue';

const allRows = ref<PlatformAdminUserRow[]>([]);
const listLoaded = ref(false);
const roleOptions = ref<{ role_id: number; role_name: string }[]>([]);
const editOpen = ref(false);
const addOpen = ref(false);
const currentUser = ref<PlatformAdminUserRow | undefined>();

function filterRows(rows: PlatformAdminUserRow[], keyword: string) {
  const kw = keyword.trim().toLowerCase();
  if (!kw) {
    return rows;
  }
  return rows.filter((row) => {
    return (
      row.user_name?.toLowerCase().includes(kw) ||
      row.role_names?.some((name) => name.toLowerCase().includes(kw))
    );
  });
}

async function loadList() {
  const res = await PlatformAdminUserApi.userList(true);
  allRows.value = parseApiList<PlatformAdminUserRow>(res.data);
  roleOptions.value = res.data?.roleList ?? [];
}

const formOptions: VbenFormProps = {
  showCollapseButton: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '筛选用户名或角色',
      },
      fieldName: 'keyword',
      label: '用户名',
    },
  ],
};

const gridOptions: VxeGridProps<PlatformAdminUserRow> = {
  border: true,
  columns: [
    { field: 'admin_user_id', title: '用户ID', width: 90 },
    { field: 'user_name', minWidth: 140, title: '用户名' },
    {
      field: 'is_super',
      slots: { default: 'isSuper' },
      title: '超管',
      width: 80,
    },
    {
      field: 'role_names',
      minWidth: 220,
      slots: { default: 'roles' },
      title: '角色',
    },
    {
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 120,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 15, pageSizes: [10, 15, 20, 30, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!listLoaded.value) {
          try {
            await loadList();
          } catch {
            allRows.value = [];
            roleOptions.value = [];
          }
          listLoaded.value = true;
        }
        const filtered = filterRows(
          allRows.value,
          String(formValues?.keyword ?? ''),
        );
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: filtered.slice(start, start + page.pageSize),
          total: filtered.length,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'admin_user_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function reload() {
  listLoaded.value = false;
  try {
    await loadList();
  } catch {
    allRows.value = [];
    roleOptions.value = [];
  }
  listLoaded.value = true;
  gridApi.reload();
}

function openEdit(row: PlatformAdminUserRow) {
  currentUser.value = row;
  editOpen.value = true;
}

function openAdd() {
  addOpen.value = true;
}
</script>

<template>
  <Page>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-access:code="'platform:adminUser:add'"
          :icon="Plus"
          type="primary"
          @click="openAdd"
        >
          添加账号
        </ElButton>
      </template>

      <template #isSuper="{ row }">
        <ElTag :type="row.is_super ? 'success' : 'info'" size="small">
          {{ row.is_super ? '是' : '否' }}
        </ElTag>
      </template>

      <template #roles="{ row }">
        <ElTag
          v-for="name in row.role_names || []"
          :key="name"
          class="mr-1"
          size="small"
        >
          {{ name }}
        </ElTag>
        <span v-if="!row.role_names?.length">—</span>
      </template>

      <template #action="{ row }">
        <ElButton
          v-access:code="'platform:adminUser:edit'"
          link
          type="primary"
          @click="openEdit(row)"
        >
          分配角色
        </ElButton>
      </template>
    </Grid>

    <AdminUserAddModal
      v-model:open="addOpen"
      :role-options="roleOptions"
      @success="reload"
    />

    <AdminUserEditModal
      v-model:open="editOpen"
      :role-options="roleOptions"
      :user="currentUser"
      @success="reload"
    />
  </Page>
</template>
