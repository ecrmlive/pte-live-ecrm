<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import {
  listMerchantCouponRecordsApi,
  listMerchantCouponSendsApi,
  type MerchantCouponRecord,
  type MerchantCouponSend,
} from '#/api/core/merchant-promotion';

const active = ref('records');
const recordLoading = ref(false);
const sendLoading = ref(false);
const records = ref<MerchantCouponRecord[]>([]);
const sends = ref<MerchantCouponSend[]>([]);
const recordTotal = ref(0);
const sendTotal = ref(0);
const recordQuery = reactive({ limit: 20, page: 1 });
const sendQuery = reactive({ limit: 20, page: 1 });

function statusInfo(status: number) {
  return ({ 0: { label: '未使用', type: 'warning' }, 1: { label: '已使用', type: 'success' }, 2: { label: '已过期', type: 'info' } }[status] || { label: '未知', type: 'info' }) as { label: string; type: 'info' | 'success' | 'warning' };
}

async function loadRecords() {
  recordLoading.value = true;
  try {
    const result = await listMerchantCouponRecordsApi(recordQuery);
    records.value = result.list;
    recordTotal.value = result.total;
  } finally { recordLoading.value = false; }
}

async function loadSends() {
  sendLoading.value = true;
  try {
    const result = await listMerchantCouponSendsApi(sendQuery);
    sends.value = result.list;
    sendTotal.value = result.total;
  } finally { sendLoading.value = false; }
}

onMounted(() => void Promise.all([loadRecords(), loadSends()]));
</script>

<template>
  <Page title="优惠券记录" description="记录仅限当前商户的店铺券；用户领取、后台定向发送和实际使用状态均按商户隔离。">
    <el-tabs v-model="active">
      <el-tab-pane label="领取与使用记录" name="records"><el-card shadow="never"><el-table v-loading="recordLoading" :data="records" row-key="coupon_user_id"><el-table-column label="记录 ID" prop="coupon_user_id" width="96" /><el-table-column label="优惠券" min-width="150" prop="coupon_title" /><el-table-column label="用户 ID" prop="uid" width="100" /><el-table-column label="面额" width="106"><template #default="{ row }">¥{{ Number(row.coupon_price).toFixed(2) }}</template></el-table-column><el-table-column label="获取方式" width="100"><template #default="{ row }">{{ row.type === 'send' ? '后台发送' : '用户领取' }}</template></el-table-column><el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="statusInfo(row.status).type">{{ statusInfo(row.status).label }}</el-tag></template></el-table-column><el-table-column label="领取时间" min-width="170" prop="create_time" /><el-table-column label="使用时间" min-width="170"><template #default="{ row }">{{ row.use_time || '—' }}</template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="recordQuery.page" :page-size="recordQuery.limit" :total="recordTotal" background layout="total, prev, pager, next" @current-change="(page) => { recordQuery.page = page; loadRecords(); }" /></div></el-card></el-tab-pane>
      <el-tab-pane label="发送记录" name="sends"><el-card shadow="never"><el-table v-loading="sendLoading" :data="sends" row-key="coupon_send_id"><el-table-column label="批次 ID" prop="coupon_send_id" width="96" /><el-table-column label="优惠券 ID" prop="coupon_id" width="104" /><el-table-column label="发送数量" prop="coupon_num" width="104" /><el-table-column label="发送说明" min-width="220" prop="mark" show-overflow-tooltip /><el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="row.status === 1 ? 'success' : 'warning'">{{ row.status === 1 ? '已完成' : '发送中' }}</el-tag></template></el-table-column><el-table-column label="发送时间" min-width="180" prop="create_time" /></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="sendQuery.page" :page-size="sendQuery.limit" :total="sendTotal" background layout="total, prev, pager, next" @current-change="(page) => { sendQuery.page = page; loadSends(); }" /></div></el-card></el-tab-pane>
    </el-tabs>
  </Page>
</template>
