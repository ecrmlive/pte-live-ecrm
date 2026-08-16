<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { DiyLinkScope, DiyPageCategory } from '#/api/core/diy';

import { computed, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElForm, ElFormItem, ElInput, ElInputNumber, ElMessage, ElOption, ElSelect, ElSwitch } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { createDiyPageCategoryApi, deleteDiyPageCategoryApi, listDiyPageCategoriesApi, updateDiyPageCategoryApi } from '#/api/core/diy';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { LIST_ENABLE_STATUS_FIELD, listFormOptionsDefaults } from '#/utils/list-form-defaults';

const route = useRoute();
const scope = computed<DiyLinkScope>(() => route.path.includes('/merchant/') ? 'merchant' : 'platform');
const scopeLabel = computed(() => scope.value === 'merchant' ? '商户页面分类' : '平台页面分类');
const treeRows = ref<DiyPageCategory[]>([]);
const editingID = ref<number>();
const form = reactive({ pid: 0, name: '', sort: 0, status: true });

function flatten(list: DiyPageCategory[], out: DiyPageCategory[] = []) {
  for (const item of list) {
    out.push(item);
    flatten(item.children || [], out);
  }
  return out;
}

function filterTree(list: DiyPageCategory[], keyword: string, status?: number): DiyPageCategory[] {
  return list.reduce<DiyPageCategory[]>((result, item) => {
    const children = filterTree(item.children || [], keyword, status);
    const matchesName = !keyword || item.name.toLowerCase().includes(keyword);
    const matchesStatus = status === undefined || item.status === status;
    if (matchesName && matchesStatus || children.length) result.push({ ...item, children });
    return result;
  }, []);
}

const parentOptions = computed(() => flatten(treeRows.value).filter((item) => item.level < 3));
const formOptions = listFormOptionsDefaults([
  { component: 'Input', componentProps: { clearable: true, placeholder: '请输入分类名称' }, fieldName: 'name', label: '分类搜索' },
  LIST_ENABLE_STATUS_FIELD('显示状态'),
]);

const gridOptions: VxeGridProps<DiyPageCategory> = {
  columns: [
    { field: 'id', title: 'ID', width: 82 },
    { field: 'name', minWidth: 220, slots: { default: 'name' }, title: '名称', treeNode: true },
    { field: 'type', formatter: ({ cellValue }) => cellValue || 'link', minWidth: 140, title: '类型' },
    { field: 'status', slots: { default: 'status' }, title: '是否显示', width: 120 },
    { field: 'add_time', formatter: ({ cellValue }) => cellValue ? formatShanghaiDateTime(cellValue) : '—', minWidth: 175, title: '添加时间' },
    platformListActionColumn({ width: 150 }),
  ],
  pagerConfig: { enabled: false },
  proxyConfig: {
    ajax: {
      query: async (_ctx, values) => {
        const list = (await listDiyPageCategoriesApi(scope.value)).list || [];
        treeRows.value = list;
        const keyword = String(values?.name ?? '').trim().toLowerCase();
        const rawStatus = values?.status;
        const status = rawStatus === 0 || rawStatus === 1 ? Number(rawStatus) : undefined;
        const filtered = keyword || status !== undefined ? filterTree(list, keyword, status) : list;
        return { items: filtered, total: filtered.length };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  treeConfig: { childrenField: 'children', expandAll: true, indent: 12 },
  toolbarConfig: { custom: false, export: false, refresh: false, search: false, zoom: false },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });
const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]', confirmText: '保存', cancelText: '取消', placement: 'right', onConfirm: async () => submit(),
});

function resetForm(row?: DiyPageCategory, parentID = 0) {
  editingID.value = row?.id;
  Object.assign(form, { pid: row?.pid ?? parentID, name: row?.name || '', sort: row?.sort || 0, status: row ? row.status === 1 : true });
}

function openCreate(parentID = 0) {
  resetForm(undefined, parentID);
  formDrawerApi.setState({ title: `添加${scopeLabel.value}` }).open();
}

function openEdit(row: DiyPageCategory) {
  resetForm(row);
  formDrawerApi.setState({ title: `编辑${scopeLabel.value}` }).open();
}

async function submit() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写分类名称');
    return;
  }
  const body = { ...form, name: form.name.trim(), status: form.status ? 1 : 0, type: 'link' };
  formDrawerApi.lock();
  try {
    if (editingID.value) {
      await updateDiyPageCategoryApi(editingID.value, scope.value, body);
      ElMessage.success('页面分类已更新');
    } else {
      await createDiyPageCategoryApi(scope.value, body);
      ElMessage.success('页面分类已添加');
    }
    formDrawerApi.close();
    await gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function updateStatus(row: DiyPageCategory, enabled: boolean) {
  try {
    await updateDiyPageCategoryApi(row.id, scope.value, { name: row.name, pid: row.pid, sort: row.sort, status: enabled ? 1 : 0, type: row.type || 'link' });
    row.status = enabled ? 1 : 0;
    ElMessage.success(enabled ? '已显示' : '已隐藏');
  } catch {
    row.status = enabled ? 0 : 1;
  }
}

async function remove(row: DiyPageCategory) {
  try {
    await confirm({ content: `确认删除「${row.name}」吗？仅允许删除没有子分类和链接的分类。`, icon: 'warning', title: '删除页面分类' });
  } catch {
    return;
  }
  await deleteDiyPageCategoryApi(row.id, scope.value);
  ElMessage.success('已删除');
  await gridApi.reload();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate()">添加分类</ElButton>
      </template>
      <template #name="{ row }">
        <span>{{ row.name }}</span>
      </template>
      <template #status="{ row }">
        <ElSwitch :model-value="row.status === 1" active-text="显示" inactive-text="隐藏" inline-prompt @change="(value) => updateStatus(row, Boolean(value))" />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>
    <FormDrawer>
      <ElForm label-width="96px">
        <ElFormItem label="上级分类">
          <ElSelect v-model="form.pid" class="w-full" filterable>
            <ElOption :value="0" label="顶级分类" />
            <ElOption v-for="item in parentOptions.filter((item) => item.id !== editingID)" :key="item.id" :value="item.id" :label="`${'— '.repeat(Math.max(0, item.level - 1))}${item.name}`" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="分类名称" required><ElInput v-model="form.name" maxlength="50" show-word-limit /></ElFormItem>
        <ElFormItem label="类型"><ElInput model-value="link" disabled /></ElFormItem>
        <ElFormItem label="排序"><ElInputNumber v-model="form.sort" :min="0" :max="99999" /></ElFormItem>
        <ElFormItem label="是否显示"><ElSwitch v-model="form.status" active-text="显示" inactive-text="隐藏" inline-prompt /></ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
