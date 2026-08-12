<script lang="ts" setup>
import type { VbenFormProps } from "#/adapter/form";
import type { VxeGridProps } from "#/adapter/vxe-table";

import { reactive, ref } from "vue";

import { confirm, Page, useVbenDrawer } from "@vben/common-ui";
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from "element-plus";

import { useVbenVxeGrid } from "#/adapter/vxe-table";
import {
  deleteExpress,
  fetchExpressList,
  syncExpressCatalog,
  updateExpress,
  type ExpressRow,
} from "#/api/core/ecrm";
import { platformListActionColumn, platformListPagerConfig } from "#/constants/platform-list-grid";
import { listFormOptionsDefaults } from "#/utils/list-form-defaults";

const form = reactive({
  code: "",
  isShow: true,
  name: "",
  sort: 0,
});
const editingId = ref(0);

const [ExpressDrawer, expressDrawerApi] = useVbenDrawer({
  class: "w-[1000px] max-w-[96vw]",
  confirmText: "保存",
  cancelText: "取消",
  placement: "right",
  onConfirm: saveExpress,
});

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    {
      component: "Input",
      componentProps: {
        clearable: true,
        placeholder: "请输入物流公司名称或者编码",
      },
      fieldName: "keyword",
      label: "搜索",
    },
  ],
  { wrapperClass: "grid-cols-1 md:grid-cols-[520px_auto]" },
);

const gridOptions: VxeGridProps<ExpressRow> = {
  columns: [
    { field: "express_id", title: "ID", width: 88 },
    { field: "name", minWidth: 260, title: "物流公司名称" },
    { field: "code", minWidth: 220, title: "编码" },
    { field: "sort", sortable: true, title: "排序", width: 120 },
    {
      align: "center",
      field: "is_show",
      slots: { default: "isShow" },
      title: "是否显示",
      width: 130,
    },
    platformListActionColumn({ width: 150 }),
  ],
  formOptions,
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page, sorts }, formValues) => {
        const sort = Array.isArray(sorts) ? sorts.find((item) => item.field === "sort") : undefined;
        const data = await fetchExpressList({
          keyword: String(formValues?.keyword ?? "").trim() || undefined,
          limit: page.pageSize,
          page: page.currentPage,
          sort_order: sort?.order === "asc" || sort?.order === "desc" ? sort.order : undefined,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
    sort: true,
  },
  rowConfig: { isHover: true, keyField: "express_id" },
  sortConfig: { multiple: false, remote: true, trigger: "cell" },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function openEdit(row: ExpressRow) {
  editingId.value = row.express_id;
  form.name = row.name;
  form.code = row.code;
  form.sort = row.sort;
  form.isShow = row.is_show === 1;
  expressDrawerApi.setState({ title: "编辑物流公司" });
  expressDrawerApi.open();
}

async function saveExpress() {
  const name = form.name.trim();
  const code = form.code.trim();
  if (!name || !code) {
    ElMessage.warning("请填写物流公司名称和编码");
    return;
  }
  await updateExpress(editingId.value, {
    code,
    is_show: form.isShow ? 1 : 0,
    name,
    sort: form.sort,
  });
  ElMessage.success("已保存");
  expressDrawerApi.close();
  await gridApi.reload();
}

async function toggleShow(row: ExpressRow, isShow: boolean) {
  try {
    await updateExpress(row.express_id, {
      code: row.code,
      is_show: isShow ? 1 : 0,
      name: row.name,
      sort: row.sort,
    });
    ElMessage.success("已保存");
    await gridApi.reload();
  } catch {
    await gridApi.reload();
  }
}

async function onSync() {
  try {
    await confirm({
      content: "将同步平台内置物流公司目录，已有物流公司的排序和显示状态不会被覆盖。",
      icon: "warning",
      title: "同步物流公司",
    });
  } catch {
    return;
  }
  const result = await syncExpressCatalog();
  ElMessage.success(`同步完成：新增 ${result.created} 条，更新 ${result.updated} 条`);
  await gridApi.reload();
}

async function onDelete(row: ExpressRow) {
  try {
    await confirm({
      content: `删除物流公司“${row.name}”后将不再用于后续发货选择。`,
      icon: "warning",
      title: "删除物流公司",
    });
  } catch {
    return;
  }
  await deleteExpress(row.express_id);
  ElMessage.success("已删除");
  await gridApi.reload();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton type="primary" @click="onSync">同步物流公司</ElButton>
      </template>
      <template #isShow="{ row }">
        <ElSwitch :model-value="row.is_show === 1" @update:model-value="toggleShow(row, $event)" />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="onDelete(row)">删除</ElButton>
      </template>
    </Grid>

    <ExpressDrawer>
      <ElForm label-width="120px">
        <ElFormItem label="物流公司名称" required>
          <ElInput v-model="form.name" placeholder="请输入物流公司名称" />
        </ElFormItem>
        <ElFormItem label="编码" required>
          <ElInput v-model="form.code" placeholder="请输入物流公司编码" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="是否显示">
          <ElSwitch v-model="form.isShow" />
        </ElFormItem>
      </ElForm>
    </ExpressDrawer>
  </Page>
</template>
