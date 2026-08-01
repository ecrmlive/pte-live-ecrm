<script setup lang="ts">
import type { AttachmentCategory, AttachmentItem, AttachmentKind } from '#/api/core/attachment';

import { Page } from '@vben/common-ui';
import { ElButton, ElEmpty, ElMessage, ElMessageBox, ElPagination, ElTag, ElUpload } from 'element-plus';
import { computed, onMounted, ref } from 'vue';

import {
  createAttachmentCategoryApi,
  deleteAttachmentApi,
  deleteAttachmentCategoryApi,
  listAttachmentCategoriesApi,
  listAttachmentsApi,
  updateAttachmentCategoryApi,
  uploadAttachmentApi,
} from '#/api/core/attachment';

const categories = ref<AttachmentCategory[]>([]);
const files = ref<AttachmentItem[]>([]);
const categoryID = ref(0);
const kind = ref<AttachmentKind>('image');
const loading = ref(false);
const page = ref(1);
const pageSize = 24;
const total = ref(0);
const accept = computed(() => kind.value === 'image' ? 'image/jpeg,image/png,image/webp,image/gif' : 'video/mp4,video/quicktime,video/webm');

async function loadCategories() {
  const result = await listAttachmentCategoriesApi();
  categories.value = result.list ?? [];
}

async function loadFiles() {
  loading.value = true;
  try {
    const result = await listAttachmentsApi({
      category_id: categoryID.value || undefined,
      limit: pageSize,
      page: page.value,
      type: kind.value,
    });
    files.value = result.list ?? [];
    total.value = result.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function selectCategory(id: number) {
  categoryID.value = id;
  page.value = 1;
  void loadFiles();
}

function selectKind(value: AttachmentKind) {
  kind.value = value;
  page.value = 1;
  void loadFiles();
}

async function addCategory() {
  const { value } = await ElMessageBox.prompt('分类名称', '新增素材分类', {
    inputPattern: /\S+/,
    inputErrorMessage: '请输入分类名称',
  });
  await createAttachmentCategoryApi({ attachment_category_name: value.trim() });
  await loadCategories();
  ElMessage.success('分类已创建');
}

async function editCategory(row: AttachmentCategory) {
  const { value } = await ElMessageBox.prompt('分类名称', '编辑素材分类', {
    inputValue: row.attachment_category_name,
    inputPattern: /\S+/,
    inputErrorMessage: '请输入分类名称',
  });
  await updateAttachmentCategoryApi(row.attachment_category_id, {
    attachment_category_enname: row.attachment_category_enname,
    attachment_category_name: value.trim(),
    pid: row.pid,
    sort: row.sort,
  });
  await loadCategories();
  ElMessage.success('分类已更新');
}

async function removeCategory(row: AttachmentCategory) {
  await ElMessageBox.confirm(
    `删除「${row.attachment_category_name}」后，分类内素材会移动到“全部素材”，确认继续？`,
    '删除素材分类',
    { type: 'warning' },
  );
  await deleteAttachmentCategoryApi(row.attachment_category_id);
  if (categoryID.value === row.attachment_category_id) categoryID.value = 0;
  await Promise.all([loadCategories(), loadFiles()]);
  ElMessage.success('分类已删除');
}

async function upload({ file }: { file: File }) {
  const isImage = file.type.startsWith('image/');
  const isVideo = file.type.startsWith('video/');
  if ((kind.value === 'image' && !isImage) || (kind.value === 'video' && !isVideo)) {
    ElMessage.warning(kind.value === 'image' ? '请选择图片文件' : '请选择视频文件');
    return;
  }
  await uploadAttachmentApi(file, categoryID.value);
  await loadFiles();
  ElMessage.success('素材已上传');
}

async function removeFile(row: AttachmentItem) {
  await ElMessageBox.confirm(`确认删除素材「${row.attachment_name}」？`, '删除素材', { type: 'warning' });
  await deleteAttachmentApi(row.attachment_id);
  if (files.value.length === 1 && page.value > 1) page.value -= 1;
  await loadFiles();
  ElMessage.success('素材已删除');
}

onMounted(async () => {
  await loadCategories();
  await loadFiles();
});
</script>

<template>
  <Page title="素材库" description="平台素材按分类管理；上传使用后台保存的腾讯云 COS 配置，未启用时遵循当前上传配置。">
    <div class="attachment-page">
      <aside class="attachment-page__sidebar">
        <div class="attachment-page__side-head">
          <span>素材分类</span>
          <ElButton link type="primary" @click="addCategory">新增</ElButton>
        </div>
        <button class="attachment-page__category" :class="{ active: categoryID === 0 }" @click="selectCategory(0)">全部素材</button>
        <div v-for="row in categories" :key="row.attachment_category_id" class="attachment-page__category-row">
          <button class="attachment-page__category" :class="{ active: categoryID === row.attachment_category_id }" @click="selectCategory(row.attachment_category_id)">{{ row.attachment_category_name }}</button>
          <span class="attachment-page__category-actions">
            <ElButton link size="small" @click="editCategory(row)">编辑</ElButton>
            <ElButton link size="small" type="danger" @click="removeCategory(row)">删除</ElButton>
          </span>
        </div>
      </aside>
      <main class="attachment-page__main">
        <div class="attachment-page__toolbar">
          <div class="space-x-2">
            <ElButton :type="kind === 'image' ? 'primary' : 'default'" @click="selectKind('image')">图片</ElButton>
            <ElButton :type="kind === 'video' ? 'primary' : 'default'" @click="selectKind('video')">视频</ElButton>
          </div>
          <ElUpload :accept="accept" :http-request="upload" :show-file-list="false">
            <ElButton type="primary">上传{{ kind === 'image' ? '图片' : '视频' }}</ElButton>
          </ElUpload>
        </div>
        <div v-loading="loading" class="attachment-page__grid">
          <ElEmpty v-if="!loading && files.length === 0" description="暂无素材，请上传" />
          <article v-for="row in files" :key="row.attachment_id" class="attachment-card">
            <img v-if="row.attachment_type === 0" :src="row.attachment_src" :alt="row.attachment_name" />
            <video v-else :src="row.attachment_src" controls preload="metadata" />
            <div class="attachment-card__footer">
              <div class="min-w-0">
                <p :title="row.attachment_name">{{ row.attachment_name }}</p>
                <ElTag size="small" :type="row.attachment_type === 0 ? 'success' : 'warning'">{{ row.attachment_type === 0 ? '图片' : '视频' }}</ElTag>
              </div>
              <ElButton link type="danger" @click="removeFile(row)">删除</ElButton>
            </div>
          </article>
        </div>
        <ElPagination v-if="total > pageSize" v-model:current-page="page" :page-size="pageSize" :total="total" background layout="prev, pager, next" @current-change="loadFiles" />
      </main>
    </div>
  </Page>
</template>

<style scoped lang="scss">
.attachment-page { display: flex; min-height: 620px; overflow: hidden; border: 1px solid hsl(var(--border)); border-radius: 8px; background: hsl(var(--background)); }
.attachment-page__sidebar { width: 220px; flex: 0 0 220px; padding: 14px 10px; border-right: 1px solid hsl(var(--border)); }
.attachment-page__side-head, .attachment-page__category-row, .attachment-card__footer, .attachment-page__toolbar { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.attachment-page__side-head { padding: 0 8px 10px; font-weight: 600; }
.attachment-page__category-row { position: relative; }
.attachment-page__category { width: 100%; overflow: hidden; padding: 8px; color: hsl(var(--foreground)); text-align: left; text-overflow: ellipsis; white-space: nowrap; border-radius: 4px; }
.attachment-page__category:hover, .attachment-page__category.active { color: hsl(var(--primary)); background: hsl(var(--accent)); }
.attachment-page__category-actions { display: none; position: absolute; right: 4px; background: hsl(var(--accent)); }
.attachment-page__category-row:hover .attachment-page__category-actions { display: inline-flex; }
.attachment-page__main { display: flex; flex: 1; flex-direction: column; gap: 18px; min-width: 0; padding: 18px; }
.attachment-page__toolbar { padding-bottom: 14px; border-bottom: 1px solid hsl(var(--border)); }
.attachment-page__grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 16px; align-content: start; min-height: 440px; }
.attachment-card { overflow: hidden; border: 1px solid hsl(var(--border)); border-radius: 6px; }
.attachment-card > img, .attachment-card > video { display: block; width: 100%; height: 150px; object-fit: cover; background: hsl(var(--muted)); }
.attachment-card__footer { min-width: 0; padding: 8px; }
.attachment-card__footer p { overflow: hidden; margin: 0 0 4px; font-size: 12px; text-overflow: ellipsis; white-space: nowrap; }
@media (max-width: 768px) { .attachment-page { display: block; } .attachment-page__sidebar { width: auto; border-right: 0; border-bottom: 1px solid hsl(var(--border)); } }
</style>
