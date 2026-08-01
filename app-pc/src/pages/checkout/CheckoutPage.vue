<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  createAddress,
  fetchAddresses,
  orderCheck,
  orderCreate,
  type Address,
  type CheckoutMerchant,
} from "@/api/trade";
import { fetchUsableCoupons, type CouponUser } from "@/api/coupon";
import { ApiError } from "@/utils/request";

const route = useRoute();
const router = useRouter();
const addresses = ref<Address[]>([]);
const addressId = ref(0);
const payPrice = ref("0.00");
const couponPrice = ref("0.00");
const totalPrice = ref("0.00");
const postagePrice = ref("0.00");
const totalNum = ref(0);
const merchants = ref<CheckoutMerchant[]>([]);
const usable = ref<CouponUser[]>([]);
const selectedCouponIds = ref<number[]>([]);
const hint = ref("");
const loading = ref(false);
const buyerRemark = ref("");
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
const selectedAddress = computed(() => addresses.value.find((item) => item.address_id === addressId.value));

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
    merchants.value = [];
    const usableRes = await fetchUsableCoupons(cartIds.value);
    usable.value = usableRes.list || [];
    if (addressId.value) {
      const check = await orderCheck(cartIds.value, addressId.value, selectedCouponIds.value);
      payPrice.value = Number(check.pay_price).toFixed(2);
      couponPrice.value = Number(check.coupon_price || 0).toFixed(2);
      totalPrice.value = Number(check.total_price || 0).toFixed(2);
      postagePrice.value = Number(check.total_postage || 0).toFixed(2);
      totalNum.value = check.total_num;
      merchants.value = check.merchants || [];
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
    const g = await orderCreate(cartIds.value, addressId.value, buyerRemark.value.trim(), selectedCouponIds.value);
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
  <div class="checkout-page">
    <div class="pc-container checkout-container">
      <nav class="crumb" aria-label="当前位置"><RouterLink to="/">首页</RouterLink><span>›</span><RouterLink to="/cart">购物车</RouterLink><span>›</span><span>提交订单</span></nav>
      <p v-if="hint" class="hint">{{ hint }}</p>

      <section class="checkout-block address-block">
        <h1>收货地址</h1>
        <div class="address-list">
          <button
            v-for="address in addresses"
            :key="address.address_id"
            type="button"
            class="address-card"
            :class="{ active: address.address_id === addressId }"
            @click="addressId = address.address_id; load()"
          >
            <b v-if="address.is_default === 1">默认</b>
            <strong>{{ address.real_name }}</strong><span>{{ address.phone }}</span>
            <p>{{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}</p>
          </button>
          <button type="button" class="address-add" @click="showNew = !showNew"><i>＋</i>{{ showNew ? "收起地址表单" : "添加新地址" }}</button>
        </div>
        <div v-if="showNew" class="address-form">
          <input v-model.trim="form.real_name" placeholder="收货人姓名" />
          <input v-model.trim="form.phone" inputmode="tel" placeholder="手机号" />
          <input v-model.trim="form.province" placeholder="省份" />
          <input v-model.trim="form.city" placeholder="城市" />
          <input v-model.trim="form.district" placeholder="区县" />
          <input v-model.trim="form.detail" class="detail" placeholder="详细地址，如道路、门牌号" />
          <button class="pc-btn" type="button" @click="saveAddress">保存并使用</button>
        </div>
        <p v-else-if="!selectedAddress" class="empty-address">请选择或添加收货地址后继续结算</p>
      </section>

      <section v-for="merchant in merchants" :key="merchant.mer_id" class="checkout-block order-block">
        <header class="order-store"><b>{{ merchant.mer_name || `商户 #${merchant.mer_id}` }}</b><RouterLink :to="`/store/${merchant.mer_id}`">联系商家</RouterLink></header>
        <div v-for="item in merchant.items" :key="item.cart_id" class="order-item">
          <img v-if="item.image" :src="item.image" :alt="item.store_name" />
          <div v-else class="item-image-empty">商品</div>
          <div class="item-detail"><strong>{{ item.store_name }}</strong><span>规格：{{ item.product_attr_unique || "默认规格" }}</span></div>
          <b>¥{{ Number(item.price).toFixed(2) }} <small>×{{ item.cart_num }}</small></b>
        </div>
        <div class="checkout-row coupon-row">
          <b>使用优惠券</b>
          <div v-if="usable.length" class="coupon-options">
            <button
              v-for="coupon in usable"
              :key="coupon.coupon_user_id"
              type="button"
              :class="{ active: selectedCouponIds.includes(coupon.coupon_user_id) }"
              @click="toggleCoupon(coupon.coupon_user_id)"
            >{{ coupon.coupon_title }} · {{ couponLabel(coupon) }}<small>满{{ Number(coupon.use_min_price).toFixed(0) }}可用</small></button>
          </div>
          <span v-else class="muted">暂无可用优惠券</span>
          <strong v-if="merchant.coupon_price">-¥{{ Number(merchant.coupon_price).toFixed(2) }}</strong>
        </div>
        <div class="checkout-row"><b>运费</b><span class="muted">快递配送</span><strong>+¥{{ postagePrice }}</strong></div>
        <div class="checkout-row"><b>配送方式</b><span class="delivery-chip">快递配送</span></div>
        <div class="checkout-row invoice-row-line"><b>发票类型</b><span>{{ invoiceLabel(invoice) }}</span><button type="button" @click="openInvoice">{{ invoice ? "修改发票" : "填写发票" }}</button></div>
        <label class="remark-row"><b>买家留言</b><textarea v-model.trim="buyerRemark" maxlength="150" placeholder="填写内容与商家协商并确认，限 150 字以内"></textarea></label>
      </section>
      <section v-if="!loading && selectedAddress && !merchants.length" class="checkout-block no-items">没有可结算商品，请返回购物车重新选择。</section>

      <aside class="checkout-summary">
        <p><span>{{ totalNum }} 件商品，商品金额：</span><b>¥{{ totalPrice }}</b></p>
        <p><span>平台与店铺优惠：</span><b class="discount">-¥{{ couponPrice }}</b></p>
        <p><span>运费：</span><b>+¥{{ postagePrice }}</b></p>
        <p class="payable"><span>应付总额：</span><strong>¥{{ payPrice }}</strong></p>
        <button class="submit-order" type="button" :disabled="loading || !addressId || !merchants.length" @click="submit">{{ loading ? "提交中…" : "提交订单" }}</button>
      </aside>
    </div>

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
.checkout-page { min-height: 700px; padding: 18px 0 62px; background: #f6f6f6; }
.checkout-container { display: grid; gap: 14px; }
.crumb { display: flex; align-items: center; gap: 9px; height: 36px; color: #777; font-size: 13px; }.crumb a:hover { color: #f13728; }
.hint { margin: 0; color: #d9362b; font-size: 13px; }
.checkout-block { background: #fff; border: 1px solid #eee; padding: 20px 24px; }
.checkout-block > h1 { margin: 0 0 16px; color: #333; font-size: 18px; font-weight: 600; }
.address-block { position: relative; border-top-color: #eee; }.address-block::before { position: absolute; top: -2px; right: -1px; left: -1px; height: 3px; background: repeating-linear-gradient(90deg,#ef472e 0 10px,#f4af36 10px 20px,#2b82e0 20px 30px,#fff 30px 40px); content: ""; }
.address-list { display: flex; flex-wrap: wrap; gap: 16px; }
.address-card, .address-add { position: relative; width: 240px; min-height: 110px; padding: 16px 18px; border: 1px solid #e5e5e5; color: #555; background: #fff; text-align: left; }
.address-card.active { border-color: #f13728; box-shadow: inset 0 0 0 1px #f13728; }.address-card.active::after { position: absolute; right: 0; bottom: 0; width: 0; height: 0; border-right: 18px solid #f13728; border-top: 18px solid transparent; content: ""; }
.address-card b { position: absolute; top: 0; right: 0; padding: 3px 9px; color: #fff; background: #f13728; font-size: 11px; font-weight: 500; }.address-card strong { margin-right: 12px; color: #333; font-size: 14px; }.address-card span { color: #666; font-size: 13px; }.address-card p { margin: 12px 0 0; color: #888; font-size: 12px; line-height: 1.5; }
.address-add { display: grid; place-items: center; gap: 7px; color: #aaa; text-align: center; }.address-add i { color: #bbb; font-size: 31px; font-style: normal; font-weight: 200; }.address-add:hover { border-color: #f13728; color: #f13728; }.empty-address { margin: 16px 0 0; color: #999; font-size: 13px; }
.address-form { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; margin-top: 16px; padding-top: 16px; border-top: 1px dashed #e6e6e6; }.address-form input { min-width: 0; height: 36px; padding: 0 11px; border: 1px solid #ddd; outline-color: #f13728; font-size: 13px; }.address-form .detail { grid-column: span 2; }.address-form .pc-btn { padding: 0; border-radius: 0; }
.order-block { padding: 0 24px; }.order-store { display: flex; align-items: center; justify-content: space-between; min-height: 56px; border-bottom: 1px solid #eee; color: #555; font-size: 14px; }.order-store a { color: #777; font-size: 12px; }.order-store a:hover { color: #f13728; }
.order-item { display: grid; grid-template-columns: 76px minmax(0, 1fr) auto; align-items: center; gap: 14px; min-height: 116px; border-bottom: 1px solid #f0f0f0; }.order-item img, .item-image-empty { width: 76px; height: 76px; object-fit: cover; background: #f6f6f6; }.item-image-empty { display: grid; place-items: center; color: #aaa; font-size: 12px; }.item-detail { display: grid; gap: 8px; }.item-detail strong { color: #444; font-size: 14px; font-weight: 500; }.item-detail span { color: #999; font-size: 12px; }.order-item > b { color: #444; font-size: 14px; white-space: nowrap; }.order-item small { color: #888; font-size: 12px; font-weight: 400; }
.checkout-row { display: grid; grid-template-columns: 104px minmax(0, 1fr) auto; align-items: center; gap: 12px; min-height: 58px; border-bottom: 1px solid #f0f0f0; color: #555; font-size: 13px; }.checkout-row > b, .remark-row > b { color: #555; font-size: 13px; font-weight: 600; }.checkout-row > strong { color: #f13728; font-weight: 500; }.coupon-options { display: flex; flex-wrap: wrap; gap: 10px; padding: 8px 0; }.coupon-options button { display: grid; gap: 3px; min-width: 132px; padding: 7px 10px; border: 1px solid #f3c2bd; color: #e64a3b; background: #fff; text-align: left; font-size: 12px; }.coupon-options button.active { color: #fff; background: #f13728; border-color: #f13728; }.coupon-options small { font-size: 11px; color: inherit; opacity: .8; }.muted { color: #999; font-size: 12px; }.delivery-chip { width: max-content; padding: 6px 28px; border: 1px solid #f13728; color: #f13728; }.invoice-row-line button { justify-self: start; padding: 6px 16px; border: 1px solid #dedede; color: #666; background: #fff; font-size: 12px; }.invoice-row-line button:hover { border-color: #f13728; color: #f13728; }
.remark-row { display: grid; grid-template-columns: 104px minmax(0, 1fr); gap: 12px; min-height: 108px; padding: 17px 0; color: #555; }.remark-row textarea { min-height: 72px; padding: 10px 12px; resize: vertical; border: 0; background: #f7f7f7; outline-color: #f13728; font-size: 13px; }.no-items { color: #999; text-align: center; }
.checkout-summary { justify-self: stretch; margin-top: 0; padding: 22px 24px; background: #fff; border: 1px solid #eee; text-align: right; }.checkout-summary p { display: flex; justify-content: flex-end; gap: 30px; margin: 8px 0; color: #666; font-size: 13px; }.checkout-summary p span { min-width: 180px; }.checkout-summary p b { min-width: 100px; color: #555; font-weight: 500; }.checkout-summary p .discount { color: #f13728; }.checkout-summary .payable { margin-top: 14px; }.checkout-summary .payable strong { min-width: 100px; color: #f13728; font-size: 22px; }.submit-order { width: 180px; height: 44px; margin-top: 16px; border: 0; color: #fff; background: #f13728; font-size: 16px; font-weight: 600; }.submit-order:hover:not(:disabled) { background: #e52f21; }.submit-order:disabled { opacity: .52; cursor: not-allowed; }
.muted { color: var(--pc-muted); font-size: 0.9rem; }
.invoice-line { display: flex; align-items: center; justify-content: space-between; gap: 1rem; border-top: 1px solid #f0f0f0; padding-top: 1.2rem; }
.invoice-line h2 { margin: 0; }.invoice-line p { margin: 0.45rem 0 0; }
.invoice-mask { position: fixed; z-index: 90; inset: 0; display: grid; place-items: center; padding: 24px; background: rgb(0 0 0 / 48%); }
.invoice-modal { position: relative; width: min(680px, 100%); padding: 38px 48px 30px; background: #fff; box-shadow: 0 18px 48px rgb(0 0 0 / 20%); }
.invoice-modal h2 { margin: 0 0 30px; color: #2c2c2c; font-size: 24px; }.close { position: absolute; top: 16px; right: 20px; border: 0; background: transparent; color: #777; font-size: 32px; cursor: pointer; line-height: 1; }
.invoice-row { display: flex; align-items: center; flex-wrap: wrap; gap: 12px; margin: 18px 0; }.invoice-row > span, .invoice-field > span { width: 94px; color: #555; }.invoice-row button { min-width: 142px; padding: 9px 14px; border: 1px solid #dedede; background: #fff; color: #555; cursor: pointer; }.invoice-row button.active { border-color: #f13728; background: #fff5f4; color: #f13728; }.invoice-tip { margin: -7px 0 20px 94px; color: #999; font-size: 13px; }.invoice-field { display: flex; align-items: center; gap: 12px; margin: 16px 0; }.invoice-field input { min-width: 0; flex: 1; border: 1px solid #dedede; padding: 11px 12px; font-size: 14px; outline-color: #f13728; }.invoice-modal footer { display: flex; justify-content: flex-end; gap: 12px; margin-top: 30px; }.invoice-modal footer button { min-width: 92px; padding: 10px 22px; border: 1px solid #dedede; cursor: pointer; }.invoice-modal .cancel { background: #fff; color: #666; }.invoice-modal .confirm { border-color: #f13728; background: #f13728; color: #fff; }
@media (max-width: 860px) { .address-form { grid-template-columns: 1fr 1fr; }.address-form .detail { grid-column: span 2; }.order-block { padding-inline: 16px; }.checkout-block { padding-inline: 16px; }.checkout-row { grid-template-columns: 90px minmax(0,1fr); }.checkout-row > strong { grid-column: 2; }.invoice-row-line button { grid-column: 2; }.remark-row { grid-template-columns: 90px minmax(0,1fr); } }
@media (max-width: 600px) { .invoice-modal { padding: 32px 20px 24px; }.invoice-row > span, .invoice-field > span { width: 82px; }.invoice-row button { min-width: 0; }.invoice-tip { margin-left: 82px; } }
</style>
