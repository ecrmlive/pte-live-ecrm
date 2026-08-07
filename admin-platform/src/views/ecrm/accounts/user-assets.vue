<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElCard, ElCol, ElRow } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getPlatformUserAssetSummaryApi,
  listPlatformUserAssetsApi,
  type UserAssetLedgerRow,
  type UserAssetSummary,
  type UserAssetType,
} from '#/api/core/platform-user-assets';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);
const summary = ref<UserAssetSummary[]>([]);

const assetLabels: Record<UserAssetType, string> = {
  balance: '余额',
  commission: '佣金',
  points: '积分',
};

const summaryByType = computed(
  () => new Map(summary.value.map((item) => [item.asset_type, item])),
);

function formatAmount(row: UserAssetLedgerRow) {
  const prefix = row.amount > 0 ? '+' : '';
  return `${prefix}${row.amount.toFixed(2)}`;
}

function summaryText(type: UserAssetType) {
  const item = summaryByType.value.get(type);
  if (!item) return '暂无流水';
  return `入账 ${item.income.toFixed(2)} · 支出 ${item.expense.toFixed(2)} · ${item.count} 笔`;
}

function buildQueryParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const userIdRaw = String(formValues?.user_id ?? '').trim();
  return {
    page: page.currentPage,
    limit: page.pageSize,
    asset_type: (String(formValues?.asset_type ?? '').trim() ||
      undefined) as UserAssetType | undefined,
    user_id: userIdRaw ? Number(userIdRaw) : undefined,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    date_from: range[0],
    date_to: range[1],
  };
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
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(assetLabels).map(([value, label]) => ({
        label,
        value,
      })),
      placeholder: '全部类型',
    },
    fieldName: 'asset_type',
    label: '资产类型',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '业务单号 / 参考类型关键词',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
]);

const gridOptions: VxeGridProps<UserAssetLedgerRow> = {
  columns: [
    { field: 'id', title: '流水 ID', width: 100 },
    { field: 'user_id', title: '用户 ID', width: 100 },
    {
      field: 'asset_type',
      formatter: ({ cellValue }) =>
        assetLabels[cellValue as UserAssetType] || cellValue,
      title: '资产类型',
      width: 100,
    },
    {
      field: 'amount',
      slots: { default: 'amount' },
      title: '变动金额',
      width: 130,
    },
    {
      field: 'reference_type',
      minWidth: 150,
      showOverflow: false,
      title: '业务来源',
    },
    {
      field: 'reference_id',
      minWidth: 160,
      showOverflow: false,
      title: '业务引用',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '创建时间',
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const params = buildQueryParams(page, formValues);
        const [pageData, summaryData] = await Promise.all([
          listPlatformUserAssetsApi(params),
          getPlatformUserAssetSummaryApi({
            asset_type: params.asset_type,
            user_id: params.user_id,
            keyword: params.keyword,
            date_from: params.date_from,
            date_to: params.date_to,
          }),
        ]);
        summary.value = summaryData.list || [];
        return {
          items: pageData.list || [],
          total: pageData.total || 0,
        };
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
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canRead.value =
    profile.roles.includes('platform') &&
    permissions.includes('accounts.user_assets.read');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有查看用户资产流水的权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <ElRow :gutter="16" class="mb-4">
        <ElCol
          v-for="type in (['balance', 'points', 'commission'] as UserAssetType[])"
          :key="type"
          :md="8"
          :xs="24"
        >
          <ElCard shadow="never">
            <div class="text-sm text-gray-500">{{ assetLabels[type] }}</div>
            <div class="mt-2 text-sm">{{ summaryText(type) }}</div>
          </ElCard>
        </ElCol>
      </ElRow>
      <Grid>
        <template #amount="{ row }">
          <span :class="row.amount < 0 ? 'text-red-500' : 'text-green-600'">
            {{ formatAmount(row) }}
          </span>
        </template>
      </Grid>
    </template>
  </Page>
</template>
