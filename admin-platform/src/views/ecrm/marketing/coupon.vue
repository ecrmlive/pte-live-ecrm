<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElSwitch,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformCouponApi,
  deletePlatformCouponApi,
  listPlatformCouponsApi,
  setPlatformCouponStatusApi,
  updatePlatformCouponApi,
  type PlatformCoupon,
  type PlatformCouponSaveInput,
} from '#/api/core/platform-promotion';
import { getAccessCodesApi } from '#/api/core/auth';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const editingID = ref<number>();
const canManage = ref(false);

const form = reactive<PlatformCouponSaveInput>({
  coupon_price: 0,
  coupon_time: 30,
  is_limited: 0,
  sort: 1,
  status: 1,
  title: '',
  total_count: 0,
  use_min_price: 0,
});

function formatTime(value?: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const statusRaw = formValues?.status;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    status:
      statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '优惠券名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '已启用', value: 1 },
        { label: '已停用', value: 0 },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '启用状态',
  },
]);

const gridOptions: VxeGridProps<PlatformCoupon> = {
  columns: [
    { field: 'coupon_id', title: 'ID', width: 80 },
    { field: 'title', minWidth: 180, showOverflow: false, title: '优惠券名称' },
    {
      field: 'coupon_price',
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
      title: '面额',
      width: 108,
    },
    {
      field: 'use_min_price',
      formatter: ({ cellValue }) => `满 ¥${Number(cellValue || 0).toFixed(2)} 可用`,
      title: '使用门槛',
      width: 140,
    },
    {
      field: 'coupon_time',
      formatter: ({ cellValue }) => `领取后 ${cellValue} 天`,
      title: '有效期',
      width: 112,
    },
    {
      field: 'remain_count',
      formatter: ({ row }) =>
        row.is_limited === 1 ? `${row.remain_count}/${row.total_count}` : '不限量',
      title: '发放数量',
      width: 112,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 88,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 166 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listPlatformCouponsApi(buildListParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'coupon_id' },
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

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    coupon_price: 0,
    coupon_time: 30,
    is_limited: 0,
    sort: 1,
    status: 1,
    title: '',
    total_count: 0,
    use_min_price: 0,
  });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增平台券' }).open();
}

function openEdit(row: PlatformCoupon) {
  editingID.value = row.coupon_id;
  Object.assign(form, {
    coupon_price: row.coupon_price,
    coupon_time: row.coupon_time,
    is_limited: row.is_limited,
    sort: row.sort,
    status: row.status,
    title: row.title,
    total_count: row.total_count,
    use_min_price: row.use_min_price,
  });
  formDrawerApi.setState({ title: '编辑平台券' }).open();
}

async function save() {
  if (
    !form.title.trim() ||
    form.coupon_price <= 0 ||
    form.use_min_price < 0 ||
    form.coupon_time <= 0
  ) {
    ElMessage.warning('请完整填写优惠券名称、金额、门槛和有效天数');
    return;
  }
  if (form.is_limited === 1 && form.total_count <= 0) {
    ElMessage.warning('限量发放时必须填写发放总数');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const body = {
      ...form,
      title: form.title.trim(),
      total_count: form.is_limited === 1 ? form.total_count : 0,
    };
    if (editingID.value) await updatePlatformCouponApi(editingID.value, body);
    else await createPlatformCouponApi(body);
    formDrawerApi.close();
    ElMessage.success(editingID.value ? '平台券已更新' : '平台券已创建');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function toggle(row: PlatformCoupon) {
  const next = row.status === 1 ? 0 : 1;
  await setPlatformCouponStatusApi(row.coupon_id, next);
  row.status = next;
  ElMessage.success(next === 1 ? '平台券已启用' : '平台券已停用');
}

async function remove(row: PlatformCoupon) {
  try {
    await ElMessageBox.confirm(
      `删除平台券“${row.title}”后不可恢复，是否继续？`,
      '删除确认',
      { type: 'warning' },
    );
    await deletePlatformCouponApi(row.coupon_id);
    ElMessage.success('平台券已删除');
    gridApi.reload();
  } catch {
    /* 用户取消或统一请求层处理 */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('marketing.coupon.manage');
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
          新增平台券
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '已启用' : '已停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="warning" @click="toggle(row)">
            {{ row.status === 1 ? '停用' : '启用' }}
          </ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm class="grid grid-cols-2 gap-x-4" label-width="90px">
        <ElFormItem class="col-span-2" label="优惠券名称" required>
          <ElInput v-model="form.title" maxlength="40" show-word-limit />
        </ElFormItem>
        <ElFormItem label="优惠金额" required>
          <ElInputNumber
            v-model="form.coupon_price"
            :min="0.01"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="使用门槛">
          <ElInputNumber
            v-model="form.use_min_price"
            :min="0"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="有效天数" required>
          <ElInputNumber v-model="form.coupon_time" :min="1" class="w-full" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="限量发放">
          <ElSwitch v-model="form.is_limited" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem v-if="form.is_limited === 1" label="发放总数" required>
          <ElInputNumber v-model="form.total_count" :min="1" class="w-full" />
        </ElFormItem>
        <ElFormItem label="初始状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
