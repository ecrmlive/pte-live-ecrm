<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import QRCode from "qrcode";
import { useRoute, useRouter } from "vue-router";
import {
  fetchOrderDetail,
	fetchOrderPaymentChannels,
  applyInvoice,
  payOrder,
  type GroupOrder,
	type PaymentChannel,
  type PayType,
} from "@/api/trade";
import { ApiError } from "@/utils/request";

const route = useRoute();
const router = useRouter();
const order = ref<GroupOrder | null>(null);
const hint = ref("");
const paying = ref(false);
const paymentChannels = ref<PaymentChannel[]>([]);
const invoiceDraft = ref<{
  invoiceType: "ordinary" | "special";
  headerType: "personal" | "company";
  header: string;
  taxNo: string;
  email: string;
} | null>(null);
const invoiceApplying = ref(false);
const invoiceApplied = ref(false);
const wechatQRCode = ref("");
const wechatExpiresAt = ref("");
let paymentPollTimer: number | undefined;

const id = () => Number(route.params.id);

async function load() {
  try {
    const [detail, channels] = await Promise.all([
      fetchOrderDetail(id()),
      fetchOrderPaymentChannels(id()),
    ]);
    order.value = detail;
    const cached = sessionStorage.getItem(`qixi:invoice:${id()}`);
    if (cached) {
      try { invoiceDraft.value = JSON.parse(cached); } catch { sessionStorage.removeItem(`qixi:invoice:${id()}`); }
    }
    paymentChannels.value = channels.list ?? [];
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  }
}

const invoiceName = computed(() => {
  if (!invoiceDraft.value) return "";
  return `${invoiceDraft.value.invoiceType === "special" ? "增值税专用发票" : "增值税电子普通发票"} · ${invoiceDraft.value.header}`;
});

async function requestInvoice() {
  if (!order.value || !invoiceDraft.value || !order.value.orders?.length) {
    hint.value = "订单明细尚未加载，无法申请发票";
    return;
  }
  invoiceApplying.value = true;
  try {
    await Promise.all(order.value.orders.map((item) => applyInvoice({
      order_id: item.order_id,
      invoice_type: invoiceDraft.value?.invoiceType === "special" ? 2 : 1,
      header_type: invoiceDraft.value?.headerType === "company" ? 2 : 1,
      header: invoiceDraft.value?.header || "",
      tax_no: invoiceDraft.value?.taxNo || "",
      email: invoiceDraft.value?.email || "",
    })));
    invoiceApplied.value = true;
    sessionStorage.removeItem(`qixi:invoice:${id()}`);
    hint.value = "发票申请已提交，请在我的发票中查看处理进度";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "发票申请失败";
  } finally {
    invoiceApplying.value = false;
  }
}

onMounted(() => void load());
onBeforeUnmount(() => stopWechatPaymentPolling());

function channelEnabled(channel: "wechat" | "alipay") {
  return paymentChannels.value.some((item) => item.channel === channel && item.enabled);
}

async function pay(type: PayType) {
  if (order.value?.pay_status !== "pending") {
    hint.value = "当前订单不可支付";
    return;
  }
  paying.value = true;
  try {
    const res = await payOrder(id(), type);
    if (type === "wechat" && "code_url" in res && res.code_url) {
      wechatQRCode.value = await QRCode.toDataURL(res.code_url, { width: 252, margin: 1, errorCorrectionLevel: "M" });
      wechatExpiresAt.value = res.expires_at || "";
      hint.value = "请使用微信扫描二维码完成支付";
      startWechatPaymentPolling();
    } else {
      order.value = res as GroupOrder;
      hint.value = "支付成功";
    }
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "支付失败";
  } finally {
    paying.value = false;
  }
}

function closeWechatPay() {
  stopWechatPaymentPolling();
  wechatQRCode.value = "";
  wechatExpiresAt.value = "";
}

function stopWechatPaymentPolling() {
  if (paymentPollTimer !== undefined) {
    window.clearTimeout(paymentPollTimer);
    paymentPollTimer = undefined;
  }
}

function startWechatPaymentPolling() {
  stopWechatPaymentPolling();
  const poll = async () => {
    try {
      const detail = await fetchOrderDetail(id());
      order.value = detail;
      if (detail.paid === 1) {
        hint.value = "支付成功";
        closeWechatPay();
        return;
      }
      if (wechatExpiresAt.value && new Date(wechatExpiresAt.value).getTime() <= Date.now()) {
        hint.value = "支付二维码已过期，请重新发起微信支付";
        closeWechatPay();
        return;
      }
    } catch {
      // A temporary polling failure must not interrupt the customer's scan.
    }
    paymentPollTimer = window.setTimeout(poll, 3000);
  };
  paymentPollTimer = window.setTimeout(poll, 3000);
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>支付订单</h1>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div v-if="order">
        <p>主单号 {{ order.group_order_sn }}</p>
        <p>应付 <strong>¥{{ Number(order.pay_price).toFixed(2) }}</strong> · {{ order.total_num }} 件</p>
        <p>状态：{{ order.pay_status === "closed" ? "已取消" : order.pay_status === "refunded" ? "已退款" : order.paid === 1 ? "已支付" : "待支付" }}</p>
        <div v-if="order.pay_status === 'pending'" class="actions">
          <button v-if="channelEnabled('wechat')" class="pc-btn ghost" type="button" :disabled="paying" @click="pay('wechat')">微信支付</button>
          <button v-if="channelEnabled('alipay')" class="pc-btn ghost" type="button" :disabled="paying" @click="pay('alipay')">支付宝支付</button>
          <p v-if="!channelEnabled('wechat') && !channelEnabled('alipay')" class="hint">当前订单暂无可用的第三方支付方式，请联系店铺或稍后重试。</p>
        </div>
        <div v-else-if="order.paid === 1" class="actions">
          <RouterLink class="pc-btn" to="/orders">查看订单</RouterLink>
          <button class="pc-btn ghost" type="button" @click="router.push('/goods')">继续购物</button>
        </div>
        <div v-else class="actions"><RouterLink class="pc-btn" to="/orders">返回订单列表</RouterLink><button class="pc-btn ghost" type="button" @click="router.push('/goods')">继续购物</button></div>
        <div v-if="order.paid === 1 && order.orders?.length" class="after-sales">
          <b>订单售后</b>
          <div v-for="item in order.orders" :key="item.order_id"><span>{{ item.mer_name || `店铺订单 #${item.order_id}` }}</span><RouterLink v-if="['paid', 'fulfilling', 'shipped'].includes(item.status)" :to="`/refunds?order_id=${item.order_id}`">申请售后</RouterLink><span v-else class="muted">当前订单状态：{{ item.status }}</span></div>
        </div>
        <div v-if="order.paid === 1 && invoiceDraft && !invoiceApplied" class="invoice-apply">
          <div><b>发票信息</b><p>{{ invoiceName }}，将发送至 {{ invoiceDraft.email }}</p></div>
          <button class="pc-btn" type="button" :disabled="invoiceApplying" @click="requestInvoice">{{ invoiceApplying ? "提交中…" : "申请开票" }}</button>
        </div>
        <RouterLink v-else-if="invoiceApplied" class="invoice-link" to="/user/invoices">查看我的发票</RouterLink>
      </div>
    </section>
    <div v-if="wechatQRCode" class="wechat-mask" role="dialog" aria-modal="true" aria-label="微信扫码支付">
      <section class="wechat-dialog">
        <button class="wechat-close" type="button" aria-label="关闭" @click="closeWechatPay">×</button>
        <h2>微信扫码支付</h2>
        <p>请使用微信扫一扫完成付款</p>
        <img :src="wechatQRCode" alt="微信支付二维码" class="wechat-qr" />
        <strong>¥{{ order ? Number(order.pay_price).toFixed(2) : "0.00" }}</strong>
        <small v-if="wechatExpiresAt">支付二维码有效至 {{ new Date(wechatExpiresAt).toLocaleTimeString() }}</small>
        <button class="pc-btn" type="button" @click="load(); closeWechatPay()">我已完成支付</button>
      </section>
    </div>
  </div>
</template>

<style scoped>
.panel {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: calc(var(--pc-radius) + 4px);
  padding: 1.5rem;
  box-shadow: var(--pc-shadow);
}
.hint { color: var(--pc-muted); }
.actions { display: flex; flex-wrap: wrap; gap: 0.8rem; margin-top: 1.2rem; }
.ghost { background: transparent; color: inherit; border: 1px solid var(--pc-line); }
.invoice-apply { display: flex; align-items: center; justify-content: space-between; gap: 1rem; margin-top: 24px; padding: 18px 20px; border: 1px solid #f0d2cf; background: #fff9f8; }.invoice-apply b { color: #333; }.invoice-apply p { margin: 7px 0 0; color: #777; font-size: 14px; }.invoice-link { display: inline-block; margin-top: 22px; color: #f13728; }
.after-sales { display: grid; gap: 10px; margin-top: 24px; padding: 18px 20px; border: 1px solid #eee; }.after-sales > div { display: flex; justify-content: space-between; gap: 16px; color: #666; font-size: 14px; }.after-sales a { color: #f13728; }.muted { color: #999; }
.wechat-mask { position: fixed; inset: 0; z-index: 100; display: grid; place-items: center; padding: 24px; background: rgba(0, 0, 0, .55); }
.wechat-dialog { position: relative; width: min(420px, 100%); padding: 34px 42px 38px; text-align: center; color: #2f2f2f; background: #fff; border-radius: 4px; box-shadow: 0 24px 70px rgba(0, 0, 0, .24); }.wechat-dialog h2 { margin: 0; font-size: 24px; }.wechat-dialog p { margin: 12px 0 20px; color: #777; }.wechat-qr { display: block; width: 252px; height: 252px; margin: 0 auto 16px; background: #fff; }.wechat-dialog strong, .wechat-dialog small { display: block; }.wechat-dialog strong { margin-bottom: 8px; font-size: 22px; color: #f13728; }.wechat-dialog small { margin-bottom: 22px; color: #999; }.wechat-close { position: absolute; top: 12px; right: 14px; border: 0; background: transparent; color: #888; font-size: 30px; line-height: 1; cursor: pointer; }
</style>
