<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { confirm, Page, useVbenDrawer } from '@vben/common-ui';

import { Plus } from '@element-plus/icons-vue';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformAdmin,
  deletePlatformAdmin,
  fetchBusinessZoneAgentOptions,
  fetchPlatformAdminRegions,
  fetchPlatformAdmins,
  fetchPlatformMerchants,
  fetchPlatformRoles,
  type PlatformAdminRow,
  type PlatformRoleRow,
  updatePlatformAdmin,
} from '#/api/core/ecrm';
import { listPlatformProductStoresApi } from '#/api/core/platform-catalog';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

type SelectOption = { label: string; value: number };

const editing = ref<PlatformAdminRow>();
const roles = ref<PlatformRoleRow[]>([]);
const regions = ref<SelectOption[]>([]);
const merchants = ref<SelectOption[]>([]);
const stores = ref<SelectOption[]>([]);
const agentOptions = ref<
  Array<{
    circle_agent_id: number;
    name: string;
    phone?: string;
    type: 0 | 1;
  }>
>([]);

const form = reactive({
  account: '',
  circle_agent_id: 0,
  merchant_ids: [] as number[],
  password: '',
  phone: '',
  real_name: '',
  region_ids: [] as number[],
  role_codes: [] as string[],
  service_store_ids: [] as number[],
  status: 1,
});
const passwordForm = reactive({
  confirm_password: '',
  password: '',
});
const passwordTarget = ref<PlatformAdminRow>();

const isRegionAccount = computed(
  () => form.role_codes.includes('region') && !form.role_codes.includes('platform'),
);
const hasMerchantRole = computed(() => form.role_codes.includes('merchant'));
const hasServiceRole = computed(() => form.role_codes.includes('customer_service'));
const modalTitle = computed(() => (editing.value ? '编辑管理员' : '新增管理员'));

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_ENABLE_STATUS_FIELD('账号状态'),
  LIST_KEYWORD_FIELD('请输入账号/昵称'),
]);

const gridOptions: VxeGridProps<PlatformAdminRow> = {
  columns: [
    { field: 'admin_id', title: 'ID', width: 78 },
    { field: 'real_name', minWidth: 150, title: '管理员姓名' },
    { field: 'role_names', minWidth: 160, title: '身份' },
    {
      field: 'region_names',
      minWidth: 180,
      slots: { default: 'regions' },
      title: '所属区域',
    },
    { field: 'account', minWidth: 150, title: '账号' },
    { field: 'phone', minWidth: 140, title: '手机号' },
    {
      align: 'center',
      field: 'status',
      slots: { default: 'status' },
      title: '账号状态',
      width: 112,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 190,
      title: '创建时间',
    },
    platformListActionColumn({ width: 230 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = Number(formValues?.status);
        const result = await fetchPlatformAdmins({
          date_from: range[0] ? String(range[0]) : undefined,
          date_to: range[1] ? String(range[1]) : undefined,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          limit: page.pageSize,
          page: page.currentPage,
          status: status === 0 || status === 1 ? status : undefined,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'admin_id' },
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
  cancelText: '取消',
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: save,
});
const [PasswordDrawer, passwordDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[720px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  title: '修改密码',
  onConfirm: savePassword,
});

function csv(values: number[]) {
  return values.map(Number).filter(Boolean).join(',');
}

function parseCSV(value?: string) {
  return String(value || '')
    .split(',')
    .map(Number)
    .filter(Boolean);
}

function resetForm() {
  Object.assign(form, {
    account: '',
    circle_agent_id: 0,
    merchant_ids: [],
    password: '',
    phone: '',
    real_name: '',
    region_ids: [],
    role_codes: [],
    service_store_ids: [],
    status: 1,
  });
}

function openCreate() {
  editing.value = undefined;
  resetForm();
  formDrawerApi.setState({ title: '新增管理员' }).open();
}

function openEdit(row: PlatformAdminRow) {
  editing.value = row;
  Object.assign(form, {
    account: row.account,
    circle_agent_id: row.circle_agent_id || 0,
    merchant_ids: parseCSV(row.merchant_ids),
    password: '',
    phone: row.phone,
    real_name: row.real_name,
    region_ids: parseCSV(row.region_ids),
    role_codes: row.role_codes?.length
      ? [...row.role_codes]
      : String(row.roles || '').split(',').filter(Boolean),
    service_store_ids: parseCSV(row.service_store_ids),
    status: row.status,
  });
  formDrawerApi.setState({ title: '编辑管理员' }).open();
}

async function save() {
  if (!editing.value && (!form.account.trim() || form.password.length < 8)) {
    ElMessage.warning('请填写账号和至少 8 位初始密码');
    return;
  }
  if (!form.real_name.trim() || !form.role_codes.length) {
    ElMessage.warning('请填写管理员姓名并选择身份');
    return;
  }
  if (hasMerchantRole.value && !form.merchant_ids.length) {
    ElMessage.warning('请选择授权商户');
    return;
  }
  if (isRegionAccount.value && (!form.region_ids.length || !form.circle_agent_id)) {
    ElMessage.warning('请选择所属区域和关联代理');
    return;
  }
  if (hasServiceRole.value && !form.service_store_ids.length) {
    ElMessage.warning('请选择客服授权店铺');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      account: form.account.trim(),
      circle_agent_id: form.circle_agent_id,
      merchant_ids: csv(form.merchant_ids),
      password: form.password,
      phone: form.phone.trim(),
      real_name: form.real_name.trim(),
      region_ids: csv(form.region_ids),
      role_codes: [...form.role_codes],
      service_store_ids: csv(form.service_store_ids),
      status: form.status,
    };
    if (editing.value) await updatePlatformAdmin(editing.value.admin_id, payload);
    else await createPlatformAdmin(payload);
    formDrawerApi.close();
    ElMessage.success('管理员已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

function openPassword(row: PlatformAdminRow) {
  passwordTarget.value = row;
  Object.assign(passwordForm, { confirm_password: '', password: '' });
  passwordDrawerApi.open();
}

async function savePassword() {
  const target = passwordTarget.value;
  if (!target || passwordForm.password.length < 8) {
    ElMessage.warning('新密码至少 8 位');
    return;
  }
  if (passwordForm.password !== passwordForm.confirm_password) {
    ElMessage.warning('两次输入的密码不一致');
    return;
  }
  passwordDrawerApi.lock();
  try {
    await updatePlatformAdmin(target.admin_id, {
      account: target.account,
      avatar_url: target.avatar_url || '',
      linked_user_id: target.linked_user_id || 0,
      merchant_ids: target.merchant_ids,
      password: passwordForm.password,
      real_name: target.real_name,
      region_ids: target.region_ids,
      role_codes: target.role_codes || target.roles.split(',').filter(Boolean),
      service_store_ids: target.service_store_ids,
      status: target.status,
    });
    passwordDrawerApi.close();
    ElMessage.success('密码已保存，旧登录会话已失效');
  } finally {
    passwordDrawerApi.unlock();
  }
}

async function toggleStatus(row: PlatformAdminRow, enabled: boolean) {
  try {
    await updatePlatformAdmin(row.admin_id, {
      account: row.account,
      avatar_url: row.avatar_url || '',
      linked_user_id: row.linked_user_id || 0,
      merchant_ids: row.merchant_ids,
      real_name: row.real_name,
      region_ids: row.region_ids,
      role_codes: row.role_codes || row.roles.split(',').filter(Boolean),
      service_store_ids: row.service_store_ids,
      status: enabled ? 1 : 0,
    });
    row.status = enabled ? 1 : 0;
    ElMessage.success('账号状态已保存');
  } catch {
    gridApi.reload();
  }
}

async function remove(row: PlatformAdminRow) {
  try {
    await confirm({
      content: `删除管理员“${row.real_name || row.account}”后将无法登录，是否继续？`,
      icon: 'warning',
      title: '提示',
    });
    await deletePlatformAdmin(row.admin_id);
    ElMessage.success('管理员已删除');
    gridApi.reload();
  } catch {
    // 取消删除无需提示。
  }
}

async function loadOptions() {
  const [roleResult, regionResult, merchantResult, storeResult, agentsResult] =
    await Promise.all([
      fetchPlatformRoles({ limit: 100, page: 1 }),
      fetchPlatformAdminRegions(),
      fetchPlatformMerchants({ limit: 100, page: 1 }),
      listPlatformProductStoresApi(),
      fetchBusinessZoneAgentOptions(0),
    ]);
  roles.value = (roleResult.list || []).filter(
    (role) => role.code !== 'customer_service',
  );
  regions.value = regionResult.list || [];
  merchants.value = (merchantResult.list || []).map((item) => ({
    label: item.mer_name,
    value: item.mer_id,
  }));
  stores.value = (storeResult.list || []).map((item) => ({
    label: `${item.store_name}（${item.merchant_name}）`,
    value: item.store_id,
  }));
  agentOptions.value = agentsResult.list || [];
}

onMounted(async () => {
  try {
    await loadOptions();
  } catch {
    ElMessage.error('关联数据加载失败，请刷新后重试');
  }
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增管理员
        </ElButton>
      </template>
      <template #regions="{ row }">
        <div v-if="row.region_names" class="flex flex-wrap gap-1">
          <ElTag v-for="name in row.region_names.split('、')" :key="name" effect="plain">
            {{ name }}
          </ElTag>
        </div>
        <span v-else>—</span>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          @change="(value: boolean | string | number) => toggleStatus(row, value === true)"
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openPassword(row)">修改密码</ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer :title="modalTitle">
      <ElForm label-width="120px">
        <ElFormItem label="账号" required>
          <ElInput v-model="form.account" :disabled="!!editing" maxlength="64" />
        </ElFormItem>
        <ElFormItem v-if="!editing" label="初始密码" required>
          <ElInput
            v-model="form.password"
            autocomplete="new-password"
            show-password
            type="password"
            placeholder="至少 8 位"
          />
        </ElFormItem>
        <ElFormItem label="管理员姓名" required>
          <ElInput v-model="form.real_name" maxlength="64" />
        </ElFormItem>
        <ElFormItem label="手机号">
          <ElInput v-model="form.phone" maxlength="32" />
        </ElFormItem>
        <ElFormItem label="身份" required>
          <ElSelect v-model="form.role_codes" class="w-full" multiple placeholder="请选择身份">
            <ElOption
              v-for="role in roles"
              :key="role.code"
              :label="role.role_name"
              :value="role.code"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem v-if="hasMerchantRole" label="授权商户" required>
          <ElSelect v-model="form.merchant_ids" class="w-full" multiple placeholder="请选择商户">
            <ElOption v-for="item in merchants" :key="item.value" :label="item.label" :value="item.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem v-if="isRegionAccount" label="所属区域" required>
          <ElSelect v-model="form.region_ids" class="w-full" multiple placeholder="请选择区域">
            <ElOption v-for="item in regions" :key="item.value" :label="item.label" :value="item.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem v-if="isRegionAccount" label="关联代理" required>
          <ElSelect v-model="form.circle_agent_id" class="w-full" placeholder="请选择已审核通过的代理">
            <ElOption
              v-for="agent in agentOptions"
              :key="agent.circle_agent_id"
              :label="agent.phone ? `${agent.name} / ${agent.phone}` : agent.name"
              :value="agent.circle_agent_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem v-if="hasServiceRole" label="客服授权店铺" required>
          <ElSelect v-model="form.service_store_ids" class="w-full" multiple placeholder="请选择店铺">
            <ElOption v-for="item in stores" :key="item.value" :label="item.label" :value="item.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="账号状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>

    <PasswordDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="管理员">
          <ElInput :model-value="passwordTarget?.real_name || ''" disabled />
        </ElFormItem>
        <ElFormItem label="新密码" required>
          <ElInput v-model="passwordForm.password" autocomplete="new-password" show-password type="password" />
        </ElFormItem>
        <ElFormItem label="确认密码" required>
          <ElInput v-model="passwordForm.confirm_password" autocomplete="new-password" show-password type="password" />
        </ElFormItem>
      </ElForm>
    </PasswordDrawer>
  </Page>
</template>
