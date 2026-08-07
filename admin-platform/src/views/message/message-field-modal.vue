<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, ref, watch } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { Plus } from '@element-plus/icons-vue';
import {
  ElButton,
  ElCheckbox,
  ElInput,
  ElMessage,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import MessageApi from '#/api/core/message';

import type { MessageFieldRow, MessageItem } from './types';

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  message?: MessageItem;
}>();

const emit = defineEmits<{
  success: [];
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

const [Drawer, drawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen && props.message?.message_id) {
      deleteIds.value = [];
      void loadFields(props.message.message_id);
    }
  },
  onConfirm: () => {
    void handleSubmit();
  },
});

watch(
  open,
  (visible) => {
    if (visible) {
      drawerApi.setState({ title: title.value }).open();
      return;
    }
    drawerApi.close();
  },
  { immediate: true },
);

watch(title, (value) => {
  if (open.value) {
    drawerApi.setState({ title: value });
  }
});

watch(
  () => props.message?.message_id,
  async (messageId) => {
    if (!open.value || !messageId) return;
    deleteIds.value = [];
    await loadFields(messageId);
  },
);

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

async function handleSubmit() {
  if (!props.message?.message_id) return;
  loading.value = true;
  drawerApi.setState({ confirmLoading: true });
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
    open.value = false;
    emit('success');
  } finally {
    loading.value = false;
    drawerApi.setState({ confirmLoading: false });
  }
}
</script>

<template>
  <Drawer
    :close-on-click-modal="false"
    :confirm-loading="loading"
    :destroy-on-close="true"
    :title="title"
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
  </Drawer>
</template>
