<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { listPlatformInvoices, type PlatformInvoice } from '#/api/core/platform-invoice';
import { EcrmListPage } from '#/components/ecrm';

const rows = ref<PlatformInvoice[]>([]);
const total = ref(0);
const loading = ref(false);
const canRead = ref(false);
const query = reactive({
  limit: 20,
  order_no: '',
  page: 1,
  status: undefined as PlatformInvoice['status'] | undefined,
});

const labels: Record<PlatformInvoice['status'], string> = {
  issued: '已开票',
  rejected: '已拒绝',
  requested: '待开票',
  voided: '已作废',
};

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const result = await listPlatformInvoices({
      limit: query.limit,
      order_no: query.order_no.trim() || undefined,
      page: query.page,
      status: query.status,
    });
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
  Object.assign(query, { page: 1, status: undefined, order_no: '' });
  void load();
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.includes('platform') && permissions.includes('accounts.invoice.read');
  await load();
});
</script>

<template>
  <Page
    title="发票管理"
    description="只读监管订单发票申请与开具事实。税号、邮箱已脱敏；平台不能在此人工开票、改号、改文件、作废或删除。"
  >
    <el-alert
      v-if="!canRead"
      title="当前账号没有发票监管权限"
      type="warning"
      :closable="false"
    />
    <EcrmListPage v-else title="发票列表">
      <template #filters>
        <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
          <el-form-item label="订单号">
            <el-input v-model="query.order_no" maxlength="64" clearable />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="query.status" clearable class="w-28" placeholder="全部">
              <el-option
                v-for="(label, status) in labels"
                :key="status"
                :label="label"
                :value="status"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="search">查询</el-button>
            <el-button @click="reset">重置</el-button>
          </el-form-item>
        </el-form>
      </template>

      <el-table v-loading="loading" :data="rows" row-key="id">
        <el-table-column prop="id" label="发票 ID" width="110" />
        <el-table-column prop="order_no" label="订单号" min-width="210" />
        <el-table-column label="店铺" min-width="180">
          <template #default="{ row }">{{ row.merchant_name }} / {{ row.store_name }}</template>
        </el-table-column>
        <el-table-column prop="title" label="抬头" min-width="200" />
        <el-table-column prop="tax_no_masked" label="税号（脱敏）" min-width="160" />
        <el-table-column prop="email_masked" label="邮箱（脱敏）" min-width="180" />
        <el-table-column label="状态" width="110">
          <template #default="{ row }"><el-tag>{{ labels[row.status] }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="invoice_no" label="发票号" min-width="180" />
        <el-table-column prop="requested_at" label="申请时间" min-width="180" />
      </el-table>

      <template #pager>
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
      </template>
    </EcrmListPage>
  </Page>
</template>
