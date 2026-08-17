<script setup lang="ts">
import type { AttachmentCategory, AttachmentItem } from '#/api/core/attachment';

import { ElButton, ElDialog, ElEmpty, ElMessage, ElPagination, ElUpload } from 'element-plus';
import { computed, onMounted, ref, watch } from 'vue';

import {
  listAttachmentCategoriesApi,
  listAttachmentsApi,
  uploadAttachmentApi,
} from '#/api/core/attachment';

type PickerItem = AttachmentItem & { file_id: number; file_path: string };

/** 侧栏「系统素材」根入口（非真实分类 id） */
const SYSTEM_ROOT_ID = -1;

const open = defineModel<boolean>('open', { default: false });
const props = withDefaults(
  defineProps<{
    /** 弹窗初始素材库；DIY 场景默认从系统素材开始选择。 */
    defaultLibrary?: 'merchant' | 'system';
    kind?: 'image' | 'video';
    limit?: number;
  }>(),
  { defaultLibrary: 'merchant', kind: 'image', limit: 1 },
);
const emit = defineEmits<{ select: [PickerItem[]] }>();

const categories = ref<AttachmentCategory[]>([]);
const categoryID = ref(0);
/** all=普通库；system=侧栏「系统素材」模式（仅行级 is_system=1） */
const libraryMode = ref<'all' | 'system'>('all');
const files = ref<PickerItem[]>([]);
const selected = ref<number[]>([]);
const loading = ref(false);
const page = ref(1);
const total = ref(0);
const pageSize = 24;
/** 图片上传上限 10MB（与后端 upload 包一致） */
const IMAGE_MAX_BYTES = 10 * 1024 * 1024;

/** 与后端 SystemCategories 图片 enname 对齐 */
const IMAGE_SYSTEM_ENNAMES = new Set([
  'store_cover',
  'pay_icon',
  'logistics_icon',
  'service_icon',
  'product_image',
  'background_image',
  'list_icon',
  'other_image',
]);
/** 与后端 SystemCategories 视频 enname 对齐 */
const VIDEO_SYSTEM_ENNAMES = new Set([
  'store_video',
  'product_video',
  'other_video',
]);
/** 「系统素材」模式下上传时默认落入的系统分类 */
const DEFAULT_SYSTEM_UPLOAD_ENNAME = computed(() =>
  props.kind === 'video' ? 'other_video' : 'other_image',
);

function systemEnnames() {
  return props.kind === 'video' ? VIDEO_SYSTEM_ENNAMES : IMAGE_SYSTEM_ENNAMES;
}

function isSystemCategory(row: AttachmentCategory) {
  return (
    Number(row.is_system) === 1 ||
    systemEnnames().has(row.attachment_category_enname ?? '')
  );
}

function resolveUploadTarget() {
  if (libraryMode.value === 'system') {
    const cateID =
      categoryID.value > 0
        ? categoryID.value
        : (categories.value.find(
            (row) =>
              row.attachment_category_enname ===
              DEFAULT_SYSTEM_UPLOAD_ENNAME.value,
          )?.attachment_category_id ?? 0);
    return { cateID, isSystem: true as const };
  }
  // 全部素材：不强制落入系统分类，避免运营图被误标
  if (categoryID.value > 0) {
    return { cateID: categoryID.value, isSystem: false as const };
  }
  return { cateID: 0, isSystem: false as const };
}

const systemCategories = computed(() =>
  categories.value.filter((row) => isSystemCategory(row)),
);
const customCategories = computed(() =>
  categories.value.filter((row) => !isSystemCategory(row)),
);

const mediaLabel = computed(() => (props.kind === 'image' ? '图片' : '视频'));

async function loadCategories() {
  const result = await listAttachmentCategoriesApi({ type: props.kind });
  categories.value = result.list ?? [];
}

async function loadFiles() {
  loading.value = true;
  try {
    const params: Parameters<typeof listAttachmentsApi>[0] = {
      limit: pageSize,
      page: page.value,
      type: props.kind,
    };
    if (libraryMode.value === 'system') {
      // 行级系统预置素材；可叠加具体系统分类
      params.is_system = 1;
      if (categoryID.value > 0) {
        params.category_id = categoryID.value;
      }
    } else if (categoryID.value > 0) {
      params.category_id = categoryID.value;
    }
    const result = await listAttachmentsApi(params);
    files.value = (result.list ?? []).map((item) => ({
      ...item,
      file_id: item.attachment_id,
      file_path: item.attachment_src,
    }));
    total.value = result.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function selectCategory(id: number) {
  if (id === SYSTEM_ROOT_ID) {
    libraryMode.value = 'system';
    categoryID.value = SYSTEM_ROOT_ID;
  } else if (id === 0) {
    libraryMode.value = 'all';
    categoryID.value = 0;
  } else {
    const row = categories.value.find((item) => item.attachment_category_id === id);
    if (!row || !isSystemCategory(row)) {
      libraryMode.value = 'all';
    }
    // 系统子分类：保留当前 libraryMode（从「系统素材」点进来仍只看预置）
    categoryID.value = id;
  }
  page.value = 1;
  selected.value = [];
  void loadFiles();
}

function toggle(row: PickerItem) {
  const index = selected.value.indexOf(row.attachment_id);
  if (index >= 0) {
    selected.value.splice(index, 1);
    return;
  }
  if (selected.value.length >= props.limit) {
    ElMessage.warning(`最多选择 ${props.limit} 张${mediaLabel.value}`);
    return;
  }
  selected.value.push(row.attachment_id);
}

async function upload({ file }: { file: File }) {
  const ok =
    props.kind === 'image'
      ? file.type.startsWith('image/')
      : file.type.startsWith('video/');
  if (!ok) {
    ElMessage.warning(`请选择${mediaLabel.value}文件`);
    return;
  }
  if (props.kind === 'image' && file.size > IMAGE_MAX_BYTES) {
    ElMessage.error('图片不能超过 10MB');
    return;
  }
  const target = resolveUploadTarget();
  if (target.isSystem && !target.cateID) {
    ElMessage.warning('请先选择系统分类再上传系统素材');
    return;
  }
  await uploadAttachmentApi(file, target.cateID, { isSystem: target.isSystem });
  page.value = 1;
  await loadFiles();
  ElMessage.success(`${mediaLabel.value}已上传`);
}

function confirm() {
  const rows = files.value.filter((item) =>
    selected.value.includes(item.attachment_id),
  );
  if (!rows.length) {
    ElMessage.warning(`请选择${mediaLabel.value}`);
    return;
  }
  emit('select', rows);
  open.value = false;
}

async function initialize() {
  selected.value = [];
  // 按调用方指定的初始素材库打开；系统素材使用根入口筛选全部系统预置。
  libraryMode.value = props.defaultLibrary === 'system' ? 'system' : 'all';
  categoryID.value = props.defaultLibrary === 'system' ? SYSTEM_ROOT_ID : 0;
  page.value = 1;
  await loadCategories();
  await loadFiles();
}

watch(open, (visible) => {
  if (visible) void initialize();
});
onMounted(() => {
  if (open.value) void initialize();
});
</script>

<template>
  <ElDialog
    v-model="open"
    :title="`选择${mediaLabel}素材`"
    width="960px"
    destroy-on-close
    append-to-body
    class="image-picker-dialog"
  >
    <div class="picker">
      <aside class="picker__sidebar">
        <button
          type="button"
          class="picker__cat"
          :class="{ active: categoryID === 0 }"
          @click="selectCategory(0)"
        >
          全部素材
        </button>

        <div class="picker__scroll">
          <template v-if="customCategories.length">
            <div class="picker__section">自定义分类</div>
            <button
              v-for="row in customCategories"
              :key="row.attachment_category_id"
              type="button"
              class="picker__cat"
              :class="{ active: categoryID === row.attachment_category_id }"
              @click="selectCategory(row.attachment_category_id)"
            >
              {{ row.attachment_category_name }}
            </button>
          </template>
          <button
            v-for="row in systemCategories"
            :key="row.attachment_category_id"
            type="button"
            class="picker__cat"
            :class="{ active: categoryID === row.attachment_category_id }"
            @click="selectCategory(row.attachment_category_id)"
          >
            {{ row.attachment_category_name }}
          </button>
        </div>

        <div v-if="systemCategories.length" class="picker__footer">
          <button
            type="button"
            class="picker__cat"
            :class="{ active: libraryMode === 'system' }"
            @click="selectCategory(SYSTEM_ROOT_ID)"
          >
            系统素材
          </button>
        </div>
      </aside>

      <section class="picker__main">
        <div class="picker__toolbar">
          <span class="picker__hint">从平台素材库选择{{ mediaLabel }}</span>
          <ElUpload
            :accept="
              props.kind === 'image'
                ? 'image/jpeg,image/png,image/webp,image/gif'
                : 'video/mp4,video/quicktime,video/webm'
            "
            :http-request="upload"
            :show-file-list="false"
          >
            <ElButton type="primary">上传{{ mediaLabel }}</ElButton>
          </ElUpload>
        </div>

        <div
          v-loading="loading"
          class="picker__grid"
          :class="{ 'picker__grid--empty': !loading && files.length === 0 }"
        >
          <ElEmpty
            v-if="!loading && files.length === 0"
            :description="`暂无${mediaLabel}素材`"
          />
          <button
            v-for="row in files"
            :key="row.attachment_id"
            type="button"
            class="picker__item"
            :class="{ selected: selected.includes(row.attachment_id) }"
            @click="toggle(row)"
          >
            <img
              v-if="props.kind === 'image'"
              :src="row.attachment_src"
              :alt="row.attachment_name"
            />
            <video v-else :src="row.attachment_src" preload="metadata" />
          </button>
        </div>

        <div class="picker__pager">
          <ElPagination
            v-model:current-page="page"
            :page-size="pageSize"
            :total="total"
            background
            layout="total, prev, pager, next"
            @current-change="loadFiles"
          />
        </div>
      </section>
    </div>

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton type="primary" @click="confirm">确定</ElButton>
    </template>
  </ElDialog>
</template>

<style scoped lang="scss">
.picker {
  display: flex;
  min-height: 480px;
  overflow: hidden;
  border: 1px solid hsl(var(--border));
  border-radius: 8px;
  background: hsl(var(--background));
}

.picker__sidebar {
  display: flex;
  flex: 0 0 180px;
  flex-direction: column;
  gap: 2px;
  width: 180px;
  padding: 12px 8px;
  border-right: 1px solid hsl(var(--border));
  overflow: hidden;
}

.picker__scroll {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 2px;
  min-height: 0;
  overflow: auto;
}

.picker__footer {
  display: flex;
  flex-shrink: 0;
  flex-direction: column;
  gap: 2px;
  margin-top: auto;
  padding-top: 8px;
  border-top: 1px solid hsl(var(--border));
}

.picker__section {
  margin: 10px 8px 4px;
  color: hsl(var(--muted-foreground));
  font-size: 12px;
  font-weight: 600;
}

.picker__cat {
  display: block;
  width: 100%;
  overflow: hidden;
  padding: 8px 10px;
  border: 0;
  border-radius: 4px;
  color: hsl(var(--foreground));
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
  background: transparent;
  cursor: pointer;
}

.picker__cat:hover,
.picker__cat.active {
  color: hsl(var(--primary));
  background: hsl(var(--accent));
}

.picker__main {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 12px;
  min-width: 0;
  padding: 14px;
}

.picker__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid hsl(var(--border));
}

.picker__hint {
  color: hsl(var(--muted-foreground));
  font-size: 13px;
}

.picker__grid {
  display: grid;
  flex: 1;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 12px;
  align-content: start;
  min-height: 320px;
  min-width: 0;
}

.picker__grid--empty {
  display: flex;
  align-items: center;
  justify-content: center;
}

.picker__item {
  overflow: hidden;
  padding: 0;
  border: 2px solid transparent;
  border-radius: 6px;
  text-align: left;
  background: hsl(var(--muted));
  cursor: pointer;
}

.picker__item.selected {
  border-color: hsl(var(--primary));
}

.picker__item img,
.picker__item video {
  display: block;
  width: 100%;
  height: 96px;
  object-fit: contain;
  object-position: center;
  background: hsl(var(--muted));
}

.picker__pager {
  display: flex;
  justify-content: flex-end;
  padding-top: 4px;
}

@media (max-width: 768px) {
  .picker {
    display: block;
  }

  .picker__sidebar {
    width: auto;
    flex-basis: auto;
    border-right: 0;
    border-bottom: 1px solid hsl(var(--border));
  }

  .picker__grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}
</style>
