<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElInput, ElMessage, ElSwitch } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listStoreUserLabelsApi,
  saveStoreUserLabelsApi,
  type StoreUserLabel,
} from '#/api/core/merchant-user-label';
import { MERCHANT_LIST_GRID_LAYOUT } from '#/constants/merchant-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const rows = ref<StoreUserLabel[]>([]);
const loaded = ref(false);

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '标签名称' },
    fieldName: 'keyword',
    label: '标签搜索',
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

const gridOptions: VxeGridProps<StoreUserLabel> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    {
      field: 'name',
      minWidth: 200,
      showOverflow: false,
      slots: { default: 'name' },
      title: '标签名称',
    },
    {
      field: 'sort',
      title: '排序',
      width: 100,
      slots: { default: 'sort' },
    },
    {
      field: 'status',
      title: '启用',
      width: 100,
      slots: { default: 'status' },
    },
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_params, formValues) => {
        if (!loaded.value) {
          const data = await listStoreUserLabelsApi();
          rows.value = data.list ?? [];
          loaded.value = true;
        }
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const status = formValues?.status;
        let list = rows.value;
        if (keyword) {
          list = list.filter((item) =>
            item.name.toLowerCase().includes(keyword),
          );
        }
        if (status === 0 || status === 1) {
          list = list.filter((item) => item.status === status);
        }
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'label_id' },
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
    label_id: rows.value.length + 1,
    name: '',
    sort: rows.value.length,
    status: 1,
  });
  gridApi.reload();
}

async function save() {
  saving.value = true;
  try {
    const result = await saveStoreUserLabelsApi({
      list: rows.value.map((item, index) => ({
        ...item,
        name: item.name.trim(),
        sort: index,
      })),
    });
    rows.value = result.list ?? rows.value;
    loaded.value = true;
    ElMessage.success('用户标签已保存');
    gridApi.reload();
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton @click="addRow">新增标签</ElButton>
      <ElButton type="primary" :loading="saving" @click="save">保存</ElButton>
    </template>

    <Grid>
      <template #name="{ row }">
        <ElInput v-model="row.name" maxlength="32" />
      </template>
      <template #sort="{ rowIndex }">
        {{ rowIndex }}
      </template>
      <template #status="{ row }">
        <ElSwitch v-model="row.status" :active-value="1" :inactive-value="0" />
      </template>
    </Grid>
  </Page>
</template>
