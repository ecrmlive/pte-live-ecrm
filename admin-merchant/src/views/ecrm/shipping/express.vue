<script setup lang="ts">
import { reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { listMerchantExpressApi, type MerchantExpressRow } from '#/api/core/merchant-logistics';

const loading = ref(false);
const rows = ref<MerchantExpressRow[]>([]);
const total = ref(0);
const query = reactive({ limit: 50, page: 1 });

async function load() {
  loading.value = true;
  try {
    const result = await listMerchantExpressApi(query);
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

void load();
</script>

<template>
  <Page title="物流公司" description="只读查看平台维护的快递公司列表，发货时可选择对应物流编码。">
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="express_id">
        <el-table-column label="ID" prop="express_id" width="90" />
        <el-table-column label="名称" min-width="160" prop="name" />
        <el-table-column label="编码" min-width="140" prop="code" />
        <el-table-column label="排序" prop="sort" width="90" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_show === 1 ? 'success' : 'info'">{{ row.is_show === 1 ? '展示' : '隐藏' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination
          :current-page="query.page"
          :page-size="query.limit"
          :page-sizes="[20, 50, 100]"
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
