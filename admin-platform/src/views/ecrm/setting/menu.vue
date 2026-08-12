<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, nextTick, reactive, ref } from 'vue';

import { Plus } from '@element-plus/icons-vue';
import { confirm, Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElSwitch,
  ElTabPane,
  ElTabs,
  ElTreeSelect,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createPlatformMenu,
  deletePlatformMenu,
  fetchPlatformMenus,
  updatePlatformMenu,
  type PlatformMenuKind,
  type PlatformMenuRow,
  type PlatformMenuScope,
} from '#/api/core/ecrm';
import MenuIconPreview from '#/components/platform-menu/MenuIconPreview.vue';
import { PLATFORM_MENU_PICKER_ICONS } from '#/constants/platform-lucide-icons';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

type FormMode = 'create' | 'edit';

interface MenuTreeOption {
  children?: MenuTreeOption[];
  label: string;
  value: number;
}

const activeScope = ref<PlatformMenuScope>('platform');
const editing = ref<PlatformMenuRow>();
const formMode = ref<FormMode>('create');
const menuRows = ref<PlatformMenuRow[]>([]);
const treeExpanded = ref(false);
const form = reactive({
  code: '',
  icon: '',
  is_show: 1,
  kind: 'page' as PlatformMenuKind,
  menu_name: '',
  parent_id: 0,
  path: '',
  sort: 0,
});

const drawerTitle = computed(() =>
  formMode.value === 'create' ? '新增菜单' : '编辑菜单',
);

const parentOptions = computed<MenuTreeOption[]>(() => {
  const blocked = new Set<number>();
  if (editing.value) {
    blocked.add(editing.value.menu_id);
    let changed = true;
    while (changed) {
      changed = false;
      for (const row of menuRows.value) {
        if (blocked.has(row.pid) && !blocked.has(row.menu_id)) {
          blocked.add(row.menu_id);
          changed = true;
        }
      }
    }
  }
  const nodes = new Map<number, MenuTreeOption>();
  for (const row of menuRows.value) {
    if (!blocked.has(row.menu_id)) {
      nodes.set(row.menu_id, {
        label: `${row.menu_name} [${row.menu_id}]`,
        value: row.menu_id,
      });
    }
  }
  const roots: MenuTreeOption[] = [{ label: '顶级菜单', value: 0 }];
  for (const row of menuRows.value) {
    const node = nodes.get(row.menu_id);
    if (!node) continue;
    const parent = nodes.get(row.pid);
    if (row.pid && parent) {
      parent.children ||= [];
      parent.children.push(node);
    } else {
      roots.push(node);
    }
  }
  return roots;
});

const gridOptions: VxeGridProps<PlatformMenuRow> = {
  columns: [
    {
      field: 'menu_name',
      minWidth: 320,
      showOverflow: false,
      slots: { default: 'menu_name' },
      title: '菜单名称',
      treeNode: true,
    },
    { field: 'path', minWidth: 260, showOverflow: false, title: '菜单地址' },
    {
      align: 'center',
      field: 'icon',
      slots: { default: 'icon' },
      title: '菜单图标',
      width: 120,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      title: '创建时间',
      width: 180,
    },
    { field: 'sort', title: '排序', width: 100 },
    platformListActionColumn({ minWidth: 250 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async () => {
        menuRows.value = await fetchPlatformMenus(activeScope.value);
        return { items: menuRows.value, total: menuRows.value.length };
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm(parentID = 0) {
  Object.assign(form, {
    code: '',
    icon: '',
    is_show: 1,
    kind: 'page' as PlatformMenuKind,
    menu_name: '',
    parent_id: parentID,
    path: '',
    sort: 0,
  });
}

async function openCreate(parentID = 0) {
  formMode.value = 'create';
  editing.value = undefined;
  resetForm(parentID);
  await loadRows();
  formDrawerApi.setState({ title: drawerTitle.value }).open();
}

async function openEdit(row: PlatformMenuRow) {
  formMode.value = 'edit';
  editing.value = row;
  Object.assign(form, {
    code: row.route,
    icon: row.icon,
    is_show: row.is_show,
    kind: row.kind,
    menu_name: row.menu_name,
    parent_id: row.pid,
    path: row.path,
    sort: row.sort,
  });
  await loadRows();
  formDrawerApi.setState({ title: drawerTitle.value }).open();
}

async function loadRows() {
  menuRows.value = await fetchPlatformMenus(activeScope.value);
}

async function reload() {
  treeExpanded.value = false;
  await gridApi.reload();
}

function setTreeExpand(expand: boolean) {
  treeExpanded.value = expand;
  nextTick(() => {
    gridApi.grid?.setAllTreeExpand?.(expand);
  });
}

async function changeScope(scope: string | number) {
  activeScope.value = scope as PlatformMenuScope;
  await reload();
}

async function save() {
  const code = form.code.trim();
  const menuName = form.menu_name.trim();
  if (!menuName) {
    ElMessage.warning('菜单名称不能为空');
    return;
  }
  if (!code) {
    ElMessage.warning('菜单标识不能为空');
    return;
  }
  const payload = {
    code,
    icon: form.icon.trim(),
    is_show: form.is_show,
    kind: form.kind,
    menu_name: menuName,
    menu_scope: activeScope.value,
    parent_id: form.parent_id,
    path: form.path.trim(),
    sort: form.sort,
  };
  formDrawerApi.lock();
  try {
    if (formMode.value === 'edit' && editing.value) {
      await updatePlatformMenu(editing.value.menu_id, payload);
    } else {
      await createPlatformMenu(payload);
    }
    formDrawerApi.close();
    ElMessage.success('菜单已保存');
    await reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function remove(row: PlatformMenuRow) {
  try {
    await confirm({
      content: `确认删除菜单「${row.menu_name}」？存在子菜单时不可删除。`,
      icon: 'warning',
      title: '删除菜单',
    });
  } catch {
    return;
  }
  await deletePlatformMenu(row.menu_id);
  ElMessage.success('菜单已删除');
  await reload();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="w-full">
          <ElTabs
            v-model="activeScope"
            class="menu-manage-tabs"
            @tab-change="changeScope"
          >
            <ElTabPane label="平台" name="platform" />
            <ElTabPane label="商户" name="merchant" />
            <ElTabPane label="区域" name="region" />
          </ElTabs>
          <div class="flex items-center gap-3 pb-4">
            <ElButton :icon="Plus" type="primary" @click="openCreate()">
              新增菜单
            </ElButton>
            <ElButton @click="setTreeExpand(!treeExpanded)">
              展开/收起
            </ElButton>
          </div>
        </div>
      </template>
      <template #menu_name="{ row }">
        <span>{{ row.menu_name }} [{{ row.menu_id }}]</span>
      </template>
      <template #icon="{ row }">
        <MenuIconPreview :icon="row.icon" :size="18" />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openCreate(row.menu_id)">
          新增子菜单
        </ElButton>
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="110px">
        <ElFormItem label="上级菜单">
          <ElTreeSelect
            v-model="form.parent_id"
            :data="parentOptions"
            :props="{ children: 'children', label: 'label', value: 'value' }"
            check-strictly
            class="w-full"
            default-expand-all
            node-key="value"
          />
        </ElFormItem>
        <ElFormItem label="菜单名称" required>
          <ElInput v-model="form.menu_name" maxlength="64" show-word-limit />
        </ElFormItem>
        <ElFormItem label="菜单标识" required>
          <ElInput v-model="form.code" placeholder="如 setting.menu" />
        </ElFormItem>
        <ElFormItem label="菜单地址">
          <ElInput v-model="form.path" placeholder="如 /setting/menu" />
        </ElFormItem>
        <ElFormItem label="菜单图标">
          <ElSelect v-model="form.icon" clearable filterable class="w-full">
            <ElOption label="无图标" value="" />
            <ElOption
              v-for="icon in PLATFORM_MENU_PICKER_ICONS"
              :key="icon"
              :label="icon"
              :value="icon"
            >
              <div class="flex items-center gap-2">
                <MenuIconPreview :icon="icon" :size="16" />
                <span>{{ icon }}</span>
              </div>
            </ElOption>
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="菜单类型" required>
          <ElRadioGroup v-model="form.kind">
            <ElRadio value="directory">目录</ElRadio>
            <ElRadio value="page">菜单页</ElRadio>
            <ElRadio value="button">按钮权限</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :precision="0" class="!w-48" />
        </ElFormItem>
        <ElFormItem label="是否显示">
          <ElSwitch v-model="form.is_show" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.menu-manage-tabs :deep(.el-tabs__header) {
  margin-bottom: 14px;
}
</style>
