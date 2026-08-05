<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformTransferRecordsApi,
  type MerchantSettlementRow,
  type TransferRecordStatus,
} from '#/api/core/platform-merchant-settlement';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<MerchantSettlementRow[]>([]);
const total = ref(0);
const canRead = ref(false);
const query = reactive({
  limit: 20,
  merchant_id: undefined as number | undefined,
  page: 1,
  status: undefined as TransferRecordStatus | undefined,
});

const statusLabels: Record<TransferRecordStatus, string> = {
  approved: '待登记打款',
  paid: '已打款',
  rejected: '已拒绝',
};

function formatTime(value?: string) {
  if (!value) return '—';
  return String(value).replace('T', ' ').slice(0, 19);
}

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const page = await listPlatformTransferRecordsApi(query);
    rows.value = page.list || [];
    total.value = page.total || 0;
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.merchant_id = undefined;
  query.status = undefined;
  query.page = 1;
  void load();
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value =
    profile.roles.some((role) => role === 'platform' || role === 'region') &&
    permissions.includes('accounts.merchant_settlement.read');
  await load();
});
</script>

<template>
  <Page
    title="转账记录"
    description="只读监管店铺结算打款投影（审核通过 / 已打款 / 拒绝）；数据来自 qixi_crm_a_merchant_settlement_view，不回显收款账户或外部支付凭据。"
  >
    <el-alert
      class="mb-4"
      title="本页不发起转账、不保存打款密钥；登记内部凭证请在「店铺结算监管」完成。"
      type="warning"
      :closable="false"
    />
    <el-alert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有转账记录（结算投影）查看权限"
      type="warning"
      :closable="false"
    />
    <EcrmListPage v-else title="转账记录">
      <template #filters>
        <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
          <el-form-item label="商户 ID">
            <el-input-number v-model="query.merchant_id" :min="1" controls-position="right" />
          </el-form-item>
          <el-form-item label="打款状态">
            <el-select v-model="query.status" clearable class="w-36" placeholder="打款链路">
              <el-option
                v-for="(label, status) in statusLabels"
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

      <el-table v-loading="loading" :data="rows" row-key="settlement_id">
        <el-table-column label="结算 ID" prop="settlement_id" width="110" />
        <el-table-column label="商户" min-width="180">
          <template #default="{ row }">{{ row.merchant_name }}（{{ row.merchant_id }}）</template>
        </el-table-column>
        <el-table-column label="店铺 ID" prop="store_id" width="100" />
        <el-table-column label="结算周期" min-width="280">
          <template #default="{ row }">
            {{ formatTime(row.period_start) }} 至 {{ formatTime(row.period_end) }}
          </template>
        </el-table-column>
        <el-table-column label="金额" width="120">
          <template #default="{ row }">¥{{ Number(row.amount).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag>{{ statusLabels[row.status as TransferRecordStatus] || row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="投影更新时间" min-width="180">
          <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!loading && !rows.length" description="暂无打款链路记录" />

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
