<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import {
  ElAlert,
  ElButton,
  ElImage,
  ElMessage,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getPlatformDiscountApi,
  listPlatformDiscountsApi,
  setPlatformDiscountStatusApi,
  type PlatformDiscount,
  type PlatformDiscountProduct,
} from '#/api/core/platform-discount';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canRead = ref(false);
const canManage = ref(false);
const detail = ref<PlatformDiscount>();
const detailLoading = ref(false);
const statusBusy = ref<Record<number, boolean>>({});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '搭配套餐', value: 1 },
        { label: '固定套餐', value: 0 },
      ],
      placeholder: '请选择套餐类型',
    },
    fieldName: 'package_type',
    label: '套餐类型',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '全部', value: 'all' },
        { label: '上架', value: 1 },
        { label: '下架', value: 0 },
      ],
      placeholder: '全部',
    },
    defaultValue: 'all',
    fieldName: 'status',
    label: '套餐状态',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入套餐名称',
    },
    fieldName: 'keyword',
    label: '套餐搜索',
  },
]);

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const values = formValues || {};
  const statusRaw = values.status;
  const status =
    statusRaw === 'all' || statusRaw === '' || statusRaw === undefined
      ? undefined
      : statusRaw;
  const typeRaw = values.package_type;
  const packageType =
    typeRaw === '' || typeRaw === undefined || typeRaw === null
      ? undefined
      : typeRaw;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(values.keyword ?? '').trim() || undefined,
    status,
    package_type: packageType,
  };
}

const gridOptions: VxeGridProps<PlatformDiscount> = {
  columns: [
    { field: 'activity_id', title: 'ID', width: 80 },
    {
      field: 'name',
      minWidth: 180,
      showOverflow: false,
      title: '套餐名称',
    },
    {
      field: 'package_type_label',
      title: '套餐类型',
      width: 110,
      formatter: ({ row }) =>
        row.package_type_label ||
        (row.package_type === 1 ? '搭配套餐' : '固定套餐'),
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '套餐状态',
      width: 100,
    },
    {
      field: 'time_label',
      minWidth: 220,
      showOverflow: 'tooltip',
      title: '限时',
      formatter: ({ row }) => formatTimeRange(row),
    },
    {
      field: 'store_name',
      minWidth: 140,
      showOverflow: false,
      title: '店铺名称',
      formatter: ({ row }) =>
        row.store_name || (row.store_id ? `店铺#${row.store_id}` : '—'),
    },
    {
      field: 'create_time',
      minWidth: 170,
      showOverflow: false,
      title: '创建时间',
      formatter: ({ row }) =>
        formatShanghaiDateTime(row.create_time || row.created_at) || '—',
    },
    {
      field: 'remain_label',
      title: '剩余数量',
      width: 100,
      formatter: ({ row }) =>
        row.remain_label ||
        (row.is_limit === 1 ? String(row.limit_num ?? 0) : '不限量'),
    },
    platformListActionColumn({ width: 90 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listPlatformDiscountsApi(
          buildParams(page, formValues),
        );
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

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
});

const mainProducts = computed(
  () => detail.value?.main_products || ([] as PlatformDiscountProduct[]),
);
const comboProducts = computed(
  () => detail.value?.combo_products || ([] as PlatformDiscountProduct[]),
);
const showComboSection = computed(
  () => Number(detail.value?.package_type) === 1,
);

function formatTimeRange(row: PlatformDiscount) {
  if (row.is_time !== 1) return '不限时';
  if (row.time_label && row.time_label !== '不限时') return row.time_label;
  const start = row.starts_at ? formatShanghaiDateTime(row.starts_at) : '';
  const end = row.ends_at ? formatShanghaiDateTime(row.ends_at) : '';
  if (start || end) return `${start || '—'} ~ ${end || '—'}`;
  return '不限时';
}

function formatActivityTime(raw?: string, unlimited?: boolean) {
  if (unlimited) return '不限时';
  if (!raw) return '不限时';
  return formatShanghaiDateTime(raw) || '不限时';
}

function productImage(p?: PlatformDiscountProduct) {
  return resolveCosMediaUrl(String(p?.image || '').trim());
}

function qtyLabel(row?: PlatformDiscount) {
  if (!row) return '—';
  return (
    row.qty_label ||
    (row.is_limit === 1 ? String(row.limit_num ?? 0) : '不限量')
  );
}

async function openDetail(row: PlatformDiscount) {
  detail.value = undefined;
  detailLoading.value = true;
  detailDrawerApi
    .setState({ title: row.name || '套餐详情', loading: true })
    .open();
  try {
    detail.value = await getPlatformDiscountApi(row.activity_id);
    detailDrawerApi.setState({
      title: detail.value?.name || row.name || '套餐详情',
    });
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}

async function onStatusChange(row: PlatformDiscount, value: string | number | boolean) {
  const next = value ? 1 : 0;
  const prev = row.status === 1 ? 1 : 0;
  if (next === prev) return;
  const action = next === 1 ? '上架' : '下架';
  try {
    await confirm({
      title: '提示',
      content: `确认${action}套餐「${row.name}」？`,
      icon: 'warning',
    });
  } catch {
    row.status = prev;
    return;
  }
  statusBusy.value[row.activity_id] = true;
  try {
    await setPlatformDiscountStatusApi(row.activity_id, next as 0 | 1);
    row.status = next;
    ElMessage.success(`已${action}`);
  } catch {
    row.status = prev;
  } finally {
    statusBusy.value[row.activity_id] = false;
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const roleOK = profile.roles.some(
    (role) => role === 'platform' || role === 'operations',
  );
  canRead.value =
    roleOK && permissions.includes('marketing.discounts.read');
  canManage.value =
    roleOK && permissions.includes('marketing.discounts.manage');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有优惠套餐查看权限（marketing.discounts.read）"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage || !!statusBusy[row.activity_id]"
          @change="(v) => onStatusChange(row, v)"
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">查看</ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="detail">
        <div class="discount-detail">
          <div class="discount-detail__title">
            <IconifyIcon icon="lucide:shopping-bag" class="discount-detail__icon" />
            <span>{{ detail.name }}</span>
            <ElTag class="ml-2" size="small" type="info">
              {{
                detail.package_type_label ||
                (detail.package_type === 1 ? '搭配套餐' : '固定套餐')
              }}
            </ElTag>
          </div>

          <div class="discount-detail__meta">
            <div class="meta-item">
              <div class="meta-label">套餐数量</div>
              <div class="meta-value">{{ qtyLabel(detail) }}</div>
            </div>
            <div class="meta-item">
              <div class="meta-label">显示状态</div>
              <div class="meta-value">
                <ElTag
                  :type="detail.status === 1 ? 'success' : 'info'"
                  size="small"
                >
                  {{ detail.status_label || (detail.status === 1 ? '上架' : '下架') }}
                </ElTag>
              </div>
            </div>
            <div class="meta-item">
              <div class="meta-label">活动开始时间</div>
              <div class="meta-value">
                {{ formatActivityTime(detail.starts_at, detail.is_time !== 1) }}
              </div>
            </div>
            <div class="meta-item">
              <div class="meta-label">活动结束时间</div>
              <div class="meta-value">
                {{ formatActivityTime(detail.ends_at, detail.is_time !== 1) }}
              </div>
            </div>
            <div class="meta-item">
              <div class="meta-label">创建时间</div>
              <div class="meta-value">
                {{
                  formatShanghaiDateTime(
                    detail.create_time || detail.created_at,
                  ) || '—'
                }}
              </div>
            </div>
          </div>

          <div class="section-bar">套餐主商品</div>
          <ElTable
            :data="mainProducts"
            border
            class="mb-4 w-full"
            empty-text="暂无主商品"
          >
            <ElTableColumn label="商品名称" min-width="280">
              <template #default="{ row }">
                <div class="product-cell">
                  <ElImage
                    v-if="productImage(row)"
                    :src="productImage(row)"
                    class="product-cell__img"
                    fit="cover"
                    :preview-src-list="[productImage(row)]"
                    preview-teleported
                  />
                  <div
                    v-else
                    class="product-cell__placeholder"
                  >
                    —
                  </div>
                  <span class="product-cell__name">
                    {{ row.store_name || `商品#${row.product_id}` }}
                  </span>
                </div>
              </template>
            </ElTableColumn>
            <ElTableColumn label="参与规格" min-width="140">
              <template #default="{ row }">
                {{ row.spec || '—' }}
              </template>
            </ElTableColumn>
          </ElTable>

          <template v-if="showComboSection">
            <div class="section-bar">套餐搭配商品</div>
            <ElTable
              :data="comboProducts"
              border
              class="w-full"
              empty-text="暂无搭配商品"
            >
              <ElTableColumn label="商品名称" min-width="280">
                <template #default="{ row }">
                  <div class="product-cell">
                    <ElImage
                      v-if="productImage(row)"
                      :src="productImage(row)"
                      class="product-cell__img"
                      fit="cover"
                      :preview-src-list="[productImage(row)]"
                      preview-teleported
                    />
                    <div
                      v-else
                      class="product-cell__placeholder"
                    >
                      —
                    </div>
                    <span class="product-cell__name">
                      {{ row.store_name || `商品#${row.product_id}` }}
                    </span>
                  </div>
                </template>
              </ElTableColumn>
              <ElTableColumn label="参与规格" min-width="140">
                <template #default="{ row }">
                  {{ row.spec || '—' }}
                </template>
              </ElTableColumn>
            </ElTable>
          </template>
        </div>
      </template>
      <div v-else-if="detailLoading" class="text-muted-foreground p-4 text-sm">
        加载中…
      </div>
    </DetailDrawer>
  </Page>
</template>

<style scoped>
.discount-detail__title {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  font-size: 18px;
  font-weight: 600;
  line-height: 1.4;
}

.discount-detail__icon {
  margin-right: 8px;
  font-size: 20px;
  color: hsl(var(--primary));
}

.discount-detail__meta {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px 16px;
  padding: 14px 16px;
  margin-bottom: 16px;
  background: hsl(var(--muted) / 0.45);
  border-radius: 8px;
}

.meta-label {
  margin-bottom: 4px;
  font-size: 12px;
  color: hsl(var(--muted-foreground));
}

.meta-value {
  font-size: 14px;
  word-break: break-all;
}

.section-bar {
  padding: 8px 12px;
  margin-bottom: 10px;
  font-size: 14px;
  font-weight: 600;
  color: #1d4ed8;
  background: #eff6ff;
  border-left: 3px solid #3b82f6;
}

.product-cell {
  display: flex;
  gap: 10px;
  align-items: center;
}

.product-cell__img,
.product-cell__placeholder {
  width: 48px;
  height: 48px;
  flex-shrink: 0;
  border-radius: 4px;
}

.product-cell__placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  color: hsl(var(--muted-foreground));
  background: hsl(var(--muted));
}

.product-cell__name {
  min-width: 0;
  line-height: 1.4;
  word-break: break-all;
}
</style>
