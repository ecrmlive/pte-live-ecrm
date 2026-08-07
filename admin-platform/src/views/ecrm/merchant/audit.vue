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
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  auditMerchantIntention,
  fetchMerchantCategories,
  fetchMerchantIntentions,
  fetchMerchantTypes,
  type MerchantCategoryRow,
  type MerchantIntentionRow,
  type MerchantTypeRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const selected = ref<MerchantIntentionRow>();
const submitting = ref(false);
const canAudit = ref(false);
const categories = ref<MerchantCategoryRow[]>([]);
const types = ref<MerchantTypeRow[]>([]);

const form = reactive({
  fail_msg: '',
  mark: '',
  status: 1,
});

const statusText = (status: number) =>
  ({ 0: '待审核', 1: '审核通过', 2: '审核未通过' })[status] || '未知';
const statusType = (status: number) =>
  ({ 0: 'warning', 1: 'success', 2: 'danger' })[status] || 'info';
const formatTime = (value?: string | null) =>
  value ? formatShanghaiDateTime(value) : '—';
const displayOrDash = (value?: string | number | null) => {
  const text = String(value ?? '').trim();
  return text || '—';
};

function imageUrls(raw?: string) {
  const value = String(raw ?? '').trim();
  if (!value) return [];
  return value
    .split(/[,;|]/)
    .map((item) => resolveCosMediaUrl(item.trim()))
    .filter(Boolean)
    .slice(0, 6);
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const statusRaw = formValues?.status;
  const categoryRaw = formValues?.category_id;
  const typeRaw = formValues?.type_id;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    status:
      statusRaw === 0 || statusRaw === 1 || statusRaw === 2
        ? Number(statusRaw)
        : undefined,
    category_id: categoryRaw ? Number(categoryRaw) : undefined,
    type_id: typeRaw ? Number(typeRaw) : undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 0 },
        { label: '审核通过', value: 1 },
        { label: '审核未通过', value: 2 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '审核状态',
  },
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Select',
    componentProps: { clearable: true, options: [], placeholder: '请选择' },
    fieldName: 'category_id',
    label: '店铺分类',
  },
  {
    component: 'Select',
    componentProps: { clearable: true, options: [], placeholder: '请选择' },
    fieldName: 'type_id',
    label: '店铺类型',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入店铺名称/联系方式',
    },
    fieldName: 'keyword',
    label: '关键字',
  },
]);

const gridOptions: VxeGridProps<MerchantIntentionRow> = {
  columns: [
    { field: 'mer_intention_id', title: 'ID', width: 72 },
    {
      field: 'mer_name',
      formatter: ({ cellValue }) => displayOrDash(cellValue),
      minWidth: 140,
      showOverflow: false,
      title: '店铺名称',
    },
    {
      field: 'category_name',
      formatter: ({ cellValue }) => displayOrDash(cellValue),
      minWidth: 110,
      showOverflow: false,
      title: '店铺分类',
    },
    {
      field: 'type_name',
      formatter: ({ cellValue }) => displayOrDash(cellValue),
      minWidth: 110,
      showOverflow: false,
      title: '店铺类型',
    },
    {
      field: 'name',
      formatter: ({ cellValue }) => displayOrDash(cellValue),
      minWidth: 100,
      showOverflow: false,
      title: '店铺姓名',
    },
    {
      field: 'phone',
      formatter: ({ cellValue }) => displayOrDash(cellValue),
      minWidth: 120,
      showOverflow: false,
      title: '联系方式',
    },
    {
      field: 'images',
      slots: { default: 'images' },
      title: '资质图片',
      width: 120,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '审核状态',
      width: 108,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatTime(cellValue),
      title: '申请时间',
      width: 168,
    },
    {
      className: 'col--remark',
      field: 'mark',
      formatter: ({ cellValue }) => displayOrDash(cellValue),
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '审核备注',
      width: 180,
    },
    platformListActionColumn({ width: 120 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchMerchantIntentions(buildListParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'mer_intention_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [AuditDrawer, auditDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => submitAudit(),
});

function syncFilterSelectOptions() {
  gridApi.formApi?.updateSchema([
    {
      fieldName: 'category_id',
      componentProps: {
        clearable: true,
        options: categories.value.map((c) => ({
          label: c.category_name,
          value: c.merchant_category_id,
        })),
        placeholder: '请选择',
      },
    },
    {
      fieldName: 'type_id',
      componentProps: {
        clearable: true,
        options: types.value.map((t) => ({ label: t.name, value: t.id })),
        placeholder: '请选择',
      },
    },
  ]);
}

async function loadFilterOptions() {
  try {
    const [cats, typeList] = await Promise.all([
      fetchMerchantCategories().catch(() => ({ list: [] as MerchantCategoryRow[] })),
      fetchMerchantTypes().catch(() => ({ list: [] as MerchantTypeRow[] })),
    ]);
    categories.value = cats.list || [];
    types.value = typeList.list || [];
    syncFilterSelectOptions();
  } catch {
    /* 筛选项失败不阻断列表 */
  }
}

function openAudit(row: MerchantIntentionRow) {
  selected.value = row;
  form.status = 1;
  form.mark = row.mark || '';
  form.fail_msg = '';
  auditDrawerApi.setState({ title: '审核店铺入驻' }).open();
}

async function submitAudit() {
  if (!selected.value) return;
  if (form.status === 2 && !form.fail_msg.trim()) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  auditDrawerApi.lock();
  submitting.value = true;
  try {
    // 账号/初始密码/区域由申请资料自动带出（对齐 CRMEB：手机号作账号，区域用申请已分配 circle_id）
    await auditMerchantIntention(selected.value.mer_intention_id, {
      fail_msg: form.fail_msg.trim(),
      mark: form.mark.trim(),
      status: form.status,
    });
    ElMessage.success(form.status === 1 ? '入驻审核已通过' : '入驻申请已驳回');
    auditDrawerApi.close();
    gridApi.reload();
  } finally {
    submitting.value = false;
    auditDrawerApi.unlock();
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canAudit.value = permissions.includes('merchant.intention.audit');
  await loadFilterOptions();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #images="{ row }">
        <template v-if="imageUrls(row.images).length">
          <ElImage
            v-for="url in imageUrls(row.images)"
            :key="url"
            :preview-src-list="imageUrls(row.images)"
            :src="url"
            class="mr-1 h-8 w-8"
            fit="cover"
            preview-teleported
          />
        </template>
        <span v-else>—</span>
      </template>
      <template #status="{ row }">
        <ElTag :type="statusType(row.status)">{{ statusText(row.status) }}</ElTag>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="canAudit && row.status === 0"
          link
          type="primary"
          @click="openAudit(row)"
        >
          审核
        </ElButton>
        <span v-else>—</span>
      </template>
    </Grid>

    <AuditDrawer>
      <ElForm label-width="105px">
        <ElFormItem label="申请店铺">
          <span>{{ selected?.mer_name }}</span>
        </ElFormItem>
        <ElFormItem label="审核状态" required>
          <ElRadioGroup v-model="form.status">
            <ElRadio :value="1">同意</ElRadio>
            <ElRadio :value="2">拒绝</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="form.status === 2" label="驳回原因" required>
          <ElInput
            v-model="form.fail_msg"
            :rows="3"
            maxlength="300"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem label="审核备注">
          <ElInput
            v-model="form.mark"
            :rows="3"
            maxlength="300"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </AuditDrawer>
  </Page>
</template>
