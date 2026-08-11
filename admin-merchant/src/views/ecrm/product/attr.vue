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
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createProductAttrTemplateApi,
  deleteProductAttrTemplateApi,
  listProductAttrTemplatesApi,
  updateProductAttrTemplateApi,
  type ProductAttrTemplate,
  type ProductAttrTemplateInput,
} from '#/api/core/product-meta';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const editingID = ref<number>();
const form = reactive<ProductAttrTemplateInput>({
  template_name: '',
  template_value: '[]',
  sort: 0,
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '模板名称' },
    fieldName: 'keyword',
    label: '模板搜索',
  },
]);

const gridOptions: VxeGridProps<ProductAttrTemplate> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'template_id', title: 'ID', width: 88 },
    { field: 'template_name', minWidth: 180, showOverflow: false, title: '模板名称' },
    {
      field: 'template_value',
      minWidth: 360,
      showOverflow: false,
      title: '参数定义',
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'create_time',
      minWidth: 170,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 128 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listProductAttrTemplatesApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
        });
        return { items: data.list, total: data.total };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'template_id' },
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

function resetForm() {
  editingID.value = undefined;
  form.template_name = '';
  form.template_value = '[]';
  form.sort = 0;
}

function openCreate() {
  resetForm();
  editDrawerApi.setState({ title: '新增规格模板' }).open();
}

function openEdit(row: ProductAttrTemplate) {
  editingID.value = row.template_id;
  form.template_name = row.template_name;
  form.template_value = row.template_value || '[]';
  form.sort = row.sort;
  editDrawerApi.setState({ title: '编辑规格模板' }).open();
}

function validateTemplate() {
  try {
    const parsed = JSON.parse(form.template_value || '[]');
    return Array.isArray(parsed) || typeof parsed === 'object';
  } catch {
    return false;
  }
}

async function save() {
  if (!form.template_name.trim()) {
    ElMessage.warning('请填写模板名称');
    return;
  }
  if (!validateTemplate()) {
    ElMessage.warning('参数定义必须是合法 JSON 对象或数组');
    return;
  }
  saving.value = true;
  editDrawerApi.lock();
  try {
    if (editingID.value) {
      await updateProductAttrTemplateApi(editingID.value, form);
    } else {
      await createProductAttrTemplateApi(form);
    }
    editDrawerApi.close();
    ElMessage.success('保存成功');
    gridApi.reload();
  } finally {
    saving.value = false;
    editDrawerApi.unlock();
  }
}

async function remove(row: ProductAttrTemplate) {
  try {
    await confirm({
      content: `删除规格模板“${row.template_name}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteProductAttrTemplateApi(row.template_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" @click="openCreate">新增规格</ElButton>
    </template>

    <Grid>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <EditDrawer class="w-[680px] max-w-[96vw]">
      <ElForm label-width="92px">
        <ElFormItem label="模板名称" required>
          <ElInput v-model="form.template_name" maxlength="64" />
        </ElFormItem>
        <ElFormItem label="参数定义" required>
          <ElInput v-model="form.template_value" :rows="10" type="textarea" />
          <div class="text-xs text-gray-500">
            使用 JSON 对象或数组保存规格定义，例如
            [{"name":"颜色","values":["黑","白"]}]。
          </div>
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" />
        </ElFormItem>
      </ElForm>
    </EditDrawer>
  </Page>
</template>
