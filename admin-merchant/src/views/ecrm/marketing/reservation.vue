<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getReservationConfigApi,
  listReservationProductsApi,
  saveReservationConfigApi,
  type ReservationProduct,
  type ReservationSlot,
} from '#/api/core/merchant-reservation';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const current = ref<ReservationProduct>();
const slots = ref<ReservationSlot[]>([]);
const form = reactive({
  reservation_type: 1,
  show_reservation_days: 7,
  slotText: '',
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '商品名称 / 商品 ID' },
    fieldName: 'keyword',
    label: '商品搜索',
  },
]);

const gridOptions: VxeGridProps<ReservationProduct> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'product_id', title: '商品 ID', width: 100 },
    { field: 'store_name', minWidth: 200, showOverflow: false, title: '商品名称' },
    {
      field: 'price',
      title: '售价',
      width: 100,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    { field: 'stock', title: '库存', width: 90 },
    { field: 'show_reservation_days', title: '可预约天数', width: 110 },
    {
      field: 'reservation_type',
      title: '预约类型',
      width: 100,
      formatter: ({ cellValue }) =>
        Number(cellValue) === 2 ? '按时段' : '按日期',
    },
    merchantListActionColumn({ width: 100 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const data = await listReservationProductsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'product_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [ConfigDrawer, configDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: saveConfig,
});

function slotsToText(items: ReservationSlot[]) {
  return items
    .map((item) => `${item.start_time}-${item.end_time},${item.stock}`)
    .join('\n');
}

function parseSlots(text: string) {
  return text
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      const [timePart, stockPart] = line.split(',');
      const [start_time = '', end_time = ''] = (timePart || '').split('-');
      return {
        start_time: start_time.trim(),
        end_time: end_time.trim(),
        stock: Number(stockPart) || 0,
      };
    })
    .filter((item) => item.start_time && item.end_time);
}

async function openConfig(row: ReservationProduct) {
  current.value = row;
  const result = await getReservationConfigApi(row.product_id);
  form.reservation_type =
    result.config?.reservation_type ?? row.reservation_type ?? 1;
  form.show_reservation_days =
    result.config?.show_reservation_days ?? row.show_reservation_days ?? 7;
  slots.value = result.slots || [];
  form.slotText = slotsToText(slots.value);
  configDrawerApi.setState({ title: '预约配置' }).open();
}

async function saveConfig() {
  if (!current.value) return;
  saving.value = true;
  configDrawerApi.lock();
  try {
    await saveReservationConfigApi(current.value.product_id, {
      reservation_type: form.reservation_type,
      show_reservation_days: form.show_reservation_days,
      slots: parseSlots(form.slotText),
    });
    ElMessage.success('预约配置已保存');
    configDrawerApi.close();
    gridApi.reload();
  } finally {
    saving.value = false;
    configDrawerApi.unlock();
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openConfig(row)">配置</ElButton>
      </template>
    </Grid>

    <ConfigDrawer class="w-[640px] max-w-[96vw]">
      <template v-if="current">
        <div class="mb-4 text-base font-medium">{{ current.store_name }}</div>
        <ElForm label-width="112px">
          <ElFormItem label="预约类型">
            <ElRadioGroup v-model="form.reservation_type">
              <ElRadio :value="1">按日期</ElRadio>
              <ElRadio :value="2">按时段</ElRadio>
            </ElRadioGroup>
          </ElFormItem>
          <ElFormItem label="可预约天数">
            <ElInputNumber
              v-model="form.show_reservation_days"
              :min="1"
              :max="90"
            />
          </ElFormItem>
          <ElFormItem label="时段与库存">
            <ElInput
              v-model="form.slotText"
              :rows="6"
              placeholder="每行一条：09:00-10:00,20&#10;10:00-11:00,15"
              type="textarea"
            />
          </ElFormItem>
        </ElForm>
      </template>
    </ConfigDrawer>
  </Page>
</template>
