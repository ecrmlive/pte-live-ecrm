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
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  auditMerchantIntention,
  createMerchantIntention,
  fetchMerchantCategories,
  fetchMerchantIntention,
  fetchMerchantIntentions,
  fetchMerchantTypes,
  type MerchantCategoryRow,
  type MerchantIntentionRow,
  type MerchantTypeRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import ImageField from '#/components/shop/image-field.vue';
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
const canCreate = ref(false);
const categories = ref<MerchantCategoryRow[]>([]);
const types = ref<MerchantTypeRow[]>([]);

const auditForm = reactive({
  fail_msg: '',
  mark: '',
  status: 1,
});

const createForm = reactive({
  mer_name: '',
  name: '',
  phone: '',
  merchant_category_id: undefined as number | undefined,
  mer_type_id: undefined as number | undefined,
  images: '',
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
    platformListActionColumn({ width: 160 }),
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

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
});

const [CreateDrawer, createDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: async () => submitCreate(),
});
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

function resetCreateForm() {
  createForm.mer_name = '';
  createForm.name = '';
  createForm.phone = '';
  createForm.merchant_category_id = undefined;
  createForm.mer_type_id = undefined;
  createForm.images = '';
}

function openCreate() {
  resetCreateForm();
  createDrawerApi.setState({ title: '新增店铺入驻申请' }).open();
}

async function openDetail(row: MerchantIntentionRow) {
  try {
    selected.value = await fetchMerchantIntention(row.mer_intention_id);
  } catch {
    selected.value = row;
  }
  detailDrawerApi.setState({ title: '店铺入驻申请详情' }).open();
}

function openAudit(row: MerchantIntentionRow) {
  selected.value = row;
  auditForm.status = 1;
  auditForm.mark = row.mark || '';
  auditForm.fail_msg = '';
  auditDrawerApi.setState({ title: '审核店铺入驻' }).open();
}

async function submitCreate() {
  const merName = createForm.mer_name.trim();
  const name = createForm.name.trim();
  const phone = createForm.phone.trim();
  if (merName.length < 2) {
    ElMessage.warning('请填写店铺名称');
    return;
  }
  if (!name) {
    ElMessage.warning('请填写联系人姓名');
    return;
  }
  if (phone.length < 6) {
    ElMessage.warning('请填写有效联系方式');
    return;
  }
  if (!createForm.merchant_category_id) {
    ElMessage.warning('请选择店铺分类');
    return;
  }
  if (!createForm.mer_type_id) {
    ElMessage.warning('请选择店铺类型');
    return;
  }
  const categoryName =
    categories.value.find((c) => c.merchant_category_id === createForm.merchant_category_id)
      ?.category_name || '';
  const typeName =
    types.value.find((t) => t.id === createForm.mer_type_id)?.name || '';
  createDrawerApi.lock();
  submitting.value = true;
  try {
    await createMerchantIntention({
      mer_name: merName,
      name,
      phone,
      merchant_category_id: createForm.merchant_category_id,
      mer_type_id: createForm.mer_type_id,
      images: createForm.images.trim(),
      category_name: categoryName,
      type_name: typeName,
    });
    ElMessage.success('入驻申请已提交');
    createDrawerApi.close();
    gridApi.reload();
  } finally {
    submitting.value = false;
    createDrawerApi.unlock();
  }
}

async function submitAudit() {
  if (!selected.value) return;
  if (auditForm.status === 2 && !auditForm.fail_msg.trim()) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  auditDrawerApi.lock();
  submitting.value = true;
  try {
    await auditMerchantIntention(selected.value.mer_intention_id, {
      fail_msg: auditForm.fail_msg.trim(),
      mark: auditForm.mark.trim(),
      status: auditForm.status,
    });
    ElMessage.success(auditForm.status === 1 ? '入驻审核已通过' : '入驻申请已驳回');
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
  canCreate.value = permissions.includes('merchant.intention.create');
  await loadFilterOptions();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canCreate"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          新增
        </ElButton>
      </template>

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
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          v-if="canAudit && row.status === 0"
          link
          type="primary"
          @click="openAudit(row)"
        >
          审核
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <ElForm label-width="105px">
        <ElFormItem label="申请 ID">
          <span>{{ displayOrDash(selected?.mer_intention_id) }}</span>
        </ElFormItem>
        <ElFormItem label="店铺名称">
          <span>{{ displayOrDash(selected?.mer_name) }}</span>
        </ElFormItem>
        <ElFormItem label="店铺分类">
          <span>{{ displayOrDash(selected?.category_name) }}</span>
        </ElFormItem>
        <ElFormItem label="店铺类型">
          <span>{{ displayOrDash(selected?.type_name) }}</span>
        </ElFormItem>
        <ElFormItem label="联系人">
          <span>{{ displayOrDash(selected?.name) }}</span>
        </ElFormItem>
        <ElFormItem label="联系方式">
          <span>{{ displayOrDash(selected?.phone) }}</span>
        </ElFormItem>
        <ElFormItem label="审核状态">
          <ElTag v-if="selected" :type="statusType(selected.status)">
            {{ statusText(selected.status) }}
          </ElTag>
          <span v-else>—</span>
        </ElFormItem>
        <ElFormItem label="申请时间">
          <span>{{ formatTime(selected?.create_time) }}</span>
        </ElFormItem>
        <ElFormItem label="资质图片">
          <template v-if="imageUrls(selected?.images).length">
            <ElImage
              v-for="url in imageUrls(selected?.images)"
              :key="url"
              :preview-src-list="imageUrls(selected?.images)"
              :src="url"
              class="mr-2 h-16 w-16"
              fit="cover"
              preview-teleported
            />
          </template>
          <span v-else>—</span>
        </ElFormItem>
        <ElFormItem label="审核备注">
          <span>{{ displayOrDash(selected?.mark) }}</span>
        </ElFormItem>
        <ElFormItem v-if="selected?.status === 2" label="驳回原因">
          <span>{{ displayOrDash(selected?.fail_msg) }}</span>
        </ElFormItem>
        <ElFormItem v-if="canAudit && selected?.status === 0" label="操作">
          <ElButton type="primary" @click="openAudit(selected!)">去审核</ElButton>
        </ElFormItem>
      </ElForm>
    </DetailDrawer>

    <CreateDrawer>
      <ElForm label-width="105px">
        <ElFormItem label="店铺名称" required>
          <ElInput v-model="createForm.mer_name" maxlength="64" placeholder="请输入店铺名称" />
        </ElFormItem>
        <ElFormItem label="联系人" required>
          <ElInput v-model="createForm.name" maxlength="32" placeholder="请输入联系人姓名" />
        </ElFormItem>
        <ElFormItem label="联系方式" required>
          <ElInput v-model="createForm.phone" maxlength="32" placeholder="请输入手机号" />
        </ElFormItem>
        <ElFormItem label="店铺分类" required>
          <ElSelect
            v-model="createForm.merchant_category_id"
            class="w-full"
            clearable
            placeholder="请选择"
          >
            <ElOption
              v-for="c in categories"
              :key="c.merchant_category_id"
              :label="c.category_name"
              :value="c.merchant_category_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="店铺类型" required>
          <ElSelect
            v-model="createForm.mer_type_id"
            class="w-full"
            clearable
            placeholder="请选择"
          >
            <ElOption
              v-for="t in types"
              :key="t.id"
              :label="t.name"
              :value="t.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="资质图片">
          <ImageField v-model="createForm.images" />
        </ElFormItem>
      </ElForm>
    </CreateDrawer>

    <AuditDrawer>
      <ElForm label-width="105px">
        <ElFormItem label="申请店铺">
          <span>{{ selected?.mer_name }}</span>
        </ElFormItem>
        <ElFormItem label="审核状态" required>
          <ElRadioGroup v-model="auditForm.status">
            <ElRadio :value="1">同意</ElRadio>
            <ElRadio :value="2">拒绝</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem v-if="auditForm.status === 2" label="驳回原因" required>
          <ElInput
            v-model="auditForm.fail_msg"
            :rows="3"
            maxlength="300"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem label="审核备注">
          <ElInput
            v-model="auditForm.mark"
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
