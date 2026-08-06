<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

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
  createMerchantCategory,
  deleteMerchantCategory,
  fetchMerchantCategories,
  updateMerchantCategory,
  type MerchantCategoryRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canManage = ref(false);
const editing = ref<MerchantCategoryRow>();
const form = reactive({ category_name: '', commission_rate: 0 });

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '分类名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
]);

const gridOptions: VxeGridProps<MerchantCategoryRow> = {
  columns: [
    { field: 'merchant_category_id', title: 'ID', width: 90 },
    {
      field: 'category_name',
      minWidth: 240,
      showOverflow: false,
      title: '商户分类',
    },
    {
      field: 'commission_rate',
      formatter: ({ cellValue }) => `${Number(cellValue || 0).toFixed(2)}%`,
      title: '平台佣金比例',
      width: 170,
    },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        const keyword = String(formValues?.keyword ?? '').trim().toLowerCase();
        let list = (await fetchMerchantCategories()).list || [];
        if (keyword) {
          list = list.filter((row) =>
            row.category_name.toLowerCase().includes(keyword),
          );
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'merchant_category_id' },
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
  Object.assign(form, { category_name: '', commission_rate: 0 });
}

function openCreate() {
  resetForm();
  formModalApi.setState({ title: '新增商户分类' }).open();
}

function openEdit(row: MerchantCategoryRow) {
  editing.value = row;
  Object.assign(form, {
    category_name: row.category_name,
    commission_rate: Number(row.commission_rate),
  });
  formModalApi.setState({ title: '编辑商户分类' }).open();
}

async function save() {
  if (
    !form.category_name.trim() ||
    form.commission_rate < 0 ||
    form.commission_rate > 100
  ) {
    ElMessage.warning('请填写分类名称，并将佣金比例限制在 0% 至 100%');
    return;
  }
  formModalApi.lock();
  try {
    const input = {
      category_name: form.category_name.trim(),
      commission_rate: form.commission_rate,
    };
    if (editing.value) {
      await updateMerchantCategory(
        editing.value.merchant_category_id,
        input,
      );
    } else {
      await createMerchantCategory(input);
    }
    formModalApi.close();
    ElMessage.success('商户分类已保存');
    gridApi.reload();
  } finally {
    formModalApi.unlock();
  }
}

async function remove(row: MerchantCategoryRow) {
  try {
    await ElMessageBox.confirm(
      `删除商户分类“${row.category_name}”后不可恢复，是否继续？`,
      '删除商户分类',
      { type: 'warning' },
    );
    await deleteMerchantCategory(row.merchant_category_id);
    ElMessage.success('商户分类已删除');
    gridApi.reload();
  } catch {
    /* 用户取消或统一请求层处理 */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('merchant.category.manage');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          新增分类
        </ElButton>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <FormModal>
      <ElForm label-width="108px">
        <ElFormItem label="分类名称" required>
          <ElInput v-model="form.category_name" maxlength="128" />
        </ElFormItem>
        <ElFormItem label="佣金比例" required>
          <ElInputNumber
            v-model="form.commission_rate"
            :max="100"
            :min="0"
            :precision="2"
            class="w-full"
          />
          <span class="ml-2">%</span>
        </ElFormItem>
      </ElForm>
    </FormModal>
  </Page>
</template>
