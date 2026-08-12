<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

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
} from 'element-plus';
import { ArrowLeft, Delete, Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createDataGroupApi,
  createDataGroupItemApi,
  deleteDataGroupApi,
  deleteDataGroupItemApi,
  listDataGroupItemsApi,
  listDataGroupsApi,
  setDataGroupItemStatusApi,
  updateDataGroupApi,
  updateDataGroupItemApi,
  type DataGroup,
  type DataGroupField,
  type DataGroupItem,
} from '#/api/core/platform-data-group';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';

const READ_CODE = 'maintain.group_data';
const MANAGE_CODE = 'maintain.group_data.manage';

const canRead = ref(false);
const canManage = ref(false);
const editingGroup = ref<DataGroup>();
const activeGroup = ref<DataGroup>();
const editingItem = ref<DataGroupItem>();
const groupForm = reactive({
  description: '',
  group_key: '',
  name: '',
  sort: 0,
});
const fieldRows = ref<DataGroupField[]>([]);
const itemForm = reactive({ dataText: '{}', sort: 0, status: 1 });
const itemValues = ref<Record<string, number | string | undefined>>({});

const groupGridOptions: VxeGridProps<DataGroup> = {
  columns: [
    { field: 'id', title: 'ID', width: 100 },
    { field: 'name', minWidth: 220, title: '数据组名称' },
    { field: 'group_key', minWidth: 240, title: '数据组key' },
    {
      field: 'description',
      minWidth: 260,
      showOverflow: 'tooltip',
      title: '数据组说明',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '创建时间',
    },
    platformListActionColumn({ width: 220 }),
  ],
  emptyText: '暂无组合数据',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!canRead.value) return { items: [], total: 0 };
        const result = await listDataGroupsApi({
          limit: page.pageSize,
          page: page.currentPage,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: { custom: false, export: false, refresh: false, zoom: false },
};

const itemGridOptions: VxeGridProps<DataGroupItem> = {
  columns: [
    { field: 'id', title: 'ID', width: 100 },
    { field: 'data', minWidth: 420, showOverflow: 'tooltip', slots: { default: 'data' }, title: '数据内容' },
    { field: 'sort', title: '排序', width: 110 },
    { field: 'status', slots: { default: 'status' }, title: '是否显示', width: 130 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '创建时间',
    },
    platformListActionColumn({ width: 160 }),
  ],
  emptyText: '暂无数据',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!activeGroup.value) return { items: [], total: 0 };
        const result = await listDataGroupItemsApi(activeGroup.value.id, {
          limit: page.pageSize,
          page: page.currentPage,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: { custom: false, export: false, refresh: false, zoom: false },
};

const [GroupGrid, groupGridApi] = useVbenVxeGrid({
  gridOptions: groupGridOptions,
});
const [ItemGrid, itemGridApi] = useVbenVxeGrid({ gridOptions: itemGridOptions });

const [GroupDrawer, groupDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[720px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: async () => saveGroup(),
});
const [ItemDrawer, itemDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[720px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: async () => saveItem(),
});

const itemTitle = computed(() => activeGroup.value?.name || '数据列表');
const itemFields = computed(() => activeGroup.value?.fields || []);

function resetGroupForm() {
  editingGroup.value = undefined;
  Object.assign(groupForm, { description: '', group_key: '', name: '', sort: 0 });
  fieldRows.value = [];
}

function openCreateGroup() {
  resetGroupForm();
  groupDrawerApi.setState({ title: '新增组合数据' }).open();
}

function openEditGroup(row: DataGroup) {
  editingGroup.value = row;
  Object.assign(groupForm, {
    description: row.description,
    group_key: row.group_key,
    name: row.name,
    sort: row.sort,
  });
  fieldRows.value = (row.fields || []).map((field) => ({ ...field }));
  groupDrawerApi.setState({ title: '编辑组合数据' }).open();
}

function addField() {
  fieldRows.value.push({ field: '', name: '', type: 'input' });
}

function removeField(index: number) {
  fieldRows.value.splice(index, 1);
}

function validateFields(): DataGroupField[] | undefined {
  const keys = new Set<string>();
  for (const item of fieldRows.value) {
    const field = item.field.trim();
    const name = item.name.trim();
    if (!field || !name || !/^[A-Za-z][A-Za-z0-9_]*$/.test(field)) {
      ElMessage.warning('字段名称和字段 Key 必填，字段 Key 仅支持字母、数字和下划线');
      return;
    }
    if (keys.has(field)) {
      ElMessage.warning('字段 Key 不能重复');
      return;
    }
    keys.add(field);
  }
  return fieldRows.value.map((item) => ({
    field: item.field.trim(),
    name: item.name.trim(),
    type: item.type || 'input',
  }));
}

async function saveGroup() {
  const fields = validateFields();
  if (!fields || !groupForm.name.trim() || !groupForm.group_key.trim()) {
    if (!groupForm.name.trim() || !groupForm.group_key.trim()) {
      ElMessage.warning('请填写数据组名称和数据组 Key');
    }
    return;
  }
  groupDrawerApi.lock();
  try {
    const payload = {
      description: groupForm.description.trim(),
      fields,
      group_key: groupForm.group_key.trim(),
      name: groupForm.name.trim(),
      sort: Number(groupForm.sort || 0),
    };
    if (editingGroup.value) {
      await updateDataGroupApi(editingGroup.value.id, payload);
      ElMessage.success('组合数据已保存');
    } else {
      await createDataGroupApi(payload);
      ElMessage.success('组合数据已创建');
    }
    groupDrawerApi.close();
    groupGridApi.reload();
  } finally {
    groupDrawerApi.unlock();
  }
}

async function removeGroup(row: DataGroup) {
  try {
    await confirm({ content: `删除组合数据“${row.name}”及其数据项后不可恢复，是否继续？`, icon: 'warning', title: '提示' });
    await deleteDataGroupApi(row.id);
    ElMessage.success('组合数据已删除');
    groupGridApi.reload();
  } catch {
    // 用户取消或请求层已提示错误。
  }
}

function openItems(row: DataGroup) {
  activeGroup.value = row;
  itemGridApi.reload();
}

function closeItems() {
  activeGroup.value = undefined;
  editingItem.value = undefined;
}

function resetItemForm() {
  editingItem.value = undefined;
  Object.assign(itemForm, { dataText: '{}', sort: 0, status: 1 });
  itemValues.value = {};
}

function openCreateItem() {
  resetItemForm();
  itemDrawerApi.setState({ title: '新增数据' }).open();
}

function openEditItem(row: DataGroupItem) {
  editingItem.value = row;
  Object.assign(itemForm, { dataText: JSON.stringify(row.data || {}, null, 2), sort: row.sort, status: row.status });
  itemValues.value = Object.fromEntries(
    Object.entries(row.data || {}).map(([key, value]) => [
      key,
      typeof value === 'number' ? value : String(value ?? ''),
    ]),
  );
  itemDrawerApi.setState({ title: '编辑数据' }).open();
}

function parseItemData(): Record<string, unknown> | undefined {
  try {
    const data = JSON.parse(itemForm.dataText || '{}');
    if (!data || Array.isArray(data) || typeof data !== 'object') throw new Error('not object');
    return data;
  } catch {
    ElMessage.warning('数据内容必须是 JSON 对象');
  }
}

async function saveItem() {
  if (!activeGroup.value) return;
  const data = itemFields.value.length > 0 ? itemValues.value : parseItemData();
  if (!data) return;
  itemDrawerApi.lock();
  try {
    const payload = { data, sort: Number(itemForm.sort || 0), status: itemForm.status };
    if (editingItem.value) {
      await updateDataGroupItemApi(activeGroup.value.id, editingItem.value.id, payload);
      ElMessage.success('数据已保存');
    } else {
      await createDataGroupItemApi(activeGroup.value.id, payload);
      ElMessage.success('数据已创建');
    }
    itemDrawerApi.close();
    itemGridApi.reload();
  } finally {
    itemDrawerApi.unlock();
  }
}

async function toggleItemStatus(row: DataGroupItem, visible: boolean) {
  if (!activeGroup.value) return;
  try {
    await setDataGroupItemStatusApi(activeGroup.value.id, row.id, visible ? 1 : 0);
    ElMessage.success('显示状态已更新');
  } finally {
    itemGridApi.reload();
  }
}

async function removeItem(row: DataGroupItem) {
  if (!activeGroup.value) return;
  try {
    await confirm({ content: '删除该数据后不可恢复，是否继续？', icon: 'warning', title: '提示' });
    await deleteDataGroupItemApi(activeGroup.value.id, row.id);
    ElMessage.success('数据已删除');
    itemGridApi.reload();
  } catch {
    // 用户取消或请求层已提示错误。
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  const isPlatform = profile.roles.includes('platform');
  canRead.value = isPlatform && (permissions.includes(READ_CODE) || permissions.includes(MANAGE_CODE));
  canManage.value = isPlatform && permissions.includes(MANAGE_CODE);
  if (canRead.value) groupGridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <template v-if="!activeGroup">
      <GroupGrid>
        <template #toolbar-actions>
          <ElButton v-if="canManage" :icon="Plus" type="primary" @click="openCreateGroup">
            新增组合数据
          </ElButton>
        </template>
        <template #action="{ row }">
          <ElButton link type="primary" @click="openItems(row)">数据列表</ElButton>
          <ElButton v-if="canManage" link type="primary" @click="openEditGroup(row)">编辑</ElButton>
          <ElButton v-if="canManage" link type="danger" @click="removeGroup(row)">删除</ElButton>
        </template>
      </GroupGrid>
    </template>

    <template v-else>
      <ItemGrid>
        <template #toolbar-actions>
          <div class="flex items-center gap-3">
            <ElButton :icon="ArrowLeft" @click="closeItems">返回组合数据</ElButton>
            <span class="text-base font-medium">{{ itemTitle }}</span>
            <ElButton v-if="canManage" :icon="Plus" type="primary" @click="openCreateItem">
              新增数据
            </ElButton>
          </div>
        </template>
        <template #data="{ row }">
          <span class="block truncate text-left font-mono text-xs">{{ JSON.stringify(row.data) }}</span>
        </template>
        <template #status="{ row }">
          <ElSwitch :model-value="row.status === 1" aria-label="是否显示" @change="(visible) => toggleItemStatus(row, Boolean(visible))" />
        </template>
        <template #action="{ row }">
          <ElButton v-if="canManage" link type="primary" @click="openEditItem(row)">编辑</ElButton>
          <ElButton v-if="canManage" link type="danger" @click="removeItem(row)">删除</ElButton>
        </template>
      </ItemGrid>
    </template>

    <GroupDrawer>
      <ElForm label-width="120px">
        <ElFormItem label="数据组名称" required>
          <ElInput v-model="groupForm.name" maxlength="128" placeholder="请输入数据组名称" />
        </ElFormItem>
        <ElFormItem label="数据组key" required>
          <ElInput v-model="groupForm.group_key" maxlength="64" placeholder="请输入数据组key" />
        </ElFormItem>
        <ElFormItem label="数据组说明">
          <ElInput v-model="groupForm.description" maxlength="500" placeholder="请输入数据组说明" type="textarea" :rows="3" />
        </ElFormItem>
        <ElFormItem label="字段定义">
          <div class="w-full space-y-2">
            <div v-for="(field, index) in fieldRows" :key="index" class="flex gap-2">
              <ElInput v-model="field.name" placeholder="字段名称" />
              <ElInput v-model="field.field" placeholder="字段Key" />
              <ElSelect v-model="field.type" placeholder="字段类型">
                <ElOption label="单行文本" value="input" />
                <ElOption label="多行文本" value="textarea" />
                <ElOption label="数字" value="number" />
                <ElOption label="图片地址" value="image" />
              </ElSelect>
              <ElButton :icon="Delete" circle type="danger" @click="removeField(index)" />
            </div>
            <ElButton :icon="Plus" @click="addField">新增字段</ElButton>
          </div>
        </ElFormItem>
        <ElFormItem label="排序"><ElInputNumber v-model="groupForm.sort" :step="1" /></ElFormItem>
      </ElForm>
    </GroupDrawer>

    <ItemDrawer>
      <ElForm label-width="96px">
        <template v-if="itemFields.length > 0">
          <ElFormItem
            v-for="field in itemFields"
            :key="field.field"
            :label="field.name"
            required
          >
            <ElInputNumber
              v-if="field.type === 'number'"
              class="!w-full"
              :model-value="Number(itemValues[field.field] ?? 0)"
              @update:model-value="(value) => { itemValues[field.field] = Number(value ?? 0); }"
            />
            <ElInput
              v-else
              :model-value="String(itemValues[field.field] ?? '')"
              :placeholder="`请输入${field.name}`"
              :rows="field.type === 'textarea' ? 4 : undefined"
              :type="field.type === 'textarea' ? 'textarea' : 'text'"
              @update:model-value="(value) => { itemValues[field.field] = value; }"
            />
          </ElFormItem>
        </template>
        <ElFormItem v-else label="数据内容" required>
          <ElInput v-model="itemForm.dataText" type="textarea" :rows="14" placeholder="请输入 JSON 对象" />
        </ElFormItem>
        <ElFormItem label="排序"><ElInputNumber v-model="itemForm.sort" :step="1" /></ElFormItem>
        <ElFormItem label="是否显示"><ElSwitch v-model="itemForm.status" :active-value="1" :inactive-value="0" /></ElFormItem>
      </ElForm>
    </ItemDrawer>
  </Page>
</template>
