<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElAlert,
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElForm,
  ElFormItem,
  ElImage,
  ElInput,
  ElMessage,
  ElOption,
  ElSelect,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  listMerchantProductsApi,
  type MerchantProduct,
} from '#/api/core/merchant-catalog';
import {
  createMerchantCommunityPostApi,
  deleteMerchantCommunityPostApi,
  getMerchantCommunityPostApi,
  listMerchantCommunityCategoriesApi,
  listMerchantCommunityPostsApi,
  listMerchantCommunityRepliesApi,
  listMerchantCommunityTopicsApi,
  updateMerchantCommunityPostApi,
  type MerchantCommunityCategory,
  type MerchantCommunityPost,
  type MerchantCommunityPostInput,
  type MerchantCommunityReply,
  type MerchantCommunityTopic,
} from '#/api/core/merchant-community';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import ImageField from '#/components/shop/image-field.vue';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const saving = ref(false);
const editingID = ref<number>();
const categories = ref<MerchantCommunityCategory[]>([]);
const topics = ref<MerchantCommunityTopic[]>([]);
const products = ref<MerchantProduct[]>([]);
const current = ref<MerchantCommunityPost>();
const replies = ref<MerchantCommunityReply[]>([]);
const repliesTotal = ref(0);
const canCreate = ref(false);
const canUpdate = ref(false);
const canDelete = ref(false);
const replyQuery = reactive({ limit: 20, page: 1 });
const form = reactive<MerchantCommunityPostInput>({
  category_id: 0,
  content: '',
  image: '',
  product_id: 0,
  title: '',
  topic_id: 0,
});

const filteredTopics = computed(() =>
  topics.value.filter(
    (item) => !form.category_id || item.category_id === form.category_id,
  ),
);

function statusInfo(status: number) {
  if (status === 1) return { label: '审核通过', type: 'success' as const };
  if (status === -1) return { label: '已驳回', type: 'danger' as const };
  return { label: '待平台审核', type: 'warning' as const };
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, {
    category_id: 0,
    content: '',
    image: '',
    product_id: 0,
    title: '',
    topic_id: 0,
  });
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '标题 / 话题 / 商品' },
    fieldName: 'keyword',
    label: '帖子搜索',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待平台审核', value: 0 },
        { label: '审核通过', value: 1 },
        { label: '已驳回', value: -1 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<MerchantCommunityPost> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'community_id', title: 'ID', width: 80 },
    { field: 'title', minWidth: 220, showOverflow: false, title: '标题' },
    { field: 'topic_name', minWidth: 130, showOverflow: false, title: '话题' },
    {
      field: 'product_name',
      minWidth: 160,
      showOverflow: false,
      slots: { default: 'product' },
      title: '关联商品',
    },
    { field: 'count_reply', title: '评论', width: 76 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '审核状态',
      width: 118,
    },
    { field: 'refusal', minWidth: 160, showOverflow: false, title: '驳回原因' },
    {
      field: 'create_time',
      minWidth: 170,
      title: '发布时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 180 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = formValues?.status;
        const data = await listMerchantCommunityPostsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status:
            status === 0 || status === 1 || status === -1
              ? Number(status)
              : undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'community_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [EditDrawer, editDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  confirmText: '完成',
  cancelText: '取消',
  placement: 'right',
  onConfirm: savePost,
});

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[680px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

function openCreate() {
  resetForm();
  editDrawerApi.setState({ title: '发布帖子' }).open();
}

async function openEdit(row: MerchantCommunityPost) {
  const detail = await getMerchantCommunityPostApi(row.community_id);
  editingID.value = detail.community_id;
  Object.assign(form, {
    category_id: detail.category_id,
    content: detail.content,
    image: detail.image || '',
    product_id: detail.product_id || 0,
    title: detail.title,
    topic_id: detail.topic_id,
  });
  editDrawerApi.setState({ title: '编辑帖子' }).open();
}

function categoryChanged() {
  if (
    form.topic_id &&
    !filteredTopics.value.some((item) => item.topic_id === form.topic_id)
  ) {
    form.topic_id = 0;
  }
}

async function savePost() {
  if (!form.title.trim() || !form.content.trim()) {
    ElMessage.warning('请填写标题和正文');
    return;
  }
  saving.value = true;
  editDrawerApi.lock();
  try {
    const body: MerchantCommunityPostInput = {
      category_id: form.category_id,
      content: form.content.trim(),
      image: form.image,
      product_id: form.product_id,
      title: form.title.trim(),
      topic_id: form.topic_id,
    };
    if (editingID.value) {
      await updateMerchantCommunityPostApi(editingID.value, body);
    } else {
      await createMerchantCommunityPostApi(body);
    }
    editDrawerApi.close();
    ElMessage.success(
      editingID.value ? '帖子已更新，已重新提交平台审核' : '帖子已提交平台审核',
    );
    gridApi.reload();
  } finally {
    saving.value = false;
    editDrawerApi.unlock();
  }
}

async function openDetail(row: MerchantCommunityPost) {
  replyQuery.page = 1;
  detailDrawerApi.setState({ title: '帖子详情', loading: true }).open();
  try {
    const [post, replyResult] = await Promise.all([
      getMerchantCommunityPostApi(row.community_id),
      listMerchantCommunityRepliesApi(row.community_id, replyQuery),
    ]);
    current.value = post;
    replies.value = replyResult.list || [];
    repliesTotal.value = replyResult.total || 0;
  } finally {
    detailDrawerApi.setState({ loading: false });
  }
}

async function loadReplies() {
  if (!current.value) return;
  const result = await listMerchantCommunityRepliesApi(
    current.value.community_id,
    replyQuery,
  );
  replies.value = result.list || [];
  repliesTotal.value = result.total || 0;
}

async function remove(row: MerchantCommunityPost) {
  try {
    await confirm({
      content: `确认删除帖子“${row.title}”？删除后不可恢复。`,
      icon: 'warning',
      title: '删除确认',
    });
    await deleteMerchantCommunityPostApi(row.community_id);
    ElMessage.success('帖子已删除');
    gridApi.reload();
  } catch {
    // cancelled
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canCreate.value = permissions.includes('community/create');
  canUpdate.value = permissions.includes('community/update');
  canDelete.value = permissions.includes('community/delete');
  const [categoryResult, topicResult, productResult] = await Promise.all([
    listMerchantCommunityCategoriesApi(),
    listMerchantCommunityTopicsApi(),
    listMerchantProductsApi({ limit: 100, page: 1, status: 1 }),
  ]);
  categories.value = categoryResult.list || [];
  topics.value = topicResult.list || [];
  products.value = productResult.list || [];
});
</script>

<template>
  <Page auto-content-height>
    <template v-if="canCreate" #extra>
      <ElButton type="primary" @click="openCreate">发布帖子</ElButton>
    </template>

    <Grid>
      <template #product="{ row }">
        {{
          row.product_name ||
          (row.product_id ? `商品 #${row.product_id}` : '未关联')
        }}
      </template>
      <template #status="{ row }">
        <ElTag :type="statusInfo(row.status).type">
          {{ statusInfo(row.status).label }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">
          详情/评论
        </ElButton>
        <ElButton v-if="canUpdate" link type="primary" @click="openEdit(row)">
          编辑
        </ElButton>
        <ElButton v-if="canDelete" link type="danger" @click="remove(row)">
          删除
        </ElButton>
      </template>
    </Grid>

    <EditDrawer class="w-[720px] max-w-[96vw]">
      <ElAlert
        class="mb-4"
        :closable="false"
        show-icon
        title="发布或编辑后均需平台重新审核，审核通过后才会在 C 端展示。"
        type="info"
      />
      <ElForm class="grid grid-cols-2 gap-x-4" label-width="92px">
        <ElFormItem class="col-span-2" label="标题" required>
          <ElInput
            v-model="form.title"
            maxlength="100"
            placeholder="请输入种草标题"
            show-word-limit
          />
        </ElFormItem>
        <ElFormItem label="内容分类">
          <ElSelect
            v-model="form.category_id"
            clearable
            class="w-full"
            placeholder="可选"
            @change="categoryChanged"
          >
            <ElOption
              v-for="item in categories"
              :key="item.category_id"
              :label="item.cate_name"
              :value="item.category_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem label="关联话题">
          <ElSelect
            v-model="form.topic_id"
            clearable
            class="w-full"
            placeholder="可选"
          >
            <ElOption
              v-for="item in filteredTopics"
              :key="item.topic_id"
              :label="item.topic_name"
              :value="item.topic_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem class="col-span-2" label="关联商品">
          <ElSelect
            v-model="form.product_id"
            clearable
            filterable
            class="w-full"
            placeholder="可选，仅可关联本店已审核商品"
          >
            <ElOption
              v-for="item in products"
              :key="item.product_id"
              :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`"
              :value="item.product_id"
            />
          </ElSelect>
        </ElFormItem>
        <ElFormItem class="col-span-2" label="封面图片">
          <ImageField
            v-model="form.image"
            button-text="从素材库选择封面"
            hint="支持从本店素材库上传或选择图片"
          />
        </ElFormItem>
        <ElFormItem class="col-span-2" label="正文" required>
          <ElInput
            v-model="form.content"
            :rows="8"
            maxlength="5000"
            placeholder="分享商品使用体验、攻略或服务内容"
            show-word-limit
            type="textarea"
          />
        </ElFormItem>
      </ElForm>
    </EditDrawer>

    <DetailDrawer>
      <template v-if="current">
        <ElDescriptions :column="2" border>
          <ElDescriptionsItem :span="2" label="标题">
            {{ current.title }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="分类">
            {{ current.cate_name || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="话题">
            {{ current.topic_name || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="关联商品">
            {{
              current.product_name ||
              (current.product_id ? `商品 #${current.product_id}` : '未关联')
            }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="审核状态">
            <ElTag :type="statusInfo(current.status).type">
              {{ statusInfo(current.status).label }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="发布时间">
            {{ formatShanghaiDateTime(current.create_time) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :span="2" label="正文">
            <div class="whitespace-pre-wrap">{{ current.content }}</div>
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="current.image" :span="2" label="封面">
            <ElImage
              :preview-src-list="[current.image]"
              :src="current.image"
              class="h-24 w-24"
              fit="cover"
              preview-teleported
            />
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="current.refusal" :span="2" label="驳回原因">
            {{ current.refusal }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <div class="mb-3 mt-6 text-base font-medium">
          评论（{{ repliesTotal }}）
        </div>
        <ElTable :data="replies" border empty-text="暂无评论">
          <ElTableColumn label="用户" min-width="120" prop="nickname" />
          <ElTableColumn
            label="内容"
            min-width="250"
            prop="content"
            show-overflow-tooltip
          />
          <ElTableColumn label="时间" min-width="170">
            <template #default="{ row }">
              {{ formatShanghaiDateTime(row.create_time) }}
            </template>
          </ElTableColumn>
        </ElTable>
      </template>
    </DetailDrawer>
  </Page>
</template>
