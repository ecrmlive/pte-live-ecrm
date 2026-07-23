<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
      <span class="hint">社区种草审核 · status 1 通过 / -1 驳回</span>
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
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canAudit" type="link" @click="audit(record, 1)">通过</a-button>
          <a-button v-if="canAudit" type="link" danger @click="audit(record, -1)">驳回</a-button>
          <a-button v-if="canDelete" type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { auditPost, deletePost, fetchPosts, type CommunityPost } from '@/api/community';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canAudit = computed(() => auth.hasPerm('community/audit'));
const canDelete = computed(() => auth.hasPerm('community/delete'));

const loading = ref(false);
const list = ref<CommunityPost[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: 'ID', dataIndex: 'community_id', width: 70 },
  { title: '标题', dataIndex: 'title' },
  { title: '作者', dataIndex: 'nickname', width: 120 },
  { title: '挂货', dataIndex: 'product_name', width: 160 },
  { title: '状态', key: 'status', width: 90 },
  { title: '操作', key: 'action', width: 220 },
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

async function audit(row: CommunityPost, status: number) {
  await auditPost(row.community_id, {
    status,
    refusal: status === -1 ? '不符合规范' : '',
    is_show: status === 1 ? 1 : 0,
  });
  message.success(status === 1 ? '已通过' : '已驳回');
  void load();
}

async function onDelete(row: CommunityPost) {
  await deletePost(row.community_id);
  message.success('已删除');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
}
.hint {
  color: #999;
  font-size: 12px;
}
</style>
