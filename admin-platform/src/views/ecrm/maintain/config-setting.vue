<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref } from 'vue';

import { confirm, Page, useVbenDrawer } from '@vben/common-ui';

import { Plus } from '@element-plus/icons-vue';
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

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  type ConfigClassification,
  type ConfigFieldType,
  listConfigClassificationsApi,
} from '#/api/core/platform-config-classification';
import {
  createPlatformConfigSettingApi,
  deletePlatformConfigSettingApi,
  listPlatformConfigSettingsApi,
  type PlatformConfigSetting,
  setPlatformConfigSettingStatusApi,
  updatePlatformConfigSettingApi,
} from '#/api/core/platform-config-setting';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const READ_CODE = 'maintain.config_setting';
const canRead = ref(false);
const canManage = ref(false);
const classifications = ref<ConfigClassification[]>([]);
const editing = ref<PlatformConfigSetting>();

const form = reactive({
  backend_type: 0,
  classification_id: undefined as number | undefined,
  config_key: '',
  content: '',
  description: '',
  field_type: 'input' as ConfigFieldType,
  name: '',
  sort: 0,
  status: 1,
});

const fieldTypes: Array<{ label: string; value: ConfigFieldType }> = [
  { label: '文本框', value: 'input' },
  { label: '多行文本', value: 'textarea' },
  { label: '数字输入', value: 'number' },
  { label: '单选项', value: 'radio' },
  { label: '开关', value: 'switch' },
  { label: '图片', value: 'image' },
  { label: '文件', value: 'file' },
];

const backendTypes = [
  { label: '总后台配置', value: 0 },
  { label: '店铺后台配置', value: 1 },
];

function backendTypeLabel(value: unknown) {
  return Number(value) === 1 ? '店铺后台配置' : '总后台配置';
}

function fieldTypeLabel(value: unknown) {
  return fieldTypes.find((item) => item.value === value)?.label || '文本框';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入名称/key' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: { clearable: true, options: backendTypes, placeholder: '请选择' },
    fieldName: 'backend_type',
    label: '后台类型',
  },
]);

const gridOptions: VxeGridProps<PlatformConfigSetting> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    { field: 'name', minWidth: 170, title: '配置名称' },
    { field: 'config_key', minWidth: 190, title: '配置key' },
    { field: 'description', minWidth: 260, showOverflow: 'tooltip', title: '配置说明' },
    { field: 'field_type', formatter: ({ cellValue }) => fieldTypeLabel(cellValue), title: '类型', width: 120 },
    { field: 'backend_type', formatter: ({ cellValue }) => backendTypeLabel(cellValue), minWidth: 160, title: '后台类型' },
    { field: 'status', slots: { default: 'status' }, title: '是否显示', width: 130 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 180,
      title: '创建时间',
    },
    platformListActionColumn({ width: 160 }),
  ],
  emptyText: '暂无配置',
  formOptions,
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const backendType = Number(formValues?.backend_type);
        const result = await listPlatformConfigSettingsApi({
          backend_type: backendType === 0 || backendType === 1 ? backendType : undefined,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
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

const [ConfigGrid, configGridApi] = useVbenVxeGrid({ gridOptions });
const [ConfigDrawer, configDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[760px] max-w-[96vw]',
  confirmText: '保存',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    backend_type: 0,
    classification_id: undefined,
    config_key: '',
    content: '',
    description: '',
    field_type: 'input',
    name: '',
    sort: 0,
    status: 1,
  });
}

function openCreate() {
  resetForm();
  configDrawerApi.setState({ title: '新增配置' }).open();
}

function openEdit(row: PlatformConfigSetting) {
  editing.value = row;
  Object.assign(form, {
    backend_type: row.backend_type,
    classification_id: row.classification_id,
    config_key: row.config_key,
    content: row.content,
    description: row.description,
    field_type: row.field_type || 'input',
    name: row.name,
    sort: row.sort,
    status: row.status,
  });
  configDrawerApi.setState({ title: '编辑配置' }).open();
}

async function save() {
  const classificationID = Number(form.classification_id);
  const name = form.name.trim();
  const configKey = form.config_key.trim();
  const content = form.content.trim();
  if (!classificationID || !name || !configKey || !content) {
    ElMessage.warning('请填写配置分类、配置名称、配置key和配置内容');
    return;
  }
  configDrawerApi.lock();
  try {
    const payload = {
      backend_type: form.backend_type,
      classification_id: classificationID,
      config_key: configKey,
      content,
      description: form.description.trim(),
      field_type: form.field_type,
      name,
      sort: Number(form.sort || 0),
      status: form.status,
    };
    if (editing.value) {
      await updatePlatformConfigSettingApi(editing.value.id, payload);
    } else {
      await createPlatformConfigSettingApi(payload);
    }
    ElMessage.success('配置已保存');
    configDrawerApi.close();
    configGridApi.reload();
  } finally {
    configDrawerApi.unlock();
  }
}

async function toggleStatus(row: PlatformConfigSetting, visible: boolean) {
  try {
    await setPlatformConfigSettingStatusApi(row.id, visible ? 1 : 0);
    ElMessage.success('显示状态已保存');
  } finally {
    configGridApi.reload();
  }
}

async function remove(row: PlatformConfigSetting) {
  try {
    await confirm({ content: `删除配置“${row.name}”后不可恢复，是否继续？`, icon: 'warning', title: '提示' });
    await deletePlatformConfigSettingApi(row.id);
    ElMessage.success('配置已删除');
    configGridApi.reload();
  } catch {
    // 用户取消或统一请求层已提示错误。
  }
}

async function loadClassifications() {
  const result = await listConfigClassificationsApi({ limit: 200, page: 1 });
  classifications.value = result.list || [];
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.includes('platform') && permissions.includes(READ_CODE);
  canManage.value = canRead.value;
  if (!canRead.value) return;
  await loadClassifications();
  configGridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ConfigGrid>
      <template #toolbar-actions>
        <ElButton v-if="canManage" :icon="Plus" type="primary" @click="openCreate">
          新增配置
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          aria-label="是否显示"
          @change="(visible) => toggleStatus(row, Boolean(visible))"
        />
      </template>
      <template #action="{ row }">
        <ElButton v-if="canManage" link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton v-if="canManage" link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </ConfigGrid>

    <ConfigDrawer>
      <ElForm label-width="110px">
        <ElFormItem label="配置分类" required>
          <ElSelect v-model="form.classification_id" class="w-full" filterable placeholder="请选择配置分类">
            <ElOption
              v-for="item in classifications"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="配置名称" required>
          <ElInput v-model="form.name" maxlength="128" placeholder="请输入配置名称" />
        </ElFormItem>
        <ElFormItem label="配置key" required>
          <ElInput v-model="form.config_key" maxlength="128" placeholder="请输入配置key" />
        </ElFormItem>
        <ElFormItem label="类型" required>
          <ElSelect v-model="form.field_type" class="w-full" placeholder="请选择类型">
            <ElOption v-for="item in fieldTypes" :key="item.value" :label="item.label" :value="item.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="后台类型" required>
          <ElSelect v-model="form.backend_type" class="w-full" placeholder="请选择后台类型">
            <ElOption v-for="item in backendTypes" :key="item.value" :label="item.label" :value="item.value" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="配置内容" required>
          <ElInput v-model="form.content" :rows="5" maxlength="10000" placeholder="请输入配置内容" type="textarea" />
        </ElFormItem>
        <ElFormItem label="配置说明">
          <ElInput v-model="form.description" :rows="3" maxlength="500" placeholder="请输入配置说明" type="textarea" />
        </ElFormItem>
        <ElFormItem label="排序"><ElInputNumber v-model="form.sort" :step="1" /></ElFormItem>
        <ElFormItem label="是否显示"><ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" /></ElFormItem>
      </ElForm>
    </ConfigDrawer>
  </Page>
</template>
