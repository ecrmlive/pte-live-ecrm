<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformBrandApi,
  deletePlatformBrandApi,
  listPlatformBrandCategoriesApi,
  listPlatformBrandsApi,
  updatePlatformBrandApi,
  type PlatformBrand,
  type PlatformBrandCategory,
} from '#/api/core/platform-catalog';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const categories = ref<PlatformBrandCategory[]>([]);
const editing = ref<PlatformBrand>();
const form = reactive({ brand_name: '', category_id: 0, is_show: 1, sort: 0 });

const categoryOptions = computed(() => flatten(categories.value));

function flatten(
  nodes: PlatformBrandCategory[],
  prefix = '',
): Array<{ label: string; value: number }> {
  return nodes.flatMap((node) => [
    { label: `${prefix}${node.cate_name}`, value: node.brand_category_id },
    ...flatten(node.children || [], `${prefix}— `),
  ]);
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [],
      placeholder: '请选择',
    },
    fieldName: 'category_id',
    label: '品牌分类',
  },
]);

const gridOptions: VxeGridProps<PlatformBrand> = {
  columns: [
    { field: 'brand_id', title: 'ID', width: 90 },
    {
      field: 'brand_name',
      minWidth: 200,
      showOverflow: false,
      title: '品牌名称',
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '是否显示',
      width: 120,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const categoryIdRaw = formValues?.category_id;
        const [brandPage, categoryPage] = await Promise.all([
          listPlatformBrandsApi(
            categoryIdRaw ? { category_id: Number(categoryIdRaw) } : undefined,
          ),
          categories.value.length
            ? Promise.resolve({ list: categories.value })
            : listPlatformBrandCategoriesApi(),
        ]);
        categories.value = categoryPage.list || [];
        const list = brandPage.list || [];
        const total = list.length;
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: list.slice(start, start + page.pageSize),
          total,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'brand_id' },
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
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

async function syncCategoryFilterOptions() {
  if (!categories.value.length) {
    categories.value =
      (await listPlatformBrandCategoriesApi()).list || [];
  }
  gridApi.formApi?.updateSchema([
    {
      fieldName: 'category_id',
      componentProps: {
        clearable: true,
        options: categoryOptions.value,
        placeholder: '请选择',
      },
    },
  ]);
}

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    brand_name: '',
    category_id: 0,
    is_show: 1,
    sort: 0,
  });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加品牌' }).open();
}

function openEdit(row: PlatformBrand) {
  editing.value = row;
  Object.assign(form, {
    brand_name: row.brand_name,
    category_id: row.category_id || 0,
    is_show: row.is_show,
    sort: row.sort,
  });
  formDrawerApi.setState({ title: '编辑品牌' }).open();
}

async function save() {
  if (!form.brand_name.trim()) {
    ElMessage.warning('请填写品牌名称');
    return;
  }
  formDrawerApi.lock();
  try {
    const body = {
      brand_name: form.brand_name.trim(),
      category_id: form.category_id || 0,
      is_show: form.is_show,
      sort: form.sort,
    };
    if (editing.value) {
      await updatePlatformBrandApi(editing.value.brand_id, body);
    } else {
      await createPlatformBrandApi(body);
    }
    formDrawerApi.close();
    ElMessage.success('品牌已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: PlatformBrand, enabled: boolean) {
  const before = row.is_show === 1;
  row.is_show = enabled ? 1 : 0;
  try {
    await updatePlatformBrandApi(row.brand_id, {
      brand_name: row.brand_name,
      category_id: row.category_id || 0,
      is_show: enabled ? 1 : 0,
      sort: row.sort,
    });
  } catch {
    row.is_show = before ? 1 : 0;
  }
}

async function remove(row: PlatformBrand) {
  try {
    await ElMessageBox.confirm(
      `删除品牌“${row.brand_name}”后不可恢复，是否继续？`,
      '删除品牌',
      { type: 'warning' },
    );
    await deletePlatformBrandApi(row.brand_id);
    ElMessage.success('品牌已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}

onMounted(() => void syncCategoryFilterOptions());
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          添加品牌
        </ElButton>
      </template>
      <template #is_show="{ row }">
        <ElSwitch
          :model-value="row.is_show === 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeShow(row, Boolean(enabled))
          "
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="88px">
        <ElFormItem label="品牌名称" required>
          <ElInput v-model="form.brand_name" />
        </ElFormItem>
        <ElFormItem label="品牌分类">
          <ElSelect
            v-model="form.category_id"
            clearable
            class="w-full"
            placeholder="请选择"
          >
            <ElOption label="未分类" :value="0" />
            <ElOption
              v-for="item in categoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="是否显示">
          <ElSwitch
            v-model="form.is_show"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="显示"
            inactive-text="隐藏"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
