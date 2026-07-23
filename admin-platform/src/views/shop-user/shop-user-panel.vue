<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref, watch } from 'vue';

import {
  ElButton,
  ElMessage,
  ElMessageBox,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import PlatformShopUserApi, {
  shopUserRoleNames,
  type ShopUserRoleOption,
  type ShopUserRow,
} from '#/api/core/platform-shop-user';
import { Page } from '@vben/common-ui';
import { parseApiListPage } from '#/utils/list-response';

import ShopUserAddModal from './shop-user-add-modal.vue';
import ShopUserEditModal from './shop-user-edit-modal.vue';

const props = defineProps<{
  appId: number;
  appName?: string;
}>();

const roleOptions = ref<ShopUserRoleOption[]>([]);
const addModalOpen = ref(false);
const editModalOpen = ref(false);
const currentUserId = ref<number>();

async function fetchUserPage(pageSize: number, currentPage: number) {
  if (props.appId <= 0) {
    return { items: [], total: 0 };
  }
  const res = await PlatformShopUserApi.index(
    {
      app_id: props.appId,
      page: currentPage,
      list_rows: pageSize,
    },
    true,
  );
  const page = parseApiListPage<ShopUserRow>({ list: res.data?.list });
  roleOptions.value = res.data?.roleList ?? [];
  return { items: page.list, total: page.total };
}

const gridOptions = {
  border: true,
  columns: [
    { field: 'shop_user_id', title: 'ID', width: 72 },
    { field: 'user_name', minWidth: 100, title: '用户名' },
    { field: 'real_name', minWidth: 90, title: '姓名' },
    {
      field: 'is_super',
      slots: { default: 'userType' },
      title: '类型',
      width: 80,
    },
    {
      field: 'is_status',
      slots: { default: 'status' },
      title: '状态',
      width: 72,
    },
    {
      field: 'role_names',
      minWidth: 140,
      slots: { default: 'roles' },
      title: '所属角色',
    },
    { field: 'create_time', minWidth: 150, title: '添加时间' },
    {
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 220,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 15, 30] },
  proxyConfig: {
    ajax: {
      query: async ({ page }) =>
        fetchUserPage(page.pageSize, page.currentPage),
    },
  },
  rowConfig: { isHover: true, keyField: 'shop_user_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function reload() {
  gridApi.reload();
}

function openAdd() {
  addModalOpen.value = true;
}

function isShopSuper(row: ShopUserRow) {
  return Number(row.is_super) === 1;
}

function isShopUserEnabled(row: ShopUserRow) {
  return Number(row.is_status ?? 1) === 1;
}

function openEdit(row: ShopUserRow) {
  if (isShopSuper(row)) {
    ElMessage.info('超管账号请在「商城列表」中编辑');
    return;
  }
  currentUserId.value = row.shop_user_id;
  editModalOpen.value = true;
}

async function toggleStatus(row: ShopUserRow) {
  if (isShopSuper(row)) {
    ElMessage.info('超管账号不可禁用');
    return;
  }
  const nextStatus = isShopUserEnabled(row) ? 0 : 1;
  const action = nextStatus === 1 ? '启用' : '禁用';
  try {
    await ElMessageBox.confirm(
      `确认${action}管理员「${row.user_name}」吗？${nextStatus === 0 ? '禁用后将无法登录商户后台。' : ''}`,
      '提示',
      { cancelButtonText: '取消', confirmButtonText: '确定', type: 'warning' },
    );
  } catch {
    return;
  }
  const res = await PlatformShopUserApi.setStatus(
    props.appId,
    row.shop_user_id,
    nextStatus as 0 | 1,
    true,
  );
  if (res.code === 1) {
    ElMessage.success(res.msg || `${action}成功`);
    reload();
  }
}

async function deleteUser(row: ShopUserRow) {
  if (isShopSuper(row)) {
    ElMessage.info('超管账号不可删除');
    return;
  }
  try {
    await ElMessageBox.confirm('删除后不可恢复，确认删除该管理员吗?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
  } catch {
    return;
  }
  const res = await PlatformShopUserApi.delete(props.appId, row.shop_user_id, true);
  if (res.code === 1) {
    ElMessage.success(res.msg || '删除成功');
    reload();
  }
}

watch(
  () => props.appId,
  (appId) => {
    if (appId > 0) {
      reload();
    }
  },
);
</script>

<template>
  <div>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-access:code="'platform:shopUser:add'"
          :icon="Plus"
          size="small"
          type="primary"
          @click="openAdd"
        >
          添加管理员
        </ElButton>
      </template>

      <template #toolbar-tools>
        <span class="text-xs text-muted-foreground">
          管理商户后台（/shop）登录账号；超管请在「商城列表」编辑
        </span>
      </template>

      <template #userType="{ row }">
        <ElTag v-if="isShopSuper(row)" size="small" type="warning">超管</ElTag>
        <ElTag v-else size="small" type="info">子账号</ElTag>
      </template>

      <template #status="{ row }">
        <ElTag v-if="!isShopUserEnabled(row)" size="small" type="danger">禁用</ElTag>
        <ElTag v-else size="small" type="success">启用</ElTag>
      </template>

      <template #roles="{ row }">
        <template v-if="isShopSuper(row) && !shopUserRoleNames(row).length">
          <span class="text-muted-foreground">超级管理员</span>
        </template>
        <template v-else-if="shopUserRoleNames(row).length">
          <ElTag
            v-for="name in shopUserRoleNames(row)"
            :key="name"
            class="mr-1"
            size="small"
          >
            {{ name }}
          </ElTag>
        </template>
        <span v-else class="text-muted-foreground">未分配</span>
      </template>

      <template #action="{ row }">
        <template v-if="!isShopSuper(row)">
          <ElButton
            v-access:code="'platform:shopUser:edit'"
            link
            size="small"
            type="primary"
            @click="openEdit(row)"
          >
            编辑
          </ElButton>
          <ElButton
            v-if="!isShopUserEnabled(row)"
            v-access:code="'platform:shopUser:status'"
            link
            size="small"
            type="success"
            @click="toggleStatus(row)"
          >
            启用
          </ElButton>
          <ElButton
            v-else
            v-access:code="'platform:shopUser:status'"
            link
            size="small"
            type="warning"
            @click="toggleStatus(row)"
          >
            禁用
          </ElButton>
          <ElButton
            v-access:code="'platform:shopUser:delete'"
            link
            size="small"
            type="danger"
            @click="deleteUser(row)"
          >
            删除
          </ElButton>
        </template>
        <span v-else class="text-xs text-muted-foreground">商城列表编辑</span>
      </template>
    </Grid>

    <ShopUserAddModal
      v-model:open="addModalOpen"
      :app-id="appId"
      :role-options="roleOptions"
      @success="reload"
    />
    <ShopUserEditModal
      v-model:open="editModalOpen"
      :app-id="appId"
      :shop-user-id="currentUserId"
      @success="reload"
    />
  </div>
</template>
