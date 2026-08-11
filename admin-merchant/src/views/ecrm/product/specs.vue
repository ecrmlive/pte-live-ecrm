<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, reactive, ref, shallowRef } from 'vue';

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
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getMerchantProductCategoriesApi,
  type MerchantCategoryNode,
} from '#/api/core/merchant-catalog';
import {
  createProductParameterTemplateApi,
  deleteProductParameterTemplateApi,
  listProductParameterTemplatesApi,
  updateProductParameterTemplateApi,
  type ProductParameterItem,
  type ProductParameterTemplate,
  type ProductParameterTemplateInput,
} from '#/api/core/product-meta';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

type ParamRow = ProductParameterItem & { single: string };

const saving = ref(false);
const editingID = ref<number>();
const categoryOptions = shallowRef<Array<{ label: string; value: number }>>([]);
const paramRows = ref<ParamRow[]>([]);

const form = reactive({
  template_name: '',
  cate_id: undefined as number | undefined,
  is_required: 0 as 0 | 1,
  sort: 0,
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入参数模板名称' },
    fieldName: 'keyword',
    label: '模板名称',
  },
]);

const gridOptions: VxeGridProps<ProductParameterTemplate> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'template_id', title: 'ID', width: 88 },
    {
      field: 'template_name',
      minWidth: 180,
      showOverflow: false,
      title: '参数模板名称',
    },
    {
      field: 'cate_id',
      formatter: ({ cellValue }) => cateLabel(Number(cellValue)),
      minWidth: 160,
      title: '平台分类',
    },
    {
      field: 'is_required',
      formatter: ({ cellValue }) => (Number(cellValue) === 1 ? '是' : '否'),
      title: '是否必选',
      width: 100,
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    merchantListActionColumn({ width: 128 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listProductParameterTemplatesApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
        });
        return { items: data.list, total: data.total };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'template_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [EditDrawer, editDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: save,
});

function flattenCategories(
  nodes: MerchantCategoryNode[],
  prefix = '',
): Array<{ label: string; value: number }> {
  return nodes.flatMap((node) => [
    { label: `${prefix}${node.cate_name}`, value: node.store_category_id },
    ...flattenCategories(node.children || [], `${prefix}— `),
  ]);
}

function cateLabel(cateID: number) {
  if (!cateID) return '—';
  return categoryOptions.value.find((item) => item.value === cateID)?.label || String(cateID);
}

function emptyParamRow(): ParamRow {
  return { name: '', values: [], required: 0, sort: 0, single: '' };
}

function resetForm() {
  editingID.value = undefined;
  form.template_name = '';
  form.cate_id = undefined;
  form.is_required = 0;
  form.sort = 0;
  paramRows.value = [emptyParamRow()];
}

function openCreate() {
  resetForm();
  editDrawerApi.setState({ title: '新增参数模板' }).open();
}

function openEdit(row: ProductParameterTemplate) {
  editingID.value = row.template_id;
  form.template_name = row.template_name;
  form.cate_id = row.cate_id || undefined;
  form.is_required = Number(row.is_required) === 1 ? 1 : 0;
  form.sort = row.sort;
  paramRows.value = (row.params || []).map((item) => ({
    name: item.name,
    values: [...(item.values || [])],
    required: Number(item.required) === 1 ? 1 : 0,
    sort: item.sort || 0,
    single: '',
  }));
  if (!paramRows.value.length) {
    paramRows.value = [emptyParamRow()];
  }
  editDrawerApi.setState({ title: '编辑参数模板' }).open();
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

function buildPayload(): ProductParameterTemplateInput | null {
  if (!form.template_name.trim()) {
    ElMessage.warning('请填写参数模板名称');
    return null;
  }
  if (!form.cate_id) {
    ElMessage.warning('请选择平台分类');
    return null;
  }
  const params: ProductParameterItem[] = [];
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
    template_name: form.template_name.trim(),
    cate_id: form.cate_id,
    is_required: form.is_required,
    sort: form.sort,
    params,
  };
}

async function save() {
  const payload = buildPayload();
  if (!payload) return;
  saving.value = true;
  editDrawerApi.lock();
  try {
    if (editingID.value) {
      await updateProductParameterTemplateApi(editingID.value, payload);
    } else {
      await createProductParameterTemplateApi(payload);
    }
    editDrawerApi.close();
    ElMessage.success('保存成功');
    gridApi.reload();
  } finally {
    saving.value = false;
    editDrawerApi.unlock();
  }
}

async function remove(row: ProductParameterTemplate) {
  try {
    await confirm({
      content: `删除参数模板「${row.template_name}」后不可恢复，是否继续？`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteProductParameterTemplateApi(row.template_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

onMounted(async () => {
  try {
    const result = await getMerchantProductCategoriesApi();
    categoryOptions.value = flattenCategories(result.list || []);
  } catch {
    categoryOptions.value = [];
  }
});
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" :icon="Plus" @click="openCreate">
        新增参数模板
      </ElButton>
    </template>

    <Grid>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <EditDrawer>
      <ElForm label-width="120px" class="pr-2">
        <ElFormItem label="参数模板名称" required>
          <ElInput
            v-model="form.template_name"
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
            v-model="form.cate_id"
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
            v-model="form.is_required"
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
            v-model="form.sort"
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
            <ElButton @click="addParam">新增参数</ElButton>
          </div>
        </ElFormItem>
      </ElForm>
    </EditDrawer>
  </Page>
</template>
