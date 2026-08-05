<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { getSvipOrderSummary, listSvipOrders, type SvipOrder, type SvipOrderSummary } from '#/api/core/platform-svip-plan';

const canRead = ref(false);
const rows = ref<SvipOrder[]>([]);
const total = ref(0);
const summary = ref<SvipOrderSummary>();
const loading = ref(false);
const query = reactive({ page: 1, limit: 20, status: undefined as string | undefined });
function orderStatusText(status: SvipOrder['status']) { return ({ pending: '待支付', paid: '已支付', closed: '已关闭' }[status] || status); }
function orderStatusType(status: SvipOrder['status']) { return ({ pending: 'warning', paid: 'success', closed: 'info' }[status] || 'info') as 'info' | 'success' | 'warning'; }
async function load() { if (!canRead.value) return; loading.value = true; try { const [list, stats] = await Promise.all([listSvipOrders(query), getSvipOrderSummary()]); rows.value = list.list || []; total.value = list.total || 0; summary.value = stats; } finally { loading.value = false; } }
onMounted(async () => { const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]); canRead.value = profile.roles.some((role) => role === 'platform' || role === 'operations') && codes.includes('user.svip.record.read'); await load(); });
</script>

<template>
  <Page title="会员记录" description="只读监管 SVIP 购买记录与统计；付款状态由支付回调幂等推进，后台不提供补单、改价或人工置已支付。">
    <el-alert v-if="!canRead" title="当前账号没有会员记录查看权限" type="warning" :closable="false" />
    <template v-else><div class="mb-4 grid grid-cols-2 gap-4 md:grid-cols-5"><el-card shadow="never">总订单：{{ summary?.total || 0 }}</el-card><el-card shadow="never">待支付：{{ summary?.pending || 0 }}</el-card><el-card shadow="never">已支付：{{ summary?.paid || 0 }}</el-card><el-card shadow="never">已关闭：{{ summary?.closed || 0 }}</el-card><el-card shadow="never">已支付金额：¥{{ Number(summary?.paid_amount || 0).toFixed(2) }}</el-card></div><el-card shadow="never"><el-form inline><el-form-item label="状态"><el-select v-model="query.status" clearable><el-option label="待支付" value="pending"/><el-option label="已支付" value="paid"/><el-option label="已关闭" value="closed"/></el-select></el-form-item><el-button type="primary" @click="query.page = 1; load()">查询</el-button></el-form><el-table v-loading="loading" :data="rows"><el-table-column prop="order_no" label="会员订单号" min-width="180"/><el-table-column prop="user_id" label="用户 ID" width="100"/><el-table-column prop="plan_name" label="会员类型" min-width="130"/><el-table-column label="有效期" width="100"><template #default="{ row }">{{ row.plan_type === 'lifetime' ? '永久' : `${row.duration_days} 天` }}</template></el-table-column><el-table-column label="金额" width="100"><template #default="{ row }">¥{{ Number(row.amount).toFixed(2) }}</template></el-table-column><el-table-column label="状态" width="90"><template #default="{ row }"><el-tag :type="orderStatusType(row.status)">{{ orderStatusText(row.status) }}</el-tag></template></el-table-column><el-table-column prop="created_at" label="创建时间" min-width="170"/><el-table-column prop="paid_at" label="支付时间" min-width="170"/></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10,20,50,100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }"/></div></el-card></template>
  </Page>
</template>
