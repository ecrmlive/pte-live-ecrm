<script setup lang="ts">
import { onMounted, ref } from "vue";
import AccountFrame from "@/components/AccountFrame.vue";
import { http, ApiError } from "@/utils/request";

interface PointsLedger {
  id: number;
  amount: number;
  reference_type: string;
  reference_id: string;
  created_at: string;
}

const points = ref(0);
const totalIncome = ref(0);
const totalExpense = ref(0);
const frozenPoints = ref(0);
const rows = ref<PointsLedger[]>([]);
const hint = ref("");

function actionLabel(value: string) {
  const labels: Record<string, string> = {
    order_paid: "订单赠送",
    order_refund: "退款退回",
    points_exchange: "积分兑换",
    platform_adjust: "平台调整",
  };
  return labels[value] || value || "积分变动";
}

async function load() {
  try {
    const data = await http.get<{ summary: { points: number; total_income: number; total_expense: number; frozen_points: number }; list: PointsLedger[] }>("/account/points");
    points.value = Number(data.summary?.points || 0);
    totalIncome.value = Number(data.summary?.total_income || 0);
    totalExpense.value = Number(data.summary?.total_expense || 0);
    frozenPoints.value = Number(data.summary?.frozen_points || 0);
    rows.value = data.list || [];
    hint.value = "";
  } catch (error) {
    hint.value = error instanceof ApiError ? error.message : "积分数据加载失败";
  }
}

onMounted(() => void load());
</script>

<template>
  <AccountFrame>
    <template #crumb><span>›</span> 我的积分</template>
    <h1 class="content-title">我的积分</h1>
    <p v-if="hint" class="hint">{{ hint }}</p>
    <section class="points-summary">
      <div><span>当前积分 <i>i</i></span><strong>{{ points }}</strong></div>
      <div><span>累计积分</span><strong>{{ totalIncome }}</strong></div>
      <div><span>累计消费</span><strong>{{ totalExpense }}</strong></div>
      <div><span>冻结积分</span><strong>{{ frozenPoints }}</strong></div>
    </section>
    <section class="ledger">
      <h2>积分明细</h2>
      <div v-if="!rows.length" class="empty">暂无积分明细</div>
      <article v-for="item in rows" :key="item.id">
        <div><strong>{{ actionLabel(item.reference_type) }}{{ Math.abs(item.amount) }}积分</strong><time>{{ item.created_at }}</time></div>
        <b :class="{ plus: item.amount > 0 }">{{ item.amount > 0 ? "+" : "" }}{{ item.amount }}</b>
      </article>
    </section>
  </AccountFrame>
</template>

<style scoped>
.content-title { margin: 0; padding-bottom: 20px; border-bottom: 1px solid #eee; color: #333; font-size: 20px; }.hint { color: #d9362b; }.points-summary { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; padding: 48px 0 44px; border-bottom: 1px solid #eee; }.points-summary div { display: grid; gap: 16px; }.points-summary span { color: #888; font-size: 14px; }.points-summary i { display: inline-grid; width: 14px; height: 14px; place-items: center; border: 1px solid #f13728; border-radius: 50%; color: #f13728; font-size: 10px; font-style: normal; }.points-summary strong { color: #333; font-size: 32px; font-weight: 400; }.ledger { margin-top: 32px; }.ledger h2 { width: max-content; margin: 0; padding: 0 0 15px; border-bottom: 2px solid #f13728; color: #f13728; font-size: 16px; }.ledger { border-bottom: 1px solid #eee; }.ledger article { display: flex; align-items: center; justify-content: space-between; padding: 20px 10px; border-top: 1px dashed #e8e8e8; }.ledger article strong, .ledger time { display: block; }.ledger time { margin-top: 8px; color: #aaa; font-size: 13px; }.ledger b { color: #36a85d; }.ledger b.plus { color: #36a85d; }.empty { padding: 44px 0; color: #999; text-align: center; } @media (max-width: 680px) { .points-summary { grid-template-columns: repeat(2, 1fr); gap: 28px 18px; }.points-summary strong { font-size: 27px; } }
</style>
