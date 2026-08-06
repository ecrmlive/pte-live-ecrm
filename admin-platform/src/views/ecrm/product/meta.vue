<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page, useVbenModal } from '@vben/common-ui';
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
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';
import type { AttachmentItem } from '#/api/core/attachment';
import {
  deleteProductGuarantee,
  deleteProductLabel,
  deleteProductParameterTemplate,
  fetchProductGuarantees,
  fetchProductLabels,
  fetchProductParameterTemplates,
  saveProductGuarantee,
  saveProductLabel,
  saveProductParameterTemplate,
  type ProductGuaranteeRow,
  type ProductLabelRow,
  type ProductParameterTemplateRow,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import {
  LIST_ENABLE_STATUS_FIELD,
  LIST_KEYWORD_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

type Kind = 'guarantee' | 'label' | 'parameter';
type MetaRow = ProductLabelRow | ProductGuaranteeRow | ProductParameterTemplateRow;

const route = useRoute();
const kind = computed<Kind>(() =>
  route.path.includes('guarantee')
    ? 'guarantee'
    : route.path.includes('spec')
      ? 'parameter'
      : 'label',
);
const pageTitle = computed(
  () =>
    ({
      label: '商品标签',
      guarantee: '保障服务',
      parameter: '平台商品参数',
    })[kind.value],
);

const editing = ref<number>();
const form = reactive({
  name: '',
  description: '',
  color: '',
  content: '',
  icon_url: '',
  values: '',
  sort: 0,
  status: 1,
});
const iconPicker = ref(false);

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_KEYWORD_FIELD('名称'),
  LIST_ENABLE_STATUS_FIELD('状态'),
]);

function parameterValues(raw: string) {
  try {
    const values = JSON.parse(raw);
    return Array.isArray(values) ? values.map((value) => String(value)) : [];
  } catch {
    return [];
  }
}

function filterRows<T extends { name: string; status: number }>(
  list: T[],
  formValues?: Record<string, unknown>,
) {
  const keyword = String(formValues?.keyword ?? '')
    .trim()
    .toLowerCase();
  const statusRaw = formValues?.status;
  let rows = list;
  if (keyword) {
    rows = rows.filter((row) => row.name.toLowerCase().includes(keyword));
  }
  if (statusRaw === 0 || statusRaw === 1) {
    rows = rows.filter((row) => Number(row.status) === Number(statusRaw));
  }
  return rows;
}

const gridOptions = computed((): VxeGridProps<MetaRow> => {
  const base = {
    pagerConfig: { enabled: false },
    proxyConfig: {
      ajax: {
        query: async (_ctx: unknown, formValues?: Record<string, unknown>) => {
          if (kind.value === 'label') {
            const list = (await fetchProductLabels()).list || [];
            const items = filterRows(list, formValues);
            return { items, total: items.length };
          }
          if (kind.value === 'guarantee') {
            const list = (await fetchProductGuarantees()).list || [];
            const items = filterRows(list, formValues);
            return { items, total: items.length };
          }
          const list = (await fetchProductParameterTemplates()).list || [];
          const items = filterRows(list, formValues);
          return { items, total: items.length };
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

  if (kind.value === 'label') {
    return {
      ...base,
      columns: [
        { field: 'id', title: 'ID', width: 90 },
        { field: 'name', minWidth: 140, title: '名称' },
        {
          field: 'description',
          minWidth: 240,
          showOverflow: false,
          title: '说明',
        },
        {
          field: 'color',
          slots: { default: 'color' },
          title: '颜色',
          width: 120,
        },
        { field: 'sort', title: '排序', width: 90 },
        {
          field: 'status',
          slots: { default: 'status' },
          title: '状态',
          width: 90,
        },
        platformListActionColumn({ width: 150 }),
      ],
    };
  }

  if (kind.value === 'guarantee') {
    return {
      ...base,
      columns: [
        { field: 'id', title: 'ID', width: 90 },
        { field: 'name', minWidth: 140, title: '名称' },
        {
          field: 'content',
          minWidth: 320,
          showOverflow: false,
          title: '保障说明',
        },
        { field: 'sort', title: '排序', width: 90 },
        {
          field: 'status',
          slots: { default: 'status' },
          title: '状态',
          width: 90,
        },
        platformListActionColumn({ width: 150 }),
      ],
    };
  }

  return {
    ...base,
    columns: [
      { field: 'id', title: 'ID', width: 90 },
      { field: 'name', minWidth: 160, title: '参数名称' },
      {
        field: 'values_json',
        minWidth: 280,
        showOverflow: false,
        slots: { default: 'values' },
        title: '候选值',
      },
      { field: 'sort', title: '排序', width: 90 },
      {
        field: 'status',
        slots: { default: 'status' },
        title: '状态',
        width: 90,
      },
      platformListActionColumn({ width: 150 }),
    ],
  };
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions: gridOptions.value,
});

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => save(),
});

function resetForm() {
  editing.value = undefined;
  Object.assign(form, {
    name: '',
    description: '',
    color: '',
    content: '',
    icon_url: '',
    values: '',
    sort: 0,
    status: 1,
  });
}

function openCreate() {
  resetForm();
  formModalApi.setState({ title: `新增${pageTitle.value}` }).open();
}

function openEdit(row: MetaRow) {
  resetForm();
  editing.value = row.id;
  form.name = row.name;
  form.sort = row.sort;
  form.status = row.status;
  if ('description' in row) {
    form.description = row.description;
    form.color = row.color;
  }
  if ('content' in row) {
    form.content = row.content;
    form.icon_url = row.icon_url;
  }
  if ('values_json' in row) {
    try {
      form.values = (JSON.parse(row.values_json) as string[]).join('、');
    } catch {
      form.values = '';
    }
  }
  formModalApi.setState({ title: `编辑${pageTitle.value}` }).open();
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写名称');
    return;
  }
  formModalApi.lock();
  try {
    if (kind.value === 'label') {
      await saveProductLabel(editing.value, {
        name: form.name.trim(),
        description: form.description.trim(),
        color: form.color.trim(),
        sort: form.sort,
        status: form.status,
      });
    } else if (kind.value === 'guarantee') {
      await saveProductGuarantee(editing.value, {
        name: form.name.trim(),
        content: form.content.trim(),
        icon_url: form.icon_url.trim(),
        sort: form.sort,
        status: form.status,
      });
    } else {
      const values = form.values
        .split(/[、,，\n]/)
        .map((value) => value.trim())
        .filter(Boolean);
      if (!values.length) {
        ElMessage.warning('请至少填写一个参数值');
        return;
      }
      await saveProductParameterTemplate(editing.value, {
        name: form.name.trim(),
        values,
        sort: form.sort,
        status: form.status,
      });
    }
    formModalApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formModalApi.unlock();
  }
}

async function remove(id: number, name: string) {
  try {
    await ElMessageBox.confirm(
      `确认删除“${name}”？已关联的商品仍保留其历史快照。`,
      '删除确认',
      { type: 'warning' },
    );
    if (kind.value === 'label') await deleteProductLabel(id);
    else if (kind.value === 'guarantee') await deleteProductGuarantee(id);
    else await deleteProductParameterTemplate(id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

function onIconSelect(items: AttachmentItem[]) {
  form.icon_url = items[0]?.attachment_src || '';
}

watch(kind, () => {
  gridApi.setGridOptions(gridOptions.value);
  gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增{{ pageTitle }}
        </ElButton>
      </template>
      <template #color="{ row }">
        <span :style="{ color: row.color || undefined }">
          {{ row.color || '默认' }}
        </span>
      </template>
      <template #values="{ row }">
        <ElTag
          v-for="value in parameterValues(row.values_json)"
          :key="value"
          class="mr-1"
        >
          {{ value }}
        </ElTag>
        <span
          v-if="!parameterValues(row.values_json).length"
          class="text-gray-400"
        >
          参数值异常
        </span>
      </template>
      <template #status="{ row }">
        <ElTag :type="row.status ? 'success' : 'info'">
          {{ row.status ? '启用' : '停用' }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row.id, row.name)">
          删除
        </ElButton>
      </template>
    </Grid>

    <FormModal>
      <ElForm label-width="92px">
        <ElFormItem label="名称" required>
          <ElInput v-model="form.name" maxlength="64" show-word-limit />
        </ElFormItem>
        <template v-if="kind === 'label'">
          <ElFormItem label="说明">
            <ElInput v-model="form.description" maxlength="255" />
          </ElFormItem>
          <ElFormItem label="颜色">
            <ElInput
              v-model="form.color"
              maxlength="32"
              placeholder="#2563eb（可选）"
            />
          </ElFormItem>
        </template>
        <template v-else-if="kind === 'guarantee'">
          <ElFormItem label="保障说明">
            <ElInput
              v-model="form.content"
              :rows="4"
              maxlength="1000"
              show-word-limit
              type="textarea"
            />
          </ElFormItem>
          <ElFormItem label="保障图标">
            <ElButton @click="iconPicker = true">选择平台图片</ElButton>
            <ElImage
              v-if="form.icon_url"
              :src="form.icon_url"
              class="ml-2 h-8 w-8"
              fit="cover"
            />
            <ElButton
              v-if="form.icon_url"
              link
              type="danger"
              @click="form.icon_url = ''"
            >
              清除
            </ElButton>
          </ElFormItem>
        </template>
        <ElFormItem v-else label="参数值" required>
          <ElInput
            v-model="form.values"
            :rows="3"
            placeholder="用中文顿号、逗号或换行分隔，例如：小杯、中杯、大杯"
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber v-model="form.sort" class="w-full" />
        </ElFormItem>
        <ElFormItem label="启用">
          <ElSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </ElFormItem>
      </ElForm>
    </FormModal>

    <ImagePickerDialog
      v-model:open="iconPicker"
      default-library="system"
      kind="image"
      :limit="1"
      @select="onIconSelect"
    />
  </Page>
</template>
