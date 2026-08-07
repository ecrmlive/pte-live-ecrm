<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref, watch } from 'vue';

import { Page } from '@vben/common-ui';
import { ElTabPane, ElTabs, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listMerchantCouponRecordsApi,
  listMerchantCouponSendsApi,
  type MerchantCouponRecord,
  type MerchantCouponSend,
} from '#/api/core/merchant-promotion';
import { MERCHANT_LIST_GRID_LAYOUT } from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

type RecordsTab = 'records' | 'sends';

const activeTab = ref<RecordsTab>('records');

function recordStatusInfo(status: number) {
  return (
    (
      {
        0: { label: '未使用', type: 'warning' },
        1: { label: '已使用', type: 'success' },
        2: { label: '已过期', type: 'info' },
      } as const
    )[status] || { label: '未知', type: 'info' as const }
  );
}

const recordFormOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '优惠券名称' },
    fieldName: 'keyword',
    label: '优惠券搜索',
  },
  {
    component: 'InputNumber',
    componentProps: { clearable: true, min: 1, placeholder: '优惠券 ID' },
    fieldName: 'coupon_id',
    label: '优惠券 ID',
  },
  {
    component: 'InputNumber',
    componentProps: { clearable: true, min: 1, placeholder: '用户 ID' },
    fieldName: 'uid',
    label: '用户 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '未使用', value: 0 },
        { label: '已使用', value: 1 },
        { label: '已过期', value: 2 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '使用状态',
  },
]);

const sendFormOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'InputNumber',
    componentProps: { clearable: true, min: 1, placeholder: '优惠券 ID' },
    fieldName: 'coupon_id',
    label: '优惠券 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '发送中', value: 0 },
        { label: '已完成', value: 1 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '发送状态',
  },
]);

const recordGridOptions: VxeGridProps<MerchantCouponRecord> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'coupon_user_id', title: '记录 ID', width: 96 },
    { field: 'coupon_title', minWidth: 150, showOverflow: false, title: '优惠券' },
    { field: 'uid', title: '用户 ID', width: 100 },
    {
      field: 'coupon_price',
      title: '面额',
      width: 106,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'type',
      title: '获取方式',
      width: 100,
      formatter: ({ cellValue }) =>
        cellValue === 'send' ? '后台发送' : '用户领取',
    },
    {
      field: 'status',
      slots: { default: 'recordStatus' },
      title: '状态',
      width: 100,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '领取时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    {
      field: 'use_time',
      minWidth: 170,
      title: '使用时间',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = formValues?.status;
        const couponId = formValues?.coupon_id;
        const uid = formValues?.uid;
        const data = await listMerchantCouponRecordsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          coupon_id:
            typeof couponId === 'number' && couponId > 0
              ? Number(couponId)
              : undefined,
          uid: typeof uid === 'number' && uid > 0 ? Number(uid) : undefined,
          status:
            status === 0 || status === 1 || status === 2
              ? Number(status)
              : undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'coupon_user_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const sendGridOptions: VxeGridProps<MerchantCouponSend> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'coupon_send_id', title: '批次 ID', width: 96 },
    { field: 'coupon_id', title: '优惠券 ID', width: 104 },
    { field: 'coupon_num', title: '发送数量', width: 104 },
    { field: 'mark', minWidth: 220, showOverflow: true, title: '发送说明' },
    {
      field: 'status',
      slots: { default: 'sendStatus' },
      title: '状态',
      width: 100,
    },
    {
      field: 'create_time',
      minWidth: 180,
      title: '发送时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = formValues?.status;
        const couponId = formValues?.coupon_id;
        const data = await listMerchantCouponSendsApi({
          page: page.currentPage,
          limit: page.pageSize,
          coupon_id:
            typeof couponId === 'number' && couponId > 0
              ? Number(couponId)
              : undefined,
          status: status === 0 || status === 1 ? Number(status) : undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'coupon_send_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [RecordGrid, recordGridApi] = useVbenVxeGrid({
  formOptions: recordFormOptions,
  gridOptions: recordGridOptions,
});

const [SendGrid, sendGridApi] = useVbenVxeGrid({
  formOptions: sendFormOptions,
  gridOptions: sendGridOptions,
});

watch(activeTab, (tab) => {
  if (tab === 'records') recordGridApi.reload();
  else sendGridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElTabs v-model="activeTab">
      <ElTabPane label="领取与使用记录" name="records">
        <RecordGrid>
          <template #recordStatus="{ row }">
            <ElTag :type="recordStatusInfo(row.status).type">
              {{ recordStatusInfo(row.status).label }}
            </ElTag>
          </template>
        </RecordGrid>
      </ElTabPane>
      <ElTabPane label="发送记录" name="sends">
        <SendGrid>
          <template #sendStatus="{ row }">
            <ElTag :type="row.status === 1 ? 'success' : 'warning'">
              {{ row.status === 1 ? '已完成' : '发送中' }}
            </ElTag>
          </template>
        </SendGrid>
      </ElTabPane>
    </ElTabs>
  </Page>
</template>
