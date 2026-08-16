<script setup lang="ts">
import type { AttachmentCategory, AttachmentItem, AttachmentKind } from '#/api/core/attachment';
import type { UploadRequestOptions } from 'element-plus';

import { Page } from '@vben/common-ui';
import { Delete, Grid, List, Search, VideoPlay } from '@element-plus/icons-vue';
import { ElButton, ElCheckbox, ElEmpty, ElInput, ElMessage, ElMessageBox, ElOption, ElPagination, ElSelect, ElUpload } from 'element-plus';
import { computed, onMounted, ref } from 'vue';

import { createAttachmentCategoryApi, deleteAttachmentApi, deleteAttachmentCategoryApi, listAttachmentCategoriesApi, listAttachmentsApi, moveAttachmentsApi, updateAttachmentCategoryApi, uploadAttachmentApi } from '#/api/core/attachment';
import { getAccessCodesApi } from '#/api/core/auth';

const SYSTEM_ROOT_ID = -1;
const pageSize = 24;
const IMAGE_MAX_BYTES = 10 * 1024 * 1024;
const categories = ref<AttachmentCategory[]>([]);
const files = ref<AttachmentItem[]>([]);
const categoryID = ref(0);
const kind = ref<AttachmentKind>('image');
const libraryMode = ref<'all' | 'system'>('all');
const keyword = ref('');
const queryKeyword = ref('');
const page = ref(1);
const total = ref(0);
const loading = ref(false);
const canManage = ref(false);
const viewMode = ref<'grid' | 'list'>('grid');
const selectedIDs = ref<number[]>([]);
const moveTarget = ref<number>();

const accept = computed(() => kind.value === 'image' ? 'image/jpeg,image/png,image/webp,image/gif' : 'video/mp4,video/quicktime,video/webm');
const systemCategories = computed(() => categories.value.filter((item) => Number(item.is_system) === 1));
const customCategories = computed(() => categories.value.filter((item) => Number(item.is_system) !== 1));
const allSelected = computed(() => files.value.length > 0 && files.value.every((item) => selectedIDs.value.includes(item.attachment_id)));
const selectedCount = computed(() => selectedIDs.value.length);
const activeCategoryName = computed(() => {
  if (libraryMode.value === 'system') return '系统素材';
  if (categoryID.value === 0) return kind.value === 'image' ? '全部图片' : '全部视频';
  return categories.value.find((item) => item.attachment_category_id === categoryID.value)?.attachment_category_name || '全部素材';
});

function resetSelection() {
  selectedIDs.value = [];
  moveTarget.value = undefined;
}
async function loadCategories() {
  const result = await listAttachmentCategoriesApi({ type: kind.value });
  categories.value = result.list ?? [];
}
async function loadFiles() {
  loading.value = true;
  try {
    const params: Parameters<typeof listAttachmentsApi>[0] = { limit: pageSize, page: page.value, type: kind.value, keyword: queryKeyword.value || undefined };
    if (libraryMode.value === 'system') {
      params.is_system = 1;
      if (categoryID.value > 0) params.category_id = categoryID.value;
    } else if (categoryID.value > 0) {
      params.category_id = categoryID.value;
    }
    const result = await listAttachmentsApi(params);
    files.value = result.list ?? [];
    total.value = result.total ?? 0;
    resetSelection();
  } finally {
    loading.value = false;
  }
}
async function selectKind(value: AttachmentKind) {
  if (kind.value === value) return;
  kind.value = value;
  categoryID.value = 0;
  libraryMode.value = 'all';
  page.value = 1;
  queryKeyword.value = '';
  keyword.value = '';
  await Promise.all([loadCategories(), loadFiles()]);
}
function selectCategory(id: number) {
  libraryMode.value = id === SYSTEM_ROOT_ID ? 'system' : 'all';
  categoryID.value = id;
  page.value = 1;
  void loadFiles();
}
function search() {
  queryKeyword.value = keyword.value.trim();
  page.value = 1;
  void loadFiles();
}
function toggleSelected(id: number, checked: boolean) {
  selectedIDs.value = checked ? [...new Set([...selectedIDs.value, id])] : selectedIDs.value.filter((item) => item !== id);
}
function toggleAll(checked: boolean) {
  selectedIDs.value = checked ? files.value.map((item) => item.attachment_id) : [];
}
function uploadTarget() {
  if (libraryMode.value === 'system') {
    const fallback = systemCategories.value.find((item) => kind.value === 'image' ? item.attachment_category_enname === 'other_image' : item.attachment_category_enname === 'other_video');
    return { categoryID: categoryID.value > 0 ? categoryID.value : (fallback?.attachment_category_id ?? 0), isSystem: true };
  }
  return { categoryID: Math.max(0, categoryID.value), isSystem: false };
}
async function upload({ file }: UploadRequestOptions) {
  const rawFile = file as File;
  if (kind.value === 'image' && !rawFile.type.startsWith('image/')) return void ElMessage.warning('请选择图片文件');
  if (kind.value === 'video' && !rawFile.type.startsWith('video/')) return void ElMessage.warning('请选择视频文件');
  if (kind.value === 'image' && rawFile.size > IMAGE_MAX_BYTES) return void ElMessage.error('图片不能超过 10MB');
  const target = uploadTarget();
  if (target.isSystem && !target.categoryID) return void ElMessage.warning('请先选择系统素材分类');
  await uploadAttachmentApi(rawFile, target.categoryID, { isSystem: target.isSystem });
  page.value = 1;
  await loadFiles();
  ElMessage.success('素材已上传');
}
async function removeFiles(ids = selectedIDs.value) {
  if (!ids.length) return;
  await ElMessageBox.confirm(`确认删除已选 ${ids.length} 个素材吗？`, '删除素材', { type: 'warning' });
  await Promise.all(ids.map((id) => deleteAttachmentApi(id)));
  if (files.value.length === ids.length && page.value > 1) page.value -= 1;
  await loadFiles();
  ElMessage.success('素材已删除');
}
async function moveSelected() {
  if (!selectedIDs.value.length || !moveTarget.value) return;
  await moveAttachmentsApi(selectedIDs.value, moveTarget.value);
  await loadFiles();
  ElMessage.success('素材已移动');
}
async function addCategory() {
  const { value } = await ElMessageBox.prompt('分类名称', '新增素材分类', { inputPattern: /\S+/, inputErrorMessage: '请输入分类名称' });
  await createAttachmentCategoryApi({ attachment_category_name: value.trim() });
  await loadCategories();
  ElMessage.success('分类已创建');
}
async function editCategory(row: AttachmentCategory) {
  const { value } = await ElMessageBox.prompt('分类名称', '编辑素材分类', { inputValue: row.attachment_category_name, inputPattern: /\S+/, inputErrorMessage: '请输入分类名称' });
  await updateAttachmentCategoryApi(row.attachment_category_id, { attachment_category_enname: row.attachment_category_enname, attachment_category_name: value.trim(), pid: row.pid, sort: row.sort });
  await loadCategories();
  ElMessage.success('分类已更新');
}
async function removeCategory(row: AttachmentCategory) {
  await ElMessageBox.confirm(`删除「${row.attachment_category_name}」后，素材将移至全部素材，确认继续？`, '删除素材分类', { type: 'warning' });
  await deleteAttachmentCategoryApi(row.attachment_category_id);
  if (categoryID.value === row.attachment_category_id) categoryID.value = 0;
  await Promise.all([loadCategories(), loadFiles()]);
  ElMessage.success('分类已删除');
}
onMounted(async () => {
  const [permissions] = await Promise.all([getAccessCodesApi(), loadCategories(), loadFiles()]);
  canManage.value = permissions.includes('content.attachment.manage');
});
</script>

<template>
  <Page auto-content-height content-class="!p-0">
    <div class="asset-page">
      <header class="asset-page__tabs">
        <button :class="{ active: kind === 'image' }" @click="selectKind('image')">图片管理</button>
        <button :class="{ active: kind === 'video' }" @click="selectKind('video')">视频管理</button>
      </header>
      <div class="asset-page__body">
        <aside class="asset-page__sidebar">
          <div class="asset-page__category-row">
            <button class="asset-page__category" :class="{ active: categoryID === 0 && libraryMode === 'all' }" @click="selectCategory(0)">全部{{ kind === 'image' ? '图片' : '视频' }}</button>
            <ElButton v-if="canManage" link size="small" type="primary" @click="addCategory">新增</ElButton>
          </div>
          <div class="asset-page__category-section">自定义分组</div>
          <div v-for="row in customCategories" :key="row.attachment_category_id" class="asset-page__category-row">
            <button class="asset-page__category" :class="{ active: categoryID === row.attachment_category_id }" @click="selectCategory(row.attachment_category_id)"><span class="asset-page__caret">›</span>{{ row.attachment_category_name }}</button>
            <span v-if="canManage" class="asset-page__category-actions"><ElButton link size="small" @click="editCategory(row)">编辑</ElButton><ElButton link size="small" type="danger" @click="removeCategory(row)">删除</ElButton></span>
          </div>
          <div class="asset-page__category-section">系统分组</div>
          <button class="asset-page__category" :class="{ active: libraryMode === 'system' }" @click="selectCategory(SYSTEM_ROOT_ID)">系统素材</button>
          <button v-for="row in systemCategories" :key="row.attachment_category_id" class="asset-page__category asset-page__category--child" :class="{ active: categoryID === row.attachment_category_id && libraryMode === 'all' }" @click="selectCategory(row.attachment_category_id)">{{ row.attachment_category_name }}</button>
        </aside>
        <main class="asset-page__main">
          <div class="asset-page__toolbar">
            <div class="asset-page__actions">
              <ElUpload v-if="canManage" :accept="accept" :http-request="upload" :show-file-list="false"><ElButton type="primary">上传{{ kind === 'image' ? '图片' : '视频' }}</ElButton></ElUpload>
              <ElButton v-if="canManage" :disabled="!selectedCount" type="danger" plain @click="removeFiles()">删除{{ kind === 'image' ? '图片' : '视频' }}</ElButton>
              <ElSelect v-if="canManage" v-model="moveTarget" class="asset-page__move" clearable placeholder="素材移动至" :disabled="!selectedCount" @change="moveSelected"><ElOption v-for="row in categories" :key="row.attachment_category_id" :label="row.attachment_category_name" :value="row.attachment_category_id" :disabled="Number(row.is_system) === 1" /></ElSelect>
            </div>
            <div class="asset-page__tools">
              <ElInput v-model="keyword" clearable placeholder="请输入素材名称搜索" @keyup.enter="search"><template #suffix><Search class="cursor-pointer" @click="search" /></template></ElInput>
              <div class="asset-page__view-switch"><ElButton :type="viewMode === 'grid' ? 'primary' : 'default'" plain :icon="Grid" @click="viewMode = 'grid'" /><ElButton :type="viewMode === 'list' ? 'primary' : 'default'" plain :icon="List" @click="viewMode = 'list'" /></div>
            </div>
          </div>
          <div class="asset-page__meta"><ElCheckbox :model-value="allSelected" :indeterminate="selectedCount > 0 && !allSelected" @change="toggleAll(Boolean($event))">全选</ElCheckbox><span>{{ activeCategoryName }}</span><span v-if="selectedCount">已选 {{ selectedCount }} 个</span></div>
          <div v-loading="loading" class="asset-page__content" :class="`is-${viewMode}`">
            <ElEmpty v-if="!loading && files.length === 0" description="暂无素材，请上传" />
            <template v-else-if="viewMode === 'grid'">
              <article v-for="row in files" :key="row.attachment_id" class="asset-card" :class="{ selected: selectedIDs.includes(row.attachment_id) }" @click="toggleSelected(row.attachment_id, !selectedIDs.includes(row.attachment_id))">
                <ElCheckbox class="asset-card__checkbox" :model-value="selectedIDs.includes(row.attachment_id)" @click.stop @change="toggleSelected(row.attachment_id, Boolean($event))" />
                <img v-if="row.attachment_type === 0" :src="row.attachment_src" :alt="row.attachment_name" />
                <div v-else class="asset-card__video"><video :src="row.attachment_src" preload="metadata" /><VideoPlay /></div>
                <p :title="row.attachment_name">{{ row.attachment_name }}</p>
              </article>
            </template>
            <template v-else>
              <article v-for="row in files" :key="row.attachment_id" class="asset-list-row" :class="{ selected: selectedIDs.includes(row.attachment_id) }">
                <ElCheckbox :model-value="selectedIDs.includes(row.attachment_id)" @change="toggleSelected(row.attachment_id, Boolean($event))" />
                <img v-if="row.attachment_type === 0" :src="row.attachment_src" :alt="row.attachment_name" /><div v-else class="asset-list-row__video"><VideoPlay /></div>
                <strong>{{ row.attachment_name }}</strong><span>{{ row.attachment_type === 0 ? '图片' : '视频' }}</span><span>{{ row.create_time }}</span><ElButton v-if="canManage" link type="danger" :icon="Delete" @click="removeFiles([row.attachment_id])">删除</ElButton>
              </article>
            </template>
          </div>
          <footer class="asset-page__pager"><span>共 {{ total }} 条</span><ElPagination v-model:current-page="page" :page-size="pageSize" :total="total" background layout="prev, pager, next, jumper" @current-change="loadFiles" /></footer>
        </main>
      </div>
    </div>
  </Page>
</template>

<style scoped lang="scss">
.asset-page { min-height: 720px; background: hsl(var(--background)); }
.asset-page__tabs { display: flex; gap: 40px; height: 62px; padding: 0 28px; border-bottom: 1px solid hsl(var(--border)); }
.asset-page__tabs button { position: relative; border: 0; color: hsl(var(--muted-foreground)); font-size: 18px; font-weight: 600; background: transparent; cursor: pointer; }
.asset-page__tabs button.active { color: hsl(var(--primary)); }
.asset-page__tabs button.active::after { position: absolute; right: 0; bottom: 0; left: 0; height: 3px; background: hsl(var(--primary)); content: ''; }
.asset-page__body { display: flex; min-height: 658px; }
.asset-page__sidebar { flex: 0 0 240px; padding: 22px 16px; border-right: 1px solid hsl(var(--border)); overflow: auto; }
.asset-page__category-row { display: flex; align-items: center; min-height: 42px; }
.asset-page__category { flex: 1; overflow: hidden; padding: 10px 14px; border: 0; border-radius: 2px; color: hsl(var(--foreground)); text-align: left; text-overflow: ellipsis; white-space: nowrap; background: transparent; cursor: pointer; }
.asset-page__category:hover, .asset-page__category.active { color: hsl(var(--primary)); background: hsl(var(--accent)); }
.asset-page__category--child { padding-left: 38px; font-size: 13px; }
.asset-page__caret { display: inline-block; width: 16px; color: hsl(var(--muted-foreground)); font-size: 18px; transform: rotate(90deg); }
.asset-page__category-section { margin: 20px 14px 6px; color: hsl(var(--muted-foreground)); font-size: 12px; }
.asset-page__category-actions { display: none; white-space: nowrap; }
.asset-page__category-row:hover .asset-page__category-actions { display: inline-flex; }
.asset-page__main { display: flex; flex: 1; flex-direction: column; min-width: 0; padding: 28px 34px 22px; }
.asset-page__toolbar, .asset-page__actions, .asset-page__tools, .asset-page__meta, .asset-page__pager { display: flex; align-items: center; }
.asset-page__toolbar { justify-content: space-between; gap: 20px; }
.asset-page__actions, .asset-page__tools { gap: 12px; }
.asset-page__move { width: 180px; }
.asset-page__tools :deep(.el-input) { width: 250px; }
.asset-page__view-switch { display: flex; }
.asset-page__view-switch :deep(.el-button) { width: 42px; margin-left: -1px; border-radius: 0; }
.asset-page__view-switch :deep(.el-button:first-child) { border-radius: 4px 0 0 4px; }
.asset-page__view-switch :deep(.el-button:last-child) { border-radius: 0 4px 4px 0; }
.asset-page__meta { gap: 16px; min-height: 44px; color: hsl(var(--muted-foreground)); font-size: 13px; }
.asset-page__content { flex: 1; min-height: 430px; }
.asset-page__content.is-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 22px 18px; align-content: start; }
.asset-card { position: relative; min-width: 0; cursor: pointer; }
.asset-card img, .asset-card__video { display: block; width: 100%; aspect-ratio: 1; object-fit: cover; background: hsl(var(--muted)); }
.asset-card__video { position: relative; overflow: hidden; color: #fff; }
.asset-card__video video { width: 100%; height: 100%; object-fit: cover; opacity: .58; }
.asset-card__video svg { position: absolute; top: 50%; left: 50%; width: 42px; height: 42px; transform: translate(-50%, -50%); }
.asset-card__checkbox { position: absolute; top: 8px; right: 8px; z-index: 1; opacity: 0; }
.asset-card:hover .asset-card__checkbox, .asset-card.selected .asset-card__checkbox { opacity: 1; }
.asset-card.selected::after { position: absolute; inset: 0 0 28px; border: 2px solid hsl(var(--primary)); pointer-events: none; content: ''; }
.asset-card p { overflow: hidden; margin: 8px 0 0; color: hsl(var(--muted-foreground)); text-align: center; text-overflow: ellipsis; white-space: nowrap; }
.asset-page__content.is-list { display: flex; flex-direction: column; align-content: stretch; gap: 8px; }
.asset-list-row { display: grid; grid-template-columns: 28px 62px minmax(160px, 1fr) 90px 190px 56px; align-items: center; gap: 16px; min-height: 78px; padding: 8px 14px; border: 1px solid transparent; background: hsl(var(--muted) / .35); }
.asset-list-row.selected { border-color: hsl(var(--primary)); background: hsl(var(--accent)); }
.asset-list-row img, .asset-list-row__video { width: 60px; height: 60px; object-fit: cover; border-radius: 3px; background: hsl(var(--muted)); }
.asset-list-row__video { display: grid; place-items: center; color: hsl(var(--muted-foreground)); }
.asset-list-row span { color: hsl(var(--muted-foreground)); font-size: 13px; }
.asset-page__pager { justify-content: flex-end; gap: 18px; min-height: 68px; color: hsl(var(--muted-foreground)); }
@media (max-width: 1000px) { .asset-page__sidebar { flex-basis: 190px; } .asset-page__toolbar { align-items: flex-start; flex-direction: column; } }
</style>
