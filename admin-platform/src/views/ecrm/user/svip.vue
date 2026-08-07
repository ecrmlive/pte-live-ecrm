<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElMessage,
  ElOption,
  ElSelect,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformSvipUsersApi,
  setPlatformUserSvipApi,
  type PlatformSvipUser,
} from '#/api/core/platform-svip';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canManage = ref(false);
const saving = ref(false);
const editing = ref<PlatformSvipUser>();
const form = reactive({ is_svip: 0, svip_endtime: '' });

const svipLabels: Record<number, string> = {
  '-1': '已关闭',
  0: '普通用户',
  1: '体验会员',
  2: '有效期会员',
  3: '永久会员',
};

const svipTypes: Record<number, 'info' | 'success' | 'warning'> = {
  '-1': 'info',
  0: 'info',
  1: 'warning',
  2: 'success',
  3: 'success',
};

function svipText(value: number) {
  return svipLabels[value] || '未知';
}

function svipType(value: number) {
  return svipTypes[value] || 'info';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '用户 ID / 昵称',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '普通用户', value: 0 },
        { label: '体验会员', value: 1 },
        { label: '有效期会员', value: 2 },
        { label: '永久会员', value: 3 },
        { label: '已关闭', value: -1 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'is_svip',
    label: '会员状态',
  },
]);

const gridOptions: VxeGridProps<PlatformSvipUser> = {
  columns: [
    { field: 'uid', title: '用户 ID', width: 92 },
    {
      field: 'nickname',
      minWidth: 160,
      showOverflow: false,
      slots: { default: 'user' },
      title: '用户',
    },
    {
      field: 'now_money',
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
      title: '余额',
      width: 110,
    },
    { field: 'integral', title: '积分', width: 90 },
    {
      field: 'is_svip',
      slots: { default: 'svipStatus' },
      title: '会员状态',
      width: 120,
    },
    {
      field: 'svip_endtime',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 170,
      title: '到期时间',
    },
    platformListActionColumn({ width: 76 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canManage.value) return { items: [], total: 0 };
        const result = await listPlatformSvipUsersApi({
          page: page.currentPage,
          limit: page.pageSize,
        });
        let list = result.list || [];
        const keyword = String(formValues?.keyword ?? '').trim().toLowerCase();
        const svipRaw = formValues?.is_svip;
        if (keyword) {
          list = list.filter(
            (row) =>
              String(row.uid).includes(keyword) ||
              (row.nickname || '').toLowerCase().includes(keyword),
          );
        }
        const hasSvipFilter =
          svipRaw === 0 ||
          svipRaw === 1 ||
          svipRaw === 2 ||
          svipRaw === 3 ||
          svipRaw === -1;
        if (hasSvipFilter) {
          list = list.filter((row) => row.is_svip === Number(svipRaw));
        }
        const filtered = keyword || hasSvipFilter;
        return {
          items: list,
          total: filtered ? list.length : result.total || 0,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'uid' },
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

function openEdit(row: PlatformSvipUser) {
  editing.value = row;
  Object.assign(form, {
    is_svip: row.is_svip,
    svip_endtime:
      row.svip_endtime?.slice(0, 19).replace('T', ' ') || '',
  });
  formDrawerApi.setState({ title: '设置付费会员' }).open();
}

async function save() {
  if (!editing.value) return;
  if (form.is_svip === 2 && !form.svip_endtime) {
    ElMessage.warning('有效期会员必须填写到期时间');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    await setPlatformUserSvipApi(editing.value.uid, {
      is_svip: form.is_svip,
      ...(form.is_svip === 2 ? { svip_endtime: form.svip_endtime } : {}),
    });
    formDrawerApi.close();
    ElMessage.success('会员状态已更新');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canManage.value =
    profile.roles.includes('platform') &&
    permissions.includes('user.svip.manage');
  if (canManage.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canManage"
      class="mb-4"
      title="当前账号没有 SVIP 监管权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #user="{ row }">
        <div>{{ row.nickname || '未设置昵称' }}</div>
        <div class="text-xs text-muted-foreground">
          {{ row.phone_masked || '—' }}
        </div>
      </template>
      <template #svipStatus="{ row }">
        <ElTag :type="svipType(row.is_svip)">
          {{ svipText(row.is_svip) }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">设置</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="用户">
          <span>{{
            editing?.nickname || `用户 #${editing?.uid || ''}`
          }}</span>
        </ElFormItem>
        <ElFormItem label="会员类型">
          <ElSelect v-model="form.is_svip" class="w-full">
            <ElOption label="普通用户" :value="0" />
            <ElOption label="体验会员" :value="1" />
            <ElOption label="有效期会员" :value="2" />
            <ElOption label="永久会员" :value="3" />
            <ElOption label="关闭会员" :value="-1" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem
          v-if="form.is_svip === 2"
          label="到期时间"
          required
        >
          <ElDatePicker
            v-model="form.svip_endtime"
            class="w-full"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            type="datetime"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
