<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElSwitch,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createCustomerServiceQuickReply,
  deleteCustomerServiceQuickReply,
  fetchCustomerServiceQuickReplies,
  updateCustomerServiceQuickReply,
  type CustomerServiceQuickReply,
} from '#/api/core/customer-service';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import ImageField from '#/components/shop/image-field.vue';
import StoreRelationSelect from '#/components/ecrm/StoreRelationSelect.vue';

type AutoReplyForm = {
  content: string;
  keywords: string[];
  message_type: CustomerServiceQuickReply['message_type'];
  status: CustomerServiceQuickReply['status'];
  store_id?: number;
};

const editing = ref<CustomerServiceQuickReply>();
const form = reactive<AutoReplyForm>({
  content: '',
  keywords: [],
  message_type: 'text',
  status: 'enabled',
  store_id: undefined,
});

const gridOptions: VxeGridProps<CustomerServiceQuickReply> = {
  columns: [
    { field: 'id', title: 'ID', width: 86 },
    {
      field: 'title',
      minWidth: 260,
      showOverflow: 'tooltip',
      title: '关键词',
    },
    {
      formatter: ({ row }) =>
        row.message_type === 'image' ? '图片消息' : '文字消息',
      title: '回复类型',
      width: 130,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 128,
    },
    platformListActionColumn({ width: 160 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const result = await fetchCustomerServiceQuickReplies({
          limit: page.pageSize,
          page: page.currentPage,
        });
        return { items: result.list, total: result.total };
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });
const [ReplyDrawer, replyDrawerApi] = useVbenDrawer({
  cancelText: '取消',
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  onConfirm: async () => save(),
  placement: 'right',
});

const drawerTitle = computed(() => (editing.value ? '编辑自动回复' : '新增自动回复'));

function resetForm() {
  Object.assign(form, {
    content: '',
    keywords: [],
    message_type: 'text',
    status: 'enabled',
    store_id: undefined,
  });
}

function openCreate() {
  editing.value = undefined;
  resetForm();
  replyDrawerApi.setState({ title: '新增自动回复' }).open();
}

function openEdit(row: CustomerServiceQuickReply) {
  editing.value = row;
  Object.assign(form, {
    content: row.content,
    keywords: row.title.split(',').map((keyword) => keyword.trim()).filter(Boolean),
    message_type: row.message_type || 'text',
    status: row.status,
    store_id: row.store_id,
  });
  replyDrawerApi.setState({ title: '编辑自动回复' }).open();
}

async function save() {
  const keywords = Array.from(
    new Set(form.keywords.map((keyword) => keyword.trim()).filter(Boolean)),
  );
  if (!keywords.length || !form.content.trim()) {
    ElMessage.warning('请填写关键词与回复内容');
    return;
  }
  if (!editing.value && !form.store_id) {
    ElMessage.warning('请选择关联店铺');
    return;
  }
  replyDrawerApi.lock();
  try {
    if (editing.value) {
      await updateCustomerServiceQuickReply(editing.value.id, {
        content: form.content.trim(),
        message_type: form.message_type,
        status: form.status,
        store_id: form.store_id,
        title: keywords.join(','),
      });
    } else {
      await createCustomerServiceQuickReply({
        content: form.content.trim(),
        message_type: form.message_type,
        status: form.status,
        store_id: form.store_id!,
        title: keywords.join(','),
      });
    }
    replyDrawerApi.close();
    ElMessage.success(editing.value ? '自动回复已更新' : '自动回复已新增');
    gridApi.reload();
  } finally {
    replyDrawerApi.unlock();
  }
}

async function toggleStatus(
  row: CustomerServiceQuickReply,
  status: CustomerServiceQuickReply['status'],
) {
  try {
    await updateCustomerServiceQuickReply(row.id, {
      content: row.content,
      message_type: row.message_type,
      status,
      store_id: row.store_id,
      title: row.title,
    });
    ElMessage.success('规则状态已更新');
    gridApi.reload();
  } catch {
    gridApi.reload();
  }
}

async function remove(row: CustomerServiceQuickReply) {
  try {
    await confirm({
      content: `删除“${row.title}”后将不再作为可用客服回复，历史审计记录会保留。是否继续？`,
      icon: 'warning',
      title: '删除自动回复',
    });
    await deleteCustomerServiceQuickReply(row.id);
    ElMessage.success('自动回复已删除');
    gridApi.reload();
  } catch {
    // 取消或请求失败时保留现有列表状态。
  }
}

</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="openCreate">
          新增关键词
        </ElButton>
      </template>
      <template #status="{ row }">
        <ElSwitch
          :active-value="'enabled'"
          :inactive-value="'disabled'"
          :model-value="row.status"
          @update:model-value="toggleStatus(row, $event)"
        />
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openEdit(row)">编辑</ElButton>
        <ElButton link type="danger" @click="remove(row)">删除</ElButton>
      </template>
    </Grid>

    <ReplyDrawer :title="drawerTitle">
      <ElForm label-width="120px" class="auto-reply__form">
        <ElFormItem label="关键词" required>
          <ElSelect
            v-model="form.keywords"
            allow-create
            default-first-option
            filterable
            multiple
            placeholder="输入后回车"
            :reserve-keyword="false"
          />
        </ElFormItem>
        <ElFormItem label="规则状态">
          <ElRadioGroup v-model="form.status">
            <ElRadio value="enabled">启用</ElRadio>
            <ElRadio value="disabled">禁用</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
        <ElFormItem label="消息类型" required>
          <ElSelect v-model="form.message_type" placeholder="请选择消息类型">
            <ElOption label="文字消息" value="text" />
            <ElOption label="图片消息" value="image" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="关联店铺" required>
          <StoreRelationSelect
            v-model="form.store_id"
            placeholder="请选择关联店铺"
          />
        </ElFormItem>
        <ElFormItem v-if="form.message_type === 'text'" label="规则内容" required>
          <ElInput
            v-model="form.content"
            :rows="6"
            maxlength="2000"
            placeholder="请输入文字回复内容"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
        <ElFormItem v-else label="规则内容" required>
          <ImageField
            v-model="form.content"
            default-library="system"
            hint="请选择要发送的图片"
            :preview-size="120"
          />
        </ElFormItem>
      </ElForm>
    </ReplyDrawer>
  </Page>
</template>

<style scoped>
.auto-reply__form {
  max-width: 920px;
  padding: 24px 0;
}

.auto-reply__form :deep(.el-select) {
  width: 100%;
}
</style>
