<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElCascader,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElOption,
  ElSelect,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  copyStoreParameterTemplate,
  createStoreParameterTemplate,
  fetchPlatformMerchants,
  fetchStoreParameterTemplate,
  fetchStoreParameterTemplates,
  type StoreParameterItem,
  type StoreParameterTemplateDetail,
  type StoreParameterTemplateRow,
} from '#/api/core/ecrm';
import {
  listPlatformCategoriesApi,
  type PlatformCategory,
} from '#/api/core/platform-catalog';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

type CascaderOption = {
  label: string;
  value: number;
  children?: CascaderOption[];
};

type ParamRow = StoreParameterItem & { single: string };

const merchantOptions = ref<{ label: string; value: number }[]>([]);
const categoryOptions = ref<{ label: string; value: number }[]>([]);
const categoryTree = ref<CascaderOption[]>([]);
const viewing = ref<StoreParameterTemplateDetail>();
const copySourceId = ref(0);
const saving = ref(false);

const createForm = reactive({
  mer_id: undefined as number | undefined,
  template_name: '',
  cate_id: undefined as number | undefined,
  is_required: 0 as 0 | 1,
  sort: 0,
});
const paramRows = ref<ParamRow[]>([]);

const copyForm = reactive({
  template_name: '',
  cate_ids: [] as number[],
  sort: 0,
  status: 1 as number,
  params: [] as StoreParameterItem[],
});

const formOptions = computed((): VbenFormProps =>
  listFormOptionsDefaults([
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: merchantOptions.value,
        placeholder: '请选择',
      },
      fieldName: 'mer_id',
      label: '店铺名称',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入参数模板名称',
      },
      fieldName: 'template_name',
      label: '模板名称',
    },
  ]),
);

const gridOptions: VxeGridProps<StoreParameterTemplateRow> = {
  columns: [
    { field: 'id', title: 'ID', width: 90 },
    {
      field: 'mer_name',
      minWidth: 160,
      showOverflow: false,
      title: '店铺名称',
    },
    {
      field: 'template_name',
      minWidth: 180,
      showOverflow: false,
      title: '模板名称',
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 128 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const merRaw = formValues?.mer_id;
        const merId =
          merRaw === undefined || merRaw === null || merRaw === ''
            ? undefined
            : Number(merRaw);
        const data = await fetchStoreParameterTemplates({
          page: page.currentPage,
          limit: page.pageSize,
          mer_id: Number.isFinite(merId) && merId! > 0 ? merId : undefined,
          template_name:
            String(formValues?.template_name ?? '').trim() || undefined,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: formOptions.value,
  gridOptions,
});

const [ViewDrawer, viewDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  cancelText: '关闭',
});

const [CreateDrawer, createDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: saveCreate,
});

const [CopyDrawer, copyDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: saveCopy,
});

function toCascaderOptions(nodes: PlatformCategory[]): CascaderOption[] {
  return (nodes || []).map((node) => ({
    label: node.cate_name,
    value: node.store_category_id,
    children: node.children?.length
      ? toCascaderOptions(node.children)
      : undefined,
  }));
}

function flattenCategories(
  nodes: PlatformCategory[],
  prefix = '',
): Array<{ label: string; value: number }> {
  return (nodes || []).flatMap((node) => [
    { label: `${prefix}${node.cate_name}`, value: node.store_category_id },
    ...flattenCategories(node.children || [], `${prefix}— `),
  ]);
}

async function loadOptions() {
  try {
    const [merchants, categories] = await Promise.all([
      fetchPlatformMerchants({ page: 1, limit: 200, status: 1 }),
      listPlatformCategoriesApi(),
    ]);
    merchantOptions.value = (merchants.list || []).map((row) => ({
      label: row.mer_name,
      value: row.mer_id,
    }));
    const list = categories.list || [];
    categoryTree.value = toCascaderOptions(list);
    categoryOptions.value = flattenCategories(list);
    gridApi.setState({ formOptions: formOptions.value });
  } catch {
    merchantOptions.value = [];
    categoryTree.value = [];
    categoryOptions.value = [];
  }
}

function formatParamValues(values: string[] | undefined) {
  return (values || []).join('、') || '—';
}

function cateLabel(cateID: number | undefined) {
  if (!cateID) return '—';
  return (
    categoryOptions.value.find((item) => item.value === cateID)?.label ||
    String(cateID)
  );
}

function emptyParamRow(): ParamRow {
  return { name: '', values: [], required: 0, sort: 0, single: '' };
}

function resetCreateForm() {
  createForm.mer_id = undefined;
  createForm.template_name = '';
  createForm.cate_id = undefined;
  createForm.is_required = 0;
  createForm.sort = 0;
  paramRows.value = [emptyParamRow()];
}

function openCreate() {
  resetCreateForm();
  createDrawerApi.setState({ title: '添加参数模板' }).open();
}

function addParam() {
  paramRows.value.push(emptyParamRow());
}

function removeParam(index: number) {
  paramRows.value.splice(index, 1);
  if (!paramRows.value.length) {
    paramRows.value = [emptyParamRow()];
  }
}

function addParamValue(row: ParamRow) {
  const value = row.single.trim();
  if (!value) return;
  if (row.values.includes(value)) {
    ElMessage.warning('参数值已存在');
    row.single = '';
    return;
  }
  if ([...value].length > 64) {
    ElMessage.warning('参数值最多 64 字');
    return;
  }
  row.values.push(value);
  row.single = '';
}

function removeParamValue(row: ParamRow, index: number) {
  row.values.splice(index, 1);
}

function buildCreatePayload() {
  if (!createForm.mer_id) {
    ElMessage.warning('请选择店铺');
    return null;
  }
  if (!createForm.template_name.trim()) {
    ElMessage.warning('请填写参数模板名称');
    return null;
  }
  if (!createForm.cate_id) {
    ElMessage.warning('请选择平台分类');
    return null;
  }
  const params: StoreParameterItem[] = [];
  for (const row of paramRows.value) {
    const name = row.name.trim();
    if (!name) {
      ElMessage.warning('请输入参数名称');
      return null;
    }
    if (!row.values.length) {
      ElMessage.warning(`请为「${name}」填写参数值（回车确认）`);
      return null;
    }
    params.push({
      name,
      values: [...row.values],
      required: Number(row.required) === 1 ? 1 : 0,
      sort: Number(row.sort) || 0,
    });
  }
  if (!params.length) {
    ElMessage.warning('请至少添加一个参数');
    return null;
  }
  return {
    mer_id: createForm.mer_id,
    template_name: createForm.template_name.trim(),
    cate_id: createForm.cate_id,
    is_required: createForm.is_required,
    sort: createForm.sort,
    params,
  };
}

async function saveCreate() {
  const payload = buildCreatePayload();
  if (!payload) return;
  saving.value = true;
  createDrawerApi.lock();
  try {
    await createStoreParameterTemplate(payload);
    createDrawerApi.close();
    ElMessage.success('保存成功');
    gridApi.reload();
  } finally {
    saving.value = false;
    createDrawerApi.unlock();
  }
}

async function openView(row: StoreParameterTemplateRow) {
  const detail = await fetchStoreParameterTemplate(row.id);
  viewing.value = detail;
  viewDrawerApi.setState({ title: '查看参数模板' }).open();
}

async function openCopy(row: StoreParameterTemplateRow) {
  const detail = await fetchStoreParameterTemplate(row.id);
  copySourceId.value = detail.id;
  copyForm.template_name = `${detail.template_name}（副本）`.slice(0, 64);
  copyForm.sort = detail.sort;
  copyForm.status = 1;
  copyForm.cate_ids = detail.cate_id ? [Number(detail.cate_id)] : [];
  copyForm.params = (detail.params || []).map((item) => ({
    name: item.name,
    values: [...(item.values || [])],
    required: Number(item.required) ? 1 : 0,
    sort: item.sort,
  }));
  if (!copyForm.params.length) {
    copyForm.params.push({ name: '', values: [], required: 0, sort: 0 });
  }
  copyDrawerApi.setState({ title: '复制参数模板' }).open();
}

function addParamRow() {
  copyForm.params.push({ name: '', values: [], required: 0, sort: 0 });
}

function removeParamRow(index: number) {
  copyForm.params.splice(index, 1);
}

function valuesText(item: StoreParameterItem) {
  return (item.values || []).join('、');
}

function setValuesText(item: StoreParameterItem, text: string) {
  item.values = text
    .split(/[、,，\n]/)
    .map((v) => v.trim())
    .filter(Boolean);
}

async function saveCopy() {
  if (!copySourceId.value) return;
  const name = copyForm.template_name.trim();
  if (!name) {
    ElMessage.warning('请填写参数模板名称');
    return;
  }
  if (!copyForm.cate_ids.length) {
    ElMessage.warning('请选择关联分类');
    return;
  }
  const params = copyForm.params
    .map((item) => ({
      name: item.name.trim(),
      values: (item.values || []).map((v) => String(v).trim()).filter(Boolean),
      required: Number(item.required) ? 1 : 0,
      sort: Number(item.sort) || 0,
    }))
    .filter((item) => item.name && item.values.length);
  if (!params.length) {
    ElMessage.warning('请至少填写一个完整参数项');
    return;
  }
  copyDrawerApi.lock();
  try {
    await copyStoreParameterTemplate(copySourceId.value, {
      template_name: name,
      cate_ids: copyForm.cate_ids.map(Number),
      params,
      sort: copyForm.sort,
      status: copyForm.status,
    });
    copyDrawerApi.close();
    ElMessage.success('已复制到平台商品参数');
    gridApi.reload();
  } finally {
    copyDrawerApi.unlock();
  }
}

onMounted(() => {
  void loadOptions();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          添加参数模板
        </ElButton>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openView(row)">查看</ElButton>
        <ElButton link type="primary" @click="openCopy(row)">复制</ElButton>
      </template>
    </Grid>

    <CreateDrawer>
      <ElForm label-width="120px" class="pr-2">
        <ElFormItem label="店铺名称" required>
          <ElSelect
            v-model="createForm.mer_id"
            class="w-full"
            filterable
            clearable
            placeholder="请选择店铺"
          >
            <ElOption
              v-for="item in merchantOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
        </ElFormItem>

        <ElFormItem label="参数模板名称" required>
          <ElInput
            v-model="createForm.template_name"
            maxlength="64"
            placeholder="请输入参数模板名称"
            clearable
          />
          <div class="mt-1 text-xs text-[var(--el-text-color-secondary)]">
            请输入参数模板的名称，用于标识和管理不同分类的商品参数组
          </div>
        </ElFormItem>

        <ElFormItem label="平台分类" required>
          <ElSelect
            v-model="createForm.cate_id"
            class="w-full"
            filterable
            clearable
            placeholder="请选择平台分类"
          >
            <ElOption
              v-for="item in categoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
          <div class="mt-1 text-xs text-[var(--el-text-color-secondary)]">
            请选择关联的商品平台分类，该模板将仅对所选分类下的商品生效
          </div>
        </ElFormItem>

        <ElFormItem label="是否必选">
          <ElSwitch
            v-model="createForm.is_required"
            :active-value="1"
            :inactive-value="0"
            active-text="是"
            inactive-text="否"
          />
          <div class="mt-1 text-xs text-[var(--el-text-color-secondary)]">
            开启后，店铺添加商品关联所选分类时，必须关联并使用此参数模板；关闭则为可选关联
          </div>
        </ElFormItem>

        <ElFormItem label="排序">
          <ElInputNumber
            v-model="createForm.sort"
            :min="0"
            controls-position="right"
            class="w-full"
          />
          <div class="mt-1 text-xs text-[var(--el-text-color-secondary)]">
            请输入数字，用于控制参数模板在列表中的展示顺序，数字越大越靠前
          </div>
        </ElFormItem>

        <ElFormItem label-width="0">
          <ElTable :data="paramRows" border size="small" class="w-full">
            <ElTableColumn label="参数名称" min-width="140" align="center">
              <template #default="{ row }">
                <ElInput
                  v-model="row.name"
                  size="small"
                  maxlength="32"
                  placeholder="请输入参数名称"
                />
              </template>
            </ElTableColumn>
            <ElTableColumn label="参数值" min-width="260">
              <template #default="{ row }">
                <div class="flex flex-wrap items-center gap-1">
                  <ElTag
                    v-for="(value, index) in row.values"
                    :key="`${value}-${index}`"
                    class="mb-1"
                    closable
                    size="small"
                    @close="removeParamValue(row, index)"
                  >
                    {{ value }}
                  </ElTag>
                  <ElInput
                    v-model="row.single"
                    class="min-w-[140px] flex-1"
                    size="small"
                    maxlength="64"
                    placeholder="请输入选项，回车确认"
                    @keyup.enter="addParamValue(row)"
                    @blur="addParamValue(row)"
                  />
                </div>
              </template>
            </ElTableColumn>
            <ElTableColumn label="必填" width="140" align="center">
              <template #default="{ row }">
                <ElSwitch
                  v-model="row.required"
                  :active-value="1"
                  :inactive-value="0"
                  active-text="开启"
                  inactive-text="关闭"
                />
              </template>
            </ElTableColumn>
            <ElTableColumn label="排序" width="120" align="center">
              <template #default="{ row }">
                <ElInputNumber
                  v-model="row.sort"
                  :min="0"
                  :controls="false"
                  size="small"
                  class="w-full"
                />
              </template>
            </ElTableColumn>
            <ElTableColumn label="操作" width="80" align="center">
              <template #default="{ $index }">
                <ElButton link type="primary" @click="removeParam($index)">
                  删除
                </ElButton>
              </template>
            </ElTableColumn>
          </ElTable>
          <div class="mt-3">
            <ElButton @click="addParam">添加参数</ElButton>
          </div>
        </ElFormItem>
      </ElForm>
    </CreateDrawer>

    <ViewDrawer>
      <ElForm v-if="viewing" label-width="100px">
        <ElFormItem label="店铺名称">
          <span>{{ viewing.mer_name || '—' }}</span>
        </ElFormItem>
        <ElFormItem label="模板名称">
          <span>{{ viewing.template_name }}</span>
        </ElFormItem>
        <ElFormItem label="关联分类">
          <span>{{ cateLabel(viewing.cate_id) }}</span>
        </ElFormItem>
        <ElFormItem label="是否必选">
          <span>{{ Number(viewing.is_required) === 1 ? '是' : '否' }}</span>
        </ElFormItem>
        <ElFormItem label="排序">
          <span>{{ viewing.sort }}</span>
        </ElFormItem>
        <ElFormItem label="创建时间">
          <span>{{ formatShanghaiDateTime(viewing.created_at) }}</span>
        </ElFormItem>
        <ElFormItem label="参数项">
          <ElTable :data="viewing.params || []" border class="w-full">
            <ElTableColumn label="参数名称" min-width="120" prop="name" />
            <ElTableColumn label="参数值" min-width="220">
              <template #default="{ row }">
                {{ formatParamValues(row.values) }}
              </template>
            </ElTableColumn>
            <ElTableColumn label="必填" width="90">
              <template #default="{ row }">
                <ElTag :type="row.required ? 'warning' : 'info'" size="small">
                  {{ row.required ? '是' : '否' }}
                </ElTag>
              </template>
            </ElTableColumn>
            <ElTableColumn label="排序" prop="sort" width="80" />
          </ElTable>
        </ElFormItem>
      </ElForm>
    </ViewDrawer>

    <CopyDrawer>
      <ElForm label-width="100px">
        <ElFormItem label="模板名称" required>
          <ElInput
            v-model="copyForm.template_name"
            maxlength="64"
            show-word-limit
            placeholder="请输入参数模板名称"
          />
        </ElFormItem>
        <ElFormItem label="关联分类" required>
          <ElCascader
            v-model="copyForm.cate_ids"
            class="w-full"
            clearable
            filterable
            :options="categoryTree"
            :props="{
              multiple: true,
              checkStrictly: true,
              emitPath: false,
              value: 'value',
              label: 'label',
              children: 'children',
            }"
            placeholder="请选择平台分类"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="copyForm.sort" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="启用">
          <ElSwitch
            v-model="copyForm.status"
            :active-value="1"
            :inactive-value="0"
          />
        </ElFormItem>
        <ElFormItem label="参数项" required>
          <div class="w-full space-y-3">
            <div
              v-for="(item, index) in copyForm.params"
              :key="index"
              class="rounded border border-gray-200 p-3"
            >
              <div class="mb-2 flex items-center justify-between gap-2">
                <span class="text-sm text-gray-500">参数 {{ index + 1 }}</span>
                <ElButton
                  link
                  type="danger"
                  :disabled="copyForm.params.length <= 1"
                  @click="removeParamRow(index)"
                >
                  删除
                </ElButton>
              </div>
              <ElFormItem label="参数名称" label-width="80px" class="mb-2">
                <ElInput v-model="item.name" maxlength="32" placeholder="如：面料" />
              </ElFormItem>
              <ElFormItem label="参数值" label-width="80px" class="mb-2">
                <ElInput
                  :model-value="valuesText(item)"
                  type="textarea"
                  :rows="2"
                  placeholder="用顿号、逗号或换行分隔"
                  @update:model-value="(v) => setValuesText(item, String(v))"
                />
              </ElFormItem>
              <div class="flex flex-wrap items-center gap-4">
                <ElFormItem label="必填" label-width="80px" class="mb-0">
                  <ElSwitch
                    v-model="item.required"
                    :active-value="1"
                    :inactive-value="0"
                  />
                </ElFormItem>
                <ElFormItem label="排序" label-width="80px" class="mb-0">
                  <ElInputNumber v-model="item.sort" :min="0" />
                </ElFormItem>
              </div>
            </div>
            <ElButton :icon="Plus" @click="addParamRow">添加参数</ElButton>
          </div>
        </ElFormItem>
      </ElForm>
    </CopyDrawer>
  </Page>
</template>
