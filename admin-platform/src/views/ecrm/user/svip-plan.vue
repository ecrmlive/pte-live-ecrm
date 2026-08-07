<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElCheckbox,
  ElCheckboxGroup,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createSvipPlan,
  listSvipPlans,
  updateSvipPlan,
  type SvipPlan,
  type SvipPlanInput,
} from '#/api/core/platform-svip-plan';
import {
  listSvipInterests,
  type SvipInterest,
} from '#/api/core/platform-svip-interest';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canManage = ref(false);
const saving = ref(false);
const editingId = ref<number>();
const interests = ref<SvipInterest[]>([]);
const form = reactive<SvipPlanInput>({
  name: '',
  price: 0,
  plan_type: 'period',
  duration_days: 30,
  benefits: [],
  status: 1,
  sort: 0,
});

function planTypeName(type: SvipPlan['plan_type']) {
  return (
    { trial: '体验会员', period: '期限会员', lifetime: '永久会员' }[type] ||
    type
  );
}

function parseBenefits(raw: string) {
  try {
    const parsed = JSON.parse(raw);
    return Array.isArray(parsed)
      ? parsed.filter((item): item is string => typeof item === 'string')
      : [];
  } catch {
    return [];
  }
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '会员类型名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '启用', value: 1 },
        { label: '停用', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<SvipPlan> = {
  columns: [
    { field: 'name', minWidth: 140, showOverflow: false, title: '名称' },
    {
      field: 'plan_type',
      formatter: ({ cellValue }) =>
        planTypeName(cellValue as SvipPlan['plan_type']),
      title: '类型',
      width: 110,
    },
    {
      field: 'price',
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
      title: '售价',
      width: 110,
    },
    {
      field: 'duration_days',
      formatter: ({ row }) =>
        row.plan_type === 'lifetime' ? '永久' : `${row.duration_days} 天`,
      title: '有效期',
      width: 110,
    },
    {
      field: 'benefits',
      formatter: ({ cellValue }) =>
        parseBenefits(String(cellValue ?? '')).join('、'),
      minWidth: 220,
      showOverflow: false,
      title: '权益',
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
        const keyword = String(formValues?.keyword ?? '').trim().toLowerCase();
        const statusRaw = formValues?.status;
        let list = (await listSvipPlans()).list || [];
        if (keyword) {
          list = list.filter((row) => row.name.toLowerCase().includes(keyword));
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter((row) => row.status === Number(statusRaw));
        }
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editingId.value = undefined;
  Object.assign(form, {
    name: '',
    price: 0,
    plan_type: 'period',
    duration_days: 30,
    benefits: [],
    status: 1,
    sort: 0,
  });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增会员类型' }).open();
}

function openEdit(row: SvipPlan) {
  editingId.value = row.id;
  Object.assign(form, {
    name: row.name,
    price: Number(row.price),
    plan_type: row.plan_type,
    duration_days: row.duration_days || 0,
    benefits: parseBenefits(row.benefits),
    status: row.status,
    sort: row.sort,
  });
  formDrawerApi.setState({ title: '编辑会员类型' }).open();
}

async function save() {
  if (
    !form.name.trim() ||
    !form.benefits.length ||
    (form.plan_type === 'trial' &&
      (form.price !== 0 || form.duration_days < 1)) ||
    (form.plan_type === 'period' &&
      (form.price <= 0 || form.duration_days < 1)) ||
    (form.plan_type === 'lifetime' && form.price <= 0)
  ) {
    ElMessage.warning('请完整填写符合会员类型规则的套餐信息');
    return;
  }
  if (form.plan_type === 'lifetime') form.duration_days = 0;
  formDrawerApi.lock();
  saving.value = true;
  try {
    if (editingId.value) await updateSvipPlan(editingId.value, form);
    else await createSvipPlan(form);
    formDrawerApi.close();
    ElMessage.success('会员类型已保存');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canManage.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    codes.includes('user.svip.plan.manage');
  if (canManage.value) {
    interests.value = (
      (await listSvipInterests()).list || []
    ).filter((item) => item.status === 1);
    gridApi.reload();
  }
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canManage"
      class="mb-4"
      title="当前账号没有会员类型维护权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增类型
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status ? 'success' : 'info'">
          {{ row.status ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="98px">
        <ElFormItem label="名称">
          <ElInput v-model="form.name" maxlength="64" />
        </ElFormItem>
        <ElFormItem label="会员类型">
          <ElSelect v-model="form.plan_type" class="w-full">
            <ElOption label="体验会员（免费）" value="trial" />
            <ElOption label="期限会员" value="period" />
            <ElOption label="永久会员" value="lifetime" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="售价">
          <ElInputNumber
            v-model="form.price"
            :disabled="form.plan_type === 'trial'"
            :min="0"
            :precision="2"
          />
        </ElFormItem>
        <ElFormItem v-if="form.plan_type !== 'lifetime'" label="有效天数">
          <ElInputNumber
            v-model="form.duration_days"
            :min="1"
            :max="form.plan_type === 'trial' ? 31 : 3660"
          />
        </ElFormItem>
        <ElFormItem label="会员权益" required>
          <ElCheckboxGroup v-model="form.benefits">
            <ElCheckbox
              v-for="item in interests"
              :key="item.id"
              :label="item.name"
            >
              {{ item.name }}
            </ElCheckbox>
          </ElCheckboxGroup>
          <div v-if="!interests.length" class="text-sm text-red-500">
            暂无启用权益，请先在“会员权益”页面维护。
          </div>
        </ElFormItem>
        <ElFormItem label="状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio :value="1">启用</ElRadio>
            <ElRadio :value="0">停用</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
