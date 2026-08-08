<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createMerchantCategory,
  deleteMerchantCategory,
  fetchMerchantCategories,
  updateMerchantCategory,
  type MerchantCategoryRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const canManage = ref(false);
const editing = ref<MerchantCategoryRow>();
const form = reactive({ category_name: '', commission_rate: 0 });

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '分类名称' },
    fieldName: 'keyword',
    label: '关键词',
  },
]);

const gridOptions: VxeGridProps<MerchantCategoryRow> = {
  columns: [
    { field: 'merchant_category_id', title: 'ID', width: 90 },
    {
      field: 'category_name',
      minWidth: 240,
      showOverflow: false,
      title: '分类名称',
    },
    {
      field: 'commission_rate',
      formatter: ({ cellValue }) => `${Number(cellValue || 0).toFixed(2)}%`,
      title: '手续费',
      width: 120,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 170,
      showOverflow: false,
      title: '创建时间',
    },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const keyword = String(formValues?.keyword ?? '').trim().toLowerCase();
        let list = (await fetchMerchantCategories()).list || [];
        if (keyword) {
          list = list.filter((row) =>
            row.category_name.toLowerCase().includes(keyword),
          );
        }
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: list.slice(start, start + page.pageSize),
          total: list.length,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'merchant_category_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [CategoryDrawer, categoryDrawerApi] = useVbenDrawer({
  class: 'w-[560px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, { category_name: '', commission_rate: 0 });
}

function openCreate() {
  resetForm();
  categoryDrawerApi.setState({ title: '新增店铺分类' }).open();
}

function openEdit(row: MerchantCategoryRow) {
  editing.value = row;
  Object.assign(form, {
    category_name: row.category_name,
    commission_rate: Number(row.commission_rate),
  });
  categoryDrawerApi.setState({ title: '编辑店铺分类' }).open();
}

async function save() {
  if (
    !form.category_name.trim() ||
    form.commission_rate < 0 ||
    form.commission_rate > 100
  ) {
    ElMessage.warning('请填写分类名称，并将手续费比例限制在 0% 至 100%');
    return;
  }
  categoryDrawerApi.lock();
  try {
    const input = {
      category_name: form.category_name.trim(),
      commission_rate: form.commission_rate,
    };
    if (editing.value) {
      await updateMerchantCategory(
        editing.value.merchant_category_id,
        input,
      );
    } else {
      await createMerchantCategory(input);
    }
    categoryDrawerApi.close();
    ElMessage.success('店铺分类已保存');
    gridApi.reload();
  } finally {
    categoryDrawerApi.unlock();
  }
}

async function remove(row: MerchantCategoryRow) {
  try {
    await confirm({
      content: `删除店铺分类“${row.category_name}”后不可恢复，是否继续？`,
      icon: 'warning',
      title: '删除店铺分类',
    });
    await deleteMerchantCategory(row.merchant_category_id);
    ElMessage.success('店铺分类已删除');
    gridApi.reload();
  } catch {
    /* 用户取消或统一请求层处理 */
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('merchant.category.manage');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="categories-toolbar">
          <ElAlert
            class="categories-tip"
            type="warning"
            show-icon
            :closable="false"
            :title="
              '用于区分店铺所属行业类型，后台可按服装、百货等分类设置不同的手续费比例，例如“美妆护肤、餐饮美食”等'
            "
          />
          <div class="categories-toolbar__actions">
            <ElButton
              v-if="canManage"
              :icon="Plus"
              type="primary"
              @click="openCreate"
            >
              新增店铺分类
            </ElButton>
          </div>
        </div>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
        <span v-else>—</span>
      </template>
    </Grid>

    <CategoryDrawer>
      <ElForm label-width="108px">
        <ElFormItem label="分类名称" required>
          <ElInput
            v-model="form.category_name"
            maxlength="128"
            placeholder="请输入分类名称"
          />
        </ElFormItem>
        <ElFormItem label="手续费比例" required>
          <div class="flex w-full items-center gap-2">
            <ElInputNumber
              v-model="form.commission_rate"
              :max="100"
              :min="0"
              :precision="2"
              class="flex-1"
            />
            <span>%</span>
          </div>
        </ElFormItem>
      </ElForm>
    </CategoryDrawer>
  </Page>
</template>

<style scoped>
.categories-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
  min-width: 0;
}

.categories-tip {
  width: 100%;
}

.categories-toolbar__actions {
  display: flex;
  justify-content: flex-start;
}
</style>
