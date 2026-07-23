<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-upload v-if="canUpload" :show-upload-list="false" :before-upload="onUpload" accept="image/*">
        <a-button type="primary">上传图片</a-button>
      </a-upload>
      <a-button v-if="canUpload" @click="openCate">新建分类</a-button>
      <a-button @click="reload">刷新</a-button>
      <a-select
        v-model:value="cateFilter"
        allow-clear
        placeholder="全部分类"
        style="width: 160px"
        :options="cateOptions"
        @change="reload"
      />
    </div>
    <a-row :gutter="16">
      <a-col :span="6">
        <a-card size="small" title="分类">
          <a-list size="small" :data-source="categories">
            <template #renderItem="{ item }">
              <a-list-item>
                <a @click="filterBy(item.attachment_category_id)">{{ item.attachment_category_name }}</a>
                <template v-if="canDelete" #actions>
                  <a-button type="link" danger size="small" @click="removeCate(item)">删</a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
      <a-col :span="18">
        <a-table
          row-key="attachment_id"
          :loading="loading"
          :columns="columns"
          :data-source="list"
          :pagination="pagination"
          @change="onTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'preview'">
              <img v-if="record.attachment_src" class="thumb" :src="record.attachment_src" alt="" />
            </template>
            <template v-else-if="column.key === 'action'">
              <a-button type="link" @click="copy(record.attachment_src)">复制地址</a-button>
              <a-button v-if="canDelete" type="link" danger @click="onDelete(record)">删除</a-button>
            </template>
          </template>
        </a-table>
      </a-col>
    </a-row>

    <a-modal v-model:open="cateOpen" title="新建分类" :confirm-loading="saving" @ok="submitCate">
      <a-input v-model:value="cateName" placeholder="分类名称" />
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createCategory,
  deleteAttachment,
  deleteCategory,
  fetchAttachments,
  fetchCategories,
  uploadAttachment,
  type Attachment,
  type AttachmentCategory,
} from '@/api/attachment';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canUpload = computed(() => auth.hasPerm('attachment/upload'));
const canDelete = computed(() => auth.hasPerm('attachment/delete'));

const loading = ref(false);
const saving = ref(false);
const list = ref<Attachment[]>([]);
const categories = ref<AttachmentCategory[]>([]);
const cateFilter = ref<number | undefined>();
const cateOpen = ref(false);
const cateName = ref('');
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: '预览', key: 'preview', width: 90 },
  { title: '名称', dataIndex: 'attachment_name' },
  { title: '地址', dataIndex: 'attachment_src' },
  { title: '操作', key: 'action', width: 180 },
];

const cateOptions = computed(() =>
  categories.value.map((c) => ({ label: c.attachment_category_name, value: c.attachment_category_id })),
);

async function loadCates() {
  const { data } = await fetchCategories();
  categories.value = data.list || [];
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchAttachments({
      page: pagination.current,
      limit: pagination.pageSize,
      category_id: cateFilter.value || 0,
    });
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

function filterBy(id: number) {
  cateFilter.value = id;
  reload();
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  void load();
}

function openCate() {
  cateName.value = '';
  cateOpen.value = true;
}

async function submitCate() {
  if (!cateName.value.trim()) {
    message.warning('请填写名称');
    return;
  }
  saving.value = true;
  try {
    await createCategory({ attachment_category_name: cateName.value.trim() });
    message.success('已创建');
    cateOpen.value = false;
    await loadCates();
  } finally {
    saving.value = false;
  }
}

async function removeCate(row: AttachmentCategory) {
  await deleteCategory(row.attachment_category_id);
  message.success('已删除');
  if (cateFilter.value === row.attachment_category_id) cateFilter.value = undefined;
  await loadCates();
  void load();
}

function onUpload(file: File) {
  void (async () => {
    await uploadAttachment(file, cateFilter.value);
    message.success('上传成功');
    void load();
  })();
  return false;
}

async function onDelete(row: Attachment) {
  await deleteAttachment(row.attachment_id);
  message.success('已删除');
  void load();
}

async function copy(src: string) {
  try {
    await navigator.clipboard.writeText(src);
    message.success('已复制');
  } catch {
    message.info(src);
  }
}

onMounted(async () => {
  await loadCates();
  void load();
});
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.thumb {
  width: 56px;
  height: 56px;
  object-fit: cover;
  border-radius: 6px;
  background: #f5f5f5;
}
</style>
