<script setup lang="ts">
import type { VxeGridProps } from "#/adapter/vxe-table";

import { onMounted, reactive, ref } from "vue";

import { Page, confirm, useVbenDrawer } from "@vben/common-ui";
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
} from "element-plus";
import { Plus } from "@element-plus/icons-vue";

import { useVbenVxeGrid } from "#/adapter/vxe-table";
import { getAccessCodesApi, getUserInfoApi } from "#/api/core/auth";
import {
  createConfigItemApi,
  deleteConfigItemApi,
  listConfigItemsApi,
  setConfigItemStatusApi,
  updateConfigItemApi,
  type ConfigItem,
} from "#/api/core/platform-config-item";
import { platformListActionColumn, platformListPagerConfig } from "#/constants/platform-list-grid";
import { formatShanghaiDateTime } from "#/utils/date-time";

const HOT_SEARCH_TYPE = "hot_search" as const;
const HOT_SEARCH_READ_CODE = "setting.shop.hot.read";
const HOT_SEARCH_MANAGE_CODE = "setting.shop.hot.manage";

const canRead = ref(false);
const canManage = ref(false);
const editing = ref<ConfigItem>();
const form = reactive({ keyword: "", sort: 0, status: 1 });

const gridOptions: VxeGridProps<ConfigItem> = {
  columns: [
    { field: "id", title: "ID", width: 120 },
    { field: "name", minWidth: 220, title: "关键词" },
    {
      field: "created_at",
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || "—",
      minWidth: 190,
      title: "添加时间",
    },
    {
      field: "status",
      slots: { default: "status" },
      title: "是否显示",
      width: 150,
    },
    platformListActionColumn({ width: 180 }),
  ],
  emptyText: "暂无热门搜索",
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const result = await listConfigItemsApi(HOT_SEARCH_TYPE, {
          limit: page.pageSize,
          page: page.currentPage,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: "id" },
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
  cancelText: "取消",
  class: "w-[1000px] max-w-[96vw]",
  confirmText: "保存",
  placement: "right",
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, { keyword: "", sort: 0, status: 1 });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: "新增热门搜索" }).open();
}

function openEdit(row: ConfigItem) {
  editing.value = row;
  Object.assign(form, {
    keyword: row.name,
    sort: row.sort,
    status: row.status,
  });
  formDrawerApi.setState({ title: "编辑热门搜索" }).open();
}

async function save() {
  const keyword = form.keyword.trim();
  if (!keyword) {
    ElMessage.warning("请填写关键词");
    return;
  }
  formDrawerApi.lock();
  try {
    const payload = { name: keyword, sort: form.sort, status: form.status };
    if (editing.value) {
      await updateConfigItemApi(HOT_SEARCH_TYPE, editing.value.id, payload);
      ElMessage.success("热门搜索已更新");
    } else {
      await createConfigItemApi(HOT_SEARCH_TYPE, payload);
      ElMessage.success("热门搜索已创建");
    }
    formDrawerApi.close();
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function toggleStatus(row: ConfigItem, status: number) {
  try {
    await setConfigItemStatusApi(HOT_SEARCH_TYPE, row.id, status === 1 ? 1 : 0);
    ElMessage.success("显示状态已更新");
    gridApi.reload();
  } catch {
    gridApi.reload();
  }
}

async function remove(row: ConfigItem) {
  try {
    await confirm({
      content: `删除热门搜索“${row.name}”后不可恢复，是否继续？`,
      icon: "warning",
      title: "提示",
    });
    await deleteConfigItemApi(HOT_SEARCH_TYPE, row.id);
    ElMessage.success("热门搜索已删除");
    gridApi.reload();
  } catch {
    // 用户取消或统一请求层已提示错误。
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  const isPlatform = profile.roles.includes("platform");
  canRead.value =
    isPlatform && (codes.includes(HOT_SEARCH_READ_CODE) || codes.includes(HOT_SEARCH_MANAGE_CODE));
  canManage.value = isPlatform && codes.includes(HOT_SEARCH_MANAGE_CODE);
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton v-if="canManage" :icon="Plus" type="primary" @click="openCreate">
          新增热门搜索
        </ElButton>
      </template>

      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          aria-label="是否显示"
          @change="(visible) => toggleStatus(row, visible ? 1 : 0)"
        />
      </template>

      <template #action="{ row }">
        <ElButton v-if="canManage" link type="primary" @click="openEdit(row)"> 编辑 </ElButton>
        <ElButton v-if="canManage" link type="danger" @click="remove(row)"> 删除 </ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="100px">
        <ElFormItem label="关键词" required>
          <ElInput v-model="form.keyword" maxlength="128" placeholder="请输入关键词" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :step="1" />
        </ElFormItem>
        <ElFormItem label="是否显示">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
