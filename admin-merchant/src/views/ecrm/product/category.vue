<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
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
  ElOption,
  ElSelect,
  ElSwitch,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import type {
  MerchantCategory,
  MerchantCategoryInput,
} from '#/api/core/merchant-category';
import {
  createMerchantCategoryApi,
  deleteMerchantCategoryApi,
  listMerchantCategoriesApi,
  updateMerchantCategoryApi,
} from '#/api/core/merchant-category';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const editID = ref<number>();
const allRows = ref<MerchantCategory[]>([]);
const form = reactive<MerchantCategoryInput>({
  parent_id: 0,
  name: '',
  sort: 0,
  status: 1,
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '分类名称' },
    fieldName: 'keyword',
    label: '分类搜索',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '已启用', value: 1 },
        { label: '已停用', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<MerchantCategory> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'category_id', title: 'ID', width: 80 },
    { field: 'name', minWidth: 220, showOverflow: false, title: '分类名称', treeNode: true },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '启用',
      width: 100,
    },
    merchantListActionColumn({ width: 130 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_params, formValues) => {
        const data = await listMerchantCategoriesApi({
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status:
            formValues?.status === 0 || formValues?.status === 1
              ? Number(formValues.status)
              : undefined,
        });
        allRows.value = data.list ?? [];
        return { items: allRows.value, total: allRows.value.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'category_id' },
  treeConfig: {
    expandAll: true,
    parentField: 'parent_id',
    rowField: 'category_id',
    transform: true,
  },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [EditDrawer, editDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: save,
});

function reset() {
  editID.value = undefined;
  Object.assign(form, { parent_id: 0, name: '', sort: 0, status: 1 });
}

function create() {
  reset();
  editDrawerApi.setState({ title: '新增商品分类' }).open();
}

function edit(row: MerchantCategory) {
  editID.value = row.category_id;
  Object.assign(form, {
    parent_id: row.parent_id,
    name: row.name,
    sort: row.sort,
    status: row.status,
  });
  editDrawerApi.setState({ title: '编辑商品分类' }).open();
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  saving.value = true;
  editDrawerApi.lock();
  try {
    if (editID.value) {
      await updateMerchantCategoryApi(editID.value, {
        ...form,
        name: form.name.trim(),
      });
    } else {
      await createMerchantCategoryApi({ ...form, name: form.name.trim() });
    }
    editDrawerApi.close();
    ElMessage.success('商品分类已保存');
    gridApi.reload();
  } finally {
    saving.value = false;
    editDrawerApi.unlock();
  }
}

async function changeStatus(row: MerchantCategory) {
  try {
    await updateMerchantCategoryApi(row.category_id, {
      parent_id: row.parent_id,
      name: row.name,
      sort: row.sort,
      status: row.status,
    });
    ElMessage.success('状态已更新');
  } catch {
    gridApi.reload();
  }
}

async function remove(row: MerchantCategory) {
  try {
    await confirm({
      content: `确定删除商品分类“${row.name}”吗？存在子分类或商品时将拒绝删除。`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteMerchantCategoryApi(row.category_id);
    ElMessage.success('商品分类已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" @click="create">新增分类</ElButton>
    </template>

    <Grid>
      <template #status="{ row }">
        <ElSwitch
          v-model="row.status"
          :active-value="1"
          :inactive-value="0"
          @change="changeStatus(row)"
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="edit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <EditDrawer class="w-[500px] max-w-[96vw]">
      <ElForm label-width="88px">
        <ElFormItem label="上级分类">
          <ElSelect v-model="form.parent_id" class="w-full">
            <ElOption :value="0" label="顶级分类" />
            <ElOption
              v-for="item in allRows"
              :key="item.category_id"
              :value="item.category_id"
              :disabled="item.category_id === editID"
              :label="item.name"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="分类名称" required>
          <ElInput v-model="form.name" maxlength="128" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" />
        </ElFormItem>
        <ElFormItem label="启用">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </EditDrawer>
  </Page>
</template>
