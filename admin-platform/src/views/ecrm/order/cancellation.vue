<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';

import {
  getPlatformOrderApi,
  listPlatformOrdersApi,
  type PlatformOrder,
} from '#/api/core/platform-trade';

const loading = ref(false);
const detailOpen = ref(false);
const detailLoading = ref(false);
const rows = ref<PlatformOrder[]>([]);
const detail = ref<PlatformOrder>();
const total = ref(0);
const query = reactive({ order_id: '', page: 1, limit: 20 });

function orderStatus(value: number) {
  return ({ '-1': '已退款/取消', 0: '待发货', 1: '待收货', 2: '待评价', 3: '已完成' } as Record<number, string>)[value] || '未知';
}

function payType(value: number) {
  return ({ 0: '余额', 1: '微信', 2: '支付宝', 7: '模拟支付', 8: '积分' } as Record<number, string>)[value] || '—';
}

async function load() {
  if (query.order_id.trim()) {
    const id = Number(query.order_id);
    if (!Number.isInteger(id) || id <= 0) return;
    await openDetail(id);
    return;
  }
  loading.value = true;
  try {
    const data = await listPlatformOrdersApi({ limit: query.limit, page: query.page, status: -1 });
    rows.value = data.list || [];
    total.value = data.total || 0;
  } finally {
    loading.value = false;
  }
}

async function openDetail(id: number) {
  detailOpen.value = true;
  detailLoading.value = true;
  detail.value = undefined;
  try {
    detail.value = await getPlatformOrderApi(id);
  } finally {
    detailLoading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.order_id = '';
  query.page = 1;
  void load();
}

onMounted(() => {
  void load();
});
</script>

<template>
  <Page title="取消/退款订单" description="监管已取消或退款关闭的订单；只读，不直接修改履约或退款状态。">
    <el-card shadow="never">
      <el-form inline @submit.prevent="search">
        <el-form-item label="订单 ID"><el-input v-model="query.order_id" clearable placeholder="精确查询订单 ID" @keyup.enter="search" /></el-form-item>
        <el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button></el-form-item>
      </el-form>
    </el-card>
    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="order_id">
        <el-table-column label="订单 ID" prop="order_id" width="94" />
        <el-table-column label="订单号" min-width="178" prop="order_sn" show-overflow-tooltip />
        <el-table-column label="商户" min-width="140"><template #default="{ row }">{{ row.mer_name || `商户 #${row.mer_id}` }}</template></el-table-column>
        <el-table-column label="商品数" prop="total_num" width="80" />
        <el-table-column label="实付金额" width="110"><template #default="{ row }">¥{{ Number(row.pay_price).toFixed(2) }}</template></el-table-column>
        <el-table-column label="订单状态" width="120"><template #default="{ row }">{{ orderStatus(row.status) }}</template></el-table-column>
        <el-table-column label="下单时间" min-width="170" prop="create_time" />
        <el-table-column fixed="right" label="操作" width="78"><template #default="{ row }"><el-button link type="primary" @click="openDetail(row.order_id)">详情</el-button></template></el-table-column>
      </el-table>
    </el-card>
    <div class="mt-4 flex justify-end">
      <el-pagination
        :current-page="query.page"
        :page-size="query.limit"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        background
        layout="total, sizes, prev, pager, next"
        @current-change="(page) => { query.page = page; load(); }"
        @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }"
      />
    </div>
    <el-dialog v-model="detailOpen" title="订单详情（只读监管）" width="820px" destroy-on-close>
      <el-skeleton :loading="detailLoading" animated :rows="8">
        <template #default>
          <template v-if="detail">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="订单号">{{ detail.order_sn }}</el-descriptions-item>
              <el-descriptions-item label="商户">{{ detail.mer_name || `商户 #${detail.mer_id}` }}</el-descriptions-item>
              <el-descriptions-item label="支付方式">{{ payType(detail.pay_type) }}</el-descriptions-item>
              <el-descriptions-item label="订单状态">{{ orderStatus(detail.status) }}</el-descriptions-item>
              <el-descriptions-item label="支付金额">¥{{ Number(detail.pay_price).toFixed(2) }}</el-descriptions-item>
              <el-descriptions-item label="支付时间">{{ detail.pay_time || '—' }}</el-descriptions-item>
              <el-descriptions-item label="物流信息" :span="2">{{ detail.delivery_name || '未发货' }} {{ detail.delivery_id || '' }}</el-descriptions-item>
              <el-descriptions-item label="收货信息" :span="2">{{ detail.user_phone }} {{ detail.user_address }}</el-descriptions-item>
            </el-descriptions>
            <el-divider content-position="left">商品明细</el-divider>
            <el-table :data="detail.products || []" border>
              <el-table-column label="商品" min-width="240" prop="product_info" show-overflow-tooltip />
              <el-table-column label="单价" width="110"><template #default="{ row }">¥{{ Number(row.product_price).toFixed(2) }}</template></el-table-column>
              <el-table-column label="数量" prop="product_num" width="80" />
              <el-table-column label="小计" width="110"><template #default="{ row }">¥{{ Number(row.total_price).toFixed(2) }}</template></el-table-column>
            </el-table>
          </template>
        </template>
      </el-skeleton>
    </el-dialog>
  </Page>
</template>
