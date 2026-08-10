<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElSkeleton,
  ElTag,
} from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  auditPlatformTransferApi,
  exportPlatformTransfersApi,
  getPlatformTransferApi,
  getPlatformTransferTitleApi,
  listPlatformTransfersApi,
  markPlatformTransferApi,
  payPlatformTransferApi,
  type PlatformTransferRow,
  type PlatformTransferTitle,
} from '#/api/core/platform-transfer';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const emptyTitle = (): PlatformTransferTitle => ({
  applying_amount: 0,
  applying_merchant_count: 0,
  freeze_amount: 0,
  payable_amount: 0,
  pending_amount: 0,
  withdrawable_amount: 0,
});

const title = ref<PlatformTransferTitle>(emptyTitle());
const lastFormValues = ref<Record<string, unknown>>({});
const canManage = ref(false);
const canExport = ref(false);
const exporting = ref(false);
const detailLoading = ref(false);
const current = ref<PlatformTransferRow | null>(null);

const auditForm = reactive({ status: 1 as number, refusal: '' });
const markForm = reactive({ admin_mark: '' });
const payForm = reactive({ image: '' });
const auditing = ref(false);
const marking = ref(false);
const paying = ref(false);

const summaryCards = computed(() => [
  {
    key: 'payable_amount',
    label: '应付商户金额',
    value: formatMoney(title.value.payable_amount),
    icon: 'lucide:wallet',
    tone: 'blue',
  },
  {
    key: 'withdrawable_amount',
    label: '商户可提现金额',
    value: formatMoney(title.value.withdrawable_amount),
    icon: 'lucide:circle-dollar-sign',
    tone: 'orange',
  },
  {
    key: 'applying_merchant_count',
    label: '申请转账的商户数',
    value: String(title.value.applying_merchant_count || 0),
    icon: 'lucide:store',
    tone: 'green',
  },
  {
    key: 'applying_amount',
    label: '申请转账的总金额',
    value: formatMoney(title.value.applying_amount),
    icon: 'lucide:badge-dollar-sign',
    tone: 'pink',
  },
  {
    key: 'pending_amount',
    label: '待审核的总金额',
    value: formatMoney(title.value.pending_amount),
    icon: 'lucide:clipboard-list',
    tone: 'purple',
  },
  {
    key: 'freeze_amount',
    label: '商户冻结金额',
    value: formatMoney(title.value.freeze_amount),
    icon: 'lucide:cloud',
    tone: 'sky',
  },
]);

function formatMoney(v?: number) {
  return Number(v || 0).toFixed(2);
}

function accountType(type: number) {
  return (
    ({ 1: '银行卡', 2: '微信', 3: '支付宝' } as Record<number, string>)[type] ||
    '未知'
  );
}

function traderLabel(v: number) {
  return Number(v) === 1 ? '自营' : '非自营';
}

function auditInfo(status: number) {
  if (status === 0) return { label: '待审核', type: 'warning' as const };
  if (status === 1) return { label: '审核通过', type: 'success' as const };
  if (status === -1) return { label: '审核未通过', type: 'danger' as const };
  return { label: '未知', type: 'info' as const };
}

function arrivalInfo(status: number) {
  return Number(status) === 1
    ? { label: '已到账', type: 'success' as const }
    : { label: '未到账', type: 'info' as const };
}

function canPay(row: PlatformTransferRow) {
  return canManage.value && row.status === 1 && row.financial_status === 0;
}

function parseAccount(raw?: string) {
  if (!raw) return null;
  try {
    return JSON.parse(raw) as Record<string, string>;
  } catch {
    return null;
  }
}

function accountLines(row: PlatformTransferRow) {
  const acc = parseAccount(row.financial_account);
  if (!acc) return [row.financial_account || '—'];
  if (row.financial_type === 1) {
    return [
      `姓名：${acc.name || '—'}`,
      `开户银行：${acc.bank || '—'}`,
      `银行卡号：${acc.bank_code || '—'}`,
    ];
  }
  if (row.financial_type === 2) {
    return [`姓名：${acc.name || '—'}`, `微信号：${acc.wechat || '—'}`];
  }
  if (row.financial_type === 3) {
    return [`姓名：${acc.name || '—'}`, `支付宝账号：${acc.alipay || '—'}`];
  }
  return [row.financial_account];
}

function buildFilterParams(formValues?: Record<string, unknown>) {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  const statusRaw = formValues?.status;
  const traderRaw = formValues?.is_trader;
  const typeRaw = formValues?.financial_type;
  const arrivalRaw = formValues?.financial_status;
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    status:
      statusRaw === 0 || statusRaw === 1 || statusRaw === -1
        ? Number(statusRaw)
        : undefined,
    mer_name: String(formValues?.mer_name ?? '').trim() || undefined,
    is_trader:
      traderRaw === 0 || traderRaw === 1 ? Number(traderRaw) : undefined,
    financial_type:
      typeRaw === 1 || typeRaw === 2 || typeRaw === 3
        ? Number(typeRaw)
        : undefined,
    financial_status:
      arrivalRaw === 0 || arrivalRaw === 1 ? Number(arrivalRaw) : undefined,
    admin_keyword: String(formValues?.admin_keyword ?? '').trim() || undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    {
      ...LIST_DATE_RANGE_FIELD,
      componentProps: {
        ...LIST_DATE_RANGE_FIELD.componentProps,
        startPlaceholder: '开始时间',
        endPlaceholder: '结束时间',
      },
      label: '时间选择',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '待审核', value: 0 },
          { label: '已审核', value: 1 },
          { label: '审核失败', value: -1 },
        ],
        placeholder: '请选择',
      },
      fieldName: 'status',
      label: '审核状态',
    },
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '请输入店铺名称' },
      fieldName: 'mer_name',
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
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '银行卡', value: 1 },
          { label: '微信', value: 2 },
          { label: '支付宝', value: 3 },
        ],
        placeholder: '请选择',
      },
      fieldName: 'financial_type',
      label: '收款方式',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '未到账', value: 0 },
          { label: '已到账', value: 1 },
        ],
        placeholder: '请选择',
      },
      fieldName: 'financial_status',
      label: '转账状态',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入管理员ID/姓名',
      },
      fieldName: 'admin_keyword',
      label: '管理员搜索',
    },
  ],
  {
    commonConfig: { componentProps: { class: 'w-full' } },
    submitButtonOptions: { content: '搜索' },
    handleSubmit: async (values) => {
      lastFormValues.value = { ...values };
      await gridApi.reload(values);
    },
    handleReset: async () => {
      await formApi.resetForm();
      const values = (await formApi.getValues()) ?? {};
      lastFormValues.value = { ...values };
      await gridApi.reload(values);
    },
  },
);

const [Form, formApi] = useVbenForm(formOptions);

const gridOptions: VxeGridProps<PlatformTransferRow> = {
  columns: [
    { type: 'seq', title: '序号', width: 70 },
    {
      field: 'is_trader',
      formatter: ({ cellValue }) => traderLabel(Number(cellValue)),
      minWidth: 90,
      title: '店铺类别',
    },
    {
      field: 'mer_name',
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '店铺名称',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '申请时间',
    },
    {
      field: 'extract_money',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 120,
      title: '转账金额(元)',
    },
    {
      field: 'admin_name',
      formatter: ({ cellValue }) =>
        cellValue === undefined || cellValue === null || cellValue === ''
          ? '—'
          : String(cellValue),
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '平台管理员姓名',
    },
    {
      field: 'financial_type',
      formatter: ({ cellValue }) => accountType(Number(cellValue)),
      minWidth: 100,
      title: '收款方式',
    },
    {
      field: 'status',
      slots: { default: 'auditStatus' },
      minWidth: 120,
      title: '审核状态',
    },
    {
      field: 'financial_status',
      slots: { default: 'arrivalStatus' },
      minWidth: 100,
      title: '到账状态',
    },
    {
      field: 'mer_money',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 120,
      title: '店铺余额(元)',
    },
    platformListActionColumn({ minWidth: 180, width: 220 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const values =
          formValues && Object.keys(formValues).length > 0
            ? formValues
            : lastFormValues.value;
        const filters = buildFilterParams(values);
        const [data, stats] = await Promise.all([
          listPlatformTransfersApi({
            page: page.currentPage,
            limit: page.pageSize,
            ...filters,
          }),
          getPlatformTransferTitleApi(filters).catch(() => emptyTitle()),
        ]);
        title.value = stats || emptyTitle();
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  cancelText: '关闭',
  title: '转账信息',
});

const [MarkDrawer, markDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  confirmText: '保存备注',
  cancelText: '取消',
  title: '备注',
  onConfirm: async () => saveMark(),
});

const [PayDrawer, payDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  confirmText: '确认转账',
  cancelText: '取消',
  title: '登记转账',
  onConfirm: async () => savePay(),
});

async function reloadGrid() {
  await gridApi.reload(lastFormValues.value);
}

async function openDetail(row: PlatformTransferRow) {
  current.value = null;
  detailLoading.value = true;
  auditForm.status = 1;
  auditForm.refusal = '';
  detailDrawerApi
    .setState({
      title: row.status === 0 ? '审核' : '转账信息',
      loading: true,
      showConfirmButton: false,
    })
    .open();
  try {
    current.value = await getPlatformTransferApi(row.financial_id);
  } catch {
    ElMessage.error('加载转账信息失败');
    detailDrawerApi.close();
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}

async function submitAudit() {
  if (!current.value || !canManage.value || auditing.value) return;
  if (auditForm.status === -1 && !auditForm.refusal.trim()) {
    ElMessage.warning('请输入拒绝理由');
    return;
  }
  auditing.value = true;
  try {
    await auditPlatformTransferApi(current.value.financial_id, {
      status: auditForm.status,
      refusal: auditForm.refusal.trim(),
    });
    ElMessage.success('审核完成');
    detailDrawerApi.close();
    await reloadGrid();
  } catch {
    // requestClient 已提示
  } finally {
    auditing.value = false;
  }
}

function openMark(row: PlatformTransferRow) {
  current.value = row;
  markForm.admin_mark = row.admin_mark || '';
  markDrawerApi.open();
}

async function saveMark() {
  if (!current.value || !canManage.value || marking.value) return;
  const mark = markForm.admin_mark.trim();
  if (!mark) {
    ElMessage.warning('请输入备注');
    return;
  }
  marking.value = true;
  try {
    await markPlatformTransferApi(current.value.financial_id, mark);
    ElMessage.success('备注成功');
    markDrawerApi.close();
    await reloadGrid();
  } catch {
    // requestClient 已提示
  } finally {
    marking.value = false;
  }
}

function openPay(row: PlatformTransferRow) {
  current.value = row;
  payForm.image = '';
  payDrawerApi.open();
}

async function savePay() {
  if (!current.value || !canManage.value || paying.value) return;
  if (!payForm.image.trim()) {
    ElMessage.warning('请上传转账凭证');
    return;
  }
  paying.value = true;
  try {
    await payPlatformTransferApi(current.value.financial_id, payForm.image.trim());
    ElMessage.success('转账登记完成');
    payDrawerApi.close();
    await reloadGrid();
  } catch {
    // requestClient 已提示
  } finally {
    paying.value = false;
  }
}

async function exportRows() {
  if (!canExport.value || exporting.value) return;
  exporting.value = true;
  try {
    const filters = buildFilterParams(lastFormValues.value);
    const result = await exportPlatformTransfersApi(filters);
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name || '转账记录.csv';
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      `已导出 ${result.row_count} 条${result.truncated ? '（已截断）' : ''}`,
    );
  } catch {
    ElMessage.error('导出失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canManage.value = codes.includes('accounts.transfer.manage');
  canExport.value = codes.includes('accounts.transfer.export');
});
</script>

<template>
  <Page auto-content-height>
    <div class="transfer-filter">
      <Form />
    </div>

    <div class="transfer-summary">
      <div class="transfer-summary__grid">
        <div
          v-for="card in summaryCards"
          :key="card.key"
          class="transfer-summary__card"
          :class="`transfer-summary__card--${card.tone}`"
        >
          <div class="transfer-summary__icon">
            <IconifyIcon :icon="card.icon" />
          </div>
          <div class="transfer-summary__body">
            <div class="transfer-summary__value">{{ card.value }}</div>
            <div class="transfer-summary__label">{{ card.label }}</div>
          </div>
        </div>
      </div>
    </div>

    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canExport"
          type="primary"
          :loading="exporting"
          @click="exportRows"
        >
          导出列表
        </ElButton>
      </template>
      <template #auditStatus="{ row }">
        <div>
          <ElTag :type="auditInfo(row.status).type" size="small">
            {{ auditInfo(row.status).label }}
          </ElTag>
          <div
            v-if="row.status === -1 && row.refusal"
            class="transfer-refusal"
          >
            原因：{{ row.refusal }}
          </div>
        </div>
      </template>
      <template #arrivalStatus="{ row }">
        <ElTag :type="arrivalInfo(row.financial_status).type" size="small">
          {{ arrivalInfo(row.financial_status).label }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">
          转账信息
        </ElButton>
        <ElButton
          v-if="canManage"
          link
          type="primary"
          @click="openMark(row)"
        >
          备注
        </ElButton>
        <ElButton
          v-if="canPay(row)"
          link
          type="primary"
          @click="openPay(row)"
        >
          转账
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="10">
        <div v-if="current" class="transfer-detail">
          <ElDescriptions :column="2" border>
            <ElDescriptionsItem label="单号">
              {{ current.financial_sn }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="店铺名称">
              {{ current.mer_name }}（{{ current.mer_id }}）
            </ElDescriptionsItem>
            <ElDescriptionsItem label="店铺类别">
              {{ traderLabel(current.is_trader) }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="转账金额">
              ¥{{ formatMoney(current.extract_money) }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="店铺余额">
              ¥{{ formatMoney(current.mer_money) }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="收款方式">
              {{ accountType(current.financial_type) }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="申请时间">
              {{ formatShanghaiDateTime(current.create_time) }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="审核状态">
              {{ auditInfo(current.status).label }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="到账状态">
              {{ arrivalInfo(current.financial_status).label }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="平台管理员">
              {{ current.admin_name || '—' }}
            </ElDescriptionsItem>
            <ElDescriptionsItem
              v-if="current.status_time"
              label="审核时间"
            >
              {{ formatShanghaiDateTime(current.status_time) }}
            </ElDescriptionsItem>
            <ElDescriptionsItem
              v-if="current.status === -1 && current.refusal"
              label="拒绝理由"
              :span="2"
            >
              {{ current.refusal }}
            </ElDescriptionsItem>
            <ElDescriptionsItem label="收款信息" :span="2">
              <div
                v-for="(line, idx) in accountLines(current)"
                :key="idx"
              >
                {{ line }}
              </div>
            </ElDescriptionsItem>
            <ElDescriptionsItem
              v-if="current.mark"
              label="商户备注"
              :span="2"
            >
              {{ current.mark }}
            </ElDescriptionsItem>
            <ElDescriptionsItem
              v-if="current.admin_mark"
              label="平台备注"
              :span="2"
            >
              {{ current.admin_mark }}
            </ElDescriptionsItem>
            <ElDescriptionsItem
              v-if="current.image"
              label="转账凭证"
              :span="2"
            >
              <ElImage
                :src="current.image"
                :preview-src-list="[current.image]"
                fit="cover"
                class="transfer-voucher"
              />
            </ElDescriptionsItem>
          </ElDescriptions>

          <div
            v-if="canManage && current.status === 0"
            class="transfer-audit"
          >
            <ElForm label-width="96px">
              <ElFormItem label="审核状态" required>
                <ElRadioGroup v-model="auditForm.status">
                  <ElRadio :label="1">通过</ElRadio>
                  <ElRadio :label="-1">拒绝</ElRadio>
                </ElRadioGroup>
              </ElFormItem>
              <ElFormItem
                v-if="auditForm.status === -1"
                label="原因"
                required
              >
                <ElInput
                  v-model="auditForm.refusal"
                  type="textarea"
                  :rows="3"
                  maxlength="200"
                  show-word-limit
                  placeholder="请输入拒绝理由"
                />
              </ElFormItem>
              <ElFormItem>
                <ElButton
                  type="primary"
                  :loading="auditing"
                  @click="submitAudit"
                >
                  确认审核
                </ElButton>
              </ElFormItem>
            </ElForm>
          </div>
        </div>
      </ElSkeleton>
    </DetailDrawer>

    <MarkDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="备注" required>
          <ElInput
            v-model="markForm.admin_mark"
            type="textarea"
            :rows="4"
            maxlength="255"
            show-word-limit
            placeholder="请输入备注"
          />
        </ElFormItem>
      </ElForm>
    </MarkDrawer>

    <PayDrawer>
      <div v-if="current" class="transfer-pay">
        <ElDescriptions :column="1" border class="mb-4">
          <ElDescriptionsItem label="店铺">
            {{ current.mer_name }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="转账金额">
            ¥{{ formatMoney(current.extract_money) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="收款方式">
            {{ accountType(current.financial_type) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="收款信息">
            <div
              v-for="(line, idx) in accountLines(current)"
              :key="idx"
            >
              {{ line }}
            </div>
          </ElDescriptionsItem>
        </ElDescriptions>
        <ElForm label-width="96px">
          <ElFormItem label="转账凭证" required>
            <ImageField v-model="payForm.image" />
          </ElFormItem>
        </ElForm>
      </div>
    </PayDrawer>
  </Page>
</template>

<style scoped>
.transfer-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.transfer-summary {
  padding: 16px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.transfer-summary__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  width: 100%;
}

.transfer-summary__card {
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

.transfer-summary__icon {
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

.transfer-summary__card--blue .transfer-summary__icon {
  background: #409eff;
}

.transfer-summary__card--orange .transfer-summary__icon {
  background: #e6a23c;
}

.transfer-summary__card--green .transfer-summary__icon {
  background: #67c23a;
}

.transfer-summary__card--pink .transfer-summary__icon {
  background: #f56c6c;
}

.transfer-summary__card--purple .transfer-summary__icon {
  background: #9b59b6;
}

.transfer-summary__card--sky .transfer-summary__icon {
  background: #36cfc9;
}

.transfer-summary__value {
  color: var(--el-text-color-primary);
  font-size: 24px;
  font-weight: 600;
  line-height: 1.2;
}

.transfer-summary__label {
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.2;
}

.transfer-refusal {
  margin-top: 4px;
  color: var(--el-color-danger);
  font-size: 12px;
  line-height: 1.3;
}

.transfer-detail {
  padding-bottom: 8px;
}

.transfer-audit {
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--el-border-color-lighter);
}

.transfer-voucher {
  width: 120px;
  height: 120px;
  border-radius: 6px;
}

.transfer-pay {
  padding-bottom: 8px;
}

@media (min-width: 1600px) {
  .transfer-summary__grid {
    grid-template-columns: repeat(6, minmax(0, 1fr));
  }
}

@media (max-width: 1100px) {
  .transfer-summary__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .transfer-summary__grid {
    grid-template-columns: 1fr;
  }
}
</style>
