<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createMerchantCouponApi,
  deleteMerchantCouponApi,
  listMerchantCouponsApi,
  sendMerchantCouponApi,
  setMerchantCouponStatusApi,
  updateMerchantCouponApi,
  type MerchantCoupon,
  type MerchantCouponSaveInput,
} from '#/api/core/merchant-promotion';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const editingID = ref<number>();
const sendingCoupon = ref<MerchantCoupon>();
const form = reactive<MerchantCouponSaveInput>({
  coupon_price: 0,
  coupon_time: 30,
  is_limited: 0,
  sort: 1,
  status: 1,
  title: '',
  total_count: 0,
  use_min_price: 0,
});
const sendForm = reactive({ mark: '', uidsText: '' });

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

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '优惠券名称' },
    fieldName: 'keyword',
    label: '优惠券搜索',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '已启用', value: 1 },
        { label: '已停用', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<MerchantCoupon> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'coupon_id', title: 'ID', width: 80 },
    { field: 'title', minWidth: 180, showOverflow: false, title: '优惠券名称' },
    {
      field: 'coupon_price',
      title: '面额',
      width: 104,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'use_min_price',
      minWidth: 130,
      title: '使用门槛',
      formatter: ({ cellValue }) => `满 ¥${Number(cellValue || 0).toFixed(2)} 可用`,
    },
    {
      field: 'coupon_time',
      title: '有效期',
      width: 108,
      formatter: ({ cellValue }) => `领取后 ${cellValue} 天`,
    },
    {
      field: 'remain_count',
      minWidth: 112,
      showOverflow: false,
      slots: { default: 'quota' },
      title: '发放数量',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 88,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 222 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = formValues?.status;
        const data = await listMerchantCouponsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status: status === 0 || status === 1 ? Number(status) : undefined,
          date_from: range[0],
          date_to: range[1],
        });
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

const [CouponDrawer, couponDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => {
    if (!form.title.trim()) {
      ElMessage.warning('请填写优惠券名称');
      return;
    }
    if (
      form.coupon_price <= 0 ||
      form.use_min_price < 0 ||
      form.coupon_time <= 0
    ) {
      ElMessage.warning('请检查优惠金额、使用门槛和有效天数');
      return;
    }
    if (form.is_limited === 1 && form.total_count <= 0) {
      ElMessage.warning('限量发放时必须填写发放总数');
      return;
    }
    saving.value = true;
    couponDrawerApi.lock();
    try {
      const body = {
        ...form,
        title: form.title.trim(),
        total_count: form.is_limited === 1 ? form.total_count : 0,
      };
      if (editingID.value) {
        await updateMerchantCouponApi(editingID.value, body);
      } else {
        await createMerchantCouponApi(body);
      }
      ElMessage.success(editingID.value ? '优惠券已更新' : '优惠券已创建');
      couponDrawerApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      couponDrawerApi.unlock();
    }
  },
});

const [SendDrawer, sendDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onConfirm: async () => {
    if (!sendingCoupon.value) return;
    const uids = [
      ...new Set(
        sendForm.uidsText
          .split(/[，,\s]+/)
          .map((value) => Number(value))
          .filter((value) => Number.isSafeInteger(value) && value > 0),
      ),
    ];
    if (uids.length === 0 || uids.length > 100) {
      ElMessage.warning('请填写 1 至 100 个用户 ID，并用逗号或换行分隔');
      return;
    }
    saving.value = true;
    sendDrawerApi.lock();
    try {
      await sendMerchantCouponApi(sendingCoupon.value.coupon_id, {
        mark: sendForm.mark.trim(),
        uids,
      });
      ElMessage.success('优惠券已发送');
      sendDrawerApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      sendDrawerApi.unlock();
    }
  },
});

function openCreate() {
  resetForm();
  couponDrawerApi.setState({ title: '新增优惠券' }).open();
}

function openEdit(row: MerchantCoupon) {
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
  couponDrawerApi.setState({ title: '编辑优惠券' }).open();
}

function openSend(row: MerchantCoupon) {
  sendingCoupon.value = row;
  Object.assign(sendForm, { mark: '', uidsText: '' });
  sendDrawerApi.setState({ title: '定向发送优惠券' }).open();
}

async function toggle(row: MerchantCoupon) {
  const next = row.status === 1 ? 0 : 1;
  await setMerchantCouponStatusApi(row.coupon_id, next);
  ElMessage.success(next === 1 ? '优惠券已启用' : '优惠券已停用');
  gridApi.reload();
}

async function remove(row: MerchantCoupon) {
  try {
    await confirm({
      content: `删除优惠券“${row.title}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteMerchantCouponApi(row.coupon_id);
    ElMessage.success('优惠券已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" @click="openCreate">新增优惠券</ElButton>
    </template>

    <Grid>
      <template #quota="{ row }">
        {{
          row.is_limited === 1
            ? `${row.remain_count}/${row.total_count}`
            : '不限量'
        }}
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status === 1 ? 'success' : 'info'">
          {{ row.status === 1 ? '已启用' : '已停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton
          link
          type="success"
          :disabled="row.status !== 1"
          @click="openSend(row)"
        >
          发券
        </ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="warning" @click="toggle(row)">
          {{ row.status === 1 ? '停用' : '启用' }}
        </ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <CouponDrawer class="w-[580px] max-w-[96vw]">
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
          <ElSwitch
            v-model="form.is_limited"
            :active-value="1"
            :inactive-value="0"
          />
        </ElFormItem>
        <ElFormItem
          v-if="form.is_limited === 1"
          label="发放总数"
          required
        >
          <ElInputNumber
            v-model="form.total_count"
            :min="1"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="初始状态">
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
          />
        </ElFormItem>
      </ElForm>
    </CouponDrawer>

    <SendDrawer class="w-[540px] max-w-[96vw]">
      <ElAlert
        :closable="false"
        class="mb-4"
        title="仅允许向本商户存在已支付订单的用户发送；同一用户不可重复领取同一张券。"
        type="warning"
      />
      <ElForm label-width="88px">
        <ElFormItem label="优惠券">
          <span>{{ sendingCoupon?.title }}</span>
        </ElFormItem>
        <ElFormItem label="用户 ID" required>
          <ElInput
            v-model="sendForm.uidsText"
            :rows="4"
            placeholder="输入用户 ID，多个 ID 用逗号或换行分隔，最多 100 个"
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem label="发送说明">
          <ElInput
            v-model="sendForm.mark"
            :rows="3"
            maxlength="200"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </SendDrawer>
  </Page>
</template>
