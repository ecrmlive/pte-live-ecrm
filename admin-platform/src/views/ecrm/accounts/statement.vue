<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformUserAssetSummaryApi,
  type UserAssetSummary,
  type UserAssetType,
} from '#/api/core/platform-user-assets';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canRead = ref(false);

const assetLabels: Record<UserAssetType, string> = {
  balance: '余额',
  commission: '佣金',
  points: '积分',
};

function buildSummaryParams(formValues?: Record<string, unknown>) {
  const userIdRaw = String(formValues?.user_id ?? '').trim();
  return {
    asset_type: (String(formValues?.asset_type ?? '').trim() ||
      undefined) as UserAssetType | undefined,
    user_id: userIdRaw ? Number(userIdRaw) : undefined,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
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

const gridOptions: VxeGridProps<UserAssetSummary> = {
  columns: [
    {
      field: 'asset_type',
      formatter: ({ cellValue }) =>
        assetLabels[cellValue as UserAssetType] || cellValue,
      minWidth: 120,
      title: '资产类型',
    },
    {
      field: 'income',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      title: '入账',
      width: 140,
    },
    {
      field: 'expense',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      title: '支出',
      width: 140,
    },
    { field: 'count', title: '笔数', width: 100 },
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await getPlatformUserAssetSummaryApi(
          buildSummaryParams(formValues),
        );
        const list = data.list || [];
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'asset_type' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canRead.value = permissions.includes('accounts.user_assets.read');
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      class="mb-4"
      title="不含用户姓名、手机号、收款账户或外部支付凭据。"
      type="warning"
      :closable="false"
    />
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号无查看资产汇总权限。"
      type="info"
      :closable="false"
    />
    <Grid v-else />
  </Page>
</template>
