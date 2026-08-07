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
  getPlatformAssistApi,
  listPlatformAssistApi,
  updatePlatformAssistApi,
  type PlatformAssistActive,
} from '#/api/core/platform-assist';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD,
  LIST_MER_ID_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canManage = ref(false);
const detail = ref<PlatformAssistActive>();

function time(value: string) {
  return formatShanghaiDateTime(value);
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('活动 / 商品名称'),
  LIST_MER_ID_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '已上架', value: 1 },
        { label: '已下架', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '展示状态',
  },
]);

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const base = buildStandardListParams(page, formValues);
  const statusRaw = formValues?.status;
  return {
    ...base,
    status: undefined,
    is_show: statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
  };
}

const gridOptions: VxeGridProps<PlatformAssistActive> = {
  columns: [
    { field: 'product_assist_id', title: '活动 ID', width: 96 },
    {
      field: 'store_name',
      minWidth: 180,
      showOverflow: false,
      title: '活动 / 商品',
      formatter: ({ row }) => row.store_name || `商品 #${row.product_id}`,
    },
    {
      field: 'mer_name',
      minWidth: 130,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => row.mer_name || `商户 #${row.mer_id}`,
    },
    {
      field: 'assist_price',
      title: '助力价',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'assist_count',
      minWidth: 130,
      showOverflow: false,
      title: '助力规则',
      formatter: ({ row }) => `${row.assist_count} 人 / 每人最多 ${row.assist_user_count} 次`,
    },
    { field: 'stock', title: '活动库存', width: 100 },
    {
      field: 'start_time',
      minWidth: 240,
      showOverflow: false,
      title: '活动时间',
      formatter: ({ row }) => `${time(row.start_time)} 至 ${time(row.end_time)}`,
    },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '展示状态',
      width: 100,
    },
    platformListActionColumn({ width: 146 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformAssistApi(buildParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_assist_id' },
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
  class: 'w-[640px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

async function showDetail(row: PlatformAssistActive) {
  detail.value = await getPlatformAssistApi(row.product_assist_id);
  detailDrawerApi.setState({ title: '好友助力活动详情' }).open();
}

async function setVisible(row: PlatformAssistActive, isShow: number) {
  const action = isShow === 1 ? '上架' : '下架';
  try {
    await ElMessageBox.confirm(
      `确认${action}好友助力活动“${row.store_name || `#${row.product_assist_id}`}”吗？`,
      `${action}确认`,
      { cancelButtonText: '取消', confirmButtonText: `确认${action}`, type: 'warning' },
    );
    await updatePlatformAssistApi(row.product_assist_id, { is_show: isShow });
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
    permissions.includes('marketing.assist.manage');
});
</script>

<template>
  <Page
    auto-content-height
    description="查看各商户好友助力活动及完整规则；具备运营权限可上架或下架。删除、改价、库存与时间调整须完成订单影响审计后另行开放。"
    title="好友助力监管"
  >
    <Grid>
      <template #is_show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '已上架' : '已下架' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="showDetail(row)">详情</ElButton>
        <ElButton
          v-if="canManage"
          link
          :type="row.is_show === 1 ? 'danger' : 'success'"
          @click="setVisible(row, row.is_show === 1 ? 0 : 1)"
        >
          {{ row.is_show === 1 ? '下架' : '上架' }}
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
          <ElDescriptionsItem label="助力价">
            ¥{{ Number(detail.assist_price).toFixed(2) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="所需助力">{{ detail.assist_count }} 人</ElDescriptionsItem>
          <ElDescriptionsItem label="单人助力次数">
            最多 {{ detail.assist_user_count }} 次
          </ElDescriptionsItem>
          <ElDescriptionsItem label="活动库存">{{ detail.stock }}</ElDescriptionsItem>
          <ElDescriptionsItem label="活动时间" :span="2">
            {{ time(detail.start_time) }} 至 {{ time(detail.end_time) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="前台展示">
            {{ detail.is_show === 1 ? '上架' : '下架' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="活动状态">
            {{ detail.status === 1 ? '启用' : '停用' }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <ElAlert
          class="mt-4"
          type="warning"
          :closable="false"
          title="详情不展示参与用户资料。已发起助力单的价格、库存与完成条件须保持订单快照，本页不提供这些字段编辑。"
        />
      </template>
    </DetailDrawer>
  </Page>
</template>
