<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchAddresses, type Address } from "@/api/trade";
import { http } from "@/utils/request";
import { ApiError } from "@/utils/request";

const route = useRoute();
const router = useRouter();
const addresses = ref<Address[]>([]);
const addressId = ref(0);
const payPrice = ref("0.00");
const totalIntegral = ref(0);
const userIntegral = ref(0);
const hint = ref("");
const loading = ref(false);

const cartIds = computed(() =>
  String(route.query.cart_ids || "")
    .split(",")
    .map((x) => Number(x))
    .filter((x) => x > 0),
);

async function load() {
  if (!cartIds.value.length) {
    hint.value = "请从积分商城选择商品";
    return;
  }
  loading.value = true;
  try {
    const addr = await fetchAddresses();
    addresses.value = addr.list || [];
    const def = addresses.value.find((a) => a.is_default === 1) || addresses.value[0];
    addressId.value = def?.address_id || 0;
    const check = await http.post<{
      pay_price: number;
      total_integral: number;
      user_integral: number;
    }>("/order/v3/check", { cart_ids: cartIds.value });
    payPrice.value = Number(check.pay_price).toFixed(2);
    totalIntegral.value = check.total_integral || 0;
    userIntegral.value = check.user_integral || 0;
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());

async function submit() {
  if (!addressId.value) {
    hint.value = "请选择收货地址";
    return;
  }
  loading.value = true;
  try {
    const g = await http.post<{ group_order_id: number }>("/order/v3/create", {
      cart_ids: cartIds.value,
      address_id: addressId.value,
    });
    await http.post(`/order/points/pay/${g.group_order_id}`, {});
    router.replace({ name: "pay-result", params: { id: String(g.group_order_id) } });
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "兑换失败";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>积分兑换确认</h1>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div class="block">
        <h2>收货地址</h2>
        <label v-for="a in addresses" :key="a.address_id" class="addr">
          <input v-model="addressId" type="radio" :value="a.address_id" />
          <span>{{ a.real_name }} {{ a.phone }} · {{ a.province }}{{ a.city }}{{ a.district }} {{ a.detail }}</span>
        </label>
      </div>
      <div class="block">
        <p>需积分 <strong>{{ totalIntegral }}</strong>（余额 {{ userIntegral }}）</p>
        <p>现金应付 <strong>¥{{ payPrice }}</strong></p>
      </div>
      <button class="pc-btn" type="button" :disabled="loading" @click="submit">确认兑换并支付</button>
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
.addr { display: flex; gap: 0.6rem; margin: 0.5rem 0; line-height: 1.5; }
</style>
