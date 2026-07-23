<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate">新建文章</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="article_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">{{ record.status === 1 ? '展示' : '隐藏' }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑文章' : '新建文章'"
      :confirm-loading="saving"
      width="720px"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="标题" required>
          <a-input v-model:value="form.title" />
        </a-form-item>
        <a-form-item label="作者">
          <a-input v-model:value="form.author" />
        </a-form-item>
        <a-form-item label="简介">
          <a-input v-model:value="form.synopsis" />
        </a-form-item>
        <a-form-item label="内容" required>
          <a-textarea v-model:value="form.content" :rows="8" />
        </a-form-item>
        <a-form-item label="分类ID">
          <a-input-number v-model:value="form.cid" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="展示">
          <a-switch v-model:checked="form.show" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createArticle,
  deleteArticle,
  fetchArticles,
  updateArticle,
  type Article,
} from '@/api/article';

const loading = ref(false);
const saving = ref(false);
const list = ref<Article[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  title: '',
  author: '',
  synopsis: '',
  content: '',
  cid: 1,
  sort: 0,
  show: true,
});

const columns = [
  { title: 'ID', dataIndex: 'article_id', width: 80 },
  { title: '标题', dataIndex: 'title' },
  { title: '作者', dataIndex: 'author', width: 120 },
  { title: '浏览', dataIndex: 'visit', width: 80 },
  { title: '状态', key: 'status', width: 90 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchArticles({ page: pagination.current, limit: pagination.pageSize });
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
  form.author = '';
  form.synopsis = '';
  form.content = '';
  form.cid = 1;
  form.sort = 0;
  form.show = true;
  modalOpen.value = true;
}

function openEdit(row: Article) {
  editingId.value = row.article_id;
  form.title = row.title;
  form.author = row.author;
  form.synopsis = row.synopsis;
  form.content = row.content;
  form.cid = row.cid;
  form.sort = row.sort;
  form.show = row.status === 1;
  modalOpen.value = true;
}

async function submit() {
  if (!form.title.trim() || !form.content.trim()) {
    message.warning('请填写标题与内容');
    return;
  }
  saving.value = true;
  try {
    const payload = {
      title: form.title.trim(),
      author: form.author.trim(),
      synopsis: form.synopsis.trim(),
      content: form.content,
      cid: form.cid,
      sort: form.sort,
      status: form.show ? 1 : 0,
    };
    if (editingId.value) {
      await updateArticle(editingId.value, payload);
    } else {
      await createArticle(payload);
    }
    message.success('已保存');
    modalOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function onDelete(row: Article) {
  await deleteArticle(row.article_id);
  message.success('已删除');
  await load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
</style>
