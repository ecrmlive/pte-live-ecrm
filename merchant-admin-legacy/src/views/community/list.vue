<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canCreate" type="primary" @click="openCreate">发帖</a-button>
      <a-button @click="reload">刷新</a-button>
      <span class="hint">本店种草（mer_id 隔离）；发帖/改帖后待平台审核</span>
    </div>
    <a-table
      row-key="community_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          {{ statusText(record.status) }}
          <span v-if="record.status === -1 && record.refusal" class="refusal">（{{ record.refusal }}）</span>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canUpdate" type="link" @click="openEdit(record)">编辑</a-button>
          <a-button type="link" @click="openReplies(record)">评论</a-button>
          <a-popconfirm v-if="canDelete" title="确认删除？" @confirm="onDelete(record)">
            <a-button type="link" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑帖子' : '发帖'"
      :confirm-loading="saving"
      width="560px"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="标题" required>
          <a-input v-model:value="form.title" />
        </a-form-item>
        <a-form-item label="内容" required>
          <a-textarea v-model:value="form.content" :rows="5" />
        </a-form-item>
        <a-form-item label="封面图 URL">
          <a-input v-model:value="form.image" placeholder="可选，素材库地址" />
        </a-form-item>
        <a-form-item label="挂货商品 ID">
          <a-input-number v-model:value="form.product_id" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="replyOpen" :title="`评论 · #${replyPostId}`" :footer="null" width="520px">
      <a-spin :spinning="replyLoading">
        <a-list v-if="replies.length" size="small" :data-source="replies">
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta :title="item.nickname || '用户'" :description="item.content" />
            </a-list-item>
          </template>
        </a-list>
        <a-empty v-else description="暂无评论" />
      </a-spin>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createPost,
  deletePost,
  fetchPosts,
  fetchReplies,
  updatePost,
  type CommunityPost,
  type CommunityReply,
} from '@/api/community';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canCreate = computed(() => auth.hasPerm('community/create'));
const canUpdate = computed(() => auth.hasPerm('community/update'));
const canDelete = computed(() => auth.hasPerm('community/delete'));

const loading = ref(false);
const saving = ref(false);
const list = ref<CommunityPost[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({ title: '', content: '', image: '', product_id: 0 });
const replyOpen = ref(false);
const replyLoading = ref(false);
const replyPostId = ref(0);
const replies = ref<CommunityReply[]>([]);

const columns = [
  { title: 'ID', dataIndex: 'community_id', width: 70 },
  { title: '标题', dataIndex: 'title' },
  { title: '作者', dataIndex: 'nickname', width: 100 },
  { title: '挂货', dataIndex: 'product_name', width: 140 },
  { title: '评论', dataIndex: 'count_reply', width: 70 },
  { title: '状态', key: 'status', width: 160 },
  { title: '操作', key: 'action', width: 200 },
];

function statusText(s: number) {
  if (s === 1) return '已通过';
  if (s === -1) return '已驳回';
  return '待审';
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchPosts({ page: pagination.current, limit: pagination.pageSize });
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
  form.title = '';
  form.content = '';
  form.image = '';
  form.product_id = 0;
  modalOpen.value = true;
}

function openEdit(row: CommunityPost) {
  editingId.value = row.community_id;
  form.title = row.title;
  form.content = row.content;
  form.image = row.image || '';
  form.product_id = row.product_id || 0;
  modalOpen.value = true;
}

async function submit() {
  if (!form.title.trim() || !form.content.trim()) {
    message.warning('请填写标题和内容');
    return;
  }
  saving.value = true;
  const body = {
    title: form.title.trim(),
    content: form.content.trim(),
    image: form.image.trim(),
    product_id: form.product_id || 0,
  };
  try {
    if (editingId.value) {
      await updatePost(editingId.value, body);
    } else {
      await createPost(body);
    }
    message.success('已保存，待平台审核');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function onDelete(row: CommunityPost) {
  await deletePost(row.community_id);
  message.success('已删除');
  void load();
}

async function openReplies(row: CommunityPost) {
  replyPostId.value = row.community_id;
  replyOpen.value = true;
  replyLoading.value = true;
  try {
    const { data } = await fetchReplies(row.community_id, { page: 1, limit: 50 });
    replies.value = data.list || [];
  } finally {
    replyLoading.value = false;
  }
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.hint {
  color: #999;
  font-size: 13px;
}
.refusal {
  color: #c0392b;
  font-size: 12px;
}
</style>
