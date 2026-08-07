<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElButton, ElMessage, ElMessageBox, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformDiscountsApi,
  setPlatformDiscountStatusApi,
  type PlatformDiscount,
} from '#/api/core/platform-discount';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);
const canManage = ref(false);

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('套餐名称'),
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '店铺 ID' },
    fieldName: 'store_id',
    label: '店铺 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '上架', value: 1 },
        { label: '下架', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const base = buildStandardListParams(page, formValues);
  const storeIdRaw = String(formValues?.store_id ?? '').trim();
  return {
    page: base.page,
    limit: base.limit,
    keyword: base.keyword,
    status: base.status,
    date_from: base.date_from,
    date_to: base.date_to,
    store_id: storeIdRaw ? Number(storeIdRaw) : undefined,
  };
}

const gridOptions: VxeGridProps<PlatformDiscount> = {
  columns: [
    { field: 'activity_id', title: '活动 ID', width: 100 },
    { field: 'store_id', title: '店铺 ID', width: 100 },
    { field: 'name', minWidth: 180, showOverflow: false, title: '名称' },
    {
      field: 'package_price',
      title: '套餐价',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'product_ids',
      title: '商品数',
      width: 90,
      formatter: ({ row }) => (row.product_ids || []).length,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    {
      field: 'starts_at',
      minWidth: 220,
      showOverflow: false,
      title: '有效期',
      formatter: ({ row }) => `${row.starts_at || '—'} ~ ${row.ends_at || '—'}`,
    },
    platformListActionColumn({ width: 160 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listPlatformDiscountsApi(buildParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'activity_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function setStatus(row: PlatformDiscount, status: 0 | 1) {
  const action = status === 1 ? '上架投影' : '下架投影';
  try {
    await ElMessageBox.confirm(
      `确认对「${row.name}」执行${action}？仅更新业务投影，不直连商户库。`,
      action,
      { type: 'warning' },
    );
    await setPlatformDiscountStatusApi(row.activity_id, status);
    ElMessage.success(`已${action}`);
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  const roleOK = profile.roles.some((role) => role === 'platform' || role === 'operations');
  canRead.value = roleOK && permissions.includes('marketing.discounts.read');
  canManage.value = roleOK && permissions.includes('marketing.discounts.manage');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page
    auto-content-height
    description="只读监管店铺优惠套餐投影（qixi_crm_b_marketing_activity_view）；创建与编辑由商户后台完成。"
    title="优惠套餐"
  >
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有优惠套餐监管权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '上架' : '下架' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton
            v-if="row.status !== 1"
            link
            type="success"
            @click="setStatus(row, 1)"
          >
            上架投影
          </ElButton>
          <ElButton v-else link type="warning" @click="setStatus(row, 0)">
            下架投影
          </ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>
  </Page>
</template>
