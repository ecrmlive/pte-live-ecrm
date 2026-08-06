<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  assignMerchantIntentionRegion,
  auditMerchantIntention,
  fetchMerchantIntentions,
  type MerchantIntentionRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const selected = ref<MerchantIntentionRow>();
const submitting = ref(false);
const canAudit = ref(false);
const canAssignRegion = ref(false);

const form = reactive({
  account: '',
  fail_msg: '',
  mark: '',
  password: '',
  region_id: 0,
  status: 1,
});
const assignment = reactive({ region_id: 0 });

const statusText = (status: number) =>
  ({ 0: '待审核', 1: '已通过', 2: '已驳回' })[status] || '未知';
const statusType = (status: number) =>
  ({ 0: 'warning', 1: 'success', 2: 'danger' })[status] || 'info';
const formatTime = (value?: string | null) =>
  value ? formatShanghaiDateTime(value) : '—';

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const statusRaw = formValues?.status;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    status:
      statusRaw === 0 || statusRaw === 1 || statusRaw === 2
        ? Number(statusRaw)
        : undefined,
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
      placeholder: '店铺名称 / 联系人 / 手机号',
    },
    fieldName: 'keyword',
    label: '申请搜索',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 0 },
        { label: '已通过', value: 1 },
        { label: '已驳回', value: 2 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<MerchantIntentionRow> = {
  columns: [
    { field: 'mer_intention_id', title: '申请 ID', width: 88 },
    {
      field: 'mer_name',
      minWidth: 150,
      showOverflow: false,
      title: '店铺名称',
    },
    { field: 'name', minWidth: 110, title: '联系人' },
    { field: 'phone', minWidth: 130, title: '联系电话' },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 170,
      title: '申请时间',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 96,
    },
    {
      field: 'circle_id',
      formatter: ({ cellValue }) => cellValue || '未分配',
      title: '分配区域',
      width: 108,
    },
    {
      field: 'fail_msg',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 160,
      showOverflow: false,
      title: '驳回原因',
    },
    platformListActionColumn({ width: 205 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
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

const [AuditModal, auditModalApi] = useVbenModal({
  onConfirm: async () => submitAudit(),
});

const [RegionModal, regionModalApi] = useVbenModal({
  onConfirm: async () => assignRegion(),
});

function openAudit(row: MerchantIntentionRow, status: number) {
  selected.value = row;
  form.status = status;
  form.mark = '';
  form.fail_msg = '';
  form.account = '';
  form.password = '';
  form.region_id = row.circle_id || 0;
  auditModalApi
    .setState({
      title: status === 1 ? '通过店铺入驻' : '驳回店铺入驻',
    })
    .open();
}

function openRegionAssignment(row: MerchantIntentionRow) {
  selected.value = row;
  assignment.region_id = row.circle_id || 0;
  regionModalApi.setState({ title: '分配入驻审核区域' }).open();
}

async function assignRegion() {
  if (!selected.value || assignment.region_id <= 0) {
    ElMessage.warning('请填写有效区域 ID');
    return;
  }
  regionModalApi.lock();
  submitting.value = true;
  try {
    await assignMerchantIntentionRegion(
      selected.value.mer_intention_id,
      assignment.region_id,
    );
    ElMessage.success('入驻申请已分配区域');
    regionModalApi.close();
    gridApi.reload();
  } finally {
    submitting.value = false;
    regionModalApi.unlock();
  }
}

async function submitAudit() {
  if (!selected.value) return;
  if (form.status === 2 && !form.fail_msg.trim()) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  if (form.status === 1 && (!form.account.trim() || !form.password)) {
    ElMessage.warning('通过入驻时必须设置商户管理账号和初始密码');
    return;
  }
  auditModalApi.lock();
  submitting.value = true;
  try {
    await auditMerchantIntention(selected.value.mer_intention_id, {
      account: form.account.trim(),
      fail_msg: form.fail_msg.trim(),
      mark: form.mark.trim(),
      password: form.password,
      region_id: form.region_id,
      status: form.status,
    });
    ElMessage.success(form.status === 1 ? '入驻审核已通过' : '入驻申请已驳回');
    auditModalApi.close();
    gridApi.reload();
  } finally {
    submitting.value = false;
    auditModalApi.unlock();
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canAudit.value = permissions.includes('merchant.intention.audit');
  canAssignRegion.value = permissions.includes('merchant.intention.assign_region');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #status="{ row }">
        <ElTag :type="statusType(row.status)">{{ statusText(row.status) }}</ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="row.status === 0">
          <ElButton
            v-if="canAssignRegion"
            link
            type="primary"
            @click="openRegionAssignment(row)"
          >
            分配区域
          </ElButton>
          <ElButton
            v-if="canAudit"
            link
            type="primary"
            @click="openAudit(row, 1)"
          >
            通过
          </ElButton>
          <ElButton
            v-if="canAudit"
            link
            type="danger"
            @click="openAudit(row, 2)"
          >
            驳回
          </ElButton>
          <span v-if="!canAssignRegion && !canAudit">—</span>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <AuditModal>
      <ElForm label-width="105px">
        <ElFormItem label="申请店铺">
          <span>{{ selected?.mer_name }}</span>
        </ElFormItem>
        <template v-if="form.status === 1">
          <ElFormItem label="商户管理账号" required>
            <ElInput v-model="form.account" autocomplete="off" />
          </ElFormItem>
          <ElFormItem label="初始密码" required>
            <ElInput
              v-model="form.password"
              autocomplete="new-password"
              show-password
              type="password"
            />
          </ElFormItem>
          <ElFormItem label="所属区域ID">
            <ElInputNumber v-model="form.region_id" :min="0" />
            <div class="text-xs text-gray-500">
              填写区域商圈对应的商户 region_id，区域管理员将按此范围管理商户。
            </div>
          </ElFormItem>
        </template>
        <ElFormItem v-else label="驳回原因" required>
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
    </AuditModal>

    <RegionModal>
      <ElForm label-width="100px">
        <ElFormItem label="申请店铺">
          <span>{{ selected?.mer_name }}</span>
        </ElFormItem>
        <ElFormItem label="区域 ID" required>
          <ElInputNumber v-model="assignment.region_id" :min="1" />
        </ElFormItem>
        <p class="text-xs text-gray-500">
          分配后仅对应区域管理员可查看和审核该申请。
        </p>
      </ElForm>
    </RegionModal>
  </Page>
</template>
