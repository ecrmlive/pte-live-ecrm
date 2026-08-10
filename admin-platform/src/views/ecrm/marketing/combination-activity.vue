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
  getPlatformCombinationBuyingApi,
  listPlatformCombinationBuyingsApi,
  type PlatformCombinationBuying,
} from '#/api/core/platform-combination';
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
const detail = ref<PlatformCombinationBuying>();

function time(value?: string) {
  return formatShanghaiDateTime(value);
}

function stopTime(row: PlatformCombinationBuying) {
  if (row.stop_time) return row.stop_time;
  if (row.end_time > 0) {
    return formatShanghaiDateTime(new Date(row.end_time * 1000));
  }
  return '—';
}

function traderText(row: PlatformCombinationBuying) {
  return row.trader_name || (row.is_trader === 1 ? '自营' : '非自营');
}

function statusText(status: number) {
  if (status === 10) return '已完成';
  if (status === -1) return '已失败';
  return '未完成';
}

function statusTagType(status: number) {
  if (status === 10) return 'success';
  if (status === -1) return 'info';
  return 'warning';
}

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const values = formValues || {};
  const range = Array.isArray(values.date_range) ? values.date_range : [];
  const merRaw = values.mer_id;
  const traderRaw = values.is_trader;
  const statusRaw = values.status;
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
    status:
      statusRaw === 0 || statusRaw === 10 || statusRaw === -1
        ? Number(statusRaw)
        : undefined,
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
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '未完成', value: 0 },
        { label: '已完成', value: 10 },
        { label: '已失败', value: -1 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'status',
    label: '拼团状态',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '昵称/ID',
    },
    fieldName: 'user_name',
    label: '团长搜索',
  },
]);

const gridOptions: VxeGridProps<PlatformCombinationBuying> = {
  columns: [
    { field: 'group_buying_id', title: 'ID', width: 90 },
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
      field: 'nickname',
      minWidth: 110,
      showOverflow: false,
      title: '开团团长',
      formatter: ({ row }) => row.nickname || (row.uid ? `用户#${row.uid}` : '—'),
    },
    {
      field: 'image',
      slots: { default: 'image' },
      title: '拼团商品图片',
      width: 110,
    },
    {
      field: 'store_name',
      minWidth: 160,
      showOverflow: false,
      title: '拼团商品',
      formatter: ({ row }) =>
        row.store_name || (row.product_id ? `商品#${row.product_id}` : '—'),
    },
    {
      field: 'create_time',
      minWidth: 180,
      showOverflow: false,
      slots: { default: 'comboTime' },
      title: '拼团时间',
    },
    {
      field: 'buying_count_num',
      title: '几人团',
      width: 90,
    },
    {
      field: 'yet_buying_num',
      title: '参与人次',
      width: 100,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ width: 100 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformCombinationBuyingsApi(
          buildParams(page, formValues),
        );
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'group_buying_id' },
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

async function showDetail(row: PlatformCombinationBuying) {
  detail.value = await getPlatformCombinationBuyingApi(row.group_buying_id);
  detailDrawerApi.setState({ title: '拼团活动详情' }).open();
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
      <template #comboTime="{ row }">
        <div class="leading-5">
          <div>{{ time(row.create_time) }}</div>
          <div class="text-muted-foreground">{{ stopTime(row) }}</div>
        </div>
      </template>
      <template #status="{ row }">
        <ElTag :type="statusTagType(row.status)" size="small">
          {{ row.status_text || statusText(row.status) }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="showDetail(row)">
          查看详情
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="detail">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem label="开团 ID">
            {{ detail.group_buying_id }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="拼团商品 ID">
            {{ detail.product_group_id }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="店铺名称">
            {{ detail.mer_name || `店铺#${detail.mer_id}` }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="店铺类别">
            {{ traderText(detail) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="开团团长">
            {{
              detail.nickname ||
              (detail.uid ? `用户#${detail.uid}` : '—')
            }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="状态">
            <ElTag :type="statusTagType(detail.status)" size="small">
              {{ detail.status_text || statusText(detail.status) }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="拼团商品" :span="2">
            {{
              detail.store_name ||
              (detail.product_id ? `商品#${detail.product_id}` : '—')
            }}
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
          <ElDescriptionsItem label="拼团价">
            ¥{{ Number(detail.price || 0).toFixed(2) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="几人团">
            {{ detail.buying_count_num }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="参与人次">
            {{ detail.yet_buying_num ?? 0 }} / {{ detail.buying_count_num }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="发起时间">
            {{ time(detail.create_time) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="结束时间">
            {{ stopTime(detail) }}
          </ElDescriptionsItem>
        </ElDescriptions>

        <div class="mt-4 text-sm font-medium">参团成员</div>
        <div
          v-if="detail.members?.length"
          class="mt-2 grid grid-cols-1 gap-2 md:grid-cols-2"
        >
          <div
            v-for="member in detail.members"
            :key="member.id"
            class="flex items-center gap-3 rounded border px-3 py-2"
          >
            <ElImage
              v-if="member.avatar"
              :src="resolveCosMediaUrl(member.avatar)"
              class="h-9 w-9 rounded-full"
              fit="cover"
            />
            <div
              v-else
              class="bg-muted text-muted-foreground flex h-9 w-9 items-center justify-center rounded-full text-xs"
            >
              {{ member.is_initiator || member.is_leader ? '团' : '员' }}
            </div>
            <div class="min-w-0">
              <div class="truncate">
                {{ member.nickname || `用户#${member.uid}` }}
                <ElTag
                  v-if="member.is_initiator || member.is_leader"
                  class="ml-1"
                  size="small"
                  type="warning"
                >
                  团长
                </ElTag>
              </div>
              <div class="text-muted-foreground text-xs">
                {{ time(member.create_time) }}
              </div>
            </div>
          </div>
        </div>
        <div v-else class="text-muted-foreground mt-2 text-sm">暂无参团成员</div>
      </template>
    </DetailDrawer>
  </Page>
</template>
