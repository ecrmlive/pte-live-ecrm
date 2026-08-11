<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElSkeleton,
} from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  exportPlatformUserExtractsApi,
  getPlatformUserExtractApi,
  listPlatformUserExtractsApi,
  switchPlatformUserExtractStatusApi,
  type PlatformUserExtractQuery,
  type PlatformUserExtractRow,
} from '#/api/core/platform-user-extract';
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

const lastFormValues = ref<Record<string, unknown>>({});
const canReview = ref(false);
const canExport = ref(false);
const exporting = ref(false);
const detailLoading = ref(false);
const auditing = ref(false);
const current = ref<PlatformUserExtractRow | null>(null);
const auditForm = reactive({
  status: 1 as 1 | -1,
  fail_msg: '',
  mark: '',
});

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || String(v).trim() === '') return '-';
  return String(v);
}

function formatMoney(v?: number) {
  return Number(v || 0).toFixed(2);
}

function mediaUrl(url?: string) {
  return resolveCosMediaUrl(String(url || '').trim());
}

function typeLabel(t: number) {
  return (
    ({ 0: '银行卡', 1: '微信', 2: '支付宝', 3: '微信零钱', 4: '余额' } as Record<
      number,
      string
    >)[t] || '未知'
  );
}

function statusLabel(status: number) {
  if (status === 1) return '已通过';
  if (status === -1) return '已拒绝';
  return '审核中';
}

function accountOf(row: PlatformUserExtractRow) {
  if (row.account) return row.account;
  switch (row.extract_type) {
    case 0:
      return row.bank_code || '';
    case 1:
    case 3:
      return row.wechat || '';
    case 2:
      return row.alipay_code || '';
    default:
      return '';
  }
}

function buildFilterParams(
  formValues?: Record<string, unknown>,
): PlatformUserExtractQuery {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  const statusRaw = formValues?.status;
  const typeRaw = formValues?.extract_type;
  const allowedStatus = [0, 1, -1];
  const allowedType = [0, 1, 2, 3, 4];
  const userSearch = parseUserSearch(formValues);
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    status: allowedStatus.includes(Number(statusRaw))
      ? Number(statusRaw)
      : undefined,
    extract_type: allowedType.includes(Number(typeRaw))
      ? Number(typeRaw)
      : undefined,
    user_type: userSearch.type || 'nickname',
    user_keyword: userSearch.keyword || undefined,
    account_keyword:
      String(formValues?.account_keyword ?? '').trim() || undefined,
  };
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
          { label: '审核中', value: 0 },
          { label: '已通过', value: 1 },
          { label: '已拒绝', value: -1 },
        ],
        placeholder: '全部',
      },
      fieldName: 'status',
      label: '提现状态',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '银行卡', value: 0 },
          { label: '微信', value: 1 },
          { label: '支付宝', value: 2 },
          { label: '微信零钱', value: 3 },
          { label: '余额', value: 4 },
        ],
        placeholder: '全部',
      },
      fieldName: 'extract_type',
      label: '提现方式',
    },
    listUserSearchFormField(),
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '银行卡号/支付宝账号/微信号',
      },
      fieldName: 'account_keyword',
      label: '账号搜索',
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
      await gridApi.reload(lastFormValues.value);
    },
  },
);

const [Form, formApi] = useVbenForm(formOptions);

const gridOptions: VxeGridProps<PlatformUserExtractRow> = {
  columns: [
    { type: 'seq', title: '序号', width: 70 },
    {
      field: 'extract_pic',
      slots: { default: 'qrcode' },
      title: '二维码',
      width: 90,
    },
    {
      field: 'nickname',
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '用户信息',
    },
    {
      field: 'uid',
      minWidth: 100,
      title: '用户UID',
    },
    {
      field: 'real_name',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 100,
      title: '户名',
    },
    {
      field: 'extract_price',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 100,
      title: '提现金额',
    },
    {
      field: 'extract_type',
      formatter: ({ cellValue }) => typeLabel(Number(cellValue)),
      minWidth: 100,
      title: '提现方式',
    },
    {
      field: 'bank_name',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 110,
      showOverflow: 'tooltip',
      title: '银行名称',
    },
    {
      field: 'account',
      formatter: ({ row }) => dash(accountOf(row)),
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '账号',
    },
    {
      field: 'status',
      formatter: ({ cellValue }) => statusLabel(Number(cellValue)),
      minWidth: 90,
      title: '提现状态',
    },
    {
      field: 'fail_msg',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '拒绝原因',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '添加时间',
    },
    platformListActionColumn({ minWidth: 100, width: 120 }),
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
        const data = await listPlatformUserExtractsApi({
          page: page.currentPage,
          limit: page.pageSize,
          ...filters,
        });
        return {
          items: data.list || [],
          total: data.total || 0,
        };
      },
    },
  },
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
  footer: false,
  title: '提现详情',
});

const [AuditDrawer, auditDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  confirmText: '提交审核',
  cancelText: '取消',
  title: '提现审核',
  onConfirm: async () => submitAudit(),
});

async function reloadGrid() {
  await gridApi.reload(lastFormValues.value);
}

async function openDetail(row: PlatformUserExtractRow) {
  current.value = row;
  detailLoading.value = true;
  detailDrawerApi.open();
  try {
    current.value = await getPlatformUserExtractApi(row.extract_id);
  } catch {
    // requestClient 已提示
  } finally {
    detailLoading.value = false;
  }
}

async function openAudit(row: PlatformUserExtractRow) {
  if (!canReview.value || row.status !== 0) return;
  current.value = row;
  auditForm.status = 1;
  auditForm.fail_msg = '';
  auditForm.mark = row.mark || '';
  detailLoading.value = true;
  auditDrawerApi.open();
  try {
    current.value = await getPlatformUserExtractApi(row.extract_id);
    auditForm.mark = current.value.mark || '';
  } catch {
    // requestClient 已提示
  } finally {
    detailLoading.value = false;
  }
}

async function submitAudit() {
  if (!current.value || !canReview.value || auditing.value) return;
  if (auditForm.status === -1 && !auditForm.fail_msg.trim()) {
    ElMessage.warning('请填写拒绝原因');
    return;
  }
  auditing.value = true;
  try {
    await switchPlatformUserExtractStatusApi(current.value.extract_id, {
      status: auditForm.status,
      fail_msg: auditForm.status === -1 ? auditForm.fail_msg.trim() : '',
      mark: auditForm.mark.trim(),
    });
    ElMessage.success('审核成功');
    auditDrawerApi.close();
    await reloadGrid();
  } catch {
    // requestClient 已提示
  } finally {
    auditing.value = false;
  }
}

async function exportRows() {
  if (!canExport.value || exporting.value) return;
  exporting.value = true;
  try {
    const filters = buildFilterParams(lastFormValues.value);
    const result = await exportPlatformUserExtractsApi(filters);
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name || '提现管理.csv';
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
  canReview.value = codes.includes('accounts.withdraw.review');
  canExport.value = codes.includes('accounts.withdraw.export');
});
</script>

<template>
  <Page auto-content-height>
    <div class="withdraw-filter">
      <Form />
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
      <template #qrcode="{ row }">
        <ElImage
          v-if="mediaUrl(row.extract_pic)"
          :src="mediaUrl(row.extract_pic)"
          :preview-src-list="[mediaUrl(row.extract_pic)]"
          fit="cover"
          class="withdraw-qr"
        />
        <span v-else>-</span>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="row.status === 0 && canReview"
          link
          type="primary"
          @click="openAudit(row)"
        >
          审核
        </ElButton>
        <ElButton v-else link type="primary" @click="openDetail(row)">
          详情
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <template #default>
          <div v-if="current" class="withdraw-detail">
            <section class="withdraw-section">
              <div class="withdraw-section__title">用户信息</div>
              <div class="withdraw-kv-grid">
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">用户昵称：</span>
                  <span class="withdraw-kv__value">{{
                    dash(current.nickname)
                  }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">用户ID：</span>
                  <span class="withdraw-kv__value">{{ current.uid }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">提现金额：</span>
                  <span class="withdraw-kv__value">{{
                    formatMoney(current.extract_price)
                  }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">申请时间：</span>
                  <span class="withdraw-kv__value">{{
                    formatShanghaiDateTime(current.create_time) || '-'
                  }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">账号：</span>
                  <span class="withdraw-kv__value">{{
                    dash(accountOf(current))
                  }}</span>
                </div>
                <div v-if="current.real_name" class="withdraw-kv">
                  <span class="withdraw-kv__label">户名：</span>
                  <span class="withdraw-kv__value">{{
                    dash(current.real_name)
                  }}</span>
                </div>
              </div>
            </section>

            <section class="withdraw-section withdraw-section--last">
              <div class="withdraw-section__title">提现方式</div>
              <div class="withdraw-kv-grid">
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">提现方式：</span>
                  <span class="withdraw-kv__value">{{
                    typeLabel(current.extract_type)
                  }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">审核状态：</span>
                  <span class="withdraw-kv__value">{{
                    statusLabel(current.status)
                  }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">审核时间：</span>
                  <span class="withdraw-kv__value">{{
                    formatShanghaiDateTime(current.status_time) || '-'
                  }}</span>
                </div>
                <div
                  v-if="current.extract_type === 0 && current.bank_name"
                  class="withdraw-kv"
                >
                  <span class="withdraw-kv__label">银行名称：</span>
                  <span class="withdraw-kv__value">{{
                    dash(current.bank_name)
                  }}</span>
                </div>
                <div v-if="current.status === -1" class="withdraw-kv">
                  <span class="withdraw-kv__label">拒绝原因：</span>
                  <span class="withdraw-kv__value">{{
                    dash(current.fail_msg)
                  }}</span>
                </div>
                <div v-if="current.mark" class="withdraw-kv">
                  <span class="withdraw-kv__label">管理员备注：</span>
                  <span class="withdraw-kv__value">{{ dash(current.mark) }}</span>
                </div>
              </div>
              <div
                v-if="mediaUrl(current.extract_pic)"
                class="withdraw-kv withdraw-kv--block"
              >
                <span class="withdraw-kv__label">收款码：</span>
                <ElImage
                  :src="mediaUrl(current.extract_pic)"
                  :preview-src-list="[mediaUrl(current.extract_pic)]"
                  fit="cover"
                  class="withdraw-qr withdraw-qr--lg"
                />
              </div>
            </section>
          </div>
        </template>
      </ElSkeleton>
    </DetailDrawer>

    <AuditDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <template #default>
          <div v-if="current" class="withdraw-detail">
            <section class="withdraw-section">
              <div class="withdraw-section__title">用户信息</div>
              <div class="withdraw-kv-grid">
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">用户昵称：</span>
                  <span class="withdraw-kv__value">{{
                    dash(current.nickname)
                  }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">用户ID：</span>
                  <span class="withdraw-kv__value">{{ current.uid }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">提现金额：</span>
                  <span class="withdraw-kv__value">{{
                    formatMoney(current.extract_price)
                  }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">申请时间：</span>
                  <span class="withdraw-kv__value">{{
                    formatShanghaiDateTime(current.create_time) || '-'
                  }}</span>
                </div>
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">账号：</span>
                  <span class="withdraw-kv__value">{{
                    dash(accountOf(current))
                  }}</span>
                </div>
                <div v-if="current.real_name" class="withdraw-kv">
                  <span class="withdraw-kv__label">户名：</span>
                  <span class="withdraw-kv__value">{{
                    dash(current.real_name)
                  }}</span>
                </div>
              </div>
            </section>

            <section class="withdraw-section">
              <div class="withdraw-section__title">提现方式</div>
              <div class="withdraw-kv-grid">
                <div class="withdraw-kv">
                  <span class="withdraw-kv__label">提现方式：</span>
                  <span class="withdraw-kv__value">{{
                    typeLabel(current.extract_type)
                  }}</span>
                </div>
                <div
                  v-if="current.extract_type === 0 && current.bank_name"
                  class="withdraw-kv"
                >
                  <span class="withdraw-kv__label">银行名称：</span>
                  <span class="withdraw-kv__value">{{
                    dash(current.bank_name)
                  }}</span>
                </div>
              </div>
            </section>

            <section class="withdraw-section withdraw-section--last">
              <div class="withdraw-section__title">审核操作</div>
              <ElForm class="withdraw-audit-form" label-width="96px">
                <ElFormItem label="审核结果" required>
                  <ElRadioGroup v-model="auditForm.status">
                    <ElRadio :value="1">通过</ElRadio>
                    <ElRadio :value="-1">拒绝</ElRadio>
                  </ElRadioGroup>
                </ElFormItem>
                <ElFormItem
                  v-if="auditForm.status === -1"
                  label="拒绝原因"
                  required
                >
                  <ElInput
                    v-model="auditForm.fail_msg"
                    type="textarea"
                    :rows="3"
                    maxlength="200"
                    show-word-limit
                    placeholder="请填写拒绝原因"
                  />
                </ElFormItem>
                <ElFormItem label="管理员备注">
                  <ElInput
                    v-model="auditForm.mark"
                    type="textarea"
                    :rows="3"
                    maxlength="500"
                    show-word-limit
                    placeholder="选填"
                  />
                </ElFormItem>
              </ElForm>
            </section>
          </div>
        </template>
      </ElSkeleton>
    </AuditDrawer>
  </Page>
</template>

<style scoped>
.withdraw-filter {
  margin-bottom: 12px;
}

.withdraw-qr {
  width: 40px;
  height: 40px;
  border-radius: 4px;
}

.withdraw-qr--lg {
  width: 120px;
  height: 120px;
}

.withdraw-detail {
  padding: 4px 0 8px;
}

.withdraw-section {
  padding: 2px 0 18px;
  border-bottom: 1px dashed var(--el-border-color);
}

.withdraw-section--last {
  border-bottom: none;
}

.withdraw-section__title {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 14px;
  color: var(--el-text-color-primary);
  font-size: 14px;
  font-weight: 600;
  line-height: 22px;
}

.withdraw-section__title::before {
  content: '';
  width: 3px;
  height: 14px;
  border-radius: 1px;
  background: var(--el-color-primary);
}

.withdraw-kv-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px 24px;
}

.withdraw-kv {
  display: flex;
  gap: 4px;
  min-width: 0;
  font-size: 13px;
  line-height: 22px;
}

.withdraw-kv--block {
  margin-top: 12px;
  align-items: flex-start;
}

.withdraw-kv__label {
  flex-shrink: 0;
  color: var(--el-text-color-regular);
}

.withdraw-kv__value {
  min-width: 0;
  color: var(--el-text-color-primary);
  word-break: break-all;
}

.withdraw-audit-form {
  max-width: 640px;
}
</style>
