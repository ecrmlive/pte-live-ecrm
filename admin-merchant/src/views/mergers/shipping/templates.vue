<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElOption,
  ElSelect,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createShippingTemplate,
  deleteShippingTemplate,
  fetchShippingTemplates,
  getShippingTemplate,
  updateShippingTemplate,
  type ShippingTemplate,
} from '#/api/core/mergers';

const form = reactive({
  name: '',
  type: 3,
  sort: 0,
  first: 1,
  first_price: 0,
  continue: 1,
  continue_price: 0,
});
const editingId = ref(0);

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => {
    if (!form.name.trim()) {
      ElMessage.warning('请填写名称');
      return;
    }
    const payload = {
      name: form.name.trim(),
      type: form.type,
      sort: form.sort,
      regions: [
        {
          city_ids: '',
          first: form.first,
          first_price: form.first_price,
          continue: form.continue,
          continue_price: form.continue_price,
        },
      ],
    };
    if (editingId.value) await updateShippingTemplate(editingId.value, payload);
    else await createShippingTemplate(payload);
    ElMessage.success('已保存');
    formModalApi.close();
    gridApi.reload();
  },
});

const gridOptions: VxeGridProps<ShippingTemplate> = {
  border: true,
  columns: [
    { field: 'template_id', title: 'ID', width: 80 },
    { field: 'name', minWidth: 160, title: '名称' },
    {
      field: 'type',
      title: '计费',
      width: 100,
      formatter: ({ cellValue }) =>
        ({ 1: '按件数', 2: '按重量', 3: '包邮' } as Record<number, string>)[
          cellValue
        ] || String(cellValue),
    },
    { field: 'sort', title: '排序', width: 80 },
    { fixed: 'right', slots: { default: 'action' }, title: '操作', width: 160 },
  ],
  height: 'auto',
  pagerConfig: { enabled: true, pageSize: 20 },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const data = await fetchShippingTemplates({
          page: page.currentPage,
          limit: page.pageSize,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'template_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function openCreate() {
  editingId.value = 0;
  Object.assign(form, {
    name: '',
    type: 3,
    sort: 0,
    first: 1,
    first_price: 0,
    continue: 1,
    continue_price: 0,
  });
  formModalApi.setState({ title: '新建运费模板' });
  formModalApi.open();
}

async function openEdit(row: ShippingTemplate) {
  editingId.value = row.template_id;
  const data = await getShippingTemplate(row.template_id);
  const r = data.regions?.[0];
  Object.assign(form, {
    name: data.name,
    type: data.type,
    sort: data.sort,
    first: r?.first ?? 1,
    first_price: r?.first_price ?? 0,
    continue: r?.continue ?? 1,
    continue_price: r?.continue_price ?? 0,
  });
  formModalApi.setState({ title: '编辑运费模板' });
  formModalApi.open();
}

async function onDelete(row: ShippingTemplate) {
  try {
    await ElMessageBox.confirm(`删除模板「${row.name}」？`, '提示', {
      type: 'warning',
    });
  } catch {
    return;
  }
  await deleteShippingTemplate(row.template_id);
  ElMessage.success('已删除');
  gridApi.reload();
}
</script>

<template>
  <Page auto-content-height title="运费模板">
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">新建模板</ElButton>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="onDelete(row)">删除</ElButton>
      </template>
    </Grid>
    <FormModal>
      <ElForm label-position="top">
        <ElFormItem label="名称" required>
          <ElInput v-model="form.name" />
        </ElFormItem>
        <ElFormItem label="计费方式">
          <ElSelect v-model="form.type" class="w-full">
            <ElOption :value="1" label="按件数" />
            <ElOption :value="2" label="按重量" />
            <ElOption :value="3" label="包邮" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="首件/首重">
          <ElInputNumber v-model="form.first" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="首费">
          <ElInputNumber v-model="form.first_price" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="续件/续重">
          <ElInputNumber v-model="form.continue" :min="0" class="w-full" />
        </ElFormItem>
        <ElFormItem label="续费">
          <ElInputNumber v-model="form.continue_price" :min="0" class="w-full" />
        </ElFormItem>
      </ElForm>
    </FormModal>
  </Page>
</template>
