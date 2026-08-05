<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getDistributionSummaryApi,
  listDistributionCommissionsApi,
  listDistributionPromotersApi,
  type CommissionStatus,
  type DistributionCommission,
  type DistributionPromoter,
  type DistributionSummary,
} from '#/api/core/platform-spread';

const activeTab = ref<'commissions' | 'promoters'>('promoters');
const canRead = ref(false);
const loading = ref(false);
const promoters = ref<DistributionPromoter[]>([]);
const commissions = ref<DistributionCommission[]>([]);
const promoterTotal = ref(0);
const commissionTotal = ref(0);
const summary = ref<DistributionSummary>();
const promoterQuery = reactive({ limit: 20, page: 1, status: undefined as 0 | 1 | undefined, user_id: undefined as number | undefined });
const commissionQuery = reactive({ limit: 20, page: 1, status: undefined as CommissionStatus | undefined, user_id: undefined as number | undefined });

const commissionLabels: Record<CommissionStatus, string> = { available: '可结算', pending: '待结算', settled: '已结算', voided: '已作废' };

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const [summaryData, promoterData, commissionData] = await Promise.all([
      getDistributionSummaryApi(),
      listDistributionPromotersApi(promoterQuery),
      listDistributionCommissionsApi(commissionQuery),
    ]);
    summary.value = summaryData;
    promoters.value = promoterData.list || [];
    promoterTotal.value = promoterData.total || 0;
    commissions.value = commissionData.list || [];
    commissionTotal.value = commissionData.total || 0;
  } finally {
    loading.value = false;
  }
}

function searchPromoters() { promoterQuery.page = 1; void load(); }
function searchCommissions() { commissionQuery.page = 1; void load(); }
function resetPromoters() { promoterQuery.status = undefined; promoterQuery.user_id = undefined; promoterQuery.page = 1; void load(); }
function resetCommissions() { commissionQuery.status = undefined; commissionQuery.user_id = undefined; commissionQuery.page = 1; void load(); }

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.some((role) => role === 'platform' || role === 'operations') && permissions.includes('marketing.spread.read');
  await load();
});
</script>

<template>
  <Page title="分销监管" description="监管推广资格、直推关系和佣金状态；本页不修改佣金、提现或用户绑定关系，资金变更必须由业务域状态机处理。">
    <el-alert v-if="!canRead" class="mb-4" title="当前账号没有查看分销监管的权限" type="warning" :closable="false" />
    <template v-else>
      <el-row :gutter="16" class="mb-4"><el-col :md="6" :xs="12"><el-card shadow="never"><div class="text-sm text-gray-500">推广员</div><div class="mt-2 text-lg">{{ summary?.promoter_count || 0 }}（启用 {{ summary?.active_promoter_count || 0 }}）</div></el-card></el-col><el-col :md="6" :xs="12"><el-card shadow="never"><div class="text-sm text-gray-500">待结算佣金</div><div class="mt-2 text-lg">{{ (summary?.pending_commission || 0).toFixed(2) }}</div></el-card></el-col><el-col :md="6" :xs="12"><el-card shadow="never"><div class="text-sm text-gray-500">可结算佣金</div><div class="mt-2 text-lg">{{ (summary?.available_commission || 0).toFixed(2) }}</div></el-card></el-col><el-col :md="6" :xs="12"><el-card shadow="never"><div class="text-sm text-gray-500">已结算佣金</div><div class="mt-2 text-lg">{{ (summary?.settled_commission || 0).toFixed(2) }}</div></el-card></el-col></el-row>
      <el-card shadow="never"><el-tabs v-model="activeTab"><el-tab-pane label="推广员" name="promoters"><el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="searchPromoters"><el-form-item label="用户 ID"><el-input-number v-model="promoterQuery.user_id" :min="1" controls-position="right" /></el-form-item><el-form-item label="状态"><el-select v-model="promoterQuery.status" clearable placeholder="全部"><el-option :value="1" label="启用" /><el-option :value="0" label="停用" /></el-select></el-form-item><el-form-item><el-button type="primary" @click="searchPromoters">查询</el-button><el-button @click="resetPromoters">重置</el-button></el-form-item></el-form><el-table v-loading="loading" :data="promoters" class="mt-3" row-key="user_id"><el-table-column label="用户 ID" prop="user_id" width="110" /><el-table-column label="资格" width="100"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '停用' }}</el-tag></template></el-table-column><el-table-column label="直推用户" prop="direct_user_count" width="110" /><el-table-column label="待结算" width="120"><template #default="{ row }">{{ row.pending_commission.toFixed(2) }}</template></el-table-column><el-table-column label="可结算" width="120"><template #default="{ row }">{{ row.available_commission.toFixed(2) }}</template></el-table-column><el-table-column label="已结算" width="120"><template #default="{ row }">{{ row.settled_commission.toFixed(2) }}</template></el-table-column><el-table-column label="更新时间" min-width="180" prop="updated_at" /></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="promoterQuery.page" :page-size="promoterQuery.limit" :page-sizes="[10, 20, 50, 100]" :total="promoterTotal" background layout="total, sizes, prev, pager, next" @current-change="(page: number) => { promoterQuery.page = page; load(); }" @size-change="(limit: number) => { promoterQuery.limit = limit; promoterQuery.page = 1; load(); }" /></div></el-tab-pane><el-tab-pane label="佣金流水" name="commissions"><el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="searchCommissions"><el-form-item label="用户 ID"><el-input-number v-model="commissionQuery.user_id" :min="1" controls-position="right" /></el-form-item><el-form-item label="佣金状态"><el-select v-model="commissionQuery.status" clearable placeholder="全部"><el-option v-for="(label, value) in commissionLabels" :key="value" :label="label" :value="value" /></el-select></el-form-item><el-form-item><el-button type="primary" @click="searchCommissions">查询</el-button><el-button @click="resetCommissions">重置</el-button></el-form-item></el-form><el-table v-loading="loading" :data="commissions" class="mt-3" row-key="commission_id"><el-table-column label="流水 ID" prop="commission_id" width="110" /><el-table-column label="用户 ID" prop="user_id" width="110" /><el-table-column label="订单 ID" prop="order_id" width="110" /><el-table-column label="佣金金额" width="120"><template #default="{ row }">{{ row.amount.toFixed(2) }}</template></el-table-column><el-table-column label="状态" width="110"><template #default="{ row }"><el-tag>{{ commissionLabels[row.status] }}</el-tag></template></el-table-column><el-table-column label="可结算时间" min-width="180" prop="available_at" /><el-table-column label="创建时间" min-width="180" prop="created_at" /></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="commissionQuery.page" :page-size="commissionQuery.limit" :page-sizes="[10, 20, 50, 100]" :total="commissionTotal" background layout="total, sizes, prev, pager, next" @current-change="(page: number) => { commissionQuery.page = page; load(); }" @size-change="(limit: number) => { commissionQuery.limit = limit; commissionQuery.page = 1; load(); }" /></div></el-tab-pane></el-tabs></el-card>
    </template>
  </Page>
</template>
