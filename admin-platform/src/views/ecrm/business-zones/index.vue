<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElRadio,
  ElRadioGroup,
  ElSwitch,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createBusinessZone,
  deleteBusinessZone,
  fetchBusinessZones,
  updateBusinessZone,
  type BusinessZoneRow,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const editingID = ref<number>();
const form = reactive({
  pid: 0,
  name: '',
  circle_agent_id: 0,
  commission_type: 0,
  commission_rate: 0,
  remark: '',
  sort: 0,
  status: 1,
  type: 0,
  role_id: 0,
});

function statusText(value: number) {
  return value === 1 ? '启用' : '禁用';
}

function typeText(value: number) {
  return value === 1 ? '商户型商圈' : '区域商圈';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('商圈名称'),
  LIST_ENABLE_STATUS_FIELD('状态'),
]);

const gridOptions: VxeGridProps<BusinessZoneRow> = {
  columns: [
    { field: 'circle_id', title: 'ID', width: 72 },
    { field: 'name', minWidth: 150, title: '区域/商圈' },
    { field: 'path', minWidth: 150, title: '层级路径' },
    {
      field: 'type',
      formatter: ({ cellValue }) => typeText(Number(cellValue)),
      title: '类型',
      width: 120,
    },
    {
      field: 'commission_rate',
      formatter: ({ row }) =>
        row.commission_type === 1
          ? `${row.commission_rate}%（独立）`
          : '平台默认',
      title: '提成',
      width: 150,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const statusRaw = formValues?.status;
        const result = await fetchBusinessZones({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status:
            statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
        });
        return { items: result.list, total: result.total };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'circle_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  Object.assign(form, {
    pid: 0,
    name: '',
    circle_agent_id: 0,
    commission_type: 0,
    commission_rate: 0,
    remark: '',
    sort: 0,
    status: 1,
    type: 0,
    role_id: 0,
  });
}

function openCreate() {
  editingID.value = undefined;
  resetForm();
  formDrawerApi.setState({ title: '新建区域' }).open();
}

function openEdit(row: BusinessZoneRow) {
  editingID.value = row.circle_id;
  Object.assign(form, {
    pid: row.pid,
    name: row.name,
    circle_agent_id: row.circle_agent_id,
    commission_type: row.commission_type,
    commission_rate: row.commission_rate,
    remark: row.remark,
    sort: row.sort,
    status: row.status,
    type: row.type,
    role_id: row.role_id,
  });
  formDrawerApi.setState({ title: '编辑区域' }).open();
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写商圈名称');
    return;
  }
  formDrawerApi.lock();
  try {
    if (editingID.value) {
      await updateBusinessZone(editingID.value, form);
    } else {
      await createBusinessZone(form);
    }
    formDrawerApi.close();
    ElMessage.success('保存成功');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function remove(row: BusinessZoneRow) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.name}”后不可恢复，是否继续？`,
      '删除区域',
    );
    await deleteBusinessZone(row.circle_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新建区域
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ statusText(row.status) }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="110px">
        <ElFormItem label="上级区域ID">
          <ElInputNumber v-model="form.pid" :min="0" />
        </ElFormItem>
        <ElFormItem label="区域名称" required>
          <ElInput v-model="form.name" maxlength="64" />
        </ElFormItem>
        <ElFormItem label="类型">
          <ElRadioGroup v-model="form.type">
            <ElRadio :value="0">区域商圈</ElRadio>
            <ElRadio :value="1">商户型商圈</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="提成规则">
          <ElRadioGroup v-model="form.commission_type">
            <ElRadio :value="0">平台默认</ElRadio>
            <ElRadio :value="1">独立比例</ElRadio>
          </ElRadioGroup>
          <ElInputNumber
            v-if="form.commission_type === 1"
            v-model="form.commission_rate"
            class="ml-2"
            :max="100"
            :min="0"
            :precision="2"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem label="说明">
          <ElInput v-model="form.remark" type="textarea" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
