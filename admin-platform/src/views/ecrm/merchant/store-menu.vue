<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElRadio,
  ElRadioGroup,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createMerchantStoreMenu,
  deleteMerchantStoreMenu,
  fetchMerchantStoreMenus,
  updateMerchantStoreMenu,
  type MerchantStoreMenuRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import MenuIconPreview from '#/components/platform-menu/MenuIconPreview.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

type FormMode = 'create' | 'create-child' | 'edit';

const canManage = ref(true);
const editing = ref<MerchantStoreMenuRow>();
const formMode = ref<FormMode>('create');
const parentId = ref(0);

const form = reactive({
  code: '',
  name: '',
  path: '',
  component: '',
  icon: '',
  is_menu: 1,
  is_route: 1,
  sort: 0,
  status: 1,
});

const gridOptions: VxeGridProps<MerchantStoreMenuRow> = {
  columns: [
    {
      field: 'menu_name',
      minWidth: 260,
      showOverflow: false,
      slots: { default: 'menu_name' },
      title: '菜单名称',
      treeNode: true,
    },
    {
      field: 'path',
      minWidth: 220,
      showOverflow: false,
      title: '菜单地址',
    },
    {
      field: 'icon',
      minWidth: 180,
      showOverflow: false,
      slots: { default: 'icon' },
      title: '菜单图标',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      title: '创建时间',
      width: 180,
    },
    platformListActionColumn({ minWidth: 240 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async () => {
        const list = await fetchMerchantStoreMenus();
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => save(),
});

function resetForm() {
  form.code = '';
  form.name = '';
  form.path = '';
  form.component = '';
  form.icon = '';
  form.is_menu = 1;
  form.is_route = 1;
  form.sort = 0;
  form.status = 1;
}

function openCreate() {
  formMode.value = 'create';
  parentId.value = 0;
  editing.value = undefined;
  resetForm();
  formModalApi.setState({ title: '添加菜单' }).open();
}

function openCreateChild(row: MerchantStoreMenuRow) {
  formMode.value = 'create-child';
  parentId.value = row.menu_id;
  editing.value = undefined;
  resetForm();
  form.is_route = row.is_menu === 1 ? 1 : 0;
  formModalApi.setState({ title: `添加子菜单：${row.menu_name}` }).open();
}

function openEdit(row: MerchantStoreMenuRow) {
  formMode.value = 'edit';
  parentId.value = row.pid;
  editing.value = row;
  form.code = row.code;
  form.name = row.menu_name;
  form.path = row.path;
  form.component = row.component || row.route || '';
  form.icon = row.icon;
  form.is_menu = row.is_menu;
  form.is_route = row.is_route;
  form.sort = row.sort;
  form.status = row.is_show;
  formModalApi.setState({ title: '编辑菜单' }).open();
}

async function save() {
  const name = form.name.trim();
  const code = form.code.trim();
  if (!name) {
    ElMessage.warning('菜单名称不能为空');
    return;
  }
  if (!code) {
    ElMessage.warning('菜单权限码不能为空');
    return;
  }
  const payload = {
    parent_id: parentId.value,
    code,
    name,
    path: form.path.trim(),
    component: form.component.trim(),
    icon: form.icon.trim(),
    is_menu: form.is_menu,
    is_route: form.is_route,
    sort: form.sort,
    status: form.status,
  };
  formModalApi.lock();
  try {
    if (formMode.value === 'edit' && editing.value) {
      await updateMerchantStoreMenu(editing.value.menu_id, payload);
    } else {
      await createMerchantStoreMenu(payload);
    }
    formModalApi.close();
    ElMessage.success('菜单已保存');
    gridApi.reload();
  } finally {
    formModalApi.unlock();
  }
}

async function onDelete(row: MerchantStoreMenuRow) {
  try {
    await ElMessageBox.confirm(
      `确认删除菜单「${row.menu_name}」？存在子菜单时不可删除。`,
      '删除菜单',
      { type: 'warning' },
    );
  } catch {
    return;
  }
  await deleteMerchantStoreMenu(row.menu_id);
  ElMessage.success('菜单已删除');
  gridApi.reload();
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canManage.value = permissions.includes('store.menu');
});
</script>

<template>
  <Page auto-content-height title="店铺菜单">
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          添加菜单
        </ElButton>
      </template>
      <template #menu_name="{ row }">
        <span>{{ row.menu_name }} [{{ row.menu_id }}]</span>
      </template>
      <template #icon="{ row }">
        <span v-if="row.icon" class="inline-flex items-center gap-2">
          <MenuIconPreview :icon="row.icon" :size="16" />
          <span class="text-xs text-muted-foreground">{{ row.icon }}</span>
        </span>
        <span v-else class="text-muted-foreground">—</span>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton
            v-if="row.is_menu === 1"
            link
            type="primary"
            @click="openCreateChild(row)"
          >
            添加子菜单
          </ElButton>
          <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
          <ElButton link type="danger" @click="onDelete(row)">删除</ElButton>
        </template>
      </template>
    </Grid>

    <FormModal>
      <ElForm label-width="100px">
        <ElFormItem v-if="parentId" label="父级 ID">
          <ElInput :model-value="String(parentId)" disabled />
        </ElFormItem>
        <ElFormItem label="权限码" required>
          <ElInput v-model="form.code" placeholder="如 product.list" />
        </ElFormItem>
        <ElFormItem label="菜单名称" required>
          <ElInput v-model="form.name" />
        </ElFormItem>
        <ElFormItem label="菜单地址">
          <ElInput v-model="form.path" placeholder="/product/list" />
        </ElFormItem>
        <ElFormItem label="组件路径">
          <ElInput
            v-model="form.component"
            placeholder="views/ecrm/product/list.vue"
          />
        </ElFormItem>
        <ElFormItem label="菜单图标">
          <ElInput v-model="form.icon" placeholder="lucide:package" />
        </ElFormItem>
        <ElFormItem label="菜单类型">
          <ElRadioGroup v-model="form.is_menu">
            <ElRadio :value="1">页面</ElRadio>
            <ElRadio :value="2">按钮</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="是否路由">
          <ElSwitch v-model="form.is_route" :active-value="1" :inactive-value="0" />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" />
        </ElFormItem>
        <ElFormItem label="显示状态">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormModal>
  </Page>
</template>
