<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElForm, ElFormItem, ElInput, ElInputNumber, ElMessage, ElMessageBox, ElSwitch } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createExpress,
  deleteExpress,
  fetchExpressList,
  updateExpress,
  type ExpressRow,
} from '#/api/core/ecrm';

const form = reactive({
  name: '',
  code: '',
  sort: 0,
  show: true,
});
const editingId = ref(0);

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (!form.name.trim()) {
      ElMessage.warning('请填写名称');
      return;
    }
    const payload = {
      name: form.name.trim(),
      code: form.code.trim(),
      sort: form.sort,
      is_show: form.show ? 1 : 0,
    };
    if (editingId.value) {
      await updateExpress(editingId.value, payload);
    } else {
      await createExpress(payload);
    }
    ElMessage.success('已保存');
    formDrawerApi.close();
    gridApi.reload();
  },
});

const gridOptions: VxeGridProps<ExpressRow> = {
  border: true,
  columns: [
    { field: 'express_id', title: 'ID', width: 80 },
    { field: 'name', minWidth: 160, title: '名称' },
    { field: 'code', title: '编码', width: 140 },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'is_show',
      title: '状态',
      width: 90,
      formatter: ({ cellValue }) => (cellValue === 1 ? '展示' : '隐藏'),
    },
    { fixed: 'right', slots: { default: 'action' }, title: '操作', width: 160 },
  ],
  height: 'auto',
  pagerConfig: { enabled: true, pageSize: 10 },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const data = await fetchExpressList({
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'express_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function openCreate() {
  editingId.value = 0;
  form.name = '';
  form.code = '';
  form.sort = 0;
  form.show = true;
  formDrawerApi.setState({ title: '新建快递' });
  formDrawerApi.open();
}

function openEdit(row: ExpressRow) {
  editingId.value = row.express_id;
  form.name = row.name;
  form.code = row.code;
  form.sort = row.sort;
  form.show = row.is_show === 1;
  formDrawerApi.setState({ title: '编辑快递' });
  formDrawerApi.open();
}

async function onDelete(row: ExpressRow) {
  try {
    await ElMessageBox.confirm(`删除快递「${row.name}」？`, '提示', {
      type: 'warning',
    });
  } catch {
    return;
  }
  await deleteExpress(row.express_id);
  ElMessage.success('已删除');
  gridApi.reload();
}
</script>

<template>
  <Page auto-content-height title="物流公司">
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">新建快递</ElButton>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="onDelete(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-position="top">
        <ElFormItem label="名称" required>
          <ElInput v-model="form.name" />
        </ElFormItem>
        <ElFormItem label="编码">
          <ElInput v-model="form.code" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="展示">
          <ElSwitch v-model="form.show" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
