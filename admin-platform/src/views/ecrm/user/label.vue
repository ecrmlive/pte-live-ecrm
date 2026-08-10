<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElForm, ElFormItem, ElInput, ElMessage } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createUserLabel,
  deleteUserLabel,
  fetchUserLabels,
  updateUserLabel,
  type UserLabelRow,
} from '#/api/core/ecrm';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

const canRead = ref(false);
const canManage = ref(false);
const saving = ref(false);
const editing = ref<UserLabelRow>();
const form = reactive({ label_name: '' });

const gridOptions: VxeGridProps<UserLabelRow> = {
  columns: [
    { field: 'label_id', title: 'ID', width: 88 },
    {
      field: 'label_name',
      minWidth: 220,
      showOverflow: false,
      title: '标签名称',
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
        const result = await fetchUserLabels({
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'label_id' },
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
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  form.label_name = '';
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加用户标签' }).open();
}

function openEdit(row: UserLabelRow) {
  editing.value = row;
  form.label_name = row.label_name;
  formDrawerApi.setState({ title: '编辑用户标签' }).open();
}

async function save() {
  const name = form.label_name.trim();
  if (!name) {
    ElMessage.warning('请填写标签名称');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const body = { label_name: name, sort: editing.value?.sort ?? 0 };
    if (editing.value) {
      await updateUserLabel(editing.value.label_id, body);
    } else {
      await createUserLabel(body);
    }
    formDrawerApi.close();
    ElMessage.success(editing.value ? '用户标签已更新' : '用户标签已创建');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function remove(row: UserLabelRow) {
  try {
    await confirm({
      content: `删除用户标签“${row.label_name}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteUserLabel(row.label_id);
    ElMessage.success('用户标签已删除');
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
    (codes.includes('user.label.read') || codes.includes('user.label.manage'));
  canManage.value = roleOK && codes.includes('user.label.manage');
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
          添加用户标签
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
        <ElFormItem label="标签名称" required>
          <ElInput
            v-model="form.label_name"
            maxlength="32"
            placeholder="请输入用户标签名称"
            show-word-limit
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
