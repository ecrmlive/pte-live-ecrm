<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElCol, ElRow, ElTabPane, ElTabs, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getDistributionSummaryApi,
  listDistributionCommissionsApi,
  listDistributionPromotersApi,
  type CommissionStatus,
  type DistributionCommission,
  type DistributionPromoter,
  type DistributionSummary,
} from '#/api/core/platform-spread';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  LIST_ENABLE_STATUS_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const activeTab = ref<'commissions' | 'promoters'>('promoters');
const canRead = ref(false);
const summary = ref<DistributionSummary>();

const commissionLabels: Record<CommissionStatus, string> = {
  available: '可结算',
  pending: '待结算',
  settled: '已结算',
  voided: '已作废',
};

const promoterFormOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '用户 ID' },
    fieldName: 'user_id',
    label: '用户 ID',
  },
  LIST_ENABLE_STATUS_FIELD('推广资格'),
]);

function buildPromoterParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const base = buildStandardListParams(page, formValues);
  const userIdRaw = String(formValues?.user_id ?? '').trim();
  return {
    page: base.page,
    limit: base.limit,
    status: base.status as 0 | 1 | undefined,
    user_id: userIdRaw ? Number(userIdRaw) : undefined,
    date_from: base.date_from,
    date_to: base.date_to,
  };
}

const promoterGridOptions: VxeGridProps<DistributionPromoter> = {
  columns: [
    { field: 'user_id', title: '用户 ID', width: 110 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '资格',
      width: 100,
    },
    { field: 'direct_user_count', title: '直推用户', width: 110 },
    {
      field: 'pending_commission',
      title: '待结算',
      width: 120,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
    },
    {
      field: 'available_commission',
      title: '可结算',
      width: 120,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
    },
    {
      field: 'settled_commission',
      title: '已结算',
      width: 120,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
    },
    {
      field: 'updated_at',
      minWidth: 180,
      title: '更新时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const [summaryData, promoterData] = await Promise.all([
          getDistributionSummaryApi(),
          listDistributionPromotersApi(buildPromoterParams(page, formValues)),
        ]);
        summary.value = summaryData;
        return { items: promoterData.list || [], total: promoterData.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'user_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const commissionFormOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '用户 ID' },
    fieldName: 'user_id',
    label: '用户 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(commissionLabels).map(([value, label]) => ({ label, value })),
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '佣金状态',
  },
]);

function buildCommissionParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const userIdRaw = String(formValues?.user_id ?? '').trim();
  const statusRaw = formValues?.status;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    user_id: userIdRaw ? Number(userIdRaw) : undefined,
    status:
      typeof statusRaw === 'string' && statusRaw in commissionLabels
        ? (statusRaw as CommissionStatus)
        : undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

const commissionGridOptions: VxeGridProps<DistributionCommission> = {
  columns: [
    { field: 'commission_id', title: '流水 ID', width: 110 },
    { field: 'user_id', title: '用户 ID', width: 110 },
    { field: 'order_id', title: '订单 ID', width: 110 },
    {
      field: 'amount',
      title: '佣金金额',
      width: 120,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
    },
    {
      field: 'status',
      slots: { default: 'commission_status' },
      title: '状态',
      width: 110,
    },
    {
      field: 'available_at',
      minWidth: 180,
      title: '可结算时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    {
      field: 'created_at',
      minWidth: 180,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listDistributionCommissionsApi(buildCommissionParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'commission_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [PromoterGrid, promoterGridApi] = useVbenVxeGrid({
  formOptions: promoterFormOptions,
  gridOptions: promoterGridOptions,
});

const [CommissionGrid, commissionGridApi] = useVbenVxeGrid({
  formOptions: commissionFormOptions,
  gridOptions: commissionGridOptions,
});

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    permissions.includes('marketing.spread.read');
  if (canRead.value) promoterGridApi.reload();
});
</script>

<template>
  <Page
    auto-content-height
    description="监管推广资格、直推关系和佣金状态；本页不修改佣金、提现或用户绑定关系，资金变更必须由业务域状态机处理。"
    title="分销监管"
  >
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有查看分销监管的权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <ElRow :gutter="16" class="mb-4">
        <ElCol :md="6" :xs="12">
          <div class="rounded border p-4">
            <div class="text-sm text-gray-500">推广员</div>
            <div class="mt-2 text-lg">
              {{ summary?.promoter_count || 0 }}（启用 {{ summary?.active_promoter_count || 0 }}）
            </div>
          </div>
        </ElCol>
        <ElCol :md="6" :xs="12">
          <div class="rounded border p-4">
            <div class="text-sm text-gray-500">待结算佣金</div>
            <div class="mt-2 text-lg">{{ (summary?.pending_commission || 0).toFixed(2) }}</div>
          </div>
        </ElCol>
        <ElCol :md="6" :xs="12">
          <div class="rounded border p-4">
            <div class="text-sm text-gray-500">可结算佣金</div>
            <div class="mt-2 text-lg">{{ (summary?.available_commission || 0).toFixed(2) }}</div>
          </div>
        </ElCol>
        <ElCol :md="6" :xs="12">
          <div class="rounded border p-4">
            <div class="text-sm text-gray-500">已结算佣金</div>
            <div class="mt-2 text-lg">{{ (summary?.settled_commission || 0).toFixed(2) }}</div>
          </div>
        </ElCol>
      </ElRow>

      <ElTabs
        v-model="activeTab"
        @tab-change="
          (name) => {
            if (name === 'commissions') commissionGridApi.reload();
          }
        "
      >
        <ElTabPane label="推广员" name="promoters">
          <PromoterGrid>
            <template #status="{ row }">
              <ElTag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '启用' : '停用' }}
              </ElTag>
            </template>
          </PromoterGrid>
        </ElTabPane>
        <ElTabPane label="佣金流水" name="commissions">
          <CommissionGrid>
            <template #commission_status="{ row }">
              <ElTag>{{ commissionLabels[row.status] }}</ElTag>
            </template>
          </CommissionGrid>
        </ElTabPane>
      </ElTabs>
    </template>
  </Page>
</template>
