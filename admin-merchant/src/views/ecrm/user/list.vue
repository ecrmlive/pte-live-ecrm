<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { listStoreCustomersApi, type StoreCustomer } from '#/api/core/merchant-store-customer';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<StoreCustomer[]>([]);
const total = ref(0);
const query = reactive({ page: 1, limit: 20, keyword: '' });

async function load() {
  loading.value = true;
  try {
    const data = await listStoreCustomersApi({ page: query.page, limit: query.limit, keyword: query.keyword.trim() || undefined });
    rows.value = data.list ?? [];
    total.value = data.total ?? 0;
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="用户列表" description="展示在本店产生过有效订单的客户（脱敏手机号）。">
    <EcrmListPage title="店铺客户">
      <template #filters>
        <el-form inline @submit.prevent="query.page = 1; load()">
          <el-form-item label="关键词"><el-input v-model="query.keyword" placeholder="用户 ID 或昵称" clearable /></el-form-item>
        </el-form>
      </template>
      <template #actions><el-button type="primary" @click="query.page = 1; load()">查询</el-button></template>
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="user_id" label="用户 ID" width="90" />
        <el-table-column prop="nickname" label="昵称" min-width="140" />
        <el-table-column prop="mobile" label="手机号" width="140" />
        <el-table-column prop="order_count" label="订单数" width="90" />
        <el-table-column label="累计消费" width="120"><template #default="{ row }">¥{{ Number(row.total_pay).toFixed(2) }}</template></el-table-column>
        <el-table-column prop="last_order_at" label="最近下单" width="180" />
        <el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="row.status ? 'success' : 'info'">{{ row.status ? '正常' : '停用' }}</el-tag></template></el-table-column>
      </el-table>
      <template #pager><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total,prev,pager,next" @current-change="(page) => { query.page = page; load(); }" /></template>
    </EcrmListPage>
  </Page>
</template>
