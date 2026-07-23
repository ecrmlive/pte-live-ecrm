<script setup lang="ts">
import type { LiveVodVideoItem } from '#/api/core/live-vod';
import type { VbenFormSchema } from '#/adapter/form';

import { useVbenModal } from '@vben/common-ui';
import { Select } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { deleteLiveVodApi, getLiveVodListApi } from '#/api/core/live-vod';
import { uploadVodVideo, type VodUploadProgress } from '#/utils/vod-upload';

const props = withDefaults(
  defineProps<{
    picker?: boolean;
  }>(),
  { picker: false },
);

const emit = defineEmits<{
  select: [LiveVodVideoItem];
  'upload-state': [boolean];
}>();

const loading = ref(false);
const uploading = ref(false);
const uploadDialogVisible = ref(false);
const uploadProgress = ref<VodUploadProgress>({
  fileName: '',
  percent: 0,
  phase: '',
});
const list = ref<LiveVodVideoItem[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 12;
const selectedFileId = ref('');
const uploadInputRef = ref<HTMLInputElement>();

const searchSchema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: 'FileId',
      style: 'width: 200px',
    },
    fieldName: 'keyword',
    label: '搜索',
  },
]);

const [SearchForm, searchFormApi] = useVbenForm(
  reactive({
    actionLayout: 'inline',
    commonConfig: {
      componentProps: { size: 'default' },
    },
    handleSubmit: async () => {
      reload();
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema: searchSchema,
    showDefaultActions: true,
    submitButtonOptions: {
      content: '查询',
    },
    wrapperClass: 'grid-cols-1 md:grid-cols-2',
  }),
);

const uploadPhaseText = computed(() => {
  const map: Record<string, string> = {
    commit: '正在确认入库…',
    cover: '正在上传封面…',
    done: '上传完成',
    preparing: '正在生成封面…',
    uploading: '正在上传视频…',
  };
  return map[uploadProgress.value.phase] || '准备上传…';
});

function displayName(item: LiveVodVideoItem) {
  return item.media_name || item.file_id || '未命名';
}

function formatDuration(sec?: number) {
  const s = Math.max(0, Math.round(Number(sec) || 0));
  const m = Math.floor(s / 60);
  const r = s % 60;
  return `${m}:${String(r).padStart(2, '0')}`;
}

async function loadData() {
  loading.value = true;
  try {
    const values = await searchFormApi.getValues();
    const res = await getLiveVodListApi({
      keyword: String(values.keyword ?? '').trim(),
      list_rows: pageSize,
      page: page.value,
    });
    list.value = res.list?.data ?? [];
    total.value = res.list?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function reload() {
  page.value = 1;
  void loadData();
}

function onPageChange(p: number) {
  page.value = p;
  void loadData();
}

function onItemClick(item: LiveVodVideoItem) {
  if (!props.picker) return;
  selectedFileId.value = item.file_id;
  emit('select', item);
}

function pickUpload() {
  if (uploading.value) return;
  uploadInputRef.value?.click();
}

function onBeforeUnload(e: BeforeUnloadEvent) {
  if (!uploading.value) return;
  e.preventDefault();
  e.returnValue = '';
}

async function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement;
  const file = input.files?.[0];
  input.value = '';
  if (!file || uploading.value) return;
  uploading.value = true;
  uploadDialogVisible.value = true;
  uploadProgress.value = {
    fileName: file.name,
    percent: 0,
    phase: 'preparing',
  };
  try {
    await uploadVodVideo(file, {
      onProgress: (payload) => {
        uploadProgress.value = payload;
      },
    });
    ElMessage.success('上传成功');
    reload();
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '上传失败');
  } finally {
    uploading.value = false;
    uploadDialogVisible.value = false;
    uploadProgress.value = { fileName: '', percent: 0, phase: '' };
  }
}

async function deleteItem(item: LiveVodVideoItem) {
  if (uploading.value) return;
  try {
    await ElMessageBox.confirm(`确认删除「${displayName(item)}」？`, '提示', {
      type: 'warning',
    });
    await deleteLiveVodApi({ file_id: item.file_id });
    ElMessage.success('已删除');
    await loadData();
  } catch {
    // cancelled
  }
}

const [UploadModal, uploadModalApi] = useVbenModal({
  onOpenChange(isOpen) {
    uploadDialogVisible.value = isOpen;
  },
});

watch(uploadDialogVisible, (visible) => {
  if (visible) {
    uploadModalApi.open();
    return;
  }
  uploadModalApi.close();
});

watch(uploading, (val) => {
  emit('upload-state', val);
});

onMounted(() => {
  void loadData();
  window.addEventListener('beforeunload', onBeforeUnload);
});

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', onBeforeUnload);
  emit('upload-state', false);
});

defineExpose({
  clearSelection() {
    selectedFileId.value = '';
  },
  getSelected() {
    return list.value.find((item) => item.file_id === selectedFileId.value) ?? null;
  },
});
</script>

<template>
  <div class="vod-library" :class="{ 'vod-library--picker': picker }">
    <div class="vod-toolbar">
      <SearchForm />
      <div v-if="!picker" class="vod-toolbar__actions">
        <input
          ref="uploadInputRef"
          accept="video/mp4,video/quicktime,.mp4,.mov"
          class="vod-upload-input"
          type="file"
          @change="onFileChange"
        />
        <el-button
          :disabled="uploading"
          :loading="uploading"
          type="primary"
          @click="pickUpload"
        >
          上传视频
        </el-button>
      </div>
    </div>

    <div v-loading="loading" class="vod-grid-wrap">
      <ul v-if="list.length" class="vod-grid">
        <li
          v-for="item in list"
          :key="item.file_id"
          :class="{ 'is-selected': picker && selectedFileId === item.file_id }"
          class="vod-item"
          @click="onItemClick(item)"
        >
          <div
            v-if="picker && selectedFileId === item.file_id"
            class="vod-item__check"
          >
            <el-icon><Select /></el-icon>
          </div>
          <div class="vod-item__cover">
            <img v-if="item.cover_url" :src="item.cover_url" alt="" />
            <div v-else class="vod-item__cover-placeholder">视频</div>
            <span v-if="item.duration" class="vod-item__duration">
              {{ formatDuration(item.duration) }}
            </span>
          </div>
          <p class="vod-item__name" :title="displayName(item)">
            {{ displayName(item) }}
          </p>
          <p class="vod-item__meta">{{ item.create_time_text || '' }}</p>
          <div v-if="!picker" class="vod-item__actions" @click.stop>
            <el-button link size="small" type="danger" @click="deleteItem(item)">
              删除
            </el-button>
          </div>
        </li>
      </ul>
      <el-empty
        v-else
        :description="
          picker ? '暂无云点播视频' : '暂无云点播视频，请点击右上角上传'
        "
      />
    </div>

    <div class="pagination-wrap">
      <el-pagination
        :current-page="page"
        :disabled="uploading"
        background
        layout="total, prev, pager, next"
        :page-size="pageSize"
        size="default"
        :total="total"
        @current-change="onPageChange"
      />
    </div>

    <UploadModal
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
      class="w-[420px]"
      title="上传云点播视频"
    >
      <div class="vod-upload-progress">
        <p class="vod-upload-progress__name" :title="uploadProgress.fileName">
          {{ uploadProgress.fileName || '视频文件' }}
        </p>
        <p class="vod-upload-progress__phase">{{ uploadPhaseText }}</p>
        <el-progress
          :percentage="uploadProgress.percent"
          :status="uploadProgress.percent >= 100 ? 'success' : undefined"
          :stroke-width="10"
        />
      </div>
    </UploadModal>
  </div>
</template>

<style scoped lang="scss">
.vod-library {
  --vod-card-bg: hsl(var(--card));
  --vod-border: hsl(var(--border));
  --vod-muted: hsl(var(--muted-foreground));
  --vod-fg: hsl(var(--foreground));
  --vod-surface: hsl(var(--muted) / 18%);
  --vod-shadow: 0 1px 4px hsl(var(--foreground) / 0.04);
}

.vod-library--picker .vod-toolbar {
  margin-bottom: 8px;
}

.vod-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 12px;
  padding: 12px 16px;
  background: var(--vod-card-bg);
  border: 1px solid var(--vod-border);
  border-radius: 10px;
  box-shadow: var(--vod-shadow);
}

.vod-toolbar__actions {
  flex-shrink: 0;
}

.vod-upload-input {
  display: none;
}

.vod-grid-wrap {
  min-height: 200px;
  padding: 12px 16px;
  background: var(--vod-card-bg);
  border: 1px solid var(--vod-border);
  border-radius: 10px;
  box-shadow: var(--vod-shadow);
}

.vod-grid {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.vod-item {
  width: 140px;
  cursor: default;
  position: relative;
}

.vod-library--picker .vod-item {
  cursor: pointer;
}

.vod-item.is-selected .vod-item__cover {
  outline: 2px solid hsl(var(--primary));
}

.vod-item__check {
  position: absolute;
  top: 4px;
  right: 4px;
  z-index: 2;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: hsl(var(--primary));
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.vod-item__cover {
  position: relative;
  width: 140px;
  height: 100px;
  border-radius: 6px;
  overflow: hidden;
  background: var(--vod-surface);
  border: 1px solid var(--vod-border);
}

.vod-item__cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.vod-item__cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--vod-muted);
  font-size: 13px;
}

.vod-item__duration {
  position: absolute;
  right: 6px;
  bottom: 4px;
  padding: 0 4px;
  font-size: 11px;
  color: #fff;
  background: hsl(var(--foreground) / 55%);
  border-radius: 2px;
}

.vod-item__name {
  margin: 6px 0 0;
  font-size: 12px;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--vod-fg);
}

.vod-item__meta {
  margin: 2px 0 0;
  font-size: 11px;
  color: var(--vod-muted);
}

.vod-item__actions {
  margin-top: 4px;
}

.pagination-wrap {
  margin-top: 8px;
  display: flex;
  justify-content: flex-end;
  padding: 0 4px;
}

.vod-upload-progress__name {
  margin: 0 0 8px;
  font-size: 13px;
  color: var(--vod-fg);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.vod-upload-progress__phase {
  margin: 0 0 12px;
  font-size: 12px;
  color: var(--vod-muted);
}
</style>
