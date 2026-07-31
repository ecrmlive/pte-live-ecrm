<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import AccountFrame from "@/components/AccountFrame.vue";
import { fetchBalance, type BalanceLedger } from "@/api/trade";
import { ApiError } from "@/utils/request";

const summary = ref({ balance: 0, total_income: 0, total_expense: 0 });
const list = ref<BalanceLedger[]>([]);
const hint = ref("");
const billTab = ref<"all" | "income" | "expense">("all");
async function load() { try { const data = await fetchBalance(); summary.value = data.summary; list.value = data.list || []; hint.value = ""; } catch (e) { hint.value = e instanceof ApiError ? e.message : "余额加载失败"; } }
function label(item: BalanceLedger) { return item.reference_type === "recharge" ? "余额充值" : item.reference_type === "order_pay" ? "订单支付" : item.reference_type || "余额变动"; }
const displayedList = computed(() => list.value.filter((item) => billTab.value === "all" || (billTab.value === "income" ? item.amount > 0 : item.amount < 0)));
onMounted(() => void load());
</script>

<template>
  <AccountFrame>
    <template #crumb><span>›</span> 我的余额</template>
    <header class="balance-header"><h1>我的余额</h1><button type="button" @click="load">刷新</button></header>
    <div class="metrics"><div><span>总资产(元)</span><strong>{{ Number(summary.balance).toFixed(2) }}</strong></div><div><span>累计充值(元)</span><strong>{{ Number(summary.total_income).toFixed(2) }}</strong></div><div><span>累计消费(元)</span><strong>{{ Number(summary.total_expense).toFixed(2) }}</strong></div></div>
    <p v-if="hint" class="hint">{{ hint }}</p>
    <section class="ledger"><nav><button :class="{ active: billTab === 'all' }" type="button" @click="billTab = 'all'">账单明细</button><button :class="{ active: billTab === 'income' }" type="button" @click="billTab = 'income'">收入明细</button><button :class="{ active: billTab === 'expense' }" type="button" @click="billTab = 'expense'">支出明细</button></nav><div v-if="!displayedList.length" class="empty">暂无余额变动记录</div><article v-for="item in displayedList" :key="item.id"><div><b>{{ label(item) }}</b><span>{{ item.created_at }}</span></div><strong :class="{ income: item.amount > 0 }">{{ item.amount > 0 ? "+" : "" }}¥{{ Number(item.amount).toFixed(2) }}</strong></article></section>
  </AccountFrame>
</template>

<style scoped>
.balance-header { display: flex; align-items: center; justify-content: space-between; padding-bottom: 19px; border-bottom: 1px solid #eee; }.balance-header h1 { margin: 0; font-size: 20px; }.balance-header button { border: 0; background: transparent; color: #777; cursor: pointer; }.metrics { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; padding: 48px 0 40px; border-bottom: 1px solid #eee; }.metrics div { display: grid; gap: 16px; }.metrics span { color: #999; font-size: 14px; }.metrics strong { color: #333; font-size: 32px; font-weight: 400; }.hint { color: #d9362b; }.ledger nav { display: flex; gap: 44px; margin-top: 26px; border-bottom: 1px solid #eee; }.ledger nav button { position: relative; border: 0; padding: 0 10px 16px; color: #999; background: transparent; cursor: pointer; }.ledger nav button.active { color: #f13728; }.ledger nav button.active::after { position: absolute; bottom: -1px; left: 0; width: 100%; height: 2px; background: #f13728; content: ""; }.ledger article { display: flex; align-items: center; justify-content: space-between; padding: 18px 10px; border-bottom: 1px dashed #e7e7e7; }.ledger article div { display: grid; gap: 8px; }.ledger article b { color: #555; font-weight: 500; }.ledger article span { color: #aaa; font-size: 13px; }.ledger article strong { color: #d9362b; }.ledger article strong.income { color: #34a469; }.empty { padding: 45px 0; color: #aaa; text-align: center; }@media (max-width: 620px) { .metrics { grid-template-columns: 1fr; }.metrics strong { font-size: 26px; }.ledger nav { gap: 18px; } }
</style>
