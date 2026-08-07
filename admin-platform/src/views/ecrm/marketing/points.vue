<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElCol,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElRow,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getPlatformPointsSummaryApi,
  listPlatformPointsOrdersApi,
  listPlatformPointsProductsApi,
  updatePlatformPointsProductApi,
  type PlatformPointsOrder,
  type PlatformPointsProduct,
} from '#/api/core/platform-points';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canManage = ref(false);
const saving = ref(false);
const summary = ref({ total: 0, on_sale: 0, stock: 0 });
const editing = ref<PlatformPointsProduct>();
const form = reactive({ points_required: 1, sale_status: 1, stock: 0, version: 0 });

const productFormOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('商品名称'),
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '商户 ID' },
    fieldName: 'merchant_id',
    label: '商户 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '上架', value: 1 },
        { label: '下架', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '上架状态',
  },
]);

function buildProductParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const base = buildStandardListParams(page, formValues, { merField: 'merchant_id' });
  const statusRaw = formValues?.status;
  const merchantId = base.merchant_id;
  return {
    page: base.page,
    limit: base.limit,
    keyword: base.keyword,
    merchant_id:
      merchantId !== undefined && merchantId !== null && merchantId !== ''
        ? Number(merchantId)
        : undefined,
    sale_status: statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
    date_from: base.date_from,
    date_to: base.date_to,
  };
}

const productGridOptions: VxeGridProps<PlatformPointsProduct> = {
  columns: [
    { field: 'title', minWidth: 180, showOverflow: false, title: '商品' },
    {
      field: 'merchant_name',
      minWidth: 140,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => row.merchant_name || `商户 #${row.merchant_id}`,
    },
    { field: 'points_required', title: '所需积分', width: 110 },
    { field: 'stock', title: '库存', width: 90 },
    {
      field: 'sale_status',
      slots: { default: 'sale_status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ width: 80 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const [products, stats] = await Promise.all([
          listPlatformPointsProductsApi(buildProductParams(page, formValues)),
          getPlatformPointsSummaryApi(),
        ]);
        summary.value = stats;
        return { items: products.list || [], total: products.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const orderFormOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待支付', value: 'pending' },
        { label: '已支付', value: 'paid' },
        { label: '已关闭', value: 'closed' },
      ],
      placeholder: '全部',
    },
    fieldName: 'pay_status',
    label: '支付状态',
  },
]);

function buildOrderParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const payStatus = formValues?.pay_status;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    pay_status: typeof payStatus === 'string' && payStatus ? payStatus : undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

const orderGridOptions: VxeGridProps<PlatformPointsOrder> = {
  columns: [
    { field: 'order_no', minWidth: 190, showOverflow: false, title: '订单号' },
    { field: 'user_id', title: '用户 ID', width: 100 },
    { field: 'points_amount', title: '应付积分', width: 100 },
    { field: 'total_quantity', title: '件数', width: 80 },
    { field: 'pay_status', title: '支付状态', width: 100 },
    {
      field: 'created_at',
      minWidth: 180,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const data = await listPlatformPointsOrdersApi(buildOrderParams(page, formValues));
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

const [ProductGrid, productGridApi] = useVbenVxeGrid({
  formOptions: productFormOptions,
  gridOptions: productGridOptions,
});

const [OrderGrid] = useVbenVxeGrid({
  formOptions: orderFormOptions,
  gridOptions: orderGridOptions,
});

const [EditDrawer, editDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function edit(row: PlatformPointsProduct) {
  editing.value = row;
  Object.assign(form, {
    points_required: row.points_required,
    sale_status: row.sale_status,
    stock: row.stock,
    version: row.version,
  });
  editDrawerApi.setState({ title: '编辑积分商品' }).open();
}

async function save() {
  if (!editing.value || form.points_required < 1 || form.stock < 0) {
    ElMessage.warning('请填写正积分值和非负库存');
    return;
  }
  editDrawerApi.lock();
  saving.value = true;
  try {
    await updatePlatformPointsProductApi(editing.value.product_id, { ...form });
    editDrawerApi.close();
    ElMessage.success('积分商品已更新，新订单将使用新配置');
    productGridApi.reload();
  } finally {
    saving.value = false;
    editDrawerApi.unlock();
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canManage.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    codes.includes('marketing.points.manage');
  if (canManage.value) productGridApi.reload();
});
</script>

<template>
  <Page
    auto-content-height
    description="监管积分商品与积分订单。商品积分、库存、上下架仅影响新建订单；已创建订单的积分应付快照不可改写。"
    title="积分商城监管"
  >
    <ElAlert
      v-if="!canManage"
      title="当前账号没有积分商城监管权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <ElRow :gutter="16" class="mb-4">
        <ElCol :span="8">
          <div class="rounded border p-4">积分商品：{{ summary.total }}</div>
        </ElCol>
        <ElCol :span="8">
          <div class="rounded border p-4">上架商品：{{ summary.on_sale }}</div>
        </ElCol>
        <ElCol :span="8">
          <div class="rounded border p-4">可用库存：{{ summary.stock }}</div>
        </ElCol>
      </ElRow>

      <ProductGrid>
        <template #sale_status="{ row }">
          <ElTag :type="row.sale_status ? 'success' : 'info'">
            {{ row.sale_status ? '上架' : '下架' }}
          </ElTag>
        </template>
        <template #action="{ row }">
          <ElButton link type="primary" @click="edit(row)">编辑</ElButton>
        </template>
      </ProductGrid>

      <div class="mb-2 mt-6 text-base font-medium">积分订单（不展示个人资料）</div>
      <OrderGrid />

      <EditDrawer class="w-[500px]">
        <ElAlert
          type="warning"
          :closable="false"
          title="保存只影响后续新建订单；已有积分订单的应付积分、库存扣减和支付状态保持原快照。"
        />
        <ElForm class="mt-4" label-width="90px">
          <ElFormItem label="所需积分">
            <ElInputNumber v-model="form.points_required" :min="1" class="w-full" />
          </ElFormItem>
          <ElFormItem label="库存">
            <ElInputNumber v-model="form.stock" :min="0" class="w-full" />
          </ElFormItem>
          <ElFormItem label="上架状态">
            <ElRadioGroup v-model="form.sale_status">
              <ElRadio :value="1">上架</ElRadio>
              <ElRadio :value="0">下架</ElRadio>
            </ElRadioGroup>
          </ElFormItem>
        </ElForm>
      </EditDrawer>
    </template>
  </Page>
</template>
