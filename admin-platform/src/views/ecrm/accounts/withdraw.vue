<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElSkeleton,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  approvePlatformWithdrawApi,
  getPlatformWithdrawApi,
  listPlatformWithdrawsApi,
  markPlatformWithdrawPaidApi,
  rejectPlatformWithdrawApi,
  type PlatformWithdraw,
} from '#/api/core/platform-finance';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const auditStatus: Record<
  number,
  { label: string; type: 'danger' | 'info' | 'success' | 'warning' }
> = {
  [-1]: { label: '审核拒绝', type: 'danger' },
  0: { label: '待平台审核', type: 'warning' },
  1: { label: '审核通过', type: 'success' },
};

const current = ref<PlatformWithdraw>();
const detailLoading = ref(false);
const rejecting = ref(false);
const canReview = ref(false);
const rejectForm = reactive({ refusal: '' });

function auditInfo(status: number) {
  return auditStatus[status] || { label: '未知状态', type: 'info' as const };
}

function accountType(type: number) {
  return (
    ({ 1: '银行卡', 2: '微信', 3: '支付宝' } as Record<number, string>)[type] ||
    '未知'
  );
}

function transferInfo(row: PlatformWithdraw) {
  return row.financial_status === 1
    ? { label: '已打款', type: 'success' as const }
    : { label: '未打款', type: 'info' as const };
}

function canAudit(row: PlatformWithdraw) {
  return canReview.value && row.status === 0;
}

function canMarkPaid(row: PlatformWithdraw) {
  return canReview.value && row.withdrawal_status === 'approved';
}

function idempotencyKey(id: number) {
  return `withdraw-paid-${id}-${crypto.randomUUID()}`;
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const statusRaw = formValues?.status;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    status:
      statusRaw === 0 || statusRaw === 1 || statusRaw === -1
        ? Number(statusRaw)
        : undefined,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    financial_sn: String(formValues?.financial_sn ?? '').trim() || undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '用户 ID / 关键词',
    },
    fieldName: 'keyword',
    label: '用户搜索',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '申请单号' },
    fieldName: 'financial_sn',
    label: '申请单号',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待平台审核', value: 0 },
        { label: '审核通过', value: 1 },
        { label: '审核拒绝', value: -1 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<PlatformWithdraw> = {
  columns: [
    {
      field: 'financial_sn',
      minWidth: 180,
      showOverflow: false,
      title: '申请单号',
    },
    { field: 'user_id', title: '用户 ID', width: 100 },
    {
      field: 'extract_money',
      title: '提现金额',
      width: 116,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'financial_type',
      title: '收款方式',
      width: 104,
      formatter: ({ cellValue }) => accountType(Number(cellValue)),
    },
    {
      field: 'financial_account',
      minWidth: 150,
      showOverflow: false,
      title: '收款账户',
    },
    {
      field: 'status',
      slots: { default: 'auditStatus' },
      title: '审核状态',
      width: 116,
    },
    {
      field: 'financial_status',
      slots: { default: 'transferStatus' },
      title: '打款状态',
      width: 100,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '申请时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    platformListActionColumn({ width: 192 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformWithdrawsApi(buildListParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'financial_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [RejectDrawer, rejectDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '确认拒绝',
  cancelText: '取消',
  placement: 'right',
  title: '拒绝提现',
  onConfirm: async () => reject(),
});

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[560px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

async function openDetail(row: PlatformWithdraw) {
  current.value = undefined;
  detailLoading.value = true;
  detailDrawerApi.setState({ title: '提现详情', loading: true }).open();
  try {
    current.value = await getPlatformWithdrawApi(row.financial_id);
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}

async function reloadGrid() {
  await gridApi.reload();
}

async function approve(row: PlatformWithdraw) {
  try {
    await ElMessageBox.confirm(
      '确认审核通过该用户提现申请？审核通过不等于已打款，必须另行登记内部打款凭证。',
      '审核通过确认',
      {
        confirmButtonText: '确认通过',
        cancelButtonText: '取消',
        type: 'warning',
      },
    );
    await approvePlatformWithdrawApi(row.financial_id);
    ElMessage.success('提现申请已审核通过');
    await reloadGrid();
  } catch {
    // 用户取消或接口已返回错误时，requestClient 统一处理提示。
  }
}

async function markPaid(row: PlatformWithdraw) {
  try {
    const { value } = await ElMessageBox.prompt(
      '请输入内部打款凭证编号，不录入银行卡、账户或密钥。',
      '登记打款凭证',
      {
        inputPattern: /\S{3,}/,
        inputErrorMessage: '凭证编号至少 3 个字符',
      },
    );
    await markPlatformWithdrawPaidApi(row.financial_id, {
      idempotency_key: idempotencyKey(row.financial_id),
      payout_reference: value.trim(),
    });
    ElMessage.success('打款凭证已登记，提现状态已更新为已打款');
    await reloadGrid();
  } catch {
    // 用户取消或接口已返回错误时，requestClient 统一处理提示。
  }
}

function openReject(row: PlatformWithdraw) {
  current.value = row;
  rejectForm.refusal = '';
  rejectDrawerApi.open();
}

async function reject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写拒绝原因');
    return;
  }
  if (!current.value) return;
  rejecting.value = true;
  try {
    await rejectPlatformWithdrawApi(current.value.financial_id, refusal);
    rejectDrawerApi.close();
    ElMessage.success('提现申请已拒绝；资金释放由业务资金域处理。');
    await reloadGrid();
  } finally {
    rejecting.value = false;
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canReview.value =
    profile.roles.includes('platform') &&
    permissions.includes('accounts.withdraw.review');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #auditStatus="{ row }">
        <ElTag :type="auditInfo(row.status).type">
          {{ auditInfo(row.status).label }}
        </ElTag>
      </template>
      <template #transferStatus="{ row }">
        <ElTag :type="transferInfo(row).type">
          {{ transferInfo(row).label }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <template v-if="canAudit(row)">
          <ElButton link type="success" @click="approve(row)">通过</ElButton>
          <ElButton link type="danger" @click="openReject(row)">拒绝</ElButton>
        </template>
        <ElButton
          v-else-if="canMarkPaid(row)"
          link
          type="primary"
          @click="markPaid(row)"
        >
          登记打款
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <template #default>
          <template v-if="current">
            <ElDescriptions :column="1" border>
              <ElDescriptionsItem label="申请单号">
                {{ current.financial_sn }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="用户 ID">
                {{ current.user_id }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="提现金额">
                ¥{{ Number(current.extract_money).toFixed(2) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="收款方式">
                {{ accountType(current.financial_type) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="收款账户">
                {{ current.financial_account }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="审核状态">
                <ElTag :type="auditInfo(current.status).type">
                  {{ auditInfo(current.status).label }}
                </ElTag>
              </ElDescriptionsItem>
              <ElDescriptionsItem label="打款状态">
                <ElTag :type="transferInfo(current).type">
                  {{ transferInfo(current).label }}
                </ElTag>
              </ElDescriptionsItem>
              <ElDescriptionsItem
                v-if="current.payout_reference"
                label="内部打款凭证"
              >
                {{ current.payout_reference }}
              </ElDescriptionsItem>
              <ElDescriptionsItem v-if="current.paid_at" label="打款登记时间">
                {{ formatShanghaiDateTime(current.paid_at) }}
              </ElDescriptionsItem>
              <ElDescriptionsItem v-if="current.refusal" label="拒绝原因">
                {{ current.refusal }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="申请备注">
                {{ current.mark || '—' }}
              </ElDescriptionsItem>
              <ElDescriptionsItem label="申请时间">
                {{ formatShanghaiDateTime(current.create_time) }}
              </ElDescriptionsItem>
            </ElDescriptions>
          </template>
        </template>
      </ElSkeleton>
    </DetailDrawer>

    <RejectDrawer>
      <ElForm label-width="84px">
        <ElFormItem label="拒绝原因" required>
          <ElInput
            v-model="rejectForm.refusal"
            :rows="4"
            maxlength="200"
            placeholder="请向商户说明拒绝原因"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </RejectDrawer>
  </Page>
</template>
