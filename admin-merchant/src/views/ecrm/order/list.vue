<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
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
  deliverMerchantOrderApi,
  getMerchantOrderApi,
  listMerchantOrdersApi,
  verifyMerchantOrderApi,
  type MerchantOrder,
} from '#/api/core/merchant-trade';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canDeliver = ref(false);
const canVerifyAction = ref(false);
const detail = ref<MerchantOrder>();
const deliveryRow = ref<MerchantOrder>();
const delivering = ref(false);
const deliveryForm = reactive({
  delivery_id: '',
  delivery_name: '',
  delivery_type: 'express',
});

function statusInfo(row: MerchantOrder) {
  if (row.paid !== 1) return { label: '待支付', type: 'info' as const };
  return (
    {
      0: { label: '待发货', type: 'warning' as const },
      1: { label: '待收货', type: 'primary' as const },
      3: { label: '已完成', type: 'success' as const },
    }[row.status] || { label: '未知状态', type: 'info' as const }
  );
}

function payType(value: number) {
  return (
    { 0: '余额', 1: '微信', 2: '支付宝', 7: '模拟支付', 8: '积分' }[value] ||
    '未知'
  );
}

function rowCanDeliver(row: MerchantOrder) {
  return canDeliver.value && row.paid === 1 && row.status === 0;
}

function canVerify(row: MerchantOrder) {
  return canVerifyAction.value && Boolean(row.can_verify);
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '订单号 / 手机号 / 收货人 / 订单ID',
    },
    fieldName: 'keyword',
    label: '订单搜索',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '订单号' },
    fieldName: 'order_sn',
    label: '订单号',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待支付', value: 0 },
        { label: '已支付', value: 1 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'paid',
    label: '支付状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待发货', value: 0 },
        { label: '待收货', value: 1 },
        { label: '已完成', value: 3 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'status',
    label: '订单状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '余额', value: 0 },
        { label: '微信', value: 1 },
        { label: '支付宝', value: 2 },
        { label: '模拟支付', value: 7 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'pay_type',
    label: '支付方式',
  },
]);

const gridOptions: VxeGridProps<MerchantOrder> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'order_id', title: 'ID', width: 90 },
    { field: 'order_sn', minWidth: 170, showOverflow: false, title: '订单号' },
    {
      field: 'real_name',
      minWidth: 140,
      showOverflow: false,
      slots: { default: 'user' },
      title: '用户',
    },
    { field: 'total_num', title: '商品数', width: 80 },
    {
      field: 'pay_price',
      title: '实付',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'pay_type',
      title: '支付方式',
      width: 100,
      formatter: ({ cellValue }) => payType(Number(cellValue)),
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '订单状态',
      width: 100,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '下单时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 180 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const data = await listMerchantOrdersApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          order_sn: String(formValues?.order_sn ?? '').trim() || undefined,
          paid:
            formValues?.paid === 0 || formValues?.paid === 1
              ? Number(formValues.paid)
              : undefined,
          status:
            formValues?.status === 0 ||
            formValues?.status === 1 ||
            formValues?.status === 3
              ? Number(formValues.status)
              : undefined,
          pay_type:
            formValues?.pay_type === 0 ||
            formValues?.pay_type === 1 ||
            formValues?.pay_type === 2 ||
            formValues?.pay_type === 7
              ? Number(formValues.pay_type)
              : undefined,
          date_from: range[0],
          date_to: range[1],
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
  class: 'w-[900px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

const [DeliveryDrawer, deliveryDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  title: '订单发货',
  confirmText: '确认发货',
  onConfirm: async () => {
    if (!deliveryRow.value) return;
    if (!deliveryForm.delivery_name.trim() || !deliveryForm.delivery_id.trim()) {
      ElMessage.warning('请填写物流公司和运单号');
      return;
    }
    delivering.value = true;
    deliveryDrawerApi.lock();
    try {
      await deliverMerchantOrderApi(deliveryRow.value.order_id, {
        delivery_id: deliveryForm.delivery_id.trim(),
        delivery_name: deliveryForm.delivery_name.trim(),
        delivery_type: deliveryForm.delivery_type,
      });
      ElMessage.success('订单已发货');
      deliveryDrawerApi.close();
      gridApi.reload();
    } finally {
      delivering.value = false;
      deliveryDrawerApi.unlock();
    }
  },
});

async function openDetail(row: MerchantOrder) {
  detailDrawerApi.setState({ title: '订单详情', loading: true }).open();
  try {
    detail.value = await getMerchantOrderApi(row.order_id);
  } finally {
    detailDrawerApi.setState({ loading: false });
  }
}

function openDelivery(row: MerchantOrder) {
  deliveryRow.value = row;
  Object.assign(deliveryForm, {
    delivery_id: '',
    delivery_name: '',
    delivery_type: 'express',
  });
  deliveryDrawerApi.open();
}

async function verify(row: MerchantOrder) {
  try {
    await confirm({
      content: '确认完成该订单核销？核销后订单将变为已完成，操作不可撤销。',
      icon: 'warning',
      title: '核销确认',
    });
    await verifyMerchantOrderApi(row.order_id);
    ElMessage.success('订单已核销');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canDeliver.value = codes.includes('order.deliver');
  canVerifyAction.value = codes.includes('order.verify.action');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #user="{ row }">
        <div>{{ row.real_name || '—' }}</div>
        <div class="text-xs text-muted-foreground">
          {{ row.user_phone || '—' }}
        </div>
      </template>
      <template #status="{ row }">
        <ElTag :type="statusInfo(row).type">{{ statusInfo(row).label }}</ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          v-if="rowCanDeliver(row)"
          link
          type="success"
          @click="openDelivery(row)"
        >
          发货
        </ElButton>
        <ElButton
          v-if="canVerify(row)"
          link
          type="warning"
          @click="verify(row)"
        >
          核销
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="detail">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="订单号">
            {{ detail.order_sn }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="订单状态">
            <ElTag :type="statusInfo(detail).type">
              {{ statusInfo(detail).label }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="收货人">
            {{ detail.real_name || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="联系电话">
            {{ detail.user_phone || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="收货地址">
            {{ detail.user_address || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="支付方式">
            {{ payType(detail.pay_type) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="实付金额">
            ¥{{ Number(detail.pay_price).toFixed(2) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="物流公司">
            {{ detail.delivery_name || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="运单号">
            {{ detail.delivery_id || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="用户备注">
            {{ detail.mark || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="下单时间">
            {{ formatShanghaiDateTime(detail.create_time) }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <div class="mb-3 mt-6 text-base font-medium">商品明细</div>
        <ElTable :data="detail.products || []" border>
          <ElTableColumn label="商品 ID" prop="product_id" width="92" />
          <ElTableColumn
            label="商品信息"
            min-width="180"
            prop="product_info"
            show-overflow-tooltip
          />
          <ElTableColumn label="规格" min-width="120" prop="product_sku" />
          <ElTableColumn label="单价" width="104">
            <template #default="{ row }">
              ¥{{ Number(row.product_price).toFixed(2) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="数量" prop="product_num" width="76" />
          <ElTableColumn label="小计" width="104">
            <template #default="{ row }">
              ¥{{ Number(row.total_price).toFixed(2) }}
            </template>
          </ElTableColumn>
        </ElTable>
      </template>
    </DetailDrawer>

    <DeliveryDrawer>
      <ElForm label-width="84px">
        <ElFormItem label="物流类型">
          <ElSelect v-model="deliveryForm.delivery_type" class="w-full">
            <ElOption label="快递" value="express" />
            <ElOption label="同城配送" value="local" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="物流公司" required>
          <ElInput
            v-model="deliveryForm.delivery_name"
            placeholder="例如：顺丰速运"
          />
        </ElFormItem>
        <ElFormItem label="运单号" required>
          <ElInput
            v-model="deliveryForm.delivery_id"
            placeholder="填写物流运单号"
          />
        </ElFormItem>
      </ElForm>
    </DeliveryDrawer>
  </Page>
</template>
