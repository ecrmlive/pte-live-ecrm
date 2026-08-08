<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  fetchBusinessZoneAgents,
  resetBusinessZoneAgentPassword,
  revokeBusinessZoneAgent,
  type BusinessZoneAgentRow,
} from '#/api/core/ecrm';
import MerchantAdminFormDrawer from '#/components/ecrm/merchant-admin-form-drawer.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const formDrawerRef = ref<{
  openCreate: () => void;
  openEdit: (row: BusinessZoneAgentRow) => void;
}>();

const resetTarget = ref<BusinessZoneAgentRow>();
const passwordReset = reactive({
  password: '',
  confirmPassword: '',
  reason: '',
});

/** 负责商户：circles[].name，多商户顿号连接。 */
function formatResponsibleMerchants(row: BusinessZoneAgentRow) {
  const names = (row.circles || [])
    .map((item) => String(item?.name ?? '').trim())
    .filter(Boolean);
  return names.length > 0 ? names.join('、') : '—';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入管理员姓名',
    },
    fieldName: 'name',
    label: '管理员搜索',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入联系电话',
    },
    fieldName: 'phone',
    label: '联系电话',
  },
]);

const gridOptions: VxeGridProps<BusinessZoneAgentRow> = {
  columns: [
    { field: 'circle_agent_id', title: 'ID', width: 88 },
    { field: 'name', minWidth: 120, title: '管理员姓名' },
    { field: 'phone', title: '联系电话', width: 140 },
    {
      field: 'uid',
      formatter: ({ row }) =>
        row.uid
          ? `${row.nickname || '用户'} | ${row.uid}`
          : '—',
      minWidth: 160,
      title: '用户信息',
    },
    {
      field: 'circles',
      minWidth: 180,
      slots: { default: 'circles' },
      title: '负责商户',
    },
    platformListActionColumn({ width: 220 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const result = await fetchBusinessZoneAgents({
          page: page.currentPage,
          limit: page.pageSize,
          type: 1,
          status: 1,
          name: String(formValues?.name ?? '').trim() || undefined,
          phone: String(formValues?.phone ?? '').trim() || undefined,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'circle_agent_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [ResetDrawer, resetDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '确认重置',
  cancelText: '取消',
  placement: 'right',
  title: '重置商户管理员密码',
  onConfirm: async () => submitPasswordReset(),
});

function openCreate() {
  formDrawerRef.value?.openCreate();
}

function openEdit(row: BusinessZoneAgentRow) {
  formDrawerRef.value?.openEdit(row);
}

function onSaved() {
  gridApi.reload();
}

function openPasswordReset(row: BusinessZoneAgentRow) {
  resetTarget.value = row;
  Object.assign(passwordReset, {
    password: '',
    confirmPassword: '',
    reason: '',
  });
  resetDrawerApi.open();
}

async function submitPasswordReset() {
  const reason = passwordReset.reason.trim();
  const target = resetTarget.value;
  if (
    !target ||
    passwordReset.password.length < 12 ||
    passwordReset.password.length > 72 ||
    passwordReset.password !== passwordReset.confirmPassword ||
    reason.length < 2
  ) {
    ElMessage.warning('请填写两次一致的 12 至 72 位新密码和重置原因');
    return;
  }
  resetDrawerApi.lock();
  try {
    await resetBusinessZoneAgentPassword(target.circle_agent_id, {
      password: passwordReset.password,
      reason,
      idempotency_key: `merchant-admin-password-${target.circle_agent_id}-${crypto.randomUUID()}`,
    });
    ElMessage.success('密码已重置，旧后台会话已失效');
    resetDrawerApi.close();
  } finally {
    resetDrawerApi.unlock();
  }
}

async function remove(row: BusinessZoneAgentRow) {
  try {
    await ElMessageBox.confirm(
      `确定删除商户管理员「${row.name}」吗？删除后将解绑负责商户并停用后台登录。`,
      '提示',
      {
        type: 'warning',
        confirmButtonText: '删除',
        cancelButtonText: '取消',
      },
    );
    await revokeBusinessZoneAgent(row.circle_agent_id, {
      reason: `删除商户管理员 ${row.name}`,
      idempotency_key: `merchant-admin-delete-${row.circle_agent_id}-${Date.now()}`,
    });
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增管理
        </ElButton>
      </template>
      <template #circles="{ row }">
        {{ formatResponsibleMerchants(row) }}
      </template>
      <template #action="{ row }">
        <ElButton link type="warning" @click="openPasswordReset(row)">
          重置密码
        </ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <MerchantAdminFormDrawer
      ref="formDrawerRef"
      :show-responsible-merchants="true"
      @saved="onSaved"
    />

    <ResetDrawer>
      <ElAlert
        class="mb-4"
        type="warning"
        :closable="false"
        title="提交后旧后台会话立即失效；密码不会回显。"
      />
      <ElForm label-width="108px" @submit.prevent="submitPasswordReset">
        <ElFormItem label="管理员">
          <ElInput :model-value="resetTarget?.name || ''" disabled />
        </ElFormItem>
        <ElFormItem label="新密码" required>
          <ElInput
            v-model="passwordReset.password"
            autocomplete="new-password"
            show-password
            type="password"
          />
        </ElFormItem>
        <ElFormItem label="确认新密码" required>
          <ElInput
            v-model="passwordReset.confirmPassword"
            autocomplete="new-password"
            show-password
            type="password"
          />
        </ElFormItem>
        <ElFormItem label="重置原因" required>
          <ElInput
            v-model="passwordReset.reason"
            maxlength="500"
            show-word-limit
            type="textarea"
            :rows="3"
          />
        </ElFormItem>
      </ElForm>
    </ResetDrawer>
  </Page>
</template>
