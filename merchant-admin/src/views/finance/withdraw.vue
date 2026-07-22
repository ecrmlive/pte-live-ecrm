<template>
  <a-card :bordered="false" class="page-card">
    <a-alert
      type="info"
      show-icon
      :message="`当前余额 ¥${balance.toFixed(2)}`"
      style="margin-bottom: 16px"
    />
    <a-form layout="inline" class="form" @finish="submit">
      <a-form-item label="提现金额" name="extract_money" :rules="[{ required: true }]">
        <a-input-number v-model:value="form.extract_money" :min="0.01" :precision="2" style="width: 160px" />
      </a-form-item>
      <a-form-item label="账户类型">
        <a-select v-model:value="form.financial_type" style="width: 120px">
          <a-select-option :value="1">银行卡</a-select-option>
          <a-select-option :value="2">微信</a-select-option>
          <a-select-option :value="3">支付宝</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="收款账户" name="financial_account" :rules="[{ required: true }]">
        <a-input v-model:value="form.financial_account" style="width: 220px" placeholder="演示账户信息" />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit" :loading="saving">申请提现</a-button>
      </a-form-item>
    </a-form>

    <a-table
      style="margin-top: 24px"
      row-key="financial_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'money'">¥{{ Number(record.extract_money).toFixed(2) }}</template>
        <template v-else-if="column.key === 'status'">{{ statusText(record.status) }}</template>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { createWithdraw, fetchBalance, fetchWithdraws, type Financial } from '@/api/finance';

const balance = ref(0);
const loading = ref(false);
const saving = ref(false);
const list = ref<Financial[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const form = reactive({
  extract_money: 100,
  financial_type: 1,
  financial_account: '演示银行卡 6222****0001',
});

const columns = [
  { title: 'ID', dataIndex: 'financial_id', width: 80 },
  { title: '单号', dataIndex: 'financial_sn' },
  { title: '金额', key: 'money', width: 120 },
  { title: '状态', key: 'status', width: 120 },
  { title: '拒绝原因', dataIndex: 'refusal' },
];

function statusText(s: number) {
  if (s === 0) return '待审核';
  if (s === 1) return '已通过';
  if (s === -1) return '已拒绝';
  return String(s);
}

async function loadBalance() {
  const { data } = await fetchBalance();
  balance.value = Number(data.mer_money || 0);
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchWithdraws({
      page: pagination.current,
      limit: pagination.pageSize,
    });
    list.value = data.list || [];
    pagination.total = data.total;
  } finally {
    loading.value = false;
  }
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  load();
}

async function submit() {
  saving.value = true;
  try {
    await createWithdraw({ ...form });
    message.success('已提交提现申请');
    await loadBalance();
    pagination.current = 1;
    await load();
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  await loadBalance();
  await load();
});
</script>

<style scoped>
.page-card {
  border-radius: 14px;
}
.form {
  row-gap: 12px;
}
</style>
