<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElRadio,
  ElRadioGroup,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  deletePlatformCombinationApi,
  getPlatformCombinationApi,
  listPlatformCombinationsApi,
  updatePlatformCombinationApi,
  type PlatformCombination,
  type PlatformCombinationInput,
} from '#/api/core/platform-combination';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  buildStandardListParams,
  LIST_DATE_RANGE_FIELD,
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  LIST_MER_ID_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const canManage = ref(false);
const editingID = ref<number>();
const form = reactive<Required<PlatformCombinationInput>>({
  price: 0,
  buying_count_num: 2,
  time: 24,
  start_time: '',
  end_time: '',
  is_show: 1,
  status: 1,
});

function dateTime(value: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  LIST_KEYWORD_FIELD('商品名称'),
  LIST_MER_ID_FIELD,
  LIST_ENABLE_STATUS_FIELD(),
]);

const gridOptions: VxeGridProps<PlatformCombination> = {
  columns: [
    { field: 'product_group_id', title: 'ID', width: 80 },
    {
      field: 'store_name',
      minWidth: 180,
      showOverflow: false,
      title: '商品',
      formatter: ({ row }) => row.store_name || `商品 #${row.product_id}`,
    },
    {
      field: 'mer_name',
      minWidth: 130,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => row.mer_name || `商户 #${row.mer_id}`,
    },
    {
      field: 'price',
      title: '拼团价',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    { field: 'buying_count_num', title: '成团人数', width: 100 },
    {
      field: 'start_time',
      minWidth: 220,
      showOverflow: false,
      title: '活动时间',
      formatter: ({ row }) => `${dateTime(row.start_time)} 至 ${dateTime(row.end_time)}`,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 88,
    },
    platformListActionColumn({ width: 172 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformCombinationsApi(buildStandardListParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_group_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [EditDrawer, editDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function toFormDateTime(value: string) {
  return value ? value.replace('T', ' ').slice(0, 19) : '';
}

async function edit(row: PlatformCombination) {
  const detail = await getPlatformCombinationApi(row.product_group_id);
  editingID.value = row.product_group_id;
  Object.assign(form, {
    price: Number(detail.price),
    buying_count_num: detail.buying_count_num,
    time: detail.time || 24,
    start_time: toFormDateTime(detail.start_time),
    end_time: toFormDateTime(detail.end_time),
    is_show: detail.is_show,
    status: detail.status,
  });
  editDrawerApi.setState({ title: '编辑拼团活动' }).open();
}

async function save() {
  if (
    !editingID.value ||
    form.price <= 0 ||
    form.buying_count_num < 2 ||
    form.time < 1 ||
    !form.start_time ||
    !form.end_time ||
    new Date(form.end_time).valueOf() < new Date(form.start_time).valueOf()
  ) {
    ElMessage.warning('请填写正数拼团价、至少 2 人、有效时长和正确的活动时间');
    return;
  }
  editDrawerApi.lock();
  saving.value = true;
  try {
    await updatePlatformCombinationApi(editingID.value, { ...form });
    editDrawerApi.close();
    ElMessage.success('拼团活动已更新');
    gridApi.reload();
  } finally {
    saving.value = false;
    editDrawerApi.unlock();
  }
}

async function setStatus(row: PlatformCombination, status: number) {
  const action = status === 1 ? '启用' : '停用';
  try {
    await ElMessageBox.confirm(
      `确认${action}商品 #${row.product_id} 的拼团活动？`,
      `${action}确认`,
      { cancelButtonText: '取消', confirmButtonText: `确认${action}`, type: 'warning' },
    );
    await updatePlatformCombinationApi(row.product_group_id, { status });
    ElMessage.success(`活动已${action}`);
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

async function remove(row: PlatformCombination) {
  try {
    await ElMessageBox.confirm(
      `删除商品 #${row.product_id} 的拼团配置只会软删除活动，进行中或已完成团单不会被改写。是否继续？`,
      '删除拼团活动',
      { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' },
    );
    await deletePlatformCombinationApi(row.product_group_id);
    ElMessage.success('拼团活动已软删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canManage.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    permissions.includes('marketing.combination.manage');
});
</script>

<template>
  <Page
    auto-content-height
    description="监管各商户拼团活动；运营可维护活动价格、人数与时间配置或软删除，商品归属与已产生团单不可在此变更。"
    title="拼团监管"
  >
    <Grid>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="edit(row)">编辑</ElButton>
          <ElButton
            link
            :type="row.status === 1 ? 'danger' : 'success'"
            @click="setStatus(row, row.status === 1 ? 0 : 1)"
          >
            {{ row.status === 1 ? '停用' : '启用' }}
          </ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <EditDrawer class="w-[620px]">
      <ElAlert
        class="mb-4"
        type="warning"
        :closable="false"
        title="仅维护拼团配置；商品、商户与已产生团单的成员、价格快照不可在此修改。"
      />
      <ElForm label-width="118px">
        <ElFormItem label="拼团价" required>
          <ElInputNumber v-model="form.price" :min="0.01" :precision="2" :step="1" />
        </ElFormItem>
        <ElFormItem label="成团人数" required>
          <ElInputNumber v-model="form.buying_count_num" :min="2" :max="9999" />
        </ElFormItem>
        <ElFormItem label="成团时限（小时）" required>
          <ElInputNumber v-model="form.time" :min="1" :max="720" />
        </ElFormItem>
        <ElFormItem label="活动时间" required>
          <ElInput v-model="form.start_time" class="!w-48" placeholder="YYYY-MM-DD HH:mm:ss" />
          <span class="mx-2">至</span>
          <ElInput v-model="form.end_time" class="!w-48" placeholder="YYYY-MM-DD HH:mm:ss" />
        </ElFormItem>
        <ElFormItem label="前台展示">
          <ElRadioGroup v-model="form.is_show">
            <ElRadio :value="1">上架</ElRadio>
            <ElRadio :value="0">下架</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="活动状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio :value="1">启用</ElRadio>
            <ElRadio :value="0">停用</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
      </ElForm>
    </EditDrawer>
  </Page>
</template>
