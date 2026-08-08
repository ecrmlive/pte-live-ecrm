<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteProductGuarantee,
  fetchProductGuarantees,
  saveProductGuarantee,
  updateProductGuaranteeStatus,
  type ProductGuaranteeRow,
} from '#/api/core/ecrm';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const editing = ref<number>();
const form = reactive({
  name: '',
  content: '',
  icon_url: '',
  sort: 0,
  status: 1,
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    ...LIST_DATE_RANGE_FIELD,
    label: '时间选择',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入服务条款搜索',
    },
    fieldName: 'keyword',
    label: '服务条款',
  },
]);

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

const gridOptions: VxeGridProps<ProductGuaranteeRow> = {
  columns: [
    { type: 'seq', title: '序号', width: 70 },
    {
      field: 'name',
      minWidth: 140,
      showOverflow: false,
      title: '服务条款',
    },
    {
      field: 'icon_url',
      slots: { default: 'icon' },
      title: '服务条款图标',
      width: 120,
    },
    {
      className: 'col--remark',
      field: 'content',
      minWidth: 260,
      showOverflow: 'tooltip',
      title: '服务内容描述',
      width: 360,
    },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 110,
    },
    { field: 'mer_count', title: '使用的店铺数', width: 120 },
    { field: 'product_count', title: '使用商品数', width: 110 },
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
        const data = await fetchProductGuarantees(
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    name: '',
    content: '',
    icon_url: '',
    sort: 0,
    status: 1,
  });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加保障服务' }).open();
}

function openEdit(row: ProductGuaranteeRow) {
  editing.value = row.id;
  Object.assign(form, {
    name: row.name,
    content: row.content,
    icon_url: row.icon_url || '',
    sort: row.sort,
    status: row.status,
  });
  formDrawerApi.setState({ title: '编辑保障服务' }).open();
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写服务条款');
    return;
  }
  if (!form.content.trim()) {
    ElMessage.warning('请填写服务内容描述');
    return;
  }
  if (!form.icon_url.trim()) {
    ElMessage.warning('请选择服务条款图标');
    return;
  }
  formDrawerApi.lock();
  try {
    await saveProductGuarantee(editing.value, {
      name: form.name.trim(),
      content: form.content.trim(),
      icon_url: form.icon_url.trim(),
      sort: form.sort,
      status: form.status,
    });
    formDrawerApi.close();
    ElMessage.success('保障服务已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeShow(row: ProductGuaranteeRow, enabled: boolean) {
  const before = row.status === 1;
  row.status = enabled ? 1 : 0;
  try {
    await updateProductGuaranteeStatus(row.id, enabled ? 1 : 0);
  } catch {
    row.status = before ? 1 : 0;
  }
}

async function remove(row: ProductGuaranteeRow) {
  try {
    await ElMessageBox.confirm(
      `确认删除保障服务“${row.name}”？删除后不可恢复。`,
      '删除确认',
      { type: 'warning' },
    );
    await deleteProductGuarantee(row.id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          添加保障服务
        </ElButton>
      </template>
      <template #icon="{ row }">
        <ElImage
          v-if="row.icon_url"
          :src="resolveCosMediaUrl(row.icon_url)"
          :preview-src-list="[resolveCosMediaUrl(row.icon_url)]"
          fit="contain"
          class="h-10 w-10 rounded"
          preview-teleported
        />
        <span v-else class="text-gray-400">—</span>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeShow(row, Boolean(enabled))
          "
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="primary" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElForm label-width="110px">
        <ElFormItem label="服务条款" required>
          <ElInput
            v-model="form.name"
            maxlength="64"
            show-word-limit
            placeholder="请输入服务条款"
          />
        </ElFormItem>
        <ElFormItem label="内容描述" required>
          <ElInput
            v-model="form.content"
            :rows="5"
            maxlength="1000"
            show-word-limit
            type="textarea"
            placeholder="请输入内容描述"
          />
        </ElFormItem>
        <ElFormItem label="条款图标" required>
          <div>
            <ImageField
              v-model="form.icon_url"
              default-library="system"
              :preview-size="80"
            />
            <div class="mt-1 text-xs text-gray-400">建议尺寸：100×100px</div>
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
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" :max="99999" class="w-full" />
        </ElFormItem>
      </ElForm>
    </FormDrawer>
  </Page>
</template>
