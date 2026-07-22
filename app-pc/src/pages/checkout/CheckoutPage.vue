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
  const i = selectedCouponIds.value.indexOf(id);
  if (i >= 0) {
    selectedCouponIds.value = selectedCouponIds.value.filter((x) => x !== id);
  } else {
    selectedCouponIds.value = [...selectedCouponIds.value, id];
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
    router.replace({ name: "pay-result", params: { id: String(g.group_order_id) } });
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "下单失败";
  } finally {
    loading.value = false;
  }
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
          <span>{{ c.coupon_title }} · ¥{{ Number(c.coupon_price).toFixed(2) }}（满{{ c.use_min_price }}）</span>
        </label>
        <p v-if="!usable.length" class="muted">暂无可用券（可先去领券中心）</p>
      </div>
      <div class="block">
        <h2>应付</h2>
        <p>{{ totalNum }} 件 · 优惠 ¥{{ couponPrice }} · <strong>¥{{ payPrice }}</strong></p>
      </div>
      <button class="pc-btn" type="button" :disabled="loading" @click="submit">提交订单</button>
    </section>
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
</style>
