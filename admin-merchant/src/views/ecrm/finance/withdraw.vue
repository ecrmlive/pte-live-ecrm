<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  applyMerchantWithdrawApi,
  getMerchantBalanceApi,
  getMerchantWithdrawApi,
  listMerchantWithdrawsApi,
  type MerchantWithdraw,
} from '#/api/core/merchant-finance';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const balance = ref(0);
const saving = ref(false);
const detail = ref<MerchantWithdraw>();
const form = reactive({
  extract_money: 0,
  financial_account: '',
  financial_type: 1,
  mark: '',
});

const auditStatus: Record<
  number,
  { label: string; type: 'danger' | 'info' | 'success' | 'warning' }
> = {
  [-1]: { label: '审核拒绝', type: 'danger' },
  0: { label: '待平台审核', type: 'warning' },
  1: { label: '审核通过', type: 'success' },
};

function auditInfo(status: number) {
  return auditStatus[status] || { label: '未知状态', type: 'info' as const };
}

function accountType(type: number) {
  return ({ 1: '银行卡', 2: '微信', 3: '支付宝' }[type] || '未知');
}

function transferInfo(row: MerchantWithdraw) {
  return row.financial_status === 1
    ? { label: '已打款', type: 'success' as const }
    : { label: '未打款', type: 'info' as const };
}

async function refreshBalance() {
  const currentBalance = await getMerchantBalanceApi();
  balance.value = currentBalance.mer_money;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
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
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '审核状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '未打款', value: 0 },
        { label: '已打款', value: 1 },
      ],
      placeholder: '全部',
    },
    fieldName: 'financial_status',
    label: '打款状态',
  },
]);

const gridOptions: VxeGridProps<MerchantWithdraw> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    {
      field: 'financial_sn',
      minWidth: 180,
      showOverflow: false,
      title: '申请单号',
    },
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
      slots: { default: 'audit' },
      title: '审核状态',
      width: 116,
    },
    {
      field: 'financial_status',
      slots: { default: 'transfer' },
      title: '打款状态',
      width: 100,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '申请时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 76 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = formValues?.status;
        const financialStatus = formValues?.financial_status;
        await refreshBalance();
        const data = await listMerchantWithdrawsApi({
          page: page.currentPage,
          limit: page.pageSize,
          financial_sn:
            String(formValues?.financial_sn ?? '').trim() || undefined,
          status:
            status === -1 || status === 0 || status === 1
              ? Number(status)
              : undefined,
          financial_status:
            financialStatus === 0 || financialStatus === 1
              ? Number(financialStatus)
              : undefined,
          date_from: range[0],
          date_to: range[1],
        });
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

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[560px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

const [ApplyModal, applyModalApi] = useVbenModal({
  title: '申请提现',
  confirmText: '提交申请',
  onConfirm: async () => {
    if (form.extract_money <= 0) {
      ElMessage.warning('提现金额必须大于 0');
      return;
    }
    if (form.extract_money > balance.value) {
      ElMessage.warning('提现金额不能超过可用余额');
      return;
    }
    if (!form.financial_account.trim()) {
      ElMessage.warning('请填写收款账户');
      return;
    }
    saving.value = true;
    applyModalApi.lock();
    try {
      await applyMerchantWithdrawApi({
        ...form,
        financial_account: form.financial_account.trim(),
        mark: form.mark.trim(),
      });
      ElMessage.success('提现申请已提交，等待平台审核');
      applyModalApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      applyModalApi.unlock();
    }
  },
});

function openApply() {
  Object.assign(form, {
    extract_money: 0,
    financial_account: '',
    financial_type: 1,
    mark: '',
  });
  applyModalApi.open();
}

async function openDetail(row: MerchantWithdraw) {
  detailDrawerApi.setState({ title: '提现详情', loading: true }).open();
  try {
    detail.value = await getMerchantWithdrawApi(row.financial_id);
  } finally {
    detailDrawerApi.setState({ loading: false });
  }
}

onMounted(() => void refreshBalance());
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="flex flex-wrap items-center gap-3">
          <div>
            <span class="text-sm text-muted-foreground">当前可用余额</span>
            <span class="ml-3 text-xl font-semibold text-primary">
              ¥{{ Number(balance).toFixed(2) }}
            </span>
          </div>
          <ElButton type="primary" @click="openApply">申请提现</ElButton>
        </div>
      </template>
      <template #audit="{ row }">
        <ElTag :type="auditInfo(row.status).type" size="small">
          {{ auditInfo(row.status).label }}
        </ElTag>
      </template>
      <template #transfer="{ row }">
        <ElTag :type="transferInfo(row).type" size="small">
          {{ transferInfo(row).label }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="detail">
        <ElDescriptions :column="1" border>
          <ElDescriptionsItem label="申请单号">
            {{ detail.financial_sn }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="提现金额">
            ¥{{ Number(detail.extract_money).toFixed(2) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="收款方式">
            {{ accountType(detail.financial_type) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="收款账户">
            {{ detail.financial_account }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="审核状态">
            <ElTag :type="auditInfo(detail.status).type">
              {{ auditInfo(detail.status).label }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="打款状态">
            <ElTag :type="transferInfo(detail).type">
              {{ transferInfo(detail).label }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="detail.refusal" label="拒绝原因">
            {{ detail.refusal }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="申请备注">
            {{ detail.mark || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="申请时间">
            {{ formatShanghaiDateTime(detail.create_time) }}
          </ElDescriptionsItem>
        </ElDescriptions>
      </template>
    </DetailDrawer>

    <ApplyModal>
      <ElForm label-width="96px">
        <ElFormItem label="可用余额">
          <span>¥{{ Number(balance).toFixed(2) }}</span>
        </ElFormItem>
        <ElFormItem label="提现金额" required>
          <ElInputNumber
            v-model="form.extract_money"
            :min="0.01"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="收款方式" required>
          <ElSelect v-model="form.financial_type" class="w-full">
            <ElOption :value="1" label="银行卡" />
            <ElOption :value="2" label="微信" />
            <ElOption :value="3" label="支付宝" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="收款账户" required>
          <ElInput
            v-model="form.financial_account"
            placeholder="银行卡号、微信号或支付宝账号"
          />
        </ElFormItem>
        <ElFormItem label="申请备注">
          <ElInput
            v-model="form.mark"
            :rows="3"
            maxlength="200"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </ApplyModal>
  </Page>
</template>
