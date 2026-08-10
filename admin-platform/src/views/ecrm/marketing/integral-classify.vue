<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformIntegralCategoryApi,
  deletePlatformIntegralCategoryApi,
  listPlatformIntegralCategoriesApi,
  updatePlatformIntegralCategoryApi,
  updatePlatformIntegralCategoryStatusApi,
  type PlatformIntegralCategoryRow,
} from '#/api/core/platform-integral-category';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

const editing = ref<PlatformIntegralCategoryRow>();
const form = reactive({
  cate_name: '',
  is_show: 1,
  sort: 0,
});

const gridOptions: VxeGridProps<PlatformIntegralCategoryRow> = {
  columns: [
    {
      field: 'cate_name',
      formatter: ({ row }) => formatCateName(row),
      minWidth: 220,
      showOverflow: false,
      title: '分类名称',
    },
    { field: 'sort', title: '排序', width: 100 },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '是否显示',
      width: 120,
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 180,
      showOverflow: false,
      title: '创建时间',
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const list = (await listPlatformIntegralCategoriesApi()).list || [];
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: list.slice(start, start + page.pageSize),
          total: list.length,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'store_category_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function formatCateName(row: PlatformIntegralCategoryRow) {
  const name = String(row.cate_name || '').trim() || '—';
  const id = Number(row.store_category_id || 0);
  // CRMEB：分类名称后附 ID，如 ces  [ 72 ]
  return id > 0 ? `${name}  [ ${id}  ]` : name;
}

function resetForm() {
  editing.value = undefined;
  Object.assign(form, { cate_name: '', is_show: 1, sort: 0 });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加积分商品分类' }).open();
}

function openEdit(row: PlatformIntegralCategoryRow) {
  editing.value = row;
  Object.assign(form, {
    cate_name: row.cate_name,
    is_show: row.is_show === 1 ? 1 : 0,
    sort: Number(row.sort || 0),
  });
  formDrawerApi.setState({ title: '编辑积分商品分类' }).open();
}

async function save() {
  const name = form.cate_name.trim();
  if (!name) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  const sort = Math.max(0, Math.min(99999, Number(form.sort) || 0));
  formDrawerApi.lock();
  try {
    const payload = {
      cate_name: name,
      is_show: form.is_show === 1 ? 1 : 0,
      sort,
    };
    if (editing.value) {
      await updatePlatformIntegralCategoryApi(
        editing.value.store_category_id,
        payload,
      );
    } else {
      await createPlatformIntegralCategoryApi(payload);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: PlatformIntegralCategoryRow, enabled: boolean) {
  const before = row.is_show;
  row.is_show = enabled ? 1 : 0;
  try {
    await updatePlatformIntegralCategoryStatusApi(
      row.store_category_id,
      enabled ? 1 : 0,
    );
  } catch {
    row.is_show = before;
  }
}

async function remove(row: PlatformIntegralCategoryRow) {
  try {
    const tip =
      row.has_product === 1
        ? '该分类下有商品，删除后不可恢复，请确认是否删除'
        : '确定删除该分类吗';
    await confirm({
      content: tip,
      icon: 'warning',
      title: '提示',
    });
    await deletePlatformIntegralCategoryApi(row.store_category_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          添加积分商品分类
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
      <ElForm label-width="100px">
        <ElFormItem label="分类名称" required>
          <ElInput
            v-model="form.cate_name"
            maxlength="100"
            show-word-limit
            placeholder="请输入分类名称"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber
            v-model="form.sort"
            :min="0"
            :max="99999"
            :precision="0"
            class="w-full"
          />
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
