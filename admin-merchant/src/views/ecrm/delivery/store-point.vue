<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page, confirm, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElSwitch,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createStorePickupPointApi,
  deleteStorePickupPointApi,
  listStorePickupPointsApi,
  updateStorePickupPointApi,
  type StorePickupPoint,
} from '#/api/core/merchant-pickup-point';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const saving = ref(false);
const editID = ref<number>();
const form = ref<Omit<StorePickupPoint, 'id'>>({
  contact_name: '',
  mobile: '',
  region_code: '',
  detail: '',
  is_default: 0,
});

function resetForm() {
  editID.value = undefined;
  form.value = {
    contact_name: '',
    mobile: '',
    region_code: '',
    detail: '',
    is_default: 0,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '联系人 / 手机号 / 地址',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
]);

const gridOptions: VxeGridProps<StorePickupPoint> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    { field: 'contact_name', title: '联系人', width: 120 },
    { field: 'mobile', title: '手机号', width: 140 },
    { field: 'detail', minWidth: 220, showOverflow: true, title: '地址' },
    {
      field: 'is_default',
      title: '默认',
      width: 80,
      formatter: ({ cellValue }) => (cellValue ? '是' : '否'),
    },
    merchantListActionColumn({ width: 130 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listStorePickupPointsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [PointModal, pointModalApi] = useVbenModal({
  onConfirm: async () => {
    saving.value = true;
    pointModalApi.lock();
    try {
      if (editID.value) {
        await updateStorePickupPointApi(editID.value, form.value);
      } else {
        await createStorePickupPointApi(form.value);
      }
      ElMessage.success('自提点已保存');
      pointModalApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      pointModalApi.unlock();
    }
  },
});

function openCreate() {
  resetForm();
  pointModalApi.setState({ title: '新增自提点' }).open();
}

function openEdit(row: StorePickupPoint) {
  editID.value = row.id;
  form.value = {
    contact_name: row.contact_name,
    mobile: row.mobile,
    region_code: row.region_code,
    detail: row.detail,
    is_default: row.is_default,
  };
  pointModalApi.setState({ title: '编辑自提点' }).open();
}

async function remove(row: StorePickupPoint) {
  try {
    await confirm({
      content: `确定删除自提点“${row.contact_name}”吗？`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteStorePickupPointApi(row.id);
    ElMessage.success('自提点已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" @click="openCreate">新增自提点</ElButton>
    </template>

    <Grid>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <PointModal class="w-[520px] max-w-[96vw]">
      <ElForm label-width="88px">
        <ElFormItem label="联系人" required>
          <ElInput v-model="form.contact_name" />
        </ElFormItem>
        <ElFormItem label="手机号" required>
          <ElInput v-model="form.mobile" />
        </ElFormItem>
        <ElFormItem label="区划代码">
          <ElInput v-model="form.region_code" />
        </ElFormItem>
        <ElFormItem label="详细地址" required>
          <ElInput v-model="form.detail" :rows="2" type="textarea" />
        </ElFormItem>
        <ElFormItem label="默认">
          <ElSwitch
            v-model="form.is_default"
            :active-value="1"
            :inactive-value="0"
          />
        </ElFormItem>
      </ElForm>
    </PointModal>
  </Page>
</template>
