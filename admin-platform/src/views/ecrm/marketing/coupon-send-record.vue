<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElAvatar,
  ElButton,
  ElPagination,
  ElTag,
} from 'element-plus';
import { CloseBold } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getPlatformCouponSendDetailApi,
  listPlatformCouponSendUsersApi,
  listPlatformCouponSendsApi,
  type PlatformCouponSendDetail,
  type PlatformCouponSendRecord,
  type PlatformCouponSendUser,
} from '#/api/core/platform-coupon-command';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);
const detail = ref<PlatformCouponSendDetail | null>(null);
const detailLoading = ref(false);
const usageLoading = ref(false);
const usageSendID = ref(0);
const usagePager = reactive({ page: 1, limit: 10, total: 0 });
const usageRows = ref<PlatformCouponSendUser[]>([]);

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const typeRaw = formValues?.coupon_type;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    coupon_type:
      typeRaw === 10 || typeRaw === 11 || typeRaw === 12
        ? Number(typeRaw)
        : undefined,
    coupon_name: String(formValues?.coupon_name ?? '').trim() || undefined,
  };
}

function formatTime(value?: string | null) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function formatMinPrice(value?: number) {
  const n = Number(value || 0);
  return n > 0 ? `最低消费${n}` : '无门槛';
}

function formatPublishCount(row: PlatformCouponSendDetail) {
  return row.is_limited === 1 ? String(row.total_count) : '不限量';
}

function formatRemainCount(row: PlatformCouponSendDetail) {
  return row.is_limited === 1 ? String(row.remain_count) : '不限量';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '平台通用券', value: 10 },
        { label: '平台品类券', value: 11 },
        { label: '平台跨店券', value: 12 },
      ],
      placeholder: '全部',
    },
    fieldName: 'coupon_type',
    label: '优惠券类型',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入优惠券名称',
    },
    fieldName: 'coupon_name',
    label: '优惠券名称',
  },
]);

const gridOptions: VxeGridProps<PlatformCouponSendRecord> = {
  columns: [
    { field: 'coupon_send_id', title: 'ID', width: 80 },
    { field: 'title', minWidth: 140, showOverflow: false, title: '优惠券名称' },
    {
      field: 'coupon_type_name',
      minWidth: 120,
      showOverflow: false,
      title: '优惠券类型',
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '发送日期',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    {
      field: 'validity_text',
      title: '使用有效期',
      width: 110,
    },
    {
      field: 'filter_text',
      minWidth: 140,
      showOverflow: false,
      title: '筛选条件',
    },
    {
      field: 'coupon_num',
      minWidth: 160,
      slots: { default: 'usage' },
      title: '使用情况',
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listPlatformCouponSendsApi(
          buildParams(page, formValues),
        );
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'coupon_send_id' },
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

const [UsageDrawer, usageDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
});

async function openDetail(row: PlatformCouponSendRecord) {
  detail.value = null;
  detailDrawerApi.setState({ title: '优惠券详情', loading: true }).open();
  detailLoading.value = true;
  try {
    detail.value = await getPlatformCouponSendDetailApi(row.coupon_send_id);
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}

async function loadUsageRows() {
  if (!usageSendID.value) return;
  usageLoading.value = true;
  try {
    const data = await listPlatformCouponSendUsersApi(usageSendID.value, {
      page: usagePager.page,
      limit: usagePager.limit,
    });
    usageRows.value = data.list || [];
    usagePager.total = data.total || 0;
  } finally {
    usageLoading.value = false;
  }
}

async function openUsage(row: PlatformCouponSendRecord) {
  usageSendID.value = row.coupon_send_id;
  usagePager.page = 1;
  usageRows.value = [];
  usageDrawerApi.setState({ title: '使用记录' }).open();
  await loadUsageRows();
}

function onUsagePageChange(page: number) {
  usagePager.page = page;
  void loadUsageRows();
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canRead.value =
    profile.roles.includes('platform') &&
    codes.includes('marketing.coupon.send.read');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      title="当前账号没有优惠券发送记录查看权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #usage="{ row }">
        <div class="send-usage">
          <div>发放数量：{{ row.coupon_num }}</div>
          <div class="send-usage__used">
            发放使用数量：{{ row.use_count }}
          </div>
        </div>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton link type="primary" @click="openUsage(row)">
          使用记录
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <div v-loading="detailLoading" class="send-detail">
        <template v-if="detail">
          <div class="send-detail__grid">
            <div class="send-detail__item">
              <span class="label">优惠券名称</span>
              <span class="value">{{ detail.title || '—' }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">优惠券类型</span>
              <span class="value">{{
                detail.coupon_type_name || '通用券'
              }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">优惠券面值</span>
              <span class="value">{{
                Number(detail.coupon_price || 0).toFixed(2)
              }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">使用门槛</span>
              <span class="value">{{
                formatMinPrice(detail.use_min_price)
              }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">使用有效期</span>
              <span class="value">{{ detail.validity_text || '—' }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">是否限量</span>
              <span class="value">
                <ElTag
                  v-if="detail.is_limited === 1"
                  type="success"
                  size="small"
                >
                  是
                </ElTag>
                <span v-else class="limit-no">
                  <CloseBold class="limit-no__icon" />
                </span>
              </span>
            </div>
            <div class="send-detail__item">
              <span class="label">领取时间</span>
              <span class="value">
                <template v-if="Number(detail.is_timeout || 0) === 1">
                  {{ formatTime(detail.start_time) }} ~
                  {{ formatTime(detail.end_time) }}
                </template>
                <template v-else>不限时</template>
              </span>
            </div>
            <div class="send-detail__item">
              <span class="label">剩余总数</span>
              <span class="value">{{ formatRemainCount(detail) }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">发送类型</span>
              <span class="value">{{
                detail.send_type_name || '后台赠送'
              }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">已使用总数</span>
              <span class="value">{{ detail.used_total ?? 0 }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">已发布总数</span>
              <span class="value">{{ formatPublishCount(detail) }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">排序</span>
              <span class="value">{{ detail.sort ?? 0 }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">已发送总数</span>
              <span class="value">{{ detail.sent_total ?? 0 }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">发放筛选条件</span>
              <span class="value">{{ detail.filter_text || '无' }}</span>
            </div>
            <div class="send-detail__item">
              <span class="label">状态</span>
              <span class="value">{{
                detail.status === 1 ? '开启' : '关闭'
              }}</span>
            </div>
          </div>
        </template>
      </div>
    </DetailDrawer>

    <UsageDrawer>
      <div v-loading="usageLoading" class="usage-drawer">
        <table class="usage-table">
          <thead>
            <tr>
              <th>用户名</th>
              <th>用户头像</th>
              <th>优惠券获取方式</th>
              <th>使用情况</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in usageRows" :key="`${row.user_id}-${row.status}`">
              <td>{{ row.nickname || '—' }}</td>
              <td>
                <ElAvatar :size="36" :src="row.avatar_url || undefined">
                  {{ (row.nickname || 'U').slice(0, 1) }}
                </ElAvatar>
              </td>
              <td>{{ row.source_name || '—' }}</td>
              <td>{{ row.status_name || '—' }}</td>
            </tr>
            <tr v-if="!usageRows.length">
              <td colspan="4" class="usage-empty">暂无使用记录</td>
            </tr>
          </tbody>
        </table>
        <div v-if="usagePager.total > 0" class="usage-pager">
          <ElPagination
            background
            layout="prev, pager, next, jumper"
            :current-page="usagePager.page"
            :page-size="usagePager.limit"
            :total="usagePager.total"
            @current-change="onUsagePageChange"
          />
        </div>
      </div>
    </UsageDrawer>
  </Page>
</template>

<style scoped>
.send-usage {
  line-height: 1.6;
  font-size: 13px;
}

.send-usage__used {
  color: hsl(var(--destructive));
}

.send-detail {
  min-height: 200px;
  padding: 8px 4px 16px;
}

.send-detail__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px 32px;
}

.send-detail__item {
  display: grid;
  grid-template-columns: 120px minmax(0, 1fr);
  gap: 8px;
  align-items: start;
  font-size: 13px;
  line-height: 1.6;
}

.send-detail__item .label {
  color: hsl(var(--muted-foreground));
}

.send-detail__item .value {
  color: hsl(var(--foreground));
  word-break: break-all;
}

.limit-no {
  display: inline-flex;
  color: hsl(var(--primary));
}

.limit-no__icon {
  width: 14px;
  height: 14px;
}

.usage-drawer {
  min-height: 200px;
}

.usage-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.usage-table th,
.usage-table td {
  padding: 12px 10px;
  border-bottom: 1px solid hsl(var(--border));
  text-align: left;
  vertical-align: middle;
}

.usage-table th {
  color: hsl(var(--muted-foreground));
  font-weight: 600;
  background: hsl(var(--muted) / 0.35);
}

.usage-empty {
  padding: 28px 0 !important;
  color: hsl(var(--muted-foreground));
  text-align: center !important;
}

.usage-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
