<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  deliverMerchantOrderApi,
  getMerchantOrderApi,
  listMerchantOrdersApi,
  verifyMerchantOrderApi,
  type MerchantOrder,
} from '#/api/core/merchant-trade';
import { EcrmFormDialog, EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<MerchantOrder[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const deliveryOpen = ref(false);
const detail = ref<MerchantOrder>();
const deliveryRow = ref<MerchantOrder>();
const delivering = ref(false);
const canDeliver = ref(false);
const canVerifyAction = ref(false);
const query = reactive({
  limit: 20,
  page: 1,
  paid: undefined as number | undefined,
  status: undefined as number | undefined,
});
const deliveryForm = reactive({ delivery_id: '', delivery_name: '', delivery_type: 'express' });

function statusInfo(row: MerchantOrder) {
  if (row.paid !== 1) return { label: '待支付', type: 'info' as const };
  return (
    {
      0: { label: '待发货', type: 'warning' as const },
      1: { label: '待收货', type: 'primary' as const },
      3: { label: '已完成', type: 'success' as const },
    }[row.status] || { label: '未知状态', type: 'info' as const }
  );
}

function payType(value: number) {
  return { 0: '余额', 1: '微信', 2: '支付宝', 7: '模拟支付', 8: '积分' }[value] || '未知';
}

function rowCanDeliver(row: MerchantOrder) {
  return canDeliver.value && row.paid === 1 && row.status === 0;
}

function canVerify(row: MerchantOrder) {
  return canVerifyAction.value && Boolean(row.can_verify);
}

async function load() {
  loading.value = true;
  try {
    const result = await listMerchantOrdersApi(query);
    rows.value = result.list;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  Object.assign(query, { page: 1, paid: undefined, status: undefined });
  void load();
}

async function openDetail(row: MerchantOrder) {
  detail.value = await getMerchantOrderApi(row.order_id);
  detailOpen.value = true;
}

function openDelivery(row: MerchantOrder) {
  deliveryRow.value = row;
  Object.assign(deliveryForm, { delivery_id: '', delivery_name: '', delivery_type: 'express' });
  deliveryOpen.value = true;
}

async function deliver() {
  if (!deliveryRow.value) return;
  if (!deliveryForm.delivery_name.trim() || !deliveryForm.delivery_id.trim()) {
    ElMessage.warning('请填写物流公司和运单号');
    return;
  }
  delivering.value = true;
  try {
    await deliverMerchantOrderApi(deliveryRow.value.order_id, {
      delivery_id: deliveryForm.delivery_id.trim(),
      delivery_name: deliveryForm.delivery_name.trim(),
      delivery_type: deliveryForm.delivery_type,
    });
    deliveryOpen.value = false;
    ElMessage.success('订单已发货');
    await load();
  } finally {
    delivering.value = false;
  }
}

async function verify(row: MerchantOrder) {
  try {
    await ElMessageBox.confirm(
      '确认完成该订单核销？核销后订单将变为已完成，操作不可撤销。',
      '核销确认',
      { confirmButtonText: '确认核销', cancelButtonText: '取消', type: 'warning' },
    );
    await verifyMerchantOrderApi(row.order_id);
    ElMessage.success('订单已核销');
    await load();
  } catch {
    // 用户取消或接口错误时由统一请求层处理。
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canDeliver.value = codes.includes('order.deliver');
  canVerifyAction.value = codes.includes('order.verify.action');
  await load();
});
</script>

<template>
  <EcrmListPage title="订单管理" description="仅展示当前店铺订单；发货按钮受 order.deliver 权限与订单状态共同约束。">
    <template #filters>
      <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
        <el-form-item label="支付状态">
          <el-select v-model="query.paid" clearable class="w-36" placeholder="全部">
            <el-option label="待支付" :value="0" />
            <el-option label="已支付" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="query.status" clearable class="w-36" placeholder="全部">
            <el-option label="待发货" :value="0" />
            <el-option label="待收货" :value="1" />
            <el-option label="已完成" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">查询</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </template>

    <el-table v-loading="loading" :data="rows" row-key="order_id">
      <el-table-column label="订单号" min-width="180" prop="order_sn" />
      <el-table-column label="用户" min-width="130">
        <template #default="{ row }">
          <div>{{ row.real_name || '—' }}</div>
          <div class="text-xs text-muted-foreground">{{ row.user_phone || '—' }}</div>
        </template>
      </el-table-column>
      <el-table-column label="商品数" width="88" prop="total_num" />
      <el-table-column label="实付金额" width="116">
        <template #default="{ row }">¥{{ Number(row.pay_price).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="支付方式" width="104">
        <template #default="{ row }">{{ payType(row.pay_type) }}</template>
      </el-table-column>
      <el-table-column label="订单状态" width="104">
        <template #default="{ row }">
          <el-tag :type="statusInfo(row).type">{{ statusInfo(row).label }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="下单时间" min-width="170" prop="create_time" />
      <el-table-column fixed="right" label="操作" width="184">
        <template #default="{ row }">
          <el-button link type="primary" @click="openDetail(row)">详情</el-button>
          <el-button v-if="rowCanDeliver(row)" link type="success" @click="openDelivery(row)">发货</el-button>
          <el-button v-if="canVerify(row)" link type="warning" @click="verify(row)">核销</el-button>
        </template>
      </el-table-column>
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
        <el-descriptions-item label="物流公司">{{ detail.delivery_name || '—' }}</el-descriptions-item>
        <el-descriptions-item label="运单号">{{ detail.delivery_id || '—' }}</el-descriptions-item>
        <el-descriptions-item :span="2" label="用户备注">{{ detail.mark || '—' }}</el-descriptions-item>
      </el-descriptions>
      <div class="mb-3 mt-6 text-base font-medium">商品明细</div>
      <el-table :data="detail.products || []" border>
        <el-table-column label="商品 ID" prop="product_id" width="92" />
        <el-table-column label="商品信息" min-width="180" prop="product_info" show-overflow-tooltip />
        <el-table-column label="规格" min-width="120" prop="product_sku" />
        <el-table-column label="单价" width="104">
          <template #default="{ row }">¥{{ Number(row.product_price).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="数量" prop="product_num" width="76" />
        <el-table-column label="小计" width="104">
          <template #default="{ row }">¥{{ Number(row.total_price).toFixed(2) }}</template>
        </el-table-column>
      </el-table>
    </template>
  </el-drawer>

  <EcrmFormDialog v-model="deliveryOpen" title="订单发货" width="480px">
    <el-form label-width="84px">
      <el-form-item label="物流类型">
        <el-select v-model="deliveryForm.delivery_type" class="w-full">
          <el-option label="快递" value="express" />
          <el-option label="同城配送" value="local" />
        </el-select>
      </el-form-item>
      <el-form-item label="物流公司" required>
        <el-input v-model="deliveryForm.delivery_name" placeholder="例如：顺丰速运" />
      </el-form-item>
      <el-form-item label="运单号" required>
        <el-input v-model="deliveryForm.delivery_id" placeholder="填写物流运单号" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="deliveryOpen = false">取消</el-button>
      <el-button :loading="delivering" type="primary" @click="deliver">确认发货</el-button>
    </template>
  </EcrmFormDialog>
</template>
