<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { PlatformProduct } from '#/api/core/platform-catalog';

import { nextTick, onActivated, onMounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElRate,
} from 'element-plus';

import ImageField from '#/components/shop/image-field.vue';
import ImagesField from '#/components/shop/images-field.vue';
import ProductPickerDialog from '#/components/shop/product-picker-dialog.vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createVirtualProductComment,
  deleteVirtualProductComment,
  fetchProductComments,
  moderateProductComment,
  updateVirtualProductComment,
  type ProductCommentRow,
  type ProductCommentStatus,
} from '#/api/core/ecrm';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const route = useRoute();
const saving = ref(false);
const editing = ref<ProductCommentRow>();
const productPicker = ref(false);
const selectedImageUrls = ref<string[]>([]);
const selectedProduct = ref<{
  cover: string;
  product_id: number;
  title: string;
} | null>(null);
const scoreSortOrder = ref<'asc' | 'desc' | undefined>();

const form = reactive({
  product_id: undefined as number | undefined,
  score: 5,
  content: '',
  virtual_author_name: '',
  virtual_author_avatar: '',
  sort: 0,
  attachment_ids: [] as number[],
});

const key = (action: string, id = 'new') =>
  `comment-${action}-${id}-${crypto.randomUUID()}`;

function mediaUrls(raw: string) {
  try {
    const value = JSON.parse(raw);
    return Array.isArray(value)
      ? value.filter((url) => typeof url === 'string').slice(0, 9)
      : [];
  } catch {
    return [];
  }
}

function formatTime(value?: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function displayUserName(row: ProductCommentRow) {
  const name = String(row.user_name || '').trim();
  if (name) return name;
  if (row.source === 'virtual') {
    return String(row.virtual_author_name || '').trim() || '—';
  }
  return '—';
}

function coverUrl(row: ProductCommentRow) {
  return resolveCosMediaUrl(String(row.product_cover || '').trim());
}

function selectedProductCover() {
  return resolveCosMediaUrl(String(selectedProduct.value?.cover || '').trim());
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    user_name: String(formValues?.user_name ?? '').trim() || undefined,
    status: (String(formValues?.status ?? '').trim() ||
      undefined) as ProductCommentStatus | undefined,
    date_from: range[0],
    date_to: range[1],
    sort_field: scoreSortOrder.value ? ('score' as const) : undefined,
    sort_order: scoreSortOrder.value,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'DatePicker',
    componentProps: {
      type: 'daterange',
      valueFormat: 'YYYY-MM-DD',
      startPlaceholder: '开始时间',
      endPlaceholder: '结束时间',
      class: 'w-full',
    },
    fieldName: 'date_range',
    label: '时间选择',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入用户名称',
    },
    fieldName: 'user_name',
    label: '用户名称',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入商品ID或者商品信息',
    },
    fieldName: 'keyword',
    label: '关键字',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待审核', value: 'pending' },
        { label: '已展示', value: 'published' },
        { label: '已隐藏', value: 'hidden' },
      ],
      placeholder: '请选择',
    },
    fieldName: 'status',
    label: '评价状态',
  },
]);

const gridOptions: VxeGridProps<ProductCommentRow> = {
  columns: [
    { align: 'left', field: 'id', title: 'ID', width: 82 },
    {
      align: 'left',
      field: 'product_cover',
      slots: { default: 'cover' },
      title: '商品图',
      width: 88,
    },
    {
      align: 'left',
      field: 'product_title',
      minWidth: 200,
      showOverflow: false,
      slots: { default: 'productTitle' },
      title: '商品名称',
    },
    {
      align: 'left',
      field: 'user_name',
      minWidth: 120,
      slots: { default: 'userName' },
      title: '用户名称',
    },
    {
      align: 'left',
      field: 'score',
      sortable: true,
      title: '产品评分',
      width: 110,
    },
    { align: 'left', field: 'sort', title: '排序', width: 80 },
    {
      align: 'left',
      field: 'content',
      minWidth: 180,
      showOverflow: false,
      slots: { default: 'content' },
      title: '评价内容',
    },
    {
      align: 'left',
      field: 'reply_content',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 140,
      showOverflow: false,
      title: '回复内容',
    },
    {
      align: 'left',
      field: 'created_at',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 170,
      title: '评价时间',
    },
    platformListActionColumn({ width: 160 }),
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page, sorts }, formValues) => {
        const scoreSort = Array.isArray(sorts)
          ? sorts.find((item) => item.field === 'score')
          : undefined;
        const order = String(scoreSort?.order || '').toLowerCase();
        scoreSortOrder.value =
          order === 'asc' || order === 'desc' ? order : undefined;
        const data = await fetchProductComments(
          buildListParams(page, formValues),
        );
        return { items: data.list || [], total: data.total || 0 };
      },
    },
    sort: true,
  },
  rowConfig: { isHover: true, keyField: 'id' },
  sortConfig: { remote: true, trigger: 'cell', multiple: false },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

function routeKeyword() {
  return (
    String(route.query.keyword ?? '').trim() ||
    String(route.query.product_id ?? '').trim()
  );
}

async function applyKeywordFromRoute() {
  const keyword = routeKeyword();
  if (!keyword) return;
  await nextTick();
  if (!gridApi.formApi?.setValues) return;
  await gridApi.formApi.setValues({ keyword });
  const values = await gridApi.formApi.getValues();
  gridApi.formApi.setLatestSubmissionValues?.(values);
  await gridApi.reload(values);
}

onMounted(() => {
  void applyKeywordFromRoute();
});

onActivated(() => {
  void applyKeywordFromRoute();
});

watch(
  () => [route.query.keyword, route.query.product_id] as const,
  () => {
    void applyKeywordFromRoute();
  },
);

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '保存',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => save(),
});

function open(row?: ProductCommentRow) {
  editing.value = row;
  if (row) {
    selectedProduct.value = {
      product_id: row.product_id,
      title: String(row.product_title || '').trim() || `商品 #${row.product_id}`,
      cover: String(row.product_cover || '').trim(),
    };
    selectedImageUrls.value = mediaUrls(String(row.media || ''));
    Object.assign(form, {
      product_id: row.product_id,
      score: row.score,
      content: row.content,
      virtual_author_name: row.virtual_author_name,
      virtual_author_avatar: String(row.virtual_author_avatar || '').trim(),
      sort: row.sort,
      attachment_ids: [],
    });
  } else {
    selectedProduct.value = null;
    selectedImageUrls.value = [];
    Object.assign(form, {
      product_id: undefined,
      score: 5,
      content: '',
      virtual_author_name: '',
      virtual_author_avatar: '',
      sort: 0,
      attachment_ids: [],
    });
  }
  formDrawerApi
    .setState({ title: row ? '编辑自评' : '新增自评' })
    .open();
}

function onProductSelect(product: PlatformProduct) {
  const productId = Number(product.product_id || 0);
  if (!productId) return;
  form.product_id = productId;
  selectedProduct.value = {
    product_id: productId,
    title:
      String(product.store_name || product.title || '').trim() ||
      `商品 #${productId}`,
    cover: String(product.image || '').trim(),
  };
}

async function save() {
  if (!form.product_id || !form.content.trim() || !form.virtual_author_name.trim()) {
    ElMessage.warning('请选择商品，并填写显示昵称和评论内容');
    return;
  }
  formDrawerApi.lock();
  saving.value = true;
  try {
    const value = {
      score: form.score,
      content: form.content.trim(),
      virtual_author_name: form.virtual_author_name.trim(),
      virtual_author_avatar: form.virtual_author_avatar.trim(),
      sort: form.sort,
      attachment_ids: form.attachment_ids.length
        ? form.attachment_ids
        : undefined,
    };
    if (editing.value) {
      await updateVirtualProductComment(editing.value.id, {
        ...value,
        idempotency_key: key('update', String(editing.value.id)),
      });
    } else {
      await createVirtualProductComment({
        product_id: form.product_id,
        ...value,
        idempotency_key: key('create'),
      });
    }
    formDrawerApi.close();
    ElMessage.success(editing.value ? '自评已更新' : '自评已新增并展示');
    gridApi.reload();
  } finally {
    saving.value = false;
    formDrawerApi.unlock();
  }
}

async function moderate(row: ProductCommentRow, action: 'hide' | 'publish') {
  const label = action === 'publish' ? '展示' : '隐藏';
  try {
    const { value } = await ElMessageBox.prompt(
      `确认${label}该评论？可填写审核说明，用户原始评论不可篡改。`,
      `${label}评论`,
      {
        inputPlaceholder: '审核说明（可选，最多 500 字）',
        inputValidator: (v) => [...v].length <= 500 || '审核说明不能超过 500 字',
        confirmButtonText: `确认${label}`,
        cancelButtonText: '取消',
      },
    );
    await moderateProductComment(row.id, {
      action,
      note: value.trim(),
      idempotency_key: key(action, String(row.id)),
    });
    ElMessage.success(`评论已${label}`);
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

async function remove(row: ProductCommentRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      '删除仅对虚拟评论生效，将保留不可恢复的审核审计记录。',
      '删除自评',
      {
        inputPlaceholder: '删除说明（可选）',
        inputValidator: (v) => [...v].length <= 500 || '删除说明不能超过 500 字',
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning',
      },
    );
    await deleteVirtualProductComment(row.id, {
      note: value.trim(),
      idempotency_key: key('delete', String(row.id)),
    });
    ElMessage.success('自评已删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton type="primary" @click="open()">新增自评</ElButton>
      </template>
      <template #cover="{ row }">
        <ElImage
          v-if="coverUrl(row)"
          :src="coverUrl(row)"
          fit="cover"
          class="comment-product-cover"
          :preview-src-list="[coverUrl(row)]"
          preview-teleported
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>
      <template #productTitle="{ row }">
        <div class="comment-product-title">
          {{ row.product_title || `商品 #${row.product_id}` }}
        </div>
      </template>
      <template #userName="{ row }">
        {{ displayUserName(row) }}
      </template>
      <template #content="{ row }">
        <div class="comment-content-cell">
          <div>{{ row.content || '—' }}</div>
          <div v-if="mediaUrls(row.media).length" class="comment-content-media">
            <ElImage
              v-for="url in mediaUrls(row.media)"
              :key="url"
              :src="resolveCosMediaUrl(url)"
              fit="cover"
              class="comment-content-media__img"
              :preview-src-list="mediaUrls(row.media).map((u) => resolveCosMediaUrl(u))"
              preview-teleported
            >
              <template #error>
                <span class="text-xs text-gray-400">加载失败</span>
              </template>
            </ElImage>
          </div>
        </div>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="row.status !== 'published'"
          link
          type="primary"
          @click="moderate(row, 'publish')"
        >
          展示
        </ElButton>
        <ElButton
          v-if="row.status !== 'hidden'"
          link
          type="primary"
          @click="moderate(row, 'hide')"
        >
          隐藏
        </ElButton>
        <ElButton
          v-if="row.source === 'virtual'"
          link
          type="primary"
          @click="open(row)"
        >
          编辑
        </ElButton>
        <ElButton
          v-if="row.source === 'virtual'"
          link
          type="primary"
          @click="remove(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>

    <FormDrawer>
      <ElAlert
        class="mb-4"
        type="warning"
        :closable="false"
        title="自评必须显著标识为虚拟来源；不得冒充真实用户或用于虚假交易宣传。"
      />
      <ElForm label-width="90px">
        <ElFormItem label="商品" required>
          <div class="comment-product-field">
            <div
              v-if="selectedProduct"
              class="comment-product-card"
            >
              <ElImage
                v-if="selectedProductCover()"
                :src="selectedProductCover()"
                fit="cover"
                class="comment-product-card__cover"
              >
                <template #error>
                  <span class="comment-product-card__cover-fallback">—</span>
                </template>
              </ElImage>
              <div
                v-else
                class="comment-product-card__cover comment-product-card__cover--empty"
              >
                —
              </div>
              <div class="comment-product-card__meta">
                <div class="comment-product-card__title">
                  {{ selectedProduct.title }}
                </div>
                <div class="comment-product-card__id">
                  ID：{{ selectedProduct.product_id }}
                </div>
              </div>
              <ElButton
                v-if="!editing"
                type="primary"
                link
                class="comment-product-card__action"
                @click="productPicker = true"
              >
                重新选择
              </ElButton>
            </div>
            <ElButton
              v-else
              type="primary"
              plain
              @click="productPicker = true"
            >
              选择商品
            </ElButton>
          </div>
        </ElFormItem>
        <ElFormItem label="显示昵称" required>
          <ElInput
            v-model="form.virtual_author_name"
            maxlength="64"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="显示头像">
          <ImageField
            v-model="form.virtual_author_avatar"
            default-library="system"
            :preview-size="64"
          />
        </ElFormItem>
        <ElFormItem label="评分" required>
          <ElRate v-model="form.score" />
        </ElFormItem>
        <ElFormItem label="评论内容" required>
          <ElInput
            v-model="form.content"
            type="textarea"
            :rows="4"
            maxlength="2000"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="评论图片">
          <ImagesField
            v-model="selectedImageUrls"
            v-model:attachment-ids="form.attachment_ids"
            default-library="system"
            :limit="9"
            :preview-size="64"
          />
        </ElFormItem>
        <ElFormItem label="排序">
          <ElInputNumber
            v-model="form.sort"
            :min="0"
            :max="999999"
            controls-position="right"
          />
        </ElFormItem>
      </ElForm>
    </FormDrawer>

    <ProductPickerDialog
      v-model:open="productPicker"
      @select="onProductSelect"
    />
  </Page>
</template>

<style scoped>
.comment-product-cover {
  display: block;
  width: 48px;
  height: 48px;
  border-radius: 4px;
}

.comment-product-title {
  line-height: 1.5;
  white-space: normal;
  word-break: break-word;
}

.comment-content-cell {
  line-height: 1.5;
  white-space: normal;
  word-break: break-word;
}

.comment-content-media {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 4px;
}

.comment-content-media__img {
  width: 32px;
  height: 32px;
  border-radius: 2px;
}

.comment-product-field {
  display: flex;
  align-items: flex-start;
}

.comment-product-card {
  display: inline-flex;
  gap: 10px;
  align-items: center;
  max-width: 420px;
  padding: 8px 12px 8px 8px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  background: var(--el-fill-color-lighter);
}

.comment-product-card__cover {
  width: 48px;
  height: 48px;
  border-radius: 4px;
  flex-shrink: 0;
  background: var(--el-fill-color);
}

.comment-product-card__cover--empty,
.comment-product-card__cover-fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-placeholder);
  font-size: 12px;
}

.comment-product-card__meta {
  min-width: 0;
  flex: 1;
}

.comment-product-card__title {
  line-height: 1.4;
  word-break: break-word;
}

.comment-product-card__id {
  margin-top: 2px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.comment-product-card__action {
  flex-shrink: 0;
  margin-left: 4px;
}
</style>
