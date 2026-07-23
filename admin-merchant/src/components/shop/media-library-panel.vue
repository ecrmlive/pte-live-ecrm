<script setup lang="ts">
import type { ShopFileGroupItem, ShopFileItem } from '#/api/core/file';

import {
  CaretBottom,
  Delete,
  Edit,
  Plus,
  Upload as UploadIcon,
} from '@element-plus/icons-vue';
import {
  ElButton,
  ElDropdown,
  ElDropdownItem,
  ElDropdownMenu,
  ElEmpty,
  ElLoading,
  ElMessage,
  ElMessageBox,
  ElPagination,
  ElUpload,
} from 'element-plus';
import { computed, onMounted, ref, watch } from 'vue';

import {
  deleteShopFileGroupApi,
  deleteShopFilesApi,
  getShopFileCategoryApi,
  getShopFileListApi,
  getShopSystemImageCategoryApi,
  getShopSystemImageListApi,
  moveShopFilesApi,
  uploadShopFileDirect,
} from '#/api/core/file';

/** Sidebar entry that opens platform system icon bank (`/shop/file.image/*`). */
const SYSTEM_LIBRARY_GROUP_ID = -1;
/** Return from system icon bank to merchant upload groups. */
const BACK_TO_MERCHANT_GROUP_ID = -2;

import FileCategoryDialog from './file-category-dialog.vue';

const props = withDefaults(
  defineProps<{
    fileType: 'image' | 'video';
    /** 弹窗选择模式：单击受 selectionLimit 约束，可通过 expose 读取已选 */
    pickerMode?: boolean;
    selectionLimit?: number;
    /** 图片选择器：侧栏展示「系统素材」入口（COS pte-live/system/*） */
    enableSystemLibrary?: boolean;
    /** 弹窗打开 / resetPickerState 时的默认素材库 */
    defaultLibrary?: 'merchant' | 'system';
  }>(),
  {
    pickerMode: false,
    selectionLimit: 1,
    enableSystemLibrary: false,
    defaultLibrary: 'merchant',
  },
);

const loading = ref(false);
const uploading = ref(false);
const libraryMode = ref<'merchant' | 'system'>(props.defaultLibrary);
const merchantTypeList = ref<ShopFileGroupItem[]>([{ group_id: null, group_name: '全部' }]);
const systemTypeList = ref<ShopFileGroupItem[]>([{ group_id: null, group_name: '全部' }]);
const activeGroupId = ref<null | number>(null);

const isSystemLibrary = computed(() => libraryMode.value === 'system');

const sidebarTypeList = computed(() => {
  if (isSystemLibrary.value) {
    return [
      { group_id: BACK_TO_MERCHANT_GROUP_ID, group_name: '我的素材' },
      ...systemTypeList.value,
    ];
  }
  return merchantTypeList.value;
});

const showSystemLibraryButton = computed(
  () => props.enableSystemLibrary && props.fileType === 'image',
);

const typeList = computed(() => sidebarTypeList.value);

const showManageToolbar = computed(() => !isSystemLibrary.value);
const fileList = ref<ShopFileItem[]>([]);
const total = ref(0);
const curPage = ref(1);
const pageSize = ref(36);
const categoryDialogOpen = ref(false);
const editingCategory = ref<ShopFileGroupItem | null>(null);

const accept = computed(() =>
  props.fileType === 'image'
    ? 'image/jpeg,image/png,image/jpg'
    : 'video/mp4,video/quicktime,.mp4,.mov',
);

const uploadLabel = computed(() =>
  props.fileType === 'video' ? '视频上传中,请等待' : '图片上传中,请等待',
);

const selectedCount = computed(
  () => fileList.value.filter((item) => item.selected).length,
);

function onCategorySaved() {
  void loadCategories();
  void loadFiles();
}

async function loadMerchantCategories() {
  try {
    const res = await getShopFileCategoryApi(props.fileType);
    merchantTypeList.value = [
      { group_id: null, group_name: '全部' },
      ...(res.group_list ?? []),
    ];
  } catch {
    merchantTypeList.value = [{ group_id: null, group_name: '全部' }];
  }
}

async function loadSystemCategories() {
  try {
    const res = await getShopSystemImageCategoryApi();
    systemTypeList.value = [
      { group_id: null, group_name: '全部' },
      ...(res.list ?? []).map((item) => ({
        group_id: item.category_id,
        group_name: item.name,
      })),
    ];
  } catch {
    systemTypeList.value = [{ group_id: null, group_name: '全部' }];
  }
}

async function loadCategories() {
  await loadMerchantCategories();
  if (props.enableSystemLibrary && props.fileType === 'image') {
    await loadSystemCategories();
  }
}

async function loadMerchantFiles() {
  const res = await getShopFileListApi({
    group_id: activeGroupId.value,
    list_rows: pageSize.value,
    page: curPage.value,
    type: props.fileType,
  });
  fileList.value = (res.file_list.data ?? []).map((item) => ({
    ...item,
    selected: false,
  }));
  total.value = res.file_list.total ?? 0;
}

async function loadSystemFiles() {
  const res = await getShopSystemImageListApi({
    list_rows: pageSize.value,
    page: curPage.value,
    parentId: activeGroupId.value,
  });
  fileList.value = (res.list?.data ?? []).map((item) => ({
    file_id: item.category_id,
    file_path: item.image,
    file_name: item.name,
    real_name: item.name,
    selected: false,
  }));
  total.value = res.list?.total ?? 0;
}

async function loadFiles() {
  loading.value = true;
  try {
    if (isSystemLibrary.value) {
      await loadSystemFiles();
      return;
    }
    await loadMerchantFiles();
  } finally {
    loading.value = false;
  }
}

function resetPickerState() {
  libraryMode.value = props.defaultLibrary;
  activeGroupId.value = null;
  curPage.value = 1;
  clearPickerSelection();
  void loadCategories();
  void loadFiles();
}

function selectGroup(groupId: null | number) {
  if (groupId === BACK_TO_MERCHANT_GROUP_ID) {
    libraryMode.value = 'merchant';
    activeGroupId.value = null;
    curPage.value = 1;
    clearPickerSelection();
    void loadFiles();
    return;
  }
  if (groupId === SYSTEM_LIBRARY_GROUP_ID) {
    libraryMode.value = 'system';
    activeGroupId.value = null;
    curPage.value = 1;
    clearPickerSelection();
    void loadFiles();
    return;
  }
  activeGroupId.value = groupId;
  curPage.value = 1;
  void loadFiles();
}

function isGroupActive(groupId: null | number) {
  if (groupId === BACK_TO_MERCHANT_GROUP_ID || groupId === SYSTEM_LIBRARY_GROUP_ID) {
    return false;
  }
  return activeGroupId.value === groupId;
}

function canManageCategory(groupId: null | number) {
  return (
    groupId != null &&
    groupId !== SYSTEM_LIBRARY_GROUP_ID &&
    groupId !== BACK_TO_MERCHANT_GROUP_ID &&
    !isSystemLibrary.value
  );
}

function clearPickerSelection() {
  fileList.value.forEach((row) => {
    row.selected = false;
  });
}

function getSelectedFiles() {
  return fileList.value.filter((row) => row.selected);
}

function toggleSelect(item: ShopFileItem) {
  if (props.pickerMode) {
    const alreadySelected = !!item.selected;
    if (props.selectionLimit === 1) {
      clearPickerSelection();
      item.selected = !alreadySelected;
      return;
    }
    if (alreadySelected) {
      item.selected = false;
      return;
    }
    const pickedCount = fileList.value.filter((row) => row.selected).length;
    if (pickedCount >= props.selectionLimit) {
      ElMessage.warning(`最多选择 ${props.selectionLimit} 个文件`);
      return;
    }
    item.selected = true;
    return;
  }
  item.selected = !item.selected;
}

defineExpose({
  clearPickerSelection,
  getSelectedFiles,
  loadFiles,
  resetPickerState,
});

function openAddCategory() {
  editingCategory.value = null;
  categoryDialogOpen.value = true;
}

function openEditCategory(item: ShopFileGroupItem) {
  editingCategory.value = item;
  categoryDialogOpen.value = true;
}

async function deleteCategory(groupId: number) {
  try {
    await ElMessageBox.confirm('此操作将永久删除该分类, 是否继续?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
    await deleteShopFileGroupApi(groupId);
    ElMessage.success('删除成功');
    if (activeGroupId.value === groupId) {
      activeGroupId.value = null;
    }
    await loadCategories();
    await loadFiles();
  } catch {
    // cancelled or failed
  }
}

async function moveSelected(groupId: null | number | string) {
  const targetGroupId =
    groupId === 'all' || groupId === '' || groupId == null ? null : Number(groupId);
  const fileIds = fileList.value
    .filter((item) => item.selected)
    .map((item) => item.file_id);
  if (!fileIds.length) {
    ElMessage.warning('请选择文件');
    return;
  }
  try {
    await ElMessageBox.confirm('确定移动选中的文件吗?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
    await moveShopFilesApi({
      fileIds,
      group_id: targetGroupId,
    });
    ElMessage.success('移动成功');
    await loadCategories();
    await loadFiles();
  } catch {
    // cancelled
  }
}

async function deleteFiles(single?: ShopFileItem) {
  const fileIds = single
    ? [single.file_id]
    : fileList.value.filter((item) => item.selected).map((item) => item.file_id);
  if (!fileIds.length) {
    ElMessage.warning('请选择文件');
    return;
  }
  try {
    await ElMessageBox.confirm('此操作将永久删除该记录, 是否继续?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
    const loadingSvc = ElLoading.service({
      background: 'rgba(0, 0, 0, 0.7)',
      lock: true,
      text: '处理中,请等待',
    });
    try {
      await deleteShopFilesApi(fileIds);
      ElMessage.success('删除成功');
      await loadFiles();
    } finally {
      loadingSvc.close();
    }
  } catch {
    // cancelled
  }
}

async function onUpload(file: File) {
  const loadingSvc = ElLoading.service({
    background: 'rgba(0, 0, 0, 0.7)',
    lock: true,
    text: uploadLabel.value,
  });
  uploading.value = true;
  try {
    await uploadShopFileDirect(file, {
      file_type: props.fileType,
      group_id: activeGroupId.value ?? 0,
    });
    ElMessage.success('上传成功');
    await loadFiles();
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '上传失败');
  } finally {
    uploading.value = false;
    loadingSvc.close();
  }
  return false;
}

function onPageChange(page: number) {
  curPage.value = page;
  void loadFiles();
}

function onSizeChange(size: number) {
  pageSize.value = size;
  curPage.value = 1;
  void loadFiles();
}

watch(
  () => props.fileType,
  () => {
    libraryMode.value = 'merchant';
    activeGroupId.value = null;
    curPage.value = 1;
    void loadCategories();
    void loadFiles();
  },
);

onMounted(() => {
  void loadCategories();
  void loadFiles();
});
</script>

<template>
  <div class="media-library-panel" :class="{ 'media-library-panel--picker': pickerMode }">
    <div v-if="showManageToolbar" class="media-library-panel__toolbar">
      <div class="media-library-panel__toolbar-left">
        <ElButton :icon="Plus" type="default" @click="openAddCategory">添加分类</ElButton>
        <ElDropdown trigger="click" @command="moveSelected">
          <ElButton :icon="CaretBottom" type="default">
            移动至
          </ElButton>
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem
                v-for="item in typeList"
                :key="`move-${String(item.group_id)}`"
                :command="item.group_id ?? 'all'"
              >
                {{ item.group_name }}
              </ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
        <ElButton :icon="Delete" type="danger" @click="deleteFiles()">
          批量删除
          <span v-if="selectedCount">({{ selectedCount }})</span>
        </ElButton>
      </div>
      <ElUpload
        class="media-library-panel__toolbar-upload"
        :accept="accept"
        :before-upload="onUpload"
        :disabled="uploading"
        :show-file-list="false"
        multiple
      >
        <ElButton :icon="UploadIcon" :loading="uploading" type="primary">
          点击上传
        </ElButton>
      </ElUpload>
    </div>

    <div class="media-library-panel__body">
      <aside class="media-library-panel__groups">
        <ul class="media-library-panel__group-list">
          <li
            v-for="item in sidebarTypeList"
            :key="String(item.group_id)"
            :class="{
              active: isGroupActive(item.group_id),
              'is-back-entry': item.group_id === BACK_TO_MERCHANT_GROUP_ID,
            }"
            @click="selectGroup(item.group_id)"
          >
            <span>{{ item.group_name }}</span>
            <div
              v-if="canManageCategory(item.group_id)"
              class="media-library-panel__group-actions"
              @click.stop
            >
              <ElButton
                class="media-library-panel__group-action-btn"
                link
                type="primary"
                @click="openEditCategory(item)"
              >
                <el-icon><Edit /></el-icon>
              </ElButton>
              <ElButton
                class="media-library-panel__group-action-btn"
                link
                type="danger"
                @click="deleteCategory(item.group_id!)"
              >
                <el-icon><Delete /></el-icon>
              </ElButton>
            </div>
          </li>
        </ul>
        <div v-if="showSystemLibraryButton" class="media-library-panel__groups-footer">
          <ElButton
            class="media-library-panel__system-btn"
            :type="isSystemLibrary ? 'primary' : 'default'"
            @click="selectGroup(SYSTEM_LIBRARY_GROUP_ID)"
          >
            系统素材
          </ElButton>
        </div>
      </aside>

      <div class="media-library-panel__main">
        <div v-loading="loading" class="media-library-panel__grid-wrap">
          <ul v-if="fileList.length" class="media-library-panel__grid">
            <li
              v-for="(item, index) in fileList"
              :key="item.file_id ?? index"
              :class="{ 'is-selected': item.selected }"
              class="media-library-panel__item"
              @click="toggleSelect(item)"
            >
              <div v-if="item.selected" class="media-library-panel__selected">
                <span>✓</span>
              </div>
              <img
                v-if="fileType === 'image'"
                :src="item.file_path"
                alt=""
                class="media-library-panel__thumb"
              />
              <video
                v-else
                :src="item.file_path"
                class="media-library-panel__thumb"
                muted
                preload="metadata"
              />
              <p class="media-library-panel__name" :title="item.real_name || item.file_name">
                {{ item.real_name || item.file_name || '-' }}
              </p>
              <ElButton
                v-if="!isSystemLibrary"
                class="media-library-panel__delete"
                link
                type="danger"
                @click.stop="deleteFiles(item)"
              >
                <el-icon><Delete /></el-icon>
              </ElButton>
            </li>
          </ul>
          <ElEmpty
            v-else
            :description="
              isSystemLibrary ? '暂无系统图片' : '暂无文件，请点击右上角上传'
            "
          />
        </div>

        <div class="media-library-panel__pager">
          <ElPagination
            v-model:current-page="curPage"
            background
            layout="total, sizes, prev, pager, next, jumper"
            :page-size="pageSize"
            :page-sizes="[12, 24, 36, 42, 48]"
            :total="total"
            @current-change="onPageChange"
            @size-change="onSizeChange"
          />
        </div>
      </div>
    </div>

    <FileCategoryDialog
      v-model:open="categoryDialogOpen"
      :category="editingCategory"
      :file-type="fileType"
      @success="onCategorySaved"
    />
  </div>
</template>

<style lang="scss">
/* Unscoped: must apply inside Vben Modal teleports and legacy pages alike */
.media-library-panel {
  --ml-card-bg: hsl(var(--card));
  --ml-border: hsl(var(--border));
  --ml-muted: hsl(var(--muted-foreground));
  --ml-fg: hsl(var(--foreground));
  --ml-primary: hsl(var(--primary));
  --ml-surface: hsl(var(--muted) / 18%);
  --ml-shadow: 0 1px 4px hsl(var(--foreground) / 0.04);

  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 420px;
}

.media-library-panel__toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 12px 16px;
  background: var(--ml-card-bg);
  border: 1px solid var(--ml-border);
  border-radius: 10px;
  box-shadow: var(--ml-shadow);

  /* Card toolbar shares bg with default ElButton in dark theme — restore visible chrome */
  --el-button-size: 32px;

  .el-button {
    font-weight: 500;
  }

  .el-button--default {
    --el-button-bg-color: hsl(var(--background));
    --el-button-border-color: hsl(var(--border));
    --el-button-text-color: hsl(var(--foreground));
    --el-button-hover-bg-color: hsl(var(--accent));
    --el-button-hover-border-color: hsl(var(--primary) / 40%);
    --el-button-hover-text-color: hsl(var(--foreground));
    --el-button-active-bg-color: hsl(var(--accent));
    --el-button-active-border-color: hsl(var(--primary) / 55%);
    --el-button-active-text-color: hsl(var(--foreground));
  }

  .el-button--primary {
    --el-button-bg-color: hsl(var(--primary));
    --el-button-border-color: hsl(var(--primary));
    --el-button-text-color: hsl(var(--primary-foreground));
    --el-button-hover-bg-color: hsl(var(--primary) / 88%);
    --el-button-hover-border-color: hsl(var(--primary) / 88%);
    --el-button-hover-text-color: hsl(var(--primary-foreground));
    --el-button-active-bg-color: hsl(var(--primary) / 78%);
    --el-button-active-border-color: hsl(var(--primary) / 78%);
    --el-button-active-text-color: hsl(var(--primary-foreground));
  }

  .el-button--danger {
    --el-button-bg-color: hsl(var(--destructive));
    --el-button-border-color: hsl(var(--destructive));
    --el-button-text-color: hsl(var(--destructive-foreground));
    --el-button-hover-bg-color: hsl(var(--destructive) / 88%);
    --el-button-hover-border-color: hsl(var(--destructive) / 88%);
    --el-button-hover-text-color: hsl(var(--destructive-foreground));
    --el-button-active-bg-color: hsl(var(--destructive) / 78%);
    --el-button-active-border-color: hsl(var(--destructive) / 78%);
    --el-button-active-text-color: hsl(var(--destructive-foreground));
  }
}

.media-library-panel__toolbar-left {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.media-library-panel__toolbar-upload {
  flex-shrink: 0;
}

.media-library-panel__body {
  display: flex;
  min-height: 380px;
  background: var(--ml-card-bg);
  border: 1px solid var(--ml-border);
  border-radius: 10px;
  box-shadow: var(--ml-shadow);
  overflow: hidden;
}

.media-library-panel__groups {
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  width: 128px;
  border-right: 1px solid var(--ml-border);
  overflow: hidden;
  background: hsl(var(--muted) / 10%);
}

.media-library-panel__group-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  list-style: none;
  margin: 0;
  padding: 8px 0;
}

.media-library-panel__groups-footer {
  flex-shrink: 0;
  padding: 10px 8px;
  border-top: 1px solid var(--ml-border);
  background: hsl(var(--card));
}

.media-library-panel__system-btn {
  width: 100%;
  font-weight: 600;
}

.media-library-panel__main {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
}

.media-library-panel__grid-wrap {
  flex: 1;
  min-height: 0;
  min-width: 0;
  overflow-x: hidden;
  overflow-y: auto;
  padding: 12px 16px;
}

.media-library-panel__groups li {
  position: relative;
  padding: 10px 24px 10px 10px;
  cursor: pointer;
  text-align: center;
  font-size: 13px;
  min-height: 20px;
  color: var(--ml-muted);
  transition:
    background-color 0.2s ease,
    color 0.2s ease;
}

.media-library-panel__groups li:hover {
  color: var(--ml-fg);
  background: hsl(var(--muted) / 16%);
}

.media-library-panel__groups li.active {
  background: hsl(var(--primary) / 12%);
  color: var(--ml-primary);
  font-weight: 600;
}

.media-library-panel__groups li:hover .media-library-panel__group-actions {
  display: flex;
}

.media-library-panel__group-actions {
  display: none;
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 22px;
  flex-direction: column;
  background: hsl(var(--foreground) / 55%);
}

.media-library-panel__group-action-btn {
  flex: 1;
  border: 0;
  background: transparent;
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  height: auto;
  min-height: 0;

  &.el-button.is-link {
    color: #fff;
  }
}

.media-library-panel__pager {
  display: flex;
  flex-shrink: 0;
  justify-content: flex-end;
  padding: 8px 12px 10px;
  border-top: 1px solid var(--ml-border);
  background: var(--ml-card-bg);
}

/* 弹窗选择模式：固定高度，仅图片网格区域滚动 */
.media-library-panel--picker {
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 520px;
  max-height: calc(72vh - 120px);
  min-height: 420px;
  overflow: hidden;
}

.media-library-panel--picker .media-library-panel__toolbar {
  flex-shrink: 0;
}

.media-library-panel--picker .media-library-panel__body {
  flex: 1;
  min-height: 0;
}

.media-library-panel--picker .media-library-panel__group-list {
  overflow-y: auto;
}

.media-library-panel__groups li.is-back-entry {
  font-weight: 600;
}

.media-library-panel__grid {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.media-library-panel__item {
  position: relative;
  width: 100px;
  cursor: pointer;
}

.media-library-panel__item.is-selected .media-library-panel__thumb {
  outline: 2px solid var(--ml-primary);
}

.media-library-panel__selected {
  position: absolute;
  top: 4px;
  right: 4px;
  z-index: 2;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--ml-primary);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.media-library-panel__thumb {
  width: 100px;
  height: 100px;
  border-radius: 6px;
  object-fit: cover;
  background: var(--ml-surface);
  display: block;
  border: 1px solid var(--ml-border);
}

.media-library-panel__name {
  margin: 6px 0 0;
  font-size: 12px;
  line-height: 1.35;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--ml-muted);
}

.media-library-panel__delete {
  position: absolute;
  right: 4px;
  bottom: 28px;
  border: 0;
  background: hsl(var(--foreground) / 55%);
  color: #fff;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  cursor: pointer;
  display: none;
  align-items: center;
  justify-content: center;
  padding: 0;
  min-height: 0;

  &.el-button.is-link {
    color: #fff;
  }
}

.media-library-panel__item:hover .media-library-panel__delete {
  display: flex;
}
</style>
