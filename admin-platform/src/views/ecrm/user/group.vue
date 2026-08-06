<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformUserGroupApi,
  deletePlatformUserGroupApi,
  listPlatformUserGroupsApi,
  updatePlatformUserGroupApi,
  type PlatformUserGroup,
} from '#/api/core/platform-user-group';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const editing = ref<PlatformUserGroup>();
const form = reactive({ group_name: '', sort: 0 });

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '分组名称关键词' },
    fieldName: 'keyword',
    label: '关键词',
  },
]);

const gridOptions: VxeGridProps<PlatformUserGroup> = {
  columns: [
    { field: 'group_id', title: 'ID', width: 88 },
    { field: 'group_name', minWidth: 220, showOverflow: false, title: '分组名称' },
    { field: 'sort', title: '排序', width: 96 },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '创建时间',
    },
    platformListActionColumn({ width: 130 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const keyword = String(formValues?.keyword ?? '').trim().toLowerCase();
        const result = await listPlatformUserGroupsApi({
          page: page.currentPage,
          limit: page.pageSize,
        });
        let list = result.list || [];
        if (keyword) {
          list = list.filter((row) =>
            row.group_name.toLowerCase().includes(keyword),
          );
        }
        return { items: list, total: keyword ? list.length : result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'group_id' },
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

function resetForm() {
  editing.value = undefined;
  Object.assign(form, { group_name: '', sort: 0 });
}

function openCreate() {
  resetForm();
  formModalApi.setState({ title: '新增用户分组' }).open();
}

function openEdit(row: PlatformUserGroup) {
  editing.value = row;
  Object.assign(form, { group_name: row.group_name, sort: row.sort });
  formModalApi.setState({ title: '编辑用户分组' }).open();
}

async function save() {
  if (!form.group_name.trim()) {
    ElMessage.warning('请填写分组名称');
    return;
  }
  formModalApi.lock();
  saving.value = true;
  try {
    const body = { group_name: form.group_name.trim(), sort: form.sort };
    if (editing.value) {
      await updatePlatformUserGroupApi(editing.value.group_id, body);
    } else {
      await createPlatformUserGroupApi(body);
    }
    formModalApi.close();
    ElMessage.success(editing.value ? '用户分组已更新' : '用户分组已创建');
    gridApi.reload();
  } finally {
    saving.value = false;
    formModalApi.unlock();
  }
}

async function remove(row: PlatformUserGroup) {
  try {
    await ElMessageBox.confirm(
      `删除用户分组“${row.group_name}”后不可恢复，是否继续？`,
      '删除确认',
      { type: 'warning' },
    );
    await deletePlatformUserGroupApi(row.group_id);
    ElMessage.success('用户分组已删除');
    gridApi.reload();
  } catch {
    /* 用户取消或统一请求层处理 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增分组
        </ElButton>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormModal>
      <ElForm label-width="80px">
        <ElFormItem label="分组名称" required>
          <ElInput v-model="form.group_name" maxlength="32" show-word-limit />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
        </ElFormItem>
      </ElForm>
    </FormModal>
  </Page>
</template>
