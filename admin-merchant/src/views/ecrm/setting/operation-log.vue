<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { listMerchantOperationLogsApi, type MerchantOperationLog } from '#/api/core/merchant-operation-log';

const loading = ref(false);
const rows = ref<MerchantOperationLog[]>([]);
const total = ref(0);
const query = reactive({ page: 1, limit: 20 });

async function load() {
  loading.value = true;
  try {
    const data = await listMerchantOperationLogsApi({ ...query });
    rows.value = data.list ?? [];
    total.value = data.total ?? 0;
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="操作日志" description="店铺后台写操作审计记录。">
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="account_id" label="操作人" width="90" />
        <el-table-column prop="action" label="动作" min-width="140" />
        <el-table-column prop="resource_type" label="资源类型" width="120" />
        <el-table-column prop="resource_id" label="资源 ID" width="120" />
        <el-table-column prop="request_id" label="请求 ID" min-width="160" show-overflow-tooltip />
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total,prev,pager,next" @current-change="(page) => { query.page = page; load(); }" /></div>
    </el-card>
  </Page>
</template>
