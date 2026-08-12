<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createServiceMeal,
  deleteServiceMeal,
  fetchServiceMeals,
  updateServiceMeal,
  updateServiceMealStatus,
  type ServiceMealInput,
  type ServiceMealRow,
} from '#/api/core/platform-service-meal';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

type ServiceMealForm = ServiceMealInput;

const editing = ref<ServiceMealRow>();
const form = reactive<ServiceMealForm>({
  name: '',
  type: 1,
  price: 0,
  num: 0,
  sort: 0,
  status: 1,
});

const typeLabel = (type: number) => (type === 2 ? '电子面单' : '复制商品');
const drawerTitle = computed(() => (editing.value ? '编辑套餐' : '新增套餐'));

const gridOptions: VxeGridProps<ServiceMealRow> = {
  columns: [
    { title: '序号', type: 'seq' as const, width: 86 },
    { field: 'name', minWidth: 240, showOverflow: 'tooltip', title: '业务名称' },
    {
      field: 'type',
      formatter: ({ cellValue }) => typeLabel(Number(cellValue)),
      title: '类型',
      width: 150,
    },
    {
      field: 'price',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      title: '价格(元)',
      width: 130,
    },
    { field: 'num', title: '购买数量(次数)', width: 150 },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '添加时间',
    },
    { field: 'status', slots: { default: 'status' }, title: '是否显示', width: 120 },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const result = await fetchServiceMeals({
          limit: page.pageSize,
          page: page.currentPage,
        });
        return { items: result.list, total: result.total };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'meal_id' },
  toolbarConfig: { custom: false, export: false, refresh: false, zoom: false },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });
const [MealDrawer, mealDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[720px] max-w-[96vw]',
  confirmText: '保存',
  onConfirm: async () => save(),
  placement: 'right',
});

function resetForm() {
  Object.assign(form, {
    name: '',
    type: 1,
    price: 0,
    num: 0,
    sort: 0,
    status: 1,
  });
}

function openCreate() {
  editing.value = undefined;
  resetForm();
  mealDrawerApi.setState({ title: '新增套餐' }).open();
}

function openEdit(row: ServiceMealRow) {
  editing.value = row;
  Object.assign(form, {
    name: row.name,
    type: row.type,
    price: Number(row.price),
    num: Number(row.num),
    sort: Number(row.sort),
    status: row.status,
  });
  mealDrawerApi.setState({ title: '编辑套餐' }).open();
}

async function save() {
  const payload: ServiceMealInput = {
    ...form,
    name: form.name.trim(),
  };
  if (!payload.name) {
    ElMessage.warning('请填写套餐名称');
    return;
  }
  mealDrawerApi.lock();
  try {
    if (editing.value) {
      await updateServiceMeal(editing.value.meal_id, payload);
    } else {
      await createServiceMeal(payload);
    }
    ElMessage.success(editing.value ? '套餐已保存' : '套餐已新增');
    mealDrawerApi.close();
    gridApi.reload();
  } finally {
    mealDrawerApi.unlock();
  }
}

async function toggleStatus(row: ServiceMealRow, status: 0 | 1) {
  try {
    await updateServiceMealStatus(row.meal_id, status);
    ElMessage.success('显示状态已保存');
    gridApi.reload();
  } catch {
    gridApi.reload();
  }
}

async function remove(row: ServiceMealRow) {
  try {
    await confirm({
      content: `确定删除套餐“${row.name}”吗？`,
      icon: 'warning',
      title: '删除套餐',
    });
    await deleteServiceMeal(row.meal_id);
    ElMessage.success('套餐已删除');
    gridApi.reload();
  } catch {
    // 用户取消时不改变当前列表。
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增套餐
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :active-value="1"
          :inactive-value="0"
          :model-value="row.status"
          @update:model-value="toggleStatus(row, Number($event) as 0 | 1)"
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <MealDrawer :title="drawerTitle">
      <ElForm label-width="108px">
        <ElFormItem label="套餐名称" required>
          <ElInput v-model="form.name" :maxlength="30" placeholder="请输入套餐名称" />
        </ElFormItem>
        <ElFormItem label="套餐类型" required>
          <ElRadioGroup v-model="form.type">
            <ElRadio :value="1">一号通商品采集</ElRadio>
            <ElRadio :value="2">一号通电子面单</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="价格" required>
          <ElInputNumber v-model="form.price" :min="0" :precision="2" class="w-full" />
        </ElFormItem>
        <ElFormItem label="数量" required>
          <ElInputNumber v-model="form.num" :min="0" :precision="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio :value="1">开启</ElRadio>
            <ElRadio :value="0">关闭</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="排序" required>
          <ElInputNumber v-model="form.sort" :precision="0" class="w-full" />
        </ElFormItem>
      </ElForm>
    </MealDrawer>
  </Page>
</template>
