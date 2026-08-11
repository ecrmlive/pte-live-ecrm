<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createSvipPlan,
  deleteSvipPlan,
  listSvipPlans,
  updateSvipPlan,
  updateSvipPlanStatus,
  type SvipPlan,
  type SvipPlanInput,
} from '#/api/core/platform-svip-plan';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';

const canRead = ref(false);
const canManage = ref(false);
const saving = ref(false);
const editingId = ref<number>();
const form = reactive<SvipPlanInput>({
  name: '',
  cost_price: 0,
  price: 0,
  plan_type: 'period',
  duration_days: 30,
  benefits: [],
  status: 1,
  sort: 0,
});

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

function money(value: unknown) {
  return Number(value || 0).toFixed(2);
}

const gridOptions: VxeGridProps<SvipPlan> = {
  columns: [
    { field: 'id', title: 'ID', width: 88 },
    {
      field: 'name',
      minWidth: 140,
      showOverflow: false,
      title: '会员名',
    },
    {
      field: 'duration_days',
      formatter: ({ row }) =>
        row.plan_type === 'lifetime' ? '永久' : String(row.duration_days ?? 0),
      minWidth: 110,
      title: '有效期(天)',
    },
    {
      field: 'cost_price',
      formatter: ({ cellValue }) => money(cellValue),
      minWidth: 110,
      title: '原价',
    },
    {
      field: 'price',
      formatter: ({ cellValue }) => money(cellValue),
      minWidth: 110,
      title: '优惠价',
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否开启',
      width: 120,
    },
    platformListActionColumn({ width: 130 }),
  ],
  emptyText: '暂无数据',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const result = await listSvipPlans({
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: result.list || [], total: result.total || 0 };
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editingId.value = undefined;
  Object.assign(form, {
    name: '',
    cost_price: 0,
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
    cost_price: Number(row.cost_price || 0),
    price: Number(row.price || 0),
    plan_type: row.plan_type,
    duration_days:
      row.plan_type === 'lifetime' ? 0 : Number(row.duration_days || 0),
    benefits: parseBenefits(row.benefits),
    status: row.status === 1 ? 1 : 0,
    sort: Number(row.sort || 0),
  });
  formDrawerApi.setState({ title: '编辑会员类型' }).open();
}

async function save() {
  if (!canManage.value) {
    ElMessage.warning('当前账号没有管理权限');
    return;
  }
  const name = form.name.trim();
  if (!name) {
    ElMessage.warning('请填写会员名');
    return;
  }
  if (form.cost_price < 0 || form.price < 0) {
    ElMessage.warning('请填写有效的原价与优惠价');
    return;
  }
  if (
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
  if (form.plan_type === 'trial') form.price = 0;
  formDrawerApi.lock();
  saving.value = true;
  try {
    const payload: SvipPlanInput = {
      name,
      cost_price: Number(form.cost_price),
      price: Number(form.price),
      plan_type: form.plan_type,
      duration_days: form.duration_days,
      benefits: form.benefits,
      status: form.status === 1 ? 1 : 0,
      sort: Math.max(0, Math.min(999999, Number(form.sort) || 0)),
    };
    if (editingId.value) await updateSvipPlan(editingId.value, payload);
    else await createSvipPlan(payload);
    formDrawerApi.close();
    ElMessage.success('会员类型已保存');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: SvipPlan, enabled: boolean) {
  if (!canManage.value) {
    ElMessage.warning('当前账号没有管理权限');
    return;
  }
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await updateSvipPlanStatus(row.id, enabled ? 1 : 0);
  } catch {
    row.status = before;
  }
}

async function remove(row: SvipPlan) {
  if (!canManage.value) {
    ElMessage.warning('当前账号没有管理权限');
    return;
  }
  try {
    await confirm({
      content: `确定删除会员类型“${row.name}”吗？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteSvipPlan(row.id);
    ElMessage.success('会员类型已删除');
    gridApi.reload();
  } catch {
    /* 取消或请求层已提示 */
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const roleOK = profile.roles.some(
    (role) => role === 'platform' || role === 'operations',
  );
  canRead.value =
    roleOK &&
    (codes.includes('user.svip.plan.read') ||
      codes.includes('user.svip.plan.manage'));
  canManage.value = roleOK && codes.includes('user.svip.plan.manage');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有会员类型查看权限，请确认已绑定权限并重新登录"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          新增会员类型
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage"
          inline-prompt
          active-text="开启"
          inactive-text="关闭"
          @change="
            (enabled: string | number | boolean) =>
              changeStatus(row, Boolean(enabled))
          "
        />
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="canManage"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-if="canManage"
          link
          type="danger"
          @click="remove(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="110px">
        <ElFormItem label="会员名" required>
          <ElInput
            v-model="form.name"
            maxlength="64"
            placeholder="请输入会员名"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="会员类别" required>
          <ElSelect v-model="form.plan_type" class="w-full">
            <ElOption label="试用期" value="trial" />
            <ElOption label="有限期" value="period" />
            <ElOption label="永久期" value="lifetime" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem
          v-if="form.plan_type !== 'lifetime'"
          label="有效期(天)"
          required
        >
          <ElInputNumber
            v-model="form.duration_days"
            :min="1"
            :max="form.plan_type === 'trial' ? 31 : 3660"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="原价" required>
          <ElInputNumber
            v-model="form.cost_price"
            :min="0"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="优惠价" required>
          <ElInputNumber
            v-model="form.price"
            :disabled="form.plan_type === 'trial'"
            :min="0"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber
            v-model="form.sort"
            :min="0"
            :max="999999"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="是否开启">
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="开启"
            inactive-text="关闭"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
