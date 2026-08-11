<script setup lang="ts">
import type { ImageUploadOptions } from '@vben/plugins/tiptap';

import { computed, reactive, ref, watch } from 'vue';

import { confirm } from '@vben/common-ui';
import { VbenTiptap } from '@vben/plugins/tiptap';
import { Picture as IconPicture } from '@element-plus/icons-vue';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElIcon,
  ElInput,
  ElMessage,
} from 'element-plus';

import { uploadAttachmentApi } from '#/api/core/attachment';
import {
  createWechatNewsApi,
  getWechatNewsApi,
  updateWechatNewsApi,
  type WechatNewsItem,
} from '#/api/core/platform-wechat-news';
import ImageField from '#/components/shop/image-field.vue';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const MAX_ITEMS = 8;

const props = withDefaults(
  defineProps<{
    canManage?: boolean;
    newsId?: number;
  }>(),
  {
    canManage: false,
    newsId: 0,
  },
);
const emit = defineEmits<{ saved: [] }>();

const loading = ref(false);
const saving = ref(false);
const activeIndex = ref(0);
const items = ref<WechatNewsItem[]>([emptyItem()]);

const canManage = computed(() => props.canManage);
const editingId = computed(() => Number(props.newsId) || 0);

const form = reactive({
  title: '',
  author: '',
  synopsis: '',
  image: '',
  content: '',
});

const imageUpload: ImageUploadOptions = {
  accept: 'image/jpeg,image/png,image/gif,image/webp',
  maxSize: 5 * 1024 * 1024,
  upload: async (file) => {
    const row = await uploadAttachmentApi(file);
    return row.attachment_src;
  },
  onUploadError: () => {
    ElMessage.error('图片上传失败');
  },
};

function emptyItem(): WechatNewsItem {
  return {
    title: '',
    author: '',
    synopsis: '',
    image: '',
    content: '',
  };
}

function syncFormFromActive() {
  const item = items.value[activeIndex.value] || emptyItem();
  form.title = item.title || '';
  form.author = item.author || '';
  form.synopsis = item.synopsis || '';
  form.image = item.image || '';
  form.content = item.content || '';
}

function applyFormToActive() {
  const idx = activeIndex.value;
  if (!items.value[idx]) return;
  items.value[idx] = {
    title: form.title,
    author: form.author,
    synopsis: form.synopsis,
    image: form.image,
    content: form.content,
  };
}

function selectItem(index: number) {
  if (index === activeIndex.value) return;
  applyFormToActive();
  activeIndex.value = index;
  syncFormFromActive();
}

function addItem() {
  if (!canManage.value) return;
  if (items.value.length >= MAX_ITEMS) {
    ElMessage.warning(`最多新增 ${MAX_ITEMS} 条图文`);
    return;
  }
  applyFormToActive();
  items.value.push(emptyItem());
  activeIndex.value = items.value.length - 1;
  syncFormFromActive();
}

async function removeItem(index: number) {
  if (!canManage.value) return;
  if (items.value.length <= 1) {
    ElMessage.warning('至少保留一条图文');
    return;
  }
  try {
    await confirm({
      title: '删除确认',
      content: '确定删除该条图文吗？',
    });
  } catch {
    return;
  }
  applyFormToActive();
  items.value.splice(index, 1);
  if (activeIndex.value >= items.value.length) {
    activeIndex.value = items.value.length - 1;
  }
  syncFormFromActive();
}

function validate(): boolean {
  applyFormToActive();
  for (let i = 0; i < items.value.length; i += 1) {
    const item = items.value[i]!;
    if (!item.title.trim()) {
      ElMessage.warning(`请填写第 ${i + 1} 条图文标题`);
      selectItem(i);
      return false;
    }
    if (!item.author.trim()) {
      ElMessage.warning(`请填写第 ${i + 1} 条图文作者`);
      selectItem(i);
      return false;
    }
    if (!item.synopsis.trim()) {
      ElMessage.warning(`请填写第 ${i + 1} 条图文摘要`);
      selectItem(i);
      return false;
    }
    if (!item.image.trim()) {
      ElMessage.warning(`请上传第 ${i + 1} 条图文封面`);
      selectItem(i);
      return false;
    }
    const content = item.content.trim();
    if (!content || content === '<p></p>') {
      ElMessage.warning(`请填写第 ${i + 1} 条图文正文`);
      selectItem(i);
      return false;
    }
  }
  return true;
}

async function loadDetail(id: number) {
  loading.value = true;
  try {
    const data = await getWechatNewsApi(id);
    const list = data.article?.length ? data.article : [emptyItem()];
    items.value = list.map((item) => ({
      title: item.title || '',
      author: item.author || '',
      synopsis: item.synopsis || '',
      image: item.image || '',
      content: item.content || '',
    }));
    activeIndex.value = 0;
    syncFormFromActive();
  } finally {
    loading.value = false;
  }
}

function resetEditor() {
  items.value = [emptyItem()];
  activeIndex.value = 0;
  syncFormFromActive();
}

async function save(): Promise<boolean> {
  if (!canManage.value || !validate()) return false;
  saving.value = true;
  try {
    const payload = {
      status: 1,
      data: items.value.map((item) => ({
        title: item.title.trim(),
        author: item.author.trim(),
        synopsis: item.synopsis.trim(),
        image: item.image.trim(),
        content: item.content.trim(),
      })),
    };
    if (editingId.value) {
      await updateWechatNewsApi(editingId.value, payload);
      ElMessage.success('编辑成功');
    } else {
      await createWechatNewsApi(payload);
      ElMessage.success('新增成功');
    }
    emit('saved');
    return true;
  } finally {
    saving.value = false;
  }
}

watch(
  () => [form.title, form.author, form.synopsis, form.image, form.content],
  () => applyFormToActive(),
);

watch(
  () => editingId.value,
  (id) => {
    if (id > 0) {
      void loadDetail(id);
    } else {
      resetEditor();
    }
  },
  { immediate: true },
);

defineExpose({ submit: save });
</script>

<template>
  <div v-loading="loading || saving" class="wechat-news-save">
    <div class="wechat-news-save__body">
      <!-- 左侧预览 / 多图文切换 -->
      <div class="wechat-news-save__left">
        <div
          v-for="(item, index) in items"
          :key="index"
          class="news-preview"
          :class="{ 'is-active': activeIndex === index }"
          @click="selectItem(index)"
        >
            <template v-if="index === 0">
              <div class="news-preview__cover">
                <img
                  v-if="item.image"
                  :src="resolveCosMediaUrl(item.image)"
                  alt=""
                />
                <div v-else class="news-preview__cover-empty">
                  <ElIcon :size="36" color="#c0c4cc">
                    <IconPicture />
                  </ElIcon>
                </div>
                <div v-if="item.title" class="news-preview__cover-title">
                  {{ item.title }}
                </div>
              </div>
            </template>
            <template v-else>
              <div class="news-preview__row">
                <span class="news-preview__row-title">
                  {{ item.title || '标题' }}
                </span>
                <div class="news-preview__row-thumb">
                  <img
                    v-if="item.image"
                    :src="resolveCosMediaUrl(item.image)"
                    alt=""
                  />
                  <ElIcon v-else :size="18" color="#c0c4cc">
                    <IconPicture />
                  </ElIcon>
                </div>
              </div>
            </template>
            <ElButton
              v-if="canManage && items.length > 1"
              class="news-preview__del"
              type="danger"
              circle
              size="small"
              @click.stop="removeItem(index)"
            >
              ×
            </ElButton>
        </div>

        <div v-if="canManage" class="news-preview__add-wrap">
          <ElButton type="primary" size="small" @click="addItem">
            新增图文
          </ElButton>
        </div>
      </div>

      <!-- 右侧表单 -->
      <div class="wechat-news-save__right">
        <ElForm
          label-width="96px"
          :disabled="!canManage"
          @submit.prevent
        >
            <ElFormItem label="标题" required>
              <ElInput
                v-model="form.title"
                maxlength="64"
                show-word-limit
                placeholder="请输入文章标题"
              />
            </ElFormItem>
            <ElFormItem label="作者" required>
              <ElInput
                v-model="form.author"
                maxlength="32"
                show-word-limit
                placeholder="请输入作者名称"
              />
            </ElFormItem>
            <ElFormItem label="摘要" required>
              <ElInput
                v-model="form.synopsis"
                type="textarea"
                :rows="3"
                maxlength="128"
                show-word-limit
                placeholder="请输入摘要"
              />
            </ElFormItem>
            <ElFormItem label="图文封面" required>
              <ImageField v-model="form.image" :preview-size="88" />
            </ElFormItem>
            <ElFormItem label="正文" required>
              <VbenTiptap
                v-model="form.content"
                class="wechat-news-save__editor"
                :editable="canManage"
                :image-upload="imageUpload"
                :max-height="420"
                :min-height="280"
                :previewable="false"
                placeholder="请输入正文内容…"
              />
            </ElFormItem>
        </ElForm>
      </div>
    </div>
  </div>
</template>

<style scoped>
.wechat-news-save__body {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr);
  gap: 24px;
  align-items: start;
}

.wechat-news-save__left {
  padding: 16px;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.news-preview {
  position: relative;
  margin-bottom: 12px;
  overflow: hidden;
  background: #fff;
  border: 1px dashed #c0c4cc;
  border-radius: 4px;
  cursor: pointer;
}

.news-preview.is-active {
  border-color: var(--el-color-primary);
  border-style: solid;
}

.news-preview__cover {
  position: relative;
  height: 148px;
  background: #f5f7fa;
}

.news-preview__cover img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.news-preview__cover-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.news-preview__cover-title {
  position: absolute;
  right: 0;
  bottom: 0;
  left: 0;
  padding: 8px 10px;
  overflow: hidden;
  color: #fff;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
  background: linear-gradient(transparent, rgb(0 0 0 / 55%));
}

.news-preview__row {
  display: flex;
  gap: 10px;
  align-items: center;
  justify-content: space-between;
  min-height: 64px;
  padding: 10px;
}

.news-preview__row-title {
  flex: 1;
  overflow: hidden;
  color: #303133;
  font-size: 12px;
  line-height: 18px;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.news-preview__row-thumb {
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  overflow: hidden;
  background: #f5f7fa;
  border-radius: 4px;
}

.news-preview__row-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.news-preview__del {
  position: absolute;
  top: 6px;
  right: 6px;
  z-index: 2;
  width: 22px;
  height: 22px;
  padding: 0;
  font-size: 14px;
  line-height: 20px;
}

.news-preview__add-wrap {
  display: flex;
  justify-content: center;
  padding-top: 4px;
}

.wechat-news-save__right {
  padding: 20px 24px;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.wechat-news-save__editor {
  width: 100%;
}

@media (max-width: 960px) {
  .wechat-news-save__body {
    grid-template-columns: 1fr;
  }
}
</style>
