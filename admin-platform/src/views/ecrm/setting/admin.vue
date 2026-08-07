<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformAdmin,
  deletePlatformAdmin,
  fetchBusinessZoneAgentOptions,
  fetchPlatformAdmins,
  updatePlatformAdmin,
  type PlatformAdminRow,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const editing = ref<PlatformAdminRow>();
const agentOptions = ref<
  Array<{ circle_agent_id: number; name: string; type: 0 | 1 }>
>([]);
const form = reactive({
  account: '',
  password: '',
  real_name: '',
  phone: '',
  roles: '',
  status: 1,
  merchant_ids: '',
  region_ids: '',
  service_store_ids: '',
  circle_agent_id: 0,
});

const isRegionAccount = computed(() => {
  const roles = form.roles.split(',').map((value) => value.trim());
  return roles.includes('region') && !roles.includes('platform');
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('账号 / 姓名 / 手机'),
  LIST_ENABLE_STATUS_FIELD('状态'),
]);

const gridOptions: VxeGridProps<PlatformAdminRow> = {
  columns: [
    { field: 'admin_id', title: 'ID', width: 72 },
    { field: 'account', minWidth: 130, title: '账号' },
    { field: 'real_name', minWidth: 110, title: '姓名' },
    { field: 'phone', title: '手机号', width: 140 },
    { field: 'roles', minWidth: 144, title: '角色' },
    {
      field: 'merchant_ids',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 120,
      title: '授权商户',
    },
    {
      field: 'region_ids',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 132,
      title: '可管理区域',
    },
    {
      field: 'service_store_ids',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 148,
      title: '客服授权店铺',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ width: 130 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const statusRaw = formValues?.status;
        const result = await fetchPlatformAdmins({
          page: page.currentPage,
          limit: page.pageSize,
        });
        let list = result.list || [];
        if (keyword) {
          list = list.filter(
            (row) =>
              row.account.toLowerCase().includes(keyword) ||
              row.real_name.toLowerCase().includes(keyword) ||
              row.phone.includes(keyword),
          );
        }
        if (statusRaw === 0 || statusRaw === 1) {
          list = list.filter((row) => row.status === Number(statusRaw));
        }
        return {
          items: list,
          total:
            keyword || statusRaw === 0 || statusRaw === 1
              ? list.length
              : result.total,
        };
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
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

const modalTitle = computed(() =>
  editing.value ? '编辑管理员' : '新增管理员',
);

function resetForm() {
  Object.assign(form, {
    account: '',
    password: '',
    real_name: '',
    phone: '',
    roles: '',
    status: 1,
    merchant_ids: '',
    region_ids: '',
    service_store_ids: '',
    circle_agent_id: 0,
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
    password: '',
    real_name: row.real_name,
    phone: row.phone,
    roles: row.roles,
    status: row.status,
    merchant_ids: row.merchant_ids,
    region_ids: row.region_ids,
    service_store_ids: row.service_store_ids,
    circle_agent_id: row.circle_agent_id,
  });
  formDrawerApi.setState({ title: '编辑管理员' }).open();
}

async function save() {
  if (
    !editing.value &&
    (!form.account.trim() || form.password.length < 8)
  ) {
    ElMessage.warning('账号与至少 8 位的初始密码必填');
    return;
  }
  if (
    form.roles
      .split(',')
      .map((value) => value.trim())
      .includes('merchant') &&
    !form.merchant_ids.trim()
  ) {
    ElMessage.warning('商户角色必须填写授权商户 ID，多个以逗号分隔');
    return;
  }
  if (isRegionAccount.value && !form.region_ids.trim()) {
    ElMessage.warning('区域管理员必须填写可管理区域 ID，多个以逗号分隔');
    return;
  }
  if (isRegionAccount.value && !form.circle_agent_id) {
    ElMessage.warning('区域管理员必须关联一名已审核通过的区域代理');
    return;
  }
  if (
    form.roles
      .split(',')
      .map((value) => value.trim())
      .includes('customer_service') &&
    !form.service_store_ids.trim()
  ) {
    ElMessage.warning('客服账号必须填写授权店铺 ID，多个以逗号分隔');
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = {
      ...form,
      account: form.account.trim(),
      real_name: form.real_name.trim(),
      phone: form.phone.trim(),
      roles: form.roles.trim(),
      merchant_ids: form.merchant_ids.trim(),
      region_ids: form.region_ids.trim(),
      service_store_ids: form.service_store_ids.trim(),
    };
    if (editing.value) {
      await updatePlatformAdmin(editing.value.admin_id, payload);
    } else {
      await createPlatformAdmin(payload);
    }
    formDrawerApi.close();
    ElMessage.success('管理员已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function remove(row: PlatformAdminRow) {
  try {
    await ElMessageBox.confirm(
      `删除“${row.real_name || row.account}”后，该账号将无法登录，历史操作与审计记录会保留。是否继续？`,
      '逻辑删除管理员',
      {
        type: 'warning',
        confirmButtonText: '删除',
        cancelButtonText: '取消',
      },
    );
    await deletePlatformAdmin(row.admin_id);
    ElMessage.success('管理员已逻辑删除并强制失效');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

onMounted(async () => {
  agentOptions.value =
    (await fetchBusinessZoneAgentOptions()).list || [];
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
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'danger'">
          {{ row.status === 1 ? '启用' : '禁用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer :title="modalTitle">
      <ElForm label-width="120px">
        <ElFormItem label="登录账号" required>
          <ElInput v-model="form.account" :disabled="!!editing" />
        </ElFormItem>
        <ElFormItem
          :label="editing ? '重置密码' : '初始密码'"
          :required="!editing"
        >
          <ElInput
            v-model="form.password"
            autocomplete="new-password"
            show-password
            type="password"
            :placeholder="editing ? '留空则不修改' : '至少 8 位'"
          />
        </ElFormItem>
        <ElFormItem label="姓名">
          <ElInput v-model="form.real_name" />
        </ElFormItem>
        <ElFormItem label="手机号">
          <ElInput v-model="form.phone" />
        </ElFormItem>
        <ElFormItem label="角色代码" required>
          <ElInput
            v-model="form.roles"
            placeholder="例如 merchant；多个以逗号分隔"
          />
        </ElFormItem>
        <ElFormItem
          v-if="
            form.roles
              .split(',')
              .map((value) => value.trim())
              .includes('merchant')
          "
          label="授权商户"
          required
        >
          <ElInput
            v-model="form.merchant_ids"
            placeholder="例如 2001,2002（merchant_id）"
          />
        </ElFormItem>
        <ElFormItem v-if="isRegionAccount" label="可管理区域" required>
          <ElInput
            v-model="form.region_ids"
            placeholder="例如 1,2,3（商户 region_id）"
          />
        </ElFormItem>
        <ElFormItem
          v-if="
            form.roles
              .split(',')
              .map((value) => value.trim())
              .includes('customer_service')
          "
          label="客服授权店铺"
          required
        >
          <ElInput
            v-model="form.service_store_ids"
            placeholder="例如 1001,1002（店铺 store_id）"
          />
        </ElFormItem>
        <ElFormItem v-if="isRegionAccount" label="关联代理" required>
          <ElSelect
            v-model="form.circle_agent_id"
            class="w-full"
            filterable
            placeholder="请选择已审核通过的代理"
          >
            <ElOption
              v-for="agent in agentOptions"
              :key="agent.circle_agent_id"
              :label="`${agent.name}（ID ${agent.circle_agent_id}）`"
              :value="agent.circle_agent_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="启用状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
