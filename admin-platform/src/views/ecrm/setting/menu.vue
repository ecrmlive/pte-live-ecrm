<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
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
  fetchPlatformMenus,
  updatePlatformMenu,
  type PlatformMenuRow,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { LIST_KEYWORD_FIELD, listFormOptionsDefaults } from '#/utils/list-form-defaults';

const allRows = ref<PlatformMenuRow[]>([]);
const editing = ref<PlatformMenuRow>();
const form = reactive({ menu_name: '', sort: 0, is_show: 1 });

function typeText(value: number) {
  return value === 2 ? '按钮权限' : '菜单页面';
}

async function loadMenus() {
  allRows.value = await fetchPlatformMenus();
  return allRows.value;
}

async function saveMenu(
  id: number,
  data: { menu_name?: string; sort?: number; is_show?: number },
) {
  return updatePlatformMenu(id, data);
}

/** 筛选时保留命中节点及其祖先，便于树形仍能展开到目标项 */
function filterMenuRows(
  rows: PlatformMenuRow[],
  keyword: string,
  typeRaw: unknown,
): PlatformMenuRow[] {
  let matched = rows;
  if (keyword) {
    matched = matched.filter((row) =>
      `${row.menu_name}${row.path}${row.route}`
        .toLowerCase()
        .includes(keyword),
    );
  }
  if (typeRaw === 1 || typeRaw === 2) {
    matched = matched.filter((row) => row.is_menu === Number(typeRaw));
  }
  if (matched.length === 0 || matched.length === rows.length) {
    return matched.length === 0 ? [] : rows;
  }
  const byId = new Map(rows.map((row) => [row.menu_id, row]));
  const keep = new Set<number>();
  for (const row of matched) {
    let cur: PlatformMenuRow | undefined = row;
    while (cur) {
      if (keep.has(cur.menu_id)) break;
      keep.add(cur.menu_id);
      cur = cur.pid ? byId.get(cur.pid) : undefined;
    }
  }
  return rows.filter((row) => keep.has(row.menu_id));
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('名称 / 路径 / 路由'),
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '菜单页面', value: 1 },
        { label: '按钮权限', value: 2 },
      ],
      placeholder: '全部类型',
    },
    fieldName: 'type',
    label: '类型',
  },
]);

const gridOptions: VxeGridProps<PlatformMenuRow> = {
  columns: [
    { field: 'menu_id', title: 'ID', width: 76 },
    {
      field: 'menu_name',
      minWidth: 220,
      title: '名称',
      treeNode: true,
    },
    {
      field: 'path',
      minWidth: 210,
      showOverflow: false,
      title: '路径',
    },
    {
      field: 'route',
      minWidth: 150,
      showOverflow: false,
      title: '路由',
    },
    {
      field: 'is_menu',
      slots: { default: 'type' },
      title: '类型',
      width: 100,
    },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'is_show',
      slots: { default: 'is_show' },
      title: '显示',
      width: 80,
    },
    platformListActionColumn({ width: 90 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, formValues) => {
        await loadMenus();
        const list = filterMenuRows(
          allRows.value,
          String(formValues?.keyword ?? '')
            .trim()
            .toLowerCase(),
          formValues?.type,
        );
        return { items: list, total: list.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'menu_id' },
  treeConfig: {
    expandAll: false,
    indent: 20,
    parentField: 'pid',
    rowField: 'menu_id',
    transform: true,
  },
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
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function openEdit(row: PlatformMenuRow) {
  editing.value = row;
  Object.assign(form, {
    menu_name: row.menu_name,
    sort: row.sort,
    is_show: row.is_show,
  });
  formDrawerApi.setState({ title: '编辑菜单' }).open();
}

async function save() {
  if (!editing.value || !form.menu_name.trim()) {
    ElMessage.warning('菜单名称不能为空');
    return;
  }
  formDrawerApi.lock();
  try {
    await saveMenu(editing.value.menu_id, {
      menu_name: form.menu_name.trim(),
      sort: form.sort,
      is_show: form.is_show,
    });
    formDrawerApi.close();
    ElMessage.success('菜单已保存');
    await loadMenus();
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

onMounted(() => {
  /* grid loads on mount */
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #type="{ row }">
        <ElTag :type="row.is_menu === 2 ? 'warning' : 'success'">
          {{ typeText(row.is_menu) }}
        </ElTag>
      </template>
      <template #is_show="{ row }">
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '显示' : '隐藏' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="100px">
        <ElFormItem label="菜单 ID">
          <ElInput :model-value="String(editing?.menu_id || '')" disabled />
        </ElFormItem>
        <ElFormItem label="菜单名称">
          <ElInput v-model="form.menu_name" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" />
        </ElFormItem>
        <ElFormItem label="显示状态">
          <ElSwitch v-model="form.is_show" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
