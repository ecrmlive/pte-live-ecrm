<script setup lang="ts">
import type { MerchantStaff, MerchantStaffSaveInput } from '#/api/core/staff';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElMessageBox,
  ElPagination,
  ElSwitch,
  ElTable,
  ElTableColumn,
} from 'element-plus';
import { reactive, ref } from 'vue';

import {
  createMerchantStaffApi,
  listMerchantStaffApi,
  removeMerchantStaffApi,
  updateMerchantStaffApi,
} from '#/api/core/staff';

const rows = ref<MerchantStaff[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = 20;
const total = ref(0);
const saving = ref(false);
const editingID = ref<number | null>(null);
const form = reactive<Required<MerchantStaffSaveInput>>({
  account: '',
  is_goods: 1,
  is_open: 1,
  is_verify: 1,
  nickname: '',
  password: '',
  phone: '',
  status: 1,
});
const asStaff = (row: unknown) => row as MerchantStaff;

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

async function load() {
  loading.value = true;
  try {
    const data = await listMerchantStaffApi({ page: page.value, limit: pageSize });
    rows.value = data.list ?? [];
    total.value = data.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  form.account = '';
  form.nickname = '';
  form.password = '';
  form.phone = '';
  form.status = 1;
  form.is_open = 1;
  form.is_verify = 1;
  form.is_goods = 1;
}

function openCreate() {
  editingID.value = null;
  resetForm();
  formDrawerApi.setState({ title: '新增员工' }).open();
}

function openEdit(row: MerchantStaff) {
  editingID.value = row.service_id;
  form.account = row.account;
  form.nickname = row.nickname;
  form.password = '';
  form.phone = row.phone;
  form.status = row.status;
  form.is_open = row.is_open;
  form.is_verify = row.is_verify;
  form.is_goods = row.is_goods;
  formDrawerApi.setState({ title: '编辑员工' }).open();
}

async function save() {
  if (!form.nickname.trim()) {
    ElMessage.warning('请填写员工昵称');
    return;
  }
  if (
    editingID.value == null &&
    (!form.account.trim() || form.password.length < 8)
  ) {
    ElMessage.warning('请填写账号和至少 8 位密码');
    return;
  }
  saving.value = true;
  formDrawerApi.lock();
  try {
    const data: MerchantStaffSaveInput = { ...form };
    if (editingID.value != null) {
      delete data.account;
      if (!data.password) delete data.password;
      await updateMerchantStaffApi(editingID.value, data);
    } else {
      await createMerchantStaffApi(data);
    }
    formDrawerApi.close();
    await load();
    ElMessage.success(
      editingID.value == null ? '员工已创建' : '员工资料已更新',
    );
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function toggle(
  row: MerchantStaff,
  key: 'is_goods' | 'is_open' | 'is_verify' | 'status',
) {
  await updateMerchantStaffApi(row.service_id, {
    [key]: row[key],
  } as MerchantStaffSaveInput);
  ElMessage.success('设置已保存');
}

async function remove(row: MerchantStaff) {
  try {
    await ElMessageBox.confirm(
      `确定移除员工“${row.nickname}”吗？该账号将立即停用并失效。`,
      '移除员工',
      { type: 'warning', confirmButtonText: '移除', cancelButtonText: '取消' },
    );
    await removeMerchantStaffApi(row.service_id);
    ElMessage.success('员工已移除');
    await load();
  } catch {
    /* 取消操作无需提示 */
  }
}

void load();
</script>

<template>
  <Page auto-content-height>
    <div class="mb-4">
      <ElButton type="primary" @click="openCreate">新增员工</ElButton>
    </div>
    <ElTable v-loading="loading" :data="rows" border>
      <ElTableColumn prop="service_id" label="ID" width="80" />
      <ElTableColumn prop="account" label="账号" min-width="120" />
      <ElTableColumn prop="nickname" label="昵称" min-width="130" />
      <ElTableColumn prop="phone" label="手机号" min-width="130" />
      <ElTableColumn label="状态" width="100">
        <template #default="{ row }">
          <ElSwitch
            v-model="row.status"
            :active-value="1"
            :inactive-value="0"
            @change="toggle(asStaff(row), 'status')"
          />
        </template>
      </ElTableColumn>
      <ElTableColumn label="接单" width="100">
        <template #default="{ row }">
          <ElSwitch
            v-model="row.is_open"
            :active-value="1"
            :inactive-value="0"
            @change="toggle(asStaff(row), 'is_open')"
          />
        </template>
      </ElTableColumn>
      <ElTableColumn label="核销" width="100">
        <template #default="{ row }">
          <ElSwitch
            v-model="row.is_verify"
            :active-value="1"
            :inactive-value="0"
            @change="toggle(asStaff(row), 'is_verify')"
          />
        </template>
      </ElTableColumn>
      <ElTableColumn label="发货" width="100">
        <template #default="{ row }">
          <ElSwitch
            v-model="row.is_goods"
            :active-value="1"
            :inactive-value="0"
            @change="toggle(asStaff(row), 'is_goods')"
          />
        </template>
      </ElTableColumn>
      <ElTableColumn label="创建时间" min-width="170">
        <template #default="{ row }">{{ row.create_time }}</template>
      </ElTableColumn>
      <ElTableColumn label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <ElButton link type="primary" @click="openEdit(asStaff(row))"
            >编辑</ElButton
          >
          <ElButton link type="danger" @click="remove(asStaff(row))"
            >移除</ElButton
          >
        </template>
      </ElTableColumn>
    </ElTable>
    <ElPagination
      v-if="total > pageSize"
      v-model:current-page="page"
      class="mt-4"
      :page-size="pageSize"
      :total="total"
      background
      layout="prev, pager, next"
      @current-change="load"
    />

    <FormDrawer>
      <ElForm label-position="top">
        <ElFormItem label="登录账号" required>
          <ElInput v-model="form.account" :disabled="editingID != null" />
        </ElFormItem>
        <ElFormItem label="员工昵称" required>
          <ElInput v-model="form.nickname" />
        </ElFormItem>
        <ElFormItem label="手机号">
          <ElInput v-model="form.phone" />
        </ElFormItem>
        <ElFormItem
          :label="editingID == null ? '登录密码' : '登录密码（留空不修改）'"
          :required="editingID == null"
        >
          <ElInput v-model="form.password" type="password" show-password />
        </ElFormItem>
        <div class="grid grid-cols-2 gap-x-4">
          <ElFormItem label="启用账号">
            <ElSwitch
              v-model="form.status"
              :active-value="1"
              :inactive-value="0"
            />
          </ElFormItem>
          <ElFormItem label="接单权限">
            <ElSwitch
              v-model="form.is_open"
              :active-value="1"
              :inactive-value="0"
            />
          </ElFormItem>
          <ElFormItem label="核销权限">
            <ElSwitch
              v-model="form.is_verify"
              :active-value="1"
              :inactive-value="0"
            />
          </ElFormItem>
          <ElFormItem label="发货权限">
            <ElSwitch
              v-model="form.is_goods"
              :active-value="1"
              :inactive-value="0"
            />
          </ElFormItem>
        </div>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
