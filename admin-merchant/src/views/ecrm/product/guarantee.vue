<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, confirm, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createProductGuaranteeApi,
  deleteProductGuaranteeApi,
  listProductGuaranteesApi,
  updateProductGuaranteeApi,
  type ProductGuarantee,
  type ProductGuaranteeInput,
} from '#/api/core/product-meta';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const editingID = ref<number>();
const form = reactive<ProductGuaranteeInput>({
  name: '',
  content: '',
  sort: 0,
  status: 1,
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '保障名称' },
    fieldName: 'keyword',
    label: '保障搜索',
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

const gridOptions: VxeGridProps<ProductGuarantee> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'guarantee_id', title: 'ID', width: 88 },
    { field: 'name', minWidth: 160, showOverflow: false, title: '保障名称' },
    { field: 'content', minWidth: 260, showOverflow: false, title: '服务说明' },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '启用',
      width: 100,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 128 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const status = formValues?.status;
        const data = await listProductGuaranteesApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status: status === 0 || status === 1 ? Number(status) : undefined,
        });
        return { items: data.list, total: data.total };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'guarantee_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [EditModal, editModalApi] = useVbenModal({
  onConfirm: save,
});

function resetForm() {
  editingID.value = undefined;
  form.name = '';
  form.content = '';
  form.sort = 0;
  form.status = 1;
}

function openCreate() {
  resetForm();
  editModalApi.setState({ title: '新增保障服务' }).open();
}

function openEdit(row: ProductGuarantee) {
  editingID.value = row.guarantee_id;
  form.name = row.name;
  form.content = row.content;
  form.sort = row.sort;
  form.status = row.status;
  editModalApi.setState({ title: '编辑保障服务' }).open();
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写保障名称');
    return;
  }
  saving.value = true;
  editModalApi.lock();
  try {
    if (editingID.value) {
      await updateProductGuaranteeApi(editingID.value, form);
    } else {
      await createProductGuaranteeApi(form);
    }
    editModalApi.close();
    ElMessage.success('保存成功');
    gridApi.reload();
  } finally {
    saving.value = false;
    editModalApi.unlock();
  }
}

async function remove(row: ProductGuarantee) {
  try {
    await confirm({
      content: `删除保障服务“${row.name}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteProductGuaranteeApi(row.guarantee_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

async function changeStatus(row: ProductGuarantee) {
  try {
    await updateProductGuaranteeApi(row.guarantee_id, {
      name: row.name,
      content: row.content,
      sort: row.sort,
      status: row.status,
    });
    ElMessage.success('状态已更新');
  } catch {
    gridApi.reload();
  }
}
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" @click="openCreate">新增保障</ElButton>
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
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <EditModal class="w-[560px] max-w-[96vw]">
      <ElForm label-width="88px">
        <ElFormItem label="保障名称" required>
          <ElInput v-model="form.name" maxlength="64" />
        </ElFormItem>
        <ElFormItem label="服务说明">
          <ElInput
            v-model="form.content"
            :rows="5"
            maxlength="1000"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" />
        </ElFormItem>
        <ElFormItem label="启用">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </EditModal>
  </Page>
</template>
