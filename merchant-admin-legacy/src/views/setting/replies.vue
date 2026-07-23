<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canWrite" type="primary" @click="openCreate">新建话术</a-button>
      <a-button @click="reload">刷新</a-button>
      <span class="hint">客服工作台按 mer_id 加载已启用项</span>
    </div>
    <a-table
      row-key="service_reply_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          {{ record.status === 1 ? '启用' : '停用' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canWrite" type="link" @click="openEdit(record)">编辑</a-button>
          <a-button v-if="canWrite" type="link" @click="toggle(record)">
            {{ record.status === 1 ? '停用' : '启用' }}
          </a-button>
          <a-popconfirm v-if="canWrite" title="确认删除？" @confirm="onDelete(record)">
            <a-button type="link" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑话术' : '新建话术'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="标题/关键字" required>
          <a-input v-model:value="form.keyword" maxlength="64" />
        </a-form-item>
        <a-form-item label="回复内容" required>
          <a-textarea v-model:value="form.content" :rows="4" maxlength="512" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model:checked="form.on" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createReply,
  deleteReply,
  fetchReplies,
  updateReply,
  type ServiceReply,
} from '@/api/reply';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canWrite = computed(() => auth.hasPerm('reply/write'));

const loading = ref(false);
const saving = ref(false);
const list = ref<ServiceReply[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({ keyword: '', content: '', on: true });

const columns = [
  { title: 'ID', dataIndex: 'service_reply_id', width: 70 },
  { title: '标题', dataIndex: 'keyword', width: 120 },
  { title: '内容', dataIndex: 'content' },
  { title: '状态', key: 'status', width: 80 },
  { title: '操作', key: 'action', width: 220 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchReplies({ page: pagination.current, limit: pagination.pageSize });
    list.value = data.list || [];
    pagination.total = data.total || 0;
  } finally {
    loading.value = false;
  }
}

function reload() {
  pagination.current = 1;
  void load();
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  void load();
}

function openCreate() {
  editingId.value = 0;
  form.keyword = '';
  form.content = '';
  form.on = true;
  modalOpen.value = true;
}

function openEdit(row: ServiceReply) {
  editingId.value = row.service_reply_id;
  form.keyword = row.keyword;
  form.content = row.content;
  form.on = row.status === 1;
  modalOpen.value = true;
}

async function submit() {
  if (!form.keyword.trim() || !form.content.trim()) {
    message.warning('请填写标题与内容');
    return;
  }
  saving.value = true;
  try {
    const body = {
      keyword: form.keyword.trim(),
      content: form.content.trim(),
      type: 1,
      status: form.on ? 1 : 0,
    };
    if (editingId.value) {
      await updateReply(editingId.value, body);
    } else {
      await createReply(body);
    }
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function toggle(row: ServiceReply) {
  const next = row.status === 1 ? 0 : 1;
  await updateReply(row.service_reply_id, {
    keyword: row.keyword,
    content: row.content,
    type: row.type || 1,
    status: next,
  });
  message.success(next === 1 ? '已启用' : '已停用');
  void load();
}

async function onDelete(row: ServiceReply) {
  await deleteReply(row.service_reply_id);
  message.success('已删除');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.page-card {
  border-radius: 14px;
}
.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}
.hint {
  color: #8a8f98;
  font-size: 13px;
}
</style>
