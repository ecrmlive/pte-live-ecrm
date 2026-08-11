<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElForm, ElFormItem, ElInput, ElMessage } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createPlatformUserGroupApi,
  deletePlatformUserGroupApi,
  listPlatformUserGroupsApi,
  updatePlatformUserGroupApi,
  type PlatformUserGroup,
} from '#/api/core/platform-user-group';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

const canRead = ref(false);
const canManage = ref(false);
const saving = ref(false);
const editing = ref<PlatformUserGroup>();
const form = reactive({ group_name: '' });

const gridOptions: VxeGridProps<PlatformUserGroup> = {
  columns: [
    { field: 'group_id', title: 'ID', width: 88 },
    {
      field: 'group_name',
      minWidth: 220,
      showOverflow: false,
      title: '分组名称',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '创建时间',
    },
    platformListActionColumn({ width: 130 }),
  ],
  emptyText: '暂无数据',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const result = await listPlatformUserGroupsApi({
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'group_id' },
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
  editing.value = undefined;
  form.group_name = '';
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增用户分组' }).open();
}

function openEdit(row: PlatformUserGroup) {
  editing.value = row;
  form.group_name = row.group_name;
  formDrawerApi.setState({ title: '编辑用户分组' }).open();
}

async function save() {
  const name = form.group_name.trim();
  if (!name) {
    ElMessage.warning('请填写分组名称');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const body = { group_name: name, sort: editing.value?.sort ?? 0 };
    if (editing.value) {
      await updatePlatformUserGroupApi(editing.value.group_id, body);
    } else {
      await createPlatformUserGroupApi(body);
    }
    formDrawerApi.close();
    ElMessage.success(editing.value ? '用户分组已更新' : '用户分组已创建');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function remove(row: PlatformUserGroup) {
  try {
    await confirm({
      content: `删除用户分组“${row.group_name}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '提示',
    });
    await deletePlatformUserGroupApi(row.group_id);
    ElMessage.success('用户分组已删除');
    gridApi.reload();
  } catch {
    /* 用户取消或统一请求层处理 */
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
    (codes.includes('user.group.read') || codes.includes('user.group.manage'));
  canManage.value = roleOK && codes.includes('user.group.manage');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          新增用户分组
        </ElButton>
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
      <ElForm label-width="100px">
        <ElFormItem label="分组名称" required>
          <ElInput
            v-model="form.group_name"
            maxlength="32"
            placeholder="请输入用户分组名称"
            show-word-limit
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
