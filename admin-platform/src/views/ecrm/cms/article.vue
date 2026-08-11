<script setup lang="ts">
import type { ImageUploadOptions } from '@vben/plugins/tiptap';
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import { VbenTiptap } from '@vben/plugins/tiptap';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElOption,
  ElRadio,
  ElRadioGroup,
  ElSelect,
  ElSwitch,
  ElTabPane,
  ElTabs,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { uploadAttachmentApi } from '#/api/core/attachment';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  createArticle,
  deleteArticle,
  fetchArticles,
  getArticle,
  updateArticle,
  type ArticleRow,
} from '#/api/core/ecrm';
import {
  getArticleCategoryListApi,
  type ArticleCategoryOption,
} from '#/api/core/plus-article';
import ImageField from '#/components/shop/image-field.vue';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const canRead = ref(false);
const canManage = ref(false);
const categories = ref<ArticleCategoryOption[]>([]);
const editingId = ref(0);
const formTab = ref('info');

const form = reactive({
  title: '',
  author: '',
  cid: undefined as number | undefined,
  synopsis: '',
  image: '',
  status: 1,
  content: '',
});

const categoryMap = computed(() => {
  const map = new Map<number, string>();
  for (const item of categories.value) {
    map.set(Number(item.cid), item.title);
  }
  return map;
});

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入文章标题',
    },
    fieldName: 'title',
    label: '文章标题',
  },
]);

const gridOptions: VxeGridProps<ArticleRow> = {
  columns: [
    { field: 'article_id', title: 'ID', width: 80 },
    {
      field: 'image',
      slots: { default: 'image' },
      title: '文章图片',
      width: 110,
    },
    {
      field: 'title',
      minWidth: 220,
      showOverflow: 'tooltip',
      slots: { default: 'title' },
      title: '文章标题',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue) || '—',
      minWidth: 170,
      showOverflow: false,
      title: '时间',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '是否显示',
      width: 110,
    },
    platformListActionColumn({ width: 140 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const title = String(formValues?.title ?? '').trim();
        const data = await fetchArticles({
          page: page.currentPage,
          limit: page.pageSize,
          title: title || undefined,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'article_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

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

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function mediaUrl(url?: string) {
  return resolveCosMediaUrl(String(url || '').trim());
}

function categoryTitleOf(cid: number) {
  return categoryMap.value.get(Number(cid)) || '';
}

function displayTitle(row: ArticleRow) {
  const cat = categoryTitleOf(row.cid);
  return cat ? `[${cat}]${row.title}` : row.title;
}

function resetForm() {
  editingId.value = 0;
  formTab.value = 'info';
  Object.assign(form, {
    title: '',
    author: '',
    cid: categories.value[0]?.cid,
    synopsis: '',
    image: '',
    status: 1,
    content: '',
  });
}

function openCreate() {
  resetForm();
  formDrawerApi.setState({ title: '新增文章' }).open();
}

async function openEdit(row: ArticleRow) {
  editingId.value = row.article_id;
  formTab.value = 'info';
  formDrawerApi.setState({ title: '编辑文章' }).open();
  formDrawerApi.lock();
  try {
    const detail = await getArticle(row.article_id);
    Object.assign(form, {
      title: detail.title || '',
      author: detail.author || '',
      cid: detail.cid || undefined,
      synopsis: detail.synopsis || '',
      image: detail.image || '',
      status: detail.status === 1 ? 1 : 0,
      content: detail.content || '',
    });
  } finally {
    formDrawerApi.unlock();
  }
}

function isEmptyRichText(html: string) {
  const text = html
    .replace(/<[^>]+>/g, '')
    .replace(/&nbsp;/g, ' ')
    .trim();
  return !text;
}

async function save() {
  const title = form.title.trim();
  const author = form.author.trim();
  const image = form.image.trim();
  const content = form.content.trim();
  if (!title) {
    ElMessage.warning('请填写标题');
    formTab.value = 'info';
    return;
  }
  if (!author) {
    ElMessage.warning('请填写作者');
    formTab.value = 'info';
    return;
  }
  if (!form.cid) {
    ElMessage.warning('请选择文章分类');
    formTab.value = 'info';
    return;
  }
  if (!image) {
    ElMessage.warning('请上传图文封面');
    formTab.value = 'info';
    return;
  }
  if (!content || isEmptyRichText(content)) {
    ElMessage.warning('请填写文章内容');
    formTab.value = 'content';
    return;
  }

  const payload = {
    title,
    author,
    cid: Number(form.cid),
    synopsis: form.synopsis.trim(),
    image,
    content: form.content,
    status: form.status === 1 ? 1 : 0,
  };

  formDrawerApi.lock();
  try {
    if (editingId.value) {
      await updateArticle(editingId.value, payload);
    } else {
      await createArticle(payload);
    }
    formDrawerApi.close();
    ElMessage.success('已保存');
    gridApi.reload();
  } finally {
    formDrawerApi.unlock();
  }
}

async function changeStatus(row: ArticleRow, enabled: boolean) {
  const before = row.status;
  row.status = enabled ? 1 : 0;
  try {
    await updateArticle(row.article_id, { status: enabled ? 1 : 0 });
  } catch {
    row.status = before;
  }
}

async function removeRow(row: ArticleRow) {
  try {
    await confirm({
      content: `确定删除文章「${row.title}」吗？`,
      icon: 'warning',
      title: '提示',
    });
    await deleteArticle(row.article_id);
    ElMessage.success('已删除');
    gridApi.reload();
  } catch {
    /* cancel */
  }
}

async function loadCategories() {
  const data = await getArticleCategoryListApi();
  categories.value = (data.list || []).filter((item) => item.status === 1);
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
    loadCategories(),
  ]);
  const roleOK = profile.roles.some(
    (role) => role === 'platform' || role === 'operations',
  );
  canRead.value =
    roleOK &&
    (codes.includes('content.article.read') ||
      codes.includes('content.article.manage'));
  canManage.value = roleOK && codes.includes('content.article.manage');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canManage"
          :icon="Plus"
          type="primary"
          @click="openCreate"
        >
          新增文章
        </ElButton>
      </template>

      <template #image="{ row }">
        <ElImage
          v-if="mediaUrl(row.image)"
          :src="mediaUrl(row.image)"
          fit="cover"
          class="article-thumb"
          :preview-src-list="[mediaUrl(row.image)]"
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>

      <template #title="{ row }">
        <span class="article-title-cell">{{ displayTitle(row) }}</span>
      </template>

      <template #status="{ row }">
        <ElSwitch
          :model-value="row.status === 1"
          :disabled="!canManage"
          inline-prompt
          active-text="显示"
          inactive-text="隐藏"
          @change="
            (enabled: string | number | boolean) =>
              changeStatus(row, Boolean(enabled))
          "
        />
      </template>

      <template #action="{ row }">
        <ElButton
          v-if="canManage"
          link
          type="primary"
          @click="openEdit(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-if="canManage"
          link
          type="primary"
          @click="removeRow(row)"
        >
          删除
        </ElButton>
        <span v-if="!canManage">—</span>
      </template>
    </Grid>

    <FormDrawer>
      <ElTabs v-model="formTab" class="article-tabs">
        <ElTabPane label="文章信息" name="info">
          <ElForm label-width="100px" class="article-form">
            <ElFormItem label="标题" required>
              <ElInput
                v-model="form.title"
                maxlength="100"
                show-word-limit
                placeholder="请输入文章标题"
              />
            </ElFormItem>
            <ElFormItem label="作者" required>
              <ElInput
                v-model="form.author"
                maxlength="32"
                show-word-limit
                placeholder="请输入作者"
              />
            </ElFormItem>
            <ElFormItem label="文章分类" required>
              <ElSelect
                v-model="form.cid"
                class="w-full"
                filterable
                clearable
                placeholder="请选择"
              >
                <ElOption
                  v-for="item in categories"
                  :key="item.cid"
                  :label="item.title"
                  :value="item.cid"
                />
              </ElSelect>
            </ElFormItem>
            <ElFormItem label="文章简介">
              <ElInput
                v-model="form.synopsis"
                type="textarea"
                :rows="3"
                maxlength="200"
                show-word-limit
                placeholder="请输入文章简介"
              />
            </ElFormItem>
            <ElFormItem label="图文封面" required>
              <ImageField v-model="form.image" :preview-size="88" />
            </ElFormItem>
            <ElFormItem label="是否显示" required>
              <ElRadioGroup v-model="form.status">
                <ElRadio :label="1">显示</ElRadio>
                <ElRadio :label="0">不显示</ElRadio>
              </ElRadioGroup>
            </ElFormItem>
          </ElForm>
        </ElTabPane>
        <ElTabPane label="文章内容" name="content">
          <ElForm label-width="100px" class="article-form">
            <ElFormItem label="文章内容" required>
              <VbenTiptap
                v-model="form.content"
                :editable="true"
                :image-upload="imageUpload"
                :max-height="480"
                :min-height="320"
                :previewable="false"
                placeholder="请输入文章内容…"
              />
            </ElFormItem>
          </ElForm>
        </ElTabPane>
      </ElTabs>
    </FormDrawer>
  </Page>
</template>

<style scoped>
.article-thumb {
  display: block;
  width: 56px;
  height: 56px;
  border-radius: 6px;
}

.article-title-cell {
  line-height: 1.5;
  word-break: break-all;
}

.article-tabs {
  margin-top: 4px;
}

.article-form {
  padding-right: 8px;
}

.article-form :deep(.el-form-item) {
  margin-bottom: 18px;
}
</style>
