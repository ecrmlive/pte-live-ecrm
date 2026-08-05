<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  deleteCommunityReplyApi,
  listAllCommunityRepliesApi,
  type CommunityReply,
} from '#/api/core/platform-community';

const loading = ref(false);
const rows = ref<CommunityReply[]>([]);
const total = ref(0);
const canDelete = ref(false);
const query = reactive({ keyword: '', limit: 20, page: 1 });

async function load() {
  loading.value = true;
  try {
    const result = await listAllCommunityRepliesApi(query);
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
  query.page = 1;
  void load();
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
    if (rows.value.length === 1 && query.page > 1) query.page -= 1;
    await load();
  } catch {
    // 用户取消或 requestClient 已统一提示接口异常。
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  const isContentOperator = profile.roles.includes('platform') || profile.roles.includes('operations');
  canDelete.value = isContentOperator && permissions.includes('content.community.delete');
  await load();
});
</script>

<template>
  <Page title="社区评论" description="监管全站社区评论；有删帖权限的账号可删除违规评论。">
    <el-card shadow="never">
      <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
        <el-form-item label="关键词"><el-input v-model="query.keyword" clearable placeholder="评论内容" /></el-form-item>
        <el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button></el-form-item>
      </el-form>
    </el-card>

    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="reply_id">
        <el-table-column label="ID" prop="reply_id" width="90" />
        <el-table-column label="帖子" min-width="180" prop="post_title" show-overflow-tooltip />
        <el-table-column label="帖子 ID" width="100" prop="community_id" />
        <el-table-column label="用户" min-width="110" prop="nickname" />
        <el-table-column label="内容" min-width="260" prop="content" show-overflow-tooltip />
        <el-table-column label="时间" min-width="170" prop="create_time" />
        <el-table-column v-if="canDelete" fixed="right" label="操作" width="72">
          <template #default="{ row }"><el-button link type="danger" @click="deleteReply(row)">删除</el-button></template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination
          :current-page="query.page"
          :page-size="query.limit"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          background
          layout="total, sizes, prev, pager, next"
          @current-change="(page: number) => { query.page = page; load(); }"
          @size-change="(limit: number) => { query.limit = limit; query.page = 1; load(); }"
        />
      </div>
    </el-card>
  </Page>
</template>
