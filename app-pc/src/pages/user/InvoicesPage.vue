<script setup lang="ts">
import { onMounted, ref } from "vue";
import AccountFrame from "@/components/AccountFrame.vue";
import { fetchInvoices, type Invoice } from "@/api/trade";
import { ApiError } from "@/utils/request";

const list = ref<Invoice[]>([]);
const hint = ref("");
const loading = ref(false);

function invoiceType(value: number) { return value === 2 ? "增值税专用发票" : "增值税电子普通发票"; }
function headerType(value: number) { return value === 2 ? "企业" : "个人"; }
function statusText(value: number) { return value === 1 ? "已开具" : value === -1 ? "已驳回" : "处理中"; }

async function load() {
  loading.value = true;
  try {
    const res = await fetchInvoices();
    list.value = res.list || [];
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "发票加载失败";
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <AccountFrame>
    <template #crumb><span>›</span> 我的发票</template>
    <header class="invoice-header"><h1>我的发票</h1><button type="button" @click="load">刷新</button></header>
    <p v-if="hint" class="hint">{{ hint }}</p>
    <p v-else-if="loading" class="empty">正在加载发票…</p>
    <div v-else-if="!list.length" class="empty">暂无发票申请记录。支付完成后可在订单页申请开票。</div>
    <div v-else class="invoice-list">
      <article v-for="item in list" :key="item.invoice_id" class="invoice-card">
        <header><b>{{ invoiceType(item.invoice_type) }}</b><span :class="{ rejected: item.status === -1, issued: item.status === 1 }">{{ statusText(item.status) }}</span></header>
        <dl>
          <div><dt>发票抬头</dt><dd>{{ item.header }}（{{ headerType(item.header_type) }}）</dd></div>
          <div><dt>订单编号</dt><dd>#{{ item.order_id }}</dd></div>
          <div><dt>收票邮箱</dt><dd>{{ item.email }}</dd></div>
          <div v-if="item.tax_no"><dt>纳税人识别号</dt><dd>{{ item.tax_no }}</dd></div>
        </dl>
        <p v-if="item.mark" class="mark">处理说明：{{ item.mark }}</p>
      </article>
    </div>
  </AccountFrame>
</template>

<style scoped>
.invoice-header { display: flex; align-items: center; justify-content: space-between; padding-bottom: 19px; border-bottom: 1px solid #eee; }.invoice-header h1 { margin: 0; font-size: 20px; }.invoice-header button { border: 0; background: transparent; color: #777; cursor: pointer; }.invoice-header button:hover { color: #f13728; }.hint { color: #d9362b; }.empty { padding: 76px 0; color: #aaa; text-align: center; }.invoice-list { display: grid; gap: 16px; padding-top: 24px; }.invoice-card { padding: 19px 22px; border: 1px solid #eee; }.invoice-card > header { display: flex; justify-content: space-between; padding-bottom: 14px; border-bottom: 1px dashed #eee; color: #444; }.invoice-card > header span { color: #e79a20; }.invoice-card > header .issued { color: #35a26d; }.invoice-card > header .rejected { color: #d9362b; }.invoice-card dl { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 14px 32px; margin: 17px 0 0; }.invoice-card dl div { display: flex; gap: 11px; min-width: 0; }.invoice-card dt { flex: none; color: #999; }.invoice-card dd { overflow: hidden; margin: 0; color: #555; text-overflow: ellipsis; white-space: nowrap; }.mark { margin: 15px 0 0; color: #d9362b; font-size: 13px; }@media (max-width: 650px) { .invoice-card dl { grid-template-columns: 1fr; } }
</style>
