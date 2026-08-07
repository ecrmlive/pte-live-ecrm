<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listMerchantProductsApi,
  type MerchantProduct,
} from '#/api/core/merchant-catalog';
import {
  createMerchantAssistActiveApi,
  deleteMerchantAssistActiveApi,
  listMerchantAssistActivesApi,
  setMerchantAssistShowApi,
  updateMerchantAssistActiveApi,
  type MerchantAssistActive,
  type MerchantAssistSaveInput,
} from '#/api/core/merchant-assist';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const editingID = ref<number>();
const products = ref<MerchantProduct[]>([]);
const form = reactive({
  assist_count: 2,
  assist_price: 0,
  assist_user_count: 1,
  dates: [] as string[],
  is_show: 1,
  product_id: 0,
  stock: 100,
  store_info: '',
  store_name: '',
});

function formatTime(value: string) {
  return formatShanghaiDateTime(value) || '—';
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    assist_count: 2,
    assist_price: 0,
    assist_user_count: 1,
    dates: [],
    is_show: 1,
    product_id: 0,
    stock: 100,
    store_info: '',
    store_name: '',
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '活动名称 / 商品 ID',
    },
    fieldName: 'keyword',
    label: '活动搜索',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '已上架', value: 1 },
        { label: '已下架', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'is_show',
    label: '上架状态',
  },
]);

const gridOptions: VxeGridProps<MerchantAssistActive> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'product_assist_id', title: 'ID', width: 80 },
    { field: 'store_name', minWidth: 170, showOverflow: false, title: '活动 / 商品' },
    { field: 'product_id', title: '商品 ID', width: 92 },
    {
      field: 'assist_price',
      title: '助力价',
      width: 105,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'assist_count',
      minWidth: 150,
      showOverflow: false,
      slots: { default: 'rule' },
      title: '助力规则',
    },
    { field: 'stock', title: '库存', width: 80 },
    {
      field: 'start_time',
      minWidth: 220,
      showOverflow: false,
      slots: { default: 'dates' },
      title: '活动时间',
    },
    {
      field: 'is_show',
      slots: { default: 'show' },
      title: '展示',
      width: 88,
    },
    merchantListActionColumn({ width: 190 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const isShow = formValues?.is_show;
        const data = await listMerchantAssistActivesApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          is_show: isShow === 0 || isShow === 1 ? Number(isShow) : undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_assist_id' },
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
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (!form.product_id || form.assist_price <= 0 || form.dates.length !== 2) {
      ElMessage.warning('请选择商品、填写助力价并设置活动时间');
      return;
    }
    const body: MerchantAssistSaveInput = {
      assist_count: form.assist_count,
      assist_price: form.assist_price,
      assist_user_count: form.assist_user_count,
      end_time: form.dates[1]!,
      is_show: form.is_show,
      product_id: form.product_id,
      start_time: form.dates[0]!,
      stock: form.stock,
      store_info: form.store_info.trim(),
      store_name: form.store_name.trim(),
    };
    saving.value = true;
    editDrawerApi.lock();
    try {
      if (editingID.value) {
        await updateMerchantAssistActiveApi(editingID.value, body);
      } else {
        await createMerchantAssistActiveApi(body);
      }
      ElMessage.success(editingID.value ? '助力活动已更新' : '助力活动已创建');
      editDrawerApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      editDrawerApi.unlock();
    }
  },
});

function openCreate() {
  resetForm();
  editDrawerApi.setState({ title: '新建助力活动' }).open();
}

function openEdit(row: MerchantAssistActive) {
  editingID.value = row.product_assist_id;
  Object.assign(form, {
    assist_count: row.assist_count,
    assist_price: row.assist_price,
    assist_user_count: row.assist_user_count,
    dates: [formatTime(row.start_time), formatTime(row.end_time)],
    is_show: row.is_show,
    product_id: row.product_id,
    stock: row.stock,
    store_info: row.store_info || '',
    store_name: row.store_name,
  });
  editDrawerApi.setState({ title: '编辑助力活动' }).open();
}

async function toggleShow(row: MerchantAssistActive) {
  const isShow = row.is_show === 1 ? 0 : 1;
  await setMerchantAssistShowApi(row.product_assist_id, isShow);
  ElMessage.success(isShow ? '活动已上架' : '活动已下架');
  gridApi.reload();
}

async function remove(row: MerchantAssistActive) {
  try {
    await confirm({
      content: `删除助力活动「${row.store_name}」？删除后不可恢复。`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteMerchantAssistActiveApi(row.product_assist_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

onMounted(async () => {
  const productPage = await listMerchantProductsApi({
    limit: 100,
    page: 1,
    status: 1,
  });
  products.value = productPage.list || [];
});
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" @click="openCreate">新建助力活动</ElButton>
    </template>

    <Grid>
      <template #rule="{ row }">
        {{ row.assist_count }} 人 / 每人 {{ row.assist_user_count }} 次
      </template>
      <template #dates="{ row }">
        {{ formatTime(row.start_time) }} 至 {{ formatTime(row.end_time) }}
      </template>
      <template #show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '已上架' : '已下架' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="warning" @click="toggleShow(row)">
          {{ row.is_show === 1 ? '下架' : '上架' }}
        </ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <EditDrawer class="w-[640px] max-w-[96vw]">
      <ElForm class="grid grid-cols-2 gap-x-4" label-width="96px">
        <ElFormItem class="col-span-2" label="活动名称">
          <ElInput
            v-model="form.store_name"
            maxlength="60"
            placeholder="留空则使用商品名"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="活动简介">
          <ElInput
            v-model="form.store_info"
            :rows="2"
            maxlength="200"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem label="参与商品" required>
          <ElSelect
            v-model="form.product_id"
            :disabled="!!editingID"
            filterable
            class="w-full"
            placeholder="选择本店已审核商品"
          >
            <ElOption
              v-for="item in products"
              :key="item.product_id"
              :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`"
              :value="item.product_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="助力价" required>
          <ElInputNumber
            v-model="form.assist_price"
            :min="0.01"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="所需人数" required>
          <ElInputNumber
            v-model="form.assist_count"
            :min="1"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="单人次数" required>
          <ElInputNumber
            v-model="form.assist_user_count"
            :min="1"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="活动库存" required>
          <ElInputNumber
            v-model="form.stock"
            :min="1"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="前台展示">
          <ElSwitch v-model="form.is_show" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="活动时间" required>
          <ElDatePicker
            v-model="form.dates"
            class="w-full"
            end-placeholder="结束时间"
            start-placeholder="开始时间"
            type="datetimerange"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </ElFormItem>
      </ElForm>
    </EditDrawer>
  </Page>
</template>
