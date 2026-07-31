<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  createAddress,
  fetchAddresses,
  orderCheck,
  orderCreate,
  type Address,
} from "@/api/trade";
import { fetchUsableCoupons, type CouponUser } from "@/api/coupon";
import { ApiError } from "@/utils/request";

const route = useRoute();
const router = useRouter();
const addresses = ref<Address[]>([]);
const addressId = ref(0);
const payPrice = ref("0.00");
const couponPrice = ref("0.00");
const totalNum = ref(0);
const usable = ref<CouponUser[]>([]);
const selectedCouponIds = ref<number[]>([]);
const hint = ref("");
const loading = ref(false);
const showNew = ref(false);
const showInvoice = ref(false);
const invoiceDraft = ref({
  invoiceType: "ordinary" as "ordinary" | "special",
  headerType: "personal" as "personal" | "company",
  header: "",
  taxNo: "",
  email: "",
});
const invoice = ref<typeof invoiceDraft.value | null>(null);
const form = ref({
  real_name: "",
  phone: "",
  province: "上海市",
  city: "上海市",
  district: "浦东新区",
  detail: "",
  is_default: 1 as number,
});

const cartIds = computed(() =>
  String(route.query.cart_ids || "")
    .split(",")
    .map((x) => Number(x))
    .filter((x) => x > 0),
);

async function load() {
  if (!cartIds.value.length) {
    hint.value = "请从购物车选择商品";
    return;
  }
  loading.value = true;
  try {
    const addr = await fetchAddresses();
    addresses.value = addr.list || [];
    const def = addresses.value.find((a) => a.is_default === 1) || addresses.value[0];
    addressId.value = def?.address_id || 0;
    const usableRes = await fetchUsableCoupons(cartIds.value);
    usable.value = usableRes.list || [];
    if (addressId.value) {
      const check = await orderCheck(cartIds.value, addressId.value, selectedCouponIds.value);
      payPrice.value = Number(check.pay_price).toFixed(2);
      couponPrice.value = Number(check.coupon_price || 0).toFixed(2);
      totalNum.value = check.total_num;
    }
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());

function toggleCoupon(id: number) {
  const current = usable.value.find((coupon) => coupon.coupon_user_id === id);
  if (!current) return;
  const i = selectedCouponIds.value.indexOf(id);
  if (i >= 0) {
    selectedCouponIds.value = selectedCouponIds.value.filter((x) => x !== id);
  } else {
    // 规则与服务端一致：当前普通订单可叠加一张店铺券和一张平台券，
    // 选择同类型新券时直接替换，避免用户看到无意义的服务端冲突提示。
    selectedCouponIds.value = [
      ...selectedCouponIds.value.filter((selectedID) => {
        const selected = usable.value.find((coupon) => coupon.coupon_user_id === selectedID);
        return selected?.coupon_kind !== current.coupon_kind;
      }),
      id,
    ];
  }
  void load();
}

async function saveAddress() {
  try {
    const a = await createAddress({
      ...form.value,
      is_default: 1,
    });
    showNew.value = false;
    addressId.value = a.address_id;
    await load();
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "保存地址失败";
  }
}

async function submit() {
  if (!addressId.value) {
    hint.value = "请选择或新增收货地址";
    return;
  }
  loading.value = true;
  try {
    const g = await orderCreate(cartIds.value, addressId.value, "", selectedCouponIds.value);
    // 发票只能针对已支付的店铺订单申请；先保存本次选择，支付成功页会带入申请入口。
    if (invoice.value) {
      sessionStorage.setItem(
        `qixi:invoice:${g.group_order_id}`,
        JSON.stringify(invoice.value),
      );
    }
    router.replace({ name: "pay-result", params: { id: String(g.group_order_id) } });
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "下单失败";
  } finally {
    loading.value = false;
  }
}

function invoiceLabel(value: typeof invoiceDraft.value | null) {
  if (!value) return "不开发票";
  const kind = value.invoiceType === "special" ? "增值税专用发票" : "增值税电子普通发票";
  return `${kind} · ${value.headerType === "company" ? "企业" : "个人"}`;
}

function openInvoice() {
  if (invoice.value) invoiceDraft.value = { ...invoice.value };
  showInvoice.value = true;
}

function confirmInvoice() {
  const header = invoiceDraft.value.header.trim();
  const email = invoiceDraft.value.email.trim();
  if (!header) {
    hint.value = "请填写发票抬头";
    return;
  }
  if (!email || !/^\S+@\S+\.\S+$/.test(email)) {
    hint.value = "请填写正确的收票邮箱";
    return;
  }
  if (invoiceDraft.value.headerType === "company" && !invoiceDraft.value.taxNo.trim()) {
    hint.value = "企业抬头请填写纳税人识别号";
    return;
  }
  invoice.value = { ...invoiceDraft.value, header, email };
  showInvoice.value = false;
  hint.value = "";
}

function couponLabel(c: CouponUser) {
  return c.discount_type === "rate" ? `${Number(c.discount_value / 10).toFixed(1)}折` : `¥${Number(c.discount_value).toFixed(2)}`;
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>确认订单</h1>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div class="block">
        <h2>收货地址</h2>
        <label v-for="a in addresses" :key="a.address_id" class="addr">
          <input v-model="addressId" type="radio" :value="a.address_id" @change="load" />
          <span>
            {{ a.real_name }} {{ a.phone }}<br />
            {{ a.province }}{{ a.city }}{{ a.district }} {{ a.detail }}
          </span>
        </label>
        <button class="pc-btn ghost" type="button" @click="showNew = !showNew">
          {{ showNew ? "取消" : "新增地址" }}
        </button>
        <div v-if="showNew" class="form">
          <input v-model="form.real_name" placeholder="收货人" />
          <input v-model="form.phone" placeholder="手机号" />
          <input v-model="form.detail" placeholder="详细地址" />
          <button class="pc-btn" type="button" @click="saveAddress">保存</button>
        </div>
      </div>
      <div class="block">
        <h2>优惠券</h2>
        <label v-for="c in usable" :key="c.coupon_user_id" class="coupon">
          <input
            type="checkbox"
            :checked="selectedCouponIds.includes(c.coupon_user_id)"
            @change="toggleCoupon(c.coupon_user_id)"
          />
          <span>{{ c.coupon_title }} · {{ couponLabel(c) }}（满{{ c.use_min_price }}）</span>
        </label>
        <p v-if="!usable.length" class="muted">暂无可用券（可先去领券中心）</p>
      </div>
      <div class="block invoice-line">
        <div>
          <h2>发票信息</h2>
          <p class="muted">{{ invoiceLabel(invoice) }}<span v-if="invoice">，付款后可在订单中申请开票</span></p>
        </div>
        <button class="pc-btn ghost" type="button" @click="openInvoice">{{ invoice ? "修改" : "填写发票" }}</button>
      </div>
      <div class="block">
        <h2>应付</h2>
        <p>{{ totalNum }} 件 · 优惠 ¥{{ couponPrice }} · <strong>¥{{ payPrice }}</strong></p>
      </div>
      <button class="pc-btn" type="button" :disabled="loading" @click="submit">提交订单</button>
    </section>

    <div v-if="showInvoice" class="invoice-mask" @click.self="showInvoice = false">
      <section class="invoice-modal" role="dialog" aria-modal="true" aria-label="发票信息">
        <button class="close" type="button" aria-label="关闭" @click="showInvoice = false">×</button>
        <h2>发票信息</h2>
        <div class="invoice-row">
          <span>发票类型</span>
          <button :class="{ active: invoiceDraft.invoiceType === 'ordinary' }" type="button" @click="invoiceDraft.invoiceType = 'ordinary'">增值税电子普通发票</button>
          <button :class="{ active: invoiceDraft.invoiceType === 'special' }" type="button" @click="invoiceDraft.invoiceType = 'special'">增值税专用发票</button>
        </div>
        <p class="invoice-tip">电子发票默认发送至所提供的电子邮箱</p>
        <div class="invoice-row">
          <span>发票抬头</span>
          <button :class="{ active: invoiceDraft.headerType === 'personal' }" type="button" @click="invoiceDraft.headerType = 'personal'">个人</button>
          <button :class="{ active: invoiceDraft.headerType === 'company' }" type="button" @click="invoiceDraft.headerType = 'company'">企业</button>
        </div>
        <label class="invoice-field"><span>发票抬头</span><input v-model="invoiceDraft.header" :placeholder="invoiceDraft.headerType === 'company' ? '请填写企业名称' : '请填写个人姓名'" /></label>
        <label v-if="invoiceDraft.headerType === 'company' || invoiceDraft.invoiceType === 'special'" class="invoice-field"><span>纳税人识别号</span><input v-model="invoiceDraft.taxNo" placeholder="请填写纳税人识别号" /></label>
        <label class="invoice-field"><span>收票邮箱</span><input v-model="invoiceDraft.email" placeholder="请填写电子邮箱" /></label>
        <footer><button class="cancel" type="button" @click="showInvoice = false">取消</button><button class="confirm" type="button" @click="confirmInvoice">确定</button></footer>
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
.hint { color: #c0392b; }
.block { margin: 1.2rem 0; }
.addr, .coupon {
  display: flex;
  gap: 0.6rem;
  margin: 0.5rem 0;
  line-height: 1.5;
}
.form { display: grid; gap: 0.5rem; max-width: 360px; margin-top: 0.8rem; }
.form input {
  border: 1px solid var(--pc-line);
  border-radius: 6px;
  padding: 0.5rem 0.7rem;
}
.ghost { margin-top: 0.5rem; background: transparent; color: inherit; border: 1px solid var(--pc-line); }
.muted { color: var(--pc-muted); font-size: 0.9rem; }
.invoice-line { display: flex; align-items: center; justify-content: space-between; gap: 1rem; border-top: 1px solid #f0f0f0; padding-top: 1.2rem; }
.invoice-line h2 { margin: 0; }.invoice-line p { margin: 0.45rem 0 0; }
.invoice-mask { position: fixed; z-index: 90; inset: 0; display: grid; place-items: center; padding: 24px; background: rgb(0 0 0 / 48%); }
.invoice-modal { position: relative; width: min(680px, 100%); padding: 38px 48px 30px; background: #fff; box-shadow: 0 18px 48px rgb(0 0 0 / 20%); }
.invoice-modal h2 { margin: 0 0 30px; color: #2c2c2c; font-size: 24px; }.close { position: absolute; top: 16px; right: 20px; border: 0; background: transparent; color: #777; font-size: 32px; cursor: pointer; line-height: 1; }
.invoice-row { display: flex; align-items: center; flex-wrap: wrap; gap: 12px; margin: 18px 0; }.invoice-row > span, .invoice-field > span { width: 94px; color: #555; }.invoice-row button { min-width: 142px; padding: 9px 14px; border: 1px solid #dedede; background: #fff; color: #555; cursor: pointer; }.invoice-row button.active { border-color: #f13728; background: #fff5f4; color: #f13728; }.invoice-tip { margin: -7px 0 20px 94px; color: #999; font-size: 13px; }.invoice-field { display: flex; align-items: center; gap: 12px; margin: 16px 0; }.invoice-field input { min-width: 0; flex: 1; border: 1px solid #dedede; padding: 11px 12px; font-size: 14px; outline-color: #f13728; }.invoice-modal footer { display: flex; justify-content: flex-end; gap: 12px; margin-top: 30px; }.invoice-modal footer button { min-width: 92px; padding: 10px 22px; border: 1px solid #dedede; cursor: pointer; }.invoice-modal .cancel { background: #fff; color: #666; }.invoice-modal .confirm { border-color: #f13728; background: #f13728; color: #fff; }
@media (max-width: 600px) { .invoice-modal { padding: 32px 20px 24px; }.invoice-row > span, .invoice-field > span { width: 82px; }.invoice-row button { min-width: 0; }.invoice-tip { margin-left: 82px; } }
</style>
