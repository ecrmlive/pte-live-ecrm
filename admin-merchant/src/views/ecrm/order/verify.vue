<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import {
  getMerchantOrderApi,
  listMerchantOrdersApi,
  verifyMerchantOrderApi,
  type MerchantOrder,
} from '#/api/core/merchant-trade';

type VerifyTab = 'pending' | 'verified';

const loading = ref(false);
const rows = ref<MerchantOrder[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const detail = ref<MerchantOrder>();
const activeTab = ref<VerifyTab>('pending');
const query = reactive({ limit: 20, page: 1, paid: 1 as number | undefined, status: undefined as number | undefined });

function statusInfo(row: MerchantOrder) {
  if (row.paid !== 1) return { label: '待支付', type: 'info' as const };
  return {
    0: { label: '待核销', type: 'warning' as const },
    1: { label: '待收货', type: 'primary' as const },
    3: { label: '已核销', type: 'success' as const },
  }[row.status] || { label: '处理中', type: 'info' as const };
}

function payType(value: number) {
  return ({ 0: '余额', 1: '微信', 2: '支付宝', 7: '模拟支付', 8: '积分' }[value] || '未知');
}

function canVerify(row: MerchantOrder) {
  return activeTab.value === 'pending' && row.paid === 1 && row.status !== 3;
}

function syncQueryFromTab() {
  query.paid = 1;
  query.status = activeTab.value === 'verified' ? 3 : undefined;
  query.page = 1;
}

async function load() {
  loading.value = true;
  try {
    const result = await listMerchantOrdersApi(query);
    rows.value = activeTab.value === 'pending'
      ? (result.list || []).filter((row) => row.status !== 3)
      : (result.list || []);
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

async function openDetail(row: MerchantOrder) {
  detail.value = await getMerchantOrderApi(row.order_id);
  detailOpen.value = true;
}

async function verify(row: MerchantOrder) {
  try {
    await ElMessageBox.confirm(
      `确认核销订单 ${row.order_sn}？核销后订单将标记为已完成，操作不可撤销。`,
      '核销确认',
      { confirmButtonText: '确认核销', cancelButtonText: '取消', type: 'warning' },
    );
    await verifyMerchantOrderApi(row.order_id, row.verify_code);
    ElMessage.success('订单已核销');
    await load();
  } catch {
    // 用户取消或接口错误由统一请求层处理。
  }
}

watch(activeTab, () => {
  syncQueryFromTab();
  void load();
});

onMounted(() => {
  syncQueryFromTab();
  void load();
});
</script>

<template>
  <Page title="核销记录" description="查看待核销与已核销订单；核销操作通过独立接口提交，列表数据来自本店订单。">
    <el-card shadow="never">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="待核销" name="pending" />
        <el-tab-pane label="已核销" name="verified" />
      </el-tabs>
    </el-card>

    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="order_id">
        <el-table-column label="订单号" min-width="180" prop="order_sn" />
        <el-table-column label="用户" min-width="130">
          <template #default="{ row }">
            <div>{{ row.real_name || '—' }}</div>
            <div class="text-xs text-muted-foreground">{{ row.user_phone || '—' }}</div>
          </template>
        </el-table-column>
        <el-table-column label="商品数" prop="total_num" width="88" />
        <el-table-column label="实付金额" width="116">
          <template #default="{ row }">¥{{ Number(row.pay_price).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="支付方式" width="104">
          <template #default="{ row }">{{ payType(row.pay_type) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="104">
          <template #default="{ row }">
            <el-tag :type="statusInfo(row).type">{{ statusInfo(row).label }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="下单时间" min-width="170" prop="create_time" />
        <el-table-column fixed="right" label="操作" width="148">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
            <el-button v-if="canVerify(row)" link type="warning" @click="verify(row)">核销</el-button>
          </template>
        </el-table-column>
      </el-table>
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
    </el-card>

    <el-drawer v-model="detailOpen" :with-header="false" size="680px">
      <template v-if="detail">
        <div class="mb-5 text-lg font-medium">订单详情</div>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ detail.order_sn }}</el-descriptions-item>
          <el-descriptions-item label="订单状态">
            <el-tag :type="statusInfo(detail).type">{{ statusInfo(detail).label }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="收货人">{{ detail.real_name || '—' }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ detail.user_phone || '—' }}</el-descriptions-item>
          <el-descriptions-item :span="2" label="收货地址">{{ detail.user_address || '—' }}</el-descriptions-item>
          <el-descriptions-item label="支付方式">{{ payType(detail.pay_type) }}</el-descriptions-item>
          <el-descriptions-item label="实付金额">¥{{ Number(detail.pay_price).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item v-if="detail.verify_code" label="核销码">{{ detail.verify_code }}</el-descriptions-item>
        </el-descriptions>
        <div class="mb-3 mt-6 text-base font-medium">商品明细</div>
        <el-table :data="detail.products || []" border>
          <el-table-column label="商品 ID" prop="product_id" width="92" />
          <el-table-column label="商品信息" min-width="180" prop="product_info" show-overflow-tooltip />
          <el-table-column label="数量" prop="product_num" width="76" />
          <el-table-column label="小计" width="104">
            <template #default="{ row }">¥{{ Number(row.total_price).toFixed(2) }}</template>
          </el-table-column>
        </el-table>
      </template>
    </el-drawer>
  </Page>
</template>
