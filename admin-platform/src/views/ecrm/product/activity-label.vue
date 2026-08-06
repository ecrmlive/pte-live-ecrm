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
  ElMessage,
  ElMessageBox,
  ElSwitch,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  fetchActivityLabels,
  saveActivityLabels,
  type ProductCacheListItem,
} from '#/api/core/platform-product-cache';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const allRows = ref<ProductCacheListItem[]>([]);
const editingIndex = ref<number>();
const form = reactive<ProductCacheListItem>({
  id: '',
  name: '',
  enabled: true,
  remark: '',
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('名称 / 标识'),
  {
    ...LIST_ENABLE_STATUS_FIELD('启用状态'),
    componentProps: {
      clearable: true,
      options: [
        { label: '启用', value: 1 },
        { label: '停用', value: 0 },
      ],
      placeholder: '全部状态',
    },
  },
]);

const gridOptions: VxeGridProps<ProductCacheListItem & { _index: number }> = {
  columns: [
    { field: 'id', minWidth: 120, title: '标识' },
    { field: 'name', minWidth: 160, title: '名称' },
    {
      field: 'remark',
      minWidth: 220,
      showOverflow: false,
      title: '备注',
    },
    {
      field: 'enabled',
      slots: { default: 'enabled' },
      title: '启用',
      width: 90,
    },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const result = await fetchActivityLabels();
        allRows.value = result.list || [];
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const statusRaw = formValues?.status;
        let list = allRows.value.map((item, index) => ({
          ...item,
          _index: index,
        }));
        if (keyword) {
          list = list.filter(
            (row) =>
              row.name.toLowerCase().includes(keyword) ||
              row.id.toLowerCase().includes(keyword),
          );
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter(
            (row) => row.enabled === (Number(statusRaw) === 1),
          );
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

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => save(),
});

function resetForm() {
  editingIndex.value = undefined;
  Object.assign(form, { id: '', name: '', enabled: true, remark: '' });
}

function openCreate() {
  resetForm();
  formModalApi.setState({ title: '新增活动标签' }).open();
}

function openEdit(row: ProductCacheListItem & { _index: number }) {
  editingIndex.value = row._index;
  Object.assign(form, row);
  formModalApi.setState({ title: '编辑活动标签' }).open();
}

async function save() {
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请填写名称');
    return;
  }
  formModalApi.lock();
  try {
    const next = allRows.value.map((item) => ({ ...item }));
    const payload: ProductCacheListItem = {
      id: form.id.trim() || name,
      name,
      enabled: form.enabled,
      remark: form.remark.trim(),
    };
    if (editingIndex.value === undefined) next.push(payload);
    else next[editingIndex.value] = payload;
    const result = await saveActivityLabels(next);
    allRows.value = result.list || [];
    formModalApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formModalApi.unlock();
  }
}

async function remove(index: number, name: string) {
  try {
    await ElMessageBox.confirm(`确认删除“${name}”？`, '删除确认', {
      type: 'warning',
    });
    const next = allRows.value.filter((_, i) => i !== index);
    const result = await saveActivityLabels(next);
    allRows.value = result.list || [];
    ElMessage.success('已删除');
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
          新增标签
        </ElButton>
      </template>
      <template #enabled="{ row }">
        <ElTag :type="row.enabled ? 'success' : 'info'">
          {{ row.enabled ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row._index, row.name)">
          删除
        </ElButton>
      </template>
    </Grid>
    <FormModal>
      <ElForm label-width="84px">
        <ElFormItem label="名称" required>
          <ElInput v-model="form.name" maxlength="64" show-word-limit />
        </ElFormItem>
        <ElFormItem label="标识">
          <ElInput
            v-model="form.id"
            maxlength="64"
            placeholder="留空则使用名称"
          />
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput
            v-model="form.remark"
            maxlength="255"
            type="textarea"
            :rows="3"
          />
        </ElFormItem>
        <ElFormItem label="启用">
          <ElSwitch v-model="form.enabled" />
        </ElFormItem>
      </ElForm>
    </FormModal>
  </Page>
</template>
