<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
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
import { getAccessCodesApi } from '#/api/core/auth';
import {
  createMerchantDiscountApi,
  deleteMerchantDiscountApi,
  listMerchantDiscountsApi,
  setMerchantDiscountStatusApi,
  updateMerchantDiscountApi,
  type MerchantDiscount,
  type MerchantDiscountStatus,
} from '#/api/core/merchant-discount';
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
const canManage = ref(false);
const form = reactive({
  ends_at: '',
  free_shipping: false,
  name: '',
  package_price: 0,
  product_ids_text: '',
  remark: '',
  starts_at: '',
  status: 'draft' as MerchantDiscountStatus,
});

const statusLabels: Record<MerchantDiscountStatus, string> = {
  active: '进行中',
  closed: '已关闭',
  draft: '草稿',
  pending: '待审核',
  rejected: '已拒绝',
};

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '套餐名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(statusLabels).map(([value, label]) => ({
        label,
        value,
      })),
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<MerchantDiscount> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'activity_id', title: 'ID', width: 90 },
    { field: 'name', minWidth: 160, showOverflow: false, title: '名称' },
    {
      field: 'package_price',
      title: '套餐价',
      width: 110,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'product_ids',
      title: '商品数',
      width: 90,
      formatter: ({ cellValue }) => (cellValue || []).length,
    },
    {
      field: 'free_shipping',
      title: '包邮',
      width: 80,
      formatter: ({ cellValue }) => (cellValue ? '是' : '否'),
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    {
      field: 'starts_at',
      minWidth: 220,
      showOverflow: false,
      slots: { default: 'period' },
      title: '有效期',
    },
    merchantListActionColumn({ width: 200 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const data = await listMerchantDiscountsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status: formValues?.status || undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'activity_id' },
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
  onConfirm: async () => {
    const productIDs = parseProductIDs();
    if (
      !form.name.trim() ||
      form.package_price <= 0 ||
      productIDs.length === 0
    ) {
      ElMessage.warning('请填写名称、正数套餐价，并至少填写一个商品 ID');
      return;
    }
    saving.value = true;
    editDrawerApi.lock();
    try {
      const body = {
        ends_at: form.ends_at || undefined,
        free_shipping: form.free_shipping,
        name: form.name.trim(),
        package_price: form.package_price,
        product_ids: productIDs,
        remark: form.remark.trim(),
        starts_at: form.starts_at || undefined,
        status: form.status,
      };
      if (editingID.value) {
        await updateMerchantDiscountApi(editingID.value, body);
        ElMessage.success('优惠套餐已更新');
      } else {
        await createMerchantDiscountApi(body);
        ElMessage.success('优惠套餐已创建');
      }
      editDrawerApi.close();
      gridApi.reload();
    } finally {
      saving.value = false;
      editDrawerApi.unlock();
    }
  },
});

function parseProductIDs() {
  return form.product_ids_text
    .split(/[,，\s]+/)
    .map((part) => Number(part.trim()))
    .filter((id) => Number.isFinite(id) && id > 0);
}

function openCreate() {
  editingID.value = undefined;
  Object.assign(form, {
    ends_at: '',
    free_shipping: false,
    name: '',
    package_price: 0,
    product_ids_text: '',
    remark: '',
    starts_at: '',
    status: 'draft',
  });
  editDrawerApi.setState({ title: '新增优惠套餐' }).open();
}

function openEdit(row: MerchantDiscount) {
  editingID.value = row.activity_id;
  Object.assign(form, {
    ends_at: row.ends_at,
    free_shipping: row.free_shipping,
    name: row.name,
    package_price: row.package_price,
    product_ids_text: (row.product_ids || []).join(','),
    remark: row.remark,
    starts_at: row.starts_at,
    status: row.status,
  });
  editDrawerApi.setState({ title: '编辑优惠套餐' }).open();
}

async function toggleStatus(row: MerchantDiscount) {
  const next: MerchantDiscountStatus =
    row.status === 'active' ? 'closed' : 'active';
  const action = next === 'active' ? '上架' : '关闭';
  try {
    await confirm({
      content: `确认${action}「${row.name}」？`,
      icon: 'warning',
      title: `${action}确认`,
    });
    await setMerchantDiscountStatusApi(row.activity_id, next);
    ElMessage.success(`已${action}`);
    gridApi.reload();
  } catch {
    // cancelled
  }
}

async function remove(row: MerchantDiscount) {
  try {
    await confirm({
      content: `删除「${row.name}」后不可恢复，确认继续？`,
      icon: 'warning',
      title: '删除优惠套餐',
    });
    await deleteMerchantDiscountApi(row.activity_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('marketing.discounts.manage');
});
</script>

<template>
  <Page auto-content-height>
    <template v-if="canManage" #extra>
      <ElButton type="primary" @click="openCreate">新增套餐</ElButton>
    </template>

    <Grid>
      <template #status="{ row }">
        <ElTag>{{ statusLabels[row.status] || row.status }}</ElTag>
      </template>
      <template #period="{ row }">
        {{ row.starts_at || '—' }} ~ {{ row.ends_at || '—' }}
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="warning" @click="toggleStatus(row)">
            {{ row.status === 'active' ? '关闭' : '上架' }}
          </ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else class="text-muted-foreground">—</span>
      </template>
    </Grid>

    <EditDrawer class="w-[560px] max-w-[96vw]">
      <ElForm label-width="96px">
        <ElFormItem label="名称" required>
          <ElInput v-model="form.name" maxlength="128" show-word-limit />
        </ElFormItem>
        <ElFormItem label="套餐价" required>
          <ElInputNumber
            v-model="form.package_price"
            :min="0.01"
            :precision="2"
            :step="1"
          />
        </ElFormItem>
        <ElFormItem label="商品 ID" required>
          <ElInput
            v-model="form.product_ids_text"
            placeholder="多个 ID 用逗号分隔，如 1001,1006"
          />
        </ElFormItem>
        <ElFormItem label="状态">
          <ElSelect v-model="form.status" class="w-40">
            <ElOption
              v-for="(label, status) in statusLabels"
              :key="status"
              :label="label"
              :value="status"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="开始时间">
          <ElInput v-model="form.starts_at" placeholder="YYYY-MM-DD HH:mm:ss" />
        </ElFormItem>
        <ElFormItem label="结束时间">
          <ElInput v-model="form.ends_at" placeholder="YYYY-MM-DD HH:mm:ss" />
        </ElFormItem>
        <ElFormItem label="包邮">
          <ElSwitch v-model="form.free_shipping" />
        </ElFormItem>
        <ElFormItem label="备注">
          <ElInput v-model="form.remark" :rows="2" maxlength="255" type="textarea" />
        </ElFormItem>
      </ElForm>
    </EditDrawer>
  </Page>
</template>
