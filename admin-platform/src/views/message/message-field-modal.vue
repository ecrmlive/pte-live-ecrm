<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, ref, watch } from 'vue';

import {
  ElButton,
  ElCheckbox,
  ElDialog,
  ElInput,
  ElMessage,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import MessageApi from '#/api/core/message';

import type { MessageFieldRow, MessageItem } from './types';

const props = defineProps<{
  message?: MessageItem;
  open: boolean;
}>();

const emit = defineEmits<{
  success: [];
  'update:open': [value: boolean];
}>();

const loading = ref(true);
const fieldData = ref<MessageFieldRow[]>([]);
const deleteIds = ref<number[]>([]);

const title = computed(() => {
  const name = props.message?.message_name || '';
  return name ? `字段管理(${name})` : '字段管理';
});

const gridOptions: VxeGridProps<MessageFieldRow> = {
  border: true,
  columns: [
    {
      field: 'field_name',
      minWidth: 120,
      slots: { default: 'fieldName' },
      title: '字段名称',
    },
    {
      field: 'field_ename',
      minWidth: 130,
      slots: { default: 'fieldEname' },
      title: '字段英文名',
    },
    {
      field: 'filed_value',
      minWidth: 140,
      slots: { default: 'fieldValue' },
      title: '字段默认值',
    },
    {
      field: 'is_var',
      slots: { default: 'isVar' },
      title: '是否变量',
      width: 90,
    },
    {
      field: 'sort',
      slots: { default: 'sort' },
      title: '排序',
      width: 90,
    },
    {
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 80,
    },
  ],
  height: 320,
  loading: false,
  pagerConfig: { enabled: false },
  rowConfig: { isHover: true, keyField: 'message_field_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

watch(
  fieldData,
  (rows) => {
    gridApi.setGridOptions({ data: rows });
  },
  { deep: true },
);

watch(loading, (value) => {
  gridApi.setLoading(value);
});

watch(
  () => [props.open, props.message?.message_id] as const,
  async ([open, messageId]) => {
    if (!open || !messageId) return;
    deleteIds.value = [];
    await loadFields(messageId);
  },
);

async function loadFields(messageId: number) {
  loading.value = true;
  try {
    const res = await MessageApi.fieldList({ message_id: messageId }, true);
    fieldData.value =
      (res.data as { list?: MessageFieldRow[] })?.list ?? [];
  } catch {
    fieldData.value = [];
  } finally {
    loading.value = false;
  }
}

function addField() {
  if (!props.message?.message_id) return;
  fieldData.value.push({
    field_ename: '',
    field_name: '',
    filed_value: '',
    message_field_id: 0,
    message_id: props.message.message_id,
    sort: 100,
  });
}

function deleteField(row: MessageFieldRow) {
  const index = fieldData.value.indexOf(row);
  if (index < 0) return;
  const field = fieldData.value[index];
  if (field && field.message_field_id > 0) {
    deleteIds.value.push(field.message_field_id);
  }
  fieldData.value.splice(index, 1);
}

function checkRow(checked: boolean, row: MessageFieldRow) {
  row.is_var = checked ? 1 : 0;
}

function handleClose() {
  emit('update:open', false);
}

async function handleSubmit() {
  if (!props.message?.message_id) return;
  loading.value = true;
  try {
    await MessageApi.saveField(
      {
        deleteIds: deleteIds.value,
        fieldData: fieldData.value,
        message_id: props.message.message_id,
      },
      true,
    );
    ElMessage.success('恭喜你，修改成功');
    emit('update:open', false);
    emit('success');
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <ElDialog
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :model-value="open"
    :title="title"
    width="760px"
    @close="handleClose"
    @update:model-value="emit('update:open', $event)"
  >
    <div class="mb-3">
      <ElButton :icon="Plus" type="primary" @click="addField">添加字段</ElButton>
    </div>

    <Grid>
      <template #fieldName="{ row }">
        <ElInput v-model="row.field_name" />
      </template>
      <template #fieldEname="{ row }">
        <ElInput v-model="row.field_ename" />
      </template>
      <template #fieldValue="{ row }">
        <ElInput v-model="row.filed_value" />
      </template>
      <template #isVar="{ row }">
        <ElCheckbox
          :model-value="row.is_var === 1"
          @change="(checked) => checkRow(Boolean(checked), row)"
        />
      </template>
      <template #sort="{ row }">
        <ElInput v-model="row.sort" />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="deleteField(row)">删除</ElButton>
      </template>
    </Grid>

    <template #footer>
      <ElButton @click="handleClose">取消</ElButton>
      <ElButton :loading="loading" type="primary" @click="handleSubmit">确定</ElButton>
    </template>
  </ElDialog>
</template>
