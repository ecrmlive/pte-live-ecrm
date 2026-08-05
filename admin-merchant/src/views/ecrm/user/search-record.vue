<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { listUserSearchRecordsApi, type UserSearchRecord } from '#/api/core/merchant-search-record';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<UserSearchRecord[]>([]);
const total = ref(0);
const query = reactive({ page: 1, limit: 20 });

async function load() {
  loading.value = true;
  try {
    const data = await listUserSearchRecordsApi({ ...query });
    rows.value = data.list ?? [];
    total.value = data.total ?? 0;
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="搜索记录" description="用户在本店商品的浏览记录（业务库 browse_history 投影）。">
    <EcrmListPage title="浏览记录">
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户 ID" width="90" />
        <el-table-column label="商品" min-width="200"><template #default="{ row }">{{ row.product_title || `#${row.product_id}` }}</template></el-table-column>
        <el-table-column prop="viewed_at" label="浏览时间" width="180" />
      </el-table>
      <template #pager><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total,prev,pager,next" @current-change="(page) => { query.page = page; load(); }" /></template>
    </EcrmListPage>
  </Page>
</template>
