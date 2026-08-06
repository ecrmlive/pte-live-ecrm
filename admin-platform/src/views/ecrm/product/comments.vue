<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { reactive, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
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
  ElTag,
} from 'element-plus';
import { Plus } from '@element-plus/icons-vue';

import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';
import type { AttachmentItem } from '#/api/core/attachment';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  createVirtualProductComment,
  deleteVirtualProductComment,
  fetchProductComments,
  moderateProductComment,
  sortVirtualProductComment,
  updateVirtualProductComment,
  type ProductCommentRow,
  type ProductCommentStatus,
} from '#/api/core/ecrm';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const editing = ref<ProductCommentRow>();
const picker = ref(false);
const selectedImages = ref<AttachmentItem[]>([]);

const form = reactive({
  product_id: undefined as number | undefined,
  score: 5,
  content: '',
  virtual_author_name: '',
  sort: 0,
  attachment_ids: undefined as number[] | undefined,
});

const statusInfo = (status: ProductCommentStatus) =>
  ({
    pending: { label: '待审核', type: 'warning' as const },
    published: { label: '已展示', type: 'success' as const },
    hidden: { label: '已隐藏', type: 'info' as const },
  })[status];

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

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const productIdRaw = String(formValues?.product_id ?? '').trim();
  const replyStatus = String(formValues?.reply_status ?? '').trim();
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    product_id: productIdRaw ? Number(productIdRaw) : undefined,
    status: (String(formValues?.status ?? '').trim() ||
      undefined) as ProductCommentStatus | undefined,
    reply_status: (replyStatus === 'has_reply' || replyStatus === 'no_reply'
      ? replyStatus
      : undefined) as 'has_reply' | 'no_reply' | undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '商品名 / 评论内容 / 昵称',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '商品 ID' },
    fieldName: 'product_id',
    label: '商品 ID',
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
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '展示状态',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '已回复', value: 'has_reply' },
        { label: '未回复', value: 'no_reply' },
      ],
      placeholder: '全部',
    },
    fieldName: 'reply_status',
    label: '商家回复',
  },
]);

const gridOptions: VxeGridProps<ProductCommentRow> = {
  columns: [
    { field: 'id', title: 'ID', width: 82 },
    {
      field: 'product_title',
      minWidth: 190,
      showOverflow: false,
      slots: { default: 'product' },
      title: '商品',
    },
    {
      field: 'source',
      slots: { default: 'source' },
      title: '来源',
      width: 130,
    },
    {
      field: 'score',
      slots: { default: 'score' },
      title: '评分',
      width: 120,
    },
    {
      field: 'content',
      minWidth: 240,
      showOverflow: false,
      title: '评论内容',
    },
    {
      field: 'media',
      slots: { default: 'media' },
      title: '图片',
      width: 120,
    },
    {
      field: 'reply_content',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 150,
      showOverflow: false,
      title: '商家回复',
    },
    { field: 'sort', title: '排序', width: 80 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ width: 205 }),
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchProductComments(buildListParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [FormModal, formModalApi] = useVbenModal({
  onConfirm: async () => save(),
});

function open(row?: ProductCommentRow) {
  editing.value = row;
  selectedImages.value = [];
  Object.assign(
    form,
    row
      ? {
          product_id: row.product_id,
          score: row.score,
          content: row.content,
          virtual_author_name: row.virtual_author_name,
          sort: row.sort,
          attachment_ids: undefined,
        }
      : {
          product_id: undefined,
          score: 5,
          content: '',
          virtual_author_name: '',
          sort: 0,
          attachment_ids: [] as number[],
        },
  );
  formModalApi
    .setState({ title: row ? '编辑虚拟评论' : '新增虚拟评论' })
    .open();
}

async function save() {
  if (!form.product_id || !form.content.trim() || !form.virtual_author_name.trim()) {
    ElMessage.warning('请填写商品 ID、显示昵称和评论内容');
    return;
  }
  formModalApi.lock();
  saving.value = true;
  try {
    const value = {
      score: form.score,
      content: form.content.trim(),
      virtual_author_name: form.virtual_author_name.trim(),
      sort: form.sort,
      attachment_ids: form.attachment_ids,
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
    formModalApi.close();
    ElMessage.success(editing.value ? '虚拟评论已更新' : '虚拟评论已新增并展示');
    gridApi.reload();
  } finally {
    saving.value = false;
    formModalApi.unlock();
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

async function sort(row: ProductCommentRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      '数值越大越靠前，仅影响虚拟评论展示顺序。',
      '调整排序',
      {
        inputValue: String(row.sort),
        inputPattern: /^\d{1,6}$/,
        inputErrorMessage: '请输入 0 到 999999 的整数',
        confirmButtonText: '保存',
        cancelButtonText: '取消',
      },
    );
    await sortVirtualProductComment(row.id, {
      sort: Number(value),
      idempotency_key: key('sort', String(row.id)),
    });
    ElMessage.success('排序已更新');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

async function remove(row: ProductCommentRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      '删除仅对虚拟评论生效，将保留不可恢复的审核审计记录。',
      '删除虚拟评论',
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
    ElMessage.success('虚拟评论已删除');
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

function onImageSelect(items: AttachmentItem[]) {
  selectedImages.value = items;
  form.attachment_ids = items.map((item) => item.attachment_id);
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <ElButton :icon="Plus" type="primary" @click="open()">
          新增虚拟评论
        </ElButton>
      </template>
      <template #product="{ row }">
        {{ row.product_title || `商品 #${row.product_id}` }}
        <div class="text-xs text-gray-400">
          商品 #{{ row.product_id }} · 店铺 #{{ row.store_id }}
        </div>
      </template>
      <template #source="{ row }">
        <ElTag :type="row.source === 'virtual' ? 'warning' : 'info'">
          {{
            row.source === 'virtual'
              ? `虚拟：${row.virtual_author_name}`
              : '真实用户'
          }}
        </ElTag>
      </template>
      <template #score="{ row }">
        <ElRate :model-value="row.score" disabled />
      </template>
      <template #media="{ row }">
        <ElImage
          v-for="url in mediaUrls(row.media)"
          :key="url"
          :src="url"
          fit="cover"
          class="mr-1 h-8 w-8"
          :preview-src-list="mediaUrls(row.media)"
          preview-teleported
        />
        <span v-if="!mediaUrls(row.media).length">—</span>
      </template>
      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status).type">
          {{ statusInfo(row.status).label }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="row.status !== 'published'"
          link
          type="success"
          @click="moderate(row, 'publish')"
        >
          展示
        </ElButton>
        <ElButton
          v-if="row.status !== 'hidden'"
          link
          type="warning"
          @click="moderate(row, 'hide')"
        >
          隐藏
        </ElButton>
        <template v-if="row.source === 'virtual'">
          <ElButton link type="primary" @click="open(row)">编辑</ElButton>
          <ElButton link @click="sort(row)">排序</ElButton>
          <ElButton link type="danger" @click="remove(row)">删除</ElButton>
        </template>
      </template>
    </Grid>

    <FormModal>
      <ElAlert
        class="mb-4"
        type="warning"
        :closable="false"
        title="虚拟评论必须显著标识为虚拟来源；不得冒充真实用户或用于虚假交易宣传。"
      />
      <ElForm label-width="90px">
        <ElFormItem label="商品 ID" required>
          <ElInputNumber
            v-model="form.product_id"
            :disabled="Boolean(editing)"
            :min="1"
            class="w-full"
            controls-position="right"
          />
        </ElFormItem>
        <ElFormItem label="显示昵称" required>
          <ElInput
            v-model="form.virtual_author_name"
            maxlength="64"
            show-word-limit
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
          <ElButton @click="picker = true">从素材库选择（最多 9 张）</ElButton>
          <ElImage
            v-for="image in selectedImages"
            :key="image.attachment_id"
            :src="image.attachment_src"
            class="ml-2 h-8 w-8"
            fit="cover"
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
    </FormModal>

    <ImagePickerDialog
      v-model:open="picker"
      default-library="system"
      kind="image"
      :limit="9"
      @select="onImageSelect"
    />
  </Page>
</template>
