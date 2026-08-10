<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { h, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  deletePlatformIntegralOrderApi,
  deliverPlatformIntegralOrderApi,
  exportPlatformIntegralOrdersApi,
  getPlatformIntegralOrderApi,
  listPlatformIntegralOrdersApi,
  type PlatformIntegralOrderRow,
} from '#/api/core/platform-integral-order';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canManage = ref(false);
const exporting = ref(false);
const lastFormValues = ref<Record<string, unknown>>({});
const detail = ref<PlatformIntegralOrderRow>();
const deliveryRow = ref<PlatformIntegralOrderRow>();
const deliveryForm = reactive({
  delivery_type: 'express',
  delivery_name: '',
  delivery_id: '',
  remark: '',
});

const SEARCH_OPTIONS = [
  { label: '全部', value: 'all' },
  { label: '订单号', value: 'order_sn' },
  { label: '收货人', value: 'real_name' },
  { label: '收货电话', value: 'phone' },
  { label: '用户ID', value: 'uid' },
];

const STATUS_OPTIONS = [
  { label: '全部', value: 'all' },
  { label: '待发货', value: '0' },
  { label: '待收货', value: '1' },
  { label: '待评价', value: '2' },
  { label: '已完成', value: '3' },
  { label: '已退款/取消', value: '-1' },
  { label: '已删除', value: '-10' },
];

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || v === '' || v === '--') return '—';
  return String(v);
}

function money(v?: number) {
  return `¥${Number(v || 0).toFixed(2)}`;
}

function statusTagType(row: PlatformIntegralOrderRow) {
  if (row.user_deleted || row.status === -10) return 'danger';
  if (row.status === 0) return 'warning';
  if (row.status === 1) return 'primary';
  if (row.status === 3) return 'success';
  return 'info';
}

function buildFilterParams(formValues?: Record<string, unknown>) {
  const values = formValues || {};
  lastFormValues.value = values;
  const range = Array.isArray(values.date_range) ? values.date_range : [];
  const status = String(values.status ?? 'all');
  return {
    status: status === 'all' || status === '' ? undefined : status,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    search_type: String(values.search_type ?? 'all') || 'all',
    keyword: String(values.keyword ?? '').trim() || undefined,
  };
}

function renderSearchPrepend(
  field: 'search_type',
  options: Array<{ label: string; value: string }>,
  defaultValue: string,
) {
  return (values: Record<string, any>, api: { setFieldValue: Function }) => ({
    prepend: () =>
      h(
        ElSelect,
        {
          modelValue: values[field] || defaultValue,
          style: { width: '110px' },
          'onUpdate:modelValue': (v: string) => api.setFieldValue(field, v),
        },
        () =>
          options.map((opt) =>
            h(ElOption, {
              label: opt.label,
              value: opt.value,
              key: opt.value,
            }),
          ),
      ),
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: STATUS_OPTIONS,
        placeholder: '全部',
      },
      defaultValue: 'all',
      fieldName: 'status',
      label: '订单状态',
    },
    {
      ...LIST_DATE_RANGE_FIELD,
      label: '创建时间',
    },
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '请输入内容' },
      defaultValue: '',
      fieldName: 'keyword',
      label: '搜索',
      renderComponentContent: renderSearchPrepend(
        'search_type',
        SEARCH_OPTIONS,
        'all',
      ),
    },
    {
      component: 'Input',
      defaultValue: 'all',
      dependencies: { show: () => false, triggerFields: [''] },
      fieldName: 'search_type',
      label: '搜索类型',
    },
  ],
  {
    commonConfig: { componentProps: { class: 'w-full' } },
    submitButtonOptions: { content: '搜索' },
    handleSubmit: async (values) => {
      lastFormValues.value = { ...values };
      await gridApi.reload(values);
    },
    handleReset: async () => {
      await gridApi.formApi?.resetForm();
      const values = (await gridApi.formApi?.getValues()) ?? {};
      lastFormValues.value = { ...values };
      await gridApi.reload(values);
    },
  },
);

const gridOptions: VxeGridProps<PlatformIntegralOrderRow> = {
  columns: [
    {
      type: 'expand',
      width: 48,
      slots: { content: 'expandContent' },
    },
    {
      field: 'order_sn',
      minWidth: 200,
      showOverflow: false,
      slots: { default: 'orderSn' },
      title: '订单编号',
    },
    {
      field: 'real_name',
      minWidth: 110,
      showOverflow: false,
      title: '收货人',
      formatter: ({ row }) => dash(row.real_name),
    },
    {
      field: 'product',
      minWidth: 260,
      showOverflow: false,
      slots: { default: 'productInfo' },
      title: '商品信息',
    },
    {
      field: 'points_amount',
      title: '兑换积分',
      width: 100,
      formatter: ({ cellValue }) => String(Number(cellValue || 0)),
    },
    {
      field: 'pay_amount',
      title: '兑换金额',
      width: 100,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
    },
    {
      field: 'status_label',
      slots: { default: 'status' },
      title: '订单状态',
      width: 110,
    },
    {
      field: 'create_time',
      minWidth: 170,
      showOverflow: 'tooltip',
      title: '下单时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    platformListActionColumn({ width: 180 }),
  ],
  expandConfig: { trigger: 'default' },
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const values =
          formValues && Object.keys(formValues).length > 0
            ? formValues
            : lastFormValues.value;
        const data = await listPlatformIntegralOrdersApi({
          page: page.currentPage,
          limit: page.pageSize,
          ...buildFilterParams(values),
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
  title: '订单详情',
});

const [DeliveryDrawer, deliveryDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  title: '发送货',
  confirmText: '确认发货',
  onConfirm: async () => {
    if (!deliveryRow.value) return;
    if (
      deliveryForm.delivery_type !== 'pickup' &&
      deliveryForm.delivery_type !== 'service' &&
      (!deliveryForm.delivery_name.trim() || !deliveryForm.delivery_id.trim())
    ) {
      ElMessage.warning('请填写物流公司和运单号');
      return;
    }
    deliveryDrawerApi.lock();
    try {
      await deliverPlatformIntegralOrderApi(deliveryRow.value.order_id, {
        delivery_type: deliveryForm.delivery_type,
        delivery_name: deliveryForm.delivery_name.trim(),
        delivery_id: deliveryForm.delivery_id.trim(),
        remark: deliveryForm.remark.trim() || undefined,
      });
      ElMessage.success('发货成功');
      deliveryDrawerApi.close();
      await gridApi.reload();
    } finally {
      deliveryDrawerApi.unlock();
    }
  },
});

async function openDetail(row: PlatformIntegralOrderRow) {
  detailDrawerApi.setState({ title: '订单详情', loading: true }).open();
  try {
    detail.value = await getPlatformIntegralOrderApi(row.order_id);
  } catch {
    ElMessage.error('加载订单详情失败');
    detailDrawerApi.close();
  } finally {
    detailDrawerApi.setState({ loading: false });
  }
}

function openDelivery(row: PlatformIntegralOrderRow) {
  deliveryRow.value = row;
  Object.assign(deliveryForm, {
    delivery_type: 'express',
    delivery_name: '',
    delivery_id: '',
    remark: '',
  });
  deliveryDrawerApi.open();
}

async function removeOrder(row: PlatformIntegralOrderRow) {
  try {
    await confirm({
      content: `确认删除积分订单 ${row.order_sn}？删除后列表默认不再展示。`,
      icon: 'warning',
      title: '提示',
    });
    await deletePlatformIntegralOrderApi(row.order_id);
    ElMessage.success('已删除');
    await gridApi.reload();
  } catch {
    // cancelled
  }
}

async function exportRows() {
  if (exporting.value) return;
  exporting.value = true;
  try {
    const filters = buildFilterParams(lastFormValues.value);
    const result = await exportPlatformIntegralOrdersApi(filters);
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name || '积分订单列表.csv';
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      `已导出 ${result.row_count} 条${result.truncated ? '（已截断）' : ''}`,
    );
  } catch {
    ElMessage.error('导出失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canManage.value =
    codes.includes('marketing.integral.orders.manage') ||
    codes.includes('marketing.points.manage');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton type="primary" :loading="exporting" @click="exportRows">
          导出
        </ElButton>
      </template>

      <template #orderSn="{ row }">
        <div class="order-sn" :class="{ 'is-deleted': row.user_deleted }">
          <div>{{ row.order_sn }}</div>
          <div v-if="row.user_deleted" class="order-sn__tip">用户已删除</div>
        </div>
      </template>

      <template #productInfo="{ row }">
        <div v-if="row.product" class="order-product">
          <ElImage
            class="order-product__thumb"
            :src="resolveCosMediaUrl(row.product.product_image || '')"
            fit="cover"
          />
          <div class="order-product__text">
            <div class="order-product__title">
              {{ row.product.product_info || '—' }}
            </div>
            <div class="order-product__meta">
              {{ dash(row.product.product_sku) }}
            </div>
            <div class="order-product__meta">
              {{ money(row.product.product_price) }} ×
              {{ row.product.product_num || 0 }}
            </div>
          </div>
        </div>
        <span v-else>—</span>
      </template>

      <template #status="{ row }">
        <ElTag :type="statusTagType(row)" effect="plain">
          {{ row.status_label || '未知' }}
        </ElTag>
      </template>

      <template #action="{ row }">
        <ElButton
          v-if="canManage && row.can_deliver"
          link
          type="primary"
          @click="openDelivery(row)"
        >
          发送货
        </ElButton>
        <ElButton link type="primary" @click="openDetail(row)">
          订单详情
        </ElButton>
        <ElButton
          v-if="canManage && row.can_delete"
          link
          type="danger"
          @click="removeOrder(row)"
        >
          删除
        </ElButton>
      </template>

      <template #expandContent="{ row }">
        <div class="order-expand">
          <div class="order-expand__item">
            <span class="label">用户备注</span>
            <span class="value">{{ dash(row.user_remark) }}</span>
          </div>
          <div class="order-expand__item">
            <span class="label">店铺备注</span>
            <span class="value">{{ dash(row.merchant_remark) }}</span>
          </div>
        </div>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="detail">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="订单编号">
            <span :class="{ 'text-danger': detail.user_deleted }">
              {{ detail.order_sn }}
            </span>
            <template v-if="detail.user_deleted">（用户已删除）</template>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="订单状态">
            <ElTag :type="statusTagType(detail)" effect="plain">
              {{ detail.status_label }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="收货人">
            {{ dash(detail.real_name) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="联系电话">
            {{ dash(detail.user_phone) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="收货地址">
            {{ dash(detail.user_address) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="兑换积分">
            {{ detail.points_amount || 0 }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="兑换金额">
            {{ money(detail.pay_amount) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="物流公司">
            {{ dash(detail.delivery_name) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="运单号">
            {{ dash(detail.delivery_id) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="用户备注">
            {{ dash(detail.user_remark) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="店铺备注">
            {{ dash(detail.merchant_remark) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="下单时间">
            {{ formatShanghaiDateTime(detail.create_time) }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <div class="mb-3 mt-6 text-base font-medium">商品明细</div>
        <ElTable :data="detail.products || []" border>
          <ElTableColumn label="商品" min-width="220">
            <template #default="{ row }">
              <div class="order-product">
                <ElImage
                  class="order-product__thumb"
                  :src="resolveCosMediaUrl(row.product_image || '')"
                  fit="cover"
                />
                <div class="order-product__text">
                  {{ row.product_info || '—' }}
                </div>
              </div>
            </template>
          </ElTableColumn>
          <ElTableColumn
            label="规格"
            min-width="120"
            prop="product_sku"
            show-overflow-tooltip
          />
          <ElTableColumn label="单价" width="104">
            <template #default="{ row }">
              {{ money(row.product_price) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="数量" prop="product_num" width="76" />
          <ElTableColumn label="小计" width="104">
            <template #default="{ row }">
              {{ money(row.total_price) }}
            </template>
          </ElTableColumn>
        </ElTable>
      </template>
    </DetailDrawer>

    <DeliveryDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="物流类型">
          <ElSelect v-model="deliveryForm.delivery_type" class="w-full">
            <ElOption label="快递发货" value="express" />
            <ElOption label="同城配送" value="city" />
            <ElOption label="到店自提" value="pickup" />
            <ElOption label="虚拟发货" value="service" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem
          v-if="
            deliveryForm.delivery_type !== 'pickup' &&
            deliveryForm.delivery_type !== 'service'
          "
          label="物流公司"
          required
        >
          <ElInput
            v-model="deliveryForm.delivery_name"
            placeholder="例如：顺丰速运"
          />
        </ElFormItem>
        <ElFormItem
          v-if="
            deliveryForm.delivery_type !== 'pickup' &&
            deliveryForm.delivery_type !== 'service'
          "
          label="运单号"
          required
        >
          <ElInput
            v-model="deliveryForm.delivery_id"
            placeholder="填写物流运单号"
          />
        </ElFormItem>
        <ElFormItem label="店铺备注">
          <ElInput
            v-model="deliveryForm.remark"
            type="textarea"
            :rows="3"
            placeholder="可选，写入店铺备注"
          />
        </ElFormItem>
      </ElForm>
    </DeliveryDrawer>
  </Page>
</template>

<style scoped>
.order-sn.is-deleted,
.text-danger {
  color: var(--el-color-danger);
}

.order-sn__tip {
  margin-top: 2px;
  font-size: 12px;
  color: var(--el-color-danger);
}

.order-product {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.order-product__thumb {
  width: 52px;
  height: 52px;
  flex-shrink: 0;
  border-radius: 4px;
  background: var(--el-fill-color-light);
}

.order-product__text {
  min-width: 0;
  line-height: 1.4;
}

.order-product__title {
  word-break: break-all;
}

.order-product__meta {
  margin-top: 2px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.order-expand {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px 24px;
  padding: 8px 12px;
}

.order-expand__item {
  display: flex;
  gap: 8px;
  min-width: 0;
}

.order-expand__item .label {
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
}

.order-expand__item .value {
  min-width: 0;
  word-break: break-all;
}
</style>
