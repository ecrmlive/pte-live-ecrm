<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElSwitch,
  ElTag,
  ElTree,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformRole,
  fetchPlatformMenuTree,
  fetchPlatformRoles,
  updatePlatformRole,
  type PlatformMenuNode,
  type PlatformRoleRow,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const menus = ref<PlatformMenuNode[]>([]);
const editing = ref<PlatformRoleRow>();
const menuTree = ref();
const form = reactive({
  code: '',
  role_name: '',
  status: 1,
  is_agent: 0,
  circle_id: 0,
});

function roleType(value: number) {
  return value === 1 ? '区域角色' : '平台角色';
}

function roleIDs(row?: PlatformRoleRow) {
  return row?.rules ? row.rules.split(',').map(Number).filter(Boolean) : [];
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('角色名称 / 代码'),
  LIST_ENABLE_STATUS_FIELD('状态'),
]);

const gridOptions: VxeGridProps<PlatformRoleRow> = {
  columns: [
    { field: 'role_id', title: 'ID', width: 72 },
    { field: 'code', minWidth: 140, title: '角色代码' },
    { field: 'role_name', minWidth: 160, title: '角色名称' },
    {
      field: 'is_agent',
      slots: { default: 'is_agent' },
      title: '类型',
      width: 120,
    },
    {
      field: 'rules',
      formatter: ({ row }) => String(roleIDs(row).length),
      title: '权限数',
      width: 100,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ width: 90 }),
  ],
  pagerConfig: { enabled: true, pageSize: 50, pageSizes: [20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const statusRaw = formValues?.status;
        const result = await fetchPlatformRoles({
          page: page.currentPage,
          limit: page.pageSize,
        });
        let list = result.list || [];
        if (keyword) {
          list = list.filter(
            (row) =>
              row.code.toLowerCase().includes(keyword) ||
              row.role_name.toLowerCase().includes(keyword),
          );
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter((row) => row.status === Number(statusRaw));
        }
        return {
          items: list,
          total:
            keyword || statusRaw === 0 || statusRaw === 1
              ? list.length
              : result.total,
        };
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => save(),
});

const modalTitle = computed(() =>
  editing.value ? '编辑角色' : '新增角色',
);

async function ensureMenus() {
  if (!menus.value.length) menus.value = await fetchPlatformMenuTree();
}

async function openCreate() {
  editing.value = undefined;
  Object.assign(form, {
    code: '',
    role_name: '',
    status: 1,
    is_agent: 0,
    circle_id: 0,
  });
  await ensureMenus();
  formModalApi.setState({ title: '新增角色' }).open();
  requestAnimationFrame(() => menuTree.value?.setCheckedKeys([]));
}

async function openEdit(row: PlatformRoleRow) {
  editing.value = row;
  Object.assign(form, {
    code: row.code,
    role_name: row.role_name,
    status: row.status,
    is_agent: row.is_agent,
    circle_id: row.circle_id,
  });
  await ensureMenus();
  formModalApi.setState({ title: '编辑角色' }).open();
  requestAnimationFrame(() => menuTree.value?.setCheckedKeys(roleIDs(row)));
}

async function save() {
  if (
    !form.role_name.trim() ||
    (!editing.value && !/^[a-z0-9_]{1,32}$/.test(form.code))
  ) {
    ElMessage.warning('请填写角色名称和小写角色代码');
    return;
  }
  formModalApi.lock();
  try {
    const menu_ids = menuTree.value?.getCheckedKeys(false) || [];
    if (editing.value) {
      await updatePlatformRole(editing.value.role_id, {
        role_name: form.role_name.trim(),
        status: form.status,
        menu_ids,
      });
    } else {
      await createPlatformRole({
        code: form.code.trim(),
        role_name: form.role_name.trim(),
        status: form.status,
        menu_ids,
      });
    }
    formModalApi.close();
    ElMessage.success('角色权限已保存');
    gridApi.reload();
  } finally {
    formModalApi.unlock();
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增角色
        </ElButton>
      </template>
      <template #is_agent="{ row }">
        <ElTag :type="row.is_agent === 1 ? 'warning' : 'success'">
          {{ roleType(row.is_agent) }}
        </ElTag>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'danger'">
          {{ row.status === 1 ? '启用' : '禁用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
      </template>
    </Grid>

    <FormModal :title="modalTitle">
      <ElForm label-width="110px">
        <ElFormItem label="角色代码" required>
          <ElInput
            v-model="form.code"
            :disabled="!!editing"
            placeholder="如 operations、region_custom"
          />
        </ElFormItem>
        <ElFormItem label="角色名称" required>
          <ElInput v-model="form.role_name" />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem label="菜单与按钮">
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
    </FormModal>
  </Page>
</template>
