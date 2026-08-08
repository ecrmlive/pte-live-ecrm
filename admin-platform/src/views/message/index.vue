<script lang="ts" setup>
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import {
  ElButton,
  ElMessage,
  ElMessageBox,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import MessageApi from '#/api/core/message';
import { Page } from '@vben/common-ui';
import { deepClone } from '#/utils/base';
import { parseApiList } from '#/utils/list-response';

import MessageAddModal from './message-add-modal.vue';
import MessageEditModal from './message-edit-modal.vue';
import MessageFieldModal from './message-field-modal.vue';
import type { MessageFormModel, MessageItem } from './types';

const allRows = ref<MessageItem[]>([]);
const listLoaded = ref(false);

const addModalOpen = ref(false);
const editModalOpen = ref(false);
const fieldModalOpen = ref(false);
const messageModel = ref<MessageFormModel | MessageItem | undefined>();

function messageToText(row: MessageItem) {
  const val = row.message_to;
  return typeof val === 'object' ? val?.text ?? '—' : '—';
}

function messageTypeText(row: MessageItem) {
  const val = row.message_type;
  return typeof val === 'object' ? val?.text ?? '—' : '—';
}

async function loadList() {
  const res = await MessageApi.messageList({}, true);
  allRows.value = parseApiList<MessageItem>(res.data);
}

const gridOptions = {
  border: true,
  columns: [
    {
      field: 'message_name',
      minWidth: 140,
      title: '消息名称',
      treeNode: true,
    },
    { field: 'message_ename', minWidth: 140, title: '消息名称(英文)' },
    {
      field: 'message_to',
      minWidth: 120,
      slots: { default: 'messageTo' },
      title: '通知对象',
    },
    {
      field: 'message_type',
      minWidth: 120,
      slots: { default: 'messageType' },
      title: '消息类型',
    },
    {
      field: 'field',
      slots: { default: 'field' },
      title: '字段',
      width: 100,
    },
    { field: 'remark', minWidth: 160, title: '消息描述' },
    { field: 'sort', title: '排序', width: 80 },
    { field: 'create_time', minWidth: 160, title: '添加时间' },
    {
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 140,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 15, 20, 30, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!listLoaded.value) {
          try {
            await loadList();
          } catch {
            allRows.value = [];
          }
          listLoaded.value = true;
        }
        const start = (page.currentPage - 1) * page.pageSize;
        return {
          items: allRows.value.slice(start, start + page.pageSize),
          total: allRows.value.length,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'message_id' },
  treeConfig: {
    childrenField: 'children',
    expandAll: true,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

async function reload() {
  listLoaded.value = false;
  try {
    await loadList();
  } catch {
    allRows.value = [];
  }
  listLoaded.value = true;
  gridApi.reload();
}

function openEdit(item: MessageItem) {
  const model = deepClone(item) as MessageItem;
  messageModel.value = {
    ...model,
    message_to:
      typeof model.message_to === 'object'
        ? model.message_to.value
        : model.message_to,
    message_type:
      typeof model.message_type === 'object'
        ? model.message_type.value
        : model.message_type,
  };
  editModalOpen.value = true;
}

function openField(item: MessageItem) {
  messageModel.value = item;
  fieldModalOpen.value = true;
}

async function deleteMessage(row: MessageItem) {
  try {
    await ElMessageBox.confirm('删除后不可恢复，确认删除该记录吗?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
  } catch {
    return;
  }

  const res = await MessageApi.deleteMessage({ message_id: row.message_id });
  ElMessage.success(res.msg || '删除成功');
  await reload();
}
</script>

<template>
  <Page>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-access:code="'platform:message:add'"
          :icon="Plus"
          type="primary"
          @click="addModalOpen = true"
        >
          新增消息
        </ElButton>
      </template>

      <template #messageTo="{ row }">
        {{ messageToText(row) }}
      </template>

      <template #messageType="{ row }">
        {{ messageTypeText(row) }}
      </template>

      <template #field="{ row }">
        <ElButton link type="primary" @click="openField(row)">字段管理</ElButton>
      </template>

      <template #action="{ row }">
        <ElButton
          v-access:code="'platform:message:edit'"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-access:code="'platform:message:delete'"
          link
          type="primary"
          @click="deleteMessage(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <MessageAddModal v-model:open="addModalOpen" @success="reload" />
    <MessageEditModal
      v-model:open="editModalOpen"
      :form="messageModel as MessageFormModel"
      @success="reload"
    />
    <MessageFieldModal
      v-model:open="fieldModalOpen"
      :message="messageModel as MessageItem"
      @success="reload"
    />
  </Page>
</template>
