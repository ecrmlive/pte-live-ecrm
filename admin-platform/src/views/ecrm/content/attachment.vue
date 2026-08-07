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
import { getAccessCodesApi } from '#/api/core/auth';

/** 侧栏「系统素材」根入口（非真实分类 id） */
const SYSTEM_ROOT_ID = -1;

const categories = ref<AttachmentCategory[]>([]);
const files = ref<AttachmentItem[]>([]);
const categoryID = ref(0);
const kind = ref<AttachmentKind>('image');
const loading = ref(false);
const page = ref(1);
const pageSize = 24;
const total = ref(0);
const canManage = ref(false);
const accept = computed(() =>
  kind.value === 'image'
    ? 'image/jpeg,image/png,image/webp,image/gif'
    : 'video/mp4,video/quicktime,video/webm',
);
/** 图片上传上限 10MB（与后端 upload 包一致） */
const IMAGE_MAX_BYTES = 10 * 1024 * 1024;

/** 「全部素材」下上传时默认落入系统分类「其他图片」 */
const DEFAULT_UPLOAD_CATEGORY_ENNAME = 'other_image';

const systemCategories = computed(() =>
  categories.value.filter((row) => row.is_system === 1),
);
const customCategories = computed(() =>
  categories.value.filter((row) => row.is_system !== 1),
);

/** 当前选中分类；「全部素材」/「系统素材」根入口时回退到系统「其他图片」 */
function resolveUploadCategoryID() {
  if (categoryID.value > 0) return categoryID.value;
  const other = categories.value.find(
    (row) => row.attachment_category_enname === DEFAULT_UPLOAD_CATEGORY_ENNAME,
  );
  return other?.attachment_category_id ?? 0;
}

async function loadCategories() {
  const result = await listAttachmentCategoriesApi();
  categories.value = result.list ?? [];
}

async function loadFiles() {
  loading.value = true;
  try {
    const systemRoot = categoryID.value === SYSTEM_ROOT_ID;
    const result = await listAttachmentsApi({
      category_id: categoryID.value > 0 ? categoryID.value : undefined,
      is_system: systemRoot ? 1 : undefined,
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
  const { value } = await ElMessageBox.prompt('分类名称', '新增自定义素材分类', {
    inputPattern: /\S+/,
    inputErrorMessage: '请输入分类名称',
  });
  await createAttachmentCategoryApi({ attachment_category_name: value.trim() });
  await loadCategories();
  ElMessage.success('分类已创建');
}

async function editCategory(row: AttachmentCategory) {
  if (row.is_system === 1) {
    ElMessage.warning('系统素材分类不可修改');
    return;
  }
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
  if (row.is_system === 1) {
    ElMessage.warning('系统素材分类不可删除');
    return;
  }
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
  if (kind.value === 'image' && file.size > IMAGE_MAX_BYTES) {
    ElMessage.error('图片不能超过 10MB');
    return;
  }
  await uploadAttachmentApi(file, resolveUploadCategoryID());
  page.value = 1;
  await loadFiles();
  ElMessage.success('素材已上传');
}

async function removeFile(row: AttachmentItem) {
  await ElMessageBox.confirm(`确认删除素材「${row.attachment_name}」？`, '删除素材', {
    type: 'warning',
  });
  await deleteAttachmentApi(row.attachment_id);
  if (files.value.length === 1 && page.value > 1) page.value -= 1;
  await loadFiles();
  ElMessage.success('素材已删除');
}

onMounted(async () => {
  const [permissions] = await Promise.all([
    getAccessCodesApi(),
    loadCategories(),
    loadFiles(),
  ]);
  canManage.value = permissions.includes('content.attachment.manage');
});
</script>

<template>
  <Page title="素材库" description="平台素材按系统预设分类管理；可在分类内上传图片/视频。">
    <div class="attachment-page">
      <aside class="attachment-page__sidebar">
        <div class="attachment-page__side-head">
          <span>素材分类</span>
          <ElButton v-if="canManage" link type="primary" @click="addCategory">
            新增自定义
          </ElButton>
        </div>
        <button
          class="attachment-page__category"
          :class="{ active: categoryID === 0 }"
          @click="selectCategory(0)"
        >
          全部素材
        </button>

        <div class="attachment-page__scroll">
          <template v-if="customCategories.length">
            <div class="attachment-page__section">自定义分类</div>
            <div
              v-for="row in customCategories"
              :key="row.attachment_category_id"
              class="attachment-page__category-row"
            >
              <button
                class="attachment-page__category"
                :class="{ active: categoryID === row.attachment_category_id }"
                @click="selectCategory(row.attachment_category_id)"
              >
                {{ row.attachment_category_name }}
              </button>
              <span v-if="canManage" class="attachment-page__category-actions">
                <ElButton link size="small" @click="editCategory(row)">编辑</ElButton>
                <ElButton link size="small" type="danger" @click="removeCategory(row)">
                  删除
                </ElButton>
              </span>
            </div>
          </template>
          <div
            v-for="row in systemCategories"
            :key="row.attachment_category_id"
            class="attachment-page__category-row"
          >
            <button
              class="attachment-page__category"
              :class="{ active: categoryID === row.attachment_category_id }"
              @click="selectCategory(row.attachment_category_id)"
            >
              {{ row.attachment_category_name }}
            </button>
          </div>
        </div>

        <div v-if="systemCategories.length" class="attachment-page__footer">
          <button
            class="attachment-page__category"
            :class="{ active: categoryID === SYSTEM_ROOT_ID }"
            @click="selectCategory(SYSTEM_ROOT_ID)"
          >
            系统素材
          </button>
        </div>
      </aside>
      <main class="attachment-page__main">
        <div class="attachment-page__toolbar">
          <div class="space-x-2">
            <ElButton
              :type="kind === 'image' ? 'primary' : 'default'"
              @click="selectKind('image')"
            >
              图片
            </ElButton>
            <ElButton
              :type="kind === 'video' ? 'primary' : 'default'"
              @click="selectKind('video')"
            >
              视频
            </ElButton>
          </div>
          <ElUpload
            v-if="canManage"
            :accept="accept"
            :http-request="upload"
            :show-file-list="false"
          >
            <ElButton type="primary">
              上传{{ kind === 'image' ? '图片' : '视频' }}
            </ElButton>
          </ElUpload>
        </div>
        <div v-loading="loading" class="attachment-page__grid">
          <ElEmpty v-if="!loading && files.length === 0" description="暂无素材，请上传" />
          <article v-for="row in files" :key="row.attachment_id" class="attachment-card">
            <img
              v-if="row.attachment_type === 0"
              :src="row.attachment_src"
              :alt="row.attachment_name"
            />
            <video v-else :src="row.attachment_src" controls preload="metadata" />
            <div class="attachment-card__footer">
              <ElTag
                size="small"
                :type="row.attachment_type === 0 ? 'success' : 'warning'"
              >
                {{ row.attachment_type === 0 ? '图片' : '视频' }}
              </ElTag>
              <ElButton
                v-if="canManage"
                link
                type="danger"
                @click="removeFile(row)"
              >
                删除
              </ElButton>
            </div>
          </article>
        </div>
        <ElPagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          background
          layout="total, prev, pager, next"
          @current-change="loadFiles"
        />
      </main>
    </div>
  </Page>
</template>

<style scoped lang="scss">
.attachment-page {
  display: flex;
  min-height: 620px;
  overflow: hidden;
  border: 1px solid hsl(var(--border));
  border-radius: 8px;
  background: hsl(var(--background));
}
.attachment-page__sidebar {
  display: flex;
  flex: 0 0 220px;
  flex-direction: column;
  gap: 2px;
  width: 220px;
  padding: 14px 10px;
  border-right: 1px solid hsl(var(--border));
  overflow: hidden;
}
.attachment-page__scroll {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 2px;
  min-height: 0;
  overflow: auto;
}
.attachment-page__footer {
  display: flex;
  flex-shrink: 0;
  flex-direction: column;
  gap: 2px;
  margin-top: auto;
  padding-top: 8px;
  border-top: 1px solid hsl(var(--border));
}
.attachment-page__side-head,
.attachment-page__category-row,
.attachment-card__footer,
.attachment-page__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.attachment-page__side-head {
  flex-shrink: 0;
  padding: 0 8px 10px;
  font-weight: 600;
}
.attachment-page__section {
  margin: 10px 8px 4px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  font-weight: 600;
}
.attachment-page__category-row {
  position: relative;
}
.attachment-page__category {
  width: 100%;
  overflow: hidden;
  padding: 8px;
  border: 0;
  border-radius: 4px;
  color: hsl(var(--foreground));
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
  background: transparent;
  cursor: pointer;
}
.attachment-page__category:hover,
.attachment-page__category.active {
  color: hsl(var(--primary));
  background: hsl(var(--accent));
}
.attachment-page__category-actions {
  display: none;
  position: absolute;
  right: 4px;
  background: hsl(var(--accent));
}
.attachment-page__category-row:hover .attachment-page__category-actions {
  display: inline-flex;
}
.attachment-page__main {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 18px;
  min-width: 0;
  padding: 18px;
}
.attachment-page__toolbar {
  padding-bottom: 14px;
  border-bottom: 1px solid hsl(var(--border));
}
.attachment-page__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
  align-content: start;
  min-height: 440px;
}
.attachment-card {
  overflow: hidden;
  border: 1px solid hsl(var(--border));
  border-radius: 6px;
}
.attachment-card > img,
.attachment-card > video {
  display: block;
  width: 100%;
  height: 150px;
  object-fit: contain;
  object-position: center;
  background: hsl(var(--muted));
}
.attachment-card__footer {
  min-width: 0;
  padding: 8px;
}
@media (max-width: 768px) {
  .attachment-page {
    display: block;
  }
  .attachment-page__sidebar {
    width: auto;
    border-right: 0;
    border-bottom: 1px solid hsl(var(--border));
  }
}
</style>
