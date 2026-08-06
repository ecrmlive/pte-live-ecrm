<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElDatePicker,
  ElForm,
  ElFormItem,
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
  createMerchantCombinationGroupApi,
  deleteMerchantCombinationGroupApi,
  listMerchantCombinationGroupsApi,
  setMerchantCombinationShowApi,
  updateMerchantCombinationGroupApi,
  type MerchantCombinationGroup,
  type MerchantCombinationSaveInput,
} from '#/api/core/merchant-combination';
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
  buying_count_num: 2,
  dates: [] as string[],
  is_show: 1,
  price: 0,
  product_id: 0,
  time: 24,
});

function formatTime(value: string) {
  return formatShanghaiDateTime(value) || '—';
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    buying_count_num: 2,
    dates: [],
    is_show: 1,
    price: 0,
    product_id: 0,
    time: 24,
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '商品名称 / 商品 ID',
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

const gridOptions: VxeGridProps<MerchantCombinationGroup> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'product_group_id', title: 'ID', width: 80 },
    {
      field: 'store_name',
      minWidth: 170,
      showOverflow: false,
      slots: { default: 'product' },
      title: '商品',
    },
    {
      field: 'price',
      title: '拼团价',
      width: 105,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    { field: 'buying_count_num', title: '成团人数', width: 100 },
    {
      field: 'time',
      title: '成团时限',
      width: 100,
      formatter: ({ cellValue }) => `${cellValue} 小时`,
    },
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
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const isShow = formValues?.is_show;
        const data = await listMerchantCombinationGroupsApi({
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
  rowConfig: { isHover: true, keyField: 'product_group_id' },
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
  onConfirm: async () => {
    if (
      !form.product_id ||
      form.price <= 0 ||
      form.buying_count_num < 2 ||
      form.dates.length !== 2
    ) {
      ElMessage.warning('请选择商品、填写拼团价、成团人数与活动时间');
      return;
    }
    const body: MerchantCombinationSaveInput = {
      buying_count_num: form.buying_count_num,
      end_time: form.dates[1]!,
      is_show: form.is_show,
      price: form.price,
      product_id: form.product_id,
      start_time: form.dates[0]!,
      time: form.time,
    };
    saving.value = true;
    editModalApi.lock();
    try {
      if (editingID.value) {
        await updateMerchantCombinationGroupApi(editingID.value, body);
      } else {
        await createMerchantCombinationGroupApi(body);
      }
      ElMessage.success(editingID.value ? '拼团活动已更新' : '拼团活动已创建');
      editModalApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      editModalApi.unlock();
    }
  },
});

function openCreate() {
  resetForm();
  editModalApi.setState({ title: '新建拼团活动' }).open();
}

function openEdit(row: MerchantCombinationGroup) {
  editingID.value = row.product_group_id;
  Object.assign(form, {
    buying_count_num: row.buying_count_num,
    dates: [formatTime(row.start_time), formatTime(row.end_time)],
    is_show: row.is_show,
    price: row.price,
    product_id: row.product_id,
    time: row.time || 24,
  });
  editModalApi.setState({ title: '编辑拼团活动' }).open();
}

async function toggleShow(row: MerchantCombinationGroup) {
  const isShow = row.is_show === 1 ? 0 : 1;
  await setMerchantCombinationShowApi(row.product_group_id, isShow);
  ElMessage.success(isShow ? '活动已上架' : '活动已下架');
  gridApi.reload();
}

async function remove(row: MerchantCombinationGroup) {
  try {
    await confirm({
      content: `删除商品 #${row.product_id} 的拼团活动？进行中团单不会被改写。`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteMerchantCombinationGroupApi(row.product_group_id);
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
      <ElButton type="primary" @click="openCreate">新建拼团活动</ElButton>
    </template>

    <Grid>
      <template #product="{ row }">
        {{ row.store_name || `商品 #${row.product_id}` }}
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

    <EditModal class="w-[640px] max-w-[96vw]">
      <ElForm class="grid grid-cols-2 gap-x-4" label-width="108px">
        <ElFormItem class="col-span-2" label="参与商品" required>
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
        <ElFormItem label="拼团价" required>
          <ElInputNumber
            v-model="form.price"
            :min="0.01"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="成团人数" required>
          <ElInputNumber
            v-model="form.buying_count_num"
            :min="2"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="成团时限" required>
          <ElInputNumber
            v-model="form.time"
            :max="720"
            :min="1"
            :precision="0"
            class="w-full"
          />
          <span class="ml-2 text-sm text-muted-foreground">小时</span>
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
    </EditModal>
  </Page>
</template>
