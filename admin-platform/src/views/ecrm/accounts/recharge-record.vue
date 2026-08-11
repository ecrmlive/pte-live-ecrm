<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import {
  ElAvatar,
  ElButton,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElMessage,
} from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  getPlatformUserRechargeTotalApi,
  listPlatformUserRechargesApi,
  refundPlatformUserRechargeApi,
  type PlatformUserRechargeQuery,
  type PlatformUserRechargeRow,
  type PlatformUserRechargeTotal,
} from '#/api/core/platform-user-recharge';
import {
  listUserSearchFormField,
  parseUserSearch,
} from '#/components/ecrm/user-search-field';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const emptyTotal = (): PlatformUserRechargeTotal => ({
  total_pay_price: 0,
  total_refund_price: 0,
  total_routine_price: 0,
  total_wx_price: 0,
});

const lastFormValues = ref<Record<string, unknown>>({});
const canRead = ref(false);
const canRefund = ref(false);
const total = ref<PlatformUserRechargeTotal>(emptyTotal());
const refunding = ref(false);
const refundRow = ref<PlatformUserRechargeRow | null>(null);
const refundForm = reactive({
  amount: 0,
});

const USER_SEARCH_OPTIONS = [
  { label: '昵称', value: 'nickname' },
  { label: '用户ID', value: 'uid' },
  { label: '手机号', value: 'phone' },
  { label: '姓名', value: 'real_name' },
];

const summaryCards = computed(() => [
  {
    key: 'total_pay_price',
    label: '充值总金额',
    value: formatMoney(total.value.total_pay_price),
    icon: 'lucide:shopping-bag',
    tone: 'blue',
  },
  {
    key: 'total_routine_price',
    label: '小程序充值金额',
    value: formatMoney(total.value.total_routine_price),
    icon: 'lucide:handshake',
    tone: 'orange',
  },
  {
    key: 'total_wx_price',
    label: '公众号充值金额',
    value: formatMoney(total.value.total_wx_price),
    icon: 'lucide:wallet',
    tone: 'green',
  },
]);

function formatMoney(v?: number) {
  return Number(v || 0).toFixed(2);
}

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || String(v).trim() === '') return '-';
  return String(v);
}

function mediaUrl(url?: string) {
  return resolveCosMediaUrl(String(url || '').trim());
}

function buildFilterParams(
  formValues?: Record<string, unknown>,
): PlatformUserRechargeQuery {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  const paidRaw = formValues?.paid;
  const typeRaw = String(formValues?.recharge_type ?? '').trim();
  const userSearch = parseUserSearch(formValues);
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    paid: paidRaw === 0 || paidRaw === 1 ? Number(paidRaw) : undefined,
    recharge_type: typeRaw || undefined,
    user_type: userSearch.type || 'nickname',
    user_keyword: userSearch.keyword || undefined,
    order_id: String(formValues?.order_id ?? '').trim() || undefined,
  };
}

async function loadTotal() {
  try {
    total.value = (await getPlatformUserRechargeTotalApi()) || emptyTotal();
  } catch {
    total.value = emptyTotal();
  }
}

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    {
      ...LIST_DATE_RANGE_FIELD,
      label: '时间选择',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '已支付', value: 1 },
          { label: '未支付', value: 0 },
        ],
        placeholder: '全部',
      },
      fieldName: 'paid',
      label: '是否支付',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '微信', value: '1' },
          { label: '支付宝', value: '2' },
        ],
        placeholder: '全部',
      },
      fieldName: 'recharge_type',
      label: '充值类型',
    },
    listUserSearchFormField({
      options: USER_SEARCH_OPTIONS,
      typeWidth: '96px',
    }),
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入订单号',
      },
      fieldName: 'order_id',
      label: '订单号',
    },
  ],
  {
    commonConfig: { componentProps: { class: 'w-full' } },
    submitButtonOptions: { content: '搜索' },
    handleSubmit: async (values) => {
      lastFormValues.value = { ...values };
      await Promise.all([gridApi.reload(values), loadTotal()]);
    },
    handleReset: async () => {
      await formApi.resetForm();
      const values = (await formApi.getValues()) ?? {};
      lastFormValues.value = { ...values };
      await Promise.all([gridApi.reload(values), loadTotal()]);
    },
  },
);

const [Form, formApi] = useVbenForm(formOptions);

const gridOptions: VxeGridProps<PlatformUserRechargeRow> = {
  columns: [
    { field: 'recharge_id', title: 'ID', width: 90 },
    {
      field: 'avatar',
      slots: { default: 'avatar' },
      title: '头像',
      width: 80,
    },
    {
      field: 'nickname',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '用户昵称',
    },
    {
      field: 'order_id',
      minWidth: 200,
      showOverflow: 'tooltip',
      title: '订单号',
    },
    {
      field: 'price',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 110,
      title: '支付金额',
    },
    {
      field: 'give_price',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 110,
      title: '赠送金额',
    },
    {
      field: 'paid_name',
      formatter: ({ row }) => row.paid_name || (row.paid === 1 ? '已支付' : '未支付'),
      minWidth: 90,
      title: '是否支付',
    },
    {
      field: 'recharge_type_name',
      formatter: ({ cellValue, row }) =>
        dash(cellValue || row.recharge_type),
      minWidth: 100,
      title: '充值类型',
    },
    {
      field: 'pay_time',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '-',
      minWidth: 170,
      title: '支付时间',
    },
    platformListActionColumn({ width: 100 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const values =
          formValues && Object.keys(formValues).length > 0
            ? formValues
            : lastFormValues.value;
        const filters = buildFilterParams(values);
        const data = await listPlatformUserRechargesApi({
          page: page.currentPage,
          limit: page.pageSize,
          ...filters,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'recharge_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [RefundDrawer, refundDrawerApi] = useVbenDrawer({
  class: 'w-[520px] max-w-[96vw]',
  placement: 'right',
  title: '充值退款',
  onConfirm: async () => {
    const row = refundRow.value;
    if (!row) return;
    const remain = Number(
      (row.price - Number(row.refund_price || 0)).toFixed(2),
    );
    const amount = Number(Number(refundForm.amount || 0).toFixed(2));
    if (!(amount > 0) || amount > remain) {
      ElMessage.warning('请输入有效退款金额');
      return;
    }
    refunding.value = true;
    refundDrawerApi.lock();
    try {
      await refundPlatformUserRechargeApi(row.recharge_id, {
        amount,
        idempotency_key: `user-recharge-refund-${row.recharge_id}-${crypto.randomUUID()}`,
      });
      ElMessage.success('退款成功');
      refundDrawerApi.close();
      await Promise.all([
        gridApi.reload(lastFormValues.value),
        loadTotal(),
      ]);
    } finally {
      refunding.value = false;
      refundDrawerApi.unlock();
    }
  },
});

function openRefund(row: PlatformUserRechargeRow) {
  if (!canRefund.value || !row.can_refund) return;
  refundRow.value = row;
  refundForm.amount = Number(
    (row.price - Number(row.refund_price || 0)).toFixed(2),
  );
  refundDrawerApi.setState({
    title: `充值退款 · #${row.recharge_id}`,
  });
  refundDrawerApi.open();
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canRead.value = codes.includes('accounts.recharge_record.read');
  canRefund.value = codes.includes('accounts.recharge_record.refund');
  if (canRead.value) {
    await Promise.all([gridApi.reload(), loadTotal()]);
  }
});
</script>

<template>
  <Page auto-content-height>
    <div class="recharge-filter">
      <Form />
    </div>

    <div v-if="canRead" class="recharge-summary">
      <div class="recharge-summary__grid">
        <div
          v-for="card in summaryCards"
          :key="card.key"
          class="recharge-summary__card"
          :class="`recharge-summary__card--${card.tone}`"
        >
          <div class="recharge-summary__icon">
            <IconifyIcon :icon="card.icon" />
          </div>
          <div class="recharge-summary__body">
            <div class="recharge-summary__value">{{ card.value }}</div>
            <div class="recharge-summary__label">{{ card.label }}</div>
          </div>
        </div>
      </div>
    </div>

    <Grid>
      <template #avatar="{ row }">
        <ElAvatar :size="36" :src="mediaUrl(row.avatar) || undefined">
          {{ (row.nickname || '用').slice(0, 1) }}
        </ElAvatar>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="canRefund && row.can_refund"
          link
          type="primary"
          @click="openRefund(row)"
        >
          退款
        </ElButton>
        <span v-else class="recharge-action-empty">-</span>
      </template>
    </Grid>

    <RefundDrawer>
      <ElForm v-if="refundRow" label-width="100px" @submit.prevent>
        <ElFormItem label="订单号">
          <span>{{ refundRow.order_id }}</span>
        </ElFormItem>
        <ElFormItem label="用户">
          <span>
            {{ dash(refundRow.nickname) }}（UID {{ refundRow.uid }}）
          </span>
        </ElFormItem>
        <ElFormItem label="支付金额">
          <span>¥{{ formatMoney(refundRow.price) }}</span>
        </ElFormItem>
        <ElFormItem label="已退款">
          <span>¥{{ formatMoney(refundRow.refund_price) }}</span>
        </ElFormItem>
        <ElFormItem label="退款金额" required>
          <ElInputNumber
            v-model="refundForm.amount"
            :min="0.01"
            :max="
              Number(
                (
                  refundRow.price - Number(refundRow.refund_price || 0)
                ).toFixed(2),
              )
            "
            :precision="2"
            :step="0.01"
            class="w-full"
            controls-position="right"
          />
        </ElFormItem>
        <ElFormItem label="说明">
          <span class="recharge-refund-tip">
            退款将按支付占比回扣用户余额（含赠送部分），并写入资产流水。
          </span>
        </ElFormItem>
      </ElForm>
    </RefundDrawer>
  </Page>
</template>

<style scoped>
.recharge-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.recharge-summary {
  padding: 16px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.recharge-summary__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  width: 100%;
}

.recharge-summary__card {
  display: flex;
  gap: 16px;
  align-items: center;
  min-height: 88px;
  padding: 20px 22px;
  background: hsl(var(--background));
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  box-shadow: 0 1px 2px rgb(0 0 0 / 3%);
}

.recharge-summary__icon {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  color: #fff;
  font-size: 24px;
  border-radius: 50%;
}

.recharge-summary__card--blue .recharge-summary__icon {
  background: #409eff;
}

.recharge-summary__card--orange .recharge-summary__icon {
  background: #e6a23c;
}

.recharge-summary__card--green .recharge-summary__icon {
  background: #67c23a;
}

.recharge-summary__value {
  color: var(--el-text-color-primary);
  font-size: 24px;
  font-weight: 600;
  line-height: 1.2;
}

.recharge-summary__label {
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.2;
}

.recharge-action-empty {
  color: var(--el-text-color-placeholder);
}

.recharge-refund-tip {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
}

@media (max-width: 1100px) {
  .recharge-summary__grid {
    grid-template-columns: 1fr;
  }
}
</style>
