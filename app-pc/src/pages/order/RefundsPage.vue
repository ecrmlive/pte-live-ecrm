<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import AccountFrame from "@/components/AccountFrame.vue";
import { applyRefund, cancelRefund, fetchRefunds, requestPlatformIntervention, type RefundOrder } from "@/api/aftersale";
import { ApiError } from "@/utils/request";

const route = useRoute();
const router = useRouter();
const list = ref<RefundOrder[]>([]);
const hint = ref("");
const loading = ref(false);
const applying = ref(false);
const reason = ref("");
const orderID = computed(() => Number(route.query.order_id || 0));

const statusText: Record<RefundOrder["status_code"], string> = {
  applied: "等待商家处理", merchant_handling: "商家处理中", platform_intervene: "平台介入中", refunding: "退款处理中", refunded: "已退款", rejected: "已拒绝", cancelled: "已取消",
};

async function load() {
  loading.value = true;
  try { list.value = (await fetchRefunds()).list || []; hint.value = ""; }
  catch (error) { hint.value = error instanceof ApiError ? error.message : "售后单加载失败"; }
  finally { loading.value = false; }
}

async function submitApply() {
  if (!orderID.value) { hint.value = "请选择需要售后的订单"; return; }
  if (!reason.value.trim()) { hint.value = "请填写售后原因"; return; }
  applying.value = true;
  try {
    await applyRefund(orderID.value, reason.value.trim(), `pc-refund-${orderID.value}-${Date.now()}`);
    reason.value = "";
    hint.value = "售后申请已提交，等待商家处理";
    await router.replace({ path: "/refunds" });
    await load();
  } catch (error) { hint.value = error instanceof ApiError ? error.message : "提交售后申请失败"; }
  finally { applying.value = false; }
}

async function cancel(item: RefundOrder) { try { await cancelRefund(item.refund_order_id); await load(); } catch (error) { hint.value = error instanceof ApiError ? error.message : "取消售后失败"; } }
async function intervene(item: RefundOrder) { try { await requestPlatformIntervention(item.refund_order_id); await load(); } catch (error) { hint.value = error instanceof ApiError ? error.message : "申请平台介入失败"; } }

onMounted(() => void load());
</script>

<template>
  <AccountFrame>
    <template #crumb><span>›</span> 售后/退款</template>
    <header class="refund-header"><h1>售后/退款</h1><button type="button" @click="load">刷新</button></header>
    <section v-if="orderID" class="apply-panel">
      <div><b>申请售后</b><p>订单 #{{ orderID }} 的退款金额将由系统按该订单实际支付金额核算。</p></div>
      <textarea v-model="reason" maxlength="500" placeholder="请填写退款/售后原因，最多 500 字" />
      <footer><button type="button" @click="router.push('/orders')">返回订单</button><button class="pc-btn" type="button" :disabled="applying" @click="submitApply">{{ applying ? "提交中…" : "提交售后申请" }}</button></footer>
    </section>
    <p v-if="hint" class="hint">{{ hint }}</p>
    <p v-else-if="loading" class="empty">正在加载售后记录…</p>
    <p v-else-if="!list.length" class="empty">暂无售后记录。可从订单详情发起退款或退货申请。</p>
    <section v-else class="refund-list">
      <article v-for="item in list" :key="item.refund_order_id" class="refund-card">
        <header><b>售后单号：{{ item.refund_order_sn }}</b><span :class="item.status_code">{{ statusText[item.status_code] }}</span></header>
        <dl><div><dt>原订单</dt><dd>#{{ item.order_id }}</dd></div><div><dt>退款金额</dt><dd class="price">¥{{ Number(item.refund_price).toFixed(2) }}</dd></div><div><dt>申请时间</dt><dd>{{ item.create_time }}</dd></div><div><dt>售后原因</dt><dd>{{ item.refund_message }}</dd></div></dl>
        <footer><button v-if="item.status_code === 'applied' || item.status_code === 'merchant_handling'" type="button" @click="cancel(item)">取消申请</button><button v-if="item.status_code === 'applied' || item.status_code === 'merchant_handling'" type="button" @click="intervene(item)">申请平台介入</button></footer>
      </article>
    </section>
  </AccountFrame>
</template>

<style scoped>
.refund-header { display: flex; align-items: center; justify-content: space-between; padding-bottom: 19px; border-bottom: 1px solid #eee; }.refund-header h1 { margin: 0; font-size: 20px; }.refund-header button { border: 0; color: #777; background: transparent; cursor: pointer; }.refund-header button:hover { color: #f13728; }.apply-panel { display: grid; gap: 14px; margin-top: 24px; padding: 20px 22px; border: 1px solid #f0d5d2; background: #fff9f8; }.apply-panel b { font-size: 16px; }.apply-panel p { margin: 7px 0 0; color: #777; font-size: 13px; }.apply-panel textarea { min-height: 92px; padding: 12px; border: 1px solid #ddd; resize: vertical; outline: none; }.apply-panel textarea:focus { border-color: #f13728; }.apply-panel footer { display: flex; justify-content: flex-end; gap: 12px; }.apply-panel footer button:first-child { border: 1px solid #ddd; padding: 8px 16px; color: #777; background: #fff; cursor: pointer; }.hint { color: #d9362b; }.empty { padding: 76px 0; color: #aaa; text-align: center; }.refund-list { display: grid; gap: 16px; padding-top: 24px; }.refund-card { padding: 19px 22px; border: 1px solid #eee; }.refund-card header { display: flex; justify-content: space-between; padding-bottom: 14px; border-bottom: 1px dashed #eee; color: #444; }.refund-card header span { color: #e79a20; }.refund-card header .refunded { color: #36a367; }.refund-card header .rejected, .refund-card header .cancelled { color: #999; }.refund-card dl { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 13px 32px; margin: 17px 0 0; }.refund-card dl div { display: flex; gap: 10px; min-width: 0; }.refund-card dt { flex: none; color: #999; }.refund-card dd { overflow: hidden; margin: 0; color: #555; text-overflow: ellipsis; white-space: nowrap; }.refund-card .price { color: #f13728; font-weight: 600; }.refund-card footer { display: flex; justify-content: flex-end; gap: 16px; margin-top: 18px; }.refund-card footer button { border: 0; padding: 0; color: #777; background: transparent; cursor: pointer; }.refund-card footer button:hover { color: #f13728; }@media (max-width: 650px) { .refund-card dl { grid-template-columns: 1fr; } }
</style>
