<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { ElLink } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import ShopApi from '#/api/core/shop';
import LiveTrafficApi from '#/api/core/live-traffic';
import { Page } from '@vben/common-ui';
import { parseApiListPage } from '#/utils/list-response';

interface TrafficRow {
  app_id: number;
  app_name: string;
  lvb_play_used_gb: null | number;
  platform_operator_name?: string;
  remain_gb: null | number;
  total_gb: null | number;
  user_name: string;
  vod_play_used_gb: null | number;
}

const router = useRouter();
const rows = ref<TrafficRow[]>([]);
const listLoaded = ref(false);

function filterRows(list: TrafficRow[], formValues?: Record<string, unknown>) {
  const idKw = String(formValues?.appId ?? '').trim();
  const nameKw = String(formValues?.keyword ?? '')
    .trim()
    .toLowerCase();
  return list.filter((r) => {
    if (idKw && !String(r.app_id).includes(idKw)) {
      return false;
    }
    if (!nameKw) {
      return true;
    }
    return (
      (r.app_name && r.app_name.toLowerCase().includes(nameKw)) ||
      (r.user_name && r.user_name.toLowerCase().includes(nameKw)) ||
      (r.platform_operator_name &&
        r.platform_operator_name.toLowerCase().includes(nameKw))
    );
  });
}

function fmtGB(v: null | number | undefined) {
  if (v == null) return '—';
  const n = Number(v || 0);
  return n.toFixed(n >= 100 ? 1 : 2);
}

async function loadData() {
  const res = await ShopApi.shopList({ list_rows: 200, page: 1 });
  const shops = parseApiListPage<{
    app_id: number;
    app_name: string;
    platform_operator_name?: string;
    user_name: string;
  }>(res.data).list;

  rows.value = await Promise.all(
    shops.map(async (shop) => {
      try {
        const accRes = await LiveTrafficApi.account({ app_id: shop.app_id });
        const acc = accRes.data || {};
        return {
          app_id: shop.app_id,
          app_name: shop.app_name,
          platform_operator_name: shop.platform_operator_name,
          user_name: shop.user_name,
          total_gb: acc.total_gb ?? null,
          lvb_play_used_gb: acc.lvb_play_used_gb ?? null,
          vod_play_used_gb: acc.vod_play_used_gb ?? null,
          remain_gb: acc.remain_gb ?? null,
        };
      } catch {
        return {
          app_id: shop.app_id,
          app_name: shop.app_name,
          platform_operator_name: shop.platform_operator_name,
          user_name: shop.user_name,
          total_gb: null,
          lvb_play_used_gb: null,
          vod_play_used_gb: null,
          remain_gb: null,
        };
      }
    }),
  );
}

const formOptions: VbenFormProps = {
  showCollapseButton: false,
  schema: [
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '筛选 ID' },
      fieldName: 'appId',
      formItemClass: 'pb-0',
      label: '商城ID',
    },
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '筛选名称' },
      fieldName: 'keyword',
      formItemClass: 'pb-0',
      label: '商城名称',
    },
  ],
};

const gridOptions = {
  border: true,
  columns: [
    { field: 'app_id', title: '商城ID', width: 90 },
    { field: 'app_name', minWidth: 140, title: '商城名称' },
    {
      field: 'total_gb',
      slots: { default: 'totalGb' },
      title: '充值(GB)',
      width: 100,
    },
    {
      field: 'lvb_play_used_gb',
      slots: { default: 'lvbGb' },
      title: 'LVB已用(GB)',
      width: 110,
    },
    {
      field: 'vod_play_used_gb',
      slots: { default: 'vodGb' },
      title: 'VOD已用(GB)',
      width: 110,
    },
    {
      field: 'remain_gb',
      slots: { default: 'remainGb' },
      title: '剩余(GB)',
      width: 100,
    },
    { field: 'user_name', title: '超管账号', width: 120 },
    {
      field: 'platform_operator_name',
      formatter: ({ cellValue }) => String(cellValue ?? '') || '—',
      title: '操作员',
      width: 100,
    },
    {
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 100,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 15, pageSizes: [10, 15, 20, 30, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!listLoaded.value) {
          try {
            await loadData();
          } catch {
            rows.value = [];
          }
          listLoaded.value = true;
        }
        const filtered = filterRows(rows.value, formValues);
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: filtered.slice(start, start + page.pageSize),
          total: filtered.length,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'app_id' },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });

function openAccount(appId: number) {
  router.push({ path: '/shop/traffic/account', query: { app_id: String(appId) } });
}
</script>

<template>
  <Page>
    <Grid>
      <template #totalGb="{ row }">
        {{ fmtGB(row.total_gb) }}
      </template>

      <template #lvbGb="{ row }">
        {{ fmtGB(row.lvb_play_used_gb) }}
      </template>

      <template #vodGb="{ row }">
        {{ fmtGB(row.vod_play_used_gb) }}
      </template>

      <template #remainGb="{ row }">
        <span :class="Number(row.remain_gb) < 0 ? 'font-semibold text-[#e6a23c]' : ''">
          {{ fmtGB(row.remain_gb) }}
        </span>
      </template>

      <template #action="{ row }">
        <ElLink type="primary" @click="openAccount(row.app_id)">管理</ElLink>
      </template>
    </Grid>
  </Page>
</template>
