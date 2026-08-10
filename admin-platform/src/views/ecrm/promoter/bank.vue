<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createWithdrawBankApi,
  deleteWithdrawBankApi,
  listWithdrawBanksApi,
  setWithdrawBankStatusApi,
  updateWithdrawBankApi,
  type WithdrawBank,
} from '#/api/core/platform-spread';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

type DrawerMode = 'create' | 'edit';

const drawerMode = ref<DrawerMode>('create');
const editingId = ref(0);
const form = reactive({
  name: '',
  sort: 0,
  status: 1 as 0 | 1,
});

const gridOptions: VxeGridProps<WithdrawBank> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    {
      field: 'name',
      minWidth: 200,
      showOverflow: false,
      title: '银行名称',
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '添加时间',
    },
    {
      align: 'center',
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 120,
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const data = await listWithdrawBanksApi({
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: data.list || [], total: data.total || 0 };
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
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editingId.value = 0;
  form.name = '';
  form.sort = 0;
  form.status = 1;
}

function openCreate() {
  drawerMode.value = 'create';
  resetForm();
  formDrawerApi.setState({ title: '添加数据', confirmText: '确定' }).open();
}

function openEdit(row: WithdrawBank) {
  drawerMode.value = 'edit';
  editingId.value = row.id;
  form.name = row.name || '';
  form.sort = row.sort ?? 0;
  form.status = row.status === 1 ? 1 : 0;
  formDrawerApi.setState({ title: '编辑数据', confirmText: '确定' }).open();
}

async function save() {
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请输入银行名称');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      name,
      sort: form.sort ?? 0,
      status: form.status,
    };
    if (drawerMode.value === 'edit' && editingId.value) {
      await updateWithdrawBankApi(editingId.value, payload);
    } else {
      await createWithdrawBankApi(payload);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: WithdrawBank, enabled: boolean) {
  const before = row.status === 1;
  row.status = enabled ? 1 : 0;
  try {
    await setWithdrawBankStatusApi(row.id, enabled ? 1 : 0);
  } catch {
    row.status = before ? 1 : 0;
  }
}

async function onDelete(row: WithdrawBank) {
  try {
    await confirm({
      content: `确定删除银行「${row.name}」吗？`,
      icon: 'warning',
      title: '提示',
    });
  } catch {
    return;
  }
  await deleteWithdrawBankApi(row.id);
  ElMessage.success('已删除');
  gridApi.reload();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          添加数据
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeStatus(row, Boolean(enabled))
          "
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="onDelete(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="100px">
        <ElFormItem label="银行名称" required>
          <ElInput
            v-model="form.name"
            maxlength="64"
            placeholder="请输入银行名称"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" :precision="0" />
        </ElFormItem>
        <ElFormItem label="是否显示">
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="开启"
            inactive-text="关闭"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
