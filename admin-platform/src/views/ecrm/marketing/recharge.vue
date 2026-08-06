<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createRechargePlan,
  listRechargeOrders,
  listRechargePlans,
  updateRechargePlan,
  type RechargeOrder,
  type RechargePlan,
} from '#/api/core/platform-recharge';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canManage = ref(false);
const saving = ref(false);
const editing = ref<number>();
const form = reactive({
  name: '',
  amount: 0,
  bonus_amount: 0,
  status: 1,
  sort: 0,
  version: 0,
});

const planFormOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('计划名称'),
  LIST_ENABLE_STATUS_FIELD('启用状态'),
]);

const planGridOptions: VxeGridProps<RechargePlan> = {
  columns: [
    { field: 'name', minWidth: 160, showOverflow: false, title: '计划' },
    {
      field: 'amount',
      title: '充值金额',
      width: 120,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'bonus_amount',
      title: '赠送金额',
      width: 120,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    { field: 'sort', title: '排序', width: 80 },
    platformListActionColumn({ width: 80 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const keyword = String(formValues?.keyword ?? '').trim() || undefined;
        const statusRaw = formValues?.status;
        const data = await listRechargePlans({
          keyword,
          status: statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
        });
        const list = data.list || [];
        return { items: list, total: list.length };
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

const orderFormOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('充值单号'),
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待支付', value: 'pending' },
        { label: '已支付', value: 'paid' },
        { label: '已关闭', value: 'closed' },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '订单状态',
  },
]);

function buildOrderParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const statusRaw = formValues?.status;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    status: typeof statusRaw === 'string' && statusRaw ? statusRaw : undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

const orderGridOptions: VxeGridProps<RechargeOrder> = {
  columns: [
    { field: 'recharge_no', minWidth: 180, showOverflow: false, title: '充值单号' },
    { field: 'user_id', title: '用户 ID', width: 100 },
    { field: 'amount', title: '金额', width: 100 },
    { field: 'bonus_amount', title: '赠送', width: 100 },
    { field: 'status', title: '状态', width: 100 },
    {
      field: 'created_at',
      minWidth: 170,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const data = await listRechargeOrders(buildOrderParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
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

const [PlanGrid, planGridApi] = useVbenVxeGrid({
  formOptions: planFormOptions,
  gridOptions: planGridOptions,
});

const [OrderGrid] = useVbenVxeGrid({
  formOptions: orderFormOptions,
  gridOptions: orderGridOptions,
});

const [PlanModal, planModalApi] = useVbenModal({
  onConfirm: async () => save(),
});

function edit(row?: RechargePlan) {
  editing.value = row?.id;
  form.name = String(row?.name ?? '');
  form.amount = Number(row?.amount ?? 0);
  form.bonus_amount = Number(row?.bonus_amount ?? 0);
  form.status = Number(row?.status ?? 1);
  form.sort = Number(row?.sort ?? 0);
  form.version = Number(row?.version ?? 0);
  planModalApi.setState({ title: editing.value ? '编辑充值计划' : '新增充值计划' }).open();
}

async function save() {
  if (!form.name.trim() || form.amount <= 0 || form.bonus_amount < 0) {
    ElMessage.warning('请填写计划名称和有效金额');
    return;
  }
  planModalApi.lock();
  saving.value = true;
  try {
    if (editing.value) await updateRechargePlan(editing.value, { ...form });
    else await createRechargePlan({ ...form });
    planModalApi.close();
    ElMessage.success('充值计划已保存');
    planGridApi.reload();
  } finally {
    saving.value = false;
    planModalApi.unlock();
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canManage.value =
    profile.roles.some((r) => r === 'platform' || r === 'operations') &&
    codes.includes('marketing.recharge.manage');
  if (canManage.value) {
    planGridApi.reload();
  }
});
</script>

<template>
  <Page
    auto-content-height
    description="维护可售充值计划，监管充值订单；后台不能改写已支付订单、支付回调或用户余额。"
    title="用户充值监管"
  >
    <ElAlert
      v-if="!canManage"
      title="当前账号没有充值监管权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <PlanGrid>
        <template #toolbar-actions>
          <ElButton :icon="Plus" type="primary" @click="edit()">新增计划</ElButton>
        </template>
        <template #status="{ row }">
          <ElTag :type="row.status ? 'success' : 'info'">
            {{ row.status ? '启用' : '停用' }}
          </ElTag>
        </template>
        <template #action="{ row }">
          <ElButton link type="primary" @click="edit(row)">编辑</ElButton>
        </template>
      </PlanGrid>

      <div class="mb-2 mt-6 text-base font-medium">充值订单（仅监管）</div>
      <OrderGrid />

      <PlanModal class="w-[480px]">
        <ElForm label-width="100px">
          <ElFormItem label="计划名称">
            <ElInput v-model="form.name" />
          </ElFormItem>
          <ElFormItem label="充值金额">
            <ElInputNumber v-model="form.amount" :min="0.01" :precision="2" class="w-full" />
          </ElFormItem>
          <ElFormItem label="赠送金额">
            <ElInputNumber v-model="form.bonus_amount" :min="0" :precision="2" class="w-full" />
          </ElFormItem>
          <ElFormItem label="状态">
            <ElRadioGroup v-model="form.status">
              <ElRadio :value="1">启用</ElRadio>
              <ElRadio :value="0">停用</ElRadio>
            </ElRadioGroup>
          </ElFormItem>
          <ElFormItem label="排序">
            <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
          </ElFormItem>
        </ElForm>
      </PlanModal>
    </template>
  </Page>
</template>
