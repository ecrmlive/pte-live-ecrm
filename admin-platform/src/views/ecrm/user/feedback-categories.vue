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
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createUserFeedbackCategory,
  deleteUserFeedbackCategory,
  fetchUserFeedbackCategories,
  setUserFeedbackCategoryStatus,
  updateUserFeedbackCategory,
  type UserFeedbackCategory,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const form = reactive({ id: 0, name: '', sort: 0, status: 1 as 0 | 1 });

function idempotencyKey(action: string, id = 0) {
  return `feedback-category-${action}-${id}-${crypto.randomUUID()}`;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '分类名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '启用', value: 1 },
        { label: '停用', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<UserFeedbackCategory> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    { field: 'name', minWidth: 180, showOverflow: false, title: '分类名称' },
    { field: 'sort', title: '排序', width: 100 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 110,
    },
    {
      field: 'updated_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      title: '更新时间',
      width: 180,
    },
    platformListActionColumn({ width: 220 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const keyword = String(formValues?.keyword ?? '').trim().toLowerCase();
        const statusRaw = formValues?.status;
        let list = (await fetchUserFeedbackCategories()).list || [];
        if (keyword) {
          list = list.filter((row) => row.name.toLowerCase().includes(keyword));
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter((row) => row.status === Number(statusRaw));
        }
        return { items: list, total: list.length };
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  Object.assign(form, { id: 0, name: '', sort: 0, status: 1 });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增反馈分类' }).open();
}

function openEdit(row: UserFeedbackCategory) {
  Object.assign(form, {
    id: row.id,
    name: row.name,
    sort: row.sort,
    status: row.status,
  });
  formDrawerApi.setState({ title: '编辑反馈分类' }).open();
}

async function save() {
  const name = form.name.trim();
  if (!name || [...name].length > 32) {
    ElMessage.warning('请填写不超过 32 字的分类名称');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const payload = {
      name,
      sort: Math.max(0, Math.min(9999, Number(form.sort) || 0)),
      status: form.status,
      idempotency_key: idempotencyKey(form.id ? 'update' : 'create', form.id),
    };
    if (form.id) await updateUserFeedbackCategory(form.id, payload);
    else await createUserFeedbackCategory(payload);
    formDrawerApi.close();
    ElMessage.success('已保存');
    resetForm();
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function switchStatus(row: UserFeedbackCategory) {
  await setUserFeedbackCategoryStatus(row.id, {
    status: row.status === 1 ? 0 : 1,
    idempotency_key: idempotencyKey('status', row.id),
  });
  ElMessage.success('状态已更新');
  gridApi.reload();
}

async function remove(row: UserFeedbackCategory) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.name}”仅从可选分类中移除，既有反馈仍保留原分类文本。`,
      '删除反馈分类',
      { type: 'warning' },
    );
    await deleteUserFeedbackCategory(row.id, {
      idempotency_key: idempotencyKey('delete', row.id),
    });
    ElMessage.success('已删除');
    if (form.id === row.id) resetForm();
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增分类
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status ? 'success' : 'info'">
          {{ row.status ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link @click="switchStatus(row)">
          {{ row.status ? '停用' : '启用' }}
        </ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="80px">
        <ElFormItem label="分类名称" required>
          <ElInput v-model="form.name" maxlength="32" show-word-limit />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" :max="9999" />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio :value="1">启用</ElRadio>
            <ElRadio :value="0">停用</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
