<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';
import { listMerchantProductsApi, type MerchantProduct } from '#/api/core/merchant-catalog';
import { createProxyOrderApi } from '#/api/core/merchant-order-proxy';
import { listStoreCustomersApi, type StoreCustomer } from '#/api/core/merchant-store-customer';

const customers = ref<StoreCustomer[]>([]);
const products = ref<MerchantProduct[]>([]);
const submitting = ref(false);
const lastResult = ref<{ order_sn: string; pay_amount: number; replayed: boolean } | null>(null);
const form = reactive({
  user_id: undefined as number | undefined,
  product_id: undefined as number | undefined,
  quantity: 1,
  remark: '',
});

async function loadOptions() {
  const [customerPage, productPage] = await Promise.all([
    listStoreCustomersApi({ page: 1, limit: 100 }),
    listMerchantProductsApi({ page: 1, limit: 100, status: 1 }),
  ]);
  customers.value = customerPage.list ?? [];
  products.value = productPage.list ?? [];
}

async function submit() {
  if (!form.user_id || !form.product_id || form.quantity < 1) {
    ElMessage.warning('请选择客户、商品并填写数量');
    return;
  }
  submitting.value = true;
  lastResult.value = null;
  try {
    const result = await createProxyOrderApi({
      user_id: form.user_id,
      product_id: form.product_id,
      quantity: form.quantity,
      remark: form.remark.trim(),
    });
    lastResult.value = {
      order_sn: result.order_sn,
      pay_amount: result.pay_amount,
      replayed: result.replayed,
    };
    ElMessage.success(result.replayed ? '已按幂等键返回既有订单' : '代客下单成功，待客户支付');
  } finally {
    submitting.value = false;
  }
}

onMounted(() => void loadOptions().catch(() => {}));
</script>

<template>
  <Page title="代客下单" description="为在本店有过有效订单的客户创建待支付订单（单规格普通商品）；创建后写入库存预留命令。">
    <el-card shadow="never">
      <el-form label-width="96px" class="max-w-xl">
        <el-form-item label="客户" required>
          <el-select v-model="form.user_id" filterable class="w-full" placeholder="选择店铺客户">
            <el-option
              v-for="item in customers"
              :key="item.user_id"
              :value="item.user_id"
              :label="`${item.nickname || '用户'} (#${item.user_id}) ${item.mobile}`"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="商品" required>
          <el-select v-model="form.product_id" filterable class="w-full" placeholder="选择已上架单规格商品">
            <el-option
              v-for="item in products"
              :key="item.product_id"
              :value="item.product_id"
              :label="`${item.store_name} (#${item.product_id}) ¥${item.price}`"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="数量" required>
          <el-input-number v-model="form.quantity" :min="1" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" maxlength="200" show-word-limit />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submit">提交代客订单</el-button>
        </el-form-item>
      </el-form>
      <el-alert
        v-if="lastResult"
        class="mt-4 max-w-xl"
        type="success"
        :closable="false"
        :title="`订单号 ${lastResult.order_sn} · 应付 ¥${lastResult.pay_amount}${lastResult.replayed ? '（幂等重放）' : ''}`"
      />
    </el-card>
  </Page>
</template>
