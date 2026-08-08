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
  ElMessageBox,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteProductParameterTemplate,
  fetchProductParameterTemplate,
  fetchProductParameterTemplates,
  saveProductParameterTemplate,
  type ProductParameterItem,
  type ProductParameterTemplateRow,
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

type DrawerMode = 'create' | 'edit' | 'view';
type CascaderOption = {
  label: string;
  value: number;
  children?: CascaderOption[];
};

const drawerMode = ref<DrawerMode>('create');
const editingId = ref(0);
const categoryOptions = ref<{ label: string; value: number }[]>([]);
const categoryTree = ref<CascaderOption[]>([]);
const paramDraft = ref('');

const form = reactive({
  name: '',
  cate_ids: [] as number[],
  params: [] as ProductParameterItem[],
  sort: 0,
  status: 1,
});

const isReadonly = computed(() => drawerMode.value === 'view');
const drawerTitle = computed(() => {
  if (drawerMode.value === 'create') return '添加参数模板';
  if (drawerMode.value === 'edit') return '编辑参数模板';
  return '查看参数模板';
});

const cascaderProps = {
  multiple: true,
  checkStrictly: false,
  emitPath: false,
  value: 'value',
  label: 'label',
  children: 'children',
};

function toCascaderOptions(rows: PlatformCategory[] = []): CascaderOption[] {
  const out: CascaderOption[] = [];
  for (const row of rows) {
    const children = toCascaderOptions(row.children || []);
    const option: CascaderOption = {
      label: row.cate_name,
      value: Number(row.store_category_id),
    };
    if (children.length) option.children = children;
    out.push(option);
  }
  return out;
}

function flattenCategories(
  rows: PlatformCategory[] = [],
  acc: { label: string; value: number }[] = [],
) {
  for (const row of rows) {
    acc.push({
      label: row.cate_name,
      value: Number(row.store_category_id),
    });
    if (row.children?.length) flattenCategories(row.children, acc);
  }
  return acc;
}

const formOptions = computed((): VbenFormProps =>
  listFormOptionsDefaults([
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: categoryOptions.value,
        placeholder: '请选择',
      },
      fieldName: 'cate_id',
      label: '平台分类',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入模板名称',
      },
      fieldName: 'name',
      label: '模板名称',
    },
  ]),
);

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const cateId = Number(formValues?.cate_id || 0);
  return {
    page: page.currentPage,
    limit: page.pageSize,
    name: String(formValues?.name ?? '').trim() || undefined,
    cate_id: cateId > 0 ? cateId : undefined,
  };
}

const gridOptions: VxeGridProps<ProductParameterTemplateRow> = {
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'name',
      minWidth: 140,
      showOverflow: false,
      title: '模板名称',
    },
    {
      field: 'cate_names_text',
      minWidth: 220,
      showOverflow: 'tooltip',
      title: '关联分类',
    },
    {
      field: 'is_required',
      formatter: ({ cellValue }) =>
        Number(cellValue) === 1 ? '必选' : '非必选',
      title: '是否必选',
      width: 100,
    },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 180 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchProductParameterTemplates(
          buildListParams(page, formValues),
        );
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

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function emptyParam(): ProductParameterItem {
  return { name: '', values: [], required: 0, sort: 0 };
}

function resetForm() {
  editingId.value = 0;
  Object.assign(form, {
    name: '',
    cate_ids: [],
    params: [emptyParam()],
    sort: 0,
    status: 1,
  });
  paramDraft.value = '';
}

function applyRow(row: ProductParameterTemplateRow) {
  editingId.value = row.id;
  form.name = row.name || '';
  form.cate_ids = [...(row.cate_ids || [])];
  form.params =
    row.params?.length > 0
      ? row.params.map((item) => ({
          name: item.name || '',
          values: [...(item.values || [])],
          required: Number(item.required) === 1 ? 1 : 0,
          sort: Number(item.sort) || 0,
        }))
      : [emptyParam()];
  form.sort = row.sort || 0;
  form.status = row.status ?? 1;
}

function openCreate() {
  resetForm();
  drawerMode.value = 'create';
  formDrawerApi
    .setState({
      title: drawerTitle.value,
      showConfirmButton: true,
      confirmText: '保存',
    })
    .open();
}

async function openEdit(row: ProductParameterTemplateRow) {
  resetForm();
  drawerMode.value = 'edit';
  try {
    const detail = await fetchProductParameterTemplate(row.id);
    applyRow(detail);
  } catch {
    applyRow(row);
  }
  formDrawerApi
    .setState({
      title: drawerTitle.value,
      showConfirmButton: true,
      confirmText: '保存',
    })
    .open();
}

async function openView(row: ProductParameterTemplateRow) {
  resetForm();
  drawerMode.value = 'view';
  try {
    const detail = await fetchProductParameterTemplate(row.id);
    applyRow(detail);
  } catch {
    applyRow(row);
  }
  formDrawerApi
    .setState({
      title: drawerTitle.value,
      showConfirmButton: false,
    })
    .open();
}

function addParam() {
  form.params.push(emptyParam());
}

function removeParam(index: number) {
  if (form.params.length <= 1) {
    ElMessage.warning('请至少保留一个参数项');
    return;
  }
  form.params.splice(index, 1);
}

function addParamValue(row: ProductParameterItem, raw: string) {
  const value = raw.trim();
  if (!value) return;
  if (row.values.includes(value)) {
    ElMessage.warning('参数值已存在');
    return;
  }
  if ([...row.values, value].join('').length > 0 && value.length > 64) {
    ElMessage.warning('参数值过长');
    return;
  }
  row.values.push(value);
}

function removeParamValue(row: ProductParameterItem, index: number) {
  row.values.splice(index, 1);
}

async function save() {
  if (isReadonly.value) {
    formDrawerApi.close();
    return;
  }
  if (!form.name.trim()) {
    ElMessage.warning('请填写模板名称');
    return;
  }
  if (!form.cate_ids.length) {
    ElMessage.warning('请选择关联分类');
    return;
  }
  const params = form.params
    .map((item) => ({
      name: item.name.trim(),
      values: item.values.map((v) => v.trim()).filter(Boolean),
      required: Number(item.required) === 1 ? 1 : 0,
      sort: Number(item.sort) || 0,
    }))
    .filter((item) => item.name || item.values.length);
  if (!params.length || params.some((item) => !item.name || !item.values.length)) {
    ElMessage.warning('请完善参数名称与参数值');
    return;
  }
  formDrawerApi.lock();
  try {
    await saveProductParameterTemplate(
      drawerMode.value === 'edit' ? editingId.value : undefined,
      {
        name: form.name.trim(),
        cate_ids: [...form.cate_ids],
        params,
        sort: form.sort,
        status: 1,
      },
    );
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function remove(row: ProductParameterTemplateRow) {
  try {
    await ElMessageBox.confirm(
      `确认删除参数模板“${row.name}”？删除后不可恢复。`,
      '删除确认',
      { type: 'warning' },
    );
    await deleteProductParameterTemplate(row.id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

onMounted(async () => {
  try {
    const res = await listPlatformCategoriesApi();
    const list = Array.isArray(res) ? res : res?.list || [];
    categoryTree.value = toCascaderOptions(list);
    categoryOptions.value = flattenCategories(list);
    gridApi.setState({ formOptions: formOptions.value });
  } catch {
    categoryTree.value = [];
    categoryOptions.value = [];
  }
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
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="openView(row)">查看</ElButton>
        <ElButton link type="primary" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="110px">
        <ElFormItem label="参数模板名称" required>
          <ElInput
            v-model="form.name"
            maxlength="64"
            show-word-limit
            placeholder="请输入参数模板名称"
            :disabled="isReadonly"
          />
        </ElFormItem>
        <ElFormItem label="平台分类" required>
          <ElCascader
            v-model="form.cate_ids"
            class="w-full"
            :options="categoryTree"
            :props="cascaderProps"
            clearable
            filterable
            collapse-tags
            collapse-tags-tooltip
            placeholder="请选择"
            :disabled="isReadonly"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber
            v-model="form.sort"
            :min="0"
            :max="99999"
            class="w-full"
            :disabled="isReadonly"
          />
        </ElFormItem>
        <ElFormItem label="参数项" required>
          <div class="w-full">
            <ElTable :data="form.params" border size="small">
              <ElTableColumn label="参数名称" min-width="140">
                <template #default="{ row }">
                  <ElInput
                    v-model="row.name"
                    maxlength="32"
                    placeholder="参数名称"
                    :disabled="isReadonly"
                  />
                </template>
              </ElTableColumn>
              <ElTableColumn label="参数值" min-width="260">
                <template #default="{ row }">
                  <div class="flex flex-wrap items-center gap-1">
                    <ElTag
                      v-for="(value, idx) in row.values"
                      :key="`${value}-${idx}`"
                      :closable="!isReadonly"
                      size="small"
                      @close="removeParamValue(row, idx)"
                    >
                      {{ value }}
                    </ElTag>
                    <ElInput
                      v-if="!isReadonly"
                      class="w-28"
                      size="small"
                      placeholder="回车添加"
                      @keyup.enter="
                        (e: KeyboardEvent) => {
                          const el = e.target as HTMLInputElement;
                          addParamValue(row, el.value);
                          el.value = '';
                        }
                      "
                    />
                  </div>
                </template>
              </ElTableColumn>
              <ElTableColumn label="必填" width="110">
                <template #default="{ row }">
                  <ElSwitch
                    v-model="row.required"
                    :active-value="1"
                    :inactive-value="0"
                    inline-prompt
                    active-text="开启"
                    inactive-text="关闭"
                    :disabled="isReadonly"
                  />
                </template>
              </ElTableColumn>
              <ElTableColumn label="排序" width="120">
                <template #default="{ row }">
                  <ElInputNumber
                    v-model="row.sort"
                    :min="0"
                    :max="99999"
                    size="small"
                    class="w-full"
                    :disabled="isReadonly"
                  />
                </template>
              </ElTableColumn>
              <ElTableColumn v-if="!isReadonly" label="操作" width="80" fixed="right">
                <template #default="{ $index }">
                  <ElButton link type="danger" @click="removeParam($index)">
                    删除
                  </ElButton>
                </template>
              </ElTableColumn>
            </ElTable>
            <ElButton
              v-if="!isReadonly"
              class="mt-3"
              @click="addParam"
            >
              添加参数
            </ElButton>
          </div>
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
