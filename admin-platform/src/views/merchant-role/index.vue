<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { ElButton, ElMessage, ElMessageBox } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import MerchantRoleApi from '#/api/core/merchant-role';
import type { MerchantTemplateRoleRow } from '#/api/core/merchant-role';
import { Page } from '@vben/common-ui';
import { parseApiList } from '#/utils/list-response';

import RoleFormModal from './role-form-modal.vue';

const allRows = ref<MerchantTemplateRoleRow[]>([]);
const listLoaded = ref(false);
const formOpen = ref(false);
const formMode = ref<'add' | 'edit'>('add');
const currentRole = ref<MerchantTemplateRoleRow | undefined>();

function filterRows(rows: MerchantTemplateRoleRow[], keyword: string) {
  const kw = keyword.trim().toLowerCase();
  if (!kw) {
    return rows;
  }
  return rows.filter((row) => {
    return (
      row.role_name?.toLowerCase().includes(kw) ||
      row.remark?.toLowerCase().includes(kw)
    );
  });
}

async function loadList() {
  const res = await MerchantRoleApi.roleList(true);
  allRows.value = parseApiList<MerchantTemplateRoleRow>(res.data);
}

const formOptions: VbenFormProps = {
  showCollapseButton: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '筛选名称或备注',
      },
      fieldName: 'keyword',
      formItemClass: 'pb-0',
      label: '角色名称',
    },
  ],
};

const gridOptions = {
  border: true,
  columns: [
    { field: 'role_id', title: '角色ID', width: 90 },
    { field: 'role_name', minWidth: 160, title: '角色名称' },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'remark',
      minWidth: 200,
      showOverflow: true,
      title: '备注',
    },
    {
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 160,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 15, 20, 30, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!listLoaded.value) {
          try {
            await loadList();
          } catch {
            allRows.value = [];
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
  rowConfig: { isHover: true, keyField: 'role_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function reload() {
  listLoaded.value = false;
  try {
    await loadList();
  } catch {
    allRows.value = [];
  }
  listLoaded.value = true;
  gridApi.reload();
}

function openAdd() {
  formMode.value = 'add';
  currentRole.value = undefined;
  formOpen.value = true;
}

function openEdit(row: MerchantTemplateRoleRow) {
  formMode.value = 'edit';
  currentRole.value = row;
  formOpen.value = true;
}

async function handleDelete(row: MerchantTemplateRoleRow) {
  try {
    await ElMessageBox.confirm(`确定删除角色「${row.role_name}」？`, '提示', {
      type: 'warning',
    });
  } catch {
    return;
  }
  const res = await MerchantRoleApi.roleDelete(row.role_id, true);
  if (res.code === 1) {
    ElMessage.success(res.msg || '删除成功');
    await reload();
  }
}
</script>

<template>
  <Page>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-access:code="'platform:merchantRole:add'"
          :icon="Plus"
          type="primary"
          @click="openAdd"
        >
          新增角色
        </ElButton>
      </template>

      <template #action="{ row }">
        <ElButton
          v-access:code="'platform:merchantRole:edit'"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-if="row.role_id !== 1"
          v-access:code="'platform:merchantRole:delete'"
          link
          type="danger"
          @click="handleDelete(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <RoleFormModal
      v-model:open="formOpen"
      :mode="formMode"
      :role="currentRole"
      @success="reload"
    />
  </Page>
</template>
