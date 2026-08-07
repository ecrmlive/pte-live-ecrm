<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listCouponCommands,
  type CouponCommandRecord,
} from '#/api/core/platform-coupon-command';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const userIdRaw = String(formValues?.user_id ?? '').trim();
  const couponIdRaw = String(formValues?.coupon_id ?? '').trim();
  const actionRaw = formValues?.action;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    user_id: userIdRaw ? Number(userIdRaw) : undefined,
    coupon_id: couponIdRaw ? Number(couponIdRaw) : undefined,
    action:
      actionRaw === 'issue' || actionRaw === 'revoke'
        ? (actionRaw as CouponCommandRecord['action'])
        : undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

function actionLabel(action: CouponCommandRecord['action']) {
  return action === 'issue' ? '发放' : '撤销';
}

function stateText(value: string) {
  return value || '—';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '用户 ID' },
    fieldName: 'user_id',
    label: '用户 ID',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '优惠券 ID' },
    fieldName: 'coupon_id',
    label: '优惠券 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '发放', value: 'issue' },
        { label: '撤销', value: 'revoke' },
      ],
      placeholder: '全部操作',
    },
    fieldName: 'action',
    label: '操作',
  },
]);

const gridOptions: VxeGridProps<CouponCommandRecord> = {
  columns: [
    { field: 'id', title: '审计 ID', width: 100 },
    { field: 'user_id', title: '用户 ID', width: 100 },
    {
      field: 'coupon_name',
      minWidth: 220,
      showOverflow: false,
      title: '优惠券',
      formatter: ({ row }) => `${row.coupon_name}（#${row.coupon_id}）`,
    },
    {
      field: 'store_id',
      title: '归属',
      width: 110,
      formatter: ({ cellValue }) => (cellValue ? `店铺 #${cellValue}` : '平台券'),
    },
    {
      field: 'action',
      slots: { default: 'action' },
      title: '操作',
      width: 90,
    },
    {
      field: 'from_status',
      minWidth: 150,
      showOverflow: false,
      title: '状态变化',
      formatter: ({ row }) => `${stateText(row.from_status)} → ${stateText(row.to_status)}`,
    },
    { field: 'reason', minWidth: 240, showOverflow: false, title: '操作原因' },
    { field: 'operator_admin_id', title: '操作管理员 ID', width: 140 },
    {
      field: 'created_at',
      minWidth: 180,
      title: '操作时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listCouponCommands(buildParams(page, formValues));
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.includes('platform') && codes.includes('marketing.coupon.send.read');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page
    auto-content-height
    description="只读展示平台人工发券与撤销的不可变审计。订单锁券、核销和支付回调不在本页处理，且不会显示账号、手机号或幂等键。"
    title="优惠券发送记录"
  >
    <ElAlert
      v-if="!canRead"
      title="当前账号没有优惠券发送记录查看权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <Grid>
        <template #action="{ row }">
          <ElTag :type="row.action === 'issue' ? 'success' : 'warning'">
            {{ actionLabel(row.action) }}
          </ElTag>
        </template>
      </Grid>
      <ElAlert
        class="mt-4"
        type="info"
        :closable="false"
        title="审计记录不提供编辑、删除、补写或重放。若模板后来被删除，记录仍保留并标记为“已删除优惠券模板”。"
      />
    </template>
  </Page>
</template>
