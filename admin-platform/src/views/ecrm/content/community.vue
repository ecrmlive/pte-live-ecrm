<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  auditCommunityPostApi,
  deleteCommunityPostApi,
  deleteCommunityReplyApi,
  getCommunityPostApi,
  listCommunityPostsApi,
  listCommunityRepliesApi,
  type CommunityPost,
  type CommunityReply,
} from '#/api/core/platform-community';

const loading = ref(false);
const saving = ref(false);
const rows = ref<CommunityPost[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const rejectOpen = ref(false);
const current = ref<CommunityPost>();
const replies = ref<CommunityReply[]>([]);
const repliesTotal = ref(0);
const canAudit = ref(false);
const canDelete = ref(false);
const query = reactive({ keyword: '', limit: 20, page: 1, status: undefined as number | undefined });
const replyQuery = reactive({ limit: 20, page: 1 });
const rejectForm = reactive({ refusal: '' });

function statusInfo(status: number) {
  if (status === 1) return { label: '审核通过', type: 'success' as const };
  if (status === -1) return { label: '已驳回', type: 'danger' as const };
  return { label: '待审核', type: 'warning' as const };
}

function canAuditRow(row: CommunityPost) {
  return canAudit.value && row.status === 0;
}

function canUpdatePresentation(row: CommunityPost) {
  return canAudit.value && row.status === 1;
}

async function load() {
  loading.value = true;
  try {
    const result = await listCommunityPostsApi(query);
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.keyword = '';
  query.status = undefined;
  query.page = 1;
  void load();
}

async function loadReplies() {
  if (!current.value) return;
  const result = await listCommunityRepliesApi(current.value.community_id, replyQuery);
  replies.value = result.list || [];
  repliesTotal.value = result.total || 0;
}

async function openDetail(row: CommunityPost) {
  replyQuery.page = 1;
  const [post, replyResult] = await Promise.all([
    getCommunityPostApi(row.community_id),
    listCommunityRepliesApi(row.community_id, replyQuery),
  ]);
  current.value = post;
  replies.value = replyResult.list || [];
  repliesTotal.value = replyResult.total || 0;
  detailOpen.value = true;
}

async function approve(row: CommunityPost) {
  try {
    await ElMessageBox.confirm('确认通过该社区帖子审核？通过后将按当前显示状态在 C 端生效。', '审核通过确认', {
      cancelButtonText: '取消',
      confirmButtonText: '确认通过',
      type: 'warning',
    });
    await auditCommunityPostApi(row.community_id, { is_show: row.is_show === 1 ? 1 : 0, status: 1 });
    ElMessage.success('帖子已审核通过');
    await load();
  } catch {
    // 用户取消或 requestClient 已统一提示接口异常。
  }
}

function openReject(row: CommunityPost) {
  current.value = row;
  rejectForm.refusal = '';
  rejectOpen.value = true;
}

async function reject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写驳回原因');
    return;
  }
  if (!current.value) return;
  saving.value = true;
  try {
    await auditCommunityPostApi(current.value.community_id, { refusal, status: -1 });
    rejectOpen.value = false;
    ElMessage.success('帖子已驳回');
    await load();
  } finally {
    saving.value = false;
  }
}

async function updatePresentation(row: CommunityPost, data: { is_hot?: 0 | 1; is_show?: 0 | 1 }, action: string) {
  try {
    await ElMessageBox.confirm(`确认${action}该已审核通过的帖子？`, `${action}确认`, {
      cancelButtonText: '取消',
      confirmButtonText: `确认${action}`,
      type: 'warning',
    });
    await auditCommunityPostApi(row.community_id, { ...data, status: 0 });
    ElMessage.success(`帖子已${action}`);
    await load();
  } catch {
    // 用户取消或 requestClient 已统一提示接口异常。
  }
}

async function deletePost(row: CommunityPost) {
  try {
    await ElMessageBox.confirm('确认删除该帖子？删除后会从 C 端隐藏且不可恢复。', '删除帖子', {
      cancelButtonText: '取消',
      confirmButtonText: '确认删除',
      type: 'error',
    });
    await deleteCommunityPostApi(row.community_id);
    ElMessage.success('帖子已删除');
    if (rows.value.length === 1 && query.page > 1) query.page -= 1;
    await load();
  } catch {
    // 用户取消或 requestClient 已统一提示接口异常。
  }
}

async function deleteReply(row: CommunityReply) {
  try {
    await ElMessageBox.confirm('确认删除该评论？删除后不可恢复。', '删除评论', {
      cancelButtonText: '取消',
      confirmButtonText: '确认删除',
      type: 'error',
    });
    await deleteCommunityReplyApi(row.reply_id);
    ElMessage.success('评论已删除');
    await loadReplies();
    if (current.value) current.value = await getCommunityPostApi(current.value.community_id);
    await load();
  } catch {
    // 用户取消或 requestClient 已统一提示接口异常。
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  const isContentOperator = profile.roles.includes('platform') || profile.roles.includes('operations');
  canAudit.value = isContentOperator && permissions.includes('content.community.audit');
  canDelete.value = isContentOperator && permissions.includes('content.community.delete');
  await load();
});
</script>

<template>
  <Page title="社区种草监管" description="有内容审核权限的平台账号可审核、驳回、控制展示与置顶；有删帖权限的账号可删除违规帖子及评论。">
    <el-card shadow="never">
      <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
        <el-form-item label="关键词"><el-input v-model="query.keyword" clearable placeholder="标题或正文" /></el-form-item>
        <el-form-item label="审核状态">
          <el-select v-model="query.status" class="w-36" clearable placeholder="全部状态">
            <el-option label="待审核" :value="0" />
            <el-option label="审核通过" :value="1" />
            <el-option label="已驳回" :value="-1" />
          </el-select>
        </el-form-item>
        <el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button></el-form-item>
      </el-form>
    </el-card>

    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="community_id">
        <el-table-column label="标题" min-width="200" prop="title" show-overflow-tooltip />
        <el-table-column label="作者" min-width="110" prop="nickname" />
        <el-table-column label="商户 ID" width="100" prop="mer_id" />
        <el-table-column label="话题" min-width="120" prop="topic_name" />
        <el-table-column label="评论" width="76" prop="count_reply" />
        <el-table-column label="浏览" width="76" prop="pv" />
        <el-table-column label="审核" width="104"><template #default="{ row }"><el-tag :type="statusInfo(row.status).type">{{ statusInfo(row.status).label }}</el-tag></template></el-table-column>
        <el-table-column label="展示" width="76"><template #default="{ row }"><el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '显示' : '隐藏' }}</el-tag></template></el-table-column>
        <el-table-column label="置顶" width="76"><template #default="{ row }"><el-tag :type="row.is_hot === 1 ? 'warning' : 'info'">{{ row.is_hot === 1 ? '是' : '否' }}</el-tag></template></el-table-column>
        <el-table-column label="发布时间" min-width="170" prop="create_time" />
        <el-table-column fixed="right" label="操作" width="290">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
            <template v-if="canAuditRow(row)"><el-button link type="success" @click="approve(row)">通过</el-button><el-button link type="danger" @click="openReject(row)">驳回</el-button></template>
            <template v-else-if="canUpdatePresentation(row)">
              <el-button link type="warning" @click="updatePresentation(row, { is_hot: row.is_hot === 1 ? 0 : 1 }, row.is_hot === 1 ? '取消置顶' : '置顶')">{{ row.is_hot === 1 ? '取消置顶' : '置顶' }}</el-button>
              <el-button link type="warning" @click="updatePresentation(row, { is_show: row.is_show === 1 ? 0 : 1 }, row.is_show === 1 ? '隐藏' : '显示')">{{ row.is_show === 1 ? '隐藏' : '显示' }}</el-button>
            </template>
            <el-button v-if="canDelete" link type="danger" @click="deletePost(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page: number) => { query.page = page; load(); }" @size-change="(limit: number) => { query.limit = limit; query.page = 1; load(); }" /></div>
    </el-card>

    <el-drawer v-model="detailOpen" :with-header="false" size="680px">
      <template v-if="current">
        <div class="mb-5 text-lg font-medium">社区帖子详情</div>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="标题" :span="2">{{ current.title }}</el-descriptions-item>
          <el-descriptions-item label="作者">{{ current.nickname || '—' }}</el-descriptions-item>
          <el-descriptions-item label="商户 ID">{{ current.mer_id || '平台内容' }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{ current.cate_name || '—' }}</el-descriptions-item>
          <el-descriptions-item label="话题">{{ current.topic_name || '—' }}</el-descriptions-item>
          <el-descriptions-item label="关联商品">{{ current.product_name || (current.product_id ? `商品 #${current.product_id}` : '—') }}</el-descriptions-item>
          <el-descriptions-item label="审核状态"><el-tag :type="statusInfo(current.status).type">{{ statusInfo(current.status).label }}</el-tag></el-descriptions-item>
          <el-descriptions-item label="发布时间">{{ current.create_time }}</el-descriptions-item>
          <el-descriptions-item label="审核时间">{{ current.status_time || '—' }}</el-descriptions-item>
          <el-descriptions-item :span="2" label="正文"><div class="whitespace-pre-wrap">{{ current.content }}</div></el-descriptions-item>
          <el-descriptions-item v-if="current.image" :span="2" label="图片"><el-image :preview-src-list="[current.image]" :src="current.image" class="h-24 w-24" fit="cover" preview-teleported /></el-descriptions-item>
          <el-descriptions-item v-if="current.refusal" :span="2" label="驳回原因">{{ current.refusal }}</el-descriptions-item>
        </el-descriptions>
        <div class="mb-3 mt-6 text-base font-medium">评论（{{ repliesTotal }}）</div>
        <el-table :data="replies" border empty-text="暂无评论">
          <el-table-column label="用户" min-width="110" prop="nickname" />
          <el-table-column label="内容" min-width="220" prop="content" show-overflow-tooltip />
          <el-table-column label="时间" min-width="170" prop="create_time" />
          <el-table-column v-if="canDelete" fixed="right" label="操作" width="72"><template #default="{ row }"><el-button link type="danger" @click="deleteReply(row)">删除</el-button></template></el-table-column>
        </el-table>
        <div class="mt-3 flex justify-end"><el-pagination small :current-page="replyQuery.page" :page-size="replyQuery.limit" :total="repliesTotal" layout="prev, pager, next" @current-change="(page: number) => { replyQuery.page = page; loadReplies(); }" /></div>
      </template>
    </el-drawer>

    <el-dialog v-model="rejectOpen" destroy-on-close title="驳回帖子" width="480px">
      <el-form label-width="84px"><el-form-item label="驳回原因" required><el-input v-model="rejectForm.refusal" :rows="4" maxlength="200" placeholder="请填写可供发布者查看的驳回原因" show-word-limit type="textarea" /></el-form-item></el-form>
      <template #footer><el-button @click="rejectOpen = false">取消</el-button><el-button :loading="saving" type="danger" @click="reject">确认驳回</el-button></template>
    </el-dialog>
  </Page>
</template>
