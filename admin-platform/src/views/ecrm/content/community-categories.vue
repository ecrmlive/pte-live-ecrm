<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import { listCommunityCategoriesApi, type CommunityCategory } from '#/api/core/platform-community';

const loading = ref(false);
const rows = ref<CommunityCategory[]>([]);
const router = useRouter();

async function load() {
  loading.value = true;
  try {
    const result = await listCommunityCategoriesApi();
    rows.value = result.list || [];
  } finally {
    loading.value = false;
  }
}

function openPosts(row: CommunityCategory) {
  void router.push({ path: '/community/list', query: { category_id: String(row.category_id) } });
}

onMounted(() => void load());
</script>

<template>
  <Page title="社区分类" description="只读查看 C 端社区分类；分类维护入口待后续补齐，可先跳转社区内容页监管帖子。">
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="category_id">
        <el-table-column label="ID" prop="category_id" width="90" />
        <el-table-column label="分类名称" min-width="160" prop="cate_name" />
        <el-table-column label="上级 ID" prop="pid" width="100" />
        <el-table-column label="排序" prop="sort" width="90" />
        <el-table-column label="展示" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '显示' : '隐藏' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="120">
          <template #default="{ row }">
            <el-button link type="primary" @click="openPosts(row)">查看帖子</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </Page>
</template>
