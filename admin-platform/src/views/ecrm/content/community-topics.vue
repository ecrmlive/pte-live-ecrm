<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import { listCommunityTopicsApi, type CommunityTopic } from '#/api/core/platform-community';

const loading = ref(false);
const rows = ref<CommunityTopic[]>([]);
const router = useRouter();

async function load() {
  loading.value = true;
  try {
    const result = await listCommunityTopicsApi();
    rows.value = result.list || [];
  } finally {
    loading.value = false;
  }
}

function openPosts(row: CommunityTopic) {
  void router.push({ path: '/community/list', query: { topic_id: String(row.topic_id) } });
}

onMounted(() => void load());
</script>

<template>
  <Page title="社区话题" description="只读查看社区话题与引用次数；话题维护入口待后续补齐，可先跳转社区内容页监管帖子。">
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="topic_id">
        <el-table-column label="ID" prop="topic_id" width="90" />
        <el-table-column label="话题名称" min-width="180" prop="topic_name" />
        <el-table-column label="分类 ID" prop="category_id" width="100" />
        <el-table-column label="引用次数" prop="count_use" width="100" />
        <el-table-column label="排序" prop="sort" width="90" />
        <el-table-column label="推荐" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_hot === 1 ? 'warning' : 'info'">{{ row.is_hot === 1 ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" min-width="170" prop="create_time" />
        <el-table-column fixed="right" label="操作" width="120">
          <template #default="{ row }">
            <el-button link type="primary" @click="openPosts(row)">查看帖子</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </Page>
</template>
