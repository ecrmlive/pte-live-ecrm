<script lang="ts" setup>
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import {
  ElButton,
  ElIcon,
  ElLoading,
  ElMessage,
  ElMessageBox,
  ElUpload,
} from 'element-plus';
import { Check, Delete, Edit, Plus, Upload } from '@element-plus/icons-vue';
import type { UploadRequestOptions } from 'element-plus';

import FileApi from '#/api/core/file';
import PlatformListPager from '#/components/platform-list/PlatformListPager.vue';

import CategoryModal from './category-modal.vue';
import type { PictureCategory, PictureFile } from './types';

import './system-picture.scss';

const fileType = 'image';
const accept = 'image/jpeg,image/png,image/jpg';

const loading = ref(false);
const activeType = ref(0);
const typeList = ref<PictureCategory[]>([]);
const fileList = ref<PictureFile[]>([]);
const curPage = ref(1);
const pageSize = ref(36);
const totalDataNumber = ref(0);

const categoryModalOpen = ref(false);
const editingCategory = ref<null | PictureCategory>(null);

async function loadCategories() {
  try {
    const res = await FileApi.PictureIndex({ type: fileType }, true);
    const list = (res.data as { list?: PictureCategory[] })?.list ?? [];
    typeList.value = [{ category_id: 0, name: '全部' }, ...list];
  } catch {
    typeList.value = [{ category_id: 0, name: '全部' }];
  }
}

async function loadFiles() {
  loading.value = true;
  try {
    const res = await FileApi.SystemPictureList(
      {
        fileType,
        list_rows: pageSize.value,
        page: curPage.value,
        parentId: activeType.value,
      },
      true,
    );
    const page = (res.data as { list?: { data?: PictureFile[]; total?: number } })
      ?.list;
    fileList.value = (page?.data ?? []).map((item) => ({
      ...item,
      selected: false,
    }));
    totalDataNumber.value = page?.total ?? 0;
  } catch {
    fileList.value = [];
    totalDataNumber.value = 0;
  } finally {
    loading.value = false;
  }
}

function selectType(id: number) {
  activeType.value = id;
  curPage.value = 1;
  void loadFiles();
}

function openAddCategory() {
  editingCategory.value = null;
  categoryModalOpen.value = true;
}

function openEditCategory(item: PictureCategory) {
  editingCategory.value = item;
  categoryModalOpen.value = true;
}

async function refreshAll() {
  await loadCategories();
  await loadFiles();
}

async function deleteCategory(id: number) {
  try {
    await ElMessageBox.confirm('此操作将永久删除该记录, 是否继续?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
  } catch {
    return;
  }

  try {
    await FileApi.deleteCategory({ category_id: id });
    ElMessage.success('删除成功');
    activeType.value = 0;
    await refreshAll();
  } catch {
    ElMessage.error('删除失败');
  }
}

function toggleSelect(item: PictureFile, index: number) {
  const next = [...fileList.value];
  next[index] = { ...item, selected: !item.selected };
  fileList.value = next;
}

async function deleteFiles(single?: PictureFile) {
  try {
    await ElMessageBox.confirm('此操作将永久删除该记录, 是否继续?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
  } catch {
    return;
  }

  const ids: number[] = [];
  if (single) {
    ids.push(single.category_id);
  } else {
    fileList.value.forEach((item) => {
      if (item.selected) {
        ids.push(item.category_id);
      }
    });
  }

  if (ids.length === 0) {
    ElMessage.warning('请选择文件');
    return;
  }

  const overlay = ElLoading.service({
    background: 'rgba(0, 0, 0, 0.7)',
    lock: true,
    text: '正在删除',
  });

  try {
    const res = await FileApi.deleteFiles({ imageIds: ids.join(',') }, true);
    if ((res as { code?: number }).code === 1) {
      ElMessage.success('删除成功');
      await loadFiles();
    }
  } finally {
    overlay.close();
  }
}

function handlePageSizeChange() {
  curPage.value = 1;
  void loadFiles();
}

async function uploadImage(options: UploadRequestOptions) {
  if (!activeType.value) {
    ElMessage.error('请选择分类');
    return;
  }

  const overlay = ElLoading.service({
    background: 'rgba(0, 0, 0, 0.7)',
    lock: true,
    text: '文件上传中,请等待',
  });

  try {
    await FileApi.uploadCosDirect(options.file as File, {
      category_id: activeType.value,
      file_type: 'image',
    });
    ElMessage.success('上传文件成功');
    await loadFiles();
    options.onSuccess?.({});
  } catch (error) {
    const err = error as { message?: string; msg?: string };
    ElMessage.warning(err?.msg || err?.message || '本次上传文件失败');
    options.onError?.(error as Error);
  } finally {
    overlay.close();
  }
}

onMounted(async () => {
  await loadCategories();
  await loadFiles();
});
</script>

<template>
  <Page
    auto-content-height
    content-class="system-picture-page flex flex-col overflow-hidden min-h-0"
  >
    <div class="system-picture__toolbar shrink-0">
      <div class="flex flex-wrap gap-2">
        <ElButton
          v-access:code="'platform:picture:addCategory'"
          :icon="Plus"
          @click="openAddCategory"
        >
          新增分类
        </ElButton>
        <ElButton
          v-access:code="'platform:picture:deleteFiles'"
          :icon="Delete"
          type="danger"
          @click="deleteFiles()"
        >
          批量删除
        </ElButton>
      </div>
      <ElUpload
        v-access:code="'platform:picture:upload'"
        :accept="accept"
        :http-request="uploadImage"
        :show-file-list="false"
        multiple
      >
        <ElButton :icon="Upload" type="primary">点击上传</ElButton>
      </ElUpload>
    </div>

    <div v-loading="loading" class="system-picture__body flex-1 min-h-0">
      <aside class="system-picture__sidebar">
        <ul>
          <li
            v-for="item in typeList"
            :key="item.category_id"
            :class="[
              'system-picture__category',
              activeType === item.category_id ? 'is-active' : '',
            ]"
            @click="selectType(item.category_id)"
          >
            <span>{{ item.name }}</span>
            <div
              v-if="item.category_id !== 0"
              class="system-picture__category-actions"
              @click.stop
            >
              <button
                v-access:code="'platform:picture:editCategory'"
                type="button"
                @click="openEditCategory(item)"
              >
                <ElIcon><Edit /></ElIcon>
              </button>
              <button
                v-access:code="'platform:picture:deleteCategory'"
                type="button"
                @click="deleteCategory(item.category_id)"
              >
                <ElIcon><Delete /></ElIcon>
              </button>
            </div>
          </li>
        </ul>
      </aside>

      <div class="system-picture__grid-wrap">
        <ul v-if="fileList.length" class="system-picture__grid">
          <li
            v-for="(item, index) in fileList"
            :key="item.category_id"
            class="system-picture__file"
            @click="toggleSelect(item, index)"
          >
            <span v-if="item.selected" class="system-picture__selected">
              <ElIcon><Check /></ElIcon>
            </span>
            <img
              :src="item.image"
              :style="{ backgroundImage: `url(${item.image})` }"
              alt=""
            />
            <p class="system-picture__file-name">{{ item.name }}</p>
            <div class="system-picture__file-delete" @click.stop>
              <button
                v-access:code="'platform:picture:deleteFiles'"
                type="button"
                @click="deleteFiles(item)"
              >
                <ElIcon><Delete /></ElIcon>
              </button>
            </div>
          </li>
        </ul>
        <div v-else class="system-picture__empty">暂无图片</div>
      </div>
    </div>

    <div class="system-picture__pager shrink-0">
      <PlatformListPager
        v-model:current-page="curPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 36, 42, 48]"
        :total="totalDataNumber"
        @change="loadFiles"
        @size-change="handlePageSizeChange"
      />
    </div>

    <CategoryModal
      v-model:open="categoryModalOpen"
      :category="editingCategory"
      :file-type="fileType"
      @success="refreshAll"
    />
  </Page>
</template>
