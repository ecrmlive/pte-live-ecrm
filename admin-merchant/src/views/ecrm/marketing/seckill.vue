<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElDatePicker,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listMerchantProductsApi,
  type MerchantProduct,
} from '#/api/core/merchant-catalog';
import {
  createMerchantSeckillActiveApi,
  deleteMerchantSeckillActiveApi,
  listMerchantSeckillActivesApi,
  listMerchantSeckillTimesApi,
  setMerchantSeckillStatusApi,
  updateMerchantSeckillActiveApi,
  type MerchantSeckillActive,
  type MerchantSeckillSaveInput,
  type MerchantSeckillTime,
} from '#/api/core/merchant-seckill';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const editingID = ref<number>();
const products = ref<MerchantProduct[]>([]);
const slots = ref<MerchantSeckillTime[]>([]);
const form = reactive({
  dates: [] as string[],
  name: '',
  once_pay_count: 1,
  product_id: 0,
  seckill_price: 0,
  slot_ids: [] as number[],
  status: 1,
});

function statusType(status: number) {
  return status === 1 ? 'success' : 'info';
}

function slotsText(ids: string) {
  const selected = new Set(ids.split(',').map((id) => Number(id)));
  return (
    slots.value
      .filter((slot) => selected.has(slot.seckill_time_id))
      .map((slot) => slot.title)
      .join('、') || '全部场次'
  );
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    dates: [],
    name: '',
    once_pay_count: 1,
    product_id: 0,
    seckill_price: 0,
    slot_ids: [],
    status: 1,
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '活动名称 / 商品名称',
    },
    fieldName: 'keyword',
    label: '活动搜索',
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

const gridOptions: VxeGridProps<MerchantSeckillActive> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'seckill_active_id', title: 'ID', width: 80 },
    { field: 'name', minWidth: 150, showOverflow: false, title: '活动名称' },
    { field: 'store_name', minWidth: 170, showOverflow: false, title: '商品' },
    {
      field: 'price',
      title: '原价',
      width: 100,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'seckill_price',
      title: '秒杀价',
      width: 105,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'once_pay_count',
      title: '限购',
      width: 88,
      formatter: ({ cellValue }) => `${cellValue} 件`,
    },
    {
      field: 'start_day',
      minWidth: 190,
      showOverflow: false,
      slots: { default: 'dates' },
      title: '活动日期',
    },
    {
      field: 'seckill_time_ids',
      minWidth: 150,
      showOverflow: false,
      slots: { default: 'slots' },
      title: '场次',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 88,
    },
    merchantListActionColumn({ width: 190 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = formValues?.status;
        const data = await listMerchantSeckillActivesApi({
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
  rowConfig: { isHover: true, keyField: 'seckill_active_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [EditModal, editModalApi] = useVbenModal({
  onConfirm: async () => {
    const selectedProduct = products.value.find(
      (item) => item.product_id === form.product_id,
    );
    if (
      !form.name.trim() ||
      !form.product_id ||
      form.seckill_price <= 0 ||
      form.dates.length !== 2
    ) {
      ElMessage.warning('请填写活动名称、商品、秒杀价和活动日期');
      return;
    }
    if (selectedProduct && form.seckill_price >= selectedProduct.price) {
      ElMessage.warning('秒杀价应低于商品销售价');
      return;
    }
    const body: MerchantSeckillSaveInput = {
      end_day: form.dates[1]!,
      name: form.name.trim(),
      once_pay_count: form.once_pay_count,
      product_id: form.product_id,
      seckill_price: form.seckill_price,
      seckill_time_ids: form.slot_ids.join(','),
      start_day: form.dates[0]!,
      status: form.status,
    };
    saving.value = true;
    editModalApi.lock();
    try {
      if (editingID.value) {
        await updateMerchantSeckillActiveApi(editingID.value, body);
      } else {
        await createMerchantSeckillActiveApi(body);
      }
      ElMessage.success(editingID.value ? '秒杀活动已更新' : '秒杀活动已创建');
      editModalApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      editModalApi.unlock();
    }
  },
});

function openCreate() {
  resetForm();
  editModalApi.setState({ title: '新建秒杀活动' }).open();
}

function openEdit(row: MerchantSeckillActive) {
  editingID.value = row.seckill_active_id;
  Object.assign(form, {
    dates: [row.start_day, row.end_day],
    name: row.name,
    once_pay_count: row.once_pay_count,
    product_id: row.product_id,
    seckill_price: row.seckill_price,
    slot_ids: row.seckill_time_ids
      .split(',')
      .map((id) => Number(id))
      .filter((id) => id > 0),
    status: row.status,
  });
  editModalApi.setState({ title: '编辑秒杀活动' }).open();
}

async function toggle(row: MerchantSeckillActive) {
  const status = row.status === 1 ? 0 : 1;
  await setMerchantSeckillStatusApi(row.seckill_active_id, status);
  ElMessage.success(status ? '活动已启用' : '活动已停用');
  gridApi.reload();
}

async function remove(row: MerchantSeckillActive) {
  try {
    await confirm({
      content: `删除秒杀活动「${row.name}」？删除后不可恢复。`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteMerchantSeckillActiveApi(row.seckill_active_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

async function loadOptions() {
  const [times, productPage] = await Promise.all([
    listMerchantSeckillTimesApi(),
    listMerchantProductsApi({ limit: 100, page: 1, status: 1 }),
  ]);
  slots.value = times.list || [];
  products.value = productPage.list || [];
}

onMounted(loadOptions);
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" @click="openCreate">新建秒杀活动</ElButton>
    </template>

    <Grid>
      <template #dates="{ row }">
        {{ row.start_day }} 至 {{ row.end_day }}
      </template>
      <template #slots="{ row }">
        {{ slotsText(row.seckill_time_ids) }}
      </template>
      <template #status="{ row }">
        <ElTag :type="statusType(row.status)">
          {{ row.status === 1 ? '已启用' : '已停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="warning" @click="toggle(row)">
          {{ row.status === 1 ? '停用' : '启用' }}
        </ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <EditModal class="w-[640px] max-w-[96vw]">
      <ElForm class="grid grid-cols-2 gap-x-4" label-width="88px">
        <ElFormItem class="col-span-2" label="活动名称" required>
          <ElInput v-model="form.name" maxlength="60" show-word-limit />
        </ElFormItem>
        <ElFormItem label="参与商品" required>
          <ElSelect
            v-model="form.product_id"
            :disabled="!!editingID"
            filterable
            class="w-full"
            placeholder="选择本店已审核商品"
          >
            <ElOption
              v-for="item in products"
              :key="item.product_id"
              :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`"
              :value="item.product_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="秒杀价" required>
          <ElInputNumber
            v-model="form.seckill_price"
            :min="0.01"
            :precision="2"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="活动日期" required>
          <ElDatePicker
            v-model="form.dates"
            class="w-full"
            end-placeholder="结束日期"
            start-placeholder="开始日期"
            type="daterange"
            value-format="YYYY-MM-DD"
          />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="参与场次">
          <ElSelect
            v-model="form.slot_ids"
            multiple
            class="w-full"
            placeholder="不选则使用默认场次"
          >
            <ElOption
              v-for="slot in slots"
              :key="slot.seckill_time_id"
              :label="slot.title"
              :value="slot.seckill_time_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="单次限购">
          <ElInputNumber
            v-model="form.once_pay_count"
            :min="1"
            :precision="0"
            class="w-full"
          />
        </ElFormItem>
        <ElFormItem label="初始状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </EditModal>
  </Page>
</template>
