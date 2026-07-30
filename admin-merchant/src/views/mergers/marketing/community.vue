<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
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
import { listMerchantProductsApi, type MerchantProduct } from '#/api/core/merchant-catalog';
import ImageField from '#/components/shop/image-field.vue';

const loading = ref(false);
const saving = ref(false);
const rows = ref<MerchantCommunityPost[]>([]);
const total = ref(0);
const categories = ref<MerchantCommunityCategory[]>([]);
const topics = ref<MerchantCommunityTopic[]>([]);
const products = ref<MerchantProduct[]>([]);
const dialogOpen = ref(false);
const detailOpen = ref(false);
const editingID = ref<number>();
const current = ref<MerchantCommunityPost>();
const replies = ref<MerchantCommunityReply[]>([]);
const repliesTotal = ref(0);
const canCreate = ref(false);
const canUpdate = ref(false);
const canDelete = ref(false);
const query = reactive({ limit: 20, page: 1 });
const replyQuery = reactive({ limit: 20, page: 1 });
const form = reactive<MerchantCommunityPostInput>({ category_id: 0, content: '', image: '', product_id: 0, title: '', topic_id: 0 });

const filteredTopics = computed(() => topics.value.filter((item) => !form.category_id || item.category_id === form.category_id));

function statusInfo(status: number) {
  if (status === 1) return { label: '审核通过', type: 'success' as const };
  if (status === -1) return { label: '已驳回', type: 'danger' as const };
  return { label: '待平台审核', type: 'warning' as const };
}

function resetForm() {
  editingID.value = undefined;
  Object.assign(form, { category_id: 0, content: '', image: '', product_id: 0, title: '', topic_id: 0 });
}

function openCreate() {
  resetForm();
  dialogOpen.value = true;
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
  dialogOpen.value = true;
}

async function load() {
  loading.value = true;
  try {
    const result = await listMerchantCommunityPostsApi(query);
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

async function loadOptions() {
  const [categoryResult, topicResult, productResult] = await Promise.all([
    listMerchantCommunityCategoriesApi(),
    listMerchantCommunityTopicsApi(),
    listMerchantProductsApi({ limit: 100, page: 1, status: 1 }),
  ]);
  categories.value = categoryResult.list || [];
  topics.value = topicResult.list || [];
  products.value = productResult.list || [];
}

function categoryChanged() {
  if (form.topic_id && !filteredTopics.value.some((item) => item.topic_id === form.topic_id)) form.topic_id = 0;
}

async function save() {
  if (!form.title.trim() || !form.content.trim()) {
    ElMessage.warning('请填写标题和正文');
    return;
  }
  saving.value = true;
  try {
    const body: MerchantCommunityPostInput = {
      category_id: form.category_id,
      content: form.content.trim(),
      image: form.image,
      product_id: form.product_id,
      title: form.title.trim(),
      topic_id: form.topic_id,
    };
    if (editingID.value) await updateMerchantCommunityPostApi(editingID.value, body);
    else await createMerchantCommunityPostApi(body);
    dialogOpen.value = false;
    ElMessage.success(editingID.value ? '帖子已更新，已重新提交平台审核' : '帖子已提交平台审核');
    await load();
  } finally {
    saving.value = false;
  }
}

async function openDetail(row: MerchantCommunityPost) {
  replyQuery.page = 1;
  const [post, replyResult] = await Promise.all([
    getMerchantCommunityPostApi(row.community_id),
    listMerchantCommunityRepliesApi(row.community_id, replyQuery),
  ]);
  current.value = post;
  replies.value = replyResult.list || [];
  repliesTotal.value = replyResult.total || 0;
  detailOpen.value = true;
}

async function loadReplies() {
  if (!current.value) return;
  const result = await listMerchantCommunityRepliesApi(current.value.community_id, replyQuery);
  replies.value = result.list || [];
  repliesTotal.value = result.total || 0;
}

async function remove(row: MerchantCommunityPost) {
  try {
    await ElMessageBox.confirm(`确认删除帖子“${row.title}”？删除后不可恢复。`, '删除确认', {
      cancelButtonText: '取消',
      confirmButtonText: '确认删除',
      type: 'warning',
    });
    await deleteMerchantCommunityPostApi(row.community_id);
    ElMessage.success('帖子已删除');
    if (rows.value.length === 1 && query.page > 1) query.page -= 1;
    await load();
  } catch {
    // 用户取消或 requestClient 已统一提示接口异常。
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canCreate.value = permissions.includes('community/create');
  canUpdate.value = permissions.includes('community/update');
  canDelete.value = permissions.includes('community/delete');
  await Promise.all([load(), loadOptions()]);
});
</script>

<template>
  <Page title="逛逛社区" description="发布本店种草图文并关联本店商品；新建和编辑均会进入待平台审核，审核通过后才会在 C 端展示。">
    <template #extra><el-button v-if="canCreate" type="primary" @click="openCreate">发布帖子</el-button></template>
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="community_id">
        <el-table-column label="标题" min-width="220" prop="title" show-overflow-tooltip />
        <el-table-column label="话题" min-width="130" prop="topic_name" />
        <el-table-column label="关联商品" min-width="160"><template #default="{ row }">{{ row.product_name || (row.product_id ? `商品 #${row.product_id}` : '未关联') }}</template></el-table-column>
        <el-table-column label="评论" width="76" prop="count_reply" />
        <el-table-column label="审核状态" width="118"><template #default="{ row }"><el-tag :type="statusInfo(row.status).type">{{ statusInfo(row.status).label }}</el-tag></template></el-table-column>
        <el-table-column label="驳回原因" min-width="160" prop="refusal" show-overflow-tooltip />
        <el-table-column label="发布时间" min-width="170" prop="create_time" />
        <el-table-column fixed="right" label="操作" width="180">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情/评论</el-button>
            <el-button v-if="canUpdate" link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button v-if="canDelete" link type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div>
    </el-card>

    <el-dialog v-model="dialogOpen" :title="editingID ? '编辑帖子' : '发布帖子'" destroy-on-close width="720px">
      <el-alert class="mb-4" :closable="false" show-icon title="发布或编辑后均需平台重新审核，审核通过后才会在 C 端展示。" type="info" />
      <el-form class="grid grid-cols-2 gap-x-4" label-width="92px">
        <el-form-item class="col-span-2" label="标题" required><el-input v-model="form.title" maxlength="100" placeholder="请输入种草标题" show-word-limit /></el-form-item>
        <el-form-item label="内容分类"><el-select v-model="form.category_id" clearable class="w-full" placeholder="可选" @change="categoryChanged"><el-option v-for="item in categories" :key="item.category_id" :label="item.cate_name" :value="item.category_id" /></el-select></el-form-item>
        <el-form-item label="关联话题"><el-select v-model="form.topic_id" clearable class="w-full" placeholder="可选"><el-option v-for="item in filteredTopics" :key="item.topic_id" :label="item.topic_name" :value="item.topic_id" /></el-select></el-form-item>
        <el-form-item class="col-span-2" label="关联商品"><el-select v-model="form.product_id" clearable filterable class="w-full" placeholder="可选，仅可关联本店已审核商品"><el-option v-for="item in products" :key="item.product_id" :label="`${item.store_name}（¥${Number(item.price).toFixed(2)}）`" :value="item.product_id" /></el-select></el-form-item>
        <el-form-item class="col-span-2" label="封面图片"><ImageField v-model="form.image" button-text="从素材库选择封面" hint="支持从本店素材库上传或选择图片" /></el-form-item>
        <el-form-item class="col-span-2" label="正文" required><el-input v-model="form.content" :rows="8" maxlength="5000" placeholder="分享商品使用体验、攻略或服务内容" show-word-limit type="textarea" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="dialogOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="save">提交审核</el-button></template>
    </el-dialog>

    <el-drawer v-model="detailOpen" :with-header="false" size="680px">
      <template v-if="current">
        <div class="mb-5 text-lg font-medium">帖子详情</div>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="标题" :span="2">{{ current.title }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{ current.cate_name || '—' }}</el-descriptions-item>
          <el-descriptions-item label="话题">{{ current.topic_name || '—' }}</el-descriptions-item>
          <el-descriptions-item label="关联商品" :span="2">{{ current.product_name || (current.product_id ? `商品 #${current.product_id}` : '未关联') }}</el-descriptions-item>
          <el-descriptions-item label="审核状态"><el-tag :type="statusInfo(current.status).type">{{ statusInfo(current.status).label }}</el-tag></el-descriptions-item>
          <el-descriptions-item label="发布时间">{{ current.create_time }}</el-descriptions-item>
          <el-descriptions-item :span="2" label="正文"><div class="whitespace-pre-wrap">{{ current.content }}</div></el-descriptions-item>
          <el-descriptions-item v-if="current.image" :span="2" label="封面"><el-image :preview-src-list="[current.image]" :src="current.image" class="h-24 w-24" fit="cover" preview-teleported /></el-descriptions-item>
          <el-descriptions-item v-if="current.refusal" :span="2" label="驳回原因">{{ current.refusal }}</el-descriptions-item>
        </el-descriptions>
        <div class="mb-3 mt-6 text-base font-medium">评论（{{ repliesTotal }}）</div>
        <el-table :data="replies" border empty-text="暂无评论"><el-table-column label="用户" min-width="120" prop="nickname" /><el-table-column label="内容" min-width="250" prop="content" show-overflow-tooltip /><el-table-column label="时间" min-width="170" prop="create_time" /></el-table>
        <div class="mt-3 flex justify-end"><el-pagination small :current-page="replyQuery.page" :page-size="replyQuery.limit" :total="repliesTotal" layout="prev, pager, next" @current-change="(page) => { replyQuery.page = page; loadReplies(); }" /></div>
      </template>
    </el-drawer>
  </Page>
</template>
