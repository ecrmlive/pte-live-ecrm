<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  selectPlatformIntegralCategoriesApi,
  type PlatformIntegralCategoryRow,
} from '#/api/core/platform-integral-category';
import {
  copyPlatformPointsProductApi,
  createPlatformPointsProductApi,
  deletePlatformPointsProductApi,
  listPlatformPointsExchangesApi,
  listPlatformPointsProductsApi,
  quickAddPlatformPointsProductApi,
  updatePlatformPointsProductApi,
  updatePlatformPointsProductStatusApi,
  type PlatformPointsExchange,
  type PlatformPointsProduct,
} from '#/api/core/platform-points';
import type { PlatformProduct } from '#/api/core/platform-catalog';
import ImageField from '#/components/shop/image-field.vue';
import ProductPickerDialog from '#/components/shop/product-picker-dialog.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

type DrawerMode = 'create' | 'edit' | 'quick';

const canManage = ref(false);
const categories = ref<PlatformIntegralCategoryRow[]>([]);
const pickerOpen = ref(false);
const drawerMode = ref<DrawerMode>('create');
const editing = ref<PlatformPointsProduct>();
const exchangeProduct = ref<PlatformPointsProduct>();

const form = reactive({
  title: '',
  cover_url: '',
  cate_id: undefined as number | undefined,
  original_price: 0,
  points_required: 1,
  stock: 0,
  sort: 0,
  sale_status: 1,
  source_product_id: 0,
  merchant_id: 0,
  store_id: 0,
  merchant_name: '',
  store_name: '',
  version: 0,
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    ...LIST_DATE_RANGE_FIELD,
    label: '创建时间',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '上架', value: 1 },
        { label: '下架', value: 0 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'status',
    label: '上架状态',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入商品名称/ID',
    },
    fieldName: 'keyword',
    label: '商品搜索',
  },
]);

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const base = buildStandardListParams(page, formValues);
  const statusRaw = formValues?.status;
  return {
    page: base.page,
    limit: base.limit,
    keyword: base.keyword,
    sale_status:
      statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
    date_from: base.date_from,
    date_to: base.date_to,
  };
}

const gridOptions: VxeGridProps<PlatformPointsProduct> = {
  columns: [
    { field: 'product_id', title: 'ID', width: 80 },
    {
      field: 'cover_url',
      slots: { default: 'cover' },
      title: '商品图片',
      width: 100,
    },
    {
      field: 'title',
      minWidth: 200,
      showOverflow: false,
      title: '商品标题',
    },
    { field: 'points_required', title: '兑换积分', width: 100 },
    {
      field: 'original_price',
      formatter: ({ cellValue }) =>
        Number(cellValue || 0).toFixed(2),
      title: '兑换金额',
      width: 100,
    },
    { field: 'stock', title: '库存', width: 80 },
    {
      field: 'sales',
      formatter: ({ cellValue }) => Number(cellValue || 0),
      title: '已兑换数量',
      width: 100,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 170,
      showOverflow: false,
      title: '创建时间',
    },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'sale_status',
      slots: { default: 'sale_status' },
      title: '状态',
      width: 100,
    },
    platformListActionColumn({ width: 240 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const data = await listPlatformPointsProductsApi(
          buildParams(page, formValues),
        );
        return { items: data.list || [], total: data.total || 0 };
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

const exchangeGridOptions: VxeGridProps<PlatformPointsExchange> = {
  columns: [
    { field: 'order_no', minWidth: 180, showOverflow: false, title: '订单号' },
    { field: 'user_id', title: '用户 ID', width: 100 },
    { field: 'quantity', title: '数量', width: 80 },
    { field: 'points_amount', title: '积分', width: 90 },
    {
      field: 'pay_status',
      slots: { default: 'pay_status' },
      title: '支付状态',
      width: 100,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 170,
      title: '创建时间',
    },
  ],
  pagerConfig: platformListPagerConfig({ pageSizes: [10, 20, 50] }),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const id = Number(exchangeProduct.value?.product_id || 0);
        if (!id) return { items: [], total: 0 };
        const data = await listPlatformPointsExchangesApi(id, {
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'order_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
});

const [ExchangeGrid, exchangeGridApi] = useVbenVxeGrid({
  gridOptions: exchangeGridOptions,
});

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

const [ExchangeDrawer, exchangeDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
  title: '兑换记录',
});

function coverOf(row: PlatformPointsProduct) {
  return resolveCosMediaUrl(String(row.cover_url || '').trim());
}

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    title: '',
    cover_url: '',
    cate_id: undefined,
    original_price: 0,
    points_required: 1,
    stock: 0,
    sort: 0,
    sale_status: 1,
    source_product_id: 0,
    merchant_id: 0,
    store_id: 0,
    merchant_name: '',
    store_name: '',
    version: 0,
  });
}

function openCreate() {
  drawerMode.value = 'create';
  resetForm();
  formDrawerApi.setState({ title: '添加积分商品' }).open();
}

function openEdit(row: PlatformPointsProduct) {
  drawerMode.value = 'edit';
  editing.value = row;
  Object.assign(form, {
    title: row.title || '',
    cover_url: row.cover_url || '',
    cate_id: row.cate_id || undefined,
    original_price: Number(row.original_price || 0),
    points_required: Number(row.points_required || 1),
    stock: Number(row.stock || 0),
    sort: Number(row.sort || 0),
    sale_status: row.sale_status === 1 ? 1 : 0,
    source_product_id: Number(row.source_product_id || 0),
    merchant_id: Number(row.merchant_id || 0),
    store_id: Number(row.store_id || 0),
    merchant_name: row.merchant_name || '',
    store_name: row.store_name || '',
    version: Number(row.version || 0),
  });
  formDrawerApi.setState({ title: '编辑积分商品' }).open();
}

function openQuickAdd() {
  pickerOpen.value = true;
}

function onProductPicked(product: PlatformProduct) {
  drawerMode.value = 'quick';
  editing.value = undefined;
  Object.assign(form, {
    title: product.store_name || product.title || '',
    cover_url: product.image || '',
    cate_id: undefined,
    original_price: Number(product.price || 0),
    points_required: Math.max(1, Math.round(Number(product.ot_price || 0)) || 1),
    stock: Math.max(0, Number(product.stock || 0)),
    sort: 0,
    sale_status: 1,
    source_product_id: Number(product.product_id || 0),
    merchant_id: Number(product.mer_id || 0),
    store_id: 0,
    merchant_name: product.mer_name || '',
    store_name: '',
    version: 0,
  });
  formDrawerApi.setState({ title: '快速添加积分商品' }).open();
}

async function save() {
  const title = form.title.trim();
  if (!title) {
    ElMessage.warning('请填写商品标题');
    return;
  }
  if (form.points_required < 1) {
    ElMessage.warning('兑换积分须大于 0');
    return;
  }
  if (form.stock < 0 || form.original_price < 0) {
    ElMessage.warning('兑换金额与库存不能为负');
    return;
  }
  formDrawerApi.lock();
  try {
    if (drawerMode.value === 'edit' && editing.value) {
      await updatePlatformPointsProductApi(editing.value.product_id, {
        title,
        cover_url: form.cover_url,
        cate_id: form.cate_id || 0,
        original_price: form.original_price,
        points_required: form.points_required,
        stock: form.stock,
        sort: form.sort,
        sale_status: form.sale_status === 1 ? 1 : 0,
        version: form.version,
      });
    } else if (drawerMode.value === 'quick') {
      await quickAddPlatformPointsProductApi({
        source_product_id: form.source_product_id,
        title,
        cover_url: form.cover_url,
        cate_id: form.cate_id || 0,
        original_price: form.original_price,
        points_required: form.points_required,
        stock: form.stock,
        sort: form.sort,
        sale_status: form.sale_status === 1 ? 1 : 0,
      });
    } else {
      await createPlatformPointsProductApi({
        title,
        cover_url: form.cover_url,
        cate_id: form.cate_id || 0,
        original_price: form.original_price,
        points_required: form.points_required,
        stock: form.stock,
        sort: form.sort,
        sale_status: form.sale_status === 1 ? 1 : 0,
        merchant_id: form.merchant_id || undefined,
        store_id: form.store_id || undefined,
        merchant_name: form.merchant_name || undefined,
        store_name: form.store_name || undefined,
      });
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: PlatformPointsProduct, enabled: boolean) {
  const before = row.sale_status;
  row.sale_status = enabled ? 1 : 0;
  try {
    const updated = await updatePlatformPointsProductStatusApi(
      row.product_id,
      enabled ? 1 : 0,
    );
    row.version = updated.version;
  } catch {
    row.sale_status = before;
  }
}

async function openExchanges(row: PlatformPointsProduct) {
  exchangeProduct.value = row;
  exchangeDrawerApi
    .setState({ title: `兑换记录 · ${row.title || `#${row.product_id}`}` })
    .open();
  await exchangeGridApi.reload();
}

async function copyRow(row: PlatformPointsProduct) {
  try {
    await confirm({
      content: '确定复制该积分商品吗？复制后默认下架。',
      icon: 'warning',
      title: '提示',
    });
    await copyPlatformPointsProductApi(row.product_id);
    ElMessage.success('已复制');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

async function removeRow(row: PlatformPointsProduct) {
  try {
    await confirm({
      content: '确定删除该积分商品吗？删除后不可恢复。',
      icon: 'warning',
      title: '提示',
    });
    await deletePlatformPointsProductApi(row.product_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

function payStatusLabel(status: string) {
  if (status === 'paid') return '已支付';
  if (status === 'pending') return '待支付';
  if (status === 'closed') return '已关闭';
  return status || '—';
}

onMounted(async () => {
  const [profile, codes, cateRes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
    selectPlatformIntegralCategoriesApi().catch(() => ({ list: [] })),
  ]);
  categories.value = cateRes.list || [];
  canManage.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    (codes.includes('marketing.points.manage') ||
      codes.includes('marketing.integral.products'));
  if (canManage.value) gridApi.reload();
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
          添加积分商品
        </ElButton>
        <ElButton v-if="canManage" type="success" @click="openQuickAdd">
          快速添加
        </ElButton>
      </template>
      <template #cover="{ row }">
        <ElImage
          v-if="coverOf(row)"
          :src="coverOf(row)"
          fit="cover"
          class="integral-product-cover"
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>
      <template #sale_status="{ row }">
        <ElSwitch
          :model-value="row.sale_status === 1"
          inline-prompt
          active-text="上架"
          inactive-text="下架"
          @change="
            (enabled: string | number | boolean) =>
              changeStatus(row, Boolean(enabled))
          "
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openExchanges(row)">
          兑换记录
        </ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="copyRow(row)">复制</ElButton>
        <ElButton link type="primary" @click="removeRow(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="110px">
        <ElFormItem
          v-if="drawerMode === 'quick' && form.source_product_id"
          label="来源商品"
        >
          <span>ID {{ form.source_product_id }}</span>
          <span v-if="form.merchant_name" class="ml-2 text-gray-500">
            {{ form.merchant_name }}
          </span>
        </ElFormItem>
        <ElFormItem label="商品图片">
          <ImageField v-model="form.cover_url" />
        </ElFormItem>
        <ElFormItem label="商品标题" required>
          <ElInput
            v-model="form.title"
            maxlength="100"
            show-word-limit
            placeholder="请输入商品标题"
          />
        </ElFormItem>
        <ElFormItem label="积分分类">
          <ElSelect
            v-model="form.cate_id"
            clearable
            class="w-full"
            placeholder="请选择分类"
          >
            <ElOption
              v-for="item in categories"
              :key="item.store_category_id"
              :label="item.cate_name"
              :value="item.store_category_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="兑换积分" required>
          <ElInputNumber
            v-model="form.points_required"
            :min="1"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="兑换金额">
          <ElInputNumber
            v-model="form.original_price"
            :min="0"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="库存" required>
          <ElInputNumber
            v-model="form.stock"
            :min="0"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber
            v-model="form.sort"
            :min="0"
            :max="99999"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="上架状态">
          <ElSwitch
            v-model="form.sale_status"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="上架"
            inactive-text="下架"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>

    <ExchangeDrawer>
      <ExchangeGrid>
        <template #pay_status="{ row }">
          <ElTag
            :type="
              row.pay_status === 'paid'
                ? 'success'
                : row.pay_status === 'pending'
                  ? 'warning'
                  : 'info'
            "
          >
            {{ payStatusLabel(row.pay_status) }}
          </ElTag>
        </template>
      </ExchangeGrid>
    </ExchangeDrawer>

    <ProductPickerDialog v-model:open="pickerOpen" @select="onProductPicked" />
  </Page>
</template>

<style scoped>
.integral-product-cover {
  display: block;
  width: 36px;
  height: 36px;
  border-radius: 4px;
}
</style>
