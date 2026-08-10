<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElImage,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { fetchPlatformMerchants } from '#/api/core/ecrm';
import {
  getPlatformAssistSetApi,
  listPlatformAssistSetsApi,
  type PlatformAssistSet,
} from '#/api/core/platform-assist';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const merchantOptions = ref<{ label: string; value: number }[]>([]);
const detail = ref<PlatformAssistSet>();

function time(value?: string) {
  return formatShanghaiDateTime(value);
}

function traderText(row: PlatformAssistSet) {
  return row.trader_name || (row.is_trader === 1 ? '自营' : '非自营');
}

function statusText(status: number) {
  const map: Record<number, string> = {
    [-1]: '已失败',
    1: '进行中',
    10: '已满员',
    20: '已支付',
  };
  return map[status] || `状态 ${status}`;
}

function statusTagType(status: number) {
  if (status === 10) return 'success';
  if (status === 20) return 'warning';
  if (status === -1) return 'info';
  return 'primary';
}

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const values = formValues || {};
  const range = Array.isArray(values.date_range) ? values.date_range : [];
  const merRaw = values.mer_id;
  const traderRaw = values.is_trader;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    date_from: (range[0] as string | undefined) || undefined,
    date_to: (range[1] as string | undefined) || undefined,
    mer_id:
      merRaw === 0 || merRaw === undefined || merRaw === null || merRaw === ''
        ? undefined
        : Number(merRaw),
    is_trader:
      traderRaw === 0 || traderRaw === 1 ? Number(traderRaw) : undefined,
    keyword: String(values.keyword ?? '').trim() || undefined,
    user_name: String(values.user_name ?? '').trim() || undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      filterable: true,
      options: [],
      placeholder: '请选择',
    },
    fieldName: 'mer_id',
    label: '店铺名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '自营', value: 1 },
        { label: '非自营', value: 0 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'is_trader',
    label: '店铺类别',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '商品名称/ID',
    },
    fieldName: 'keyword',
    label: '商品搜索',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '发起人昵称',
    },
    fieldName: 'user_name',
    label: '发起人',
  },
]);

const gridOptions: VxeGridProps<PlatformAssistSet> = {
  columns: [
    { field: 'product_assist_set_id', title: 'ID', width: 90 },
    {
      field: 'mer_name',
      minWidth: 120,
      showOverflow: false,
      title: '店铺名称',
      formatter: ({ row }) => row.mer_name || `店铺#${row.mer_id}`,
    },
    {
      field: 'trader_name',
      title: '店铺类别',
      width: 90,
      formatter: ({ row }) => traderText(row),
    },
    {
      field: 'image',
      slots: { default: 'image' },
      title: '助力商品图片',
      width: 110,
    },
    {
      field: 'store_name',
      minWidth: 160,
      showOverflow: false,
      title: '商品名称',
      formatter: ({ row }) => row.store_name || `商品#${row.product_id}`,
    },
    {
      field: 'assist_price',
      title: '助力价格',
      width: 100,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'yet_assist_count',
      title: '助力人数',
      width: 100,
      formatter: ({ row }) =>
        `${row.yet_assist_count ?? 0}/${row.assist_count ?? 0}`,
    },
    {
      field: 'nickname',
      minWidth: 110,
      showOverflow: false,
      title: '发起人',
      formatter: ({ row }) => row.nickname || `用户#${row.uid}`,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '发起时间',
      formatter: ({ cellValue }) => time(cellValue),
    },
    {
      field: 'start_time',
      minWidth: 200,
      showOverflow: false,
      title: '活动时间',
      formatter: ({ row }) => `${time(row.start_time)} / ${time(row.end_time)}`,
    },
    platformListActionColumn({ width: 100 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformAssistSetsApi(
          buildParams(page, formValues),
        );
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_assist_set_id' },
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
  class: 'w-[1000px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

async function showDetail(row: PlatformAssistSet) {
  detail.value = await getPlatformAssistSetApi(row.product_assist_set_id);
  detailDrawerApi.setState({ title: '助力活动详情' }).open();
}

async function loadMerchants() {
  const data = await fetchPlatformMerchants({ page: 1, limit: 200 });
  merchantOptions.value = (data.list || []).map((row) => ({
    label: row.mer_name || `店铺#${row.mer_id}`,
    value: row.mer_id,
  }));
  await gridApi.formApi?.updateSchema?.([
    {
      fieldName: 'mer_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options: merchantOptions.value,
        placeholder: '请选择',
      },
    },
  ]);
}

onMounted(() => {
  void loadMerchants();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #image="{ row }">
        <ElImage
          v-if="row.image"
          :src="resolveCosMediaUrl(row.image)"
          class="h-14 w-14 rounded"
          fit="cover"
          :preview-src-list="[resolveCosMediaUrl(row.image)]"
          preview-teleported
        />
        <span v-else class="text-muted-foreground">—</span>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="showDetail(row)">详情</ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="detail">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="实例 ID">
            {{ detail.product_assist_set_id }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="活动商品 ID">
            {{ detail.product_assist_id }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="店铺名称">
            {{ detail.mer_name || `店铺#${detail.mer_id}` }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="店铺类别">
            {{ traderText(detail) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="商品名称" :span="2">
            {{ detail.store_name || `商品#${detail.product_id}` }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="商品图">
            <ElImage
              v-if="detail.image"
              :src="resolveCosMediaUrl(detail.image)"
              class="h-16 w-16 rounded"
              fit="cover"
              :preview-src-list="[resolveCosMediaUrl(detail.image)]"
              preview-teleported
            />
            <span v-else>—</span>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="助力价格">
            ¥{{ Number(detail.assist_price || 0).toFixed(2) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="助力人数">
            {{ detail.yet_assist_count ?? 0 }} / {{ detail.assist_count ?? 0 }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="状态">
            <ElTag :type="statusTagType(detail.status)" size="small">
              {{ statusText(detail.status) }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="发起人">
            {{ detail.nickname || `用户#${detail.uid}` }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="发起时间">
            {{ time(detail.create_time) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="活动时间" :span="2">
            {{ time(detail.start_time) }} 至 {{ time(detail.end_time) }}
          </ElDescriptionsItem>
        </ElDescriptions>

        <div class="mt-4 text-sm font-medium">助力好友</div>
        <div
          v-if="detail.helpers?.length"
          class="mt-2 grid grid-cols-1 gap-2 md:grid-cols-2"
        >
          <div
            v-for="helper in detail.helpers"
            :key="helper.product_assist_user_id"
            class="flex items-center gap-3 rounded border px-3 py-2"
          >
            <ElImage
              v-if="helper.avatar_img"
              :src="resolveCosMediaUrl(helper.avatar_img)"
              class="h-9 w-9 rounded-full"
              fit="cover"
            />
            <div
              v-else
              class="bg-muted text-muted-foreground flex h-9 w-9 items-center justify-center rounded-full text-xs"
            >
              友
            </div>
            <div class="min-w-0">
              <div class="truncate">{{ helper.nickname || `用户#${helper.uid}` }}</div>
              <div class="text-muted-foreground text-xs">
                {{ time(helper.create_time) }}
              </div>
            </div>
          </div>
        </div>
        <div v-else class="text-muted-foreground mt-2 text-sm">暂无助力好友</div>
      </template>
    </DetailDrawer>
  </Page>
</template>
