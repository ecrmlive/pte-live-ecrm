<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getPlatformUserAssetSummaryApi,
  listPlatformUserAssetsApi,
  type UserAssetLedgerRow,
  type UserAssetSummary,
  type UserAssetType,
} from '#/api/core/platform-user-assets';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<UserAssetLedgerRow[]>([]);
const total = ref(0);
const summary = ref<UserAssetSummary[]>([]);
const canRead = ref(false);
const query = reactive({
  asset_type: undefined as UserAssetType | undefined,
  limit: 20,
  page: 1,
  user_id: undefined as number | undefined,
});

const summaryByType = computed(() => new Map(summary.value.map((item) => [item.asset_type, item])));

const assetLabels: Record<UserAssetType, string> = {
  balance: '余额',
  commission: '佣金',
  points: '积分',
};

function formatAmount(row: UserAssetLedgerRow) {
  const prefix = row.amount > 0 ? '+' : '';
  return `${prefix}${row.amount.toFixed(2)}`;
}

function summaryText(type: UserAssetType) {
  const item = summaryByType.value.get(type);
  if (!item) return '暂无流水';
  return `入账 ${item.income.toFixed(2)} · 支出 ${item.expense.toFixed(2)} · ${item.count} 笔`;
}

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const [pageData, summaryData] = await Promise.all([
      listPlatformUserAssetsApi(query),
      getPlatformUserAssetSummaryApi(),
    ]);
    rows.value = pageData.list || [];
    total.value = pageData.total || 0;
    summary.value = summaryData.list || [];
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.asset_type = undefined;
  query.user_id = undefined;
  query.page = 1;
  void load();
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.includes('platform') && permissions.includes('accounts.user_assets.read');
  await load();
});
</script>

<template>
  <Page
    title="用户资产流水"
    description="仅平台角色可查看业务库不可变资产流水；页面不展示用户姓名、手机号、收款账户或外部支付凭据。"
  >
    <el-alert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有查看用户资产流水的权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <el-row :gutter="16" class="mb-4">
        <el-col
          v-for="type in (['balance', 'points', 'commission'] as UserAssetType[])"
          :key="type"
          :md="8"
          :xs="24"
        >
          <el-card shadow="never">
            <div class="text-sm text-gray-500">{{ assetLabels[type] }}</div>
            <div class="mt-2 text-sm">{{ summaryText(type) }}</div>
          </el-card>
        </el-col>
      </el-row>

      <EcrmListPage title="资产流水">
        <template #filters>
          <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
            <el-form-item label="用户 ID">
              <el-input-number v-model="query.user_id" :min="1" controls-position="right" />
            </el-form-item>
            <el-form-item label="资产类型">
              <el-select v-model="query.asset_type" clearable class="w-32" placeholder="全部">
                <el-option
                  v-for="(label, type) in assetLabels"
                  :key="type"
                  :label="label"
                  :value="type"
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
          <el-table-column label="流水 ID" prop="id" width="100" />
          <el-table-column label="用户 ID" prop="user_id" width="100" />
          <el-table-column label="资产类型" width="100">
            <template #default="{ row }">{{ assetLabels[row.asset_type] }}</template>
          </el-table-column>
          <el-table-column label="变动金额" width="130">
            <template #default="{ row }">
              <span :class="row.amount < 0 ? 'text-red-500' : 'text-green-600'">
                {{ formatAmount(row) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="业务来源" min-width="150" prop="reference_type" />
          <el-table-column label="业务引用" min-width="160" prop="reference_id" />
          <el-table-column label="创建时间" min-width="180" prop="created_at" />
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
    </template>
  </Page>
</template>
