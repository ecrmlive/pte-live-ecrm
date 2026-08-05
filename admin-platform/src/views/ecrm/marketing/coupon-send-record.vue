<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import { listCouponCommands, type CouponCommandRecord } from '#/api/core/platform-coupon-command';

const rows = ref<CouponCommandRecord[]>([]);
const total = ref(0);
const loading = ref(false);
const canRead = ref(false);
const query = reactive({ page: 1, limit: 20, user_id: undefined as number | undefined, coupon_id: undefined as number | undefined, action: undefined as 'issue' | 'revoke' | undefined });

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const result = await listCouponCommands({ ...query });
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}
function reset() {
  Object.assign(query, { page: 1, user_id: undefined, coupon_id: undefined, action: undefined });
  void load();
}
function actionLabel(action: CouponCommandRecord['action']) { return action === 'issue' ? '发放' : '撤销'; }
function stateText(value: string) { return value || '—'; }

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.includes('platform') && codes.includes('marketing.coupon.send.read');
  await load();
});
</script>

<template>
  <Page title="优惠券发送记录" description="只读展示平台人工发券与撤销的不可变审计。订单锁券、核销和支付回调不在本页处理，且不会显示账号、手机号或幂等键。">
    <el-alert v-if="!canRead" title="当前账号没有优惠券发送记录查看权限" type="warning" :closable="false" />
    <template v-else>
      <el-card shadow="never"><el-form inline @submit.prevent="query.page = 1; load()"><el-form-item label="用户 ID"><el-input-number v-model="query.user_id" :min="1" controls-position="right" /></el-form-item><el-form-item label="优惠券 ID"><el-input-number v-model="query.coupon_id" :min="1" controls-position="right" /></el-form-item><el-form-item label="操作"><el-select v-model="query.action" clearable class="w-28"><el-option label="发放" value="issue" /><el-option label="撤销" value="revoke" /></el-select></el-form-item><el-button type="primary" @click="query.page = 1; load()">查询</el-button><el-button @click="reset">重置</el-button></el-form></el-card>
      <el-card class="mt-4" shadow="never"><el-table v-loading="loading" :data="rows"><el-table-column prop="id" label="审计 ID" width="100" /><el-table-column prop="user_id" label="用户 ID" width="100" /><el-table-column label="优惠券" min-width="220"><template #default="{ row }">{{ row.coupon_name }}（#{{ row.coupon_id }}）</template></el-table-column><el-table-column label="归属" width="110"><template #default="{ row }">{{ row.store_id ? `店铺 #${row.store_id}` : '平台券' }}</template></el-table-column><el-table-column label="操作" width="90"><template #default="{ row }"><el-tag :type="row.action === 'issue' ? 'success' : 'warning'">{{ actionLabel(row.action) }}</el-tag></template></el-table-column><el-table-column label="状态变化" min-width="150"><template #default="{ row }">{{ stateText(row.from_status) }} → {{ stateText(row.to_status) }}</template></el-table-column><el-table-column prop="reason" label="操作原因" min-width="240" show-overflow-tooltip /><el-table-column prop="operator_admin_id" label="操作管理员 ID" width="140" /><el-table-column prop="created_at" label="操作时间" min-width="180" /></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
      <el-alert class="mt-4" type="info" :closable="false" title="审计记录不提供编辑、删除、补写或重放。若模板后来被删除，记录仍保留并标记为“已删除优惠券模板”。" />
    </template>
  </Page>
</template>
