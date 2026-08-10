<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createRechargePlan,
  deleteRechargePlan,
  listRechargePlans,
  updateRechargePlan,
  updateRechargePlanStatus,
  type RechargePlan,
} from '#/api/core/platform-recharge';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

const canRead = ref(false);
const canManage = ref(false);
const editing = ref<RechargePlan>();
const form = reactive({
  amount: 1,
  bonus_amount: 0,
  sort: 0,
  status: 1,
});

const gridOptions: VxeGridProps<RechargePlan> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    {
      field: 'amount',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      minWidth: 120,
      title: '充值金额',
    },
    {
      field: 'bonus_amount',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      minWidth: 120,
      title: '赠送金额',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 180,
      showOverflow: false,
      title: '添加时间',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 120,
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listRechargePlans({
          page: page.currentPage,
          limit: page.pageSize,
        });
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '确定保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    amount: 1,
    bonus_amount: 0,
    sort: 0,
    status: 1,
  });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加充值档位' }).open();
}

function openEdit(row: RechargePlan) {
  editing.value = row;
  Object.assign(form, {
    amount: Number(row.amount || 0),
    bonus_amount: Number(row.bonus_amount || 0),
    sort: Number(row.sort || 0),
    status: row.status === 1 ? 1 : 0,
  });
  formDrawerApi.setState({ title: '编辑充值档位' }).open();
}

async function save() {
  if (!canManage.value) {
    ElMessage.warning('当前账号没有管理权限');
    return;
  }
  const amount = Number(form.amount);
  const bonus = Number(form.bonus_amount);
  const sort = Math.max(0, Math.min(999999, Number(form.sort) || 0));
  if (!(amount > 0) || Number.isNaN(amount) || bonus < 0 || Number.isNaN(bonus)) {
    ElMessage.warning('请填写有效的充值金额与赠送金额');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      amount,
      bonus_amount: bonus,
      sort,
      status: form.status === 1 ? 1 : 0,
    };
    if (editing.value) {
      await updateRechargePlan(editing.value.id, {
        ...payload,
        version: Number(editing.value.version || 1),
      });
    } else {
      await createRechargePlan(payload);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: RechargePlan, enabled: boolean) {
  if (!canManage.value) {
    ElMessage.warning('当前账号没有管理权限');
    return;
  }
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await updateRechargePlanStatus(row.id, enabled ? 1 : 0);
    row.version = Number(row.version || 0) + 1;
  } catch {
    row.status = before;
  }
}

async function remove(row: RechargePlan) {
  if (!canManage.value) {
    ElMessage.warning('当前账号没有管理权限');
    return;
  }
  try {
    await confirm({
      content: '确定删除该充值档位吗？',
      icon: 'warning',
      title: '提示',
    });
    await deleteRechargePlan(row.id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const roleOK = profile.roles.some(
    (r) => r === 'platform' || r === 'operations',
  );
  canRead.value = roleOK && codes.includes('marketing.balance.config.read');
  canManage.value =
    roleOK && codes.includes('marketing.balance.config.manage');
  if (canRead.value) {
    gridApi.reload();
  }
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      title="当前账号没有余额充值设置查看权限，请确认已绑定权限并重新登录"
      type="warning"
      :closable="false"
      class="mb-3"
    />
    <Grid v-else>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          添加数据
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeShow(row, Boolean(enabled))
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
          type="primary"
          @click="remove(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="100px">
        <ElFormItem label="充值金额" required>
          <ElInputNumber
            v-model="form.amount"
            :min="0.01"
            :max="1000000"
            :precision="2"
            :step="1"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="赠送金额" required>
          <ElInputNumber
            v-model="form.bonus_amount"
            :min="0"
            :max="1000000"
            :precision="2"
            :step="1"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber
            v-model="form.sort"
            :min="0"
            :max="999999"
            :precision="0"
            :step="1"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="是否显示">
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="显示"
            inactive-text="隐藏"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
