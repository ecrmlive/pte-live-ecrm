<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElMessage,
  ElMessageBox,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getPlatformPresellApi,
  listPlatformPresellsApi,
  updatePlatformPresellApi,
  type PlatformPresell,
} from '#/api/core/platform-presell';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  LIST_MER_ID_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canManage = ref(false);
const detail = ref<PlatformPresell>();

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('预售商品名称'),
  LIST_MER_ID_FIELD,
  LIST_ENABLE_STATUS_FIELD(),
]);

const gridOptions: VxeGridProps<PlatformPresell> = {
  columns: [
    { field: 'product_presell_id', title: 'ID', width: 80 },
    { field: 'store_name', minWidth: 180, showOverflow: false, title: '预售商品' },
    {
      field: 'mer_name',
      minWidth: 130,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => row.mer_name || `商户 #${row.mer_id}`,
    },
    {
      field: 'price',
      title: '预售价',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'start_time',
      minWidth: 220,
      showOverflow: false,
      title: '时间',
      formatter: ({ row }) => `${row.start_time} 至 ${row.end_time}`,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 88,
    },
    platformListActionColumn({ width: 146 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformPresellsApi(buildStandardListParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_presell_id' },
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
  class: 'w-[660px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

async function showDetail(row: PlatformPresell) {
  detail.value = await getPlatformPresellApi(row.product_presell_id);
  detailDrawerApi.setState({ title: '预售活动详情' }).open();
}

async function setStatus(row: PlatformPresell, status: number) {
  const action = status === 1 ? '启用' : '停用';
  try {
    await ElMessageBox.confirm(`确认${action}预售活动“${row.store_name}”？`, `${action}确认`, {
      cancelButtonText: '取消',
      confirmButtonText: `确认${action}`,
      type: 'warning',
    });
    await updatePlatformPresellApi(row.product_presell_id, { status });
    ElMessage.success(`活动已${action}`);
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canManage.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    permissions.includes('marketing.presell.manage');
});
</script>

<template>
  <Page
    auto-content-height
    description="监管各商户预售活动；可查看完整价格、库存、发货与尾款时间窗。具备运营活动权限的账号可启停，资金字段调整待订单快照影响审计闭环后开放。"
    title="预售监管"
  >
    <Grid>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="showDetail(row)">详情</ElButton>
        <ElButton
          v-if="canManage"
          link
          :type="row.status === 1 ? 'danger' : 'success'"
          @click="setStatus(row, row.status === 1 ? 0 : 1)"
        >
          {{ row.status === 1 ? '停用' : '启用' }}
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="detail">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="活动名称" :span="2">
            {{ detail.store_name }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="商品 / 商户">
            #{{ detail.product_id }} / {{ detail.mer_name || `商户 #${detail.mer_id}` }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="活动类型">
            {{ detail.presell_type === 2 ? '定金预售' : '全款预售' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="预售价">
            ¥{{ Number(detail.price).toFixed(2) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="定金 / 尾款">
            {{
              detail.presell_type === 2
                ? `¥${Number(detail.down_price || 0).toFixed(2)} / ¥${Number(detail.final_price || 0).toFixed(2)}`
                : '不适用'
            }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="活动时间" :span="2">
            {{ detail.start_time }} 至 {{ detail.end_time }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="尾款时间" :span="2">
            {{
              detail.presell_type === 2
                ? `${detail.final_start_time || '—'} 至 ${detail.final_end_time || '—'}`
                : '不适用'
            }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="库存 / 已售">
            {{ detail.stock || 0 }} / {{ detail.seles || 0 }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="发货">
            {{ detail.delivery_type === 2 ? '发货后' : '付款后'
            }}{{ detail.delivery_day ? ` ${detail.delivery_day} 天` : '' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="前台展示">
            {{ detail.is_show === 1 ? '上架' : '下架' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="活动状态">
            {{ detail.status === 1 ? '启用' : '停用' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="活动说明" :span="2">
            {{ detail.store_info || '—' }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <ElAlert
          class="mt-4"
          type="warning"
          :closable="false"
          title="已产生定金或尾款订单的金额、时间窗须保持订单快照。本页不提供资金字段编辑，待订单影响审计闭环后再开放。"
        />
      </template>
    </DetailDrawer>
  </Page>
</template>
