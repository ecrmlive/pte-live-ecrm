<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElAvatar,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';
import { Plus, UserFilled } from '@element-plus/icons-vue';
import { useRouter } from 'vue-router';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getUserInfoApi } from '#/api/core/auth';
import {
  fetchCustomerServiceAgentUsers,
  fetchCustomerServiceAgents,
  type CustomerServiceAgent,
  type CustomerServiceAgentUser,
} from '#/api/core/customer-service';
import {
  createPlatformAdmin,
  deletePlatformAdmin,
  updatePlatformAdmin,
} from '#/api/core/ecrm';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';
import ImageField from '#/components/shop/image-field.vue';

import CustomerUserPicker, {
  type PickedCustomerUser,
} from './customer-user-picker.vue';

type AgentForm = {
  account: string;
  avatar_url: string;
  confirm_password: string;
  linked_user_id: number;
  password: string;
  phone: string;
  real_name: string;
  roles: string;
  service_store_ids: string;
  status: number;
};

const router = useRouter();
const canManage = ref(false);
const editing = ref<CustomerServiceAgent>();
const historyRows = ref<CustomerServiceAgentUser[]>([]);
const historyLoading = ref(false);
const userPickerOpen = ref(false);
const selectedCustomerUser = ref<PickedCustomerUser>();
const form = reactive<AgentForm>({
  account: '',
  avatar_url: '',
  confirm_password: '',
  linked_user_id: 0,
  password: '',
  phone: '',
  real_name: '',
  roles: 'customer_service',
  service_store_ids: '',
  status: 1,
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_ENABLE_STATUS_FIELD('显示状态'),
  LIST_KEYWORD_FIELD('请输入关键字'),
]);

const gridOptions: VxeGridProps<CustomerServiceAgent> = {
  columns: [
    { field: 'id', title: 'ID', width: 86 },
    { field: 'wechat_username', minWidth: 190, title: '微信用户名' },
    {
      slots: { default: 'avatar' },
      title: '客服头像',
      width: 112,
    },
    { field: 'display_name', minWidth: 150, title: '客服名称' },
    { field: 'status', slots: { default: 'status' }, title: '账号状态', width: 108 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 180,
      title: '添加时间',
    },
    platformListActionColumn({ minWidth: 280 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const result = await fetchCustomerServiceAgents({
          keyword: String(formValues?.keyword || '').trim() || undefined,
          status: [0, 1].includes(Number(formValues?.status))
            ? Number(formValues?.status)
            : undefined,
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: result.list, total: result.total };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: { custom: false, export: false, refresh: false, search: false, zoom: false },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });
const [AgentDrawer, agentDrawerApi] = useVbenDrawer({
  class: 'w-[720px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});
const [HistoryDrawer, historyDrawerApi] = useVbenDrawer({
  class: 'w-[860px] max-w-[96vw]',
  footer: false,
  placement: 'right',
});

const drawerTitle = computed(() => (editing.value ? '编辑客服' : '新增客服'));

function resetForm() {
  selectedCustomerUser.value = undefined;
  Object.assign(form, {
    account: '', avatar_url: '', confirm_password: '', linked_user_id: 0, password: '', phone: '', real_name: '', roles: 'customer_service', service_store_ids: '', status: 1,
  });
}

function openCreate() {
  editing.value = undefined;
  resetForm();
  agentDrawerApi.setState({ title: '新增客服' }).open();
}

function openEdit(row: CustomerServiceAgent) {
  editing.value = row;
  selectedCustomerUser.value = row.linked_user_id
    ? {
        id: row.linked_user_id,
        nickname: row.wechat_username || row.account,
        avatar_url: row.avatar_url || '',
        mobile: '',
      }
    : undefined;
  Object.assign(form, {
    account: row.account,
    avatar_url: row.avatar_url || '',
    confirm_password: '',
    linked_user_id: row.linked_user_id || 0,
    password: '',
    phone: row.phone,
    real_name: row.display_name,
    roles: row.roles,
    service_store_ids: (row.service_store_ids || []).join(','),
    status: row.status,
  });
  agentDrawerApi.setState({ title: '编辑客服' }).open();
}

async function save() {
  if (!editing.value && !form.linked_user_id) {
    ElMessage.warning('请选择用户');
    return;
  }
  if (!editing.value && (!form.account.trim() || form.password.length < 8)) {
    ElMessage.warning('请填写客服账号及至少 8 位客服密码');
    return;
  }
  if (form.password && form.password !== form.confirm_password) {
    ElMessage.warning('两次输入的密码不一致');
    return;
  }
  if (!form.avatar_url.trim()) {
    ElMessage.warning('请上传客服头像');
    return;
  }
  if (!form.real_name.trim()) {
    ElMessage.warning('请填写客服昵称');
    return;
  }
  agentDrawerApi.lock();
  try {
    const payload = {
      ...form,
      account: form.account.trim(),
      avatar_url: form.avatar_url.trim(),
      linked_user_id: form.linked_user_id,
      phone: form.phone.trim(),
      real_name: form.real_name.trim(),
      roles: form.roles.trim(),
      service_store_ids: form.service_store_ids.trim(),
    };
    if (editing.value) await updatePlatformAdmin(editing.value.id, payload);
    else await createPlatformAdmin(payload);
    agentDrawerApi.close();
    ElMessage.success('客服账号已保存');
    gridApi.reload();
  } finally {
    agentDrawerApi.unlock();
  }
}

function selectCustomerUser(user: PickedCustomerUser) {
  selectedCustomerUser.value = user;
  form.linked_user_id = user.id;
  form.avatar_url = user.avatar_url;
  if (!form.real_name.trim()) form.real_name = user.nickname;
  if (!form.phone.trim()) form.phone = user.mobile;
  if (!form.account.trim() && user.mobile) form.account = user.mobile;
}

async function toggleStatus(row: CustomerServiceAgent, enabled: boolean) {
  try {
    await updatePlatformAdmin(row.id, {
      account: row.account,
      avatar_url: row.avatar_url || '',
      linked_user_id: row.linked_user_id || 0,
      phone: row.phone,
      real_name: row.display_name,
      roles: row.roles,
      service_store_ids: row.service_store_ids.join(','),
      status: enabled ? 1 : 0,
    });
    ElMessage.success(enabled ? '客服账号已启用' : '客服账号已停用');
    gridApi.reload();
  } catch {
    gridApi.reload();
  }
}

async function remove(row: CustomerServiceAgent) {
  try {
    await confirm({
      title: '删除客服账号',
      content: `删除“${row.display_name || row.account}”后，该账号无法登录；历史会话和审计记录会保留。是否继续？`,
      icon: 'warning',
    });
    await deletePlatformAdmin(row.id);
    ElMessage.success('客服账号已删除');
    gridApi.reload();
  } catch {
    // 取消或请求失败时，列表状态保持不变。
  }
}

function openWorkbench() {
  void router.push('/service');
}

async function openHistory(row: CustomerServiceAgent) {
  historyRows.value = [];
  historyDrawerApi.setState({ title: `${row.display_name || row.account}的聊天记录` }).open();
  historyLoading.value = true;
  try {
    const result = await fetchCustomerServiceAgentUsers(row.id, { page: 1, limit: 100 });
    historyRows.value = result.list;
  } finally {
    historyLoading.value = false;
  }
}

function isDedicatedAgent(row: CustomerServiceAgent) {
  return row.roles === 'customer_service';
}

onMounted(async () => {
  const profile = await getUserInfoApi();
  canManage.value = profile.roles.includes('platform');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton v-if="canManage" :icon="Plus" type="primary" @click="openCreate">新增客服</ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage || !isDedicatedAgent(row)"
          @update:model-value="toggleStatus(row, $event)"
        />
      </template>
      <template #avatar="{ row }">
        <ElAvatar :icon="UserFilled" :src="row.avatar_url || undefined" shape="square" :size="48" />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openWorkbench">进入工作台</ElButton>
        <ElButton link type="primary" @click="openHistory(row)">聊天记录</ElButton>
        <ElButton v-if="canManage && isDedicatedAgent(row)" link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton v-if="canManage && isDedicatedAgent(row)" link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <AgentDrawer :title="drawerTitle">
      <ElForm label-width="120px">
        <ElFormItem label="用户" required>
          <button class="customer-user-selector" type="button" @click="userPickerOpen = true">
            <ElAvatar
              :icon="UserFilled"
              :size="64"
              :src="selectedCustomerUser?.avatar_url || undefined"
              shape="square"
            />
            <span class="customer-user-selector__content">
              <strong>{{ selectedCustomerUser?.nickname || '选择用户' }}</strong>
              <small v-if="selectedCustomerUser">用户 ID：{{ selectedCustomerUser.id }}</small>
              <small v-else>点击关联用户</small>
            </span>
          </button>
        </ElFormItem>
        <ElFormItem label="客服头像" required>
          <ImageField v-model="form.avatar_url" :empty-icon="UserFilled" :preview-size="96" />
        </ElFormItem>
        <ElFormItem label="客服昵称" required>
          <ElInput v-model="form.real_name" maxlength="64" placeholder="请输入客服昵称" />
        </ElFormItem>
        <ElFormItem label="客服账号" required>
          <ElInput v-model="form.account" :disabled="Boolean(editing)" maxlength="64" />
        </ElFormItem>
        <ElFormItem :label="editing ? '客服密码' : '客服密码'" :required="!editing">
          <ElInput v-model="form.password" type="password" show-password autocomplete="new-password" :placeholder="editing ? '留空则不修改' : '至少 8 位'" />
        </ElFormItem>
        <ElFormItem label="确认密码" :required="Boolean(form.password)">
          <ElInput v-model="form.confirm_password" type="password" show-password autocomplete="new-password" placeholder="请输入确认密码" />
        </ElFormItem>
        <ElFormItem label="通知电话">
          <ElInput v-model="form.phone" maxlength="32" />
        </ElFormItem>
        <ElFormItem label="账号状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </AgentDrawer>

    <CustomerUserPicker v-model:open="userPickerOpen" @select="selectCustomerUser" />

    <HistoryDrawer>
      <el-alert class="mb-4" type="info" :closable="false" title="仅展示该客服已领取会话的脱敏用户信息与业务会话状态；聊天正文由 pte-live-im 提供。" />
      <ElTable v-loading="historyLoading" :data="historyRows">
        <ElTableColumn prop="binding_id" label="会话编号" width="100" />
        <ElTableColumn prop="nickname" label="用户昵称" min-width="150" />
        <ElTableColumn prop="mobile" label="手机号" width="140" />
        <ElTableColumn prop="store_name" label="店铺" min-width="150" />
        <ElTableColumn label="会话状态" width="100">
          <template #default="{ row }"><ElTag :type="row.status === 'open' ? 'warning' : 'info'">{{ row.status === 'open' ? '进行中' : '已关闭' }}</ElTag></template>
        </ElTableColumn>
        <ElTableColumn label="最近更新" min-width="180">
          <template #default="{ row }">{{ formatShanghaiDateTime(row.updated_at) }}</template>
        </ElTableColumn>
      </ElTable>
    </HistoryDrawer>
  </Page>
</template>

<style scoped>
.customer-user-selector {
  display: inline-flex;
  width: 260px;
  height: 96px;
  align-items: center;
  justify-content: flex-start;
  gap: 12px;
  overflow: hidden;
  padding: 12px;
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  color: var(--el-text-color-secondary);
  background: var(--el-fill-color-blank);
  cursor: pointer;
}

.customer-user-selector__content {
  display: flex;
  min-width: 0;
  flex: 1;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  color: var(--el-text-color-primary);
  text-align: left;
}

.customer-user-selector__content strong,
.customer-user-selector__content small {
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.customer-user-selector__content small {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  font-weight: 400;
}

.customer-user-selector:hover {
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
}
</style>
