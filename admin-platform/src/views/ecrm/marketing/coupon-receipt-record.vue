<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { listCouponReceiptRecords, type CouponReceiptRecord } from '#/api/core/platform-coupon-command';

const rows = ref<CouponReceiptRecord[]>([]);
const total = ref(0);
const loading = ref(false);
const canRead = ref(false);
const query = reactive({ page: 1, limit: 20, user_id: undefined as number | undefined, coupon_id: undefined as number | undefined, status: undefined as CouponReceiptRecord['status'] | undefined });
const statusLabel: Record<CouponReceiptRecord['status'], string> = { unused: '未使用', locked: '已锁定', used: '已使用', expired: '已过期' };
const statusType: Record<CouponReceiptRecord['status'], 'danger' | 'info' | 'success' | 'warning'> = { unused: 'success', locked: 'warning', used: 'info', expired: 'danger' };

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const result = await listCouponReceiptRecords({ ...query });
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}
function reset() {
  Object.assign(query, { page: 1, user_id: undefined, coupon_id: undefined, status: undefined });
  void load();
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.includes('platform') && codes.includes('marketing.coupon.record.read');
  await load();
});
</script>

<template>
  <Page title="优惠券领取记录" description="只读监管用户已领取优惠券及其真实状态；页面不返回账号、手机号或地址，不提供人工核销、解锁、删除或状态修复。">
    <el-alert v-if="!canRead" title="当前账号没有优惠券领取记录查看权限" type="warning" :closable="false" />
    <template v-else>
      <el-card shadow="never"><el-form inline @submit.prevent="query.page = 1; load()"><el-form-item label="用户 ID"><el-input-number v-model="query.user_id" :min="1" controls-position="right" /></el-form-item><el-form-item label="优惠券 ID"><el-input-number v-model="query.coupon_id" :min="1" controls-position="right" /></el-form-item><el-form-item label="状态"><el-select v-model="query.status" clearable class="w-28"><el-option v-for="(label, value) in statusLabel" :key="value" :label="label" :value="value" /></el-select></el-form-item><el-button type="primary" @click="query.page = 1; load()">查询</el-button><el-button @click="reset">重置</el-button></el-form></el-card>
      <el-card class="mt-4" shadow="never"><el-table v-loading="loading" :data="rows"><el-table-column prop="id" label="用户券 ID" width="110" /><el-table-column prop="user_id" label="用户 ID" width="100" /><el-table-column label="优惠券" min-width="220"><template #default="{ row }">{{ row.coupon_name }}（#{{ row.coupon_id }}）</template></el-table-column><el-table-column label="归属" width="110"><template #default="{ row }">{{ row.store_id ? `店铺 #${row.store_id}` : '平台券' }}</template></el-table-column><el-table-column prop="source" label="来源" width="150" /><el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="statusType[row.status]">{{ statusLabel[row.status] }}</el-tag></template></el-table-column><el-table-column label="关联订单 ID" width="140"><template #default="{ row }">{{ row.used_order_id || '—' }}</template></el-table-column><el-table-column prop="obtained_at" label="领取时间" min-width="180" /></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
      <el-alert class="mt-4" type="info" :closable="false" title="“已锁定、已使用、已过期”等状态由下单、支付、取消和券规则状态机写入；平台本页仅监管事实，不能绕过订单状态机修改。" />
    </template>
  </Page>
</template>
