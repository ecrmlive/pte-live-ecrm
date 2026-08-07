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
  ElRadio,
  ElRadioGroup,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  fetchProfitsharingApplication,
  fetchProfitsharingApplications,
  reviewProfitsharingApplication,
  saveProfitsharingApplicationNote,
  type ProfitsharingApplication,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const can = ref(false);
const statusFilter = ref('');
const submitting = ref(false);
const selected = ref<ProfitsharingApplication>();
const detail = reactive<Partial<ProfitsharingApplication>>({});
const reviewForm = reactive({
  note: '',
  status: 'auditing',
});

/** 对齐 CRMEB 平台审核：通过→审核中，拒绝→平台驳回 */
const reviewStatusOptions = [
  { label: '通过', value: 'auditing' },
  { label: '拒绝', value: 'platform_rejected' },
] as const;

/** CRMEB: 0待审 / -1平台驳回 / 10审核中 / 11店铺验证 / 20完成 / 30冻结 / 40微信驳回 */
const statusTabs = [
  { label: '全部', value: '' },
  { label: '待审核', value: 'applied' },
  { label: '平台驳回', value: 'platform_rejected' },
  { label: '审核中', value: 'auditing' },
  { label: '店铺验证', value: 'shop_verify' },
  { label: '已完成', value: 'completed' },
  { label: '已冻结', value: 'frozen' },
  { label: '微信驳回', value: 'wechat_rejected' },
] as const;

const statusLabels: Record<string, string> = {
  applied: '待审核',
  approved: '已完成',
  rejected: '平台驳回',
  platform_rejected: '平台驳回',
  auditing: '审核中',
  shop_verify: '店铺验证',
  completed: '已完成',
  frozen: '已冻结',
  wechat_rejected: '微信驳回',
};

const statusTagType: Record<string, 'success' | 'warning' | 'danger' | 'info'> =
  {
    applied: 'warning',
    auditing: 'info',
    shop_verify: 'info',
    completed: 'success',
    approved: 'success',
    frozen: 'danger',
    platform_rejected: 'danger',
    rejected: 'danger',
    wechat_rejected: 'danger',
  };

function formatTime(value?: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    status: statusFilter.value || undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

function setStatusTab(value: string) {
  if (statusFilter.value === value) return;
  statusFilter.value = value;
  gridApi.reload();
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入店铺名称',
    },
    fieldName: 'keyword',
    label: '店铺名称',
  },
]);

const gridOptions: VxeGridProps<ProfitsharingApplication> = {
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'applyment_id',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 180,
      title: '微信支付申请单号',
    },
    { field: 'application_no', minWidth: 160, title: '业务申请编号' },
    { field: 'merchant_id', title: '店铺ID', width: 90 },
    {
      field: 'merchant_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 140,
      title: '店铺名称',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 110,
    },
    {
      field: 'message',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 160,
      showOverflow: false,
      title: '审核说明',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 170,
      title: '申请时间',
    },
    {
      field: 'review_note',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 140,
      showOverflow: false,
      title: '备注',
    },
    platformListActionColumn({ width: 180 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchProfitsharingApplications(
          buildListParams(page, formValues),
        );
        const list = data.list || [];
        return { items: list, total: Number(data.total ?? 0) };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
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
  footer: false,
  title: '分账申请详情',
});

const [ReviewDrawer, reviewDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => submitReview(),
  title: '审核分账申请',
});

async function openDetail(row: ProfitsharingApplication) {
  const data = await fetchProfitsharingApplication(row.id);
  Object.assign(detail, data);
  detailDrawerApi.open();
}

async function openNote(row: ProfitsharingApplication) {
  try {
    const { value } = await ElMessageBox.prompt('填写备注。', '备注', {
      inputValue: row.review_note || '',
      inputValidator: (v: string) => {
        const note = String(v ?? '').trim();
        if (!note) return '备注不能为空';
        return [...note].length <= 500 ? true : '备注不能超过 500 个字符';
      },
    });
    await saveProfitsharingApplicationNote(row.id, value.trim());
    ElMessage.success('备注已保存');
    gridApi.reload();
  } catch (error) {
    if (error === 'cancel' || error === 'close' || error === 'escape') return;
    throw error;
  }
}

function openReview(row: ProfitsharingApplication) {
  selected.value = row;
  reviewForm.status = 'auditing';
  reviewForm.note = '';
  reviewDrawerApi.open();
}

function noteRequired() {
  return reviewForm.status === 'platform_rejected';
}

async function submitReview() {
  if (!selected.value) return;
  const note = reviewForm.note.trim();
  if (noteRequired() && !note) {
    ElMessage.warning('驳回时必须填写审核说明');
    return;
  }
  if ([...note].length > 500) {
    ElMessage.warning('审核说明不能超过 500 个字符');
    return;
  }
  reviewDrawerApi.lock();
  submitting.value = true;
  try {
    await reviewProfitsharingApplication(selected.value.id, {
      note,
      status: reviewForm.status,
    });
    ElMessage.success('审核已保存');
    reviewDrawerApi.close();
    gridApi.reload();
  } finally {
    submitting.value = false;
    reviewDrawerApi.unlock();
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  can.value = codes.includes('merchant.profitsharing.review');
});
</script>

<template>
  <Page auto-content-height>
    <div class="applyments-page">
      <div class="applyments-status-row">
        <span class="applyments-status-row__label">审核状态：</span>
        <div class="applyments-status-tabs" role="tablist">
          <button
            v-for="tab in statusTabs"
            :key="tab.value || 'all'"
            type="button"
            role="tab"
            class="applyments-status-tabs__item"
            :aria-selected="statusFilter === tab.value"
            :class="{ 'is-active': statusFilter === tab.value }"
            @click="setStatusTab(tab.value)"
          >
            {{ tab.label }}
          </button>
        </div>
      </div>

      <Grid>
        <template #status="{ row }">
          <ElTag :type="statusTagType[row.status] || 'info'" size="small">
            {{ statusLabels[row.status] || row.status }}
          </ElTag>
        </template>

        <template #action="{ row }">
          <template v-if="can">
            <ElButton link type="primary" @click="openNote(row)">
              备注
            </ElButton>
            <ElButton link type="primary" @click="openDetail(row)">
              详情
            </ElButton>
            <ElButton
              v-if="row.status === 'applied'"
              link
              type="primary"
              @click="openReview(row)"
            >
              审核
            </ElButton>
          </template>
          <span v-else>—</span>
        </template>
      </Grid>
    </div>

    <ReviewDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="申请店铺">
          <span>{{ selected?.merchant_name || '—' }}</span>
        </ElFormItem>
        <ElFormItem label="审核状态" required>
          <ElRadioGroup v-model="reviewForm.status">
            <ElRadio
              v-for="opt in reviewStatusOptions"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem :required="noteRequired()" label="审核说明">
          <ElInput
            v-model="reviewForm.note"
            :rows="3"
            maxlength="500"
            show-word-limit
            type="textarea"
            :placeholder="
              noteRequired() ? '请填写驳回原因' : '选填，可填写审核说明'
            "
          />
        </ElFormItem>
      </ElForm>
    </ReviewDrawer>

    <DetailDrawer class="w-[560px]">
      <ElDescriptions :column="1" border>
        <ElDescriptionsItem label="ID">{{ detail.id }}</ElDescriptionsItem>
        <ElDescriptionsItem label="微信支付申请单号">
          {{ detail.applyment_id || '—' }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="业务申请编号">
          {{ detail.application_no || '—' }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="店铺ID">
          {{ detail.merchant_id }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="店铺名称">
          {{ detail.merchant_name || '—' }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="状态">
          {{ statusLabels[detail.status || ''] || detail.status || '—' }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="审核说明">
          {{ detail.message || '—' }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="申请说明">
          {{ detail.description || '—' }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="备注">
          {{ detail.review_note || '—' }}
        </ElDescriptionsItem>
        <ElDescriptionsItem label="申请时间">
          {{ formatTime(detail.created_at) }}
        </ElDescriptionsItem>
      </ElDescriptions>
    </DetailDrawer>
  </Page>
</template>

<style scoped>
.applyments-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
}

.applyments-status-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: none;
  margin-bottom: 12px;
  padding: 8px 16px;
  border-radius: 8px;
  background: hsl(var(--card));
  border: 1px solid hsl(var(--border));
}

.applyments-status-row__label {
  flex: none;
  color: hsl(var(--foreground) / 80%);
  font-size: 14px;
  line-height: 22px;
  white-space: nowrap;
}

.applyments-status-tabs {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  gap: 8px;
  min-width: 0;
}

.applyments-status-tabs__item {
  padding: 4px 12px;
  border: 1px solid transparent;
  border-radius: 4px;
  background: transparent;
  color: hsl(var(--foreground) / 75%);
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
}

.applyments-status-tabs__item:hover {
  color: hsl(var(--primary));
}

.applyments-status-tabs__item.is-active {
  border-color: hsl(var(--primary));
  color: hsl(var(--primary));
  font-weight: 600;
}

.applyments-page :deep(.vben-vxe-grid),
.applyments-page :deep(.bg-card) {
  flex: 1;
  min-height: 0;
}
</style>
