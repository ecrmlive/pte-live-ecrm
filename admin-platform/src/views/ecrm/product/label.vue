<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
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
import {
  deleteProductLabel,
  fetchProductLabels,
  saveProductLabel,
  updateProductLabelStatus,
  type ProductLabelRow,
} from '#/api/core/ecrm';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const editing = ref<ProductLabelRow>();
const form = reactive({
  name: '',
  description: '',
  sort: 0,
  status: 1,
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('标签名称'),
]);

const gridOptions: VxeGridProps<ProductLabelRow> = {
  columns: [
    {
      field: 'name',
      minWidth: 140,
      showOverflow: false,
      title: '标签名称',
    },
    {
      field: 'description',
      minWidth: 220,
      showOverflow: false,
      title: '标签说明',
    },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 120,
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
        const list = (await fetchProductLabels()).list || [];
        const keyword = String(formValues?.keyword ?? '')
          .trim()
          .toLowerCase();
        const filtered = keyword
          ? list.filter(
              (row) =>
                row.name.toLowerCase().includes(keyword) ||
                String(row.description || '')
                  .toLowerCase()
                  .includes(keyword),
            )
          : list;
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: filtered.slice(start, start + page.pageSize),
          total: filtered.length,
        };
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
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    name: '',
    description: '',
    sort: 0,
    status: 1,
  });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '添加商品标签' }).open();
}

function openEdit(row: ProductLabelRow) {
  editing.value = row;
  Object.assign(form, {
    name: row.name,
    description: row.description || '',
    sort: row.sort,
    status: row.status,
  });
  formDrawerApi.setState({ title: '编辑商品标签' }).open();
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写标签名称');
    return;
  }
  formDrawerApi.lock();
  try {
    await saveProductLabel(editing.value?.id, {
      name: form.name.trim(),
      description: form.description.trim(),
      color: editing.value?.color || '',
      sort: form.sort,
      status: form.status,
    });
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: ProductLabelRow, enabled: boolean) {
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await updateProductLabelStatus(row.id, enabled ? 1 : 0);
  } catch {
    row.status = before;
  }
}

async function remove(row: ProductLabelRow) {
  try {
    await ElMessageBox.confirm(
      `删除商品标签“${row.name}”后不可恢复，是否继续？`,
      '删除确认',
      { type: 'warning' },
    );
    await deleteProductLabel(row.id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 取消 */
  }
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          添加商品标签
        </ElButton>
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
      <ElForm label-width="88px">
        <ElFormItem label="标签名称" required>
          <ElInput
            v-model="form.name"
            maxlength="64"
            show-word-limit
            placeholder="请输入标签名称"
          />
        </ElFormItem>
        <ElFormItem label="标签说明">
          <ElInput
            v-model="form.description"
            :rows="3"
            maxlength="255"
            show-word-limit
            type="textarea"
            placeholder="请输入标签说明"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" :min="0" class="w-full" />
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
      </ElForm>
    </FormDrawer>
  </Page>
</template>
