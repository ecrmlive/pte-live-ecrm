<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElInput, ElMessage } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listProductUnitsApi,
  saveProductUnitsApi,
  type ProductUnit,
} from '#/api/core/merchant-product-unit';
import { MERCHANT_LIST_GRID_LAYOUT } from '#/constants/merchant-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const rows = ref<ProductUnit[]>([]);
const loaded = ref(false);

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '单位名称' },
    fieldName: 'keyword',
    label: '单位搜索',
  },
]);

const gridOptions: VxeGridProps<ProductUnit> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    {
      field: 'name',
      minWidth: 220,
      showOverflow: false,
      slots: { default: 'name' },
      title: '单位名称',
    },
    {
      field: 'sort',
      title: '排序',
      width: 120,
      slots: { default: 'sort' },
    },
    {
      field: 'action',
      fixed: 'right',
      showOverflow: false,
      slots: { default: 'action' },
      title: '操作',
      width: 100,
    },
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_params, formValues) => {
        if (!loaded.value) {
          const data = await listProductUnitsApi();
          rows.value = data.list ?? [];
          loaded.value = true;
        }
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        let list = rows.value;
        if (keyword) {
          list = list.filter((item) =>
            item.name.toLowerCase().includes(keyword),
          );
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'unit_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

function addRow() {
  rows.value.push({
    unit_id: rows.value.length + 1,
    name: '',
    sort: rows.value.length,
  });
  gridApi.reload();
}

function removeRow(index: number) {
  rows.value.splice(index, 1);
  gridApi.reload();
}

async function save() {
  if (rows.value.some((item) => !item.name.trim())) {
    ElMessage.warning('请填写全部单位名称');
    return;
  }
  saving.value = true;
  try {
    const result = await saveProductUnitsApi({
      list: rows.value.map((item, index) => ({
        ...item,
        name: item.name.trim(),
        sort: index,
      })),
    });
    rows.value = result.list ?? rows.value;
    loaded.value = true;
    ElMessage.success('商品单位已保存');
    gridApi.reload();
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton @click="addRow">新增单位</ElButton>
      <ElButton type="primary" :loading="saving" @click="save">保存</ElButton>
    </template>

    <Grid>
      <template #name="{ row }">
        <ElInput v-model="row.name" maxlength="16" />
      </template>
      <template #sort="{ rowIndex }">
        {{ rowIndex }}
      </template>
      <template #action="{ rowIndex }">
        <ElButton link type="danger" @click="removeRow(rowIndex)">删除</ElButton>
      </template>
    </Grid>
  </Page>
</template>
