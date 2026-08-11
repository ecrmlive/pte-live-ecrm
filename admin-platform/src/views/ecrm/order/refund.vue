<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import {
  ElButton,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElPagination,
  ElRadio,
  ElRadioGroup,
  ElSkeleton,
  ElTabPane,
  ElTable,
  ElTableColumn,
  ElTabs,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  approvePlatformRefundApi,
  exportPlatformRefundsApi,
  getPlatformRefundApi,
  getPlatformRefundTabCountsApi,
  listPlatformRefundEventsApi,
  listPlatformRefundsApi,
  rejectPlatformRefundApi,
  type PlatformRefundEvent,
  type PlatformRefundListParams,
  type PlatformRefundOrder,
  type PlatformRefundTabCounts,
} from '#/api/core/platform-aftersale';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import {
  listUserSearchFormField,
  parseUserSearch,
} from '#/components/ecrm/user-search-field';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

import OrderDetailDrawer from './components/OrderDetailDrawer.vue';

const STATUS_TABS: Array<{ key: keyof PlatformRefundTabCounts; name: string }> =
  [
    { key: 'all', name: '全部' },
    { key: 'applied', name: '待审核' },
    { key: 'rejected', name: '审核未通过' },
    { key: 'approved', name: '审核通过' },
    { key: 'awaiting_receipt', name: '待收货' },
    { key: 'dispute', name: '纠纷中' },
    { key: 'completed', name: '已完成' },
  ];

const tabStatus = ref<keyof PlatformRefundTabCounts>('all');
const tabCounts = ref<PlatformRefundTabCounts>({
  all: 0,
  applied: 0,
  rejected: 0,
  approved: 0,
  awaiting_receipt: 0,
  dispute: 0,
  completed: 0,
});
const lastFormValues = ref<Record<string, unknown>>({});
const lastQueryParams = ref<Partial<PlatformRefundListParams>>({});

const refundDetail = ref<PlatformRefundOrder>();
const refundDetailLoading = ref(false);
const refundDetailTab = ref('info');
const auditTarget = ref<PlatformRefundOrder>();
const auditSubmitting = ref(false);
const canApprove = ref(false);
const canReject = ref(false);
const auditForm = reactive<{
  platform_mark: string;
  /** 1=同意 0=拒绝 — CRMEB 维权审核 */
  status: 0 | 1;
}>({
  status: 1,
  platform_mark: '',
});

const orderDetailDrawerRef = ref<InstanceType<typeof OrderDetailDrawer>>();

const logLoading = ref(false);
const logs = ref<PlatformRefundEvent[]>([]);
const logTotal = ref(0);
const logPage = ref(1);
const logLimit = ref(10);
const logTerminal = ref('');
const logDates = ref<string[]>([]);

const statusTagType: Record<
  string,
  'danger' | 'info' | 'success' | 'warning'
> = {
  applied: 'warning',
  merchant_handling: 'warning',
  awaiting_return: 'warning',
  awaiting_receipt: 'warning',
  refunding: 'warning',
  refunded: 'success',
  platform_intervene: 'danger',
  rejected: 'danger',
  cancelled: 'info',
};

function money(v?: number) {
  return `¥${Number(v || 0).toFixed(2)}`;
}

function moneyPlain(v?: number) {
  return Number(v || 0).toFixed(2);
}

/** List / expand empty placeholder */
function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || v === '' || v === '--') return '--';
  return String(v);
}

/** Detail empty → blank (match CRMEB) */
function blank(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || v === '' || v === '--' || v === '-') {
    return '';
  }
  return String(v);
}

function statusText(row?: PlatformRefundOrder) {
  return row?.status_label || '未知状态';
}

function refundTypeTitle(row?: PlatformRefundOrder) {
  const label =
    row?.refund_type_label ||
    (row?.refund_type === 2 ? '退货退款' : '仅退款');
  return `退款单信息-${label}`;
}

function isUserDeleted(row?: PlatformRefundOrder) {
  if (!row) return false;
  if (row.user_deleted) return true;
  if (Number(row.uid) < 0) return true;
  return String(row.nickname || '').trim() === '用户已被删除';
}

/** CRMEB 平台列表：仅「纠纷中/平台介入」可维权审核 */
function isDisputeStatus(row?: PlatformRefundOrder) {
  if (!row) return false;
  return (
    row.status_code === 'platform_intervene' || Number(row.status) === 4
  );
}

function rowCanDisputeReview(row?: PlatformRefundOrder) {
  return isDisputeStatus(row) && (canApprove.value || canReject.value);
}

function userNickText(row?: PlatformRefundOrder) {
  if (!row) return '';
  if (isUserDeleted(row)) {
    const uid = Number(row.uid);
    return `用户已被删除 | ${Number.isFinite(uid) ? uid : -1}`;
  }
  return blank(row.nickname);
}

function productTotalPrice(row: PlatformRefundOrder) {
  if (row.product_total_price != null) {
    return Number(row.product_total_price);
  }
  const products = row.products || [];
  return products.reduce((sum, p) => {
    const qty = Number(p.refund_num || p.product_num || 1);
    return sum + Number(p.product_price || 0) * qty;
  }, 0);
}

function refundProductCount(row: PlatformRefundOrder) {
  if (row.refund_num != null && row.refund_num > 0) return row.refund_num;
  return (row.products || []).reduce(
    (sum, p) => sum + Number(p.refund_num || p.product_num || 0),
    0,
  );
}

function setStatusTab(key: keyof PlatformRefundTabCounts) {
  if (tabStatus.value === key) return;
  tabStatus.value = key;
  gridApi.reload();
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
): PlatformRefundListParams {
  const values = formValues || {};
  lastFormValues.value = values;
  const range = Array.isArray(values.date_range) ? values.date_range : [];
  const isTrader = values.is_trader;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    tab_status: tabStatus.value === 'all' ? undefined : tabStatus.value,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    is_trader:
      isTrader === 0 || isTrader === 1 || isTrader === '0' || isTrader === '1'
        ? isTrader
        : undefined,
    refund_order_sn: String(values.refund_order_sn ?? '').trim() || undefined,
    order_sn: String(values.order_sn ?? '').trim() || undefined,
    user_search_type: parseUserSearch(values).type || 'nickname',
    user_search_keyword: parseUserSearch(values).keyword || undefined,
  };
}

async function loadTabCounts(formValues?: Record<string, unknown>) {
  const params = buildListParams(
    { currentPage: 1, pageSize: 10 },
    formValues ?? lastFormValues.value,
  );
  const {
    page: _p,
    limit: _l,
    tab_status: _t,
    ...filters
  } = params;
  try {
    tabCounts.value = await getPlatformRefundTabCountsApi(filters);
  } catch {
    tabCounts.value = {
      all: 0,
      applied: 0,
      rejected: 0,
      approved: 0,
      awaiting_receipt: 0,
      dispute: 0,
      completed: 0,
    };
  }
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  { ...LIST_DATE_RANGE_FIELD, label: '时间选择' },
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
    componentProps: { clearable: true, placeholder: '请输入退款单号' },
    fieldName: 'refund_order_sn',
    label: '退款单号',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入订单号' },
    fieldName: 'order_sn',
    label: '订单编号',
  },
  listUserSearchFormField(),
]);

const gridOptions: VxeGridProps<PlatformRefundOrder> = {
  columns: [
    {
      type: 'expand',
      width: 48,
      slots: { content: 'expandContent' },
    },
    {
      field: 'refund_order_sn',
      minWidth: 180,
      showOverflow: false,
      title: '退款单号',
    },
    {
      field: 'nickname',
      minWidth: 120,
      showOverflow: false,
      slots: { default: 'userInfo' },
      title: '用户信息',
    },
    {
      field: 'store_name',
      minWidth: 120,
      showOverflow: false,
      title: '店铺名称',
      formatter: ({ row }) =>
        row.store_name || row.mer_name || `店铺 #${row.mer_id}`,
    },
    {
      field: 'store_category_name',
      minWidth: 90,
      title: '店铺类别',
      formatter: ({ row }) => dash(row.store_category_name),
    },
    {
      field: 'refund_price',
      minWidth: 96,
      title: '退款金额',
      formatter: ({ cellValue }) => moneyPlain(Number(cellValue)),
    },
    {
      field: 'products',
      minWidth: 240,
      showOverflow: false,
      slots: { default: 'productInfo' },
      title: '商品信息',
    },
    {
      field: 'status_code',
      slots: { default: 'status' },
      title: '退款单状态',
      width: 110,
    },
    {
      field: 'refund_message',
      minWidth: 180,
      showOverflow: false,
      slots: { default: 'remark' },
      title: '说明',
    },
    platformListActionColumn({ width: 220 }),
  ],
  expandConfig: { trigger: 'default' },
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const params = buildListParams(page, formValues);
        lastQueryParams.value = params;
        const [data] = await Promise.all([
          listPlatformRefundsApi(params),
          loadTabCounts(formValues),
        ]);
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'refund_order_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    // false：关闭 Vxe 右上角圆形搜索开关（icon 字体缺失时会露出蓝色残片）
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [RefundDrawer, refundDrawerApi] = useVbenDrawer({
  class: 'w-[1100px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

const [AuditDrawer, auditDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  title: '维权审核',
  confirmText: '保存',
  onConfirm: async () => {
    await submitDisputeReview();
  },
});

async function exportRows() {
  try {
    const { page: _page, limit: _limit, ...filters } = lastQueryParams.value;
    const result = await exportPlatformRefundsApi({
      ...filters,
      reason: '平台退款订单列表导出',
    });
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      `已导出 ${result.row_count} 条${result.truncated ? '（已截断）' : ''}`,
    );
  } catch {
    ElMessage.info('导出列表接口暂不可用或无权限，已保留按钮占位');
  }
}

function resetLogs() {
  logs.value = [];
  logTotal.value = 0;
  logPage.value = 1;
  logLimit.value = 10;
  logTerminal.value = '';
  logDates.value = [];
}

async function loadLogs() {
  if (!refundDetail.value?.refund_order_id) return;
  logLoading.value = true;
  try {
    const range = Array.isArray(logDates.value) ? logDates.value : [];
    const data = await listPlatformRefundEventsApi(
      refundDetail.value.refund_order_id,
      {
        page: logPage.value,
        limit: logLimit.value,
        terminal: logTerminal.value || undefined,
        date_from: range[0],
        date_to: range[1],
      },
    );
    logs.value = data.list || [];
    logTotal.value = data.total || 0;
  } catch {
    logs.value = [];
    logTotal.value = 0;
  } finally {
    logLoading.value = false;
  }
}

async function openRefundDetail(row: PlatformRefundOrder) {
  refundDetail.value = undefined;
  refundDetailTab.value = 'info';
  refundDetailLoading.value = true;
  resetLogs();
  refundDrawerApi.setState({ title: '退款单详情', loading: true }).open();
  try {
    refundDetail.value = await getPlatformRefundApi(row.refund_order_id);
    await loadLogs();
  } finally {
    refundDetailLoading.value = false;
    refundDrawerApi.setState({ loading: false });
  }
}

function openOrderDetail(row: PlatformRefundOrder) {
  orderDetailDrawerRef.value?.open(row.order_id);
}

function onRefundDetailTabChange(name: string | number) {
  if (String(name) === 'logs') {
    void loadLogs();
  }
}

function openDisputeReview(row: PlatformRefundOrder) {
  if (!rowCanDisputeReview(row)) return;
  auditTarget.value = row;
  auditForm.status = canApprove.value ? 1 : 0;
  auditForm.platform_mark = '';
  auditDrawerApi.open();
}

async function submitDisputeReview() {
  const row = auditTarget.value;
  if (!row) return;
  const agree = auditForm.status === 1;
  if (agree && !canApprove.value) {
    ElMessage.warning('无同意退款权限');
    return;
  }
  if (!agree && !canReject.value) {
    ElMessage.warning('无拒绝退款权限');
    return;
  }
  const reason = auditForm.platform_mark.trim();
  if (!agree) {
    if (!reason) {
      ElMessage.warning('请输入拒绝原因');
      return;
    }
    if ([...reason].length > 500) {
      ElMessage.warning('拒绝原因不能超过 500 字');
      return;
    }
  }
  try {
    await confirm({
      content: agree
        ? '确认同意该纠纷退款？仅退款将进入支付渠道退款阶段；退货退款将等待用户寄回。已验证的支付回调才会标记为已退款。'
        : '确认拒绝该纠纷退款申请？拒绝后订单将恢复售后前状态。',
      icon: 'warning',
      title: agree ? '同意退款确认' : '拒绝退款确认',
    });
  } catch {
    return;
  }
  auditSubmitting.value = true;
  auditDrawerApi.lock();
  try {
    if (agree) {
      await approvePlatformRefundApi(row.refund_order_id);
      ElMessage.success('维权审核已同意');
    } else {
      await rejectPlatformRefundApi(row.refund_order_id, reason);
      ElMessage.success('维权审核已拒绝');
    }
    auditDrawerApi.close();
    if (
      refundDetail.value?.refund_order_id === row.refund_order_id
    ) {
      refundDetail.value = await getPlatformRefundApi(row.refund_order_id);
      await loadLogs();
    }
    gridApi.reload();
  } finally {
    auditSubmitting.value = false;
    auditDrawerApi.unlock();
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canApprove.value = permissions.includes('order.refund.approve');
  canReject.value = permissions.includes('order.refund.reject');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="refund-toolbar">
          <div class="refund-status-tabs" role="tablist">
            <button
              v-for="tab in STATUS_TABS"
              :key="tab.key"
              type="button"
              role="tab"
              class="refund-status-tabs__item"
              :aria-selected="tabStatus === tab.key"
              :class="{ 'is-active': tabStatus === tab.key }"
              @click="setStatusTab(tab.key)"
            >
              {{ tab.name }}({{ tabCounts[tab.key] || 0 }})
            </button>
          </div>
          <div class="refund-toolbar__actions">
            <ElButton type="primary" @click="exportRows">导出列表</ElButton>
          </div>
        </div>
      </template>

      <template #userInfo="{ row }">
        <span>{{ userNickText(row) || dash(row.nickname) }}</span>
      </template>

      <template #productInfo="{ row }">
        <div v-if="row.products?.length" class="refund-products">
          <div
            v-for="(p, idx) in row.products"
            :key="`${row.refund_order_id}-${idx}`"
            class="refund-product"
          >
            <ElImage
              class="refund-product__thumb"
              :src="resolveCosMediaUrl(p.product_image || '')"
              fit="cover"
            >
              <template #error>
                <div class="refund-product__thumb-fallback" />
              </template>
            </ElImage>
            <div class="refund-product__text">
              <div class="refund-product__name">
                {{ p.product_info || '—' }}
              </div>
              <div class="refund-product__price">
                {{ money(p.product_price) }} x
                {{ p.refund_num || p.product_num || 1 }}
              </div>
            </div>
          </div>
        </div>
        <span v-else>—</span>
      </template>

      <template #status="{ row }">
        <ElTag
          :type="statusTagType[row.status_code] || 'info'"
          effect="plain"
          class="refund-status-tag"
        >
          {{ statusText(row) }}
        </ElTag>
      </template>

      <template #remark="{ row }">
        <div class="refund-remark">
          <div>退款原因: {{ dash(row.refund_message) }}</div>
          <div>
            状态变更: {{ formatShanghaiDateTime(row.status_time) }}
          </div>
        </div>
      </template>

      <template #action="{ row }">
        <ElButton
          v-if="rowCanDisputeReview(row)"
          link
          type="primary"
          @click="openDisputeReview(row)"
        >
          维权审核
        </ElButton>
        <ElButton link type="primary" @click="openOrderDetail(row)">
          订单详情
        </ElButton>
        <ElButton link type="primary" @click="openRefundDetail(row)">
          退款单详情
        </ElButton>
      </template>

      <template #expandContent="{ row }">
        <div class="refund-expand">
          <div class="refund-expand__col">
            <div class="refund-expand__item">
              <span class="label">退款商品总价:</span>
              <span class="value">{{ moneyPlain(productTotalPrice(row)) }}</span>
            </div>
            <div class="refund-expand__item">
              <span class="label">用户备注:</span>
              <span class="value">{{ dash(row.user_remark) }}</span>
            </div>
          </div>
          <div class="refund-expand__col">
            <div class="refund-expand__item">
              <span class="label">退款商品总数:</span>
              <span class="value">{{ refundProductCount(row) }}</span>
            </div>
            <div class="refund-expand__item">
              <span class="label">店铺备注:</span>
              <span class="value">{{ dash(row.merchant_remark) }}</span>
            </div>
          </div>
          <div class="refund-expand__col">
            <div class="refund-expand__item">
              <span class="label">申请退款时间:</span>
              <span class="value">{{
                formatShanghaiDateTime(row.create_time)
              }}</span>
            </div>
          </div>
        </div>
      </template>
    </Grid>

    <RefundDrawer>
      <ElSkeleton :loading="refundDetailLoading" animated :rows="10">
        <template #default>
          <div v-if="refundDetail" class="refund-detail">
            <div class="refund-detail__header">
              <div class="refund-detail__identity">
                <div class="refund-detail__icon">
                  <IconifyIcon icon="ant-design:file-text-outlined" />
                </div>
                <div class="refund-detail__titles">
                  <div class="refund-detail__type">
                    {{ refundTypeTitle(refundDetail) }}
                  </div>
                  <div class="refund-detail__sn">
                    订单编号: {{ blank(refundDetail.order_sn) }}
                  </div>
                </div>
                <div
                  v-if="rowCanDisputeReview(refundDetail)"
                  class="refund-detail__header-actions"
                >
                  <ElButton
                    type="primary"
                    @click="openDisputeReview(refundDetail)"
                  >
                    维权审核
                  </ElButton>
                </div>
              </div>
              <div class="refund-detail__status">
                <div class="refund-detail__status-item">
                  <span class="label">退款单状态</span>
                  <span class="value is-warn">{{
                    statusText(refundDetail)
                  }}</span>
                </div>
                <div class="refund-detail__status-item">
                  <span class="label">实际退款</span>
                  <span class="value">{{
                    money(refundDetail.refund_price)
                  }}</span>
                </div>
                <div class="refund-detail__status-item">
                  <span class="label">退回方式</span>
                  <span class="value">{{
                    blank(refundDetail.refund_method) || '原路返回'
                  }}</span>
                </div>
                <div class="refund-detail__status-item">
                  <span class="label">申请退款时间</span>
                  <span class="value">{{
                    formatShanghaiDateTime(refundDetail.create_time)
                  }}</span>
                </div>
              </div>
            </div>

            <ElTabs
              v-model="refundDetailTab"
              class="refund-detail__tabs"
              @tab-change="onRefundDetailTabChange"
            >
              <ElTabPane label="订单信息" name="info">
                <div class="refund-section refund-section--user">
                  <div class="refund-section__title">用户信息</div>
                  <div class="refund-kv-grid">
                    <div class="refund-kv">
                      <span class="refund-kv__label">用户昵称:</span>
                      <span class="refund-kv__value">{{
                        userNickText(refundDetail)
                      }}</span>
                    </div>
                    <div class="refund-kv">
                      <span class="refund-kv__label">用户电话:</span>
                      <span class="refund-kv__value">{{
                        blank(refundDetail.user_phone_mask)
                      }}</span>
                    </div>
                  </div>
                </div>

                <div class="refund-section refund-section--aftersale">
                  <div class="refund-section__title">售后提交信息</div>
                  <div class="refund-kv-grid refund-kv-grid--3">
                    <div class="refund-kv">
                      <span class="refund-kv__label">退款总件数:</span>
                      <span class="refund-kv__value">{{
                        refundDetail.refund_num ?? ''
                      }}</span>
                    </div>
                    <div class="refund-kv">
                      <span class="refund-kv__label">退款原因:</span>
                      <span class="refund-kv__value">{{
                        blank(refundDetail.refund_message)
                      }}</span>
                    </div>
                    <div class="refund-kv">
                      <span class="refund-kv__label">退款发起方:</span>
                      <span class="refund-kv__value">{{
                        blank(refundDetail.refund_initiator)
                      }}</span>
                    </div>
                  </div>
                  <div class="refund-kv-stack">
                    <div class="refund-kv">
                      <span class="refund-kv__label">用户备注:</span>
                      <span class="refund-kv__value">{{
                        blank(refundDetail.user_remark)
                      }}</span>
                    </div>
                    <div class="refund-kv">
                      <span class="refund-kv__label">退款凭证:</span>
                      <span class="refund-kv__value">
                        <template
                          v-if="refundDetail.refund_evidence?.length"
                        >
                          <ElImage
                            v-for="(url, i) in refundDetail.refund_evidence"
                            :key="i"
                            class="refund-evidence-thumb"
                            :src="resolveCosMediaUrl(url)"
                            :preview-src-list="
                              refundDetail.refund_evidence.map((u) =>
                                resolveCosMediaUrl(u),
                              )
                            "
                            fit="cover"
                          >
                            <template #error>
                              <div class="refund-product__thumb-fallback" />
                            </template>
                          </ElImage>
                        </template>
                      </span>
                    </div>
                  </div>
                </div>

                <div class="refund-section refund-section--last">
                  <div class="refund-section__title">店铺备注</div>
                  <div class="refund-kv">
                    <span class="refund-kv__label">店铺备注:</span>
                    <span class="refund-kv__value">{{
                      blank(refundDetail.merchant_remark)
                    }}</span>
                  </div>
                </div>
              </ElTabPane>

              <ElTabPane label="商品信息" name="products">
                <ElTable
                  :data="refundDetail.products || []"
                  border
                  class="refund-detail-table"
                >
                  <ElTableColumn
                    label="商品ID"
                    min-width="90"
                    prop="product_id"
                  />
                  <ElTableColumn label="商品信息" min-width="280">
                    <template #default="{ row }">
                      <div class="refund-product refund-product--detail">
                        <ElImage
                          class="refund-product__thumb"
                          :src="resolveCosMediaUrl(row.product_image || '')"
                          fit="cover"
                        >
                          <template #error>
                            <div class="refund-product__thumb-fallback" />
                          </template>
                        </ElImage>
                        <div class="refund-product__text">
                          <div class="refund-product__name">
                            {{ row.product_info || '' }}
                          </div>
                          <div class="refund-product__sku">
                            规格：{{ row.product_sku || '默认' }}
                          </div>
                        </div>
                      </div>
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="售价(元)" min-width="100">
                    <template #default="{ row }">
                      {{ moneyPlain(row.product_price) }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn
                    label="退款数量"
                    min-width="90"
                    prop="refund_num"
                  />
                  <ElTableColumn label="实付金额(元)" min-width="110">
                    <template #default="{ row }">
                      {{ moneyPlain(row.pay_price) }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="退款金额(元)" min-width="110">
                    <template #default="{ row }">
                      {{ moneyPlain(row.refund_price) }}
                    </template>
                  </ElTableColumn>
                </ElTable>
              </ElTabPane>

              <ElTabPane label="订单记录" name="logs">
                <div class="refund-log-filters">
                  <div class="refund-log-filters__item">
                    <span class="filter-label">操作端</span>
                    <ElSelect
                      v-model="logTerminal"
                      clearable
                      placeholder="请选择"
                      style="width: 160px"
                      @change="
                        () => {
                          logPage = 1;
                          loadLogs();
                        }
                      "
                    >
                      <ElOption label="用户" value="user" />
                      <ElOption label="店铺" value="merchant" />
                      <ElOption label="平台" value="platform" />
                      <ElOption label="系统" value="system" />
                    </ElSelect>
                  </div>
                  <div class="refund-log-filters__item">
                    <span class="filter-label">操作时间</span>
                    <ElDatePicker
                      v-model="logDates"
                      type="datetimerange"
                      value-format="YYYY-MM-DD HH:mm:ss"
                      start-placeholder="开始时间"
                      end-placeholder="结束时间"
                      range-separator="至"
                      @change="
                        () => {
                          logPage = 1;
                          loadLogs();
                        }
                      "
                    />
                  </div>
                </div>
                <ElTable
                  v-loading="logLoading"
                  :data="logs"
                  border
                  class="refund-detail-table"
                >
                  <ElTableColumn label="订单编号" min-width="190">
                    <template #default="{ row }">
                      {{ row.order_sn || refundDetail.refund_order_sn }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="操作记录" min-width="160">
                    <template #default="{ row }">
                      {{ row.content || row.reason || '' }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="操作角色" width="100">
                    <template #default="{ row }">
                      {{ row.role || '' }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="操作人" min-width="180">
                    <template #default="{ row }">
                      {{ row.operator || `ID:${row.actor_id}` }}
                    </template>
                  </ElTableColumn>
                  <ElTableColumn label="操作时间" min-width="170">
                    <template #default="{ row }">
                      {{
                        formatShanghaiDateTime(
                          row.operate_time || row.created_at,
                        )
                      }}
                    </template>
                  </ElTableColumn>
                </ElTable>
                <div class="refund-log-pager">
                  <ElPagination
                    v-model:current-page="logPage"
                    v-model:page-size="logLimit"
                    background
                    layout="prev, pager, next, jumper"
                    :total="logTotal"
                    @current-change="loadLogs"
                    @size-change="
                      () => {
                        logPage = 1;
                        loadLogs();
                      }
                    "
                  />
                </div>
              </ElTabPane>
            </ElTabs>
          </div>
        </template>
      </ElSkeleton>
    </RefundDrawer>

    <AuditDrawer>
      <ElForm label-width="96px" @submit.prevent>
        <ElFormItem label="审核结果：" required>
          <ElRadioGroup v-model="auditForm.status">
            <ElRadio v-if="canApprove" :value="1">同意</ElRadio>
            <ElRadio v-if="canReject" :value="0">拒绝</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem
          v-if="auditForm.status === 0"
          label="填写原因："
          required
        >
          <ElInput
            v-model="auditForm.platform_mark"
            type="textarea"
            :rows="4"
            maxlength="500"
            show-word-limit
            placeholder="请输入拒绝原因"
            :disabled="auditSubmitting"
          />
        </ElFormItem>
      </ElForm>
    </AuditDrawer>

    <OrderDetailDrawer ref="orderDetailDrawerRef" />
  </Page>
</template>

<style scoped>
.refund-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  min-width: 0;
}

.refund-status-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 28px;
  border-bottom: 1px solid hsl(var(--border));
}

.refund-status-tabs__item {
  margin-bottom: -1px;
  padding: 10px 2px 12px;
  color: hsl(var(--foreground) / 70%);
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
  background: transparent;
  border: 0;
  border-bottom: 2px solid transparent;
}

.refund-status-tabs__item:hover {
  color: hsl(var(--primary));
}

.refund-status-tabs__item.is-active {
  color: hsl(var(--primary));
  font-weight: 600;
  border-bottom-color: hsl(var(--primary));
}

.refund-toolbar__actions {
  display: flex;
  gap: 8px;
  justify-content: flex-start;
}

.refund-products {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 4px 0;
}

.refund-product {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.refund-product__thumb {
  flex: 0 0 48px;
  width: 48px;
  height: 48px;
  overflow: hidden;
  border-radius: 4px;
  background: var(--el-fill-color-light);
}

.refund-product__thumb-fallback {
  width: 100%;
  height: 100%;
  min-height: 48px;
  background: var(--el-fill-color-light);
}

.refund-product__text {
  min-width: 0;
  line-height: 1.45;
  word-break: break-all;
}

.refund-product__name {
  color: var(--el-text-color-primary);
  font-size: 13px;
}

.refund-product__price,
.refund-product__sku {
  margin-top: 4px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.refund-product--detail .refund-product__thumb {
  flex: 0 0 56px;
  width: 56px;
  height: 56px;
}

.refund-status-tag.el-tag--success {
  --el-tag-border-color: var(--el-color-success);
}

.refund-remark {
  line-height: 1.55;
  font-size: 13px;
}

.refund-expand {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px 32px;
  padding: 10px 16px;
}

.refund-expand__col {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.refund-expand__item {
  display: flex;
  gap: 6px;
  font-size: 13px;
  line-height: 22px;
}

.refund-expand__item .label {
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
}

.refund-expand__item .value {
  color: var(--el-text-color-primary);
  word-break: break-all;
}

.refund-detail {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.refund-detail__header {
  padding-bottom: 4px;
}

.refund-detail__identity {
  display: flex;
  gap: 12px;
  align-items: center;
}

.refund-detail__header-actions {
  margin-left: auto;
}

.refund-detail__icon {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  color: #fff;
  font-size: 22px;
  background: var(--el-color-primary);
  border-radius: 8px;
}

.refund-detail__type {
  color: var(--el-text-color-primary);
  font-size: 18px;
  font-weight: 600;
  line-height: 26px;
}

.refund-detail__sn {
  margin-top: 2px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 20px;
}

.refund-detail__status {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px 24px;
  margin-top: 16px;
  margin-bottom: 4px;
}

.refund-detail__status-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.refund-detail__status-item .label {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 18px;
}

.refund-detail__status-item .value {
  color: var(--el-text-color-primary);
  font-size: 14px;
  line-height: 22px;
}

.refund-detail__status-item .value.is-warn {
  color: #ed6a0c;
  font-weight: 500;
}

.refund-detail__tabs {
  min-height: 420px;
}

.refund-detail__tabs :deep(.el-tabs__header) {
  margin-bottom: 18px;
}

.refund-detail__tabs :deep(.el-tabs__nav-wrap::after) {
  height: 1px;
  background-color: var(--el-border-color-lighter);
}

.refund-detail__tabs :deep(.el-tabs__item) {
  height: 40px;
  padding: 0 18px;
  color: var(--el-text-color-secondary);
  font-size: 14px;
  font-weight: 400;
}

.refund-detail__tabs :deep(.el-tabs__item.is-active) {
  color: var(--el-color-primary);
  font-weight: 500;
}

.refund-detail__tabs :deep(.el-tabs__active-bar) {
  height: 2px;
  background-color: var(--el-color-primary);
}

.refund-section {
  padding: 2px 0 16px;
}

.refund-section--user {
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.refund-section--aftersale {
  margin-top: 16px;
  border-bottom: 1px dashed var(--el-border-color);
}

.refund-section--last {
  margin-top: 16px;
}

.refund-section__title {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 14px;
  color: var(--el-text-color-primary);
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

.refund-section__title::before {
  content: '';
  width: 3px;
  height: 14px;
  border-radius: 1px;
  background: var(--el-color-primary);
}

.refund-kv-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px 20px;
}

.refund-kv-grid--3 {
  grid-template-columns: repeat(3, minmax(0, 1fr));
  margin-bottom: 10px;
}

.refund-kv-stack {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.refund-kv {
  display: flex;
  gap: 4px;
  min-width: 0;
  font-size: 13px;
  line-height: 22px;
}

.refund-kv__label {
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
}

.refund-kv__value {
  min-width: 0;
  color: var(--el-text-color-primary);
  word-break: break-all;
}

.refund-evidence-thumb {
  width: 56px;
  height: 56px;
  margin-right: 8px;
  border-radius: 4px;
}

.refund-detail-table :deep(.el-table__header th) {
  background: #f5f7fa;
  color: var(--el-text-color-regular);
  font-weight: 500;
}

.refund-detail-table :deep(.el-table__cell) {
  padding: 12px 0;
}

.refund-log-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 16px 24px;
  margin-bottom: 14px;
}

.refund-log-filters__item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.filter-label {
  color: var(--el-text-color-regular);
  font-size: 13px;
  white-space: nowrap;
}

.refund-log-pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 14px;
}

@media (max-width: 960px) {
  .refund-detail__status,
  .refund-kv-grid,
  .refund-kv-grid--3,
  .refund-expand {
    grid-template-columns: 1fr;
  }
}
</style>
