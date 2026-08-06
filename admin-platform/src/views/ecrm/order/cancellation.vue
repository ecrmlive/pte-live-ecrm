<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElSkeleton,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getPlatformOrderApi,
  listPlatformOrdersApi,
  type PlatformOrder,
} from '#/api/core/platform-trade';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const detail = ref<PlatformOrder>();
const detailLoading = ref(false);

function orderStatus(value: number) {
  return (
    (
      {
        '-1': '已取消/关闭',
        0: '待发货',
        1: '待收货',
        2: '待评价',
        3: '已完成',
      } as Record<number, string>
    )[value] || '未知'
  );
}

function payType(value: number) {
  return (
    ({ 0: '余额', 1: '微信', 2: '支付宝', 7: '模拟支付', 8: '积分' } as Record<
      number,
      string
    >)[value] || '—'
  );
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '订单号 / 手机号 / 收货人 / 订单 ID',
    },
    fieldName: 'keyword',
    label: '订单搜索',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '精确或模糊订单号' },
    fieldName: 'order_sn',
    label: '订单号',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '商户 ID' },
    fieldName: 'mer_id',
    label: '商户 ID',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '收货人手机号' },
    fieldName: 'phone',
    label: '手机号',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '收货人姓名' },
    fieldName: 'real_name',
    label: '收货人',
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

const gridOptions: VxeGridProps<PlatformOrder> = {
  columns: [
    { field: 'order_id', title: '订单 ID', width: 94 },
    { field: 'order_sn', minWidth: 178, showOverflow: false, title: '订单号' },
    {
      field: 'mer_name',
      minWidth: 140,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => row.mer_name || `商户 #${row.mer_id}`,
    },
    {
      field: 'real_name',
      minWidth: 100,
      showOverflow: false,
      title: '收货人',
      formatter: ({ row }) => row.real_name || '—',
    },
    { field: 'user_phone', minWidth: 120, title: '联系电话' },
    { field: 'total_num', title: '商品数', width: 80 },
    {
      field: 'pay_price',
      title: '实付金额',
      width: 110,
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
      title: '订单状态',
      width: 120,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '下单时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    platformListActionColumn({ width: 88 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const keyword =
          String(formValues?.keyword ?? '').trim() ||
          String(formValues?.phone ?? '').trim() ||
          String(formValues?.real_name ?? '').trim() ||
          undefined;
        const data = await listPlatformOrdersApi({
          page: page.currentPage,
          limit: page.pageSize,
          status: -1,
          keyword,
          order_sn: String(formValues?.order_sn ?? '').trim() || undefined,
          mer_id: formValues?.mer_id ? Number(formValues.mer_id) : undefined,
          phone: String(formValues?.phone ?? '').trim() || undefined,
          real_name: String(formValues?.real_name ?? '').trim() || undefined,
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

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[900px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

async function openDetail(id: number) {
  detail.value = undefined;
  detailLoading.value = true;
  detailDrawerApi
    .setState({ title: '订单详情（只读监管）', loading: true })
    .open();
  try {
    detail.value = await getPlatformOrderApi(id);
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #status="{ row }">
        <ElTag type="info">{{ orderStatus(row.status) }}</ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row.order_id)">
          详情
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <template #default>
          <template v-if="detail">
            <ElDescriptions :column="2" border>
              <ElDescriptionsItem label="订单号">
                {{ detail.order_sn }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="商户">
                {{ detail.mer_name || `商户 #${detail.mer_id}` }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="支付方式">
                {{ payType(detail.pay_type) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="订单状态">
                {{ orderStatus(detail.status) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="支付金额">
                ¥{{ Number(detail.pay_price).toFixed(2) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="支付时间">
                {{ formatShanghaiDateTime(detail.pay_time) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="收货人">
                {{ detail.real_name || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="联系电话">
                {{ detail.user_phone || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem :span="2" label="物流信息">
                {{ detail.delivery_name || '未发货' }}
                {{ detail.delivery_id || '' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem :span="2" label="收货地址">
                {{ detail.user_address || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem :span="2" label="下单时间">
                {{ formatShanghaiDateTime(detail.create_time) }}
              </ElDescriptionsItem>
            </ElDescriptions>
            <div class="mb-2 mt-4 text-sm font-medium">商品明细</div>
            <ElTable :data="detail.products || []" border>
              <ElTableColumn
                label="商品"
                min-width="240"
                prop="product_info"
                show-overflow-tooltip
              />
              <ElTableColumn label="单价" width="110">
                <template #default="{ row }">
                  ¥{{ Number(row.product_price).toFixed(2) }}
                </template>
              </ElTableColumn>
              <ElTableColumn label="数量" prop="product_num" width="80" />
              <ElTableColumn label="小计" width="110">
                <template #default="{ row }">
                  ¥{{ Number(row.total_price).toFixed(2) }}
                </template>
              </ElTableColumn>
            </ElTable>
          </template>
        </template>
      </ElSkeleton>
    </DetailDrawer>
  </Page>
</template>
