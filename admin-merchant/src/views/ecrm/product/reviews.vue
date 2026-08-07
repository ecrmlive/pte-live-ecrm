<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElRate,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listMerchantProductCommentsApi,
  replyMerchantProductCommentApi,
  type MerchantProductComment,
} from '#/api/core/merchant-product-comment';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const replySaving = ref(false);
const current = ref<MerchantProductComment>();
const replyContent = ref('');

const statusLabels: Record<string, string> = {
  hidden: '已隐藏',
  pending: '待审核',
  published: '已展示',
};

const statusTypes: Record<
  string,
  'danger' | 'info' | 'success' | 'warning'
> = {
  hidden: 'info',
  pending: 'warning',
  published: 'success',
};

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '评论内容 / 商品名称',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'InputNumber',
    componentProps: { clearable: true, min: 1, placeholder: '商品 ID' },
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
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<MerchantProductComment> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'product_title',
      minWidth: 180,
      showOverflow: false,
      slots: { default: 'product' },
      title: '商品',
    },
    { field: 'user_id', title: '用户 ID', width: 90 },
    {
      field: 'score',
      slots: { default: 'score' },
      title: '评分',
      width: 120,
    },
    { field: 'content', minWidth: 220, showOverflow: true, title: '评论' },
    { field: 'reply_content', minWidth: 160, showOverflow: true, title: '回复' },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 96,
    },
    {
      field: 'created_at',
      minWidth: 170,
      title: '时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 76 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const productId = formValues?.product_id;
        const status = String(formValues?.status ?? '').trim();
        const data = await listMerchantProductCommentsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          product_id:
            typeof productId === 'number' && productId > 0
              ? Number(productId)
              : undefined,
          status: (status || undefined) as
            | 'hidden'
            | 'pending'
            | 'published'
            | undefined,
          date_from: range[0],
          date_to: range[1],
        });
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

const [ReplyDrawer, replyDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: async () => {
    if (!current.value || !replyContent.value.trim()) {
      ElMessage.warning('请填写回复内容');
      return;
    }
    replySaving.value = true;
    replyDrawerApi.lock();
    try {
      await replyMerchantProductCommentApi(current.value.id, {
        reply_content: replyContent.value.trim(),
      });
      ElMessage.success('商家回复已保存');
      replyDrawerApi.close();
      gridApi.reload();
    } finally {
      replySaving.value = false;
      replyDrawerApi.unlock();
    }
  },
});

function openReply(row: MerchantProductComment) {
  current.value = row;
  replyContent.value = row.reply_content || '';
  replyDrawerApi.setState({ title: '商家回复' }).open();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #product="{ row }">
        {{ row.product_title || `#${row.product_id}` }}
      </template>
      <template #score="{ row }">
        <ElRate :model-value="row.score" disabled />
      </template>
      <template #status="{ row }">
        <ElTag :type="statusTypes[row.status] || 'info'">
          {{ statusLabels[row.status] || row.status }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openReply(row)">回复</ElButton>
      </template>
    </Grid>

    <ReplyDrawer class="w-[520px] max-w-[96vw]">
      <ElForm label-width="72px">
        <ElFormItem label="回复内容" required>
          <ElInput
            v-model="replyContent"
            :rows="4"
            maxlength="500"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </ReplyDrawer>
  </Page>
</template>
