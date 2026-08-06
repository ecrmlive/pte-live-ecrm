<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref, watch } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElMessage,
  ElRadioButton,
  ElRadioGroup,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
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

type VerifyTab = 'pending' | 'verified';

const activeTab = ref<VerifyTab>('pending');
const detail = ref<MerchantOrder>();
const canVerifyAction = ref(false);

function statusInfo(row: MerchantOrder) {
  if (row.paid !== 1) return { label: '待支付', type: 'info' as const };
  return (
    {
      0: { label: '待核销', type: 'warning' as const },
      1: { label: '待收货', type: 'primary' as const },
      3: { label: '已核销', type: 'success' as const },
    }[row.status] || { label: '处理中', type: 'info' as const }
  );
}

function payType(value: number) {
  return (
    { 0: '余额', 1: '微信', 2: '支付宝', 7: '模拟支付', 8: '积分' }[value] ||
    '未知'
  );
}

function rowCanVerify(row: MerchantOrder) {
  return (
    canVerifyAction.value &&
    activeTab.value === 'pending' &&
    Boolean(row.can_verify)
  );
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
]);

const gridOptions: VxeGridProps<MerchantOrder> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'order_sn', minWidth: 170, showOverflow: false, title: '订单号' },
    {
      field: 'real_name',
      minWidth: 140,
      showOverflow: false,
      slots: { default: 'user' },
      title: '用户',
    },
    { field: 'total_num', title: '商品数', width: 88 },
    {
      field: 'pay_price',
      title: '实付金额',
      width: 116,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'pay_type',
      title: '支付方式',
      width: 104,
      formatter: ({ cellValue }) => payType(Number(cellValue)),
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 104,
    },
    {
      field: 'verify_status',
      slots: { default: 'verify' },
      title: '核销',
      width: 96,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '下单时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 148 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const data = await listMerchantOrdersApi({
          page: page.currentPage,
          limit: page.pageSize,
          paid: 1,
          verify_tab: activeTab.value,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
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

async function openDetail(row: MerchantOrder) {
  detailDrawerApi.setState({ title: '订单详情', loading: true }).open();
  try {
    detail.value = await getMerchantOrderApi(row.order_id);
  } finally {
    detailDrawerApi.setState({ loading: false });
  }
}

async function verify(row: MerchantOrder) {
  try {
    await confirm({
      content: `确认核销订单 ${row.order_sn}？核销后订单将标记为已完成，操作不可撤销。`,
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

watch(activeTab, () => {
  gridApi.reload();
});

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canVerifyAction.value = codes.includes('order.verify.action');
});
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElRadioGroup v-model="activeTab" class="mr-3">
        <ElRadioButton value="pending">待核销</ElRadioButton>
        <ElRadioButton value="verified">已核销</ElRadioButton>
      </ElRadioGroup>
    </template>

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
      <template #verify="{ row }">
        <ElTag v-if="row.verify_status === 'used'" type="success">
          已核销
        </ElTag>
        <ElTag
          v-else-if="row.has_verify_code || row.can_verify"
          type="warning"
        >
          待核销
        </ElTag>
        <span v-else class="text-muted-foreground">—</span>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          v-if="rowCanVerify(row)"
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
          <ElDescriptionsItem label="配送方式">
            {{ detail.delivery_type || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="核销状态">
            <span v-if="detail.verify_status === 'used'">已核销</span>
            <span v-else-if="detail.can_verify || detail.has_verify_code">
              待核销
            </span>
            <span v-else>—</span>
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
  </Page>
</template>
