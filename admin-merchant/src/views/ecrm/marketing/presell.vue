<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref, watch } from 'vue';

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
  ElRadio,
  ElRadioGroup,
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
  createMerchantPresellActiveApi,
  deleteMerchantPresellActiveApi,
  listMerchantPresellActivesApi,
  setMerchantPresellShowApi,
  updateMerchantPresellActiveApi,
  type MerchantPresellActive,
  type MerchantPresellSaveInput,
} from '#/api/core/merchant-presell';
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
  activity_dates: [] as string[],
  down_price: 0,
  final_dates: [] as string[],
  final_price: 0,
  is_show: 1,
  presell_type: 1,
  price: 0,
  product_id: 0,
  stock: 100,
  store_info: '',
  store_name: '',
});

function dateText(value: unknown) {
  return formatShanghaiDateTime(value);
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    activity_dates: [],
    down_price: 0,
    final_dates: [],
    final_price: 0,
    is_show: 1,
    presell_type: 1,
    price: 0,
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
    componentProps: { clearable: true, placeholder: '活动名称' },
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

const gridOptions: VxeGridProps<MerchantPresellActive> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'product_presell_id', title: 'ID', width: 80 },
    { field: 'store_name', minWidth: 180, showOverflow: false, title: '预售名称' },
    {
      field: 'presell_type',
      slots: { default: 'type' },
      title: '类型',
      width: 100,
    },
    {
      field: 'price',
      title: '预售价',
      width: 100,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'down_price',
      minWidth: 130,
      showOverflow: false,
      slots: { default: 'deposit' },
      title: '定金/尾款',
    },
    {
      field: 'stock',
      title: '库存/销量',
      width: 110,
      formatter: ({ row }) => `${row.stock} / ${row.seles}`,
    },
    {
      field: 'start_time',
      minWidth: 210,
      showOverflow: false,
      slots: { default: 'dates' },
      title: '活动时间',
    },
    {
      field: 'is_show',
      slots: { default: 'show' },
      title: '上架',
      width: 80,
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
        const data = await listMerchantPresellActivesApi({
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
  rowConfig: { isHover: true, keyField: 'product_presell_id' },
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
    const selectedProduct = products.value.find(
      (item) => item.product_id === form.product_id,
    );
    if (
      !form.product_id ||
      !form.store_name.trim() ||
      form.price <= 0 ||
      form.stock <= 0 ||
      form.activity_dates.length !== 2
    ) {
      ElMessage.warning('请填写活动名称、商品、预售价、库存和活动时间');
      return;
    }
    if (selectedProduct && form.price > selectedProduct.price) {
      ElMessage.warning('预售价不能高于商品销售价');
      return;
    }
    if (
      form.presell_type === 2 &&
      (!validDeposit() || form.final_dates.length !== 2)
    ) {
      ElMessage.warning('定金预售需填写尾款支付期，且定金加尾款必须等于预售价');
      return;
    }
    const body: MerchantPresellSaveInput = {
      down_price: form.presell_type === 2 ? form.down_price : 0,
      end_time: form.activity_dates[1]!,
      final_end_time: form.presell_type === 2 ? form.final_dates[1]! : '',
      final_price: form.presell_type === 2 ? form.final_price : 0,
      final_start_time: form.presell_type === 2 ? form.final_dates[0]! : '',
      is_show: form.is_show,
      presell_type: form.presell_type,
      price: form.price,
      product_id: form.product_id,
      stock: form.stock,
      store_info: form.store_info.trim(),
      store_name: form.store_name.trim(),
      start_time: form.activity_dates[0]!,
    };
    saving.value = true;
    editDrawerApi.lock();
    try {
      if (editingID.value) {
        await updateMerchantPresellActiveApi(editingID.value, body);
      } else {
        await createMerchantPresellActiveApi(body);
      }
      ElMessage.success(editingID.value ? '预售活动已更新' : '预售活动已创建');
      editDrawerApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      editDrawerApi.unlock();
    }
  },
});

function validDeposit() {
  return (
    form.down_price > 0 &&
    form.final_price > 0 &&
    Math.abs(form.down_price + form.final_price - form.price) < 0.01
  );
}

function openCreate() {
  resetForm();
  editDrawerApi.setState({ title: '新建预售活动' }).open();
}

function openEdit(row: MerchantPresellActive) {
  editingID.value = row.product_presell_id;
  Object.assign(form, {
    activity_dates: [dateText(row.start_time), dateText(row.end_time)],
    down_price: row.down_price,
    final_dates: [
      dateText(row.final_start_time),
      dateText(row.final_end_time),
    ].filter(Boolean),
    final_price: row.final_price,
    is_show: row.is_show,
    presell_type: row.presell_type,
    price: row.price,
    product_id: row.product_id,
    stock: row.stock,
    store_info: row.store_info,
    store_name: row.store_name,
  });
  editDrawerApi.setState({ title: '编辑预售活动' }).open();
}

function selectedProductChanged(productID: number) {
  const product = products.value.find((item) => item.product_id === productID);
  if (!product || editingID.value) return;
  form.store_name = `${product.store_name} · 预售`;
  form.price = product.price;
  form.stock = product.stock;
}

async function toggle(row: MerchantPresellActive) {
  const isShow = row.is_show === 1 ? 0 : 1;
  await setMerchantPresellShowApi(row.product_presell_id, isShow);
  ElMessage.success(isShow ? '活动已上架' : '活动已下架');
  gridApi.reload();
}

async function remove(row: MerchantPresellActive) {
  try {
    await confirm({
      content: `删除预售活动「${row.store_name}」？删除后不可恢复。`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteMerchantPresellActiveApi(row.product_presell_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

watch(
  () => form.presell_type,
  (value) => {
    if (value === 1) {
      form.down_price = 0;
      form.final_price = 0;
      form.final_dates = [];
    }
  },
);

onMounted(async () => {
  const data = await listMerchantProductsApi({ limit: 100, page: 1, status: 1 });
  products.value = data.list || [];
});
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" @click="openCreate">新建预售活动</ElButton>
    </template>

    <Grid>
      <template #type="{ row }">
        <ElTag :type="row.presell_type === 2 ? 'warning' : 'success'">
          {{ row.presell_type === 2 ? '定金预售' : '全款预售' }}
        </ElTag>
      </template>
      <template #deposit="{ row }">
        <span v-if="row.presell_type === 2">
          ¥{{ Number(row.down_price).toFixed(2) }} /
          ¥{{ Number(row.final_price).toFixed(2) }}
        </span>
        <span v-else>—</span>
      </template>
      <template #dates="{ row }">
        {{ dateText(row.start_time) }} 至 {{ dateText(row.end_time) }}
      </template>
      <template #show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '上架' : '下架' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="warning" @click="toggle(row)">
          {{ row.is_show === 1 ? '下架' : '上架' }}
        </ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <EditDrawer class="w-[720px] max-w-[96vw]">
      <ElForm class="grid grid-cols-2 gap-x-4" label-width="96px">
        <ElFormItem class="col-span-2" label="预售名称" required>
          <ElInput v-model="form.store_name" maxlength="120" show-word-limit />
        </ElFormItem>
        <ElFormItem label="参与商品" required>
          <ElSelect
            v-model="form.product_id"
            :disabled="!!editingID"
            filterable
            class="w-full"
            placeholder="选择本店已审核商品"
            @change="selectedProductChanged"
          >
            <ElOption
              v-for="item in products"
              :key="item.product_id"
              :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`"
              :value="item.product_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="预售类型">
          <ElRadioGroup v-model="form.presell_type">
            <ElRadio :value="1">全款</ElRadio>
            <ElRadio :value="2">定金</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="预售价" required>
          <ElInputNumber
            v-model="form.price"
            :min="0.01"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="预售库存" required>
          <ElInputNumber
            v-model="form.stock"
            :min="1"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="活动时间" required>
          <ElDatePicker
            v-model="form.activity_dates"
            class="w-full"
            end-placeholder="活动结束"
            start-placeholder="活动开始"
            type="datetimerange"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </ElFormItem>
        <template v-if="form.presell_type === 2">
          <ElFormItem label="定金" required>
            <ElInputNumber
              v-model="form.down_price"
              :min="0.01"
              :precision="2"
              class="w-full"
            />
          </ElFormItem>
          <ElFormItem label="尾款" required>
            <ElInputNumber
              v-model="form.final_price"
              :min="0.01"
              :precision="2"
              class="w-full"
            />
          </ElFormItem>
          <ElFormItem class="col-span-2" label="尾款支付期" required>
            <ElDatePicker
              v-model="form.final_dates"
              class="w-full"
              end-placeholder="尾款截止"
              start-placeholder="尾款开始"
              type="datetimerange"
              value-format="YYYY-MM-DD HH:mm:ss"
            />
          </ElFormItem>
        </template>
        <ElFormItem class="col-span-2" label="活动说明">
          <ElInput
            v-model="form.store_info"
            :rows="3"
            maxlength="500"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem label="初始上架">
          <ElSwitch v-model="form.is_show" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </EditDrawer>
  </Page>
</template>
