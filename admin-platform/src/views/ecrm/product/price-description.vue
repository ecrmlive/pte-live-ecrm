<script setup lang="ts">
import type { ImageUploadOptions } from '@vben/plugins/tiptap';
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { VbenTiptap } from '@vben/plugins/tiptap';
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
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { uploadAttachmentApi } from '#/api/core/attachment';
import {
  deleteProductPriceRule,
  fetchProductPriceRule,
  fetchProductPriceRules,
  saveProductPriceRule,
  updateProductPriceRuleStatus,
  type ProductPriceRuleRow,
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

type DrawerMode = 'create' | 'edit';
type CascaderOption = {
  label: string;
  value: number;
  children?: CascaderOption[];
};

const drawerMode = ref<DrawerMode>('create');
const editingId = ref(0);
const categoryOptions = ref<{ label: string; value: number }[]>([]);
const categoryTree = ref<CascaderOption[]>([]);
const editorReady = ref(false);

const form = reactive({
  name: '',
  cate_ids: [] as number[],
  content: '',
  sort: 0,
  status: 1,
});

const drawerTitle = computed(() =>
  drawerMode.value === 'create' ? '添加说明' : '编辑说明',
);

const cascaderProps = {
  multiple: true,
  checkStrictly: true,
  emitPath: false,
  value: 'value',
  label: 'label',
  children: 'children',
};

const imageUpload: ImageUploadOptions = {
  accept: 'image/jpeg,image/png,image/gif,image/webp',
  maxSize: 5 * 1024 * 1024,
  upload: async (file) => {
    const row = await uploadAttachmentApi(file);
    return row.attachment_src;
  },
  onUploadError: () => {
    ElMessage.error('图片上传失败');
  },
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

function isEmptyRichText(html: string) {
  const text = html
    .replace(/<[^>]*>/g, '')
    .replace(/&nbsp;/g, ' ')
    .trim();
  return !text;
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
      label: '商品分类',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '显示', value: 1 },
          { label: '隐藏', value: 0 },
        ],
        placeholder: '请选择',
      },
      fieldName: 'status',
      label: '显示状态',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入说明名称',
      },
      fieldName: 'keyword',
      label: '搜索',
    },
  ]),
);

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const cateId = Number(formValues?.cate_id || 0);
  const statusRaw = formValues?.status;
  const status =
    statusRaw === 0 || statusRaw === 1 || statusRaw === '0' || statusRaw === '1'
      ? Number(statusRaw)
      : undefined;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    cate_id: cateId > 0 ? cateId : undefined,
    status,
  };
}

const gridOptions: VxeGridProps<ProductPriceRuleRow> = {
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'name',
      minWidth: 140,
      showOverflow: false,
      title: '名称',
    },
    {
      field: 'cate_names_text',
      minWidth: 220,
      showOverflow: 'tooltip',
      title: '使用分类',
    },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 120,
    },
    {
      field: 'updated_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '更新时间',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchProductPriceRules(
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
  confirmText: '确定',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
  onOpenChange(isOpen) {
    if (!isOpen) {
      editorReady.value = false;
    }
  },
});

function resetForm() {
  editingId.value = 0;
  Object.assign(form, {
    name: '',
    cate_ids: [],
    content: '',
    sort: 0,
    status: 1,
  });
}

function applyRow(row: ProductPriceRuleRow) {
  editingId.value = row.id;
  form.name = row.name || '';
  form.cate_ids = [...(row.cate_ids || [])];
  form.content = row.content || '';
  form.sort = row.sort || 0;
  form.status = row.status ?? 1;
}

function openCreate() {
  resetForm();
  drawerMode.value = 'create';
  editorReady.value = true;
  formDrawerApi
    .setState({
      title: drawerTitle.value,
      showConfirmButton: true,
      confirmText: '确定',
    })
    .open();
}

async function openEdit(row: ProductPriceRuleRow) {
  resetForm();
  drawerMode.value = 'edit';
  editorReady.value = false;
  try {
    const detail = await fetchProductPriceRule(row.id);
    applyRow(detail);
  } catch {
    applyRow(row);
  }
  editorReady.value = true;
  formDrawerApi
    .setState({
      title: drawerTitle.value,
      showConfirmButton: true,
      confirmText: '确定',
    })
    .open();
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写名称');
    return;
  }
  if (isEmptyRichText(form.content)) {
    ElMessage.warning('请填写价格说明详情');
    return;
  }
  formDrawerApi.lock();
  try {
    await saveProductPriceRule(
      drawerMode.value === 'edit' ? editingId.value : undefined,
      {
        name: form.name.trim(),
        cate_ids: [...form.cate_ids],
        content: form.content,
        sort: form.sort,
        status: form.status,
      },
    );
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: ProductPriceRuleRow, enabled: boolean) {
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await updateProductPriceRuleStatus(row.id, enabled ? 1 : 0);
  } catch {
    row.status = before;
  }
}

async function remove(row: ProductPriceRuleRow) {
  try {
    await ElMessageBox.confirm(
      `确认删除价格说明“${row.name}”？删除后不可恢复。`,
      '删除确认',
      { type: 'warning' },
    );
    await deleteProductPriceRule(row.id);
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
        <div class="toolbar-row">
          <ElButton :icon="Plus" type="primary" @click="openCreate">
            添加价格说明
          </ElButton>
          <span class="toolbar-hint">
            填写价格说明，明确优惠、规格、运费等差异，避免误解与售后纠纷
          </span>
        </div>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeStatus(row, Boolean(enabled))
          "
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="120px">
        <ElFormItem label="名称" required>
          <ElInput
            v-model="form.name"
            maxlength="64"
            show-word-limit
            placeholder="请输入名称"
          />
        </ElFormItem>
        <ElFormItem label="使用商品分类">
          <div class="w-full">
            <ElCascader
              v-model="form.cate_ids"
              class="w-full"
              :options="categoryTree"
              :props="cascaderProps"
              clearable
              filterable
              placeholder="请选择"
            />
            <div class="field-hint">
              注：当不选择任何分类时，默认全部商品
            </div>
          </div>
        </ElFormItem>
        <ElFormItem label="是否显示">
          <ElSwitch
            v-model="form.status"
            :active-value="1"
            :inactive-value="0"
            inline-prompt
            active-text="显示"
            inactive-text="隐藏"
          />
        </ElFormItem>
        <ElFormItem label="排序" required>
          <ElInputNumber v-model="form.sort" :min="0" :max="99999" />
        </ElFormItem>
        <ElFormItem label="价格说明详情" required>
          <div class="w-full">
            <VbenTiptap
              v-if="editorReady"
              v-model="form.content"
              :image-upload="imageUpload"
              :max-height="420"
              :min-height="280"
              :previewable="false"
              placeholder="请输入价格说明详情…"
            />
          </div>
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.toolbar-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.toolbar-hint {
  color: hsl(var(--foreground) / 65%);
  font-size: 13px;
  line-height: 1.4;
}

.field-hint {
  margin-top: 6px;
  color: hsl(var(--foreground) / 55%);
  font-size: 12px;
  line-height: 1.4;
}
</style>
