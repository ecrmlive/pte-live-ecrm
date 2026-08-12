<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
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
  ElSwitch,
} from 'element-plus';
import { ArrowLeft, Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createConfigClassificationApi,
  createConfigClassificationItemApi,
  deleteConfigClassificationApi,
  deleteConfigClassificationItemApi,
  listConfigClassificationItemsApi,
  listConfigClassificationsApi,
  setConfigClassificationItemStatusApi,
  setConfigClassificationStatusApi,
  updateConfigClassificationApi,
  updateConfigClassificationItemApi,
  type ConfigClassification,
  type ConfigClassificationItem,
} from '#/api/core/platform-config-classification';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_ENABLE_STATUS_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const READ_CODE = 'maintain.config_classify';
const canRead = ref(false);
const canManage = ref(false);
const activeClassification = ref<ConfigClassification>();
const editingClassification = ref<ConfigClassification>();
const editingItem = ref<ConfigClassificationItem>();

const classificationForm = reactive({
  classify_key: '',
  description: '',
  icon: '',
  name: '',
  sort: 0,
  status: 1,
});
const itemForm = reactive({
  config_key: '',
  content: '',
  description: '',
  name: '',
  sort: 0,
  status: 1,
});

const classificationFormOptions: VbenFormProps = listFormOptionsDefaults([
  {
    ...LIST_ENABLE_STATUS_FIELD('是否显示'),
    componentProps: {
      clearable: true,
      options: [
        { label: '显示', value: 1 },
        { label: '隐藏', value: 0 },
      ],
      placeholder: '请选择',
    },
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入分类名称' },
    fieldName: 'name',
    label: '分类名称',
  },
]);

const classificationGridOptions: VxeGridProps<ConfigClassification> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    { field: 'name', minWidth: 180, title: '配置分类名称' },
    { field: 'classify_key', minWidth: 180, title: '配置分类key' },
    { field: 'description', minWidth: 220, showOverflow: 'tooltip', title: '配置分类说明' },
    { field: 'icon', minWidth: 150, title: '图标' },
    { field: 'status', slots: { default: 'status' }, title: '是否显示', width: 130 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '创建时间',
    },
    platformListActionColumn({ width: 230 }),
  ],
  emptyText: '暂无配置分类',
  formOptions: classificationFormOptions,
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const status = Number(formValues?.status);
        const result = await listConfigClassificationsApi({
          limit: page.pageSize,
          name: String(formValues?.name ?? '').trim() || undefined,
          page: page.currentPage,
          status: status === 0 || status === 1 ? status : undefined,
        });
        return { items: result.list || [], total: result.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: { custom: false, export: false, refresh: false, zoom: false },
};

const itemGridOptions: VxeGridProps<ConfigClassificationItem> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    { field: 'name', minWidth: 180, title: '配置名称' },
    { field: 'config_key', minWidth: 220, title: '配置key' },
    { field: 'content', minWidth: 260, showOverflow: 'tooltip', title: '配置内容' },
    { field: 'description', minWidth: 220, showOverflow: 'tooltip', title: '配置说明' },
    { field: 'status', slots: { default: 'itemStatus' }, title: '是否显示', width: 130 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '创建时间',
    },
    platformListActionColumn({ width: 160 }),
  ],
  emptyText: '暂无配置项',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!activeClassification.value) return { items: [], total: 0 };
        const result = await listConfigClassificationItemsApi(activeClassification.value.id, {
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

const [ClassificationGrid, classificationGridApi] = useVbenVxeGrid({
  gridOptions: classificationGridOptions,
});
const [ItemGrid, itemGridApi] = useVbenVxeGrid({ gridOptions: itemGridOptions });

const [ClassificationDrawer, classificationDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[720px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: async () => saveClassification(),
});
const [ItemDrawer, itemDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[720px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: async () => saveItem(),
});

const itemTitle = computed(() => `${activeClassification.value?.name ?? ''}配置列表`);

function resetClassificationForm() {
  editingClassification.value = undefined;
  Object.assign(classificationForm, {
    classify_key: '', description: '', icon: '', name: '', sort: 0, status: 1,
  });
}

function openCreateClassification() {
  resetClassificationForm();
  classificationDrawerApi.setState({ title: '新增配置分类' }).open();
}

function openEditClassification(row: ConfigClassification) {
  editingClassification.value = row;
  Object.assign(classificationForm, {
    classify_key: row.classify_key,
    description: row.description,
    icon: row.icon,
    name: row.name,
    sort: row.sort,
    status: row.status,
  });
  classificationDrawerApi.setState({ title: '编辑配置分类' }).open();
}

async function saveClassification() {
  const name = classificationForm.name.trim();
  const key = classificationForm.classify_key.trim();
  if (!name || !key) {
    ElMessage.warning('请填写配置分类名称和配置分类Key');
    return;
  }
  classificationDrawerApi.lock();
  try {
    const payload = {
      classify_key: key,
      description: classificationForm.description.trim(),
      icon: classificationForm.icon.trim(),
      name,
      sort: Number(classificationForm.sort || 0),
      status: classificationForm.status,
    };
    if (editingClassification.value) {
      await updateConfigClassificationApi(editingClassification.value.id, payload);
      ElMessage.success('配置分类已保存');
    } else {
      await createConfigClassificationApi(payload);
      ElMessage.success('配置分类已创建');
    }
    classificationDrawerApi.close();
    classificationGridApi.reload();
  } finally {
    classificationDrawerApi.unlock();
  }
}

async function toggleClassificationStatus(row: ConfigClassification, visible: boolean) {
  try {
    await setConfigClassificationStatusApi(row.id, visible ? 1 : 0);
    ElMessage.success('显示状态已保存');
  } finally {
    classificationGridApi.reload();
  }
}

async function removeClassification(row: ConfigClassification) {
  try {
    await confirm({
      content: `删除配置分类“${row.name}”及其配置项后不可恢复，是否继续？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteConfigClassificationApi(row.id);
    ElMessage.success('配置分类已删除');
    classificationGridApi.reload();
  } catch {
    // 用户取消或统一请求层已提示错误。
  }
}

function openItems(row: ConfigClassification) {
  activeClassification.value = row;
  itemGridApi.reload();
}

function closeItems() {
  activeClassification.value = undefined;
  editingItem.value = undefined;
}

function resetItemForm() {
  editingItem.value = undefined;
  Object.assign(itemForm, {
    config_key: '', content: '', description: '', name: '', sort: 0, status: 1,
  });
}

function openCreateItem() {
  resetItemForm();
  itemDrawerApi.setState({ title: '新增配置项' }).open();
}

function openEditItem(row: ConfigClassificationItem) {
  editingItem.value = row;
  Object.assign(itemForm, {
    config_key: row.config_key,
    content: row.content,
    description: row.description,
    name: row.name,
    sort: row.sort,
    status: row.status,
  });
  itemDrawerApi.setState({ title: '编辑配置项' }).open();
}

async function saveItem() {
  if (!activeClassification.value) return;
  const name = itemForm.name.trim();
  const key = itemForm.config_key.trim();
  const content = itemForm.content.trim();
  if (!name || !key || !content) {
    ElMessage.warning('请填写配置名称、配置Key和配置内容');
    return;
  }
  itemDrawerApi.lock();
  try {
    const payload = {
      config_key: key,
      content,
      description: itemForm.description.trim(),
      name,
      sort: Number(itemForm.sort || 0),
      status: itemForm.status,
    };
    if (editingItem.value) {
      await updateConfigClassificationItemApi(activeClassification.value.id, editingItem.value.id, payload);
      ElMessage.success('配置项已保存');
    } else {
      await createConfigClassificationItemApi(activeClassification.value.id, payload);
      ElMessage.success('配置项已创建');
    }
    itemDrawerApi.close();
    itemGridApi.reload();
  } finally {
    itemDrawerApi.unlock();
  }
}

async function toggleItemStatus(row: ConfigClassificationItem, visible: boolean) {
  if (!activeClassification.value) return;
  try {
    await setConfigClassificationItemStatusApi(activeClassification.value.id, row.id, visible ? 1 : 0);
    ElMessage.success('显示状态已保存');
  } finally {
    itemGridApi.reload();
  }
}

async function removeItem(row: ConfigClassificationItem) {
  if (!activeClassification.value) return;
  try {
    await confirm({ content: `删除配置项“${row.name}”后不可恢复，是否继续？`, icon: 'warning', title: '提示' });
    await deleteConfigClassificationItemApi(activeClassification.value.id, row.id);
    ElMessage.success('配置项已删除');
    itemGridApi.reload();
  } catch {
    // 用户取消或统一请求层已提示错误。
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.includes('platform') && permissions.includes(READ_CODE);
  canManage.value = canRead.value;
  if (canRead.value) classificationGridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <template v-if="!activeClassification">
      <ClassificationGrid>
        <template #toolbar-actions>
          <ElButton v-if="canManage" :icon="Plus" type="primary" @click="openCreateClassification">
            新增配置分类
          </ElButton>
        </template>
        <template #status="{ row }">
          <ElSwitch
            :model-value="row.status === 1"
            aria-label="是否显示"
            @change="(visible) => toggleClassificationStatus(row, Boolean(visible))"
          />
        </template>
        <template #action="{ row }">
          <ElButton link type="primary" @click="openItems(row)">配置列表</ElButton>
          <ElButton v-if="canManage" link type="primary" @click="openEditClassification(row)">编辑</ElButton>
          <ElButton v-if="canManage" link type="danger" @click="removeClassification(row)">删除</ElButton>
        </template>
      </ClassificationGrid>
    </template>

    <template v-else>
      <ItemGrid>
        <template #toolbar-actions>
          <div class="flex items-center gap-3">
            <ElButton :icon="ArrowLeft" @click="closeItems">返回配置分类</ElButton>
            <span class="text-base font-medium">{{ itemTitle }}</span>
            <ElButton v-if="canManage" :icon="Plus" type="primary" @click="openCreateItem">
              新增配置项
            </ElButton>
          </div>
        </template>
        <template #itemStatus="{ row }">
          <ElSwitch
            :model-value="row.status === 1"
            aria-label="是否显示"
            @change="(visible) => toggleItemStatus(row, Boolean(visible))"
          />
        </template>
        <template #action="{ row }">
          <ElButton v-if="canManage" link type="primary" @click="openEditItem(row)">编辑</ElButton>
          <ElButton v-if="canManage" link type="danger" @click="removeItem(row)">删除</ElButton>
        </template>
      </ItemGrid>
    </template>

    <ClassificationDrawer>
      <ElForm label-width="120px">
        <ElFormItem label="配置分类名称" required>
          <ElInput v-model="classificationForm.name" maxlength="128" placeholder="请输入配置分类名称" />
        </ElFormItem>
        <ElFormItem label="配置分类key" required>
          <ElInput v-model="classificationForm.classify_key" maxlength="64" placeholder="请输入配置分类key" />
        </ElFormItem>
        <ElFormItem label="配置分类说明">
          <ElInput v-model="classificationForm.description" :rows="3" maxlength="500" placeholder="请输入配置分类说明" type="textarea" />
        </ElFormItem>
        <ElFormItem label="图标">
          <ElInput v-model="classificationForm.icon" maxlength="96" placeholder="例如 lucide:settings" />
        </ElFormItem>
        <ElFormItem label="排序"><ElInputNumber v-model="classificationForm.sort" :step="1" /></ElFormItem>
        <ElFormItem label="是否显示"><ElSwitch v-model="classificationForm.status" :active-value="1" :inactive-value="0" /></ElFormItem>
      </ElForm>
    </ClassificationDrawer>

    <ItemDrawer>
      <ElForm label-width="100px">
        <ElFormItem label="配置名称" required>
          <ElInput v-model="itemForm.name" maxlength="128" placeholder="请输入配置名称" />
        </ElFormItem>
        <ElFormItem label="配置key" required>
          <ElInput v-model="itemForm.config_key" maxlength="128" placeholder="请输入配置key" />
        </ElFormItem>
        <ElFormItem label="配置内容" required>
          <ElInput v-model="itemForm.content" :rows="6" maxlength="10000" placeholder="请输入配置内容" type="textarea" />
        </ElFormItem>
        <ElFormItem label="配置说明">
          <ElInput v-model="itemForm.description" :rows="3" maxlength="500" placeholder="请输入配置说明" type="textarea" />
        </ElFormItem>
        <ElFormItem label="排序"><ElInputNumber v-model="itemForm.sort" :step="1" /></ElFormItem>
        <ElFormItem label="是否显示"><ElSwitch v-model="itemForm.status" :active-value="1" :inactive-value="0" /></ElFormItem>
      </ElForm>
    </ItemDrawer>
  </Page>
</template>
